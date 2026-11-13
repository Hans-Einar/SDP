# Need-driven Activities across successive Sessions

| Field | Value |
| --- | --- |
| id | KB-SDP-061 |
| project | SDP |
| type | Proposal |
| CardState | completed |
| created | 2026-11-11 |
| PlanId | PLAN-SDP-0027 |
| Systems | ProjectGovernance, SDPTOOL |
| source | Owner proposes Activity as the link between a customer need and successive Sessions |

[Study](../../02--Requirements/ProjectGovernance/Activities/Study.md).
[Session0015](../../Sessions/session-%230015--SDL_browsing_planning.md).

Recommend a need-level Activity linking Sessions, plans, cards and model/Blueprint
evidence, with one selected active Activity in an explicitly scoped authority.
Keep plan milestones, card state and PG runs authoritative in their existing homes.
Distinguish management Activity from SDL activity. Reconcile KB042 and PG
KB038/KB052 before schema/runtime adoption; preserve existing history/work states.

2026-11-11 owner selects Activity integration before KB060 browsing. Current scope
is complete design, tracked by [Session0016](../../Sessions/session-%230016--Activity_design.md).
[Design candidate](../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md)
covers lifecycle/history, six-phase Activity folders, producer-owned Activities
tab replacing Sessions, event timelines, PG ownership and migration. No runtime
schema adoption or implementation is claimed. Next: model/protocol reconciliation
and an implementation plan; qualify Activity integration before resuming KB060.

2026-11-12 owner refinements (Session0016 T2): generic Activity README, need/outcome
in Mandate.md, Activity-root Sessions, shared ledger confirmed, free unrelated
discussion/capture with Activity-scoped implementation, and attributed commit/ledger
markers over lifecycle intervals. Design and acceptance updated; no runtime delivery.

2026-11-12 Session0016 T3: owner selects SDL design for tool-created Activities and
Sessions. [Commands](../../04--Design/SDPTool/ProjectGovernance/Activities/Commands.md)
maps existing source boundaries and proposed operations. The new ActivityWorkflow
ModelGovernance WORK validates canonically; checkpoint remains preliminary.
Session creation does not dispatch a PG run. [Evidence](../../04--Design/SDPTool/ProjectGovernance/Activities/Evidence-T3.md).

2026-11-12 owner sequencing decision: hold Activity feature work while
ACTIVITY-SDP-0001 restructures SDPTool without new functionality. Resume by manually
creating the second Activity after restructuring acceptance. Preserve all existing
study/design/model evidence. Review trigger: PLAN-SDP-0026 M5 acceptance or an
explicit owner priority change; no arbitrary restart date is inferred.

Owner accepts Activity0001 and selects ACTIVITY-SDP-0002 for implementation.
[PLAN0027](../../Activities/Activity-0002--Activity-Management/05-Implementation/Plan.md)
owns delivery; Session0018 coordinates execution. The prior hold is resolved by
actual predecessor acceptance. Activity0003/browsing remains outside this scope.

2026-11-12 Activity0002 technical delivery: producer eb11ac7 and isolated XFMD
consumer a558c31 implement Activity/Session creation, guarded lifecycle/recovery,
producer-owned Activities pages/timeline and explicit governance admission.
[Qualification](../../Activities/Activity-0002--Activity-Management/06-Verification/Verification.md)
records full regression, final delta race, native workflows and independent review.
CardState gate-review; actual owner acceptance remains pending. No Activity0003,
live admission, installation, merge or release is inferred.

## Delivered outcome

Completed under the owner's autonomous implementation/publication authority.
Signed public XFMD0.10.0-76 delivers the selected capability with exact private
producer c3463b7. Independent source/package/public verification and actual signed
DNF upgrade/fresh-install native workflows pass. Activity3/06-Verification holds
evidence. Earlier queue/gate-review/no-adoption statements are historical. Personal
owner review was not performed; no main merge or Toolkit installer release claimed.
