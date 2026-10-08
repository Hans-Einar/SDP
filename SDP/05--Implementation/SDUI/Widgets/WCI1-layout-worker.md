# WCI1 layout Worker handoff

## Assignment and candidate

Authorized SDP Worker layout lane for WCI1-M1 / PLAN-SDP-0022 / KB-SDUI-003.
Implementation is confined to `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci1`.
Candidate: HEAD `d1c5c88d1989b173d610dad916f6f493f838ab7b` plus the eight layout source/test files below.
No commits, staging, branch switches, runtime/parser/host edits or original-workspace
writes were performed. Other lanes concurrently own their changes in this clone.
The test results are working-candidate results, not unchanged-commit evidence.

Loaded Skills/sdp/SKILL.md (1.1.1), Skills/sdp-worker/SKILL.md (2.0.0), and the
shared document-workflow. Read the pruned Collections.md contract, prerequisite,
plan, SDUI instructions/layout contracts, and the original Session0010 roadmap.
The clone lacks Session0010; the coordinator owns its update, as explicitly assigned.

API proposal was sent to main before dependent implementation and approved as a
bounded implementation choice. Read James's WCI1-runtime-API.md before implementing
snapshot integration; use his runtime Snapshot and ViewportState. Dalton supplied
parser.EffectiveProfile; no competing parser/profile or runtime offset types exist.
Noether received the collection measurement and routing API. Native adapter
integration/measurement acceptance is not established by this Worker report.

## Delivery

- Profile-aware layout admission uses parser.EffectiveProfile and verifies 0.3
  scroll owners, including hidden branches. Scroll owners are frame/group/tree/list;
  scroll axes require fill/fr/scale or a resolved frame ratio, and start justification.
  Unknown 0.3 instance/widget kinds reject. Legacy 0.2 unsupported-scroll diagnostics
  and geometry paths remain; existing legacy layout/SVG/prototype tests pass.
- Extended existing Box geometry for 0.3 container scroll: finite parent references,
  actual direct-child extents, effective offsets, composed screen clips, one ancestor
  translation per level, nested viewport ancestry, and outer-box-only contribution
  of nested hidden content. Regions scroll with their owner; oversized footer starts
  nonnegative. Content/operation bounds remain enforced.
- SnapshotLayout adds only outer geometry and viewport facts, not a row scene.
  LayoutSnapshot consumes the prospective detached runtime snapshot, returns matching
  geometry and EffectiveOffsets for the pure state gate, and never mutates Session.
  Nonfinite/negative requests reject; excessive requests clamp; removed/hidden owners
  are omitted; non-scroll axes have zero effective offset.
- CollectionMeasurer receives the actual assigned outer size and reports native
  outer Minimum, local Content, and inset Viewport. It separates full row extent
  from native minimum; title/gutters enter geometry once. No guessed row renderer or
  fallback collection metrics were added. Missing/invalid metrics reject explicitly.
  Collection outer Box moves only with ancestors. Its native adapter applies that
  collection's own accepted offset once within the returned viewport/clip.
- RouteScroll hit-tests the topmost branch, consumes each axis inner-to-outer and
  passes only unused delta to ancestors. Siblings and obscured/disabled branches
  cannot receive the routed scroll. EnsureVisible adjusts minimally inner-to-outer,
  translates the target before each parent, and aligns oversized targets at start.
  Both return detached complete offset maps; neither publishes state or events.

## Integration API

```go
func (*layout.Engine) LayoutSnapshot(runtime.Snapshot, layout.Size) (*layout.SnapshotLayout, error)
// SnapshotLayout: Root *Box; Viewports map[string]Viewport
// Viewport: Path/Parent, screen Rect/Clip, local Content,
//           ScrollX/ScrollY, Offset/Maximum runtime.ViewportState.
func (*layout.SnapshotLayout) EffectiveOffsets() map[string]runtime.ViewportState
func (*layout.SnapshotLayout) RouteScroll(x,y,dx,dy float64) (map[string]runtime.ViewportState, runtime.ViewportState, error)
func (*layout.SnapshotLayout) EnsureVisible(path string, target layout.Rect) (map[string]runtime.ViewportState, error)

type CollectionMeasurer interface {
    MeasureCollection(*parser.Instance, float64, layout.Size) (layout.CollectionMetrics, error)
}
// CollectionMetrics: Minimum, Content Size; Viewport Rect local to outer control.
```

A native measurer must implement ordinary Measurer plus CollectionMeasurer and
capture the SAME prospective Snapshot used by LayoutSnapshot, rather than inspect
live Session state. It reports rows/indent/status content, excluding fixed title;
Viewport excludes title and native gutters. Positive routed deltas move toward
content end. EnsureVisible target rectangles are current screen coordinates,
including the adapter's current own offset for collection rows.

Host must validate captured bundle/viewport identity before publishing the complete
returned map atomically with runtime.SetViewports (or runtime's checked batch API).
Do not publish each changed ancestor separately. Native sync must be muted and
mirror runtime offsets; no second offset authority. Stage native changes from the
pure gate and commit them only with the accepted runtime candidate. Layout does
not invoke loaders, bindings, callbacks, resource publication or disposal.

## Verification

Environment: `go version go1.27.1 linux/amd64`.
From `/tmp/sdp-sdui-widgets/SDUI/go`:

```sh
go test ./layout -count=1
go test -race ./layout -count=1
go test -race ./layout ./svg ./prototype -count=1
```

All PASS. Final affected-consumer run: layout 1.825s, svg 1.112s, prototype 1.166s.
`git diff --check -- SDUI/go/layout` also PASS.
One initial `go test ./layout` was invoked from the repository root and returned
“cannot find main module”; it was rerun from SDUI/go and passed. No product test
failure was suppressed or removed.

New independent geometry oracles cover:

- Two nested viewports plus sibling; exact translated rectangles/clips, half-open
  hits, nested content exclusion, non-scroll zero offset and returned-map isolation.
- Forward/backward delta remainder, sibling/overlay/disabled isolation, zero ranges,
  nonfinite deltas, inner-to-outer reveal and oversized leading-edge behavior.
- Clamp on shrink, removal/hidden owners, invalid offsets, root-only Engine reuse,
  mixed/unknown profiles, unsupported/indefinite hidden scroll owners and justify.
- Finite scale/fr/minimum allocation without fixed-track shrink; group scroll,
  ratio-resolved dimensions, relative dependency failure, regions/footer movement.
- Collection title/gutter viewport geometry, minimum versus long content,
  nested ancestor clip, one collection offset application, error/clip/scroll,
  missing/nonfinite/out-of-bounds metrics and content-sized chrome accounting.
- Actual runtime.CheckStateWith/SetViewports/ReplaceCollection: failed prospective
  geometry leaves runtime state and staged geometry unchanged; data shrink publishes
  the matching clamped offsets and geometry atomically.

These are deterministic layout/state-gate tests with synthetic adapter metrics.
They do NOT establish native Fyne row metrics, thumb dragging, real pointer/keyboard
input, painting, collection identity mapping or publication/disposal acceptance.
Full SDUI/SDL/SDPTool integration and independent review belong to main after sibling
lanes finish; this report does not close WCI1 or KB-SDUI-003.

## Changed files and SHA-256

```text
add919eaa11348fb9d960074a5b4d59cb05458362c5f6ce8db882a44575ff99c  SDUI/go/layout/engine.go
b9ea8e016f70e5b7ebb623b07126bfc19821048eb6c45356263328da40f8e3fb  SDUI/go/layout/types.go
fd7e1e7535a4fc4c1cce56bda18ddaf45066e9652d67dbb624c828c5675da5db  SDUI/go/layout/contents.go
08497a3a9dcbe4afbbcef1c2f55bcab0865a4a641ddb075cf5921f893740f411  SDUI/go/layout/viewport.go
67027fc7c585f165e20192c995d3135ac0396f3453f22f0479d182bb2dfd05d7  SDUI/go/layout/scroll.go
fcbad63445d5d07a17475d7c6887cf9d62b8c70090f87e616ee61623f53dc7d4  SDUI/go/layout/collection.go
5d0fd73b4b46c4b32045e80c200f0e0e663686ce7cbf34537306b130265f5e5d  SDUI/go/layout/viewport_test.go
33aa9a5b86c3e8bb8cd6c4dc22248b4e67d0fa53e416559ffa657de0497c5d19  SDUI/go/layout/collection_test.go
```

## Coordinator Session/document handoff

Session0010 S2 / WCI1 remains active. Record as current work summary: explicit
Worker layout assignment after independent design approval; scope only layout and
this report; routines loaded as above; main approved minimal snapshot/viewport API;
James owns runtime snapshots/offset identity and Dalton owns effective profile;
layout implementation and the stated tests passed; native integration and independent
review remain. Next: Noether integrates native metrics/routing, main verifies the
complete candidate and records acceptance evidence. No exact owner transcript or
unknown timestamps/host IDs are reconstructed here.

Coordinator should update SDUI/docs/go-layout-contract.md and affected architecture/
requirements/acceptance records to distinguish preserved 0.2 rejection from the
new 0.3 snapshot/scroll path, with this actual evidence and remaining native scope.
No new owner decision or scope change is claimed.


## Addendum — layout contract and read-only host inspection

The subsequent bounded assignment authorized only SDUI/docs/go-layout-contract.md
and this addendum. Reused the already loaded Worker/document workflow; reread the
root English-documentation rule, previous layout contract, neighboring profile-0.3
contract and current host integration. No layout or host implementation was edited.
All eight layout file hashes above were checked and remained unchanged.

Updated the English layout contract with preserved G2/0.2 behavior, 0.3 admission,
finite scroll geometry, actual SnapshotLayout/Viewport/CollectionMeasurer signatures,
coordinate ownership, pure state-gate integration, routing/revelation and explicit
native acceptance limits. No WCI2 or later widget-family capability is claimed.
Local Markdown links, fence balance and scoped git diff --check pass. No product
tests were rerun for the documentation-only change.

Read-only integration findings were sent directly to main and Noether. These
observations concern the evolving host candidate, not independent final approval:

1. **Measured focus target:** initial document_events.go ensureItem used viewport
   width at an indented row origin. A 50-unit-wide row starting at x=18 inside a
   200-unit viewport was supplied as width=200, unnecessarily scrolling 18 units.
   Later readback shows nativeText(rowText(...))+8 instead. Code adjustment observed;
   target minimality and native focus behavior still need integrated evidence.
2. **Keyboard scroll coordinate:** initial PageUp/PageDown used viewport.Rect+1;
   with Rect.Y=-30 and Clip.Y=0 that point is outside the visible window and the
   shared hit-based router consumes nothing. Later readback uses the effective
   Clip center. Code adjustment observed; nested keyboard proof remains pending.
3. **Row text extent:** initial MeasureCollection used unprefixed label width plus
   fixed 32, while rendering included branch/selection/focus markers at arbitrary
   font sizes. Later readback uses shared rowText, nativeText and the renderer's
   regular font. Code adjustment observed; exact native metric tests remain pending.
4. **Thumb mapping remains open in inspected candidate:** Dragged divides absolute
   pointer position by full viewport length, while painting uses thumb travel
   (viewport length minus thumb length). For height=100, thumb=20, max=100 and
   offset=25, the painted thumb starts at 20; grabbing that top position with zero
   delta requests 20 rather than 25. Preserve the grab offset and use actual travel.
5. **Minimum mismatch remains open in inspected candidate:** collectionRenderer
   MinSize returns width 70 and height 2*rowHeight; MeasureCollection reports width
   3*rowHeight+gutter and height 2*rowHeight+gutter. Admission and native minimum need
   one shared calculation, particularly across font sizes and scroll-axis choices.
6. **Frame/group native routing coverage gap:** only CollectionControl implements
   Scrolled/Dragged in inspected host sources. DocumentHost/outer frame/group
   backgrounds have no scroll target or visible thumbs. Existing wrappers use
   ScrollNone; pinned Fyne 2.8.1 internal/widget/scroller.go Scrolled ignores input
   in that mode. A collection can route remainder to its parent, but wheel over
   frame/group background or ordinary controls has no equivalent route here. An
   Alt+arrow route was also absent. Main/Noether were asked to cover these bounded
   native interactions; no geometry-code defect or new product API is assumed.

Positive coordinate observation: the host stores shared outer Box bounds, applies
ancestor clip offsets in its wrapper, leaves the title outside row content, and
subtracts the collection's own runtime offset within row painting/hit tests. This
is structurally consistent with the layout boundary; it is not native paint/input
proof. No host file was modified by this Worker.

Readback identities (SHA-256; host files may change further in their owning lane):

```text
7cbb894339e7132636e3c0612f1df3124559d1e2ee6f57a137a64a09a962bba5  SDUI/docs/go-layout-contract.md
5480f1e2315289ce8082477a0effa57891dce43925a74fe945f5529232ae5827  SDUI/go/host/fynehost/collection_metrics.go
8869826a24fd964ccb46814eb25a3f882c4c2c313b6b48b50c7d775834abd47d  SDUI/go/host/fynehost/collection_control.go
3eceed0d5122ad537263fe865a4b54215cfdb22ab578bb6305c65ad855bb2b50  SDUI/go/host/fynehost/document_events.go
adc914333b54af40ad3214241c279075e5e86a65b44b4c21d71664fa6f142d32  SDUI/go/host/fynehost/document_sync.go
```

Session0010 S2 work-summary payload: owner requested bounded current-layout
contract upkeep and read-only native integration inspection; documented implemented
APIs with 0.2 preservation, sent six concrete findings, observed ongoing fixes to
first three, held layout code unchanged. Native proof, remaining host fixes and
integrated review remain pending. Coordinator owns Session/plan/card projection;
this addendum is a handoff, not a replacement journal or completion claim.


## Addendum — delegated native fixture controls and nested variant

Main assigned a new bounded Worker lane for SDL/go/examples/collections after
Noether's explicit whole-subtree freeze/handoff. That confirmation was observed in
Noether's active-thread commentary and acknowledged to both main and Noether before
any fixture edit. Noether had already fixed canonical SDL ordering; his real SDL
fixture/bridge test passed before handoff. This Worker preserved that action source
and registration and did not edit host, runtime, parser, bridge or layout code.

The English fixture README is part of the fixture-only handoff. Main retains all
actual XTest/OS screenshot work. No commits or branch switches were performed.

### Delivered fixture behavior

- Added NestedSource: same tree/list paths under a scrolling body, with collections
  1.6 times the finite body height. Body and collection are two nested viewport
  levels; the fixed footer preview is an independent sibling. Default Source and
  the real retained SDL engine remain available. The window title stays
  `SDUI WCI1 Collections`.
- Added a fixture-only Controller for strict stdin commands using existing
  DocumentHost.Mutate/runtime operations. New controls cover resize, shrink,
  reorder, reset, hide/show, disable/enable, pending requests and phase-specific
  failed reloads. No product control-channel API was introduced.
- Added source/action SHA-256 fields and pending-barrier keys to state JSON, retaining
  row/viewport screen rectangles and clips from Bundle.Inspect. `command-result`
  correlates monotonic commandId, exact command line, status ok/error, candidate
  sequence, current source/model revision, action count, bundleChanged and error.
  Expected failure commands return status error; they are not silent success.
- Completion accepts only success/error/empty/invalid and releases one named
  started barrier. Invalid outcomes no longer consume a request. Pending keys are
  owned/sorted snapshots. Command acknowledgement is not completion acceptance;
  load-return and accepted/rejected runtime state remain separately observable.
- Main reported that programmatic w.Close bypassed CloseIntercept. Updated native
  CLI close and OS-close intercept to invoke the same sync.Once cleanup: host token/
  context revocation, fixture barrier/engine cleanup, closed JSON, then window quit.
  `closed` includes sessionClosed, remaining pending keys and action count. Native
  verification of this fix remains main's responsibility.
- Main reported black Fyne Canvas.Capture output despite correct OS screenshots in
  pilot2. Kept the capture command as a diagnostic helper and explicitly documented
  its limitation. OS ImageMagick import is the acceptance route; no Fyne capture
  investigation or screenshot-based success claim was added.

### Commands and cases for main

Build from SDL/go:

```sh
go build -tags desktop -o ../../wci1-fixture-native ./examples/collections/cmd/native
```

Run `wci1-fixture-native --nested`, `--empty`, or both. Feed stdin commands:

```text
state
pending
resize 900 600
reorder list
shrink list 5
reset list
hide tree
show tree
disable tree
enable tree
fail profile
fail provider
fail binding
fail layout
fail resource
fail guard
fail stale
reload
cancel tree
complete KEY success
complete KEY error
complete KEY empty
complete KEY invalid
close
```

Replace KEY with the exact started load/pending key; do not replay the same consumed
key or choose it with timing sleeps. Drive normal selection/activation/scrolling,
focus, disclosure, retry and cancellation with actual XTest events. `fail resource`
uses the existing resource-preparation hook, not a real allocator failure claim.
`fail stale` deliberately records a newer draft while its detached candidate is
pending, so old-bundle identity should stay fixed but that draft must survive.
No other failure/mutation command should execute a domain action.

### Verification and identity

From SDL/go, Go 1.27.1 linux/amd64:

```sh
go test ./examples/collections -count=1 -timeout=60s
go test -race ./examples/collections ./bridge -count=1 -timeout=60s
go build -tags desktop -o ../../wci1-fixture-native ./examples/collections/cmd/native
```

All PASS. Race results: collections 39.395s, bridge 1.435s. Desktop build completed
successfully after the shared-close/status-schema update. Scoped git diff --check
also passes. Product race tests do not compile the desktop-tagged main; desktop
build establishes compilation, while main owns external process/close evidence.

The new headless tests prove: source-defined nested geometry, stable selection on
reorder with old-target invalidation, deletion clearing selection and clamping,
reset/hide/show/disable/enable, finite resize and too-small last-valid preservation,
all six early/final injected failure stages preserving the published bundle,
newer-draft stale rejection and successful subsequent reload, strict command
rejection, exact request-key barriers, canceled late result rejection, explicit
error recovery and loaded-empty success without an automatic retry. Real SDL
activation/reload parity remains covered by the unchanged prior fixture test.
No native A08/A10/A11/A14/A15 completion claim follows from these tests.

Worker edits: fixture.go (Pending/strict Complete), new controls.go/controls_test.go,
new README.md, and cmd/native/main.go (variant/control wiring, logs, shared close).
fixture_test.go was retained unchanged and is listed only to identify the tested
fixture. The root binary is a local verification artifact, not a source delivery
or a file to stage. The sibling host candidate was still being integrated; a later
host change requires a new exact-candidate build/native identity.

SHA-256 at build/handoff:

```text
71cd65e6f32f3e419670908f8d1ad7494f11484d6d5b8f75bde4e76f067c4c7e  SDL/go/examples/collections/controls.go
dc7e1e2c7e6db5f7c778a6daa2f41ee1df3406b3aa7f29c626e5d62093229fce  SDL/go/examples/collections/controls_test.go
ebd31aee6ca726da58baf20093a8bb72e63fe332799abd7fb2ce2446f52a57fe  SDL/go/examples/collections/fixture.go
067212671843638b85784f788ae0549d37dfceebd584fc627e2953b4927b037f  SDL/go/examples/collections/fixture_test.go
0b993a80f81018ccedbc44c288f9fc75876644eb872ab5ff27be2e1c303327e6  SDL/go/examples/collections/cmd/native/main.go
bc9abff3deac206e71f0e14f11fcb849e2dd39021aee114d39876468a04f6930  SDL/go/examples/collections/README.md
fdda68a3f86e0bc29f8eccd66768e2002599875925907d2144c5cf73452052c5  wci1-fixture-native
```

Session0010 S2 payload: main requested bounded fixture instrumentation after
Noether's explicit handoff, then reported programmatic-close cleanup bypass and
Canvas.Capture limitation. Worker added nested source and deterministic commands,
fixed shared cleanup, ran the above tests/build, and handed command schema/cases
and binary hash to main. Actual native runs, remaining host findings, integrated
review and milestone acceptance remain pending. Fixture sources are held stable
pending main's native feedback; no unrequested next lane is started.

## Addendum — headless nested geometry sample for native harness

Main requested concrete nested geometry and command-result examples without host
edits. Ran an ephemeral read-only application probe from SDL/go:
`go run /tmp/wci1-geometry-probe.go`. It imports the current fixture/host, creates a
headless Fyne test app and eager NestedSource at 1000x650, then uses existing runtime
mutations. No product, fixture or host file was edited. This is headless measured
geometry for harness expectations, not an actual native run or a new binary build.

Initial geometry (x,y,width,height, logical units):

| Object | Rectangle | Effective clip | Content / range |
| --- | --- | --- | --- |
| page/body viewport | 0,26.1,1000,584.4 | same | content height 935.04; max Y 350.64 |
| page/body/entries outer Box | 505,26.1,495,935.04 | 505,26.1,495,584.4 | viewport parent page/body |
| page/body/entries inner viewport | 505,50.2875,483,898.8525 | 505,50.2875,483,560.2125 | content height 2338.125; max Y 1439.2725 |
| page/body/nav inner viewport | 0,50.2875,483,898.8525 | 0,50.2875,483,560.2125 | content height 72.5625; max Y 0 initially |
| page/footer/preview Box | 0,617,1000,33 | same | independent fixed sibling |

With entries at max offset, an additional +120 logical vertical delta leaves its
inner offset 1439.2725, moves body offset to 120 and consumes the full delta. Entries
outer Y becomes -93.9 and its inner viewport Y becomes -69.7125. Its effective clip
becomes y=26.1,height=584.4. Preview remains unchanged. Native harness should assert
these relationships against reported actual metrics with tolerance; font/driver
metrics are not frozen pixel fixtures.

Probe commands shrink list 5, reorder list, hide/show tree, disable/enable tree and
resize 900 600 all succeeded with unchanged bundle/model revision and zero actions.
Failures profile/provider/binding/layout/resource/guard/stale all returned errors,
kept the same bundle/model revision and executed zero actions. The stale case kept
its intentionally newer draft. Subsequent reload succeeded, changed the bundle and
advanced model revision 1 to 2, still with zero actions. Failure codes/messages were
version, collection-data invalid identity, unbound-module, native-minimum, injected
resource failure, injected final guard failure, and stale published state.

These numbers and outcomes were sent directly to main for its harness. No claim is
made that the complete command's acknowledgement proves provider delivery; the
separate request-key/load-return/state barrier rule still applies. Main retains
native proof and Session0010 S2 journal upkeep.


## Addendum — approved native gutter refinement implemented

The coordinator reported independent design approval and authorized this bounded
layout, layout-contract and report change. Governing authority is the final native
pilot refinement in SDP/04--Design/SDUI/Widgets/Collections.md, including the
reviewer's axis, corner, content-clip, single-cost and exhausted-space conditions.
The native overlap symptom and earlier headless coordinates above describe the
pre-gutter candidate; they are not current acceptance coordinates.

Implemented the optional ViewportMeasurer.MeasureViewport(*parser.Instance,
float64) (ViewportInsets,error), with Right/Bottom logical units and effective
font argument. Only 0.3 scrolling frames/groups consult it, once per owner per
layout run. Right requires declared vertical scrolling, Bottom horizontal;
reservation does not depend on range. Invalid, nonfinite, negative, excessive or
space-exhausting metrics reject before successful candidate geometry is returned.
Padding is applied first; gutters reduce child finite references and contribute
once to content-sized non-scroll outer dimensions. Absent adjunct and all 0.2
geometry retain their existing behavior, including unsupported-scroll diagnostics.

Viewport.HorizontalGutter and VerticalGutter are screen-space effective strips.
They exclude the bottom-right corner, intersect the owner allocation and ancestor
content/window clips, and do not intersect the owner's own content clip. Ancestor
offsets move them once; own offsets do not. Viewport.Rect supplies full unclipped
track length. Child clips exclude ancestor gutters. RouteScroll accepts the
owner's strips and retains inner-to-outer remainder routing. Collections derive
strips from existing CollectionMetrics.Viewport/outer allocation, excluding title,
without invoking the container adjunct or introducing another offset authority.
No source syntax, runtime types, host files, fixture files or renderer was added
or changed in this refinement. No branch switch or commit was made.

Verification on Go go1.27.1 linux/amd64, from SDUI/go:

- go test ./layout -count=1: PASS after initial gutter tests and the state-gate
  resize/shrink regression.
- go test -race ./layout ./svg ./prototype -count=1: PASS on final code/tests;
  layout 1.246s, svg 1.031s, prototype 1.047s.
- git diff --check -- SDUI/go/layout SDUI/docs/go-layout-contract.md: PASS.

Exact oracles cover padded dual-axis strips and excluded corner; nested frame/group
translation and clipping with independently routable gutters; zero-range reserved
space; content-derived width/height with one gutter charge; one adjunct call per
owner; invalid metrics/axis/exhaustion rejection; absent-adjunct and exact 0.2
compatibility; collection title/chrome/dual-axis strips and ancestor clipping;
remaining finite 1fr:3fr tracks; failed resize preserving runtime state and staged
geometry; successful resize/clamp; and hidden-content shrink retaining gutters
while clamping offsets to zero. These are synthetic metrics and state-gate tests,
not actual Fyne painting or external pointer evidence.

Candidate is the uncommitted shared-clone branch sdui/widgets-wci1 at base HEAD
d1c5c88d1989b173d610dad916f6f493f838ab7b. The scoped file hashes below identify this
refinement's resulting bytes, not an unchanged-commit result. Other lanes retain
ownership of their concurrent changes.

| File | SHA-256 |
| --- | --- |
| SDUI/go/layout/engine.go | 9fdbcc98f43c64ea1de72757dea24218088cfb8dc59c01e2dbc1fab2b1aa90f9 |
| SDUI/go/layout/types.go | c46914d323585a25990a5ed9449e9a9ad72eb51bc22acfd605db45c5bdc2f29a |
| SDUI/go/layout/insets.go | 3e10b7cfa619450f03b25c3e78963e2d209b63fcd92453b75c41332c80fdfd39 |
| SDUI/go/layout/insets_test.go | 2f8248f743b4b1f6e6cbcd6274e9b3071c75dc19595ff81a6e994f612d50f5b6 |
| SDUI/go/layout/viewport.go | d659123416b56a3a5187cf04359854b729978b0999dd4783166b8966361c31e9 |
| SDUI/go/layout/collection.go | 6cee39ba1ac958bc2c96b3d7c086401dfcb102e23ec24fc4f0bc91b3a09c0ced |
| SDUI/go/layout/scroll.go | e31f8436943df0a9c57e24858b1ccf7bff8ec9c337021eb6a53b269857cde410 |
| SDUI/docs/go-layout-contract.md | d6b7ae09e9e214f7ea321e42323e9408a6acffa8e4f8b87b75556d3c3eb20c22 |

API and passing results were sent directly to main and Noether. Noether owns the
native adjunct and use of explicit strips. Main must rebuild the native candidate
and repeat separated inner/outer thumb drags and wheel routing. No native binary
was rebuilt in this refinement; the earlier binary hash is historical. Native
proof remains pending, and Canvas.Capture remains a diagnostic helper with the
reported black-image limitation; acceptance uses OS capture.

Session0010 S2 handoff for coordinator-owned journal/roadmap/plan projections:
recovered current Session and refined stage design; applied SDP entrypoint 1.1.1,
Worker 2.0.0 and document-workflow within the exclusive lane. Owner input was
independent approval and code authorization with the listed gutter conditions.
Worker delivered scoped implementation, passing geometry/state-gate regression
checks, English contract reconciliation and this report. Main/Noether received
the stable API. Next step is native integration, actual separated drag evidence
and independent candidate review; no WCI1 completion or future-family claim.
