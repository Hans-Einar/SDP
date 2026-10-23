#!/usr/bin/env python3
"""Build one source-identified, release-independent PG owner trial. No model launch."""
import argparse
import hashlib
import json
import os
import re
from pathlib import Path
import shutil
import subprocess
import tarfile


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--destination', required=True, type=Path)
    parser.add_argument('--revision', required=True)
    parser.add_argument('--source', type=Path, default=Path(__file__).resolve().parents[3])
    parser.add_argument('--go', default='go')
    args = parser.parse_args()
    source = args.source.resolve()
    revision = subprocess.check_output(['git', 'rev-parse', '--verify', args.revision + '^{commit}'], cwd=source, text=True).strip()
    dest = args.destination.expanduser().absolute()
    if dest.exists():
        parser.error('destination exists; choose a new directory to preserve prior trials')
    dest.mkdir(parents=True, mode=0o700)
    candidate = dest / 'candidate'
    snapshot = candidate / 'source'
    binary = candidate / 'bin'
    snapshot.mkdir(parents=True)
    binary.mkdir()
    archive = candidate / 'source.tar'
    subprocess.run(['git', 'archive', '--format=tar', '--output', str(archive), revision,
                    'SDPTool', 'SDL', 'SDUI', 'Skills'], cwd=source, check=True)
    with tarfile.open(archive) as bundle:
        bundle.extractall(snapshot, filter='data')
    go = shutil.which(args.go)
    if go is None:
        raise RuntimeError('Go executable unavailable')
    with (candidate / 'build.log').open('w') as log:
        for name in ['sdptool', 'sdp-controller', 'sdp-mcp']:
            subprocess.run([go, 'build', '-buildvcs=false', '-o', str(binary / name), './cmd/' + name],
                           cwd=snapshot / 'SDPTool', check=True, stdout=log, stderr=subprocess.STDOUT)
    sha = lambda p: hashlib.sha256(p.read_bytes()).hexdigest()
    harness = Path(__file__).resolve().parent
    for name in ['console.py', 'seed.py', 'README.md']:
        shutil.copyfile(harness / name, dest / name)
    # Describe the archived engine, not the possibly newer preparation harness.
    compatibility = snapshot / 'SDPTool/codexcompat/version.go'
    if compatibility.exists():
        matches = re.findall(r'const MinimumVersion = "([0-9]+\.[0-9]+\.[0-9]+)"', compatibility.read_text())
        if len(matches) != 1:
            raise RuntimeError('Selected source has no unambiguous minimum host version')
        requirement = 'stable codex-cli >= ' + matches[0] + '; bounded wire compatibility required'
    else:
        matches = re.findall(r'const HostVersion = "([0-9]+\.[0-9]+\.[0-9]+)"', (snapshot / 'SDPTool/governancehost/client.go').read_text())
        if len(matches) != 1:
            raise RuntimeError('Selected source has no unambiguous pinned host version')
        requirement = 'codex-cli ' + matches[0]
    identity = {'sourceCommit': revision, 'sourceRepository': str(source), 'archiveSHA256': sha(archive),
                'go': subprocess.check_output([go, 'version'], text=True).strip(),
                'hostRequirement': requirement, 'releasedInstallation': False,
                'binaries': {name: sha(binary / name) for name in ['sdptool', 'sdp-controller', 'sdp-mcp']},
                'harness': {name: sha(dest / name) for name in ['console.py', 'seed.py', 'README.md']}}
    (dest / 'candidate.json').write_text(json.dumps(identity, indent=2) + '\n')
    env = dict(os.environ)
    env['PATH'] = str(Path(go).parent) + os.pathsep + env.get('PATH', '')
    subprocess.run(['python3', str(dest / 'seed.py'), str(dest)], env=env, check=True)
    (dest / 'lab').write_text('#!/bin/sh\ncd -- "$(dirname -- "$0")" || exit 1\nexec python3 ./console.py "$@"\n')
    (dest / 'lab').chmod(0o700)
    print('Ready. Inspect with:', dest / 'lab', 'status')
    print('No model task has been started. Read README.md before running the trial.')


if __name__ == '__main__':
    main()
