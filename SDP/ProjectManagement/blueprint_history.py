"""Read-only canonical assignment replay, matching SDPTool/blueprintstate."""
import json
import re
from pathlib import Path, PurePosixPath
from datetime import datetime
from jsonschema import Draft202012Validator, FormatChecker

SCHEMA = Draft202012Validator(json.loads((Path(__file__).parent / 'blueprint-assignment.schema.json').read_text()), format_checker=FormatChecker())

def require(ok, message):
    if not ok:
        raise ValueError('blueprint history: ' + message)

def relative(p):
    return p not in ('', '.', '..') and not p.startswith(('/', '../')) and not any(c in p for c in '\\:\x00') and str(PurePosixPath(p)) == p and '..' not in PurePosixPath(p).parts

def instant(value):
    match = re.fullmatch(r"(\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2})(?:\.(\d{1,9}))?(Z|[+-]\d{2}:\d{2})", value)
    require(match, 'timestamp syntax')
    seconds = int(datetime.fromisoformat(match[1]+match[3]).timestamp())
    return seconds, int((match[2] or '').ljust(9, '0'))

def validate(events):
    states = {}
    for e in events:
        if not e['eventType'].startswith('x-blueprint:'):
            continue
        require(set(e) == {'schemaVersion', 'eventId', 'eventType', 'occurredAt', 'actor', 'commit', 'subjectId', 'payload'}, 'event fields')
        p = e['payload']
        SCHEMA.validate(p)
        r, proof = p['request'], p.get('proof')
        ident, action = r['assignmentId'], r['action']
        old = states.get(ident, {})
        origin = old.get('state', '')
        require(e['schemaVersion'] == '1.0' and e['eventId'] == r['eventId'] and e['subjectId'] == ident and e['eventType'] == 'x-blueprint:' + action and e['occurredAt'] == r['occurredAt'], 'envelope binding')
        require(e['actor'].strip() and r['reason'].strip(), 'empty attribution')
        require(r['expectedEvent'] == old.get('event', ''), 'predecessor mismatch')
        require(origin not in ('completed', 'canceled', 'superseded'), 'terminal state')
        instant(e['occurredAt'])
        if old:
            require(instant(e['occurredAt']) >= instant(old['time']), 'time reversal')
        required = 'assignee' if action in ('start', 'submit') else 'reviewer' if action in ('reject', 'accept-review') else 'controller'
        require(p['authority'] == required, 'authority')
        if required == 'assignee':
            require(old.get('assignee') == e['actor'], 'wrong assignee')
        if required == 'reviewer':
            require(old.get('assignee') and old['assignee'] != e['actor'], 'review independence')
        for key, wanted in [('binding', action == 'create'), ('assignee', action == 'assign'), ('successor', action == 'supersede'), ('evidence', action in ('adopt-readiness', 'submit')), ('traceEvent', action == 'submit')]:
            require(bool(r.get(key)) == wanted, 'unexpected ' + key)
        if 'evidence' in r:
            require(relative(r['evidence']), 'evidence path')
        needs = action in ('adopt-readiness', 'submit', 'accept-review', 'complete')
        require((proof is not None) == needs, 'proof presence')
        if proof:
            if action == 'adopt-readiness':
                require(not any(proof.get(k) for k in ('codeDigest', 'traceEvent', 'traceDigest')), 'readiness proof')
            else:
                require(all(proof.get(k) for k in ('codeDigest', 'traceEvent', 'traceDigest')), 'implementation proof')
        n = dict(old, event=e['eventId'], time=e['occurredAt'])
        if action == 'create':
            b = r['binding']
            require(not old and not r['expectedEvent'], 'creation')
            require(relative(b['bundle']) and (b['modelArea'] == '.' or relative(b['modelArea'])), 'binding paths')
            n.update(state='draft', binding=b)
        elif action == 'adopt-readiness':
            require(origin == 'draft', 'readiness origin')
            n.update(state='ready', readiness=proof)
        elif action == 'assign':
            require(origin == 'ready' and r['assignee'].strip(), 'assign origin')
            n.update(state='assigned', assignee=r['assignee'])
        elif action == 'start':
            require(origin == 'assigned', 'start origin')
            n['state'] = 'in-progress'
        elif action == 'submit':
            require(origin == 'in-progress' and proof['traceEvent'] == r['traceEvent'], 'submission')
            n.update(state='review', submission=proof, accepted='')
        elif action == 'reject':
            require(origin == 'review', 'reject origin')
            n.update(state='in-progress', submission=None, accepted='')
        elif action == 'accept-review':
            require(origin == 'review' and not old.get('accepted') and old.get('submission') == proof, 'accept proof')
            n['accepted'] = e['eventId']
        elif action == 'complete':
            require(origin == 'review' and old.get('accepted') and old.get('submission') == proof, 'complete proof')
            n['state'] = 'completed'
        elif action == 'hold':
            require(origin and origin != 'on-hold', 'hold origin')
            n.update(state='on-hold', paused=origin)
        elif action == 'resume':
            require(origin == 'on-hold' and old.get('paused'), 'resume origin')
            n.update(state=old['paused'], paused='')
        elif action in ('cancel', 'supersede'):
            require(bool(old), 'disposition origin')
            if action == 'supersede':
                require(r['successor'] != ident and r['successor'] in states and states[r['successor']]['binding']['system'] == old['binding']['system'], 'successor')
            n['state'] = 'canceled' if action == 'cancel' else 'superseded'
        states[ident] = n
    return states
