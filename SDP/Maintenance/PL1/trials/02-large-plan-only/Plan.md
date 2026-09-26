# PLAN-FILES-0007 — Durable upload recovery

| Field | Value |
| --- | --- |
| id | PLAN-FILES-0007 |
| project | FILES |
| state | planned |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| source | Owner request; KB-FILES-012; approved upload recovery design |
| Systems | FILES |

## Outcome and authority
Revise the existing authoritative plan so an interrupted upload resumes across desktop, API and worker. Only planning is authorized in this request. The approved design and existing branch commitments remain authoritative.

## Scope and baseline
API owns durable upload identity and account authorization; worker owns idempotent completion; desktop owns persisted identity and restart/resume interaction. Preserve ordinary successful upload behavior. No new container, storage redesign, publication or unrelated upload feature is included.

## Git policy
Keep the promised stack: upload/p1-contracts from main, upload/p2-recovery from P1 delivery, upload/p3-client from P2 delivery. Commit each milestone with its ID and concrete delivery. This planning revision is prepared on feature/upload-planning; future implementation branches are not created merely because they appear here. No push, merge or release permission is inferred.

## Phases and milestones
| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| P1 — Durable API contract | M1 | API persists an upload ID through restart; authorized resume observes its state; another account's ID reveals no state | Planned |
| P2 — Worker recovery | M2 | Retried processing after interruption produces one recorded completion and no duplicate final effect; API exposes that durable result | Planned |
| P3 — Desktop resume | M3 | Start upload, interrupt client/API/worker at selected boundaries, restart and resume the same upload to one successful result; ordinary upload still succeeds | Planned |

## Dependencies and verification
P2 depends on the P1 persisted identity contract; P3 depends on P2 recovery semantics. Use contract and worker integration checks at M1/M2, then the integrated desktop/API/worker recovery workflow at M3. Record exact candidate and actual results when implementation is authorized. None of these proposed checks has passed as part of planning.

## Remaining work
Owner execution authorization and all implementation/verification remain outstanding. The plan revision is the deliverable of this request.
