#!/usr/bin/env python3
"""Read-only Activity/admission schema and causal replay validation (0.1/0.2)."""
import hashlib
import json
import posixpath
import re
from pathlib import Path, PurePosixPath

import yaml
from jsonschema import Draft202012Validator, FormatChecker

LIMIT = 16 << 20
STATES = {'planned', 'active', 'on-hold', 'completed', 'canceled'}
TRANSITIONS = {
    'planned': {'active', 'canceled'},
    'active': {'on-hold', 'completed', 'canceled'},
    'on-hold': {'active', 'canceled'},
    'completed': {'planned'}, 'canceled': {'planned'},
}


def require(condition, message):
    if not condition:
        raise ValueError(message)


def unique_pairs(pairs):
    result = {}
    for key, value in pairs:
        require(key not in result, f'duplicate JSON key: {key}')
        result[key] = value
    return result


class UniqueYAML(yaml.SafeLoader):
    pass


def yaml_mapping(loader, node, deep=False):
    return unique_pairs([(loader.construct_object(k, deep=deep),
                          loader.construct_object(v, deep=deep)) for k, v in node.value])


UniqueYAML.add_constructor(yaml.resolver.BaseResolver.DEFAULT_MAPPING_TAG, yaml_mapping)


def load_events(path):
    data = Path(path).read_bytes()
    require(len(data) <= LIMIT, 'canonical history exceeds capacity')
    require(not data or data.endswith(b'\n'), 'incomplete canonical history tail')
    return [json.loads(line, object_pairs_hook=unique_pairs)
            for line in data.decode('utf-8').splitlines() if line.strip()]


def portable(path):
    require(isinstance(path, str) and path and '\\' not in path
            and not any(ord(c) < 32 for c in path) and ':' not in path
            and not path.startswith('/') and str(PurePosixPath(path)) == path,
            f'invalid Activity path: {path!r}')
    for part in path.split('/'):
        require(part not in ('.', '..') and not part.endswith((' ', '.'))
                and not re.fullmatch(r'(CON|PRN|AUX|NUL|COM[1-9]|LPT[1-9])', part.split('.')[0], re.I),
                f'invalid Activity path component: {part!r}')


def record(sdp, path):
    portable(path)
    target = sdp / path
    for part in [target, *target.parents]:
        require(not part.is_symlink(), f'symlink in Activity record: {path}')
        if part == sdp:
            break
    require(target.is_file(), f'Activity record unavailable: {path}')
    require(target.stat().st_size <= LIMIT, f'Activity record too large: {path}')
    return target.read_bytes()


def sha(data):
    return hashlib.sha256(data).hexdigest()


def metadata(data):
    return dict((key.strip(), value.strip()) for key, value in
                re.findall(r'^\|\s*([^|]+)\|\s*([^|]*)\|\s*$', data.decode().strip(), re.M))


def validate(events, sdp, schema_root=None):
    """Validate all Activity domains; return replayed records and policy presence."""
    sdp = Path(sdp).absolute()
    schemas = Path(schema_root or sdp / 'ProjectManagement')
    validators = {}

    def schema(name, value):
        if name not in validators:
            schema_path = schemas / name
            if not schema_path.exists() and name == 'ledger-event.schema.json':
                schema_path = Path(__file__).resolve().parents[2] / 'Toolkit/schemas/ledger-event.schema.json'
            validators[name] = Draft202012Validator(json.loads(schema_path.read_text()),
                                                    format_checker=FormatChecker())
        validators[name].validate(value)

    last, all_ids, operations, sessions, adoption = {}, set(), set(), set(), None
    namespace = None
    activity_events = [e for e in events if e['eventType'].startswith('x-activity:')]
    index = 0
    while index < len(activity_events):
        first = activity_events[index]
        batch = first['payload'].get('batch')
        if batch:
            require(batch == {'index': 0, 'count': 2} and index + 1 < len(activity_events), 'incomplete Activity switch')
            second = activity_events[index + 1]
            require(second['payload'].get('batch') == {'index': 1, 'count': 2}
                    and second['payload']['operationId'] == first['payload']['operationId']
                    and second['subjectId'] != first['subjectId']
                    and second['payload'].get('requestDigest') == first['payload'].get('requestDigest')
                    and first['payload']['from'] == 'active' and first['payload']['to'] == 'on-hold'
                    and second['payload']['to'] == 'active', 'invalid Activity switch batch')
            # A batch must also be adjacent in the canonical mixed-domain stream.
            position = events.index(first)
            require(position + 1 < len(events) and events[position + 1] == second, 'interleaved Activity switch')
            index += 2
        else:
            index += 1

    for event in events:
        schema('ledger-event.schema.json', event)
        require(event['eventId'] not in all_ids, 'duplicate canonical event ID')
        all_ids.add(event['eventId'])
        kind, payload = event['eventType'], event['payload']
        if kind.startswith('x-activity-admission:'):
            require(kind == 'x-activity-admission:adopted' and adoption is None, 'unknown/duplicate Activity admission')
            schema('activity-admission-payload-0.1.schema.json', payload)
            require(event['subjectId'] == 'ACTIVITY-ADMISSION-' + payload['projectId'], 'admission project mismatch')
            portable(payload['authority']['source'])
            adoption = event
            continue
        if not kind.startswith('x-activity:'):
            continue
        version = payload.get('schemaVersion')
        require(version in ('0.1', '0.2'), 'unsupported Activity payload')
        schema(f'activity-payload-{version}.schema.json', payload)
        project = payload['projectId']
        require(namespace in (None, project), 'ambiguous Activity project')
        namespace = project
        suffix = r'\d{6,}' if version == '0.1' else r'(?:\d{6,}|[A-F0-9]{32}(?:-TARGET)?)'
        require(re.fullmatch(r'EVT-ACTIVITY-' + re.escape(project) + '-' + suffix, event['eventId']), 'Activity event identity mismatch')
        operation = payload['operationId']
        require(operation not in operations or payload.get('batch', {}).get('index') == 1, 'Activity operation identity reused')
        operations.add(operation)
        portable(payload['authority']['source'])
        pin = payload['descriptor']
        path, expected = pin['path'], pin['sha256']
        portable(path)
        require(PurePosixPath(path).name == 'activity.yaml' and str(PurePosixPath(path).parent.parent) == 'Activities', 'invalid Activity home')
        live = record(sdp, path)
        data = live if sha(live) == expected else record(sdp, 'Activities/.descriptors/' + expected + '.yaml')
        require(sha(data) == expected, 'descriptor pin mismatch')
        descriptor = yaml.load(data, Loader=UniqueYAML)
        descriptor_version = descriptor.get('schema', '').removeprefix('sdp-activity/')
        require(descriptor_version in ('0.1', '0.2'), 'unsupported Activity descriptor')
        schema(f'activity-descriptor-{descriptor_version}.schema.json', descriptor)
        ident = descriptor['id']
        require(ident == event['subjectId'] and descriptor['projectId'] == project
                and re.fullmatch('ACTIVITY-' + re.escape(project) + r'-\d{4,}', ident), 'Activity subject/project mismatch')
        require(PurePosixPath(path).parent.name.split('--')[0] == 'Activity-' + ident.rsplit('-', 1)[1], 'Activity folder/identity mismatch')
        portable(descriptor['mandate'])
        for scope in descriptor.get('scope', []):
            portable(scope)
            require(not any(c in scope for c in '*?[]'), 'scope must use literal prefixes')
        prior = last.get(ident)
        before, after = payload['from'], payload['to']
        if kind == 'x-activity:created':
            require(prior is None and payload['previousEventId'] is None and before is None and after == 'planned' and not payload.get('session')
                    and not payload.get('batch'), 'invalid Activity creation')
        else:
            require(prior and payload['previousEventId'] == prior['event']['eventId']
                    and before == prior['state'] and prior['path'] == path, 'broken Activity predecessor')
            if kind == 'x-activity:transitioned':
                require(after in TRANSITIONS[before] and not payload.get('session'), 'invalid Activity lifecycle transition')
                if version == '0.2' and after == 'on-hold':
                    require(bool(payload.get('resumeCondition', '').strip()), 'hold lacks resume condition')
                if version == '0.2' and after == 'completed':
                    require(bool(payload.get('outcome', '').strip()) and payload.get('evidence'), 'completion lacks outcome/evidence')
            elif kind == 'x-activity:updated':
                require(version == '0.2' and before == after and not payload.get('session'), 'metadata changed lifecycle state')
            elif kind == 'x-activity:session-created':
                session = payload.get('session')
                require(version == '0.2' and before == after and session and session['id'] not in sessions, 'invalid Session creation')
                require(session['path'].startswith(str(PurePosixPath(path).parent) + '/Sessions/')
                        and re.fullmatch('SESSION-' + re.escape(project) + r'-\d{4,}', session['id']), 'Session outside Activity/project root')
                record(sdp, session['path'])
                require(sha(record(sdp, 'Activities/.sessions/' + session['sha256'] + '.md')) == session['sha256'], 'Session creation snapshot mismatch')
                sessions.add(session['id'])
            else:
                raise ValueError('unsupported Activity event: ' + kind)
        selected = {r['id'] for r in descriptor['references'] if r.get('membership')}
        if version == '0.2' and kind == 'x-activity:transitioned' and after in ('completed', 'canceled'):
            dispositions = payload.get('contributionDispositions', [])
            require({d['referenceId'] for d in dispositions} == selected
                    and len(dispositions) == len(selected), 'selected contribution lacks disposition')
            require(after == 'canceled' or all(d['disposition'] in ('completed', 'not-applicable') for d in dispositions), 'unmet mandatory contribution')
        completion = event['eventId'] if kind == 'x-activity:transitioned' and after == 'completed' else (prior or {}).get('completion')
        last[ident] = {'event': event, 'state': after, 'descriptor': descriptor, 'path': path, 'completion': completion}
        if not payload.get('batch') or payload['batch']['index'] == 1:
            require(sum(a['state'] == 'active' for a in last.values()) <= 1, 'multiple active Activities')

    live_ids = set()
    for path in (sdp / 'Activities').glob('*/activity.yaml'):
        require(not path.is_symlink(), 'linked Activity descriptor')
        descriptor = yaml.load(record(sdp, path.relative_to(sdp).as_posix()), Loader=UniqueYAML)
        ident = descriptor.get('id')
        require(ident in last and ident not in live_ids, 'unregistered/duplicate Activity home')
        live_ids.add(ident)
        current = last[ident]
        require(path.relative_to(sdp).as_posix() == current['path'] and sha(path.read_bytes()) == current['event']['payload']['descriptor']['sha256'], 'live descriptor differs from current history')
        record(sdp, str(PurePosixPath(current['path']).parent / descriptor['mandate']))
        for ref in descriptor['references']:
            portable(ref['path'])
            if ref.get('availability') == 'unavailable':
                continue
            require(ref.get('projectId', namespace) == namespace, 'external reference needs explicit unavailable disposition')
            content = record(sdp, ref['path'])
            require(ref['id'] in content.decode(), 'typed reference identity mismatch')
            if ref.get('revision'):
                require(sha(content) == ref['revision'], 'reviewed reference revision changed')
    require(live_ids == set(last), 'registered Activity home missing')
    graph = {ident: [d['activityId'] for d in a['descriptor'].get('dependencies', [])] for ident, a in last.items()}
    visited, walking = set(), set()
    def visit(ident):
        require(ident in graph and ident not in walking, 'missing/cyclic Activity dependency')
        if ident in visited:
            return
        walking.add(ident)
        for dependency in graph[ident]:
            visit(dependency)
        walking.remove(ident)
        visited.add(ident)
    for ident, activity in last.items():
        visit(ident)
        if activity['state'] == 'active':
            for dependency in activity['descriptor'].get('dependencies', []):
                target = last[dependency['activityId']]
                require(target['state'] == 'completed' and dependency.get('satisfiedByEventId') == target['completion'], 'stale/unsatisfied Activity dependency')
    if adoption:
        p = adoption['payload']
        require(namespace in (None, p['projectId']), 'admission namespace mismatch')
        stores = set()
        for path in (sdp / 'Activities/.admission/stores').glob('*.json'):
            value = json.loads(record(sdp, path.relative_to(sdp).as_posix()), object_pairs_hook=unique_pairs)
            require(set(value) == {'statePath'} and Path(value['statePath']).is_absolute()
                    and posixpath.normpath(value['statePath']) == value['statePath']
                    and not value['statePath'].startswith('//'), 'invalid execution store locator')
            require(path.name == sha(value['statePath'].encode()) + '.json' and value['statePath'] not in stores, 'duplicate/mismatched store locator')
            stores.add(value['statePath'])
        require({s['statePath'] for s in p['stores']} <= stores and len(stores) <= 128, 'adopted inventory missing/too large')
        pending = sdp / 'Activities/.admission/pending.json'
        if pending.exists():
            value = json.loads(record(sdp, pending.relative_to(sdp).as_posix()), object_pairs_hook=unique_pairs)
            require(set(value) <= {'operationId', 'activityId', 'targetActivityId', 'action'}
                    and {'operationId', 'activityId', 'action'} <= set(value)
                    and value['activityId'] in last and value['action'], 'invalid pending lifecycle')
    return last, adoption is not None


def main():
    import argparse
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('sdp', type=Path)
    args = parser.parse_args()
    records, adopted = validate(load_events(args.sdp / 'ProjectManagement/Ledger.ndjson'), args.sdp, args.sdp / 'ProjectManagement')
    print(f'PASS: {len(records)} Activities; admission adopted={adopted}; read-only records/replay')


if __name__ == '__main__':
    main()
