import os,subprocess,json,time,concurrent.futures,hashlib,datetime
from pathlib import Path
root=Path('/home/warloc/git/SDP');out=Path(__file__).parent;bin=out/'bin';bin.mkdir(exist_ok=True)
go='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go'
env=dict(os.environ,PATH=str(Path(go).parent)+os.pathsep+os.environ['PATH'],GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache')
records=[]
def run(label,args):
 start=time.time();p=subprocess.run(args,cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,text=True)
 log=out/(label+'.log');log.write_text(p.stdout)
 r={'label':label,'command':args,'exitCode':p.returncode,'seconds':round(time.time()-start,3),'log':str(log),'logSHA256':hashlib.sha256(log.read_bytes()).hexdigest(),'summary':p.stdout[-3500:]}
 print(label,p.returncode,flush=True);return r
records.append(run('build-current',[go,'-C','SDPTool','build','-trimpath','-ldflags','-X github.com/Hans-Einar/SDP/SDPTool.BuildVersion=development-WR2 -X github.com/Hans-Einar/SDP/SDPTool.BuildRevision=4174c68','-o',str(bin),'./cmd/sdptool','./cmd/sdp-controller','./cmd/sdp-mcp']))
if records[0]['exitCode']==0:
 with concurrent.futures.ThreadPoolExecutor(max_workers=3) as pool:
  jobs=[('sdptool-race',[go,'-C','SDPTool','test','-race','-count=1','./...']),('sdl-race',[go,'-C','SDL/go','test','-race','-count=1','./...']),('sdui-race',[go,'-C','SDUI/go','test','-race','-count=1','./...']),('bootstrap-race',[go,'-C','SDPTool/bootstrap','test','-race','-count=1','./...']),('vet-sdptool',[go,'-C','SDPTool','vet','./...']),('vet-sdl',[go,'-C','SDL/go','vet','./...']),('vet-sdui',[go,'-C','SDUI/go','vet','./...'])]
  for r in pool.map(lambda j:run(*j),jobs):records.append(r)
report={'candidate':subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip(),'startedSource':'4174c68','observedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'go':subprocess.check_output([go,'version'],text=True).strip(),'checks':records,'binaries':{p.name:{'path':str(p),'sha256':hashlib.sha256(p.read_bytes()).hexdigest()} for p in bin.iterdir() if p.is_file()}}
(out/'checks.json').write_text(json.dumps(report,indent=2)+'\n')
print('FINAL',str(out/'checks.json'),flush=True)
