#!/usr/bin/env python3
"""Verify the pinned local GLFW source and exact upstream patch provenance."""
from pathlib import Path
import hashlib
import json
import shutil
import subprocess
import tempfile

policy = Path(__file__).resolve().parent
selected = policy.parent / "glfw"
manifest = json.loads((policy / "source-inventory.json").read_text())
expected = {row["path"]: row for row in manifest["files"]}
actual = {str(path.relative_to(selected)) for path in selected.rglob("*") if path.is_file()}
assert actual == set(expected), "Missing or extra dependency files"
for name, row in expected.items():
    assert hashlib.sha256((selected / name).read_bytes()).hexdigest() == row["selectedSHA256"], name
with tempfile.TemporaryDirectory(prefix="sdui-glfw-verify-") as directory:
    root = Path(directory)
    shutil.copytree(selected, root / "source")
    # Module-cache copies can retain read-only directory modes. Only the temporary
    # verification copy needs write permission for reverse patching.
    for path in (root / "source").rglob("*"):
        path.chmod(path.stat().st_mode | 0o200)
    (root / "source").chmod((root / "source").stat().st_mode | 0o200)
    subprocess.run(["patch", "--batch", "--reverse", "-p1", "-i", str(policy / "xim-filter.patch")],
                   cwd=root / "source", check=True, stdout=subprocess.DEVNULL)
    for name, row in expected.items():
        assert hashlib.sha256((root / "source" / name).read_bytes()).hexdigest() == row["upstreamSHA256"], name
print(f"PASS: {len(expected)} pinned files; selected bytes and exact upstream patch reversal")
