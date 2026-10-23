#!/usr/bin/env python3
"""Run a separate compiled CLI revision/retirement trial; no host or model calls."""
import argparse
import hashlib
import json
from pathlib import Path
import subprocess
import sys
import uuid


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--destination", type=Path, required=True)
    parser.add_argument("--controller", type=Path, required=True)
    parser.add_argument("--old-controller", type=Path, required=True)
    parser.add_argument("--sdptool", type=Path, required=True,
                        help="pre-C1 binder, for the existing preparation trial")
    args = parser.parse_args()
    base = args.destination.resolve()
    subprocess.run([sys.executable, str(Path(__file__).with_name("preparation_trial.py")),
                    "create", "--destination", str(base), "--controller", str(args.controller.resolve()),
                    "--sdptool", str(args.sdptool.resolve())], check=True)
    root, state = base / "project", base / "runtime/state"
    binary = base / "sdp-controller"
    evidence = json.loads((base / "evidence.json").read_text())
    cases = evidence["cases"]
    master_case = next(c for c in cases if c["case"] == "prepare" and
                       c.get("request", {}).get("role") == "Master")
    original, first = master_case["request"], master_case["response"]

    def call(action, request=None, expected=None):
        cmd = [str(binary), "assignment", action, "--project", str(root), "--state", str(state)]
        if request is not None:
            path = base / f"request-{len(cases)}.json"
            path.write_text(json.dumps(request, indent=2) + "\n")
            cmd += ["--input", str(path)]
        before = (state / "HEAD").read_bytes()
        process = subprocess.run(cmd, capture_output=True, text=True)
        record = {"case": action, "request": request}
        if expected:
            assert process.returncode != 0 and expected in process.stderr, process
            assert (state / "HEAD").read_bytes() == before, "refusal changed HEAD"
            record.update(expectedFailure=expected, stderr=process.stderr)
            cases.append(record)
            return None
        assert process.returncode == 0, process.stderr
        result = json.loads(process.stdout)
        record["response"] = result
        cases.append(record)
        return result

    def change(p):
        view = call("inspect")
        return {"schema": "sdp-governance-assignment/2", "projectRef": original["projectRef"],
                "checkoutId": view["checkoutId"], "operationId": str(uuid.uuid4()),
                "expectedGeneration": view["generation"], "assignmentId": p["id"],
                "expectedRevision": p["revision"], "reason": "Fixture owner revises the work"}

    context = call("context", {"role": "Master", "selection": original["selection"]})
    revision = change(first["assignment"])
    revision.update(inputRevision=context["context"]["revision"],
                    assignment=dict(original["assignment"], writablePaths=["code-v2"]),
                    artifacts=original["artifacts"], returnCriteria=original["returnCriteria"])
    call("revise", revision, "stale_revision:")
    revision["artifacts"] = [dict(pin, sha256=hashlib.sha256((root / pin["path"]).read_bytes()).hexdigest())
                             for pin in original["artifacts"]]
    second = call("revise", revision)
    assert second["assignment"]["id"] == first["assignment"]["id"]
    assert second["assignment"]["revision"] == 2 and second["assignment"]["state"] == "prepared"
    assert call("revise", revision) == second
    before = (state / "HEAD").read_bytes()
    old = subprocess.run([str(args.old_controller.resolve()), "assignment", "inspect", "--project",
                          str(root), "--state", str(state)], capture_output=True, text=True)
    assert old.returncode != 0 and (state / "HEAD").read_bytes() == before
    cases.append({"case": "C1-reader-refuses-0.3", "stderr": old.stderr,
                  "binarySHA256": hashlib.sha256(args.old_controller.read_bytes()).hexdigest()})
    conflict = dict(revision, **change(second["assignment"]))
    conflict["assignment"] = dict(revision["assignment"], writablePaths=["model-work/nested"])
    call("revise", conflict, "conflict:")
    retired_request = change(second["assignment"])
    retired_request["reason"] = "Fixture owner no longer wants this prepared work"
    target = root / "artifacts/target.json"
    target_bytes = target.read_bytes()
    target.unlink()
    third = call("retire", retired_request)
    assert third["assignment"]["revision"] == 3 and third["assignment"]["state"] == "retired"
    assert call("retire", retired_request) == third
    assert call("revise", revision) == second
    assert call("prepare", original) == first
    terminal = dict(revision, **change(third["assignment"]))
    call("revise", terminal, "invalid_transition:")
    view = call("inspect")
    assert next(v for v in view["assignments"] if v["assignment"]["id"] == first["assignment"]["id"])["readiness"] == "retired"
    target.write_bytes(target_bytes)
    new = dict(original, operationId=str(uuid.uuid4()), expectedGeneration=view["generation"],
               assignment=revision["assignment"], artifacts=revision["artifacts"])
    fresh = call("prepare", new)
    assert fresh["assignment"]["id"] != first["assignment"]["id"]
    final = call("inspect")
    assert final["schema"] == "sdp-governance-assignment/2", "prepare downgraded profile"
    assert len(final["assignments"]) == 3
    assert all(v["dispatch"] == "unavailable" and v["designAcceptance"] == "not-assessed"
               for v in final["assignments"])
    evidence["lifecycle"] = {"originalID": first["assignment"]["id"], "replacementID": fresh["assignment"]["id"],
                             "retiredRevision": 3, "profile": "sdp-project-governance-state/0.3"}
    (base / "evidence.json").write_text(json.dumps(evidence, indent=2) + "\n")
    (base / "README.md").write_text("# Assignment lifecycle trial\n\nSynthetic fixture; no model runs. "
        "A stale Master was revised with refreshed pins, then retired despite a missing TARGET. "
        "A new Master claims its released scope. Historical retries retain their original results. "
        "The C1 reader refuses the upgraded store. Dispatch remains unavailable; design acceptance "
        "is not assessed. evidence.json records the actual subprocess requests/results.\n")
    print("PASS: revise, retire, immutable retries, scope release, terminal/conflict refusal and C1 reader refusal.")


if __name__ == "__main__":
    main()
