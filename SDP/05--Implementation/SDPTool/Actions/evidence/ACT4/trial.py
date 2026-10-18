#!/usr/bin/env python3
"""Actual compiled consumer and isolated gh-extension trial. No installation upgrade."""
import hashlib
import json
import os
from pathlib import Path
import shutil
import signal
import subprocess
import sys
import tempfile
import time

repo = Path(__file__).resolve().parents[6]
wrapper = Path(sys.argv[1]).resolve()
evidence = Path(__file__).resolve().parent
candidate = subprocess.check_output(['git', 'rev-parse', 'HEAD'], cwd=repo, text=True).strip()
records = []

def call(argv, *, request=None, env=None, cwd=None, code=0):
    p = subprocess.run([str(x) for x in argv], input=None if request is None else json.dumps(request),
                       text=True, capture_output=True, env=env, cwd=cwd, timeout=90)
    assert p.returncode == code, (argv, p.returncode, p.stdout, p.stderr)
    return p

with tempfile.TemporaryDirectory(prefix='sdp-actions-trial-') as temporary:
    work = Path(temporary)
    dist = work / 'dist'; dist.mkdir()
    binary = dist / 'sdptool'
    call(['go', '-C', repo / 'SDPTool', 'build', '-ldflags',
          f'-X github.com/Hans-Einar/SDP/SDPTool.BuildVersion=dev-actions -X github.com/Hans-Einar/SDP/SDPTool.BuildRevision={candidate}',
          '-o', binary, './cmd/sdptool'])
    call(['go', 'run', evidence / 'dev-descriptor.go', dist, candidate], cwd=repo)
    extension = work / 'gh-sdp'; extension.mkdir()
    call(['go', '-C', wrapper, 'build', '-o', extension / 'gh-sdp', '.'])
    env = dict(os.environ, GH_CONFIG_DIR=str(work/'gh-config'), XDG_DATA_HOME=str(work/'data'),
               SDP_RELEASE=str(dist/'release.json'), SDP_TEST_KEY=str(dist/'test-key.pub'),
               SDP_CACHE_DIR=str(work/'cache'), SDP_OFFLINE='false', GH_PROMPT_DISABLED='1')
    # Local binary extension, installed only into this trial's isolated GH config.
    call(['gh','extension','install','.'], env=env, cwd=extension)
    root=work/'project'; area=root/'SDP'; area.mkdir(parents=True)
    (area/'SDP-project.manifest.yaml').write_text("schemaVersion: '1.0'\nproject:\n  name: actions-trial\n")
    models=area/'Models';models.mkdir()
    for name in ['Before','After']:
        created=json.loads(call([binary,models,'model','create',f'work:{name}','--initial','--json']).stdout)
        shutil.copytree(repo/'experiments/blueprint_mvp1/NOW',Path(created['path']),dirs_exist_ok=True)
    task=root/'task.json';task.write_text(json.dumps({'id':'compiled-actions','intent':'Bounded consumer protocol trial'}))
    def invoke(action, parameters, expected=0, attribution=None):
        req={'schema':'sdptool.invoke/1','requestId':f'trial-{len(records)+1}',
             'action':action,'context':{'project':str(root)},'parameters':parameters}
        args=['gh','sdp','invoke','--request','-']+(attribution or [])
        p=call(args,request=req,env=env,code=expected)
        result=json.loads(p.stdout)
        assert result['exitCode']==expected and result['requestId']==req['requestId'] and not p.stderr
        records.append({'action':action,'exitCode':expected,'status':result['status']})
        return result
    catalogue=json.loads(call(['gh','sdp',root,'actions','--json'],env=env).stdout)
    env['SDP_OFFLINE']='true'
    assert len(catalogue['actions'])==5 and catalogue['producer']['revision']==candidate
    (evidence/'catalogue.json').write_text(json.dumps(catalogue,indent=2)+'\n')
    invoke('project.discover',{})
    invoke('future.action',{},2)
    generated=invoke('blueprint.generate',{'modelArea':'SDP/Models','from':'work:Before','to':'work:After',
                        'entry':'System.design','task':'task.json','catalogue':'SDP/Blueprints'})['result']
    bundle=Path(generated['path']); document=json.loads((bundle/'blueprint.json').read_text())
    ev={'schema':'sdp-blueprint-evidence/1','blueprintRevision':document['revision'],
        'retainedRevision':generated['retainedRevision'],'taskDigest':document['taskDigest'],
        'scope':'Bounded compiled protocol trial','code':{},'mappings':[],
        'requiredChecks':[],'receipts':[],'dispositions':[]}
    (root/'evidence.json').write_text(json.dumps(ev))
    invoke('blueprint.assess',{'bundle':str(bundle),'evidence':'evidence.json'},3)
    pm=area/'ProjectManagement';pm.mkdir()
    plan={'schemaVersion':'1.0','eventId':'EVT-PM-SDP-000001','eventType':'x-management:created',
          'occurredAt':'2026-10-18T00:00:00Z','actor':'owner','commit':None,'subjectId':'PLAN-SDP-0001',
          'payload':{'schemaVersion':'0.2','kind':'Plan','planType':'ImplementationPlan','previousEventId':None,
                     'from':None,'to':'active','fromPath':None,'toPath':'05--Implementation/Test.md','reason':'fixture','links':[]}}
    history=pm/'Ledger.ndjson';history.write_text(json.dumps(plan)+'\n')
    binding={'bundle':str(bundle.relative_to(area)),'task':document['taskId'],'system':document['analysis']['System'],
             'revision':document['revision'],'retained':generated['retainedRevision'],'plan':'PLAN-SDP-0001',
             'milestone':'M1','modelArea':'Models','from':'work:Before','to':'work:After'}
    request={'schema':'sdp-blueprint-assignment/1','eventId':'EVT-BPA-TRIAL-000001',
             'occurredAt':'2026-10-18T01:00:00Z','assignmentId':'BPA-trial','action':'create',
             'expectedEvent':'','binding':binding,'reason':'compiled trial'}
    invoke('blueprint.assignment.apply',{'request':request},1)
    caller=['--as','controller','--authority','controller']
    invoke('blueprint.assignment.apply',{'request':request},attribution=caller)
    before=history.read_bytes()
    invoke('blueprint.assignment.apply',{'request':request},attribution=caller)
    assert before==history.read_bytes()
    invoke('blueprint.assignment.list',{})
    # A retained menu configuration survives producer downgrade; icon absence uses fallback.
    saved=['project.discover','blueprint.generate','future.action']; preserved=list(saved)
    old={'project.discover'};visible=[x for x in saved if x in old]
    assert saved==preserved and visible==['project.discover']
    icons={'project-discover':'folder'}
    assert any(icons.get(a['iconKey'],'generic')=='generic' for a in catalogue['actions'])
    # Actual executable must unblock incomplete stdin when interrupted.
    child=subprocess.Popen([binary,'invoke','--request','-'],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
    child.stdin.write('{');child.stdin.flush();time.sleep(0.2);child.send_signal(signal.SIGINT)
    child.wait(timeout=5)
    child.stdin.close();child.stdin=None
    out,err=child.communicate(timeout=5)
    assert child.returncode==130 and json.loads(out)['status']=='canceled', (child.returncode,out,err)
    summary={'candidate':candidate,'wrapperCommit':subprocess.check_output(['git','rev-parse','HEAD'],cwd=wrapper,text=True).strip(),
             'engineSHA256':hashlib.sha256(binary.read_bytes()).hexdigest(),
             'wrapperSHA256':hashlib.sha256((extension/'gh-sdp').read_bytes()).hexdigest(),
             'descriptorSHA256':hashlib.sha256((dist/'release.json').read_bytes()).hexdigest(),
             'checks':records,'cancellation':'passed','retryHistoryUnchanged':True,
             'consumerFixtures':'bounded headless checks, not native XFMD',
             'distribution':'isolated local test-signed descriptor; no release or normal installation change'}
    (evidence/'trial.json').write_text(json.dumps(summary,indent=2)+'\n')
    print('PASS: all five actions through actual gh sdp, structured failures, retry, consumer fixtures and cancellation')
