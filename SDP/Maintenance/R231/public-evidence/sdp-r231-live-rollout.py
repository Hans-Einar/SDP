from pathlib import Path
import os,json,subprocess,hashlib,shutil
base=Path('/tmp/sdp-r231-live');base.mkdir(exist_ok=False);root=Path('/home/warloc/git/HEOS');expected='04f78532a4e7ae8653df5b1e623e1407312713a90043f2ac36e7f68e0dc0495b'
assert json.load(open('/tmp/sdp-r231-public-framework/results.json'))['target']=='2.3.1'
assert len(json.load(open('/tmp/sdp-r231-public-consumers/results.json'))['outcomes'])==3
correction=root/'SDP/Maintenance/SDP-Upgrade-20261114/evidence/Correction.json';cr=json.load(open(correction));ledger=root/'SDP/ProjectManagement/Ledger.ndjson';assert hashlib.sha256(ledger.read_bytes()).hexdigest()==cr['correctedSHA256']
env={k:v for k,v in os.environ.items() if k not in ['SDP_RELEASE','SDP_TEST_KEY','SDP_OFFLINE','SDP_CACHE_DIR','GH_DEBUG']}
commands=[]
def run(name,args,cwd=root):
 p=subprocess.run(args,env=env,cwd=cwd,capture_output=True,text=True,timeout=240)
 commands.append(dict(name=name,args=[str(x) for x in args],exitCode=p.returncode,stdout=p.stdout,stderr=p.stderr));(base/'commands.json').write_text(json.dumps(commands,indent=2)+'\n');assert p.returncode==0,name+': '+p.stderr+p.stdout[:300];return p.stdout
def snapshot():
 d={}
 for parent,dirs,files in os.walk(root,followlinks=False):
  dirs[:]=[n for n in dirs if n not in ['.git','build','.cache']]
  for name in files:
   p=Path(parent)/name;d[str(p.relative_to(root))]='symlink:'+os.readlink(p) if p.is_symlink() else hashlib.sha256(p.read_bytes()).hexdigest()
 return d
before=snapshot();(base/'before.json').write_text(json.dumps(before,indent=2)+'\n');oldledger=ledger.read_bytes()
old=Path.home()/'.local/share/gh/extensions/gh-sdp';shutil.copytree(old,base/'previous-extension')
run('live-extension-upgrade',['gh','extension','upgrade','sdp'])
setup=json.loads(run('live-setup',['gh','sdp','setup','--json']))
alias=Path.home()/'.local/bin/sdptool';assert alias.is_symlink() and hashlib.sha256(alias.read_bytes()).hexdigest()==expected
version=json.loads(run('gh-version',['gh','sdp','--version','--json']));assert version['version']=='2.3.1' and version['revision']=='84d5e46012fdc45a3690bd0b26b617501fd1caa7'
assert json.loads(run('direct-version',[str(alias),'--version','--json']))==version
assert snapshot()==before,'tool update unexpectedly mutated HEOS'
# Preserve original signed predecessor in the persistent default verified cache.
seed=base/'empty-cache-seed';seed.mkdir();run('original-public-descriptor',['gh','sdp',str(seed),'install','--release','https://github.com/Hans-Einar/SDP/releases/download/v2.2.0/sdp-release.json','--plan-output',str(base/'original-plan.json'),'--json'])
run('live-public-default-preview',['gh','sdp',str(root),'upgrade','--plan-output',str(base/'upgrade-plan.json'),'--json'])
plan=json.load(open(base/'upgrade-plan.json'));assert plan['canApply'] and not plan['conflicts'];assert len(plan['actions'])==14
allowed={f['path'] for f in json.load(open('/tmp/sdp-r231-public/canonical/sdp-release.json'))['files'] if f['ownership']=='managed'}
descriptorFiles={f['path']:f for f in json.load(open('/tmp/sdp-r231-public/canonical/sdp-release.json'))['files']}
assert all(a['action']=='write' and (a['path'] in allowed or (descriptorFiles[a['path']]['ownership']=='initialize-if-missing' and a['path'] not in before)) for a in plan['actions'])
assert snapshot()==before,'live tree changed since preview'
run('live-public-default-apply',['gh','sdp',str(root),'upgrade','--apply',str(base/'upgrade-plan.json'),'--json'])
receipt=json.loads((root/'SDP/Framework/installed-toolkit.manifest.yaml').read_text());assert receipt['release']=='2.3.1' and receipt['provenance']=='signed' and receipt['descriptorDigest']=='9c83225b3a5eae1762e5609965378b8c14d89fa676ae242856b323c1dd9facf6'
after=snapshot();changed=[p for p,h in before.items() if after.get(p)!=h];assert not set(changed)-(allowed|{'SDP/Framework/installed-toolkit.manifest.yaml','SDP/ProjectManagement/Ledger.ndjson'}),changed
assert ledger.read_bytes().startswith(oldledger)
# Check each installed managed file against the actual signed descriptor.
import base64
for f in json.load(open('/tmp/sdp-r231-public/canonical/sdp-release.json'))['files']:
 if f['ownership']=='managed' or f['path'] in {a['path'] for a in plan['actions']}:assert (root/f['path']).read_bytes()==base64.b64decode(f['content']),f['path']
run('live-repeat-preview',['gh','sdp',str(root),'upgrade','--plan-output',str(base/'noop-plan.json'),'--json']);noop=json.load(open(base/'noop-plan.json'));assert noop['noChange'] and not noop['actions']
(base/'receipt.json').write_text(json.dumps(receipt,indent=2)+'\n');(base/'result.json').write_text(json.dumps(dict(receipt=receipt,version=version,registeredPath=str(alias),registeredTarget=str(alias.resolve()),engineSHA256=expected,actions=len(plan['actions']),changedExistingPaths=changed,existingOwnerFilesPreserved=True,allManagedBytesMatchSignedDescriptor=True,correctedHistoryPrefixPreserved=True,repeatNoChange=True,defaultRelease=True,priorUncommittedWorkPreserved=True),indent=2)+'\n')
print('Actual gh extension and HEOS signed2.3.1 upgrade verified; owner files preserved; repeat no-op')
