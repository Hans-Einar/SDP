import concurrent.futures,json,os,subprocess,time
from pathlib import Path
base=Path(__file__).resolve().parent
root=Path('/home/warloc/git/SDP')
go='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go'
env=os.environ.copy();env.update(GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache');env['PATH']=str(Path(go).parent)+os.pathsep+env['PATH']
env.update(PG_EVIDENCE_WIRE=str(base/'core-wire.json'),PG_EVIDENCE_CLI_WIRE=str(base/'cli-wire.json'),PG_EVIDENCE_OLD_CONTROLLER='/home/warloc/.cache/sdp-workstation/PGL3-C8/bin/sdp-controller')
cases=[('sdptool-race',[go,'-C','SDPTool','test','-race','-count=1','./...']),('sdptool-vet',[go,'-C','SDPTool','vet','./...'])]
def run(case):
 name,cmd=case;start=time.monotonic()
 with (base/(name+'.log')).open('w') as f:r=subprocess.run(cmd,cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT)
 result=dict(name=name,command=cmd,exitCode=r.returncode,seconds=round(time.monotonic()-start,2));print(json.dumps(result),flush=True);return result
with concurrent.futures.ThreadPoolExecutor(max_workers=2) as pool:results=list(pool.map(run,cases))
(base/'checks.json').write_text(json.dumps(results,indent=2)+'\n')
raise SystemExit(any(r['exitCode'] for r in results))
