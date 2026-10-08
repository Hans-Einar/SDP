# WCI4 providers and package preparation — provisional stage contract

| Field | Value |
| --- | --- |
| Assignment | SDP Architect; draft only, 2026-10-08 |
| Status | Provisional; requires WCI1–WCI3 pilot reconciliation and independent review before code |
| Authority | KB-SDUI-003 full inventory; PLAN-SDP-0022 WCI4-M1/M2; Session0010 |
| Parents | [Design](Design.md), [Acceptance](Acceptance.md), [Preparation](Preparation.md), [Collections](Collections.md), [Panes](Panes-and-commands.md), [Values](Values-and-text.md) |
| Outcome | Truthful SVG/Markdown previews, explicit provider/resource ownership and matching all-family producer/helper package preparation |

## 1. Boundary and inspected baseline

This stage completes the preview/provider and distribution-preparation rows; it does not implement an editor,
fetch URLs, add a generic resource system, rewrite XFMD, install packages or publish a release.
Source remains the unreleased `sdui 0.3;` profile with AST `sdui-ast/0.3` and generator `sdui-go-model/2`.
Preserve frozen 0.2 output/APIs and WCI1–WCI3 contracts; no extra version per family.

Observed in the original workspace, including its existing uncommitted integration files:

| Actual owner | Existing behavior and consequence |
| --- | --- |
| `SDUI/go/parser` | SVG has symbolic source Reference and optional label; Markdown is a string component. Neither resolves resources. |
| `SDUI/go/markdown` | Goldmark-derived bounded blocks/tables; 32768 bytes, 256 blocks, eight diagrams. HTML/images reject. Prepare accepts an optional Renderer; Measure/Render perform no I/O. |
| `markdown.Mmdr` | Explicit executable, stdin source, five-second deadline, 4 MiB output cap; flowchart/graph only. Resource validation rejects active/external SVG. Nil renderer currently gives a labelled Mermaid placeholder. |
| `SDUI/go/svg`, Fyne admission | Ordinary svg is a labelled placeholder; not vector execution. Native background uses measured geometry and exact prepared-control inventory. |
| `preparation`, DocumentHost | Detached checks and application-owned resources; state/geometry and publication must use the WCI2 finalized-snapshot ticket contract, not the older WCI1 speculative pending path. |
| `SDPTool/sdui.go` | Delegates discovery/combined static preview to SDUI; source snapshots and revisions are consumer boundaries, not runtime readiness. |
| `SDUI/go/cmd/sdui-preview`, `sdui-fyne` | Private helper sources live here; static/check and desktop native launcher are distinct. Standalone launch cannot invent collection providers or SDL bindings. |
| `SDPTool/package.sh` | Builds only sdptool, its build manifest and checksums into a new directory. It does not package Fyne helpers. |
| `Toolkit` | Process/template installation and shell CLI distribution; its CLI installer copies scripts without compilation. This is not an SDUI native dependency manager. |

The KB005 producer plan describes privately bundled XFMD helpers. No external XFMD checkout/build was inspected;
that historical arrangement proves neither its current module versions nor 0.3/native capability.

## 2. Small opt-in source contract

Keep bare Markdown strings and legacy `svg(source, label?)` unchanged, including their explicitly described
legacy placeholder behavior. Add only the following 0.3 local schemas; ordinary call grammar already suffices.

| Form | Exact semantics |
| --- | --- |
| `svg(source, label?, description, fallback)` | Presence of description/fallback opts into resource preview; both then required. Source is a declared-module Reference with member `resource`, e.g. `art.Chart.@resource`. Optional label remains a short caption. |
| `markdown(text, description, fallback)` | New call form for explicitly governed preview content; text is a string, not a file/resource URL. Description and fallback required. |

At most one positional argument (source or text), before named arguments. Reject duplicate/unknown fields,
wrong types, callbacks, malformed references and unsupported profiles. Description is nonblank valid UTF-8,
at most 4096 bytes; text retains the 32768-byte Markdown limit. Fallback is exactly `"reject"` or `"label"`.
Description without fallback, or fallback without description, rejects; no implicit opt-in or enrichment of old source.

```text
sdui 0.3;
ref: art "application-resources";
Main=[
  chart=svg(art.Chart.@resource,description="Monthly production by site",fallback="reject");
  explanation=markdown("# Production\n\nMeasured totals.",description="Production summary",fallback="label")
];
```

The module path remains symbolic; parser/codegen never open it. `@resource` is not an SDL action or callback.
Source AST calls use existing widget/Arguments fields. Normalize explicit markdown calls to the existing
`Kind="markdown"`, empty Widget, Text equal to the text argument, and retain typed Arguments for text,
description and fallback. Bare strings keep empty Arguments. SVG retains Kind widget/Widget svg and its arguments.
No new AST fields, generic annotations or resource-expression grammar. Generated constructors preserve these
distinctions, spans/use paths and metadata; no resource bytes, renderer objects or live registrations in generated Go.

This opt-in avoids tightening already accepted 0.3 placeholder calls and makes fallback reviewable in source.
No named style/theme/font registry, Markdown image syntax, imported files or navigation behavior is selected.

## 3. Application resources and preparation

Use the existing `markdown.Resource{SVG, Width, Height}` value as the shared immutable prepared image payload.
Add one concrete application input, keyed by exact normalized SVG instance path:

```go
type PreparedSVG struct {
    Source parser.Reference // module/object/member must match the declared source
    ProviderID string       // nonempty implementation identity, not a dynamic service locator
    SHA256 string           // full lowercase digest of Resource.SVG
    Resource markdown.Resource
}
type MarkdownRenderer struct {
    ProviderID, Revision string // nonempty application-owned implementation/configuration identity
    Renderer markdown.Renderer // existing Render(diagramSource) (Resource, error); nonnil when bound
}
// Document-host request inputs:
// SVGResources map[string]PreparedSVG
// MarkdownRenderers map[string]MarkdownRenderer // exact normalized explicit-markdown instance path
```

Reused instances require explicit path bindings; payload bytes may be shared internally only after defensive
copying. Missing Markdown renderer binding means diagrams are unavailable, not an implicit global renderer lookup;
plain supported Markdown needs none. A supplied binding must address an explicit Markdown node containing diagrams.
Keep legacy bare-string/global-renderer APIs unchanged. Pass only each fenced diagram's source to Render, never the
description or fallback policy. Source cannot select an executable. The request owns no renderer output until preparation.

First validate binding identity: extra/wrong-kind/unused bindings, mismatched source triples/digests, empty identity
fields or nil bound renderers are fatal regardless of fallback. Then validate content under the per-node policy in §4.
No runtime loader, event bus, async provider protocol or new SDL transport. Application owns acquiring bytes and
registering any existing local Mermaid renderer before preparation.

Bound each resource to 4 MiB and prepared distinct bytes to 32 MiB per candidate. Require a single complete SVG
document, finite positive viewBox width/height <=32768, and matching finite supplied dimensions. Reuse the bounded
SVG validation in markdown.Mmdr for all supplied/custom-renderer results; do not trust arbitrary Renderer output.
Directives, active content, external references/CSS and malformed/trailing documents are content failures: never
render their bytes; choose that node's label/reject policy after binding identity checks. This is a preview
subset, not a claim of full SVG conformance: unsupported backend features must reject or take the declared fallback.
For supplied resource SVG, initially allow only svg/g/path/rect/circle/ellipse/line/polyline/polygon/title/desc;
root viewBox/namespace, finite shape geometry, path data, transforms, solid fill/stroke, stroke width and opacity.
Treat other elements/attributes, CSS, text/fonts, gradients and external references as unsupported content under
that same policy, never silent backend omissions. The existing registered Mermaid profile remains separately verified, not licensed by this
shape subset. No automatic URL/image/font retrieval. Inspect actual native/export fixtures for each advertised
route; broader glyph/fidelity work remains KB-SDUI-004 and cannot be claimed by merely accepting XML.

Prepare all declarations, including hidden pages and closed surfaces, before activation. Freeze bundle-owned
prepared outcomes keyed by normalized instance path, for both SVG and explicit Markdown. Each record captures
description, fallback policy, source identity and the rendered/label outcome with its bounded diagnostic. Markdown
records retain ordered parsed blocks and per-diagram-ID rendered Resource or label diagnostic; identical text at
different paths cannot share a policy, description or outcome record. Parsed documents and validated immutable bytes
may share internally; old text-keyed Documents/Resources caches are not authoritative for these per-instance decisions.
Copy custom renderer output, validate it, compute its digest and retain the result; do not keep mutable renderer-owned
buffers. Renderer objects stay preparation inputs, not paint-time dependencies. Measurement, gates and native tickets
consume only these frozen records; repeated geometry probes/resize do not rerender or reconsider fallback.

Capture source/model,
provider identity and byte digest with the candidate; its application Guard rechecks current resource identity
and Markdown renderer Revision before publication. A changed resource or renderer identity requires a new guarded
candidate even if the source hash is unchanged.
Retain resources for that bundle only. Failed/abandoned candidates dispose only their prepared native objects;
replacement/disposal revokes old callbacks and releases old references after successful publication.

Measurement uses immutable prepared dimensions and existing relative layout; preserve aspect ratio, fit within
the allocated content rectangle, center unused space and apply the shared clip. Fallback uses measured text.
Reuse Fyne image/content controls and the existing background/provider route; no custom scene or independent hit geometry.
Never invoke a renderer or resolve a resource from Measure, paint, resize, StateGate or the non-yielding commit.
Prepare/cache content before pure measurement; WCI2 resource tickets receive exactly the finalized snapshot.

## 4. Fallback, capabilities and exports

`reject` retains the current bundle and reports source path/span plus provider/resource identity when content
cannot be prepared or represented. `label` may substitute only a visible description plus an English status:
`Preview unavailable: <description> (<bounded reason>)`. Status cannot be suppressed by a nonempty caption.
Retain the failure in the preparation report, distinguish fallback from rendered content, and expose the same
description/status to native accessibility. No silent success, empty image or rich-preview readiness claim.
Precedence is identity, then content, then publication: source/schema/profile errors, extra/wrong/unused bindings,
source/provider identity or supplied digest mismatch, aggregate candidate-budget failure, and stale publication are
fatal for the candidate, even with label policy. An incorrect digest cannot be disguised as malformed-content fallback.
With valid bindings, missing provider, malformed/unsupported supplied SVG or custom renderer output, invalid resource
dimensions/per-resource limit, unsupported content/backend and renderer failure use that node's label/reject policy.
Invalid bytes are never rendered; a labelled content failure does not excuse an identity error elsewhere in the candidate.
Keep ordinary Markdown prose in a valid document; an unavailable diagram gets the description/status at its block.
If the document itself is unsupported (e.g. HTML/image syntax), label replaces the whole explicit preview.

Use exactly frontend/layout/widget/viewport/provider/host dimensions; exact major matches remain mandatory.
Existing WCI1–WCI3 family capabilities remain independent and are not implied by this table:

| Prepared outcome | Required facts in addition to frontend sdui/0.3 and layout relative, major 1 |
| --- | --- |
| Actual resource SVG | widget svg/1, provider svg-resource/1, layout preview-resource/1; native additionally host svg-resource/1 |
| Explicit bounded Markdown | provider markdown/1; native additionally host markdown/1 |
| Rendered Mermaid block | provider mermaid-flowchart/1 plus bounded Markdown facts; only registered verified flowchart/graph output |
| SVG labelled fallback | Existing widget svg-placeholder/1 and provider svg-placeholder/1, reporting fallback rather than svg-resource capability |
| Markdown labelled fallback | provider markdown/1, with explicit unavailable status; no mermaid-flowchart fact inferred |

Resolve actual prepared outcome before advertising its facts; a declaration of fallback is not a provider registration.
Missing native capability can select an authorized labelled fallback, otherwise rejects. No bridge/resource/export
dimension, generic provider discovery or readiness boolean. Existing bare-content/0.2 capability behavior stays unchanged.

Static composition/text show source reference, description, policy and rendered/unprovided/fallback status only
when known; structural source-only output says resources not supplied. Public SVG may reject rich preview with
`unsupported-resource-export`, or emit the explicit labelled fallback; it need not gain a general image snapshot
renderer. A provider-backed SVG export is supported only with inspected positive output for its declared subset.
Existing collection/pane/value/transient rejection rules remain; this stage cannot export them by hiding controls.
Native SkipControls may omit an opt-in resource SVG only for its exact actually prepared native inventory entry;
legacy placeholder SVG stays rendered as before. Validate all descendant controls and unknown kinds as in M1.
All export failures preserve the prior artifact; asset writes use the existing owning staged bundle, not paint-time files.

## 5. Matching producer/consumer package preparation

Prepare one reviewable candidate set from inventoried source, not an inferred installed environment:

| Build/root or artifact | Required preparation evidence |
| --- | --- |
| SDUI/go | Parser/layout/runtime/providers/host and generated-model compatibility; sdui, sdui-preview and desktop-tagged sdui-fyne built from the same candidate. |
| SDL/go | Its matching SDUI replacement, action runtime/bridge/codegen and all-family connected fixture built/tested against that exact candidate. |
| SDPTool | Matching SDL/SDUI replacements and existing bootstrap module; discovery/UIPreview/profile metadata and source-snapshot compatibility; run existing package.sh into a new staging directory. |
| Private helper payload | Separately stage prebuilt sdui-preview and sdui-fyne, checksums, module/build facts and required native library/license inventory. SDPTool/package.sh is not assumed to include them. |
| Toolkit | Validate relevant process/distribution descriptors remain truthful; no widget binaries inserted into process templates or shell-script installer. No Toolkit version bump inferred from this draft. |
| External XFMD package | Matching payload handed to its package owner using its actual supplied build contract; record uninspected/unavailable roots as pending, never guess paths, versions or capabilities. No application rewrite. |

Use `GOWORK=off`, explicit build roots, locked go.mod/go.sum/replacements, `-mod=readonly` and a recorded compiler,
OS/architecture, tags and native dependency inventory. Inspect `go version -m` and hashes of each built binary.
Build in a new staging directory; no PATH edits, installed-helper replacement, network fetch during viewing or
compilation on launch. Provider executables are separately explicit registrations; do not silently bundle/find mmdr.
Package evidence records source commit plus dirty-file hashes where needed, generator/profile identities, module
graph/replacements, binary hashes, helper command/protocol fixtures and supported/fallback/unsupported outcomes.
Use a small candidate evidence manifest, not a new runtime discovery protocol or installer framework.

Exercise every family through discovery, normalized model, generated constructors, structural text and native
or explicit export rejection. Standalone helpers without application providers/bindings must reject/report that
fact; all-family connected evidence uses the real application fixture, not fabricated empty providers or a hidden input.
Keep protocol envelopes distinct from source profile: `sdptool/0.2` helper protocol does not mean `sdui/0.2` source.
Verify original source snapshots, UTF-8 source links, stale revisions, paths with spaces and exact profile metadata.
Compare static/readiness protocols with existing 0.2 clients; any necessary protocol change requires review before code.

## 6. Conditional IME dependency policy

The [IME probe](ime-probe/README.md) records a bounded X11/GLFW filter trial, not selected product support.
If WCI3 review selects that patch, preserve the exact upstream version, minimal patch, licenses and hashes in the
candidate. Explicitly select the replacement at each maintained native build root: at least SDUI/go and SDL/go,
and any actual consumer root compiling Fyne helpers. Go dependency replacements do not propagate transitively.
Check the resolved module source, compiled binary build info and an actual composed-text fixture from each distinct
native build recipe; a generator import or successful build does not prove the patch is present or IME works.
Other platforms retain upstream behavior unless independently selected/tested. Do not patch the global module cache.
If the policy is unselected, record upstream dependency identity and IME gap; neither package nor WCI4 may claim
full text-family acceptance while the required WCI3 evidence is pending. External roots remain pending until supplied.

## 7. Acceptance and handoff

| Pending ID | Exact proof |
| --- | --- |
| WCI4-A01 | Opt-in schemas/normalization/generated constructors; symbolic no-I/O references; 0.2 and legacy 0.3 bytes/behavior retained; unknown/duplicate/type errors and missing description/policy reject. |
| WCI4-A02 | Real supplied SVG and bounded Markdown preview with inspected native output, descriptions and shared relative geometry/clips; supported subset explicitly named; no full SVG/editor claim. |
| WCI4-A03 | Identity/digest/extra/wrong binding failures stay fatal under label policy; malformed/unsupported supplied SVG and custom output follow per-node label/reject consistently. Include missing/failed renderer, hidden nodes, unsupported HTML/images and limits; prior bundle retained on rejection; visible/accessibility fallback never claims rich readiness. |
| WCI4-A04 | Equal Markdown text at two paths with different descriptions/policies has independent immutable outcomes; parsed bytes may share. Count renderer calls through repeated gates/resize/tickets (zero rerenders); mutated returned buffers cannot change the bundle. Reload/disposal/resource or renderer-revision changes and guard failures retain publication/identity rules. |
| WCI4-A05 | All-family AST/profile/codegen/discovery/text and explicit export support/rejection; exact native inventory; failed export preserves artifact; generated 0.2 bytes unchanged. |
| WCI4-A06 | Matching module/build-root/private-helper inventory, clean staging, checksums/build info/licenses/native dependencies; actual packaged helper commands and 0.2/0.3 protocol/stale-source fixtures pass. |
| WCI4-A07 | If selected, exact IME patch at every native build recipe plus real configured IME composition/commit/cancel and no preedit SDL actions; unselected/unverified roots remain pending. |
| WCI4-A08 | Exact-candidate SDUI race suite, affected SDL/SDPTool suites and all-family native workflows; independent review closes every matrix row or records an explicit unresolved obligation. No publication inferred. |

Main records Session0010 and stage/matrix disposition. This draft changes only this file; no WCI4 code,
dependency selection, package production, external checkout mutation, publication or all-family completion is authorized here.

Review handoff: corrected independent review of ab4b1ed with explicit per-instance Markdown renderer inputs and
immutable prepared outcomes, plus identity-first fatal failures versus per-node content fallback. Architect context
reused; only this draft revised. Main owns Session0010 and requests re-review; no implementation approval inferred.
