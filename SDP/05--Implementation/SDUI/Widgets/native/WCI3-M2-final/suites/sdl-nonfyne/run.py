import subprocess,json,pathlib,hashlib,time,os
root=pathlib.Path('/tmp/wci3-m2-sdl-nonfyne-regressions')
module=pathlib.Path('/tmp/sdp-sdui-widgets/SDL/go')
rows=json.loads((root/'packages.json').read_text())
selected=[x['package'] for x in rows if not x['fyneTestClosure']]
primary=['./runtime','./codegen']
remaining=[p for p in selected if p not in primary]
paths=list(module.rglob('*.go'))+[module/'go.mod',module/'go.sum']
def hashes():return {str(p.relative_to(module)):hashlib.sha256(p.read_bytes()).hexdigest() for p in sorted(paths)}
before=hashes();(root/'source-before.json').write_text(json.dumps(before,indent=2)+'\n')
head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=module,text=True).strip()
version=subprocess.check_output(['go','version'],cwd=module,text=True).strip()
flags=subprocess.check_output(['go','env','GOFLAGS'],cwd=module,text=True).strip()
env=os.environ.copy();env['GOFLAGS']=(flags+' -mod=readonly').strip()
meta={'cwd':str(module),'head':head,'goVersion':version,'GOFLAGS':env['GOFLAGS'],'selected':selected,'excluded':[r for r in rows if r['fyneTestClosure']],'runs':[]}
(root/'results.json').write_text(json.dumps(meta,indent=2)+'\n')
for name,packages in [('runtime-codegen',primary),('remaining-nonfyne',remaining)]:
 cmd=['go','test','-race','-count=1','-p','1']+packages
 (root/(name+'.command.json')).write_text(json.dumps({'argv':cmd,'cwd':str(module),'GOFLAGS':env['GOFLAGS']},indent=2)+'\n')
 started=time.time()
 with (root/(name+'.log')).open('w') as log:
  p=subprocess.Popen(cmd,cwd=module,env=env,stdout=log,stderr=subprocess.STDOUT)
  code=p.wait()
 duration=time.time()-started
 (root/(name+'.exit')).write_text(str(code)+'\n')
 entry={'name':name,'command':cmd,'exit':code,'seconds':round(duration,3),'logSHA256':hashlib.sha256((root/(name+'.log')).read_bytes()).hexdigest()}
 meta['runs'].append(entry);(root/'results.json').write_text(json.dumps(meta,indent=2)+'\n')
 print(json.dumps(entry),flush=True)
 print((root/(name+'.log')).read_text(),flush=True)
after=hashes();(root/'source-after.json').write_text(json.dumps(after,indent=2)+'\n')
meta['sourceBytesUnchanged']=before==after
meta['changedSourcePaths']=[p for p in before if before[p]!=after.get(p)]
(root/'results.json').write_text(json.dumps(meta,indent=2)+'\n')
print(json.dumps({'sourceBytesUnchanged':meta['sourceBytesUnchanged'],'changedSourcePaths':meta['changedSourcePaths']}),flush=True)
raise SystemExit(1 if any(r['exit'] for r in meta['runs']) else 0)
