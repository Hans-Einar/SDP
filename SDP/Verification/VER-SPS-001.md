# VER-SPS-001 — Canonical installation contract verification

Status: passed for independent Slice review
Slice: `SPS-001`
Release context: `REL-0.2.0` (`0.2.0`, unreleased)
Candidate commit: `ce9278b8b78f9c320a65799aefd13101582d1eb8`
Verified at: `2026-07-13T21:48:40Z`

This is Slice verification, not the final Toolkit `0.2.0` release gate. It does
not claim a release-preparation commit, tag, GitHub Release or publication.

## Environment

- OS: Microsoft Windows `10.0.26200`
- PowerShell: Windows PowerShell `5.1.26100.8655`
- Python: `3.11.5`
- Git: `2.43.0.windows.1`
- Source: clean branch `codex/sdp-install-contract-v1`
- Base: `origin/main` at `bc110bb5fd60009ba67015cf640ad6ddbfe1b04b`

## Exact commands and results

### Toolkit contracts

```powershell
python Toolkit\scripts\validate_sdp.py --mode toolkit --base-ref origin/main
```

Result: exit `0`, `SDP Toolkit validation passed`.

The same candidate also passed the backward-compatible bare invocation:

```powershell
python Toolkit\scripts\validate_sdp.py
```

### Python tests

```powershell
python -m unittest discover -s Toolkit\tests -p 'test_*.py' -v
```

Result: `49` tests passed in `1.567s` on the Master rerun. Coverage includes:

- installation-manifest schema, Toolkit-version agreement and generator facts;
- unique IDs/destinations, normalized relative paths, traversal/absolute-path
  rejection, missing sources, exclusions and complete managed inventory;
- neutral template content and absence of active repository records;
- installation-plan schema and mutation truth;
- valid consuming project without Toolkit-repository layout;
- malformed/unsupported project and installed manifests;
- installed skill version and minimum-Toolkit agreement;
- CurrentIndex, Relations and deterministic path/ID checks;
- empty, malformed, generic, specialized release and namespaced project Ledger
  events;
- release/Fix records and Unreleased-note structure; and
- canonical generic-Ledger and install-plan examples.

### PowerShell installer fixtures

```powershell
.\Toolkit\tests\Install-SDP.Tests.ps1
```

Result: exit `0`, `Installer fixture tests passed`.

The suite exercised empty preview/apply, normal repeat, non-empty preservation,
managed refresh, force and backup, legacy/pre-versioning migration, AGENTS
migration conflict, unsupported schemas, downgrade prevention and sibling/child
path handling. It additionally proved:

- `-PlanJson` conforms to `SDP-install-plan.schema.json` and is mutation-free;
- human preview is mutation-free;
- PowerShell consumes the canonical 40-entry manifest and rejects unsafe policy,
  generator and semantic-destination drift before mutation;
- a `gh-sdp`-like bootstrap preserves Mandate, project manifest, release notes,
  CurrentIndex, Relations and Ledger;
- exactly one `.git` remains at the project root;
- initialization installs only neutral templates;
- `SDP/Releases/REL-0.2.0.yaml` is absent;
- no active Toolkit Sprint, release, review or verification history is copied;
- repeated initialization preview and plan do not mention `REL-0.2.0`; and
- a copied source-archive fixture with no `.git` installs successfully, records
  `sourceCommit: null` and passes consuming-project validation.

### Diff, state and publication truth

```powershell
git diff --check origin/main...HEAD
git status --porcelain
git tag --list
gh release list --repo Hans-Einar/SDP --limit 20
```

Results: diff check passed; Python caches were removed; the worktree was clean;
no tags and no GitHub Releases existed.

## Contract evidence

- `Toolkit/SDP-install.manifest.json`: schema `1.0`, 40 explicit entries, two
  generators and 22 explicit exclusions.
- `Toolkit/project-templates/`: 22 neutral source files and no Ledger, Release,
  Sprint/Refactor, review or verification record.
- `Toolkit/scripts/Install-SDP.ps1`: manifest-driven inventory and policy,
  deterministic `-PlanJson`, compatibility entry point retained.
- `Toolkit/scripts/validate_sdp.py`: explicit Toolkit and consuming-project modes.
- Normal GitHub source archive is sufficient; no custom release asset is needed
  based on the passing no-`.git` fixture.

## Limitations and pending evidence

- Local evidence is Windows-based. Linux contract execution remains required from
  the draft pull request's GitHub Actions run.
- Independent review is pending. This record alone does not complete `SPS-001`.
- Final Toolkit `0.2.0` release-gate verification remains future work; release
  state and publication identities remain unchanged.
