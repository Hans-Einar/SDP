# SDP installation contract v1 conformance package

This directory is the checked-in, language-neutral conformance authority for
`sdp-install` schema v1. A client can consume `scenarios.json`, the scenario
schema and committed `expected/` outcomes without executing or inspecting
PowerShell.

## Outcome model

Each scenario declares portable source mode, options, a deterministic before
state, one expected outcome and preservation assertions:

- `applicable-plan` references a plan with `canApply: true`;
- `blocked-plan` references a schema-valid plan with `canApply: false` and
  exactly one sequence-1 `block` action;
- `fatal` references `{ "kind": "fatal", "failureClass": "..." }` from the
  closed vocabulary in the index.

Before-state `profile` values are normative shortcuts:

- `empty` creates an empty project root;
- `installed-default` materializes every default manifest entry;
- `installed-initialize` materializes every default and initialize-only entry.

Profiles use copied bytes from the canonical manifest, an empty-ledger generator
and an installed-manifest generator with fixed
`toolkitInstalledAt: "2026-01-01T00:00:00Z"` and `sourceCommit: null`.
`installedManifestFacts` overrides declared facts, then explicit `files` and
`directories` overlay the profile. `contentBase64` carries exact bytes; clients
must not normalize line endings or decode/re-encode those bytes.

All scenarios use `archive-no-git`. No scenario, expected outcome or assertion
contains an absolute path, temporary path or volatile output. The invalid
manifest scenario uses one declared JSON replacement rather than a private
PowerShell mutation.

## Canonical ordering and AGENTS behavior

Every expected plan declares
`orderingPolicy: "migration-first-manifest-order-v1"`. The expected files are
authoritative reviewed artifacts. Normal validation compares implementation
output to them and never regenerates them.

AGENTS migration actions carry:

- `targetSource: "AGENTS.md"`;
- `targetSourceSha256`, computed over the exact source bytes;
- `destinationPrecondition: "absent"`.

An already-preserved identical regular file emits no migration action.
Different bytes and unsupported objects are fatal pre-plan outcomes. A source
or destination change during apply uses the separate race failure classes
listed in `scenarios.json`.

## Run

Validate the language-neutral package without PowerShell:

```powershell
python Toolkit/conformance/install-v1/run_conformance.py --validate-only
```

Run every scenario against a PowerShell reference:

```powershell
python Toolkit/conformance/install-v1/run_conformance.py --powershell powershell.exe
python Toolkit/conformance/install-v1/run_conformance.py --powershell pwsh
```

The harness materializes disposable trees, invokes PowerShell through a wrapper
that serializes `Exception.Data["sdpFailureClass"]` across the process boundary,
validates and compares plans, proves planning is mutation-free, applies
applicable plans, checks preservation assertions and cleans up.

Maintainers may intentionally write candidate expected files with
`--write-candidates`. That command is never used by normal tests or CI. Review
the complete diff, rerun without `--write-candidates`, and obtain independent
review before accepting changed authorities.
