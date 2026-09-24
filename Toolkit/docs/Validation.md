# Validation

`Toolkit/scripts/validate_sdp.py` has explicit Toolkit-repository and installed-
project modes. Install its Python dependencies first:

```powershell
python -m pip install -r Toolkit\tests\requirements.txt
```

## Toolkit mode

```powershell
python Toolkit\scripts\validate_sdp.py --mode toolkit
python Toolkit\scripts\validate_sdp.py --mode toolkit --repo C:\path\to\SDP
python Toolkit\scripts\validate_sdp.py --mode toolkit --base-ref origin/main
```

Toolkit mode is the default for backward compatibility. It validates all Draft
2020-12 schemas, release and installation manifest agreement, explicit source
inventory and exclusions, normalized paths, generator and capability facts,
skill metadata, neutral templates, examples, records, release notes and Toolkit
traceability. `--base-ref` additionally checks that released release-note
sections have not changed; it is valid only in Toolkit mode and needs the named
Git baseline.

Installation-plan validation is both structural and semantic. It checks
contiguous ordered sequences, top-level/action version agreement, the normative
reason/action/mutation table, `canApply`/block agreement and exact
entry/source/generator/destination/ownership agreement with the canonical
installation manifest, including the two target-to-target AGENTS migrations.
It also enforces `migration-first-manifest-order-v1`: at most one v1 migration
forms the action prefix, ordinary entries follow exact manifest-array order,
and each backup is adjacent to and identity-equal with its matching replacement
or regeneration. Migration source hashes and destination preconditions are
validated.

Toolkit mode also validates the language-neutral installation package at
`Toolkit/conformance/install-v1/`: scenario schema, unique portable IDs/paths,
closed failure classes, promised category coverage, archive/no-`.git` inputs,
and every committed plan structurally and semantically.

## Installation conformance

Validate the language-neutral package on any platform without invoking
PowerShell:

```powershell
python Toolkit\conformance\install-v1\run_conformance.py --validate-only
```

When Windows PowerShell or PowerShell 7 is available, materialize every scenario,
compare the reference `PlanJson` or stable fatal class to the committed
authority, verify planning made zero mutations, apply applicable plans and check
preservation assertions:

```powershell
python Toolkit\conformance\install-v1\run_conformance.py --powershell powershell
python Toolkit\conformance\install-v1\run_conformance.py --powershell pwsh
```

Normal tests are comparison-only. They never derive expected output from the
PowerShell implementation. A maintainer may explicitly produce reviewed
candidate updates with:

```powershell
python Toolkit\conformance\install-v1\run_conformance.py `
  --powershell powershell `
  --write-candidates
```

This command changes contract authorities and must not run in CI. Review the
complete expected-outcome diff before committing it.

## Consuming-project mode

```powershell
python Toolkit\scripts\validate_sdp.py `
  --mode project `
  --project-root C:\path\to\Project
```

`--project-root` is required in project mode. Schemas come from the Toolkit
running the validator; the consuming project need not contain the Toolkit
repository layout. Project mode validates:

- `SDP/SDP-project.manifest.yaml` and its installed-manifest reference;
- `SDP/Framework/installed-toolkit.manifest.yaml`, supported versions and every
  declared installed skill's metadata/version compatibility;
- CurrentIndex and Relations when present, including deterministic path/ID links,
  review/verification Ledger subjects and required reverse links for
  Slice/Fix-review, Slice/Fix-verification, release-migration and review
  resolution pairs;
- every non-empty Ledger line when the Ledger is present; an empty Ledger is valid;
- project release notes and release/Fix records, including deterministic links,
  record/Relations identity agreement, state-specific publication fields and
  tag/version/release-ID agreement.

Record and relation paths beginning with `SDP/` are project-root-relative.
Other portable record paths are resolved from the project `SDP/` root first,
with project-root fallback for legacy records. This prevents `SDP/SDP/...`
double-prefixing while retaining existing unprefixed project paths. Release and
Fix record filenames and relation paths use canonical `.yaml`; a `.yml` record
is an error and cannot evade discovery.

Relations retain the stable core categories and explicitly support the project
extension categories exercised by the pinned `gh-sdp` consumer: `mandates`,
`outcomes`, `boundaries`, `successCriteria`, `assumptions` and `questions`.
`CurrentIndex.release.targetStatus` is a supported project status, and Slice
`release`, `reviews` and `verification` links may be lists. Further local
categories use the documented `x-*` namespace rather than unqualified keys.

Missing required manifests or release notes fail. Validation never invents
evidence or treats an absent work/review/verification record as completion.

## Ledger schemas and extensions

Every event satisfies `ledger-event.schema.json`. Canonical `release-*` events
also satisfy the stricter `release-event.schema.json`. Generic canonical event
types use `work-*`, `review-*` or `verification-*`.

Project extensions use exactly:

```text
x-<namespace>:<event-name>
```

For example, `x-acme.example:deployment-approved`. Namespaced events still obey
the common envelope. Unknown schema versions are unsupported; extensibility is
limited to boundaries explicitly permitted by each schema.

The offline compatibility regression is a deterministic export of exact
`Hans-Einar/gh-sdp` closure
`ed205c1ef193ab8a6e5cd1c50e558c3049ce6def`. It contains no `.git`, validates in
project mode without network access, and asserts normalized source tree
`54f0e5854fd34e5d8bcb301f4921b956a2030e61`.

## Truthful Markdown boundary

Validation checks canonical-path existence, release-note structure and
deterministic references where a machine-readable contract exists. It does not
claim to understand arbitrary Markdown semantics, infer approval from prose or
prove that a Mandate, design or review is substantively complete.

Success exits `0` and prints a mode-specific pass message. Validation errors are
listed on standard error and exit `1`; command-line misuse exits through the
argument parser.
