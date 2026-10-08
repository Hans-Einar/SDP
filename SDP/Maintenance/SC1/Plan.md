# SC1 — Session continuity correction

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0014 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

## Outcome and authority

Owner instruction 2026-10-02: immediately require Session updates in agent and
project-governance instructions after missing journal entries during KB048 work.
Use existing sdp governance entrypoint; no project_governance skill exists in the
current collection. Recover Session0005 honestly, without invented transcripts.

## Scope and Git policy

Edit AGENTS, shared sdp skill, its exact install inventory version and Session guide;
register/link Session0005 and record owner blueprint corrections. Current branch
sdp/blueprint-model-history-proposal; commit SC1-M1. No backend implementation,
release, installed-project mutation or unrelated untracked files.

## Phase SC1

| Milestone | State | Acceptance / evidence |
| --- | --- | --- |
| SC1-M1 | completed | Per-turn discussion upkeep explicit; retrospective journal labeled; owner corrections recorded; skill, management and Toolkit validators pass |

## Verification and limits

Skill quick_validate, ProjectManagement validate, Toolkit validate_sdp and
git diff --check pass. Initial Toolkit validation caught the additional canonical
SDP.manifest.yaml and installed-manifest example skill-version references; these
were aligned with the install inventory
and validation repeated. These checks establish structure/consistency,
not guaranteed future agent behavior. This maintenance delivers local instructions,
not runtime enforcement or automatic conversation capture. Session0005 continues.
