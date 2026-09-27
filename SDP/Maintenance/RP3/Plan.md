# RP3 — Release canonical SDP templates and upgrade XFMD

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

SDP 0.2.2 corrects template organization and source-authoring guidance without
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
