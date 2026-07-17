#!/usr/bin/env python3
"""Portable materializer and PowerShell-reference harness for install contract v1."""

from __future__ import annotations

import argparse
import base64
import hashlib
import importlib.util
import json
import os
import shutil
import subprocess
import sys
import tempfile
from pathlib import Path, PurePosixPath
from typing import Any, Sequence

from jsonschema import Draft202012Validator


ROOT = Path(__file__).resolve().parents[3]
PACKAGE_ROOT = Path(__file__).resolve().parent
VALIDATOR_PATH = ROOT / "Toolkit/scripts/validate_sdp.py"
SPEC = importlib.util.spec_from_file_location("sdp_validate_for_conformance", VALIDATOR_PATH)
assert SPEC and SPEC.loader
VALIDATE = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = VALIDATE
SPEC.loader.exec_module(VALIDATE)


def load_json(path: Path) -> Any:
    return json.loads(path.read_text(encoding="utf-8"))


def portable_path(root: Path, value: str) -> Path:
    return root / Path(*PurePosixPath(value).parts)


def archive_ignore(_directory: str, names: list[str]) -> set[str]:
    ignored = {name for name in names if name in {".git", "__pycache__"}}
    ignored.update(name for name in names if name.endswith((".pyc", ".pyo")))
    return ignored


def create_archive_source(destination: Path) -> None:
    shutil.copytree(ROOT, destination, ignore=archive_ignore)
    if any(path.name == ".git" for path in destination.rglob("*")):
        raise AssertionError("archive source unexpectedly contains .git")


def installed_manifest_bytes(
    contract: dict[str, Any],
    overrides: dict[str, Any] | None = None,
) -> bytes:
    generator = next(
        item
        for item in contract["generators"]
        if item["id"] == "installed-toolkit-manifest"
    )
    facts = dict(generator["facts"])
    if overrides:
        facts.update(overrides)
    lines = [
        f"schemaVersion: {json.dumps(facts['schemaVersion'])}",
        f"toolkitVersion: {json.dumps(facts['toolkitVersion'])}",
        f"frameworkVersion: {json.dumps(facts['frameworkVersion'])}",
        f"agentsContractVersion: {json.dumps(facts['agentsContractVersion'])}",
        f"installerVersion: {json.dumps(facts['installerVersion'])}",
        'toolkitInstalledAt: "2026-01-01T00:00:00Z"',
        "sourceCommit: null",
        "skills:",
    ]
    lines.extend(
        f"  {skill_id}: {json.dumps(version)}"
        for skill_id, version in facts["skills"].items()
    )
    lines.append("capabilities:")
    lines.extend(f"  - {capability}" for capability in facts["capabilities"])
    return ("\n".join(lines) + "\n").encode("utf-8")


def materialize_profile(
    project: Path,
    source: Path,
    contract: dict[str, Any],
    before: dict[str, Any],
) -> None:
    project.mkdir(parents=True, exist_ok=True)
    profile = before["profile"]
    initialize = profile == "installed-initialize"
    if profile != "empty":
        for entry in contract["entries"]:
            if entry["selectionPolicy"] == "initialize-only" and not initialize:
                continue
            destination = portable_path(project, entry["destination"])
            destination.parent.mkdir(parents=True, exist_ok=True)
            if entry["kind"] == "copied":
                destination.write_bytes(portable_path(source, entry["source"]).read_bytes())
            elif entry["generator"] == "empty-ledger":
                destination.write_bytes(b"")
            elif entry["generator"] == "installed-toolkit-manifest":
                destination.write_bytes(
                    installed_manifest_bytes(
                        contract,
                        before.get("installedManifestFacts"),
                    )
                )
            else:  # pragma: no cover - contract validation closes this vocabulary.
                raise AssertionError(f"unsupported generator {entry['generator']}")

    for relative in before["directories"]:
        portable_path(project, relative).mkdir(parents=True, exist_ok=True)
    for declaration in before["files"]:
        path = portable_path(project, declaration["path"])
        path.parent.mkdir(parents=True, exist_ok=True)
        if "contentBase64" in declaration:
            content = base64.b64decode(declaration["contentBase64"], validate=True)
        else:
            content = declaration["contentUtf8"].encode("utf-8")
        path.write_bytes(content)


def tree_fingerprint(root: Path) -> str:
    rows: list[bytes] = []
    for path in sorted(root.rglob("*"), key=lambda item: item.relative_to(root).as_posix()):
        relative = path.relative_to(root).as_posix().encode("utf-8")
        if path.is_dir():
            rows.append(b"D:" + relative)
        else:
            rows.append(
                b"F:"
                + relative
                + b"="
                + hashlib.sha256(path.read_bytes()).hexdigest().encode("ascii")
            )
    return hashlib.sha256(b"\n".join(rows)).hexdigest()


def powershell_base_command(executable: str) -> list[str]:
    return [
        executable,
        "-NoProfile",
        "-NonInteractive",
        "-ExecutionPolicy",
        "Bypass",
    ]


def option_arguments(scenario: dict[str, Any]) -> list[str]:
    arguments: list[str] = []
    options = scenario["options"]
    if options["initializeProjectStructure"]:
        arguments.append("-InitializeProjectStructure")
    if options["forceManagedFiles"]:
        arguments.append("-ForceManagedFiles")
    return arguments


def invoke_plan(
    executable: str,
    installer: Path,
    project: Path,
    scenario: dict[str, Any],
) -> dict[str, Any]:
    command = powershell_base_command(executable) + [
        "-File",
        str(installer),
        "-ProjectRoot",
        str(project),
        "-PlanJson",
        *option_arguments(scenario),
    ]
    result = subprocess.run(command, capture_output=True, text=True)
    if result.returncode != 0:
        raise AssertionError(
            f"{scenario['id']}: PlanJson failed ({result.returncode})\n"
            f"stdout:\n{result.stdout}\nstderr:\n{result.stderr}"
        )
    return json.loads(result.stdout)


def invoke_fatal(
    executable: str,
    installer: Path,
    project: Path,
    scenario: dict[str, Any],
) -> dict[str, Any]:
    wrapper = r"""
$ErrorActionPreference = 'Stop'
$parameters = @{
    ProjectRoot = $env:SDP_CONFORMANCE_PROJECT
    PlanJson = $true
}
if ($env:SDP_CONFORMANCE_INITIALIZE -ceq 'true') {
    $parameters.InitializeProjectStructure = $true
}
if ($env:SDP_CONFORMANCE_FORCE -ceq 'true') {
    $parameters.ForceManagedFiles = $true
}
try {
    & $env:SDP_CONFORMANCE_INSTALLER @parameters | Out-Null
    [ordered]@{ kind = 'unexpected-success'; failureClass = $null } |
        ConvertTo-Json -Compress
} catch {
    $failureClass = $null
    $current = $_.Exception
    while (($null -ne $current) -and [string]::IsNullOrWhiteSpace($failureClass)) {
        $failureClass = [string]$current.Data['sdpFailureClass']
        $current = $current.InnerException
    }
    [ordered]@{ kind = 'fatal'; failureClass = $failureClass } |
        ConvertTo-Json -Compress
}
"""
    environment = os.environ.copy()
    environment.update(
        {
            "SDP_CONFORMANCE_INSTALLER": str(installer),
            "SDP_CONFORMANCE_PROJECT": str(project),
            "SDP_CONFORMANCE_INITIALIZE": str(
                scenario["options"]["initializeProjectStructure"]
            ).lower(),
            "SDP_CONFORMANCE_FORCE": str(
                scenario["options"]["forceManagedFiles"]
            ).lower(),
        }
    )
    result = subprocess.run(
        powershell_base_command(executable) + ["-Command", wrapper],
        capture_output=True,
        text=True,
        env=environment,
    )
    if result.returncode != 0:
        raise AssertionError(
            f"{scenario['id']}: fatal wrapper failed ({result.returncode})\n"
            f"stdout:\n{result.stdout}\nstderr:\n{result.stderr}"
        )
    return json.loads(result.stdout)


def invoke_apply(
    executable: str,
    installer: Path,
    project: Path,
    scenario: dict[str, Any],
) -> None:
    command = powershell_base_command(executable) + [
        "-File",
        str(installer),
        "-ProjectRoot",
        str(project),
        *option_arguments(scenario),
    ]
    result = subprocess.run(command, capture_output=True, text=True)
    if result.returncode != 0:
        raise AssertionError(
            f"{scenario['id']}: apply failed ({result.returncode})\n"
            f"stdout:\n{result.stdout}\nstderr:\n{result.stderr}"
        )


def apply_source_mutations(source: Path, scenario: dict[str, Any]) -> None:
    mutations = scenario.get("sourceManifestMutations", [])
    if not mutations:
        return
    manifest_path = source / "Toolkit/SDP-install.manifest.json"
    manifest = load_json(manifest_path)
    for mutation in mutations:
        if mutation != {
            "op": "replace",
            "path": "/orderingPolicy",
            "value": mutation["value"],
        }:
            raise AssertionError(f"unsupported manifest mutation: {mutation}")
        manifest["orderingPolicy"] = mutation["value"]
    manifest_path.write_text(json.dumps(manifest, indent=2) + "\n", encoding="utf-8")


def verify_after_apply(project: Path, scenario: dict[str, Any]) -> None:
    after_apply = scenario["assertions"].get("afterApply", {})
    for assertion in after_apply.get("files", []):
        path = portable_path(project, assertion["path"])
        if not path.is_file():
            raise AssertionError(f"{scenario['id']}: expected file is missing: {assertion['path']}")
        actual = hashlib.sha256(path.read_bytes()).hexdigest()
        if actual != assertion["sha256"]:
            raise AssertionError(
                f"{scenario['id']}: SHA-256 differs for {assertion['path']}: {actual}"
            )
    for relative in after_apply.get("absentPaths", []):
        if portable_path(project, relative).exists():
            raise AssertionError(f"{scenario['id']}: forbidden path exists: {relative}")
    for assertion in after_apply.get("utf8Contains", []):
        text = portable_path(project, assertion["path"]).read_text(encoding="utf-8")
        if assertion["value"] not in text:
            raise AssertionError(
                f"{scenario['id']}: {assertion['path']} omits {assertion['value']!r}"
            )


def select_powershell(explicit: str | None) -> str | None:
    if explicit:
        resolved = shutil.which(explicit)
        if resolved:
            return resolved
        path = Path(explicit)
        return str(path.resolve()) if path.is_file() else None
    for candidate in ("pwsh", "powershell.exe", "powershell"):
        resolved = shutil.which(candidate)
        if resolved:
            return resolved
    return None


def validate_index_shape(index: Any, schema: Any) -> list[str]:
    validator = Draft202012Validator(schema)
    return [
        error.message
        for error in sorted(
            validator.iter_errors(index),
            key=lambda error: tuple(str(part) for part in error.absolute_path),
        )
    ]


def run(argv: Sequence[str] | None = None) -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--validate-only", action="store_true")
    parser.add_argument("--powershell", help="PowerShell executable or absolute path")
    parser.add_argument("--scenario", action="append", dest="scenario_ids")
    parser.add_argument(
        "--write-candidates",
        action="store_true",
        help="Explicit maintainer-only update of expected outcomes; never used by CI",
    )
    args = parser.parse_args(argv)
    if args.validate_only and args.write_candidates:
        parser.error("--validate-only cannot be combined with --write-candidates")
    if args.write_candidates and os.environ.get("CI"):
        parser.error("--write-candidates is maintainer-only and cannot run in CI")

    index = load_json(PACKAGE_ROOT / "scenarios.json")
    index_schema = load_json(PACKAGE_ROOT / "scenario-index.schema.json")
    shape_errors = validate_index_shape(index, index_schema)
    if shape_errors:
        for error in shape_errors:
            print(f"install-v1 scenario index: {error}", file=sys.stderr)
        return 1

    if not args.write_candidates:
        errors = VALIDATE.validate_install_conformance_package(ROOT)
        if errors:
            for error in errors:
                print(error, file=sys.stderr)
            return 1
    if args.validate_only:
        print(f"install-v1 package validation passed ({len(index['scenarios'])} scenarios)")
        return 0

    executable = select_powershell(args.powershell)
    if executable is None:
        print("No PowerShell executable found; use --powershell PATH", file=sys.stderr)
        return 2
    selected_ids = set(args.scenario_ids or [])
    scenarios = [
        scenario
        for scenario in index["scenarios"]
        if not selected_ids or scenario["id"] in selected_ids
    ]
    missing_ids = selected_ids - {scenario["id"] for scenario in scenarios}
    if missing_ids:
        print("Unknown scenario IDs: " + ", ".join(sorted(missing_ids)), file=sys.stderr)
        return 2

    contract = load_json(ROOT / "Toolkit/SDP-install.manifest.json")
    plan_schema = load_json(ROOT / "Toolkit/schemas/SDP-install-plan.schema.json")
    with tempfile.TemporaryDirectory(prefix="sdp-install-v1-") as temporary:
        temporary_root = Path(temporary)
        shared_source = temporary_root / "source-archive"
        create_archive_source(shared_source)
        for scenario in scenarios:
            scenario_source = shared_source
            if scenario.get("sourceManifestMutations"):
                scenario_source = temporary_root / f"source-{scenario['id']}"
                shutil.copytree(shared_source, scenario_source)
                apply_source_mutations(scenario_source, scenario)
            project = temporary_root / f"project-{scenario['id']}"
            materialize_profile(project, scenario_source, contract, scenario["before"])
            before = tree_fingerprint(project)
            expected_path = portable_path(PACKAGE_ROOT, scenario["expected"]["path"])
            if scenario["expected"]["kind"] == "fatal":
                actual = invoke_fatal(
                    executable,
                    scenario_source / "Toolkit/scripts/Install-SDP.ps1",
                    project,
                    scenario,
                )
                if (
                    not isinstance(actual, dict)
                    or set(actual) != {"kind", "failureClass"}
                    or actual.get("kind") != "fatal"
                    or actual.get("failureClass") not in VALIDATE.INSTALL_FAILURE_CLASSES
                ):
                    raise AssertionError(
                        f"{scenario['id']}: reference emitted an invalid fatal outcome "
                        f"{actual!r}"
                    )
            else:
                actual = invoke_plan(
                    executable,
                    scenario_source / "Toolkit/scripts/Install-SDP.ps1",
                    project,
                    scenario,
                )
                plan_errors = VALIDATE.validate_install_plan(
                    actual,
                    plan_schema,
                    f"{scenario['id']} reference plan",
                    contract,
                )
                if plan_errors:
                    raise AssertionError("\n".join(plan_errors))
                expected_can_apply = (
                    scenario["expected"]["kind"] == "applicable-plan"
                )
                if actual.get("canApply") is not expected_can_apply:
                    raise AssertionError(
                        f"{scenario['id']}: expected kind contradicts reference canApply"
                    )
            after_plan = tree_fingerprint(project)
            if before != after_plan:
                raise AssertionError(f"{scenario['id']}: planning mutated the project tree")

            if args.write_candidates:
                expected_path.parent.mkdir(parents=True, exist_ok=True)
                expected_path.write_text(
                    json.dumps(actual, indent=2) + "\n",
                    encoding="utf-8",
                )
            else:
                expected = load_json(expected_path)
                if actual != expected:
                    raise AssertionError(
                        f"{scenario['id']}: reference outcome differs from committed authority"
                    )

            if scenario["expected"]["kind"] == "applicable-plan":
                invoke_apply(
                    executable,
                    scenario_source / "Toolkit/scripts/Install-SDP.ps1",
                    project,
                    scenario,
                )
                verify_after_apply(project, scenario)
            print(f"PASS {scenario['id']}")

    print(
        f"install-v1 PowerShell reference conformance passed "
        f"({len(scenarios)} scenarios; {Path(executable).name})"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(run())
