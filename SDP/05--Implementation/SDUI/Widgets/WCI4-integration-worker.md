# WCI4-M1 host and preparation worker

Status: component implementation and scoped verification complete; lane files frozen.
The predecessor scoped host/preparation review approved its exact 15-file manifest.
A confirmed native owner-loss finding required the authorized lifecycle delta below;
independent delta review subsequently approved it; actual native reruns remain
with main.
No stage completion or package/publication claim.

## Authority and scope

Coordinator selected WCI4-M1 after M2 delivery 69d0a33 and phase handoff 90a94b5.
Worked only in `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci4`, baseline HEAD
`90a94b5daa4b36635be102858b2ce4a8012150c6`. Read the ORIGINAL selected
Providers-and-packaging contract (SHA-256
`91285ab840cc7ddab2c571a513899c969ee109781d7ff4ff276cfae254b7a4d8`), final four
implementation API handoffs, assignment audit and applicable AGENTS. Reused SDP
1.1.1, Worker 2.0.0 and shared document workflow. Main retains Session/canonical/
management, module/dependency/build roots, native/integration and final disposition.
No commits, branch switches, module or other-lane edits.

## Result

DocumentRequest accepts exact-path SVGResources/MarkdownRenderers. Detached host
preparation normalizes and prepares immutable content once before session/gates;
pure preparation independently checks the actual normalized-root fingerprint before
capabilities/session creation. Hidden pages/closed dialogs participate. Native facts
follow actual outcomes; placeholder fallback never implies svg-resource and this
native Markdown route never advertises rendered Mermaid. Original legacy Check
capability-origin behavior and all other family requirements remain intact.

Direct native SVG resource and a separate text-only glyph overlay consume the
provider's single SVGRects/band plan. Shared clipping follows aspect fit, including
negative positions after scrolling. Caption is one measured ellipsized row; complete
fallback status wraps in its independently admitted region. Complete overlay strict
decode/serialization failures reject the prospective ticket. Actual Markdown prose
stays in background content; its transparent native adapter supplies accessibility
without a fabricated export inventory kind. Existing accepted-presentation tickets
control native resource publication; rejected tickets preserve old images/captions.

The host caches diagnostic metadata and native immutable resources separately;
application request maps/renderer objects are not retained. Native resource Content
returns defensive copies. Close/reload releases this bundle's references. Surface
close/hide clears mounted preview resources; bundle declarations remain for reopen.
No immediate GPU/cache reclamation guarantee is claimed. Context/wheel routing uses
existing adapters, with no runtime preview state, renderer framework or new scene.

Inspect.previews reports prepared/mounted distinctions, actual native canvas/clip,
shared image/caption/status rectangles and the mounted fyne.Accessible methods.
Full description/status is independent of caption. Linux OS accessibility remains
unsupported; glyph coverage remains the existing Go-Regular boundary. Component
inspection is not OS screenshot/screen-reader proof.

## API coordination

James owns the actual markdown exports, including Previews and RenderSVGText;
Gibbs owns PreviewMeasurer/FitPreview; Dalton owns PreviewOptions and native/export
checks. No competing stubs were created. Lorentz received the actual request and
inspector shape and can use the connected preview fixture. Frontend resource-SVG
omission is backed by an actual prepared native adapter; Markdown remains the
prepared Content route, legacy SVG remains background-painted.

## Verification

All Go commands use cwd `/tmp/sdp-sdui-widgets/SDUI/go` and prefix
`GOWORK=off GOCACHE=/tmp/sdp-wci0-gocache GOFLAGS=-mod=readonly`.

| Command | Result |
| --- | --- |
| `go test ./host/fynehost ./preparation -run 'TestPreparedPreview\|TestPreviewFingerprint\|TestHiddenPreview' -count=1` | Initial integration PASS host 0.509s, preparation 0.005s. |
| `go test ./host/fynehost ./preparation -count=1` | Host PASS 10.449s; preparation found new source-schema check had changed one legacy capability-origin diagnostic. Fixed by applying PreviewOptions only at explicit argument presence; existing test retained. |
| `go test ./host/fynehost ./preparation -run 'TestPreparedPreview\|TestPreviewFingerprint\|TestHiddenPreview\|TestCapabilitiesInspectHiddenNodes' -count=1` | Final expanded targeted PASS host 1.097s, preparation 0.007s. |
| `go test -race ./host/fynehost ./preparation -count=1` | Predecessor PASS, exit 0: host 200.731s, preparation 1.067s. Its 15-file hashes matched; this predates the owner-loss delta below. |
| `git diff --check -- SDUI/go/host/fynehost SDUI/go/preparation` | PASS. |

Eight new meaningful tests cover: input/native-accessor copying and resource
release; actual accessibility and exact native inventory; caption/fit/overlay;
label-vs-rich facts, fatal digest and stale Guard; preparation-only renderer calls
with native diagram fallback; rejected-ticket caption/resource retention; ancestor
scroll clipping without refit; hidden resources, nonmodal canvas ownership and
close/reopen; fingerprint-before-capability/gate ordering and each hidden native
resource capability. The scroll test initially used invalid relative-layout syntax;
corrected the test source to existing scale-x/scale-y, with no product-rule change.

No broad SDUI suite, native binary build, OS input, package staging or install was
performed by this lane. Main owns those checks on the integrated frozen candidate.
Parallel provider/frontend/fixture results are not substituted for this lane's
reported commands. No network/sandbox block encountered.

## Candidate inventory

The following 16 product/test/doc files are the current host/preparation lane changes.
This report is the additional permitted evidence file. Manifest SHA-256:
`ba7d99e360f174f12c74ca090980d8616d1086cb61ffd334dee09e8e6e2d82ce`.
This identifies current lane bytes, not a committed or whole-repository candidate.

```text
5197c39a6f29d2e0c8d231904d6d78c26245f85b6fe900697d8eebfc5cd2f5af  SDUI/go/host/fynehost/README.md
514c602d8985a3e58830b6d5e6f596f82ec2687e4edfed6ce5c1e43f9f2f5891  SDUI/go/host/fynehost/collection_metrics.go
a2b5eaba7582fdf65a67ca5ad7f42fbc41a70724c2c2b68a6d3c0258451eb208  SDUI/go/host/fynehost/document.go
047dd7033307f8ec8450a5bd9ca6a673026f94f51ca1f580a9585c33d8784b8c  SDUI/go/host/fynehost/document_commands.go
3c6a7e08e4e75f308ccf066a2ce4ff942d756011898e56b5831915a1c8bd2fb7  SDUI/go/host/fynehost/document_previews.go
3c62b4cafa0e31fadadd94b7319c62ea8a71ab9dc5041c3a9a62e5d263f8db8b  SDUI/go/host/fynehost/document_surfaces.go
eac9c2b1994324dc4b81b8cb72ac52aaebf43fb517b5848e8b875b5e51a0f97e  SDUI/go/host/fynehost/document_sync.go
2b6dd893bc7150b704851b5420562cdf71eb31b6bb0ef085374f5291c6c605d8  SDUI/go/host/fynehost/inspection.go
cf35ffcf9d1b598ccffdb12560d9cdff09af3ae1b62322046fe8a3f73a446e13  SDUI/go/host/fynehost/preview_control.go
d41335452d8cd08ed81726c2bb95f5bafa0f96d60bca61e0a698e7fc47da0599  SDUI/go/host/fynehost/preview_inspection.go
19616bd936c73c133a1667ae4a7e3f1ac6007df9fec26fb9f80c0cb6946e27dc  SDUI/go/host/fynehost/preview_owner_test.go
524cc7aecb60a858d24a45b855262c45edc0b68f6a24dc0b2491515951b6dfe8  SDUI/go/host/fynehost/preview_test.go
df2c210b88503b2c94379cb62dfc9d9d2ef6801a636b7b36cf5c55ac2e862c2a  SDUI/go/host/fynehost/view.go
c540dcc61ba24789478206c3b413565de16516c7598ea6fd45d497966263601a  SDUI/go/preparation/capabilities.go
8abf082c3b9f33fdd7d9644530f11d81cd6ac7d62fa985272f0a415071aaa856  SDUI/go/preparation/prepare.go
9f50a14b4d91b18e7fc2e8287aa90e1265d349a32274991d0547c3e9c06cdfec  SDUI/go/preparation/previews_test.go
```

## Native owner-loss correction and new freeze

Main's actual `/tmp/wci4-pilot-forms3` evidence showed runtime detail Open=false and
one parent-hidden sequence-0 receipt, but Inspect.previews still claimed mounted
and visible. Source diagnosis: forced RevokeSurfaces deliberately bypasses gates
and apply; syncSurfaces destroyed the owner window/overlay but retained the prior
presentation's preview/canvas records and native image/raster references. Inspector
then used that stale frame and fell back to main-canvas origin. Normal source
Cancel already went through accepted presentation and did not reproduce it.

Main explicitly stopped its affected whole-SDUI run before authorizing this delta;
`/tmp/wci4-final-suites1/aborted.json` records that main-owned aborted run. The first
finished package/consumer artifacts remain historical candidate-1 evidence.
No worker edits occurred during the preceding read-only diagnosis.

The bounded correction retires only records whose native frame equals the exact
retired opening. It hides/empties its published preview controls, clears resource
and raster references, removes stale preview/canvas records from owned presentations,
and detaches the old native background/objects (Markdown pixels live there).
Frozen bundle declarations and empty control identities remain available for reopen.
A different/new native opening cannot be cleared by old-target cleanup. Inspector
also requires current Open, exact target, shown/non-retiring owner and matching
published native frame before reporting mounted accessibility/geometry.
No new presentation gate, provider call, source Render semantics or runtime API is
introduced; forced native owner loss remains ungated.

Only these five files changed relative to the prior approved lane bytes:

- `host/fynehost/document_previews.go`: exact-opening cleanup and reference release.
- `host/fynehost/document_surfaces.go`: invoke cleanup during owner retirement.
- `host/fynehost/preview_inspection.go`: actual live-owner/target check.
- `host/fynehost/preview_owner_test.go`: new modal/nonmodal hide/close regression.
- `host/fynehost/README.md`: authorized explicit light-theme acceptance limitation.

The report is the additional permitted evidence delta. Prior scoped approval of
manifest `5ab915b4886425a68c7427d1151d87e08cf6373222ba2efdcc791f4f44316f14`
and reviewer race 21.348s/1.060s apply to the predecessor only, not this correction.

Commands use the same cwd/environment prefix listed above:

1. Before fix: `go test ./host/fynehost -run '^TestPreparedPreviewNativeOwnerLoss$' -count=1`
   failed all four modal/nonmodal × parent-hidden/parent-closed cases, exit 1,
   host 0.830s. Raw output: `/tmp/wci4-host-owner-before.log`.
2. After fix: the same command passed, exit 0, host 1.887s. Raw output:
   `/tmp/wci4-host-owner-after.log`. A further background-reference assertion was
   included in the final race run below.
3. Final: `go test -race ./host/fynehost ./preparation -run 'TestPreparedPreview|TestPreviewFingerprint|TestHiddenPreview|TestCapabilitiesInspectHiddenNodes|TestDynamicSurfaceParentOwnsModalCanvas|TestNativeOwnerCloseRevokesModalChildWithoutGateOrAcceptance' -count=1`
   PASS, actual exit 0, host 54.553s and preparation 1.091s. Raw output:
   `/tmp/wci4-host-owner-race.log`.

The regression verifies SVG and explicit Markdown, actual resource/raster and
background/object detachment (not merely inspector masking), exact sequence-0
receipts, zero calls to a deliberately failing gate/resource hook, preservation
of the main preview, retained declaration resources, reopen with a fresh target
and duplicate obsolete native callbacks that cannot dispose the replacement.
The prior preview/foreign-root/capability tests and dynamic-owner lifecycle tests
also passed in this final selection. All 16 current hashes were rechecked at freeze.

## Remaining work

The lifecycle delta is frozen and independently approved in component scope;
main's exact integrated native/package/full-suite reruns remain. No component test
failure remains. WCI4-M2/final
delivery, package publication and OS accessibility are not claimed here.

The declared native acceptance profile is FYNE_THEME=light. Inherited fixed
caption/glyph colors have poor default-dark contrast; dark/system-theme contrast
is unverified. KB-SDUI-004 retains this fidelity backlog (main's EVT27 registration).
This explicitly bounded profile adds no theme implementation.

## Independent owner-loss delta review

Mendel independently verified all 16 hashes in manifest
`ba7d99e360f174f12c74ca090980d8616d1086cb61ffd334dee09e8e6e2d82ce`
and approved the five-path correction in component scope. Independent race covering
owner loss, dynamic parent and WM/modal-child lifecycle passed in 25.454s, exit 0.
A supplemental reviewer overlay directly invoked retireSurfacePreviews(old) after
reopen/apply, beyond obsolete callback guards: all four cases passed in 1.285s,
exit 0. Overlay evidence resides under `/tmp/wci4-provider-independent/` in
`owner-overlay.json` and `preview_owner_overlay_test.go`. No repository edits were
made by that review. These are attributed independent results, not worker reruns.
Final exact native owner-loss and whole-candidate/provider approval remain separate.
