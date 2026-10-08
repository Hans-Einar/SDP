from pathlib import Path
import subprocess,os,json,time,sys,hashlib
root=Path('/home/warloc/git/SDP-vNow');scripts=root/'SDP/05--Implementation/SDUI/Widgets/native';package=Path(sys.argv[1]);out=Path(sys.argv[2]);out.mkdir(exist_ok=False)
env=dict(os.environ,FYNE_THEME='light');rows=[]
jobs=[]
def add(family,script,variants):
 for v in variants:
  variant=v.removeprefix('nonmodal-');extra=['--nonmodal'] if v.startswith('nonmodal-') else []
  jobs.append((family+'-'+v,script,package/'libexec/sdui'/family,variant,extra))
add('previews','verify_wci4_previews.py',['forms','nonmodal-forms','geometry','fallback','lifetime'])
add('collections','verify_wci1.py',['basic','nested'])
add('panes','verify_wci2.py',['panes','split'])
add('commands','verify_wci2_commands.py',['commands','context','dialog','nonmodal-nested','lifecycle'])
add('values','verify_wci3_values.py',['boolean','numeric','slider','choice','forms'])
add('text','verify_wci3_text.py',['single','multiline','history','scroll','command_load','ime','nonmodal-forms'])
jobs.append(('helper-ime','verify_wci4_helper_ime.py',package/'libexec/xfmd/sdui-fyne',None,[]))
for name,script,binary,variant,extra in jobs:
 cmd=['python3',str(scripts/script),'--binary',str(binary),'--display',':191','--out',str(out/name)]
 if variant:cmd+=['--variant',variant]
 cmd+=extra
 if name in ('text-ime','helper-ime'):
  cmd=['python3',str(scripts/'with_ibus.py'),'--display',':191','--installation','/tmp/sdui-ime-probe/root','--components','/tmp/sdui-ime-probe/components','--out',str(out/(name+'-environment')),'--',*cmd]
 begin=time.monotonic()
 with (out/(name+'-launcher.log')).open('w') as f:r=subprocess.run(cmd,cwd=root,env=env,stdout=f,stderr=subprocess.STDOUT)
 row=dict(name=name,argv=cmd,cwd=str(root),exit=r.returncode,seconds=round(time.monotonic()-begin,3),FYNE_THEME='light',binarySHA256=hashlib.sha256(binary.read_bytes()).hexdigest(),harnessSHA256=hashlib.sha256((scripts/script).read_bytes()).hexdigest())
 rows.append(row);(out/'invocations.json').write_text(json.dumps(rows,indent=2)+'\n');print(name,r.returncode,flush=True)
 if r.returncode:raise SystemExit(r.returncode)
