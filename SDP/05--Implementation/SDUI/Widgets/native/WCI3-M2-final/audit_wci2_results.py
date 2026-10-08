#!/usr/bin/env python3
"""Check exact native receipts from retained command-fixture event logs."""
import argparse
import collections
import gzip
import json
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument("logs", nargs="+")
args = parser.parse_args()
results = {}
for filename in args.logs:
    path = Path(filename)
    raw = gzip.decompress(path.read_bytes()) if path.suffix == ".gz" else path.read_bytes()
    events = [json.loads(line) for line in raw.splitlines()]
    opened = set()
    receipts = collections.Counter()
    for event in events:
        if event["event"] == "state":
            for surface in (event["data"]["snapshot"].get("Surfaces") or {}).values():
                if surface["Open"]:
                    opened.add(json.dumps(surface["Target"], sort_keys=True))
        if event["event"] == "dialog-result":
            receipts[json.dumps(event["data"]["Surface"], sort_keys=True)] += 1
    if set(receipts) != opened or any(count != 1 for count in receipts.values()):
        raise SystemExit("Missing, duplicate, or unpublished terminal result: " + str(path))
    closed = [e["data"] for e in events if e["event"] == "closed"]
    if len(closed) != 1 or not closed[0]["sessionClosed"] or closed[0]["pending"]:
        raise SystemExit("Incomplete native teardown: " + str(path))
    results[str(path)] = {"publishedOpenings": len(opened), "terminalResults": sum(receipts.values()), "exactlyOnce": True}
print(json.dumps(results, indent=2))
