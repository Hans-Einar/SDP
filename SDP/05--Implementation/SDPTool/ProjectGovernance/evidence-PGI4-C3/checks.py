"""Final candidate build/race/vet and observed wire capture. No model calls."""
import hashlib,json,os,subprocess,time
from pathlib import Path
root=Path(__file__).resolve().parents[5]
base=Path(__file__).resolve().parent
cache=Path('/home/warloc/.cache/sdp-workstation')
go=cache/'go1.27.1/go/bin/go'
env={**os.environ,'PATH':str(go.parent)+':'+os.environ.get('PATH',''),'GOTOOLCHAIN':'local','GOCACHE':str(cache/'gocache'),'GOMODCACHE':str(cache/'gomodcache'),'PG_PUBLICATION_WIRE':str(base/'core-wire.json'),'PG_PUBLICATION_COORDINATION_WIRE':str(base/'coordination-wire.json'),'PG_PUBLICATION_INSPECT_ROOT':str(root)}
checks=[]
for name,args in [('build-controller',['build','-o',str(cache/'PGI4-C3/bin/sdp-controller'),'./cmd/sdp-controller']),('build-mcp',['build','-o',str(cache/'PGI4-C3/bin/sdp-mcp'),'./cmd/sdp-mcp']),('sdptool-race',['test','-race','-count=1','./...']),('sdptool-vet',['vet','./...'])]:
 start=time.monotonic()
 with (base/(name+'.log')).open('w') as out:
  r=subprocess.run([str(go),'-C',str(root/'SDPTool'),*args],env=env,stdout=out,stderr=subprocess.STDOUT)
 checks.append({'name':name,'command':[str(go),'-C',str(root/'SDPTool'),*args],'exitCode':r.returncode,'elapsedSeconds':round(time.monotonic()-start,3)})
 (base/'checks.json').write_text(json.dumps(checks,indent=2)+'\n')
 print(name,r.returncode,flush=True)
 if r.returncode:raise SystemExit(r.returncode)
