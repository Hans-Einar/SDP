import hashlib,json,re,subprocess
from pathlib import Path
from urllib.parse import unquote,urlsplit
base=Path(__file__).resolve().parent
root=Path.cwd()
h=lambda p:hashlib.sha256(p.read_bytes()).hexdigest()
candidate=json.loads((base/'candidate.json').read_text())
for name,digest in candidate['sourceFiles'].items():
 assert h(root/name)==digest, name+' differs from tested candidate'
manifest=json.loads((base/'manifest.json').read_text())
for name,digest in manifest['artifacts'].items():
 assert h(base/name)==digest,name+' artifact changed'
prefixes={}
for name in ('SDP/ProjectManagement/Ledger.ndjson','SDP/Traceability/Ledger.ndjson'):
 old=subprocess.check_output(['git','show',candidate['productBaseline']+':'+name])
 current=(root/name).read_bytes();assert current.startswith(old),name+' rewrote history'
 rows=[json.loads(line) for line in current.splitlines()]
 assert len({row['eventId'] for row in rows})==len(rows),name+' duplicate event ID'
 prefixes[name]={'bytes':len(old),'sha256':hashlib.sha256(old).hexdigest(),'preserved':True}
changed=set(subprocess.check_output(['git','diff','--name-only','HEAD'],text=True).splitlines())
changed.update(subprocess.check_output(['git','ls-files','--others','--exclude-standard'],text=True).splitlines())
links=0
for name in changed:
 p=root/name
 if p.suffix!='.md' or not p.is_file():continue
 for match in re.finditer(r'(?<!!)\[[^\]]*\]\(([^)]+)\)',p.read_text()):
  target=match.group(1).split(' "',1)[0].strip('<>')
  if target.startswith(('#','http:','https:','mailto:')):continue
  path=unquote(urlsplit(target).path)
  if not path:continue
  assert (p.parent/path).exists(),name+' -> '+target
  links+=1
checks={}
for name,cmd in [('management',['python3','SDP/ProjectManagement/validate.py']),('toolkit',['python3','Toolkit/scripts/validate_sdp.py','--base-ref','main']),('whitespace',['git','diff','--check'])]:
 r=subprocess.run(cmd,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
 checks[name]={'exitCode':r.returncode,'output':r.stdout};assert r.returncode==0,r.stdout
result={'sourceFilesMatched':len(candidate['sourceFiles']),'artifactsMatched':len(manifest['artifacts']),'relativeLinks':links,'ledgerPrefixes':prefixes,'checks':checks}
(base/'final-record-checks.json').write_text(json.dumps(result,indent=2)+'\n')
print(json.dumps(result,indent=2))
