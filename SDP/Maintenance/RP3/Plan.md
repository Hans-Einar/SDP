# RP3 — Release Go-only SDP 1.0.0 and upgrade XFMD

Current authority: the owner's Go-only correction and explicit 1.0.0 selection
below supersede earlier 0.2.2 and legacy-runtime release-gate text. No 0.2.2
release was published. Only the revised Go candidate may proceed.

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0010 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

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
release. The requested version decision is pending because removal retires former
source-distributed entry points; do not reuse the old candidate evidence as final.

## Current release selection

Owner explicitly selected SDP 1.0.0 following the major-version rule for removal
of the former public source-distributed installer entry points. No 0.2.2 release
was published. Go command/receipt compatibility and supported signed predecessors
0.2.0/0.2.1 are retained. gh-sdp 0.1.2 updates its immutable default to 1.0.0.
Only Go-based release gates apply to this revised candidate; previous runtime
results describe the superseded candidate and do not require rerunning retired
engines. All .ps1/.psm1/.psd1 files and executable invocation paths are absent.
