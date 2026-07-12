#!/usr/bin/env python3
"""Deterministic validation for the SDP Toolkit release contracts."""

from __future__ import annotations

import argparse
import json
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any, Iterable

import yaml
from jsonschema import Draft202012Validator, FormatChecker

SEMVER_PATTERN = re.compile(
    r"^(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)"
    r"(?:-([0-9A-Za-z.-]+))?(?:\+([0-9A-Za-z.-]+))?$"
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
                if not identifier or (identifier.isdigit() and len(identifier) > 1 and identifier.startswith("0")):
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
    except (OSError, yaml.YAMLError) as exc:
        raise ValueError(f"Cannot parse YAML {path}: {exc}") from exc


def load_json(path: Path) -> Any:
    try:
        return json.loads(path.read_text(encoding="utf-8"))
    except (OSError, json.JSONDecodeError) as exc:
        raise ValueError(f"Cannot parse JSON {path}: {exc}") from exc


def validate_json(instance: Any, schema: Any, label: str) -> list[str]:
    validator = Draft202012Validator(schema, format_checker=FormatChecker())
    return [f"{label}: {error.message}" for error in sorted(validator.iter_errors(instance), key=lambda e: list(e.path))]


def parse_front_matter(path: Path) -> dict[str, Any]:
    text = path.read_text(encoding="utf-8")
    if not text.startswith("---\n"):
        raise ValueError(f"Missing YAML front matter: {path}")
    try:
        _, front, _ = text.split("---\n", 2)
    except ValueError as exc:
        raise ValueError(f"Malformed YAML front matter: {path}") from exc
    data = yaml.safe_load(front)
    if not isinstance(data, dict):
        raise ValueError(f"Front matter is not an object: {path}")
    return data


def release_sections(text: str) -> list[tuple[str, str]]:
    headings = list(re.finditer(r"(?m)^## \[([^\]]+)\](?: - \d{4}-\d{2}-\d{2})?\s*$", text))
    sections: list[tuple[str, str]] = []
    for index, heading in enumerate(headings):
        end = headings[index + 1].start() if index + 1 < len(headings) else len(text)
        sections.append((heading.group(1), text[heading.start():end].rstrip() + "\n"))
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
        except ValueError as exc:
            errors.append(f"{label}: {exc}")
    return errors


def validate_repository(repo: Path, base_ref: str | None = None) -> list[str]:
    errors: list[str] = []
    manifest_path = repo / "SDP.manifest.yaml"
    manifest_schema_path = repo / "Toolkit/schemas/SDP-manifest.schema.json"
    manifest = load_yaml(manifest_path)
    manifest_schema = load_json(manifest_schema_path)
    errors += validate_json(manifest, manifest_schema, "SDP.manifest.yaml")

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
    except (KeyError, ValueError):
        pass

    expected_skills = manifest.get("skills", {})
    actual_skill_dirs = {path.parent.name for path in (repo / "Toolkit/skills").glob("*/SKILL.md")}
    if set(expected_skills) != actual_skill_dirs:
        errors.append(
            "Skill set differs between manifest and Toolkit/skills: "
            f"manifest={sorted(expected_skills)}, files={sorted(actual_skill_dirs)}"
        )
    for skill_id, expected_version in sorted(expected_skills.items()):
        path = repo / "Toolkit/skills" / skill_id / "SKILL.md"
        try:
            metadata = parse_front_matter(path)
        except (OSError, ValueError) as exc:
            errors.append(str(exc))
            continue
        if metadata.get("skillId") != skill_id:
            errors.append(f"{path}: skillId must be {skill_id}")
        if metadata.get("skillVersion") != expected_version:
            errors.append(f"{path}: skillVersion differs from SDP.manifest.yaml")
        if metadata.get("minimumToolkitVersion") != toolkit.get("version"):
            errors.append(f"{path}: minimumToolkitVersion differs from Toolkit version")
        if not isinstance(metadata.get("capabilities"), list) or not metadata["capabilities"]:
            errors.append(f"{path}: capabilities must be a non-empty list")
        if not metadata.get("compatibilityNotes"):
            errors.append(f"{path}: compatibilityNotes is required")

    notes_path = repo / manifest.get("releaseNotesPath", "RELEASE-NOTES.md")
    notes = notes_path.read_text(encoding="utf-8")
    sections = release_sections(notes)
    if not sections or sections[0][0] != "Unreleased":
        errors.append("RELEASE-NOTES.md must begin with ## [Unreleased]")
    elif "Release-Date: unreleased" not in sections[0][1]:
        errors.append("Unreleased release notes must contain Release-Date: unreleased")
    errors += validate_released_immutability(repo, notes, base_ref)

    project_schema = load_json(repo / "Toolkit/schemas/SDP-project-manifest.schema.json")
    project_template = load_yaml(repo / "Toolkit/payload/sdp-root/Framework/templates/SDP-project.manifest.yaml")
    errors += validate_json(project_template, project_schema, "project manifest template")

    installed_schema = load_json(repo / "Toolkit/schemas/installed-toolkit-manifest.schema.json")
    installed_example = load_yaml(repo / "examples/installed-toolkit.manifest.example.yaml")
    errors += validate_json(installed_example, installed_schema, "installed manifest example")

    build_schema = load_json(repo / "Toolkit/schemas/build-identity.schema.json")
    build_example = load_json(repo / "examples/build-identity.example.json")
    errors += validate_json(build_example, build_schema, "build identity example")

    event_schema = load_json(repo / "Toolkit/schemas/release-event.schema.json")
    event_path = repo / "examples/release-events.ndjson.example"
    for line_number, line in enumerate(event_path.read_text(encoding="utf-8").splitlines(), start=1):
        if not line.strip():
            continue
        try:
            event = json.loads(line)
        except json.JSONDecodeError as exc:
            errors.append(f"{event_path}:{line_number}: invalid JSON: {exc}")
            continue
        errors += validate_json(event, event_schema, f"{event_path}:{line_number}")

    for trace_path in (repo / "Traceability/CurrentIndex.yaml", repo / "Traceability/Relations.yaml"):
        data = load_yaml(trace_path)
        if not isinstance(data, dict):
            errors.append(f"{trace_path}: must be a YAML object")

    return errors


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--repo", type=Path, default=Path(__file__).resolve().parents[2])
    parser.add_argument("--base-ref", help="Git ref used to verify immutable released note sections")
    args = parser.parse_args()

    try:
        errors = validate_repository(args.repo.resolve(), args.base_ref)
    except (OSError, ValueError, KeyError, TypeError) as exc:
        errors = [str(exc)]

    if errors:
        print("SDP validation failed:", file=sys.stderr)
        for error in errors:
            print(f"- {error}", file=sys.stderr)
        return 1
    print("SDP validation passed")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
