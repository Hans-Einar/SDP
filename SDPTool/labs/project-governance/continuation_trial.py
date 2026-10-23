#!/usr/bin/env python3
"""Native two-attempt attention trial in a new disposable project; no release needed."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import signal
import subprocess
import time
import uuid


def main():
    p = argparse.ArgumentParser(description=__doc__)
    for name in ['destination', 'controller', 'sdptool', 'mcp', 'host', 'old-controller']:
        p.add_argument('--'+name, type=Path, required=True)
    a = p.parse_args()
    base = a.destination.resolve()
    base.mkdir(parents=True, exist_ok=False)
    root, state = base/'project', base/'runtime/state'
    controller, tool, mcp = base/'sdp-controller', base/'sdptool', base/'sdp-mcp'
    for source,dest in [(a.controller,controller),(a.sdptool,tool),(a.mcp,mcp)]: shutil.copy2(source,dest)
    files = {
        'AGENTS.md': '# Bounded native attention trial\nRead SDP/authority.md and selected work. Only label.go and label_test.go may change. Preserve all other files and do not create additional project files. Put caches in /tmp. No delegation, network, installation, commits, publication or host configuration changes. Ask the required question before editing.\n',
        'SDP/SDP-project.manifest.yaml': "schemaVersion: '1.0'\nproject:\n  name: pg-continuation-lab\n",
        'SDP/KanBan/active/card.md': '| Field | Value |\n| --- | --- |\n| id | KB-PGCONT-001 |\n| PlanId | MAINT-PGCONT-0001 |\n| CardState | in-progress |\n',
        'SDP/Maintenance/Plan.md': '# Label prefix trial\n\n| Field | Value |\n| --- | --- |\n| id | MAINT-PGCONT-0001 |\n| PlanType | MaintenancePlan |\n| state | active |\n\nM1: Ask which prefix to use, then implement the recorded choice in Label while retaining Unicode edge trimming and interior content. Change only label.go and label_test.go; verify Go tests.\n',
        'SDP/Sessions/session.md': '| State | Step |\n| --- | --- |\n| active | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-PGCONT-0001 |\n| Status | active |\n',
        'SDP/authority.md': 'Owner authorizes this isolated development trial of PG attention continuation. The scripted fixture-owner reply is test input, not actual production design approval. Ask before modifying code; continue after the recorded choice within the same assignment. No delegation or canonical record writes.\n',
        'go.mod': 'module pgattentionfixture\n\ngo 1.26.0\n',
        'label.go': 'package pgattentionfixture\n\nimport "strings"\n\nfunc Label(s string) string { return strings.TrimSpace(s) }\n',
        'label_test.go': 'package pgattentionfixture\n\nimport "testing"\n\nfunc TestInitialLabel(t *testing.T) { if got := Label(" Keep  Case "); got != "Keep  Case" { t.Fatal(got) } }\n',
        'unrelated.txt': 'Preserve these exact bytes.\n'
    }
    ledger=[]
    for subject,kind,path in [('KB-PGCONT-001','x-kanban:created','active/card.md'),('MAINT-PGCONT-0001','x-management:started','Maintenance/Plan.md')]:
        ledger.append({'schemaVersion':'1.0','eventId':f'FIXTURE-{len(ledger)+1}','eventType':kind,'subjectId':subject,'payload':{'schemaVersion':'0.2','previousEventId':None,'to':'active','toPath':path}})
    files['SDP/ProjectManagement/Ledger.ndjson']=''.join(json.dumps(x)+'\n' for x in ledger)
    for name,value in files.items():
        path=root/name;path.parent.mkdir(parents=True,exist_ok=True);path.write_text(value)
    sha=lambda path:hashlib.sha256(path.read_bytes()).hexdigest()
    before={name:sha(root/name) for name in files}
    cases=[]
    def core(action,request=None):
        cmd=[str(tool),str(root),'governance',action,'--state',str(state),'--json']
        if request is not None:
            path=base/f'request-{len(cases)}.json';path.write_text(json.dumps(request)+'\n');cmd+=['--input',str(path)]
        r=subprocess.run(cmd,capture_output=True,text=True,check=True,timeout=30)
        value=json.loads(r.stdout);cases.append({'operation':action,'request':request,'response':value});return value
    sel={'card':'SDP/KanBan/active/card.md','plan':'SDP/Maintenance/Plan.md','session':'SDP/Sessions/session.md','authority':'SDP/authority.md','step':'S1'}
    schema='sdp-project-governance/0.1'
    context=core('context',sel);assert not context['conditions'],context
    route=core('route',{'schema':schema,'categories':['maintenance'],'selection':sel,'routineId':'maintenance','routineVersion':'0.1'})
    binding=core('bind')
    run=core('start',{'schema':schema,'projectRef':binding['projectRef'],'checkoutId':binding['checkoutId'],'operationId':str(uuid.uuid4()),'requestId':str(uuid.uuid4()),'inputRevision':context['revision'],'routineDigest':route['routineDigest'],'selection':sel,'assignment':{'outcome':'Ask the label-prefix choice, implement it, verify and report','writablePaths':['label.go','label_test.go'],'excludedPaths':['SDP','AGENTS.md','go.mod','unrelated.txt']}})
    run_id=run['run']['id']
    common=['--project',str(root),'--state',str(state)]
    def inspect():
        for _ in range(20):
            result=subprocess.run([str(controller),'attention','list',*common,'--json'],capture_output=True,text=True,timeout=30)
            if result.returncode==0:return json.loads(result.stdout)
            if 'stale_revision: attention changed during inspection' not in result.stderr:raise AssertionError(result.stderr)
            time.sleep(0.1)
        raise AssertionError('Attention inspection remained unstable')
    def current():
        head=(state/'HEAD').read_text();return json.loads((state/'generations'/f'{head}.json').read_text())['runs'][run_id]
    golden=base/'verify';golden.mkdir()
    (golden/'go.mod').write_text('module pgverify\n\ngo 1.26.0\nrequire pgattentionfixture v0.0.0\nreplace pgattentionfixture => '+json.dumps(str(root))+'\n')
    (golden/'main.go').write_text('package main\nimport ("fmt";"os"; f "pgattentionfixture")\nfunc main(){for _,c:=range [][2]string{{"\\t Keep  Case\\u00a0","OK: Keep  Case"},{"A\\tB","OK: A\\tB"},{"\\t\\n","OK: "},{"Mixed CASE","OK: Mixed CASE"}}{if got:=f.Label(c[0]);got!=c[1]{fmt.Printf("FAIL: %q -> %q; want %q\\n",c[0],got,c[1]);os.Exit(1)}};fmt.Println("PASS: selected prefix, Unicode trimming and interior preservation")}\n')
    initial=subprocess.run(['go','run','.'],cwd=golden,capture_output=True,text=True,timeout=60);assert initial.returncode!=0 and 'FAIL:' in initial.stdout
    message=base/'message.txt'
    message.write_text('This is an owner-authorized isolated PG attention/continuation test. First inspect context and run. If the launch payload has no attentionReferences, ask exactly one blocking attention_request question with options short (prefix "OK: ") and detailed (prefix "Result: "). Do not edit or report before the answer. Return a waiting summary immediately after the question is stored. On a fresh continuation with attentionReferences, use the selected recorded reply from attentionContext; the earlier question basis becoming historical after a new attempt is expected. Do not ask it again. Implement the selected prefix in Label in label.go, retain Unicode edge trimming, interior whitespace and case; update label_test.go accordingly. Preserve all other files. Run go test ./..., then use run_read and run_transition report with a fresh UUID and current identities. No delegation, network, installations, host config changes, commits or record publication. A report is not owner acceptance.\n')
    cmd=[str(controller),'run',*common,'--run',run_id,'--codex',str(a.host.resolve()),'--mcp-binary',str(mcp),'--message',str(message),'--sandbox','workspace-write','--approvals','never','--reporting','automatic','--attention-continuations','1','--timeout','6m']
    evidence={'status':'running','hostKind':'native','modelTaskSubmitted':False,'actualOwnerApproval':False,'runId':run_id,'cases':cases,'goldenBefore':initial.stdout,'controllerSHA256':sha(controller),'sdptoolSHA256':sha(tool),'mcpSHA256':sha(mcp),'hostBinarySHA256':sha(a.host.resolve())}
    record=lambda:(base/'evidence.json').write_text(json.dumps(evidence,indent=2)+'\n')
    record();print('Starting native attention/continuation trial: '+str(base),flush=True)
    process=None
    try:
        with (base/'controller.ndjson').open('w') as out,(base/'controller.err').open('w') as err:
            process=subprocess.Popen(cmd,stdin=subprocess.DEVNULL,stdout=out,stderr=err)
            deadline=time.monotonic()+300
            while time.monotonic()<deadline:
                if process.poll() is not None:raise AssertionError('Controller ended before waiting: '+(base/'controller.err').read_text())
                r=current()
                if r.get('attempts') and r['attempts'][-1]['status']=='completed':
                    items=inspect()['items']
                    if items:break
                time.sleep(1)
            else:raise AssertionError('No completed waiting question within deadline')
            evidence['modelTaskSubmitted']=True
            assert len(items)==1 and items[0]['status']=='open' and items[0]['applicability']=='current',items
            assert before=={name:sha(root/name) for name in files},'Code changed before answer'
            evidence['waitingBeforeReply']={'attempts':len(r['attempts']),'question':items[0],'allProjectBytesUnchanged':True}
            # Delayed input is a separate CLI process. The waiting loop must not write receipts.
            head=(state/'HEAD').read_text();time.sleep(2);assert head==(state/'HEAD').read_text(),'Waiting appended state'
            q=items[0]['question'];reply={'schema':'sdp-governance-attention/1','projectRef':run['projectRef'],'checkoutId':run['checkoutId'],'operationId':str(uuid.uuid4()),'runId':run_id,'questionId':q['id'],'basis':q['basis']}
            path=base/'reply.json';path.write_text(json.dumps(reply)+'\n')
            reply_result=subprocess.run([str(controller),'attention','reply',*common,'--input',str(path),'--owner','fixture-owner'],input='short\nUse the prefix OK: followed by the trimmed input. Synthetic trial input only.\n',capture_output=True,text=True,timeout=30)
            assert reply_result.returncode==0,reply_result.stderr
            evidence['replyReceipt']=json.loads(reply_result.stdout);record()
            print('Recorded delayed fixture reply; waiting for fresh attempt and report.',flush=True)
            process.wait(timeout=240)
            assert process.returncode==0,(base/'controller.err').read_text()
        r=current();view=inspect()
        evidence.update(finalRun=r,attention=view,controllerExitCode=process.returncode)
        assert len(r['attempts'])==2 and all(v['status']=='completed' for v in r['attempts']) and r['state']=='verifying',r
        assert len(view['items'])==1 and view['items'][0]['consumption']=='dispatch-acknowledged',view
        assert len(r['inputs'])==2 and r['inputs'][1]['attentionReplies'][0]['replyOperationId']==reply['operationId']
        after={str(x.relative_to(root)):sha(x) for x in root.rglob('*') if x.is_file()}
        assert set(after)==set(before),'Unexpected project file additions/removals'
        changed=[n for n in before if before[n]!=after[n]];assert set(changed)=={'label.go','label_test.go'},changed
        golden_after=subprocess.run(['go','run','.'],cwd=golden,capture_output=True,text=True,timeout=60);assert golden_after.returncode==0,golden_after.stdout+golden_after.stderr
        tests=subprocess.run(['go','test','./...'],cwd=root,capture_output=True,text=True,timeout=60);assert tests.returncode==0,tests.stdout+tests.stderr
        head=(state/'HEAD').read_text()
        old=subprocess.run([str(a.old_controller.resolve()),'attention','list',*common,'--json'],capture_output=True,text=True,timeout=30);assert old.returncode!=0 and head==(state/'HEAD').read_text()
        evidence.update(status='passed',changedFiles=changed,protectedFiles=len(before)-len(changed),goldenAfter=golden_after.stdout,projectTests=tests.stdout,oldReader={'exitCode':old.returncode,'stderr':old.stderr},controllerTranscriptSHA256=sha(base/'controller.ndjson'))
        record();(base/'README.md').write_text('# Native attention continuation trial\n\nThe first agent asked and waited without editing. A separate CLI captured a synthetic reply. The watcher started one fresh bounded attempt in the same run; it implemented the selected prefix and reported. External behavior tests and project tests passed. The run is verifying, not accepted. evidence.json retains the actual candidate and results.\n')
        print('PASS: native ask, delayed reply, one fresh attempt, selected behavior, automatic report and protected files.',flush=True)
    except BaseException as error:
        evidence.update(status='failed',failure=str(error),finalRun=current());evidence['modelTaskSubmitted']=any(v.get('turnId') for v in evidence['finalRun'].get('attempts',[]));evidence['dispatchMayBeUncertain']=any(v.get('state') in ['sending','uncertain'] for v in evidence['finalRun'].get('inputs',[]));record();raise
    finally:
        if process and process.poll() is None:
            process.send_signal(signal.SIGINT)
            try:process.wait(timeout=10)
            except subprocess.TimeoutExpired:process.kill();process.wait()

if __name__=='__main__':main()
