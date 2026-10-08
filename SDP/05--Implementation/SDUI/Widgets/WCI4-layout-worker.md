# WCI4 layout worker handoff

Status: bounded layout implementation complete and frozen for integration/review.
Date: 2026-10-08. Lane: layout (Gibbs). Clone `/tmp/sdp-sdui-widgets`, branch
`sdui/widgets-wci4`, HEAD `90a94b5daa4b36635be102858b2ce4a8012150c6` plus the
nine listed uncommitted product/document paths. Other active lanes own their
concurrent changes. No commits, branch changes, module edits, canonical/management
writes, native process actions or M2 archive changes were made in this assignment.

## Authority and implementation

Reused SDP 1.1.1, Worker 2.0.0 and the shared document workflow. Recovered original
Session0010 S5 read-only. Read original selected Providers-and-packaging.md, the
four final WCI4 API memos under SDP/05--Implementation/SDUI/Widgets, assignment
audit and local instructions. Earlier memo language saying unselected describes
its historical preparation; the coordinator's WCI4-M1 selection governs this work.
Main owns Session, canonical projections, integration/package/native proof and
independent review. No material departure from the selected seam was found.

Implemented actual `layout.PreviewMeasurer.MeasurePreview(node,font,width)` with
separate natural/minimum sizes and pure `layout.FitPreview(Rect,Size)(Rect,error)`.
Coordinated actual exports directly with Dalton, James and Noether; consumed
Dalton's actual `parser.PreviewOptions`, without a competing helper or stub.
Natural/minimum measurement is independent; constrained natural size need not
exceed a required text minimum. Both use existing finite/nonnegative 1e7 bounds.

Only explicitly opted-in SVG takes the adjunct. Desired allocation uses natural
size for unassigned axes; pane solving uses measured minimum plus source minima.
Final arrangement checks the minimum at final width after source maximums/tracks,
so maximums cannot silently erase required caption/status. Missing active adjunct,
invalid metrics and provider errors reject. Existing source/scroll/pane geometry
and runtime presentation data remain authoritative. Hidden/inactive/closed nodes
are not measured; their source policy is checked, while resource readiness for
all declarations remains the provider's preparation responsibility.

Fit centers the aspect-preserving image in full canvas allocation and permits
upscaling. Zero-area allocation returns exact Rect{} after validating inputs.
Nonpositive intrinsic sizes, nonfinite/negative/out-of-bound sizes, nonfinite
positions and positive fits that underflow to no representable area reject.
Normalizing the aspect along its larger intrinsic axis avoids overflowing scale
for tiny positive dimensions. Shared scrolling translates only once; clipping is
applied afterward and cannot refit the visible fragment. There is no new offset,
scene, resource-byte validation, Markdown import, native font callback or renderer.
Provider SVGRects/RenderSVGText owns the three bands and exact text measurement.

The maintained layout contract and new `SDUI/docs/wci4-layout-api.md` document the
actual calls and bounds, preserved legacy path, final minimum check and native
proof boundary. No canonical architecture amendment is needed for these exports.

## Meaningful verification

All commands ran in `/tmp/sdp-sdui-widgets/SDUI/go` with Go 1.27.1 linux/amd64,
`GOWORK=off GOFLAGS=-mod=readonly GOCACHE=/tmp/sdp-wci0-gocache`.

- `go test ./layout`: exit 0, PASS 0.052s after the implementation seam landed.
- `go test ./layout -run 'Test(FitPreview|Preview)' -count=1`: exit 0, PASS 0.005s.
  Nine new tests cover asymmetric/portrait/upscaled/zero/tiny fitting; invalid
  nonfinite/negative/bounded inputs; natural/source allocation and required text
  minimum under both assigned dimensions and source maxima; missing/error metrics;
  source-relative split bounds; preserved exact legacy .2/basic .3 geometry;
  nested translations/extents/gutters/hits and clipping without refitting;
  inactive/closed geometry and invalid hidden policy; actual runtime presentation
  gates rejecting a nonmodal resize without state/ticket/publication mutation.
  An initial fixture compilation failure used duplicate local names; distinct
  fixture names fixed it. No product rule or expected outcome was weakened.
- `go test -race ./layout ./svg ./prototype -count=1`: overall exit 1 because the
  concurrently implemented provider did not yet define `projectMarkdown` referenced
  by markdown/preview_prepare.go. Layout PASS 1.406s and SVG PASS 1.051s under race.
  Prototype did not build; no prototype pass is claimed. James was notified;
  no provider files were edited. The retry below closes this transient compile gap.
- After James confirmed the actual provider exports compiled,
  `go test -race ./prototype -count=1`: exit 0, PASS 1.151s. This closes the
  previously blocked downstream check; all nine frozen layout product/doc hashes
  were reverified unchanged.
- `git diff --check -- SDUI/go/layout SDUI/docs`: exit 0.

These are component/synthetic geometry and real pure-runtime-gate tests, not OS
native painting, resource backend, text glyph, accessibility or SDL action evidence.
No test can prove arbitrary third-party Measurer implementations are pure; the
interface contract requires immutable prepared inputs, and provider/host tests must
establish zero renderer calls through repeated gates and tickets.

## Exact changed-file freeze

The nine paths below are the complete product/test/document lane change. This root
report is the tenth written path. SHA-256 values identify bytes rather than an
unchanged HEAD commit; parser/provider/host concurrent work is not included.

```text
a4bfe6e22d3328bc4e3e3db1f9cf1df841289705fca9fa462ab0d3ed5e39adc6  SDUI/go/layout/previews.go
821b0b653603c9c55af97c2dde5466fb7adb9d56314eef0b16ea3c52d899d2c2  SDUI/go/layout/collection.go
7864cce7438383d1b3efe252082125177a520e2fff9c6fcc3d3ed9c587b71fff  SDUI/go/layout/pane_minimum.go
7ef09a499113170db482cfde0da4cd29ad10b782dfa300182d2a7b093b6bb578  SDUI/go/layout/engine.go
229123bdcac81457b8e31e7b88bbf5352351710d1e076f74e3de84860b9454f1  SDUI/go/layout/viewport.go
67ef364b4bb4c2510c517859ab5645262508816c90eb70fd0f18610ec99cb407  SDUI/go/layout/previews_test.go
9adca5a88f07b75153d83b2a085eed0fd434fd85f36d564dfb32ff0db5d0856f  SDUI/go/layout/preview_snapshot_test.go
5aa4c48e823c1ef49dda8039a123cca33dd5a0ad053cb8fa52b961b1e199ea17  SDUI/docs/go-layout-contract.md
79e444a5e40d7c9346a1f8d82db4602016d2588509f68e144919c1027e9dbe0e  SDUI/docs/wci4-layout-api.md
```

## Remaining integration and Session handoff

- James owns immutable outcomes, caption/status band minimums, SVGRects and
  RenderSVGText and must retain full fallback status. Noether delegates the adjunct
  through native metrics and uses the fitted rectangle directly with shared clips.
- The previously blocked prototype race check now passes; integrated candidate
  verification remains main-owned.
- Main owns native positive SVG/prose/fallback/clip/accessibility/lifecycle proof,
  all-family regression, matching package preparation, final aggregate tests and
  independent Mendel review. No overall WCI4 delivery/release claim is made here.
- Session work summary: owner selected WCI4-M1; implemented the approved geometry
  seam in exclusive layout scope, preserved legacy behavior, verified the bounded
  component tests and froze hashes. Next: provider/host integration and independent
  exact-candidate review. Main records this in Session0010; worker did not edit it.
