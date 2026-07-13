from __future__ import annotations

import copy
import importlib.util
import json
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
        "reviews": {"REV-001": {"path": "CodeReview/REV-001.md"}},
        "verification": {"VER-001": {"path": "Verification/VER-001.md"}},
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
        for value in ("1", "v1.2.3", "01.2.3", "1.2.3-01", "1.2"):
            with self.subTest(value=value), self.assertRaises(ValueError):
                VALIDATE.SemVer.parse(value)


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

    def validate_contract(self, contract: object) -> list[str]:
        return VALIDATE.validate_installation_contract(
            ROOT,
            self.toolkit_manifest,
            contract_data=contract,
            check_installer_integration=False,
        )

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

    def test_contract_detects_toolkit_version_disagreement(self) -> None:
        contract = copy.deepcopy(self.contract)
        contract["toolkitVersion"] = "0.3.0"
        self.assertTrue(
            any(
                "toolkitVersion differs from SDP.manifest.yaml" in error
                for error in self.validate_contract(contract)
            )
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
