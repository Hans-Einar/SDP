#!/usr/bin/env python3
"""Archive WCI3-M2 ONLY when main executes this script after evidence completion.

Preparation of this file is not permission to execute it. No tests/builds are run.
Default output: ORIGINAL native/WCI3-M2-final and WCI3-M2-prior-candidates.
Existing output is never overwritten. All checks finish in staging before promotion.

Required --suite-manifest JSON (main supplies actual completion receipts):
{
  "runs": [
    {"name":"sdui-full", "status":"passed", "exit":0,
     "command":"<exact recorded command>", "cwd":"<actual cwd>",
     "receipt_origin":"captured result.json",
     "log":"/absolute/test.log", "receipt":"/absolute/result.json",
     "artifacts":["/absolute/metadata.json", "/absolute/source-before.sha256",
                  "/absolute/source-after.sha256"]},
    {"name":"sdl-frozen-group", "status":"superseded-timeout", "exit":1,
     "command":"<exact command>", "cwd":"<actual cwd>",
     "receipt_origin":"coordinator/worker observed actual tool return",
     "log":"/absolute/prior-600s.log", "artifacts":[]},
    {"name":"sdl-text-final", "status":"passed", "exit":0,
     "command":"go test -race ./examples/text -count=1 -v",
     "cwd":"/tmp/sdp-sdui-widgets/SDL/go",
     "receipt_origin":"worker observed actual tool return",
     "log":"/tmp/wci3-m2-text-race-final.log", "artifacts":[]},
    {"name":"sdl-nonfyne", "status":"passed", "exit":0,
     "command":"<two exact argv sequences; keep actual GOFLAGS in artifacts>",
     "cwd":"/tmp/sdp-sdui-widgets/SDL/go", "receipt_origin":"results.json",
     "log":"/tmp/wci3-m2-sdl-nonfyne-regressions/remaining-nonfyne.log",
     "receipt":"/tmp/wci3-m2-sdl-nonfyne-regressions/results.json",
     "artifacts":["<each log, .exit, .command.json, README, SHA256SUMS, before/after>"]},
    {"name":"sdl-collections", "status":"passed", "exit":0,
     "command":"go test -race ./examples/collections/... -count=1",
     "cwd":"/tmp/sdp-sdui-widgets/SDL/go", "receipt_origin":"result.json",
     "log":"/tmp/wci3-m2-collections-regression/test.log",
     "receipt":"/tmp/wci3-m2-collections-regression/result.json",
     "artifacts":["<command.json, exit, SHA256SUMS>"]},
    {"name":"sdl-current-group", "status":"pending", "exit":null,
     "command":"<actual command>", "cwd":"/tmp/sdp-sdui-widgets/SDL/go",
     "receipt_origin":"pending actual completion",
     "log":"<actual log>", "receipt":"<actual result.json>",
     "source_before":"<actual before>", "source_after":"<actual after>"},
    {"name":"sdptool-original", "status":"passed", "exit":0,
     "command":"go test ./...", "cwd":"/home/warloc/git/SDP-vNow/SDPTool",
     "receipt_origin":"coordinator reported actual tool return; not script-captured",
     "log":"/tmp/wci3-m2-original-sdptool-test.log", "artifacts":[],
     "notes":"Default environment; concurrent unrelated work and cached unchanged packages."}
  ]
}

Never infer a missing exit or receipt. Coordinator/worker attestation is recorded
as such, not retroactively represented as a subprocess captured by this script.
The initial SDL group's nonzero exit is retained, not relabelled as suite success;
its older four completed packages are prior-candidate evidence only. Current
bridge/values/commands/panes receipts are required separately.
"""
from pathlib import Path, PurePosixPath
import argparse
import gzip
import hashlib
import io
import json
import os
import re
import shutil
import subprocess
import sys
import tempfile

REPO = Path('/home/warloc/git/SDP-vNow')
SOURCE = Path('/tmp/wci3-m2-final-evidence')
NATIVE = REPO / 'SDP/05--Implementation/SDUI/Widgets/native'
RUNS = ('single', 'multiline', 'history', 'readonly', 'pasteguard', 'failures',
        'retention', 'forms', 'required_empty', 'constraints', 'scroll', 'ime',
        'nonmodal-forms')
BINARY_SHA = '0078f72eec7fe61bd41a1b090bb7ca2f26ba48c3a1e54c4eff6bd240c010f001'
INVENTORY_SHA = 'f551e2e3a7f652b93fce8bce076750c768d4a448577c19476e0bac6c1d06c3c6'
EXPECTED_CHECKS, EXPECTED_FILES = 86, 222
PRIOR_ROOTS = tuple(Path('/tmp/wci3-m2-' + n) for n in
                    ('pilot1', 'pilot2', 'redo-pilot', 'fast-pilot'))
EXTENSIONS = {'.ndjson', '.png', '.txt', '.log'}
EXCLUDED = {'cache', '.cache', 'config', 'runtime', 'data', '__pycache__', '.git'}
SUITES = {'sdui-full', 'sdl-frozen-group', 'sdl-text-final', 'sdl-nonfyne',
          'sdl-collections', 'sdptool-original', 'sdl-current-group'}


def require(ok, message):
    if not ok:
        raise ValueError(message)


def raw(path):
    path = Path(path)
    require(path.is_file() and not path.is_symlink(), f'Missing/nonregular evidence: {path}')
    return path.read_bytes()


def sha(data):
    return hashlib.sha256(data).hexdigest()


def load(path):
    return json.loads(raw(path))


def write_json(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, ensure_ascii=False) + '\n')


def records(path):
    data = raw(path)
    if Path(path).suffix == '.gz':
        data = gzip.decompress(data)
    return [json.loads(line) for line in data.splitlines() if line.strip()]


def gzip_bytes(data):
    # No timestamp or original filename; fixed compression and OS-independent header.
    out = io.BytesIO()
    with gzip.GzipFile(filename='', mode='wb', fileobj=out, mtime=0, compresslevel=9) as f:
        f.write(data)
    return out.getvalue()


def retain(src, dest, provenance):
    data = raw(src)
    if src.name == 'events.ndjson':
        dest = dest.with_name(dest.name + '.gz')
        payload = gzip_bytes(data)
        require(gzip.decompress(payload) == data, f'Compression mismatch: {src}')
    else:
        payload = data
    require(not dest.exists(), f'Archive collision: {dest}')
    dest.parent.mkdir(parents=True, exist_ok=True)
    dest.write_bytes(payload)
    provenance.append({'source': str(src), 'path': str(dest), 'rawSHA256': sha(data),
                       'archivedSHA256': sha(payload)})


def retained_tree(src, dest, provenance):
    require(src.is_dir() and not src.is_symlink(), f'Missing evidence directory: {src}')
    count = 0
    for directory, dirs, files in os.walk(src, followlinks=False):
        base = Path(directory)
        dirs[:] = sorted(d for d in dirs if d not in EXCLUDED
                         and not d.startswith('.') and not (base/d).is_symlink())
        for name in sorted(files):
            f = base/name
            if f.suffix in EXTENSIONS and not f.is_symlink():
                retain(f, dest/f.relative_to(src), provenance)
                count += 1
    require(count, f'No retainable evidence: {src}')


def audit_run(path):
    actions = records(path/'actions.ndjson')
    candidates = [a['value'] for a in actions if a.get('kind') == 'candidate']
    require(len(candidates) == 1 and candidates[0]['sha256'] == BINARY_SHA,
            f'Wrong/ambiguous binary: {path}')
    expected = 'forms' if path.name == 'nonmodal-forms' else path.name
    require(candidates[0]['variant'] == 'wci3-' + expected, f'Wrong variant: {path}')
    require([a['value'] for a in actions if a.get('kind') == 'result'] == ['passed'],
            f'Incomplete/failed run: {path}')
    checks = [a['value'] for a in actions if a.get('kind') == 'check']
    require(checks and all(c.get('passed') is True for c in checks), f'Failed checks: {path}')
    require(raw(path/'stderr.txt') == b'', f'Nonempty fixture stderr: {path}')
    require(records(path/'events.ndjson'), f'Missing native events: {path}')
    for item in actions:
        if item.get('kind') == 'screenshot':
            name = item['value']
            require(isinstance(name, str) and Path(name).name == name and name.endswith('.png'),
                    f'Unsafe screenshot reference: {name}')
            require(raw(path/name).startswith(b'\x89PNG\r\n\x1a\n'), f'Missing PNG: {path/name}')
    return len(checks)


def inventory_equivalence(inventory, repo):
    require(sha(raw(inventory)) == INVENTORY_SHA, 'Tested 222-file inventory identity changed')
    manifest = load(inventory)
    rows = manifest['files']
    require(len(rows) == EXPECTED_FILES, 'Wrong tested inventory count')
    seen = set()
    for row in rows:
        name = row['path']; path = PurePosixPath(name)
        require(not path.is_absolute() and '..' not in path.parts and name not in seen,
                f'Unsafe/duplicate inventory path: {name}')
        seen.add(name)
        require(sha(raw(repo/name)) == row['sha256'], f'Original does not match tested file: {name}')
    return {'paths': len(rows), 'allEqual': True, 'testedInventorySHA256': INVENTORY_SHA,
            'scope': 'Exactly the retained tested-file inventory, not the whole dirty repository.'}


def require_pass_log(path):
    data = raw(path).decode('utf-8')
    require(bool(data.strip()) and re.search(r'^ok\s+', data, re.M), f'No completed package: {path}')
    require(not re.search(r'^(FAIL(?:\s|$)|--- FAIL:|panic:|WARNING: DATA RACE)', data, re.M),
            f'Failure marker in successful suite log: {path}')
    return data


def audit_suites(manifest, dest, provenance):
    spec = load(manifest)
    rows = spec['runs']
    require(len(rows) == len(SUITES) and {r['name'] for r in rows} == SUITES,
            'Suite manifest must cover all seven named evidence records exactly')
    for row in rows:
        name = row['name']; base = dest/name
        require(type(row.get('exit')) is int, f'Missing actual exit: {name}')
        require(row.get('command') and row.get('cwd') and row.get('receipt_origin'),
                f'Missing command/cwd/receipt provenance: {name}')
        if name == 'sdl-frozen-group':
            require(row['status'] == 'superseded-timeout' and row['exit'] != 0,
                    'Initial 600s timeout must remain a failure record')
            text = raw(row['log']).decode('utf-8')
            require('test timed out after 10m0s' in text, 'Missing unchanged 600s timeout evidence')
            for package in ('bridge', 'examples/values', 'examples/commands', 'examples/panes'):
                require(re.search(r'^ok\s+\S+/' + re.escape(package) + r'\s+', text, re.M),
                        f'No frozen completed package: {package}')
        else:
            require(row['status'] == 'passed' and row['exit'] == 0, f'Nonpassing suite: {name}')
            text = require_pass_log(row['log'])
        if name == 'sdl-current-group':
            require(row.get('receipt'), 'Current SDL group requires an actual receipt')
            for package in ('bridge', 'examples/values', 'examples/commands', 'examples/panes'):
                require(re.search(r'^ok\s+\S+/' + re.escape(package) + r'\s+', text, re.M),
                        f'No current completed package: {package}')
        if name == 'sdl-text-final':
            require('508.408s' in text and 'go test -race ./examples/text -count=1 -v' == row['command'],
                    'Final text evidence/unchanged timeout command mismatch')
        if name == 'sdptool-original':
            require(row['command'] == 'go test ./...' and row['cwd'] == str(REPO/'SDPTool'),
                    'Do not invent SDPTool flags/environment')
            require('coordinator' in row['receipt_origin'].lower(), 'SDPTool receipt must be attributed')
        if row.get('receipt'):
            receipt = load(row['receipt'])
            exits = [receipt[k] for k in ('exit', 'exit_code') if k in receipt]
            if 'runs' in receipt:
                exits += [r.get('exit') for r in receipt['runs']]
            require(exits and all(type(x) is int and x == row['exit'] for x in exits),
                    f'Receipt exit mismatch: {name}')
            if receipt.get('logSHA256'):
                require(sha(raw(row['log'])) == receipt['logSHA256'], f'Log digest mismatch: {name}')
            if name == 'sdl-current-group':
                require(receipt.get('sourceBytesUnchanged') is True or receipt.get('source_unchanged') is True,
                        'Current SDL group lacks unchanged-source receipt')
                if 'before' in receipt and 'after' in receipt:
                    require(receipt['before'] == receipt['after'], 'Current SDL group source mismatch')
                else:
                    require(row.get('source_before') and row.get('source_after'),
                            'Current SDL group requires source before/after paths')
                    require(raw(row['source_before']) == raw(row['source_after']),
                            'Current SDL group source manifests differ')
        paths = [Path(row['log'])] + ([Path(row['receipt'])] if row.get('receipt') else [])
        paths += [Path(f) for f in row.get('artifacts', [])]
        paths += [Path(row[k]) for k in ('source_before', 'source_after') if row.get(k)]
        for f in dict.fromkeys(paths):
            retain(f, base/f.name, provenance)
        write_json(base/'archive-record.json', row)
    # Verify retained source stability receipts, not only their success assertions.
    sdui = dest/'sdui-full'
    before, after = raw(sdui/'source-before.sha256'), raw(sdui/'source-after.sha256')
    result, meta = load(sdui/'result.json'), load(sdui/'metadata.json')
    require(before == after and result.get('source_unchanged') is True,
            'SDUI source changed during suite')
    require(sha(before) == meta['source_manifest_sha256'] == result['after_manifest_sha256'],
            'SDUI source manifest digest mismatch')
    nf = dest/'sdl-nonfyne'; result = load(nf/'results.json')
    require(result.get('sourceBytesUnchanged') is True and not result.get('changedSourcePaths'),
            'SDL non-Fyne source stability failed')
    require(load(nf/'source-before.json') == load(nf/'source-after.json'), 'SDL source manifests differ')
    for run in result['runs']:
        require(sha(raw(nf/(run['name']+'.log'))) == run['logSHA256'], 'SDL log digest mismatch')
        require(raw(nf/(run['name']+'.exit')).strip() == b'0', 'SDL captured exit not zero')
    cr = dest/'sdl-collections'; result = load(cr/'result.json')
    require(result.get('sourceBytesUnchanged') is True and result['before'] == result['after'],
            'Collections source/module stability failed')
    require(sha(raw(cr/'test.log')) == result['logSHA256'] and raw(cr/'exit').strip() == b'0',
            'Collections log/exit mismatch')
    return rows


def final_inventory(directory):
    rows = [{'path': str(p.relative_to(directory)), 'sha256': sha(raw(p))}
            for p in sorted(directory.rglob('*')) if p.is_file() and p.name != 'inventory.json']
    write_json(directory/'inventory.json', rows)


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument('--suite-manifest', required=True, type=Path)
    parser.add_argument('--source', type=Path, default=SOURCE)
    parser.add_argument('--binary', type=Path, default=Path('/tmp/wci3-m2-final-native'))
    parser.add_argument('--reviewer-binary', type=Path, default=Path('/tmp/wci3-m2-reviewer-final-native'))
    parser.add_argument('--tested-inventory', type=Path, default=Path('/tmp/wci3-m2-candidate-inventory.json'))
    args = parser.parse_args()
    out, old = NATIVE/'WCI3-M2-final', NATIVE/'WCI3-M2-prior-candidates'
    require(not out.exists() and not old.exists(), 'Destination already exists; no overwrite permitted')
    require(NATIVE.is_dir(), 'Original native evidence parent missing')
    for p in (args.binary, args.reviewer_binary):
        require(sha(raw(p)) == BINARY_SHA, f'Expected final 0078 binary: {p}')
    original = inventory_equivalence(args.tested_inventory, REPO)
    checks = {name: audit_run(args.source/name) for name in RUNS}
    require(sum(checks.values()) == EXPECTED_CHECKS, f'Expected exactly 86 native checks, got {checks}')
    invocations = load(args.source/'invocations.json')
    require([r['variant'] for r in invocations] == list(RUNS), 'Incomplete/changed final invocation list')
    for row in invocations:
        name, argv = row['variant'], row['argv']
        require(type(row.get('exit')) is int and row['exit'] == 0, f'Launcher incomplete/failed: {name}')
        def arg(flag):
            positions = [i for i, v in enumerate(argv) if v == flag]
            # IME wrapper has its own --out/display before the nested command.
            require(positions, f'Missing invocation {flag}: {name}')
            return argv[positions[-1] + 1]
        require(arg('--binary') == str(args.binary) and arg('--out') == str(args.source/name),
                f'Invocation source/binary mismatch: {name}')
        require(arg('--variant') == ('forms' if name == 'nonmodal-forms' else name), 'Wrong invocation variant')
        require(('--nonmodal' in argv) == (name == 'nonmodal-forms'), 'Nonmodal invocation mismatch')
        require(row.get('FYNE_THEME') == 'light', 'Unexpected final-run theme')
        require(bool(raw(args.source/(name+'-launcher.log'))), f'Missing launcher transcript: {name}')
    for name in ('ibus.log', 'xim.log', 'engine.log', 'invocation.txt'):
        raw(args.source/'ime-environment'/name)
    require(any(Path(a).name == 'with_ibus.py' for a in invocations[11]['argv']), 'IME lacks real launcher')
    # Keep complete pre-publication staging separate from either final destination.
    with tempfile.TemporaryDirectory(prefix='.wci3-m2-stage-', dir=NATIVE) as tmp:
        stage = Path(tmp); final, prior = stage/out.name, stage/old.name
        final.mkdir(); prior.mkdir(); provenance = []
        for name in RUNS:
            retained_tree(args.source/name, final/name, provenance)
            retain(args.source/(name+'-launcher.log'), final/(name+'-launcher.log'), provenance)
        retain(args.source/'invocations.json', final/'invocations.json', provenance)
        retained_tree(args.source/'ime-environment', final/'ime-environment', provenance)
        audit = NATIVE/'audit_wci2_results.py'
        retain(audit, final/'audit_wci2_results.py', provenance)
        result = subprocess.run([sys.executable, str(final/'audit_wci2_results.py')] +
                                [str(final/n/'events.ndjson.gz') for n in RUNS],
                                check=True, capture_output=True, text=True)
        receipts = {Path(k).parent.name: v for k, v in json.loads(result.stdout).items()}
        require(set(receipts) == set(RUNS) and all(v['exactlyOnce'] is True for v in receipts.values()),
                'Receipt audit did not cover all final runs')
        # Receipt total is derived, never hardcoded from the coordinator estimate of five.
        write_json(final/'terminal-results.json', receipts)
        write_json(final/'checks.json', checks)
        retain(args.tested_inventory, final/'tested-source-inventory.json', provenance)
        write_json(final/'original-equivalence.json', original)
        retain(args.suite_manifest, final/'suite-manifest.json', provenance)
        suites = audit_suites(args.suite_manifest, final/'suites', provenance)
        for source in PRIOR_ROOTS:
            retained_tree(source, prior/source.name, provenance)
        # Preserve exact captured failure/pass bytes; never audit pilots as final acceptance.
        (prior/'README.md').write_text('''# WCI3-M2 prior candidates and diagnostic runs

These records preserve pilot1, pilot2, redo-pilot and fast-pilot provenance.
Failures, timeouts, intermediate passes and changed harness/native candidates are
not final acceptance. Their original actions/events/launcher logs and OS images
are retained selectively. Events are losslessly gzip-compressed with mtime zero.
No caches, preferences, runtime sockets, binaries or configuration are included.
The retranscribed retained-tool-output excerpt (not original stdout) for the earlier
SDL 600-second timeout is retained under ../WCI3-M2-final/suites/
sdl-frozen-group with its actual nonzero exit, alongside the unchanged-limit final
text run. No failure was erased or retroactively relabelled as a pass.
''')
        build = subprocess.check_output(['go', 'version', '-m', str(args.binary)])
        (final/'build-info.txt').write_bytes(build)
        write_json(final/'binary-identity.json', {'sha256': BINARY_SHA, 'binary': str(args.binary),
                   'reviewerBinary': str(args.reviewer_binary), 'sameBytes': True,
                   'note': 'Byte identity verified; this script did not build either binary.'})
        summary = {'runs': len(RUNS), 'checks': sum(checks.values()), 'binarySHA256': BINARY_SHA,
                   'fixtureStderrEmpty': True, 'terminalResults': sum(v['terminalResults'] for v in receipts.values()),
                   'receiptAudit': 'Every published opening has exactly one terminal result; each run closes.',
                   'suiteEvidence': [{'name': r['name'], 'status': r['status'], 'exit': r['exit'],
                                      'receipt_origin': r['receipt_origin']} for r in suites]}
        write_json(final/'verification.json', summary)
        (final/'README.md').write_text(f'''# WCI3-M2 final native evidence

Frozen native binary SHA-256 `{BINARY_SHA}`; main/reviewer bytes verified equal.
Thirteen named native workflows pass {sum(checks.values())} checks with empty
fixture stderr. Exact receipt audit derives {summary['terminalResults']} terminal
results; no estimated receipt count substitutes for that audit. Actual invocations,
launcher logs, isolated IBus environment logs and OS screenshots are retained.
Launcher/environment diagnostics are not silently represented as empty stderr.
Raw event bytes are losslessly compressed, with deterministic gzip headers.

The exact tested {EXPECTED_FILES}-file inventory matches the original workspace at
archival. This scoped comparison is not a claim that the entire repository was
unchanged. Build information was read from the tested binary; no rebuild occurred.

Suite manifests retain actual commands/status provenance. Original SDPTool used
`go test ./...` with its default environment, concurrent unrelated work and cached
unchanged packages. No explicit GOWORK/readonly flags are claimed for that command.
SDL current coverage combines the separately captured current bridge/values/commands/
panes group, final 508.408s text run and non-Fyne/collection follow-ups.
The old group's four passes concern an earlier host and are prior-candidate evidence only. The earlier whole-group 600-second timeout remains
a nonpassing record; the final text command keeps the unchanged default time limit.
SDUI's complete race run retains before/after source manifests and its actual result.
No native/desktop or full-suite claim is inferred solely from package compilation.

This archive script validates recorded evidence; it does not execute native tests
or suites. Independent final archive review remains pending. WCI4/package/release
obligations are not closed by this archive. Prior pilot evidence is retained in
../WCI3-M2-prior-candidates and excluded from the final 86-check count.
''')
        # All retained source bytes must still match after copying/auditing.
        for item in provenance:
            require(sha(raw(item['source'])) == item['rawSHA256'], f'Evidence changed during archive: {item["source"]}')
            item['path'] = str(Path(item['path']).relative_to(stage))
        write_json(final/'retention-provenance.json', provenance)
        inventory_equivalence(args.tested_inventory, REPO)
        for p in (args.binary, args.reviewer_binary):
            require(sha(raw(p)) == BINARY_SHA, 'Binary changed during archive')
        final_inventory(prior); final_inventory(final)
        require(not out.exists() and not old.exists(), 'Destination appeared during staging')
        # Two directories cannot be atomically renamed as one: if the second move
        # fails, first remains complete, never overwritten; report for main review.
        os.rename(prior, old)
        os.rename(final, out)
    print(json.dumps(summary, indent=2))


if __name__ == '__main__':
    try:
        main()
    except (OSError, ValueError, KeyError, TypeError, IndexError, subprocess.CalledProcessError) as exc:
        print(f'Archive refused: {exc}', file=sys.stderr)
        raise SystemExit(1)
