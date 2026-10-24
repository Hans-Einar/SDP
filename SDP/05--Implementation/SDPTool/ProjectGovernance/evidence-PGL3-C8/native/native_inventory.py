from pathlib import Path
import subprocess,json
base=Path(__file__).parent
rows=[]
for isolated in (False,True):
 p=subprocess.Popen([str(base/'codex-lab') if isolated else '/home/warloc/.local/bin/codex','app-server','--listen','stdio://'],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.DEVNULL)
 def call(i,method,params):
  p.stdin.write((json.dumps({'id':i,'method':method,'params':params})+'\n').encode());p.stdin.flush()
  while True:
   raw=p.stdout.readline(16<<20);v=json.loads(raw)
   if v.get('id')!=i:continue
   row={'isolated':isolated,'method':method,'bytes':len(raw),'error':v.get('error')}
   if method=='mcpServerStatus/list':row['servers']=[{'name':x['name'],'tools':len(x.get('tools',{}))} for x in v.get('result',{}).get('data',[])];row['nextCursor']=v.get('result',{}).get('nextCursor')
   rows.append(row);return v.get('result',{})
 try:
  call(1,'initialize',{'clientInfo':{'name':'sdp_governance','version':'C8-diagnostic'},'capabilities':{'experimentalApi':False}});p.stdin.write(b'{"method":"initialized","params":{}}\n');p.stdin.flush()
  t=call(2,'thread/start',{'cwd':'/home/warloc/git/pg_lab','ephemeral':True,'sandbox':'read-only','approvalPolicy':'never','config':{'mcp_servers':{'sdp_governance':{'command':str(base/'bin/sdp-mcp'),'args':['--project','/home/warloc/git/pg_lab','--state',str(base/'state'),'--coordinator','native-lab-leader']}}}})['thread']['id']
  call(3,'mcpServerStatus/list',{'limit':64,'threadId':t})
 finally:p.kill();p.wait()
(base/'native-inventory.json').write_text(json.dumps(rows,indent=2)+'\n');print(json.dumps(rows,indent=2))
