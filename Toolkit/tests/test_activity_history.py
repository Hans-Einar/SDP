"""Read-only schema/replay checks on fictional Activity records and retained bytes."""
import copy
import hashlib
import importlib.util
import json
import tempfile
import unittest
import uuid
from pathlib import Path

import yaml

ROOT = Path(__file__).resolve().parents[2]
SCHEMAS = ROOT / 'Template/sdp-root/ProjectManagement'
spec = importlib.util.spec_from_file_location('activity_history', ROOT / 'Toolkit/scripts/activity_history.py')
v = importlib.util.module_from_spec(spec)
spec.loader.exec_module(v)
TIME = '2026-11-12T00:00:00Z'


def digest(data):
    return hashlib.sha256(data).hexdigest()


class ActivityHistory(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory()
        self.addCleanup(self.tmp.cleanup)
        self.sdp = Path(self.tmp.name) / 'SDP'
        self.sdp.mkdir()
        self.events = []

    def write(self, path, data):
        target = self.sdp / path
        target.parent.mkdir(parents=True, exist_ok=True)
        target.write_bytes(data)
        return target

    def create(self, number):
        home = f'Activities/Activity-{number:04}--Need'
        descriptor = dict(schema='sdp-activity/0.2', id=f'ACTIVITY-TEST-{number:04}',
                          projectId='TEST', title='Need', owner='owner', mandate='01-Mandate/Mandate.md',
                          sourceCards=[], references=[], scope=['Source'])
        data = yaml.safe_dump(descriptor).encode()
        self.write(home + '/activity.yaml', data)
        self.write(home + '/01-Mandate/Mandate.md', b'Need')
        return descriptor['id'], dict(path=home + '/activity.yaml', sha256=digest(data))

    def event(self, record, kind='created', before=None, after='planned', previous=None, **extra):
        number = len(self.events) + 1
        ident, pin = record
        payload = dict(schemaVersion='0.2', projectId='TEST', operationId=str(uuid.UUID(int=number)),
                       recordedAt=TIME, authority=dict(role='project-owner', source='authority.md'),
                       previousEventId=previous, to=after, reason='Owner selects fictional work',
                       descriptor=pin, links=[], requestDigest='a' * 64, **extra)
        payload['from'] = before
        event = dict(schemaVersion='1.0', eventId=f'EVT-ACTIVITY-TEST-{number:06}',
                     eventType='x-activity:' + kind, occurredAt=TIME, actor='fixture', commit=None,
                     subjectId=ident, payload=payload)
        self.events.append(event)
        return event

    def check(self):
        return v.validate(self.events, self.sdp, SCHEMAS)

    def test_schema_replay_and_mutable_session_document(self):
        record = self.create(1)
        created = self.event(record)
        active = self.event(record, 'transitioned', 'planned', 'active', created['eventId'])
        path = str(Path(record[1]['path']).parent / 'Sessions/session.md')
        data = b'SESSION-TEST-0001\nCreated objective\n'
        self.write(path, data)
        self.write('Activities/.sessions/' + digest(data) + '.md', data)
        self.event(record, 'session-created', 'active', 'active', active['eventId'],
                   session=dict(id='SESSION-TEST-0001', path=path, sha256=digest(data)))
        self.write(path, data + b'Actual later work summary\n')
        records, adopted = self.check()
        self.assertEqual(records[record[0]]['state'], 'active')
        self.assertFalse(adopted)
        self.events[-1]['payload']['unexpected'] = True
        with self.assertRaises(Exception):
            self.check()

    def test_switch_is_complete_ordered_and_adjacent(self):
        first, second = self.create(1), self.create(2)
        created = self.event(first)
        active = self.event(first, 'transitioned', 'planned', 'active', created['eventId'])
        planned = self.event(second)
        held = self.event(first, 'transitioned', 'active', 'on-hold', active['eventId'],
                          resumeCondition='Owner resumes', batch=dict(index=0, count=2))
        self.event(second, 'transitioned', 'planned', 'active', planned['eventId'], batch=dict(index=1, count=2))
        self.events[-1]['payload']['operationId'] = held['payload']['operationId']
        self.check()
        self.events.pop()
        with self.assertRaisesRegex(ValueError, 'incomplete Activity switch'):
            self.check()

    def test_stale_descriptor_and_duplicate_keys_fail(self):
        record = self.create(1)
        self.event(record)
        self.check()
        self.write(record[1]['path'], b'schema: unsupported\n')
        with self.assertRaises(Exception):
            self.check()
        path = self.write('ProjectManagement/Ledger.ndjson', b'{"eventId":"one","eventId":"two"}\n')
        with self.assertRaisesRegex(ValueError, 'duplicate JSON key'):
            v.load_events(path)
        path.write_bytes(b'{}')
        with self.assertRaisesRegex(ValueError, 'incomplete'):
            v.load_events(path)

    def test_missing_adopted_registry_fails(self):
        record = self.create(1)
        self.event(record)
        store = '/explicit/known-state'
        event = dict(schemaVersion='1.0', eventId='EVT-ACTIVITY-ADMISSION-000001',
                     eventType='x-activity-admission:adopted', occurredAt=TIME, actor='fixture', commit=None,
                     subjectId='ACTIVITY-ADMISSION-TEST', payload=dict(schemaVersion='0.1', projectId='TEST',
                     operationId=str(uuid.UUID(int=100)), requestDigest='b' * 64,
                     authority=dict(role='project-owner', source='authority.md'), stores=[dict(statePath=store)]))
        self.events.append(event)
        with self.assertRaisesRegex(ValueError, 'adopted inventory missing'):
            self.check()
        self.write('Activities/.admission/stores/' + digest(store.encode()) + '.json', json.dumps(dict(statePath=store)).encode())
        self.assertTrue(self.check()[1])

    def test_existing_manual_schemas_remain_byte_preserved(self):
        for name in ['activity-descriptor-0.1.schema.json', 'activity-payload-0.1.schema.json']:
            self.assertEqual((SCHEMAS / name).read_bytes(), (ROOT / 'SDP/ProjectManagement' / name).read_bytes())
        self.assertEqual((ROOT / 'Toolkit/scripts/activity_history.py').read_bytes(),
                         (ROOT / 'SDP/ProjectManagement/activity_history.py').read_bytes())


if __name__ == '__main__':
    unittest.main()
