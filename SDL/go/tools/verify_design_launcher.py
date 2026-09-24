"""Native launcher acceptance; use an isolated display, never the user desktop."""
import argparse,ctypes,os,pathlib,subprocess,tempfile,time
ROOT=pathlib.Path(__file__).resolve().parents[3]
cli=argparse.ArgumentParser()
cli.add_argument('--xfmd',required=True);cli.add_argument('--sdl',required=True);cli.add_argument('--renderer',required=True)
cli.add_argument('--launcher',default=str(ROOT/'SDL/scripts/sdl-design'))
args=cli.parse_args();XFMD=str(pathlib.Path(args.xfmd).resolve());SDL=str(pathlib.Path(args.sdl).resolve());RENDERER=str(pathlib.Path(args.renderer).resolve())

with tempfile.TemporaryDirectory(prefix='sdl-launch-test-') as temp:
 d=pathlib.Path(temp);(d/'runtime').mkdir(mode=0o700);(d/'.cache').mkdir();(d/'.config').mkdir();(d/'bin').mkdir()
 (d/'bin/go').write_text('#!/bin/sh\necho unexpected-build > "$HOME/go-called"\nexit 99\n');(d/'bin/go').chmod(0o700)
 x=subprocess.Popen(['Xvfb',':93','-screen','0','1600x1000x24','-nolisten','tcp'],stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
 launch=None
 try:
  time.sleep(.5);assert x.poll() is None
  env=dict(os.environ,DISPLAY=':93',HOME=temp,XDG_RUNTIME_DIR=str(d/'runtime'),XDG_CACHE_HOME=str(d/'.cache'),XDG_CONFIG_HOME=str(d/'.config'),SDP_SDL_TOOL=SDL,SDP_XFMD=XFMD,SDP_MMDR=RENDERER,PATH=str(d/'bin')+os.pathsep+os.environ['PATH'])
  log=open(d/'log','w');launch=subprocess.Popen([args.launcher],cwd='/tmp',env=env,stdout=log,stderr=log)
  window=f'sdl-design-{launch.pid}'
  for _ in range(200):
   if (d/'runtime/xfmd'/f'{window}.sock').exists():break
   assert launch.poll() is None,(d/'log').read_text();time.sleep(.05)
  info=subprocess.check_output([XFMD,'--window',window,'--info'],env=env,text=True).strip().split('\t')
  assert info[0]=='OK',info
  nav=pathlib.Path(info[2]);session=nav.parent.parent
  assert nav.name=='navigator.md' and nav.exists()
  assert 'sdl-view://sdui/' in nav.read_text()
  assert not list(nav.parent.rglob('*.svg')) and not list(nav.parent.rglob('*.mmd'))
  assert not (nav.parent/'viewpoints').exists()
  # Exercise the exact registered tool with current source; no pre-generated detail.
  result=subprocess.run([SDL,'view',str(ROOT/'SDUI/design/architecture.design'),'--uri','sdl-view://sdui/VP02?diagram=VP02-roots','--output',str(session/'selected'),'--renderer',RENDERER],env=env,capture_output=True,text=True,timeout=60)
  assert result.returncode==0,result.stderr
  assert list((session/'selected/diagrams').glob('*.svg'))
  # Close only this private X server's FOX top-level window normally.
  lib=ctypes.CDLL('libX11.so.6');lib.XOpenDisplay.restype=ctypes.c_void_p;display=lib.XOpenDisplay(b':93')
  lib.XDefaultRootWindow.argtypes=[ctypes.c_void_p];lib.XDefaultRootWindow.restype=ctypes.c_ulong
  lib.XQueryTree.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.POINTER(ctypes.c_ulong),ctypes.POINTER(ctypes.c_ulong),ctypes.POINTER(ctypes.POINTER(ctypes.c_ulong)),ctypes.POINTER(ctypes.c_uint)]
  root,parent=ctypes.c_ulong(),ctypes.c_ulong();children=ctypes.POINTER(ctypes.c_ulong)();count=ctypes.c_uint();lib.XQueryTree(display,lib.XDefaultRootWindow(display),ctypes.byref(root),ctypes.byref(parent),ctypes.byref(children),ctypes.byref(count))
  lib.XFetchName.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.POINTER(ctypes.c_char_p)];lib.XFree.argtypes=[ctypes.c_void_p];xid=None
  for i in range(count.value):
   name=ctypes.c_char_p()
   if lib.XFetchName(display,children[i],ctypes.byref(name)) and name.value:
    if b'xfmd' in name.value:xid=children[i]
    lib.XFree(name)
  lib.XFree(children);assert xid
  lib.XInternAtom.argtypes=[ctypes.c_void_p,ctypes.c_char_p,ctypes.c_int];lib.XInternAtom.restype=ctypes.c_ulong
  class Data(ctypes.Union):_fields_=[('l',ctypes.c_long*5),('b',ctypes.c_char*20)]
  class Client(ctypes.Structure):_fields_=[('type',ctypes.c_int),('serial',ctypes.c_ulong),('send_event',ctypes.c_int),('display',ctypes.c_void_p),('window',ctypes.c_ulong),('message_type',ctypes.c_ulong),('format',ctypes.c_int),('data',Data)]
  class Event(ctypes.Union):_fields_=[('xclient',Client),('pad',ctypes.c_long*24)]
  e=Event();e.xclient.type=33;e.xclient.display=display;e.xclient.window=xid;e.xclient.message_type=lib.XInternAtom(display,b'WM_PROTOCOLS',0);e.xclient.format=32;e.xclient.data.l[0]=lib.XInternAtom(display,b'WM_DELETE_WINDOW',0)
  lib.XSendEvent.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.c_int,ctypes.c_long,ctypes.POINTER(Event)];lib.XSendEvent(display,xid,0,0,ctypes.byref(e));lib.XFlush.argtypes=[ctypes.c_void_p];lib.XFlush(display);lib.XCloseDisplay.argtypes=[ctypes.c_void_p];lib.XCloseDisplay(display)
  assert launch.wait(timeout=10)==0,(d/'log').read_text();assert not session.exists();assert not (d/'go-called').exists()
  failed=subprocess.run([args.launcher],env=dict(env,SDP_SDL_TOOL=str(d/'absent')),capture_output=True,text=True)
  assert failed.returncode==2 and 'Ferdigbygd SDL-verktøy mangler' in failed.stderr
  assert not list((d/'runtime').glob('sdl-design.*'))
  print('PASS: prebuilt tools only, main page + navigator in native XFMD, on-demand SVG, missing-tool diagnostic, normal close cleanup')
 finally:
  if launch and launch.poll() is None:launch.terminate();launch.wait(timeout=10)
  x.terminate();x.wait(timeout=5)
