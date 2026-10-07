#!/usr/bin/env python3
"""Recreate the authored BP2-C worker candidate and injected negative control."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess

p = argparse.ArgumentParser()
p.add_argument('--sdl', required=True)
p.add_argument('--output', required=True)
a = p.parse_args()
root = Path(a.output)
root.mkdir(exist_ok=False)
base = Path(__file__).resolve().parent / 'NOW'
report = {'toolCommit': '043c59c',
          'trial': 'authored temporary model assignment, not generated production bundle',
          'candidates': {}}
for name in ['CandidateA', 'CandidateB']:
    dest = root / name
    shutil.copytree(base, dest)
    f = dest / 'Containers/MachineService.design'
    s = f.read_text().replace('container MachineService.', 'container MachineService.\nunit CalibrationProcessor.')
    s = s.replace('MachineService owns CalibrateMeasurement.',
                  'MachineService contains CalibrationProcessor.\nCalibrationProcessor owns CalibrateMeasurement.')
    if name == 'CandidateB':
        s = s.replace('MachineService owns ReduceMachineState.', 'BuckingWeb owns ReduceMachineState.')
    f.write_text(s)
    formatted = json.loads(subprocess.check_output([a.sdl, 'format', str(dest / 'System.design'), '--file-map']))
    for rel, content in formatted['files'].items():
        (dest / rel).write_text(content)
    check = json.loads(subprocess.check_output([a.sdl, 'check', str(dest / 'System.design')]))
    if not check.get('valid'):
        raise SystemExit(check)
    report['candidates'][name] = {
        'check': check,
        'files': {str(f.relative_to(dest)): hashlib.sha256(f.read_bytes()).hexdigest()
                  for f in sorted(dest.rglob('*.design'))},
        'changedFiles': sorted(str(f.relative_to(dest)) for f in dest.rglob('*.design')
                               if f.read_bytes() != (base / f.relative_to(dest)).read_bytes())}
(root / 'evidence.json').write_text(json.dumps(report, indent=2) + '\n')
