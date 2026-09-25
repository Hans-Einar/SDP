#!/usr/bin/env python3
"""Positive local workflow and corruption checks; all created records are fictional."""
import copy
import hashlib
import json
import tempfile
import unittest
from pathlib import Path
import validate as v

TIME = '2026-09-25T00:00:00Z'


def card(n):
    return {'schemaVersion': '1.0', 'eventId': f'EVT-KB-SDP-{n:06}', 'eventType': 'x-kanban:created', 'occurredAt': TIME, 'actor': 'fixture', 'commit': None, 'subjectId': f'KB-SDP-{n:03}', 'payload': {'schemaVersion': '0.2', 'projectId': 'SDP', 'previousEventId': None, 'from': None, 'to': 'backlog', 'fromPath': None, 'toPath': f'backlog/#{n:03}--Example.md', 'reason': 'Fictional example', 'links': []}}


def work(n, ident, kind, links, state='active', prev=None, action='created', members=None):
    old = prev['payload'] if prev else None
    event = {'schemaVersion': '1.0', 'eventId': f'EVT-PM-SDP-{n:06}', 'eventType': 'x-management:'+action, 'occurredAt': TIME, 'actor': 'fixture', 'commit': None, 'subjectId': ident, 'payload': {'schemaVersion': '0.1', 'kind': kind, 'previousEventId': prev['eventId'] if prev else None, 'from': old['to'] if old else None, 'to': state, 'fromPath': old['toPath'] if old else None, 'toPath': (old['toPath'] if old else kind+'/'+ident+'.md'), 'reason': 'Fictional workflow', 'links': links}}
    if members is not None:
        event['payload']['members'] = members
    return event


class Management(unittest.TestCase):
    def fixture(self):
        cards = [card(1), card(2)]
        scrum = work(1, 'SCRUM-SDP-0001', 'Scrum', ['KB-SDP-001', 'KB-SDP-002'])
        maintenance = work(2, 'MAINT-SDP-0001', 'Maintenance', ['SCRUM-SDP-0001'])
        end = work(3, maintenance['subjectId'], 'Maintenance', ['SCRUM-SDP-0001'], 'completed', maintenance, 'completed')
        return cards + [scrum, maintenance, end]

    def test_supported_origins_and_work_kinds(self):
        v.history(self.fixture())  # multiple cards -> Scrum -> Maintenance
        v.history([work(1, 'SCRUM-SDP-0001', 'Scrum', []), work(2, 'MAINT-SDP-0001', 'Maintenance', ['SCRUM-SDP-0001'])])
        v.history([card(1), work(1, 'MAINT-SDP-0001', 'Maintenance', ['KB-SDP-001'])])
        for kind, prefix in [('CodeReview', 'REVIEW'), ('Refactor', 'REFACTOR')]:
            v.history([card(1), work(1, prefix+'-SDP-0001', kind, ['KB-SDP-001'])])

    def test_invalid_histories(self):
        cases = []
        def change(label, edit):
            events = self.fixture(); edit(events); cases.append((label, events))
        change('duplicate ID', lambda es: es.append(es[-1]))
        change('missing predecessor', lambda es: es[-1]['payload'].update(previousEventId='EVT-PM-SDP-999999'))
        change('wrong origin', lambda es: es[-1]['payload'].update(fromPath='wrong.md'))
        change('wrong kind', lambda es: es[-1]['payload'].update(kind='Sprint'))
        change('missing link', lambda es: es[-1]['payload'].update(links=['KB-SDP-999']))
        change('traversal', lambda es: es[-1]['payload'].update(toPath='../outside.md'))
        change('absolute path', lambda es: es[-1]['payload'].update(toPath='/outside.md'))
        change('time reversal', lambda es: es[-1].update(occurredAt='2020-01-01T00:00:00Z'))
        change('unknown event', lambda es: es[-1].update(eventType='x-management:magic'))
        change('wrong completion', lambda es: es[-1]['payload'].update(to='planned'))
        change('no closure evidence', lambda es: es[-1]['payload'].update(links=[]))
        change('recreated subject', lambda es: es[-1].update(eventType='x-management:created'))
        for label, events in cases:
            with self.subTest(label=label), self.assertRaises(Exception):
                v.history(events)

    def test_sprint_membership_and_start_invariant(self):
        c = card(1)
        sprint = work(1, 'SPR-SDP-0001', 'Sprint', ['KB-SDP-001'], 'planned', members=['KB-SDP-001'])
        es = [c, sprint]
        v.history(es)
        with tempfile.TemporaryDirectory() as tmp:
            sdp = Path(tmp); board = sdp/'KanBan'; (board/'backlog').mkdir(parents=True)
            record = sdp/sprint['payload']['toPath'];record.parent.mkdir()
            path = board/c['payload']['toPath']
            def fill(state='planned', tag='SPR-SDP-0001', members='KB-SDP-001', cardstate='backlog'):
                path.write_text(f'# Card\n\n| Field | Value |\n| --- | --- |\n| id | KB-SDP-001 |\n| project | SDP |\n| CardState | {cardstate} |\n| SprintId | {tag} |\n')
                record.write_text(f'# Sprint\n\n| Field | Value |\n| --- | --- |\n| id | SPR-SDP-0001 |\n| state | {state} |\n| Members | {members} |\n')
                (board/'README.md').write_text('| KB-SDP-001 | Change | backlog | example |\n')
            fill();v.validate_current(es,board,sdp)
            fill(tag='SPR-SDP-9999')
            with self.assertRaises(ValueError):v.validate_current(es,board,sdp)
            fill(members='')
            with self.assertRaises(ValueError):v.validate_current(es,board,sdp)
            started=work(2,'SPR-SDP-0001','Sprint',['KB-SDP-001'],'active',sprint,'started',['KB-SDP-001'])
            fill(state='active')
            with self.assertRaisesRegex(ValueError,'backlog members'):v.validate_current(es+[started],board,sdp)
            (board/'active').mkdir();new=board/'active/#001--Example.md';path.rename(new);path=new
            moved=copy.deepcopy(c);moved.update(eventId='EVT-KB-SDP-000002',eventType='x-kanban:moved')
            moved['payload'].update(previousEventId=c['eventId'], **{'from':'backlog','fromPath':c['payload']['toPath'],'to':'active','toPath':'active/#001--Example.md'})
            fill(state='active',cardstate='ready');(board/'README.md').write_text('| KB-SDP-001 | Change | active | example |\n')
            v.validate_current(es+[moved,started],board,sdp)

    def test_typed_plan_and_plan_only_sprint(self):
        plan=work(1,'PLAN-SDP-0001','Plan',[],state='planned')
        plan['payload'].update(schemaVersion='0.2',planType='ImplementationPlan')
        sprint=work(2,'SPR-SDP-0001','Sprint',['PLAN-SDP-0001'],'planned',members=[])
        sprint['payload'].update(schemaVersion='0.2',plans=['PLAN-SDP-0001'])
        with tempfile.TemporaryDirectory() as tmp:
            sdp=Path(tmp);board=sdp/'KanBan';board.mkdir();(board/'README.md').write_text('')
            def record(event, extra):
                path=sdp/event['payload']['toPath'];path.parent.mkdir(exist_ok=True)
                path.write_text('# Work\n\n| Field | Value |\n| --- | --- |\n| id | '+event['subjectId']+' |\n| state | '+event['payload']['to']+' |\n'+extra)
            record(plan,'| PlanType | ImplementationPlan |\n| BranchPolicy | phase |\n| CommitPolicy | milestone |\n| SprintId | SPR-SDP-0001 |\n')
            record(sprint,'| Members | |\n| Plans | PLAN-SDP-0001 |\n')
            v.validate_current([plan,sprint],board,sdp)
            path=sdp/plan['payload']['toPath'];original=path.read_text()
            for before,after in [('ImplementationPlan','DesignPlan'),('phase','magic'),('SPR-SDP-0001','SPR-SDP-9999')]:
                path.write_text(original.replace(before,after))
                with self.assertRaises(ValueError):v.validate_current([plan,sprint],board,sdp)
            path.write_text(original)
        start=work(3,'SPR-SDP-0001','Sprint',['PLAN-SDP-0001'],'active',sprint,'started',[])
        start['payload'].update(schemaVersion='0.2',plans=['PLAN-SDP-0001'])
        with self.assertRaisesRegex(ValueError,'planned plans'):v.history([plan,sprint,start])
        begun=work(3,'PLAN-SDP-0001','Plan',[],'active',plan,'started')
        begun['payload'].update(schemaVersion='0.2',planType='ImplementationPlan')
        start['eventId']='EVT-PM-SDP-000004'
        v.history([plan,sprint,begun,start])
        closed=work(5,'SPR-SDP-0001','Sprint',['PLAN-SDP-0001'],'completed',start,'completed',[])
        closed['payload'].update(schemaVersion='0.2',plans=['PLAN-SDP-0001'])
        with self.assertRaisesRegex(ValueError,'unfinished plans'):v.history([plan,sprint,begun,start,closed])
        removed=work(5,'SPR-SDP-0001','Sprint',['PLAN-SDP-0001'],'active',start,'updated',[])
        removed['payload'].update(schemaVersion='0.2',plans=[],reason='Defer the sole plan to later work; retain its identity.')
        end=work(6,'SPR-SDP-0001','Sprint',['PLAN-SDP-0001'],'completed',removed,'completed',[])
        end['payload'].update(schemaVersion='0.2',plans=[])
        v.history([plan,sprint,begun,start,removed,end])

    def test_typed_plan_corruption_and_legacy_adoption(self):
        old=work(1,'MAINT-SDP-0001','Maintenance',[])
        update=work(2,'MAINT-SDP-0001','Maintenance',[],prev=old,action='updated')
        update['payload'].update(schemaVersion='0.2',planType='MaintenancePlan')
        v.history([old,update])
        for key,value in [('planType','ArchitecturePlan'),('schemaVersion','9.0'),('plans',[]),('members',[])]:
            bad=copy.deepcopy(update);bad['payload'][key]=value
            with self.subTest(key=key),self.assertRaises(Exception):v.history([old,bad])
        for type_name in ['MaintenancePlan','RequirementPlan','ArchitecturePlan','DesignPlan','ImplementationPlan','VerificationPlan']:
            plan=work(1,'PLAN-SDP-0001','Plan',[]);plan['payload'].update(schemaVersion='0.2',planType=type_name)
            v.history([plan])
            changed=work(2,'PLAN-SDP-0001','Plan',[],prev=plan,action='updated')
            changed['payload'].update(schemaVersion='0.2',planType='DesignPlan' if type_name!='DesignPlan' else 'RequirementPlan')
            with self.assertRaisesRegex(ValueError,'type changed'):v.history([plan,changed])
        empty=work(1,'SPR-SDP-0001','Sprint',[],members=[]);empty['payload'].update(schemaVersion='0.2',plans=[])
        with self.assertRaises(Exception):v.history([empty])

    def test_actual_import_is_exact_and_read_only(self):
        manifest=json.loads((v.AREA/'History/import.json').read_text())
        ledger=(v.AREA/'Ledger.ndjson').read_bytes()
        for item in manifest['imports']:
            data=(v.ROOT/item['archive']).read_bytes()
            self.assertEqual(hashlib.sha256(data).hexdigest(),item['sha256'])
            self.assertEqual(ledger[item['offset']:item['offset']+item['bytes']],data)
            self.assertFalse((v.ROOT/item['source']).exists())
        self.assertEqual(sum(i['events'] for i in manifest['imports']),113)


if __name__ == '__main__':
    unittest.main(verbosity=2)
