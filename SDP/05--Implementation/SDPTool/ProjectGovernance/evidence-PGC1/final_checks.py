import os,subprocess,json,time,hashlib,concurrent.futures,datetime
from pathlib import Path
root=Path('/home/warloc/git/SDP');base=Path(__file__).parent;go='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go';binary=base/'bin-final';binary.mkdir(exist_ok=True)
env=dict(os.environ,PATH=str(Path(go).parent)+os.pathsep+os.environ['PATH'],GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache',PGI_SCHEMA_EVIDENCE=str(base/'schema-cases-final'))
head=subprocess.check_output(['git','rev-parse','HEAD'],cwd=root,text=True).strip();assert head.startswith('98b449c')
def run(name,args):
 start=time.time();p=subprocess.run(args,cwd=root,env=env,capture_output=True,text=True);log=base/(name+'.log');log.write_text(p.stdout+p.stderr)
 r={'name':name,'command':args,'exitCode':p.returncode,'seconds':round(time.time()-start,3),'log':str(log),'logSHA256':hashlib.sha256(log.read_bytes()).hexdigest()};print(name,p.returncode,flush=True);return r
jobs=[('full-sdptool-race',[go,'-C','SDPTool','test','-race','-count=1','./...']),('full-sdptool-vet',[go,'-C','SDPTool','vet','./...']),('toolkit-validator',['python','Toolkit/scripts/validate_sdp.py','--base-ref','main']),('toolkit-tests',['python','-m','unittest','discover','-s','Toolkit/tests','-p','test_*.py']),('management-tests',['python','-m','unittest','discover','-s','SDP/ProjectManagement','-p','test_*.py']),('management-validator',['python','SDP/ProjectManagement/validate.py']),('history-import',['python','SDP/ProjectManagement/History/Session0011/verify.py'])]
with concurrent.futures.ThreadPoolExecutor(max_workers=3) as pool:records=list(pool.map(lambda j:run(*j),jobs))
records.append(run('final-build',[go,'-C','SDPTool','build','-trimpath','-ldflags','-X github.com/Hans-Einar/SDP/SDPTool.BuildVersion=development-PGC1 -X github.com/Hans-Einar/SDP/SDPTool.BuildRevision='+head,'-o',str(binary),'./cmd/sdp-controller','./cmd/sdp-mcp']))
if records[-1]['exitCode']==0:
 records.append(run('live-codex-0.164-probe',[str(binary/'sdp-controller'),'probe','--project',str(root),'--codex','/home/warloc/.local/bin/codex','--mcp-binary',str(binary/'sdp-mcp'),'--state',str(base/'final-no-state')]))
records.append(run('prepare-no-model-lab',['python','SDPTool/labs/project-governance/prepare.py','--source',str(root),'--revision',head,'--destination',str(base/'lab-final'),'--go',go]))
r={'candidate':head,'observedAt':datetime.datetime.now(datetime.timezone.utc).isoformat(),'go':subprocess.check_output([go,'version'],text=True).strip(),'checks':records,'binaries':{p.name:hashlib.sha256(p.read_bytes()).hexdigest() for p in binary.iterdir() if p.is_file()},'modelTaskSubmitted':False,'independentReview':False,'nativeFullWorkflow':'not verified'}
(base/'final-checks.json').write_text(json.dumps(r,indent=2)+'\n');print('FINAL',flush=True)
