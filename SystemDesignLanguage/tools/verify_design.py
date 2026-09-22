#!/usr/bin/env python3
"""Regenerate and verify the shared SDL design using the public toolkit commands."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[2]


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def run(args, cwd=ROOT, env=None):
    result = subprocess.run(args, cwd=cwd, env=env, text=True, capture_output=True, check=True)
    return result.stdout + result.stderr


def verify(renderer, phase):
    counts = {}
    suites = [('sdl_parser', ROOT, 'experiments/design_core', None),
              ('viewpoint_tool', ROOT, 'SystemDesignLanguage/tools', None),
              ('sdui', ROOT / 'SDUI', 'tests', dict(os.environ, PYTHONPATH='src'))]
    for name, cwd, directory, env in suites:
        output = run([sys.executable, '-m', 'unittest', 'discover', '-s', directory, '-q'], cwd, env)
        counts[name] = int(re.search(r'Ran (\d+) tests?', output).group(1))
    run([sys.executable, 'SDUI/tools/export_design.py'])
    command = [sys.executable, 'SystemDesignLanguage/tools/sdl.py', 'viewpoints',
               'SDUI/design/architecture.design', '--output', 'SDUI/design/viewpoints', '--renderer', str(renderer)]
    summary = json.loads(run(command))
    output = ROOT / 'SDUI/design/viewpoints'
    snapshot = lambda: {str(p.relative_to(output)): digest(p) for p in output.rglob('*') if p.is_file()}
    before = snapshot()
    run(command)
    assert snapshot() == before, 'Repeat export differs'
    manifest = json.loads((output / 'manifest.json').read_text())
    facts = {f['id']: f for f in manifest['facts']}
    source = (ROOT / manifest['source']).read_text()
    assert digest(ROOT / manifest['source']) == manifest['source_sha256']
    for fact in facts.values():
        span = fact['span']
        assert source[span['start']:span['end']] == fact['text'], fact['id']
    for name, sha in manifest['outputs'].items():
        assert digest(output / name) == sha, name
    for name, sha in manifest['generator_sha256'].items():
        assert digest(ROOT / name) == sha, name
    for diagram in manifest['diagrams']:
        svg = output / 'diagrams' / (diagram['id'] + '.svg')
        labels = ''.join(''.join(ET.parse(svg).getroot().itertext()).split())
        for node in diagram['nodes'].values():
            assert node['model_id'] in labels, (diagram['id'], node['model_id'])
            span = node['span']
            assert node['model_id'] in source[span['start']:span['end']]
        assert set(diagram['source_facts']) <= set(facts)
        for edge in diagram['edges']:
            f = facts[edge['fact']]
            assert (edge['source'], edge['relation'], edge['target']) == ('n_' + f['subject'], f['verb'], 'n_' + f['object'])
        if diagram['kind'] == 'packet':
            for e in diagram['elements']:
                f = facts[e['fact']]
                assert (e['field'], e['first_bit'], e['last_bit']) == (f['field'], f['offset'], f['offset'] + f['width'] - 1)
    validation = ROOT / 'SDUI/design/architecture.validation.json'
    for name, sha in json.loads(validation.read_text())['files'].items():
        assert digest(ROOT / name) == sha, name
    paths = list(manifest['generator_sha256']) + ['SDUI/design/architecture.design',
            'SDUI/design/architecture.validation.json', 'SDUI/design/viewpoints/manifest.json']
    report = dict(phase=phase, scope='Generated SDL model/views; no Go/runtime execution', tests=counts,
                  rendered_diagrams=summary['diagrams'], source_declarations=summary['declarations'],
                  source_facts=summary['facts'], unsupported_viewpoints=summary['unsupported_viewpoints'],
                  model_gap_count=len(manifest['model_gaps']),
                  checks=['all SVG node labels', 'exact source facts/spans', 'output/generator hashes',
                          'packet bit ranges', 'byte-identical repeat export'],
                  files={p: digest(ROOT / p) for p in paths})
    (Path(__file__).with_name('verification.json')).write_text(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
    return report


if __name__ == '__main__':
    cli = argparse.ArgumentParser(description=__doc__)
    cli.add_argument('--renderer', type=Path, required=True)
    cli.add_argument('--phase', required=True)
    args = cli.parse_args()
    try:
        print(json.dumps(verify(args.renderer.resolve(), args.phase), ensure_ascii=False))
    except (OSError, ValueError, AssertionError, subprocess.SubprocessError) as error:
        cli.exit(1, str(error) + '\n')
