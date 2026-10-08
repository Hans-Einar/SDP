#!/usr/bin/env python3
"""Real configured XIM discrimination on the staged SDUI-root legacy helper."""
import argparse,hashlib,json,os,subprocess,sys,time
from pathlib import Path

TITLE='SDUI · page · native prototype'
SOURCE='sdui 0.2;\npage=[field=input("Text",value="A") {x=fill}] {x=fill,y=fill};\n'

def main():
 p=argparse.ArgumentParser();p.add_argument('--binary',required=True);p.add_argument('--display',required=True);p.add_argument('--out',required=True);args=p.parse_args()
 out=Path(args.out).resolve();out.mkdir(parents=True,exist_ok=False)
 for name in ['config','cache']:(out/name).mkdir()
 src=out/'input.sdui';src.write_text(SOURCE);digest=hashlib.sha256(src.read_bytes()).hexdigest()
 journal=(out/'actions.ndjson').open('w');log=(out/'stdout.txt').open('w');errors=(out/'stderr.txt').open('w');owner=None;process=None
 def record(kind,value):journal.write(json.dumps({'kind':kind,'value':value},ensure_ascii=False)+'\n');journal.flush()
 def native(action,*values):
  record('x11-input',{'action':action,'values':values})
  return subprocess.run([sys.executable,str(Path(__file__).with_name('x11_input.py')),'--display',args.display,'--title',TITLE,action,*values],capture_output=True,text=True,timeout=10)
 def key(action,*values):
  r=native(action,*values);assert r.returncode==0,r.stderr
 def check(name,passed):
  record('check',{'name':name,'passed':bool(passed)});assert passed,name;print('PASS',name,flush=True)
 def screen(name):
  time.sleep(.2);subprocess.run(['import','-display',args.display,'-window','root',str(out/(name+'.png'))],check=True,timeout=10);record('screenshot',name+'.png')
 def copy(expected):
  nonlocal owner
  if owner is not None and owner.poll() is None:owner.terminate();owner.wait(timeout=3)
  owner=subprocess.Popen(['xclip','-display',args.display,'-selection','clipboard','-target','UTF8_STRING','-in','-quiet'],stdin=subprocess.PIPE,stdout=subprocess.DEVNULL,stderr=subprocess.DEVNULL)
  owner.stdin.write(b'clipboard-observation-sentinel');owner.stdin.close()
  deadline=time.monotonic()+3
  while True:
   r=subprocess.run(['xclip','-display',args.display,'-selection','clipboard','-target','UTF8_STRING','-out'],capture_output=True,timeout=3)
   if r.stdout==b'clipboard-observation-sentinel':break
   assert time.monotonic()<deadline,'Clipboard sentinel did not become current';time.sleep(.05)
  key('chord','Control_L','a');key('chord','Control_L','c')
  deadline=time.monotonic()+5
  while True:
   r=subprocess.run(['xclip','-display',args.display,'-selection','clipboard','-target','UTF8_STRING','-out'],capture_output=True,timeout=3)
   if r.returncode==0 and r.stdout.decode('utf-8')==expected:record('clipboard-observed',expected);return True
   if time.monotonic()>deadline:record('clipboard-mismatch',r.stdout.decode('utf-8',errors='replace'));return False
   time.sleep(.05)
 def dirty():key('key','End');key('chord','Shift_L','b')
 def compose(hexvalue):key('chord','Control_L','Shift_L','u');key('key',*hexvalue)
 try:
  argv=[str(Path(args.binary).resolve()),'-watch=false','-entry','page','-revision',digest,str(src)]
  record('candidate',{'binary':argv[0],'binarySHA256':hashlib.sha256(Path(argv[0]).read_bytes()).hexdigest(),'sourceSHA256':digest,'argv':argv,'display':args.display,'XMODIFIERS':os.environ.get('XMODIFIERS'),'IBUS_ENABLE_CTRL_SHIFT_U':os.environ.get('IBUS_ENABLE_CTRL_SHIFT_U')})
  env=dict(os.environ,DISPLAY=args.display,XDG_CONFIG_HOME=str(out/'config'),XDG_CACHE_HOME=str(out/'cache'),FYNE_THEME='light')
  process=subprocess.Popen(argv,env=env,stdout=log,stderr=errors)
  deadline=time.monotonic()+45
  while True:
   located=native('locate')
   if located.returncode==0 and json.loads(located.stdout)['viewable']:break
   assert process.poll() is None,'Helper exited before its native window';assert time.monotonic()<deadline,'Helper window did not become viewable';time.sleep(.1)
  key('key','Tab');check('actual staged helper starts with accepted A',copy('A'));screen('helper-initial')
  dirty();compose('4e2d');screen('helper-composition');key('key','Return')
  check('real IME Return inserts composition after the dirty suffix',copy('AB中'))
  key('key','Escape');check('composition Return did not prematurely accept the dirty suffix',copy('A'))
  dirty();compose('65b0');key('key','Escape');check('IME-consumed Escape preserves the prior dirty suffix',copy('AB'))
  key('key','Escape');check('ordinary Escape restores the accepted baseline',copy('A'))
  dirty();compose('4e2d');key('key','Return');check('second real composition produces the expected text',copy('AB中'))
  key('key','Return');key('key','x');key('key','Escape');check('ordinary Return subsequently accepts composed text',copy('AB中'));screen('helper-composed-accepted')
  key('close-window');process.wait(timeout=10);check('helper closes cleanly through the native window protocol',process.returncode==0)
  record('result','passed')
 except Exception as e:record('result',{'failed':str(e)});raise
 finally:
  if process is not None and process.poll() is None:process.terminate();process.wait(timeout=5)
  if owner is not None and owner.poll() is None:owner.terminate();owner.wait(timeout=3)
  log.close();errors.close();journal.close()

if __name__=='__main__':main()
