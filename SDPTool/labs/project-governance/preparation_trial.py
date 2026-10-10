#!/usr/bin/env python3
"""Exercise preparation with a separate synthetic project; never start a model."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import uuid


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("mode", choices=["create", "inspect"])
    parser.add_argument("--destination", type=Path, required=True)
    parser.add_argument("--controller", type=Path)
    parser.add_argument("--sdptool", type=Path)
    args = parser.parse_args()
    base = args.destination.resolve()
    if args.mode == "inspect":
        result = subprocess.check_output([str(base / "sdp-controller"), "assignment", "inspect",
            "--project", str(base / "project"), "--state", str(base / "runtime/state")])
        for row in json.loads(result)["assignments"]:
            work = row["assignment"]
            print(f'{work["role"]}: {work["id"]} — {row["readiness"]}')
            for condition in row["conditions"]:
                print("  " + json.dumps(condition, ensure_ascii=True))
        print("Preparation only. Design acceptance is not assessed; dispatch is unavailable.")
        return
    if not args.controller or not args.sdptool:
        parser.error("create requires --controller and --sdptool")
    base.mkdir(parents=True, exist_ok=False)
    root, state = base / "project", base / "runtime/state"
    binary = base / "sdp-controller"
    shutil.copy2(args.controller, binary)
    root.mkdir()

    def write(name, text):
        path = root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(text)

    write("AGENTS.md", "# Preparation fixture\nNo agent execution or delegation is authorized.\n")
    write("SDP/SDP-project.manifest.yaml", "schemaVersion: '1.0'\nproject:\n  name: pg-preparation-lab\n")
    write("SDP/authority.md", "# Fixture authority\nController may prepare and inspect these synthetic assignments. This is not approval to implement or launch.\n")
    ledger, selections = [], {}
    for index, role in enumerate(["SAD", "Master"], 1):
        name = role.lower()
        card, plan, session = f"KB-PGPREP-{index:03}", f"PLAN-PGPREP-{index:04}", f"SESSION-PGPREP-{index:04}"
        sel = {"card": f"SDP/KanBan/active/{name}.md", "plan": f"SDP/Plans/{name}.md",
               "session": f"SDP/Sessions/{name}.md", "authority": "SDP/authority.md", "step": "S1"}
        selections[role] = sel
        write(sel["card"], f"# Fixture work\n\n| Field | Value |\n| --- | --- |\n| id | {card} |\n| PlanId | {plan} |\n| CardState | in-progress |\n")
        kind = "DesignPlan" if role == "SAD" else "ImplementationPlan"
        write(sel["plan"], f"# Fixture plan\n\n| Field | Value |\n| --- | --- |\n| id | {plan} |\n| PlanType | {kind} |\n| state | active |\n")
        write(sel["session"], f"# Fixture goal\n\n| State | Step |\n| --- | --- |\n| active | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | {session} |\n| Status | active |\n")
        for subject, event_type, path in [(card, "x-kanban:created", f"active/{name}.md"), (plan, "x-management:created", f"Plans/{name}.md")]:
            ledger.append({"schemaVersion": "1.0", "eventId": f"FIXTURE-{len(ledger)+1}", "eventType": event_type, "subjectId": subject,
                           "payload": {"schemaVersion": "0.2", "previousEventId": None, "to": "active", "toPath": path}})
    write("SDP/ProjectManagement/Ledger.ndjson", "".join(json.dumps(e) + "\n" for e in ledger))
    for kind in ["now", "target", "task", "blueprint", "design-disposition"]:
        write(f"artifacts/{kind}.json", json.dumps({"fixture": kind, "approval": "not-assessed"}) + "\n")
    subprocess.run([str(args.sdptool.resolve()), str(root), "governance", "bind", "--state", str(state), "--json"], check=True, capture_output=True)
    cases = []

    def call(action, request=None, expected=None):
        command = [str(binary), "assignment", action, "--project", str(root), "--state", str(state)]
        if request is not None:
            path = base / f"request-{len(cases)}.json"
            path.write_text(json.dumps(request, indent=2) + "\n")
            command += ["--input", str(path)]
        process = subprocess.run(command, capture_output=True, text=True)
        if expected:
            assert process.returncode != 0 and expected in process.stderr, process
            cases.append({"case": action, "expectedFailure": expected, "stderr": process.stderr})
            return None
        assert process.returncode == 0, process.stderr
        value = json.loads(process.stdout)
        cases.append({"case": action, "request": request, "response": value})
        return value

    requests, results = {}, {}
    for role in ["SAD", "Master"]:
        context = call("context", {"role": role, "selection": selections[role]})
        assert context["context"]["conditions"] == [], context
        kinds = ["now"] if role == "SAD" else ["now", "target", "task", "blueprint", "design-disposition"]
        pins = [{"kind": kind, "path": f"artifacts/{kind}.json", "sha256": hashlib.sha256((root / f"artifacts/{kind}.json").read_bytes()).hexdigest()} for kind in kinds]
        q = {"schema": "sdp-governance-assignment/1", "projectRef": context["context"]["projectRef"], "checkoutId": context["context"]["checkoutId"],
             "operationId": str(uuid.uuid4()), "expectedGeneration": context["generation"], "inputRevision": context["context"]["revision"],
             "role": role, "selection": selections[role], "assignment": {"outcome": f"Prepare bounded {role} fixture", "writablePaths": ["model-work" if role == "SAD" else "code"], "excludedPaths": ["SDP", "artifacts"]},
             "artifacts": pins, "returnCriteria": ["Return exact candidate evidence and unresolved questions"]}
        requests[role] = q
        results[role] = call("prepare", q)
        assert call("prepare", q) == results[role], "retry allocated a new assignment"
    before = call("inspect")
    assert len(before["assignments"]) == 2 and all(v["readiness"] == "prepared" for v in before["assignments"])
    conflict = dict(requests["Master"], operationId=str(uuid.uuid4()), expectedGeneration=before["generation"])
    conflict["assignment"] = dict(conflict["assignment"], writablePaths=["model-work/subdir"])
    call("prepare", conflict, "conflict:")
    write("artifacts/target.json", '{"fixture":"changed TARGET"}\n')
    after = call("inspect")
    assert {v["assignment"]["role"]: v["readiness"] for v in after["assignments"]} == {"SAD": "prepared", "Master": "stale"}, after
    assert call("prepare", requests["Master"]) == results["Master"], "stale retry lost original identity"
    stale = dict(requests["Master"], operationId=str(uuid.uuid4()), expectedGeneration=after["generation"])
    call("prepare", stale, "stale_revision:")
    # The pre-preparation tool must refuse the upgraded store without rewriting it.
    head = (state / "HEAD").read_bytes()
    old = subprocess.run([str(args.sdptool.resolve()), str(root), "governance", "bind", "--state", str(state), "--json"], capture_output=True, text=True)
    assert old.returncode != 0 and (state / "HEAD").read_bytes() == head, "old reader silently accepted/rewrote 0.2"
    cases.append({"case": "old-binary-refuses-0.2", "stderr": old.stderr})
    (base / "evidence.json").write_text(json.dumps({"fixtureOnly": True, "modelSubmitted": False, "binarySHA256": hashlib.sha256(binary.read_bytes()).hexdigest(), "cases": cases}, indent=2) + "\n")
    (base / "README.md").write_text("# Assignment preparation trial\n\nSynthetic fixture; no model runs. Two assignments were prepared, then TARGET was deliberately changed. SAD remains prepared; Master is stale. No design acceptance or dispatch is granted. The original owner lab is untouched. evidence.json contains actual subprocess requests/results.\n")
    print("PASS: two assignments, retry/restart, conflict refusal, changed TARGET, old-reader refusal. No model started.")


if __name__ == "__main__":
    main()
