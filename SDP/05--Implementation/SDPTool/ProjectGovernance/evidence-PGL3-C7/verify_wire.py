import json
from pathlib import Path
from jsonschema import Draft202012Validator, FormatChecker
from referencing import Registry, Resource
base=Path(__file__).resolve().parent
schemas=Path.cwd()/'SDPTool/governance'
registry=Registry()
for name in ('schema.json','leader.schema.json','coordination.schema.json','overview.schema.json','inbox.schema.json'):
 p=schemas/name;data=json.loads(p.read_text());resource=Resource.from_contents(data)
 registry=registry.with_resource(p.as_uri(),resource)
 if '$id' in data: registry=registry.with_resource(data['$id'],resource)
results=[]
def check(label,value,file,definition):
 validator=Draft202012Validator({'$ref':(schemas/file).as_uri()+'#/$defs/'+definition},registry=registry,format_checker=FormatChecker())
 validator.validate(value);results.append(dict(object=label,schema=file,definition=definition,passed=True))
for p in sorted((base/'wire').iterdir()):
 if p.suffix not in ('.json','.ndjson'):continue
 values=[json.loads(p.read_text())] if p.suffix=='.json' else [json.loads(line) for line in p.read_text().splitlines()]
 for i,value in enumerate(values):
  if value['schema']=='sdp-governance-coordination-event/1':check(f'{p.name}:{i}',value,'coordination.schema.json','CoordinationEvent')
  elif 'action' in value:check(f'{p.name}:{i}',value,'leader.schema.json','LeaderResult')
  else:
   check(f'{p.name}:{i}',value,'leader.schema.json','LeaderView')
   for j,d in enumerate(value['leader']['dispatches']):
    check(f'{p.name}:{i} payload {j}',json.loads(d['payload']),'leader.schema.json','LeaderContext')
(base/'schema-checks.json').write_text(json.dumps(results,indent=2)+'\n')
print(str(len(results))+' captured wire objects and context payloads passed')
