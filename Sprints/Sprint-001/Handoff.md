# Sprint-001 handoff

Active Sprint: none
Active Iteration: none
Active Slice: none
Status: `Sprint-001 / SPI-001 / SPS-001` complete

## Completion evidence

- Slice contract: `Sprints/Sprint-001/ScrumIterations.md`
- Contract study: `Sprints/Sprint-001/InstallationContractStudy.md`
- Implementation notes: `Sprints/Sprint-001/implementationNotes.md`
- Final verification: `Verification/VER-SPS-001-005.md`
- Final review: `CodeReview/REV-SPS-001-005.md`
- Current state: `Traceability/CurrentIndex.yaml`
- Relations and history: `Traceability/Relations.yaml` and
  `Traceability/Ledger.ndjson`
- Draft PR: `https://github.com/Hans-Einar/SDP/pull/4`
- Verified product head: `7c78132f31788d7494eb93c6cba5310c5503e416`
- GitHub Actions run: `https://github.com/Hans-Einar/SDP/actions/runs/29298949100`

The final product head passed Linux contracts and the hosted Windows installer
suite. Normal PR checks remain required on every later commit, including this
documentation-only closure integration; GitHub is authoritative for the latest
head status.

## Delivered boundary

`Toolkit/SDP-install.manifest.json` schema v1 is the authoritative 40-entry
installation inventory and policy contract. PowerShell consumes it, exposes a
deterministic JSON plan, initializes only neutral project templates and supports
explicit Toolkit and consuming-project validator modes. Normal GitHub source
archives are sufficient; `.git` is optional and unavailable commit identity is
recorded as null.

Toolkit `0.2.0` remains unreleased. No tag, GitHub Release or release asset was
created, and all publication identities remain null.

## Cross-repository handoff to `gh-sdp`

The `gh-sdp` Steering Group can now consume the public contract rather than
duplicate PowerShell behavior:

- locate `Toolkit/SDP-install.manifest.json` in an extracted SDP source archive;
- implement schema v1 entry selection, ownership, generation, refresh, backup,
  force, initialization, downgrade and preview semantics;
- resolve sources repository-relatively and honor all exclusions;
- compare its normalized action plan with the PowerShell `-PlanJson` reference;
- validate the offline fixture pinned to `gh-sdp`
  `ed205c1ef193ab8a6e5cd1c50e558c3049ce6def` and normalized tree
  `54f0e5854fd34e5d8bcb301f4921b956a2030e61`; and
- keep `gh-sdp` client `0.1.0` identity separate from Toolkit `0.2.0`.

That integration belongs in `Hans-Einar/gh-sdp`; this Sprint did not modify it.

## Exact next step before Toolkit `0.2.0` publication

Keep PR `#4` in draft until the Steering Group accepts the contract handoff and
normal PR review/merge policy is satisfied. Before publishing Toolkit `0.2.0`,
open a separate release-candidate gate: verify the exact release candidate,
create release-specific verification and review records, obtain explicit human
publication authorization, and only then populate release commit/tag/GitHub
Release identities and perform the two-phase tag/Release publication. Do not
reuse this Slice approval as release authorization.
