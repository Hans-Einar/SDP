#!/usr/bin/env python3
"""Generate portable SDP build identity JSON from project and Git state."""

from __future__ import annotations

import argparse
import json
import subprocess
from datetime import datetime, timezone
from pathlib import Path
from typing import Any

import yaml


def git(project_root: Path, *args: str) -> str | None:
    try:
        result = subprocess.run(
            ["git", "-C", str(project_root), *args],
            check=True,
            capture_output=True,
            text=True,
        )
    except (FileNotFoundError, subprocess.CalledProcessError):
        return None
    return result.stdout.strip()


def load_manifest(project_root: Path) -> dict[str, Any]:
    path = project_root / "SDP" / "SDP-project.manifest.yaml"
    if not path.is_file():
        raise FileNotFoundError(f"Project manifest not found: {path}")
    data = yaml.safe_load(path.read_text(encoding="utf-8"))
    if not isinstance(data, dict):
        raise ValueError(f"Project manifest is not an object: {path}")
    return data


def generate_identity(
    project_root: Path,
    timestamp: str | None = None,
) -> dict[str, Any]:
    manifest = load_manifest(project_root)
    release = manifest.get("release") or {}
    development = manifest.get("development") or {}

    sprint_id = development.get("sprintId")
    refactor_id = development.get("refactorId")
    if sprint_id and refactor_id:
        raise ValueError("sprintId and refactorId cannot both be active")

    current_version = release.get("currentVersion")
    target_version = release.get("nextTargetVersion") or current_version
    if not target_version:
        raise ValueError("release.currentVersion or nextTargetVersion is required")

    release_state = release.get("state", "unreleased")
    allowed_states = {"unreleased", "prerelease", "released", "yanked"}
    if release_state not in allowed_states:
        raise ValueError(f"Unsupported release state: {release_state}")
    full_commit = git(project_root, "rev-parse", "HEAD")
    short_commit = git(project_root, "rev-parse", "--short=7", "HEAD")
    status = git(project_root, "status", "--porcelain")
    dirty = bool(status) if status is not None else None

    if timestamp is None:
        timestamp = datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")

    coordinate_parts = [
        value
        for value in (
            sprint_id or refactor_id,
            development.get("iterationId"),
            development.get("sliceId") or development.get("fixId"),
        )
        if value
    ]
    revision = development.get("revision")

    released = release_state == "released"
    if release_state == "released":
        version_label = str(target_version)
    elif release_state == "yanked":
        version_label = f"{target_version}-yanked"
    elif release_state == "prerelease":
        version_label = f"{target_version}-prerelease"
    else:
        version_label = f"{target_version}-dev"
    display_parts = [version_label]
    if coordinate_parts:
        display_parts.append(" / ".join(coordinate_parts))
    if revision is not None:
        display_parts.append(f"r{revision}")
    if short_commit:
        display_parts.append(short_commit)

    machine_parts = [f"v{target_version}"]
    if release_state == "unreleased":
        machine_parts.append("dev")
    elif release_state != "released":
        machine_parts.append(release_state)
    if sprint_id:
        machine_parts.append("s" + "".join(ch for ch in str(sprint_id) if ch.isdigit()).zfill(3))
    elif refactor_id:
        ref = "".join(ch for ch in str(refactor_id) if ch.isalnum()).lower()
        machine_parts.append(f"rf.{ref}")
    if development.get("iterationId"):
        iteration = "".join(ch for ch in str(development["iterationId"]) if ch.isdigit()).zfill(3)
        machine_parts.append(f"i{iteration}")
    if development.get("sliceId"):
        slice_number = "".join(ch for ch in str(development["sliceId"]) if ch.isdigit()).zfill(3)
        machine_parts.append(f"sl{slice_number}")
    elif development.get("fixId"):
        machine_parts.append(f"fix.{str(development['fixId']).lower()}")
    if revision is not None:
        machine_parts.append(f"r{int(revision):03d}")
    development_id = ".".join(machine_parts)
    if short_commit:
        development_id += f"+sha.{short_commit}"

    return {
        "schemaVersion": "1.0",
        "releaseVersion": str(target_version),
        "releaseState": release_state,
        "sprintId": sprint_id,
        "refactorId": refactor_id,
        "iterationId": development.get("iterationId"),
        "sliceId": development.get("sliceId"),
        "fixId": development.get("fixId"),
        "revision": revision,
        "commit": full_commit,
        "shortCommit": short_commit,
        "buildTimestamp": timestamp,
        "dirty": dirty,
        "developmentId": development_id,
        "displayVersion": " · ".join(display_parts),
    }


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--project-root", type=Path, default=Path.cwd())
    parser.add_argument("--output", type=Path)
    parser.add_argument("--timestamp", help="Deterministic UTC timestamp for testing")
    args = parser.parse_args()

    try:
        identity = generate_identity(args.project_root.resolve(), args.timestamp)
    except (FileNotFoundError, ValueError, yaml.YAMLError) as exc:
        parser.error(str(exc))

    output = json.dumps(identity, indent=2, sort_keys=True) + "\n"
    if args.output:
        args.output.parent.mkdir(parents=True, exist_ok=True)
        args.output.write_text(output, encoding="utf-8")
    else:
        print(output, end="")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
