# RP2 — Publish the external-reference navigation correction

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0008 |
| project | SDP |
| state | completed |
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

### Release predecessor correction

The first signed 0.2.0 → 0.2.1 rehearsal correctly rejected the candidate because
the descriptor builder always emitted an empty upgradesFrom list. RP2 therefore
adds the exact published 0.2.0 descriptor digest to the authored Go profile and
passes that list through the existing descriptor validator/signature. No implicit
SemVer upgrade permission or installer-policy relaxation is introduced. The signed
upgrade rehearsal must pass before release. An existing process integration test
also now derives its current version from artifact facts before exercising the
99.0.0 downgrade rejection; its previous 0.2.0 literal stopped mutating the fixture.

### RP2-M2 — upstream published

PR43 passed all four CI jobs: contracts 4m31s, Go installation 1m41s, Windows
installer 11m52s and Linux process/recovery 6m47s. All local review conditions passed;
[final independent review](evidence/review.md) approves preparation 60dbbe8 with no
unresolved blocker/high/medium findings. Ordinary conformance passes 19 scenarios.
The corrected real downgrade fixture passes.

PR43 merged to 4bacfce05f92f0dab9680456e297214727b54bc6. Its exact clean package was
production-signed, installed, repeated and used for read-only XFMD navigation.
Signed 0.2.0 → 0.2.1 upgrade preserves owner content and history; unchanged payload
requires zero file actions while the signed receipt records 0.2.1. The extracted
GitHub source archive installs successfully and records null sourceCommit.

The annotated v0.2.1 tag and final GitHub Release now exist:
https://github.com/Hans-Einar/SDP/releases/tag/v0.2.1

Descriptor SHA256: 66590e8e967ede6b36d8fa45cdbee1cd80f69505cca698b0e4a9bd960842735a.
[Publication](evidence/publication.json), [final package](evidence/final-package.json),
[upgrade](evidence/final-upgrade.json) and [archive](evidence/final-archive-install.txt)
record observed outcomes. Prior 0.2.0 tag/assets/released notes remain unchanged.
gh-sdp 0.1.1 is published at
https://github.com/Hans-Einar/gh-sdp/releases/tag/v0.1.1
from 8cbef9693e353cb04dc6d98021aca454b9078e3f, with its own independent review.
The installed extension was upgraded from 0.1.0 to 0.1.1. With no release or test-key
override, it resolves SDPTool 0.2.1 at the exact released source revision and returns
all 18 XFMD cards with validated KanBan and externalReference KB-SDP-014.
[Installed identity](evidence/installed-version.json) and
[actual tree](evidence/installed-xfmd-tree.json) record the completed user workflow.
XFMD remains clean and its project process files are unchanged. The other gh-tree
extension remains installed. No native GUI badge or wider platform claim is made.

Downloaded upstream assets pass SHA256SUMS and production signature/install checks;
see [download verification](evidence/downloaded-package.json). All selected RP2
milestones are complete. Publication reconciliation changes only current records
and their derived profile artifact; the published tagged artifacts remain immutable.
