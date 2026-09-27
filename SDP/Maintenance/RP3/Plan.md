# RP3 — Release Go-only SDP 1.0.0 and upgrade XFMD

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0010 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

Current authority: the owner's Go-only correction and explicit 1.0.0 selection
below supersede earlier 0.2.2 and legacy-runtime release-gate text. No 0.2.2
release was published. Only the revised Go candidate may proceed.

Owner authorization, 2026-09-28: commit, integrate and release SDP plus gh-sdp,
then use the client to upgrade xfmd-sdl-navigation. Git write access is restored.
Use sdp/mvp1-source-ui and milestone commits. PR45 carries the earlier verified
MVP1 pilot and KB035 as well as TS1; preserve unrelated sourceinput/node files.

## Version and scope

The initially selected SDP 0.2.2 candidate corrected template organization and source-authoring guidance without
changing installed paths, receipt protocols or project-owned overwrite policy.
Both published 0.2.0/0.2.1 descriptor digests are declared predecessors. gh-sdp
0.1.2 updates its immutable default via the shared bootstrap module. Linux amd64
is the existing published platform scope. No new parser/runtime capability.

## Milestones

- RP3-M1: release notes/identities, current conformance expectations, Go and Toolkit
  tests, preserved history, exact clean signed package and upgrade rehearsal;
  independent review where required by the conformance/client contracts.
- RP3-M2: integrate approved candidate, package exact merged source, publish tag
  and immutable signed assets; release/install client selecting that descriptor.
- RP3-M3: inspect XFMD receipt/worktree, preview/apply signed upgrade, verify new
  files, preserved local work, no-op repeat and actual navigation. Reconcile true
  publication and installed identities. No XFMD application code changes.

## Evidence boundaries

Production signing requires a clean exact commit. Use an isolated clean checkout
for packaging rather than including unrelated untracked files. Test-only artifacts
never establish production provenance. Keep all old release assets/ledger bytes.
Legacy conformance authority changes are limited to relocated source paths and
current-version scalars; replay and review before publication. Record evidence
under this plan. Existing project prose remains preserved and is not reported as
updated merely because a newer template exists.

## Owner correction — Go-only delivery

2026-09-28: owner explicitly rejects retaining or running the old shell installer.
Release execution paused before merge/publication; previous candidate 0.2.2 is
not published. RP3-GO removes executable legacy installation/build/test paths,
moves current payload authoring into SDPTool, and makes CI use only Go for
installation/recovery. Historical prose, immutable release assets and data contracts
remain historical; no runnable shell copies remain, including the bootstrap archive.
Verify actual Go signed install/upgrade/recovery and preservation before resuming
release. At that pause the version decision was pending because removal retires former
source-distributed entry points. The subsequent owner decision below selects 1.0.0;
the old candidate evidence is not final evidence.

## Current release selection

Owner explicitly selected SDP 1.0.0 following the major-version rule for removal
of the former public source-distributed installer entry points. No 0.2.2 release
was published. Go command/receipt compatibility and supported signed predecessors
0.2.0/0.2.1 are retained. gh-sdp 0.1.2 updates its immutable default to 1.0.0.
Only Go-based release gates apply to this revised candidate; previous runtime
results describe the superseded candidate and do not require rerunning retired
engines. All .ps1/.psm1/.psd1 files and executable invocation paths are absent.

## RP3-M1 revised candidate verified

Go-only candidate 705df7e passes both CI jobs (contracts 11 seconds; Go race,
packaged signed child and descriptor build 1m37s). Independent review is in
evidence/review-go-only.md. Clean production-signed 1.0.0 rehearsal from both
published predecessors passes with four payload actions, preserved owner prose
and repeat no-op; see evidence/go-only-signed-upgrades.json. The minor signing
example now names 1.0.0. Final release artifacts must be rebuilt from merged HEAD.

## RP3-M2 — exact-source publication

Exact-head ff959935 CI passed contracts and Go installation/race/package gates.
PR45 merged to fede327d6f3af35fe7aad323e2c485a27134d20d. Its clean package
was production-signed and rehearsed from both published predecessors: owner prose
preserved, four payload actions and repeat no-op. Downloaded release assets pass
SHA256SUMS. Production bootstrap verifies the signature and exact executable.

SDP 1.0.0 was published at 2026-09-27T23:39:07Z (2026-09-28 locally):
https://github.com/Hans-Einar/SDP/releases/tag/v1.0.0

Descriptor digest: 767527e0d7f54866bab95f4642ffffb2a344024c1805f44b3d5ea9c024829125.
Client 0.1.2 was independently reviewed and published from clean merged
a4988846432c7f7f7b3f722f2786359c9b127372:
https://github.com/Hans-Einar/gh-sdp/releases/tag/v0.1.2

The global extension was upgraded from 0.1.1 to 0.1.2. A fresh-cache installed
`gh sdp --version`, with no release/trust overrides, returns SDP 1.0.0 at the
exact tagged source. Publication reconciliation leaves tagged assets immutable.
No PowerShell command ran after the owner correction.

## RP3-M3 — live XFMD upgrade completed

The existing signed 0.2.0 receipt permits normal upgrade; no custom adoption
manifest is needed. The installed client previewed and applied the root-bound
plan to /home/warloc/git/xfmd-sdl-navigation. Four payload actions update managed
AGENTS.md and Framework/README.md, and initialize SDP/SDL/README.md and AGENTS.md.
The installer wrote signed receipt 1.0.0, MAINT-XFMD-0003 and a ledger event.

All 1229 original tracked/nonignored files were hash-checked. Only the declared
managed files, receipt and appended management ledger changed; the complete
previous ledger bytes remain a prefix. Pre-existing KanBan README/card changes
and application code are preserved. The upgrade changes remain uncommitted in
the XFMD worktree alongside that unrelated local work. Existing project prose
was deliberately preserved, so old dated installation narratives are not rewritten.

A repeated upgrade is noChange=true with zero actions. Actual `gh sdp tree`
returns validated KanBan, including the external primary reference and new local
card. SDL/SDUI remain absent because XFMD has no registered models yet; the
process upgrade does not implement its native sidebar or application design.

Evidence: [publication](evidence/publication.json),
[final package](evidence/final-package.json),
[signed upgrades](evidence/final-signed-upgrades.json),
[installed default](evidence/installed-version.json),
[upgrade summary](evidence/xfmd-upgrade-summary.json),
[apply result](evidence/xfmd-apply.json),
[preservation](evidence/xfmd-preservation.json) and
[actual tree](evidence/installed-xfmd-tree.json).
All selected milestones are delivered; no wider platform or GUI claim is made.
