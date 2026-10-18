from pathlib import Path
import os,sys,subprocess,json,hashlib,time,shutil
stage=Path(sys.argv[1]).resolve();out=Path(sys.argv[2]).resolve();out.mkdir(exist_ok=False)
consumer=Path('/home/warloc/git/xfmd-sdl-navigation');original=consumer/'build-ui/SduiWorkflowGuiTest';target=stage/'bin/SduiWorkflowGuiTest'
def sha(p):return hashlib.sha256(p.read_bytes()).hexdigest()
assert not target.exists(), target
shutil.copy2(original,target)
files=['tools/build_sdui_tools.py','tools/run_with_xvfb.py','cmake/Sdp.cmake','cmake/Tests.cmake','src/application/sdp/SduiTools.cpp','src/application/sdp/SduiWorkflow.cpp','src/application/sdp/SdpProtocol.cpp','tests/gui/SduiWorkflowGuiTest.cpp']
before={n:sha(consumer/n) for n in files}
record=dict(consumerRoot=str(consumer),consumerHEAD=subprocess.check_output(['git','rev-parse','HEAD'],cwd=consumer,text=True).strip(),originalBinary=str(original),binarySHA256=sha(original),binaryBuildSource='Unknown exact C++ build source. This verifies the actual supplied binary, not current consumer source or release.',sourceBefore=before)
env=dict(os.environ,FYNE_THEME='light',XFMD_SDP_TOOL=str(stage/'bin/sdptool'),XFMD_UI_EVIDENCE=str(out))
for name in ('XFMD_SDUI_FYNE','XFMD_SDUI_PREVIEW'):env.pop(name,None)
argv=['strace','-f','-e','trace=execve','-s','8192','-o',str(out/'execve.log'),'python3',str(consumer/'tools/run_with_xvfb.py'),str(target)]
start=time.monotonic()
with (out/'consumer.log').open('w') as f:r=subprocess.run(argv,cwd=consumer,env=env,stdout=f,stderr=subprocess.STDOUT)
target.unlink() # Remove only this test's owned copy; it is not part of the prepared payload.
record.update(argv=argv,exit=r.returncode,seconds=round(time.monotonic()-start,3),environment={k:env[k] for k in ('FYNE_THEME','XFMD_SDP_TOOL','XFMD_UI_EVIDENCE')},helperOverrides='Both absent; executable-relative discovery required.',sourceAfter={n:sha(consumer/n) for n in files},binaryAfterSHA256=sha(original))
trace=(out/'execve.log').read_text();required=[stage/'bin/sdptool',stage/'libexec/xfmd/sdui-preview',stage/'libexec/xfmd/sdui-fyne']
record['executedStagedPaths']={str(p):('execve("'+str(p)+'",') in trace for p in required}
(out/'receipt.json').write_text(json.dumps(record,indent=2)+'\n')
print(json.dumps({k:record[k] for k in ('exit','seconds','executedStagedPaths')},indent=2),flush=True)
assert r.returncode==0
assert record['sourceBefore']==record['sourceAfter']
assert record['binarySHA256']==record['binaryAfterSHA256']
assert all(record['executedStagedPaths'].values())
