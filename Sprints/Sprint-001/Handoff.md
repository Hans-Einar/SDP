# Sprint-001 handoff

Active Sprint: none
Active Iteration: none
Active Slice: none
Status: `Sprint-001 / SPI-001 / SPS-001 revision 1` complete

## Completion evidence

- Slice contract: `Sprints/Sprint-001/ScrumIterations.md`
- Contract study: `Sprints/Sprint-001/InstallationContractStudy.md`
- Implementation notes: `Sprints/Sprint-001/implementationNotes.md`
- Final verification: `Verification/VER-SPS-001-006.md`
- Final review: `CodeReview/REV-SPS-001-007.md`
- Downstream trigger review: `CodeReview/REV-SPS-001-006.md`
- Current state: `Traceability/CurrentIndex.yaml`
- Relations and history: `Traceability/Relations.yaml` and
  `Traceability/Ledger.ndjson`
- Draft PR: `https://github.com/Hans-Einar/SDP/pull/4`
- Reviewed and verified contract head:
  `a37a8fa298fa3ace2ec5826bc3a30f26f90a64ac`
- GitHub Actions:
  `https://github.com/Hans-Einar/SDP/actions/runs/29590324256`

The exact contract head passed Linux contracts, the hosted Windows installer
suite and all 17 reference-conformance scenarios. Normal PR checks remain
required on every later documentation/traceability integration commit; GitHub
is authoritative for the latest draft-PR head.

## Delivered revision-1 boundary

The required plan/manifest `orderingPolicy` is
`migration-first-manifest-order-v1`. A blocked plan has one sequence-1 block;
migrations precede ordinary entries; ordinary entries follow exact manifest
array order; backup/mutation pairs are adjacent and identity-equal; clients do
not privately sort; sequence is assigned last and is contiguous.

AGENTS conflict preservation hashes exact bytes into
`AGENTS-project.migration-sha256-<sha256>.md`. An absent destination emits one
migration, an identical regular file emits none, and different content or an
unsupported object fails with a stable machine-readable class before mutation.
Apply rechecks both assumptions, uses exclusive creation and never overwrites
project-owned content.

`Toolkit/conformance/install-v1/` is the shared language-neutral authority. Its
17 scenarios cover empty, initialize, repeat, legacy, upgrade, force,
project-owned preservation, archive/no-Git, unsupported schemas, downgrade and
malformed-manifest outcomes. The 11 applicable plans, 3 blocked plans and 3
fatal outcomes are committed authorities. Normal CI validates and compares
them; only the explicit maintainer candidate command can write replacements.

Toolkit `0.2.0` remains unreleased. No tag, GitHub Release, release asset or
merge was created, and every publication identity remains null.

## Exact cross-repository handoff

Send this message to the `gh-sdp` Steering Group:

> Please perform a new bounded installation-contract v1 acceptance pass for
> `Hans-Einar/SDP` draft PR #4. The exact independently reviewed and verified
> contract candidate is
> `a37a8fa298fa3ace2ec5826bc3a30f26f90a64ac`; exact-head GitHub Actions run
> `29590324256` passed Linux contracts, the hosted Windows full installer suite
> and all 17 reference scenarios. Consume
> `Toolkit/SDP-install.manifest.json`, both install schemas, the normative
> installation/migration documentation and
> `Toolkit/conformance/install-v1/scenarios.json` plus its 17 committed expected
> outcomes. Reassess canonical manifest-order planning, backup adjacency,
> exact-byte AGENTS collision/idempotency/race behavior and language-neutral
> fixture consumption. Toolkit `0.2.0` remains unreleased; no tag or GitHub
> Release exists, and PR #4 must remain draft. Please return `ACCEPTED`,
> `ACCEPTED_WITH_LOW_FINDINGS` with no merge-blocking item, or
> `UPSTREAM_REWORK_REQUIRED`, naming the exact assessed head and every finding.

That acceptance work belongs in `Hans-Einar/gh-sdp`; this Sprint did not modify
the downstream repository or implement its client.

## Next gate before Toolkit `0.2.0` publication

Keep PR #4 draft until the Steering Group returns `ACCEPTED` or
`ACCEPTED_WITH_LOW_FINDINGS` with no merge-blocking item and normal PR
review/merge policy is satisfied. Publication then requires a separate
release-candidate gate, exact release-candidate verification, release-specific
review records and explicit human authorization before any release commit, tag
or GitHub Release identity is populated. Do not reuse this Slice approval as
release authorization.
