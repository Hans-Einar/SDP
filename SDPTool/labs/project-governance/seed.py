import json,subprocess,uuid,hashlib
from pathlib import Path
import sys
base=Path(sys.argv[1]).resolve();root=base/'project';root.mkdir(parents=True,exist_ok=False)
files={
'AGENTS.md':'# Pilot assignment\nRead SDP/authority.md and selected work. No delegation. Only normalize.go and normalize_test.go may change. Never edit SDP records, go.mod or unrelated.txt. This is an owner development trial of PLAN-SDP-0018 PGI3-U1. Use the available Go toolchain.\n',
'SDP/SDP-project.manifest.yaml':"schemaVersion: '1.0'\nproject:\n  name: pg-owner-lab\n",
'SDP/KanBan/active/card.md':'| Field | Value |\n| --- | --- |\n| id | KB-PGLAB-001 |\n| PlanId | MAINT-PGLAB-0001 |\n| CardState | in-progress |\n',
'SDP/Maintenance/Plan.md':'# Bounded label-normalization fix\n\n| Field | Value |\n| --- | --- |\n| id | MAINT-PGLAB-0001 |\n| PlanType | MaintenancePlan |\n| state | active |\n\nM1: NormalizeLabel must trim surrounding Unicode whitespace while preserving interior whitespace and case. Fix only normalize.go and normalize_test.go. Verify go test ./...; retain the unrelated file. No delegation, publication or records editing.\n',
'SDP/Sessions/session.md':'# Live pilot goal\n\n| State | Step | Work |\n| --- | --- | --- |\n| active | S1 | Fix normalization and verify preservation |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-PGLAB-0001 |\n| Status | active |\n',
'SDP/authority.md':'Owner request on 2026-10-09 selects a local PG development trial. The owner starts this bounded assignment explicitly through the lab run command. This lab uses minimal PG coverage records, not a complete released SDP installation. Fix Unicode boundary whitespace normalization; do not expand scope or invoke child agents. Canonical SDP documents must remain byte-identical.\n',
'SDP/ProjectManagement/Ledger.ndjson':json.dumps({'schemaVersion':'1.0','eventId':'EVT-KB-PGLAB-000001','eventType':'x-kanban:created','subjectId':'KB-PGLAB-001','payload':{'schemaVersion':'0.2','previousEventId':None,'to':'active','toPath':'active/card.md'}})+'\n'+json.dumps({'schemaVersion':'1.0','eventId':'EVT-PM-PGLAB-000001','eventType':'x-management:started','subjectId':'MAINT-PGLAB-0001','payload':{'schemaVersion':'0.2','previousEventId':None,'to':'active','toPath':'Maintenance/Plan.md'}})+'\n',
'cmd/label/main.go':'package main\nimport (\"fmt\";\"os\"; labels \"pgi3pilot\")\nfunc main(){if len(os.Args)!=2{fmt.Fprintln(os.Stderr,\"usage: label TEXT\");os.Exit(2)};fmt.Printf(\"%q\\n\",labels.NormalizeLabel(os.Args[1]))}\n',
'go.mod':'module pgi3pilot\n\ngo 1.26.0\n',
'normalize.go':'package pgi3pilot\n\nimport "strings"\n\nfunc NormalizeLabel(s string) string { return strings.Trim(s, " ") }\n',
'normalize_test.go':'package pgi3pilot\n\nimport "testing"\n\nfunc TestNormalizeLabel(t *testing.T) { if got := NormalizeLabel("  Keep  Case  "); got != "Keep  Case" { t.Fatal(got) } }\n',
'unrelated.txt':'This file is outside the assignment; preserve these exact bytes.\n'
}
for name,content in files.items():
 p=root/name;p.parent.mkdir(parents=True,exist_ok=True);p.write_text(content)
# Install the actual maintained skill sources into this disposable fixture only.
import shutil
source=base/'candidate/source/Skills'
shutil.copytree(source,root/'Skills')
(root/'.agents/skills').mkdir(parents=True)
for entry in source.iterdir():
 if entry.is_dir() and (entry/'SKILL.md').exists():
  (root/'.agents/skills'/entry.name).symlink_to('../../Skills/'+entry.name)
for entry in (root/'Skills').rglob('*'):
 if entry.is_file():files[str(entry.relative_to(root))]=entry.read_text()
sel={'card':'SDP/KanBan/active/card.md','plan':'SDP/Maintenance/Plan.md','session':'SDP/Sessions/session.md','authority':'SDP/authority.md','step':'S1'}
state=base/'runtime/state';binary=str(base/'candidate/bin/sdptool')
def call(op,q=None):
 args=[binary,str(root),'governance',op,'--state',str(state),'--json']
 if q is not None:
  f=base/(op+'.json');f.write_text(json.dumps(q));args+=['--input',str(f)]
 return json.loads(subprocess.check_output(args))
c=call('context',sel);assert c['conditions']==[],c
r=call('route',{'schema':'sdp-project-governance/0.1','categories':['maintenance'],'selection':sel,'routineId':'maintenance','routineVersion':'0.1'});assert r['disposition']=='matched',r
b=call('bind')
run=call('start',{'schema':'sdp-project-governance/0.1','projectRef':b['projectRef'],'checkoutId':b['checkoutId'],'operationId':str(uuid.uuid4()),'requestId':str(uuid.uuid4()),'inputRevision':c['revision'],'routineDigest':r['routineDigest'],'selection':sel,'assignment':{'outcome':'Fix NormalizeLabel Unicode boundary whitespace while preserving interior content; verify with Go tests','writablePaths':['normalize.go','normalize_test.go'],'excludedPaths':['SDP','Skills','AGENTS.md','go.mod','unrelated.txt','cmd']}})
(base/'run.json').write_text(json.dumps(run,indent=2)+'\n')
(base/'discovery-links.json').write_text(json.dumps({str(p.relative_to(root)):str(p.readlink()) for p in (root/'.agents/skills').iterdir()},indent=2)+'\n')
(base/'preimages.json').write_text(json.dumps({name:hashlib.sha256((root/name).read_bytes()).hexdigest() for name in files},indent=2)+'\n')
(base/'owner.txt').write_text('Fix NormalizeLabel so the label command removes surrounding Unicode whitespace while preserving interior whitespace and letter case. Follow the selected skills and work records. Work only on normalize.go and normalize_test.go. Run go test ./... and summarize the observable change, tests and limitations. Do not delegate, commit, publish records or claim independent review. After testing, call run_read and then one run_transition action report using its exact current identities and expected revision with a fresh UUID operationId. Wait for the owner to answer the resulting per-request prompt. A report means ready for checking, not completed work or owner acceptance. Do not bypass a rejected report.\n')
subprocess.run(['git','init','-b','trial/initial'],cwd=root,check=True,capture_output=True)
subprocess.run(['git','add','.'],cwd=root,check=True)
subprocess.run(['git','-c','user.name=PG Lab','-c','user.email=pg-lab@localhost','commit','-m','PGI3-U1: initial owner trial project'],cwd=root,check=True,capture_output=True)
(base/'FEEDBACK.md').write_text('# Owner trial feedback\n\nCandidate and run identity are recorded in candidate.json and run.json.\n\n| Observation | Expected workflow | Actual experience | Classification | Disposition |\n| --- | --- | --- | --- | --- |\n\nRecord confusing decisions, missing context, difficult controls and desired changes.\nClassify as implementation defect, workflow mismatch, new requirement or understood limitation.\nNo owner trial or acceptance is recorded yet.\n')
print('Prepared owner trial; run',run['run']['id'])
