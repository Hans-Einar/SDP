# Release And Versioning Guide

## Baseline and target

The reusable Toolkit introduced before formal version metadata is treated as the
migration baseline `0.1.0`. This is a documented baseline, not a claim that a
Git tag or GitHub Release was published. The release/versioning capability is a
backward-compatible method addition, so the pending Toolkit version is `0.2.0`.
It remains `unreleased` until an explicitly authorized publication transaction.

Released software and released SDP Toolkits use Semantic Versioning:

```text
MAJOR.MINOR.PATCH
```

- MAJOR: incompatible public contract or SDP method change
- MINOR: backward-compatible capability or method addition
- PATCH: backward-compatible correction

Sprint, Refactor, Iteration, Slice, Fix and revision identifiers never appear in
the SemVer core.

The full accepted form is `MAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`. Numeric core
and numeric prerelease identifiers have no leading zeroes; dot-separated
prerelease and build identifiers contain only ASCII alphanumerics and hyphens.
Prereleases have lower precedence than their matching final release. Build
metadata is part of version identity but is ignored when deciding upgrade or
downgrade precedence.

## Separate identities

The public release identity and the work identity are separate records:

```yaml
releaseVersion: 0.2.0
releaseState: unreleased
sprintId: Sprint-026
refactorId: null
iterationId: SPI-002
sliceId: SPS-003
fixId: null
revision: 1
commit: abc1234
```

Exactly one of `sprintId` and `refactorId` may be non-null. A generated display
may be:

```text
v0.2.0-dev.s026.i002.sl003.r001+sha.abc1234
```

That string is a display identity, not an official release version.

## Version selection

1. Compare the proposed public contract or Toolkit method with the latest
   released version.
2. Select MAJOR, MINOR or PATCH from compatibility impact, not work size.
3. Set the target under `Unreleased` and record `release-version-selected`.
4. Keep state `unreleased` during implementation.
5. A candidate may move to `prerelease` only after the deterministic gate passes.
6. State becomes `released` only after the tag and GitHub Release actually exist
   and a post-publication record has been committed.

## Project and Toolkit ownership

`SDP.manifest.yaml` is Toolkit-owned. In a consuming project,
`SDP/Framework/installed-toolkit.manifest.yaml` is generated and Toolkit-owned.
`SDP/SDP-project.manifest.yaml` and `SDP/RELEASE-NOTES.md` are project-owned and
are only created when missing.

## Deferred skills

- `sdp-code-review` is deferred because `sdp-reviewer` already owns that workflow.
- `sdp-planner` is deferred because planning belongs to `sdp-master` and
  `sdp-architect` until a distinct repeatable workflow emerges.
- `sdp-migration` is deferred because migration is currently executable installer
  behavior rather than a separate agent procedure.
