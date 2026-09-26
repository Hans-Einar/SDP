#!/usr/bin/env python3
"""Read an explicit XFMD checkout; install only into a disposable snapshot."""
import argparse
import hashlib
import importlib.util
import json
import os
from pathlib import Path
import shutil
import subprocess
import tempfile

ROOT=Path(__file__).resolve().parents[3]
s=importlib.util.spec_from_file_location('process_tests',ROOT/'Toolkit/tests/test_process_install.py')
m=importlib.util.module_from_spec(s);s.loader.exec_module(m)
p=argparse.ArgumentParser();p.add_argument('source',type=Path);p.add_argument('--provenance',required=True,type=Path);a=p.parse_args()
source=a.source.resolve()
before=subprocess.check_output(['git','-C',str(source),'status','--porcelain'])
commit=subprocess.check_output(['git','-C',str(source),'rev-parse','HEAD'],text=True).strip()
case=m.ProcessInstall();case.setUp()
try:
    copied={}
    for base,dirs,files in os.walk(source,followlinks=False):
        dirs[:]=[d for d in dirs if d not in {'.git','build','node_modules','.venv','vendor','.cache'} and not (Path(base)/d).is_symlink()]
        for name in files:
            src=Path(base)/name;rel=src.relative_to(source)
            if src.is_symlink() or not src.is_file():continue
            if rel.parts[0]!='SDP' and src.suffix.lower()!='.md':continue
            dest=case.root/rel;dest.parent.mkdir(parents=True,exist_ok=True);dest.write_bytes(src.read_bytes())
            copied[str(rel)]=hashlib.sha256(src.read_bytes()).hexdigest()
    history=(case.root/'SDP/Agents/KanBan/Ledger.ndjson').read_bytes()
    plan=case.plan('-ForceManagedFiles')
    assert plan['canApply'],plan['conflicts']
    result=json.loads(case.apply(plan,'-ForceManagedFiles').stdout)
    assert (case.root/'SDP/ProjectManagement/Ledger.ndjson').read_bytes().startswith(history)
    assert case.plan()['actions']==[]
    a.provenance.write_text(json.dumps(dict(source=str(source),commit=commit,dirty=bool(before),files=copied,result=result,status='passed'),indent=2)+'\n')
    assert subprocess.check_output(['git','-C',str(source),'status','--porcelain'])==before
    for rel,digest in copied.items():assert hashlib.sha256((source/rel).read_bytes()).hexdigest()==digest
    print(f'PASS: disposable XFMD snapshot {commit}, {len(copied)} files; original bytes/status unchanged')
finally:
    case.tearDown()
