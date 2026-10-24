#!/usr/bin/env python3
"""Native Codex transport observer and optional one-shot terminal-read disconnect.

Never simulates model replies or changes their content. Records allowlisted method,
identity, size and tool-name metadata only. Test paths are supplied by the launcher.
"""
import json
import os
from pathlib import Path
import subprocess
import sys
import threading


def main():
    binary = os.environ.get('PG_LAB_CODEX', '/home/warloc/.local/bin/codex')
    argv = [binary, '--disable', 'apps', '--disable', 'multi_agent', '--disable', 'remote_plugin', *sys.argv[1:]]
    if '--version' in sys.argv:
        os.execv(binary, argv)
    record = Path(os.environ['PG_LAB_PROTOCOL_LOG'])
    marker = Path(os.environ['PG_LAB_FAULT_MARKER']) if os.environ.get('PG_LAB_FAULT_MARKER') else None
    child = subprocess.Popen(argv, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.DEVNULL)
    pending = {}
    held_turn = None
    fault_mode = os.environ.get("PG_LAB_FAULT_MODE", "terminal-read")
    mutex = threading.Lock()

    def log(value):
        with mutex:
            with record.open('a') as output:
                output.write(json.dumps(value) + '\n')

    def send():
        try:
            for raw in sys.stdin.buffer:
                value = json.loads(raw)
                method = value.get('method', '')
                if 'id' in value and method:
                    pending[str(value['id'])] = method
                    row = {'direction': 'request', 'method': method}
                    params = value.get('params', {})
                    for key in ('threadId', 'clientUserMessageId', 'sandbox', 'approvalPolicy'):
                        if key in params: row[key] = params[key]
                    log(row)
                child.stdin.write(raw)
                child.stdin.flush()
        except (BrokenPipeError, ValueError):
            pass
        finally:
            try: child.stdin.close()
            except BrokenPipeError: pass

    threading.Thread(target=send, daemon=True).start()
    try:
        for raw in child.stdout:
            value = json.loads(raw)
            method = pending.pop(str(value.get('id')), '') if 'method' not in value else ''
            result = value.get('result', {})
            if method:
                row = {'direction': 'response', 'method': method, 'bytes': len(raw)}
                if 'error' in value: row['errorCode'] = value['error'].get('code')
                if method in ('thread/start', 'thread/resume', 'thread/read'):
                    thread = result.get('thread', {})
                    row['thread'] = {k: thread[k] for k in ('id', 'cwd', 'cliVersion', 'status') if k in thread}
                    row['turns'] = [{'id': turn.get('id'), 'status': turn.get('status'), 'itemsView': turn.get('itemsView')} for turn in thread.get('turns', [])]
                    for key in ('model', 'reasoningEffort'):
                        if key in result: row[key] = result[key]
                if method == 'turn/start': row['turn'] = {k: result.get('turn', {}).get(k) for k in ('id', 'status')}
                if method == 'mcpServerStatus/list':
                    row['servers'] = [{'name': server['name'], 'tools': sorted(server.get('tools', {}))} for server in result.get('data', [])]
                    row['nextCursor'] = result.get('nextCursor')
                log(row)
                if marker and fault_mode == 'terminal-read' and not marker.exists() and method == 'thread/read' and any(t.get('status') == 'completed' for t in result.get('thread', {}).get('turns', [])):
                    marker.write_text(json.dumps({'fault': 'terminal-read-response-withheld', 'observation': row}, indent=2) + '\n')
                    # Host has already returned completed history. Close transport;
                    # the controller must recover its mapped turn without resending.
                    break
                if marker and fault_mode == 'turn-start' and not marker.exists() and method == 'turn/start' and 'error' not in value:
                    held_turn = result.get('turn', {}).get('id')
                    log({'fault': 'turn-start-response-withheld', 'turnId': held_turn})
                    continue
            if held_turn and value.get('method') == 'turn/completed' and value.get('params', {}).get('turn', {}).get('id') == held_turn:
                marker.write_text(json.dumps({'fault': 'turn-start-response-withheld-until-completion', 'turnId': held_turn, 'status': value['params']['turn']['status']}, indent=2) + '\n')
                break
            sys.stdout.buffer.write(raw)
            sys.stdout.buffer.flush()
    except BrokenPipeError:
        pass
    finally:
        child.terminate()
        try: child.wait(timeout=5)
        except subprocess.TimeoutExpired: child.kill(); child.wait()


if __name__ == '__main__':
    main()
