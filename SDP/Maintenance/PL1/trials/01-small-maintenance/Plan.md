# MAINT-DOCS-0001 — Correct three documentation links

| Field | Value |
| --- | --- |
| id | MAINT-DOCS-0001 |
| project | DOCS |
| state | planned |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| source | Owner request; KB-DOCS-001 |

## Outcome and authority
Correct the three selected relative links. The owner authorized planning and execution, with no publication. This is the sole maintenance work record; no wrapper plan or Sprint is needed.

## Scope and baseline
Apply the exact replacements in the request to docs/install.md, docs/usage.md and docs/faq.md. Preserve surrounding prose and link labels. Replacement destinations exist according to the fixture; verify them against the actual candidate during execution. There are no architecture or dependency changes.

## Git policy
Use feature/docs-link-repair after checking status; preserve unrelated changes. Commit the delivery and its plan/evidence updates at M1, with P1-M1 and the concrete link correction in the commit message. No push, PR, merge or release is authorized.

## Delivery
| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| P1 — Link repair | M1 | All three links resolve relative to their containing files; diff changes only the authorized links and necessary work records | Planned |

## Verification and remaining work
Inspect the three relative destinations and diff; run the installed document/management checks relevant to touched records. No new software tests are needed for these link-only changes. No checks have run in this simulation. Remaining: execute, inspect, record real evidence and commit the delivery.
