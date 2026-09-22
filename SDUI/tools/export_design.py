#!/usr/bin/env python3
"""Inspect the shared SDL/SDUI design using the existing design-core frontend."""
from collections import Counter, defaultdict
import hashlib
import json
from pathlib import Path
import sys

ROOT = Path(__file__).resolve().parents[2]
sys.path.insert(0, str(ROOT / 'experiments/design_core'))
import design_core as core


def export():
    source_path = ROOT / 'SDUI/design/architecture.design'
    source = source_path.read_text(encoding='utf-8')
    model, diagnostics = core.check(source)
    if diagnostics:
        print(json.dumps(core.to_json(diagnostics), indent=2), file=sys.stderr)
        return 1
    # Refuse drift in canonical source, before replacing any derived artifact.
    assert core.canonicalize(model) == source
    design_dir = source_path.parent
    artifact = {'valid': True, 'diagnostics': [], 'ast': core.to_json(model),
                'symbols': core.to_json(core.symbol_table(model))}
    (design_dir / 'architecture.ast.json').write_text(
        json.dumps(artifact, indent=2, ensure_ascii=False) + '\n', encoding='utf-8')
    owners, children, interfaces = defaultdict(list), defaultdict(list), defaultdict(list)
    for statement in model.statements:
        if isinstance(statement, core.Relation):
            target = {'owns': owners, 'contains': children, 'consumes': interfaces}.get(statement.verb)
            if target is not None:
                target[statement.subject.name].append(statement.object.name)
    lines = ['# Parser/runtime — generert ansvarsoversikt', '',
             'Generert fra [architecture.design](architecture.design) med eksisterende design-core-parser.',
             'Dette er en visning av målstruktur, ikke implementasjonsbevis. Ikke rediger denne filen manuelt.', '',
             '| Unit / Container | Underliggende Units | Eide Functionality-er | Brukte grensesnitt |',
             '| --- | --- | --- | --- |']
    for declaration in model.declarations:
        if declaration.kind in ('unit', 'container'):
            name = declaration.name.name
            cells = [name, ', '.join(children[name]) or '—', ', '.join(owners[name]) or '—',
                     ', '.join(interfaces[name]) or '—']
            lines.append('| ' + ' | '.join(cells) + ' |')
    (design_dir / 'architecture.catalog.md').write_text('\n'.join(lines) + '\n', encoding='utf-8')
    counts = Counter(d.kind for d in model.declarations)
    statements = Counter(core.sentence(s).split()[1] for s in model.statements)
    files = [source_path, ROOT / 'experiments/design_core/design_core.py',
             ROOT / 'experiments/design_core/data_core.py',
             ROOT / 'docs/Design-Language-Definition.md', Path(__file__).resolve(),
             design_dir / 'architecture.ast.json', design_dir / 'architecture.catalog.md']
    report = {
        'profile': 'design-core ' + model.header.version, 'valid': True, 'diagnostics': [],
        'canonical': True, 'declarations': dict(sorted(counts.items())),
        'statements': dict(sorted(statements.items())),
        'scope': 'Structural SDL model only. No runtime, UI or Go code executed.',
        'files': {str(p.relative_to(ROOT)): hashlib.sha256(p.read_bytes()).hexdigest() for p in files},
    }
    (design_dir / 'architecture.validation.json').write_text(
        json.dumps(report, indent=2, ensure_ascii=False) + '\n', encoding='utf-8')
    print(json.dumps({'valid': True, 'declarations': sum(counts.values()),
                      'facts': len(model.statements), 'kinds': dict(counts)}, ensure_ascii=False))
    return 0


if __name__ == '__main__':
    raise SystemExit(export())
