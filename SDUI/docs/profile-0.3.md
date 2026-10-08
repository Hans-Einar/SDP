# SDUI 0.3 — development widgets profile

The WCI1 frontend adds exact `sdui 0.3;` alongside the preserved
[0.2 profile](language.md). Source acceptance, static export and native connected
readiness are separate claims. This document specifies the source/export boundary;
the [WCI1 contract](../../SDP/04--Design/SDUI/Widgets/Collections.md) governs
provider, runtime, layout, native interaction and activation acceptance. WCI2-M1
adds the tabs/page/split source contract below within this same unreleased profile;
its [reviewed stage contract](../../SDP/04--Design/SDUI/Widgets/Panes-and-commands.md)
governs runtime, geometry and native acceptance. M2 extends the same development
profile with commands, menus and dialogs below; frontend support is not native proof.

## Grammar and calls

The [0.3 EBNF](../grammar/sdui-0.3.ebnf) retains the 0.2 lexical, composition,
reuse, formatting, resource-limit and symbolic-reference productions, adding
bounded pane call bodies below. Reject alternative
spellings such as `0.30` and `3e-1`. Input/svg schemas remain unchanged; the 0.3
button extension below preserves legacy dispatch unless explicitly opted in.

```text
sdui 0.3;
ref: nav "collections.sdl";
items=<nodes=tree("Navigation", callback=nav.Activate.@invoke),
       rows=list(label="Entries")>;
page=[body=items, footer=<preview=input("Preview")>];
nav.Activate.setHandle(page.footer.preview);
```

Tree/list require a nonblank accessible string label, either positional first or
named `label=...`. The optional `callback` is a symbolic member reference; it
requires a widget name and declared module alias. Only collection activation uses
this callback. The frontend does not resolve modules or execute actions. Connected
preflight separately requires the applicable `@invoke` action and signature.

Reject duplicates, extra arguments, wrong types and positional arguments after
named arguments. No inline items, provider references, arrays, event lambdas or
new callback property names are added. Applications supply typed collection data
and providers independently for each normalized instance. Existing scroll syntax
expresses intent; parsing it does not prove layout or host support.

## M1 pane bodies and identity

`tabs`, `page` and `split` require a bracketed body after their call; M2 also
allows the menu/menuGroup/dialog bodies below. Formatting follows the body. Body `header/body/footer` assignments are
ordinary named children; use a nested frame to declare frame regions.

```text
sdui 0.3;
ref: actions "actions.sdl";
Main=[
  panes=split(axis="horizontal",minFirst=0.15,minSecond=0.2)[
    navigation=tree("Navigation");
    work=tabs("Workspace",selected="overview",callback=actions.Page.@invoke)[
      overview=page("Overview")[draft=input("Draft",value="")];
      notes=page("Notes")[]
    ]
  ];
  preview=input("Preview",value="")
];
actions.Page.setHandle(Main.preview);
```

| Form | Source arguments and body |
| --- | --- |
| tabs | Required nonblank label, optional selected string and callback member reference. Callback requires a named owner and means ActivatePage only. One or more named pages, exactly one per row. |
| page | Required nonblank label, optional nonblank symbolic icon string. Ordinary content rows, including empty. Direct child of tabs after reuse expansion. |
| split | Required axis string horizontal/vertical; optional numeric proportion/minFirst/minSecond and boolean collapsible. Exactly two named children, one per row. |

One positional argument precedes named arguments (label, or split axis). Unknown,
duplicate and wrong-type fields reject. Names come from explicit assignments or
named definition uses, never labels/indices. A standalone page definition is a
reusable template; embedding it outside tabs rejects. A selected page names an
eligible direct page; its enabled/visible declaration intent must allow selection.
Inactive ancestor tabs do not invalidate their initial local selection. Without
selected, runtime chooses the first eligible page or no page if all are ineligible.

Split defaults normalize to typed Literals: proportion .5, minFirst/minSecond 0,
collapsible true. Numbers are finite; minima are nonnegative with sum <1;
proportion must lie in [minFirst,1-minSecond]. Measured minima, collapse and resize
are runtime/layout operations, not parser guesses. Icon resources and SDL references
are symbolic; frontend/codegen never open them.

Bodies reuse existing Node/Instance fields: Kind composition, Widget tabs/page/split,
typed Arguments and Rows. No new AST format or family version. All source/reuse
spans, paths and Profile propagate through expansion and generated constructors.
`parser.PaneChildren(owner)` returns ordered `PaneChild{ID,Node}` for normalized
tabs/split, validating direct named shape. ID is local to its owner; Node.Path is
the exact instance identity. M2 supplies the command reference resolver below. Existing setHandle paths still
resolve ordinary widget receivers.

## M2 commands, menus and dialogs

| Form | Arguments and structure |
| --- | --- |
| command | Named declaration, required nonblank label. Optional icon/tooltip strings, toggle/checked booleans, exclusive/key/context/target/effect strings and callback reference. No body or geometry; enabled/visible formatting only. |
| button | Label required unless command string supplied. Same behavior fields as command, or command reference with only label/icon/tooltip overrides. No body. |
| menu | Required nonblank label; mode defaults to bar, or context/submenu. Context requires target; other modes forbid it. Body contains one item/separator/menuGroup/submenu per row. Nested menus explicitly use mode=submenu. Context accepts enabled/visible formatting only. |
| menuGroup | Nonblank label and menu-child body; non-invokable heading. |
| item | Required named command string argument; no positional argument or body. |
| separator | No arguments, body, formatting, action or state. |
| dialog | Named declaration, nonblank label, modal boolean defaults true, optional callback. Ordinary composed body; initially closed auxiliary surface. |

One positional label precedes named arguments. Unknown/duplicate/wrong-type fields
reject. Icon is a symbolic resource, never a fetched URL. Definitions can supply
stable names to reused command/dialog roots. Direct split children must be visual;
commands, dialogs and context menus cannot serve as the two split panes.

checked/exclusive require toggle=true. Exclusive groups use the innermost definition
instance and permit at most one initially checked member. Context defaults none;
widget/item require target, and item requires a tree/list. A context menu's captured
control must match each contextual command's target. Ordinary no-callback commands
are unbound; pure toggles are local. Effects are open/accept/cancel/close, forbid
callback/toggle, and require context none. Open needs a dialog target; other effects
forbid target and require an enclosing dialog. M2 callbacks use symbolic @invoke.
No parser operation opens SDL modules or executes an action.

A basic button (including callback/icon/tooltip) retains legacy Activate/Handler.
Presence of command/toggle/checked/exclusive/key/context/target/effect opts into
InvokeCommand/InteractionHandler; even toggle=false or context="none" is explicit.
A referring button has one canonical command owner; it cannot override behavior.
Key syntax is ordered Primary/Ctrl/Alt/Shift plus A–Z/0–9/F1–F12; letters/digits
require a modifier. Reject duplicate/unordered modifiers and Primary+Ctrl. Exact
same-surface duplicates reject; platform key normalization/reservations belong to
preparation. Native focus, captures, dismissal and dialog results remain runtime
and host responsibilities under the reviewed stage contract.

Normalization preserves source arguments and adds only the closed `$scope` string
Literal on M2 interaction nodes, recording the innermost expanded definition
instance. Entering reuse resets that scope; inline containers do not. Source cannot
spell the derived key. Node/Instance fields and AST/version tags remain unchanged;
legacy and M1 nodes acquire no derived fields.

`parser.ResolveInteractions(root)` strictly validates the selected expanded tree,
including hidden declarations. It returns `map[string]InteractionIdentity`, keyed
by exact Instance.Path, with Scope/Command/Target/Dialog strings. Command is the
canonical declaration or implicit button path; Target is a context/open target;
Dialog is nearest lexical enclosing dialog (empty means entry canvas). A dialog
record names its parent dialog. Owned inputs carry Dialog for acceptance mapping.
Source references resolve named paths relative to Scope, or leading `/` from the
selected entry. Anonymous container segments are skipped for name lookup, but
results retain exact paths. Reject parent/index/anonymous segments, missing,
ambiguous or wrong-kind targets. No parent-scope fallback or runtime resolver.

Normalize returns all definitions, including reusable templates whose absolute
entry references are unresolved in isolation. This is not selected-entry validity:
runtime, preparation, generation and exports must call ResolveInteractions before
claiming it. Static AST `local-profile` reports schema/placement validation only.
`ResolveDialogField(root, dialogPath, fieldPath)` resolves a named relative owned
input, rejecting nested-dialog fields; bridge uses its exact Path. `IsCommandButton`
exposes opt-in, and `IsAuxiliary` exposes non-flow declarations without state guesses.

## WCI3-M1 scalar fields

The same development profile adds four leaf controls. WCI3-M2 extends input
separately, as described below.

| Call | Arguments/defaults |
| --- | --- |
| checkbox(label) | Nonblank string label; optional boolean value=false and readOnly=false; callback reference. |
| slider(label,min,max,step,value) | Label and all four numeric tokens required; optional readOnly=false and callback. |
| number(label,min,max,step,value) | Same numeric contract; optional readOnly=false, placeholder string="" and callback. |
| select(label) | Nonblank label; optional option-ID string value="", required=false, readOnly=false and callback. |

Only label may be positional; all other arguments must be named. Duplicate,
unknown and wrong-type arguments reject. No bodies, inline options, expressions,
source onChange or file/provider I/O. Callbacks require named widgets and declared
module aliases; connected Commit preflight requires @invoke. enabled/visible stay
formatting properties. A named definition can supply a reused scalar root's name.
No source/normalized default arguments are inserted. Consumers apply these closed
defaults; codegen preserves whether an optional argument was present.

```text
sdui 0.3;
ref: settings "settings.sdl";
Main=[
  enabledFlag=checkbox("Active", value=false);
  level=slider("Level", min=0, max=100, step=1, value=50);
  count=number("Count", min=-0, max=1e1, step=1.00e-1, value=0.30);
  mode=select("Mode", value="compact", callback=settings.AcceptMode.@invoke)
];
settings.AcceptMode.setHandle(Main.mode);
```

Scalar calls use Kind widget and existing Arguments/spans/reuse representation.
Only slider/number min/max/step/value numeric tokens become
`Literal{Kind:"number-lexeme", Value:<exact token string>, Span:<original span>}`.
Quoted numeric strings reject. Existing .2, split and formatting numeric literals
remain Kind number with float64 values. No lexer, AST field/tag or format-version
migration occurs. Generated constructors retain the exact spelling, including
signed zero, exponent syntax and significant digits; activation needs no source file.

`parser.NumericArguments(n)` returns `(min,max,step,value string,err error)` and
checks the closed .3 numeric-widget literal representation. It does not format
binary64 values back to source or validate a user's text draft. The stdlib-only
`numeric` package owns exact decimal/grid arithmetic; parser and runtime can both
import it without a dependency cycle. Parser validation uses NewGrid and Parse for
source initial admission, reporting errors at the relevant numeric token span.

Constraints require finite operands, min<max, step>0, initial range membership and
exact rational `(value-min)/step` integrality. Text/source admission never snaps
or uses a tolerance. Source decimal 0.3 on step 0.1 is valid; an off-grid decimal
that rounds to the same binary64 is still invalid. The reviewed bounded decimal
scan admits at most 32768 mantissa digits and effective exponent magnitude 4096
for nonzero coefficients before constructing powers. Exact safe53 integral
min/max/step grids use uint64 ticks without the fractional interval cap. All other
Go grids require at most 2^26 intervals and step strictly larger than adjacent
binary64 spacing at the rounded endpoints, checking both directions. Connected
SDL separately requires exact integers within ±9007199254740991; parser does not
infer transport compatibility from a symbolic callback. See the
[reviewed values contract](../../SDP/04--Design/SDUI/Widgets/Values-and-text.md)
for the complete arithmetic/state boundary.

Choice options are supplied by exact normalized instance path, outside the AST:
unique stable IDs, duplicate labels allowed, <=4096 options. Static source can
preserve an initial ID but cannot establish its existence/eligibility without the
provider. Required empty choice is an invalid editable state; no implicit first
option. Codegen contains no options, live generations or provider closures.

ResolveInteractions includes nearest-dialog ownership for these scalar fields,
without adding command identities or $scope arguments to them. Context widget
targets may reference scalars; item contexts remain tree/list-only.
ResolveDialogField remains input-only for WCI2 text mappings. Mixed Go Accept uses
runtime typed field capture; scalar data never enters legacy Widget.Value/Draft
strings or overwrites retained numeric lexemes in projected source Arguments.

Static descriptions identify every scalar, label/path/span, initial values,
constraints and symbolic Commit binding; select reports unsupplied options.
Public SVG rejects all four with unsupported-value-export, including hidden fields
and supplied measured state. The CLI preserves an existing artifact on failure.
Native SkipControls requires a matching prepared checkbox/slider/select/number
entry for every omitted field; InteractionRoot still validates dialog canvases.
Standalone prototype preflight reports unsupported-value for unavailable scalar
adapters, or unsupported-provider / choice-options for select. Source acceptance
and static descriptions do not prove native/provider/SDL readiness.

## WCI3-M2 extended input

`input(text, value="", multiline=false, readOnly=false, placeholder?, required=false,
callback?)` remains a leaf. text/value/placeholder are strings; multiline/readOnly/
required are booleans; callback is the existing symbolic Commit reference. Only text
may be positional, first. Unknown, duplicate and wrong-type arguments reject.
Callbacks require a named widget (an extended-input definition root is named by its
definition); connected execution requires @invoke. No textarea, source onChange,
validator, file I/O, wrap or scrolling arguments are introduced.

Presence of any new argument opts into extended behavior, even explicit false or
empty placeholder. Basic .3 inputs without these arguments retain legacy semantics;
.2 rejects every new argument. No profile migration is required. Omitted placeholder
uses text; explicit empty suppresses it. Empty text remains legal. Source defaults
are never inserted into Arguments, and boolean/string literals retain existing AST
shapes and spans. Reuse and generated constructors preserve argument presence.

`parser.InputOptions(*Instance) (InputPolicy, error)` is the shared source facts
helper. InputPolicy contains Extended, Multiline, ReadOnly, Required, Placeholder
and PlaceholderSet. It checks the profile, leaf shape and closed typed arguments;
it does not validate live text, mutate source or establish native/SDL readiness.
Selected-root ResolveInteractions also validates input schemas, including hidden
nodes and input dialog ownership. Widget.Value/Draft remain the runtime text store;
Snapshot.Fields projects text with validation/metadata, not new source arguments.

Extended text Commit uses typed String capture; Change is a local Go notification.
ControlText maps exact text to SDL; TextResult remains text/input. Extended Commit
requires a self receiver and exact returned echo, with revision guards. Basic input
and legacy nonself TextResult remain unchanged; Load may target readOnly input.
Validation, draft retention and native editing belong to runtime/document-host
adapters, not parser normalization. required-empty is an editable invalid state,
not a grammar error. The single-line newline predicate is exactly CR or LF;
text is not normalized or silently trimmed. Actual IME behavior needs native proof.

Structural descriptions show explicit/default properties, initial text, source
path/span and unverified symbolic binding. Public SVG rejects any extended input
with unsupported-text-export, including false/empty opt-in, hidden input or supplied
geometry/state. Native background omission requires an exact prepared `input`
inventory entry; there is no new widget kind. Standalone prototype reports
unsupported-text for extended input: it has no document-host text adapter. This
does not require migration of the legacy launcher. The connected document-host
route must be explicitly prepared and tested with its actual capabilities/binder;
static support or generated constructors alone cannot claim that readiness.

## Profile identity and compatibility

Document.Profile remains explicit: `sdui/0.2` or `sdui/0.3`. AST JSON keeps the
same three envelope keys: `astFormat`, `validation`, `document`. The validation
value is `syntax-only` or `local-profile`. `parser.ASTFormat(profile)` returns
`sdui-ast/0.2` or `sdui-ast/0.3`; unknown document profiles reject. Source AST
field/tag shapes and spans are unchanged. Syntax-only mode can retain unknown
calls; local validation is required before any profile-support claim.

Every normalized 0.3 Instance has `Profile: "sdui/0.3"`. Legacy 0.2 instances retain
an empty Profile sentinel omitted from both ordinary and tagged JSON, preserving
their prior representation. `parser.EffectiveProfile(root)` validates the entire
bounded normalized tree, including hidden descendants, returning the effective
profile or an error for mixed, unknown or malformed trees. This helper does not
establish widget/provider/native capability support. Clone/reuse/reload consumers
must preserve the field and validate correspondence to the Document profile.

Generated 0.3 Go constructors use `sdui-go-model/2`, exposing `UISourceSHA256`,
`UIProfile`, `UIASTFormat`, `Document()` and `Root()`. They contain models only;
application providers and bindings remain supplied at activation. Generated 0.2
bytes retain `sdui-go-model/1` and omit the new empty Instance.Profile field and
0.3 metadata constants. Frozen 0.2 fixtures are not migrated to the new profile.

## Structural descriptions and SVG

Text/Markdown and combined composition explicitly identify Tree/List, accessible
labels, source/instance identity, symbolic activation and absent provider data.
They are structural descriptions, not interactive widgets or fetched item lists.
Unknown widget/node kinds reject, including hidden unsupported content.
Pane descriptions include every declaration, label/path, initial selection, split
constraints, symbolic callback/icon and visibility intent. They do not claim live
selection or native geometry; combined composition retains source and reuse links.

Public collection SVG export rejects with `unsupported-collection-export`, the
instance path and original source span, even when geometry is supplied. No collection
snapshot renderer or automatic static fallback is provided. The CLI checks this
before measurement/provider work and preserves an existing output on failure.
Public pane SVG likewise rejects with `unsupported-pane-export`, including hidden
panes and supplied measured geometry. Standalone prototype preflight rejects panes
with `unsupported-pane`; static source support cannot establish native readiness.

Native hosts may use SVG as a background behind actual controls. `svg.Options`
then requires `SkipControls: true` and `NativeControls`, a map of exact normalized
paths to kinds for controls actually prepared in that bundle. Every omitted control
must match; extra/mismatched entries and unknown widgets reject. This inventory
is a trusted host-adapter assertion, not public collection-export support. Merely
setting SkipControls cannot make a collection disappear successfully. Ordinary
0.2 static SVG remains supported.

For pane backgrounds, prepared inventory includes each tabs/split path and every
descendant button/input/tree/list control, including inactive pages. Page has no
inventory entry: its header/shell belongs to the enclosing prepared tabs adapter.
Omitting native pane chrome never skips descendant validation or supported content
painting. Missing, extra or wrong-kind entries reject; unknown compositions cannot
be licensed by inventing an inventory key.

M2 structural text describes every command/menu/dialog declaration, source path,
initial checked/exclusive intent, symbolic callbacks/icons and closed surfaces,
including hidden declarations. It does not simulate live menus or acceptance.
Public M2 SVG rejects with `unsupported-interaction-export`, even with supplied
geometry. Basic button icon/tooltip decorations also reject in public SVG rather
than silently disappearing; that does not promote their dispatch. Standalone
prototype preflight rejects unavailable M2 adapters with
`unsupported-interaction`; parsing cannot claim connected readiness.

Native background inventory additionally includes each prepared menu and dialog
adapter (including context/submenu and initially closed declarations). Menu-owned
items/groups/separators have no separate inventory entry; command declarations
have no native control and require strict entry resolution before omission.
Buttons still use exact button entries. `Options.InteractionRoot` supplies the
full selected snapshot for a native dialog canvas; the rendered root must be the
identical node within that tree. It is valid only with SkipControls and never
relaxes public export. Inventory remains scoped to that canvas's source subtree.
