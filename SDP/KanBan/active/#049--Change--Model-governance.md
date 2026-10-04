# ModelGovernance — local model lifecycle and recovery

| Field | Value |
| --- | --- |
| id | KB-SDP-049 |
| project | SDP |
| type | Change |
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

Implement a bounded model lifecycle independent of Git: WORK-local checkpoints,
recovery and merge; optional proposals, immutable candidates/releases, retained
messages/lineage after dropping undo payloads. Cover SDL and SDUI source trees.

[Study](../../04--Design/SDPTool/ModelGovernance/Study.md),
[Design](../../04--Design/SDPTool/ModelGovernance/Design.md),
[PLAN-SDP-0016](../../04--Design/SDPTool/ModelGovernance/Plan.md),
[Session0006](../../Sessions/session-%230006--Model_governance.md).

## Work and remaining delivery

2026-10-03 MG1-M1: Owner selects this workstream; study and initial design delivered.
Card active/in-progress for MG2 contract design. No product code delivered. Next:
resolve minimal schema/command/recovery rules, author supported SDL model, then
exercise a bounded filesystem proof and produce the ImplementationPlan.

Completion requires implemented and verified bounded workflow under that successor
plan; finishing the study alone does not close this card. Semantic blueprint impact
selection/diagrams/assignment generation belongs to KB050, not this implementation.

2026-10-03 MG2-M1: [Contract](../../04--Design/SDPTool/ModelGovernance/Contract.md)
and a four-file SDL concern model delivered; valid design-core/0.6 check with no
warnings and generated VP01/VP02/VP08 (9 diagrams). Evidence.md distinguishes
structural validation from runtime proof. Next MG3-M1 filesystem experiment; card
remains in-progress. No ModelGovernance command is implemented.

2026-10-04 MG3-M1: Bounded Go filesystem experiment passes 13 top-level tests with
race instrumentation and vet. [Proof](../../04--Design/SDPTool/ModelGovernance/Proof.md)
records recovery, merge, interrupted publication and promotion limits. No production
CLI delivered; full schema/concurrency/transaction gates remain. Next MG4-M1
ImplementationPlan; CardState stays in-progress.
