from pathlib import Path
import sys,subprocess,json,os,tempfile,hashlib
binary,target,mode,out=sys.argv[1:]
base=Path(tempfile.mkdtemp(prefix='sdp-22-trial-'))
env={**os.environ,'SDP_CACHE_DIR':str(base/'cache')}
for name in ['SDP_RELEASE','SDP_TEST_KEY','SDP_OFFLINE']:env.pop(name,None)
results=[]
def run(root,op,*args):
 p=subprocess.run([binary,str(root),op,*args,'--json'],env=env,text=True,capture_output=True)
 if p.returncode:raise RuntimeError(p.stdout+p.stderr)
 return json.loads(p.stdout)['result']
opts=['--release',target] if mode=='signed' else ['--artifact',target,'--allow-unreleased']
for version in ['0.2.0','0.2.1','1.0.0','2.0.0','2.1.0']:
 for custom in [False,True]:
  root=base/(version+('-custom' if custom else '-fresh'));root.mkdir()
  oldpath=Path('/tmp/sdp22-predecessors')/version/'sdp-release.json'
  ip=base/(root.name+'-install.json')
  old=run(root,'install','--release',str(oldpath),'--plan-output',str(ip));run(root,'install','--apply',str(ip))
  original=b'legacy navigation is inert even if invalid JSON\n'
  (root/'SDP/navigation.json').write_bytes(original)
  if custom:
   (root/'SDP/Sessions').mkdir(exist_ok=True);(root/'SDP/Sessions/README.md').write_text('owner guide\n');(root/'SDP/Sessions/session-#0001--Owner.md').write_text('owner session\n')
  program=b'{"schemaVersion":"sdp-programs/1","programs":[]}\n'
  (root/'SDP/programs.json').write_bytes(program)
  (root/'SDP/Owner.sdui').write_text('sdui 0.2; Page=[<"Owner">];\n')
  prefix=(root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes()
  up=base/(root.name+'-upgrade.json');plan=run(root,'upgrade',*opts,'--plan-output',str(up));result=run(root,'upgrade','--apply',str(up))
  assert (root/'SDP/navigation.json').read_bytes()==original
  assert (root/'SDP/programs.json').read_bytes()==program
  assert (root/'SDP/Owner.sdui').read_text()=='sdui 0.2; Page=[<"Owner">];\n'
  assert (root/'SDP/Sessions/Session-template.md').exists()
  assert (root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes().startswith(prefix)
  if custom:
   assert (root/'SDP/Sessions/README.md').read_text()=='owner guide\n'
   assert (root/'SDP/Sessions/session-#0001--Owner.md').read_text()=='owner session\n'
  repeat=run(root,'upgrade',*opts);assert repeat['noChange'] and not repeat['actions']
  receipt=json.loads((root/'SDP/Framework/installed-toolkit.manifest.yaml').read_text());assert 'sdp.sessions.manual.v1' in receipt['capabilities']
  discovery=subprocess.run([binary,str(root),'discover','--json'],env=env,text=True,capture_output=True);assert discovery.returncode==0,discovery.stderr
  if mode=='signed':assert receipt['release']=='2.2.0' and receipt['provenance']=='signed'
  results.append({'previous':version,'custom':custom,'previousDigest':old['release']['sha256'],'previousProvenance':old['release']['provenance'],'status':result['status'],'actions':len(plan['actions']),'preservedLegacyNavigation':True,'preservedSessionsAndLedger':True,'preservedProgramsAndSDUI':True,'repeatNoChange':True,'receipt':receipt})
# New installation must not create a registry.
root=base/'new';root.mkdir();ip=base/'new-install.json';run(root,'install',*opts,'--plan-output',str(ip));run(root,'install','--apply',str(ip));assert not (root/'SDP/navigation.json').exists()
proof={'binarySHA256':hashlib.sha256(Path(binary).read_bytes()).hexdigest(),'descriptorSHA256':hashlib.sha256(Path(target).read_bytes()).hexdigest(),'targetProvenance':mode,'newInstallWithoutRegistry':True,'cases':results}
Path(out).write_text(json.dumps(proof,indent=2)+'\n');print('PASS ten predecessor upgrade/preservation/no-op cases and clean install',out)
