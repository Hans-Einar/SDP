#!/usr/bin/env python3
"""Generate selected viewpoints from validated SDL; no inferred architecture or new grammar."""
import argparse
import hashlib
import json
from pathlib import Path
import re
import shutil
import subprocess
import tempfile
import xml.etree.ElementTree as ET
from viewpoints import CATALOG, Views


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def read_previous(output):
    try:
        return json.loads((output / 'manifest.json').read_text(encoding='utf-8'))
    except (OSError, ValueError):
        return {}


def build(args, views, output, previous):
    renderer = args.renderer.resolve() if args.renderer else None
    diagrams = output / 'diagrams'
    diagrams.mkdir()
    (output / 'viewpoints.md').write_text(views.markdown(), encoding='utf-8')
    if 'VP08' in views.selected:
        (output / 'message-sets.json').write_text(json.dumps(views.message_sets, ensure_ascii=False, indent=2) + '\n', encoding='utf-8')
    manifest = {'source': str(args.source), 'source_sha256': digest(args.source),
                'scope': f'design-core {views.model.header.version} structural projections; missing concepts are not inferred',
                'viewpoints': [{'id': i, 'title': t, 'status': s, 'note': n} for i, t, s, n in CATALOG if i in views.selected],
                'facts': views.facts, 'model_gaps': views.gaps, 'diagrams': [], 'renderer': None}
    if renderer:
        manifest['renderer'] = {'path': str(renderer), 'sha256': digest(renderer),
                               'version': subprocess.check_output([str(renderer), '--version'], text=True).strip()}
    prior_diagrams = {d['id']: d for d in previous.get('diagrams', [])}
    for d in views.diagrams:
        mmd = diagrams / (d.ident + '.mmd')
        mmd.write_text(d.mermaid(views.kinds), encoding='utf-8')
        item = {'id': d.ident, 'kind': d.kind, 'mmd_sha256': digest(mmd), 'nodes': views.node_map(d),
                'source_facts': sorted(set(d.source_facts + [e[3] for e in d.edges])), 'elements': d.elements,
                'edges': [{'source': f'n_{a}', 'relation': b, 'target': f'n_{c}', 'fact': f}
                          for a, b, c, f in sorted(d.edges)]}
        if renderer:
            svg = diagrams / (d.ident + '.svg')
            cached = args.output / 'diagrams' / svg.name
            prior = prior_diagrams.get(d.ident, {})
            if (previous.get('renderer') == manifest['renderer'] and
                    prior.get('mmd_sha256') == item['mmd_sha256'] and cached.is_file() and
                    prior.get('svg_sha256') == digest(cached)):
                shutil.copyfile(cached, svg)
            else:
                subprocess.run([str(renderer), '-i', str(mmd), '-o', str(svg)], check=True,
                               capture_output=True, text=True, timeout=60)
            tree = ET.parse(svg)
            if tree.getroot().tag != '{http://www.w3.org/2000/svg}svg':
                raise ValueError('Renderer did not produce SVG: ' + d.ident)
            item.update(svg_sha256=digest(svg), svg_size={k: tree.getroot().get(k) for k in ['width', 'height', 'viewBox']})
        manifest['diagrams'].append(item)
    if renderer:
        (output / 'printout.md').write_text(views.markdown(rendered=True), encoding='utf-8')
    root = Path(__file__).resolve().parents[2]
    sources = [p for directory in (Path(__file__).parent, root / 'experiments/design_core')
               for p in sorted(directory.glob('*.py')) if not p.name.startswith('test_')]
    manifest['generator_sha256'] = {str(p.resolve().relative_to(root)): digest(p) for p in sources}
    manifest['outputs'] = {str(p.relative_to(output)): digest(p) for p in sorted(output.rglob('*')) if p.is_file()}
    (output / 'manifest.json').write_text(json.dumps(manifest, indent=2, ensure_ascii=False) + '\n', encoding='utf-8')
    return manifest


def publish(stage, output, previous, manifest):
    output.mkdir(parents=True, exist_ok=True)
    # Remove only known generated names from the previous successful export.
    old_names = set(previous.get('outputs', {}))
    if previous:  # Compatibility with the first generator manifest, no semantic fallback.
        old_names.add('printout.md')
        for d in previous.get('diagrams', []):
            old_names.update('diagrams/' + d['id'] + ext for ext in ('.mmd', '.svg'))
    for name in sorted(old_names - set(manifest['outputs'])):
        if name not in ('printout.md', 'viewpoints.md', 'message-sets.json') and not re.fullmatch(r'diagrams/VP\d+-[A-Za-z0-9-]+\.(mmd|svg)', name):
            continue
        p = output / name
        if p.is_file() and p.resolve().is_relative_to(output.resolve()):
            p.unlink()
    for p in sorted(stage.rglob('*')):
        if p.is_file():
            target = output / p.relative_to(stage)
            target.parent.mkdir(parents=True, exist_ok=True)
            p.replace(target)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    command = parser.add_subparsers(dest='command', required=True)
    export = command.add_parser('viewpoints', help='Generate Markdown/Mermaid from a validated SDL source')
    export.add_argument('source', type=Path)
    export.add_argument('--output', type=Path, required=True)
    export.add_argument('--viewpoint', action='append', choices=[v[0] for v in CATALOG], help='Repeat to select viewpoints; default is all')
    export.add_argument('--renderer', type=Path, help='Optional explicit mmdr executable for SVG printout')
    args = parser.parse_args()
    try:
        if args.source.resolve().is_relative_to(args.output.resolve()):
            raise ValueError('Keep the source outside the generated output directory')
        views = Views(args.source.read_text(encoding='utf-8'), args.viewpoint)
        previous = read_previous(args.output)
        args.output.parent.mkdir(parents=True, exist_ok=True)
        with tempfile.TemporaryDirectory(prefix='.sdl-viewpoints-', dir=args.output.parent) as temp:
            stage = Path(temp)
            manifest = build(args, views, stage, previous)
            publish(stage, args.output, previous, manifest)
    except (OSError, UnicodeError, ValueError, subprocess.SubprocessError, ET.ParseError) as error:
        parser.exit(2, str(error) + '\n')
    print(json.dumps({'diagrams': len(views.diagrams), 'rendered': bool(args.renderer),
                      'declarations': len(views.model.declarations), 'facts': len(views.facts),
                      'unsupported_viewpoints': [i for i, _, s, _ in CATALOG if s == 'blocked' and i in views.selected]}))


if __name__ == '__main__':
    main()
