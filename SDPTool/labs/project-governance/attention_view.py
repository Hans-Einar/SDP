#!/usr/bin/env python3
"""Inspect an isolated attention trial; optionally capture one attributed reply."""
import argparse
import json
from pathlib import Path
import subprocess
import tempfile
import uuid


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--preview', type=Path, required=True)
    p.add_argument('--reply', action='store_true')
    p.add_argument('--owner')
    p.add_argument('--question')
    p.add_argument('--json', action='store_true')
    a = p.parse_args()
    if (a.reply and not a.owner) or (not a.reply and (a.owner or a.question)):
        p.error('--reply requires --owner; --question selects among open questions')
    base = a.preview.resolve()
    common = ['--project',str(base/'project'),'--state',str(base/'runtime/state')]
    binary = str(base/'sdp-controller')
    print('Separate synthetic attention preview. This helper starts no agent itself; an authorized watcher may continue after a reply. No design approval is recorded.', file=__import__('sys').stderr)
    if not a.reply:
        subprocess.run([binary,'attention','list',*common,*(['--json'] if a.json else [])],check=True)
        return
    view = json.loads(subprocess.check_output([binary,'attention','list',*common,'--json']))
    items = [v for v in view['items'] if v['status']=='open' and v['applicability']=='current' and (not a.question or v['question']['id']==a.question)]
    if len(items)!=1:
        p.error('select exactly one current open question with --question ID; list questions first')
    item = items[0]
    request = {'schema':view['schema'],'projectRef':view['projectRef'],'checkoutId':view['checkoutId'],'operationId':str(uuid.uuid4()),'runId':item['runId'],'questionId':item['question']['id'],'basis':item['question']['basis']}
    with tempfile.NamedTemporaryFile(mode='w',suffix='.json',prefix='pg-attention-') as f:
        json.dump(request,f); f.flush()
        subprocess.run([binary,'attention','reply',*common,'--input',f.name,'--owner',a.owner],check=True)

if __name__ == '__main__':
    main()
