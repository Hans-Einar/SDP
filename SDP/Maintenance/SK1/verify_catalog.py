#!/usr/bin/env python3
"""Measure catalog discovery in fresh local hosts; no model or agent turns."""
import json,pathlib,selectors,subprocess,sys,tempfile,shutil,time,hashlib
ROOT=pathlib.Path(__file__).resolve().parents[3]
EXPECTED=sorted(p.parent.name for p in (ROOT/'Skills').glob('*/SKILL.md'))
def query(cwd):
    p=subprocess.Popen(['codex','app-server'],cwd=cwd,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.DEVNULL,text=True)
    sel=selectors.DefaultSelector();sel.register(p.stdout,selectors.EVENT_READ)
    def rpc(i,m,args):
        p.stdin.write(json.dumps({'id':i,'method':m,'params':args})+'\n');p.stdin.flush();until=time.monotonic()+20
        while time.monotonic()<until:
            if not sel.select(1):continue
            line=p.stdout.readline()
            if not line:raise RuntimeError('app-server exited')
            r=json.loads(line)
            if r.get('id')==i:
                if 'error' in r:raise RuntimeError(r['error'])
                return r['result']
        raise TimeoutError(m)
    try:
        rpc(1,'initialize',{'clientInfo':{'name':'sdp-skills-catalog-check','version':'1.0'}})
        p.stdin.write('{"method":"initialized"}\n');p.stdin.flush()
        return rpc(2,'skills/list',{'cwds':[str(cwd)],'forceReload':True})['data'][0]
    finally:
        p.terminate()
        try:p.wait(timeout=5)
        except subprocess.TimeoutExpired:p.kill();p.wait()
        sel.close()

def check(label,cwd,source):
    r=query(cwd);skills=[s for s in r['skills'] if s['name']=='sdp' or s['name'].startswith('sdp-')]
    assert sorted(s['name'] for s in skills)==EXPECTED,(label,[(s['name'],s['path']) for s in skills])
    assert not [e for e in r.get('errors',[]) if str(source) in e['path']],r.get('errors')
    for s in skills:
        assert pathlib.Path(s['path']).resolve()==(source/s['name']/'SKILL.md').resolve(),s
        assert s['enabled'] is True,s
    return {'case':label,'count':len(skills),'names':EXPECTED,'duplicates':False,'allEnabled':True,'sourcePathsMatch':True}

def git(repo,*args):return subprocess.run(['git',*args],cwd=repo,check=True,capture_output=True)
results=[check('project-root',ROOT,ROOT/'Skills'),check('nested-cwd',ROOT/'Toolkit/SDPTool',ROOT/'Skills')]
with tempfile.TemporaryDirectory(prefix='sk1-catalog-') as d:
    d=pathlib.Path(d);repo=d/'repository';repo.mkdir();git(repo,'init','-q')
    shutil.copytree(ROOT/'Skills',repo/'Skills');a=repo/'.agents/skills';a.mkdir(parents=True)
    for n in EXPECTED:(a/n).symlink_to('../../Skills/'+n,target_is_directory=True)
    git(repo,'add','.');git(repo,'-c','user.name=Fixture','-c','user.email=fixture@example.invalid','commit','-qm','Skills fixture')
    worktree=d/'worktree';git(repo,'worktree','add','-q','--detach',str(worktree))
    results.append(check('separate-git-worktree',worktree,worktree/'Skills'))
    project=d/'installed';project.mkdir();git(project,'init','-q')
    r=subprocess.run([sys.argv[1],'-NoProfile','-File',str(ROOT/'Toolkit/scripts/Install-SDP.ps1'),'-ProjectRoot',str(project)],capture_output=True,text=True)
    assert r.returncode==0,r.stdout+r.stderr
    results.append(check('actual-installer-codex-destination',project,project/'.codex/skills'))
    # Independent nested repository must not inherit its parent's local collection.
    sub=repo/'independent';sub.mkdir();git(sub,'init','-q');r=query(sub)
    names=[s['name'] for s in r['skills'] if s['name'] in EXPECTED]
    results.append({'case':'independent-nested-repository','inheritedNames':names})
    assert not names,names
hashes={str(p.relative_to(ROOT)):hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted((ROOT/'Skills').rglob('*.md')) if p.name!='README.md'}
print(json.dumps({'host':subprocess.check_output(['codex','--version'],text=True).strip(),'baseCommit':subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip(),'sourceHashes':hashes,'cases':results,'limit':'Catalog only; behavior and reference loading have separate fresh-context evaluation.'},indent=2))
