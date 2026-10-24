"""Read and retry the final compiled publication; validate canonical output and MCP.
No model calls. Requires the separately created disposable fixture from the Go test.
"""
import hashlib,json,os,runpy,subprocess
from pathlib import Path
root=Path(__file__).resolve().parents[5];base=Path(__file__).resolve().parent
cache=Path('/home/warloc/.cache/sdp-workstation/PGI4-C3');lab=Path('/home/warloc/git/pg_lab');project=lab/'publication-c3-final-schema';state=cache/'lab-state-final-schema'
controller=cache/'bin/sdp-controller';mcp=cache/'bin/sdp-mcp'
wire=json.loads((base/'lab-wire.json').read_text());publication=wire[-1]['publication'];ident=publication['request']['operationId']
original={p:hashlib.sha256((project/p).read_bytes()).hexdigest() for p in [x['path'] for x in publication['changes']]}
def command(args):
 r=subprocess.run([str(controller),*args],capture_output=True,text=True);assert r.returncode==0,r.stderr;return json.loads(r.stdout)
current=command(['publication','inspect','--project',str(project),'--state',str(state),'--publication',ident]);assert current['publication']['status']=='published' and current['applicable']
retried=command(['publication','apply','--project',str(project),'--state',str(state),'--publication',ident]);assert retried==current
coord=command(['coordinate','--project',str(project),'--state',str(state),'--json']);reviews=[x['review'] for x in coord['items'] if x['kind']=='review'];assert len(reviews)==1 and reviews[0]['publicationStatus']=='published' and reviews[0]['publicationApplicable'];assert not reviews[0]['qualified'] # Original authority bytes changed through recorded publication, not a new review.
(base/'lab-coordination.json').write_text(json.dumps(coord,indent=2)+'\n')
module=runpy.run_path(str(root/'SDP/ProjectManagement/validate.py'));events=[json.loads(x) for x in (project/'SDP/ProjectManagement/Ledger.ndjson').read_text().splitlines()]
counts=module['validate_current'](events,board=project/'SDP/KanBan',sdp=project/'SDP')
from jsonschema import Draft202012Validator,FormatChecker
schema=json.loads((root/'Toolkit/schemas/ledger-event.schema.json').read_text());ev=Draft202012Validator(schema,format_checker=FormatChecker())
for line in (project/'SDP/Traceability/Ledger.ndjson').read_text().splitlines():ev.validate(json.loads(line))
# Actual local stdio MCP exchange with the final binary, no Codex/app-server.
p=subprocess.Popen([str(mcp),'--project',str(project),'--state',str(state),'--coordinator','fixture-leader'],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
protocol=[]
def rpc(i,method,params):
 req={'jsonrpc':'2.0','id':i,'method':method,'params':params};p.stdin.write(json.dumps(req)+'\n');p.stdin.flush();line=p.stdout.readline();assert line,'MCP EOF';reply=json.loads(line);assert reply.get('id')==i and 'error' not in reply,reply;protocol.append({'request':req,'response':reply});return reply['result']
try:
 rpc(1,'initialize',{'protocolVersion':'2025-11-25','capabilities':{},'clientInfo':{'name':'pgi4-c3-fixture','version':'1'}})
 p.stdin.write(json.dumps({'jsonrpc':'2.0','method':'notifications/initialized'})+'\n');p.stdin.flush()
 tools=rpc(2,'tools/list',{});assert {x['name'] for x in tools['tools']}=={'coordination_read','report_acknowledge'}
 result=rpc(3,'tools/call',{'name':'coordination_read','arguments':{}});assert not result.get('isError');observed=result['structuredContent'];assert observed['schema']=='sdp-governance-coordination/2';assert [x['review'] for x in observed['items'] if x['kind']=='review']==reviews
finally:
 p.stdin.close();p.wait(timeout=10)
(base/'mcp-protocol.json').write_text(json.dumps(protocol,indent=2)+'\n')
# Observe actual retained C2 native reviews, without reopening or modifying their stores.
native=[]
for name in ['positive','negative']:
 ns=Path('/home/warloc/.cache/sdp-workstation/PGI4-C2')/('state-final-'+name);before=(ns/'HEAD').read_bytes()
 view=command(['coordinate','--project',str(lab),'--state',str(ns),'--json']);found=[x['review'] for x in view['items'] if x['kind']=='review'];assert len(found)==1;assert before==(ns/'HEAD').read_bytes();native.append({'case':name,'head':before.decode(),'review':found[0]})
(base/'retained-native-review-notices.json').write_text(json.dumps(native,indent=2)+'\n')
assert original=={p:hashlib.sha256((project/p).read_bytes()).hexdigest() for p in original}
baseline=json.loads((base/'lab-baseline.json').read_text());assert all(hashlib.sha256((lab/p).read_bytes()).hexdigest()==h for p,h in baseline['files'].items() if p not in ['pg','README.md'])
result={'publicationId':ident,'canonicalValidation':{'cards':counts[0],'managementRecords':counts[1],'lineageOperations':counts[2],'events':len(events)},'traceEvents':1,'retryFilesUnchanged':len(original),'priorLabFilesPreserved':len(baseline['files'])-2,'intentionalOwnerCommandUpdates':['pg','README.md'],'mcpTools':[t['name'] for t in tools['tools']],'newNativeModelCalls':0,'reviewProvenance':'publication fixture uses synthetic protocol observations; C2 native outcomes inspected read-only, not re-reviewed','binaries':{str(p):hashlib.sha256(p.read_bytes()).hexdigest() for p in [controller,mcp]}}
(base/'lab-checks.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result,indent=2))
