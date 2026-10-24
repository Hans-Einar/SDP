"""Validate final observed core, CLI, coordination and MCP wire objects."""
import json
from pathlib import Path
from jsonschema import Draft202012Validator,FormatChecker
from referencing import Registry,Resource
root=Path(__file__).resolve().parents[5];base=Path(__file__).resolve().parent
folder=root/'SDPTool/governance'
registry=Registry().with_resources((p.as_uri(),Resource.from_contents(json.loads(p.read_text()))) for p in folder.glob('*.json') if '$schema' in json.loads(p.read_text()))
core=json.loads((folder/'schema.json').read_text());coord=json.loads((folder/'coordination.schema.json').read_text())
# Register relative references under the file URI; relative JSON references in this contract are local.
def validator(schema,name,file):
 selected={**schema,'$id':file.as_uri(),'$ref':'#/$defs/'+name}
 return Draft202012Validator(selected,registry=registry,format_checker=FormatChecker())
cases=json.loads((base/'core-wire.json').read_text())
state=Path('/home/warloc/.cache/sdp-workstation/PGI4-C3/lab-state-final')
snapshot=json.loads((state/'generations'/((state/'HEAD').read_text()+'.json')).read_text())
for receipt in snapshot['receipts'].values():
 if receipt['result']['operation']=='publication.prepare':cases.append({'definition':'Result','value':receipt['result']})
for item in cases:validator(core,item['definition'],folder/'schema.json').validate(item['value'])
for v in json.loads((base/'lab-wire.json').read_text()):validator(core,'PublicationView',folder/'schema.json').validate(v);cases.append({'definition':'PublicationView','value':v})
views=[json.loads((base/n).read_text()) for n in ['coordination-wire.json','lab-coordination.json']]
views.append(json.loads((base/'mcp-protocol.json').read_text())[-1]['response']['result']['structuredContent'])
for v in views:validator(coord,'CoordinationView',folder/'coordination.schema.json').validate(v)
result={'coreAndPublicationObjects':len(cases),'coordinationObjects':len(views),'total':len(cases)+len(views)}
(base/'schema-checks.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result))
