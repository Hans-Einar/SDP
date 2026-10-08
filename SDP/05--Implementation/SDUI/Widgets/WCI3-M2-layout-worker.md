# WCI3-M2 layout worker handoff

Status: component implementation complete and lane frozen for integration/review,
2026-10-08. No whole-M2/native acceptance or independent review is claimed.

## Authority and candidate

Owner selected WCI3-M2 after approved M1. Reused SDP 1.1.1, Worker 2.0.0 and shared
document workflow; recovered Session0010 S4 read-only. ORIGINAL governing contract
`SDP/04--Design/SDUI/Widgets/Values-and-text.md` SHA-256
`8199c275fecff3f6d03835540a3e586dd418ff70470f641590d2d75ace8f51be`, including
explicit-presence opt-in, exact CR/LF, exact-empty baseline and reviewed native
editing recovery refinements. Read the frontend/runtime/host/layout M2 preparation
memos and actual exported parser/runtime APIs before dependent implementation.

Clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, HEAD
`ea49991f2065679c93e39fe02a3d458fe72803df`. Tests ran on this HEAD plus concurrent
M2 lane changes, not on unchanged HEAD. Go 1.27.1 linux/amd64. No commits, branch
switch, modules, host/parser/runtime/bridge, management, Session or canonical
architecture edits by this lane. M1 manifests/evidence and prior reports retained.

## Delivered behavior

The existing FieldMeasurer API now handles extended input selected solely by
parser.InputOptions. Basic 0.2/0.3 inputs retain generic measurement and geometry.
The new internal predicate participates in desired size, recursive pane minimums
and final arrangement. Input policy is validated on inactive/closed source branches;
only active fields require snapshot data and metrics.

Active fields require normalized InstancePath/kind and InputState matching immutable
Multiline, Required and effective Placeholder. Mutable ReadOnly remains snapshot
state. James confirmed omitted-placeholder fallback updates with runtime Label;
a real runtime Apply test covers it. Input metadata and raw draft passed to adapters
are copied. No authoritative text, parsing, newline validation or native scroll
state was added to layout.

Empty input text may omit Label as the zero rectangle; nonempty text cannot.
Input labels may not overlap the Entry, active feedback is mandatory and separate,
and numeric button regions must be absent. Existing finite positive minima, source
max bounds, bounded intrinsic probes, final allocation checks and source diagnostics
remain applicable. Control covers the entire native Entry including its own chrome.
No exported layout type, measurement signature or PresentationState shape changed.

Shared clips/transforms/gutters/reveal and independent canvas sizing continue to
apply. Long text cannot contribute another shared viewport or text-content extent.
The native three-row word-wrap policy stays in Noether's actual adapter. Layout
cannot certify native hit routing or history. The docs explicitly retain the host's
reviewed exceptional recovery for an already-delegated rejected native edit while
requiring failed geometry probes themselves to remain pure.

## Verification

Working directory: `/tmp/sdp-sdui-widgets/SDUI/go`.

* `go test -mod=readonly ./layout -run 'TestInput|TestScalar' -count=1`
  passed during implementation; final expanded test set is included below.
* `go test -mod=readonly -race ./layout ./svg ./prototype -count=1`
  PASS: layout 1.399s, svg 1.065s, prototype 1.063s.
* `git diff --check -- SDUI/go/layout SDUI/docs/go-layout-contract.md`
  passed. All changed/new Go files formatted with gofmt; new documentation/report
  whitespace and linked local paths checked separately.

Ten new meaningful test cases (plus adversarial subtests) cover:

1. Explicit false/empty opt-in versus legacy geometry JSON, retained legacy draft,
   no promotion from incidental Fields metadata, and 0.2 argument rejection.
2. Empty/nonempty labels, required feedback, inherited font, readonly/disabled hits
   and copied Input/RawDraft metadata that cannot alias the snapshot.
3. Exact nested scroll translation/clips, both independently reachable ancestor
   gutters, outer EnsureVisible, bounded extents and no Entry viewport authority.
4. Relative padding/native split minimums, undersized intrinsic probes, reused and
   anonymous runtime identity, and rejected finite allocations without mutation.
5. Missing/wrong snapshot policy/identity, absent adapter, source-only rejection,
   post-source-max native minimum, direct input scroll rejection and hidden malformed
   source policy admission.
6. Malformed/overlapping/absent input parts, nonfinite/zero minima, input-only empty
   label exception and rejection of numeric subcontrols.
7. Actual runtime Label/ReadOnly Apply preserving effective placeholder projection
   and independent readonly eligibility.
8. A forced prospective field minimum failure before draft/observer/ticket change;
   valid and invalid-but-editable drafts subsequently publish once with fixed geometry.
9. Hidden page offsets/drafts, collapse/resize/failed restore purity and reveal clamp
   using real runtime PresentationGate and preparation tickets.
10. Closed/hidden/nested-surface omission, exact open nonmodal Entry coordinates,
    failed child resize preserving runtime/ticket and independent parent resize.

Synthetic measured rows are test oracles, not native font or word-wrap proof. The
existing M1 scalar suite and legacy layout suite ran in the race command; downstream
SVG/prototype packages passed with the shared candidate. No OS input, SDL action,
clipboard, history or IME proof was attempted by this worker.

## Exact changed paths and hash freeze

The following ten paths are this lane's product/document/test changes. This root
report is the sole additional handoff file. Dependency bytes from other lanes are
not claimed as owned or frozen by this worker.

| Path | SHA-256 |
| --- | --- |
| SDUI/go/layout/fields.go | 4a7370e8aa775917efa1954df9b27f3f02a56cb610270347f4bc94ac6d55c1f5 |
| SDUI/go/layout/collection.go | c9ae6e16c60997d90404956f63aa86abcf9ced0e211e63c9909d39d71ac1041c |
| SDUI/go/layout/engine.go | 5f5405af50fc821c155d9289e85a306fa65c369a465841168e6d4a363efdd8c6 |
| SDUI/go/layout/pane_minimum.go | 3c94125f1adc21bb830a328d3540e06b1da9e61ddae7e7649f4450d19fb475ab |
| SDUI/go/layout/viewport.go | c46f813beb94f2f00e40c855c7b28aa6655136a43ea2e760bdb70464b3f9eb49 |
| SDUI/go/layout/input_test.go | d010ba5474d18d6336dec55959092f4a0b353882d85fc732e10d5abf81e06545 |
| SDUI/go/layout/input_validation_test.go | 286da84c8dce713efb92a259a9671e67016c169b6017d7f5d1fac6799b30a68d |
| SDUI/go/layout/input_gate_test.go | 2f5941093448497802d28533d1936389c58afb9acad91ac30a2fe76f11cd6790 |
| SDUI/docs/go-layout-contract.md | 4bb6821efd9286628228263ae9ee6d8bd3cdfd72d16f7fd7e76317f03644ed24 |
| SDUI/docs/wci3-m2-layout-api.md | 4debc79503411f0a669a371265919955fcafa4c4d5b1584f74ea8393c4ec3d82 |

Inspected API seam hashes at handoff:

* `SDUI/go/parser/input.go`:
  `38c105f5a831877137e4ddc2a044b0323a8adc8a987ceab287359af54e600f83`
* `SDUI/go/runtime/field_types.go`:
  `5c571c969a27bcc815609a42e9667d0d0525b906c51e42a3d1c370e77e751027`

## Remaining work and coordinator handoff

No unresolved layout API gap or material contract departure identified. Noether
has the exact unchanged FieldMeasurer signature, empty-label/part rules and snapshot
policy checks. Dalton supplied the actual InputOptions helper; James supplied
InputState and confirmed mutable-label fallback semantics. No competing stubs.

Independent final review and integrated native proof remain required. Host/main
must establish actual three-row themed minima, long word-wrapped editing, text-area
wheel consumption without outer chaining, reachable gutters, label/feedback paint,
retention/recovery policy, clipboard, typed Commit and actual OS IME. A pure layout
failure proves neither native OnChanged rollback nor exact history preservation.
The declared recovery exception is host policy, not a layout relaxation.

Session0010 S4 work summary for main: owner selected M2 layout after M1; existing
field seam extended without new geometry/state types, scoped tests/race passed,
English layout contract/API docs updated, ten-path hashes frozen. Next step is
host integration and independent review/native acceptance. Main owns Session,
plan/card/traceability and exact integrated-candidate evidence. No management record
was edited here, and this report does not close M2 or authorize release/merge.
