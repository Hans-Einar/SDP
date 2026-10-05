# ModelGovernance — local model lifecycle and recovery

| Field | Value |
| --- | --- |
| id | KB-SDP-049 |
| project | SDP |
| type | Change |
| CardState | completed |
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

2026-10-05 MG4-M1: DesignPlan PLAN-SDP-0016 completed;
[ImplementationPlan PLAN-SDP-0019](../../05--Implementation/SDPTool/ModelGovernance/Plan.md)
is prepared, not started. CardState in-progress -> ready for MGI1-M1. Five phases
cover creation, recovery, integration, frozen delivery and consumer/closeout. MG3
proof gaps map to concrete acceptance tests; separate Blueprint scope remains KB050.
No implementation milestone or installed/released functionality is claimed.

2026-10-05T00:16:32.082240+00:00 MGI1-M1: WORK creation/status and strict metadata foundation implemented, with staged area-locked publication. SDL/SDUI validation and recovery library scaffolding included but later CLI operations not yet exposed. Full SDPTool test suite passes.

2026-10-05T00:17:57.226155+00:00 MGI2-M1: Commit/history/whole-state restore and explicit recover resume/abort CLI implemented. Race tests pass, including child-process interruption at prepared/backup/installed boundaries, dirty preservation, corrupt payload and external edit refusal. Backups retained; physical power-loss and non-Linux mutation not claimed.

2026-10-05T00:20:56.033163+00:00 MGI3-M1: Native bounded Go three-way merge, combined WORK creation, persistent conflict inventory and resolved commits implemented. Tests cover disjoint and overlapping same-file edits, unrelated bases, dirty inputs, repeat integration without extra events, three archive generations and whole-state rollback. No external Git merge dependency.

2026-10-05T00:28:20.994585+00:00 MGI4-M1: Frozen candidate/proposal CLI now validates real SDL source graphs and SDUI using existing parsers, preserves metadata lineage and drops undo payloads. Model tests pass. Independent review found dirty-capture identity and file/directory restore defects; regression fixes and stricter domain validation are included, with final re-review pending.

2026-10-05T00:30:41.202655+00:00 MGI4-M2: Release promotion requires candidate validation and explicit model-only or verified evidence attribution. Default WORK resolves the unique accepted head; stale and competing releases fail closed. Model tests pass including actual two-clone Git transport. Review fixes add abort of unjournaled staging and base64 conflict values with bounded YAML round-trip validation before publication.

2026-10-05T00:36:50.759403+00:00 MGI5-M1: Artifact-aware discovery exposes kind, UUID and preliminary role, prunes only owned histories and transaction staging, and preserves ordinary projects. Read-only snapshot returns captured bytes/digest without a commit. Full SDPTool tests pass, including compiled no-Git CLI lifecycle with real SDL/SDUI validation and copied-area inspection. Provenance now includes local author/acceptor attribution, original names and merge/restore references.

2026-10-05T00:39:38.822402+00:00 MGI5-M2: Integrated candidate fd7033b passes SDPTool race suite and vet, SDL parser/sourcegraph and SDUI parser tests, compiled CLI lifecycle, and Windows amd64/macOS arm64 cross-builds. Independent fresh-context review approves bounded Linux implementation after regression fixes. Child-process recovery covers dirty restore at six boundaries. Canonical SDL activity is implemented and nine viewpoint diagrams were regenerated; product release/main merge remain excluded.

## Completion — MGI5-M2

All five PLAN-SDP-0019 phases delivered with independent review and integrated
verification. [Evidence](../../05--Implementation/SDPTool/ModelGovernance/Evidence.md)
and [review](../../05--Implementation/SDPTool/ModelGovernance/Review.md) identify
claims and limits. Session0006 closes; KB050 remains separate semantic blueprint
work. No main merge, SDP release, installation migration or XFMD change performed.
