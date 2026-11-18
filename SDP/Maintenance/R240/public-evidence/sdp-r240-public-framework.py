from pathlib import Path
import os,json,subprocess,hashlib
base=Path('/tmp/sdp-r240-public-framework');base.mkdir()
home=Path('/tmp/sdp-r240-public-consumers/clean-install');engine=home/'.local/bin/sdptool';root=base/'project';root.mkdir()
env={k:v for k,v in os.environ.items() if k not in ['SDP_RELEASE','SDP_TEST_KEY','SDP_OFFLINE','SDP_CACHE_DIR','GH_TOKEN','GITHUB_TOKEN','GH_ENTERPRISE_TOKEN','GITHUB_ENTERPRISE_TOKEN','GH_DEBUG']};env['HOME']=str(home);env['SDP_CACHE_DIR']=str(base/'cache')
commands=[]
def run(name,args):
 p=subprocess.run([str(engine),str(root)]+args+['--json'],env=env,capture_output=True,text=True,timeout=180)
 commands.append(dict(name=name,args=args,exitCode=p.returncode,stdout=p.stdout,stderr=p.stderr));(base/'commands.json').write_text(json.dumps(commands,indent=2)+'\n')
 assert p.returncode==0,name+': '+p.stderr
 return json.loads(p.stdout)
run('public-default-install-preview',['install','--plan-output',str(base/'install-plan.json')])
plan=json.loads((base/'install-plan.json').read_text());assert plan['canApply']
run('public-default-install-apply',['install','--apply',str(base/'install-plan.json')])
receipt=json.loads((root/'SDP/Framework/installed-toolkit.manifest.yaml').read_text());assert receipt['release']=='2.4.0' and receipt['provenance']=='signed' and receipt['descriptorDigest']=='7184efb8b5a4b1fa9280f61f5e9c839b4d6dc14ee58173a81478a7eaaa815908'
owner=root/'SDP/owner-file.md';owner.write_bytes(b'# Owner file\nKeep exact bytes.\n');original=owner.read_bytes()
run('public-default-upgrade-noop',['upgrade','--plan-output',str(base/'noop-plan.json')])
noop=json.loads((base/'noop-plan.json').read_text());assert noop['noChange'] and not noop['actions'];assert owner.read_bytes()==original
(base/'receipt.json').write_text(json.dumps(receipt,indent=2)+'\n');(base/'results.json').write_text(json.dumps(dict(target='2.4.0',provenance=receipt['provenance'],descriptorDigest=receipt['descriptorDigest'],publicDefault=True,releaseOverridesAbsent=True,testTrustAbsent=True,previewApply=True,repeatNoChange=True,ownerFilePreserved=True,ownerSHA256=hashlib.sha256(original).hexdigest()),indent=2)+'\n')
