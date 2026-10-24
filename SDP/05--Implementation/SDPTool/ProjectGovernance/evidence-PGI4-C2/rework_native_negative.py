"""Return the actual failed native review to C1 rework, preserving provenance.

The attached JSON is an attributed artifact. Independent-context observation remains
owned by the dedicated review record, never inferred from this attachment.
"""
import hashlib,json,subprocess,uuid
from pathlib import Path
base=Path(__file__).resolve().parent;root=Path('/home/warloc/git/pg_lab')
summary=json.loads((base/'native-summary.json').read_text());case=next(x for x in summary['nativeReviews'] if x['case']=='negative')
state=Path(case['state']);binary=summary['binary'];wire=[]
def call(args,definition):
 p=subprocess.run([binary,*args,'--project',str(root),'--state',str(state)],capture_output=True,text=True);assert p.returncode==0,p.stderr
 v=json.loads(p.stdout);wire.append({'definition':definition,'value':v});return v
source=root/'format.go';original=source.read_bytes()
try:
 source.write_bytes(original.replace(b'return prefix + ": " + label',b'return prefix + "! " + label'))
 view=call(['review','inspect','--review',case['reviewId']],'ReviewView')
 assert view['applicable'] and view['review']['outcome']['verdict']=='failed'
 candidate=view['review']['spec']['candidate'];run=candidate['runId']
 current=call(['evidence','inspect','--run',run,'--json'],'EvidenceView')
 artifact=root/'evidence/PGI4-C2/negative-review-final.json'
 q={'schema':'sdp-governance-evidence/1','projectRef':view['projectRef'],'checkoutId':view['checkoutId'],'operationId':str(uuid.uuid4()),'runId':run,'expectedRevision':current['revision'],'inputRevision':candidate['inputRevision'],'reportId':view['review']['spec']['reportId'],'candidate':candidate,'kind':'review','verdict':'failed','issuer':'PGI4-C2 native review controller','note':'Recorded native review '+case['reviewId']+' found the deliberate separator defect. This artifact attachment does not itself establish reviewer independence.','artifact':{'path':str(artifact.relative_to(root)),'sha256':hashlib.sha256(artifact.read_bytes()).hexdigest()}}
 wire.append({'definition':'EvidenceRequest','value':q});req=state.parent/'native-finding.json';req.write_text(json.dumps(q));call(['evidence','attach','--input',str(req)],'Result')
 rq={'schema':'sdp-governance-evidence/1','projectRef':view['projectRef'],'checkoutId':view['checkoutId'],'operationId':str(uuid.uuid4()),'runId':run,'expectedRevision':current['revision'],'inputRevision':candidate['inputRevision'],'evidenceId':q['operationId'],'reason':'Return the observed native finding for bounded rework in the same Session; source restored by the test coordinator afterward.'}
 wire.append({'definition':'EvidenceReworkRequest','value':rq});req=state.parent/'native-rework.json';req.write_text(json.dumps(rq));r=call(['evidence','rework','--input',str(req)],'Result');assert r['run']['state']=='ready'
finally:source.write_bytes(original)
final=call(['review','inspect','--review',case['reviewId']],'ReviewView');assert not final['qualified'] and not final['applicable'] and final['review']['outcome']['verdict']=='failed'
(base/'native-rework-wire.json').write_text(json.dumps(wire,indent=2)+'\n')
out={'reviewId':case['reviewId'],'findingArtifact':q['artifact'],'evidenceReceipt':q['operationId'],'reworkReceipt':rq['operationId'],'finalRunState':r['run']['state'],'sameSession':r['run']['selection']['session'],'sourceRestored':True,'limits':'Review provenance remains in the native review record; the C1 artifact is attributed, not an independent-review import.'}
(base/'native-rework.json').write_text(json.dumps(out,indent=2)+'\n');(root/'evidence/PGI4-C2/rework.json').write_text(json.dumps(out,indent=2)+'\n');print(json.dumps(out,indent=2))
