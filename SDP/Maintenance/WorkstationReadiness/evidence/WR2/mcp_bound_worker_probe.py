import subprocess,json,select,time,hashlib
from pathlib import Path
base=Path(__file__).parent;root=Path('/home/warloc/git/SDP');state=base/'bound-prelude-inspection-state'
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
import uuid
binding=json.loads((base/'bound-worker-binding.json').read_text());lab=base/'bound-worker-trial';root=lab/'project';state=lab/'runtime/state'
# Replace the initially created read-only server with the explicitly bound fixture server.
proc.stdin.close();proc.wait(timeout=3)
r=binding['run'];proc=subprocess.Popen([str(base/'bin/sdp-mcp'),'--project',str(root),'--state',str(state),'--run',r['run']['id'],'--attempt',binding['attemptID'],'--epoch',binding['epoch']],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,bufsize=0)
frames=[];seq=0;buffer=b''
pre=json.loads((lab/'preimages.json').read_text())
def preserved():return all(hashlib.sha256((root/name).read_bytes()).hexdigest()==digest for name,digest in pre.items())
assert preserved()
try:
 init=call('initialize',{'protocolVersion':'2025-11-25','capabilities':{},'clientInfo':{'name':'WR2-bound-worker-probe','version':'1'}})
 assert init['result']['protocolVersion']=='2025-11-25'
 proc.stdin.write(b'{"jsonrpc":"2.0","method":"notifications/initialized"}\n');proc.stdin.flush()
 listed=call('tools/list',{});names=sorted(t['name'] for t in listed['result']['tools']);assert names==sorted(['context_read','route_evaluate','run_read','events_read','run_transition','attention_request','attention_read']),names
 c=tool('context_read',r['run']['selection']);assert not c['result'].get('isError'),c
 current=tool('run_read',{});assert not current['result'].get('isError'),current
 data=current['result']['structuredContent']
 ev=tool('events_read',{'cursor':'','limit':10});assert not ev['result'].get('isError'),ev
 attention=tool('attention_read',{});assert not attention['result'].get('isError'),attention
 q={'schema':'sdp-project-governance/0.1','projectRef':data['projectRef'],'checkoutId':data['checkoutId'],'operationId':str(uuid.uuid4()),'runId':data['run']['id'],'expectedRevision':data['run']['revision'],'inputRevision':data['run']['inputRevision'],'action':'report','reason':'WR2 scripted fixture report only; no implementation, native agent or acceptance'}
 report=tool('run_transition',q);assert not report['result'].get('isError'),report
 retry=tool('run_transition',q);assert retry['result']==report['result'],'retry changed receipt'
 last=tool('run_read',{})['result']['structuredContent'];assert last['run']['state']=='verifying',last
 q['action']='complete';q['operationId']=str(uuid.uuid4());q['expectedRevision']=last['run']['revision'];refused=tool('run_transition',q);assert 'error' in refused or refused.get('result',{}).get('isError'),refused
 forbidden=tool('run_start',{});assert 'error' in forbidden or forbidden.get('result',{}).get('isError'),forbidden
 assert preserved(),'protected fixture bytes changed'
 result={'status':'passed','realStdioServer':True,'syntheticHostBinding':True,'modelTaskSubmitted':False,'ownerApprovalRecorded':False,'independentReview':False,'toolNames':names,'finalState':last['run']['state'],'duplicateReportSameReceipt':True,'completionRefused':True,'controllerStartNotExposed':True,'protectedFileCount':len(pre),'protectedFilesUnchanged':True,'frames':frames}
 print('PASS: seven real stdio worker tools; context/events/attention reads; idempotent report; verifying, not completed; completion/controller start refuse; protected fixture files unchanged')
finally:
 proc.stdin.close()
 try:proc.wait(timeout=3)
 except subprocess.TimeoutExpired:proc.terminate();proc.wait(timeout=3)
(base/'mcp-bound-worker.json').write_text(json.dumps(result,indent=2)+'\n')
