import json,subprocess,pathlib,time,re
root='/home/warloc/git/HEOS';tool='/tmp/sdptool-activity3-final';out=pathlib.Path('/tmp/activity3-producer-matrix');out.mkdir(exist_ok=True)
def cmd(*args):
 r=subprocess.run([tool,root,*args,'--json'],capture_output=True,text=True)
 if r.returncode: raise RuntimeError(r.stderr)
 return json.loads(r.stdout)
p=cmd('discover','--navigation','lazy-sdl');results=[];index=0
for source in p['sources']:
 if source['state']!='validated' or source['kind']!='sdl':continue
 tree=cmd('tree','--model-only','--model',source['id'],'--revision',source['revision'])
 queries=[];objects=set()
 for n in tree['nodes']:
  if n.get('target',{}).get('operation')!='select':continue
  if n['kind'] in ('viewpoint','collection','file'):queries.append((n['id'],n['label'],n['target']))
  elif '/collection/' in n['id']:
   group=n['id'].rsplit('/',1)[0]
   if group not in objects:queries.append((n['id'],n['label'],n['target']));objects.add(group)
 for node,label,target in queries:
  index+=1;dest=out/str(index);r=cmd('select','--model',target['model'],'--uri',target['uri'],'--revision',target['revision'],'--output',str(dest));text=pathlib.Path(r['entry']).read_text()
  stripped=re.sub(r'\]\([^)]*\)',']',text)
  if re.search(r'\b[0-9a-f]{32,64}\b',stripped):raise RuntimeError('visible machine identity '+node)
  results.append(dict(model=source['source'],node=node,label=label,uri=target['uri'],entry=r['entry'],sourceLinks=len(r.get('sourceLinks',[]))))
 print(source['source'],len(queries),'passed',flush=True)
(out/'results.json').write_text(json.dumps(dict(producer='7370e29a46b392b146f26e694f1da55153a09647',count=len(results),results=results),indent=2)+'\n')
print('PASS',len(results),'reading selections; full-source snapshots and no visible machine IDs',flush=True)
