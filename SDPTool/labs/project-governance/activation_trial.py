#!/usr/bin/env python3
"""Activate an approved fixture into the real PG controller; host kind is explicit."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import sys
import uuid


def main():
    p = argparse.ArgumentParser(description=__doc__)
    for name in ['destination', 'controller', 'sdptool', 'mcp', 'source', 'host', 'old-controller']:
        p.add_argument('--' + name, type=Path, required=True)
    p.add_argument('--host-kind', choices=['protocol-fixture', 'native'], required=True)
    p.add_argument("--reporting", choices=["manual", "automatic"], default="manual")
    a = p.parse_args()
    base = a.destination.resolve()
    subprocess.run([sys.executable, str(Path(__file__).with_name('design_trial.py')), '--destination', str(base),
                    '--controller', str(a.controller.resolve()), '--sdptool', str(a.sdptool.resolve()),
                    '--source', str(a.source.resolve()), '--execution-fixture'], check=True)
    root, state = base / 'project', base / 'runtime/state'
    controller = base / 'sdp-controller'
    target = root / 'SDP/Models/WORK--After/System.design'
    original = (a.source / 'System.design').read_bytes()
    target.write_bytes(original)
    evidence = json.loads((base / 'evidence.json').read_text())
    evidence.update(hostKind=a.host_kind, actualOwnerApprovalRecorded=False, modelTaskSubmitted=False)
    cases = evidence['cases']
    check = json.loads((base / 'design-check.json').read_text())

    def inventory():
        return {str(p.relative_to(root)): hashlib.sha256(p.read_bytes()).hexdigest()
                for p in root.rglob('*') if p.is_file()}

    def run(action, request=None, stdin='', success=True, owner=None, binary=None):
        command = [str(binary or controller), 'assignment', action, '--project', str(root), '--state', str(state)]
        if request is not None:
            f = base / ('activation-request-' + str(uuid.uuid4()) + '.json')
            f.write_text(json.dumps(request, indent=2) + '\n')
            command += ['--input', str(f)]
        if owner:
            command += ['--owner', owner]
        r = subprocess.run(command, input=stdin, text=True, capture_output=True, timeout=60)
        assert (r.returncode == 0) == success, r.stderr
        result = json.loads(r.stdout) if success else None
        cases.append({'operation': action, 'request': request, 'result': result, 'exitCode': r.returncode, 'stderr': r.stderr})
        return result

    basis = run('design-check', check)
    binding = json.loads((state / 'generations' / ((state / 'HEAD').read_text() + '.json')).read_text())
    decide = {'schema': 'sdp-governance-design-decision/1', 'projectRef': binding['project'], 'checkoutId': binding['checkoutId'],
              'operationId': str(uuid.uuid4()), 'review': check, 'reviewRevision': basis['reviewRevision']}
    decision = run('design-decide', decide, 'approved\nScripted fixture authority for the isolated PG runtime test only.\n', owner='fixture-owner')
    check['expectedGeneration'] = (state / 'HEAD').read_text()
    activate = {'schema': 'sdp-governance-activation/1', 'projectRef': binding['project'], 'checkoutId': binding['checkoutId'],
                'operationId': str(uuid.uuid4()), 'requestId': str(uuid.uuid4()), 'review': dict(check),
                'reviewRevision': basis['reviewRevision'], 'decisionOperationId': decision['operationId']}
    result = run('activate', activate)
    assert result['run']['routine']['id'] == 'implementation'
    assert run('activate', activate) == result
    other = dict(activate, operationId=str(uuid.uuid4()))
    run('activate', other, success=False)
    before_head = (state / 'HEAD').read_text()
    run('inspect', binary=a.old_controller.resolve(), success=False)
    assert (state / 'HEAD').read_text() == before_head
    run_id = result['run']['id']
    message = base / 'message.txt'
    message.write_text('This is an owner-authorized isolated PG integration test. Use the approved fixture package. '
                       'Fix Unicode whitespace trimming in code/normalize.go and add focused tests in code/normalize_test.go. '
                       'Preserve interior spaces, letter case, code/go.mod, unrelated.txt, all models and all SDP records. '
                       'Use the assignment-bound context_read and run_read tools, then report via run_transition if successful. '
                       'Do not spawn agents, use network, install anything, change host settings, commit or publish. '
                       'No actual owner project approval is represented by the scripted fixture decision.\n')
    command = [str(controller), 'run', '--project', str(root), '--state', str(state), '--run', run_id,
               '--codex', str(a.host.resolve()), '--mcp-binary', str(a.mcp.resolve()), '--message', str(message),
               '--sandbox', 'workspace-write' if a.host_kind == 'native' else 'read-only', '--approvals', 'never', '--reporting', a.reporting, '--timeout', '5m']
    env = dict(os.environ)
    if a.host_kind == 'protocol-fixture':
        env['SDP_PGI3_FAKE_HOST'] = '1'
    # Stale design must refuse before any model turn. With this protocol fixture
    # even thread creation is excluded by the durable input claim guard.
    target.write_bytes(original + b'\n// stale before activation launch\n')
    rejected = subprocess.run(command, input='', text=True, capture_output=True, env=env, timeout=60)
    assert rejected.returncode != 0 and 'stale_revision' in rejected.stderr, rejected.stderr
    target.write_bytes(original)
    evidence['status'] = 'preparing-host-trial'
    (base / 'evidence.json').write_text(json.dumps(evidence, indent=2) + '\n')
    before = inventory()
    golden = base / 'verify'
    if a.host_kind == 'native':
        golden.mkdir()
        (golden / 'go.mod').write_text('module pgverify\n\ngo 1.26.0\nrequire pgfixture v0.0.0\nreplace pgfixture => ' + json.dumps(str(root / 'code')) + '\n')
        (golden / 'main.go').write_text('package main\nimport ("fmt"; "os"; fixture "pgfixture")\nfunc main() { for _, c := range [][2]string{{"\\t  Keep  Case\\u00a0", "Keep  Case"}, {"A\\tB", "A\\tB"}, {"\\u2003Mixed CASE\\n", "Mixed CASE"}, {"\\t\\n", ""}} { if got := fixture.Normalize(c[0]); got != c[1] { fmt.Printf("FAIL: %q -> %q; want %q\\n", c[0], got, c[1]); os.Exit(1) } }; fmt.Println("PASS: independent Unicode and interior-preservation cases") }\n')
        initial = subprocess.run(['go', 'run', '.'], cwd=golden, text=True, capture_output=True, timeout=60)
        assert initial.returncode != 0 and 'FAIL:' in initial.stdout, (initial.stdout, initial.stderr)
        evidence['goldenBefore'] = initial.stdout
    stdout_path, stderr_path = base / 'controller.ndjson', base / 'controller.err'
    print('Starting ' + a.host_kind + ' controller trial in ' + str(base), flush=True)
    with stdout_path.open('w') as stdout, stderr_path.open('w') as stderr:
        process = subprocess.run(command, input='', text=True, stdout=stdout, stderr=stderr, env=env, timeout=330)
    observed = [json.loads(line) for line in stdout_path.read_text().splitlines() if line.strip()]
    evidence.update(status='verification-incomplete', controllerExitCode=process.returncode,
                    modelTaskSubmitted=a.host_kind == 'native' and any(v.get('observation', {}).get('kind') == 'turn.started' for v in observed))
    (base / 'evidence.json').write_text(json.dumps(evidence, indent=2) + '\n')
    assert process.returncode == 0, stderr_path.read_text()
    status_cmd = [str(controller), 'status', '--project', str(root), '--state', str(state), '--run', run_id]
    final = json.loads(subprocess.check_output(status_cmd, text=True))
    assert final['run']['routine']['id'] == 'implementation' and final['coverage'] == 'known'
    assert final['run']['state'] != 'completed', 'host completion must not complete governance work'
    assert final['run']['attempts'][-1]['status'] == 'completed', final
    if a.reporting == 'automatic':
        inbox = json.loads(subprocess.check_output([str(controller), 'inbox', '--project', str(root),
                                                   '--state', str(state), '--run', run_id, '--json'], text=True))
        (base / 'inbox.json').write_text(json.dumps(inbox, indent=2) + '\n')
        evidence['inbox'] = inbox
        evidence['reportingMode'] = a.reporting
        (base / 'evidence.json').write_text(json.dumps(evidence, indent=2) + '\n')
        assert final['run']['state'] == 'verifying', final
        assert len(inbox['reports']) == 1 and inbox['reports'][0]['runId'] == run_id, inbox
        assert not any(v.get('observation', {}).get('kind') == 'approval.requested' for v in observed), 'unexpected approval prompt'
    # A second launch cannot send another turn into the same completed attempt.
    retry = subprocess.run(command, input='', text=True, capture_output=True, env=env, timeout=60)
    assert retry.returncode != 0, 'duplicate launch accepted'
    after = inventory()
    changed = sorted(k for k in set(before) | set(after) if before.get(k) != after.get(k))
    assert set(changed) <= {'code/normalize.go', 'code/normalize_test.go'}, changed
    if a.host_kind == 'native':
        assert changed, 'native fixture did no work'
        tested = subprocess.run(['go', 'test', './...'], cwd=root / 'code', text=True, capture_output=True, timeout=60)
        assert tested.returncode == 0, tested.stderr
        evidence['goTestOutput'] = tested.stdout
        verified = subprocess.run(['go', 'run', '.'], cwd=golden, text=True, capture_output=True, timeout=60)
        assert verified.returncode == 0, (verified.stdout, verified.stderr)
        evidence['goldenAfter'] = verified.stdout
    else:
        assert not changed, changed
    evidence.update(status='passed', modelTaskSubmitted=a.host_kind == 'native', activationRequest=activate, activationResult=result,
                    finalRun=final, changedFiles=changed, protectedFilesPreserved=True,
                    controllerLogSHA256=hashlib.sha256(stdout_path.read_bytes()).hexdigest(),
                    hostBinarySHA256=hashlib.sha256(a.host.resolve().read_bytes()).hexdigest(),
                    mcpBinarySHA256=hashlib.sha256(a.mcp.resolve().read_bytes()).hexdigest())
    (base / 'activation.json').write_text(json.dumps(activate, indent=2) + '\n')
    (base / 'evidence.json').write_text(json.dumps(evidence, indent=2) + '\n')
    print('PASS: approved implementation activation through ' + a.host_kind + '; stale and duplicate launch refused; protected files preserved.')


if __name__ == '__main__':
    main()
