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
