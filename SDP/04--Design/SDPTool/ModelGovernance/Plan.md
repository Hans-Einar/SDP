# MG — ModelGovernance design and delivery preparation

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0016 |
| project | SDP |
| state | active |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDL, SDUI |
| source | KB-SDP-049; owner instruction 2026-10-03 |

## Outcome and authority

Owner selects a separate ModelGovernance workstream: create study/design and a new
Session, activate its card and begin planning. Deliver a bounded implementable
contract and subsequent ImplementationPlan. This turn establishes the study and
initial design; production implementation is not claimed by document creation.

## Scope and baseline

[Study](Study.md) consolidates decisions. [Design](Design.md) defines responsibility,
WORK-local recovery, immutable promotion and retained lineage. Semantic blueprints
remain KB050/BP2. No XFMD changes, global VCS, distributed locking or mandatory Git.
Session0006 owns the roadmap; this plan owns detailed design milestones.

## Git policy

Use current sdp/blueprint-model-history-proposal branch; commit each milestone with
its ID and evidence. Do not rewrite its historical name/commits merely to rename the
feature. No main merge or release authorized by this preparation. Preserve unrelated
untracked files. Existing general phase-push policy remains separate from publication.

## Phases and milestones

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| MG1 Study | MG1-M1 | Consolidated decisions, explicit scope split, current capability limits, active card/session and initial design | completed |
| MG2 Contract | MG2-M1 | Resolve minimal YAML/history schema, command grammar, naming, recovery and promotion policy; model supported behavior in SDL | next |
| MG3 Proof | MG3-M1 | Disposable filesystem experiment validates reconstruction, merge and metadata-only promotion including negative cases | planned |
| MG4 Handoff | MG4-M1 | ImplementationPlan defines vertical slices, tests, migration/discovery impact and review boundary | planned |

## Evidence and remaining work

MG1-M1: Study/Design created; KB048 scope split into KB049 ModelGovernance and KB050
Blueprints; Session0005 handed off, Session0006 active. Management/Toolkit validators
and diff checks are run for this delivery. These are document consistency checks,
not proof of the proposed storage/merge design. MG2–MG4 remain undelivered; this
DesignPlan stays active. No automatic owner acceptance of design recommendations.
