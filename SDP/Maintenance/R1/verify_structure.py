#!/usr/bin/env python3
"""Verify the R1 migration boundaries, immutable artifacts and current doc links."""
import hashlib
import json
import re
import subprocess
from pathlib import Path
from urllib.parse import unquote, urlsplit

ROOT = Path(__file__).resolve().parents[3]
AREA = ROOT / 'SDP/Maintenance/R1'


def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()


def main():
    migration = json.loads((AREA / 'Migration-map.json').read_text())
    for old, record in migration['files'].items():
        target = ROOT / record['target']
        assert target.is_file(), record['target']
        assert not (ROOT / old).exists(), old
        if record['target'].startswith('SDP/History/legacy-bootstrap/') or old.endswith(('.go', '/go.mod', '/go.sum')):
            assert digest(target) == record['sha256_before'], old
    # Historical evidence addresses/hashes are not silently reissued.
    for old in ['docs/checkpoint#1/source-index.json',
                'SystemDesignLanguage/go/evidence/release.json',
                'SystemDesignLanguage/tools/verification.json']:
        record = migration['files'][old]
        assert digest(ROOT / record['target']) == record['sha256_before'], old
    old_ledger = subprocess.check_output(
        ['git', 'show', migration['baseline_commit'] + ':Traceability/Ledger.ndjson'], cwd=ROOT)
    assert (ROOT / 'SDP/Traceability/Ledger.ndjson').read_bytes().startswith(old_ledger)
    generated = 0
    for directory in ['viewpoints', 'navigation', 'runtime-preview']:
        base = ROOT / 'SDUI/design' / directory
        manifest = json.loads((base / 'manifest.json').read_text())
        for relative, expected in manifest['outputs'].items():
            assert digest(base / relative) == expected, relative
            generated += 1
    # Public installation destinations/policies remain the same; only source roots move.
    before = json.loads(subprocess.check_output(
        ['git', 'show', migration['baseline_commit'] + ':Toolkit/SDP-install.manifest.json'], cwd=ROOT))
    after = json.loads((ROOT / 'Toolkit/SDP-install.manifest.json').read_text())
    expected = json.loads(json.dumps(before['entries']).replace('Toolkit/project-templates/', 'Template/'))
    assert after['entries'] == expected
    assert before['generators'] == after['generators']
    links = 0
    for path in ROOT.rglob('*.md'):
        if any(part in path.parts for part in ['.git', 'fixtures', 'testdata', 'legacy-bootstrap']):
            continue
        text = re.sub(r'```.*?```', '', path.read_text(errors='replace'), flags=re.S)
        for raw in re.findall(r'\]\(([^)]+)\)', text):
            url = urlsplit(raw)
            if url.scheme or not url.path or '<' in raw:
                continue
            target = (path.parent / unquote(url.path)).resolve()
            assert target.exists(), f'{path.relative_to(ROOT)}: {raw}'
            links += 1
    print(f'PASS: {len(migration["files"])} file dispositions; archived bytes/Go sources, '
          f'historic fingerprints and ledger preserved; {generated} generated outputs; '
          f'install destinations/policies unchanged; {links} local Markdown file targets')
    print('Scope: file targets only, not anchors or external URLs; frozen testdata and legacy bootstrap excluded.')


if __name__ == '__main__':
    main()
