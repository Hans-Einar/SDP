import subprocess,json,select,time,hashlib
from pathlib import Path
base=Path(__file__).parent;root=Path('/home/warloc/git/SDP');state=base/'readonly-mcp-state'
assert not state.exists()
proc=subprocess.Popen([str(base/'bin/sdp-mcp'),'--project',str(root),'--state',str(state)],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,bufsize=0)
frames=[];seq=0;buffer=b''
def call(method,params=None):
 global seq,buffer
 seq+=1;q={'jsonrpc':'2.0','id':seq,'method':method}
 if params is not None:q['params']=params
 proc.stdin.write((json.dumps(q)+'\n').encode());proc.stdin.flush()
 end=time.monotonic()+30
 while time.monotonic()<end:
  if b'\n' not in buffer:
   ready,_,_=select.select([proc.stdout],[],[],max(0,end-time.monotonic()))
   if not ready:break
   import os
   chunk=os.read(proc.stdout.fileno(),1048576)
   if not chunk:raise RuntimeError('server closed')
   buffer+=chunk
  while b'\n' in buffer:
   line,buffer=buffer.split(b'\n',1);r=json.loads(line)
   if r.get('id')==seq:frames.append({'request':q,'response':r});return r
 raise RuntimeError('timeout')
def tool(name,args):return call('tools/call',{'name':name,'arguments':args})
try:
 init=call('initialize',{'protocolVersion':'2025-11-25','capabilities':{},'clientInfo':{'name':'WR2-readonly-probe','version':'1'}})
 assert init['result']['protocolVersion']=='2025-11-25',init
 proc.stdin.write(b'{"jsonrpc":"2.0","method":"notifications/initialized"}\n');proc.stdin.flush()
 listed=call('tools/list',{});names=sorted(t['name'] for t in listed['result']['tools']);assert names==['context_read','route_evaluate','run_read'],names
 sel={'card':'SDP/KanBan/active/#056--Change--SDP-consolidation-and-PG-MCP-readiness.md','plan':'SDP/Maintenance/WorkstationReadiness/ExecutionPlan.md','session':'SDP/Sessions/session-#0013--SDP_consolidation_and_Ponsse_readiness.md','authority':'SDP/Maintenance/WorkstationReadiness/ExecutionPlan.md','step':'S5'}
 context=tool('context_read',sel);assert not context['result'].get('isError'),context
 route=tool('route_evaluate',{'schema':'sdp-project-governance/0.1','categories':['status']});assert route['result']['structuredContent']['disposition']=='no-work',route
 denied=tool('run_read',{});assert denied['result']['isError'],denied
 missing=tool('run_transition',{});assert 'error' in missing or missing.get('result',{}).get('isError'),missing
 assert not state.exists(),'inspection created state'
 result={'status':'passed','realStdioServer':True,'modelTaskSubmitted':False,'client':'direct JSON-RPC probe, not a Codex model turn','toolNames':names,'stateCreated':False,'frames':frames,'contextConditions':context['result']['structuredContent'].get('conditions')}
 print('PASS: real compiled MCP stdio; three tools; context/status; unbound/mutation refusal; no state created')
finally:
 proc.stdin.close()
 try:proc.wait(timeout=3)
 except subprocess.TimeoutExpired:proc.terminate();proc.wait(timeout=3)
(base/'mcp-readonly.json').write_text(json.dumps(result,indent=2)+'\n')
