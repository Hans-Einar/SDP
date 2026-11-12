# Complete SDPTool architecture migration

| Field | Value |
| --- | --- |
| id | KB-SDP-062 |
| project | SDP |
| type | Proposal |
| CardState | gate-review |
| created | 2026-11-12 |
| Systems | SDPTOOL |
| PlanId | PLAN-SDP-0026 |
| source | Owner challenges open-ended cleanup and loss of architectural intent between analysis and implementation |

Need: complete one coherent SDPTool migration without leaving transitional code as
the accidental architecture for later agents. [Proposed ImplementationPlan](../../05--Implementation/SDPTool/ArchitectureMigration/Plan.md)
sets whole-system inventory, SDL target, bounded Blueprint use, executable
architecture checks, per-workflow deletion and final integrated acceptance.

2026-11-12 Session0016 T5: proposal prepared using current SDP Planning and Vertical
Refactor methods. Source inspection confirms Blueprint supports pinned inputs and
evidence but cannot itself prove Go conformance. Execution is not selected. KB061
Activity design stays active and KB060 browsing stays queued. Adoption of the
proposed combined migration/Activity delivery remains a scope decision.

2026-11-12 Session0016 T6: owner selects ACTIVITY-SDP-0001 for restructuring
only. Activity implementation is removed from PLAN-SDP-0026 and reserved for the
next manually created Activity; browsing follows tool support. Card is ready for
M0; the plan remains planned. Manual registration is recorded in Session0017,
not represented as product implementation.

Session0017 T002: owner authorizes execution of Activity0001. PLAN0026 is active;
M0 baseline/inventory begins on sdp/sdptool-architecture-migration.

Session0017 T002 technical delivery: candidate 527cc9008bdf8f567f8b201b665b2960865518f1
completes the selected restructuring. Typed governance operations, direct Session
adapters, registered action policy, shared MCP framing, single CLI routing and
typed program execution are implemented; no migration-only adapter remains.
All 60 operation rows and external boundaries are reconciled. Exact-candidate full
suite, focused race/vet, 187 CLI comparisons, nested/bootstrap checks and all six
native gates pass. Independent product review approved with no material findings.
[Qualification](../../Activities/Activity-0001--SDPTool-Restructuring/06-Verification/Verification.md)
contains reproducible evidence and limitations. CardState gate-review requests the
owner disposition on this concrete result. Plan/Activity stay active until that
disposition; no new Activity functionality, merge or release is claimed.
