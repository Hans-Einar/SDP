#!/usr/bin/env python3
"""Show the isolated implementation trial and its actual recorded result, without launching work."""
import argparse
import json
from pathlib import Path
import subprocess

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('--preview', type=Path, default=Path(__file__).resolve().parent / 'previews/pgl2-c2d')
p.add_argument('--json', action='store_true')
a = p.parse_args()
base = a.preview.resolve()
evidence = json.loads((base / 'evidence.json').read_text())
activation = evidence['activationResult']
run_id = activation['run']['id']
status = subprocess.run([str(base / 'sdp-controller'), 'status', '--project', str(base / 'project'),
                         '--state', str(base / 'runtime/state'), '--run', run_id], text=True, capture_output=True)
if status.returncode:
    raise SystemExit(status.stderr)
current = json.loads(status.stdout)
if a.json:
    print(json.dumps(current, indent=2))
else:
    run = current['run']
    print('PG implementation trial')
    print('Task: ' + json.dumps(run['assignment']['outcome']))
    print('Routine: ' + run['routine']['id'] + ' / ' + run['state'])
    print('Current governing sources: ' + current['coverage'])
    print('Last recorded host state: ' + run['attempts'][-1]['status'])
    if run['state'] == 'running' and run['attempts'][-1]['status'] == 'completed':
        print('Assignment report: not recorded; inspect the recorded host result before continuing.')
    print('Trial host: ' + evidence['hostKind'])
    for name in evidence.get('changedFiles', []):
        print('Changed file:')
        print((base / 'project' / name).as_uri())
    if evidence.get('goldenAfter'):
        print(evidence['goldenAfter'].strip())
    print('Protected project files preserved: ' + str(evidence['protectedFilesPreserved']))
    print('Scripted design decision; no actual project owner acceptance or independent review.')
    print('Host completion does not complete the governance assignment. This command only reads state.')
