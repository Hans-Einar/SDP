"""One-time port oracle capture; not used by the Go implementation."""
import json
from pathlib import Path
import sys
import unittest
sys.path.insert(0,'experiments/design_core')
import design_core as dc
original=dc.parse
sources={}
def capture(source):
 if len(source)<=30000:sources[source]=None
 return original(source)
dc.parse=capture
suite=unittest.defaultTestLoader.discover('experiments/design_core',pattern='test_*.py')
result=unittest.TextTestRunner(verbosity=0).run(suite)
if not result.wasSuccessful():raise SystemExit(1)
dc.parse=original
cases=[]
for source in sources:
 model,diagnostics=dc.check(source)
 cases.append({'source':source,'ast':dc.to_json(model),'diagnostics':[{'code':d.code,'span':dc.to_json(d.span)} for d in diagnostics], 'canonical':dc.canonicalize(model) if model and not dc.validate(model) else None})
Path('SystemDesignLanguage/go/parser/testdata/python-port-cases.json').write_text(json.dumps(cases,ensure_ascii=False,indent=2)+'\n')
print(len(cases),'captured structural parser/validator cases')
