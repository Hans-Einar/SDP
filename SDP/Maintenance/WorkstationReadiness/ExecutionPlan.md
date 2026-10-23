# WR2 — Consolidated candidate verification and PG MCP adoption study

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0018 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDL, SDUI, ProjectGovernance |
| source | Owner Session0013 T002; KB-SDP-056; successor to MAINT-SDP-0017 |

## Outcome and authority

Owner T002 authorizes executing the proposed Session route with one replacement:
C4 becomes a study of what is required to start using ProjectGovernance through
MCP. Defer the MVP1 pilot; the owner confirms a pilot model already exists in
Ponsse. Do not replace it or make another pilot model a prerequisite. The broader
SDP system is the readiness concern. Continue on the existing working branch.

## Scope and milestones

| Milestone | Observable acceptance | Current result |
| --- | --- | --- |
| WR2-M1 / C1 | Pin integrated source; provision compatible isolated toolchain if needed; build current CLI/controller/MCP; verify selected executable and relevant Go/race/vet suites | delivered; private Go 1.27.1/build and four race suites/three vet checks pass; current controller host refusal retained |
| WR2-M2 / C2 | Reconcile current entrypoints and model/Blueprint/PG identity-state responsibilities; correct demonstrated documentation/contract mismatches without speculative refactoring | delivered; see result register |
| WR2-M3 / C3 | Exercise integrated provider/assignment/MCP workflow and negative/recovery checks against compiled candidate; distinguish fixtures, real SDK transport and host/native execution; preserve required independent-review boundary | delivered; see result register |
| WR2-M4 / revised C4 | Deliver evidence-backed MCP adoption study: usable modes/tools, deployment/configuration, authority, current limitations and smallest implementation sequence for ordinary use | delivered; see result register |
| WR2-M5 / C5 | Reconcile Session/cards/plan/evidence; state actual verification gaps and next bounded PG work | in-progress; final record/link checks and closeout |

## Preserved boundaries

No Ponsse application edits, MVP1 trial, new model, main merge, release or stable
gh-sdp replacement. A private pinned development build is sufficient for C1.
Global MCP registration is not required for the adoption study; prepare a concrete
configuration and test transport locally. Existing PG, Blueprint and XFMD scopes
and owner pauses remain; completion here does not close PLAN-SDP-0018. Do not
implement the full unfinished PG leader merely to make the study look ready.

Perform directly under current host rules; no sub-agents requested or spawned.
Same-context checks are not independent review. Retain that unverified boundary
if independent review cannot be supplied within authorized host/session rules.

## Concurrent-work preflight — 2026-10-23 T002

HEAD f8110e9 on sdp/workstation-completion-assessment, clean; main bfe3b40 is an
ancestor and only current documentation commits differ. One registered SDP worktree;
owner confirms full laptop handover. PG Session0007 T028 and XFMD Session0007 T012
are latest owning records; active states describe retained scope, not observed
live agents. No competing dirty/staged/untracked product edits observed.
Green for this Maintenance directory, Session0013, KB055/056/index and append-only
management registration; scope also includes current PG/MCP documentation and
isolated build/verification artifacts outside source. Check new identities across
canonical ledger/current records before allocation. Product source fixes, if
needed, require a refreshed exact write-set verdict and existing contract evidence.
No new product architecture or state ownership is selected by this plan.

## Git and evidence

BranchPolicy current; use sdp/workstation-completion-assessment with milestone
commits and separate product fix evidence if a real defect is found. Keep build
caches/toolchains/lab state outside Git. Capture command results, tool versions,
source/tree and binary digests, and fixture limitations. Preserve ledger bytes.
Update governing current documentation with behavioral findings, not just evidence
reports. No release/installability claim from a local development build.

## Result register

Updated as milestones complete. The Session roadmap is the cross-work projection;
this plan owns these bounded milestones. Do not duplicate PG implementation phases.

## WR2-M1–M4 delivered results — T002

[Execution evidence](Execution-Evidence.md) and retained logs identify exact binaries,
commands and limitations. All four Go module race suites and three vet checks pass.
The provider/assignment trial, real unbound and bound-worker MCP transport, and
actual Codex 0.164.0 read calls pass. The pinned PG controller refuses 0.164.0 as
expected; no native task or independent review is claimed. Current responsibilities
and stale MCP/action documentation are reconciled without product code changes.
[MCP adoption study](MCP-Adoption-Study.md) and a validated concrete inspection
configuration are delivered. WR2-M5 closeout is next; broader PG readiness remains
open under its existing plan/cards.
