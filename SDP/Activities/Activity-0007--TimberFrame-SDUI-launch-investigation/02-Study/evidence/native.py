import subprocess,os,pathlib,json
p=pathlib.Path('/tmp/sdui-timberframe-repro');r,w=os.pipe();xf=subprocess.Popen(['/usr/bin/Xvfb','-displayfd',str(w),'-screen','0','1760x1300x24','-nolisten','tcp'],pass_fds=(w,),stderr=open(p/'xvfb.log','w'));os.close(w)
try:
 display=os.read(r,32).decode().strip();os.close(r);env=os.environ.copy();env.update(DISPLAY=':'+display,LIBGL_ALWAYS_SOFTWARE='1',XDG_CACHE_HOME=str(p/'cache'),GOTOOLCHAIN='local',GOCACHE='/home/warloc/.cache/sdp-workstation/gocache',GOMODCACHE='/home/warloc/.cache/sdp-workstation/gomodcache');env['PATH']='/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin:'+env['PATH']
 cmd=['./tools/preview_gui.sh','-entry','Analysis','-smoke','-inspect',str(p/'analysis-inspection.json')]
 with open(p/'native.log','w') as log:
  run=subprocess.run(cmd,cwd='/home/warloc/git/freecad-timberframe',env=env,stdout=log,stderr=log,timeout=240)
 (p/'native-result.json').write_text(json.dumps(dict(command=cmd,exitCode=run.returncode,isolatedDisplay=True),indent=2)+'\n');print('native exit',run.returncode)
finally:xf.terminate();xf.wait(timeout=10)
