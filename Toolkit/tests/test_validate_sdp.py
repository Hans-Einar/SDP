from __future__ import annotations

import importlib.util
import tempfile
import unittest
import sys
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


class RepositoryValidationTests(unittest.TestCase):
    def test_repository_contracts(self) -> None:
        self.assertEqual(VALIDATE.validate_repository(ROOT), [])


if __name__ == "__main__":
    unittest.main()
