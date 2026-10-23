#!/usr/bin/env python3
"""Inspect the isolated decision fixture; --decide explicitly opens its local input adapter."""
import argparse
import json
from pathlib import Path
import subprocess
import tempfile
import uuid

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('--preview', type=Path, default=Path(__file__).resolve().parent / 'previews/pgl2-c2c')
parser.add_argument('--json', action='store_true')
parser.add_argument('--decide', action='store_true')
parser.add_argument('--owner', help='attributed owner label for an explicit fixture decision')
args = parser.parse_args()
if args.decide != bool(args.owner):
    parser.error('--decide and --owner must be supplied together')
p = args.preview.resolve()
state = p / 'runtime/state'
check = json.loads((p / 'design-review.json').read_text())
check['expectedGeneration'] = (state / 'HEAD').read_text()
base = [str(p / 'sdp-controller'), 'assignment']
scope = ['--project', str(p / 'project'), '--state', str(state)]

with tempfile.TemporaryDirectory(prefix='pg-decision-view-') as tmp:
    request = Path(tmp) / 'request.json'
    request.write_text(json.dumps(check))
    done = subprocess.run(base + ['design-review'] + scope + ['--input', str(request)], capture_output=True, text=True)
    if done.returncode:
        raise SystemExit(done.stderr)
    view = json.loads(done.stdout)
    if args.decide:
        if view['check']['status'] != 'current':
            raise SystemExit('Design is blocked. Inspect it before making another fixture decision.')
        q = json.loads((p / 'design-decision.json').read_text())
        q.update(operationId=str(uuid.uuid4()), review=check, reviewRevision=view['check']['reviewRevision'])
        request.write_text(json.dumps(q))
        print('Isolated synthetic project. Your input affects this fixture only; no agent will launch.', flush=True)
        raise SystemExit(subprocess.run(base + ['design-decide'] + scope + ['--input', str(request), '--owner', args.owner]).returncode)
    if args.json:
        print(json.dumps(view, indent=2))
    else:
        print('PG design review: ' + view['acceptance'])
        print('Current technical check: ' + view['check']['status'])
        if view.get('decision'):
            captured = view['decision']['record']['input']
            print('Recorded decision: ' + captured['decision'])
            print('Attributed actor: ' + json.dumps(captured['actor']))
            print('Comment: ' + json.dumps(captured['comment']))
        for condition in view['check']['conditions']:
            print('  ' + json.dumps(condition))
        print('Dispatch: ' + view['dispatch'])
        print('Separate synthetic preview. Scripted fixture decisions are not actual project owner approval.')
