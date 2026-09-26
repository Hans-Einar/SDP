#!/usr/bin/env python3
"""Audit this exercise's authored inventory; not an SDL parser or type checker.

Python 3.6+, standard library only. Does not execute prose, validate payloads,
prove branch predicates, or establish design/implementation completeness.
"""
import argparse
import collections
import csv
import hashlib
import json
from pathlib import Path
import re
import sys

IDENT = r"[A-Z][A-Za-z0-9]*"
QUOTED = r'"(?:\\.|[^"\\])*"'
HEADER = "language sdl-mvp1-exercise version 0.1."
KINDS = set(("system container unit subsystem library layer channel message-set "
             "contract client-tool external-device port adapter feature scenario "
             "step condition outcome invariant requirement decision gap binding "
             "evidence data-type state representation composition presentation "
             "view capability functionality interface mode").split())


def audit(root, source_root=None):
    errors = []
    def check(ok, message):
        if not ok:
            errors.append(message)

    base = root.parent.parent
    entry = root / "System.design"
    manifest = json.loads((base / "source-inventory.json").read_text())
    paths = sorted(root.rglob("*.design"))
    includes = re.findall(r'^includes "([^"]+)"\.$', entry.read_text(), re.M)
    actual = set(str(p.relative_to(root)) for p in paths if p != entry)
    check(len(includes) == len(set(includes)), "Duplicate includes")
    check(set(includes) == actual, "Entry includes differ from on-disk design files")
    kinds, locations, facts = {}, {}, []
    seen = {}
    for path in paths:
        lines = path.read_text().splitlines()
        check(bool(lines) and lines[0] == HEADER, "Wrong profile: " + str(path))
        for number, line in enumerate(lines, 1):
            line = line.strip()
            if not line or line == HEADER or line.startswith("//"):
                continue
            location = "{}:{}".format(path.relative_to(root), number)
            check(line.endswith("."), "Missing terminator: " + location)
            for literal in re.findall(QUOTED, line):
                try:
                    json.loads(literal)
                except ValueError:
                    errors.append("Invalid JSON string: " + location)
            declaration = re.fullmatch(r"([a-z-]+) (" + IDENT + r")\.", line)
            if declaration:
                kind, name = declaration.groups()
                check(kind in KINDS, "Unknown declaration kind: " + location)
                check(name not in kinds, "Duplicate declaration: " + name)
                kinds[name], locations[name] = kind, location
                continue
            if not line.startswith("includes "):
                check(line not in seen, "Duplicate fact: {} / {}".format(seen.get(line), location))
                seen[line] = location
                facts.append((line, location))

    def matching(pattern):
        return [m.groups() for line, _ in facts
                for m in [re.fullmatch(pattern, line)] if m]

    for line, location in facts:
        bare = re.sub(QUOTED, "", line)
        for name in re.findall(r"\b" + IDENT + r"\b", bare):
            check(name in kinds, "Unresolved {} at {}".format(name, location))
    systems = [n for n, k in kinds.items() if k == "system"]
    check(systems == ["MVP1"], "Expected exactly one System MVP1")
    check(locations.get("MVP1", "").startswith("System.design:"), "System not declared in entry")
    for c in manifest["constituents"]:
        check(kinds.get(c["name"]) == c["kind"], "Constituent missing/mistyped: " + c["name"])
    for kind in ("container", "library"):
        expected = set(c["name"] for c in manifest["constituents"] if c["kind"] == kind)
        check(set(n for n,k in kinds.items() if k == kind) == expected,
              "Unexpected/missing registry " + kind)

    pair = r"(" + IDENT + r") {} (" + IDENT + r")\."
    parents = {}
    for parent, child in matching(pair.format("contains")):
        check(child not in parents, "Multiple containment parents: " + child)
        parents[child] = parent
    for name in parents:
        visited, current = set(), name
        while current in parents:
            if current in visited:
                errors.append("Containment cycle at " + current)
                break
            visited.add(current)
            current = parents[current]

    owners = collections.defaultdict(list)
    for owner, function in matching(pair.format("owns")):
        owners[function].append(owner)
        check(kinds.get(function) == "functionality", "Owns target not Functionality: " + function)
        check(kinds.get(owner) in ("unit","container","library","subsystem","adapter","client-tool"),
              "Invalid functionality owner: " + owner)
    functions = set(n for n,k in kinds.items() if k == "functionality")
    for name in functions:
        check(len(owners[name]) == 1, "Functionality needs exactly one owner: " + name)

    calls = {}
    for step, function, context in matching(
            r"(" + IDENT + r") calls (" + IDENT + r")(?: within (" + IDENT + r"))?\."):
        check(step not in calls, "Multiple call targets: " + step)
        calls[step] = (function, context)
        check(kinds.get(step) == "step" and function in functions, "Invalid call: " + step)
        check(context is None or kinds.get(context) == "container", "Invalid context: " + step)

    contracts = dict(matching(pair.format("follows")))
    allows = set(matching(pair.format("allows")))
    participant_pattern = (r"(?:[a-z-]+ )?(" + IDENT + r") uses (?:channel )?(" + IDENT
        + r") as (sender|receiver) of (?:message-set )?(" + IDENT
        + r")(?: in mode (" + IDENT + r"))?\.")
    participants = matching(participant_pattern)
    for endpoint, channel, role, messages, mode in participants:
        check(kinds.get(channel) == "channel" and kinds.get(messages) == "message-set",
              "Invalid channel/message-set participation: " + endpoint)
        check((contracts.get(channel), messages) in allows, "Contract disallows: " + messages)
    traffic = matching(r"(" + IDENT + r") (receives|emits) (" + IDENT + r") through (" + IDENT + r")\.")
    for step, direction, messages, channel in traffic:
        context = calls.get(step, (None, None))[1]
        role = "sender" if direction == "emits" else "receiver"
        check(any(e == context and c == channel and r == role and m == messages
                  for e,c,r,m,_ in participants), "Traffic lacks context membership: " + step)

    starts = matching(pair.format("starts-at"))
    graph = collections.defaultdict(set)
    for start, end, condition in matching(
            r"(" + IDENT + r") proceeds-to (" + IDENT + r") when (" + IDENT + r")\."):
        graph[start].add(end)
        check(kinds.get(condition) == "condition", "Invalid branch condition: " + start)
    terminals = collections.defaultdict(set)
    for step, scenario, outcome in matching(
            r"(" + IDENT + r") finishes (" + IDENT + r") with (" + IDENT + r")\."):
        terminals[scenario].add(step)
        check(kinds.get(outcome) == "outcome", "Invalid terminal outcome: " + step)
    scenarios = set(n for n,k in kinds.items() if k == "scenario")
    features = set(n for n,k in kinds.items() if k == "feature")
    check(set(f for _,f in matching(pair.format("illustrates"))) == features,
          "Feature without a scenario, or an invalid illustrated feature")
    check(collections.Counter(s for s,_ in starts) == collections.Counter({s:1 for s in scenarios}),
          "Each scenario needs one entry")
    all_reached = set()
    for scenario, start in starts:
        reached, todo = set(), [start]
        while todo:
            step = todo.pop()
            if step in reached:
                continue
            reached.add(step)
            todo.extend(graph[step])
        all_reached.update(reached)
        check(terminals[scenario] <= reached, "Unreachable terminal: " + scenario)
        for step in reached:
            check(kinds.get(step) == "step", "Non-step in scenario: " + step)
            check(bool(graph[step]) or step in terminals[scenario], "Dead end: " + step)
            check(not (graph[step] and step in terminals[scenario]), "Terminal has successor: " + step)
        can_finish = set(terminals[scenario])
        while True:
            more = set(s for s in reached if graph[s] & can_finish)
            if more <= can_finish:
                break
            can_finish.update(more)
        check(reached <= can_finish, "A step cannot reach an outcome: " + scenario)
    check(all_reached == set(n for n,k in kinds.items() if k == "step"), "Orphan/mistyped steps")

    with (base / "requirement-coverage.csv").open(newline="") as stream:
        rows = list(csv.DictReader(stream))
    expected_ids = [i for s in manifest["sources"] for i in s["requirement_ids"]]
    check(collections.Counter(r["requirement_id"] for r in rows) == collections.Counter(expected_ids),
          "Coverage differs from pinned requirement inventory")
    check(set(r["model_id"] for r in rows) == set(n for n,k in kinds.items() if k == "requirement"),
          "Requirement declarations differ from coverage")
    constraints = set(matching(pair.format("constrains")))
    statuses = dict(matching(r"(" + IDENT + r") status ([a-z-]+)\."))
    for row in rows:
        name = row["model_id"]
        check(statuses.get(name) == row["disposition"], "Disposition mismatch: " + name)
        for target in row["design_targets"].split(";"):
            check((name,target) in constraints, "Missing obligation mapping: " + name + " -> " + target)
        check((name + " cites " + json.dumps(row["source"]) + ".") in seen,
              "Requirement source mismatch: " + name)

    if source_root:
        for source in manifest["sources"]:
            path = source_root / source["path"]
            check(path.is_file(), "Missing source: " + str(path))
            if not path.is_file():
                continue
            check(hashlib.sha256(path.read_bytes()).hexdigest() == source["sha256"],
                  "Source hash mismatch: " + source["path"])
            if source["requirement_ids"]:
                ids = re.findall(r"^(?:### |\- `)([A-Z0-9-]+-REQ-\d+)\b", path.read_text(), re.M)
                check(ids == source["requirement_ids"], "Source IDs changed: " + source["path"])

    with (base / "functionality-coverage.csv").open(newline="") as stream:
        function_rows = list(csv.DictReader(stream))
    check(collections.Counter(r["functionality"] for r in function_rows) ==
          collections.Counter({f:1 for f in functions}), "Function coverage inventory mismatch")
    for row in function_rows:
        function = row["functionality"]
        check(owners[function] == [row["owner"]], "Function coverage owner mismatch: " + function)
        actual_steps = set(s for s,(f,_) in calls.items() if f == function)
        recorded_steps = set(filter(None, row["scenario_steps"].split(";")))
        check(actual_steps == recorded_steps, "Function scenario coverage mismatch: " + function)
        check(locations.get(function, "").split(":")[0] == row["design_file"],
              "Function source location mismatch: " + function)

    return {
        "audit": "authoring-inventory-only",
        "result": "pass" if not errors else "fail",
        "files": len(paths), "declarations": len(kinds),
        "by_kind": dict(sorted(collections.Counter(kinds.values()).items())),
        "requirements": len(rows), "source_files": len(manifest["sources"]),
        "source_hashes_checked": bool(source_root),
        "traffic_facts": len(traffic), "participation_facts": len(participants),
        "functionalities_called_by_scenarios": len(set(f for f,_ in calls.values())),
        "errors": errors,
        "not_proven": ["SDL grammar/type conformance", "predicate coverage/exclusivity",
            "mode compatibility or deployment allocation", "payload/delivery correctness",
            "design sufficiency", "implementation correctness", "executable IR completeness"]
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--source-root", type=Path, help="Optional pinned Ponsse source snapshot")
    args = parser.parse_args()
    result = audit(Path(__file__).resolve().parent / "SDL" / "MVP1", args.source_root)
    print(json.dumps(result, indent=2, sort_keys=True))
    return 0 if result["result"] == "pass" else 1


if __name__ == "__main__":
    sys.exit(main())
