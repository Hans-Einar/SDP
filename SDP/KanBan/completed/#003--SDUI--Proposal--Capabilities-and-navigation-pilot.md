# Capabilities and a bounded navigation pilot

| Field | Value |
| --- | --- |
| id | KB-SDUI-003 |
| project | SDUI |
| type | Proposal |
| CardState | completed |
| Systems | SDUI |
| created | 2026-09-29T16:53:54.251554+00:00 |
| source | PLAN-SDP-0009; KB-SDP-041; external XFMD gap register |
| next_review | None for this delivered scope; integration and publication are separate |
| PlanId | PLAN-SDP-0022 |

## Delivered assignment — Session0010 T001–T003

The owner now requests execution of the work described in this card. Codex takes
coordination responsibility for the complete inventory, beginning with
[PLAN-SDP-0021](../../04--Design/SDUI/Widgets/Plan.md), then the staged
[PLAN-SDP-0022](../../05--Implementation/SDUI/Widgets/Plan.md). Earlier authorization
limits below describe the previous card-editing request, not this assignment.
The [Session](../../Sessions/session-%230010--SDUI_widgets.md) owns continuity.
The bounded inventory is now implemented, verified and independently reviewed; publication and merge remain unselected.

## Historical owner scope clarification — 2026-09-29

The 2026-10-07 owner refinement below expands the planned widget scope. The
existing prohibition on requiring full XFMD parity remains applicable.

SDUI primarily describes the concept, composition and placement of a UI using a
limited widget vocabulary. Running a basic SDUI and calling SDL remains part of
its purpose. It need not reproduce every function of a rich native application.
A tree is a possible useful addition, not an approved requirement to rebuild
XFMD in SDUI. Full editors, sophisticated dialogs, application-wide command
systems and host parity remain optional proposals requiring demonstrated need.
This clarification governs the study's successor selection; the gap inventory
records consumer observations, not a mandatory product backlog.

## Need and bounded next step

[PLAN-SDP-0009 study](../../02--Requirements/XFMD-Gaps/Study.md) evaluates
GAP-XFMD-SDUI-001–010. Own the follow-up here rather than creating ten loosely
coupled widget cards. First select a DesignPlan for required capabilities,
activation preflight and typed interaction state. The owner now requests the
concrete widget inventory below as the target for this card. Start implementation
with a collection/viewport slice, then complete the remaining families in stages.
This card revision selects planning scope; it does not start implementation or
adopt the illustrative syntax names.

Reuse SDL/go/bridge typed binding validation and SDUI/go/runtime drafts, atomic
updates, identity and reload protections. Keep frontend I/O-free. A missing SDL
file passing parsing is intentional; missing supplied runtime module/signature
must fail before activation. Scroll already parses but layout rejects it.

## Staged scope and acceptance

Follow [XGP3–XGP5](../../02--Requirements/XFMD-Gaps/Delivery-Proposal.md):

- First: separate frontend, layout, provider and native-host capabilities; no
  partial activation. Typed stable item identity plus model generation; stale,
  deleted, canceled or disposed targets cannot mutate a replacement UI.
- Pilot: bounded tree/list, expansion/loading/error, selection and scrolling,
  verified with keyboard/pointer and reload in the chosen actual host. Group
  rows are not selectable; expanding alone does not navigate.
- Subsequent required widget families: panes, commands/transient surfaces, typed
  values and basic multiline text input, as enumerated below; draft transactions
  remain application-owned. Keep relative dimensions.
- Optional research, not required SDUI parity: component source sets, editor/preview host surfaces and broader
  resource profiles. Require a real consumer and original-source diagnostics;
  do not create an unrestricted native-object escape hatch.

New profile work must address parser/AST/normalization, diagnostics, formatting
where supported, codegen, presentation, runtime and actual host verification.
Pictures/AST do not prove interaction. Unsupported exports either reject or use
an explicit labeled fallback selected by the consumer.

## Boundaries and dependencies

[KB-SDUI-004](../backlog/%23004--SDUI--Bug--Text-and-Markdown-fidelity.md) independently owns
text fidelity. [KB-SDL-005](../backlog/%23005--SDL--Change--System-and-source-sets.md) owns SDL
source sets, not SDUI component imports. The first SDUI fixture-driven pilot need
not wait for every SDL feature. Real generated-view integration needs a coherent
producer revision and existing SDPTool delegation.

External XFMD owns its C++ document/open/save/lease workflow and any native
integration. The external discovery register is at
`SDP/02--Requirements/XFMD-Modeling/Gaps.md` in its repository; no local checkout
is guaranteed. No FOX bridge, XFMD rewrite, new Rust renderer or full document-editor parity
is selected. Basic multiline Fyne-backed input is included below. [KB040](../backlog/%23040--Proposal--KanBan-TUI.md) and Ponsse are potential reuse
inputs, not permission to add another host now.

## Completion boundary

A design-only successor must explicitly hand remaining implementation to a
selected plan/card; this proposal cannot be closed as implemented from design
or a static export. Retain a disposition for every SDUI-001–010 gap; the explicit
2026-10-07 inventory below defines this card’s widget delivery scope. Defer a
listed family only through a recorded owner disposition with a linked successor;
a parser-only implementation or a silent downgrade to a static label is not
completion. Independent optional research keeps its earlier disposition.
Native workflow evidence and independent review
are needed for claims about delivered interactive behavior.

## Owner refinement — 2026-10-07, complete widget inventory

Source: external XFMD Session0004 T005. Owner asks to update KB-SDUI-003 with all
widgets it should implement and explicitly allows XM-M2 in normal model discovery.
The following inventory is the proposed concrete realization of that authorized
scope, based on GAP-XFMD-SDUI-001–010 and the XM-M2 probes. It supersedes the
older tree-only candidate/optional-widget framing for these listed families.
The owner has authorized this card update, not execution, release or a new grammar.

### Widget families and observable behavior

Names below identify capabilities. A DesignPlan must select canonical names and
whether a capability is a widget, property or composition; do not infer grammar
from rejected XM-M2 spellings. Reuse Fyne controls when they meet SDUI semantics.

| Family | Required behavior | Source / acceptance focus |
| --- | --- | --- |
| Button and toggle button (extend existing button) | Text/icon presentation, accessible name, enabled/checked state, tooltip and activation; shared command identity for toolbar/menu/keyboard; exclusive groups where needed | SDUI-007; disabled commands never fire; checked state is typed; programmatic changes do not simulate clicks |
| Single-line input (retain and extend existing input) | Label, value, placeholder, read-only/disabled state, draft/change/commit/validation behavior; focus and Unicode input | SDUI-005/006/008; preserve existing draft/revert, revision and reload behavior |
| Multiline text input (`textarea` or multiline input property) | Load supplied text into a multiline buffer, edit/select/copy/paste, line breaks, basic undo/redo, scrolling, read-only state and explicit change/commit events | SDUI-005; UTF-8/IME and actual Fyne editing tests; Load/Save are explicit application/SDL actions, not parser or widget file I/O; no syntax editor or full XFMD parity |
| Tree | Stable node identities, hierarchical children, expansion separate from activation, selection, lazy/loading/error/retry state and scrolling | SDUI-001; keyboard/pointer, stale responses, deleted items and reload; expanding alone does not invoke the selected action |
| List | Stable row identities and selection/activation, grouped sections with nonselectable headings/separators, bounded incremental data and scrolling | SDUI-001; application owns sorting/data; verify mixed folder/file examples, stale selection and grouping |
| Tabs | Named pages, selected-page state, enabled/visible state and activation events; retain compatible per-page state/focus | SDUI-002; switching pages does not recreate unrelated state or invoke hidden actions |
| Splitter / split pane | Horizontal/vertical split, relative proportions, minimum extents, drag/keyboard adjustment, collapse and restore | SDUI-002; resize/reload retains valid proportions and focus; no unbounded pixel-layout requirement |
| Menu / context menu | Command-backed items, groups/separators, enabled/checked state and basic submenus; pointer and keyboard context invocation with stable target identity | SDUI-004/007; `onContext` denotes the required event capability, not adopted spelling; dismiss/Escape does not activate; stale targets refused |
| Dialog | Titled transient content, explicit modal/nonmodal behavior, accept/cancel/close results, focus capture/restoration and parent lifetime | SDUI-004; cancellation leaves domain drafts untouched; use ordinary content composition, not a separate widget language |
| Checkbox | Typed boolean value with label, enabled/read-only state and change/commit semantics | SDUI-006; pointer/keyboard agree; programmatic updates emit no user action; tri-state is not implicitly required |
| Slider | Typed numeric value, min/max/step, continuous-changing versus committed events, label/value feedback | SDUI-006; clamping, invalid bounds, keyboard increments and rollback of a form draft |
| Choice / select | Stable option IDs and labels, selected value, empty/disabled options, dropdown selection and validation | SDUI-006; never use display text or list position as domain identity; changed option sets do not silently select a different value |
| Numeric input / spin control | Editable typed numeric value with min/max/step, validation and increment/decrement controls | SDUI-006; complements slider for precise settings; invalid intermediate draft is distinct from accepted value |
| SVG / preview content (retain existing svg and Markdown) | Explicit provider/resource identity, declared supported content, accessible description and truthful static fallback; an area may remain a labelled preview placeholder | SDUI-010; no automatic resource fetching, unrestricted native objects or full XFMD renderer parity; glyph/wrapping fidelity remains KB-SDUI-004 |

Frames/groups, rows and Markdown labels remain reusable composition primitives;
there is no need to create a new widget for every toolbar, status bar or form.
Radio-style exclusive choices may be expressed through the selected choice/toggle
contract; a second independent radio widget is not a completion requirement.
Tables/data grids, rich text editors, embedded browsers and specialist controls
are outside this bounded inventory unless a subsequent owner decision adds them.

### Cross-cutting capabilities required for the inventory

- Scrollable viewports: make already parsed scroll intent work through layout,
  clipping, runtime and Fyne; define axes, extent/offset, wheel/keyboard access,
  resizing/clamping, focus visibility and hit testing. Parser success is not proof.
- Shared semantic properties/events: enabled, visible, read-only where meaningful,
  accessible label, focus/tab order, command invocation, checked/value state,
  validation feedback, icon identity and context target. Define them centrally;
  avoid toolkit constants and per-widget incompatible callback conventions.
- Typed binding/activation: reuse SDL bridge validation and SDUI runtime handles.
  Validate module signatures, capabilities and event payloads before connected
  activation; missing/stale/disposed bindings cannot partially activate or mutate
  a replacement instance. Local prototypes must explicitly identify unbound SDL
  callbacks. Basic SDL-connected operation needs a real end-to-end fixture.
- Every added family must survive parse/profile/AST/normalization with original
  source spans, definition reuse and instance identity; cover layout, composition
  diagrams, static text, SVG/code generation where supported and actual Fyne
  interaction. Each exporter must implement or explicitly declare unsupported
  output/fallback. Do not advertise an inert export as an interactive widget.
- Preserve programmatic-update versus user-event distinction, atomic draft/form
  updates, cancellation, keyboard/pointer equivalence and reload/disposal guards.
  Domain persistence, file operations and application transactions stay outside
  the SDUI parser and generic controls.

### Recommended implementation order and completion

1. Select the versioned profile/capability/typed-event design and acceptance
   matrix; reuse KB-SDUI-005 readiness without treating its local-prototype check
   as complete SDL binding validation.
2. Implement tree/list plus scrolling as the first complete interactive slice.
3. Implement tabs/splits and command/button/menu/dialog behavior.
4. Implement checkbox/slider/select/numeric fields and basic multiline input.
5. Integrate every family with producer discovery, combined composition/text,
   native host and explicit SDL binding fixtures; test compatibility and package
   matching versions for consumers. Maintain a delivered/pending matrix per family.

Completion requires all listed families at the declared support level, meaningful
native tests, inspected generated previews and an explicit disposition of each
old SDUI gap. A design-only phase hands remaining implementation forward; it does
not complete this card. No requirement to rebuild XFMD's FOX UI in SDUI is added.

### XFMD distribution and historical fixtures

XFMD0.7 gets discovery from SDPTool/gh sdp but ships private sdui-preview and
sdui-fyne helpers. A gh extension update alone does not upgrade those helpers.
Coordinate a matching SDPTool distribution plus rebuilt/repackaged XFMD helpers
and native integration evidence; C++ changes are required only when consumer
contracts change. No release version or publication is selected by this card.

Owner explicitly permits XM-M2 fixtures in the normal model overview. Retain
those sources and their diagnostics; withdraw the preceding XFMD suggestion to
exclude them. They are dated positive/negative evidence, not an application defect
list. Do not change frozen fixtures, expected results or hashes when new widgets
are implemented: add a new versioned test/evidence baseline. Some historical
negative examples may become valid later, while unsupported combinations should
still report genuine diagnostics. Do not hide every invalid source.

External source of owner instructions:
file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230004--SDUI-workflow.md

## Worklog

- 2026-09-29T16:53:54.251554+00:00 — EVT-KB-SDUI-000012: Registered from completed producer research; remains backlog.

- 2026-09-29T17:55:11.132172+00:00 — EVT-KB-SDUI-000014: Owner clarifies limited conceptual/layout UI and basic SDL-connected execution. Rich application parity is not required; tree is a candidate. Update proposal scope; keep backlog, no implementation selected.

- 2026-10-07T21:44:06.393995+00:00 — EVT-KB-SDUI-000018: Applied owner XFMD Session0004 T005/T006 widget inventory refinement. Retain XM-M2 in ordinary discovery, backlog state and separate Fyne ownership; no implementation or publication selected.

- 2026-10-07T21:49:33.950733+00:00 — EVT-KB-SDUI-000019: Owner requests the KB-SDUI-003 job; backlog -> active/in-progress. Start WCD1 design, retain every family, link staged implementation and Session0010. No implementation completion claimed.

- 2026-10-07T21:54:14.780622+00:00 — EVT-KB-SDUI-000020: WCD1 independently reviewed and complete at design level. Start PLAN-SDP-0022 WCI0 detached preparation/admission in isolated clone; every widget family remains pending.

- 2026-10-07T22:26:58.777242+00:00 — EVT-KB-SDUI-000021: WCI0 independently reviewed at detached preparation/admission boundary; WCI1 contract next. Track existing SDL constructor/frontend mismatch as prerequisite within this card before generation claims.

- 2026-10-07T23:43:13.186726+00:00 — EVT-KB-SDUI-000022: WCI1-M1 tree/list/scroll and guarded native publication implemented, verified and independently reviewed (c39b330; 54 native checks). Keep in-progress for WCI2–WCI4. Evidence: ../../05--Implementation/SDUI/Widgets/Evidence-WCI1.md.

- 2026-10-08T08:07:43.894416+00:00 — EVT-KB-SDUI-000023: WCI2-M1 tabs/split delivered and reviewed at 403c540; remain in-progress for commands/surfaces, values/text and providers/package preparation. Evidence: ../../05--Implementation/SDUI/Widgets/Evidence-WCI2-M1.md.

- 2026-10-08T09:45:51.265777+00:00 — EVT-KB-SDUI-000024: WCI2-M2 delivered at 0fc15c8; 118 native checks, 26 exactly-once results and independent approval. WCI3-M1 selected; retain in-progress for values/text/providers/package preparation. Evidence: ../../05--Implementation/SDUI/Widgets/Evidence-WCI2-M2.md. Backlog review preserves separate KB004 and KB005 dispositions.

- 2026-10-08T10:51:53.682710+00:00 — EVT-KB-SDUI-000025: WCI3-M1 delivered at ea49991f; 113 native checks, 12 exactly-once results and independent approval. WCI3-M2 selected; remain in-progress for extended text/IME and WCI4. Evidence: ../../05--Implementation/SDUI/Widgets/Evidence-WCI3-M1.md. Backlog review retains separate KB004/KB005 dispositions.

M2 adapter reconciliation: the reviewed Values-and-text native editing refinement
retains ordinary undo/redo while explicitly limiting history preservation after an
actual native edit rejected by runtime/publication. Identical displayed bytes keep
history; different programmatic replacement resets it. Failed Commit/reload/probe
and validation-invalid admitted drafts retain history. No Fyne fork or parallel
editor is selected; native acceptance must prove these boundaries.

- 2026-10-08T12:14:43.708648+00:00 — EVT-KB-SDUI-000026: WCI3-M2 delivered at 69d0a332; 92 native checks, 5 exactly-once results and independent approval. WCI4 selected; remain in-progress for previews, consumer packages and integrated review. Evidence: ../../05--Implementation/SDUI/Widgets/Evidence-WCI3-M2.md. Backlog review retains separate KB004/KB005 dispositions.

- 2026-10-08T12:56:13.294018+00:00 — EVT-KB-SDUI-000028: WCI4 implementation and candidate-one package/protocol pilots are recorded. Native owner-loss preview cleanup and independent finite-SVG geometry findings require correction before final aggregate acceptance; remain in-progress. Matching package has ten binaries; supplied consumer GUI test and 50 protocol commands pass on that prior candidate. Final corrected-source native/all-family suites and independent review remain required. Dark-theme contrast observation is registered separately on KB-SDUI-004.

- 2026-10-08T13:07:30.392770+00:00 — EVT-KB-SDUI-000029: WCI4-M1 delivered at a026a5d1; reviewed 74-path provider/native candidate, matching ten-binary payload, 50 protocol commands, 24 native checks and ten exact receipts. WCI4-M2 stays active for all-family integrated acceptance. Backlog review retains KB-SDUI-004 fidelity and separate KB-SDUI-005 disposition.

## Final delivery — WCI4-M2

The complete selected widget inventory is delivered through WCI0–WCI4.
[Final evidence](../../05--Implementation/SDUI/Widgets/Evidence-WCI4-M2.md) records
corrected native SVG geometry, explicit supported/fallback boundaries, matching
packages, native receipts, suites and independent review. No selected family remains
open. This closes the local card on its selected baseline; it does not close the
external XFMD gap register, merge earlier dependency commits or publish a release.
Milestone backlog/onHold review retains KB-SDUI-004 backlog, KB-SDUI-005 separate
gate-review and optional source sets/rich-editor proposals unselected. No new scope.

- 2026-10-08T13:50:20.709403+00:00 — EVT-KB-SDUI-000030: active/in-progress -> completed. WCI4-M2 and full bounded widget inventory verified and independently reviewed at dcf2a741; 28 applicable native workflows, 321 checks, 25 exact terminal results, matching package/protocol/consumer and passing affected suites. No merge/publication.
