#!/usr/bin/env python3
"""Read-only Codex catalog probe in disposable repositories; no agent/model turns."""
import pathlib,tempfile,subprocess,json,selectors,shutil,time,hashlib
ROOT=pathlib.Path(__file__).resolve().parents[3]
out={'host':subprocess.check_output(['codex','--version'],text=True).strip(),'sourceCommit':'c4aed09944235f5da1d247001688f6f162ec124d','operation':'Local app-server initialize + skills/list only; no thread or model turn','cases':[]}
with tempfile.TemporaryDirectory(prefix='sdp-skills-probe-') as d:
 base=pathlib.Path(d)
 def query(repo):
  p=subprocess.Popen(['codex','app-server'],cwd=repo,stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.DEVNULL,text=True)
  sel=selectors.DefaultSelector();sel.register(p.stdout,selectors.EVENT_READ)
  def rpc(i,method,params):
   p.stdin.write(json.dumps({'id':i,'method':method,'params':params})+'\n');p.stdin.flush()
   until=time.monotonic()+20
   while time.monotonic()<until:
    if not sel.select(1):continue
    line=p.stdout.readline()
    if not line:raise RuntimeError('server exited')
    r=json.loads(line)
    if r.get('id')==i:
     if 'error' in r:raise RuntimeError(r['error'])
     return r['result']
   raise TimeoutError(method)
  try:
   rpc(1,'initialize',{'clientInfo':{'name':'sdp-skills-discovery-probe','version':'0.1'}})
   p.stdin.write('{"method":"initialized"}\n');p.stdin.flush()
   r=rpc(2,'skills/list',{'cwds':[str(repo)],'forceReload':True})
   return r
  finally:
   p.terminate()
   try:p.wait(timeout=5)
   except subprocess.TimeoutExpired:p.kill();p.wait()
   sel.close()
 for mode in ['root-only','symlink-candidate','legacy-frontmatter','root-skill-only']:
  repo=base/mode;repo.mkdir();subprocess.run(['git','init','-q',str(repo)],check=True)
  p=repo/'Skills';p.mkdir()
  source = 'Toolkit/skills/sdp-master' if mode=='legacy-frontmatter' else 'Toolkit/skills_v2/sdp'
  name = 'sdp-master' if mode=='legacy-frontmatter' else 'sdp'
  target = p/name;target.mkdir()
  data = subprocess.check_output(['git','show','c4aed09:'+source+'/SKILL.md'],cwd=ROOT)
  (target/'SKILL.md').write_bytes(data)
  if mode in ['symlink-candidate','legacy-frontmatter']:
   a=repo/'.agents/skills';a.mkdir(parents=True)
   name='sdp-master' if mode=='legacy-frontmatter' else 'sdp'
   (a/name).symlink_to('../../Skills/'+name,target_is_directory=True)
  if mode=='root-skill-only':shutil.copy2(p/'sdp/SKILL.md',repo/'SKILL.md')
  result=query(repo)
  items=result.get('data',[])
  filtered=[]
  for item in items:
   filtered.append({'skills':[{'name':s['name'],'path':s['path'].replace(str(repo),'$FIXTURE'),'scope':s.get('scope')} for s in item.get('skills',[]) if str(repo) in s.get('path','')], 'errors':[{'path':e['path'].replace(str(repo),'$FIXTURE'),'message':e['message']} for e in item.get('errors',[]) if str(repo) in e.get('path','')]})
  out['cases'].append({'case':mode,'result':filtered})
 out['limits']='Filtered to fixture-local catalog results. No project activation, model loading, routing behavior, other host or installer tested.'
print(json.dumps(out,indent=2))
