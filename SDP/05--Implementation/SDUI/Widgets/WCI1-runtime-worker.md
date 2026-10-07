# WCI1 runtime Worker handoff

Date: 2026-10-08. Role: explicitly assigned SDP Worker, runtime lane.
Working directory: `/tmp/sdp-sdui-widgets`; branch `sdui/widgets-wci1`.
Candidate: HEAD `d1c5c88d1989b173d610dad916f6f493f838ab7b` plus the uncommitted runtime sources listed below.
No commits, branch changes, Session/plan/card/ledger edits or other-lane source edits.
The coordinator subsequently assigned two documentation files to this lane: runtime README and runtime-contract; both are now updated.
Concurrent frontend/layout/host/bridge work exists in this shared clone and is not claimed here.

## Outcome

Runtime lane is implemented and ready for independent review/integration. The API was published first in `WCI1-runtime-API.md` and sent to coordinator/layout/host. Implemented typed collection/provider/target/request/state contracts, exact instance-path provider admission, bounded whole/subtree validation, atomic data mutations, stable item identity, request supersession/cancel/retry/completion, local typed events and the full post-handler collection-target guard. Runtime never invokes a provider, imports layout/GUI, or starts a goroutine. The host owns context cancellation and asynchronous delivery.

The pure state gate receives detached model/collection/viewport snapshots and returns all clamped offsets before publication. Every accepted UI-state mutation advances StateRevision; consumed callback events advance it even when unbound. Failed prospective validation preserves live state/sequence. Separate viewport owner handles plus checked model revision reject stale native scroll commands. Trusted bulk offset changes remain available for geometry/chained scrolling.

Detached Successor carries the application event-sequence watermark, non-reused live widget/viewport generations, compatible draft/value/focus, collection data/local state and request watermark, and compatible offsets. It advances collection/model revision, drops old handlers/gates/request tokens and leaves the predecessor untouched. Legacy Reload retains compatible legacy handlers through this candidate path. Snapshot cloning preserves Instance.Profile and deep-copies Uses, provider data, collection maps/request pointers and viewport maps.

Coordinator-approved clarification: an unloaded empty root after compatible reload may have AutoLoadPending=false because a previous load was already started. Explicit Retry/Load accepts that paused root; it never restarts automatically. Repeated Retry while loading consumes only its event sequence, allocating no second request and invoking no domain callback. Branch recovery is not broadened.

## Validation

All commands below ran with Go `go1.27.1 linux/amd64` from `SDUI/go` unless noted.

- `go test -race ./runtime` — PASS, final run 1.225s. Includes original 0.2 tests and new collection/state/successor coverage.
- `go vet ./runtime` — PASS.
- `go test ./runtime -coverprofile=/tmp/wci1-runtime-coverage.out` — PASS; 88.1% statement coverage. This metric is supplementary, not acceptance evidence on its own.
- `go test -run '^$' ./layout ./preparation ./host/fynehost` — PASS compilation/API integration check; **no tests executed** by this command. It is not native evidence.
- `git diff --check -- SDUI/go/runtime` from clone root — PASS.

Substantive test obligations:

| Runtime obligation | Reproducible tests |
| --- | --- |
| Invalid IDs/UTF-8/controls, limits, duplicate IDs, missing parents/cycles, list/branch flags; atomic mixed-invalid batch | TestCollectionDataValidationAndAtomicMixedBatch |
| Exact provider admission, no Load during preparation, copied data, binding identity/epoch | TestProviderAdmissionAllInstancesAndOwnedCopies |
| Subtree order/unrelated state, generation, remove/recreate identity, row-to-group | TestSubtreeReplacementIdentityAndSelection |
| Prospective gate, effective clamp, invalid/stale viewport offsets, no partial data | TestPureGateAtomicDataAndClampedViewport; TestViewportGateRejectsInvalidEffectiveOffsetsAtomically |
| One request, supersession, late success/error, empty success, error/canceled recovery | TestLoadCancellationSupersessionAndFailureRecovery; TestRootAutoLoadAndExplicitRestart |
| Invalid data/layout completion preserves old data/selection/offsets and shows recoverable error | TestCompletionRejectsInvalidDataAndLayoutWithoutPartialPublication |
| Ancestor collapse/hide/disable/refresh/provider/reload/Close revoke requests | TestLoadRevocationForAncestorHideRefreshProviderAndClose |
| Provider ignoring cancellation; deterministic channel barriers, no sleeps | TestLoaderBarrierDeliveryCannotPublishAfterCancellation |
| Typed payload/identity/sequence/visible-row validation; local navigation without SDL | TestTypedCollectionEventsLocalNavigationAndSequence; TestRetryEventsAreLocalAndOneRequest |
| Reentrant identical data/delete-recreate/kind/collapse/hide/disable/reload/Close; pointer mutation cannot weaken guard | TestCollectionResultGuardCoversReentrantMutations |
| Failed local event/loading preparation consumes neither sequence nor live request | TestLocalEventGateFailureConsumesNeitherSequenceNorState; TestFailedDataAndLoadingGateRetainActiveRequest |
| Reused instances have separate provider/state/request identities; deep source provenance | TestReusedInstancesOwnProvidersRequestsStateAndProvenance |
| Detached preparation/disposal leaves prior live state/requests intact; compatible local state | TestDetachedSuccessorRetainsCompatibleStateAndNoHandlers; TestSuccessorProviderChangeAndFailedAdmissionRetainOldBundle |
| Sequence across three successors; removed/recreated widget generations | TestSequenceAndGenerationAcrossMultiplePublishedSuccessors |
| Root auto flag/status/recovery retained across compatible reload | TestRootRecoverySurvivesCompatibleReloadWithoutAutomaticRestart; TestExplicitRootLoadAfterSuccessorAllocatesOnceWithoutDomainCall |
| Viewport removed/changed kinds/axes revoke or preserve correctly | TestViewportIdentityRemovalKindChangeAndAxisPreservation |
| Uniform/mixed/unknown/cyclic profile handling, deep Uses, state revisions, 0.2 empty profile | TestSnapshotProvenanceProfilesAndStateRevision |

## Integration boundaries and remaining evidence

No native UI, provider context lifecycle, actual SDL invocation, application-wide publication, full SDUI suite or Xvfb evidence is claimed by this lane. The coordinator/other lanes own that integration and independent review. Bind providers before installing the prospective state gate. Host compares accepted request tokens and cancels old contexts only after revocation; loader results return on the owner goroutine. Source generation and parser profile support belong to the frontend lane.

Gate purity is an API requirement. If a host gate rejects even the fallback error presentation, CompleteLoad returns that gate error and preserves the last valid runtime state; the host must surface that out-of-band failure. No runtime bypass of geometry validation is used.

Coordinator Session payload: owner assigned authorized WCI1 runtime-only implementation, then required checked viewport identity and deep Uses cloning, then approved explicit paused-root recovery. Loaded sdp 1.1.1, sdp-worker 2.0.0 and document-workflow. WCI1 runtime outcome is this tested candidate; Session/plan/card/ledger updates remain with main per exclusive assignment. Next step: independent runtime review and native/bridge integration acceptance. No exact owner-prompt timestamp is invented.

## Assigned documentation closeout

Updated `SDUI/go/runtime/README.md` and `SDUI/docs/runtime-contract.md` after explicit coordinator assignment. They now document the actual typed collection/provider/request APIs, public versus normalized path namespaces, pure Snapshot/state gate, checked viewport identity, detached Successor, StateRevision/sequence retention, paused-root recovery, and bridge CollectionItemID-to-text mapping. Original 0.2 APIs/semantics remain documented. Native acceptance and future widget families are not claimed. Relative links and scoped diff whitespace validation passed. These documentation-only changes do not change the tested runtime source hashes.

## Changed runtime files

- `SDUI/go/runtime/collection_events_test.go`
- `SDUI/go/runtime/collection_isolation_test.go`
- `SDUI/go/runtime/collection_load.go`
- `SDUI/go/runtime/collection_load_test.go`
- `SDUI/go/runtime/collection_mutation.go`
- `SDUI/go/runtime/collection_validation.go`
- `SDUI/go/runtime/collections.go`
- `SDUI/go/runtime/collections_test.go`
- `SDUI/go/runtime/events.go`
- `SDUI/go/runtime/properties.go`
- `SDUI/go/runtime/reload.go`
- `SDUI/go/runtime/session.go`
- `SDUI/go/runtime/state.go`
- `SDUI/go/runtime/successor.go`
- `SDUI/go/runtime/successor_test.go`
- `SDUI/go/runtime/types.go`

## Additional assigned documentation files

- `SDUI/go/runtime/README.md`
- `SDUI/docs/runtime-contract.md`

## Tested source SHA-256

```text
7a05658ae0aa60fc28594697c53452371c2af8ca5009ee8a6f5e714bba50447e  SDUI/go/runtime/collection_events_test.go
9e2663d1d463fc9c0d8c25421c4c8a6fb7a29e91df0bc1006da4bf86fae7f471  SDUI/go/runtime/collection_isolation_test.go
6c5629a3bafe82378ac42ae79b0a096256160b28b00ddbc36b773440bc2ca0bb  SDUI/go/runtime/collection_load.go
68b056511f20f9bbc0e3dc473a8cb2bef37ae8293f5a041a365d7c3ea32cc38f  SDUI/go/runtime/collection_load_test.go
625d5aba04e7ef3a4a8f284e22376215c02d229f4796424e03a48ee08d3a6a5e  SDUI/go/runtime/collection_mutation.go
ecf6f836d0cd48b028c8bf5770eccadc68490af6412d2d6deb8f8b62e2c6b37e  SDUI/go/runtime/collection_validation.go
55af3b03b1fe2906b61d6e37bbc49a85f4b110ce89f3c1fa3e9ce3128fadeb7d  SDUI/go/runtime/collections.go
17e0ad4be7cbe35f1e55269baa4120754d3be16bc2336e16d5f83a997842d30b  SDUI/go/runtime/collections_test.go
13be26c1e9a7900b7ad3ae25ed0635a75ace76d4a260f2e588fdfa644fe0dfa3  SDUI/go/runtime/events.go
3af1e75c3413ecda7e39ddad19527b20047d49d82877a55e578bb11f4acc7ccc  SDUI/go/runtime/properties.go
11580db0bc9e1f3eb3cb979e1a0213b131782fcf640c8acebf0949a5bbadc2fa  SDUI/go/runtime/reload.go
c8431a3ad724b56818926515262947aeb0a57b31ae9a01e554be9b1925fa02ac  SDUI/go/runtime/reload_test.go
f68a4b0d5503b8241195fb88864e043aa57b5859476363e62bfe26e6fdf1d515  SDUI/go/runtime/session.go
66dd49250795873934bee0d5f8c64acad22db29bf9e17500c7b48b026b8f22fa  SDUI/go/runtime/session_test.go
a805354fdf26fcb0072b2de7bc16088d5f741b82f9e30f5b1356af5d57ee3f7f  SDUI/go/runtime/state.go
954cfae9a28d5a5332fa7208ca9b105fdfe10a7f06c440a1ecbdaaa6f7568d7d  SDUI/go/runtime/successor.go
6ed8d32380d7b45052f5a5a2e7a3f3b25e83d9d8e944f8ba35ad27a253b6ad24  SDUI/go/runtime/successor_test.go
d2f14c287a4bafd266e65957d340479e365ecdd32354cabed3e77bc6bfce58a8  SDUI/go/runtime/types.go
cd01111b146df65c817f8242effa9a92e351d2674d03ca6a63279f5e15ea51b5  SDUI/go/parser/ast.go
39a6c4f0a2f061699773ac7cbfc588107186eb979c4270850fc4a677071d03b7  SDUI/go/parser/formatting.go
479109b1443ee490a5dc704054b25ce776518415d22d2348732658d0481a9778  SDUI/go/parser/lexer.go
253c3fee761c1966a83c3c2f78be25258684b329e51953514b5bf2d3e44d5004  SDUI/go/parser/normalize.go
a25394519dd3961d1bfd6ff20be27df2516789ab7b2b9e0eb843739dc4cd4753  SDUI/go/parser/parser.go
517498c68b9666a153d1011cb6c9d98381a27120d8b8f833a2c0e87cd33b39ac  SDUI/go/parser/profile.go
1899629a5c714c8fb889641b07ae851700d3ec7822b0b63a7e5709791a1e0468  SDUI/go/parser/validate.go
```

## Documentation SHA-256

```text
25eca9f7c74d56c43f54a7060b2cede17a2e4614038df41dcc7ad31a78edbcc5  SDUI/go/runtime/README.md
90dfd2caaa3745afd50ad7c4ee769ce23a142c3717301cf42e1d4a85704544a2  SDUI/docs/runtime-contract.md
```

## Read-only DocumentHost integration check (coordinator follow-up)

This is an implementation integration check, **not** independent review or a gate.
No host/preparation/layout/bridge sources were edited. A standalone headless Fyne
harness outside the clone exercised public APIs with channel barriers and no sleeps.

Command from `SDUI/go`: `go run /tmp/wci1-runtime-integration-check.go`.
Transcript: `/tmp/wci1-runtime-integration-check.txt`. The command completed with
exit status 0 because it reports observations rather than asserting acceptance.
It reproduced three integration mismatches; WCI1 native integration is not closed.

1. **High: accepted reentrant runtime state is not synchronized on Dispatch error.**
   `DocumentHost.Mutate` runs `after()` only for nil error. A handler performs
   ReplaceCollection successfully, then Dispatch correctly rejects its stale result.
   Observed runtime data `new`, native inspection row `old`, returned `stale-result`.
   The host must synchronize/reconcile any accepted state mutation despite the
   outer error, e.g. compare StateRevision before/after and apply accepted pending
   presentation while still reporting the error. Revoked provider contexts need
   the same reconciliation.
2. **High: failed small-window resize can strand completed provider requests.**
   Resize stores `h.size` before the prospective gate, and subsequent stages use
   that invalid size. At size 1×1 both resize and CancelCollection rejected layout;
   provider context remained uncanceled. Matching completion and its fallback error
   presentation also failed; restoring 800×500 left `Loading`, an active request
   whose goroutine had finished, and Retry returned nil without restarting.
   Keep physical clipping and last-valid measurement distinct for lifecycle recovery,
   or explicitly retain/drain failed completions. Do not silently discard a completed
   provider reply while retaining its loading token/flight.
3. **Medium: resource preparation receives requested rather than clamped offsets.**
   Bundle.stage passes its incoming Snapshot to PrepareResources and stores it in
   pending, then returns measured effective offsets to runtime. Observed hook Y=999,
   accepted runtime Y=0 and measured geometry Y=0. Built-in collection painting uses
   the measured offset and was coherent in this case; custom resource preparation
   can be inconsistent. Replace snapshot.Viewports with EffectiveOffsets before
   resource preparation/pending storage.

Positive observed boundary: changing the predecessor draft between Prepare/Commit
caused `stale: published state changed during preparation`; the old Session remained
open. Static inspection also confirms that ordinary failed detached Prepare closes
its candidate rather than the prior Session, and runtime clones do not share
collection requests/maps with it. No defect in that StateRevision guard was found.

Additional consequence to inspect after fixing #1: collectionEvent currently builds
a new target with Session.Target instead of retaining the rendered collection
generation. If stale displayed rows exist, that refresh can retarget a recreated
ID. Normal bundle identity checks protect replacement bundles, but rendered data
generation should also be validated. This consequence is from code inspection,
not a separately executed activation assertion.

The three reproduced findings and positive guard result were sent to coordinator
and Noether. Host owns fixes and integrated acceptance; runtime candidate hashes
and runtime tests above remain unchanged.

Read-only check artifact SHA-256:

```text
1586daa9ff747d816e451ef6098a5f115bd97a4ee8c741cd36610de6b4e1ccd4  /tmp/wci1-runtime-integration-check.go
f0bc228d182be169c7d74201a39cb86ed971a450be8f8227666e8420005487ed  /tmp/wci1-runtime-integration-check.txt
d2c0b4a9364cb4a8beac22f11a168e47a1ce914437ae8361ac6cb3121b7770d7  SDUI/go/host/fynehost/document.go
260b41271659accea8bb3b7335e9d1ec3b7d2f7e8e05b851bb6122311b1c2ded  SDUI/go/host/fynehost/document_sync.go
3eceed0d5122ad537263fe865a4b54215cfdb22ab578bb6305c65ad855bb2b50  SDUI/go/host/fynehost/document_events.go
a1b346d1e806433c72b93fece23835053dfb470a5aed07f1bcd7a3f17b86998b  SDUI/go/preparation/prepare.go
```

## Authorized DocumentHost lifecycle follow-up

After the read-only findings, coordinator transferred only
`SDUI/go/host/fynehost/document.go` and the new
`SDUI/go/host/fynehost/document_lifecycle_test.go` to this Worker. The documented
Noether freeze was confirmed before implementation. Starting document.go SHA-256:
`d2c0b4a9364cb4a8beac22f11a168e47a1ce914437ae8361ac6cb3121b7770d7`.
No other host file or runtime source was edited in this follow-up. No commits.
Noether supplied b.size support in document_sync and owns clipping wrapper,
rendered-target and native thumb changes.

The three reproduced host failures are corrected in document.go:

- Mutate synchronizes accepted StateRevision changes even when Dispatch returns
  a stale-result or other outer error. Valid reentrant data/requests therefore
  reach the native presentation and provider-context reconciliation.
- Bundle.size is the accepted measurement basis; host.size remains the actual
  physical window. Resize stages a new bundle size and restores the previous
  valid basis on geometry failure. Cancellations and provider completions continue
  against the retained valid presentation, clipped by the actual canvas. No
  runtime geometry gate is bypassed and completed replies do not get stranded.
- Resource preparation and the staged snapshot receive measured effective
  offsets, matching geometry and the runtime values accepted with it.

Actual native focus restoration is muted and checked around the final local
loading apply: it occurs after old bundle disposal, before publication observers,
then is rechecked without side effects if already correct. A real native pilot
had logged Fyne's “not part of canvas content” error. Inspection found that the
old routedClip wrapper was initialized with Scroll's own Base identity before
ExtendBaseWidget(wrapper); Fyne complete focus traversal skips unrendered wrappers.
Noether corrected the wrapper initialization. The new lifecycle test exercises
the actual software-canvas FocusManager with that correction, asserts actual
focused object equals the successor control, and routes R through Canvas.Focused.
It also requires no synthetic StateRevision increment during muted restoration.
Actual OS input/cache timing remains main's XTest acceptance, not this test's claim.

Validation progression:

- Before document.go changes, `go test ./host/fynehost -run '^TestDocumentLifecycle'
  -count=1 -timeout=45s` reproduced **five failures**: stale native data after
  reentrant error, unreconciled canceled context, requested/effective offset
  mismatch, cancellation blocked after failed resize, and stranded completed
  provider request. The software-canvas focus test already passed with Noether's
  corrected wrapper.
- The same targeted tests passed after the fix.
- `go test -race ./host/fynehost ./preparation ./runtime -count=1 -timeout=90s` —
  PASS (host 12.742s, preparation 1.058s, runtime 1.186s).
- Final lifecycle suite additionally asserts publication-observer native focus,
  failed candidate resource preparation retaining the predecessor flight, and
  invalid/nonfinite resize preserving prior state.
  `go test -race ./host/fynehost -run '^TestDocumentLifecycle' -count=1 -timeout=60s`
  — PASS (5.961s).
- `go vet ./host/fynehost` — PASS.
- `git diff --check -- SDUI/go/host/fynehost/document.go
  SDUI/go/host/fynehost/document_lifecycle_test.go` — PASS.

The initial review patch is `/tmp/wci1-document-lifecycle.patch`; it records the
first proposal against the frozen source, before final formatting/focus-observer
refinement. The actual final candidate is the files identified below, not that
initial proposal patch. Rebuilt native empty-root reload/focus and resize acceptance
is still required. This is implementation evidence, not independent approval.

Final owned host source SHA-256:

```text
8f4a0f908edbffa93d5170c9bc6efd72402fc418fcd9ec4cf53248158cf1c9d6  SDUI/go/host/fynehost/document.go
06753c229c5978f1038ab76093dd39e163602d4e827cfb41232002cf3754fb56  SDUI/go/host/fynehost/document_lifecycle_test.go
```

### Coordinator-reported actual native pilot4b

Main reported seven passing actual native empty-root checks against the current
document.go plus Noether's wrapper fixes. OS keyboard R after reload starts a
fresh request, the old completion is ignored, and teardown closes the Session
and clears pending work. Evidence directory:
`/tmp/sdui-wci1-harness-empty-pilot4b`. This result was supplied by the coordinator;
this Worker did not independently execute the native pilot. It supersedes the
empty-root focus uncertainty above for that candidate. A final candidate native
rerun after the pending gutter integration remains required and owned by main.
