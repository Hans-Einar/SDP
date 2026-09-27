#!/usr/bin/env python3
"""Compare captured read-only legacy/Go adoption actions, never run installers."""
import argparse
import base64
import hashlib
import json
from pathlib import Path

p = argparse.ArgumentParser(description=__doc__)
p.add_argument('--legacy', type=Path, required=True)
p.add_argument('--go', type=Path, required=True)
p.add_argument('--output', type=Path, required=True)
a = p.parse_args()
old = json.loads(a.legacy.read_bytes())
new = json.loads(a.go.read_bytes())
assert old['canApply'] and new['directAndGhApplyPassed']
# Receipt 3.0 intentionally replaces legacy facts 2.0 and is journal finalization
# rather than an ordinary Go preview action. All other actions must correspond.
excluded = {'SDP/Framework/installed-toolkit.manifest.yaml'}
x = {v['destination']: v for v in old['actions'] if v['destination'] not in excluded}
y = {v['path']: v for v in new['plannedActions']}
assert set(x) == set(y), (set(x) - set(y), set(y) - set(x))
semantic_json = []
for key in sorted(x):
    assert x[key]['action'] == y[key]['action'], key
    assert x[key]['before'] == y[key]['before'], key
    if 'jsonValue' in y[key]:
        assert json.loads(base64.b64decode(x[key]['content'])) == y[key]['jsonValue'], key
        semantic_json.append(key)
    else:
        assert x[key]['after'] == y[key]['after'], key
proof = dict(schemaVersion='gip-legacy-comparison/1', status='passed',
             legacyCaptureSHA256=hashlib.sha256(a.legacy.read_bytes()).hexdigest(),
             goEvidenceSHA256=hashlib.sha256(a.go.read_bytes()).hexdigest(),
             actionCount=len(x), exactPayloads=len(x)-len(semantic_json),
             jsonKeyOrderingOnly=semantic_json, excludedReceipt=list(excluded),
             policyDifferences=['Go requires named managed refresh for adoption; legacy uses global ForceManagedFiles',
                                'Go receipt 3.0 pins descriptor/provenance; legacy facts 2.0 pin configuration',
                                'Go observes directory types as well as file hashes'])
a.output.write_text(json.dumps(proof, indent=2, sort_keys=True)+'\n')
print(f'PASS: {len(x)} matching actions; {len(semantic_json)} JSON-order normalizations; new receipt intentionally separate')
