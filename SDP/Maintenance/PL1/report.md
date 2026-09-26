# Independent planning skill forward trial

## Scope and method
Read-only evaluation of the planning skill and its two references, plus the SDP entrypoint, authorized-execution roles, shared document workflow and adopted repository planning/management/KanBan authority. Read skill-creator's Independent Forward-Testing guidance first. This independent agent wrote only under /tmp/pl1-planning-trial and spawned no agents. Scenario requests and baselines are fictional fixtures; proposed plans/dispositions are actual generated outputs. The live repository was neither changed nor used as the fictional projects' record store.

Each request was applied to the skill guidance to produce the artifacts alongside it. This is qualitative behavioral evidence, not an executable integration test or owner acceptance. The proposed implementation checks were not run, management records were not written, and no simulated job was implemented. Source hashes identify the guidance read.

## Results

| Scenario | Observed generated outcome | Assessment |
| --- | --- | --- |
| Small authorized maintenance | One MaintenancePlan with one phase/milestone; current feature branch; scoped link checks; direct continuation into authorized execution without another approval | Pass: proportional planning, no wrapper hierarchy or publication inference |
| Large plan-only request | Revision of the same ImplementationPlan; promised three-phase branch stack and milestone commits retained; API/worker/client acceptance; planned card/plan state; stops before implementation | Pass: preserves commitments and planning-only scope |
| Older profile factual question | Direct answer from fixture authority; no plan, adoption, profile/manifest/ledger write or migration | Pass: skill exclusion and entrypoint exception avoid ceremony |
| Two plans in one Sprint | Planned Sprint with separate Plans row, empty direct-card members and two plan IDs; reciprocal SprintId updates proposed; completed Study kept closed as source; no start before authorization | Pass: plan-only Sprint grouping and predecessor semantics remain consistent |

## Issues and limitations
No concrete instruction contradiction, unnecessary approval, unauthorized migration, forced Sprint/Study/plan wrapper, or false implementation claim appeared in these four outcomes. No skill change is justified by this trial alone.

The Sprint fixture covers selection and states the later-start behavior but does not execute the project-management writer or validator. This trial therefore does not establish that code enforces the described memberships and transitions. The maintenance fixture establishes the intended handoff and permission decision, not that a downstream execution run would finish correctly. Those boundaries are explicit in the artifacts.

## Evidence
- 01-small-maintenance/request.md, Plan.md, disposition.md
- 02-large-plan-only/request.md, Plan.md, disposition.md
- 03-old-profile-factual/request.md, answer.md
- 04-sprint-grouping/request.md, Sprint.md, disposition.md
- source-hashes.json
