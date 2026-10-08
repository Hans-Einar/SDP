import os,subprocess,time,pathlib,signal
r=pathlib.Path('/tmp/sdui-ime-probe')
e=dict(os.environ,DISPLAY=':189',IBUS_ENABLE_CTRL_SHIFT_U='1',LC_ALL='en_US.UTF-8',XMODIFIERS='@im=ibus',IBUS_COMPONENT_PATH=str(r/'components'),LD_LIBRARY_PATH=str(r/'root/usr/lib64'),GSETTINGS_SCHEMA_DIR=str(r/'root/usr/share/glib-2.0/schemas'),XDG_CONFIG_HOME=str(r/'config'),XDG_CACHE_HOME=str(r/'cache'),XDG_RUNTIME_DIR=str(r/'runtime'),GSETTINGS_BACKEND='memory')
(r/'runtime').chmod(0o700)
ps=[]
def start(command,name):
 f=(r/(name+'.log')).open('w');p=subprocess.Popen(command,env=e,stdout=f,stderr=subprocess.STDOUT);ps.append(p);return p
try:
 daemon=start([str(r/'root/usr/bin/ibus-daemon'),'--single','--panel','disable','--emoji-extension','disable','--config','disable','--cache','none'],'ibus')
 for _ in range(100):
  q=subprocess.run([str(r/'root/usr/bin/ibus'),'address'],env=e,capture_output=True,text=True)
  if q.returncode==0 and q.stdout.strip().startswith('unix:'):break
  time.sleep(.1)
 print('BUS',q.stdout.strip(),flush=True)
 xim=start([str(r/'root/usr/libexec/ibus-x11'),'--debug=1','--locale=en_US.UTF-8'],'xim')
 time.sleep(.5)
 q=subprocess.run([str(r/'root/usr/bin/ibus'),'engine','xkb:us::eng'],env=e,capture_output=True,text=True);print('ENGINE',q.returncode,q.stdout,q.stderr,flush=True)
 probe=start([str(r/'app/probe-patched')],'probe')
 print('PROBE',probe.pid,flush=True)
 while probe.poll() is None:time.sleep(.25)
finally:
 for p in reversed(ps):
  if p.poll() is None:p.terminate()
 for p in ps:
  try:p.wait(timeout=3)
  except subprocess.TimeoutExpired:p.kill()
