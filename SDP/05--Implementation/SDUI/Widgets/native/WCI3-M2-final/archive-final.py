from pathlib import Path
import importlib.util,hashlib,json,shutil,subprocess,tempfile,os
repo=Path('/home/warloc/git/SDP-vNow');native=repo/'SDP/05--Implementation/SDUI/Widgets/native'
spec=importlib.util.spec_from_file_location('archive_helpers','/tmp/archive-wci3-m2.py');a=importlib.util.module_from_spec(spec);spec.loader.exec_module(a)
source=Path('/tmp/wci3-m2-guard-final-evidence');binary=Path('/tmp/wci3-m2-guard-final-native');inv=Path('/tmp/wci3-m2-guard-candidate-inventory.json')
m=json.loads(inv.read_text());a.BINARY_SHA=a.sha(binary.read_bytes());a.INVENTORY_SHA=a.sha(inv.read_bytes());a.EXPECTED_FILES=len(m['files']);a.EXPECTED_CHECKS=92;a.RUNS=('command_load',*a.RUNS)
assert a.sha(Path('/tmp/wci3-m2-reviewer-guard-final-native').read_bytes())==a.BINARY_SHA
original=a.inventory_equivalence(inv,repo);checks={n:a.audit_run(source/n) for n in a.RUNS};assert sum(checks.values())==92
calls=json.loads((source/'invocations.json').read_text());assert [x['variant'] for x in calls]==list(a.RUNS) and all(x['exit']==0 and x['FYNE_THEME']=='light' for x in calls)
for c in calls:
 argv=c['argv'];assert argv[argv.index('--binary')+1]==str(binary)
 assert ('--nonmodal' in argv)==(c['variant']=='nonmodal-forms')
out=native/'WCI3-M2-final';assert not out.exists()
with tempfile.TemporaryDirectory(prefix='.wci3-m2-final-',dir=native) as temp:
 final=Path(temp)/out.name;final.mkdir();provenance=[]
 for n in a.RUNS:
  a.retained_tree(source/n,final/n,provenance);a.retain(source/(n+'-launcher.log'),final/(n+'-launcher.log'),provenance)
 a.retained_tree(source/'ime-environment',final/'ime-environment',provenance)
 for src,dst in [(source/'invocations.json','invocations.json'),(inv,'tested-source-inventory.json'),(Path('/tmp/wci3-m2-guard-suite-manifest.json'),'suite-manifest.json'),(Path('/tmp/wci3-m2-guard-delta.json'),'post-review-delta.json'),(native/'audit_wci2_results.py','audit_wci2_results.py'),(native/'verify_wci3_text.py','verify_wci3_text.py'),(native/'verify_wci3_values.py','verify_wci3_values.py'),(Path('/tmp/run-wci3-m2-guard-final.py'),'run-final.py'),(Path(__file__),'archive-final.py')]:a.retain(src,final/dst,provenance)
 result=subprocess.check_output(['python3',str(final/'audit_wci2_results.py'),*[str(final/n/'events.ndjson.gz') for n in a.RUNS]],text=True)
 receipts={Path(k).parent.name:v for k,v in json.loads(result).items()};assert len(receipts)==14 and sum(x['terminalResults'] for x in receipts.values())==5
 suites=a.audit_suites(Path('/tmp/wci3-m2-guard-suite-manifest.json'),final/'suites',provenance)
 a.write_json(final/'checks.json',checks);a.write_json(final/'terminal-results.json',receipts);a.write_json(final/'original-equivalence.json',original)
 (final/'build-info.txt').write_bytes(subprocess.check_output(['go','version','-m',str(binary)]))
 a.write_json(final/'binary-identity.json',{'sha256':a.BINARY_SHA,'coordinatorBinary':str(binary),'reviewerBinary':'/tmp/wci3-m2-reviewer-guard-final-native','sameBytes':True})
 (final/'README.md').write_text(f'''# WCI3-M2 final native evidence after interaction guard correction

Fourteen fresh native workflows pass 92 checks on binary `{a.BINARY_SHA}`.
Five published openings have exactly one terminal receipt each; all fixture stderr
is empty and each run closes its session/providers. Actual invocation logs retain
XTest, clipboard, UI-owner completion barriers, source identities, screenshots and
real configured IBus/XIM. Launcher/environment warnings are retained separately.
No visible marked preedit, Linux OS accessibility, Wayland or decoration claim.

The {a.EXPECTED_FILES}-path inventory matches both product integration and the original
workspace. The reviewer independently rebuilt identical bytes. New command_load
proves an actual shared command updates an extended read-only receiver; the earlier
0078 archive is preserved in ../WCI3-M2-before-interaction-fix with its late defect.
No superseded approval or failed pilot is treated as final acceptance.

Fresh bridge/values/commands/panes race tests cover the changed interaction bridge.
The unchanged SDUI race suite, text fixture, collections, non-Fyne and SDPTool
proofs retain their original candidate identities. post-review-delta.json records
why they remain applicable; no exact-new-binary claim is made for those prior runs.
The text fixture defaults and typed Commit/legacy Load/DialogAccept branches remain
unchanged. The new command-load flag only selects a native fixture source variant.

Earlier timeout evidence is explicitly a retranscribed tool-output excerpt, not
an original stdout file; cancelled duplicate runs remain diagnostic records in the
prior-candidate archive. The old host group is excluded from current proof.
Native event files are losslessly gzip-compressed with deterministic headers.

Independent final correction/archive review remains pending. Provider previews,
matching consumer package preparation and whole-card closeout remain WCI4.
''')
 a.retain(Path('/tmp/resume-wci3-m2-guard-final.py'),final/'resume-final.py',provenance)
 for f in Path('/tmp/wci3-m2-interaction-review-3er3wmck').iterdir():
  if f.is_file():a.retain(f,final/'independent-interaction-review'/f.name,provenance)
 a.retain(Path('/tmp/wci3-m2-guard-targeted.log'),final/'guard-targeted.log',provenance)
 for item in provenance:
  assert a.sha(Path(item['source']).read_bytes())==item['rawSHA256'];item['path']=str(Path(item['path']).relative_to(final))
 a.write_json(final/'retention-provenance.json',provenance);a.inventory_equivalence(inv,repo);a.final_inventory(final);os.rename(final,out)
print(json.dumps({'runs':14,'checks':92,'receipts':5,'files':a.EXPECTED_FILES,'binarySHA256':a.BINARY_SHA}))
