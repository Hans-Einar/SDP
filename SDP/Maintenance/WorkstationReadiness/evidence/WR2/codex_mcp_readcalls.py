import subprocess,json,select,time,os
from pathlib import Path
base=Path(__file__).parent;root='/home/warloc/git/SDP';state=base/'codex-discovery-state'
proc=subprocess.Popen(['codex','app-server','--listen','stdio://'],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.DEVNULL,bufsize=0)
seq=0;buffer=b'';calls=[]
def call(method,params):
 global seq,buffer
 seq+=1;q={'id':seq,'method':method,'params':params};calls.append(method)
 proc.stdin.write((json.dumps(q)+'\n').encode());proc.stdin.flush();end=time.monotonic()+35
 while time.monotonic()<end:
  while b'\n' in buffer:
   line,buffer=buffer.split(b'\n',1);r=json.loads(line)
   if r.get('id')==seq:
    if 'error' in r:raise RuntimeError(json.dumps(r['error']))
    return r['result']
  r,_,_=select.select([proc.stdout],[],[],max(0,end-time.monotonic()))
  if not r:break
  chunk=os.read(proc.stdout.fileno(),1048576)
  if not chunk:raise RuntimeError('app-server closed')
  buffer+=chunk
 raise RuntimeError('host request timed out: '+method)
result={'probe':'independent read-only Codex MCP discovery, not the pinned PG controller','hostVersion':'0.164.0','modelTaskSubmitted':False,'persistentConfigChanged':False,'accountRead':False}
try:
 init=call('initialize',{'clientInfo':{'name':'wr2_discovery','version':'1'},'capabilities':{'experimentalApi':False}})
 result['userAgent']=init.get('userAgent')
 proc.stdin.write(b'{"method":"initialized","params":{}}\n');proc.stdin.flush()
 cfg={'pg_wr2_readonly':{'command':str(base/'bin/sdp-mcp'),'args':['--project',root,'--state',str(state)],'startup_timeout_sec':20,'enabled_tools':['context_read','route_evaluate','run_read']}}
 t=call('thread/start',{'cwd':root,'ephemeral':True,'sandbox':'read-only','approvalPolicy':'never','config':{'mcp_servers':cfg}})['thread'];result['threadID']=t['id'];result['threadStatus']=t.get('status')
 status=call('mcpServerStatus/list',{'limit':64,'serverName':'pg_wr2_readonly','threadId':t['id']})
 entry=next(x for x in status['data'] if x['name']=='pg_wr2_readonly');names=sorted(entry['tools']);result['tools']=names
 assert names==['context_read','route_evaluate','run_read'],entry

 selection={'card':'SDP/KanBan/active/#056--Change--SDP-consolidation-and-PG-MCP-readiness.md','plan':'SDP/Maintenance/WorkstationReadiness/ExecutionPlan.md','session':'SDP/Sessions/session-#0013--SDP_consolidation_and_Ponsse_readiness.md','authority':'SDP/Maintenance/WorkstationReadiness/ExecutionPlan.md','step':'S5'}
 def readtool(tool,args):return call('mcpServer/tool/call',{'threadId':t['id'],'server':'pg_wr2_readonly','tool':tool,'arguments':args})
 context=readtool('context_read',selection);assert not context.get('isError'),context
 result['contextConditions']=context['structuredContent']['conditions'];result['contextRevision']=context['structuredContent']['revision']
 route=readtool('route_evaluate',{'schema':'sdp-project-governance/0.1','categories':['status']});result['routeDisposition']=route['structuredContent']['disposition'];assert result['routeDisposition']=='no-work',route
 denied=readtool('run_read',{});result['unboundRunRefused']=denied['isError'];assert result['unboundRunRefused'],denied
 assert not state.exists(),'inspection unexpectedly created operational state'
 result['status']='passed';result['stateCreated']=False
except Exception as e:
 result['status']='failed';result['error']=str(e)
finally:
 result['methods']=calls;proc.stdin.close()
 try:proc.wait(timeout=2)
 except subprocess.TimeoutExpired:proc.terminate();proc.wait(timeout=3)
(base/'codex-mcp-readcalls.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result))
