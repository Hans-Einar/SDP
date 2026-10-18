#!/usr/bin/env python3
"""Run a native acceptance command with an isolated real IBus/XIM session.

Requires an explicitly supplied extracted IBus installation and components. Never
installs packages, modifies desktop configuration or selects the owner's display.
"""
import argparse
import os
from pathlib import Path
import subprocess
import sys
import time


def main():
    p = argparse.ArgumentParser()
    p.add_argument('--display', required=True)
    p.add_argument('--installation', required=True, type=Path)
    p.add_argument('--components', required=True, type=Path)
    p.add_argument('--out', required=True, type=Path)
    p.add_argument('--inside-bus', action='store_true', help=argparse.SUPPRESS)
    p.add_argument('command', nargs=argparse.REMAINDER)
    a = p.parse_args()
    command = a.command[1:] if a.command[:1] == ['--'] else a.command
    if not command:
        p.error('A native verification command is required after --')
    root, out = a.installation.resolve(), a.out.resolve()
    out.mkdir(parents=True, exist_ok=True)
    for n in ('config', 'cache', 'runtime', 'data'):
        (out/n).mkdir(mode=0o700, exist_ok=True)
    if not a.inside_bus:
        bus_env=dict(os.environ, XDG_DATA_HOME=str(out/'data'),
                     XDG_CONFIG_HOME=str(out/'config'), XDG_CACHE_HOME=str(out/'cache'),
                     XDG_RUNTIME_DIR=str(out/'runtime'))
        return subprocess.call(['dbus-run-session', '--', sys.executable,
                                str(Path(__file__).resolve()), '--inside-bus',
                                *sys.argv[1:]], env=bus_env)
    (out/'runtime').chmod(0o700)
    env = dict(os.environ, DISPLAY=a.display, IBUS_ENABLE_CTRL_SHIFT_U='1',
               LC_ALL='en_US.UTF-8', XMODIFIERS='@im=ibus',
               IBUS_COMPONENT_PATH=str(a.components.resolve()),
               LD_LIBRARY_PATH=str(root/'usr/lib64'),
               GSETTINGS_SCHEMA_DIR=str(root/'usr/share/glib-2.0/schemas'),
               XDG_CONFIG_HOME=str(out/'config'), XDG_CACHE_HOME=str(out/'cache'),
               XDG_RUNTIME_DIR=str(out/'runtime'), GSETTINGS_BACKEND='memory')
    env.pop('IBUS_ADDRESS', None)
    children, logs = [], []
    def start(argv, name):
        stream=(out/(name+'.log')).open('w')
        logs.append(stream)
        child=subprocess.Popen(argv, env=env, stdout=stream, stderr=subprocess.STDOUT)
        children.append(child)
        return child
    try:
        daemon=start([str(root/'usr/bin/ibus-daemon'), '--single', '--panel', 'disable',
                      '--emoji-extension', 'disable', '--config', 'disable',
                      '--cache', 'none'], 'ibus')
        deadline=time.monotonic()+15
        while True:
            result=subprocess.run([str(root/'usr/bin/ibus'), 'address'], env=env,
                                  capture_output=True, text=True, timeout=3)
            address=result.stdout.strip()
            if result.returncode==0 and address.startswith('unix:'):
                break
            if daemon.poll() is not None or time.monotonic()>deadline:
                raise RuntimeError('Isolated IBus did not publish a usable address')
            time.sleep(.1)
        xim=start([str(root/'usr/libexec/ibus-x11'), '--debug=1', '--locale=en_US.UTF-8'], 'xim')
        # XIM server registration happens asynchronously; this is environment
        # setup, not a substitute for application composition assertions.
        time.sleep(.5)
        result=subprocess.run([str(root/'usr/bin/ibus'), 'engine', 'xkb:us::eng'],
                              env=env, capture_output=True, text=True, timeout=5)
        (out/'engine.log').write_text(result.stdout+result.stderr)
        if result.returncode or xim.poll() is not None:
            raise RuntimeError('Isolated XIM engine setup failed')
        with (out/'invocation.txt').open('w') as log:
            log.write('Display: '+a.display+'\nEngine: xkb:us::eng\n')
            log.write('Command argv: '+repr(command)+'\n')
        return subprocess.call(command, env=env)
    finally:
        for child in reversed(children):
            if child.poll() is None:
                child.terminate()
        for child in children:
            try:
                child.wait(timeout=3)
            except subprocess.TimeoutExpired:
                child.kill(); child.wait(timeout=3)
        for stream in logs:
            stream.close()


if __name__=='__main__':
    raise SystemExit(main())
