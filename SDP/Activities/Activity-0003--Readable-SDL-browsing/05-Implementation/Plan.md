# PLAN-SDP-0028 — Readable SDL browsing and RPM delivery

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0028 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Activity | ACTIVITY-SDP-0003 |
| Source | Owner autonomous completion and RPM publication request; KB-SDP-060 |

## Outcome and boundaries

Deliver all-viewpoint readable browsing described in NavigationDesign through
SDL/SDPTool, generic XFMD consumption and signed XFMD0.10 RPM publication.
[ImplementationContract](../04-Design/ImplementationContract.md) binds ownership.
PLAN0025 is the prerequisite design route, reconciled within B0. No hand-authored
viewpoints, new SDL port grammar, unrelated refactoring or main merge.

| Milestone | Observable acceptance | State |
| --- | --- | --- |
| B0 | Tool-created Activity/Session; requirements, WORK model and handoff pinned | completed |
| B1 | System/include/fragment navigation; ordered unique actionable VP collections | implemented |
| B2 | All-VP explanatory pages, complete chapter collections, separate registers and readable exact source links | implemented |
| B3 | Faithful container/unit/channel reading and generic native consumption; HEOS and edge fixtures | implemented |
| B4 | Integrated regression, actual GUI/print checks and independent review | in-progress |
| B5 | Exact committed0.10 package with matching producer, signed public repo and DNF verification | planned |

Use SDP branch sdp/readable-browsing and XFMD isolated branch of the same name.
Milestone commits/pushes preserve the Activity2 history. Package publication is
explicitly authorized; existing main and primary XFMD dirty files remain intact.

## Verification

All VP collections and isolated declarations; HEOS all twelve fragment contexts;
orphan/broken includes; stale siblings, source spans/Unicode; no machine IDs;
large selections without omissions; channel mode/role/topology correctness; native
expansion versus activation and Back/refresh; print readability. Reuse unchanged
Activity2 evidence and run integration regressions for modified boundaries.
Consolidate independent review before signing/publication. Actual delivered
artifacts and limitations belong in 06-Verification, never preclaimed here.
