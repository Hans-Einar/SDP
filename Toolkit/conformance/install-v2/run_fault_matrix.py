#!/usr/bin/env python3
"""Force process exit at every modeled boundary, in independent parallel fixtures."""
from concurrent.futures import ThreadPoolExecutor
import importlib.util
import json
import os
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
spec=importlib.util.spec_from_file_location('process_tests',ROOT/'Toolkit/tests/test_process_install.py')
m=importlib.util.module_from_spec(spec);spec.loader.exec_module(m)
WORKERS=4

def run(part):
    case=m.ProcessInstall()
    case.setUp()
    try:
        case.xfmd()
        case.put('AGENTS.md','Owner instructions to preserve.\n')
        history=(case.root/'SDP/Agents/KanBan/Ledger.ndjson').read_bytes()
        plan=case.plan('-ForceManagedFiles')
        env=dict(os.environ,SDP_INSTALL_INTERRUPT='prepared:0',SDP_INSTALL_HARD_EXIT='1')
        case.apply(plan,'-ForceManagedFiles',env=env,ok=False)
        path=next((case.root/'SDP/.sdp-operations').glob('*/journal.json'))
        j=json.loads(path.read_text());ident=j['operationId'];count=len(j['steps'])
        start=count*part//WORKERS;end=count*(part+1)//WORKERS
        if start:
            env=dict(os.environ,SDP_INSTALL_INTERRUPT=f'backup:{start}',SDP_INSTALL_HARD_EXIT='1')
            assert case.call('-ResumeOperation',ident,env=env,ok=False).returncode==97
        for i in range(start,end):
            for boundary in ('backup','write','journal'):
                env=dict(os.environ,SDP_INSTALL_INTERRUPT=f'{boundary}:{i}',SDP_INSTALL_HARD_EXIT='1')
                p=case.call('-ResumeOperation',ident,env=env,ok=False)
                assert p.returncode==97,(i,boundary,p.stderr)
            if (i-start)%5==0: print(f'partition {part}: PASS step {i}/{end-1}',flush=True)
        env=dict(os.environ,SDP_INSTALL_INTERRUPT=f'complete:{count}',SDP_INSTALL_HARD_EXIT='1')
        case.call('-ResumeOperation',ident,env=env,ok=False)
        case.call('-ResumeOperation',ident)
        assert json.loads(path.read_text())['status']=='completed'
        final=(case.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes()
        assert final.startswith(history)
        events=[json.loads(x) for x in final.splitlines()]
        assert len(events)==3 and len({e['eventId'] for e in events})==3
        case.call('-ResumeOperation',ident)
        assert (case.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes()==final
        assert case.plan()['actions']==[]
        return end-start
    finally:
        case.tearDown()

with ThreadPoolExecutor(max_workers=WORKERS) as executor:
    total=sum(executor.map(run,range(WORKERS)))
print(f'PASS: all {total} steps at backup/write/journal boundaries plus prepare/completion exits; forward resume, preserved history and exactly-once finalization',flush=True)
