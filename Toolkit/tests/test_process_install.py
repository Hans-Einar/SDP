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
        self.put('SDP/Agents/KanBan/Ledger.ndjson',json.dumps(dict(schemaVersion='1.0',eventId='EVT-KB-FIXTURE-000001',eventType='x-kanban:created',subjectId='KB-FIXTURE-001',payload=dict(schemaVersion='0.2',previousEventId=None,toPath='backlog/#001--Idea.md')))+'\n')
        self.put('SDP/Agents/KanBan/backlog/#001--Idea.md','# Idea\n\n[Requirements](../../../02--Requirements/README.md)\n')
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

if __name__ == '__main__':
    unittest.main()
