#!/usr/bin/env python3
"""Reproduce BP2-A model snapshots with existing tools; not a blueprint generator."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess

p = argparse.ArgumentParser()
p.add_argument('--sdptool', required=True)
p.add_argument('--sdl', required=True)
p.add_argument('--output', required=True)
a = p.parse_args()
root = Path(__file__).resolve().parent
out = Path(a.output).resolve()
if out.exists():
    raise SystemExit('output must not exist')
out.mkdir(parents=True)
area = out / 'models'
area.mkdir()
def run(tool, *args):
    completed = subprocess.run([tool, *map(str, args)], check=True, capture_output=True, text=True)
    return json.loads(completed.stdout)
def model(*args):
    return run(a.sdptool, area, 'model', '--json', *args)
def files(path):
    return {str(f.relative_to(path)): hashlib.sha256(f.read_bytes()).hexdigest()
            for f in sorted(path.rglob('*')) if f.is_file()}
results = {}
for state in ('NOW', 'TARGET'):
    check = run(a.sdl, 'check', root / state / 'System.design')
    if not check.get('valid'):
        raise SystemExit(check)
    results[state] = {'check': check, 'sourceHashes': files(root / state)}
work = model('create', 'work:Baseline', '--initial')
shutil.copytree(root / 'NOW', work['path'], dirs_exist_ok=True)
model('create', 'candidate:Baseline', 'from', 'work:Baseline')
results['release'] = model('create', 'release:0.1.0', 'from', 'candidate:Baseline', '--evidence', 'model-only')
work = model('create', 'work:ExtractCalibration')
shutil.copytree(root / 'TARGET', work['path'], dirs_exist_ok=True)
results['preliminary'] = model('snapshot', 'work:ExtractCalibration')
results['candidate'] = model('create', 'candidate:ExtractCalibration', 'from', 'work:ExtractCalibration')
if results['preliminary']['liveDigest'] != results['candidate']['liveDigest']:
    raise SystemExit('candidate differs from captured WORK')
for state in ('NOW', 'TARGET'):
    results[state]['viewpoints'] = run(a.sdl, 'viewpoints', root / state / 'System.design',
        '--format', 'static', '--viewpoint', 'VP02,VP03', '--output', out / state,
        '--project', 'mvp1-blueprint-pilot')
results['toolHashes'] = {tool: hashlib.sha256(Path(tool).read_bytes()).hexdigest() for tool in (a.sdptool, a.sdl)}
results['limitations'] = ['Authored reduced MVP1 fixture, not a port of the experimental model',
 'Ordinary toolkit viewpoints, not semantic diff or generated blueprint annotations',
 'Model-only release; no Ponsse code conformance or owner acceptance implied']
(out / 'evidence.json').write_text(json.dumps(results, indent=2) + '\n')
print(json.dumps({'evidence': str(out / 'evidence.json'), 'validated': True}))
