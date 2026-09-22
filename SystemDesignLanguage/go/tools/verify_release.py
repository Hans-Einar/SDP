#!/usr/bin/env python3
"""Run Go release checks; this driver contains no SDL/SDUI language logic."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import subprocess
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[3]

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('--go', default='go')
    parser.add_argument('--renderer', required=True, help='Explicit mmdr executable')
    args = parser.parse_args()
    env = dict(os.environ, MMDR=str(Path(args.renderer).resolve()),
               SDUI_MMDR=str(Path(args.renderer).resolve()))
    report = {'toolchain': subprocess.check_output([args.go, 'version'], text=True).strip(),
              'renderer_sha256': hashlib.sha256(Path(args.renderer).read_bytes()).hexdigest(),
              'tests': {}, 'bundles': {}, 'active_sources': {}}
    for module in ['SDUI/go', 'SystemDesignLanguage/go']:
        proc = subprocess.run([args.go, '-C', str(ROOT / module), 'test', '-race', '-json', './...'],
                              env=env, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        events = [json.loads(line) for line in proc.stdout.splitlines() if line.startswith('{')]
        report['tests'][module] = {
            'exit_code': proc.returncode,
            'passed_tests_and_subtests': sum(e.get('Action') == 'pass' and bool(e.get('Test')) for e in events),
            'passed_packages': sum(e.get('Action') == 'pass' and not e.get('Test') for e in events),
            'failed': [e for e in events if e.get('Action') == 'fail'],
        }
        if proc.returncode:
            raise RuntimeError(proc.stderr + proc.stdout)
        print(module, report['tests'][module], flush=True)
    current = hashlib.sha256((ROOT / 'SDUI/design/architecture.design').read_bytes()).hexdigest()
    for name in ['viewpoints', 'navigation', 'runtime-preview']:
        directory = ROOT / 'SDUI/design' / name
        manifest = json.loads((directory / 'manifest.json').read_text())
        for path, digest in manifest['outputs'].items():
            target = directory / path
            if hashlib.sha256(target.read_bytes()).hexdigest() != digest:
                raise RuntimeError(f'Changed generated file: {target}')
            if target.suffix == '.svg':
                ET.parse(target)
        if name != 'runtime-preview' and manifest['revision'] != current:
            raise RuntimeError('Stale design export: ' + name)
        if name == 'runtime-preview':
            p = json.loads((directory / 'provenance.json').read_text())
            if p['design_sha256'] != current:
                raise RuntimeError('Stale UI design provenance')
        report['bundles'][name] = {'files': len(manifest['outputs']), 'revision': manifest['revision']}
    for module in ['SDUI/go', 'SystemDesignLanguage/go']:
        for path in sorted((ROOT / module).rglob('*.go')):
            report['active_sources'][str(path.relative_to(ROOT))] = hashlib.sha256(path.read_bytes()).hexdigest()
    for directory in ['SDUI/src', 'SDUI/tests', 'SystemDesignLanguage/tools', 'experiments/design_core']:
        if list((ROOT / directory).rglob('*.py')):
            raise RuntimeError('Replaced Python source still present: ' + directory)
    output = ROOT / 'SystemDesignLanguage/go/evidence/release.json'
    output.write_text(json.dumps(report, indent=2, sort_keys=True) + '\n')
    print(output)

if __name__ == '__main__':
    main()
