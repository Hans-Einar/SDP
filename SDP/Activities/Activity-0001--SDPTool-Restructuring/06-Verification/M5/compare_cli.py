"""Read-only wire comparison: python compare_cli.py BASELINE CANDIDATE OUTPUT.

Uses invalid arguments and read-only Session operations; never supplies a valid
mutation request. Output includes exact binary hashes and every differing case.
"""
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import tempfile

baseline, candidate, output = map(Path, sys.argv[1:])
inventory = json.loads((Path(__file__).parents[2] / "02-Study/Operations.json").read_text())
commands = "actions invoke governance run model release-log install upgrade preview discover view tree select sdui-preview".split()
with tempfile.TemporaryDirectory(prefix="sdp-cli-parity-") as directory:
    root = Path(directory)
    request = root / "invalid.json"
    request.write_text('{"unexpected":true}')
    cases = []
    for name in commands:
        cases.extend([[name, "--definitely-invalid", "--json"],
                      [str(root), name, "--definitely-invalid", "--json"],
                      ["--json", name, "--json=broken"]])
    operations = [r["operation"].split(" ", 1)[1] for r in inventory["rows"]
                  if r["operation"].startswith("governance ")]
    operations += ["session-" + r["operation"].removeprefix("governance.session.")
                   for r in inventory["rows"] if r["operation"].startswith("governance.session.")]
    for operation in operations + ["unknown"]:
        for suffix in [[], ["--input", str(request)], ["--run", "missing"],
                       ["--instance", "missing"], ["--definitely-invalid"]]:
            cases.append([str(root), "governance", operation, "--state", str(root / "state"), "--json"] + suffix)
    differences = []
    for args in cases:
        results = []
        for binary in (baseline, candidate):
            p = subprocess.run([str(binary)] + args, cwd=root, input=b"", capture_output=True, timeout=20)
            results.append([p.returncode, p.stdout.decode(), p.stderr.decode()])
        if results[0] != results[1]:
            differences.append({"args": args, "baseline": results[0], "candidate": results[1]})
    result = {"cases": len(cases), "differences": differences,
              "stateWritten": (root / "state").exists(),
              "binaries": {str(p): hashlib.sha256(p.read_bytes()).hexdigest() for p in (baseline, candidate)}}
    output.write_text(json.dumps(result, indent=2) + "\n")
    print(json.dumps(result))
    sys.exit(bool(differences or result["stateWritten"]))
