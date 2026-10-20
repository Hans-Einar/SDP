#!/usr/bin/env python3
"""Exercise durable attention with a separate synthetic project; starts no model."""
import argparse
import hashlib
import json
from pathlib import Path
import shutil
import subprocess
import uuid


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--destination', type=Path, required=True)
    p.add_argument('--controller', type=Path, required=True)
    p.add_argument('--sdptool', type=Path, required=True)
    p.add_argument('--old-controller', type=Path, required=True)
    a = p.parse_args()
    base = a.destination.resolve()
    base.mkdir(parents=True, exist_ok=False)
    root, state = base / 'project', base / 'runtime/state'
    controller, tool = base / 'sdp-controller', base / 'sdptool'
    shutil.copy2(a.controller, controller)
    shutil.copy2(a.sdptool, tool)
    files = {
        'AGENTS.md': '# Attention fixture\nNo model execution or delegation is authorized.\n',
        'SDP/SDP-project.manifest.yaml': "schemaVersion: '1.0'\nproject:\n  name: pg-attention-lab\n",
        'SDP/KanBan/active/card.md': '| Field | Value |\n| --- | --- |\n| id | KB-PGATT-001 |\n| PlanId | MAINT-PGATT-0001 |\n| CardState | in-progress |\n',
        'SDP/Maintenance/Plan.md': '| Field | Value |\n| --- | --- |\n| id | MAINT-PGATT-0001 |\n| PlanType | MaintenancePlan |\n| state | active |\n',
        'SDP/Sessions/session.md': '| State | Step |\n| --- | --- |\n| active | S1 |\n\n| Field | Value |\n| --- | --- |\n| Session reference | SESSION-PGATT-0001 |\n| Status | active |\n',
        'SDP/authority.md': 'Synthetic controller-only attention trial. No real design approval or model dispatch.\n',
        'target.txt': 'Example output label\n', 'unrelated.txt': 'Preserve these bytes.\n'
    }
    ledger = []
    for subject, kind, path in [('KB-PGATT-001', 'x-kanban:created', 'active/card.md'), ('MAINT-PGATT-0001', 'x-management:started', 'Maintenance/Plan.md')]:
        ledger.append({'schemaVersion':'1.0', 'eventId':f'FIXTURE-{len(ledger)+1}', 'eventType':kind, 'subjectId':subject, 'payload':{'schemaVersion':'0.2', 'previousEventId':None, 'to':'active', 'toPath':path}})
    files['SDP/ProjectManagement/Ledger.ndjson'] = ''.join(json.dumps(v)+'\n' for v in ledger)
    for name, value in files.items():
        path = root / name
        path.parent.mkdir(parents=True, exist_ok=True)
        path.write_text(value)
    original = {name: hashlib.sha256((root/name).read_bytes()).hexdigest() for name in files}
    cases = []
    schema = 'sdp-project-governance/0.1'
    att = 'sdp-governance-attention/1'
    uid = lambda: str(uuid.uuid4())

    def call(group, action, request=None, stdin=None, extra=(), failure=None):
        cmd = ([str(tool), str(root), 'governance', action, '--state', str(state), '--json'] if group == 'core'
               else [str(controller), 'attention', action, '--project', str(root), '--state', str(state)])
        if request is not None:
            path = base / f'request-{len(cases)}.json'
            path.write_text(json.dumps(request, indent=2)+'\n')
            cmd += ['--input', str(path)]
        cmd += list(extra)
        result = subprocess.run(cmd, input=stdin, capture_output=True, text=True)
        case = {'operation':f'{group}.{action}', 'request':request, 'exitCode':result.returncode, 'stderr':result.stderr}
        cases.append(case)
        if failure:
            assert result.returncode != 0 and failure in result.stderr, case
            case['expectedFailure'] = failure
            return
        assert result.returncode == 0, case
        case['response'] = json.loads(result.stdout)
        return case['response']

    sel = {'card':'SDP/KanBan/active/card.md', 'plan':'SDP/Maintenance/Plan.md', 'session':'SDP/Sessions/session.md', 'authority':'SDP/authority.md', 'step':'S1'}
    context = call('core','context',sel)
    assert not context['conditions'], context
    route = call('core','route',{'schema':schema, 'categories':['maintenance'], 'selection':sel, 'routineId':'maintenance', 'routineVersion':'0.1'})
    binding = call('core','bind')

    def start():
        q = {'schema':schema, 'projectRef':binding['projectRef'], 'checkoutId':binding['checkoutId'], 'operationId':uid(), 'requestId':uid(), 'inputRevision':context['revision'], 'routineDigest':route['routineDigest'], 'selection':sel, 'assignment':{'outcome':'Clarify the fixture output label', 'writablePaths':['target.txt'], 'excludedPaths':['SDP','AGENTS.md','unrelated.txt']}}
        return transition(call('core','start',q), 'start')

    def transition(r, action, failure=None):
        q = {'schema':schema, 'projectRef':r['projectRef'], 'checkoutId':r['checkoutId'], 'operationId':uid(), 'runId':r['run']['id'], 'expectedRevision':r['run']['revision'], 'inputRevision':r['run']['inputRevision'], 'action':action, 'reason':'Synthetic attention trial'}
        return call('core','transition',q,failure=failure)

    def ask(r, prompt):
        q = {'schema':att, 'projectRef':r['projectRef'], 'checkoutId':r['checkoutId'], 'operationId':uid(), 'runId':r['run']['id'], 'expectedRevision':r['run']['revision'], 'inputRevision':r['run']['inputRevision'], 'prompt':prompt, 'options':[{'id':'short','label':'Short label'}, {'id':'detailed','label':'Detailed label'}], 'blocking':True}
        return call('attention','request',q)

    def disposition(r, index=0):
        q = r['run']['attention'][index]
        return {'schema':att, 'projectRef':r['projectRef'], 'checkoutId':r['checkoutId'], 'operationId':uid(), 'runId':r['run']['id'], 'questionId':q['id'], 'basis':q['basis']}

    r = ask(start(), 'Which output label should this fixture use?')
    transition(r, 'report', 'attention_required')
    second = start()
    assert transition(second, 'report')['run']['state'] == 'verifying'
    reply = disposition(r)
    head = (state/'HEAD').read_bytes()
    call('attention','reply',reply,stdin='',extra=['--owner','fixture-owner'],failure='no reply recorded')
    call('attention','reply',reply,stdin='missing\nSynthetic response.\n',extra=['--owner','fixture-owner'],failure='invalid_request')
    assert (state/'HEAD').read_bytes() == head
    answered = call('attention','reply',reply,stdin='short\nSynthetic fixture response: use a short label.\n',extra=['--owner','fixture-owner'])
    assert call('attention','reply',reply,stdin='',extra=['--owner','fixture-owner']) == answered
    call('attention','reply',reply,stdin='',extra=['--owner','different-actor'],failure='idempotency_conflict')
    # Each call is a new process: reopened state must expose the captured answer.
    view = call('attention','list',extra=['--json'])
    assert view['items'][0]['status'] == 'answered'
    r = ask(answered, 'This deliberately obsolete question must be withdrawn.')
    stale = disposition(r,1)
    (root/'SDP/authority.md').write_text('Changed synthetic authority.\n')
    view = call('attention','list',extra=['--json'])
    assert all(v['applicability']=='stale' for v in view['items'])
    call('attention','reply',stale,stdin='short\nDo not use this answer.\n',extra=['--owner','fixture-owner'],failure='stale_revision')
    r = call('attention','withdraw',stale,extra=['--reason','Synthetic authority changed; retain history.'])
    (root/'SDP/authority.md').write_text(files['SDP/authority.md'])
    r = ask(r, 'Try the local mailbox: which label style would you choose? Your reply stays in this isolated preview.')
    view = call('attention','list',extra=['--json'])
    assert [v['status'] for v in view['items']] == ['answered','withdrawn','open']
    assert view['items'][-1]['applicability'] == 'current'
    before = (state/'HEAD').read_bytes()
    old = subprocess.run([str(a.old_controller.resolve()),'inbox','--project',str(root),'--state',str(state),'--json'],capture_output=True,text=True)
    assert old.returncode != 0 and (state/'HEAD').read_bytes() == before
    cases.append({'operation':'old-reader-refusal', 'exitCode':old.returncode, 'stderr':old.stderr})
    assert original == {name:hashlib.sha256((root/name).read_bytes()).hexdigest() for name in files}
    (base/'evidence.json').write_text(json.dumps({'fixtureOnly':True,'modelSubmitted':False,'protectedFiles':len(original),'protectedBytesUnchanged':True,'controllerSHA256':hashlib.sha256(controller.read_bytes()).hexdigest(),'sdptoolSHA256':hashlib.sha256(tool.read_bytes()).hexdigest(),'cases':cases},indent=2)+'\n')
    (base/'README.md').write_text('# Local attention trial\n\nSynthetic preview. No model started. The fixture contains one answered, one withdrawn and one open question. Replies record input; they do not launch a model or approve design. Use the attention helper to try the open question. evidence.json retains actual subprocess requests and results.\n')
    print('PASS: durable questions/replies, CLI capture/EOF/retry, stale refusal, withdrawal, independent run, old-reader refusal and protected bytes. One open question remains for an owner trial.')

if __name__ == '__main__':
    main()
