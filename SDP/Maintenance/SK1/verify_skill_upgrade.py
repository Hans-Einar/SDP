#!/usr/bin/env python3
"""Exercise actual PowerShell clean/repeat/legacy skill upgrade without live installs."""
import hashlib,json,pathlib,subprocess,sys,tempfile,yaml
ROOT=pathlib.Path(__file__).resolve().parents[3]
PWSH=sys.argv[1]
INSTALL=ROOT/'Toolkit/scripts/Install-SDP.ps1'
def run(p,*args,ok=True):
 r=subprocess.run([PWSH,'-NoProfile','-File',str(INSTALL),'-ProjectRoot',str(p),*args],text=True,capture_output=True)
 if ok and r.returncode:raise AssertionError(r.stdout+r.stderr)
 if not ok:assert r.returncode!=0,(r.stdout,r.stderr)
 return r

def fingerprint(p):return {str(f.relative_to(p)):hashlib.sha256(f.read_bytes()).hexdigest() for f in p.rglob('*') if f.is_file()}
with tempfile.TemporaryDirectory(prefix='sk1-upgrade-') as d:
 p=pathlib.Path(d)/'project';p.mkdir()
 run(p)
 before=fingerprint(p);run(p);assert before==fingerprint(p),'repeat changed files'
 # Restore real prior skill bytes/facts to exercise the version transition.
 facts=p/'SDP/Framework/installed-toolkit.manifest.yaml';current=yaml.safe_load(facts.read_text())
 old=yaml.safe_load(subprocess.check_output(['git','show','aabb359:SDP.manifest.yaml'],cwd=ROOT,text=True))
 old_contract=json.loads(subprocess.check_output(['git','show','aabb359:Toolkit/SDP-install.manifest.json'],cwd=ROOT,text=True))
 current_contract=json.loads((ROOT/'Toolkit/SDP-install.manifest.json').read_text())
 old_destinations={e['destination'] for e in old_contract['entries']}
 for e in current_contract['entries']:
  target=p/e['destination']
  if e['destination'] not in old_destinations and target.is_file():target.unlink()
 for e in old_contract['entries']:
  if e['kind']=='copied' and e['ownership']=='toolkit-managed':
   target=p/e['destination'];target.parent.mkdir(parents=True,exist_ok=True)
   target.write_bytes(subprocess.check_output(['git','show','aabb359:'+e['source']],cwd=ROOT))
 assert not (p/'.codex/skills/sdp/SKILL.md').exists()
 assert 'If not explicitly delegated' in (p/'AGENTS.md').read_text()
 current['skills']=old['skills'];current['agentsContractVersion']='1.0.0';current['capabilities']=old['capabilities']
 facts.write_text(''.join(f'{k}: {json.dumps(v)}\n' for k,v in current.items() if k not in ('skills','capabilities'))+'skills:\n'+''.join(f'  {k}: {json.dumps(v)}\n' for k,v in current['skills'].items())+'capabilities:\n'+''.join(f'  - {v}\n' for v in current['capabilities']))
 for n in old['skills']:
  target=p/'.codex/skills'/n/'SKILL.md'
  target.write_bytes(subprocess.check_output(['git','show',f'aabb359:Toolkit/skills/{n}/SKILL.md'],cwd=ROOT))
 changed=p/'.codex/skills/sdp-worker/SKILL.md';changed.write_text(changed.read_text()+'\nLocal customization to preserve in backup.\n')
 oldbytes=changed.read_bytes();owned=p/'owner.txt';owned.write_text('Owner content.\n')
 before=fingerprint(p);result=run(p,ok=False);assert before==fingerprint(p),'unforced failure mutated tree'
 assert 'ForceManagedFiles' in result.stderr+result.stdout,result.stderr+result.stdout
 plan=json.loads(run(p,'-ForceManagedFiles','-PlanJson').stdout);assert plan['canApply'];assert before==fingerprint(p),'plan wrote files'
 run(p,'-ForceManagedFiles')
 assert owned.read_text()=='Owner content.\n'
 assert any(f.read_bytes()==oldbytes for f in (p/'SDP/.sdp-backups').rglob('SKILL.md')),'old customization not backed up'
 updated=yaml.safe_load(facts.read_text());expected=yaml.safe_load((ROOT/'SDP.manifest.yaml').read_text())
 assert updated['skills']==expected['skills']
 assert (p/'AGENTS.md').read_bytes()==(ROOT/'Toolkit/payload/project-root/AGENTS.md.template').read_bytes()
 for src in (ROOT/'Skills').rglob('*.md'):
  if src.name=='README.md':continue
  assert (p/'.codex/skills'/src.relative_to(ROOT/'Skills')).read_bytes()==src.read_bytes()
 before=fingerprint(p);run(p);assert before==fingerprint(p),'post-upgrade repeat changed files'
 # Existing managed-path link must never write into an external target.
 linkproject=pathlib.Path(d)/'link-project';linkproject.mkdir();external=pathlib.Path(d)/'external';external.mkdir()
 (linkproject/'.codex').symlink_to(external,target_is_directory=True)
 before=fingerprint(linkproject);run(linkproject,ok=False)
 assert list(external.iterdir())==[] and before==fingerprint(linkproject),'unsafe link was mutated'
 print('PASS: clean/repeat; unforced version-change failure without writes; forced upgrade/backup; source equality; project content preservation; post-upgrade repeat; external-link rejection')
