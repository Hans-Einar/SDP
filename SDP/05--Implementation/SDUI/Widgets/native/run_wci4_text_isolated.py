from pathlib import Path
import os,sys,subprocess,json,time
sys.path.insert(0,'/home/warloc/git/SDP-vNow/SDP/05--Implementation/SDUI/Widgets/native')
from package_wci4 import inventory,write
root=Path('/tmp/sdp-sdui-widgets');out=Path('/tmp/wci4-text-isolated2');out.mkdir(exist_ok=False);env=dict(os.environ,GOWORK='off',GOFLAGS='-mod=readonly',GOCACHE='/tmp/sdp-wci0-gocache')
before=inventory(root);write(out/'before.json',before);cmd=['go','test','-race','./examples/text','-count=1'];start=time.monotonic()
with (out/'test.log').open('w') as f:r=subprocess.run(cmd,cwd=root/'SDL/go',env=env,stdout=f,stderr=subprocess.STDOUT)
after=inventory(root);write(out/'after.json',after);record=dict(argv=cmd,cwd=str(root/'SDL/go'),environment={k:env[k] for k in ('GOWORK','GOFLAGS','GOCACHE')},exit=r.returncode,seconds=time.monotonic()-start,sourceUnchanged=before==after)
write(out/'receipt.json',record);print(record,flush=True);raise SystemExit(0 if r.returncode==0 and before==after else 1)
