#!/usr/bin/env python3
"""Contract checks for fictional lineage events, not a production board tool.

Requires Python 3 and jsonschema. No file mutations or network access.
"""
import copy
import json
from datetime import datetime
from pathlib import Path

from jsonschema import Draft202012Validator, FormatChecker, ValidationError

BOARD = Path(__file__).resolve().parent.parent
ROOT = BOARD.parents[2]
ENVELOPE = Draft202012Validator(
    json.loads((ROOT / 'Toolkit/schemas/ledger-event.schema.json').read_text()),
    format_checker=FormatChecker(),
)
PAYLOADS = {
    version: Draft202012Validator(json.loads((BOARD / filename).read_text()))
    for version, filename in (
        ('0.1', 'ledger-payload.schema.json'),
        ('0.2', 'ledger-payload-0.2.schema.json'),
    )
}


def require(condition, message):
    if not condition:
        raise ValueError(message)


def validate(events):
    """Replay complete supplied examples and verify their lineage participants."""
    latest, event_ids, operations = {}, set(), {}
    open_states = {'backlog', 'active', 'onHold'}
    for event in events:
        ENVELOPE.validate(event)
        p = event['payload']
        require(p['schemaVersion'] in PAYLOADS, 'unknown payload version')
        PAYLOADS[p['schemaVersion']].validate(p)
        subject, typ = event['subjectId'], event['eventType']
        require(event['eventId'] not in event_ids, 'duplicate event ID')
        event_ids.add(event['eventId'])
        require(subject.startswith('KB-' + p['projectId'] + '-'), 'wrong project')
        prev = latest.get(subject)
        if prev is None:
            require(typ == 'x-kanban:created', 'missing creation')
            require(p['previousEventId'] is p['from'] is p['fromPath'] is None,
                    'invalid creation predecessor')
            require(p['to'] in {'backlog', 'active'}, 'invalid initial state')
        else:
            require(p['previousEventId'] == prev['eventId'], 'broken predecessor')
            require((p['from'], p['fromPath']) ==
                    (prev['payload']['to'], prev['payload']['toPath']), 'wrong origin')
            require(datetime.fromisoformat(event['occurredAt']) >=
                    datetime.fromisoformat(prev['occurredAt']), 'time reversal')
            same = (p['from'], p['fromPath']) == (p['to'], p['toPath'])
            require((typ == 'x-kanban:reviewed' and same) or
                    (typ == 'x-kanban:moved' and not same), 'invalid transition')
        for side in ('from', 'to'):
            if p[side] is not None:
                require(Path(p[side + 'Path']).parent.name == p[side], 'wrong directory')
                require(Path(p[side + 'Path']).name.startswith('#' + subject.split('-')[-1] + '--'),
                        'path/card mismatch')
        if p['to'] in {'completed', 'superseded'}:
            require(bool(p['links']), 'missing outcome')
        lineage = p.get('lineage')
        if lineage:
            sources = {s['id']: s for s in lineage['sources']}
            targets = {t['id'] for t in lineage['targets']}
            require(len(sources) == len(lineage['sources']), 'duplicate source')
            require(len(targets) == len(lineage['targets']), 'duplicate target')
            require(not sources.keys() & targets, 'source/target overlap')
            require(subject in sources or subject in targets, 'unrelated participant')
            opid = lineage['operationId']
            if opid not in operations:
                require(all(s in latest and latest[s]['payload']['to'] in open_states
                            for s in sources), 'missing or closed source')
                require(all(t not in latest for t in targets), 'target already exists')
                operations[opid] = {'definition': lineage, 'members': set()}
            op = operations[opid]
            require(op['definition'] == lineage, 'inconsistent operation')
            require(subject not in op['members'], 'duplicate participant')
            if subject in targets:
                require(typ == 'x-kanban:created', 'target must be new')
                require(sources.keys() <= set(p['links']), 'missing source links')
            else:
                require(targets <= op['members'], 'targets must be created first')
                source = sources[subject]
                require(targets <= set(p['links']), 'missing successor links')
                require(p['from'] in open_states, 'source already closed')
                if source['remaining']:
                    require(bool(source['remaining'].strip()), 'blank remainder')
                    require(typ == 'x-kanban:reviewed' and p['to'] in open_states,
                            'remaining work lost')
                else:
                    require(typ == 'x-kanban:moved' and p['to'] == 'superseded',
                            'fully replaced source left open')
            op['members'].add(subject)
        latest[subject] = event
    for op in operations.values():
        definition = op['definition']
        expected = {x['id'] for x in definition['sources'] + definition['targets']}
        require(op['members'] == expected, 'incomplete operation')
    return latest, operations


def main():
    fixture = [json.loads(line) for line in
               (BOARD / 'examples/lineage.ndjson').read_text().splitlines()]
    latest, operations = validate(fixture)
    require(len(operations) == 4, 'expected full/partial merge and split')
    require(latest['KB-DEMO-004']['payload']['to'] == 'active', 'remainder closed')
    require(latest['KB-DEMO-006']['payload']['to'] == 'active', 'split remainder closed')
    require(latest['KB-DEMO-005']['payload']['to'] == 'superseded', 'source not closed')
    # Historical payloads remain valid; never rewrite the real ledgers for tests.
    legacy = copy.deepcopy(fixture[:2])
    for e in legacy:
        e['payload']['schemaVersion'] = '0.1'
    validate(legacy)

    def reject(name, mutate):
        candidate = copy.deepcopy(fixture)
        mutate(candidate)
        try:
            validate(candidate)
        except (ValueError, ValidationError):
            return
        raise AssertionError('Accepted invalid fixture: ' + name)

    cases = [
        ('missing participant', lambda es: es.pop()),
        ('missing successor link', lambda es: es[7]['payload']['links'].pop()),
        ('missing source link', lambda es: es[2]['payload']['links'].pop()),
        ('different operation definition', lambda es: es[3]['payload']['lineage']['targets'][0].update(scope='Different')),
        ('wrong predecessor', lambda es: es[3]['payload'].update(previousEventId=None)),
        ('closed remainder', lambda es: es[9].update(eventType='x-kanban:moved')),
        ('duplicate event', lambda es: es.append(copy.deepcopy(es[-1]))),
        ('invalid directory', lambda es: es[7]['payload'].update(toPath='active/#003--Idea--Example.md')),
        ('wrong operation kind', lambda es: es[2]['payload']['lineage'].update(kind='rename')),
        ('unknown payload version', lambda es: es[2]['payload'].update(schemaVersion='9.0')),
    ]
    for name, mutate in cases:
        reject(name, mutate)

    def change_group(events, transform):
        for e in events[2:5]:
            transform(e['payload']['lineage'])

    grouped = [
        ('merge with one source', lambda op: op['sources'].pop()),
        ('duplicate source ID', lambda op: op['sources'][1].update(id='KB-DEMO-001')),
        ('source/target overlap', lambda op: op['targets'][0].update(id='KB-DEMO-001')),
        ('unknown source', lambda op: op['sources'][0].update(id='KB-DEMO-099')),

    ]
    for name, transform in grouped:
        reject(name, lambda es, transform=transform: change_group(es, transform))
    reject('existing split target', lambda es: [e['payload']['lineage']['targets'][0].update(id='KB-DEMO-001') for e in es[5:8]])
    print(f'PASS: full/partial merge and split; historical 0.1; {len(cases) + len(grouped) + 1} negative cases')


if __name__ == '__main__':
    main()
