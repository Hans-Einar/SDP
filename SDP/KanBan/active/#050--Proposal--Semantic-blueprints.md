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
contracts and assignment boundaries from SDL/SDUI models. Include a retained
blueprint catalogue discoverable by SDPTool, with revision-bound assignments and
evidence-backed work-state grouping for an XFMD Blueprint sub-tab. Distinguish preliminary
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


2026-10-07 BP2-C-M1: Producer-and-Handoff.md and independent Trial-review.md
delivered. A meets the authored assignment; deliberately violating B is rejected
despite parser success. [PLAN-SDP-0020](../../05--Implementation/SDPTool/Blueprints/Plan.md)
is planned for Go analysis, publication and generated-bundle verification.
Next BPI1; BP2-A owner pilot disposition remains pending. No production generator
or implemented Ponsse change is claimed.


2026-10-07 Session0008 T004: owner requests an XFMD Blueprint sub-tab with
assignment/implementation progress. Catalogue requirements and proposed lifecycle
are recorded in Producer-and-Handoff.md. Extend BPI2/BPI3 design before execution;
no duplicate manual navigation registry or inference of implementation from
parent-card status. Native XFMD integration remains separately owned.


## Catalogue delivery — selected 2026-10-07

Owner authorizes extending this existing card and plan rather than creating a
parallel feature card. PLAN-SDP-0020 owns the detailed milestone acceptance:
BPI2-M2 delivers retained revisions and discovery; BPI3-M2 delivers assignment
lifecycle and status projection. A manual navigation registry is excluded.
Blueprint freshness/readiness and assignment progress remain separate.
The producer must support a headless consumer; native XFMD changes remain with
the XFMD agent and are not claimed by this card's producer acceptance.
Session0008 T005 records the plan revision; production execution has not started.
