from pathlib import Path
import subprocess,os,json,time
root=Path('/home/warloc/git/SDP-vNow'); scripts=root/'SDP/05--Implementation/SDUI/Widgets/native'; out=Path('/tmp/wci3-m2-guard-final-evidence');out.mkdir(exist_ok=False)
env=dict(os.environ,FYNE_THEME='light'); rows=[]
for name in ['command_load','single','multiline','history','readonly','pasteguard','failures','retention','forms','required_empty','constraints','scroll','ime','nonmodal-forms']:
 cmd=['python3',str(scripts/'verify_wci3_text.py'),'--binary','/tmp/wci3-m2-guard-final-native','--display',':190','--out',str(out/name),'--variant','forms' if name=='nonmodal-forms' else name]
 if name=='nonmodal-forms':cmd+=['--nonmodal']
 if name=='ime':cmd=['python3',str(scripts/'with_ibus.py'),'--display',':190','--installation','/tmp/sdui-ime-probe/root','--components','/tmp/sdui-ime-probe/components','--out',str(out/'ime-environment'),'--',*cmd]
 start=time.monotonic()
 with (out/(name+'-launcher.log')).open('w') as f:r=subprocess.run(cmd,cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT)
 rows.append({'variant':name,'argv':cmd,'exit':r.returncode,'seconds':round(time.monotonic()-start,3),'FYNE_THEME':'light'})
 (out/'invocations.json').write_text(json.dumps(rows,indent=2)+'\n')
 print(name,r.returncode,flush=True)
 if r.returncode:raise SystemExit(r.returncode)
