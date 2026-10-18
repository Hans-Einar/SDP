# WCI3-M1 layout Worker handoff

## Assignment and authority

The coordinator selected WCI3-M1 only: checkbox, slider, select and number. The
reviewed canonical ORIGINAL Values-and-text contract and phase plan govern this
work; extended input/text/IME remains WCI3-M2. Canonical contract SHA-256 read:
`719cb7eab8375f20bbedb994fc0dc2eb40928b43903d5270238ce07629220309`.
Frontend/runtime WCI3 memos were read, then actual exports used once available.
The reviewed numeric admission and closed-dialog successor refinements remain
owned by runtime/frontend; layout supplies neither numeric parsing nor state.

Clone: `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`.
Base HEAD: `a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625`.
Evidence below concerns the uncommitted working candidate plus concurrent peer
implementations, not an unchanged HEAD or a completed WCI3 phase.

Only the twelve files listed below and this requested report were written. No
host inspection, host core, runtime, numeric, parser, bridge, fixture, module,
management or Session files were edited; no commits or branch changes occurred.
SDP 1.1.1, Worker 2.0.0 and document-workflow were reused. Session0010 S4 was
recovered read-only; main owns its journal and canonical integration evidence.

## Delivered behavior and coordinated API

Noether explicitly agreed the bounded optional `FieldMeasurer` seam before
implementation. Its actual exported signatures and semantics are recorded in
[the API note](SDUI/docs/wci3-m1-layout-api.md) and
[the maintained contract](SDUI/docs/go-layout-contract.md).

- `MeasureField(instance, runtime.FieldState, font, outer)` returns native
  `FieldMetrics` with `Minimum`, `Label`, `Control`, `Feedback`, `Decrement` and
  `Increment`. Label/control/feedback chrome is counted once; native host policy
  reserves a fixed feedback row. Intrinsic probes consume only Minimum. Final
  assignment validates all regions and rejects insufficient native space even
  after source maximum bounds.
- `SnapshotLayout.Fields` returns per-canvas `FieldLayout` rectangles, individual
  effective clips and Enabled/ReadOnly metadata. It is omitted from JSON when
  empty. Number Control identifies the Entry separately from two distinct step
  buttons. Checkbox's native label may overlap the combined control. Feedback
  cannot overlap label/controls and is mandatory when validation is present.
- Existing snapshot/canvas, scroll, fr/scale, pane minima, hit testing and
  EnsureVisible paths handle these leaf widgets. Each part receives ancestor
  translation once. Hidden/inactive pages and closed surfaces contribute no
  geometry. Actual open surfaces retain independent local coordinates and
  accepted nonmodal sizes; no additional viewport/value authority is introduced.
- Runtime Snapshot.Fields is joined by exact normalized InstancePath. Anonymous
  and reused controls can have different public Handle.Path spelling; layout
  does not implement another public-path resolver. Field identity/kind and native
  measurement are required for active scalar geometry. The source-only Layout
  entry has no typed field snapshot and rejects active scalar layout explicitly.
- Metric arguments copy RawDraft, Numeric and Options so retaining or changing an
  adapter's field argument cannot mutate the input snapshot. Source lexical
  number arguments stay untouched, including invalid live raw drafts.
- Existing PresentationState remains offsets and split bounds. Scalar geometry
  uses the same final ticket. Failed geometry returns no aggregate and publishes
  no draft, Change notification or preparation/publication side effect.

James confirmed actual FieldState/Snapshot.Fields and then real field edit/gate
methods. Dalton confirmed actual strict scalar schemas and source lexemes. The
normalized/public-path distinction was relayed to James and Noether. No competing
runtime/numeric/parser stubs or native chrome guesses were added.

## Scoped verification

Environment: Go 1.27.1, linux/amd64.
Commands run from `SDUI/go` unless explicitly indicated:

| Command | Result |
| --- | --- |
| `go test -mod=readonly ./layout -run TestScalar -count=1` | PASS after all new geometry/gate tests were added |
| `go test -mod=readonly -race ./layout ./svg ./prototype -count=1` | PASS: layout 1.341s, svg 1.055s, prototype 1.059s |
| `go test -mod=readonly -race ./layout -count=1` | Final PASS 1.597s after strengthening the unchanged 0.2 JSON/unsupported-scroll assertions |
| Root `git diff --check -- SDUI/go/layout SDUI/docs/go-layout-contract.md SDUI/docs/wci3-m1-layout-api.md` | PASS |
| Root `gofmt -l SDUI/go/layout` | Empty output |

Meaningful oracles cover all four controls, source/live-value separation, exact
label/feedback/button parts, inherited font, detached metric state, read-only
versus inherited disabled hits, nested clip/translation, independent native parts,
exact increment-button revelation, unchanged scroll routing, split minimum bounds,
undersized intrinsic probes, inactive page and closed nested-dialog exclusion,
anonymous/reused runtime identity, independent open canvases and failed child
sizes. Adversarial metrics include zero/negative/nonfinite/unbounded minima,
missing/partially empty/outside/overlapping regions, missing native adapter and
wrong/missing field projection. Invalid geometry leaves the supplied snapshot
unchanged and returns no partial candidate.

Real runtime gate tests prove rejected prospective geometry leaves the live typed
draft, accepted geometry, observer count and preparation/publication counts
unchanged. A successful edit publishes its matching geometry while retaining the
accepted value separately. Invalid raw numeric text stays visible with validation
in an open dialog, keeps its fixed measured parts and survives rejected child
resize. A parent resize preserves the accepted nonmodal field allocation.
Ordinary 0.2/0.3 input/button output remains byte-equal with no Fields JSON key;
0.2 scalar sources reject and its unsupported-scroll behavior is preserved.
Existing collection, pane, canvas, SVG and prototype regressions also passed.

## Remaining integration boundaries

This layout lane has no known unresolved implementation blocker. Host integration
must consume these exact measured parts with the same native objects/metrics.
Actual Fyne checkbox/slider/select/number minimums, visible validation feedback,
read-only/disabled gestures, number entry/button hit geometry and clipping,
focus/EnsureVisible in real windows, option popups, typed SDL calls and native
screenshots/event counts remain host/main acceptance. Synthetic metric and Go
runtime tests do not certify those native outcomes. No WCI3-M2 editing/IME or
whole-card completion is claimed. Concurrent runtime/host changes require main's
final exact-candidate integration and independent review.

## Exact owned candidate hashes

| Path | SHA-256 |
| --- | --- |
| SDUI/go/layout/fields.go | 889f8f3b78ab28515b40ad51ff73b8dd185bd81f34321349cced413f8e827d61 |
| SDUI/go/layout/fields_test.go | 48a2c9153b99571e52a39f582b83234d75fd973d91e1c29937a83dc8ca45cc48 |
| SDUI/go/layout/field_validation_test.go | 131cd1614a61f0cfe940103296e4d801a6a2027ffbc59b34eaf877e7000b4458 |
| SDUI/go/layout/field_gate_test.go | 0df6f2557126b1f03bcdaf18ef7f3db05764f7e7409d6a5084e4fdd1813353bb |
| SDUI/go/layout/canvases.go | bba40ccfb1301f9c2facc2de45fedc14719f89d9bab717271a035eeaff074ed1 |
| SDUI/go/layout/collection.go | 4baa4a4d74f0a22fb479877c1fc2d8b98bb7881518c4c02fa3ea0b688770b718 |
| SDUI/go/layout/engine.go | 8104b67f26cee049ff6bb6421e4dd0b11c61e552fab34713df77e3a1737b4d35 |
| SDUI/go/layout/pane_minimum.go | c0eda4bee8b71daaae98455bc1dbeba6e65f1a81be1c06feb138e73286d14b82 |
| SDUI/go/layout/types.go | 3e0d078d8547583ed2c87442df244b95439b53f616041760f54cf0d010a0c99c |
| SDUI/go/layout/viewport.go | 6b5d4a8b60410a0112bcabf52163e91f1ee59b4c9eaea5c31ff674e0b68a440d |
| SDUI/docs/go-layout-contract.md | 8790400cbdcfb903178b21d16499c7b84379bcd89bd3cdd3a5bf9e6ab6b60154 |
| SDUI/docs/wci3-m1-layout-api.md | a28633188d683f4fc8a331b2ecfc1df4e51613061f340b6ecc699dc723433846 |

## Coordinator Session handoff

Owner input selected WCI3-M1 at a28b3cc and assigned only layout/docs/report.
Worker recovered canonical authority and APIs, froze native field metrics with
Noether, integrated actual James/Dalton types, delivered geometry/minima/feedback
and scoped meaningful tests, and documented the normalized/public identity
correction. Result: scoped checks passing, layout ready for native integration.
Next: host uses the measured fields; main runs actual native/SDL acceptance and
independent final-candidate review. Main updates Session0010 S4 and plan/evidence;
this report is the Worker handoff, not a management transition or stage closure.
