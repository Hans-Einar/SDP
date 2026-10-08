# WCI4 frontend API preparation — read-only, no stage selection

SDP Worker/document workflow reused, 2026-10-08. Only this memo is written.
M2 frontend remains frozen. Original authority: Providers-and-packaging.md,
SHA-256 `83d54b218c7c2dfbebff17b641c2fdc67e4ca7d1f5fb73fe65851de5bee8ad46`.
Also read the original WCI4 consumer preparation/launcher clarification and actual
clone parser, codegen, presentation, SVG, Markdown, preparation and host paths.
Clone HEAD at final inspection: `d6742742f0968661d15698518df688ffb1d28946`;
concurrent integration work is not WCI4 evidence. All APIs below are proposals,
not exports, implementation approval or new native capability claims.

## 1. Closed source forms and one facts helper

Use the existing .3 call grammar and existing Literal/Reference/Argument fields:

| Form | Source rules |
| --- | --- |
| svg(source, label?, description, fallback) | source Reference; optional string label; description/fallback presence opts in and requires both; member must be resource for opt-in. |
| markdown(text, description, fallback) | New call; all three string arguments required; no legacy call variant. |

Only source/text may be positional, first. Reject bodies, callbacks, unknown or
duplicate arguments and wrong literal payload types. Description: valid UTF-8,
nonblank, <=4096 bytes; preserve the supplied bytes rather than trimming them.
Fallback is exactly "reject" or "label". Explicit Markdown text is valid UTF-8,
<=32768 bytes; empty text is legal. Markdown syntax support is a content-preparation
decision, not a second parser grammar: unsupported HTML/images can use label policy.
Missing policy/description, wrong profile or malformed source identity is fatal
schema failure; fallback cannot authorize malformed source.

Bare Markdown strings and svg(source,label?) retain their exact old schemas,
AST/output bytes and placeholder behavior under .2 and .3. The legacy source
Reference member is not tightened to resource. Neither a label, @resource member
alone, nor .3 profile alone opts legacy SVG into resource preview. .2 rejects all
new arguments and markdown calls. Module paths remain symbolic and unopened.

Proposed sole normalized source-policy accessor:

```go
type PreviewPolicy struct {
    Explicit bool
    Description, Fallback string
}
func PreviewOptions(n *Instance) (PreviewPolicy, error)
```

Accept only normalized widget/svg or Kind markdown, and validate effective profile,
leaf shape and closed typed arguments. For bare Markdown, empty Arguments and empty
Widget mean Explicit=false. For SVG, presence of description OR fallback triggers
both-required validation; return Explicit=true only for valid opt-in. For explicit
Markdown, nonempty Arguments must contain exactly text/description/fallback;
Text must equal the string text argument and Widget must be empty. Reject nil,
wrong-kind, mixed profile, malformed constructed literals/references or partial
explicit metadata with parser.Diagnostic at the offending original span/path.

The helper derives no prepared/rendered state. SVG Source remains the validated
Arguments["source"].(Reference); label remains Argument("label"). Markdown content
remains Text. No new union payload, provider registry, runtime import or renderer
dependency is needed in parser. Source validation checks declared module aliases;
the node-only helper checks reference shape/member, not a nonexistent module table.
Resource preparation matches the declared module/object/member triple, not Span.

## 2. AST, normalization, reuse and codegen

Source Node for markdown(...) remains Kind widget/Widget markdown with typed
Arguments. Normalize it to Kind markdown, Widget empty, Text from the text Literal,
retaining all three Arguments with original spans. Do not synthesize source Node.Text
or mutate the source AST. Bare strings retain their existing Node.Text representation
and normalized empty Arguments. SVG remains Kind widget/Widget svg. Use normalized
Markdown kind for formatting/combination checks, including reuse overlays; otherwise
call spelling would accidentally introduce different layout rules from bare Markdown.

Preserve Declaration, original Span, each UseSite, unique expanded path, and separate
argument maps per use. Explicit Markdown is content, not an interactive widget:
no callback, field state, binding owner, setHandle receiver or command context target.
Existing normalized-kind rules already exclude Markdown from these widget roles.
No derived $scope, resource bytes/digests or default properties enter source arguments.

Keep sdui/0.3, sdui-ast/0.3, Version03=sdui-go-model/2 and current generated metadata;
.2 keeps omitted empty Instance.Profile and existing generator bytes. Existing generic
Go literal emission can reconstruct both Document and Root without a new generator.
Add strict preview-schema validation to the selected-root ResolveInteractions pass,
including hidden content. All consumers continue using that shared validation;
no separate lexical/selected-entry resolver. Normalize validates source schemas even
when a reusable template cannot yet resolve an absolute command reference.

Generated Document must normalize back to generated Root, retaining explicit versus
bare Markdown and SVG opt-in. Constructors include no renderer, resource, prepared
outcome or application binding. Source references never cause filesystem access or
SDL invocation; current bridge binds callback owners, not arbitrary source references.
Snapshot projection must preserve Text/Arguments equality for explicit Markdown,
and original SVG source/description/policy; caption projection must not become fallback
description. Do not put prepared status or bytes into SnapshotRoot to appease exporters.

## 3. Export and native background boundary

Smallest public SVG policy: reject every explicit preview node with
unsupported-resource-export, path and original span, even fallback="label", hidden,
with geometry, or with a supplied Content renderer. The reviewed contract permits
this boundary; no mandatory resource snapshot renderer or public fallback renderer.
Keep legacy SVG/bare Markdown export unchanged. All existing collection/pane/value/
text/interaction rejection rules stay independent. CLI preflight must precede renderer,
asset or output writes; failed export retains the prior artifact.

Source-only Dump/Markdown/Combined show preview kind, source reference or text,
description, fallback policy, original path/span and "resources not supplied".
They may quote source prose but must not call it a prepared rich preview or substitute
the legacy generic Mermaid/label rendering for a known prepared outcome. Hidden
declarations remain in the structural description. No new status-injection API is
needed: callers with no prepared bundle cannot report rendered/fallback success.

For native SkipControls, opted-in SVG may be omitted ONLY with the exact actually
prepared NativeControls[path]=="svg" entry. Both Check and Render currently exclude
all SVG from omission; change that branch only for Explicit=true after selection.
Legacy SVG remains painted as its old placeholder; an invented inventory entry
must not license omission. Check missing/extra/wrong-kind entries across the canvas
subtree and retain full InteractionRoot selected-entry validation.

Explicit Markdown remains content painted by the frozen provider through existing
svg.ContentRenderer; do not invent a Markdown control/inventory kind. One small
optional interface is proposed to prevent an arbitrary legacy Content implementation
from bypassing the explicit-node prepared-outcome requirement:

```go
// package svg; existing ContentRenderer remains unchanged.
type PreparedPreviewRenderer interface {
    ContentRenderer
    CheckPreview(*parser.Instance) error
}
```

Only native SkipControls may admit explicit Markdown through this interface.
CheckPreview is pure: require the matching path and frozen source/policy/text outcome,
including a valid labelled outcome; reject absent/mismatched records even under label
policy. Mutable caption/layout are not part of the source-policy-text fingerprint. It neither
prepares nor rerenders and cannot turn public export rejection into support. A
provider-owned Previews can implement it without new svg->markdown imports (Markdown
already imports svg). James and Gibbs agree with this proposed seam; host coordination
remains before implementation. All checks scan hidden/closed descendants; inventory
cannot mask a malformed node. No competing stub is introduced.
Its assertion is trusted prepared-host evidence, like the existing native inventory,
not an untrusted renderer claiming provider capability by interface presence alone.

## 4. Preparation ownership and concrete integration gaps

James proposes existing markdown package ownership of PreparedSVG/MarkdownRenderer
request inputs and immutable Previews outcomes, keyed by exact normalized path.
Gibbs proposes a geometry-only measurement adjunct for opted-in SVG plus centered
aspect fit within the allocated rectangle, applying Box.Clip afterward rather than
fitting to the visible clip; explicit Markdown retains Measure/Content.Render. Neither proposal needs
AST geometry or layout imports of Markdown. Frontend supplies source facts only.

Actual gaps requiring implementation when selected, not architecture contradictions:

* schemaFor currently rejects markdown calls and new SVG keys; normalization does
  not perform the call-to-Markdown conversion. Strict selected-root validation must
  reject malformed explicit nodes, not silently treat them as old content.
* markdown.Provider.Documents is keyed by text and Resources by diagram ID; Prepare
  skips equal text. Those caches cannot hold independent per-path descriptions,
  policies or outcomes. Shared parsed data is permitted; outcome authority is not.
* Custom Renderer output is currently retained directly. The unexported Mmdr
  validator is not a complete WCI4 supplied-shape validator: its denylist is broader
  than the selected allowlist, and single-document/trailing-input checks need work.
  Reuse bounded validation machinery without claiming current acceptance proves
  the stricter subset. Content failures follow label/reject after identity checks.
* DocumentHost.measureCanvases calls markdown.Prepare(snapshot.Root,nil). WCI4
  preparation must move before repeated gates/tickets; gates use frozen outcomes
  and no renderer. Merely adding renderer arguments to measureCanvases is incorrect.
* preparation.Check maps every SVG to svg-placeholder. Actual outcome-dependent
  facts must replace that assumption for opt-in only; fallback="label" is not itself
  evidence of provider svg-resource. No parser readiness flag resolves this gap.

Identity/digest/extra/wrong binding errors, aggregate budget and stale publication
stay fatal regardless of fallback. Validly bound malformed/unsupported content,
dimensions/per-resource budget and renderer errors use the node's declared policy.
No renderer execution or outcome reconsideration during layout/paint/StateGate.
This memo adds no runtime state/API, resource loader or generic framework.

## 5. Consumer readiness and proposed proof

Static sdui-preview uses Combined and remains source-only. Its sdptool/0.2 protocol
is distinct from source profile; preserve the .3 profile metadata and .2 envelope
bytes, snapshots, source links and stale-revision behavior. Standalone -check/native
launcher may explicitly reject these governed preview adapters with source-linked
diagnostics; no implicit document-host migration or new XFMD Launch claim. Name and
stage the actual connected DocumentHost fixture route in WCI4 package evidence.
Helper packaging, actual XFMD build/protocol, matching module roots/IME replacement,
binary licenses/hashes and OS proof belong to main/integration, not this memo.

Proposed targeted acceptance after selection:

* A01: full schema/type/presence negatives, @resource only on opt-in, description
  UTF-8/byte limits, raw multiline text, no-I/O module path, .2 rejection, explicit
  Markdown normalization/reuse/spans and malformed constructed Text mismatch.
* A01/A05: compile/run generated constructors; Normalize(Document)==Root; independent
  reuse; unchanged .2 and legacy .3 AST/Go/text/SVG bytes; no resource payload emission.
* A03/A04: same Markdown text at two paths with different policies/descriptions;
  missing/mismatched prepared outcome fails native preflight; repeated export checks
  call only pure CheckPreview, never Renderer. Provider lane proves identity-first
  failures, immutable bytes, per-node fallback and zero gate/ticket rerenders.
* A05: hidden explicit-node public rejection with supplied geometry/Content; exact
  opted-in SVG omission inventory and legacy placeholder retention; unsupported
  export preserves old artifact; structural output never invents prepared status.
* A06: SDUI helper/SDPTool static delegation and existing protocol/stale-source tests;
  separate unsupported standalone row from real staged connected-native proof.

No product tests run and no WCI4 capability/pass claim. No material conflict with
the reviewed source contract found; the optional prepared-content check requires
lane coordination before implementation. Only this memo changed; stage selection,
canonical/Session updates and independent review remain with main.
