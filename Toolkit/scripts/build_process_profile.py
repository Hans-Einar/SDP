#!/usr/bin/env python3
"""Build a reproducible, self-contained process artifact; never mutate a project."""
import argparse
import base64
import hashlib
import json
from pathlib import Path
import re
import yaml

ROOT = Path(__file__).resolve().parents[2]

def canonical(value):
    return (json.dumps(value, ensure_ascii=True, sort_keys=True, separators=(',', ':')) + '\n').encode()

def digest(data):
    return hashlib.sha256(data).hexdigest()

def unique(pairs):
    result = {}
    for key, value in pairs:
        if key in result:
            raise ValueError(f'duplicate key: {key}')
        result[key] = value
    return result

def path(value):
    if not isinstance(value, str) or not value or any(
        not re.fullmatch(r'[A-Za-z0-9_.-]+', x) or x in ('.', '..') or
        x.lower() == '.git' or x.endswith('.') or
        re.fullmatch(r'(?i)(con|prn|aux|nul|com[0-9]|lpt[0-9])(\..*)?', x)
        for x in value.split('/')
    ):
        raise ValueError(f'unsafe portable path: {value!r}')
    return value

def reserved(value, board_mapping=False):
    protected = ['SDP/.sdp-operations', 'SDP/.sdp-backups',
                 'SDP/Framework/installed-toolkit.manifest.yaml',
                 'SDP/ProjectManagement/Ledger.ndjson', 'SDP/navigation.json',
                 'SDP/Traceability/Ledger.ndjson']
    if not board_mapping:
        protected.append('SDP/KanBan/board.json')
    key = value.lower()
    for p in protected:
        p = p.lower()
        if key == p or key.startswith(p+'/') or p.startswith(key+'/'):
            raise ValueError('engine-owned destination or ancestor/descendant')

def shape(obj, fields):
    if not isinstance(obj, dict) or set(obj) != set(fields.split()):
        raise ValueError(f'unknown/missing fields: expected {fields}')

def build(config, root=ROOT):
    shape(config, 'schemaVersion profile managementProfile prerequisites capabilities files relocations')
    if config['schemaVersion'] != '2.0' or config['profile'] != 'sdp-five-phase/0.1' or config['managementProfile'] != 'sdp-project-management/0.1':
        raise ValueError('unsupported profile/schema')
    if config['prerequisites'] != ['powershell>=7.4']:
        raise ValueError('unsupported prerequisites')
    if not isinstance(config['capabilities'], list) or len(set(config['capabilities'])) != len(config['capabilities']) or any(not re.fullmatch(r'sdp\.[a-z0-9.-]+\.v[0-9]+', x) for x in config['capabilities']):
        raise ValueError('invalid capabilities')
    facts = yaml.safe_load((root/'SDP.manifest.yaml').read_text())
    inventory = []
    seen = set()
    for item in config['files']:
        shape(item, 'source destination ownership')
        source, dest = path(item['source']), path(item['destination'])
        if item['ownership'] not in ('managed', 'project'):
            raise ValueError('unknown ownership')
        if not source.startswith(('Skills/', 'Toolkit/payload/', 'Template/')):
            raise ValueError('source outside distribution')
        if item['ownership'] == 'managed' and source.startswith('Template/'):
            raise ValueError('templates must remain project owned')
        reserved(dest)
        if dest != 'AGENTS.md' and not dest.startswith(('SDP/', '.codex/skills/')):
            raise ValueError('destination outside installation scope')
        if item['ownership'] == 'managed' and not (dest == 'AGENTS.md' or dest.startswith(('.codex/skills/', 'SDP/Framework/'))):
            raise ValueError('invalid managed destination')
        key = dest.lower()
        if any(key == p or key.startswith(p+'/') or p.startswith(key+'/') for p in seen):
            raise ValueError('duplicate, overlapping or case-colliding destination')
        seen.add(key)
        src = root/source
        if not src.is_file() or any(p.is_symlink() for p in [src, *src.parents]):
            raise ValueError('missing or linked source')
        data = src.read_bytes()
        inventory.append(dict(item, sha256=digest(data), content=base64.b64encode(data).decode()))
    # Every claimed skill must actually be in the bundle.
    for skill in facts['skills']:
        expected = f'.codex/skills/{skill}/SKILL.md'
        if not any(f['destination'] == expected and f['source'] == f'Skills/{skill}/SKILL.md'
                   and f['ownership'] == 'managed' for f in inventory):
            raise ValueError(f'missing declared skill: {skill}')
    seen_moves = set()
    for move in config['relocations']:
        shape(move, 'from to')
        a, b = path(move['from']), path(move['to'])
        reserved(a)
        reserved(b, board_mapping=(a == 'SDP/Agents/KanBan' and b == 'SDP/KanBan'))
        if not a.startswith('SDP/') or not b.startswith('SDP/') or a.lower() == b.lower() or b.lower().startswith(a.lower()+'/'):
            raise ValueError('unsafe relocation')
        for k in (a.lower(), b.lower()):
            if any(k == p or k.startswith(p+'/') or p.startswith(k+'/') for p in seen_moves):
                raise ValueError('overlapping relocation')
            seen_moves.add(k)
    payload = dict(config)
    payload['files'] = inventory
    payload['facts'] = dict(toolkitVersion=facts['toolkit']['version'], frameworkVersion=facts['framework']['version'],
        agentsContractVersion=facts['agentsContract']['version'], installerVersion=facts['toolkit']['minimumInstallerVersion'],
        skills=facts['skills'], capabilities=sorted(set(facts['capabilities']+config['capabilities'])))
    payload['sourceDigest'] = digest(canonical(config) + (root/'SDP.manifest.yaml').read_bytes())
    payload['configurationDigest'] = digest(canonical(payload))
    return payload

def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--config', type=Path, default=ROOT/'Toolkit/profiles/five-phase.json')
    p.add_argument('--output', type=Path, required=True)
    a = p.parse_args()
    artifact = build(json.loads(a.config.read_text(), object_pairs_hook=unique))
    a.output.parent.mkdir(parents=True, exist_ok=True)
    a.output.write_bytes(canonical(artifact))

if __name__ == '__main__':
    main()
