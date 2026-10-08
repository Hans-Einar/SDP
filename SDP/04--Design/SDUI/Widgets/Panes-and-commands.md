# WCI2 panes and commands — reviewed stage design

| Field | Value |
| --- | --- |
| Assignment | SDP Architect; future WCI2-M1/M2 under PLAN-SDP-0022 |
| Status | Reviewed; WCI1 completed; WCI2-M1 delivered; M2 selected after independent revised contract review |
| Authority | KB-SDUI-003 full inventory; owner’s bounded design assignment, 2026-10-08 |
| Parents | [Design](Design.md), [Acceptance](Acceptance.md), reviewed [Collections](Collections.md), [Plan](../../../05--Implementation/SDUI/Widgets/Plan.md) |
| Obligations | SDUI-R05/R12/R15–R18/R23/R26–R28; GAP-XFMD-SDUI-002/004/007/008 |
| Continuity | [Session0010](../../../Sessions/session-%230010--SDUI_widgets.md), S3 active; S2 completed |

## 1. Outcome, observations and choices

WCI2-M1 adds tabs/splits; WCI2-M2 adds shared button/menu/key commands and composed modal/nonmodal dialogs. Preserve page drafts,
distinguish Cancel from Accept and reject stale context actions.

Inspected `b61a3de` plus working changes: parser has Rows/Arguments/provenance but no call bodies; runtime has handles/drafts/Apply;
Fyne mounts button/input. Clone memo `/tmp/sdp-sdui-widgets/WCI1-runtime-API.md` proposes Snapshot/StateGate/Successor; reconcile
with delivered WCI1. No implementation tests run here.

Choose call bodies using existing brackets/Rows over new delimiters, nested component arguments or native objects. Local Fyne v2.8.1
has AppTabs, Split, menu state and custom dialog content. Prefer these with bounded adapters. Split lacks a public change callback
and its Offset differs from content fraction at minima; Button lacks checked/tooltip fields. Nonmodal adapter selection is complete
in §5. Native conformance and the other bounded adapters still require evidence.

## 2. Proposed canonical syntax and AST

Extend the still-unreleased exact `sdui 0.3;` profile; preserve 0.2 grammar, output, fixtures and rejection behavior. Within the
existing `primary` production replace `widget-call` with `call-primary` for 0.3 only:

```ebnf
call-primary = widget-call, [ "[", [ rows ], "]" ] ;
```

Only tabs/page/split/menu/menuGroup/dialog require and accept bodies; formatting follows the body. Keep existing tokens, argument
values, reuse and resource limits. Body assignments are names, not frame regions; nest an ordinary frame for header/body/footer. No
expressions, imports or new source language.

```text
sdui 0.3;
ref: actions "actions.sdl";
Main = [
  details=command("Details", effect="open", target="settings");
  compact=command("Compact", toggle=true, checked=false, key="Primary+K");
  tools=<button(command="details"), button(command="compact")>;
  file=menu("View")[item(command="details"); separator(); item(command="compact")];
  panes=split(axis="horizontal", proportion=0.35, minFirst=0.15, minSecond=0.2)[
    navigation=tree("Navigation");
    workspace=tabs("Workspace", selected="overview", callback=actions.Page.@invoke)[
      overview=page("Overview")["Overview content"];
      notes=page("Notes")["Notes content"]]
  ];
  preview=input("Preview", value="");
  settings=dialog("Settings", modal=true, callback=actions.Save.@invoke)[
    name=input("Name", value="");
    actions=<button("OK", effect="accept"), button("Cancel", effect="cancel")>
  ]
];
actions.Page.setHandle(Main.preview);
```

Proposed syntax, not current-parser evidence. Applications supply tree providers and §5 binding plans. Page returns preview text;
Save uses the typed acceptance result, with no setHandle receiver.

| Form / arguments (schema notation) | Children, defaults and validation |
| --- | --- |
| `tabs(label, selected?: string, callback?: member-ref)` | Named when callback supplied; callback handles ActivatePage only. One or more named pages, one per row; selected is a direct page name; default first eligible. Unknown/ineligible initial selection rejects. |
| `page(label, icon?: string)` | Only under tabs after reuse expansion; ordinary content rows, possibly empty. Name is stable page ID; never label/index. |
| `split(axis, proportion=0.5, minFirst=0, minSecond=0, collapsible=true)` | Exactly two named children, one per row; horizontal means left/right, vertical top/bottom. All numeric arguments finite; minima >=0, sum <1; initial proportion within [minFirst,1-minSecond]. |
| `menu(label, mode="bar", target?: string)` | Rows each contain item, separator, menuGroup or nested menu. mode is bar/context/submenu; nested menus must be submenu. Context mode requires a tree/list/widget target; other modes forbid target. |
| `menuGroup(label)` / `item(command: string)` / `separator()` | Group is a non-invokable heading with menu children; item refers to exactly one command; separator has no action/state. Empty menus/groups allowed; useless separators suppressed in native presentation only. |
| `dialog(label, modal=true, callback?: member-ref)` | Named, initially closed auxiliary surface with ordinary composed content. callback handles Accept and requires the §5 acceptance result contract. Cancel/Close are observable terminal results, never acceptance callbacks. |
| `command(label, ...)` / extended `button(label?, ...)` | Shared schema below. Command declaration must be named and has no geometry. Button may omit label only with command reference. |

Require nonempty accessible labels and at most one positional argument before named; reject unknown, duplicate/wrong-type fields.
enabled/visible remain boolean formatting; tooltip is text, icon a symbolic resource ID. Pane/dialog formatting stays relative;
split owns child axis extents. Commands/context menus accept state formatting only; bar menus stay at their content position. Modal
size references parent content; nonmodal initial size references parent host area, then its actual window content area.

Body forms use Kind `composition`, Widget naming the form, typed Arguments and Rows; leaves stay `widget`. Instance preserves
Profile/provenance; auxiliary nodes remain in walks, outside layout tracks. Validate placement after reuse. Keep AST fields/tags,
`sdui-ast/0.3`, `sdui-go-model/2` and source SHA; consumers check kinds/capabilities as well as profile. Never emit WCI2 in 0.2 or
infer profile from spelling.

Command/target paths are relative to the innermost definition instance, or leading `/` to selected entry. Reject `..`, indexes and
anonymous `$` segments; resolve named segments to one typed target after expansion. Reuse has independent local state; absolute
references share deliberately. Preserve scope/span/use chain in generated constructors rather than guessing scope from parent paths.
Frontend resolves these identities during normalization; runtime/bridge consume them rather than inventing separate lexical-scope
resolvers.

## 3. Commands, buttons, keys and context

Fields: label/icon/tooltip, toggle=false, checked=false, exclusive?:string, key?:string, context="none", callback?:member-ref,
effect?:string, target?:string. Context is none/widget/item. checked/exclusive require toggle; exclusive groups are
definition-instance scoped, with at most one checked (zero allowed). Initial conflicts reject; activating a checked exclusive member
is a no-op, another clears peers atomically. Ordinary toggles invert once. Widget/item context requires target; button/key captures
its widget/selected item (missing selection disables). Menu captures its own row from the same collection. None forbids target
except effect open.

A button owns an implicit command or references a shared command. Referring buttons may override only label/icon/tooltip;
behavior/callback/key overrides reject. Local enabled/visible gates can restrict, never enable, the command. Reject read-only on
commands; programmatic changes never emit user actions.

Compatibility boundary: existing basic buttons retain Activate/Handler dispatch,
including existing 0.3 sources. Explicit M2 behavioral arguments (`command`,
`toggle`, `checked`, `exclusive`, `key`, `context`, `target`, `effect`) select one
implicit/shared command and InvokeCommand/InteractionHandler instead. Icon or
tooltip alone does not promote dispatch. Never register or invoke both routes
for one button; frontend preserves the distinction through normalization/codegen.

`effect` is a closed local enum open/accept/cancel/close. Open requires a dialog target; the other effects require a button/command
belonging to an enclosing dialog and forbid target. Effect and callback are mutually exclusive; effects forbid toggle and require
context none. Without either, a toggle changes local state; an ordinary command is explicitly unbound and cannot establish connected
readiness. Symbolic callbacks still use declared SDL aliases and `@invoke`, with no domain action during preflight.

Dispatch checks the live command handle/model/state, origin surface and effective enabled/visible state before calling a handler.
Closed menu items may share an enabled command with keys; menu openness is not command visibility. Hidden-page or closed-dialog
commands cannot run through keys. A modal surface blocks commands outside its subtree. Toolbar/menu/key use one InvokeCommand
dispatch and the same command handler. Stage proposed checked/exclusive changes with returned updates; do not toggle in native
OnTapped first. Local effects resolve internally: Accept calls the enclosing dialog handler, never the referring button and dialog
both. Keep the original sequence through this derivation; no recursive Dispatch/new SDL action. §5 defines preflight, reply
publication and observable outcomes, including domain success/UI failure.

Key: optional ordered Primary/Ctrl/Alt/Shift joined by `+`, then A–Z/0–9/F1–F12; letters/digits need a modifier. Reject duplicate
modifiers, Primary+Ctrl, same-surface duplicate commands and host-reserved keys after platform Primary normalization.
Editing/navigation wins; inactive surfaces receive no keys. No global/chord shortcuts. Tooltip appears on hover/focus; labels remain
accessible. Missing icons reject unless an explicit reported text-only fallback was selected before activation.

Secondary click or Menu/Shift+F10 captures clicked/focused context without selection: widget handle/revision or WCI1
CollectionTarget. Blank space permits widget context only; item context requires a selectable live row. Reject
stale/hidden/disabled/group/separator targets; revalidate at invoke/result, never retarget current selection. Menus use
arrows/Home/End/Enter/Space/Escape; dismissal executes nothing, submenus restore parent focus. Capture expected
Event.StateRevision with the context after menu-opening publication; never restamp saved context at invocation. Any
intervening accepted state change makes it stale: dismiss without callback, including hide/reveal with unchanged handles.

Selection-triggered dismissal retains the original root-menu capture until its
single command dispatch validates it. Native menu hiding must not publish a
CloseMenu before that dispatch or restamp its context. Successful dispatch stages
root dismissal with the accepted command update. After any rejected/stale/failed
invocation, dismiss and revoke that opening without retry or additional action;
retain any reported domain outcome. Escape/outside dismissal revokes immediately
without invoking a command. Duplicate native dismissal callbacks are harmless.
Pinned Fyne menuItem.trigger calls parent.Dismiss before Item.Action. Its adapter
must therefore distinguish the synchronous selected-item callback from dismissal
cleanup using a bounded per-opening input scope around native pointer activation
and Enter/Space forwarding. OnDismiss hides immediately; only within that scope
does it retain capture for one synchronous Item.Action. Scope exit revokes any
unclaimed dismissed opening before forwarding returns. Escape/outside dismissal
outside this scope revokes immediately. Exact-opening checks prevent duplicate
cleanup from closing a replacement. No queue/timer heuristic, revision restamping,
second invocation or general receipt registry. The native adapter must prove this
ordering for nested pointer and keyboard selection before product acceptance.

## 4. Pane, surface and draft state

Tabs retain hidden drafts/scroll/focus; only selected eligible content lays out/receives input. Hiding revokes
transient/context/load targets under WCI1 rules. Arrows/Home/End select; header keeps focus, Tab enters remembered valid or first
content focus. Selected removal/hide/disable picks first eligible declaration, or empty body. Invalid content focus returns to
surviving header; canceled tokens never revive.

Pointer or header-key selection of a different eligible page dispatches ActivatePage to the tabs callback, carrying old/new stable
IDs and the new page handle. Same-page activation is a no-op. Stage selection/focus and callback updates together;
handler/preparation failure retains the previous page/drafts/focus. Programmatic selection, initial selection, removal fallback and
reload are silent. Eligibility checks use the page’s own/ancestor enabled/visible intent, not its current unselected-content
invisibility. Omitted callback means local navigation; supplied but unbound prototype callback is reported, never connected-ready.

Let usable axis U = allocated extent minus measured divider, U>0; first-content fraction is p. Expanded minima are a=max(minFirst*U,
measuredFirst), b=max(minSecond*U, measuredSecond). Require a+b<=U; clamp user drag/resize to [a/U,1-b/U]; invalid programmatic
values reject. Shared layout owns these bounds. Focused-divider arrows adjust 0.05; Home/End reach legal limits; Ctrl+Home/End
collapse first/second; Space restores. Collapse is an explicit none/first/second state, never inferred from p=0 or 1. With
collapsible=true, the collapsed child has zero extent and BOTH minima are suspended and it leaves measurement/layout/hit testing;
the visible child gets U and meets its own minimum. Divider/restore affordance stays visible/focusable. Save the last expanded
proportion; preserve child drafts/scroll, revoke hidden targets. Resize while collapsed keeps zero extent and the saved proportion.
Restore clamps that saved proportion against BOTH current minima; if they cannot fit, remain collapsed with diagnostic, without
partial focus or state changes. A too-small visible side likewise rejects geometry. Other layout failure retains the last valid
presentation. Compatible reload retains state; axis/kind changes reset it; changing collapsible to false requires successful
expanded preparation. Native Split minimum clamping cannot override this rule.

One opening per dialog: fresh surface generation, parent handle/generation, model revision and context. Reopen focuses without draft
reset. Modal traps subtree input/focus; nonmodal permits parent interaction. Parent closure/removal/hide/disposal revokes
descendants before teardown; nested focus follows parent stack. Successful reload closes transients without Accept; failure
preserves the entire bundle/drafts.

Opening prepares ordinary descendants/resources and focus before making the surface visible. Initial focus is first eligible control
in source order (dialog itself if empty); Tab/Shift+Tab cycle in modal scope. Explicit accept/cancel/close buttons use local
effects. WCI2 proves Escape consumes one level: menu dismissal, then dirty single-line input draft revert, then clean-surface
Cancel. Single-line Enter commits; focused-button Enter activates that button, never implicit whole-dialog Accept. Never cascade one
key through multiple levels. Future WCI3 compatibility only: choice popup joins menu priority, active IME precedes draft revert;
multiline Enter inserts newline and Primary+Enter commits. Native-consumed IME keys must not dispatch SDL. Implementing and proving
new IME/typed-control/multiline behavior belongs to WCI3 evidence, not the WCI2 gate. Window/chrome dismissal is Close. Cancel/Close
discard drafts but retain distinct result kinds. Restore opener focus if still valid, else nearest live parent’s first focusable
control; never steal focus from another active nonmodal surface when closing an unfocused one.

Accept captures all input drafts owned by this open surface (including hidden pages, excluding nested dialogs), with exact
handles/value/draft revisions. Validate before calling its acceptance adapter. Accepted=false/error keeps it open with drafts and
visible error; Accepted=true commits captured drafts and returned updates atomically, then closes once. With no callback/Go binding,
local acceptance is explicit and makes no persistence claim. A declared but unbound callback cannot silently accept as connected.
Cancel/Close reset unaccepted drafts to current accepted values (not opening-time values), invoke no acceptance handler, and publish
their distinct terminal result. Existing explicit child Commit actions are not rolled back; atomic forms bind persistence only to
Accept. Parent closure/hide, reload and disposal yield Close with the precise lifecycle reason. Programmatic close is also
observable but emits no user/domain action. App observers receive a detached DialogResult after close publication, once per opening,
including the old token and accepted fields only for Accept. Rejected Accept emits no terminal result. Queue notifications outside
state mutation; observer errors cannot veto/undo close. Teardown must deliver the notification without reviving obsolete handlers.

## 5. Interfaces, preparation and native boundary

```go
type ContextTarget struct { Widget Handle; ModelRevision uint64; Item *CollectionTarget }
type SurfaceTarget struct { Handle Handle; ModelRevision, OpenGeneration uint64 }
type DraftField struct { Handle Handle; Value Value; ValueRevision, DraftRevision uint64 } // WCI2 String input draft
type PageActivation struct { PreviousID, PageID string; Page Handle }
type SplitChange struct { Operation string; Proportion float64 } // ratio/collapse-first/collapse-second/restore
type CommandInvocation struct { Origin Handle; Via string; Surface *SurfaceTarget; Context *ContextTarget; Checked *bool }
type DialogRequest struct { Surface SurfaceTarget; Fields []DraftField } // Fields only on Accept
type AcceptDecision struct { Accepted bool; Message string }
type DomainOutcome string // not-called/succeeded/rejected/unknown; empty is unspecified
type InteractionReply struct { Updates []Update; Accept *AcceptDecision; Domain DomainOutcome }
type InteractionHandler func(Event) (InteractionReply, error)
type InteractionResult struct { Sequence uint64; Status string; Domain DomainOutcome }
type DialogResult struct { Surface SurfaceTarget; Sequence, AcceptSequence uint64; Kind, Reason string; Domain DomainOutcome; Fields []DraftField }
// DispatchInteraction(Event) (InteractionResult,error); legacy Dispatch remains compatible.
// Composition root supplies OnDialogResult func(DialogResult); observers return no updates.
```

WCI2 Event retains Handle/ModelRevision/Sequence and adds expected StateRevision plus exactly one pointer payload: PageActivation
for ActivatePage (Handle=tabs), SplitChange for AdjustSplit (split), CommandInvocation for InvokeCommand (canonical command/implicit
button), DialogRequest for Accept/Cancel/Close (dialog). Legacy button/input/collection schemas remain unchanged. Reject unknown
kinds/enums, wrong owner, extra/missing payloads, nonfinite ratios, stale revisions/generations and duplicate sequences. Via is
button/menu/key; Origin must be a current referring control (command itself for key). Checked is required only for toggles and must
equal the runtime-computed proposal. Context exactly matches the declared kind; item targets must reference the same widget. Surface
must match Origin’s active dialog, otherwise nil. PageID must match Page’s named direct child; PreviousID matches current selection.
Ratio operations alone carry Proportion; nonratio requires zero. Accept fields are runtime-captured, bounded by the existing
256-update batch limit; external copies must match exactly. Captured fields plus returned Updates share that total
256-write limit; no overlap or implicit limit increase. Before execution reject oversized captures and any statically
known mandatory update budget exceeding the combined limit. Arbitrary custom-handler outputs are unknowable before
execution: validate them after reply against the remaining 256 minus captured-field count. Oversize is ui-conflict with
retained domain outcome/replay block; never truncate or replay. No budget-registration API. DialogAcceptResult returns zero updates. WCI2 captures the existing input’s raw draft as String Value. Explicitly
require valid UTF-8 and <=32768 bytes at this WCI2 capture boundary: current legacy Draft/String validation checks length only. Do
not strengthen legacy 0.2 APIs or require new typed-control state. Future compatibility only: the [WCI3 draft](Values-and-text.md)
owns RawDraft/FieldValidation/OptionTarget extensions, typed proposed values and option-generation/eligibility checks. Invalid
numeric drafts cannot become accepted values; select requires its dedicated option validator, never a hidden list. Those extensions
and evidence are WCI3 work.

BindInteraction(Handle, InteractionHandler) binds the callback-owning tabs/command/dialog, never each presentation. Composition
callback-owner enumeration extends bridge discovery without changing legacy Widgets(). Reply.Accept is required for a successful
Accept reply; Accepted=false forbids Updates. Other kinds require nil Accept. Split is local; effects have no command handler
(Accept derives one dialog call). Preflight validates identity/signature/geometry, then consumes one sequence before calling the
synchronous handler once. Reentrant interaction dispatch rejects. Capture all current widget handles/value/draft revisions before
execution, bounded by existing tree limits; returned Update targets are not known beforehand. Bridge also captures its known
receivers and engine revision. Immediately before Apply validate actual returned targets against that capture, plus engine/bundle
identity and the internal post-sequence-consumption StateRevision baseline. Update has only ExpectedValueRevision: compare
DraftRevision explicitly; AcceptDraft/equal text is insufficient. Reentrant legacy mutations remain accepted but make the outer
reply stale. Never roll them back. Reject overlap with reserved selection/toggle/captured-draft changes. Final gate, resource
preparation and atomic publication use the same finalized state; no yield after last check; native sync muted.

Result.Status is committed/rejected/ui-conflict. Read Reply.Domain even when error is nonnil; discard Updates/Accept on error. Empty
Domain after a called handler becomes unknown, never inferred success/no-execution from an error. Explicit not-called requires proof
execution did not start. SDL adapter sets succeeded after a valid successful Execute result (DialogAcceptResult true), rejected for
a valid false decision, unknown for execution error/malformed output. Subsequent engine/target/revision failure preserves that
annotation: successful domain execution plus failed UI publication is ui-conflict/succeeded. Local success is committed/not-called;
false Accept is rejected/rejected. Reject inconsistent decision/outcome pairs before UI publication; retain unknown for an
indeterminate executed attempt.

Keep the attempt/outcome and visibly report domain success with UI failure, or an unknown domain outcome. Either blocks further
Accept for this opening. Editing/revert and token-based Cancel/Close remain available without obsolete field-revision checks. No
later reply may publish the stale batch. Closing discards unaccepted drafts only; its result retains the attempt/outcome and never
claims domain rollback. Application owns reconciliation and authoritative reload/reopen; neither automatically replays the action.
No ResolveAccept or mandatory general reconciliation API.

Snapshot adds detached Tabs/Splits/Commands/Surfaces maps keyed by exact normalized InstancePath. Tabs retain owner/page handles,
selected stable ID and per-page focus; splits retain axis/proportion/collapsed side/saved expanded proportion; surfaces retain
parent/opener tokens, opening generation and acceptance outcome. Retain declaration enabled/visible intent separately from effective
activity. Recompute descendants from intent plus selected/collapsed/open ancestors; unselected pages never acquire visible=false
intent. Header/divider focus is explicit, not a fabricated input widget. Tabs(Handle)/Split(Handle) read state;
SelectPage(Handle,string) is silent; SetSplitProportion(Handle,float64) is strict;
CollapseSplit(Handle,SplitSide)/RestoreSplit(Handle) use §4. These are bounded methods, not a reducer framework.

```go
type SplitGeometry struct { Lower, Upper, Effective float64 } // expanded content fractions
type PresentationState struct { Viewports map[string]ViewportState; Splits map[string]SplitGeometry }
type PresentationGate func(Snapshot) (PresentationState, error)
// CheckPresentationWith(PresentationGate) error; exactly one effective gate.
```

CheckStateWith keeps its WCI1 signature through an adapter, not a second gate. Split-bearing models require the typed gate. Layout
returns finite legal bounds/effective ratios for every active expanded split; missing/extra/invalid entries reject. Collapsed splits
validate the visible minimum without replacing saved ratio; inactive splits retain state. Runtime accepts clamps for
drag/resize/restore, rejects out-of-bounds programmatic ratios, and validates the returned Effective against that operation. Runtime
imports no layout/GUI. Resources receive the finalized snapshot with these proportions and offsets, identical to published runtime
and native geometry; no native-only clamp. Gate viewport results overlay active measured offsets. Retain still-existing inactive
pane/surface descendants' offsets, remove deleted owners/axes and clamp on reveal; preserve ordinary legacy hide behavior outside
these compositions.

**Preflight is not publication.** The host's speculative geometry/resource probe owns disposable preparation only; it must never
write Bundle.pending. A private per-mutation preparation ticket binds the final prepared presentation to bundle/model/source/size
and the finalized snapshot. Only successful runtime publication promotes that ticket to accepted pending in the same non-yielding
commit; rejection discards only that ticket. Host after/apply consumes only accepted tickets. Advancing Sequence/StateRevision,
recording an outcome, or returning from Mutate with an error is not promotion authority. Thus failed callbacks cannot expose
speculative page/toggle/focus changes through b.pending. Independently accepted legacy reentrant mutations retain their own accepted
ticket and synchronize even if the outer interaction fails. This refines the existing state gate/host handoff, not a new public
transaction framework.

OpenSurface(Handle,ContextTarget) (SurfaceTarget,error) uses zero context for none; CloseSurface(SurfaceTarget,kind) accepts only
Cancel/Close, independent of captured field revisions. Session owns state; host owns native objects; application owns persistence.
OpenSurfaceFrom(dialog Handle, opener Handle, context ContextTarget) is the bounded
explicit-origin operation used by native command effects: derive the parent canvas
and restoration target from that actual opener, without a synthetic Focus mutation.
OpenSurface is the programmatic convenience using current valid focus or root;
it must not guess a native button origin from another nonmodal canvas's focus.
Mark an opening published only at successful owner publication. For prospective runtime parent hide/page selection/
split collapse/reload, stage child closure and request revocation with that parent transition; publish them only if the
transition commits. Failed geometry/callback preserves live children/requests and emits no result. Once committed,
child teardown cannot be vetoed by another geometry/acceptance check. Genuine parent native hide/close or disposal
revokes published tokens and queues one result immediately, then destroys, bypassing fallible geometry.
Candidate/abandoned openings emit none.
Failed successor disposal cannot close the predecessor; successful successor starts closed and closes published predecessor surfaces
with reason reload. DrainDialogResults() returns detached queued results once; host drains after each sync/teardown outside mutation
into OnDialogResult. At most one per published opening; drain a declaration’s queued result before allowing its next opening,
bounding the queue by the existing surface/tree limit without dropping forced results. Observer failure cannot veto close. Teardown
drains too. Preflight checks declarations/resources even in hidden pages/closed surfaces; no native pointers in runtime state.

Extend WCI1 Source.EventField with closed TabPageID/TabPreviousPageID (ActivatePage→SDL text), CommandContextItemID (item-context
InvokeCommand→text), CommandChecked (toggle InvokeCommand→boolean), DialogFieldValue (Accept→text). Only DialogFieldValue permits
required companion FieldPath: named path relative to that dialog instance, resolved to an owned input handle in preflight and read
from captured Fields, never later live state. Exactly one selector remains required. New selectors never coerce text; legacy
Event/Widget conversions stay unchanged. DialogFieldValue remains WCI2 text-only; WCI3 typed capture needs separately reviewed typed
mappings. Wrong owner/event/type/path/registration rejects preflight.

Add Plan.ResultMode with closed TextResult (zero/default, existing OutputField/text/setHandle rules) and DialogAcceptResult (Accept
only). The latter requires AcceptField:boolean and MessageField:text output fields, forbids
OutputField/setHandle/RevisionField/RevisionContext, and converts them to Reply.Accept; it returns no widget update itself. No
hidden input. Validate signature, actual Go registration AND returned record types/fields before decision; malformed result
commits/closes nothing and reports unknown domain outcome. Accepted=false means no domain commit; Message is valid UTF-8 <=32768
bytes, with a standard error if empty. True commits captured drafts then closes. Reject unknown result modes; legacy 0.2/default
TextResult behavior is unchanged. Application owns persistence.

```go
// actions.Page: {PageID:text} -> {Preview:text}; actions.Save: {Name:text} -> {Accepted:boolean, Message:text}
map[string]bridge.Plan{
  "actions.Page": {Inputs: map[string]bridge.Source{"PageID": {EventField: bridge.TabPageID}}, OutputField: "Preview"},
  "actions.Save": {Inputs: map[string]bridge.Source{"Name": {EventField: bridge.DialogFieldValue, FieldPath: "name"}},
    ResultMode: bridge.DialogAcceptResult, AcceptField: "Accepted", MessageField: "Message"},
}
```

DialogResult.Kind is accept/cancel/close; Reason is user/programmatic/parent-closed/parent-hidden/reload/dispose. User results carry
the triggering sequence; automatic lifecycle closure carries zero plus the unique surface token (not an SDL sequence).
OnDialogResult is a completion observer, not another source callback; Cancel and Close therefore cannot accidentally execute Save.
Command local Accept effects use the same decision and result pipeline as a direct dialog request, including validation failure and
repeated-click rejection. AcceptSequence/Domain retain the last acceptance attempt/outcome, or zero/not-called; thus Cancel after
domain success/UI conflict cannot imply rollback. Accept results alone include the committed captured fields.

Keep exactly the existing frontend/layout/widget/viewport/provider/host dimensions. Require widget AND host
tabs/split/command/menu/dialog/button-toggle at major 1 when used; page/menuGroup/item/separator are checked by their parent family.
Add host tab-activate, command-key, context-target, tooltip, dialog-modal, dialog-nonmodal and dialog-result at major 1 for
corresponding behavior; provider icon/1 when requested. Frontend is sdui/0.3 major 1; layout relative/1 and viewport scroll-x/y
remain separate. Derive requirements from all hidden/visible declarations; missing exact facts reject with source/path diagnostics.
No interaction or bridge dimension/registry. Actual Binder signature/plan preflight establishes readiness, not a flag.

**Selected nonmodal adapter: `app.NewWindow`, with application-owned parent lifetime.** Coordinator reports review acceptance of
this design choice; the bounded probe below supports it. Product native acceptance remains pending. Prefer AppTabs, HSplit/VSplit,
Button, Menu/PopUpMenu and CustomWithoutButtons elsewhere. Pinned Fyne v2.8.1 source evidence: `app/app.go` delegates NewWindow to
CreateWindow; `window.go` provides separate Canvas, Close/Hide, SetCloseIntercept and SetOnClosed, but no parent/transient-window
relationship. Conversely, `widget/popup.go` installs outside Hide dismissal, `internal/widget/overlay_container.go` consumes
background taps, and `internal/driver/util.go` hit-tests only the top overlay when present. Ordinary NewPopUp therefore dismisses on
a parent click instead of allowing persistent simultaneous parent interaction. Reject that alternative; no custom
overlay/event-routing framework is justified. This choice introduces no OS transient/always-on-top guarantee.

The composition owner records parent/surface tokens, creates the child on the UI goroutine and never marks it master. Route native
close through SetCloseIntercept; revoke the token before Close and use idempotent SetOnClosed fallback without recursively closing.
In pinned GLFW `window.go`, Close bypasses the interceptor and invokes OnClosed before marking closing; Hide invokes neither. Thus
all parent hide/reload/close paths explicitly dispose owned children and emit one Close result; no reliance on Fyne ownership
propagation. Keep prepared opening/publication under existing gates. Restore the valid opener in its canvas; request OS focus only
for an explicit child gesture, never parent-driven or background disposal. `window_desktop.go` makes RequestFocus a no-op on
Wayland; elsewhere it remains best effort. Canvas.Focused alone is not OS-focus evidence. Native support must be truthful for the
tested desktop environment.

[Durable native selection probe](nonmodal-probe/README.md) retains source/dependency hashes, logs, OS captures and inventory;
source/binary identity matches the temporary `/tmp/wci2-surface-probe` build (Go 1.27.1/Fyne v2.8.1). Coordinator's actual XTest run
on isolated :189, without a WM, observed a dirty child remaining open during parent action, one local Enter commit, Escape revert
then Cancel, and parent Space reopening through restored opener focus. All five openings produced one terminal result each: Escape
Cancel, parent-command Close, parent-hidden Close, simulated-reload Close, parent-closed Close. Parent-driven closure added no
opener-focus request; exit 0, stderr empty. This supports NewWindow selection, not SDUI/SDL product, persistence, failed-publication
or WM/chrome conformance. Before M2 acceptance repeat A06/A08 with real runtime tokens, bridge and publication; test decoration
close and desktop focus policy where a WM exists. Keep environment-specific limits explicit. Native failures require concrete
adapter correction/review, not postponing this choice; no generalized framework, renderer or inventory expansion.

## 6. Export, acceptance and handoff

Composition/text show all panes/command links/menus/closed dialogs as static structure with source/use paths and initial state. SVG
rejects unsupported WCI2 with source-linked unsupported-pane/transient-export before writing partial artifacts; existing button
properties must be represented or rejected. No snapshot renderer required. Native omission requires prepared controls. Codegen
reconstructs Document/Root/rows/typed arguments, profile/provenance, not runtime state/resources/closures. Every kind/property must
be supported or rejected.

| Pending acceptance | Required result / parent matrix rows |
| --- | --- |
| WCI2-A01 | Positive/negative grammar, profile, schema, placement, duplicate/invalid refs, limits and reused definitions; 0.2 byte-stable output; spans/use chains; compile generated constructors and compare normalization. All frontend rows/R05/R26. |
| WCI2-A02 | Native tabs pointer/keys invoke bound callback once with stable IDs; same-page/programmatic/fallback/reload invoke zero; failure preserves old page/drafts/focus. Reuse, disable/remove/recreate, zero eligible, intent/activity and retained inactive scroll clamped on reveal. |
| WCI2-A03 | Native split drag/keyboard; collapse with positive relative/measured minima gives zero extent; resize stays collapsed; restore clamps or rejects without mutation; collapsible=false reload preflight. Gate returns clamps; runtime/resource/geometry/hit/focus agree; strict programmatic ratio rejects. |
| WCI2-A04 | Toolbar/menu/key invoke one canonical handler; exclusive updates atomic; effects do not double-dispatch. Test wrong/multiple payload, stale sequence/state/origin, reentrancy, mute, disabled, callback failure, tooltip/icon and key conflicts. |
| WCI2-A05 | Nested menus/groups, keyboard/pointer/Escape, blank/group/item contexts; stale/deleted/recreated item, changed generation, hide/reveal with unchanged handles, and disposed opener reject without action. Menu/context. |
| WCI2-A06 | Existing single-line input captures/Commit, menu→dirty draft revert→clean Cancel; composed/reused drafts, nested-field exclusion, modal/nonmodal focus/lifetime and one terminal result per opening. Include stale requests/reopening/observer failure and selected adapter evidence. New IME/multiline/typed-option proof belongs to WCI3. |
| WCI2-A07 | Real SDL tab/shared-command/Accept fixtures; false/true/error, signatures/FieldPath negatives, zero preflight calls. After draft or geometry conflict following domain success, Accept stays blocked; further edits then Cancel/Close succeed and retain attempt/outcome without stale-draft checks, replay or rollback claims. Preserve succeeded through post-Execute adapter conflicts; malformed output/unknown blocks replay; UTF-8 and combined 256-write boundary (256 fields plus one unrelated update rejects atomically). No reconciliation API. |
| WCI2-A08 | Inject missing hidden capability/resource, native preparation failure, stale state and failed reload; old bundle/transients/focus remain; rejected parent hide/collapse emits no child Close. Consumed-sequence callback failure never promotes speculative pending; accepted reentrant changes survive. Only published openings emit forced Close; abandoned candidate emits none; teardown bypasses failing gate and drains once. R16/R28. |
| WCI2-A09 | Inspect composition/text and explicit SVG diagnostics; unsupported properties never disappear; exact source/codegen identities, affected consumer regression and full SDUI race suite. Export/integration. |

Use WCI1’s separate Xvfb/XTest route, screenshots and correlated state/callback logs, not headless tests alone. All acceptance
remains pending. Reconcile with WCI1 pilot, then obtain independent stage review.

Coordinator handoff: loaded sdp 1.1.1, sdp-architect 2.0.0 and document-workflow; recovered Session0010 S2/S3. Only this draft
revised. Reconciliation/WCI3 proof pruning is complete; NewWindow selection is reported review-accepted and now linked to bounded
native evidence. James's read-only runtime memo informed the concrete gate, activity/offset, outcome and publication contracts
above; independent stage review remains pending. Coordinator reports WCI1 basic/expanded 16 PASS and empty-root 7 PASS, final gutter
pending; separate from WCI2 evidence. Coordinator records handoff in Session0010; no product code/other record edits, implementation
approval or inventory closure.

## Stage transition — Session0010 T002

WCI1 phase c39b330 is verified and independently reviewed; final native evidence
contains 54 passing checks and original integration suites pass. Reviewer confirmed
the final substantive contract hash cf93dea68510be05a626a3953159fc375ce04f55b18562c35eb3699d0859b839
including context, published closure, combined write bound and nonmodal probe.
Coordinator selects WCI2-M1 now under existing full-card owner authorization.
M1 delivers panes first; M2 command/surface implementation follows M1 evidence.
Earlier draft observations above retain their provenance and are not current
implementation results. All WCI2 product acceptance remains pending.

## M2 transition — Session0010 T003

M1 commit 403c540 is independently approved with 58 native checks. The revised
contract (pre-transition hash 9a9c8710a4245a28a27029af5cb636ae60b5d39e3c794c5ea8aea56c50e73187)
is approved for M2 entry, including basic-button compatibility, explicit opener
and synchronous menu dismissal handling. Coordinator selects bounded M2 code under
existing full-card owner authorization. Native adapter conformance remains an
implementation proof obligation; private-field access, dependency patches or
broader routing machinery require concrete design reconciliation.
