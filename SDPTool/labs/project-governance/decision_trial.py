#!/usr/bin/env python3
"""Exercise controller-attributed decisions with scripted fixture input, never owner approval."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import sys
import uuid


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    for name in ['destination', 'controller', 'sdptool', 'source', 'old-controller']:
        parser.add_argument('--' + name, type=Path, required=True)
    args = parser.parse_args()
    base = args.destination.resolve()
    subprocess.run([sys.executable, str(Path(__file__).with_name('design_trial.py')),
                    '--destination', str(base), '--controller', str(args.controller.resolve()),
                    '--sdptool', str(args.sdptool.resolve()), '--source', str(args.source.resolve())], check=True)
    root, state = base / 'project', base / 'runtime/state'
    controller = base / 'sdp-controller'
    shutil.copy2(args.old_controller.resolve(), base / 'old-controller')
    target = root / 'SDP/Models/WORK--After/System.design'
    original = (args.source / 'System.design').read_bytes()
    target.write_bytes(original)
    evidence = json.loads((base / 'evidence.json').read_text())
    evidence['fixtureOwnerInput'] = True
    evidence['ownerApprovalRecorded'] = False
    cases = evidence['cases']
    check = json.loads((base / 'design-check.json').read_text())

    def inventory():
        return {str(p.relative_to(base)): hashlib.sha256(p.read_bytes()).hexdigest()
                for d in [root, state] for p in d.rglob('*') if p.is_file()}

    def command(action, request=None, owner=None, old=False):
        cmd = [str(base / 'old-controller' if old else controller), 'assignment', action,
               '--project', str(root), '--state', str(state)]
        if request is not None:
            path = base / ('decision-request-' + str(uuid.uuid4()) + '.json')
            path.write_text(json.dumps(request, indent=2) + '\n')
            cmd += ['--input', str(path)]
        if owner:
            cmd += ['--owner', owner]
        return cmd

    def run(action, request=None, stdin='', success=True, owner=None, old=False):
        done = subprocess.run(command(action, request, owner, old), input=stdin, text=True,
                              capture_output=True, timeout=30)
        assert (done.returncode == 0) == success, (done.stdout, done.stderr)
        value = json.loads(done.stdout) if success else None
        cases.append({'operation': action, 'request': request, 'scriptedInput': stdin, 'fixtureOwner': owner,
                      'oldReader': old, 'result': value, 'exitCode': done.returncode, 'stderr': done.stderr})
        return value

    def fresh_check():
        check['expectedGeneration'] = (state / 'HEAD').read_text()
        return dict(check)

    def decision_request():
        review = fresh_check()
        basis = run('design-check', review)
        assert basis['status'] == 'current'
        binding = json.loads((state / 'generations' / ((state / 'HEAD').read_text() + '.json')).read_text())
        return {'schema': 'sdp-governance-design-decision/1', 'projectRef': binding['project'],
                'checkoutId': binding['checkoutId'], 'operationId': str(uuid.uuid4()),
                'review': review, 'reviewRevision': basis['reviewRevision']}

    q = decision_request()
    # No choice, cancellation or invalid input must never register a decision.
    for stdin in ['', '\n', 'cancel\n', 'yes\n', 'approved\n', 'approved\n\n']:
        before = inventory()
        run('design-decide', q, stdin, success=False, owner='fixture-owner')
        assert inventory() == before
    first = run('design-decide', q, 'approved\nSynthetic fixture approval; no real owner decision.\n', owner='fixture-owner')
    assert first['record']['input']['source'] == 'local-controller-stdin'
    # Lost-response retry needs no new stdin and returns the historical receipt.
    before = inventory()
    assert run('design-decide', q, owner='fixture-owner') == first
    assert inventory() == before
    run('design-decide', q, success=False, owner='different-fixture-owner')
    review = run('design-review', fresh_check())
    assert review['acceptance'] == 'approved' and review['dispatch'] == 'unavailable'
    before = inventory()
    run('inspect', success=False, old=True)
    assert inventory() == before, 'old reader mutated 0.4 store'

    target.write_bytes(original + b'\n// changed after scripted approval\n')
    before = inventory()
    review = run('design-review', fresh_check())
    assert review['acceptance'] == 'stale' and review['decision'] == first
    assert run('design-decide', q, owner='fixture-owner') == first
    assert inventory() == before
    target.write_bytes(original)

    # Change a source while the interactive adapter is waiting for its decision.
    pending = decision_request()
    process = subprocess.Popen(command('design-decide', pending, 'fixture-owner'), stdin=subprocess.PIPE,
                               stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    prompt = ''
    while True:
        line = process.stderr.readline()
        assert line, 'controller ended before owner prompt'
        prompt += line
        if 'Enter approved,' in line:
            break
    target.write_bytes(original + b'\n// changed during owner input\n')
    before = inventory()
    stdout, stderr = process.communicate('approved\nSynthetic stale decision attempt.\n', timeout=30)
    assert process.returncode != 0 and 'stale_revision' in stderr and inventory() == before
    cases.append({'operation': 'design-decide-drift-during-input', 'request': pending,
                  'scriptedInput': 'approved\nSynthetic stale decision attempt.\n', 'fixtureOwner': 'fixture-owner',
                  'result': None, 'exitCode': process.returncode, 'stderr': prompt + stderr})
    target.write_bytes(original)

    for choice in ['changes-requested', 'rejected']:
        next_request = decision_request()
        run('design-decide', next_request, choice + '\nSynthetic superseding fixture decision.\n', owner='fixture-owner')
        view = run('design-review', fresh_check())
        assert view['acceptance'] == choice and view['decision']['record']['input']['decision'] == choice
    before = inventory()
    assert run('design-decide', q, owner='fixture-owner') == first
    assert run('design-review', fresh_check())['acceptance'] == 'rejected'
    assert inventory() == before

    # Leave a new request ready for an optional manual fixture decision.
    ready = decision_request()
    (base / 'design-decision.json').write_text(json.dumps(ready, indent=2) + '\n')
    (base / 'design-review.json').write_text(json.dumps(fresh_check(), indent=2) + '\n')
    evidence['oldControllerSHA256'] = hashlib.sha256((base / 'old-controller').read_bytes()).hexdigest()
    (base / 'evidence.json').write_text(json.dumps(evidence, indent=2) + '\n')
    (base / 'README.md').write_text(
        '# Owner decision adapter trial\n\n'
        'This separate project uses real model/Blueprint services and scripted fixture-owner input. '
        'No actual project owner approval or agent launch occurred. The final decision is rejected. '
        'A fresh design-decision.json is available for an optional local CLI fixture interaction. '
        'The adapter displays the proposal before reading a choice and comment; it does not launch work. '
        'Use design-review with the current generation to inspect applicability; old receipts are historical.\n')
    print('PASS: compiled decision capture, cancel/EOF, immutable retry, stale input, supersession and old-reader refusal.')


if __name__ == '__main__':
    main()
