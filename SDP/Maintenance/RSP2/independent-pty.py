import os, pty, fcntl, termios, signal, subprocess, json, time, tempfile, pathlib, select
binary='/tmp/sdp-rereview-sdptool'
for mode in ['normal', 'sigterm', 'start-failure']:
    with tempfile.TemporaryDirectory(prefix='sdp-review-') as directory:
        root=pathlib.Path(directory)
        (root/'SDP/UI').mkdir(parents=True)
        (root/'SDP/UI/main.sdui').write_text('sdui 0.2; Page=[<"Hello", button("OK")>];')
        (root/'SDP/programs.json').write_text(json.dumps({'schemaVersion':'sdp-programs/1','programs':[{'id':'test','label':'Test','source':'SDP/UI/main.sdui','entry':'Page','command':['./runner']}]}))
        scripts={'normal':'#!/bin/sh\necho $$ > runner-pid\nread value\necho "$value" > received\nexit 7\n','sigterm':'#!/bin/sh\necho $$ > runner-pid\n(sleep 1; echo leaked > descendant) &\necho ready > ready\nwait\n','start-failure':'#!/no-such-interpreter\n'}
        (root/'runner').write_text(scripts[mode]); (root/'runner').chmod(0o700)
        master,slave=pty.openpty(); readfd,writefd=os.pipe()
        pid=os.fork()
        if pid==0:
            os.close(readfd); os.close(master)
            try:
                os.setsid(); fcntl.ioctl(slave,termios.TIOCSCTTY,0)
                for fd in [0,1,2]: os.dup2(slave,fd)
                before=os.tcgetpgrp(0)
                child=subprocess.Popen([binary,str(root),'run','--program','test'])
                active=None
                if mode != 'start-failure':
                    deadline=time.monotonic()+5
                    marker=root/('ready' if mode=='sigterm' else 'runner-pid')
                    while not marker.exists():
                        if time.monotonic()>deadline: raise RuntimeError('runner not ready')
                        time.sleep(.01)
                    active=os.tcgetpgrp(0)
                    if mode=='sigterm': child.send_signal(signal.SIGTERM)
                code=child.wait(timeout=8)
                after=os.tcgetpgrp(0)
                result={'mode':mode,'before':before,'during':active,'after':after,'exit':code,'restored':before==after}
                if mode=='normal': result['input']=(root/'received').read_text().strip()
                if mode=='sigterm':
                    time.sleep(1.1); result['descendant_leaked']=(root/'descendant').exists()
                os.write(writefd,json.dumps(result).encode())
            except BaseException as e:
                os.write(writefd,repr(e).encode())
            finally: os._exit(0)
        os.close(writefd); os.close(slave)
        if mode=='normal':
            deadline=time.monotonic()+5
            while not (root/'runner-pid').exists():
                if time.monotonic()>deadline: raise RuntimeError('input readiness timeout')
                time.sleep(.01)
            time.sleep(.1)
            os.write(master,b'terminal value\n')
        ready,_,_=select.select([readfd],[],[],15)
        if not ready:
            os.kill(pid,signal.SIGKILL); raise RuntimeError('supervisor timeout')
        result=json.loads(os.read(readfd,65536)); print(json.dumps(result),flush=True)
        os.waitpid(pid,0); os.close(master); os.close(readfd)
        assert result['restored'],result
        if mode=='normal': assert result['exit']==7 and result['input']=='terminal value'
        if mode=='sigterm': assert result['exit']!=0 and not result['descendant_leaked']
        if mode=='start-failure': assert result['exit']!=0
