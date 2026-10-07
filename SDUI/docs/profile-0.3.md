# SDUI 0.3 — bounded collection source profile

The WCI1 frontend adds exact `sdui 0.3;` alongside the preserved
[0.2 profile](language.md). Source acceptance, static export and native connected
readiness are separate claims. This document specifies the source/export boundary;
the [WCI1 contract](../../SDP/04--Design/SDUI/Widgets/Collections.md) governs
provider, runtime, layout, native interaction and activation acceptance.

## Grammar and calls

The [0.3 EBNF](../grammar/sdui-0.3.ebnf) retains the 0.2 lexical, composition,
reuse, formatting, resource-limit and symbolic-reference productions. Only its
exact header and locally supported widget vocabulary differ. Reject alternative
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

Public collection SVG export rejects with `unsupported-collection-export`, the
instance path and original source span, even when geometry is supplied. No collection
snapshot renderer or automatic static fallback is provided. The CLI checks this
before measurement/provider work and preserves an existing output on failure.

Native hosts may use SVG as a background behind actual controls. `svg.Options`
then requires `SkipControls: true` and `NativeControls`, a map of exact normalized
paths to kinds for controls actually prepared in that bundle. Every omitted control
must match; extra/mismatched entries and unknown widgets reject. This inventory
is a trusted host-adapter assertion, not public collection-export support. Merely
setting SkipControls cannot make a collection disappear successfully. Ordinary
0.2 static SVG remains supported.
