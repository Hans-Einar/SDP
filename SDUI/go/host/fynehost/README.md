# Native Fyne composition

`DocumentHost` composes a parsed SDUI document, typed collection providers, optional
connected bindings, runtime state and native presentation. It supports the bounded
0.2/0.3 admission capabilities in this package. The application owns source loading,
provider identity, domain binding and the Fyne window.

## Owner-goroutine lifecycle

Create and use the host on Fyne's UI owner goroutine with a finite positive intended
size. Mount `host.Container` as window content. The following fragment assumes an
application-created `DocumentRequest` named `request`:

```go
host := fynehost.NewDocumentHost(window.Canvas(), layout.Size{W: 1000, H: 700})
window.SetContent(host.Container)

candidate, err := host.Prepare(request)
if err != nil {
    return err // the previous published bundle remains in place
}
// Close the candidate if application logic abandons it here.
if err := host.Commit(candidate); err != nil {
    candidate.Close() // safe, including after Commit rejected/closed it
    return err
}
```

`Prepare` normalizes the supplied document itself, constructs a detached successor,
validates capabilities/bindings, and prepares measured native resources. It does not
execute domain actions or start provider loads. Do not mutate its document, provider
configuration or detached session while awaiting publication. `Commit` checks the
candidate, predecessor state, source identity, intended size and application guard
before swapping the prepared bundle. Post-publication reconciliation starts eligible
loads and restores focus. `Adopt(request)` is the immediate Prepare/Commit convenience.

Use `host.Mutate(func(s *runtime.Session) error { ... })` for application state changes
so accepted state is painted and request contexts reconciled. A handler may commit a
valid reentrant change before its outer dispatch returns an error; Mutate reconciles
that accepted change too. `OnChange` observes published state; `OnStatus` reports
errors. Keep these callbacks bounded and on the owner goroutine.

Container resizing invokes host layout. A rejected size retains the last valid
presentation measurement while the physical window clips it, allowing cancellation
and completion to continue. Call `host.Close()` during window teardown: runtime
acceptance is revoked before provider contexts are canceled. Close abandoned detached
bundles explicitly. Application-owned provider barriers/resources need their own
cleanup; host cancellation does not forcibly terminate a provider goroutine.

## Providers, binding and resources

- `DocumentRequest.Providers` maps normalized widget paths to typed
  `runtime.CollectionProvider` values. The application supplies stable item IDs,
  bounded valid data, provider ID/epoch, initial root readiness and optional `Load`.
  IDs are opaque domain identities; the host does not interpret them as file paths.
- `Load(ctx, request)` runs off the UI goroutine. Honor cancellation where possible;
  return data/errors without touching Fyne or the session. Host request tokens reject
  late/stale replies. `host.Post` delivers completion on the owner goroutine and
  defaults to `fyne.Do`. Tests may provide an explicitly drained owner queue; an
  inline worker callback is not a valid production replacement.
- Choose `preparation.Connected` with a `Bind` adapter that validates loaded modules
  and typed plans and installs handlers without executing actions. Preparation checks
  required bindings afterward. `preparation.Prototype` does not install connected
  domain bindings. Provider capability and native host capability are distinct.
- `PrepareResources` receives read-only prospective state with effective, clamped
  viewport offsets. It can run repeatedly during preparation and state gates. Do not
  mutate live state, invoke actions, or start loads there. The application owns any
  additional caches/resources it allocates and their cleanup; this hook has no
  resource-disposal return value.
- `Guard` is an application-owned, synchronous, side-effect-free identity recheck at
  preparation/publication. Capture and compare the actual source, provider epoch and
  external engine/module revisions relevant to the request. The host checks its own
  candidate identities but cannot discover external changes for the application.
  `SourceRevision` identifies the source; `Sequence` orders application candidates.

Runtime owns selection, expansion, requests and viewport offsets. Native collection
controls, clipping wrappers and viewport input surfaces project that state. They do
not maintain an independent scroll model. Existing Fyne inputs/buttons remain native;
collection metrics and rendered row labels share the same font and text definition.

## Legacy API and verification boundary

`NewRuntime(session, canvas)` / `RuntimeView` remains the legacy 0.2 adapter. It takes
an existing session, supports its existing reload controller and native draft/commit
behavior, and validates Markdown/native admission. It does not provide the document
bundle transaction, typed collection provider lifecycle or 0.3 collection admission.
Existing callers are not silently migrated to DocumentHost.

WCI1 collection/viewport behavior was delivered and reviewed before WCI2-M1.
M1 tabs/page/split code has targeted integration tests; its native acceptance is
separate and remains coordinator-owned. Headless tests do not prove OS input or painting.
The fixture's Canvas.Capture helper has produced black images in native pilots; use
actual OS screenshots for visual acceptance.

See the [runnable real SDL collection fixture](../../../../SDL/go/examples/collections/README.md),
[runtime contract](../../runtime/README.md),
[detached preparation API](../../preparation/prepare.go),
[preparation design](../../../../SDP/04--Design/SDUI/Widgets/Preparation.md) and
[collection design](../../../../SDP/04--Design/SDUI/Widgets/Collections.md).

## M1 panes and accepted preparation

DocumentHost admits tabs/page/split using `PaneCapabilities`; legacy RuntimeView
and standalone prototype callers retain their narrower capabilities. Symbolic page
icons currently reject with the missing `provider icon/1` diagnostic, including
hidden declarations. No implicit text-only fallback or M2 menu/dialog/command support
is advertised.

The runtime presentation gate probes geometry only. Finalized snapshots then enter
`PresentationPrepare`; its private ticket owns the native background/object list.
Only runtime's successful ticket Publish promotes accepted pending presentation.
Callback failure or sequence consumption alone cannot expose the speculative page.
An independently accepted reentrant draft retains its own ticket and is synchronized
even when the outer interaction fails. Application `Guard` remains a Prepare/Commit
identity check: it is not called from state preparation inside an executing domain
handler. The bridge owns its domain-result checks; resource hooks must remain pure.

The bounded header wrapper reuses AppTabs rendering and actual header rectangles,
adds one focus stop and arrows/Home/End, and sends selection through
DispatchInteraction before displaying accepted state. Tab enters remembered/first
eligible page content; Shift+Tab leaves the header. Inactive page bodies retain
runtime drafts/scroll and have no active native input. A bounded divider supplies
pointer drag and focused arrows/Home/End, Ctrl+Home/End collapse and Space restore.
Its thickness is measured from a native empty Fyne Split, but shared pane geometry
owns child extents and ratio clamps; native Split state cannot override them.

`Bundle.Inspect()` retains snapshot/rows/widgets/viewports and adds diagnostic maps:
`tabs[path]` contains header/clip/body/selected and per-page id/label/enabled/rect/clip;
`splits[path]` contains shared Divider/DividerClip/First/Second and split geometry;
`controls[path]` contains mounted rect/clip/visible for inputs, widgets, headers and
dividers. Rectangles use logical screen coordinates. Hidden control rectangles can
be retained from their last mount; consult visible before using them as input targets.
These are inspection projections, not another model mutation API.

See the [real SDL panes fixture](../../../../SDL/go/examples/panes/README.md),
[shared pane layout API](../../../../SDUI/docs/wci2-layout-api.md) and
[reviewed pane contract](../../../../SDP/04--Design/SDUI/Widgets/Panes-and-commands.md).
