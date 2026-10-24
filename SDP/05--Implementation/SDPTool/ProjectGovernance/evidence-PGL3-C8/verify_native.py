import hashlib,json
from pathlib import Path
from jsonschema import Draft202012Validator,FormatChecker
from referencing import Registry,Resource
base=Path(__file__).resolve().parent
schemas=Path.cwd()/'SDPTool/governance'
registry=Registry()
for name in ('schema.json','leader.schema.json','coordination.schema.json','overview.schema.json','inbox.schema.json'):
 p=schemas/name;v=json.loads(p.read_text());r=Resource.from_contents(v);registry=registry.with_resource(p.as_uri(),r)
 if '$id' in v:registry=registry.with_resource(v['$id'],r)
checks=[]
def check(label,value,file,definition=''):
 uri=(schemas/file).as_uri()+('#/$defs/'+definition if definition else '')
 Draft202012Validator({'$ref':uri},registry=registry,format_checker=FormatChecker()).validate(value)
 checks.append({'object':label,'schema':file,'definition':definition,'passed':True})
for name,file,definition in [('leader-view','leader.schema.json','LeaderView'),('run','schema.json','Result'),('overview','overview.schema.json',''),('report-inbox','inbox.schema.json','')]:
 v=json.loads((base/'native'/f'{name}.json').read_text());check(name,v,file,definition)
 if name=='leader-view':
  for d in v['leader']['dispatches']:check('context '+d['id'],json.loads(d['payload']),'leader.schema.json','LeaderContext')
for p in sorted((base/'native').glob('leader-enqueue-*.json')):check(p.name,json.loads(p.read_text()),'leader.schema.json','LeaderResult')
summary=json.loads((base/'native-summary.json').read_text())
for path,digest in summary['sourceFiles'].items():assert hashlib.sha256(Path(path).read_bytes()).hexdigest()==digest,path
for i in (2,3):
 original=[json.loads(x) for x in (base/'native'/f'leader-{i}-protocol.ndjson').read_text().splitlines()]
 recovered=[json.loads(x) for x in (base/'native'/f'leader-{i}-recovery-protocol.ndjson').read_text().splitlines()]
 assert sum(r.get('direction')=='request' and r.get('method')=='turn/start' for r in original)==1
 assert [r['method'] for r in recovered if r['direction']=='request']==['initialize','thread/read']
 start=next(r for r in original if r.get('direction')=='response' and r.get('method')=='thread/start')
 assert start['model']==summary['model'] and start['reasoningEffort']==summary['reasoningEffort']
 assert start['thread']['cliVersion']==summary['host']
(base/'schema-checks.json').write_text(json.dumps(checks,indent=2)+'\n')
print(len(checks),'native wire/context schema checks and exact recovery/mapping checks passed')
