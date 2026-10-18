from pathlib import Path
import json,hashlib,sys
old=Path('/tmp/wci4-final-native2');new=Path('/tmp/wci4-final-native3')
rows=[]
for parent in (old,new,Path('/tmp/wci4-final-native3-supplement')):
 for run in json.loads((parent/'invocations.json').read_text()):
  name=run['name']
  if parent==old and (name.startswith('previews-') or name in ('text-ime','helper-ime')):continue
  if parent==new and name=='previews-origins':continue
  assert run['exit']==0,(parent,name,run['exit'])
  root=parent/name
  actions=[json.loads(x) for x in (root/'actions.ndjson').read_text().splitlines()]
  checks=[x['value'] for x in actions if x['kind']=='check']
  assert checks and all(x['passed'] for x in checks),(name,'checks')
  assert any(x['kind']=='result' and x['value']=='passed' for x in actions),(name,'result')
  row=dict(name=name,evidence=str(root),candidate='unchanged-mechanism2' if parent==old else 'corrected3',checks=len(checks),invocation=run)
  if name!='helper-ime':
   events=[json.loads(x) for x in (root/'events.ndjson').read_text().splitlines()]
   closed=[e['data'] for e in events if e['event']=='closed'];assert len(closed)==1 and closed[0]['sessionClosed'] and not closed[0]['pending'],(name,'closed')
   def key(t):return(t['Handle']['Session'],t['Handle']['Path'],t['Handle']['Generation'],t['Handle']['Kind'],t['ModelRevision'],t['OpenGeneration'])
   opens=set()
   for e in events:
    if e['event']=='state':
     for s in e['data'].get('snapshot',{}).get('Surfaces',{}).values():
      if s['Open']:opens.add(key(s['Target']))
   results=[key(e['data']['Surface']) for e in events if e['event']=='dialog-result']
   assert len(results)==len(set(results)) and set(results)==opens,(name,'terminal',opens,results)
   row.update(publishedOpenings=len(opens),uniqueTerminalResults=len(results),closed=closed[0])
  stderr=root/'stderr.txt'
  row['stderr']=stderr.read_text() if stderr.exists() else 'see helper-specific log'
  rows.append(row)
assert len(rows)==28 and len({x['name'] for x in rows})==28
out=dict(workflows=len(rows),freshWorkflows=sum(x['candidate']=='corrected3' for x in rows),applicablePredecessorWorkflows=sum(x['candidate']=='unchanged-mechanism2' for x in rows),checks=sum(x['checks'] for x in rows),publishedOpenings=sum(x.get('publishedOpenings',0) for x in rows),terminalResults=sum(x.get('uniqueTerminalResults',0) for x in rows),rows=rows)
Path(sys.argv[1]).write_text(json.dumps(out,indent=2)+'\n');print({k:v for k,v in out.items() if k!='rows'})
