import sys,json
from pathlib import Path
sys.path[:0]=['SDUI/tests','SDUI/src']
import test_language
from sdui import parse,validate
from sdui.ast import to_data
from sdui.normalize import normalize
cases={}
def summary(roots):
 out={}
 def visit(n):
  out[n.path]={'kind':n.kind,'path':n.path,'widget':n.widget or '', 'variant':n.variant or '', 'text':n.text or '', 'layout':dict(n.layout),'arguments':{k:to_data(v) for k,v in n.arguments}, 'rows':[[c.path for c in row] for row in n.rows], 'regions':{k:v.path for k,v in n.regions}}
  for row in n.rows:
   for c in row:visit(c)
  for _,v in n.regions:visit(v)
 for root in roots.values():visit(root)
 return out
class Capture(test_language.LanguageTest):
 def valid(self,s):
  t=super().valid(s)
  if len(s)<16000:cases[s]={'source':s,'code':'','instances':summary(normalize(t))}
  return t
 def invalid(self,s,code=None):
  e=super().invalid(s,code)
  if len(s)<16000 and not any(0xD800<=ord(x)<=0xDFFF for x in s):cases[s]={'source':s,'code':e.code}
  return e
# Exercise original semantic cases. CLI behavior and resource stress have separate Go tests.
for name in sorted(dir(Capture)):
 if name.startswith('test_') and name not in {'test_cli_source_protection_and_exit_codes','test_truncations_and_malformed_inputs_are_structured'}:
  case=Capture(name);getattr(case,name)()
for p in Path('SDUI/examples').glob('*.sdui'):
 s=p.read_text();d=parse(s);validate(d);cases[s]={'source':s,'code':'','instances':summary(normalize(d))}
Path('SDUI/go/parser/testdata/python-port-cases.json').write_text(json.dumps(list(cases.values()),ensure_ascii=False,indent=2)+'\n')
print(len(cases),'captured parser/semantic cases')
