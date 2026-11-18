# Standard DocumentHost prototype launch

| Field | Value |
| --- | --- |
| id | KB-SDP-065 |
| project | SDP |
| type | Change |
| CardState | completed |
| Systems | SDUI, SDPTOOL |
| created | 2026-11-17 |
| PlanId | PLAN-SDP-0030 |
| source | Owner authorizes extending standard launch after TimberFrame reproduction |

Use the existing DocumentHost with explicit choice/resource context for ordinary
SDPTool launch. Preserve legacy 0.2 and truthful readiness. [Plan](../../05--Implementation/SDUI/DocumentLaunch/Plan.md).

## Worklog

2026-11-17: owner selects implementation following completed KB-SDUI-007.

2026-11-17: DL1/DL2 delivered at 2db8ab1/735c865. Explicit context, DocumentHost
launch and guarded reload verified using actual packaged TimberFrame Analysis.
[Evidence](../../05--Implementation/SDUI/DocumentLaunch/Verification.md) and
[independent review](../../05--Implementation/SDUI/DocumentLaunch/IndependentReview.md).
LaunchWorkflow model 0.1.0 released. No remaining selected implementation work;
binary publication and XFMD context-selection UI are separate scope.

2026-11-18: owner-authorized publication completed under
[MAINT-SDP-0023](../../Maintenance/R240/Plan.md). Signed SDPTool/gh-sdp2.4.0 public
assets and isolated ordinary extension upgrade verified; no remaining release work.
