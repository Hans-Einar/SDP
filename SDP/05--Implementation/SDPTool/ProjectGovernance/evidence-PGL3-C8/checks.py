import concurrent.futures,json,os,subprocess,time
from pathlib import Path
base=Path(__file__).resolve().parent
root=Path('/home/warloc/git/SDP')
go='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go'
env=os.environ.copy();env.update(GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache');env['PATH']=str(Path(go).parent)+os.pathsep+env['PATH']
cases=[('affected-race',[go,'-C','SDPTool','test','-race','-count=1','./governancecontroller','./governancehost','./cmd/sdp-controller']),('affected-vet',[go,'-C','SDPTool','vet','./governancecontroller','./governancehost','./cmd/sdp-controller'])]
def run(case):
 name,cmd=case;start=time.monotonic()
 with (base/(name+'.log')).open('w') as f:r=subprocess.run(cmd,cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT)
 result=dict(name=name,command=cmd,exitCode=r.returncode,seconds=round(time.monotonic()-start,2));print(json.dumps(result),flush=True);return result
with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:results=list(pool.map(run,cases))
(base/'checks.json').write_text(json.dumps(results,indent=2)+'\n')
raise SystemExit(any(r['exitCode'] for r in results))
