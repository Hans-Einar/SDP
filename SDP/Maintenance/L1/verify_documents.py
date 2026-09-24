#!/usr/bin/env python3
"""Check current documentation links/bundles and L1's preserved historical records.

This supplements R1's dated, byte-identical Go migration check after intentional
changes to generator code. It does not claim semantic language detection or GUI QA.
"""
import hashlib
import json
import re
import subprocess
import xml.etree.ElementTree as ET
from pathlib import Path
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[3]
BASE = '9e0d555'


def sha(data):
    return hashlib.sha256(data).hexdigest()


def anchors(text):
    result, counts = set(), {}
    for line in text.splitlines():
        match = re.match(r'^#{1,6} (.*?)(?: +#+)?$', line)
        if match:
            name = re.sub(r'[^\w\- ]', '', match[1].lower()).replace(' ', '-')
            n = counts.get(name, 0)
            counts[name] = n + 1
            result.add(name + ('-' + str(n) if n else ''))
    result.update(re.findall(r'<a (?:id|name)="([^"]+)"', text))
    return result


def main():
    names = subprocess.check_output(['git', 'ls-tree', '-r', '--name-only', BASE], cwd=ROOT, text=True).splitlines()
    preserved = 0
    for name in names:
        frozen = '/testdata/' in name or '/fixtures/' in name
        machine_evidence = '/evidence/' in name and not name.endswith('.md')
        machine_evidence |= name == 'SDP/History/checkpoint-1/source-index.json'
        archive = name.startswith('SDP/History/legacy-bootstrap/') and not name.endswith('/README.md')
        ledger = name.endswith('/Ledger.ndjson')
        if frozen or machine_evidence or archive or ledger:
            before = subprocess.check_output(['git', 'show', BASE + ':' + name], cwd=ROOT)
            current = (ROOT / name).read_bytes()
            assert current.startswith(before) if ledger else current == before, name
            preserved += 1
    generated = 0
    revision = sha((ROOT / 'SDUI/design/architecture.design').read_bytes())
    for folder in ['viewpoints', 'navigation', 'runtime-preview']:
        directory = ROOT / 'SDUI/design' / folder
        manifest = json.loads((directory / 'manifest.json').read_text())
        for name, digest in manifest['outputs'].items():
            path = directory / name
            assert sha(path.read_bytes()) == digest, str(path)
            if path.suffix == '.svg':
                ET.parse(path)
            generated += 1
        if folder != 'runtime-preview':
            assert manifest['revision'] == revision
        else:
            assert json.loads((directory / 'provenance.json').read_text())['design_sha256'] == revision
    links, fragments = 0, 0
    for path in ROOT.rglob('*.md'):
        if any(part in path.parts for part in ['.git', 'testdata', 'fixtures', 'legacy-bootstrap']):
            continue
        text = re.sub(r'```.*?```', '', path.read_text(), flags=re.S)
        for raw in re.findall(r'\]\(([^)\s]+)\)', text):
            url = urlsplit(raw)
            if url.scheme or url.netloc or '<' in raw:
                continue
            target = (path.parent / unquote(url.path)).resolve() if url.path else path
            assert target.exists(), f'{path}: {raw}'
            links += 1
            if url.fragment and target.suffix == '.md':
                assert unquote(url.fragment) in anchors(target.read_text()), f'{path}: {raw}'
                fragments += 1
    print(f'PASS: {preserved} frozen records/ledger prefixes; {generated} generated outputs; {links} file links, {fragments} fragments')
    print('External URLs, localized UI content and visual rendering are outside this check.')


if __name__ == '__main__':
    main()
