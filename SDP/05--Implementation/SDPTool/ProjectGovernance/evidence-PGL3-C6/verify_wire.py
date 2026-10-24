import json
from pathlib import Path
from jsonschema import Draft202012Validator, FormatChecker
from referencing import Registry, Resource
base=Path(__file__).parent
schemas=Path.cwd()/'SDPTool/governance'
registry=Registry()
for name in ('schema.json','coordination.schema.json','overview.schema.json','inbox.schema.json'):
 p=schemas/name;data=json.loads(p.read_text());resource=Resource.from_contents(data)
 registry=registry.with_resource(p.as_uri(),resource)
 if '$id' in data: registry=registry.with_resource(data['$id'],resource)
results=[]
def check(label,value,file,definition):
 validator=Draft202012Validator({'$ref':(schemas/file).as_uri()+'#/$defs/'+definition},registry=registry,format_checker=FormatChecker())
 validator.validate(value);results.append(dict(object=label,schema=file,definition=definition,passed=True))
wire=json.loads((base/'mcp-wire.json').read_text())
check('MCP view',wire['view'],'coordination.schema.json','CoordinationView')
check('MCP acknowledgement request',wire['acknowledgement'],'schema.json','ReportAcknowledgementRequest')
check('MCP result',wire['result'],'schema.json','Result')
for i,event in enumerate(json.loads((base/'console-wire.json').read_text())):
 check('console event '+str(i),event,'coordination.schema.json','CoordinationEvent')
(base/'schema-checks.json').write_text(json.dumps(results,indent=2)+'\n')
print(str(len(results))+' captured wire objects passed')
