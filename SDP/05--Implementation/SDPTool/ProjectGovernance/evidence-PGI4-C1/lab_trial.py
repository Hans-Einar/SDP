"""Exercise candidate evidence against C8's real completed worker in a state fork.

The fork preserves the original native history. No model is launched. One deliberate
source fault is restored in finally; the original C8 operational store is read-only.
"""
import hashlib,json,os,shutil,subprocess,uuid
from pathlib import Path
base=Path(__file__).resolve().parent
root=Path('/home/warloc/git/pg_lab')
cache=Path('/home/warloc/.cache/sdp-workstation/PGI4-C1')
original=Path('/home/warloc/.cache/sdp-workstation/PGL3-C8/state')
state=cache/'lab-state-final'
binary=cache/'bin/sdp-controller'
run='be00db9f-4e12-4d8e-8728-9fa367450065'
go='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go'
h=lambda b:hashlib.sha256(b).hexdigest()
head=(original/'HEAD').read_bytes()
assert not state.exists(), 'Never reset an existing trial; retain its result.'
shutil.copytree(original,state)
lab_evidence=root/'evidence/PGI4-C1';lab_evidence.mkdir(parents=True,exist_ok=True)
env=os.environ.copy();env.update(GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache')
wire=[]
commands=[]
def call(*args,definition=None):
 cmd=[str(binary),*args,'--project',str(root),'--state',str(state)]
 p=subprocess.run(cmd,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
 commands.append({'args':list(args),'exitCode':p.returncode})
 assert p.returncode==0,p.stderr
 value=json.loads(p.stdout)
 if definition: wire.append({'definition':definition,'value':value})
 return value

def capture():return call('evidence','capture','--run',run,'--dependency','go.mod','--dependency','normalize.go','--dependency','normalize_test.go','--dependency','cmd',definition='Candidate')
def view():return call('evidence','inspect','--run',run,'--json',definition='EvidenceView')
def test(name,expected):
 p=subprocess.run([go,'test','-count=1','./...'],cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
 (lab_evidence/name).write_bytes(p.stdout)
 assert (p.returncode==0)==expected,p.stdout.decode()
 return {'command':[go,'test','-count=1','./...'],'exitCode':p.returncode,'artifact':str((lab_evidence/name).relative_to(root)),'sha256':h(p.stdout)}
def attach(candidate,result,verdict):
 current=view()
 reports=call('inbox','--run',run,'--json')['reports'];assert len(reports)==1
 req={'schema':'sdp-governance-evidence/1','projectRef':current['projectRef'],'checkoutId':current['checkoutId'],'operationId':str(uuid.uuid4()),'runId':run,'expectedRevision':current['revision'],'inputRevision':candidate['inputRevision'],'reportId':reports[0]['id'],'candidate':candidate,'kind':'test','verdict':verdict,'issuer':'PGI4-C1 external Go verifier','note':'External Go test result. Failed case is a deliberate trial fault, not an original C8 defect. No independent review claim.','artifact':{'path':result['artifact'],'sha256':result['sha256']}}
 wire.append({'definition':'EvidenceRequest','value':req})
 path=cache/('request-'+req['operationId']+'.json');path.write_text(json.dumps(req))
 return call('evidence','attach','--input',str(path),definition='Result'),path

files_before={str(p.relative_to(root)):h(p.read_bytes()) for p in root.rglob('*') if p.is_file() and not p.is_symlink() and '.git' not in p.relative_to(root).parts and not p.is_relative_to(lab_evidence)}
source=root/'format.go';source_bytes=source.read_bytes()
unrelated=root/'unrelated.txt';unrelated_bytes=unrelated.read_bytes()
try:
 candidate=capture();passed=test('passed-final.log',True);attached,request_path=attach(candidate,passed,'passed')
 assert view()['items'][0]['applicable']
 unrelated.write_bytes(unrelated_bytes+b'\nPGI4-C1 temporary unrelated edit\n')
 assert view()['items'][0]['applicable']
 unrelated.write_bytes(unrelated_bytes)
 source.write_bytes(source_bytes.replace(b'return prefix + ": " + label',b'return prefix + "! " + label'))
 assert source.read_bytes()!=source_bytes
 stale=view();assert not stale['items'][0]['candidateCurrent'] and not stale['items'][0]['applicable']
 assert call('evidence','attach','--input',str(request_path),definition='Result')==attached
 bad_candidate=capture();failed=test('deliberate-failure-final.log',False);failed_receipt,_=attach(bad_candidate,failed,'failed')
 current=view()
 req={'schema':'sdp-governance-evidence/1','projectRef':current['projectRef'],'checkoutId':current['checkoutId'],'operationId':str(uuid.uuid4()),'runId':run,'expectedRevision':current['revision'],'inputRevision':bad_candidate['inputRevision'],'evidenceId':failed_receipt['operationId'],'reason':'Trial: return the deliberately broken separator to bounded rework; coordinator restores source after observing the transition.'}
 wire.append({'definition':'EvidenceReworkRequest','value':req})
 path=cache/'rework.json';path.write_text(json.dumps(req))
 result=call('evidence','rework','--input',str(path),definition='Result')
 assert result['run']['state']=='ready' and result['run']['revision']==current['revision']+1
 assert call('evidence','rework','--input',str(path),definition='Result')==result
finally:
 source.write_bytes(source_bytes);unrelated.write_bytes(unrelated_bytes)
 assert (original/'HEAD').read_bytes()==head,'original C8 state changed'
restored=test('restored-final.log',True)
final=view();assert len(final['items'])==2 and all(not x['applicable'] and not x['reportCurrent'] for x in final['items'])
assert all(h((root/p).read_bytes())==digest for p,digest in files_before.items()),'baseline file changed'
old=subprocess.run(['/home/warloc/.cache/sdp-workstation/PGL3-C8/bin/sdp-controller','inbox','--project',str(root),'--state',str(state),'--json'],capture_output=True,text=True)
assert old.returncode!=0 and ('incompatible_profile' in old.stderr or 'reconciliation_required' in old.stderr),old.stderr
summary={'labBaseline':subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),'runId':run,'originalStateHead':head.decode(),'originalStatePreserved':True,'stateFork':str(state),'binary':str(binary),'binarySHA256':h(binary.read_bytes()),'candidateBefore':candidate['digest'],'candidateFault':bad_candidate['digest'],'passed':passed,'deliberateFailure':failed,'restored':restored,'preservedBaselineFiles':len(files_before),'finalState':final['state'],'evidenceRecords':len(final['items']),'newModelTurns':0,'oldReaderRefused':{'exitCode':old.returncode,'diagnostic':old.stderr.strip()},'limits':'External tests on C8 native implementation; state fork is a test, not canonical disposition, independent review or publication.'}
(base/'lab-summary.json').write_text(json.dumps(summary,indent=2)+'\n')
(base/'lab-wire.json').write_text(json.dumps(wire,indent=2)+'\n')
(base/'lab-commands.json').write_text(json.dumps(commands,indent=2)+'\n')
(lab_evidence/'trial.json').write_text(json.dumps(summary,indent=2)+'\n')
print(json.dumps(summary,indent=2))
