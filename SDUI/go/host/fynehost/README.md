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

DocumentHost admits tabs/page/split and the selected M2 controls using
`CommandCapabilities` (extended by WCI3 `FieldCapabilities`); legacy RuntimeView and standalone prototype callers retain
their narrower capabilities. Symbolic icons require `DocumentRequest.Icons`,
including icons in hidden/closed declarations. Missing or unsupported resources
reject preparation. Prepared bytes retain the native theme-color contract.

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


## M2 commands and auxiliary surfaces

Use `DocumentRequest.Icons` for symbolic resources and `Bind` for typed command
and dialog interaction handlers. Connected admission rejects ordinary commands
without a handler. Prototype mode reports unbound declarations; toggle-only and
built-in effects remain local runtime operations. Plain buttons retain legacy
Activate dispatch, including icon/tooltip-only decoration. Explicit command fields
select exactly one `DispatchInteraction` route shared by button/menu/key.

Native buttons project accepted checked/exclusive state; tooltips use the current
canvas content without adding an input-stealing overlay. Resource bytes are owned
by the bundle. Native key normalization rejects duplicate/reserved keys before
publication, including hidden declarations. Fyne sends shortcuts directly to
focused Shortcutable controls and sends Shift-only/function keys through TypedKey;
the bounded adapters forward unhandled declared keys to the same runtime command
capture. Entry editing and collection/divider navigation take precedence.

Menus retain Fyne Menu/PopUpMenu items and their native hover/scroll tree. A public
pointer/keyboard wrapper authorizes only the current synchronous selection input.
Dismiss hides/removes that opening before Action, which may claim its capture once.
Escape/outside revoke immediately; late dismiss cannot remove a newer overlay.
Caller menu models are cloned before decorating Actions. Context capture preserves
selection and rejects stale item/widget/state revisions through runtime validation.

Modal surfaces use CustomWithoutButtons with ordinary composed controls and a
bounded, visible status line. Shared geometry describes the content; native chrome
is measured separately, including prospective parent resize. Empty dialogs have
one native focus stop for the dialog itself. Nonmodal surfaces use `app.NewWindow`
and independent canvases; the first outer size includes actual measured chrome,
then accepted content resize is authoritative. No OS transient/always-on-top claim
is made. RequestFocus is used for explicit opening/child gestures and remains
platform-dependent (Wayland may ignore it).

The application composition root must call `NativeParentHidden()` before hiding
its owner and `NativeParentClosed()` before closing it. Both return errors and
revoke published descendant generations before native teardown. Host Close/reload
also revoke children. WM/decoration close is native owner loss: its Close interceptor calls exact
`RevokeSurface(target, "parent-closed")` before teardown, including modal descendants.
Its receipts have reason `parent-closed` and sequence zero, preserve any prior domain
outcome and never invoke Accept. A child modal or failing presentation gate cannot
veto this owner teardown. The idempotent OnClosed fallback uses the same classification
for direct Window.Close. Source Close/Cancel remain user interactions with the
triggering sequence and modal guards; they are not relabeled as native teardown. Publication is acknowledged only after native Show;
abandoned preparation emits no terminal result. `OnDialogResult` observes each
published opening's terminal outcome once; observer panic is isolated, and an
observer may reenter without replaying a drained result. It is not an Accept handler.

`InteractionError` wraps the original error with status/domain/sequence and preserves
errors.Is/As. Failed or blocked dialog acceptance keeps drafts and a visible compact
status; full outcome/message remains in runtime. Cancel/Close remain usable.

Inspect adds actual canvas-local controls, visible menu rows and surface geometry.
A control's `canvas` is the literal `"main"` for the main canvas. A nonmodal
dialog path identifies its independent canvas; a modal uses its actual parent
canvas ID. Consult visibility/clip and actual window title before native
input. Inspect reports geometry; it is not a paint-completion acknowledgement.

See the [real SDL commands fixture](../../../../SDL/go/examples/commands/README.md),
[layout canvas contract](../../../docs/wci2-m2-layout-api.md) and
[runtime M2 contract](../../runtime/README.md). Coordinator-owned native acceptance
and independent review remain separate from these implementation tests.

## WCI3-M1 typed scalar controls

The document host adds `checkbox`, `slider`, `select` and `number` through
`FieldCapabilities`. Legacy `.2` RuntimeView, public SVG and standalone prototype
admission retain their narrower contracts. Extended input/text/IME is not part of
this milestone.

Supply `DocumentRequest.Choices` as `map[string][]runtime.ChoiceOption`, keyed by
exact normalized select paths, including hidden and closed-dialog declarations.
Every select needs an entry; an explicitly empty slice differs from an absent
provider. IDs identify options; duplicate display labels remain distinct. Runtime
owns the copied inventory, accepted/proposed values, validation, revisions and
option generations. `BindChoices`/`SuccessorWithChoices` establish readiness before
native resources, and connected binding validation/postchecks also precede the
final native preparation hook. Widget, host and choice-provider capabilities are
separate requirements.

Use the runtime `ObserveChanges`/`ValidateFieldWith` and typed handler APIs from the
application binding adapter. A user Change publishes a proposal before Commit.
Checkbox and choice gestures commit their original returned `FieldChange.Field.Target`;
a reentrant observer cannot cause the host to recapture and authorize a different
automatic Commit. Slider release/key-up and number Enter explicitly capture the
current field. Ordinary synchronization and checked application `Apply` emit no user
Change/Commit. Failed native preparation preserves the previous presentation and
restores the native number entry's raw projection. Runtime's existing private
accepted-preparation ticket remains the sole resource publication authority.

The bounded adapters retain Fyne Check/Slider/Select/Entry painting and selection.
Fyne Select identifies choices by display label and keeps its popup private, so the
host uses the existing native menu adapter with exact captured `OptionTarget`s.
Disabled entries cannot activate; replacing the option generation retires its popup,
including an empty opening. Escape/outside dismissal does not propose a choice.
The popup retains keyboard focus through ordinary synchronization.

Fyne Slider's native numeric step rounding uses a zero origin and floating-point
arithmetic. The adapter uses its normalized visual track only; shared `numeric.Grid`
and runtime `EditTick` determine each legal value. Arrow/Home/End input changes a
proposal while held and commits once on key-up. Pointer movement behaves similarly
until release; Escape reverts and suppresses the remaining drag. Sync retains the
same native object and gesture. Read-only controls remain focusable and refuse user
mutation. The number Entry retains selection/copy; typing, paste/cut, stepping and
Commit are blocked when read-only. Invalid raw numbers remain visible with fixed,
ellipsized feedback and cannot be stepped. Full feedback remains in runtime state.
Number step affordances reuse Button chrome without additional Tab stops; their
native renderer owner receives the actual allocated size. Tab/blur never commits.

`Bundle.Inspect()` adds `fields[path]`: `kind`, `canvas`, `title`, `visible`, `clip`,
`label`, `control`, `feedback`, `entry`, `decrement`, `increment`, `slidertrack` and
`sliderthumb`. Parts are actual native rectangles in the named canvas's logical
coordinates; absent parts are zero rectangles. Consult visibility and clip before
sending input. An open `choices[path]` contains its canvas/title/generation and
`items[]` with `id`, `label`, `enabled`, `rect`, `clip`, and exact captured `target`.
`Snapshot.Fields` is the authoritative typed state; these diagnostics add no model
or command API. The main canvas ID is `main`; modal controls share their actual
parent canvas and nonmodal controls identify their own window.

The targeted host tests cover rejected resource tickets, reentrant Change, held
slider gestures, exact large-origin stepping, invalid raw numbers, read-only copy,
option identity/generation, empty inventory, native step hit rectangles and fixed
feedback at several fonts. Native OS input/painting and whole-stage acceptance are
coordinator-owned; the integration worker report records the tested candidate.

See the [real SDL values fixture](../../../../SDL/go/examples/values/README.md),
[scalar layout API](../../../docs/wci3-m1-layout-api.md),
[runtime contract](../../runtime/README.md), and
[values/text design](../../../../SDP/04--Design/SDUI/Widgets/Values-and-text.md).

## WCI3-M2 extended native text

`DocumentHost` uses `TextCapabilities` and a retained native Entry for inputs
that explicitly supply any of `multiline`, `readOnly`, `placeholder` or `required`.
Explicit false/empty values opt in too. Inputs without those arguments keep
legacy `.2`/basic `.3` behavior. The frontend `InputOptions` helper selects this
policy; `FieldState.Input` supplies copied runtime metadata. Runtime Widget
Value/Draft remain the sole accepted/proposed text store.

Multiline Entry uses native word wrapping with three visible minimum rows.
Shared `FieldMeasurer` assigns the whole Entry, including native scroll chrome,
plus a separate optional label and fixed feedback row. Empty labels allocate no
label region. Detached probes use an empty, identically configured Entry: its
native scrolling minimum depends on themed character metrics and visible rows,
not draft extent. Equivalence tests compare real short/long Unicode editors at
multiple fonts, avoiding repeated full-text shaping during preparation. Live
editors still receive their exact text. Long text stays inside the assigned finite editor; Entry owns its
caret, selection, undo and internal scroll. Ordinary publication preserves the
observed native scroll offset using public `ScrollToOffset`, which clamps after
resize. Text-area wheel input uses the native subtree and consumes at its limits;
outer gutters/background retain their separate WCI1 routing. No Entry offset is
stored in runtime Viewports.

Native OnChanged calls checked `EditField` once. Enter commits single-line text;
multiline Enter/Shift+Enter inserts a newline and exact Primary+Enter commits.
Tab/Shift+Tab follow source/surface order; blur never commits. Native undo/redo
changes the draft only. Readonly remains focusable/selectable/copyable, with
typing, cut/paste, deletion, undo/redo and Commit blocked. A guarded editing menu
keeps late actions subject to current eligibility; declared context menus retain
the existing WCI2 route. Single-line keyboard/context Paste rejects CR or LF
before native mutation. Clipboard content is read once and that exact validated
snapshot is delegated, preventing Fyne's single-line LF-to-space conversion.

Identical displayed bytes preserve history, including self-echo acceptance and
explicit Apply. Different programmatic text uses one muted SetText and resets
history. Page hide/show, collapse and failed reload retain the Entry. Successful
compatible reload may recreate it under the runtime retention rules.

Fyne edits text/history before OnChanged. Under the reviewed exception, an actual
native edit rejected by EditField or publication restores the **current** runtime
draft muted on the same Entry and clears undo/redo; caret/selection/scroll may
reset. Reentrant accepted work is never overwritten with an older draft. This
exception does not apply to invalid-but-retained drafts, failed Commit/reload/probe
or pre-delegation refusal. Those paths preserve their existing history.

For extended inputs, `Inspect().fields[path]` includes actual Entry `text`,
`placeholder`, `readOnly`, `focused`, `cursorRow`, `cursorColumn`, `selectedText`
and `multiline`, alongside existing canvas/title/control/entry/label/feedback
rectangles. When its public renderer exposes a native Scroll, `scrollRect` and
`scrollOffset` report that object; offset uses Size W=X/H=Y. These are diagnostic
observations, not another state model. Consult control visibility and clip.

IME preedit and consumed Return/Escape belong to the selected pinned GLFW driver
filter, maintained by the coordinator. Host adds no composition timer or guessed
preedit state. Component tests cover editing, history, lifecycle and native object
routing; actual configured OS IME, SDL and visual acceptance require the main
native harness on the exact integrated candidate. See the
[text layout API](../../../docs/wci3-m2-layout-api.md) and
[real SDL text fixture](../../../../SDL/go/examples/text/README.md).
