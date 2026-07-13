#!/usr/bin/env python3
"""Deterministic validation for SDP Toolkit and consuming-project contracts."""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path, PurePosixPath
from typing import Any, Iterable, Sequence

import yaml
from jsonschema import Draft202012Validator, FormatChecker
from jsonschema.exceptions import SchemaError


SEMVER_PATTERN = re.compile(
    r"^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)"
    r"(?:-([0-9A-Za-z.-]+))?(?:\+([0-9A-Za-z.-]+))?$"
)
WINDOWS_DRIVE_PATTERN = re.compile(r"^[A-Za-z]:")
SKILL_ID_PATTERN = re.compile(r"^sdp-[a-z0-9]+(?:-[a-z0-9]+)*$")

SUPPORTED_PROJECT_MANIFEST_SCHEMAS = frozenset({"1.0"})
SUPPORTED_INSTALLED_MANIFEST_SCHEMAS = frozenset({"1.0"})
SUPPORTED_TRACE_EVENT_SCHEMAS = frozenset({"1.0"})
SUPPORTED_RELEASE_RECORD_SCHEMAS = frozenset({"1.0"})
SUPPORTED_FIX_RECORD_SCHEMAS = frozenset({"1.0"})

RELEASE_EVENT_TYPES = frozenset(
    {
        "release-planned",
        "release-version-selected",
        "release-candidate-opened",
        "release-verification-completed",
        "release-approved",
        "release-tag-created",
        "release-published",
        "release-yanked",
        "release-migration-applied",
        "release-notes-corrected",
    }
)
RELEASE_NOTE_CATEGORIES = frozenset(
    {"Added", "Changed", "Fixed", "Deprecated", "Removed", "Security", "Migration"}
)
RELATION_PATH_FIELDS = (
    "path",
    "study",
    "implementationNotes",
    "releaseNotes",
    "manifest",
    "releaseRecord",
)


@dataclass(frozen=True, order=True)
class SemVer:
    major: int
    minor: int
    patch: int
    prerelease_rank: int
    prerelease: tuple[tuple[int, int | str], ...]

    @classmethod
    def parse(cls, value: str) -> "SemVer":
        match = SEMVER_PATTERN.fullmatch(value)
        if not match:
            raise ValueError(f"Invalid SemVer: {value}")
        prerelease_text = match.group(4)
        identifiers: list[tuple[int, int | str]] = []
        if prerelease_text:
            for identifier in prerelease_text.split("."):
                if not identifier or (
                    identifier.isdigit() and len(identifier) > 1 and identifier.startswith("0")
                ):
                    raise ValueError(f"Invalid SemVer prerelease: {value}")
                identifiers.append((0, int(identifier)) if identifier.isdigit() else (1, identifier))
        # A final release has higher precedence than the same core prerelease.
        return cls(
            int(match.group(1)),
            int(match.group(2)),
            int(match.group(3)),
            1 if prerelease_text is None else 0,
            tuple(identifiers),
        )


def load_yaml(path: Path) -> Any:
    try:
        return yaml.safe_load(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, yaml.YAMLError) as exc:
        raise ValueError(f"Cannot parse YAML {path}: {exc}") from exc


def load_json(path: Path) -> Any:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (OSError, UnicodeError, json.JSONDecodeError) as exc:
        raise ValueError(f"Cannot parse JSON {path}: {exc}") from exc


def _json_error_key(error: Any) -> tuple[str, ...]:
    return tuple(str(part) for part in error.absolute_path)


def validate_json(instance: Any, schema: Any, label: str) -> list[str]:
    validator = Draft202012Validator(schema, format_checker=FormatChecker())
    errors: list[str] = []
    for error in sorted(validator.iter_errors(instance), key=_json_error_key):
        location = "".join(
            f"[{part}]" if isinstance(part, int) else f".{part}" for part in error.absolute_path
        )
        errors.append(f"{label}{location}: {error.message}")
    return errors


def parse_front_matter(path: Path) -> dict[str, Any]:
    try:
        text = path.read_text(encoding="utf-8")
    except (OSError, UnicodeError) as exc:
        raise ValueError(f"Cannot read skill metadata {path}: {exc}") from exc
    match = re.match(r"\A---\n(.*?)\n---(?:\n|\Z)", text, re.DOTALL)
    if not match:
        raise ValueError(f"Missing or malformed YAML front matter: {path}")
    try:
        data = yaml.safe_load(match.group(1))
    except yaml.YAMLError as exc:
        raise ValueError(f"Cannot parse YAML front matter {path}: {exc}") from exc
    if not isinstance(data, dict):
        raise ValueError(f"Front matter is not an object: {path}")
    return data


def release_sections(text: str) -> list[tuple[str, str]]:
    headings = list(re.finditer(r"(?m)^## \[([^\]]+)\](?: - \d{4}-\d{2}-\d{2})?\s*$", text))
    sections: list[tuple[str, str]] = []
    for index, heading in enumerate(headings):
        end = headings[index + 1].start() if index + 1 < len(headings) else len(text)
        sections.append((heading.group(1), text[heading.start() : end].rstrip() + "\n"))
    return sections


def git_show(repo: Path, ref: str, path: str) -> str | None:
    result = subprocess.run(
        ["git", "-C", str(repo), "show", f"{ref}:{path}"],
        capture_output=True,
        text=True,
    )
    return result.stdout if result.returncode == 0 else None


def validate_released_immutability(repo: Path, notes: str, base_ref: str | None) -> list[str]:
    if not base_ref:
        return []
    baseline = git_show(repo, base_ref, "RELEASE-NOTES.md")
    if baseline is None:
        return []
    current_sections = dict(release_sections(notes))
    errors: list[str] = []
    for version, baseline_section in release_sections(baseline):
        if version == "Unreleased":
            continue
        if current_sections.get(version) != baseline_section:
            errors.append(f"Released notes section [{version}] differs from {base_ref}")
    return errors


def validate_semver_values(values: Iterable[str], label: str) -> list[str]:
    errors: list[str] = []
    for value in values:
        try:
            SemVer.parse(value)
        except (TypeError, ValueError) as exc:
            errors.append(f"{label}: {exc}")
    return errors


def _schema_version_errors(
    data: Any, supported: frozenset[str], label: str
) -> list[str]:
    if not isinstance(data, dict) or "schemaVersion" not in data:
        return []
    version = data.get("schemaVersion")
    if not isinstance(version, str) or version not in supported:
        return [
            f"{label}: unsupported schemaVersion {version!r}; supported versions: "
            + ", ".join(sorted(supported))
        ]
    return []


def _read_yaml(path: Path, label: str, required: bool = True) -> tuple[Any | None, list[str]]:
    if not path.is_file():
        return None, [f"{label}: required file is missing"] if required else []
    try:
        return load_yaml(path), []
    except ValueError as exc:
        return None, [str(exc)]


def _read_json(path: Path, label: str, required: bool = True) -> tuple[Any | None, list[str]]:
    if not path.is_file():
        return None, [f"{label}: required file is missing"] if required else []
    try:
        return load_json(path), []
    except ValueError as exc:
        return None, [str(exc)]


def portable_relative_path_error(value: Any) -> str | None:
    if not isinstance(value, str) or not value:
        return "must be a non-empty relative path string"
    if "\\" in value:
        return "must use portable '/' separators"
    if value.startswith("/") or WINDOWS_DRIVE_PATTERN.match(value):
        return "must be relative, not absolute"
    if "\x00" in value:
        return "must not contain a NUL character"
    pure = PurePosixPath(value)
    if ".." in pure.parts:
        return "must not contain parent traversal"
    normalized = pure.as_posix()
    if normalized in {"", "."} or normalized != value:
        return "must be normalized"
    return None


def _within(path: Path, root: Path) -> bool:
    try:
        path.relative_to(root)
    except ValueError:
        return False
    return True


def validate_existing_path(
    value: Any,
    bases: Sequence[Path],
    label: str,
    *,
    require_file: bool = True,
) -> list[str]:
    path_error = portable_relative_path_error(value)
    if path_error:
        return [f"{label}: {path_error}: {value!r}"]
    assert isinstance(value, str)
    relative = Path(*PurePosixPath(value).parts)
    candidates: list[Path] = []
    for base in bases:
        resolved_base = base.resolve()
        candidate = (resolved_base / relative).resolve()
        if not _within(candidate, resolved_base):
            return [f"{label}: resolved path escapes {resolved_base}: {value}"]
        candidates.append(candidate)
    if any(candidate.is_file() if require_file else candidate.exists() for candidate in candidates):
        return []
    rendered = " or ".join(str(candidate) for candidate in candidates)
    kind = "file" if require_file else "path"
    return [f"{label}: referenced {kind} does not exist: {rendered}"]


def validate_release_notes(path: Path, label: str, *, required: bool = True) -> list[str]:
    if not path.is_file():
        return [f"{label}: required file is missing"] if required else []
    try:
        text = path.read_text(encoding="utf-8")
    except (OSError, UnicodeError) as exc:
        return [f"{label}: cannot read: {exc}"]
    errors: list[str] = []
    if not re.search(r"(?m)^# .+\S\s*$", text):
        errors.append(f"{label}: must contain a level-one title")
    sections = release_sections(text)
    names = [name for name, _ in sections]
    if not sections or names[0] != "Unreleased":
        errors.append(f"{label}: must begin its version sections with ## [Unreleased]")
    if names.count("Unreleased") != 1:
        errors.append(f"{label}: must contain exactly one [Unreleased] section")
    if sections and sections[0][0] == "Unreleased":
        if not re.search(r"(?m)^## \[Unreleased\]\s*$", sections[0][1]):
            errors.append(f"{label}: [Unreleased] heading must not carry a release date")
        if not re.search(r"(?m)^Release-Date:\s*unreleased\s*$", sections[0][1]):
            errors.append(f"{label}: [Unreleased] must contain Release-Date: unreleased")
    seen: set[str] = set()
    for name, _ in sections:
        if name in seen:
            errors.append(f"{label}: duplicate release-notes section [{name}]")
        seen.add(name)
        if name != "Unreleased":
            try:
                SemVer.parse(name)
            except ValueError as exc:
                errors.append(f"{label}: released section [{name}] is not SemVer: {exc}")
    category_headings = list(re.finditer(r"(?m)^### ([^\r\n]+)\s*$", text))
    for heading in category_headings:
        category = heading.group(1).strip()
        if category not in RELEASE_NOTE_CATEGORIES:
            errors.append(
                f"{label}: unsupported release-note category {category!r}; "
                f"allowed: {', '.join(sorted(RELEASE_NOTE_CATEGORIES))}"
            )
        next_heading = re.search(r"(?m)^#{2,3} ", text[heading.end() :])
        body_end = (
            heading.end() + next_heading.start() if next_heading is not None else len(text)
        )
        if not text[heading.end() : body_end].strip():
            errors.append(f"{label}: release-note category {category!r} is empty")
    return errors


def validate_skill_metadata(
    path: Path,
    expected_skill_id: str,
    expected_version: str,
    toolkit_version: str,
) -> list[str]:
    errors: list[str] = []
    try:
        metadata = parse_front_matter(path)
    except ValueError as exc:
        return [str(exc)]
    if metadata.get("skillId") != expected_skill_id:
        errors.append(f"{path}: skillId must be {expected_skill_id}")
    if metadata.get("skillVersion") != expected_version:
        errors.append(f"{path}: skillVersion differs from installed manifest")
    errors += validate_semver_values(
        [str(metadata.get("skillVersion", ""))], f"{path}: skillVersion"
    )
    minimum_toolkit = metadata.get("minimumToolkitVersion")
    try:
        if SemVer.parse(toolkit_version) < SemVer.parse(str(minimum_toolkit)):
            errors.append(
                f"{path}: minimumToolkitVersion {minimum_toolkit} is newer than "
                f"installed Toolkit {toolkit_version}"
            )
    except ValueError as exc:
        errors.append(f"{path}: {exc}")
    if not isinstance(metadata.get("capabilities"), list) or not metadata["capabilities"]:
        errors.append(f"{path}: capabilities must be a non-empty list")
    if not metadata.get("compatibilityNotes"):
        errors.append(f"{path}: compatibilityNotes is required")
    return errors


def _known_subject_category(subject_id: str) -> str | None:
    patterns = (
        (r"^REL-", "releases"),
        (r"^Sprint-", "sprints"),
        (r"^Refactor-", "refactors"),
        (r"^SPI-", "iterations"),
        (r"^SPS-", "slices"),
        (r"^(FIX|HOTFIX)-", "fixes"),
    )
    for pattern, category in patterns:
        if re.match(pattern, subject_id):
            return category
    return None


def validate_ledger(
    path: Path,
    schema_root: Path,
    *,
    required: bool = False,
    relations: dict[str, Any] | None = None,
    label: str | None = None,
) -> list[str]:
    ledger_label = label or str(path)
    if not path.is_file():
        return [f"{ledger_label}: required file is missing"] if required else []
    try:
        lines = path.read_text(encoding="utf-8").splitlines()
    except (OSError, UnicodeError) as exc:
        return [f"{ledger_label}: cannot read: {exc}"]
    try:
        generic_schema = load_json(schema_root / "ledger-event.schema.json")
        release_schema = load_json(schema_root / "release-event.schema.json")
    except ValueError as exc:
        return [str(exc)]
    errors: list[str] = []
    event_ids: set[str] = set()
    for line_number, line in enumerate(lines, start=1):
        if not line.strip():
            continue
        line_label = f"{ledger_label}:{line_number}"
        try:
            event = json.loads(line)
        except json.JSONDecodeError as exc:
            errors.append(f"{line_label}: invalid JSON: {exc}")
            continue
        errors += _schema_version_errors(event, SUPPORTED_TRACE_EVENT_SCHEMAS, line_label)
        errors += validate_json(event, generic_schema, line_label)
        if not isinstance(event, dict):
            continue
        event_id = event.get("eventId")
        if isinstance(event_id, str):
            if event_id in event_ids:
                errors.append(f"{line_label}: duplicate eventId {event_id}")
            event_ids.add(event_id)
        event_type = event.get("eventType")
        if event_type in RELEASE_EVENT_TYPES:
            errors += validate_json(event, release_schema, f"{line_label} release event")
        if relations:
            references: list[tuple[str, str]] = []
            release_id = event.get("releaseId")
            if isinstance(release_id, str):
                references.append(("releases", release_id))
            subject_id = event.get("subjectId")
            if isinstance(subject_id, str):
                category = _known_subject_category(subject_id)
                if category:
                    references.append((category, subject_id))
            for category, reference in references:
                category_map = relations.get(category)
                if isinstance(category_map, dict) and reference not in category_map:
                    errors.append(
                        f"{line_label}: {reference} is not defined in Relations.{category}"
                    )
    return errors


def validate_relations_semantics(relations: Any, sdp_root: Path, label: str) -> list[str]:
    if not isinstance(relations, dict):
        return []
    errors: list[str] = []
    for category, records in relations.items():
        if not isinstance(records, dict):
            continue
        for record_id, relation in records.items():
            if not isinstance(relation, dict):
                continue
            for field in RELATION_PATH_FIELDS:
                if field in relation:
                    errors += validate_existing_path(
                        relation[field],
                        [sdp_root],
                        f"{label}.{category}.{record_id}.{field}",
                    )

    reference_rules: dict[str, dict[str, str]] = {
        "sprints": {"release": "releases", "iterations": "iterations"},
        "refactors": {"release": "releases", "iterations": "iterations"},
        "iterations": {"sprint": "sprints", "refactor": "refactors", "slices": "slices"},
        "slices": {
            "sprint": "sprints",
            "refactor": "refactors",
            "iteration": "iterations",
            "release": "releases",
            "review": "reviews",
            "verification": "verification",
        },
        "fixes": {"release": "releases", "review": "reviews", "verification": "verification"},
        "releases": {
            "sprints": "sprints",
            "fixes": "fixes",
            "slices": "slices",
            "reviews": "reviews",
            "verification": "verification",
            "migrations": "migrations",
        },
    }
    for category, fields in reference_rules.items():
        records = relations.get(category)
        if not isinstance(records, dict):
            continue
        for record_id, relation in records.items():
            if not isinstance(relation, dict):
                continue
            for field, target_category in fields.items():
                value = relation.get(field)
                references = value if isinstance(value, list) else [value]
                target = relations.get(target_category)
                if not isinstance(target, dict):
                    continue
                for reference in references:
                    if isinstance(reference, str) and reference not in target:
                        errors.append(
                            f"{label}.{category}.{record_id}.{field}: dangling ID {reference}; "
                            f"not found in {target_category}"
                        )
    return errors


def validate_current_index_semantics(
    current: Any,
    relations: Any,
    label: str,
    project_manifest: Any | None = None,
) -> list[str]:
    if not isinstance(current, dict):
        return []
    errors: list[str] = []
    if isinstance(relations, dict):
        release = current.get("release")
        if isinstance(release, dict) and isinstance(release.get("activeReleaseId"), str):
            release_id = release["activeReleaseId"]
            relation_releases = relations.get("releases")
            if isinstance(relation_releases, dict) and release_id not in relation_releases:
                errors.append(f"{label}.release.activeReleaseId: dangling ID {release_id}")
        active = current.get("active")
        if isinstance(active, dict):
            category_by_field = {
                "sprint": "sprints",
                "refactor": "refactors",
                "iteration": "iterations",
                "slice": "slices",
                "fix": "fixes",
            }
            for field, category in category_by_field.items():
                reference = active.get(field)
                target = relations.get(category)
                if isinstance(reference, str) and isinstance(target, dict) and reference not in target:
                    errors.append(f"{label}.active.{field}: dangling ID {reference}")

    if isinstance(project_manifest, dict):
        index_release = current.get("release")
        project_release = project_manifest.get("release")
        if isinstance(index_release, dict) and isinstance(project_release, dict):
            comparisons = (
                ("state", "state"),
                ("targetVersion", "nextTargetVersion"),
            )
            for index_field, project_field in comparisons:
                if index_release.get(index_field) != project_release.get(project_field):
                    errors.append(
                        f"{label}.release.{index_field} differs from project manifest "
                        f"release.{project_field}"
                    )
            previous = index_release.get("previousVersion")
            if previous is not None and previous != project_release.get("currentVersion"):
                errors.append(
                    f"{label}.release.previousVersion differs from project manifest "
                    "release.currentVersion"
                )
            active_release_id = index_release.get("activeReleaseId")
            target_version = index_release.get("targetVersion")
            if isinstance(active_release_id, str) and isinstance(target_version, str):
                if active_release_id != f"REL-{target_version}":
                    errors.append(
                        f"{label}.release.activeReleaseId does not match targetVersion"
                    )
        index_active = current.get("active")
        development = project_manifest.get("development")
        if isinstance(index_active, dict) and isinstance(development, dict):
            fields = {
                "sprint": "sprintId",
                "refactor": "refactorId",
                "iteration": "iterationId",
                "slice": "sliceId",
                "fix": "fixId",
                "revision": "revision",
            }
            for index_field, project_field in fields.items():
                if index_active.get(index_field) != development.get(project_field):
                    errors.append(
                        f"{label}.active.{index_field} differs from project manifest "
                        f"development.{project_field}"
                    )
    return errors


def _record_path_bases(project_root: Path, sdp_root: Path, value: str) -> list[Path]:
    return [project_root] if value == "SDP" or value.startswith("SDP/") else [sdp_root, project_root]


def _looks_like_path(value: str) -> bool:
    return "/" in value or value.endswith((".md", ".yaml", ".yml", ".json"))


def validate_release_and_fix_records(
    project_root: Path,
    sdp_root: Path,
    schema_root: Path,
    relations: dict[str, Any] | None,
) -> list[str]:
    errors: list[str] = []
    release_schema = load_json(schema_root / "release-record.schema.json")
    fix_schema = load_json(schema_root / "fix-record.schema.json")
    relation_releases = relations.get("releases", {}) if isinstance(relations, dict) else None
    relation_fixes = relations.get("fixes", {}) if isinstance(relations, dict) else None

    release_root = sdp_root / "Releases"
    if release_root.is_dir():
        for path in sorted(release_root.rglob("*.yaml")):
            label = str(path.relative_to(project_root)).replace("\\", "/")
            record, read_errors = _read_yaml(path, label)
            errors += read_errors
            if record is None:
                continue
            errors += _schema_version_errors(
                record, SUPPORTED_RELEASE_RECORD_SCHEMAS, label
            )
            errors += validate_json(record, release_schema, label)
            if not isinstance(record, dict):
                continue
            release_id = record.get("releaseId")
            if isinstance(release_id, str):
                if path.stem != release_id:
                    errors.append(f"{label}: filename must match releaseId {release_id}")
                if isinstance(relation_releases, dict) and release_id not in relation_releases:
                    errors.append(f"{label}: {release_id} is not defined in Relations.releases")
                version = record.get("version")
                if isinstance(version, str) and release_id != f"REL-{version}":
                    errors.append(f"{label}: releaseId does not match version")
                relation_release = (
                    relation_releases.get(release_id)
                    if isinstance(relation_releases, dict)
                    else None
                )
                if isinstance(relation_release, dict):
                    for field in ("version", "state"):
                        if (
                            field in relation_release
                            and relation_release.get(field) != record.get(field)
                        ):
                            errors.append(
                                f"{label}.{field} differs from Relations.releases.{release_id}"
                            )
            release_notes_path = record.get("releaseNotesPath")
            if isinstance(release_notes_path, str):
                errors += validate_existing_path(
                    release_notes_path,
                    _record_path_bases(project_root, sdp_root, release_notes_path),
                    f"{label}.releaseNotesPath",
                )
            migration = record.get("migrationRecord")
            if isinstance(migration, str):
                errors += validate_existing_path(
                    migration,
                    _record_path_bases(project_root, sdp_root, migration),
                    f"{label}.migrationRecord",
                )
            for field, category in (
                ("verificationRecords", "verification"),
                ("reviewRecords", "reviews"),
            ):
                references = record.get(field)
                if not isinstance(references, list):
                    continue
                relation_category = relations.get(category, {}) if isinstance(relations, dict) else None
                for reference in references:
                    if not isinstance(reference, str):
                        continue
                    if _looks_like_path(reference):
                        errors += validate_existing_path(
                            reference,
                            _record_path_bases(project_root, sdp_root, reference),
                            f"{label}.{field}",
                        )
                    elif isinstance(relation_category, dict) and reference not in relation_category:
                        errors.append(f"{label}.{field}: dangling ID {reference}")

    fixes_root = sdp_root / "Fixes"
    if fixes_root.is_dir():
        for path in sorted(fixes_root.rglob("*.yaml")):
            label = str(path.relative_to(project_root)).replace("\\", "/")
            record, read_errors = _read_yaml(path, label)
            errors += read_errors
            if record is None:
                continue
            errors += _schema_version_errors(record, SUPPORTED_FIX_RECORD_SCHEMAS, label)
            errors += validate_json(record, fix_schema, label)
            if not isinstance(record, dict):
                continue
            fix_id = record.get("fixId")
            if isinstance(fix_id, str):
                if path.stem != fix_id:
                    errors.append(f"{label}: filename must match fixId {fix_id}")
                if isinstance(relation_fixes, dict) and fix_id not in relation_fixes:
                    errors.append(f"{label}: {fix_id} is not defined in Relations.fixes")
                fix_match = re.fullmatch(r"(?:FIX|HOTFIX)-(.+)-[0-9]{3}", fix_id)
                if fix_match and fix_match.group(1) != record.get("targetVersion"):
                    errors.append(f"{label}: fixId does not match targetVersion")
            review = record.get("reviewRecord")
            if isinstance(review, str):
                if _looks_like_path(review):
                    errors += validate_existing_path(
                        review,
                        _record_path_bases(project_root, sdp_root, review),
                        f"{label}.reviewRecord",
                    )
                elif isinstance(relations, dict):
                    relation_reviews = relations.get("reviews")
                    if isinstance(relation_reviews, dict) and review not in relation_reviews:
                        errors.append(f"{label}.reviewRecord: dangling ID {review}")
    return errors


def validate_project(project_root: Path, schema_root: Path | None = None) -> list[str]:
    """Validate one installed consuming project without assuming a Toolkit repo layout."""

    project_root = project_root.resolve()
    schemas = (schema_root or Path(__file__).resolve().parents[1] / "schemas").resolve()
    errors: list[str] = []
    if not project_root.is_dir():
        return [f"Project root does not exist or is not a directory: {project_root}"]
    sdp_root = project_root / "SDP"
    if not sdp_root.is_dir():
        return [f"Consuming project SDP directory is missing: {sdp_root}"]

    project_schema = load_json(schemas / "SDP-project-manifest.schema.json")
    installed_schema = load_json(schemas / "installed-toolkit-manifest.schema.json")
    current_schema = load_json(schemas / "current-index.schema.json")
    relations_schema = load_json(schemas / "relations.schema.json")

    project_manifest_path = sdp_root / "SDP-project.manifest.yaml"
    project_manifest, read_errors = _read_yaml(
        project_manifest_path, "SDP/SDP-project.manifest.yaml"
    )
    errors += read_errors
    if project_manifest is not None:
        errors += _schema_version_errors(
            project_manifest,
            SUPPORTED_PROJECT_MANIFEST_SCHEMAS,
            "SDP/SDP-project.manifest.yaml",
        )
        errors += validate_json(
            project_manifest, project_schema, "SDP/SDP-project.manifest.yaml"
        )

    installed_relative: Any = "Framework/installed-toolkit.manifest.yaml"
    if isinstance(project_manifest, dict):
        installed = project_manifest.get("installed")
        if isinstance(installed, dict) and "manifestPath" in installed:
            installed_relative = installed.get("manifestPath")
    path_error = portable_relative_path_error(installed_relative)
    if path_error:
        errors.append(
            "SDP/SDP-project.manifest.yaml.installed.manifestPath: "
            f"{path_error}: {installed_relative!r}"
        )
        installed_path = sdp_root / "Framework/installed-toolkit.manifest.yaml"
    else:
        assert isinstance(installed_relative, str)
        installed_path = sdp_root / Path(*PurePosixPath(installed_relative).parts)
        if not _within(installed_path.resolve(), sdp_root.resolve()):
            errors.append(
                "SDP/SDP-project.manifest.yaml.installed.manifestPath: "
                "resolved path escapes the SDP directory"
            )
            installed_path = sdp_root / "Framework/installed-toolkit.manifest.yaml"
    installed_manifest, read_errors = _read_yaml(
        installed_path,
        str(installed_path.relative_to(project_root)).replace("\\", "/"),
    )
    errors += read_errors
    if installed_manifest is not None:
        installed_label = str(installed_path.relative_to(project_root)).replace("\\", "/")
        errors += _schema_version_errors(
            installed_manifest,
            SUPPORTED_INSTALLED_MANIFEST_SCHEMAS,
            installed_label,
        )
        errors += validate_json(installed_manifest, installed_schema, installed_label)
        if isinstance(installed_manifest, dict):
            toolkit_version = installed_manifest.get("toolkitVersion")
            skills = installed_manifest.get("skills")
            if isinstance(toolkit_version, str) and isinstance(skills, dict):
                for skill_id, expected_version in sorted(skills.items()):
                    if not isinstance(skill_id, str) or not SKILL_ID_PATTERN.fullmatch(skill_id):
                        errors.append(f"{installed_label}.skills: invalid installed skill ID {skill_id!r}")
                        continue
                    if not isinstance(expected_version, str):
                        continue
                    skill_path = project_root / ".codex" / "skills" / skill_id / "SKILL.md"
                    if not skill_path.is_file():
                        errors.append(f"Installed skill is missing: {skill_path}")
                        continue
                    errors += validate_skill_metadata(
                        skill_path, skill_id, expected_version, toolkit_version
                    )

    trace_root = sdp_root / "Traceability"
    current_path = trace_root / "CurrentIndex.yaml"
    relations_path = trace_root / "Relations.yaml"
    ledger_path = trace_root / "Ledger.ndjson"
    current: Any | None = None
    relations: Any | None = None
    if current_path.is_file():
        current, read_errors = _read_yaml(current_path, "SDP/Traceability/CurrentIndex.yaml")
        errors += read_errors
        if current is not None:
            errors += validate_json(
                current, current_schema, "SDP/Traceability/CurrentIndex.yaml"
            )
    if relations_path.is_file():
        relations, read_errors = _read_yaml(relations_path, "SDP/Traceability/Relations.yaml")
        errors += read_errors
        if relations is not None:
            errors += validate_json(
                relations, relations_schema, "SDP/Traceability/Relations.yaml"
            )
            errors += validate_relations_semantics(
                relations, sdp_root, "SDP/Traceability/Relations.yaml"
            )
    if current is not None:
        errors += validate_current_index_semantics(
            current,
            relations,
            "SDP/Traceability/CurrentIndex.yaml",
            project_manifest,
        )
    errors += validate_ledger(
        ledger_path,
        schemas,
        required=False,
        relations=relations if isinstance(relations, dict) else None,
        label="SDP/Traceability/Ledger.ndjson",
    )

    errors += validate_release_notes(
        sdp_root / "RELEASE-NOTES.md", "SDP/RELEASE-NOTES.md", required=True
    )
    errors += validate_release_and_fix_records(
        project_root,
        sdp_root,
        schemas,
        relations if isinstance(relations, dict) else None,
    )
    return errors


def _validate_all_schemas(schema_root: Path) -> list[str]:
    errors: list[str] = []
    for path in sorted(schema_root.glob("*.schema.json")):
        try:
            schema = load_json(path)
            Draft202012Validator.check_schema(schema)
        except (ValueError, SchemaError) as exc:
            errors.append(f"{path}: invalid Draft 2020-12 schema: {exc}")
    return errors


def validate_installation_contract(
    repo: Path,
    toolkit_manifest: dict[str, Any],
    *,
    contract_data: Any | None = None,
    check_installer_integration: bool = True,
) -> list[str]:
    """Validate the portable install manifest plus deterministic source inventory."""

    errors: list[str] = []
    contract_path = repo / "Toolkit/SDP-install.manifest.json"
    schema_path = repo / "Toolkit/schemas/SDP-install-manifest.schema.json"
    if contract_data is None:
        contract, read_errors = _read_json(
            contract_path, "Toolkit/SDP-install.manifest.json"
        )
        errors += read_errors
    else:
        contract = contract_data
    schema, schema_errors = _read_json(
        schema_path, "Toolkit/schemas/SDP-install-manifest.schema.json"
    )
    errors += schema_errors
    if contract is None or schema is None:
        return errors
    errors += validate_json(contract, schema, "Toolkit/SDP-install.manifest.json")
    if not isinstance(contract, dict):
        return errors
    toolkit_facts = toolkit_manifest.get("toolkit", {})
    expected_version = toolkit_facts.get("version")
    if contract.get("toolkitVersion") != expected_version:
        errors.append(
            "Toolkit/SDP-install.manifest.json toolkitVersion differs from SDP.manifest.yaml"
        )

    sources = contract.get("sources")
    if isinstance(sources, dict) and isinstance(sources.get("repositoryRoot"), str):
        repository_root = sources["repositoryRoot"]
        repository_anchor = (contract_path.parent / repository_root).resolve()
        if repository_anchor != repo.resolve():
            errors.append(
                "Toolkit/SDP-install.manifest.json sources.repositoryRoot does not "
                "resolve to the repository root"
            )

    contract_capabilities = contract.get("capabilities")
    manifest_capabilities = toolkit_manifest.get("capabilities")
    if contract_capabilities != manifest_capabilities:
        errors.append(
            "Toolkit/SDP-install.manifest.json capabilities differ from SDP.manifest.yaml"
        )

    generators = contract.get("generators")
    generator_ids: set[str] = set()
    installed_generator: dict[str, Any] | None = None
    if isinstance(generators, list):
        for index, generator in enumerate(generators):
            label = f"Toolkit/SDP-install.manifest.json.generators[{index}]"
            if not isinstance(generator, dict):
                continue
            generator_id = generator.get("id")
            if isinstance(generator_id, str):
                if generator_id in generator_ids:
                    errors.append(f"{label}: duplicate generator ID {generator_id}")
                generator_ids.add(generator_id)
            if generator_id == "installed-toolkit-manifest":
                installed_generator = generator
    required_generators = {"installed-toolkit-manifest", "empty-ledger"}
    if generator_ids != required_generators:
        errors.append(
            "Installation generator IDs differ from the v1 contract: "
            f"expected={sorted(required_generators)}, actual={sorted(generator_ids)}"
        )
    if installed_generator is not None:
        facts = installed_generator.get("facts")
        if isinstance(facts, dict):
            fact_agreements = (
                ("toolkitVersion", expected_version),
                ("frameworkVersion", toolkit_manifest.get("framework", {}).get("version")),
                (
                    "agentsContractVersion",
                    toolkit_manifest.get("agentsContract", {}).get("version"),
                ),
                ("skills", toolkit_manifest.get("skills")),
                ("capabilities", contract_capabilities),
            )
            for field, expected in fact_agreements:
                if facts.get(field) != expected:
                    errors.append(
                        "Toolkit/SDP-install.manifest.json installed generator "
                        f"facts.{field} differs from canonical Toolkit metadata"
                    )
            try:
                if SemVer.parse(str(facts.get("installerVersion"))) < SemVer.parse(
                    str(toolkit_facts.get("minimumInstallerVersion"))
                ):
                    errors.append(
                        "Installation generator installerVersion is lower than "
                        "SDP.manifest.yaml minimumInstallerVersion"
                    )
            except ValueError as exc:
                errors.append(f"installation generator version: {exc}")

    exclusions = contract.get("exclusions")
    excluded: list[tuple[str, str]] = []
    excluded_paths: set[str] = set()
    if isinstance(exclusions, list):
        for index, exclusion in enumerate(exclusions):
            label = f"Toolkit/SDP-install.manifest.json.exclusions[{index}]"
            if not isinstance(exclusion, dict):
                continue
            excluded_path = exclusion.get("path")
            excluded_kind = exclusion.get("kind")
            path_error = portable_relative_path_error(excluded_path)
            if path_error:
                errors.append(f"{label}.path: {path_error}: {excluded_path!r}")
                continue
            assert isinstance(excluded_path, str)
            if excluded_path in excluded_paths:
                errors.append(f"{label}: duplicate exclusion path {excluded_path}")
            excluded_paths.add(excluded_path)
            if isinstance(excluded_kind, str):
                excluded.append((excluded_path, excluded_kind))
            actual = repo / Path(*PurePosixPath(excluded_path).parts)
            if excluded_kind == "file" and not actual.is_file():
                errors.append(f"{label}: excluded file does not exist: {excluded_path}")
            if excluded_kind == "tree" and not actual.is_dir():
                errors.append(f"{label}: excluded tree does not exist: {excluded_path}")

    def is_excluded(path: str) -> bool:
        return any(
            path == excluded_path
            or (kind == "tree" and path.startswith(excluded_path + "/"))
            for excluded_path, kind in excluded
        )

    required_exclusions = {
        "01--Mandate",
        "02--Study",
        "03--Requirements",
        "04--Architecture",
        "05--DesignAnalysis",
        "06--Design",
        "07--Implementation",
        "CodeReview",
        "Fixes",
        "Instructions",
        "Refactors",
        "Releases",
        "Sprints",
        "Traceability",
        "Verification",
        "RELEASE-NOTES.md",
        "SDP.manifest.yaml",
        "SDP-DOCUMENT-GUIDE.md",
        "payload",
        "skills",
        "Toolkit/payload/project-root/AGENTS-project.md.template",
        "Toolkit/payload/sdp-root/AGENT-REMINDERS.md.template",
    }
    missing_exclusions = sorted(required_exclusions - excluded_paths)
    if missing_exclusions:
        errors.append(
            "Installation contract omits required live/legacy exclusions: "
            + ", ".join(missing_exclusions)
        )

    entries = contract.get("entries")
    if not isinstance(entries, list):
        return errors
    seen_ids: set[str] = set()
    seen_destinations: dict[str, str] = {}
    listed_sources: set[str] = set()
    used_generators: set[str] = set()
    allowed_source_roots = (
        "Toolkit/payload/",
        "Toolkit/project-templates/",
        "Toolkit/skills/",
    )
    for index, entry in enumerate(entries):
        label = f"Toolkit/SDP-install.manifest.json.entries[{index}]"
        if not isinstance(entry, dict):
            continue
        entry_id = entry.get("id")
        if isinstance(entry_id, str):
            if entry_id in seen_ids:
                errors.append(f"{label}: duplicate entry ID {entry_id}")
            seen_ids.add(entry_id)
        destination = entry.get("destination")
        if destination is not None:
            path_error = portable_relative_path_error(destination)
            if path_error:
                errors.append(f"{label}.destination: {path_error}: {destination!r}")
            elif isinstance(destination, str):
                destination_key = destination.casefold()
                if destination_key in seen_destinations:
                    errors.append(
                        f"{label}: duplicate/case-colliding destination {destination}; "
                        f"already used by {seen_destinations[destination_key]}"
                    )
                seen_destinations[destination_key] = destination
        source = entry.get("source")
        if source is not None:
            path_error = portable_relative_path_error(source)
            if path_error:
                errors.append(f"{label}.source: {path_error}: {source!r}")
            elif isinstance(source, str):
                listed_sources.add(source)
                errors += validate_existing_path(source, [repo], f"{label}.source")
                if not source.startswith(allowed_source_roots):
                    errors.append(
                        f"{label}.source: copied sources must come from a canonical "
                        "Toolkit payload, project-template, or skill root"
                    )
                if is_excluded(source):
                    errors.append(f"{label}.source: source is explicitly excluded: {source}")
                if source.startswith("Toolkit/project-templates/"):
                    if entry.get("ownership") != "project-owned":
                        errors.append(
                            f"{label}: neutral project-template sources must be project-owned"
                        )
                elif entry.get("ownership") != "toolkit-managed":
                    errors.append(
                        f"{label}: managed payload and skill sources must be toolkit-managed"
                    )
        generator = entry.get("generator")
        if isinstance(generator, str):
            used_generators.add(generator)
            if generator not in generator_ids:
                errors.append(f"{label}.generator: unknown generator {generator}")
        governing = entry.get("governing")
        schema_reference = governing.get("schema") if isinstance(governing, dict) else None
        if isinstance(schema_reference, str):
            errors += validate_existing_path(
                schema_reference, [repo], f"{label}.governing.schema"
            )
        if isinstance(governing, dict):
            capability = governing.get("capability")
            if isinstance(contract_capabilities, list) and capability not in contract_capabilities:
                errors.append(
                    f"{label}.governing.capability: {capability!r} is not declared by the contract"
                )

    unused_generators = sorted(generator_ids - used_generators)
    if unused_generators:
        errors.append("Installation contract defines unused generators: " + ", ".join(unused_generators))

    inventory_roots = (
        repo / "Toolkit/payload",
        repo / "Toolkit/project-templates",
        repo / "Toolkit/skills",
    )
    expected_sources = {
        path.relative_to(repo).as_posix()
        for root in inventory_roots
        if root.is_dir()
        for path in root.rglob("*")
        if path.is_file()
        and not is_excluded(path.relative_to(repo).as_posix())
    }
    unlisted = sorted(expected_sources - listed_sources)
    if unlisted:
        errors.append(
            "Installable source inventory contains unlisted files: " + ", ".join(unlisted)
        )
    outside_inventory = sorted(listed_sources - expected_sources)
    if outside_inventory:
        errors.append(
            "Installation contract lists sources outside the canonical inventory: "
            + ", ".join(outside_inventory)
        )

    project_template_root = repo / "Toolkit/project-templates"
    forbidden_template_names = {
        "Ledger.ndjson",
        "ScrumIterations.md",
        "implementationNotes.md",
        "Handoff.md",
    }
    for path in sorted(project_template_root.rglob("*")):
        if not path.is_file():
            continue
        relative = path.relative_to(project_template_root).as_posix()
        if path.name in forbidden_template_names or re.fullmatch(r"REL-.+\.ya?ml", path.name):
            errors.append(f"Neutral project-template payload contains active record shape: {relative}")
        try:
            if "REL-0.2.0" in path.read_text(encoding="utf-8"):
                errors.append(f"Neutral project-template payload contains Toolkit release ID: {relative}")
        except UnicodeError as exc:
            errors.append(f"Cannot read neutral project template {relative}: {exc}")

    installer_path = repo / "Toolkit/scripts/Install-SDP.ps1"
    if check_installer_integration and installer_path.is_file():
        installer_text = installer_path.read_text(encoding="utf-8")
        if "SDP-install.manifest.json" not in installer_text:
            errors.append("PowerShell installer does not reference the canonical installation manifest")
    return errors


def validate_repository(repo: Path, base_ref: str | None = None) -> list[str]:
    """Validate the SDP Toolkit repository (backward-compatible public name)."""

    repo = repo.resolve()
    errors: list[str] = []
    schema_root = repo / "Toolkit/schemas"
    errors += _validate_all_schemas(schema_root)

    manifest_path = repo / "SDP.manifest.yaml"
    manifest_schema_path = schema_root / "SDP-manifest.schema.json"
    manifest = load_yaml(manifest_path)
    manifest_schema = load_json(manifest_schema_path)
    errors += validate_json(manifest, manifest_schema, "SDP.manifest.yaml")
    if not isinstance(manifest, dict):
        return errors

    toolkit = manifest.get("toolkit", {})
    errors += validate_semver_values(
        [
            toolkit.get("version", ""),
            toolkit.get("releaseTargetVersion", ""),
            toolkit.get("minimumInstallerVersion", ""),
            manifest.get("framework", {}).get("version", ""),
            manifest.get("agentsContract", {}).get("version", ""),
            *manifest.get("skills", {}).values(),
        ],
        "manifest version",
    )
    try:
        if SemVer.parse(toolkit["releaseTargetVersion"]) < SemVer.parse(toolkit["version"]):
            errors.append("Toolkit releaseTargetVersion is lower than toolkit version")
    except (KeyError, TypeError, ValueError):
        pass
    declared_project_schemas = manifest.get("supportedProjectManifestSchemaVersions")
    if isinstance(declared_project_schemas, list):
        unsupported_project_schemas = sorted(
            str(version)
            for version in declared_project_schemas
            if not isinstance(version, str)
            or version not in SUPPORTED_PROJECT_MANIFEST_SCHEMAS
        )
        if unsupported_project_schemas:
            errors.append(
                "SDP.manifest.yaml declares unsupported project manifest schemas: "
                + ", ".join(unsupported_project_schemas)
            )

    expected_skills = manifest.get("skills", {})
    actual_skill_dirs = {
        path.parent.name for path in (repo / "Toolkit/skills").glob("*/SKILL.md")
    }
    if isinstance(expected_skills, dict) and set(expected_skills) != actual_skill_dirs:
        errors.append(
            "Skill set differs between manifest and Toolkit/skills: "
            f"manifest={sorted(expected_skills)}, files={sorted(actual_skill_dirs)}"
        )
    if isinstance(expected_skills, dict):
        for skill_id, expected_version in sorted(expected_skills.items()):
            path = repo / "Toolkit/skills" / skill_id / "SKILL.md"
            if isinstance(expected_version, str) and isinstance(toolkit.get("version"), str):
                errors += validate_skill_metadata(
                    path, skill_id, expected_version, toolkit["version"]
                )

    errors += validate_installation_contract(repo, manifest)

    notes_path = repo / manifest.get("releaseNotesPath", "RELEASE-NOTES.md")
    errors += validate_release_notes(notes_path, "RELEASE-NOTES.md")
    try:
        notes = notes_path.read_text(encoding="utf-8")
        errors += validate_released_immutability(repo, notes, base_ref)
    except OSError:
        pass

    project_schema = load_json(schema_root / "SDP-project-manifest.schema.json")
    project_templates = sorted(
        path
        for path in (repo / "Toolkit").rglob("SDP-project.manifest.yaml")
        if "project-templates" in path.parts or "payload" in path.parts
    )
    if not project_templates:
        errors.append("No canonical SDP-project.manifest.yaml template exists")
    for path in project_templates:
        template = load_yaml(path)
        errors += validate_json(
            template, project_schema, str(path.relative_to(repo)).replace("\\", "/")
        )

    installed_schema = load_json(schema_root / "installed-toolkit-manifest.schema.json")
    installed_example = load_yaml(repo / "examples/installed-toolkit.manifest.example.yaml")
    errors += validate_json(
        installed_example, installed_schema, "examples/installed-toolkit.manifest.example.yaml"
    )
    if isinstance(installed_example, dict):
        installed_agreements = (
            ("toolkitVersion", toolkit.get("version")),
            ("frameworkVersion", manifest.get("framework", {}).get("version")),
            ("agentsContractVersion", manifest.get("agentsContract", {}).get("version")),
            ("skills", expected_skills),
            ("capabilities", manifest.get("capabilities")),
        )
        for field, expected in installed_agreements:
            if installed_example.get(field) != expected:
                errors.append(
                    f"examples/installed-toolkit.manifest.example.yaml {field} "
                    "differs from SDP.manifest.yaml"
                )

    build_schema = load_json(schema_root / "build-identity.schema.json")
    build_example = load_json(repo / "examples/build-identity.example.json")
    errors += validate_json(build_example, build_schema, "examples/build-identity.example.json")

    install_plan_schema = load_json(schema_root / "SDP-install-plan.schema.json")
    install_plan_example = load_json(repo / "examples/install-plan.example.json")
    errors += validate_json(
        install_plan_example,
        install_plan_schema,
        "examples/install-plan.example.json",
    )
    if (
        isinstance(install_plan_example, dict)
        and install_plan_example.get("toolkitVersion") != toolkit.get("version")
    ):
        errors.append(
            "examples/install-plan.example.json toolkitVersion differs from SDP.manifest.yaml"
        )

    release_record_schema = load_json(schema_root / "release-record.schema.json")
    release_templates = sorted((repo / "Toolkit").rglob("ReleaseRecord.yaml"))
    for path in release_templates:
        errors += validate_json(
            load_yaml(path),
            release_record_schema,
            str(path.relative_to(repo)).replace("\\", "/"),
        )
    fix_record_schema = load_json(schema_root / "fix-record.schema.json")
    fix_templates = sorted((repo / "Toolkit").rglob("FixRecord.yaml"))
    for path in fix_templates:
        errors += validate_json(
            load_yaml(path),
            fix_record_schema,
            str(path.relative_to(repo)).replace("\\", "/"),
        )
    current_schema = load_json(schema_root / "current-index.schema.json")
    relations_schema = load_json(schema_root / "relations.schema.json")
    neutral_root = repo / "Toolkit/project-templates/sdp-root"
    neutral_project_manifest = load_yaml(neutral_root / "SDP-project.manifest.yaml")
    neutral_current = load_yaml(neutral_root / "Traceability/CurrentIndex.yaml")
    neutral_relations = load_yaml(neutral_root / "Traceability/Relations.yaml")
    errors += validate_json(
        neutral_current,
        current_schema,
        "Toolkit/project-templates/sdp-root/Traceability/CurrentIndex.yaml",
    )
    errors += validate_json(
        neutral_relations,
        relations_schema,
        "Toolkit/project-templates/sdp-root/Traceability/Relations.yaml",
    )
    errors += validate_current_index_semantics(
        neutral_current,
        neutral_relations,
        "Toolkit/project-templates/sdp-root/Traceability/CurrentIndex.yaml",
        neutral_project_manifest,
    )
    errors += validate_release_notes(
        neutral_root / "RELEASE-NOTES.md",
        "Toolkit/project-templates/sdp-root/RELEASE-NOTES.md",
    )
    current = load_yaml(repo / "Traceability/CurrentIndex.yaml")
    relations = load_yaml(repo / "Traceability/Relations.yaml")
    errors += validate_json(current, current_schema, "Traceability/CurrentIndex.yaml")
    errors += validate_json(relations, relations_schema, "Traceability/Relations.yaml")
    errors += validate_relations_semantics(relations, repo, "Traceability/Relations.yaml")
    errors += validate_current_index_semantics(
        current, relations, "Traceability/CurrentIndex.yaml"
    )
    if isinstance(current, dict):
        current_release = current.get("release")
        if isinstance(current_release, dict):
            toolkit_release_agreements = (
                ("targetVersion", toolkit.get("releaseTargetVersion")),
                ("state", toolkit.get("releaseState")),
            )
            for field, expected in toolkit_release_agreements:
                if current_release.get(field) != expected:
                    errors.append(
                        f"Traceability/CurrentIndex.yaml release.{field} "
                        "differs from SDP.manifest.yaml"
                    )
            target = current_release.get("targetVersion")
            active_release_id = current_release.get("activeReleaseId")
            if isinstance(target, str) and active_release_id != f"REL-{target}":
                errors.append(
                    "Traceability/CurrentIndex.yaml activeReleaseId does not match targetVersion"
                )
    errors += validate_release_and_fix_records(
        repo,
        repo,
        schema_root,
        relations if isinstance(relations, dict) else None,
    )
    errors += validate_ledger(
        repo / "Traceability/Ledger.ndjson",
        schema_root,
        required=True,
        relations=relations if isinstance(relations, dict) else None,
        label="Traceability/Ledger.ndjson",
    )
    errors += validate_ledger(
        repo / "examples/ledger-events.ndjson.example",
        schema_root,
        required=True,
        relations=None,
        label="examples/ledger-events.ndjson.example",
    )
    errors += validate_ledger(
        repo / "examples/release-events.ndjson.example",
        schema_root,
        required=True,
        relations=None,
        label="examples/release-events.ndjson.example",
    )
    return errors


def main(argv: Sequence[str] | None = None) -> int:
    parser = argparse.ArgumentParser(
        description="Validate an SDP Toolkit repository or installed consuming project."
    )
    parser.add_argument(
        "--mode",
        choices=("toolkit", "project"),
        default="toolkit",
        help="Validation target (default: toolkit for backward compatibility)",
    )
    parser.add_argument(
        "--repo",
        type=Path,
        default=Path(__file__).resolve().parents[2],
        help="SDP Toolkit repository root in toolkit mode",
    )
    parser.add_argument("--project-root", type=Path, help="Consuming-project root in project mode")
    parser.add_argument("--base-ref", help="Git ref used to verify immutable released note sections")
    args = parser.parse_args(argv)

    if args.mode == "project" and args.project_root is None:
        parser.error("--project-root is required when --mode project")
    if args.mode == "project" and args.base_ref:
        parser.error("--base-ref is only valid in toolkit mode")

    try:
        if args.mode == "project":
            errors = validate_project(args.project_root.resolve())
            success_message = f"SDP consuming-project validation passed: {args.project_root.resolve()}"
        else:
            errors = validate_repository(args.repo.resolve(), args.base_ref)
            success_message = "SDP Toolkit validation passed"
    except (OSError, ValueError, KeyError, TypeError, SchemaError) as exc:
        errors = [str(exc)]
        success_message = ""

    if errors:
        print("SDP validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print(success_message)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
