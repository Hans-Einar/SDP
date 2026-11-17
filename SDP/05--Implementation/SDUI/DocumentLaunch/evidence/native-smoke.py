import os,subprocess,time,pathlib,ctypes as C,json,signal
base=pathlib.Path('/tmp/sdp-document-launch');binary=base/'package/sdptool';fixture=base/'timberframe';source=fixture/'TimberFrame.sdui';ctx=fixture/'Analysis.context.json';original=source.read_bytes();originalctx=ctx.read_bytes()
x=C.CDLL('libX11.so.6');x.XOpenDisplay.argtypes=[C.c_char_p];x.XOpenDisplay.restype=C.c_void_p;x.XDefaultRootWindow.argtypes=[C.c_void_p];x.XDefaultRootWindow.restype=C.c_ulong;x.XQueryTree.argtypes=[C.c_void_p,C.c_ulong,C.POINTER(C.c_ulong),C.POINTER(C.c_ulong),C.POINTER(C.POINTER(C.c_ulong)),C.POINTER(C.c_uint)];x.XFetchName.argtypes=[C.c_void_p,C.c_ulong,C.POINTER(C.c_void_p)];x.XFree.argtypes=[C.c_void_p];x.XInternAtom.argtypes=[C.c_void_p,C.c_char_p,C.c_int];x.XInternAtom.restype=C.c_ulong;x.XSendEvent.argtypes=[C.c_void_p,C.c_ulong,C.c_int,C.c_long,C.c_void_p];x.XFlush.argtypes=[C.c_void_p];x.XCloseDisplay.argtypes=[C.c_void_p]
class Data(C.Union):_fields_=[('b',C.c_char*20),('s',C.c_short*10),('l',C.c_long*5)]
class Client(C.Structure):_fields_=[('type',C.c_int),('serial',C.c_ulong),('send_event',C.c_int),('display',C.c_void_p),('window',C.c_ulong),('message_type',C.c_ulong),('format',C.c_int),('data',Data)]
class Event(C.Union):_fields_=[('client',Client),('pad',C.c_long*24)]
r,w=os.pipe();xf=subprocess.Popen(['/usr/bin/Xvfb','-displayfd',str(w),'-screen','0','1800x1400x24','-nolisten','tcp'],pass_fds=(w,),stderr=open(base/'xvfb.log','w'));os.close(w);app=None;display=None
try:
 number=os.read(r,32).decode().strip();os.close(r);env=os.environ.copy();env.update(DISPLAY=':'+number,LIBGL_ALWAYS_SOFTWARE='1',XDG_CACHE_HOME=str(base/'native-cache'),XDG_CONFIG_HOME=str(base/'native-config'),HOME=str(base/'native-home'));display=x.XOpenDisplay(env['DISPLAY'].encode());assert display
 log=open(base/'native.log','w');out=open(base/'native-result.json','w');app=subprocess.Popen([str(binary),'sdui','run','-entry','Analysis','-context',str(ctx),str(source),'--json'],env=env,stdout=out,stderr=log,start_new_session=True)
 def windows(w):
  name=C.c_void_p();title=''
  if x.XFetchName(display,w,C.byref(name)) and name.value:title=C.string_at(name).decode(errors='replace');x.XFree(name)
  yield w,title
  root=C.c_ulong();parent=C.c_ulong();children=C.POINTER(C.c_ulong)();n=C.c_uint()
  if x.XQueryTree(display,w,C.byref(root),C.byref(parent),C.byref(children),C.byref(n)):
   ids=[children[i] for i in range(n.value)]
   if children:x.XFree(children)
   for child in ids:yield from windows(child)
 win=None
 for _ in range(200):
  assert app.poll() is None,'launcher exited: '+(base/'native.log').read_text()
  win=next((w for w,title in windows(x.XDefaultRootWindow(display)) if 'native document prototype' in title),None)
  if win:break
  time.sleep(.1)
 assert win,'native window missing'
 def wait_log(marker, count=1):
  deadline=time.monotonic()+45
  while time.monotonic()<deadline:
   assert app.poll() is None,(base/'native.log').read_text()
   if (base/'native.log').read_text().count(marker)>=count:return
   time.sleep(.1)
  raise AssertionError('missing marker '+marker+' '+(base/'native.log').read_text())
 def shot(name):
  deadline=time.monotonic()+30
  while time.monotonic()<deadline:
   time.sleep(1);assert app.poll() is None
   dest=str(base/(name+'.png'));subprocess.run(['import','-window','root',dest],env=env,check=True,timeout=30)
   mean=float(subprocess.check_output(['identify','-format','%[fx:mean]',dest],text=True))
   if mean>.2:return
  raise AssertionError('native framebuffer stayed blank')
 wait_log('Document prototype ready:');shot('native-initial')
 svg=fixture/'assets/frame.svg';origsvg=svg.read_bytes();svg.write_bytes(origsvg+b'\n<!-- reload resource -->\n')
 wait_log('Document prototype reloaded:',1);shot('native-resource-reload');svg.write_bytes(origsvg)
 wait_log('Document prototype reloaded:',2)
 source.write_text('not valid SDUI');wait_log('syntax at');shot('native-rejected-reload');source.write_bytes(original)
 wait_log('Document prototype reloaded:',3);shot('native-recovered')
 assert 'pane-minimum' not in (base/'native.log').read_text()

 ev=Event();ev.client.type=33;ev.client.display=display;ev.client.window=win;ev.client.message_type=x.XInternAtom(display,b'WM_PROTOCOLS',0);ev.client.format=32;ev.client.data.l[0]=x.XInternAtom(display,b'WM_DELETE_WINDOW',0);ev.client.data.l[1]=0;x.XSendEvent(display,win,0,0,C.byref(ev));x.XFlush(display)
 code=app.wait(timeout=20);log.close();out.close();assert code==0,(base/'native.log').read_text()
 (base/'native-smoke-result.json').write_text(json.dumps(dict(exitCode=code,actualPackagedLauncher=True,windowObserved=True,resourceReloadProcessSurvived=True,invalidSourcePreservedProcess=True,recoveryProcessSurvived=True,closedByWindowProtocol=True,reloadAcknowledged=True,nonblankFramebuffer=True,noStartupMinimumDiagnostic=True,scope='Actual native markers and nonblank captures; programmatic host/watcher tests establish state assertions'),indent=2)+'\n');print('packaged native smoke passed')
finally:
 source.write_bytes(original);ctx.write_bytes(originalctx)
 if 'origsvg' in locals():svg.write_bytes(origsvg)
 if app is not None and app.poll() is None:os.killpg(app.pid,signal.SIGTERM);app.wait(timeout=10)
 if display:x.XCloseDisplay(display)
 xf.terminate();xf.wait(timeout=10)
