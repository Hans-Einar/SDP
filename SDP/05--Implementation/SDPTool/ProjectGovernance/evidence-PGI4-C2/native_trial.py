"""Real read-only review trials; no simulated host/model answers.

Use fresh state forks; never reset an uncertain review. The deliberate negative
source change is restored even if the test fails. No original C8 state is mutated.
"""
import hashlib,json,os,shutil,subprocess,uuid
from pathlib import Path
base=Path(__file__).resolve().parent
root=Path('/home/warloc/git/pg_lab');cache=Path('/home/warloc/.cache/sdp-workstation/PGI4-C2')
original=Path('/home/warloc/.cache/sdp-workstation/PGL3-C8/state');binary=cache/'bin/sdp-controller';mcp=cache/'bin/sdp-mcp'
proxy='/home/warloc/.cache/sdp-workstation/PGL3-C8/codex-proxy'
go='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go';run='be00db9f-4e12-4d8e-8728-9fa367450065'
h=lambda b:hashlib.sha256(b).hexdigest();wire=[];results=[]
source=root/'format.go';original_source=source.read_bytes();original_head=(original/'HEAD').read_bytes()
lab=root/'evidence/PGI4-C2';lab.mkdir(parents=True,exist_ok=True)
baseline={str(p.relative_to(root)):h(p.read_bytes()) for p in root.rglob('*') if p.is_file() and not p.is_symlink() and '.git' not in p.relative_to(root).parts and not p.is_relative_to(lab)}
baseenv=os.environ.copy();baseenv.update(GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache')
try:
 for name,fault,expected in [('positive','turn-start','passed'),('negative','terminal-read','failed')]:
  state=cache/('state-final-'+name);assert not state.exists(),'Preserve existing trial state; do not reset.';shutil.copytree(original,state)
  source.write_bytes(original_source if name=='positive' else original_source.replace(b'return prefix + ": " + label',b'return prefix + "! " + label'))
  env=baseenv.copy();env.update(PG_LAB_PROTOCOL_LOG=str(base/(name+'-prepare-protocol.ndjson')))
  def invoke(args,definition=None,allow_failure=False,label=None):
   cmd=[str(binary),*args,'--project',str(root),'--state',str(state)]
   p=subprocess.run(cmd,env=env,capture_output=True,text=True,timeout=600)
   if label:
    (base/(name+'-'+label+'.stdout')).write_text(p.stdout);(base/(name+'-'+label+'.stderr')).write_text(p.stderr)
   if not allow_failure:assert p.returncode==0,p.stderr
   values=[json.loads(line) for line in p.stdout.splitlines() if line.strip()]
   if definition:wire.extend({'definition':definition,'value':v} for v in values)
   return p,values
  _,values=invoke(['evidence','capture','--run',run,'--dependency','go.mod','--dependency','normalize.go','--dependency','normalize_test.go','--dependency','cmd'],'Candidate');candidate=values[-1]
  tests=subprocess.run([go,'test','-count=1','./...'],cwd=root,env=env,capture_output=True)
  log=tests.stdout+tests.stderr;(lab/(name+'-tests-final.log')).write_bytes(log)
  assert (tests.returncode==0)==(name=='positive'),log.decode()
  _,reports=invoke(['inbox','--run',run,'--json']);report=reports[0]['reports'][0]
  request={'schema':'sdp-governance-evidence/1','projectRef':candidate['projectRef'],'checkoutId':candidate['checkoutId'],'operationId':str(uuid.uuid4()),'runId':run,'expectedRevision':report['currentRevision'],'inputRevision':candidate['inputRevision'],'reportId':report['id'],'candidate':candidate,'kind':'test','verdict':expected,'issuer':'PGI4-C2 external Go verifier','note':'Actual Go test result; negative case is a deliberate separator defect. Native reviewer must inspect independently.','artifact':{'path':str((lab/(name+'-tests-final.log')).relative_to(root)),'sha256':h(log)}}
  req=cache/(name+'-evidence.json');req.write_text(json.dumps(request));invoke(['evidence','attach','--input',str(req)],'Result')
  spec={'candidate':candidate,'reportId':report['id'],'evidenceIds':[request['operationId']],'focus':'Review FormatLabel against the selected authority and LeaderTrial plan. Normalize input with NormalizeLabel; empty normalized input returns empty. Otherwise return the prefix unchanged, colon-space, and normalized input. Check exact source, existing normalization dependency and tests, including prefix whitespace/case and Unicode trimming. Inspect the supplied test artifact; do not execute tests or mutate files. Report concrete defects or missing evidence honestly.'}
  specpath=cache/(name+'-spec.json');specpath.write_text(json.dumps(spec));wire.append({'definition':'ReviewInput','value':spec})
  review=str(uuid.uuid4());args=['review','prepare','--operation',review,'--input',str(specpath),'--mcp-binary',str(mcp),'--codex',proxy,'--skill','Skills/sdp-reviewer/SKILL.md']
  invoke(args,'Result',label='prepare')
  env.update(PG_LAB_PROTOCOL_LOG=str(base/(name+'-execution-protocol.ndjson')),PG_LAB_FAULT_MODE=fault,PG_LAB_FAULT_MARKER=str(cache/(name+'-fault-final.json')))
  args=['review','run','--review',review,'--mcp-binary',str(mcp),'--codex',proxy,'--timeout','8m','--interval','1s']
  print(name+': native reviewer launched; deliberate '+fault+' response loss selected',flush=True)
  failed,values=invoke(args,'ReviewView',allow_failure=True,label='interrupted')
  assert failed.returncode!=0,'Expected injected transport failure'
  marker=cache/(name+'-fault-final.json');assert marker.exists(),failed.stderr
  (base/(name+'-fault.json')).write_bytes(marker.read_bytes())
  env['PG_LAB_PROTOCOL_LOG']=str(base/(name+'-recovery-protocol.ndjson'))
  _,values=invoke(args,'ReviewView',label='recovered');v=values[-1]
  assert v['review']['state']=='completed',v
  assert v['review'].get('outcome',{}).get('verdict')==expected,v
  assert v['applicable'] and v['qualified']==(expected=='passed'),v
  methods=[json.loads(line)['method'] for line in Path(env['PG_LAB_PROTOCOL_LOG']).read_text().splitlines() if json.loads(line).get('direction')=='request']
  assert methods==['initialize','thread/read'],methods
  original_snapshot=json.loads((original/'generations'/(original_head.decode()+'.json')).read_text())
  implementation_threads={a['threadId'] for r in original_snapshot['runs'].values() for a in r.get('attempts',[])}
  assert v['review']['threadId'] not in implementation_threads, 'Review reused implementation context'
  (lab/(name+'-review-final.json')).write_text(json.dumps(v,indent=2)+'\n')
  results.append({'case':name,'reviewId':review,'threadId':v['review']['threadId'],'turnId':v['review']['turnId'],'verdict':expected,'qualified':v['qualified'],'candidate':candidate['digest'],'fault':fault,'recoveryMethods':methods,'state':str(state),'testExitCode':tests.returncode})
  print(name+': completed '+expected+'; exact recovery without replacement',flush=True)
finally:
 source.write_bytes(original_source)
 assert (original/'HEAD').read_bytes()==original_head
 assert all(h((root/p).read_bytes())==digest for p,digest in baseline.items()),'Baseline file changed'
restored=subprocess.run([go,'test','-count=1','./...'],cwd=root,env=baseenv,capture_output=True);assert restored.returncode==0
(lab/'restored-tests-final.log').write_bytes(restored.stdout+restored.stderr)
summary={'productBaseline':subprocess.check_output(['git','rev-parse','HEAD'],text=True).strip(),'labBaseline':subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),'binary':str(binary),'binarySHA256':h(binary.read_bytes()),'mcpSHA256':h(mcp.read_bytes()),'originalC8Head':original_head.decode(),'originalStatePreserved':True,'baselineFilesPreserved':len(baseline),'nativeReviews':results,'limits':'Native product-test reviews of pg_lab only; not independent review of this PG implementation, native child hierarchy, owner acceptance or publication.'}
(base/'native-summary.json').write_text(json.dumps(summary,indent=2)+'\n');(base/'native-wire.json').write_text(json.dumps(wire,indent=2)+'\n');(lab/'trial.json').write_text(json.dumps(summary,indent=2)+'\n')
print(json.dumps(summary,indent=2),flush=True)
