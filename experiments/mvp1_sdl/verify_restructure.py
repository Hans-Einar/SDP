#!/usr/bin/env python3
"""Verify the MPV1 relocation against its recorded statement multiset."""
import hashlib
import json
from pathlib import Path
from audit_inventory import audit


def verify():
    base = Path(__file__).resolve().parent
    root = base / 'SDL/MVP1'
    record = json.loads((base / 'evidence/MPV1/migration.json').read_text())
    statements = sorted(
        line.strip() for path in root.rglob('*.design')
        for line in path.read_text().splitlines()
        if line.strip() and not line.startswith(('language ', 'includes ', '//')))
    digest = hashlib.sha256(('\n'.join(statements) + '\n').encode()).hexdigest()
    result = audit(root)
    errors = list(result['errors'])
    if digest != record['statementSha256'] or len(statements) != record['statementCount']:
        errors.append('Model statements differ from the recorded pre-migration multiset')
    for source, destinations in record['moves'].items():
        for destination in destinations:
            if not (root / destination).is_file():
                errors.append('Missing destination: ' + destination)
        if source not in destinations and (root / source).exists():
            errors.append('Legacy source still exists: ' + source)
    for container in ('BuckingUI', 'SimulatorUI'):
        if not (root / 'Containers' / container / 'Domain.design').is_file():
            errors.append('Missing separate UI container: ' + container)
    result.update(result='fail' if errors else 'pass', errors=errors,
                  preserved_statements=len(statements), statement_sha256=digest)
    return result


if __name__ == '__main__':
    result = verify()
    print(json.dumps(result, indent=2, sort_keys=True))
    raise SystemExit(result['result'] != 'pass')
