#!/usr/bin/env python3
"""Validate local management history and current records; never mutate them."""
import json
import re
import runpy
from datetime import datetime
from pathlib import Path
from jsonschema import Draft202012Validator, FormatChecker

AREA = Path(__file__).resolve().parent
SDP = AREA.parent
ROOT = SDP.parent
BOARD = SDP / 'KanBan'
ENVELOPE = Draft202012Validator(json.loads((ROOT / 'Toolkit/schemas/ledger-event.schema.json').read_text()), format_checker=FormatChecker())
PAYLOAD = Draft202012Validator(json.loads((AREA / 'management-payload.schema.json').read_text()))
KANBAN = runpy.run_path(str(BOARD / 'examples/verify_lineage.py'))['validate']
STATES = {'backlog': {'backlog', 'queued'}, 'active': {'ready', 'in-progress', 'gate-review'}, **{s: {s} for s in ('onHold', 'completed', 'canceled', 'superseded', 'irrelevant')}}
KINDS = {'Scrum': 'SCRUM', 'Sprint': 'SPR', 'Maintenance': 'MAINT', 'CodeReview': 'REVIEW', 'Refactor': 'REFACTOR'}


def require(ok, message):
    if not ok:
        raise ValueError(message)


def metadata(path):
    result = {}
    started = False
    for line in path.read_text().splitlines():
        if re.fullmatch(r'\|\s*Field\s*\|\s*Value\s*\|', line):
            started = True
            continue
        if not started:
            if not line.strip() or line.startswith('# '):
                continue
            break
        if not line.startswith('|'):
            break
        parts = line.split('|')
        require(len(parts) == 4, f'{path}: malformed metadata row')
        key, value = (s.strip() for s in parts[1:3])
        if key.startswith('---'):
            continue
        require(key not in result, f'{path}: duplicate {key}')
        result[key] = value
    return result


def history(events):
    seen, last_numbers, management = set(), {}, {}
    for e in events:
        ENVELOPE.validate(e)
        ident = e['eventId']
        require(ident not in seen, 'duplicate event ID')
        seen.add(ident)
        prefix, number = ident.rsplit('-', 1)
        number = int(number)
        require(number > last_numbers.get(prefix, -1), 'non-increasing namespace event ID')
        last_numbers[prefix] = number
        if e['eventType'].startswith('x-kanban:'):
            continue
        require(e['eventType'] in {'x-management:created', 'x-management:updated', 'x-management:started', 'x-management:completed', 'x-management:canceled'}, 'unknown event type')
        p = e['payload']
        PAYLOAD.validate(p)
        subject = e['subjectId']
        require(re.fullmatch(KINDS[p['kind']] + r'-SDP-\d{4,}', subject), 'kind/ID mismatch')
        require(prefix == 'EVT-PM-SDP', 'management event namespace')
        for key in ('fromPath', 'toPath'):
            path = p[key]
            if path is not None:
                require(not Path(path).is_absolute() and '..' not in Path(path).parts and '\\' not in path, 'invalid management path')
        prev = management.get(subject)
        action = e['eventType'].split(':')[1]
        if prev is None:
            require(action == 'created' and p['previousEventId'] is p['from'] is p['fromPath'] is None, 'invalid creation')
            require(p['to'] in ('planned', 'active'), 'invalid creation state')
        else:
            old = prev['payload']
            require(p['previousEventId'] == prev['eventId'], 'broken management chain')
            require((p['kind'], p['from'], p['fromPath']) == (old['kind'], old['to'], old['toPath']), 'wrong management origin')
            require(datetime.fromisoformat(e['occurredAt']) >= datetime.fromisoformat(prev['occurredAt']), 'time reversal')
            if action == 'updated':
                require(p['from'] == p['to'] and p['fromPath'] == p['toPath'] and p['from'] not in ('completed', 'canceled'), 'invalid update')
            elif action == 'started':
                require((p['from'], p['to']) == ('planned', 'active'), 'invalid start')
            elif action in ('completed', 'canceled'):
                require(p['from'] in ('planned', 'active') and p['to'] == action, 'invalid closure')
                require(bool(p['links']), 'closure lacks outcome/source links')
            else:
                require(False, 'recreated management subject')
        if p['kind'] == 'Sprint':
            require('members' in p, 'Sprint requires a membership snapshot')
        else:
            require('members' not in p, 'members only applies to Sprint')
        management[subject] = e
    cards, operations = KANBAN([e for e in events if e['eventType'].startswith('x-kanban:')])
    known = cards.keys() | management.keys()
    for e in management.values():
        require(set(e['payload']['links']) <= known, 'unresolved management link')
        if e['payload']['kind'] == 'Sprint':
            require(set(e['payload']['members']) <= cards.keys(), 'unresolved sprint member')
    return cards, management, operations


def validate_current(events, board=BOARD, sdp=SDP):
    cards, management, operations = history(events)
    actual, metas = {}, {}
    for folder, allowed in STATES.items():
        for path in (board / folder).glob('#*.md'):
            m = metadata(path)
            ident = m.get('id')
            require(ident not in actual, 'duplicate physical card')
            require(m.get('CardState') in allowed, f'{path}: invalid CardState/folder')
            require(m.get('project') == ident.split('-')[1], 'card namespace changed')
            actual[ident], metas[ident] = str(path.relative_to(board)), m
    require(actual == {i: e['payload']['toPath'] for i, e in cards.items()}, 'card/history placement mismatch')
    index = dict(re.findall(r'^\| (KB-[A-Z]+-\d+) \| [^|]+ \| (\w+) \|', (board / 'README.md').read_text(), re.M))
    require(index == {i: p.split('/')[0] for i, p in actual.items()}, 'board index mismatch')
    for ident, e in management.items():
        p = e['payload']
        m = metadata(sdp / p['toPath'])
        require(m.get('id') == ident and m.get('state') == p['to'], 'management document/history mismatch')
        if p['kind'] == 'Sprint':
            members = set(p['members'])
            declared = {s.strip() for s in m.get('Members', '').split(',') if s.strip()}
            tagged = {i for i, card in metas.items() if card.get('SprintId') == ident}
            require(members == declared == tagged, 'sprint membership mismatch')
            if p['to'] == 'active':
                require(all(actual[i].split('/')[0] != 'backlog' for i in members), 'started sprint has backlog members')
    for ident, m in metas.items():
        for field, kind in [('SprintId', 'Sprint'), ('ScrumId', 'Scrum')]:
            if field in m:
                require(m[field] in management and management[m[field]]['payload']['kind'] == kind, f'{ident}: unknown {field}')
    return len(cards), len(management), len(operations)


def main():
    descriptor = json.loads((BOARD / 'board.json').read_text())
    require(descriptor == {'schemaVersion': '0.2', 'projectId': 'SDP', 'namespaces': ['SDP', 'SDL', 'SDUI'], 'ledger': '../ProjectManagement/Ledger.ndjson', 'profile': 'sdp-project-management/0.1'}, 'unsupported board descriptor')
    events = [json.loads(s) for s in (AREA / 'Ledger.ndjson').read_text().splitlines()]
    c, m, o = validate_current(events)
    print(f'PASS: {c} cards, {m} management records, {o} lineage operations, {len(events)} events; schema/history/metadata/placement/grouping')


if __name__ == '__main__':
    main()
