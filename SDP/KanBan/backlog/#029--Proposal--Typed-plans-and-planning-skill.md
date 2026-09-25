# Typed plans and a shared Planning skill

| Field | Value |
| --- | --- |
| id | KB-SDP-029 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-09-25T15:19:14Z |
| source | Owner conversation 2026-09-25 during MAINT-SDP-0003 |
| next_review | After MAINT-SDP-0003 |
| tags | process, planning, skills |

## Owner direction

Make Plan an explicit concept with MaintenancePlan, RequirementPlan,
ArchitecturePlan, DesignPlan, ImplementationPlan and VerificationPlan.
The normal outcome of a KanBan card is a plan; put the card in active when
executing that plan. A Sprint can contain one or several plans. Plans have
phases and milestones, scaled to the actual size of the work.

Each plan declares its Git policy. One branch per phase and one commit per
milestone mainly suits large implementation plans. Smaller plans normally stay
on the active working branch, with commits per phase or milestone as appropriate.
Do not make elaborate branching a universal planning requirement.

Add a Planning skill that covers the plan framework together with KanBan,
Scrum and Sprint selection/grouping. Reuse the existing shared process authority
and avoid a second lifecycle or duplicating plan text in cards and Sprint records.

## Scope and acceptance

Select a minimal common plan contract with typed purpose, owner outcome,
authority, scope, phase/milestone acceptance, dependencies, evidence, lifecycle,
Git policy and references. Define where plans live and which current documents
already satisfy the contract; avoid mass renaming merely for type labels.

Reconcile planning versus execution CardState with existing Study completion,
direct card-to-Maintenance and Scrum-to-Maintenance routes. Retain lightweight
routes and historical meanings; define whether planning work needs a separate
state without inventing an owner decision. Update skills, neutral templates,
validators and installed-profile compatibility together where necessary.

Use the skill-creator procedure when actually implementing the Planning skill.
Test routing and ordinary small/large plan examples. Update AGENTS Git guidance
so an explicit plan policy governs future work while preserving historical
branch commitments.

## Worklog

Registered during IU2. MAINT-SDP-0003 continues using its already selected
stacked phase branches and milestone commits. This proposal does not restart
that maintenance or retrospectively change existing record states.
