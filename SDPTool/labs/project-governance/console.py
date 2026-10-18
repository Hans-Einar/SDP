#!/usr/bin/env python3
"""A trial-only readable adapter over existing PG commands; no workflow engine."""
import argparse
import hashlib
import json
import os
from pathlib import Path
import queue
import shutil
import signal
import subprocess
import sys
import threading

BASE = Path(__file__).resolve().parent
PROJECT = BASE / 'project'
RUNTIME = BASE / 'runtime'
BIN = BASE / 'candidate/bin'


def show(text):
    # Messages are data, never terminal controls or executable commands.
    print(''.join(c if c.isprintable() or c in '\n\t' else repr(c)[1:-1] for c in str(text)), flush=True)


def load(name):
    return json.loads((BASE / name).read_text())


def identity():
    data = load('candidate.json')
    for name, digest in data['binaries'].items():
        if hashlib.sha256((BIN / name).read_bytes()).hexdigest() != digest:
            raise RuntimeError('Candidate binary changed: ' + name)
    return data


def env():
    result = dict(os.environ)
    if not shutil.which('go') and (Path.home() / '.local/bin/go').is_file():
        result['PATH'] = str(Path.home() / '.local/bin') + os.pathsep + result.get('PATH', '')
    return result


def args(operation):
    return [str(BIN / 'sdp-controller'), operation, '--project', str(PROJECT),
            '--state', str(RUNTIME / 'state'), '--run', load('run.json')['run']['id']]


def read_controller(operation):
    p = subprocess.run(args(operation), capture_output=True, text=True, env=env())
    if p.returncode:
        raise RuntimeError(p.stderr or p.stdout)
    return [json.loads(line) for line in p.stdout.splitlines() if line]


def summary(value):
    run = value['run']
    states = {'ready': 'Ready to launch', 'running': 'Work in progress or awaiting reconciliation',
              'verifying': 'Agent report received; checking/acceptance is still open',
              'interrupted': 'Routine interrupted; inspect host outcome',
              'replan-required': 'Scope changed; a revised assignment is needed'}
    show(f"Run: {run['id']}\nState: {run['state']} — {states.get(run['state'], 'Inspect the recorded outcome')}\nRevision: {run['revision']}")
    for attempt in run.get('attempts', []):
        show(f"Attempt: {attempt['id']} | host: {attempt['status']} | gap: {attempt['gap']}\nThread: {attempt['threadId']} | turn: {attempt['turnId']}")
    for entry in run.get('inputs', []):
        show(f"Input: {entry['kind']} / {entry['state']}")


def overview():
    candidate = identity()
    show('ProjectGovernance owner lab — an unreleased workflow trial')
    show(f"Candidate: {candidate['sourceCommit']}\nProject: {PROJECT}")
    show('Goal: make the label command trim boundary Unicode whitespace without changing interior text or case.')
    result = read_controller('status')[0]
    run = result['run']
    selection = run['selection']
    show('\nSelected work (actual stored references):')
    for key in ['session', 'step', 'card', 'plan', 'authority']:
        show(f"  {key}: {selection[key]}")
    show('Writable: ' + ', '.join(run['assignment']['writablePaths']))
    show('Excluded: ' + ', '.join(run['assignment']['excludedPaths']))
    show('\nRoadmap: inspect assignment → start → observe/steer → report → inspect diff/tests → record feedback')
    summary(result)
    show('\nRoles: you select/steer; one fresh Master performs the bounded task. No child agents or independent review in this candidate.')
    show('Automatic card/Session publication and work acceptance are not implemented. The report stops at verifying.')
    show('\nCommands: ./lab overview | run | status | provenance | result | recover | feedback')
    show('Read README.md for the trial and controls. Edit owner.txt before first launch to refine wording within the stored scope.')


def render(value):
    if 'run' in value:
        summary(value)
    elif 'approval' in value:
        prompt = value['approval']
        show('\nDECISION REQUIRED — exact pending request:')
        show(json.dumps(prompt['request'], indent=2))
        show('Use /accept, /decline or /cancel. No default response is sent.')
        return prompt
    elif 'observation' in value:
        o = value['observation']
        if o['kind'] == 'agentMessage':
            show(f"\nAgent [{o['phase']}]:\n{o['text']}")
        elif o['kind'] == 'userMessage':
            show('Host observed the submitted input (raw content retained in runtime/events.ndjson).')
        elif o['kind'] in ['turn.started', 'turn.completed', 'tool.completed', 'approval.resolved']:
            show(f"[{o['kind']}] {o.get('status', '')}")
    elif 'error' in value:
        show('Controller: ' + json.dumps(value, ensure_ascii=False))
    return None


def command(line, pending):
    line = line.strip()
    if line == '/status':
        return {'action': 'status', 'text': ''}
    for prefix, action in [('/steer ', 'steer'), ('/interrupt ', 'interrupt'), ('/scope ', 'scope-change')]:
        if line.startswith(prefix) and line[len(prefix):].strip():
            return {'action': action, 'text': line[len(prefix):]}
    if line in ['/accept', '/decline', '/cancel']:
        if not pending:
            raise ValueError('No unanswered request is currently presented.')
        reply = {key: pending[key] for key in ['requestId', 'epoch', 'threadId', 'turnId']}
        reply['decision'] = line[1:]
        return {'action': 'approval', 'approval': reply}
    raise ValueError('Use /status, /steer TEXT, /interrupt REASON, /scope TEXT, /accept, /decline, /cancel or /quit.')


def run():
    identity()
    if not sys.stdin.isatty():
        raise RuntimeError('Run this owner trial in an interactive terminal; no unattended model launch.')
    overview()
    if input('\nStart this bounded Codex task now? Type start: ').strip() != 'start':
        show('No task started.')
        return
    cmd = args('run') + ['--mcp-binary', str(BIN / 'sdp-mcp'), '--message', str(BASE / 'owner.txt'),
                         '--sandbox', 'workspace-write', '--approvals', 'ask',
                         '--skill', 'Skills/sdp/SKILL.md', '--skill', 'Skills/sdp-master/SKILL.md']
    RUNTIME.mkdir(exist_ok=True, mode=0o700)
    events = queue.Queue()
    with (RUNTIME / 'controller.err').open('a') as err, (RUNTIME / 'events.ndjson').open('a') as raw, (RUNTIME / 'commands.ndjson').open('a') as sent:
        p = subprocess.Popen(cmd, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=err, text=True, bufsize=1, env=env())
        def reader(stream, kind):
            for line in stream:
                events.put((kind, line))
            events.put((kind + '-eof', ''))
        threading.Thread(target=reader, args=(p.stdout, 'host'), daemon=True).start()
        threading.Thread(target=reader, args=(sys.stdin, 'owner'), daemon=True).start()
        pending = None
        show('Live controls: /status | /steer TEXT | /interrupt REASON | /scope TEXT | /quit')
        try:
            while True:
                kind, line = events.get()
                if kind == 'host-eof':
                    break
                if kind == 'owner-eof':
                    show('Input closed: observation only; no decision or cancellation inferred.')
                elif kind == 'host':
                    raw.write(line); raw.flush()
                    prompt = render(json.loads(line))
                    if prompt:
                        pending = prompt
                elif kind == 'owner':
                    if line.strip() == '/quit':
                        p.send_signal(signal.SIGINT)
                        show('Controller stopping. Inspect status/recover; this is not a completion claim.')
                        continue
                    try:
                        value = command(line, pending)
                        encoded = json.dumps(value) + '\n'
                        sent.write(encoded); sent.flush()
                        p.stdin.write(encoded); p.stdin.flush()
                        if value['action'] == 'approval':
                            pending = None
                    except (ValueError, BrokenPipeError) as e:
                        show(e)
        except KeyboardInterrupt:
            p.send_signal(signal.SIGINT)
            show('Interrupted controller; use status/recover to inspect outcome.')
        finally:
            if p.poll() is None:
                try:
                    p.wait(timeout=10)
                except subprocess.TimeoutExpired:
                    p.terminate(); p.wait(timeout=10)
            p.stdin.close()
        show(f'Controller exit: {p.returncode}. Diagnostics: {RUNTIME / "controller.err"}')
    show('Next: ./lab result, then write observations in FEEDBACK.md. Host completion is not work acceptance.')


def result():
    identity()
    summary(read_controller('status')[0])
    show('\nActual working diff:')
    subprocess.run(['git', 'diff', '--', '.'], cwd=PROJECT, check=True)
    show('\nProject tests (local checking, not independent review):')
    p = subprocess.run(['go', 'test', './...'], cwd=PROJECT, env=env(), capture_output=True, text=True)
    show(p.stdout + p.stderr)
    (RUNTIME / 'local-check.txt').write_text(p.stdout + p.stderr + f'Exit: {p.returncode}\n')
    pre = load('preimages.json')
    changed = [path for path, digest in pre.items() if not (PROJECT / path).is_file() or hashlib.sha256((PROJECT / path).read_bytes()).hexdigest() != digest]
    unexpected = [path for path in changed if path not in ['normalize.go', 'normalize_test.go']]
    links = load('discovery-links.json')
    altered = [path for path, target in links.items() if not (PROJECT / path).is_symlink() or str((PROJECT / path).readlink()) != target]
    show('Changed baseline files: ' + ', '.join(changed))
    show('Unexpected baseline/link changes: ' + json.dumps(unexpected + altered))
    show('Also inspect newly created files:')
    subprocess.run(['git', 'status', '--short'], cwd=PROJECT, check=True)
    show('Record your experience in FEEDBACK.md; tests do not decide whether this is the workflow you want.')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('command', nargs='?', default='overview', choices=['overview', 'run', 'status', 'provenance', 'result', 'recover', 'feedback'])
    selection = parser.parse_args().command
    if selection == 'overview':
        overview()
    elif selection == 'run':
        run()
    elif selection == 'result':
        result()
    elif selection == 'feedback':
        show((BASE / 'FEEDBACK.md').read_text())
        show('Edit: ' + str(BASE / 'FEEDBACK.md'))
    elif selection == 'provenance':
        identity()
        value = read_controller('provenance')[0]
        context = value.get('context', {})
        show('Context snapshot: ' + value['status'])
        show('Assigned role: ' + context.get('assignedRole', 'No launched attempt'))
        show('Host native role: ' + (context.get('host', {}).get('nativeRole') or 'Not reported'))
        inventory = context.get('skills', {})
        show('Project skill availability: ' + ', '.join(s['name'] for s in inventory.get('project', []) if s['enabled']))
        invocation = value.get('invocation')
        if invocation:
            show('Explicit selection: ' + ', '.join(s['name'] for s in invocation['skills']))
            show('Host acceptance: ' + invocation['acceptance'])
            show('Procedural compliance: ' + invocation['compliance'])
        else:
            show('No explicit skill invocation recorded yet.')
        show('Raw snapshot is available through candidate/bin/sdp-controller provenance; availability does not establish use.')
    else:
        identity()
        for value in read_controller(selection):
            if 'run' in value:
                summary(value)
            else:
                show(json.dumps(value, indent=2, ensure_ascii=False))


if __name__ == '__main__':
    try:
        main()
    except (RuntimeError, subprocess.CalledProcessError, OSError) as e:
        show('Lab error: ' + str(e))
        sys.exit(2)
