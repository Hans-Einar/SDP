from __future__ import annotations

import copy
import hashlib
import importlib.util
import json
import shutil
import tempfile
import unittest
import sys
from contextlib import redirect_stderr, redirect_stdout
from io import StringIO
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
SPEC = importlib.util.spec_from_file_location("validate_sdp", ROOT / "Toolkit/scripts/validate_sdp.py")
assert SPEC and SPEC.loader
VALIDATE = importlib.util.module_from_spec(SPEC)
sys.modules[SPEC.name] = VALIDATE
SPEC.loader.exec_module(VALIDATE)

BUILD_SPEC = importlib.util.spec_from_file_location("build_identity", ROOT / "Toolkit/scripts/New-SDPBuildIdentity.py")
assert BUILD_SPEC and BUILD_SPEC.loader
BUILD = importlib.util.module_from_spec(BUILD_SPEC)
sys.modules[BUILD_SPEC.name] = BUILD
BUILD_SPEC.loader.exec_module(BUILD)


def write_yaml(path: Path, value: object) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(yaml.safe_dump(value, sort_keys=False), encoding="utf-8")


def read_yaml(path: Path) -> dict[str, object]:
    value = yaml.safe_load(path.read_text(encoding="utf-8"))
    assert isinstance(value, dict)
    return value


def git_tree_sha(path: Path) -> str:
    """Hash an exported all-100644 text tree like Git, independent of checkout EOLs."""

    def object_digest(kind: str, content: bytes) -> bytes:
        header = f"{kind} {len(content)}\0".encode("ascii")
        return hashlib.sha1(header + content, usedforsecurity=False).digest()

    def tree_digest(directory: Path) -> bytes:
        entries: list[tuple[str, bytes]] = []
        for child in directory.iterdir():
            name = child.name.encode("utf-8")
            if child.is_dir():
                mode = b"40000"
                digest = tree_digest(child)
                sort_key = child.name + "/"
            else:
                mode = b"100644"
                content = child.read_bytes().replace(b"\r\n", b"\n")
                digest = object_digest("blob", content)
                sort_key = child.name
            entries.append((sort_key, mode + b" " + name + b"\0" + digest))
        payload = b"".join(entry for _, entry in sorted(entries))
        return object_digest("tree", payload)

    return tree_digest(path).hex()


def create_valid_consuming_project(root: Path) -> Path:
    sdp = root / "SDP"
    project_manifest = {
        "schemaVersion": "1.0",
        "project": {"name": "Validator fixture", "capabilities": []},
        "installed": {"manifestPath": "Framework/installed-toolkit.manifest.yaml"},
        "release": {
            "currentVersion": "1.0.0",
            "nextTargetVersion": "1.1.0",
            "state": "unreleased",
            "latestTag": "v1.0.0",
            "latestCommit": "abcdef1234567890",
        },
        "development": {
            "sprintId": "Sprint-001",
            "refactorId": None,
            "iterationId": "SPI-001",
            "sliceId": "SPS-001",
            "fixId": None,
            "revision": None,
        },
        "migration": {"pendingWarnings": []},
    }
    write_yaml(sdp / "SDP-project.manifest.yaml", project_manifest)

    installed_manifest = {
        "schemaVersion": "1.0",
        "toolkitVersion": "0.2.0",
        "frameworkVersion": "1.0.0",
        "agentsContractVersion": "1.0.0",
        "installerVersion": "0.2.0",
        "toolkitInstalledAt": "2026-07-13T20:00:00Z",
        "sourceCommit": None,
        "skills": {"sdp-test": "1.0.0"},
        "capabilities": ["sdp.install.v1", "sdp.manifest.v1"],
    }
    write_yaml(sdp / "Framework/installed-toolkit.manifest.yaml", installed_manifest)
    skill_path = root / ".codex/skills/sdp-test/SKILL.md"
    skill_path.parent.mkdir(parents=True, exist_ok=True)
    skill_path.write_text(
        """---
skillId: sdp-test
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
  - sdp.test
compatibilityNotes: Deterministic validator fixture.
---

# Test skill
""",
        encoding="utf-8",
    )

    (sdp / "Sprints/Sprint-001").mkdir(parents=True)
    (sdp / "Sprints/Sprint-001/ScrumIterations.md").write_text(
        "# Sprint-001\n", encoding="utf-8"
    )
    (sdp / "CodeReview").mkdir(parents=True)
    (sdp / "CodeReview/REV-001.md").write_text("# Review\n", encoding="utf-8")
    (sdp / "Verification").mkdir(parents=True)
    (sdp / "Verification/VER-001.md").write_text("# Verification\n", encoding="utf-8")
    (sdp / "RELEASE-NOTES.md").write_text(
        """# Release Notes

## [Unreleased]

Release-Date: unreleased

### Added

- [SPS-001] Validator fixture.
""",
        encoding="utf-8",
    )

    release_record = {
        "schemaVersion": "1.0",
        "releaseId": "REL-1.1.0",
        "version": "1.1.0",
        "state": "unreleased",
        "releaseNotesPath": "SDP/RELEASE-NOTES.md",
        "includedCapabilities": [],
        "verificationRecords": ["VER-001"],
        "reviewRecords": ["REV-001"],
        "migrationRecord": None,
        "releasePreparationCommit": None,
        "gitTag": None,
        "githubReleaseUrl": None,
    }
    write_yaml(sdp / "Releases/REL-1.1.0.yaml", release_record)
    fix_record = {
        "schemaVersion": "1.0",
        "fixId": "FIX-1.1.0-001",
        "targetVersion": "1.1.0",
        "status": "planned",
        "reasonFullSliceIsDisproportionate": "Fixture coverage",
        "scope": [],
        "expectedFiles": [],
        "invariants": [],
        "nonGoals": [],
        "verification": [],
        "reviewRecord": "REV-001",
        "revision": None,
    }
    write_yaml(sdp / "Fixes/FIX-1.1.0-001.yaml", fix_record)

    current = {
        "project": {"name": "Validator fixture", "status": "active-development"},
        "release": {
            "previousVersion": "1.0.0",
            "activeReleaseId": "REL-1.1.0",
            "targetVersion": "1.1.0",
            "state": "unreleased",
        },
        "active": {
            "sprint": "Sprint-001",
            "refactor": None,
            "iteration": "SPI-001",
            "slice": "SPS-001",
            "fix": None,
            "revision": None,
        },
    }
    write_yaml(sdp / "Traceability/CurrentIndex.yaml", current)
    relations = {
        "requirements": {},
        "designs": {},
        "sprints": {
            "Sprint-001": {
                "status": "active",
                "path": "Sprints/Sprint-001/ScrumIterations.md",
                "release": "REL-1.1.0",
                "iterations": ["SPI-001"],
            }
        },
        "refactors": {},
        "iterations": {
            "SPI-001": {
                "status": "active",
                "sprint": "Sprint-001",
                "refactor": None,
                "path": "Sprints/Sprint-001/ScrumIterations.md",
                "slices": ["SPS-001"],
            }
        },
        "slices": {
            "SPS-001": {
                "status": "active",
                "sprint": "Sprint-001",
                "iteration": "SPI-001",
                "release": "REL-1.1.0",
                "path": "Sprints/Sprint-001/ScrumIterations.md",
                "review": "REV-001",
                "verification": "VER-001",
            }
        },
        "fixes": {
            "FIX-1.1.0-001": {
                "status": "planned",
                "path": "Fixes/FIX-1.1.0-001.yaml",
                "release": "REL-1.1.0",
                "review": "REV-001",
            }
        },
        "reviews": {
            "REV-001": {
                "path": "CodeReview/REV-001.md",
                "slice": "SPS-001",
                "fix": "FIX-1.1.0-001",
            }
        },
        "verification": {
            "VER-001": {"path": "Verification/VER-001.md", "slice": "SPS-001"}
        },
        "migrations": {},
        "releases": {
            "REL-1.1.0": {
                "version": "1.1.0",
                "state": "unreleased",
                "releaseNotes": "RELEASE-NOTES.md",
                "manifest": "SDP-project.manifest.yaml",
                "releaseRecord": "Releases/REL-1.1.0.yaml",
                "sprints": ["Sprint-001"],
                "slices": ["SPS-001"],
                "verification": ["VER-001"],
                "reviews": ["REV-001"],
                "releaseCommit": None,
                "gitTag": None,
                "githubRelease": None,
            }
        },
    }
    write_yaml(sdp / "Traceability/Relations.yaml", relations)
    (sdp / "Traceability/Ledger.ndjson").write_text("", encoding="utf-8")
    return root


class SemVerTests(unittest.TestCase):
    def test_ordering(self) -> None:
        self.assertLess(VALIDATE.SemVer.parse("0.1.9"), VALIDATE.SemVer.parse("0.2.0"))
        self.assertLess(VALIDATE.SemVer.parse("1.0.0-rc.1"), VALIDATE.SemVer.parse("1.0.0"))
        self.assertLess(VALIDATE.SemVer.parse("1.0.0-alpha.1"), VALIDATE.SemVer.parse("1.0.0-alpha.beta"))

    def test_invalid_semver(self) -> None:
        for value in (
            "1",
            "v1.2.3",
            "01.2.3",
            "1.2.3-01",
            "1.2",
            "1.2.3-alpha..1",
            "1.2.3-alpha.",
            "1.2.3-.alpha",
            "1.2.3+build..1",
            "1.2.3+build.",
        ):
            with self.subTest(value=value), self.assertRaises(ValueError):
                VALIDATE.SemVer.parse(value)

    def test_every_public_semver_definition_uses_the_strict_rule(self) -> None:
        invalid_values = (
            "1.2.3-01",
            "1.2.3-alpha..1",
            "1.2.3-alpha.",
            "1.2.3+build..1",
            "1.2.3+build.",
        )
        for path in sorted((ROOT / "Toolkit/schemas").glob("*.json")):
            schema = json.loads(path.read_text(encoding="utf-8"))
            semver = schema.get("$defs", {}).get("semver")
            if semver is None:
                continue
            self.assertEqual(semver.get("pattern"), VALIDATE.SEMVER_PATTERN_TEXT, path.name)
            for value in invalid_values:
                with self.subTest(schema=path.name, value=value):
                    self.assertTrue(VALIDATE.validate_json(value, semver, path.name))

    def test_release_id_and_tag_boundaries_reject_invalid_identifiers(self) -> None:
        schemas = {
            path.name: json.loads(path.read_text(encoding="utf-8"))
            for path in (ROOT / "Toolkit/schemas").glob("*.json")
        }
        release_id = "REL-1.2.3-01"
        release_id_cases = (
            (
                schemas["current-index.schema.json"]["$defs"]["releaseId"],
                release_id,
            ),
            (
                schemas["release-record.schema.json"]["properties"]["releaseId"],
                release_id,
            ),
            (
                schemas["release-event.schema.json"]["properties"]["releaseId"],
                release_id,
            ),
            (
                schemas["ledger-event.schema.json"]["properties"]["releaseId"],
                release_id,
            ),
            (
                schemas["relations.schema.json"]["$defs"]["releases"][
                    "propertyNames"
                ],
                release_id,
            ),
        )
        for boundary, value in release_id_cases:
            with self.subTest(value=value):
                self.assertTrue(VALIDATE.validate_json(value, boundary, "release ID"))

        invalid_tag = "v1.2.3+build..1"
        tag_cases = (
            schemas["SDP-manifest.schema.json"]["$defs"]["versionTag"],
            schemas["SDP-project-manifest.schema.json"]["$defs"]["versionTag"],
            schemas["release-record.schema.json"]["properties"]["gitTag"],
            schemas["relations.schema.json"]["$defs"]["relation"]["properties"][
                "gitTag"
            ],
            schemas["release-event.schema.json"]["allOf"][0]["then"]["properties"][
                "payload"
            ]["properties"]["tag"],
        )
        for boundary in tag_cases:
            self.assertTrue(
                VALIDATE.validate_json(invalid_tag, boundary, "release tag")
            )

    def test_release_and_fix_ids_preserve_full_semver_identity(self) -> None:
        schemas = {
            path.name: json.loads(path.read_text(encoding="utf-8"))
            for path in (ROOT / "Toolkit/schemas").glob("*.json")
        }
        release_id = "REL-1.2.3-alpha.1+build.5"
        release_id_schemas = (
            schemas["release-record.schema.json"]["properties"]["releaseId"],
            schemas["release-event.schema.json"]["properties"]["releaseId"],
            schemas["ledger-event.schema.json"]["properties"]["releaseId"],
            schemas["current-index.schema.json"]["$defs"]["releaseId"],
            schemas["relations.schema.json"]["$defs"]["releases"]["propertyNames"],
        )
        for index, schema in enumerate(release_id_schemas):
            with self.subTest(kind="release", index=index):
                self.assertEqual(
                    VALIDATE.validate_json(release_id, schema, "release ID"), []
                )

        fix_id = "FIX-1.2.3-alpha.1+build.5-001"
        fix_id_schemas = (
            schemas["fix-record.schema.json"]["properties"]["fixId"],
            schemas["current-index.schema.json"]["$defs"]["fixId"],
            schemas["relations.schema.json"]["$defs"]["fixes"]["propertyNames"],
        )
        for index, schema in enumerate(fix_id_schemas):
            with self.subTest(kind="fix", index=index):
                self.assertEqual(VALIDATE.validate_json(fix_id, schema, "fix ID"), [])


class ReleaseNotesTests(unittest.TestCase):
    def test_extract_sections(self) -> None:
        text = "# Notes\n\n## [Unreleased]\n\nA\n\n## [1.0.0] - 2026-01-01\n\nB\n"
        self.assertEqual([name for name, _ in VALIDATE.release_sections(text)], ["Unreleased", "1.0.0"])


class FrontMatterTests(unittest.TestCase):
    def test_all_canonical_skills_have_metadata(self) -> None:
        for path in ROOT.glob("Toolkit/skills/*/SKILL.md"):
            metadata = VALIDATE.parse_front_matter(path)
            self.assertEqual(metadata["skillId"], path.parent.name)
            VALIDATE.SemVer.parse(metadata["skillVersion"])
            self.assertTrue(metadata["capabilities"])


class InstallationContractTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls) -> None:
        cls.toolkit_manifest = yaml.safe_load(
            (ROOT / "SDP.manifest.yaml").read_text(encoding="utf-8")
        )
        cls.contract = json.loads(
            (ROOT / "Toolkit/SDP-install.manifest.json").read_text(encoding="utf-8")
        )
        cls.contract_schema = json.loads(
            (ROOT / "Toolkit/schemas/SDP-install-manifest.schema.json").read_text(
                encoding="utf-8"
            )
        )
        cls.plan_schema = json.loads(
            (ROOT / "Toolkit/schemas/SDP-install-plan.schema.json").read_text(
                encoding="utf-8"
            )
        )

    def validate_contract(self, contract: object) -> list[str]:
        return VALIDATE.validate_installation_contract(
            ROOT,
            self.toolkit_manifest,
            contract_data=contract,
            check_installer_integration=False,
        )

    def valid_plan(self) -> dict[str, object]:
        entry = next(
            item for item in self.contract["entries"] if item["id"] == "project-manifest"
        )
        return {
            "schemaVersion": "1.0",
            "manifestSchemaVersion": self.contract["schemaVersion"],
            "orderingPolicy": self.contract["orderingPolicy"],
            "mode": "plan",
            "toolkitVersion": self.contract["toolkitVersion"],
            "installedToolkitVersion": None,
            "options": {
                "initializeProjectStructure": False,
                "forceManagedFiles": False,
            },
            "canApply": True,
            "actions": [
                {
                    "sequence": 1,
                    "action": "create",
                    "entryId": entry["id"],
                    "source": entry["source"],
                    "generator": None,
                    "targetSource": None,
                    "targetSourceSha256": None,
                    "destinationPrecondition": None,
                    "destination": entry["destination"],
                    "ownership": entry["ownership"],
                    "reason": "missing-target",
                    "mutatesTarget": True,
                    "oldToolkitVersion": None,
                    "newToolkitVersion": self.contract["toolkitVersion"],
                }
            ],
        }

    def plan_action(
        self,
        entry_id: str,
        *,
        sequence: int,
        action: str = "create",
        reason: str = "missing-target",
        mutates: bool = True,
    ) -> dict[str, object]:
        entry = next(
            item for item in self.contract["entries"] if item["id"] == entry_id
        )
        return {
            "sequence": sequence,
            "action": action,
            "entryId": entry_id,
            "source": entry.get("source"),
            "generator": entry.get("generator"),
            "targetSource": None,
            "targetSourceSha256": None,
            "destinationPrecondition": None,
            "destination": entry["destination"],
            "ownership": entry["ownership"],
            "reason": reason,
            "mutatesTarget": mutates,
            "oldToolkitVersion": None,
            "newToolkitVersion": self.contract["toolkitVersion"],
        }

    def test_canonical_manifest_schema_and_semantics(self) -> None:
        self.assertEqual(
            VALIDATE.validate_json(
                self.contract, self.contract_schema, "SDP install contract"
            ),
            [],
        )
        self.assertEqual(self.validate_contract(self.contract), [])
        entries = self.contract["entries"]
        self.assertEqual(len(entries), 40)
        self.assertEqual(len({entry["id"] for entry in entries}), len(entries))
        self.assertEqual(
            len({entry["destination"].casefold() for entry in entries}), len(entries)
        )

    def test_contract_rejects_traversal_and_absolute_paths(self) -> None:
        for invalid_path in ("../outside", "C:/outside", "/outside", "a\\b"):
            with self.subTest(path=invalid_path):
                contract = copy.deepcopy(self.contract)
                contract["entries"][0]["source"] = invalid_path
                errors = self.validate_contract(contract)
                self.assertTrue(
                    any("source" in error and "must" in error for error in errors),
                    errors,
                )

    def test_portable_path_subset_rejects_windows_aliases_and_git_metadata(self) -> None:
        portable_schema = self.contract_schema["$defs"]["portablePath"]
        invalid_paths = (
            "SDP/NUL",
            "SDP/con.txt",
            "SDP/COM1.log",
            "SDP/Lpt9",
            "SDP/CONIN$",
            "SDP/trailing.",
            "SDP/trailing ",
            "SDP/trailing/",
            "SDP/alternate:data",
            "SDP/less<than",
            "SDP/question?mark",
            "SDP/RELEAS~1.MD",
            "SDP/COM¹.txt",
            "SDP/LPT²",
            "SDP/control\x1fcharacter",
            "SDP/.git/config",
        )
        for value in invalid_paths:
            with self.subTest(value=value):
                self.assertIsNotNone(VALIDATE.portable_relative_path_error(value))
                self.assertTrue(
                    VALIDATE.validate_json(value, portable_schema, "portable path")
                )

    def test_every_install_contract_path_class_uses_the_portable_subset(self) -> None:
        mutations = {
            "source": lambda contract: contract["entries"][0].update(
                source="Toolkit/payload/NUL.txt"
            ),
            "destination": lambda contract: contract["entries"][0].update(
                destination="SDP/.git/config"
            ),
            "governing.schema": lambda contract: contract["entries"][0][
                "governing"
            ].update(schema="Toolkit/schemas/CON.json"),
            "exclusion.path": lambda contract: contract["exclusions"][0].update(
                path="Releases/trailing."
            ),
        }
        for field, mutate in mutations.items():
            with self.subTest(field=field):
                contract = copy.deepcopy(self.contract)
                mutate(contract)
                errors = self.validate_contract(contract)
                self.assertTrue(
                    any(
                        field.split(".")[-1] in error
                        and any(
                            marker in error
                            for marker in (
                                "does not match",
                                "not valid under any",
                                "reserved device segment",
                            )
                        )
                        for error in errors
                    ),
                    errors,
                )

    def test_exclusion_paths_use_case_insensitive_collision_keys(self) -> None:
        contract = copy.deepcopy(self.contract)
        duplicate = copy.deepcopy(contract["exclusions"][0])
        duplicate["path"] = duplicate["path"].swapcase()
        contract["exclusions"].append(duplicate)
        errors = self.validate_contract(contract)
        self.assertTrue(any("case-colliding exclusion" in error for error in errors), errors)

    def test_contract_detects_toolkit_version_disagreement(self) -> None:
        contract = copy.deepcopy(self.contract)
        contract["toolkitVersion"] = "0.3.0"
        self.assertTrue(
            any(
                "toolkitVersion differs from SDP.manifest.yaml" in error
                for error in self.validate_contract(contract)
            )
        )

    def test_governing_schemas_require_their_canonical_capabilities(self) -> None:
        expected_pairs = {
            "Toolkit/schemas/fix-record.schema.json": "sdp.release.v1",
            "Toolkit/schemas/release-record.schema.json": "sdp.release.v1",
            "Toolkit/schemas/SDP-project-manifest.schema.json": "sdp.manifest.v1",
            "Toolkit/schemas/installed-toolkit-manifest.schema.json": "sdp.manifest.v1",
            "Toolkit/schemas/current-index.schema.json": "sdp.traceability.current-index.v1",
            "Toolkit/schemas/relations.schema.json": "sdp.traceability.relations.v1",
            "Toolkit/schemas/ledger-event.schema.json": "sdp.traceability.ledger-events.v1",
        }
        self.assertEqual(VALIDATE.CANONICAL_GOVERNING_CAPABILITIES, expected_pairs)
        for schema, capability in expected_pairs.items():
            with self.subTest(schema=schema, capability=capability):
                contract = copy.deepcopy(self.contract)
                governed_entry = next(
                    entry
                    for entry in contract["entries"]
                    if (entry["governing"].get("schema") or "").casefold()
                    == schema.casefold()
                )
                governed_entry["governing"]["capability"] = "sdp.install.v1"
                errors = self.validate_contract(contract)
                self.assertTrue(
                    any(
                        f"{schema.rsplit('/', 1)[-1]} requires capability {capability}"
                        in error
                        for error in errors
                    ),
                    errors,
                )

    def test_contract_detects_duplicate_ids_and_destinations(self) -> None:
        contract = copy.deepcopy(self.contract)
        contract["entries"][1]["id"] = contract["entries"][0]["id"]
        contract["entries"][1]["destination"] = contract["entries"][0][
            "destination"
        ].lower()
        errors = self.validate_contract(contract)
        self.assertTrue(any("duplicate entry ID" in error for error in errors), errors)
        self.assertTrue(any("case-colliding destination" in error for error in errors), errors)

    def test_contract_rejects_case_insensitive_destination_prefix_conflicts(self) -> None:
        for child_first in (False, True):
            with self.subTest(child_first=child_first):
                contract = copy.deepcopy(self.contract)
                project_agents = next(
                    entry for entry in contract["entries"] if entry["id"] == "project-agents"
                )
                project_agents["destination"] = "agents.md/README.md"
                if child_first:
                    managed_agents = next(
                        entry
                        for entry in contract["entries"]
                        if entry["id"] == "managed-agents"
                    )
                    contract["entries"] = [
                        entry
                        for entry in contract["entries"]
                        if entry["id"] != "managed-agents"
                    ] + [managed_agents]
                errors = self.validate_contract(contract)
                self.assertTrue(
                    any("ancestor/descendant destination conflict" in error for error in errors),
                    errors,
                )

    def test_contract_detects_missing_source_and_inventory_drift(self) -> None:
        contract = copy.deepcopy(self.contract)
        contract["entries"][0]["source"] = "Toolkit/payload/missing.txt"
        errors = self.validate_contract(contract)
        self.assertTrue(any("referenced file does not exist" in error for error in errors), errors)
        self.assertTrue(any("unlisted files" in error for error in errors), errors)

    def test_contract_requires_live_record_exclusions(self) -> None:
        contract = copy.deepcopy(self.contract)
        contract["exclusions"] = [
            exclusion
            for exclusion in contract["exclusions"]
            if exclusion["path"] != "Releases"
        ]
        errors = self.validate_contract(contract)
        self.assertTrue(
            any("omits required live/legacy exclusions" in error for error in errors),
            errors,
        )

    def test_contract_rejects_an_explicitly_excluded_source(self) -> None:
        contract = copy.deepcopy(self.contract)
        project_agents = next(
            entry for entry in contract["entries"] if entry["id"] == "project-agents"
        )
        project_agents["source"] = (
            "Toolkit/payload/project-root/AGENTS-project.md.template"
        )
        errors = self.validate_contract(contract)
        self.assertTrue(any("source is explicitly excluded" in error for error in errors), errors)

    def test_project_owned_policy_cannot_refresh_or_replace(self) -> None:
        contract = copy.deepcopy(self.contract)
        project_agents = next(
            entry for entry in contract["entries"] if entry["id"] == "project-agents"
        )
        project_agents["refreshPolicy"] = "always"
        project_agents["forcePolicy"] = "replace-managed"
        errors = VALIDATE.validate_json(
            contract, self.contract_schema, "SDP install contract"
        )
        self.assertTrue(any("'never' was expected" in error for error in errors), errors)
        self.assertTrue(any("'preserve' was expected" in error for error in errors), errors)

    def test_migration_and_backup_policy_combinations_are_constrained(self) -> None:
        for changes in (
            {"backupPolicy": "migration-aware", "migrationPolicy": "none"},
            {"backupPolicy": "before-replace", "refreshPolicy": "never"},
        ):
            with self.subTest(changes=changes):
                contract = copy.deepcopy(self.contract)
                entry = next(
                    item
                    for item in contract["entries"]
                    if item["id"] == "managed-framework-readme"
                )
                entry.update(changes)
                self.assertTrue(
                    VALIDATE.validate_json(
                        contract, self.contract_schema, "SDP install contract"
                    )
                )

    def test_neutral_template_payload_contains_no_active_records(self) -> None:
        template_root = ROOT / "Toolkit/project-templates"
        files = [path.relative_to(template_root).as_posix() for path in template_root.rglob("*") if path.is_file()]
        self.assertNotIn("sdp-root/Traceability/Ledger.ndjson", files)
        self.assertFalse(any(Path(path).name.startswith("REL-") for path in files))
        self.assertFalse(any(Path(path).name == "ScrumIterations.md" for path in files))
        current = read_yaml(template_root / "sdp-root/Traceability/CurrentIndex.yaml")
        self.assertTrue(all(value is None for value in current["active"].values()))
        self.assertIsNone(current["release"]["activeReleaseId"])
        relations = read_yaml(template_root / "sdp-root/Traceability/Relations.yaml")
        self.assertTrue(relations)
        self.assertTrue(all(value == {} for value in relations.values()))

    def test_install_plan_schema_enforces_action_mutation_truth(self) -> None:
        schema = json.loads(
            (ROOT / "Toolkit/schemas/SDP-install-plan.schema.json").read_text(
                encoding="utf-8"
            )
        )
        plan = {
            "schemaVersion": "1.0",
            "manifestSchemaVersion": "1.0",
            "orderingPolicy": "migration-first-manifest-order-v1",
            "mode": "plan",
            "toolkitVersion": "0.2.0",
            "installedToolkitVersion": None,
            "options": {
                "initializeProjectStructure": False,
                "forceManagedFiles": False,
            },
            "canApply": True,
            "actions": [
                {
                    "sequence": 1,
                    "action": "create",
                    "entryId": "project-manifest",
                    "source": "Toolkit/project-templates/sdp-root/SDP-project.manifest.yaml",
                    "generator": None,
                    "targetSource": None,
                    "targetSourceSha256": None,
                    "destinationPrecondition": None,
                    "destination": "SDP/SDP-project.manifest.yaml",
                    "ownership": "project-owned",
                    "reason": "missing-target",
                    "mutatesTarget": True,
                    "oldToolkitVersion": None,
                    "newToolkitVersion": "0.2.0",
                }
            ],
        }
        self.assertEqual(VALIDATE.validate_json(plan, schema, "install plan"), [])
        plan["actions"][0]["mutatesTarget"] = False
        self.assertTrue(VALIDATE.validate_json(plan, schema, "install plan"))

    def test_install_plan_reason_vocabulary_is_exact_and_normative(self) -> None:
        reason_schema = self.plan_schema["$defs"]["action"]["properties"]["reason"]
        self.assertEqual(
            set(reason_schema["enum"]), set(VALIDATE.PLAN_REASON_CONDITIONS)
        )
        plan = self.valid_plan()
        plan["actions"][0]["reason"] = "content-matches"
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("reason 'content-matches' requires" in error for error in errors))

    def test_install_plan_semantics_reject_sequence_version_and_manifest_drift(self) -> None:
        mutations = {
            "sequence": lambda plan: plan["actions"][0].update(sequence=2),
            "new version": lambda plan: plan["actions"][0].update(
                newToolkitVersion="0.3.0"
            ),
            "old version": lambda plan: plan["actions"][0].update(
                oldToolkitVersion="0.1.0"
            ),
            "source": lambda plan: plan["actions"][0].update(
                source="Toolkit/project-templates/sdp-root/RELEASE-NOTES.md"
            ),
            "generator": lambda plan: plan["actions"][0].update(
                source=None, generator="empty-ledger"
            ),
            "destination": lambda plan: plan["actions"][0].update(
                destination="SDP/elsewhere.yaml"
            ),
            "entry": lambda plan: plan["actions"][0].update(entryId="missing-entry"),
        }
        for name, mutate in mutations.items():
            with self.subTest(name=name):
                plan = self.valid_plan()
                mutate(plan)
                errors = VALIDATE.validate_install_plan(
                    plan, self.plan_schema, "install plan", self.contract
                )
                self.assertTrue(errors, name)

        plan = self.valid_plan()
        duplicate = copy.deepcopy(plan["actions"][0])
        plan["actions"].append(duplicate)
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("unique, ordered, and contiguous" in error for error in errors))

    def test_install_plan_enforces_canonical_manifest_order(self) -> None:
        plan = self.valid_plan()
        plan["actions"] = [
            self.plan_action("managed-fix-record-template", sequence=1),
            self.plan_action("managed-framework-readme", sequence=2),
        ]
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(
            any("installation-manifest array order" in error for error in errors),
            errors,
        )

        migration = {
            **self.plan_action("managed-agents", sequence=2),
            "action": "migrate",
            "source": None,
            "targetSource": "AGENTS.md",
            "targetSourceSha256": "a" * 64,
            "destinationPrecondition": "absent",
            "destination": "AGENTS-project.migration-sha256-" + "a" * 64 + ".md",
            "ownership": "project-owned",
            "reason": "preserve-existing-agents-conflict",
        }
        plan["actions"] = [
            self.plan_action("managed-agents", sequence=1),
            migration,
        ]
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("migrations must precede" in error for error in errors), errors)

    def test_install_plan_enforces_adjacent_identity_equal_backup_pair(self) -> None:
        backup = self.plan_action(
            "managed-framework-readme",
            sequence=1,
            action="backup",
            reason="backup-before-replace",
        )
        replacement = self.plan_action(
            "managed-framework-readme",
            sequence=3,
            action="replace",
            reason="refresh-managed-content",
        )
        plan = self.valid_plan()
        plan["actions"] = [
            backup,
            self.plan_action("managed-fix-record-template", sequence=2),
            replacement,
        ]
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(
            any("immediately precede" in error or "immediately preceding" in error for error in errors),
            errors,
        )

        plan["actions"] = [backup, {**replacement, "sequence": 2}]
        plan["actions"][1]["destination"] = "SDP/Framework/other.md"
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("disagrees on destination" in error for error in errors), errors)

        plan["actions"] = [
            self.plan_action("managed-framework-readme", sequence=1),
            self.plan_action("managed-framework-readme", sequence=2),
        ]
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("non-canonical second action" in error for error in errors), errors)

    def test_install_plan_requires_machine_readable_ordering_policy(self) -> None:
        plan = self.valid_plan()
        plan["orderingPolicy"] = "private-path-sort-v1"
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("orderingPolicy" in error for error in errors), errors)

        contract = copy.deepcopy(self.contract)
        contract["orderingPolicy"] = "private-path-sort-v1"
        self.assertTrue(
            any("orderingPolicy" in error for error in self.validate_contract(contract))
        )

    def test_target_to_target_migration_plan_is_explicit_and_deterministic(self) -> None:
        plan = self.valid_plan()
        action = plan["actions"][0]
        action.update(
            action="migrate",
            entryId="managed-agents",
            source=None,
            generator=None,
            targetSource="AGENTS.md",
            targetSourceSha256="a" * 64,
            destinationPrecondition="absent",
            destination="AGENTS-project.migration-sha256-"
            + "a" * 64
            + ".md",
            ownership="project-owned",
            reason="preserve-existing-agents-conflict",
        )
        self.assertEqual(
            VALIDATE.validate_install_plan(
                plan, self.plan_schema, "install plan", self.contract
            ),
            [],
        )
        action["destination"] = "AGENTS-project.migration-timestamp.md"
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("source-hash migration name" in error for error in errors))
        action["destination"] = "AGENTS-project.migration-sha256-" + "b" * 64 + ".md"
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("source-hash migration name" in error for error in errors))

    def test_blocked_plan_cannot_mix_block_and_non_block_actions(self) -> None:
        plan = self.valid_plan()
        normal_action = copy.deepcopy(plan["actions"][0])
        normal_action["sequence"] = 2
        block_action = plan["actions"][0]
        block_action.update(
            action="block",
            source=None,
            generator=None,
            reason="malformed-project-manifest",
            mutatesTarget=False,
        )
        plan["canApply"] = False
        plan["actions"].append(normal_action)
        self.assertEqual(
            VALIDATE.validate_json(plan, self.plan_schema, "install plan"), []
        )
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("and no other actions" in error for error in errors), errors)

    def test_block_reason_must_name_its_canonical_manifest_entry(self) -> None:
        plan = self.valid_plan()
        action = plan["actions"][0]
        action.update(
            action="block",
            source=None,
            generator=None,
            reason="downgrade-blocked",
            mutatesTarget=False,
        )
        plan["canApply"] = False
        errors = VALIDATE.validate_install_plan(
            plan, self.plan_schema, "install plan", self.contract
        )
        self.assertTrue(any("contradicts block reason" in error for error in errors), errors)

    def test_manifest_materializes_a_valid_neutral_project_without_live_state(self) -> None:
        generators = {
            generator["id"]: generator for generator in self.contract["generators"]
        }
        for initialize in (False, True):
            with self.subTest(initialize=initialize), tempfile.TemporaryDirectory() as directory:
                project_root = Path(directory)
                (project_root / ".git").mkdir()
                for entry in self.contract["entries"]:
                    if entry["selectionPolicy"] == "initialize-only" and not initialize:
                        continue
                    destination = project_root / Path(*entry["destination"].split("/"))
                    destination.parent.mkdir(parents=True, exist_ok=True)
                    if entry["kind"] == "copied":
                        source = ROOT / Path(*entry["source"].split("/"))
                        destination.write_bytes(source.read_bytes())
                        continue
                    generator = generators[entry["generator"]]
                    if generator["type"] == "empty-ledger":
                        destination.write_text(generator["content"], encoding="utf-8")
                    elif generator["type"] == "installed-toolkit-manifest":
                        installed = copy.deepcopy(generator["facts"])
                        installed["toolkitInstalledAt"] = "2026-07-13T20:00:00Z"
                        installed["sourceCommit"] = None
                        write_yaml(destination, installed)
                    else:  # pragma: no cover - schema and semantic checks constrain v1.
                        self.fail(f"Unknown generator: {generator['type']}")

                self.assertEqual(VALIDATE.validate_project(project_root), [])
                self.assertEqual(
                    [path for path in project_root.rglob(".git") if path.is_dir()],
                    [project_root / ".git"],
                )
                self.assertFalse((project_root / "Toolkit").exists())
                self.assertFalse((project_root / "SDP/Releases/REL-0.2.0.yaml").exists())
                self.assertEqual(
                    (project_root / "SDP/Traceability/Ledger.ndjson").read_bytes(), b""
                )
                self.assertFalse(any((project_root / "SDP/Sprints").glob("Sprint-*")))
                self.assertEqual(
                    (project_root / "SDP/Sprints/README.md").exists(), initialize
                )


class BuildIdentityTests(unittest.TestCase):
    def test_unreleased_display_and_coordinates(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SDP").mkdir()
            manifest = {
                "release": {
                    "currentVersion": "0.7.0",
                    "nextTargetVersion": "0.8.0",
                    "state": "unreleased",
                },
                "development": {
                    "sprintId": "Sprint-026",
                    "refactorId": None,
                    "iterationId": "SPI-002",
                    "sliceId": "SPS-003",
                    "fixId": None,
                    "revision": 1,
                },
            }
            (root / "SDP/SDP-project.manifest.yaml").write_text(yaml.safe_dump(manifest), encoding="utf-8")
            identity = BUILD.generate_identity(root, "2026-01-01T00:00:00Z")
            self.assertEqual(identity["releaseVersion"], "0.8.0")
            self.assertIn("0.8.0-dev", identity["displayVersion"])
            self.assertIn("s026.i002.sl003.r001", identity["developmentId"])

    def test_released_build_uses_current_not_next_target(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SDP").mkdir()
            manifest = {
                "release": {
                    "currentVersion": "1.2.3",
                    "nextTargetVersion": "1.3.0",
                    "state": "released",
                },
                "development": {},
            }
            (root / "SDP/SDP-project.manifest.yaml").write_text(yaml.safe_dump(manifest), encoding="utf-8")
            identity = BUILD.generate_identity(root, "2026-01-01T00:00:00Z")
            self.assertEqual(identity["releaseVersion"], "1.2.3")
            self.assertEqual(identity["displayVersion"], "1.2.3")

    def test_yanked_build_is_not_presented_as_development(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SDP").mkdir()
            manifest = {
                "release": {"currentVersion": "1.2.3", "state": "yanked"},
                "development": {},
            }
            (root / "SDP/SDP-project.manifest.yaml").write_text(yaml.safe_dump(manifest), encoding="utf-8")
            identity = BUILD.generate_identity(root, "2026-01-01T00:00:00Z")
            self.assertIn("1.2.3-yanked", identity["displayVersion"])
            self.assertIn(".yanked", identity["developmentId"])
            self.assertNotIn("-dev", identity["displayVersion"])

    def test_rejects_simultaneous_sprint_and_refactor(self) -> None:
        with tempfile.TemporaryDirectory() as directory:
            root = Path(directory)
            (root / "SDP").mkdir()
            manifest = {
                "release": {"currentVersion": "1.0.0", "state": "unreleased"},
                "development": {"sprintId": "Sprint-001", "refactorId": "Refactor-001"},
            }
            (root / "SDP/SDP-project.manifest.yaml").write_text(yaml.safe_dump(manifest), encoding="utf-8")
            with self.assertRaises(ValueError):
                BUILD.generate_identity(root, "2026-01-01T00:00:00Z")


class PinnedGhSdpCompatibilityTests(unittest.TestCase):
    FIXTURE = (
        ROOT
        / "Toolkit/tests/fixtures"
        / "gh-sdp-ed205c1ef193ab8a6e5cd1c50e558c3049ce6def"
    )

    def test_exact_exported_closure_validates_offline(self) -> None:
        self.assertTrue(self.FIXTURE.is_dir())
        self.assertFalse(any(path.name == ".git" for path in self.FIXTURE.rglob("*")))
        self.assertEqual(
            git_tree_sha(self.FIXTURE),
            "54f0e5854fd34e5d8bcb301f4921b956a2030e61",
        )
        self.assertEqual(VALIDATE.validate_project(self.FIXTURE), [])


class ProjectValidationTests(unittest.TestCase):
    def setUp(self) -> None:
        self.temporary_directory = tempfile.TemporaryDirectory()
        self.root = create_valid_consuming_project(Path(self.temporary_directory.name))

    def tearDown(self) -> None:
        self.temporary_directory.cleanup()

    def assert_error_contains(self, text: str) -> list[str]:
        errors = VALIDATE.validate_project(self.root)
        self.assertTrue(
            any(text in error for error in errors),
            f"Expected {text!r} in validation errors:\n" + "\n".join(errors),
        )
        return errors

    def test_valid_consuming_project_does_not_need_toolkit_repo_layout(self) -> None:
        self.assertFalse((self.root / "Toolkit").exists())
        self.assertEqual(VALIDATE.validate_project(self.root), [])

    def test_project_cli_mode(self) -> None:
        stdout = StringIO()
        stderr = StringIO()
        with redirect_stdout(stdout), redirect_stderr(stderr):
            result = VALIDATE.main(
                ["--mode", "project", "--project-root", str(self.root)]
            )
        self.assertEqual(result, 0, stderr.getvalue())
        self.assertIn("consuming-project validation passed", stdout.getvalue())

    def test_project_cli_requires_project_root(self) -> None:
        with redirect_stderr(StringIO()), self.assertRaises(SystemExit) as context:
            VALIDATE.main(["--mode", "project"])
        self.assertEqual(context.exception.code, 2)

    def test_malformed_project_manifest_fails_clearly(self) -> None:
        (self.root / "SDP/SDP-project.manifest.yaml").write_text(
            "schemaVersion: [\n", encoding="utf-8"
        )
        self.assert_error_contains("Cannot parse YAML")

    def test_malformed_installed_manifest_fails_clearly(self) -> None:
        (self.root / "SDP/Framework/installed-toolkit.manifest.yaml").write_text(
            "skills: [\n", encoding="utf-8"
        )
        self.assert_error_contains("installed-toolkit.manifest.yaml")
        self.assert_error_contains("Cannot parse YAML")

    def test_unsupported_project_manifest_schema_fails_clearly(self) -> None:
        path = self.root / "SDP/SDP-project.manifest.yaml"
        manifest = read_yaml(path)
        manifest["schemaVersion"] = "9.0"
        write_yaml(path, manifest)
        self.assert_error_contains("unsupported schemaVersion '9.0'")

    def test_unsupported_installed_manifest_schema_fails_clearly(self) -> None:
        path = self.root / "SDP/Framework/installed-toolkit.manifest.yaml"
        manifest = read_yaml(path)
        manifest["schemaVersion"] = "9.0"
        write_yaml(path, manifest)
        self.assert_error_contains("unsupported schemaVersion '9.0'")

    def test_malformed_current_index_fails_schema_validation(self) -> None:
        path = self.root / "SDP/Traceability/CurrentIndex.yaml"
        current = read_yaml(path)
        active = current["active"]
        assert isinstance(active, dict)
        active["iteration"] = None
        write_yaml(path, current)
        self.assert_error_contains("SDP/Traceability/CurrentIndex.yaml.active")

    def test_neutral_current_index_and_relations_are_valid(self) -> None:
        project_path = self.root / "SDP/SDP-project.manifest.yaml"
        project = read_yaml(project_path)
        release = project["release"]
        development = project["development"]
        assert isinstance(release, dict) and isinstance(development, dict)
        release["currentVersion"] = "0.0.0"
        release["nextTargetVersion"] = None
        release["latestTag"] = None
        release["latestCommit"] = None
        for key in development:
            development[key] = None
        write_yaml(project_path, project)
        current = {
            "project": {"name": "TBD", "status": "active-development"},
            "release": {
                "previousVersion": None,
                "activeReleaseId": None,
                "targetVersion": None,
                "state": "unreleased",
            },
            "active": {
                "sprint": None,
                "refactor": None,
                "iteration": None,
                "slice": None,
                "fix": None,
                "revision": None,
            },
        }
        write_yaml(self.root / "SDP/Traceability/CurrentIndex.yaml", current)
        relations = {
            name: {}
            for name in (
                "requirements",
                "designs",
                "sprints",
                "refactors",
                "iterations",
                "slices",
                "fixes",
                "reviews",
                "verification",
                "migrations",
                "releases",
            )
        }
        write_yaml(self.root / "SDP/Traceability/Relations.yaml", relations)
        for directory in ("Releases", "Fixes"):
            for path in (self.root / "SDP" / directory).glob("*.yaml"):
                path.unlink()
        self.assertEqual(VALIDATE.validate_project(self.root), [])

    def test_malformed_relations_fails_schema_validation(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        relations["sprints"] = []
        write_yaml(path, relations)
        self.assert_error_contains("SDP/Traceability/Relations.yaml.sprints")

    def test_namespaced_relation_category_and_detail_fields_are_extensible(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        relations["x-acme.example"] = {
            "DEPLOYMENT-001": {"environment": "test", "approved": True}
        }
        slices = relations["slices"]
        assert isinstance(slices, dict) and isinstance(slices["SPS-001"], dict)
        slices["SPS-001"]["x-acme.example:ticket"] = "ACME-123"
        write_yaml(path, relations)
        self.assertEqual(VALIDATE.validate_project(self.root), [])

    def test_invalid_ledger_json_fails_with_line_number(self) -> None:
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            "{not-json}\n", encoding="utf-8"
        )
        errors = self.assert_error_contains("invalid JSON")
        self.assertTrue(any("Ledger.ndjson:1" in error for error in errors))

    def test_invalid_generic_ledger_event_fails(self) -> None:
        event = {
            "schemaVersion": "1.0",
            "eventId": "EVT-WORK-0001",
            "eventType": "work-started",
            "occurredAt": "2026-07-13T20:00:00Z",
            "commit": None,
            "subjectId": "SPS-001",
            "payload": {},
        }
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            json.dumps(event) + "\n", encoding="utf-8"
        )
        self.assert_error_contains("'actor' is a required property")

    def test_unsupported_ledger_event_schema_fails_clearly(self) -> None:
        event = {
            "schemaVersion": "2.0",
            "eventId": "EVT-WORK-0001",
            "eventType": "work-started",
            "occurredAt": "2026-07-13T20:00:00Z",
            "actor": "fixture",
            "commit": None,
            "subjectId": "SPS-001",
            "payload": {},
        }
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            json.dumps(event) + "\n", encoding="utf-8"
        )
        self.assert_error_contains("unsupported schemaVersion '2.0'")

    def test_invalid_release_event_fails_stricter_schema(self) -> None:
        event = {
            "schemaVersion": "1.0",
            "eventId": "EVT-REL-0001",
            "eventType": "release-tag-created",
            "occurredAt": "2026-07-13T20:00:00Z",
            "releaseId": "REL-1.1.0",
            "actor": "fixture",
            "commit": None,
            "payload": {},
        }
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            json.dumps(event) + "\n", encoding="utf-8"
        )
        errors = self.assert_error_contains("release event")
        self.assertTrue(any("'tag' is a required property" in error for error in errors))

    def test_dangling_relation_path_fails(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        slices = relations["slices"]
        assert isinstance(slices, dict) and isinstance(slices["SPS-001"], dict)
        slices["SPS-001"]["path"] = "Sprints/Sprint-001/missing.md"
        write_yaml(path, relations)
        self.assert_error_contains("referenced file does not exist")

    def test_dangling_relation_id_fails(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        slices = relations["slices"]
        assert isinstance(slices, dict) and isinstance(slices["SPS-001"], dict)
        slices["SPS-001"]["iteration"] = "SPI-999"
        write_yaml(path, relations)
        self.assert_error_contains("dangling ID SPI-999")

    def test_review_and_verification_ledger_subjects_must_resolve(self) -> None:
        cases = (
            ("review-started", "REV-MISSING", "Relations.reviews"),
            ("verification-started", "VER-MISSING", "Relations.verification"),
        )
        for event_type, subject_id, expected in cases:
            with self.subTest(subject_id=subject_id):
                event = {
                    "schemaVersion": "1.0",
                    "eventId": f"EVT-{subject_id}",
                    "eventType": event_type,
                    "occurredAt": "2026-07-13T20:00:00Z",
                    "actor": "fixture",
                    "commit": None,
                    "subjectId": subject_id,
                    "payload": {},
                }
                (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
                    json.dumps(event) + "\n", encoding="utf-8"
                )
                self.assert_error_contains(expected)

    def test_review_and_verification_links_are_bidirectional(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        slices = relations["slices"]
        assert isinstance(slices, dict) and isinstance(slices["SPS-001"], dict)
        del slices["SPS-001"]["review"]
        write_yaml(path, relations)
        self.assert_error_contains("missing reverse link")

        relations = read_yaml(path)
        slices = relations["slices"]
        verification = relations["verification"]
        assert isinstance(slices, dict) and isinstance(slices["SPS-001"], dict)
        assert isinstance(verification, dict) and isinstance(verification["VER-001"], dict)
        slices["SPS-001"]["review"] = "REV-001"
        verification["VER-001"]["slice"] = "SPS-999"
        write_yaml(path, relations)
        self.assert_error_contains("dangling ID SPS-999")

    def test_migration_release_links_are_bidirectional(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        migrations = relations["migrations"]
        releases = relations["releases"]
        assert isinstance(migrations, dict) and isinstance(releases, dict)
        assert isinstance(releases["REL-1.1.0"], dict)
        migrations["MIG-001"] = {"release": "REL-1.1.0"}
        releases["REL-1.1.0"]["migrations"] = ["MIG-001"]
        write_yaml(path, relations)
        self.assertEqual(VALIDATE.validate_project(self.root), [])
        del releases["REL-1.1.0"]["migrations"]
        write_yaml(path, relations)
        self.assert_error_contains("missing reverse link")

    def test_release_and_fix_records_require_canonical_yaml_extension(self) -> None:
        cases = (
            self.root / "SDP/Releases/REL-1.1.0.yaml",
            self.root / "SDP/Fixes/FIX-1.1.0-001.yaml",
        )
        for path in cases:
            with self.subTest(path=path.name):
                renamed = path.with_suffix(".yml")
                path.rename(renamed)
                renamed.write_text("schemaVersion: [\n", encoding="utf-8")
                self.assert_error_contains("canonical .yaml extension")
                renamed.rename(path)

    def test_release_relation_cannot_substitute_non_record_content(self) -> None:
        (self.root / "SDP/Releases/REL-1.1.0.yaml").unlink()
        relation_path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(relation_path)
        relations["releases"]["REL-1.1.0"]["releaseRecord"] = "SDP/RELEASE-NOTES.md"
        write_yaml(relation_path, relations)
        self.assert_error_contains("must reference a canonical .yaml release record")

    def test_fix_relation_cannot_substitute_non_record_content(self) -> None:
        (self.root / "SDP/Fixes/FIX-1.1.0-001.yaml").unlink()
        relation_path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(relation_path)
        relations["fixes"]["FIX-1.1.0-001"]["path"] = "SDP/RELEASE-NOTES.md"
        write_yaml(relation_path, relations)
        self.assert_error_contains("must reference a canonical .yaml Fix record")

    def test_unreleased_publication_fields_must_remain_null(self) -> None:
        release_path = self.root / "SDP/Releases/REL-1.1.0.yaml"
        release = read_yaml(release_path)
        release["releasePreparationCommit"] = "abcdef1234567890"
        release["gitTag"] = "v1.1.0"
        release["githubReleaseUrl"] = "https://example.invalid/releases/1.1.0"
        write_yaml(release_path, release)
        self.assert_error_contains("must be null")
        release["releasePreparationCommit"] = None
        release["gitTag"] = None
        release["githubReleaseUrl"] = None
        write_yaml(release_path, release)

        relation_path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(relation_path)
        releases = relations["releases"]
        assert isinstance(releases, dict) and isinstance(releases["REL-1.1.0"], dict)
        releases["REL-1.1.0"].update(
            releaseCommit="abcdef1234567890",
            gitTag="v1.1.0",
            githubRelease="https://example.invalid/releases/1.1.0",
        )
        write_yaml(relation_path, relations)
        self.assert_error_contains("while release is unreleased")

    def test_release_record_and_relation_identity_must_agree(self) -> None:
        path = self.root / "SDP/Traceability/Relations.yaml"
        relations = read_yaml(path)
        releases = relations["releases"]
        assert isinstance(releases, dict) and isinstance(releases["REL-1.1.0"], dict)
        releases["REL-1.1.0"]["reviews"] = []
        write_yaml(path, relations)
        self.assert_error_contains("reviewRecords differs from Relations.releases")

    def test_project_latest_publication_identity_matches_current_version(self) -> None:
        path = self.root / "SDP/SDP-project.manifest.yaml"
        manifest = read_yaml(path)
        manifest["release"]["latestTag"] = "v9.9.9"
        write_yaml(path, manifest)
        self.assert_error_contains("latestTag does not match release.currentVersion")

    def test_release_event_tag_must_match_release_id(self) -> None:
        event = {
            "schemaVersion": "1.0",
            "eventId": "EVT-REL-TAG-0001",
            "eventType": "release-tag-created",
            "occurredAt": "2026-07-13T20:00:00Z",
            "releaseId": "REL-1.1.0",
            "actor": "fixture",
            "commit": "abcdef1234567890",
            "payload": {"tag": "v1.2.0"},
        }
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            json.dumps(event) + "\n", encoding="utf-8"
        )
        self.assert_error_contains("does not match releaseId REL-1.1.0")

    def test_empty_ledger_is_valid(self) -> None:
        ledger = self.root / "SDP/Traceability/Ledger.ndjson"
        ledger.write_bytes(b"")
        self.assertEqual(VALIDATE.validate_project(self.root), [])

    def test_namespaced_project_event_extension_is_valid(self) -> None:
        event = {
            "schemaVersion": "1.0",
            "eventId": "EVT-X-0001",
            "eventType": "x-acme.example:deployment-approved",
            "occurredAt": "2026-07-13T20:00:00Z",
            "actor": "fixture",
            "commit": None,
            "payload": {"environment": "test"},
        }
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            json.dumps(event) + "\n", encoding="utf-8"
        )
        self.assertEqual(VALIDATE.validate_project(self.root), [])

    def test_unnamespaced_project_event_extension_is_rejected(self) -> None:
        event = {
            "schemaVersion": "1.0",
            "eventId": "EVT-X-0001",
            "eventType": "deployment-approved",
            "occurredAt": "2026-07-13T20:00:00Z",
            "actor": "fixture",
            "commit": None,
            "payload": {},
        }
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            json.dumps(event) + "\n", encoding="utf-8"
        )
        self.assert_error_contains("is not valid under any of the given schemas")

    def test_duplicate_ledger_event_id_is_rejected(self) -> None:
        event = {
            "schemaVersion": "1.0",
            "eventId": "EVT-WORK-0001",
            "eventType": "work-started",
            "occurredAt": "2026-07-13T20:00:00Z",
            "actor": "fixture",
            "commit": None,
            "subjectId": "SPS-001",
            "payload": {},
        }
        line = json.dumps(event) + "\n"
        (self.root / "SDP/Traceability/Ledger.ndjson").write_text(
            line + line, encoding="utf-8"
        )
        self.assert_error_contains("duplicate eventId EVT-WORK-0001")

    def test_release_and_fix_records_are_schema_validated(self) -> None:
        release_path = self.root / "SDP/Releases/REL-1.1.0.yaml"
        release = read_yaml(release_path)
        del release["state"]
        write_yaml(release_path, release)
        self.assert_error_contains("'state' is a required property")

        write_yaml(
            release_path,
            {
                **release,
                "state": "unreleased",
            },
        )
        fix_path = self.root / "SDP/Fixes/FIX-1.1.0-001.yaml"
        fix = read_yaml(fix_path)
        fix["schemaVersion"] = "9.0"
        write_yaml(fix_path, fix)
        self.assert_error_contains("unsupported schemaVersion '9.0'")

    def test_unreleased_notes_structure_is_validated(self) -> None:
        (self.root / "SDP/RELEASE-NOTES.md").write_text(
            "# Release Notes\n\n## [1.0.0] - 2026-01-01\n", encoding="utf-8"
        )
        self.assert_error_contains("## [Unreleased]")

    def test_installed_skill_version_agreement_is_validated(self) -> None:
        skill_path = self.root / ".codex/skills/sdp-test/SKILL.md"
        skill_path.write_text(
            skill_path.read_text(encoding="utf-8").replace(
                "skillVersion: 1.0.0", "skillVersion: 2.0.0"
            ),
            encoding="utf-8",
        )
        self.assert_error_contains("skillVersion differs from installed manifest")

    def test_installed_skill_minimum_toolkit_compatibility_is_validated(self) -> None:
        skill_path = self.root / ".codex/skills/sdp-test/SKILL.md"
        skill_path.write_text(
            skill_path.read_text(encoding="utf-8").replace(
                "minimumToolkitVersion: 0.2.0", "minimumToolkitVersion: 0.3.0"
            ),
            encoding="utf-8",
        )
        self.assert_error_contains("is newer than installed Toolkit 0.2.0")

    def test_non_scalar_schema_version_fails_clearly(self) -> None:
        path = self.root / "SDP/SDP-project.manifest.yaml"
        manifest = read_yaml(path)
        manifest["schemaVersion"] = ["1.0"]
        write_yaml(path, manifest)
        self.assert_error_contains("unsupported schemaVersion ['1.0']")

    def test_installed_manifest_path_traversal_is_rejected(self) -> None:
        path = self.root / "SDP/SDP-project.manifest.yaml"
        manifest = read_yaml(path)
        installed = manifest["installed"]
        assert isinstance(installed, dict)
        installed["manifestPath"] = "../outside.yaml"
        write_yaml(path, manifest)
        self.assert_error_contains("must not contain parent traversal")


class InstallConformancePackageTests(unittest.TestCase):
    def setUp(self) -> None:
        self.contract = json.loads(
            (ROOT / "Toolkit/SDP-install.manifest.json").read_text(encoding="utf-8")
        )
        self.plan_schema = json.loads(
            (ROOT / "Toolkit/schemas/SDP-install-plan.schema.json").read_text(
                encoding="utf-8"
            )
        )
        self.package_root = ROOT / "Toolkit/conformance/install-v1"
        self.index = json.loads(
            (self.package_root / "scenarios.json").read_text(encoding="utf-8")
        )

    def test_package_is_valid_and_covers_required_scenarios(self) -> None:
        self.assertEqual(
            VALIDATE.validate_install_conformance_package(
                ROOT, self.contract, self.plan_schema
            ),
            [],
        )
        scenario_ids = {scenario["id"] for scenario in self.index["scenarios"]}
        self.assertEqual(
            scenario_ids,
            {
                "empty-default",
                "empty-initialize",
                "repeat-default",
                "repeat-initialize",
                "legacy-agents-migrate",
                "legacy-agents-conflict-new-hash-target",
                "legacy-agents-conflict-already-preserved",
                "legacy-agents-conflict-different-content",
                "legacy-agents-conflict-invalid-object",
                "upgrade-managed-content",
                "same-version-force-managed-content",
                "project-owned-content-preserved",
                "archive-source-no-git",
                "unsupported-project-schema",
                "unsupported-installed-schema",
                "downgrade-blocked",
                "malformed-manifest",
            },
        )
        covered = {
            category
            for scenario in self.index["scenarios"]
            for category in scenario["categories"]
        }
        self.assertTrue(
            VALIDATE.REQUIRED_INSTALL_CONFORMANCE_CATEGORIES.issubset(covered)
        )

    def test_all_fixture_json_uses_a_standard_parser(self) -> None:
        paths = sorted(self.package_root.rglob("*.json"))
        self.assertGreater(len(paths), 3)
        for path in paths:
            with self.subTest(path=path.relative_to(ROOT).as_posix()):
                parsed = json.loads(path.read_text(encoding="utf-8"))
                self.assertIsNotNone(parsed)

    def test_every_committed_plan_is_schema_and_semantically_valid(self) -> None:
        for scenario in self.index["scenarios"]:
            if scenario["expected"]["kind"] == "fatal":
                continue
            path = self.package_root / scenario["expected"]["path"]
            plan = json.loads(path.read_text(encoding="utf-8"))
            with self.subTest(scenario=scenario["id"]):
                self.assertEqual(
                    VALIDATE.validate_install_plan(
                        plan,
                        self.plan_schema,
                        scenario["id"],
                        self.contract,
                    ),
                    [],
                )

    def test_changed_committed_authority_is_rejected_without_regeneration(self) -> None:
        with tempfile.TemporaryDirectory() as temp:
            temp_root = Path(temp)
            copied_package = temp_root / "Toolkit/conformance/install-v1"
            shutil.copytree(self.package_root, copied_package)
            plan_path = copied_package / "expected/empty-default.plan.json"
            altered = json.loads(plan_path.read_text(encoding="utf-8"))
            altered["actions"][0]["destination"] = "SDP/wrong-location.yaml"
            plan_path.write_text(
                json.dumps(altered, indent=2) + "\n", encoding="utf-8"
            )

            errors = VALIDATE.validate_install_conformance_package(
                temp_root, self.contract, self.plan_schema
            )
            self.assertTrue(
                any(
                    "destination differs from installation entry" in error
                    for error in errors
                ),
                errors,
            )


class RepositoryValidationTests(unittest.TestCase):
    def test_generic_ledger_example_contract(self) -> None:
        example = ROOT / "examples/ledger-events.ndjson.example"
        self.assertEqual(
            len([line for line in example.read_text(encoding="utf-8").splitlines() if line.strip()]),
            4,
        )
        self.assertEqual(
            VALIDATE.validate_ledger(
                example,
                ROOT / "Toolkit/schemas",
                required=True,
                relations=None,
                label="examples/ledger-events.ndjson.example",
            ),
            [],
        )

    def test_repository_contracts(self) -> None:
        self.assertEqual(VALIDATE.validate_repository(ROOT), [])


if __name__ == "__main__":
    unittest.main()
