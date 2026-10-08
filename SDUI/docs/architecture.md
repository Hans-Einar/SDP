# SDUI — implemented Go architecture

Updated for the WCI2 development candidate, 2026-10-08. One active frontend
preserves SDUI 0.2 and adds the bounded development 0.3 profile. Executable modules/commands are in the [Go area](../go/README.md); the shared design is [described in SDL](../design/README.md).

| Package | Responsibility |
| --- | --- |
| go/parser | Lexer, recursive descent, AST/spans, local rules and normalization |
| go/layout | One measured geometry; relative dimensions, rows, wrap, ratio and clipping |
| go/markdown | Bounded Goldmark content; measurement and registered Mermaid provider |
| go/svg | Static export using shared geometry, Go Regular glyphs and native-control appearance |
| go/presentation | Structural console/Markdown dumps and bounded control gallery |
| go/runtime | Session, handles, accepted/draft, collection/provider state, typed events and atomic prospective updates |
| go/reload | Validated candidates and compatible state/identity preservation |
| go/host/fynehost | Native input/button/collection adapters, focus/keyboard and guarded bundle publication |
| go/codegen | Independent typed Go constructors for Document/Root |
| go/cmd | CLI and native composition; file access belongs here, not in parser |

Parse → Normalize/Compile produces Document and expanded Instance trees. AST JSON uses sdui-ast/0.2 or sdui-ast/0.3 to match the selected source profile. Spans are half-open UTF-8 byte ranges with one-based Unicode line/column positions. Nodes preserve groups, rows, regions and suffix formatting. Runtime copies input models; Go structs are not language-level immutable. Do not mutate models in use by a host.

Explicit widget names produce public instance paths; anonymous segments use synthetic names. Session/Path/Generation/Kind identify handles. Reload preserves compatible named instances, accepted/draft and focus; deletion/type changes invalidate old handles. Source/binding failures retain the last valid model. Source watchers publish candidates through fyne.Do.

Layout receives the host's available area; source has no pixel width/height. Children use their nearest source ancestor; `{16:9,<->}` derives height from filled width. Fonts remain logical DIP during resize. SVG and Fyne share rectangles, text measurements and clipping; native controls use host rasterization/theme.

The SDL adapter lives in SDL/go/bridge. Parser ref/callback/setHandle are data; composition registers SDL modules, Go functions and typed bridge.Plan. Action-core 0.1 provides explicitly bounded execution. An accepted Go domain transaction cannot roll back if later UI publication fails; there is no automatic replay. See the [runtime contract](runtime-contract.md).

Limits and actual trials: [G1](../go/evidence/G1.md), [G2](../go/evidence/G2.md), [G3](../go/evidence/G3.md), [G4](../../SDL/go/evidence/G4.md), [G5](../../SDL/go/evidence/G5.md). No alternative SVG/Fyne/XFMD parser, extracted Rust crate or mandatory C ABI.

## XFMD combined preview integration — 2026-10-07

presentation.Combined adds per-container composition mindmaps and a provenance/binding
table to the existing Markdown export. Normalize preserves declaration and use
sites on instances; UTF-8 byte spans remain half-open. Region and ordered-row branches preserve source composition; nested containers
have detail maps to avoid dense, overlapping containment rectangles. Hidden nodes remain in
source composition even when the text prototype omits them. Source links require
the consuming client's matching snapshot; no filesystem access is added to parser.

cmd/sdui-preview owns bounded source I/O, revision checks, caller-owned bundles
and prototype check JSON. prototype.Check validates the selected frame, layout
and runtime model, and reports unbound symbolic callbacks/connections. It executes
none. The standalone Fyne command accepts an initial revision requirement,
retains explicit local input/button handlers and does not auto-load SDL modules.
The binaries are bundled privately by XFMD 0.7; SDPTool UIPreview reuses Combined.

## WCI0 detached preparation — KB-SDUI-003

`go/preparation` checks exact profile/layout/widget/provider/host capabilities,
normalizes the selected document and builds a detached runtime Session. It imports
no GUI or SDL and performs no I/O. `host/fynehost/admission` supplies concrete
legacy host facts and measured provider checks without a Fyne dependency. Prototype
check uses this shared path; native runtime checks admission before construction
and on model checks. Connected preparation takes an explicit SDL adapter; actual
module/signature validation stays in SDL bridge. Candidate preparation executes
no domain actions and does not itself publish a native model. WCI1 owns full
publication and new controls.

## WCI1 collection and viewport candidate

The same parser and normalizer accept exact 0.3 tree/list declarations. Each
normalized node retains its profile and original declaration/use provenance.
The 0.2 normalized representation and generated bytes retain their legacy form.
Generated 0.3 constructors identify sdui-go-model/2 and sdui-ast/0.3; providers
are application-owned runtime inputs, never generated closures or parser I/O.

Runtime owns stable item IDs, collection/request generations, selection, focus,
expansion, load status and viewport identities/offsets. Pure prospective state
gates validate data and measured geometry before publication. Providers execute
off the UI goroutine only after loading state is accepted; UI-goroutine delivery
rechecks the current bundle, request and provider identity. Cancellation revokes
acceptance before stopping the provider context, including providers that ignore
cancellation. Compatible detached successors carry sequence and compatible state
without old handlers or outstanding request tokens.

Layout consumes one detached runtime snapshot and native collection metrics. It
returns composed clips, content extents and effective offsets; nested wheel
remainder and focus reveal use that same geometry. Fyne owns native controls and
resource lifetime, with bounded adapters where the pinned toolkit cannot expose
its scroll state or required collection interaction. There is no separate row
scene or second scroll authority. DocumentHost prepares a complete detached bundle
and guards its publication; native behavior is verified separately from core
unit tests.

The SDL bridge extracts a validated CollectionItemID into a typed text field and
captures the result receiver before execution. A successful domain action is not
replayed after a stale receiver or UI publication failure. The real action-core
collection fixture exercises the shared tree/list action and retained engine
across UI reloads.

Structural text/composition describes collections and absent runtime data. Public
SVG rejects unsupported collection export with source provenance; native background
omission requires an exact inventory of prepared controls. Static helper and
SDPTool metadata preserve the actual profile. A standalone prototype without
collection providers reports unsupported-provider, not connected readiness.

These are candidate implementation boundaries, not a complete widget delivery
claim. [WCI1 acceptance](../../SDP/04--Design/SDUI/Widgets/Collections.md) and the
[implementation plan](../../SDP/05--Implementation/SDUI/Widgets/Plan.md) retain the
required native, independent-review and later-family obligations.


## WCI2 panes, commands and auxiliary surfaces

The same unreleased 0.3 profile adds tabs/splits and explicit shared commands,
menus and dialogs. Frontend selected-root resolution owns canonical identities;
runtime owns page selection, split state, command checked/exclusive state,
menu captures and exact surface opening tokens. Basic buttons retain their
legacy handler unless they explicitly opt into command behavior.

Layout measures each active canvas independently and returns one aggregate
presentation state. A pure gate validates prospective geometry/resources without
publishing pending state. Only the accepted final ticket may mount native changes.
Modal Fyne dialogs and ordinary nonmodal windows use the same runtime surface
lifecycle; native destruction revokes exact published openings even when geometry
can no longer be prepared. Closing an obsolete native window cannot revoke a
replacement opening. Publication and terminal-result delivery are distinct.

A bounded synchronous adapter retains a menu capture across Fyne's dismissal
before action ordering, allows one invocation, then revokes it. It does not queue
cleanup or restamp revisions. Focused controls forward otherwise unhandled
command shortcuts through one canonical route while native editing/navigation
retains precedence. Prepared icon bytes are application supplied and immutable.

SDL command mapping carries typed checked/context fields. Dialog acceptance
captures owned text fields and uses a closed typed acceptance result. Unknown
or post-domain publication failures retain their outcome and prevent automatic
replay; Cancel cannot roll back a successful domain action. Public SVG rejects
interactive declarations; native background rendering requires the matching
per-canvas control inventory and selected interaction root.

M1 is verified in [pane evidence](../../SDP/05--Implementation/SDUI/Widgets/Evidence-WCI2-M1.md).
M2 is verified and independently reviewed in [command/surface evidence](../../SDP/05--Implementation/SDUI/Widgets/Evidence-WCI2-M2.md), including 118 native checks and exact dynamic-parent lifecycle. These
boundaries do not imply WCI3 values/text or WCI4 provider/package completion.
