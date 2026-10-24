# SDPTool Session workflow concern — T041

Owner Session0007 T041 selects the bounded design described in
[Session execution](../../../../04--Design/SDPTool/ProjectGovernance/Session-Execution.md).
This is a detail view of the existing SDPTool/SdpToolHost responsibility, not a
new software system or runtime container. Existing System.design and the separate
ModelGovernance concern remain authoritative at their scopes; no source migration
or replacement of their model histories is selected.

## Artifacts and status

- WORK--SessionLifecycle is the mutable TARGET, created through ModelGovernance
  with --initial because this new concern area has no accepted RELEASE.
- PROPOSAL--ObservedBaseline--24a4 is a frozen source-observed baseline of four
  bounded responsibilities at code 7850fdf. It is a comparison input, not an
  owner-accepted NOW release or a complete reconstruction of all existing PG.
- WORK history retains the initial empty baseline, observed-foundation checkpoint
  and target checkpoint. Do not edit metadata/checkpoints by hand.
- No CANDIDATE or RELEASE of the target exists. Model integration, exact-candidate
  review, code verification and accepted model release remain explicit obligations
  in PLAN-SDP-0018. This design turn does not close the PG assignment.

## Source ownership and language coverage

System.design selects design-core/0.6 and composes Foundation.design (bounded
existing responsibility projection), Workflow.design (planned Session ownership and
function allocation) and Contracts.design (shared command/effect/host boundaries and
four illustrative scenarios). ExternalOwnerCli and ExternalLeaderMcp are external
caller boundaries, not containers deployed by SDPTool. ExistingRunFoundation is
implemented at the four named responsibility scopes; SessionWorkflowDelivery is
planned. Stage predicates and crash invariants are specified in the linked contract.

SDL validates structure, ownership, message participation and declared scenario
ordering. Its modes are contexts, not an executable finite-state-machine grammar;
these sources do not simulate runtime guards or prove crash safety. No new language
profile or unsupported state/transition syntax is introduced.

## Reproduction and evidence

Build repository SDL and SDPTool commands. From repository root:

```sh
sdl check SDP/SDL/ProjectGovernance/SDPTool/SessionWorkflow/WORK--SessionLifecycle/System.design
sdl viewpoints SDP/SDL/ProjectGovernance/SDPTool/SessionWorkflow/WORK--SessionLifecycle/System.design --format static --output /tmp/sdp-t041-views
sdptool SDP/SDL/ProjectGovernance/SDPTool/SessionWorkflow model status work:SessionLifecycle --json
sdptool SDP/SDL/ProjectGovernance/SDPTool/SessionWorkflow model create blueprint from proposal:ObservedBaseline to work:SessionLifecycle --entry System.design --task SDP/04--Design/SDPTool/ProjectGovernance/evidence-T041/blueprint-task.json --output /tmp/sdp-t041-blueprint --json
```

[Evidence](../../../../04--Design/SDPTool/ProjectGovernance/Evidence-T041.md) pins the
actual source/tool hashes and output manifest. The diagnostic Blueprint intentionally
has no phase-specific change permissions, so ConstraintsPass is false. It is useful
overview/impact context, not an implementation-ready assignment. The baseline was
captured before canonical formatting; ModelGovernance and Blueprint parsing validate
its semantics. The current WORK passes the strict canonical SDL check.

Local .model-operations journals/backups remain on disk but are excluded from Git;
tracked WORK checkpoints and frozen provenance retain portable model history. No
pending model operation may be hidden by that exclusion. Generated views stay in
external output directories and must be rebuilt from validated source.

## T042 bounded implementation mapping

PGF1-C1 now implements the binding/persistence subset of SessionCoordinator and
RunRepository, with registration recovery through the existing publication
primitive. [Evidence](../../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C1.md)
records the actual boundary. SessionWorkflowDelivery remains planned because
selection, canonical card effects, dispatch, design/model gates and disposition
are not all implemented. WORK :00002 remains the same authoritative TARGET;
no unsupported partial-status syntax or premature CANDIDATE/RELEASE is introduced.

A separate phase-scoped diagnostic Blueprint now passes its explicit constraints.
Its preliminary status and observed comparison baseline are retained; this is
same-context scope checking, not independent review or accepted NOW evidence.
The whole-workflow model review/integration/release obligation remains open.


## T043 checked selection mapping

PGF1-C2a implements a bounded SelectCoveredStep/PlanCardLifecycleEffects path through
SessionCoordinator and CanonicalEffectPublisher, with persistent intent/recovery in
RunRepository. Scoped claim checks and canonical history coordination participate
in that selection. The [evidence](../../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2a.md)
retains code/test identity and limitations. No model source semantics changed;
WORK :00002 remains the supplied TARGET and validates unchanged. The aggregate
SessionWorkflowDelivery remains planned because design/dispatch/verification/model
release and complete disposition are not delivered. The scoped diagnostic Blueprint
was generated before implementation and its constraints pass; independent product
review and accepted model release remain unfulfilled gates.


## T044 design-admission mapping

PGF1-C2b1 implements the explicit assessment/routing subset of RequireDesignWork
and AdvanceWorkflowStage, reusing CanonicalEffectPublisher and RunRepository for
records/recovery. WORK :00002 validates unchanged; the pre-implementation diagnostic
Blueprint passes its scoped constraints. This is partial target implementation,
not an accepted NOW, independent review or model release. SessionWorkflowDelivery
stays planned until run integration, design/model handoff and verified disposition
are delivered. [Evidence](../../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2b1.md)
records the bounded candidate and remaining obligations.


## T045 start/progress mapping

PGF1-C2b2 implements the no-design subset of AdvanceWorkflowStage and
DispatchAfterCommittedEffects through the existing run/controller foundation.
Canonical in-progress records precede ready-run creation; current assigned report
and candidate evidence precede verification-stage commitment. Required skill snapshots
and the full selected-card/contribution binding follow controller context. Unknown or
design-changing work retains RequireDesignWork/PGF3 gates.

WORK--SessionLifecycle :00002 remains unchanged TARGET. The scoped preliminary
Blueprint passes constraints; this does not accept the observed baseline or release
the model. Aggregate SessionWorkflowDelivery remains planned while disposition,
model handoff/release, integrated native validation and independent review are open.
[Evidence](../../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2b2.md)
identifies the exact candidate, model validation and implementation limits.


## T046 verification/review mapping

PGF1-C2c1 implements the candidate-bound verification/review subset of
AdvanceWorkflowStage through existing evidence/review owners and ApplyCheckedRecordEffects.
All selected contributions remain in canonical journals and review context. No
CheckWholeCardDisposition or claim release is asserted. WORK :00002 stays unchanged
TARGET; the scoped diagnostic Blueprint passes constraints and source validation.
Aggregate SessionWorkflowDelivery remains planned until the remaining runtime,
independent review and model integration/release obligations are met.
[Evidence](../../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c1.md)
retains exact code/model identities and the test boundary.


## T047 implementation mapping — PGF1-C2c2a

The existing CheckWholeCardDisposition target now has checked selected-step effects
for the no-design, owner-acceptance-not-required route: whole-card completion versus
contribution remaining scope, result publication and claim release before operational
completion. The scoped Blueprint constrains this subset; WORK :00002 sources remain
unchanged. SessionWorkflowDelivery remains planned: owner decisions, rework/dependency
progression, overall closure, leader/model handoff and acceptance/release are open.
[Evidence](../../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c2a.md)
distinguishes synthetic report/review fixtures from native or independent PG review.
