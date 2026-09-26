import subprocess,tempfile,pathlib,os,time,json,socket,ctypes,re
import argparse
cli=argparse.ArgumentParser(description='Run only under an isolated Xvfb display. Kills/restarts its own daemon and closes its own FOX window.')
cli.add_argument('--xfmd',required=True);cli.add_argument('--daemon',required=True);cli.add_argument('--request',required=True);cli.add_argument('--renderer',required=True);cli.add_argument('--capture',default='/tmp/sdl-broker-native.png')
config=cli.parse_args();xfmd=config.xfmd
display=os.environ['DISPLAY']
def close_window():
    lib=ctypes.CDLL('libX11.so.6');lib.XOpenDisplay.restype=ctypes.c_void_p;d=lib.XOpenDisplay(display.encode())
    lib.XDefaultRootWindow.argtypes=[ctypes.c_void_p];lib.XDefaultRootWindow.restype=ctypes.c_ulong
    lib.XQueryTree.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.POINTER(ctypes.c_ulong),ctypes.POINTER(ctypes.c_ulong),ctypes.POINTER(ctypes.POINTER(ctypes.c_ulong)),ctypes.POINTER(ctypes.c_uint)]
    root,parent=ctypes.c_ulong(),ctypes.c_ulong();children=ctypes.POINTER(ctypes.c_ulong)();count=ctypes.c_uint()
    lib.XQueryTree(d,lib.XDefaultRootWindow(d),ctypes.byref(root),ctypes.byref(parent),ctypes.byref(children),ctypes.byref(count))
    lib.XFetchName.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.POINTER(ctypes.c_char_p)];lib.XFree.argtypes=[ctypes.c_void_p];xid=None
    for i in range(count.value):
        name=ctypes.c_char_p()
        if lib.XFetchName(d,children[i],ctypes.byref(name)) and name.value:
            if b'xfmd' in name.value:xid=children[i]
            lib.XFree(name)
    lib.XFree(children);assert xid is not None
    lib.XInternAtom.argtypes=[ctypes.c_void_p,ctypes.c_char_p,ctypes.c_int];lib.XInternAtom.restype=ctypes.c_ulong
    class Data(ctypes.Union):_fields_=[('l',ctypes.c_long*5),('b',ctypes.c_char*20)]
    class Client(ctypes.Structure):_fields_=[('type',ctypes.c_int),('serial',ctypes.c_ulong),('send_event',ctypes.c_int),('display',ctypes.c_void_p),('window',ctypes.c_ulong),('message_type',ctypes.c_ulong),('format',ctypes.c_int),('data',Data)]
    class Event(ctypes.Union):_fields_=[('xclient',Client),('pad',ctypes.c_long*24)]
    event=Event();event.xclient.type=33;event.xclient.display=d;event.xclient.window=xid;event.xclient.message_type=lib.XInternAtom(d,b'WM_PROTOCOLS',0);event.xclient.format=32;event.xclient.data.l[0]=lib.XInternAtom(d,b'WM_DELETE_WINDOW',0)
    lib.XSendEvent.argtypes=[ctypes.c_void_p,ctypes.c_ulong,ctypes.c_int,ctypes.c_long,ctypes.POINTER(Event)];lib.XSendEvent(d,xid,0,0,ctypes.byref(event));lib.XFlush.argtypes=[ctypes.c_void_p];lib.XFlush(d);lib.XCloseDisplay.argtypes=[ctypes.c_void_p];lib.XCloseDisplay(d)
with tempfile.TemporaryDirectory(prefix='sdl-lease-native-') as temp:
    root=pathlib.Path(temp);runtime=root/'runtime';runtime.mkdir(mode=0o700);source=root/'model.design';source.write_text('language design-core version 0.5.\nunit Host.\n');nav=root/'navigator.md';nav.write_text('# Navigator\n\n[Architecture](sdl-view://demo/VP02?diagram=VP02-roots)\n')
    env=dict(os.environ,DISPLAY=display,HOME=temp,XDG_RUNTIME_DIR=str(runtime));daemon=None;window=None
    def start():
        p=subprocess.Popen([config.daemon,'-source',str(source),'-project','demo','-runtime',str(runtime/'sdl'),'-xfmd',xfmd,'-renderer',config.renderer],env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
        for _ in range(100):
            if p.poll() is not None:raise RuntimeError(p.communicate())
            try:
                with socket.socket(socket.AF_UNIX,socket.SOCK_SEQPACKET) as s:s.connect(str(runtime/'sdl/views.sock'));return p
            except OSError:time.sleep(.03)
        raise RuntimeError('daemon timeout')
    def call(*args,ok=True):
        p=subprocess.run([config.request,'-socket',str(runtime/'sdl/views.sock'),*args],env=env,capture_output=True,text=True,timeout=15)
        if ok:assert p.returncode==0,p.stderr;return json.loads(p.stdout)
        assert p.returncode!=0;return p.stderr
    def select(seq):return call('-uri','sdl-view://demo/VP02?diagram=VP02-roots','-window','lease-window','-client','probe','-sequence',str(seq),'-open')['result']
    try:
        daemon=start();window=subprocess.Popen([xfmd,'--navigator',str(nav),'--window-id','lease-window','--sdl-tool',config.request,'--project','demo','--broker',str(runtime/'sdl/views.sock')],env=env,stdout=subprocess.PIPE,stderr=subprocess.PIPE,text=True)
        for _ in range(100):
            if (runtime/'xfmd/lease-window.sock').exists():break
            if window.poll() is not None:raise RuntimeError(window.communicate())
            time.sleep(.03)
        first=select(1);entry=pathlib.Path(first['Entry']);assert entry.exists();assert list((entry.parent/'diagrams').glob('*.svg'))
        daemon.kill();daemon.communicate(timeout=5);daemon=start();call('-sweep');assert entry.exists(),'lease lost after crash'
        call('-uri','sdl-view://demo/VP02','-window','lease-window','-client','probe','-sequence','1','-open',ok=False)
        source.write_text('invalid');call('-uri','sdl-view://demo/VP02','-window','lease-window','-client','probe','-sequence','2','-open',ok=False);assert entry.exists()
        source.write_text('language design-core version 0.5.\nunit Reloaded.\n');second=select(3);time.sleep(.3);call('-sweep');assert not entry.exists(),'old lease not released on replace';assert pathlib.Path(second['Entry']).exists()
        time.sleep(.4);subprocess.run(['import','-window','root',config.capture],env=env,check=True)
        close_window();window.communicate(timeout=5);time.sleep(.2);call('-sweep');assert not pathlib.Path(second['Entry']).exists(),'active lease not released on close'
        print('PASS: real SDL daemon + FOX, SVG lease, SIGKILL/restart, invalid source, replacement release, native window-close release')
    finally:
        for p in [window,daemon]:
            if p and p.poll() is None:p.terminate();p.communicate(timeout=5)
