# VER-SPS-001 — SPS-001 Documentation Verification

Status: passed — independently confirmed for SPS-001 closure
Slice: SPS-001
Mandate: MAN-001
Release context: REL-0.1.0 (gh-sdp 0.1.0, proposed and unreleased)
Reviewed candidate: 381b38ed6fc59936e7e2cb30b877e156ff57914c
Pull request: https://github.com/Hans-Einar/gh-sdp/pull/1 (draft)
Exact candidate rerun: 2026-07-13T10:06:44Z
Remediation working-tree validation: 2026-07-13T10:06:44Z
Committed remediation candidate: bb16bee815d40742f29e15eeac4254bd1310376e
Committed remediation rerun: 2026-07-13T10:09:41Z
Follow-up reviewed candidate: f9d97c5a19d7b7f035f8dbfa2dd1f9e464b4378f
Follow-up review completed: 2026-07-13T10:16:59Z

This is passed Slice evidence, not a product-release gate. It distinguishes the
initial reviewed candidate, committed remediation, and exact pushed follow-up
candidate. Fresh review approved `SPS-001` closure; no Steering approval, release
approval, tag, GitHub Release, published capability, Toolkit support declaration,
or Phase 2 authorization is claimed.

## Environment

- OS: Microsoft Windows 10.0.26200
- PowerShell: 5.1.26100.8655
- Git: 2.43.0.windows.1
- Python: 3.11.5
- GitHub CLI: 2.96.0
- Target repository: C:\Users\hanse\GIT\gh-sdp
- Canonical SDP clone: C:\Users\hanse\GIT\SDP

## Authoritative Source Evidence

Commands:

~~~powershell
git -C C:\Users\hanse\GIT\SDP fetch --prune origin
git -C C:\Users\hanse\GIT\SDP pull --ff-only origin main
git -C C:\Users\hanse\GIT\SDP rev-parse HEAD
git -C C:\Users\hanse\GIT\SDP status --short --branch
git -C C:\Users\hanse\GIT\SDP tag --list
gh release list --repo Hans-Einar/SDP --limit 20
gh run view 29207275200 --repo Hans-Einar/SDP --json databaseId,headSha,headBranch,event,status,conclusion,url,workflowName,jobs
~~~

Results:

- The clone fast-forwarded from b47c191 to clean origin/main at
  bc110bb5fd60009ba67015cf640ad6ddbfe1b04b.
- SDP.manifest.yaml and canonical documentation identify Toolkit 0.2.0 as
  unreleased, with null tag and release commit.
- git tag --list and gh release list returned no tags or GitHub Releases.
- GitHub Actions run 29207275200 at the exact source commit passed both
  contracts and installer jobs.

Canonical validation commands, run from C:\Users\hanse\GIT\SDP:

~~~powershell
python -m pip install -r Toolkit\tests\requirements.txt
python Toolkit\scripts\validate_sdp.py
python -m unittest discover -s Toolkit\tests -p "test_*.py" -v
.\Toolkit\tests\Install-SDP.Tests.ps1
~~~

Results: canonical validation passed; all 9 Python tests passed; installer
fixture tests passed. Generated Python cache directories were removed afterward,
and the canonical clone was clean. These commands validate the Toolkit source,
not this consuming project.

## Bootstrap Evidence

Commands:

~~~powershell
& 'C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1' -ProjectRoot 'C:\Users\hanse\GIT\gh-sdp' -InitializeProjectStructure -Preview
& 'C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1' -ProjectRoot 'C:\Users\hanse\GIT\gh-sdp' -InitializeProjectStructure
~~~

Results:

- Preview proposed 46 actions, applied none, and left the worktree unchanged.
- Apply reported Toolkit 0.2.0, 41 applied files, one preserved file, three
  unchanged files, and no release or tag.
- The installed manifest records source commit bc110bb5fd60009ba67015cf640ad6ddbfe1b04b.
- The initializer's Toolkit-specific REL-0.2.0 seed and copied Toolkit
  traceability facts were replaced with truthful gh-sdp REL-0.1.0 records.

## Exact Reviewed Candidate Rerun

The following command was run from the target repository after initial review.
All candidate content is read from Git objects, so current remediation edits
cannot affect the result.

~~~powershell
$candidate = '381b38ed6fc59936e7e2cb30b877e156ff57914c'
$base = '2dcc18423a10497678eb19ace8fd9668efda9551'
$resolved = git rev-parse "$candidate^{commit}"
if ($LASTEXITCODE -ne 0 -or $resolved -ne $candidate) { throw 'Candidate resolution failed' }
if ((git rev-parse HEAD) -ne $candidate) { throw 'HEAD is not the reviewed candidate' }
git diff --check $base $candidate
if ($LASTEXITCODE -ne 0) { throw 'Full PR candidate git diff --check failed' }

@'
from pathlib import Path
import hashlib, json, re, subprocess, yaml
from jsonschema import Draft202012Validator, FormatChecker

root = Path(r'C:\Users\hanse\GIT\gh-sdp')
upstream = Path(r'C:\Users\hanse\GIT\SDP')
candidate = '381b38ed6fc59936e7e2cb30b877e156ff57914c'
source = 'bc110bb5fd60009ba67015cf640ad6ddbfe1b04b'

def blob(repo, revision, path):
    return subprocess.run(['git', 'show', f'{revision}:{path}'], cwd=repo,
                          check=True, stdout=subprocess.PIPE).stdout

def target(path): return blob(root, candidate, path)
def source_blob(path): return blob(upstream, source, path)
def target_yaml(path): return yaml.safe_load(target(path).decode('utf-8'))

schema_records = (
    ('SDP/SDP-project.manifest.yaml', 'SDP-project-manifest.schema.json'),
    ('SDP/Framework/installed-toolkit.manifest.yaml', 'installed-toolkit-manifest.schema.json'),
    ('SDP/Releases/REL-0.1.0.yaml', 'release-record.schema.json'),
)
records = []
for record_path, schema_name in schema_records:
    record = target_yaml(record_path)
    schema = json.loads(source_blob(f'Toolkit/schemas/{schema_name}'))
    errors = list(Draft202012Validator(schema, format_checker=FormatChecker()).iter_errors(record))
    assert not errors, [error.message for error in errors]
    records.append(record)
project, installed, release = records
current = target_yaml('SDP/Traceability/CurrentIndex.yaml')
relations = target_yaml('SDP/Traceability/Relations.yaml')
assert project['project']['name'] == current['project']['name'] == 'gh-sdp'
assert project['release']['nextTargetVersion'] == release['version'] == '0.1.0'
assert installed['sourceCommit'] == source
assert release['releasePreparationCommit'] is None
assert release['gitTag'] is None and release['githubReleaseUrl'] is None

ledger_schema = json.loads(source_blob('Toolkit/schemas/release-event.schema.json'))
for line in target('SDP/Traceability/Ledger.ndjson').decode('utf-8').splitlines():
    if line.strip():
        errors = list(Draft202012Validator(ledger_schema, format_checker=FormatChecker()).iter_errors(json.loads(line)))
        assert not errors, [error.message for error in errors]

mandate = target('SDP/01--Mandate/mandate.md').decode('utf-8')
for section in ('outcomes', 'boundaries', 'successCriteria'):
    for stable_id, record in relations[section].items():
        assert stable_id in mandate and record['mandate'] == 'MAN-001'
        target(record['path'])
assert 'assumptions' not in relations and 'questions' not in relations

tracked = subprocess.run(['git', 'ls-tree', '-r', '--name-only', candidate],
    cwd=root, check=True, text=True, encoding='utf-8', stdout=subprocess.PIPE).stdout.splitlines()
forbidden_pattern = re.compile(r'(^|/)(go\.mod|go\.sum|.*\.go|.*\.exe|.*\.dll|.*\.sh|.*\.bat|.*\.cmd|Dockerfile|Makefile)$', re.I)
assert any(path.startswith('.codex/') for path in tracked)
assert not any(path.startswith('.github/') for path in tracked)
assert not [path for path in tracked if forbidden_pattern.search(path)]

pairs = [('AGENTS.md', 'Toolkit/payload/project-root/AGENTS.md.template')]
source_paths = subprocess.run(['git', 'ls-tree', '-r', '--name-only', source,
    '--', 'Toolkit/skills', 'Toolkit/payload/sdp-root/Framework'], cwd=upstream,
    check=True, text=True, encoding='utf-8', stdout=subprocess.PIPE).stdout.splitlines()
for path in source_paths:
    if path.endswith('/SKILL.md'):
        pairs.append((f'.codex/skills/{Path(path).parent.name}/SKILL.md', path))
    elif path.startswith('Toolkit/payload/sdp-root/Framework/'):
        pairs.append((f"SDP/Framework/{path.removeprefix('Toolkit/payload/sdp-root/Framework/')}", path))
for target_path, source_path in pairs:
    assert hashlib.sha256(target(target_path)).digest() == hashlib.sha256(source_blob(source_path)).digest()

future = ('02--Study/study.md', '03--Requirements/requirements.md',
          '04--Architecture/architecture.md', '05--DesignAnalysis/design-analysis.md',
          '06--Design/design.md', '07--Implementation/implementation-plan.md')
for path in future:
    assert hashlib.sha256(target(f'SDP/{path}')).digest() == hashlib.sha256(source_blob(path)).digest()

print(f'PASS exact candidate: {candidate}')
print(f'PASS complete tracked tree: {len(tracked)} files including hidden .codex; no .github or product/build files')
print(f'PASS managed Git-blob hashes: {len(pairs)}; later-lifecycle hashes: {len(future)}')
print('CONFIRMED REVIEW FINDING: candidate lacks assumption/question relation maps')
'@ | python -
if ($LASTEXITCODE -ne 0) { throw 'Exact committed candidate validation failed' }
~~~

Results: commit identity resolved exactly; full PR `git diff --check` from base
`2dcc18423a10497678eb19ace8fd9668efda9551` through candidate
`381b38ed6fc59936e7e2cb30b877e156ff57914c` emitted no output;
the three canonical YAML schemas and one ledger release-event record passed;
version/publication identities passed; all 46 tracked paths were checked,
including hidden `.codex`; no `.github` or product/build paths existed; all 16
managed Git blobs and six later-lifecycle templates matched the pinned source.
The rerun also reproduced the known review finding that the initial committed
candidate did not yet map Mandate assumptions or questions. That gap is
addressed in committed remediation candidate
`bb16bee815d40742f29e15eeac4254bd1310376e` below.

## Remediation Record Validation

The following exact command was run from the target repository. It uses the
canonical schemas and checks cross-record identity, publication truth, stable
IDs, paths, and later-lifecycle state.

~~~powershell
@'
from pathlib import Path
import json
import yaml
from jsonschema import Draft202012Validator, FormatChecker

root = Path(r'C:\Users\hanse\GIT\gh-sdp')
upstream = Path(r'C:\Users\hanse\GIT\SDP')

def load_yaml(path):
    return yaml.safe_load(path.read_text(encoding='utf-8'))

def validate_yaml(label, instance_path, schema_name):
    instance = load_yaml(instance_path)
    schema = json.loads((upstream/'Toolkit/schemas'/schema_name).read_text(encoding='utf-8'))
    errors = sorted(Draft202012Validator(schema, format_checker=FormatChecker()).iter_errors(instance), key=lambda e: list(e.path))
    assert not errors, f'{label}: {[e.message for e in errors]}'
    print(f'PASS schema: {label}')
    return instance

project = validate_yaml('project manifest', root/'SDP/SDP-project.manifest.yaml', 'SDP-project-manifest.schema.json')
installed = validate_yaml('installed Toolkit manifest', root/'SDP/Framework/installed-toolkit.manifest.yaml', 'installed-toolkit-manifest.schema.json')
release = validate_yaml('REL-0.1.0 release record', root/'SDP/Releases/REL-0.1.0.yaml', 'release-record.schema.json')
current = load_yaml(root/'SDP/Traceability/CurrentIndex.yaml')
relations = load_yaml(root/'SDP/Traceability/Relations.yaml')
assert isinstance(current, dict) and isinstance(relations, dict)
print('PASS YAML objects: CurrentIndex and Relations')

event_schema = json.loads((upstream/'Toolkit/schemas/release-event.schema.json').read_text(encoding='utf-8'))
ledger_lines = [line for line in (root/'SDP/Traceability/Ledger.ndjson').read_text(encoding='utf-8').splitlines() if line.strip()]
assert ledger_lines, 'Ledger is empty'
for number, line in enumerate(ledger_lines, 1):
    event = json.loads(line)
    errors = list(Draft202012Validator(event_schema, format_checker=FormatChecker()).iter_errors(event))
    assert not errors, f'ledger {number}: {[e.message for e in errors]}'
print(f'PASS release-event schema: {len(ledger_lines)} ledger line(s)')

assert project['project']['name'] == current['project']['name'] == 'gh-sdp'
assert project['release']['currentVersion'] == '0.0.0'
assert project['release']['nextTargetVersion'] == current['release']['targetVersion'] == release['version'] == '0.1.0'
assert project['release']['state'] == current['release']['state'] == release['state'] == 'unreleased'
assert current['release']['activeReleaseId'] == release['releaseId'] == 'REL-0.1.0'
assert project['development']['sprintId'] == current['active']['sprint'] == 'Sprint-001'
assert project['development']['iterationId'] == current['active']['iteration'] == 'SPI-001'
assert project['development']['sliceId'] == current['active']['slice'] == 'SPS-001'
assert installed['toolkitVersion'] == '0.2.0'
assert installed['sourceCommit'] == 'bc110bb5fd60009ba67015cf640ad6ddbfe1b04b'
assert project['release']['latestTag'] is None and project['release']['latestCommit'] is None
assert release['releasePreparationCommit'] is None and release['gitTag'] is None and release['githubReleaseUrl'] is None
assert relations['requirements'] == {} and relations['designs'] == {}
assert 'REL-0.1.0' in relations['releases'] and not (root/'SDP/Releases/REL-0.2.0.yaml').exists()
print('PASS cross-record identity, lifecycle, and publication truth')

mandate_text = (root/'SDP/01--Mandate/mandate.md').read_text(encoding='utf-8')
mandate_relation = relations['mandates']['MAN-001']
for section in ('outcomes','boundaries','successCriteria','assumptions','questions'):
    assert set(mandate_relation[section]) == set(relations[section]), section
    for stable_id, record in relations[section].items():
        assert stable_id in mandate_text, stable_id
        assert record['mandate'] == 'MAN-001'
        assert (root/record['path']).is_file(), record['path']
for category in ('mandates','sprints','iterations','slices','reviews','verification'):
    for stable_id, record in relations[category].items():
        if 'path' in record:
            assert (root/record['path']).is_file(), f'{stable_id}: {record["path"]}'
for key in ('releaseNotes','manifest','releaseRecord'):
    assert (root/relations['releases']['REL-0.1.0'][key]).is_file(), key
print('PASS stable-ID and relation path resolution')

future_templates = [
    'SDP/02--Study/study.md','SDP/03--Requirements/requirements.md',
    'SDP/04--Architecture/architecture.md','SDP/05--DesignAnalysis/design-analysis.md',
    'SDP/06--Design/design.md','SDP/07--Implementation/implementation-plan.md',
]
for rel in future_templates:
    text = (root/rel).read_text(encoding='utf-8')
    assert 'Status: template' in text, rel
print('PASS later lifecycle remains template-only')
print('PHASE 1 RECORD VALIDATION PASSED')
'@ | python -
~~~

Result: all three YAML records passed their canonical schemas; CurrentIndex and
Relations parsed; the one ledger line passed the release-event schema; identity,
lifecycle, publication, relation, ID, path, and template checks passed. The
expanded relation check resolved all six outcomes, eight boundaries, seven
success criteria, four assumptions, and seven questions to `MAN-001` and its
real path.

## Managed Files, Layout, and Scope

The following exact command compared SHA-256 hashes with the pinned source and
checked the skill set, untouched future templates, Git layout, and no-product
boundary.

~~~powershell
@'
from pathlib import Path
import hashlib

root = Path(r'C:\Users\hanse\GIT\gh-sdp')
upstream = Path(r'C:\Users\hanse\GIT\SDP')

def digest(path):
    return hashlib.sha256(path.read_bytes()).hexdigest()

pairs = [(upstream/'Toolkit/payload/project-root/AGENTS.md.template', root/'AGENTS.md')]
for source in sorted((upstream/'Toolkit/skills').glob('*/SKILL.md')):
    pairs.append((source, root/'.codex/skills'/source.parent.name/'SKILL.md'))
framework_source = upstream/'Toolkit/payload/sdp-root/Framework'
for source in sorted(path for path in framework_source.rglob('*') if path.is_file()):
    pairs.append((source, root/'SDP/Framework'/source.relative_to(framework_source)))
for source, target in pairs:
    assert target.is_file(), target
    assert digest(source) == digest(target), f'mismatch: {target}'
print(f'PASS managed SHA-256 equality: {len(pairs)} files')

expected_skills = {p.parent.name for p in (upstream/'Toolkit/skills').glob('*/SKILL.md')}
actual_skills = {p.parent.name for p in (root/'.codex/skills').glob('*/SKILL.md')}
assert actual_skills == expected_skills
print(f'PASS managed skill set: {len(actual_skills)} skills')

future_pairs = {
    'SDP/02--Study/study.md': '02--Study/study.md',
    'SDP/03--Requirements/requirements.md': '03--Requirements/requirements.md',
    'SDP/04--Architecture/architecture.md': '04--Architecture/architecture.md',
    'SDP/05--DesignAnalysis/design-analysis.md': '05--DesignAnalysis/design-analysis.md',
    'SDP/06--Design/design.md': '06--Design/design.md',
    'SDP/07--Implementation/implementation-plan.md': '07--Implementation/implementation-plan.md',
}
for target_rel, source_rel in future_pairs.items():
    assert digest(root/target_rel) == digest(upstream/source_rel), target_rel
print('PASS untouched later-lifecycle template hashes: 6 files')
print('MANAGED AND PLACEHOLDER HASH VALIDATION PASSED')
'@ | python -

$gitRoots = @(Get-ChildItem -LiteralPath 'C:\Users\hanse\GIT\gh-sdp' -Force -Recurse -Directory -Filter .git -ErrorAction Stop | ForEach-Object FullName)
if ($gitRoots.Count -ne 1 -or $gitRoots[0] -ne 'C:\Users\hanse\GIT\gh-sdp\.git') { throw "Unexpected Git roots: $($gitRoots -join ', ')" }

$files = @(git ls-files --cached --others --exclude-standard)
if ($LASTEXITCODE -ne 0) { throw 'Unable to enumerate the complete repository file set' }
$codexFiles = @($files | Where-Object { $_ -match '^\.codex[\\/]' })
if ($codexFiles.Count -eq 0) { throw 'Hidden .codex paths were not enumerated' }
$githubFiles = @($files | Where-Object { $_ -match '^\.github[\\/]' })
if ($githubFiles.Count -ne 0) { throw "Forbidden workflow files: $($githubFiles -join ', ')" }
$forbidden = @($files | Where-Object { $_ -match '(^|[\\/])(go\.mod|go\.sum|.*\.go|.*\.exe|.*\.dll|.*\.sh|.*\.bat|.*\.cmd|Dockerfile|Makefile)$' })
if ($forbidden.Count -ne 0) { throw "Forbidden product/build files: $($forbidden -join ', ')" }
if (Test-Path -LiteralPath 'SDP\.sdp-backups') { throw 'Unexpected SDP backup directory' }
git diff --check
if ($LASTEXITCODE -ne 0) { throw 'Remediation working-tree git diff --check failed' }
~~~

Results before commit: all 16 copied managed files matched the pinned source; all
10 skills matched; all 6 later-lifecycle documents matched their untouched
templates; exactly one `.git` existed at repository root; all 47 tracked or
untracked files were enumerated, including the 10 hidden `.codex` skill files and
the then-untracked initial review record; no `.github`, product/build, or backup
artifact existed; and remediation working-tree `git diff --check` emitted no
output.

After commit, the same matrix was rerun directly from Git objects at
`bb16bee815d40742f29e15eeac4254bd1310376e`. All 47 committed paths, three YAML
schemas, the release-event ledger line, cross-record and publication identities,
all 32 subordinate Mandate IDs, 16 managed blobs, 10 skills, six later-lifecycle
templates, the full PR diff, and the no-product/workflow boundary passed.

The Master and fresh follow-up Reviewer then repeated the Git-object matrix
against exact pushed head `f9d97c5a19d7b7f035f8dbfa2dd1f9e464b4378f`.
The worktree and origin matched that SHA; all 47 governance paths, schemas,
ledger, identities, relations, managed blobs, templates, repository layout,
full-PR diff, and no-product/workflow boundary passed.

## Installer Repeatability

Commands:

~~~powershell
& 'C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1' -ProjectRoot 'C:\Users\hanse\GIT\gh-sdp' -Preview
& 'C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1' -ProjectRoot 'C:\Users\hanse\GIT\gh-sdp' -InitializeProjectStructure -Preview
~~~

Results:

- Normal same-version reinstall: Proposed 0, Applied 0, five project-owned files
  preserved, and 20 managed or template files unchanged.
- Repeating the one-time initializer: Proposed 1, namely creation of
  SDP/Releases/REL-0.2.0.yaml; zero actions were applied. This is a canonical
  Toolkit payload/idempotency gap, not a gh-sdp release record. It is mitigated
  by project warnings and preview-first one-time use. Correcting it requires a
  separate authorized change in Hans-Einar/SDP.

## Authoritative Tooling Gaps and Limitations

- Toolkit/scripts/validate_sdp.py validates the Toolkit repository layout. It has
  no consuming-project mode and cannot truthfully be presented as this project's
  validator.
- Canonical schemas cover project, installed, and release manifests plus release
  ledger events, but not Mandate, CurrentIndex, Relations, or generic work and
  review events. The checks above are the strongest available without fabricating
  an authoritative validator.
- Toolkit test dependencies use compatible version ranges rather than a fully
  locked environment.
- Baseline evidence applies directly to reviewed commit
  `381b38ed6fc59936e7e2cb30b877e156ff57914c` and draft PR #1. The four review
  corrections are separately validated at committed candidate
  `bb16bee815d40742f29e15eeac4254bd1310376e`. Fresh follow-up review approved
  exact pushed head `f9d97c5a19d7b7f035f8dbfa2dd1f9e464b4378f` for Slice closure.

## Reviewer Confirmation

Initial review `REV-SPS-001-001` was completed by fresh Reviewer
`/root/phase1_reviewer` against commit
381b38ed6fc59936e7e2cb30b877e156ff57914c and draft PR #1. Its disposition was
changes required: zero blocking, zero high, four medium, and one low finding.
Committed remediation candidate `bb16bee815d40742f29e15eeac4254bd1310376e`
addressed all four medium findings. Fresh Reviewer
`/root/phase1_followup_reviewer` inspected exact pushed head
`f9d97c5a19d7b7f035f8dbfa2dd1f9e464b4378f` and approved `SPS-001` closure in
`REV-SPS-001-002`: zero blocking, zero high, zero medium, and one low finding.
The low empty-repository-description finding remains for Steering Group
disposition. Terminal confirmation of the later closure-metadata commit is
recorded externally on draft PR #1 because a tracked record cannot contain the
hash of the commit that contains itself.
