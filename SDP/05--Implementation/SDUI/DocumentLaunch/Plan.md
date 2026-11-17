# DocumentHost prototype launch

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0030 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| Systems | SDUI, SDPTOOL |

Owner authorizes implementation after the TimberFrame reproduction (KB-SDUI-007).
Primary execution card KB-SDP-065; reuse Activity0007 with a new execution Session.
Use branch sdp/document-launch, milestone commits. No public release or main merge
is selected. One bounded Sol/medium worker owns SDUI prototype context and native
composition; parent owns facade, model, docs, integration and one final review.

| Milestone | State | Acceptance |
| --- | --- | --- |
| DL1 | active | Validated WORK model and implementable design; common context and native DocumentHost route, facade CLI/actions and focused tests |
| DL2 | planned | Packaged engine native TimberFrame reproduction succeeds; preservation/error checks, final review, model release and honest closeout |

## Design decision

Use the existing DocumentHost for SDUI 0.3 and explicit context-driven runs; preserve
legacy 0.2 behavior without context. No parser grammar, callbacks or domain runtime
execution changes. Reusing the actual host avoids a second pane implementation.
Removing unsupported guards alone cannot supply options, resources or native state.

Add optional --context FILE and --context-revision SHA to check/run and corresponding
action parameters. Context schema sdui-prototype-context/1 contains optional viewport
{width,height}, choices keyed by exact instance path with {id,label,enabled} options,
and svg keyed by exact instance path with {file,width,height}. Paths are relative
to the context directory, bounded local files with containment validation. No
implicit discovery or external program execution. Unknown/wrong paths, missing
choice providers including hidden controls, malformed contexts and stale expected
identities reject. SVG bytes are explicitly associated with the model reference.

Headless check validates source, capabilities, fixture inputs and semantic readiness;
it reports host=document and nativeLayout=pending, not native geometry success.
It cannot fabricate collection providers. Launch performs actual native resource
and layout preparation before publication. Legacy check semantics remain unchanged
for 0.2. Static preview does not consume a launch context.

Native DocumentHost atomically adopts prepared snapshots and preserves compatible
state. Watch reloads source, context and resource changes; failures retain the last
valid UI and report a diagnostic. Disk snapshot guards run before UI dispatch;
DocumentRequest.Guard compares watcher-observed generation identity without I/O.
This guarantees coherent observed snapshots, not atomicity against arbitrary
filesystem writers. Initial expected source/context revisions are checked; subsequent
watch edits are intentional updates. No domain callbacks execute; values are local
prototype state, provider/resource data remain caller-owned immutable snapshots.

## Model and blueprint decision

Create and validate a bounded LaunchWorkflow WORK under SDUINativeHost before
implementation. This adds the launcher/context responsibility detail to the existing
host catalog without migrating its historical root model. Release the model after
verification, with code/evidence identity. This focused two-layer change needs the
explicit model/design contract; a separate generated task blueprint adds no useful
scope decomposition for the single bounded worker and is not selected.

## Scoped concurrent-work preflight

Producer baseline 672fcb0, only existing ledger lock; fresh consumer main 9b7af853
clean. Historical detached release worktrees 774e150,84d5e46,f7a6c9f are unchanged.
Activities 1–7 completed before selecting this continuation; active PG/blueprint
cards concern other areas. PLAN0030/KB065 identities checked against current ledger.
Worker exclusively owns SDUI/go/prototype and cmd/sdui-fyne; parent owns SDPTool,
model/records/docs and integration. No edits to consumer product/model sources.
Green for this scope; unknown remote writers remain outside observations. Recheck
on participant/head/scope change. Model and package release are distinct; installing
or publishing a product release is not implied by this implementation request.

## Verification

Targeted loader/context/stale/error tests; native DocumentHost behavior and legacy
checks; CLI/action validation and native argument forwarding; one actual packaged
engine against TimberFrame Analysis with explicit fixture context in isolated Xvfb.
Exercise tabs/splits and reload failure preservation through bounded native harness.
Consolidate independent review at DL2; do not repeat broad installation campaigns.
