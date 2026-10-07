# SDUI widget capability and interaction design

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0021 |
| project | SDP |
| state | completed |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDUI; SDL |
| source | KB-SDUI-003; Session0010 T001 |

## Outcome and authority

The owner asks to take the work described in KB-SDUI-003. This selects its full
2026-10-07 inventory and staged execution, including the prerequisite design.
Earlier statements that only card editing was authorized describe the earlier
request; they do not limit this new assignment. No publication or merge is selected.

Reuse PLAN-SDP-0009 research, the existing runtime and SDL bridge, and KB-SDUI-005
producer work. Deliver a versioned capability/interaction contract, per-family
acceptance and an executable staged handoff. Design completion does not complete
the widget card. The matching implementation plan is PLAN-SDP-0022.

## Git policy and baseline

Use current working branch `sdp/model-governance-implementation` for this design
milestone. HEAD at intake: `3d265d1`. Preserve the pre-existing dirty governance,
SDUI preview, sourceinput and external XFMD producer work. Do not switch the shared
worktree. Commit scoped design/management changes per milestone after validation;
do not stage unrelated hunks in shared ledgers/indexes. Implementation isolation
must retain the current preview additions as explicit candidate inputs.

## Milestones

| Milestone | Acceptance | State |
| --- | --- | --- |
| WCD1-M1 | Resolve current behavior, ownership, profiles, activation and typed identities | completed |
| WCD1-M2 | Every family and SDUI-001–010 gap has acceptance and an implementation destination | completed |
| WCD1-M3 | Independent design review, checked records and bounded implementation handoff | completed |

## Evidence and remaining work

[Evidence](Evidence.md) records independent review, corrected findings and limits.
WCD1 completes design readiness only; stage-specific elaboration remains governed
by the implementation plan's review gate.

See [Design](Design.md), [acceptance matrix](Acceptance.md) and
[implementation plan](../../../05--Implementation/SDUI/Widgets/Plan.md).
Native interaction, SDL-connected fixtures and consumer packaging remain
implementation obligations. Existing static pictures and successful parsing
cannot satisfy them. Review must distinguish proposed rules from implemented ones.
