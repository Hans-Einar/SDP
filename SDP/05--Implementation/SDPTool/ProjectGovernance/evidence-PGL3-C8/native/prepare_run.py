from pathlib import Path
import json,uuid,subprocess,hashlib
base=Path(__file__).parent;root=Path('/home/warloc/git/pg_lab');state=base/'state';tool=base/'bin/sdptool';controller=base/'bin/sdp-controller';cases=[]
def core(action,q=None):
 cmd=[str(tool),str(root),'governance',action,'--state',str(state),'--json']
 if q is not None:
  p=base/f'request-{len(cases)}.json';p.write_text(json.dumps(q)+'\n');cmd+=['--input',str(p)]
 r=subprocess.run(cmd,text=True,capture_output=True,timeout=30);assert r.returncode==0,r.stderr
 v=json.loads(r.stdout);cases.append({'operation':action,'request':q,'response':v});return v
sel={'card':'SDP/KanBan/active/leader-trial.md','plan':'SDP/Maintenance/LeaderTrial/Plan.md','session':'SDP/Sessions/leader-trial.md','authority':'SDP/leader-trial-authority.md','step':'S1'}
schema='sdp-project-governance/0.1';context=core('context',sel);assert not context['conditions'],context
route=core('route',{'schema':schema,'categories':['maintenance'],'selection':sel,'routineId':'maintenance','routineVersion':'0.1'})
b=core('bind')
r=core('start',{'schema':schema,'projectRef':b['projectRef'],'checkoutId':b['checkoutId'],'operationId':str(uuid.uuid4()),'requestId':str(uuid.uuid4()),'inputRevision':context['revision'],'routineDigest':route['routineDigest'],'selection':sel,'assignment':{'outcome':'Implement and test FormatLabel in the native PG lab, then report for Project Leader processing','writablePaths':['format.go','format_test.go'],'excludedPaths':['SDP','Skills','AGENTS.md','go.mod','normalize.go','normalize_test.go','unrelated.txt']}})
(base/'preparation.json').write_text(json.dumps(cases,indent=2)+'\n');(base/'run-id').write_text(r['run']['id'])
message=base/'leader-message-1.txt';message.write_text('Synthetic test-owner message L1, under actual owner-authorized pg_lab testing. Read the supplied project-leader skill and coordination_read. Identify the ready FormatLabel assignment and its allowed files. Return a brief status containing L1 and the run ID. Do not launch work, edit files, approve anything or acknowledge a report yet. This is one short coordination turn.\n')
q=[str(controller),'leader','message','--project',str(root),'--state',str(state),'--owner','test-owner','--operation',str(uuid.uuid4()),'--input',str(message)]
v=subprocess.run(q,capture_output=True,text=True,check=True);(base/'leader-enqueue-1.json').write_text(v.stdout)
paths=subprocess.check_output(['git','ls-files','-z'],cwd=root).decode().split('\0')
(base/'protected-before.json').write_text(json.dumps({p:hashlib.sha256((root/p).read_bytes()).hexdigest() for p in paths if p and not (root/p).is_symlink()},indent=2)+'\n')
print('Ready worker run:',r['run']['id'])
