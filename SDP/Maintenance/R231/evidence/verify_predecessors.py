from pathlib import Path
import subprocess,os,json,hashlib
base=Path('/tmp/sdp-r231-final-predecessors');base.mkdir()
engine=Path('/tmp/sdp-r231-final-package/sdptool'); target=Path('/tmp/sdp-r231-final-package/sdp-release.json')
results=[]
for version in ['0.2.0','0.2.1','1.0.0','2.0.0','2.1.0','2.2.0','2.2.1','2.3.0']:
 case=base/('signed-trial-'+version);case.mkdir()
 root=case/'project';root.mkdir()
 cache=case/'cache'
 old=Path.home()/('.local/state/sdp-release-2.2.1/predecessor-'+version+'/sdp-release.json')
 if version=='2.2.1':old=Path.home()/'.local/state/sdp-rollout-2.2.1/producer-final/sdp-release.json'
 if version=='2.3.0':old=Path('/tmp/activity4-public-canonical/sdp-release.json')
 env=os.environ.copy()
 for k in ['SDP_RELEASE','SDP_TEST_KEY','SDP_OFFLINE','SDP_CACHE_DIR']:env.pop(k,None)
 env['SDP_CACHE_DIR']=str(cache)
 def cmd(name,args,release,expected=0):
  env['SDP_RELEASE']=str(release)
  p=subprocess.run([str(engine),str(root)]+args+['--json'],env=env,capture_output=True,text=True,timeout=120)
  (case/(name+'.json')).write_text(p.stdout);(case/(name+'.err')).write_text(p.stderr)
  assert p.returncode==expected,version+' '+name+' '+p.stderr+p.stdout[:200]
  return p
 cmd('install-preview',['install','--plan-output',str(case/'install-plan.json')],old)
 plan=json.loads((case/'install-plan.json').read_text());assert plan['canApply']
 cmd('install-apply',['install','--apply',str(case/'install-plan.json')],old)
 fixtures={'AGENTS-project.md':'# Owner project instructions\nPreserve this.\n','SDP/Sessions/owner.md':'# Owner session\nKeep exact bytes.\n','SDP/SDL/owner.sdl':'// owner fixture, not executable\n'}
 for rel,data in fixtures.items():
  f=root/rel;f.parent.mkdir(parents=True,exist_ok=True);f.write_text(data)
 # Test managed-file conflict without changing production repositories.
 managed=root/'AGENTS.md'; original=managed.read_bytes();managed.write_bytes(original+b'\nOwner edit\n')
 p=subprocess.run([str(engine),str(root),'upgrade','--release',str(target),'--plan-output',str(case/'conflict-plan.json'),'--json'],env=env,capture_output=True,text=True,timeout=120)
 (case/'conflict-output.json').write_text(p.stdout)
 cp=json.loads((case/'conflict-plan.json').read_text())
 assert not cp['canApply'] and cp['conflicts'], 'managed conflict not blocked'
 managed.write_bytes(original)
 cmd('upgrade-preview',['upgrade','--plan-output',str(case/'upgrade-plan.json')],target)
 plan=json.loads((case/'upgrade-plan.json').read_text());assert plan['canApply'],str(plan['conflicts'])
 cmd('upgrade-apply',['upgrade','--apply',str(case/'upgrade-plan.json')],target)
 for rel,data in fixtures.items():assert (root/rel).read_text()==data,rel
 receipt=json.loads((root/'SDP/Framework/installed-toolkit.manifest.yaml').read_text())
 assert receipt['release']=='2.3.1' and receipt['provenance']=='signed'
 assert receipt['descriptorDigest']==hashlib.sha256(target.read_bytes()).hexdigest()
 cmd('noop',['upgrade','--plan-output',str(case/'noop-plan.json')],target)
 np=json.loads((case/'noop-plan.json').read_text());assert np['noChange'] and not np['actions']
 results.append({'predecessor':version,'target':'2.3.1','provenance':receipt['provenance'],'descriptorDigest':receipt['descriptorDigest'],'actions':len(plan['actions']),'preserved':True,'conflictRefused':True,'repeatNoChange':True})
 (base/'signed-upgrade-results.json').write_text(json.dumps(results,indent=2)+'\n')
 print(version,'signed upgrade/preservation/conflict/no-op passed',flush=True)
