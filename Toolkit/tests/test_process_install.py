"""Disposable process-install integration fixtures. Set SDP_TEST_PWSH to PowerShell 7.4+."""
import hashlib
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[2]
PWSH = os.environ.get('SDP_TEST_PWSH') or shutil.which('pwsh')
ARTIFACT = ROOT/'Toolkit/profiles/five-phase.artifact.json'
INSTALLER = ROOT/'Toolkit/scripts/Install-SDP.ps1'

def snapshot(root):
    return {str(p.relative_to(root)): hashlib.sha256(p.read_bytes()).hexdigest()
            for p in root.rglob('*') if p.is_file()}

@unittest.skipUnless(PWSH, 'PowerShell 7.4+ required')
class ProcessInstall(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix='sdp-process-')
        self.root = Path(self.tmp.name)/'project'
        self.root.mkdir()
    def tearDown(self):
        self.tmp.cleanup()
    def put(self, path, data):
        p = self.root/path
        p.parent.mkdir(parents=True, exist_ok=True)
        p.write_text(data)
    def call(self, *args, ok=True, env=None, artifact=ARTIFACT):
        p = subprocess.run([PWSH,'-NoProfile','-File',str(INSTALLER),'-ProjectRoot',str(self.root),
             '-ProfileArtifact',str(artifact),*map(str,args)],capture_output=True,text=True,env=env)
        if ok:
            self.assertEqual(p.returncode, 0, p.stdout+p.stderr)
        else:
            self.assertNotEqual(p.returncode, 0, p.stdout+p.stderr)
        return p
    def plan(self, *args):
        return json.loads(self.call('-PlanJson',*args).stdout)
    def save(self, plan):
        p = Path(self.tmp.name)/'plan.json'; p.write_text(json.dumps(plan)); return p
    def xfmd(self):
        # Shape observed at XFMD cf11709e; neutral fixture IDs/content, no live mutation.
        self.put('SDP/02--Requirements/README.md','# Requirements\n')
        self.put('SDP/Agents/KanBan/board.json',json.dumps(dict(schemaVersion='0.1',projectId='FIXTURE',ledger='Ledger.ndjson')))
        self.put('SDP/Agents/KanBan/Ledger.ndjson',json.dumps(dict(schemaVersion='1.0',eventId='EVT-KB-FIXTURE-000001',eventType='x-kanban:created',subjectId='KB-FIXTURE-001',occurredAt='2026-09-25T10:00:00Z',actor='fixture',commit=None,payload=dict(schemaVersion='0.2',projectId='FIXTURE',previousEventId=None,**{'from':None,'to':'backlog'},fromPath=None,toPath='backlog/#001--Idea.md',reason='Register fixture idea',links=[])))+'\n')
        self.put('SDP/Agents/KanBan/backlog/#001--Idea.md','# Idea\n\n| Field | Value |\n| --- | --- |\n| id | KB-FIXTURE-001 |\n| CardState | backlog |\n\n[Requirements](../../../02--Requirements/README.md)\n')
        self.put('SDP/README.md','[Card](Agents/KanBan/backlog/%23001--Idea.md)\n')
    def test_plan_clean_readonly_deterministic(self):
        before = snapshot(self.root)
        a,b=self.plan(),self.plan()
        self.assertEqual(a,b)
        self.assertEqual(before,snapshot(self.root))
        self.assertTrue(a['canApply'])
        self.assertEqual(a['oldFacts']['toolkitVersion'],'unknown')
    def test_plan_legacy_and_links(self):
        self.put('SDP/03--Requirements/requirements.md','[Study](../02--Study/study.md)\n')
        self.put('SDP/02--Study/study.md','# Study\n')
        plan=self.plan()
        self.assertTrue(plan['canApply'])
        self.assertEqual(plan['baseline'],'legacy-seven-phase')
        writes={x['destination']:x for x in plan['actions'] if x['action']=='write'}
        import base64
        text=base64.b64decode(writes['SDP/02--Requirements/requirements.md']['content']).decode()
        self.assertIn('../Archive/Legacy/02--Study/study.md',text)
    def test_plan_manual_history(self):
        self.xfmd()
        old=(self.root/'SDP/Agents/KanBan/Ledger.ndjson').read_bytes()
        plan=self.plan()
        self.assertTrue(plan['canApply'])
        import base64
        actions={a['destination']:a for a in plan['actions']}
        self.assertEqual(base64.b64decode(actions['SDP/ProjectManagement/Ledger.ndjson']['content']),old)
        text=base64.b64decode(actions['SDP/KanBan/backlog/#001--Idea.md']['content']).decode()
        self.assertIn('../../02--Requirements/README.md',text)
    def test_plan_conflicts_and_links(self):
        self.put('SDP/03--Requirements/a.md','old')
        self.put('SDP/02--Requirements/a.md','new')
        self.assertIn('mixed-phase-layout',self.plan()['conflicts'])
        (self.root/'SDP/link').symlink_to('/tmp')
        self.call('-PlanJson',ok=False)
    def test_corrupt_artifact(self):
        bad=Path(self.tmp.name)/'artifact.json'
        bad.write_bytes(ARTIFACT.read_bytes().replace(b'sdp-five-phase/0.1',b'sdp-five-phase/0.2'))
        self.call('-PlanJson',artifact=bad,ok=False)

    def apply(self, plan=None, *args, **kw):
        return self.call('-ApplyPlan',self.save(plan or self.plan(*args)),*args,**kw)
    def test_apply_repeat_and_history(self):
        result=json.loads(self.apply().stdout)
        self.assertEqual(result['status'],'completed')
        self.assertEqual(len((self.root/'SDP/ProjectManagement/Ledger.ndjson').read_text().splitlines()),2)
        self.assertEqual(self.plan()['actions'],[])
        before=snapshot(self.root)
        self.assertEqual(json.loads(self.apply().stdout)['status'],'no-change')
        self.assertEqual(before,snapshot(self.root))
    def test_apply_manual_preserves_ledger_prefix(self):
        self.xfmd()
        history=(self.root/'SDP/Agents/KanBan/Ledger.ndjson').read_bytes()
        self.apply()
        self.assertTrue((self.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes().startswith(history))
        self.assertFalse((self.root/'SDP/Agents/KanBan/Ledger.ndjson').exists())
        self.assertIn('KanBan/backlog/%23001', (self.root/'SDP/README.md').read_text())
        self.assertEqual(self.plan()['actions'],[])
    def test_target_drift_and_managed_backup(self):
        plan=self.plan()
        self.put('notes.md','owner edit')
        self.apply(plan,ok=False)
        self.assertFalse((self.root/'SDP/Framework/installed-toolkit.manifest.yaml').exists())
        self.apply()
        self.put('.codex/skills/sdp-worker/SKILL.md','custom worker')
        self.assertFalse(self.plan()['canApply'])
        forced=self.plan('-ForceManagedFiles')
        result=json.loads(self.apply(forced,'-ForceManagedFiles').stdout)
        backups=self.root/'SDP/.sdp-operations'/result['operationId']/'backups'
        self.assertIn(b'custom worker',[p.read_bytes() for p in backups.iterdir()])
    def test_resume_and_postfailure_edits(self):
        plan=self.plan()
        env=dict(os.environ, SDP_INSTALL_INTERRUPT='write:0', SDP_INSTALL_HARD_EXIT='1')
        self.assertEqual(self.apply(plan,ok=False,env=env).returncode,97)
        self.assertFalse(self.plan()['canApply'])
        jpath=next((self.root/'SDP/.sdp-operations').glob('*/journal.json'))
        j=json.loads(jpath.read_text())
        self.put(j['steps'][0]['destination'],'later owner edit')
        self.call('-ResumeOperation',j['operationId'],ok=False)
        import base64
        self.put(j['steps'][0]['destination'],base64.b64decode(j['steps'][0]['content']).decode())
        self.call('-ResumeOperation',j['operationId'])
        ledger=(self.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes()
        self.call('-ResumeOperation',j['operationId'])
        self.assertEqual(ledger,(self.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes())
        self.assertEqual(self.plan()['actions'],[])
    def test_concurrent_writer(self):
        import time
        path=self.root/'SDP/.sdp-operations/install.lock'
        path.parent.mkdir(parents=True)
        script=Path(self.tmp.name)/'lock.ps1'
        script.write_text("$s=[IO.FileStream]::new($args[0],[IO.FileMode]::OpenOrCreate,[IO.FileAccess]::ReadWrite,[IO.FileShare]::None); Write-Output 'locked'; Start-Sleep 30; $s.Dispose()")
        lock=subprocess.Popen([PWSH,'-NoProfile','-File',str(script),str(path)],stdout=subprocess.PIPE,text=True)
        try:
            self.assertEqual(lock.stdout.readline().strip(),'locked')
            self.apply(ok=False)
        finally:
            lock.terminate();lock.wait();lock.stdout.close()

    def test_history_rejects_broken_chain(self):
        self.xfmd()
        p=self.root/'SDP/Agents/KanBan/Ledger.ndjson'
        event=json.loads(p.read_text());event['payload']['previousEventId']='EVT-KB-FIXTURE-999999'
        p.write_text(json.dumps(event)+'\n')
        self.assertFalse(self.plan()['canApply'])
    def test_all_markdown_link_forms(self):
        self.put('SDP/03--Requirements/a file.md','# Source\n')
        self.put('SDP/03--Requirements/a(v1).md','# Source\n')
        self.put('docs/guide.md','[inline](../SDP/03--Requirements/a%20file.md "title")\n[angle](<../SDP/03--Requirements/a file.md>)\n[ref]: ../SDP/03--Requirements/a%20file.md "title"\n[balanced](../SDP/03--Requirements/a(v1).md)\n[escaped](../SDP/03--Requirements/a\\(v1\\).md)\n')
        plan=self.plan()
        import base64
        a=next(a for a in plan['actions'] if a['destination']=='docs/guide.md')
        data=base64.b64decode(a['content']).decode()
        self.assertNotIn('03--Requirements',data)
        self.assertEqual(data.count('02--Requirements/a%20file.md'),3)
        self.assertIn('"title"',data)
        self.assertEqual(data.count('02--Requirements/a%28v1%29.md'),2)
    @unittest.skipUnless(os.name=='posix','POSIX special objects')
    def test_special_file_rejected_without_blocking(self):
        self.put('SDP/02--Requirements/README.md','requirements')
        os.mkfifo(self.root/'SDP/fifo')
        p=subprocess.run([PWSH,'-NoProfile','-File',str(INSTALLER),'-ProjectRoot',str(self.root),'-ProfileArtifact',str(ARTIFACT),'-PlanJson'],capture_output=True,text=True,timeout=20)
        self.assertNotEqual(p.returncode,0)

    def test_legacy_versioned_upgrade_and_downgrade(self):
        legacy=subprocess.run([PWSH,'-NoProfile','-File',str(INSTALLER),'-ProjectRoot',str(self.root),'-InitializeProjectStructure'],capture_output=True,text=True)
        self.assertEqual(legacy.returncode,0,legacy.stderr)
        plan=self.plan()
        self.assertEqual(plan['oldFacts']['schemaVersion'],'1.0')
        self.assertTrue(plan['canApply'],plan['conflicts'])
        self.apply(plan)
        old_reader=subprocess.run([PWSH,'-NoProfile','-File',str(INSTALLER),'-ProjectRoot',str(self.root),'-PlanJson'],capture_output=True,text=True)
        self.assertEqual(old_reader.returncode,0,old_reader.stderr)
        blocked=json.loads(old_reader.stdout)
        self.assertFalse(blocked['canApply'])
        self.assertEqual(blocked['actions'][0]['reason'],'unsupported-installed-schema')
        facts=self.root/'SDP/Framework/installed-toolkit.manifest.yaml'
        data=facts.read_text()
        facts.write_text(data.replace('toolkitVersion: "0.2.0"','toolkitVersion: "99.0.0"'))
        self.assertIn('downgrade-blocked',self.plan()['conflicts'])
        facts.write_text(data.replace('schemaVersion: "2.0"','schemaVersion: "99.0"'))
        self.call('-PlanJson',ok=False)

    @unittest.skipUnless(os.environ.get('SDP_TEST_TOOL'),'prebuilt SDPTool required')
    def test_installed_consumer_workflow(self):
        self.apply()
        tool=os.environ['SDP_TEST_TOOL']
        for selected in (self.root,self.root/'SDP'):
            result=subprocess.run([tool,str(selected),'discover'],capture_output=True,text=True)
            self.assertEqual(result.returncode,0,result.stderr)
            p=json.loads(result.stdout)
            self.assertEqual(p['installation']['facts']['schemaVersion'],'2.0')
            self.assertEqual(p['installation']['state'],'declared')
            self.assertEqual(p['registration']['models'],[])
            self.assertEqual(p['registration']['sdui'],[])
        result=subprocess.run([tool,str(self.root),'tree'],capture_output=True,text=True)
        self.assertEqual(result.returncode,0,result.stderr)
        nodes={n['id']:n for n in json.loads(result.stdout)['nodes']}
        self.assertEqual(nodes['kanban']['state'],'validated')
        self.assertEqual(nodes['sdl']['state'],'absent')
        validator=subprocess.run([__import__('sys').executable,str(ROOT/'Toolkit/scripts/validate_sdp.py'),'--mode','project','--project-root',str(self.root)],capture_output=True,text=True)
        self.assertEqual(validator.returncode,0,validator.stdout+validator.stderr)

    def test_case_collision(self):
        self.put('SDP/02--Requirements/A.md','one')
        self.put('SDP/02--Requirements/a.md','two')
        self.call('-PlanJson',ok=False)

    def test_typed_plan_history_and_sprint_rejections(self):
        import copy
        self.put('SDP/02--Requirements/README.md','# Requirements\n')
        plan=dict(schemaVersion='1.0',eventId='EVT-PM-PROJECT-000001',eventType='x-management:created',subjectId='PLAN-PROJECT-0001',occurredAt='2026-09-25T18:00:00Z',actor='fixture',commit=None,payload=dict(schemaVersion='0.2',kind='Plan',planType='DesignPlan',previousEventId=None,**{'from':None,'to':'planned'},fromPath=None,toPath='04--Design/Plan.md',reason='Plan fixture',links=[]))
        text='# Plan\n\n| Field | Value |\n| --- | --- |\n| id | PLAN-PROJECT-0001 |\n| state | planned |\n| PlanType | DesignPlan |\n| BranchPolicy | current |\n| CommitPolicy | phase |\n'
        def put_events(events):self.put('SDP/ProjectManagement/Ledger.ndjson',''.join(json.dumps(e)+'\n' for e in events))
        put_events([plan]);self.put('SDP/04--Design/Plan.md',text)
        self.assertTrue(self.plan()['canApply'])
        update=copy.deepcopy(plan);update.update(eventId='EVT-PM-PROJECT-000002',eventType='x-management:updated')
        update['payload'].update(previousEventId=plan['eventId'],**{'from':'planned'},fromPath=plan['payload']['toPath'],planType='ImplementationPlan')
        put_events([plan,update]);self.put('SDP/04--Design/Plan.md',text.replace('DesignPlan','ImplementationPlan'))
        self.assertFalse(self.plan()['canApply'])
        put_events([plan]);self.put('SDP/04--Design/Plan.md',text.replace('current','unknown'))
        self.assertFalse(self.plan()['canApply'])
        self.put('SDP/04--Design/Plan.md',text+'| SprintId | SPR-PROJECT-0001 |\n')
        sprint=copy.deepcopy(plan);sprint.update(eventId='EVT-PM-PROJECT-000002',subjectId='SPR-PROJECT-0001')
        sprint['payload'].pop('planType');sprint['payload'].update(kind='Sprint',members=[],plans=['PLAN-PROJECT-0001'],links=['PLAN-PROJECT-0001'],toPath='Sprints/One.md')
        put_events([plan,sprint]);self.put('SDP/Sprints/One.md','# Sprint\n\n| Field | Value |\n| --- | --- |\n| id | SPR-PROJECT-0001 |\n| state | planned |\n| Members |  |\n| Plans | PLAN-PROJECT-0001 |\n')
        self.assertTrue(self.plan()['canApply'])
        bad=copy.deepcopy(sprint);bad['payload']['plans']=['PLAN-PROJECT-9999'];put_events([plan,bad])
        self.assertFalse(self.plan()['canApply'])
        put_events([plan,sprint]);self.put('SDP/04--Design/Plan.md',text) # missing reciprocal link
        self.assertFalse(self.plan()['canApply'])
        self.put('SDP/04--Design/Plan.md',text+'| SprintId | SPR-PROJECT-0001 |\n')
        started=copy.deepcopy(sprint);started.update(eventId='EVT-PM-PROJECT-000003',eventType='x-management:started');started['payload'].update(previousEventId=sprint['eventId'],**{'from':'planned','to':'active'},fromPath=sprint['payload']['toPath'])
        put_events([plan,sprint,started]);self.put('SDP/Sprints/One.md',(self.root/'SDP/Sprints/One.md').read_text().replace('| state | planned |','| state | active |'))
        before=snapshot(self.root)
        self.assertFalse(self.plan()['canApply']) # plan has not started
        self.assertEqual(before,snapshot(self.root))
        # All selected work may be deferred explicitly; ordinary empty Markdown rows work.
        plan['payload']['to']='active'
        removed=copy.deepcopy(started);removed.update(eventId='EVT-PM-PROJECT-000004',eventType='x-management:updated')
        removed['payload'].update(previousEventId=started['eventId'],**{'from':'active'},plans=[],reason='Defer the only plan with a preserved reference.')
        closed=copy.deepcopy(removed);closed.update(eventId='EVT-PM-PROJECT-000005',eventType='x-management:completed')
        closed['payload'].update(previousEventId=removed['eventId'],to='completed')
        put_events([plan,sprint,started,removed,closed])
        self.put('SDP/04--Design/Plan.md',text.replace('| state | planned |','| state | active |'))
        self.put('SDP/Sprints/One.md','# Sprint\n\n| Field | Value |\n| --- | --- |\n| id | SPR-PROJECT-0001 |\n| state | completed |\n| Members | |\n| Plans | |\n')
        self.assertTrue(self.plan()['canApply'])

    def test_planning_profile_upgrade(self):
        import yaml
        old=Path(self.tmp.name)/'old.artifact.json'
        old.write_bytes(subprocess.check_output(['git','show','28bf156a02e1dc3e6aa29fd8579a1273bb343a1f:Toolkit/profiles/five-phase.artifact.json'],cwd=ROOT))
        old_plan=json.loads(self.call('-PlanJson',artifact=old).stdout)
        self.call('-ApplyPlan',self.save(old_plan),artifact=old)
        history=(self.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes()
        self.put('SDP/ProjectManagement/README.md','# Owner workflow notes\n')
        self.assertFalse(self.plan()['canApply'])  # managed instructions must be explicitly refreshed
        plan=self.plan('-ForceManagedFiles');self.assertTrue(plan['canApply'],plan['conflicts'])
        self.apply(plan,'-ForceManagedFiles')
        self.assertTrue((self.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes().startswith(history))
        self.assertEqual((self.root/'SDP/ProjectManagement/README.md').read_text(),'# Owner workflow notes\n')
        facts=yaml.safe_load((self.root/'SDP/Framework/installed-toolkit.manifest.yaml').read_text())
        self.assertEqual(facts['managementProfile'],'sdp-project-management/0.2')
        self.assertIn('sdp.planning.v1',facts['capabilities'])
        self.assertEqual(facts['skills']['sdp-planning'],'1.0.0')
        self.assertTrue((self.root/'SDP/Framework/planning/Plans.md').is_file())
        self.assertTrue((self.root/'.codex/skills/sdp-planning/SKILL.md').is_file())
        board=json.loads((self.root/'SDP/KanBan/board.json').read_text())
        self.assertEqual(board['profile'],'sdp-project-management/0.2')
        self.assertEqual(self.plan()['actions'],[])
        downgrade=json.loads(self.call('-PlanJson','-ForceManagedFiles',artifact=old).stdout)
        self.assertFalse(downgrade['canApply'])
        self.assertIn('management-profile-downgrade-blocked',downgrade['conflicts'])
        # New typed events are accepted only through the new explicit payload schema.
        ledger=self.root/'SDP/ProjectManagement/Ledger.ndjson'
        event=dict(schemaVersion='1.0',eventId='EVT-PM-PROJECT-000003',eventType='x-management:created',subjectId='PLAN-PROJECT-0001',occurredAt='2026-09-25T18:00:00Z',actor='fixture',commit=None,payload=dict(schemaVersion='0.2',kind='Plan',planType='DesignPlan',previousEventId=None,**{'from':None,'to':'planned'},fromPath=None,toPath='04--Design/Plan.md',reason='Plan fixture',links=[]))
        # Installer has already allocated four events across the two operations.
        events=[json.loads(x) for x in ledger.read_text().splitlines()]
        n=max(int(e['eventId'].rsplit('-',1)[1]) for e in events)+1
        event['eventId']=f'EVT-PM-PROJECT-{n:06d}'
        with ledger.open('a') as f:f.write(json.dumps(event)+'\n')
        self.put('SDP/04--Design/Plan.md','# Plan\n\n| Field | Value |\n| --- | --- |\n| id | PLAN-PROJECT-0001 |\n| state | planned |\n| PlanType | DesignPlan |\n| BranchPolicy | current |\n| CommitPolicy | phase |\n')
        self.assertTrue(self.plan()['canApply'])
        event['payload']['planType']='InventedPlan'
        ledger.write_text('\n'.join(json.dumps(e) for e in events+[event])+'\n')
        self.assertFalse(self.plan()['canApply'])

    def test_final_artifact_release_notes_recovery(self):
        self.xfmd()
        history=(self.root/'SDP/Agents/KanBan/Ledger.ndjson').read_bytes()
        env=dict(os.environ,SDP_INSTALL_INTERRUPT='prepared:0',SDP_INSTALL_HARD_EXIT='1')
        self.assertEqual(self.apply(ok=False,env=env).returncode,97)
        path=next((self.root/'SDP/.sdp-operations').glob('*/journal.json'))
        journal=json.loads(path.read_text());ident=journal['operationId']
        index=next(i for i,s in enumerate(journal['steps']) if s['destination']=='SDP/RELEASE-NOTES.md')
        for boundary in ('backup','write','journal'):
            env=dict(os.environ,SDP_INSTALL_INTERRUPT=f'{boundary}:{index}',SDP_INSTALL_HARD_EXIT='1')
            self.assertEqual(self.call('-ResumeOperation',ident,env=env,ok=False).returncode,97)
        env=dict(os.environ,SDP_INSTALL_INTERRUPT=f'complete:{len(journal["steps"])}',SDP_INSTALL_HARD_EXIT='1')
        self.assertEqual(self.call('-ResumeOperation',ident,env=env,ok=False).returncode,97)
        self.call('-ResumeOperation',ident)
        final=(self.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes()
        self.assertTrue(final.startswith(history))
        events=[json.loads(x) for x in final.splitlines()]
        self.assertEqual(len(events),3)
        self.assertEqual(len({e['eventId'] for e in events}),3)
        import yaml
        facts=yaml.safe_load((self.root/'SDP/Framework/installed-toolkit.manifest.yaml').read_text())
        self.assertEqual(facts['configurationDigest'],json.loads(ARTIFACT.read_text())['configurationDigest'])
        self.assertEqual(self.plan()['actions'],[])
        # Project release notes are owned by the consuming project, including force mode.
        self.put('SDP/RELEASE-NOTES.md','# Owner release history\n')
        self.assertEqual(self.plan('-ForceManagedFiles')['actions'],[])
        self.assertEqual((self.root/'SDP/RELEASE-NOTES.md').read_text(),'# Owner release history\n')

if __name__ == '__main__':
    unittest.main()
