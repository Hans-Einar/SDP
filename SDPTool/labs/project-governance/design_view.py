#!/usr/bin/env python3
"""Read the isolated PGL2-C2b preview; never change its state or model."""
import argparse
import json
from pathlib import Path
import subprocess

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument("--preview", type=Path, default=Path(__file__).resolve().parent / "previews/pgl2-c2b")
parser.add_argument("--json", action="store_true")
args = parser.parse_args()
p = args.preview.resolve()
result = subprocess.run([
    str(p / "sdp-controller"), "assignment", "design-check",
    "--project", str(p / "project"), "--state", str(p / "runtime/state"),
    "--input", str(p / "design-check.json"),
], text=True, capture_output=True)
if result.returncode:
    raise SystemExit(result.stderr)
value = json.loads(result.stdout)
if args.json:
    print(json.dumps(value, indent=2))
else:
    print("PG design check: " + value["status"])
    print("Assignment: " + value["assignment"]["id"])
    print("Preparation revision: " + str(value["assignment"]["revision"]))
    for condition in value["conditions"]:
        print("  " + condition)
    if value.get("reviewRevision"):
        print("Review revision: " + value["reviewRevision"])
    print("Owner design acceptance: " + value["designAcceptance"])
    print("Dispatch: " + value["dispatch"])
    print("Synthetic assignment with real model/Blueprint services. TARGET is deliberately changed.")
