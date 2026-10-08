# WCI4 provider worker — corrected source frozen for independent recheck

Role: SDP Worker 2.0.0, SDP entrypoint 1.1.1 and shared document workflow.
Selected WCI4-M1, provider lane only. Main owns Session/canonical records,
integration/native/package proof and final gate. No commit or branch change.

## Reproduced final-allocation defect and approved correction

Before the correction, from SDUI/go:

```text
GOWORK=off go test -mod=readonly ./markdown -run TestPreviewMarkdownFallbackRejectsTruncatedAllocation
--- FAIL: TestPreviewMarkdownFallbackRejectsTruncatedAllocation (0.00s)
    preview_geometry_test.go:179: truncated fallback published or output changed <nil>
FAIL
FAIL github.com/Hans-Einar/SDP/SDUI/go/markdown 0.009s
```

The test creates explicit whole-label and per-diagram partial outcomes, measures
at the final width, then renders into a box one unit shorter. The prior provider
returned success and wrote the truncated result. Layout's ordinary Markdown path
does not enforce the widget minimum at final allocation. Main and Gibbs approved
a provider-only pure projection guard for label/partial outcomes before ticket
publication. Legacy and fully rendered prose retain existing behavior. Native
pilot evidence before this correction identifies only that earlier candidate.

The source/test manifest and final verification follow after completion.

## Implemented provider lane

Base HEAD: `90a94b5daa4b36635be102858b2ce4a8012150c6`, branch
`sdui/widgets-wci4`. The implementation is the following uncommitted file manifest,
not an unchanged-commit result. Compiler: go1.27.1 linux/amd64. Main supplied the
selected canonical contract after M2 delivery; provider/API scope is unchanged.

- Exact agreed PreparedSVG, MarkdownRenderer, PreviewBackend, immutable Previews
  and detached Outcome exports. Source-only parser PreviewOptions owns opt-in.
- Global source/binding inventory and identity/digest preflight before callbacks or
  fallback. Freeze source metadata and validate/copy supplied resources one at a
  time before application callbacks, avoiding an unbounded staged input copy.
- Strict complete XML, namespaces/attributes/elements, finite numbers, shape/path
  grammar, literal arc flags, transforms including ancestor composition, solid
  paints, dimensions and security. Separate registered Mermaid output grammar.
  No modifications to the legacy Parse/Provider/Mmdr/WriteResources files.
- Per-path and per-occurrence outcomes, copied custom output/backend arguments/
  accessors/provenance/table rows. Whole-document label versus diagram-only partial
  fallback preserves ordinary prose. No provider callbacks or policy reconsideration
  in pure gates, geometry or painting; no runtime preview state.
- Exact 4 MiB resource and full-SHA-deduplicated 32 MiB admission tally before
  backend classification, including backend-labelled failures. No refund or live
  retained-memory claim. Valid registered diagram output is processed even when
  the whole Markdown backend is unavailable. Failed backend resources are absent
  from returned outcomes; their digest/diagnostic remains.
- Shared SVG natural/minimum and final text bands, full-allocation aspect fit,
  caption ellipsis, complete status wrapping, pure glyph-only overlay, and the
  reproduced explicit-Markdown label/partial final-allocation correction above.
  A pure natural-height calculation divides before multiplying to avoid needless
  floating underflow for very small finite resource dimensions.
- New package README documents actual API, ownership, legacy boundary, security,
  budget, native backend limits and accessibility boundary. Root API memo updated
  only for main/reviewer's selected budget clarification; main owns canonical copy.

No commits, branch changes, module/runtime/other-lane/canonical/Session writes.
Export coordination completed with Dalton, Gibbs, Noether and Lorentz. API shape
matches the reviewed handoff; main/Mendel approved the aggregate-budget refinement.
Gibbs and main approved the bounded Markdown final-allocation guard. No general
resource/renderer/transaction framework or new native scene was introduced.

## Verification and practical limits

Sixteen new TestPreview tests plus one fuzz target cover closed SVG/profile/security,
complete documents/dimensions, identity precedence, caller/custom-output/backend/
accessor copies, per-path and repeated-fence independence, callback mutation of later
input/source, table/provenance copies, required document bounds, legacy .2/.3 bytes,
concurrent pure access, renderer counts, geometry/clip/ellipsis/status allocation,
exact aggregate boundaries/deduplication and unsupported-backend budget behavior.
The legacy configured-Mmdr test remains environment gated by SDUI_MMDR.

Targeted post-fix tests passed in 0.030s:
`GOWORK=off go test -mod=readonly ./markdown -run 'TestPreview(SVGClosedGrammar|MarkdownFallbackRejectsTruncatedAllocation|SharedGeometryCaptionAndStatus|LegacyProjectionPreserved)'`.
`GOWORK=off go vet -mod=readonly ./markdown` passed with no output.
Final race/fuzz results are recorded below after completion; earlier race runs are
not represented as tests of the subsequent allocation/arc-flag changes.

Component evidence does not claim native pixels, arbitrary Unicode glyph coverage,
actual mounted Accessible wrappers, OS screen-reader delivery, native ticket
lifecycle, package installation or overall WCI4 acceptance. Host/main own integrated
proof; Linux OS accessibility remains unsupported/unverified. The selected native
Mermaid route is unavailable, not a rendered capability. Registered custom Go
renderers have no new timeout/cancellation wrapper. Independent provider review is
pending; this implementation report is not an independent approval.

## Initial candidate file manifest — superseded by finite-geometry correction

| Owned changed path | SHA-256 |
| --- | --- |
| `SDUI/go/markdown/preview_attributes.go` | `5919685de853fe24649b3214c1e3cb096e7469484258220402178bb0f35d503c` |
| `SDUI/go/markdown/preview_budget_test.go` | `a6f8a14f2173a852a72696b8e91a83885f69076b221b4062dfb52b3fd7084193` |
| `SDUI/go/markdown/preview_geometry.go` | `238548e085e563ebfced94084b2c24973ee16918de8f2fc8801dcd47b928104b` |
| `SDUI/go/markdown/preview_geometry_test.go` | `846102d16c948b9decd8ce87ab601450ed7ba7f90f96ea7775c1d2e1a8859566` |
| `SDUI/go/markdown/preview_identity.go` | `f111bd2e81fa9da9fc89dfc9482f5d67fb453a90bd477694795974c1f8e9998f` |
| `SDUI/go/markdown/preview_numbers.go` | `bec22d6c823e7b8d5e06114ce4af4a2c1c668dc5a01ab38419f79e43ac515634` |
| `SDUI/go/markdown/preview_prepare.go` | `a6fe54ab839075dafeeb4b401ae8c609486b4bce6286bea504c03b2384867239` |
| `SDUI/go/markdown/preview_prepare_test.go` | `32766f865a6d04abd418d9fdf95f833fe7309247afe81cf69b8cbe875ce62708` |
| `SDUI/go/markdown/preview_projection.go` | `1a3a312ec703621836ed70d7dff305717525974dd005a7f2d508ceda1b309b4e` |
| `SDUI/go/markdown/preview_types.go` | `f3104ea9ce79683f0a4f7b237ab1de19dc4a72546f64919a730176ce89802e8f` |
| `SDUI/go/markdown/preview_validation.go` | `1d86c80e57302b595ba6544707f93b4559492ca635413028c8d514d82f9783a9` |
| `SDUI/go/markdown/preview_validation_test.go` | `ba8fa33cc40eb7e5016bca88da165dbce6b99d855d26f9a092daebd188609f9c` |
| `SDUI/go/markdown/README.md` | `fd9ef36063dd8655dcfbe671aa9572e5377877e3c7a9e2840c26bb0e20a05af8` |
| `WCI4-provider-API.md` | `02af76c946d7adf4b64b3ee7579fdebb4e47a50178544838448bba9ccc9ec673` |

Canonical ORIGINAL observed SHA-256: `91285ab840cc7ddab2c571a513899c969ee109781d7ff4ff276cfae254b7a4d8`.

## Initial verification and freeze — subsequently revoked by review

All commands ran from `/tmp/sdp-sdui-widgets/SDUI/go` on the manifest above:

```text
GOWORK=off go test -mod=readonly -race ./markdown ./layout ./preparation ./svg
ok github.com/Hans-Einar/SDP/SDUI/go/markdown    71.709s
ok github.com/Hans-Einar/SDP/SDUI/go/layout      (cached)
ok github.com/Hans-Einar/SDP/SDUI/go/preparation 1.190s
ok github.com/Hans-Einar/SDP/SDUI/go/svg         (cached)

GOWORK=off go test -mod=readonly ./markdown -run '^$' -fuzz FuzzPreviewSVGAdmission -fuzztime=3s -parallel=2
PASS: 22142 executions; 44 new interesting inputs; 67 total corpus entries
ok github.com/Hans-Einar/SDP/SDUI/go/markdown 3.309s

GOWORK=off go vet -mod=readonly ./markdown
PASS, no output
```

The 3-second fuzz run is bounded supplemental evidence, not an exhaustive security
proof. Race coverage includes repeated concurrent projections and the real 32 MiB
boundary fixtures. Gofmt reports no unformatted package files. All four legacy
Markdown source/test files compare byte-for-byte with HEAD. The candidate manifest
was checked again after the final runs and all 14 entries match.

Provider source is frozen and ready for independent review/integration. No known
provider implementation blocker remains. Independent review and actual native/
all-family/package evidence are pending with their owners; do not label the whole
WCI4 stage delivered from this report. Main can rebuild only after incorporating
this explicit-Markdown allocation correction and the exact current candidate.

## Independent finite-geometry finding and bounded correction

Mendel reproduced three admitted nonfinite shapes on the initial 14-file freeze:

```text
<g transform="matrix(1e308 0 -1e308 1 0 0)"><line x1="1" y1="-1" x2="0" y2="0"/></g>
<rect x="1e308" y="0" width="1e308" height="1"/>
<path d="M0 0 Q-1e308 0 1e308 0 T0 0"/>
```

Both the old provider validation and strict backend returned nil in the reviewer's
reproducer, while the first actual transformed endpoint was +Inf. Independent
files: `/tmp/wci4-provider-independent/finite.go` and `finite.log`.
Main explicitly authorized this bounded provider correction and stopped the earlier
candidate build/suite. The initial freeze above is historical and superseded.

Before the fix, the new owned regression command failed in 0.007s:
`GOWORK=off go test -mod=readonly ./markdown -run 'TestPreview(ActualFiniteShapeGeometry|FiniteGeometryPreservesValidControls)'`.
Observed assertions included `admitted nonfinite derived geometry`, `invalid geometry
reached backend`, `invalid geometry became rendered`, `reject policy admitted invalid
geometry`, and a valid line rejected because each scalar was checked as (n,n).

The correction removes the fictitious (n,n) point. It validates actual line endpoint
pairs, rectangle corners/extents, local ellipse extents and exact transformed ellipse
axis extents. Stroke width is a linear vector magnitude, not a translated point.
Smooth S/T commands derive and check reflected prior controls, with correct previous-
command resets, relative coordinates and repeated operands. Reflection uses x+(x-c)
to avoid a spurious 2*x overflow when the actual result is finite.

Arcs now validate actual endpoint-to-center/radius-correction arithmetic, finite
local/transformed axes and extrema traversed by the selected sweep. Zero-radius
line and coincident-endpoint empty arcs retain SVG behavior. Nonfinite numerical
intermediates reject before backend callbacks; there is no arbitrary coordinate
cap or assumption that strict backend parsing catches them. Unused portions of the
ellipse do not create synthetic swept-extrema failures. This remains a bounded
finite numerical profile, not a full SVG renderer or an arbitrary-precision claim.

Tests cover the independent repros, rectangle/circle/ellipse/transformed-arc/stroke
and reflected-control overflow, label versus reject with zero backend calls, ordinary
rotated/transformed shapes/arcs, relative and repeated curves, control resets,
zero/coincident arcs and valid large finite coordinates. These are required regressions
for the confirmed review blocker, not a new optional family or architecture scope.
The package README now describes derived finite geometry validation.

## Corrected candidate verification and final freeze

All commands ran from `/tmp/sdp-sdui-widgets/SDUI/go`, same base HEAD and toolchain:

```text
GOWORK=off go test -mod=readonly -race ./markdown -run 'TestPreview(ActualFiniteShapeGeometry|FiniteGeometryPreservesValidControls|SVGClosedGrammar|CompleteDocumentAndDimensions|MermaidProfileSeparate|CustomRendererValidationAndLimits|IdentityBeforeFallbackAndCallbacks|SharedGeometryCaptionAndStatus|MarkdownFallbackRejectsTruncatedAllocation)'
ok github.com/Hans-Einar/SDP/SDUI/go/markdown 1.278s

GOWORK=off go test -mod=readonly ./markdown
ok github.com/Hans-Einar/SDP/SDUI/go/markdown 4.192s

GOWORK=off go vet -mod=readonly ./markdown
PASS, no output
```

The corrected candidate has eighteen TestPreview tests and the existing fuzz target.
No optional additional testing/scope was started after these necessary checks passed.
The earlier full affected race and fuzz runs identify the initial candidate; the
commands above verify this delta. Main owns the fresh whole-SDUI/SDL and package/
native final candidate runs. Gofmt is clean and all four legacy files still match
HEAD byte-for-byte. No runtime, host, layout, module or canonical source changed.

Current provider source is frozen for Mendel's independent recheck. The exact
current 17-file manifest below replaces the initial manifest. No known provider
blocker remains after the reproduced correction; no independent approval or native/
whole-stage delivery is inferred. Main and dependent lanes have been notified.

| Current owned changed path | SHA-256 |
| --- | --- |
| `SDUI/go/markdown/preview_arc.go` | `92d4bb45593ed7443a724eb0268093d19a4c850c4ac462ee0649ae31213ef4c5` |
| `SDUI/go/markdown/preview_attributes.go` | `685c0044474298b981acffd3791e88d83afa8d753a879eaa13f364f8dd48f219` |
| `SDUI/go/markdown/preview_budget_test.go` | `a6f8a14f2173a852a72696b8e91a83885f69076b221b4062dfb52b3fd7084193` |
| `SDUI/go/markdown/preview_finite_test.go` | `1fcc9524374d545e992f95523229088f832471f4fed6e2c18baef66dbac18e6d` |
| `SDUI/go/markdown/preview_geometry.go` | `238548e085e563ebfced94084b2c24973ee16918de8f2fc8801dcd47b928104b` |
| `SDUI/go/markdown/preview_geometry_test.go` | `846102d16c948b9decd8ce87ab601450ed7ba7f90f96ea7775c1d2e1a8859566` |
| `SDUI/go/markdown/preview_identity.go` | `f111bd2e81fa9da9fc89dfc9482f5d67fb453a90bd477694795974c1f8e9998f` |
| `SDUI/go/markdown/preview_numbers.go` | `44c58c0e901d279454d5bdb0a96e007791c8d0aae71810280311659408ccdfd9` |
| `SDUI/go/markdown/preview_prepare.go` | `a6fe54ab839075dafeeb4b401ae8c609486b4bce6286bea504c03b2384867239` |
| `SDUI/go/markdown/preview_prepare_test.go` | `32766f865a6d04abd418d9fdf95f833fe7309247afe81cf69b8cbe875ce62708` |
| `SDUI/go/markdown/preview_primitives.go` | `74274df9e51845b76b97f3fd7584129c9a186f6ab140553adde8ee0f05b14ff6` |
| `SDUI/go/markdown/preview_projection.go` | `1a3a312ec703621836ed70d7dff305717525974dd005a7f2d508ceda1b309b4e` |
| `SDUI/go/markdown/preview_types.go` | `f3104ea9ce79683f0a4f7b237ab1de19dc4a72546f64919a730176ce89802e8f` |
| `SDUI/go/markdown/preview_validation.go` | `1d86c80e57302b595ba6544707f93b4559492ca635413028c8d514d82f9783a9` |
| `SDUI/go/markdown/preview_validation_test.go` | `ba8fa33cc40eb7e5016bca88da165dbce6b99d855d26f9a092daebd188609f9c` |
| `SDUI/go/markdown/README.md` | `f50ec0ba41f13c9d80db7a3bfdb96e3dfe324bf884cb1204c08a559effd75e58` |
| `WCI4-provider-API.md` | `02af76c946d7adf4b64b3ee7579fdebb4e47a50178544838448bba9ccc9ec673` |
