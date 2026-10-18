"""Verify byte-preserved source history and explicit canonical projections."""
import hashlib
import json
from pathlib import Path
area = Path(__file__).resolve().parent
manifest = json.loads((area / 'import.json').read_text())
ledger = (area.parents[1] / 'Ledger.ndjson').read_bytes()
canonical = {e['eventId']: e for e in map(json.loads, ledger.splitlines())}
def key(e):
    return json.dumps(e, sort_keys=True, separators=(',', ':'))
sources = {}
for name, source in manifest['sources'].items():
    raw = (area / source['path']).read_bytes()
    assert hashlib.sha256(raw).hexdigest() == source['sha256'], name
    sources[name] = list(map(json.loads, raw.splitlines()))
assert ledger.startswith((area / 'release.ndjson').read_bytes())
represented = {key(e) for e in sources['release']}
for row in manifest['events']:
    original = sources[row['source']][row['line'] - 1]
    assert original['eventId'] == row['originalEventId']
    assert original['subjectId'] == row['originalSubjectId']
    current = canonical[row['canonicalEventId']]
    assert current['subjectId'] == row['canonicalSubjectId']
    raw = json.dumps(current, separators=(',', ':')).encode()
    assert hashlib.sha256(raw).hexdigest() == row['canonicalSHA256']
    represented.add(key(original))
for name, records in sources.items():
    assert all(key(e) in represented for e in records), name
print('PASS: exact archives, release prefix, complete event representation and mapped projections')
