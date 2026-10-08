# SDUI 0.3 — development collections and panes profile

The WCI1 frontend adds exact `sdui 0.3;` alongside the preserved
[0.2 profile](language.md). Source acceptance, static export and native connected
readiness are separate claims. This document specifies the source/export boundary;
the [WCI1 contract](../../SDP/04--Design/SDUI/Widgets/Collections.md) governs
provider, runtime, layout, native interaction and activation acceptance. WCI2-M1
adds the tabs/page/split source contract below within this same unreleased profile;
its [reviewed stage contract](../../SDP/04--Design/SDUI/Widgets/Panes-and-commands.md)
governs runtime, geometry and native acceptance. M2 commands/menus/dialogs are not
part of the implemented source vocabulary.

## Grammar and calls

The [0.3 EBNF](../grammar/sdui-0.3.ebnf) retains the 0.2 lexical, composition,
reuse, formatting, resource-limit and symbolic-reference productions, adding
bounded pane call bodies below. Reject alternative
spellings such as `0.30` and `3e-1`. Existing button/input/svg call schemas remain.

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

Only `tabs`, `page` and `split` require and accept a bracketed body after their
call. Formatting follows the body. Body `header/body/footer` assignments are
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
the exact instance identity. M1 adds no command reference resolver or lexical-scope
fields. Existing setHandle paths still resolve ordinary widget receivers.

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
