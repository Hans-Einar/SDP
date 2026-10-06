# Semantic blueprints and assignment context

| Field | Value |
| --- | --- |
| id | KB-SDP-050 |
| project | SDP |
| type | Proposal |
| CardState | in-progress |
| Systems | SDPTOOL, SDL, SDUI |
| created | 2026-10-03 |
| source | Owner discussion; KB-SDP-048 split KBO-SDP-000005 |

## Source and ownership

[KB048](../superseded/%23048--Proposal--Versioned-design-reviews-and-blueprint-diffs.md)
provides historical decisions; split KBO-SDP-000005 transfers this bounded scope.
The sibling card owns the complementary scope; no feature is declared implemented
by the split.

## Outcome and scope

Generate semantic before/after differences, impacted system context, preserved
contracts and assignment boundaries from SDL/SDUI models. Distinguish preliminary
WORK views from retained reviewed targets; connect code evidence without claiming
that tags alone prove behavior. Keep analysis independent of storage backend.

[BP1 study](../../04--Design/SDPTool/Blueprints/Study.md) and
[BP2 / PLAN-SDP-0001](../../04--Design/SDPTool/Blueprints/Plan.md) remain the study
and planned design successor. KB004 owns general model/code/evidence mapping.
KB049 owns ModelGovernance source views/identities; it does not implement blueprints.

2026-10-03: Scope separated from KB048. No new Session or implementation activated
for this card. BP2 remains planned. Session0005 discussion closes by explicit
handoff; blueprint planning will be selected separately.

## BP2 activation — 2026-10-06

Owner selects continuation. [Session0008](../../Sessions/session-%230008--Semantic_blueprints.md)
executes the existing [BP2 DesignPlan](../../04--Design/SDPTool/Blueprints/Plan.md).
The [contract candidate](../../04--Design/SDPTool/Blueprints/Contract.md) uses a
reduced MVP1 calibration-ownership pilot. Real NOW/TARGET parsing and model-governance
capture/promotion pass; this is neither a full experimental-profile migration nor
an implemented blueprint producer. The specific pilot awaits owner feedback.
Next BP2-B: impact selection and explicit unknown/negative cases. BP2-C and a
production ImplementationPlan remain open.

2026-10-07 BP2-B-M1: Selection-and-Evidence.md and eight passing experimental
selection tests delivered. Actual parser acceptance of a missing channel endpoint
is recorded as a separate blueprint completeness gap. NOW/TARGET union, contract
closure, scope/frontier and code/evidence-reference rules are specified. Next BP2-C;
card remains in-progress, with no production generator or owner-approved pilot claim.
