#!/usr/bin/env python3
"""Deterministic validation for SDP Toolkit and consuming-project contracts."""

from __future__ import annotations

import argparse
import base64
import binascii
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


SEMVER_PATTERN_TEXT = (
    r"^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)"
    r"(?:-((?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*)"
    r"(?:\.(?:0|[1-9][0-9]*|[0-9]*[A-Za-z-][0-9A-Za-z-]*))*))?"
    r"(?:\+([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?$"
)
SEMVER_PATTERN = re.compile(SEMVER_PATTERN_TEXT)
WINDOWS_DRIVE_PATTERN = re.compile(r"^[A-Za-z]:")
WINDOWS_RESERVED_SEGMENT_PATTERN = re.compile(
    r"^(?:CON|PRN|AUX|NUL|CLOCK\$|CONIN\$|CONOUT\$|COM[1-9¹²³]|LPT[1-9¹²³])(?:\..*)?$",
    re.IGNORECASE,
)
WINDOWS_FORBIDDEN_PATH_CHARACTERS = frozenset('<>:"|?*~')
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
    "identifierConvention",
    "releaseNotes",
    "manifest",
    "releaseRecord",
    "migrationRecord",
)

CANONICAL_INSTALL_ORDERING_POLICY = "migration-first-manifest-order-v1"
INSTALL_FAILURE_CLASSES = frozenset(
    {
        "install-manifest-invalid",
        "agents-migration-destination-content-mismatch",
        "agents-migration-destination-unsupported-object",
        "agents-migration-source-changed",
        "agents-migration-destination-changed",
    }
)
REQUIRED_INSTALL_CONFORMANCE_CATEGORIES = frozenset(
    {
        "empty",
        "legacy",
        "upgrade",
        "force",
        "initialize",
        "repeat-initialize",
        "archive",
        "error",
    }
)

PLAN_REASON_CONDITIONS: dict[str, tuple[str, bool]] = {
    "missing-target": ("create", True),
    "missing-generated-target": ("generate", True),
    "content-matches": ("unchanged", False),
    "missing-only-content": ("preserve", False),
    "managed-content-differs": ("preserve", False),
    "backup-before-replace": ("backup", True),
    "refresh-managed-content": ("replace", True),
    "refresh-generated-content": ("generate", True),
    "migrate-existing-agents": ("migrate", True),
    "preserve-existing-agents-conflict": ("migrate", True),
    "malformed-project-manifest": ("block", False),
    "unsupported-project-schema": ("block", False),
    "malformed-installed-manifest": ("block", False),
    "unsupported-installed-schema": ("block", False),
    "downgrade-blocked": ("block", False),
}

BLOCK_REASON_ENTRY_IDS = {
    "malformed-project-manifest": "project-manifest",
    "unsupported-project-schema": "project-manifest",
    "malformed-installed-manifest": "generated-installed-toolkit-manifest",
    "unsupported-installed-schema": "generated-installed-toolkit-manifest",
    "downgrade-blocked": "generated-installed-toolkit-manifest",
}

CANONICAL_GOVERNING_CAPABILITIES = {
    "Toolkit/schemas/fix-record.schema.json": "sdp.release.v1",
    "Toolkit/schemas/release-record.schema.json": "sdp.release.v1",
    "Toolkit/schemas/SDP-project-manifest.schema.json": "sdp.manifest.v1",
    "Toolkit/schemas/installed-toolkit-manifest.schema.json": "sdp.manifest.v1",
    "Toolkit/schemas/current-index.schema.json": "sdp.traceability.current-index.v1",
    "Toolkit/schemas/relations.schema.json": "sdp.traceability.relations.v1",
    "Toolkit/schemas/ledger-event.schema.json": "sdp.traceability.ledger-events.v1",
}
CANONICAL_GOVERNING_CAPABILITIES_BY_KEY = {
    schema.casefold(): capability
    for schema, capability in CANONICAL_GOVERNING_CAPABILITIES.items()
}


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


def validate_install_plan(
    plan: Any,
    schema: Any,
    label: str,
    installation_contract: Any | None = None,
) -> list[str]:
    """Validate an install plan's JSON shape and cross-action/manifest semantics."""

    errors = validate_json(plan, schema, label)
    if not isinstance(plan, dict):
        return errors

    actions = plan.get("actions")
    if not isinstance(actions, list):
        return errors

    sequences = [
        action.get("sequence") if isinstance(action, dict) else None for action in actions
    ]
    expected_sequences = list(range(1, len(actions) + 1))
    if sequences != expected_sequences:
        errors.append(
            f"{label}.actions: sequence values must be unique, ordered, and contiguous "
            f"from 1; actual={sequences!r}"
        )

    entries_by_id: dict[str, dict[str, Any]] = {}
    manifest_index_by_id: dict[str, int] = {}
    if isinstance(installation_contract, dict):
        contract_schema_version = installation_contract.get("schemaVersion")
        if plan.get("manifestSchemaVersion") != contract_schema_version:
            errors.append(
                f"{label}.manifestSchemaVersion differs from the installation contract"
            )
        contract_toolkit_version = installation_contract.get("toolkitVersion")
        if plan.get("toolkitVersion") != contract_toolkit_version:
            errors.append(f"{label}.toolkitVersion differs from the installation contract")
        contract_ordering_policy = installation_contract.get("orderingPolicy")
        if plan.get("orderingPolicy") != contract_ordering_policy:
            errors.append(f"{label}.orderingPolicy differs from the installation contract")
        contract_entries = installation_contract.get("entries")
        if isinstance(contract_entries, list):
            entries_by_id = {
                entry["id"]: entry
                for entry in contract_entries
                if isinstance(entry, dict) and isinstance(entry.get("id"), str)
            }
            manifest_index_by_id = {
                entry["id"]: index
                for index, entry in enumerate(contract_entries)
                if isinstance(entry, dict) and isinstance(entry.get("id"), str)
            }

    if plan.get("orderingPolicy") != CANONICAL_INSTALL_ORDERING_POLICY:
        errors.append(
            f"{label}.orderingPolicy must be {CANONICAL_INSTALL_ORDERING_POLICY!r}"
        )

    toolkit_version = plan.get("toolkitVersion")
    installed_version = plan.get("installedToolkitVersion")
    block_count = 0
    migration_count = 0
    ordinary_actions_started = False
    last_manifest_index = -1
    ordinary_action_counts: dict[str, int] = {}
    for index, action in enumerate(actions):
        action_label = f"{label}.actions[{index}]"
        if not isinstance(action, dict):
            continue
        action_name = action.get("action")
        reason = action.get("reason")
        mutates_target = action.get("mutatesTarget")
        entry_id = action.get("entryId")
        if reason not in PLAN_REASON_CONDITIONS:
            errors.append(f"{action_label}.reason: unsupported v1 reason {reason!r}")
        else:
            expected_action, expected_mutation = PLAN_REASON_CONDITIONS[reason]
            if action_name != expected_action or mutates_target is not expected_mutation:
                errors.append(
                    f"{action_label}: reason {reason!r} requires action "
                    f"{expected_action!r} and mutatesTarget={expected_mutation}"
                )
        if action_name == "block":
            block_count += 1
            expected_entry_id = BLOCK_REASON_ENTRY_IDS.get(reason)
            if entry_id != expected_entry_id:
                errors.append(
                    f"{action_label}.entryId contradicts block reason {reason!r}; "
                    f"expected {expected_entry_id!r}"
                )
        if action.get("newToolkitVersion") != toolkit_version:
            errors.append(
                f"{action_label}.newToolkitVersion differs from plan toolkitVersion"
            )
        if action.get("oldToolkitVersion") != installed_version:
            errors.append(
                f"{action_label}.oldToolkitVersion differs from plan installedToolkitVersion"
            )

        entry = entries_by_id.get(entry_id) if isinstance(entry_id, str) else None
        if entries_by_id and entry is None:
            errors.append(f"{action_label}.entryId: unknown installation entry {entry_id!r}")
            continue
        if entry is None:
            continue

        source = action.get("source")
        generator = action.get("generator")
        target_source = action.get("targetSource")
        target_source_sha256 = action.get("targetSourceSha256")
        destination_precondition = action.get("destinationPrecondition")
        destination = action.get("destination")
        ownership = action.get("ownership")
        if action_name == "migrate":
            if ordinary_actions_started:
                errors.append(
                    f"{action_label}: migrations must precede every ordinary manifest action"
                )
            migration_count += 1
            if migration_count > 1:
                errors.append(
                    f"{action_label}: installation contract v1 permits at most one migration"
                )
            if source is not None or generator is not None:
                errors.append(f"{action_label}: migrate actions must not use manifest content")
            if (
                not isinstance(target_source_sha256, str)
                or re.fullmatch(r"[0-9a-f]{64}", target_source_sha256) is None
            ):
                errors.append(
                    f"{action_label}.targetSourceSha256 must identify exact source bytes"
                )
            if destination_precondition != "absent":
                errors.append(
                    f"{action_label}.destinationPrecondition must be 'absent'"
                )
            managed_agents_entry = entries_by_id.get("managed-agents", {})
            project_agents_entry = entries_by_id.get("project-agents", {})
            managed_agents_destination = managed_agents_entry.get("destination", "AGENTS.md")
            project_agents_destination = project_agents_entry.get(
                "destination", "AGENTS-project.md"
            )
            migration_contracts = {
                "migrate-existing-agents": {
                    "entryId": "project-agents",
                    "targetSource": managed_agents_destination,
                    "destination": project_agents_destination,
                },
                "preserve-existing-agents-conflict": {
                    "entryId": "managed-agents",
                    "targetSource": managed_agents_destination,
                    "destinationPattern": re.compile(
                        r"^AGENTS-project\.migration-sha256-[0-9a-f]{64}\.md$"
                    ),
                },
            }
            expected = migration_contracts.get(reason)
            if expected is not None:
                if entry_id != expected["entryId"]:
                    errors.append(
                        f"{action_label}.entryId contradicts migration reason {reason!r}"
                    )
                if target_source != expected["targetSource"]:
                    errors.append(
                        f"{action_label}.targetSource contradicts migration reason {reason!r}"
                    )
                exact_destination = expected.get("destination")
                destination_pattern = expected.get("destinationPattern")
                if exact_destination is not None and destination != exact_destination:
                    errors.append(
                        f"{action_label}.destination contradicts migration reason {reason!r}"
                    )
                if isinstance(destination_pattern, re.Pattern):
                    expected_hash_destination = (
                        "AGENTS-project.migration-sha256-"
                        f"{target_source_sha256}.md"
                    )
                    if (
                        not isinstance(destination, str)
                        or destination_pattern.fullmatch(destination) is None
                        or destination != expected_hash_destination
                    ):
                        errors.append(
                            f"{action_label}.destination must use the exact source-hash "
                            "migration name"
                        )
            if ownership != "project-owned":
                errors.append(f"{action_label}.ownership: migrations are project-owned")
            continue

        if action_name != "block":
            ordinary_actions_started = True
            manifest_index = manifest_index_by_id.get(entry_id, -1)
            if manifest_index < last_manifest_index:
                errors.append(
                    f"{action_label}: ordinary actions must follow installation-manifest "
                    "array order"
                )
            last_manifest_index = max(last_manifest_index, manifest_index)
            if isinstance(entry_id, str):
                ordinary_action_counts[entry_id] = (
                    ordinary_action_counts.get(entry_id, 0) + 1
                )
                if ordinary_action_counts[entry_id] > 2:
                    errors.append(
                        f"{action_label}: manifest entry produced more than two actions"
                    )
                if ordinary_action_counts[entry_id] == 2:
                    previous = actions[index - 1] if index else None
                    if (
                        not isinstance(previous, dict)
                        or previous.get("action") != "backup"
                        or previous.get("entryId") != entry_id
                        or action_name not in {"replace", "generate"}
                    ):
                        errors.append(
                            f"{action_label}: manifest entry produced a non-canonical "
                            "second action"
                        )

        if target_source is not None:
            errors.append(f"{action_label}.targetSource must be null for non-migrate actions")
        if target_source_sha256 is not None:
            errors.append(
                f"{action_label}.targetSourceSha256 must be null for non-migrate actions"
            )
        if destination_precondition is not None:
            errors.append(
                f"{action_label}.destinationPrecondition must be null for non-migrate actions"
            )
        if action_name == "block":
            if source is not None or generator is not None:
                errors.append(f"{action_label}: block actions must not read manifest content")
        else:
            for field in ("source", "generator"):
                if action.get(field) != entry.get(field):
                    errors.append(
                        f"{action_label}.{field} differs from installation entry {entry_id}"
                    )
        for field in ("destination", "ownership"):
            if action.get(field) != entry.get(field):
                errors.append(
                    f"{action_label}.{field} differs from installation entry {entry_id}"
                )
        if action_name == "generate" and entry.get("kind") != "generated":
            errors.append(f"{action_label}: generate action references a copied entry")
        if action_name in {"create", "replace"} and entry.get("kind") != "copied":
            errors.append(f"{action_label}: {action_name} action references a generated entry")

        if action_name == "backup":
            mutation = actions[index + 1] if index + 1 < len(actions) else None
            if (
                not isinstance(mutation, dict)
                or mutation.get("entryId") != entry_id
                or mutation.get("action") not in {"replace", "generate"}
            ):
                errors.append(
                    f"{action_label}: backup must immediately precede its matching mutation"
                )
            else:
                identity_fields = (
                    "entryId",
                    "source",
                    "generator",
                    "targetSource",
                    "targetSourceSha256",
                    "destinationPrecondition",
                    "destination",
                    "ownership",
                    "oldToolkitVersion",
                    "newToolkitVersion",
                )
                for field in identity_fields:
                    if action.get(field) != mutation.get(field):
                        errors.append(
                            f"{action_label}: backup/mutation pair disagrees on {field}"
                        )
        if (
            action_name in {"replace", "generate"}
            and reason in {"refresh-managed-content", "refresh-generated-content"}
            and entry.get("backupPolicy") == "before-replace"
        ):
            backup = actions[index - 1] if index else None
            if (
                not isinstance(backup, dict)
                or backup.get("action") != "backup"
                or backup.get("entryId") != entry_id
            ):
                errors.append(
                    f"{action_label}: mutation is missing its immediately preceding backup"
                )

    expected_can_apply = block_count == 0
    if plan.get("canApply") != expected_can_apply:
        errors.append(f"{label}.canApply contradicts the presence of block actions")
    if block_count and (block_count != 1 or len(actions) != 1):
        errors.append(
            f"{label}.actions: a v1 blocked plan must contain exactly one block action "
            "and no other actions"
        )
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
    pure = PurePosixPath(value)
    if ".." in pure.parts:
        return "must not contain parent traversal"
    normalized = pure.as_posix()
    if normalized in {"", "."} or normalized != value:
        return "must be normalized"
    for segment in pure.parts:
        if any(ord(character) < 32 or ord(character) == 127 for character in segment):
            return "must not contain control characters"
        forbidden = sorted(set(segment) & WINDOWS_FORBIDDEN_PATH_CHARACTERS)
        if forbidden:
            return f"must not contain Windows-invalid characters: {''.join(forbidden)}"
        if segment.endswith((".", " ")):
            return "segments must not end with a dot or space"
        if WINDOWS_RESERVED_SEGMENT_PATTERN.fullmatch(segment):
            return f"must not contain a Windows reserved device segment: {segment}"
        if segment.casefold() == ".git":
            return "must not address a .git segment"
    return None


def portable_path_key(value: str) -> tuple[str, ...]:
    """Return the collision key shared by supported case-insensitive filesystems."""

    return tuple(part.casefold() for part in PurePosixPath(value).parts)


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
        (r"^REV-", "reviews"),
        (r"^VER-", "verification"),
    )
    for pattern, category in patterns:
        if re.match(pattern, subject_id):
            return category
    return None


def _id_references(value: Any) -> list[str]:
    if isinstance(value, str):
        return [value]
    if isinstance(value, list):
        return [item for item in value if isinstance(item, str)]
    return []


def _release_version_from_id(release_id: Any) -> str | None:
    if not isinstance(release_id, str) or not release_id.startswith("REL-"):
        return None
    version = release_id[4:]
    try:
        SemVer.parse(version)
    except ValueError:
        return None
    return version


def validate_publication_identity(
    release_id: Any,
    version: Any,
    state: Any,
    commit: Any,
    tag: Any,
    publication: Any,
    label: str,
    *,
    commit_field: str = "releaseCommit",
    publication_field: str = "githubRelease",
) -> list[str]:
    errors: list[str] = []
    if isinstance(release_id, str) and isinstance(version, str):
        if release_id != f"REL-{version}":
            errors.append(f"{label}: release ID does not match version")
    if isinstance(version, str) and isinstance(tag, str) and tag != f"v{version}":
        errors.append(f"{label}: publication tag {tag!r} does not match version {version}")
    if state == "unreleased":
        for field, value in (
            (commit_field, commit),
            ("gitTag", tag),
            (publication_field, publication),
        ):
            if value is not None:
                errors.append(f"{label}.{field}: must be null while release is unreleased")
    if state in {"released", "yanked"}:
        for field, value in (
            (commit_field, commit),
            ("gitTag", tag),
            (publication_field, publication),
        ):
            if not isinstance(value, str) or not value:
                errors.append(f"{label}.{field}: is required for {state} releases")
    return errors


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
            release_version = _release_version_from_id(event.get("releaseId"))
            payload = event.get("payload")
            tag = payload.get("tag") if isinstance(payload, dict) else None
            if (
                event_type in {"release-tag-created", "release-published"}
                and release_version is not None
                and isinstance(tag, str)
                and tag != f"v{release_version}"
            ):
                errors.append(
                    f"{line_label} release event: tag {tag!r} does not match "
                    f"releaseId {event.get('releaseId')}"
                )
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


def validate_relations_semantics(
    relations: Any,
    sdp_root: Path,
    label: str,
    project_root: Path | None = None,
) -> list[str]:
    if not isinstance(relations, dict):
        return []
    root = project_root or sdp_root
    errors: list[str] = []
    for category, records in relations.items():
        if not isinstance(records, dict):
            continue
        for record_id, relation in records.items():
            if not isinstance(relation, dict):
                continue
            for field in RELATION_PATH_FIELDS:
                if field in relation:
                    value = relation[field]
                    if isinstance(value, str) and value.casefold().endswith(".yml"):
                        errors.append(
                            f"{label}.{category}.{record_id}.{field}: "
                            "YAML record paths must use the canonical .yaml extension"
                        )
                    errors += validate_existing_path(
                        value,
                        _record_path_bases(root, sdp_root, value)
                        if isinstance(value, str)
                        else [sdp_root],
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
            "reviews": "reviews",
            "verification": "verification",
            "mandate": "mandates",
        },
        "fixes": {
            "release": "releases",
            "review": "reviews",
            "reviews": "reviews",
            "verification": "verification",
        },
        "reviews": {
            "slice": "slices",
            "fix": "fixes",
            "release": "releases",
            "verification": "verification",
            "resolvedBy": "reviews",
            "resolves": "reviews",
        },
        "verification": {
            "slice": "slices",
            "fix": "fixes",
            "release": "releases",
            "releaseContext": "releases",
            "review": "reviews",
            "reviewerConfirmation": "reviews",
        },
        "migrations": {"slice": "slices", "fix": "fixes", "release": "releases"},
        "releases": {
            "sprints": "sprints",
            "fixes": "fixes",
            "slices": "slices",
            "reviews": "reviews",
            "verification": "verification",
            "migrations": "migrations",
            "mandate": "mandates",
        },
        "mandates": {
            "sprint": "sprints",
            "outcomes": "outcomes",
            "boundaries": "boundaries",
            "successCriteria": "successCriteria",
            "assumptions": "assumptions",
            "questions": "questions",
        },
        "outcomes": {"mandate": "mandates"},
        "boundaries": {"mandate": "mandates"},
        "successCriteria": {"mandate": "mandates"},
        "assumptions": {"mandate": "mandates"},
        "questions": {"mandate": "mandates"},
    }
    for category, fields in reference_rules.items():
        records = relations.get(category)
        if not isinstance(records, dict):
            continue
        for record_id, relation in records.items():
            if not isinstance(relation, dict):
                continue
            for field, target_category in fields.items():
                references = _id_references(relation.get(field))
                target = relations.get(target_category)
                if not isinstance(target, dict):
                    continue
                for reference in references:
                    if reference not in target:
                        errors.append(
                            f"{label}.{category}.{record_id}.{field}: dangling ID {reference}; "
                            f"not found in {target_category}"
                        )

    reverse_rules: tuple[
        tuple[str, tuple[str, ...], str, tuple[str, ...]], ...
    ] = (
        ("slices", ("review", "reviews"), "reviews", ("slice",)),
        ("reviews", ("slice",), "slices", ("review", "reviews")),
        ("slices", ("verification",), "verification", ("slice",)),
        ("verification", ("slice",), "slices", ("verification",)),
        ("fixes", ("review", "reviews"), "reviews", ("fix",)),
        ("reviews", ("fix",), "fixes", ("review", "reviews")),
        ("fixes", ("verification",), "verification", ("fix",)),
        ("verification", ("fix",), "fixes", ("verification",)),
        ("migrations", ("release",), "releases", ("migrations",)),
        ("releases", ("migrations",), "migrations", ("release",)),
        ("reviews", ("resolvedBy",), "reviews", ("resolves",)),
        ("reviews", ("resolves",), "reviews", ("resolvedBy",)),
    )
    for source_category, source_fields, target_category, reverse_fields in reverse_rules:
        source_records = relations.get(source_category)
        target_records = relations.get(target_category)
        if not isinstance(source_records, dict) or not isinstance(target_records, dict):
            continue
        for source_id, source_relation in source_records.items():
            if not isinstance(source_relation, dict):
                continue
            for source_field in source_fields:
                for target_id in _id_references(source_relation.get(source_field)):
                    target_relation = target_records.get(target_id)
                    if not isinstance(target_relation, dict):
                        continue
                    if not any(
                        source_id in _id_references(target_relation.get(reverse_field))
                        for reverse_field in reverse_fields
                    ):
                        rendered_fields = " or ".join(reverse_fields)
                        errors.append(
                            f"{label}.{source_category}.{source_id}.{source_field}: "
                            f"missing reverse link from {target_category}.{target_id}."
                            f"{rendered_fields}"
                        )

    release_records = relations.get("releases")
    if isinstance(release_records, dict):
        for release_id, release_relation in release_records.items():
            if not isinstance(release_relation, dict):
                continue
            errors += validate_publication_identity(
                release_id,
                release_relation.get("version"),
                release_relation.get("state"),
                release_relation.get("releaseCommit"),
                release_relation.get("gitTag"),
                release_relation.get("githubRelease"),
                f"{label}.releases.{release_id}",
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
    release_record_paths: dict[str, Path] = {}
    if release_root.is_dir():
        release_paths: list[Path] = []
        for path in sorted(release_root.rglob("*")):
            if not path.is_file() or path.suffix.casefold() not in {".yaml", ".yml"}:
                continue
            label = str(path.relative_to(project_root)).replace("\\", "/")
            if path.suffix != ".yaml":
                errors.append(
                    f"{label}: release records must use the canonical .yaml extension"
                )
            else:
                release_paths.append(path)
        for path in release_paths:
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
                if release_id in release_record_paths:
                    errors.append(f"{label}: duplicate release record ID {release_id}")
                release_record_paths[release_id] = path.resolve()
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
                    relation_agreements = {
                        "version": "version",
                        "state": "state",
                        "releaseNotes": "releaseNotesPath",
                        "verification": "verificationRecords",
                        "reviews": "reviewRecords",
                        "releaseCommit": "releasePreparationCommit",
                        "gitTag": "gitTag",
                        "githubRelease": "githubReleaseUrl",
                    }
                    for relation_field, record_field in relation_agreements.items():
                        relation_value = relation_release.get(relation_field)
                        record_value = record.get(record_field)
                        if (
                            relation_field == "releaseNotes"
                            and isinstance(relation_value, str)
                            and isinstance(record_value, str)
                        ):
                            relation_base = _record_path_bases(
                                project_root, sdp_root, relation_value
                            )[0]
                            record_base = _record_path_bases(
                                project_root, sdp_root, record_value
                            )[0]
                            agrees = (
                                relation_base / Path(*PurePosixPath(relation_value).parts)
                            ).resolve() == (
                                record_base / Path(*PurePosixPath(record_value).parts)
                            ).resolve()
                        else:
                            agrees = relation_value == record_value
                        if not agrees:
                            errors.append(
                                f"{label}.{record_field} differs from "
                                f"Relations.releases.{release_id}.{relation_field}"
                            )
                    relation_record_path = relation_release.get("releaseRecord")
                    if isinstance(relation_record_path, str):
                        relative_record_path = Path(
                            *PurePosixPath(relation_record_path).parts
                        )
                        relation_candidates = {
                            (base.resolve() / relative_record_path).resolve()
                            for base in _record_path_bases(
                                project_root, sdp_root, relation_record_path
                            )
                        }
                        if path.resolve() not in relation_candidates:
                            errors.append(
                                f"{label}: file identity differs from "
                                f"Relations.releases.{release_id}.releaseRecord"
                            )
                    else:
                        errors.append(
                            f"{label}: Relations.releases.{release_id}.releaseRecord "
                            "is required for record identity agreement"
                        )
                errors += validate_publication_identity(
                    release_id,
                    record.get("version"),
                    record.get("state"),
                    record.get("releasePreparationCommit"),
                    record.get("gitTag"),
                    record.get("githubReleaseUrl"),
                    label,
                    commit_field="releasePreparationCommit",
                    publication_field="githubReleaseUrl",
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

    if isinstance(relation_releases, dict):
        for release_id, relation in relation_releases.items():
            relation_label = f"Relations.releases.{release_id}.releaseRecord"
            release_record = relation.get("releaseRecord") if isinstance(relation, dict) else None
            if not isinstance(release_record, str) or not release_record.endswith(".yaml"):
                errors.append(f"{relation_label}: must reference a canonical .yaml release record")
                continue
            discovered = release_record_paths.get(release_id)
            if discovered is None:
                errors.append(f"{relation_label}: no validated release record exists")
                continue
            relative = Path(*PurePosixPath(release_record).parts)
            candidates = {
                (base.resolve() / relative).resolve()
                for base in _record_path_bases(project_root, sdp_root, release_record)
            }
            if discovered not in candidates:
                errors.append(f"{relation_label}: does not identify the governed release record")

    fixes_root = sdp_root / "Fixes"
    fix_record_paths: dict[str, Path] = {}
    if fixes_root.is_dir():
        fix_paths: list[Path] = []
        for path in sorted(fixes_root.rglob("*")):
            if not path.is_file() or path.suffix.casefold() not in {".yaml", ".yml"}:
                continue
            label = str(path.relative_to(project_root)).replace("\\", "/")
            if path.suffix != ".yaml":
                errors.append(f"{label}: Fix records must use the canonical .yaml extension")
            else:
                fix_paths.append(path)
        for path in fix_paths:
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
                if fix_id in fix_record_paths:
                    errors.append(f"{label}: duplicate Fix record ID {fix_id}")
                fix_record_paths[fix_id] = path.resolve()
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
    if isinstance(relation_fixes, dict):
        for fix_id, relation in relation_fixes.items():
            relation_label = f"Relations.fixes.{fix_id}.path"
            fix_record = relation.get("path") if isinstance(relation, dict) else None
            if not isinstance(fix_record, str) or not fix_record.endswith(".yaml"):
                errors.append(f"{relation_label}: must reference a canonical .yaml Fix record")
                continue
            discovered = fix_record_paths.get(fix_id)
            if discovered is None:
                errors.append(f"{relation_label}: no validated Fix record exists")
                continue
            relative = Path(*PurePosixPath(fix_record).parts)
            candidates = {
                (base.resolve() / relative).resolve()
                for base in _record_path_bases(project_root, sdp_root, fix_record)
            }
            if discovered not in candidates:
                errors.append(f"{relation_label}: does not identify the governed Fix record")
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
        if isinstance(project_manifest, dict):
            release = project_manifest.get("release")
            if isinstance(release, dict):
                current_version = release.get("currentVersion")
                latest_tag = release.get("latestTag")
                latest_commit = release.get("latestCommit")
                if (latest_tag is None) != (latest_commit is None):
                    errors.append(
                        "SDP/SDP-project.manifest.yaml.release: latestTag and "
                        "latestCommit must both be null or both be populated"
                    )
                if (
                    isinstance(current_version, str)
                    and isinstance(latest_tag, str)
                    and latest_tag != f"v{current_version}"
                ):
                    errors.append(
                        "SDP/SDP-project.manifest.yaml.release.latestTag does not "
                        "match release.currentVersion"
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
                relations,
                sdp_root,
                "SDP/Traceability/Relations.yaml",
                project_root,
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


def validate_install_conformance_package(
    repo: Path,
    installation_contract: Any | None = None,
    plan_schema: Any | None = None,
) -> list[str]:
    """Validate the language-neutral install-v1 scenario package and authorities."""

    root = repo / "Toolkit/conformance/install-v1"
    index_path = root / "scenarios.json"
    index_schema_path = root / "scenario-index.schema.json"
    index, index_errors = _read_json(index_path, "install-v1 scenarios")
    index_schema, schema_errors = _read_json(
        index_schema_path, "install-v1 scenario index schema"
    )
    errors = index_errors + schema_errors
    if index is None or index_schema is None:
        return errors
    try:
        Draft202012Validator.check_schema(index_schema)
    except SchemaError as exc:
        errors.append(
            f"Toolkit/conformance/install-v1/scenario-index.schema.json: "
            f"invalid Draft 2020-12 schema: {exc}"
        )
        return errors
    errors += validate_json(index, index_schema, "install-v1 scenarios")
    if not isinstance(index, dict):
        return errors

    contract = installation_contract
    if contract is None:
        contract, read_errors = _read_json(
            repo / "Toolkit/SDP-install.manifest.json",
            "Toolkit/SDP-install.manifest.json",
        )
        errors += read_errors
    schema = plan_schema
    if schema is None:
        schema, read_errors = _read_json(
            repo / "Toolkit/schemas/SDP-install-plan.schema.json",
            "Toolkit/schemas/SDP-install-plan.schema.json",
        )
        errors += read_errors
    if index.get("orderingPolicy") != CANONICAL_INSTALL_ORDERING_POLICY:
        errors.append("install-v1 scenarios orderingPolicy is not canonical")
    failure_classes = index.get("failureClasses")
    if isinstance(failure_classes, list) and set(failure_classes) != set(
        INSTALL_FAILURE_CLASSES
    ):
        errors.append(
            "install-v1 scenarios failureClasses differ from the closed v1 vocabulary"
        )

    scenarios = index.get("scenarios")
    if not isinstance(scenarios, list):
        return errors
    seen_ids: set[str] = set()
    seen_expected_paths: set[tuple[str, ...]] = set()
    referenced_expected_paths: set[str] = set()
    covered_categories: set[str] = set()
    volatile_timestamp = re.compile(
        r"\b[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}"
    )

    def check_portable(value: Any, label: str) -> None:
        path_error = portable_relative_path_error(value)
        if path_error:
            errors.append(f"{label}: {path_error}: {value!r}")

    def walk_strings(value: Any) -> Iterable[str]:
        if isinstance(value, str):
            yield value
        elif isinstance(value, dict):
            for child in value.values():
                yield from walk_strings(child)
        elif isinstance(value, list):
            for child in value:
                yield from walk_strings(child)

    for scenario_index, scenario in enumerate(scenarios):
        label = f"install-v1 scenarios.scenarios[{scenario_index}]"
        if not isinstance(scenario, dict):
            continue
        scenario_id = scenario.get("id")
        if isinstance(scenario_id, str):
            if scenario_id in seen_ids:
                errors.append(f"{label}.id: duplicate scenario ID {scenario_id}")
            seen_ids.add(scenario_id)
        categories = scenario.get("categories")
        if isinstance(categories, list):
            covered_categories.update(
                category for category in categories if isinstance(category, str)
            )
        if scenario.get("sourceMode") != "archive-no-git":
            errors.append(f"{label}.sourceMode: install-v1 fixtures must not require .git")

        before = scenario.get("before")
        if isinstance(before, dict):
            for directory_index, directory in enumerate(before.get("directories", [])):
                check_portable(
                    directory,
                    f"{label}.before.directories[{directory_index}]",
                )
                if isinstance(directory, str) and (
                    directory == "Toolkit" or directory.startswith("Toolkit/")
                ):
                    errors.append(
                        f"{label}.before.directories[{directory_index}]: "
                        "fixture input contains Toolkit live state"
                    )
            for file_index, file_declaration in enumerate(before.get("files", [])):
                file_label = f"{label}.before.files[{file_index}]"
                if not isinstance(file_declaration, dict):
                    continue
                path = file_declaration.get("path")
                check_portable(path, f"{file_label}.path")
                if isinstance(path, str) and (
                    path == "Toolkit"
                    or path.startswith("Toolkit/")
                    or any(part.casefold() == ".git" for part in PurePosixPath(path).parts)
                ):
                    errors.append(
                        f"{file_label}.path: fixture input contains Toolkit live/admin state"
                    )
                content_base64 = file_declaration.get("contentBase64")
                if isinstance(content_base64, str):
                    try:
                        base64.b64decode(content_base64, validate=True)
                    except (binascii.Error, ValueError):
                        errors.append(f"{file_label}.contentBase64: invalid base64")

        expected = scenario.get("expected")
        if not isinstance(expected, dict):
            continue
        expected_path_value = expected.get("path")
        check_portable(expected_path_value, f"{label}.expected.path")
        if not isinstance(expected_path_value, str):
            continue
        referenced_expected_paths.add(expected_path_value)
        expected_key = portable_path_key(expected_path_value)
        if expected_key in seen_expected_paths:
            errors.append(
                f"{label}.expected.path: expected outcome paths must be unique"
            )
        seen_expected_paths.add(expected_key)
        expected_path = root / Path(*PurePosixPath(expected_path_value).parts)
        outcome, outcome_errors = _read_json(
            expected_path,
            f"install-v1 {scenario_id} expected outcome",
        )
        errors += outcome_errors
        if outcome is None:
            continue
        kind = expected.get("kind")
        expected_suffix = (
            ".failure.json" if kind == "fatal" else ".plan.json"
        )
        if (
            not expected_path_value.startswith("expected/")
            or not expected_path_value.endswith(expected_suffix)
        ):
            errors.append(
                f"{label}.expected.path: {kind!r} outcomes must use an "
                f"expected/*{expected_suffix} path"
            )
        if kind in {"applicable-plan", "blocked-plan"}:
            if schema is not None:
                errors += validate_install_plan(
                    outcome,
                    schema,
                    f"install-v1 {scenario_id} expected plan",
                    contract,
                )
            if isinstance(outcome, dict):
                expected_can_apply = kind == "applicable-plan"
                if outcome.get("canApply") is not expected_can_apply:
                    errors.append(
                        f"install-v1 {scenario_id} expected kind contradicts canApply"
                    )
        elif kind == "fatal":
            if not isinstance(outcome, dict) or set(outcome) != {
                "kind",
                "failureClass",
            }:
                errors.append(
                    f"install-v1 {scenario_id} fatal outcome must contain only "
                    "kind and failureClass"
                )
            elif (
                outcome.get("kind") != "fatal"
                or outcome.get("failureClass") not in INSTALL_FAILURE_CLASSES
            ):
                errors.append(
                    f"install-v1 {scenario_id} fatal outcome uses an unsupported class"
                )
        for text in walk_strings(outcome):
            if volatile_timestamp.search(text):
                errors.append(
                    f"install-v1 {scenario_id} expected outcome contains a volatile timestamp"
                )
            if text.startswith(("/", "\\")) or WINDOWS_DRIVE_PATTERN.match(text):
                errors.append(
                    f"install-v1 {scenario_id} expected outcome contains an absolute path"
                )

        assertions = scenario.get("assertions")
        if isinstance(assertions, dict):
            after_apply = assertions.get("afterApply")
            if isinstance(after_apply, dict):
                for field in ("files", "utf8Contains"):
                    for assertion_index, assertion in enumerate(
                        after_apply.get(field, [])
                    ):
                        if isinstance(assertion, dict):
                            check_portable(
                                assertion.get("path"),
                                f"{label}.assertions.afterApply.{field}"
                                f"[{assertion_index}].path",
                            )
                for assertion_index, path in enumerate(
                    after_apply.get("absentPaths", [])
                ):
                    check_portable(
                        path,
                        f"{label}.assertions.afterApply.absentPaths"
                        f"[{assertion_index}]",
                    )

    missing_categories = sorted(
        REQUIRED_INSTALL_CONFORMANCE_CATEGORIES - covered_categories
    )
    if missing_categories:
        errors.append(
            "install-v1 scenarios omit promised categories: "
            + ", ".join(missing_categories)
        )
    expected_root = root / "expected"
    actual_expected_paths = {
        path.relative_to(root).as_posix()
        for path in expected_root.rglob("*")
        if path.is_file()
    }
    unreferenced_expected_paths = sorted(
        actual_expected_paths - referenced_expected_paths
    )
    if unreferenced_expected_paths:
        errors.append(
            "install-v1 contains unreferenced expected outcomes: "
            + ", ".join(unreferenced_expected_paths)
        )
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
    if contract.get("orderingPolicy") != CANONICAL_INSTALL_ORDERING_POLICY:
        errors.append(
            "Toolkit/SDP-install.manifest.json orderingPolicy differs from the "
            "canonical install-v1 policy"
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
    excluded: list[tuple[tuple[str, ...], str]] = []
    excluded_paths: dict[tuple[str, ...], str] = {}
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
            excluded_key = portable_path_key(excluded_path)
            if excluded_key in excluded_paths:
                errors.append(
                    f"{label}: duplicate/case-colliding exclusion path {excluded_path}; "
                    f"already used by {excluded_paths[excluded_key]}"
                )
            excluded_paths[excluded_key] = excluded_path
            if isinstance(excluded_kind, str):
                excluded.append((excluded_key, excluded_kind))
            actual = repo / Path(*PurePosixPath(excluded_path).parts)
            if excluded_kind == "file" and not actual.is_file():
                errors.append(f"{label}: excluded file does not exist: {excluded_path}")
            if excluded_kind == "tree" and not actual.is_dir():
                errors.append(f"{label}: excluded tree does not exist: {excluded_path}")

    def is_excluded(path: str) -> bool:
        path_key = portable_path_key(path)
        return any(
            path_key == excluded_key
            or (
                kind == "tree"
                and len(path_key) > len(excluded_key)
                and path_key[: len(excluded_key)] == excluded_key
            )
            for excluded_key, kind in excluded
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
    missing_exclusions = sorted(
        path
        for path in required_exclusions
        if portable_path_key(path) not in excluded_paths
    )
    if missing_exclusions:
        errors.append(
            "Installation contract omits required live/legacy exclusions: "
            + ", ".join(missing_exclusions)
        )

    entries = contract.get("entries")
    if not isinstance(entries, list):
        return errors
    seen_ids: set[str] = set()
    seen_destinations: dict[tuple[str, ...], str] = {}
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
                destination_key = portable_path_key(destination)
                if destination_key in seen_destinations:
                    errors.append(
                        f"{label}: duplicate/case-colliding destination {destination}; "
                        f"already used by {seen_destinations[destination_key]}"
                    )
                for existing_key, existing_destination in seen_destinations.items():
                    destination_contains_existing = (
                        len(destination_key) > len(existing_key)
                        and destination_key[: len(existing_key)] == existing_key
                    )
                    existing_contains_destination = (
                        len(existing_key) > len(destination_key)
                        and existing_key[: len(destination_key)] == destination_key
                    )
                    if destination_contains_existing or existing_contains_destination:
                        errors.append(
                            f"{label}: ancestor/descendant destination conflict {destination}; "
                            f"already used by {existing_destination}"
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
            expected_capability = (
                CANONICAL_GOVERNING_CAPABILITIES_BY_KEY.get(schema_reference.casefold())
                if isinstance(schema_reference, str)
                else None
            )
            if expected_capability is not None and capability != expected_capability:
                errors.append(
                    f"{label}.governing: schema {schema_reference} requires capability "
                    f"{expected_capability}, not {capability!r}"
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
    if isinstance(toolkit, dict):
        toolkit_version = toolkit.get("version")
        toolkit_tag = toolkit.get("gitTag")
        if (
            isinstance(toolkit_version, str)
            and isinstance(toolkit_tag, str)
            and toolkit_tag != f"v{toolkit_version}"
        ):
            errors.append("SDP.manifest.yaml toolkit.gitTag does not match toolkit.version")
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
    install_contract = load_json(repo / "Toolkit/SDP-install.manifest.json")
    errors += validate_install_plan(
        install_plan_example,
        install_plan_schema,
        "examples/install-plan.example.json",
        install_contract,
    )
    if (
        isinstance(install_plan_example, dict)
        and install_plan_example.get("toolkitVersion") != toolkit.get("version")
    ):
        errors.append(
            "examples/install-plan.example.json toolkitVersion differs from SDP.manifest.yaml"
        )
    errors += validate_install_conformance_package(
        repo,
        install_contract,
        install_plan_schema,
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
