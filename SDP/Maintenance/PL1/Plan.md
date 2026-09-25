# MAINT-SDP-0004 — typed plans and Planning skill

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0004 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | phase |
| source | KB-SDP-029; owner authorization 2026-09-25 to plan and execute |
| Systems | SDP |

## Outcome and scope

Make six typed plans usable without adding duplicate documents or mandatory
Scrum/Sprint layers. Coordinate cards, plans and sprints through one management
history. Add a discoverable, distributed sdp-planning skill and test small/large
planning decisions and installed compatibility. Maintain English documentation.

Baseline: 28bf156, completed MAINT-SDP-0003, management profile 0.1, thirteen
canonical skills. Preserve all historical ledger bytes/IDs and the unrelated
untracked SDL/go/sourceinput draft. No live consumer rollout, language/runtime
work or historical mass migration. KB014/018 retain their independent scopes.

## Selected design

PlanType distinguishes MaintenancePlan, RequirementPlan, ArchitecturePlan,
DesignPlan, ImplementationPlan and VerificationPlan. Existing Maintenance
records can carry MaintenancePlan without another ID or wrapper. Other new
plans use PLAN-<PROJECT>-<number> and management kind Plan. A plan is the work
record across planned/active/completed/canceled states, not a second lifecycle.

Keep planning cards in backlog/queued until their selected plan executes; active
then means ready/in-progress/gate-review. A Study whose bounded outcome is a
plan may still complete with an explicit successor, preserving lightweight and
historical routes. No new CardState or mandatory gate. A sprint may list plans
as well as legacy/direct card members, without copying their text.

Introduce explicit management payload/profile 0.2 for typed plans and plan
membership; retain 0.1 history interpretation. The five-phase layout identity
stays 0.1. Upgrade readers and installed facts before emitting profile 0.2.
Project-owned old documents remain preserved; managed planning guidance and
versioned neutral templates provide the new authority without overwriting them.

## Git policy and execution authority

One working branch sdp/maintenance-pl1-typed-plans, created from IU3's completed
branch to preserve that delivery. One commit per phase; no extra milestone
branches. Push each completed phase and open a combined PR against sdp-vNow;
no merge. This policy applies the owner's KB029 direction to this small plan.
Phases still have milestone acceptance; record evidence at each commit.

## Phases and milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| PL1-A | PL1-A-M1 | Adopt one plan contract, template and scaled Git/lifecycle guidance; activate KB029 and preserve history | Completed |
| PL1-B | PL1-B-M1 | Implement typed-plan/sprint validation, portable Planning skill, profile/readers and safe distribution | Planned |
| PL1-C | PL1-C-M1 | Verify behavior, compatibility, installed workflow and records; resolve review findings and close work | Planned |

## Verification and completion

Use old/new lifecycle fixtures and negative cases for types, membership, state,
identity and Git-policy metadata. Verify skill metadata/references/discovery and
realistic small/large/no-authorization planning behavior. Test an old-profile
installation upgraded with preservation, truthful facts, no-change repetition,
and SDPTool recognition. Rebuild the artifact reproducibly. Check management,
links, unchanged historical prefixes and the known 38 Toolkit baseline findings.
Review actual final changes, distinguishing independent skill trials from local
checks. Evidence belongs in Evidence.md. Completion requires the usable contract,
skill and distribution, not a live project upgrade or release publication.
