#!/usr/bin/env python3
"""BP2-B experimental selection over actual toolkit ASTs, not a language parser."""
import argparse
import hashlib
import json
import subprocess
from pathlib import Path

SUPPORTED = {'Relation', 'Allocation', 'Participation', 'PropertyAssignment'}


def clean(value):
    if isinstance(value, dict):
        return {k: clean(v) for k, v in value.items() if k not in {'span', 'path_span', 'path'}}
    if isinstance(value, list):
        return [clean(v) for v in value]
    return value


def identifiers(value):
    if isinstance(value, dict):
        if value.get('node') == 'Identifier':
            return {value['name']}
        return set().union(set(), *(identifiers(v) for v in value.values()))
    if isinstance(value, list):
        return set().union(set(), *(identifiers(v) for v in value))
    return set()


def key(value):
    return json.dumps(clean(value), sort_keys=True, separators=(',', ':'))


def load(tool, directory):
    entry = directory / 'System.design'
    def call(command):
        return json.loads(subprocess.check_output([tool, command, str(entry)], text=True))
    checked = call('check')
    if not checked.get('valid'):
        raise ValueError('invalid input')
    ast = call('ast')
    if not ast.get('valid'):
        raise ValueError('invalid AST')
    nodes, facts = {}, {}
    for file in ast['files']:
        for d in file['ast']['declarations']:
            nodes[d['name']['name']] = d['kind']
        for st in file['ast']['statements']:
            if st['node'] not in SUPPORTED:
                raise ValueError('unsupported relation semantics: ' + st['node'])
            facts[key(st)] = {'refs': sorted(identifiers(st)), 'source': st['span'], 'fact': clean(st)}
    return {'nodes': nodes, 'facts': facts, 'revision': ast['revision']}


def select(now, target, authored=(), limit=256, expected=None):
    if expected is not None and expected != (now['revision'], target['revision']):
        raise ValueError('stale inputs')
    nodes = dict(now['nodes'])
    nodes.update(target['nodes'])
    facts = dict(now['facts'])
    facts.update(target['facts'])
    changed = set(now['facts']) ^ set(target['facts'])
    seeds = set(now['nodes']) ^ set(target['nodes'])
    seeds.update(n for n in set(now['nodes']) & set(target['nodes']) if now['nodes'][n] != target['nodes'][n])
    for k in changed:
        seeds.update(facts[k]['refs'])
    if not set(authored) <= nodes.keys():
        raise ValueError('unknown authored inclusion')
    reasons = {n: ['changed fact/declaration'] for n in seeds}
    for n in authored:
        reasons.setdefault(n, []).append('authored context')
    selected, expanded, selected_facts = set(seeds) | set(authored), set(), set()
    while True:
        pending = sorted(selected - expanded)
        if not pending:
            break
        for n in pending:
            expanded.add(n)
            if nodes[n] in {'system', 'mode'}:
                continue  # Context labels do not imply sibling/runtime impact.
            for k in sorted(facts):
                f = facts[k]
                if n not in f['refs']:
                    continue
                selected_facts.add(k)
                for ref in f['refs']:
                    if ref not in nodes:
                        raise ValueError('dangling reference')
                    if ref not in selected:
                        reasons[ref] = [k]
                        selected.add(ref)
                if len(selected) > limit:
                    raise ValueError('selection overflow; no complete bundle')
        if len(selected) > limit:
            raise ValueError('selection overflow; no complete bundle')
    frontier = sorted(k for k, f in facts.items() if k not in selected_facts and set(f['refs']) & selected)
    gaps = []
    for side, model in [('NOW', now), ('TARGET', target)]:
        roles = {}
        for f in model['facts'].values():
            st = f.get('fact', {})
            if st.get('node') == 'Participation' and st['channel']['name'] in selected:
                group = tuple(st[k]['name'] for k in ('channel', 'message', 'mode'))
                roles.setdefault(group, set()).add(st['role'])
        for group, present in sorted(roles.items()):
            if present != {'sender', 'receiver'}:
                gaps.append(side + ' missing channel peer: ' + '/'.join(group))
    return {'nodes': sorted(selected), 'facts': sorted(selected_facts), 'reasons': reasons,
            'added': sorted(set(target['facts']) - set(now['facts'])),
            'removed': sorted(set(now['facts']) - set(target['facts'])),
            'frontier': frontier, 'excluded': sorted(set(nodes) - selected),
            'unknowns': ['No complete runtime dependency or code-to-model mapping is available'] + gaps,
            'revisions': [now['revision'], target['revision']]}


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--sdl', required=True)
    p.add_argument('--output', required=True)
    a = p.parse_args()
    root = Path(__file__).resolve().parent
    now, target = load(a.sdl, root / 'NOW'), load(a.sdl, root / 'TARGET')
    result = select(now, target, authored=['BuckingUI'])
    # Retain origins independently for both sides, including deleted edges.
    result['origins'] = {side: {k: f['source'] for k, f in model['facts'].items()
                              if k in result['facts']}
                         for side, model in [('NOW', now), ('TARGET', target)]}
    result['toolHash'] = hashlib.sha256(Path(a.sdl).read_bytes()).hexdigest()
    result['status'] = 'experimental selection; not a production blueprint'
    Path(a.output).write_text(json.dumps(result, indent=2, sort_keys=True) + '\n')

if __name__ == '__main__':
    main()
