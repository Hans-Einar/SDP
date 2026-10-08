from pathlib import Path
import os,subprocess,time,json,sys
sys.path.insert(0,'/home/warloc/git/SDP-vNow/SDP/05--Implementation/SDUI/Widgets/native')
from package_wci4 import inventory,write
root=Path('/tmp/sdp-sdui-widgets'); out=Path(sys.argv[1]);out.mkdir(exist_ok=False)
env=dict(os.environ,GOWORK='off',GOFLAGS='-mod=readonly',GOCACHE='/tmp/sdp-wci0-gocache')
suites=[('sdui','SDUI/go',['go','test','-race','./...','-count=1']),('sdl-previews','SDL/go',['go','test','-race','./examples/previews/...','-count=1'])]
rows=[]
for name,sub,cmd in suites:
 before=inventory(root);write(out/(name+'-before.json'),before);start=time.monotonic()
 with (out/(name+'.log')).open('w') as f:r=subprocess.run(cmd,cwd=root/sub,env=env,stdout=f,stderr=subprocess.STDOUT)
 after=inventory(root);write(out/(name+'-after.json'),after)
 row=dict(name=name,argv=cmd,cwd=str(root/sub),environment={k:env[k] for k in ('GOWORK','GOFLAGS','GOCACHE')},exit=r.returncode,seconds=round(time.monotonic()-start,3),sourceUnchanged=before==after)
 rows.append(row);write(out/'commands.json',rows);print(name,row,flush=True)
 if r.returncode or before!=after:raise SystemExit(1)
