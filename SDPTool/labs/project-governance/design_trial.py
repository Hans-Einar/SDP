#!/usr/bin/env python3
"""Build real model/Blueprint fixture records and check them through the PG CLI."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import uuid


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--destination", type=Path, required=True)
    parser.add_argument("--controller", type=Path, required=True)
    parser.add_argument("--sdptool", type=Path, required=True)
    parser.add_argument("--source", type=Path, required=True, help="SDL fixture directory with System.design")
    parser.add_argument("--execution-fixture", action="store_true", help="authorize only a bounded isolated runtime test task")
    args = parser.parse_args()
    base = args.destination.resolve()
    base.mkdir(parents=True, exist_ok=False)
    root, state = base / "project", base / "runtime/state"
    controller, tool = base / "sdp-controller", base / "sdptool"
    shutil.copy2(args.controller.resolve(), controller)
    shutil.copy2(args.sdptool.resolve(), tool)
    cases = []

    def write(path, content):
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(content)

    def request_file(value):
        path = base / f"request-{len(cases)}-{uuid.uuid4()}.json"
        path.write_text(json.dumps(value, indent=2) + "\n")
        return path

    def run(command, operation, request=None):
        process = subprocess.run([str(x) for x in command], capture_output=True, text=True)
        assert process.returncode == 0, process.stderr
        result = json.loads(process.stdout)
        cases.append({"operation": operation, "request": request, "result": result})
        return result

    def pg(action, request=None):
        command = [controller, "assignment", action, "--project", root, "--state", state]
        if request is not None:
            command += ["--input", request_file(request)]
        return run(command, action, request)

    def inventory():
        return {str(p.relative_to(base)): hashlib.sha256(p.read_bytes()).hexdigest()
                for directory in [root, state] for p in sorted(directory.rglob("*")) if p.is_file()}

    selection = {"card": "SDP/KanBan/active/master.md", "plan": "SDP/Plans/master.md",
                 "session": "SDP/Sessions/master.md", "authority": "SDP/authority.md", "step": "S1"}
    if args.execution_fixture:
        write(root / "AGENTS.md", "# Isolated PG runtime fixture\nThe owner authorizes testing PG execution in this isolated fixture. Only change code/normalize.go and code/normalize_test.go. Preserve all other files, models, governance records and existing unrelated content. Do not spawn agents, install tools, use network, commit, publish or modify the host configuration. Report through the existing assignment-bound MCP tools.\n")
        write(root / "code/go.mod", "module pgfixture\n\ngo 1.26.0\n")
        write(root / "code/normalize.go", 'package pgfixture\nimport "strings"\nfunc Normalize(s string) string { return strings.Trim(s, " ") }\n')
        write(root / "code/normalize_test.go", 'package pgfixture\nimport "testing"\nfunc TestExisting(t *testing.T) { if Normalize("  Keep  Case  ") != "Keep  Case" { t.Fatal("wrong normalization") } }\n')
        write(root / "unrelated.txt", "Preserve this original fixture content.\n")
    else:
        write(root / "AGENTS.md", "# Design check fixture\nNo model task or implementation is authorized.\n")
    write(root / "SDP/SDP-project.manifest.yaml", "schemaVersion: '1.0'\nproject:\n  name: pg-design-lab\n")
    write(root / selection["authority"], "# Fixture authority\n" + ("Session0007 owner continuation authorizes a bounded isolated PG execution test. The test controller may create a clearly attributed fixture decision and launch one Master to fix Unicode edge whitespace normalization and add focused tests in code. No actual owner acceptance of project design is inferred.\n" if args.execution_fixture else "Prepare and inspect a bounded design package. No owner design acceptance is recorded.\n"))
    write(root / selection["card"], "# Work\n\n| Field | Value |\n| --- | --- |\n| id | KB-PGDES-001 |\n| PlanId | PLAN-PGDES-0001 |\n| CardState | in-progress |\n")
    write(root / selection["plan"], "# Plan\n\n| Field | Value |\n| --- | --- |\n| id | PLAN-PGDES-0001 |\n| PlanType | ImplementationPlan |\n| state | active |\n")
    write(root / selection["session"], "# Goal\n\n| State | Step |\n| --- | --- |\n| active | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-PGDES-0001 |\n| Status | active |\n")
    ledger = []
    for index, (subject, event_type, path) in enumerate([
            ("KB-PGDES-001", "x-kanban:created", "active/master.md"),
            ("PLAN-PGDES-0001", "x-management:created", "Plans/master.md")], 1):
        payload = {"schemaVersion": "0.2", "previousEventId": None, "from": None, "to": "active", "fromPath": None,
                   "toPath": path, "reason": "Synthetic design check fixture", "links": []}
        if index == 2:
            payload.update(kind="Plan", planType="ImplementationPlan")
        ledger.append({"schemaVersion": "1.0", "eventId": f"EVT-PGDES-{index}", "eventType": event_type,
                       "occurredAt": "2026-10-18T00:00:00Z", "actor": "fixture-controller", "commit": None,
                       "subjectId": subject, "payload": payload})
    write(root / "SDP/ProjectManagement/Ledger.ndjson", "".join(json.dumps(v) + "\n" for v in ledger))
    models = root / "SDP/Models"
    models.mkdir()
    artifacts = []
    for name in ["Before", "After"]:
        result = run([tool, models, "model", "create", "work:" + name, "--initial", "--json"], "create-model")
        path = Path(result["path"])
        shutil.copytree(args.source.resolve(), path, dirs_exist_ok=True)
        artifacts.append(path)
    task = {"ID": "design-fixture", "Intent": ("Fix Unicode surrounding whitespace in code/normalize.go; add regression tests preserving interior spacing and case. Preserve code/go.mod and unrelated content." if args.execution_fixture else "Inspect unchanged bounded model"), "AllowedCodePaths": ["code"]}
    task_path = root / "task.json"
    write(task_path, json.dumps(task) + "\n")
    bundle = run([tool, models, "model", "create", "blueprint", "from", "work:Before", "to", "work:After",
                  "--entry", "System.design", "--task", task_path, "--catalogue", root / "SDP/Blueprints", "--json"], "create-blueprint")
    bundle_path = Path(bundle["path"])
    doc = json.loads((bundle_path / "blueprint.json").read_text())
    binding = {"bundle": str(bundle_path.relative_to(root / "SDP")), "task": doc["taskId"], "system": doc["analysis"]["System"],
               "revision": bundle["revision"], "retained": bundle["retainedRevision"], "plan": "PLAN-PGDES-0001",
               "milestone": "M1", "modelArea": "Models", "from": "work:Before", "to": "work:After"}
    create = {"schema": "sdp-blueprint-assignment/1", "eventId": "EVT-PGDES-BPA-1", "occurredAt": "2026-10-18T00:00:00Z",
              "assignmentId": "BPA-design-fixture", "action": "create", "expectedEvent": "", "binding": binding,
              "reason": "Synthetic provider fixture, no owner approval"}
    run([tool, root / "SDP", "model", "assignment", "apply", "--request", request_file(create), "--as",
         "fixture-controller", "--authority", "controller", "--json"], "create-provider-assignment", create)
    run([tool, root, "governance", "bind", "--state", state, "--json"], "bind")
    context = pg("context", {"role": "Master", "selection": selection})
    write(root / "design-disposition.json", '{"fixture":"No owner decision recorded"}\n')
    pins = []
    paths = [artifacts[0] / (artifacts[0].name + ".yaml"), artifacts[1] / (artifacts[1].name + ".yaml"),
             task_path, bundle_path / "blueprint.json", root / "design-disposition.json"]
    for kind, path in zip(["now", "target", "task", "blueprint", "design-disposition"], paths):
        pins.append({"kind": kind, "path": str(path.relative_to(root)), "sha256": hashlib.sha256(path.read_bytes()).hexdigest()})
    prepare = {"schema": "sdp-governance-assignment/1", "projectRef": context["context"]["projectRef"],
               "checkoutId": context["context"]["checkoutId"], "operationId": str(uuid.uuid4()),
               "expectedGeneration": context["generation"], "inputRevision": context["context"]["revision"],
               "role": "Master", "selection": selection, "assignment": {"outcome": task["Intent"], "writablePaths": ["code"],
               "excludedPaths": ["SDP"]}, "artifacts": pins, "returnCriteria": ["Return exact design review basis and unknowns"]}
    prepared = pg("prepare", prepare)
    view = pg("inspect")
    check = {"schema": "sdp-governance-design-check/1", "assignmentId": prepared["assignment"]["id"],
             "expectedRevision": 1, "expectedGeneration": view["generation"], "blueprintAssignment": create["assignmentId"],
             "blueprintEvent": create["eventId"]}
    before = inventory()
    current = pg("design-check", check)
    assert current["status"] == "current" and current["reviewRevision"], current
    assert current["designAcceptance"] == "not-recorded" and current["dispatch"] == "unavailable"
    assert current == pg("design-check", check) and before == inventory()
    target = artifacts[1] / "System.design"
    original = target.read_bytes()
    target.write_bytes(original + b"\n// uncommitted TARGET changed after review\n")
    before = inventory()
    simple = pg("inspect")
    assert simple["assignments"][0]["readiness"] == "prepared", "metadata should be unchanged"
    blocked = pg("design-check", check)
    assert blocked["status"] == "blocked" and "reviewRevision" not in blocked, blocked
    assert before == inventory(), "read-only checks mutated files"
    target.write_bytes(original)
    assert pg("design-check", check) == current, "restored exact sources lost deterministic basis"
    target.write_bytes(original + b"\n// uncommitted TARGET changed after review\n")
    write(base / "design-check.json", json.dumps(check, indent=2) + "\n")
    evidence = {"fixtureOnly": True, "modelTaskSubmitted": False, "ownerApprovalRecorded": False,
                "controllerSHA256": hashlib.sha256(controller.read_bytes()).hexdigest(),
                "sdptoolSHA256": hashlib.sha256(tool.read_bytes()).hexdigest(), "cases": cases}
    write(base / "evidence.json", json.dumps(evidence, indent=2) + "\n")
    write(base / "README.md", "# Provider-backed design trial\n\nReal ModelGovernance artifacts and retained Blueprint, with a synthetic assignment. "
          "No agent task or owner approval. TARGET is deliberately changed at the end; file-pin inspection still says prepared, "
          "while design-check reports blocked. Exact source restoration restores the review digest. No read-only check writes state.\n")
    print("PASS: real provider/CLI journey, stable review digest, uncommitted TARGET drift, no writes or approval.")


if __name__ == "__main__":
    main()
