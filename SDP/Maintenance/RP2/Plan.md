# RP2 — Publish the external-reference navigation correction

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0008 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| SourceCards | KB-SDP-034 |

## Outcome and authority

Owner explicitly requests merge and release on 2026-09-27. Merge PR42 and publish
SDP 0.2.1 plus the thin gh-sdp client patch selecting it. The backward-compatible
fix keeps valid boards available and adds informational externalReference metadata;
no profile, receipt or installation migration changes. Linux amd64 only, as before.
[ER1](../../05--Implementation/SDPTool/ExternalReferences/Plan.md) proves behavior.

## Milestones and Git policy

Use sdp/release-r2 with milestone commits and PRs to main. Owner authorization
covers these merges and actual releases. Preserve unrelated SDL/go/sourceinput.

- RP2-M1: freeze notes/version, pin default, verify PR checks, exact clean package,
  production signature, install/repeat, read-only XFMD tree, archive install,
  Toolkit/management/document checks. Retain published 0.2.0 bytes.
- RP2-M2: merge preparation, package exact merged source, annotate/push v0.2.1,
  publish immutable assets, verify downloaded signatures/checksums/behavior.
  Coordinate independently reviewed gh-sdp patch. Reconcile actual publication
  identities and verify the released default client sees the corrected XFMD tree.

No XFMD application or installed process files are changed. Root runtime design
is unchanged; the consumer-contract correction remains in ER1. The gh-sdp local
process separately requires Worker and fresh independent review. This release
plan supersedes RP1's default pin for new packages, not its historical evidence.

## Evidence

Preparation candidate 7ff5a4e passes full Go race tests (including bootstrap),
vet, 106 Toolkit tests (19 environment-dependent skips), schema/management and
preserved-history/document checks. The clean signed package passes install/apply,
zero-change repeat and actual read-only XFMD tree checks; see
[evidence](evidence/candidate-package.json). Production signing uses the existing
publisher key. PR42 passed all four CI jobs and merged as ae08e81.

The legacy expected-authority refresh was independently reviewed by rp2_review:
15 plans contain only 638 scalar version substitutions, with unchanged actions,
ordering, safety and failure outcomes. Same-version fixture and prerelease tests
retain their meaning. No blocker/high/medium findings; approval is conditional on
the ordinary non-regenerating 19-scenario replay passing before publication.
The replay and release-preparation CI remain pending; no tag/publication claimed.
Final merged-source signing and downloaded-package checks are still required.
