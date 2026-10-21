#!/usr/bin/env python3
"""Read a completed native attention trial without starting or changing work."""
import argparse
import json
from pathlib import Path
import subprocess


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--preview', type=Path, required=True)
    parser.add_argument('--json', action='store_true')
    args = parser.parse_args()
    base = args.preview.resolve()
    evidence = json.loads((base/'evidence.json').read_text())
    common = ['--project', str(base/'project'), '--state', str(base/'runtime/state')]
    result = subprocess.check_output([str(base/'sdp-controller'), 'status', *common, '--run', evidence['runId']], text=True)
    run = json.loads(result.splitlines()[0])['run']
    inbox = json.loads(subprocess.check_output([str(base/'sdp-controller'), 'attention', 'list', *common, '--json'], text=True))
    if args.json:
        print(json.dumps({'trial': evidence['status'], 'run': run, 'attention': inbox, 'changedFiles': evidence.get('changedFiles', [])}, indent=2))
        return
    quote = lambda value: json.dumps(value, ensure_ascii=True)
    print('PG continuation: isolated native trial; scripted answer, not owner acceptance')
    print('Trial:', evidence['status'], '| Current run:', run['state'])
    for number, attempt in enumerate(run.get('attempts', []), 1):
        print(f'Attempt {number}: {attempt["status"]}; gap={attempt.get("gap", False)}')
    for item in inbox['items']:
        question = item['question']
        print('Question:', quote(question['prompt']))
        reply = question.get('reply')
        if reply:
            print('Reply:', quote(reply.get('optionId', '')), quote(reply['text']))
        print('Status:', item['status'], '| Delivery:', item['consumption'])
    print('Changed files:', ', '.join(evidence.get('changedFiles', [])))
    print('Behavior check:', evidence.get('goldenAfter', '').strip())
    print('Report is awaiting verification. Dispatch acknowledgement is not acceptance.')
    print('This view is read-only. Source and evidence:')
    print((base/'project'/'label.go').as_uri())
    print((base/'README.md').as_uri())


if __name__ == '__main__':
    main()
