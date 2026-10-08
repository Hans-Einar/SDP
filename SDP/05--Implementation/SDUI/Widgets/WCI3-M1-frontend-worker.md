# WCI3-M1 frontend/export worker

**Status:** bounded lane implemented and scoped tests passed; ready for independent integrated review. Native/SDL whole-slice acceptance remains with main. Frontend edits freeze after this report except coordinated verified fixes.

Authority: selected WCI3-M1, canonical ORIGINAL Values-and-text reviewed refinement (SHA-256 719cb7eab8375f20bbedb994fc0dc2eb40928b43903d5270238ce07629220309). SDP Worker and document workflow reused. Clone branch sdui/widgets-wci3, phase HEAD a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625 plus concurrent disjoint worker changes. No commits, branch switches, module edits, management/Session records, original files or other-lane files modified by this worker.

## Concrete delivery

- Four closed .3 scalar leaf schemas: checkbox, slider, select and number. Nonblank labels, argument order/types/required numeric fields, named callback owners and body rejection. A named definition can supply a reused scalar's callback owner name. No inline choices, source change handler, validator expression or future input fields.
- Exact numeric-token capture only for .3 slider/number named min/max/step/value. Existing numeric lexer and all .2/layout/split literals unchanged. Original span/token spelling survives Node/Instance normalization, reuse and generated constructors as number-lexeme Literal string. Quoted numeric strings and constructed rounded/malformed literals reject.
- Actual exported parser.NumericArguments(n *Instance)(min,max,step,value string,err error), strict representation extraction only. Scalar source and selected-root validation call James's shared numeric.NewGrid/Parse for exact grid/initial admission, retaining original error spans. No parser/runtime import cycle, private duplicate arithmetic or competing numeric stub.
- ResolveInteractions extends nearest-dialog ownership to scalar controls and permits them as widget-context targets, without command identities or $scope on fields. ResolveDialogField remains input-only for WCI2 text mapping; mixed typed capture stays runtime-owned. Existing chosen-entry/reuse reference semantics remain.
- Structural text/composition includes all scalar declarations, including hidden ones; source initial values/constraints preserve numeric spelling and do not claim live state. Select explicitly reports absent choice options. Callback descriptions remain symbolic/unverified; no parser/provider/SDL I/O.
- Public SVG emits source-linked unsupported-value-export for all four controls, even hidden or supplied geometry. Native SkipControls supports exact prepared scalar kind/path entries and preserves full-entry InteractionRoot pointer checks. Unknown/missing/wrong/extra inventory still rejects. CLI failed export preserves the old artifact.
- Standalone prototype reports unsupported-value for unimplemented scalar adapters and unsupported-provider choice-options for select; it never fabricates provider/native readiness from source acceptance. Host/preparation's actual M1 capability inventory is Noether's lane, consuming these unchanged normalized node/argument shapes.
- Updated profile and Go-generation docs, grammar annotations and actual frontend API memo. Input multiline/readOnly/placeholder/required source additions remain rejected until WCI3-M2 selection.

## Compatibility and integration seams

No Node/Instance fields, public existing parser signatures, profile/envelope tags or generator versions changed. NumericArguments is additive. .2 Profile sentinel remains empty/omitted; frozen default .2 source AST, text and generated output checks pass. There is no new SVG option or change to legacy SkipControls policy in this slice. Existing M2 basic-button dispatch distinction, source identities and constructor tests remain passing.

New scalar defaults stay absent in normalized Arguments: runtime applies closed defaults. Native live accepted/proposed/raw/choice state belongs in Snapshot.Fields, never by replacing number-lexeme with binary64 or stringifying typed values into legacy input fields. Peers received exact NumericArguments, Dialog ownership, native inventory and this projection rule before dependent coding. James owns numeric implementation/tests; frontend consumes its actual API.

No material contract departure remains in this lane. Outside-safe53 Go grid bounds and closed-dialog reload policy were already selected in the canonical refinement, not silently decided here. Option existence/generation, typed SDL receiver/result readiness, Change/Commit/Accept execution, geometry and native events require the other lanes and coordinator evidence.

## Verification

Environment: Go 1.27.1 linux/amd64. Dirty-candidate evidence, identified by this lane's hashes below; not an unchanged-commit or full-stage pass claim.

Final PASS commands:

- SDUI/go: `go test -race ./parser ./codegen ./presentation ./svg ./prototype ./cmd/sdui`.
- SDL/go: `go test ./codegen` (affected .3 metadata/provenance consumer and .2 preservation).
- `git diff --check`.

Earlier combined runs encountered incomplete runtime fields/SuccessorWithChoices exports while that lane was implementing. No stubs or cross-lane edits were made; final combined race run passed after actual exports became available.

WCI3-A01: all scalar schemas, .2 negative cases, M2 input exclusion, raw numeric Kind/string/span, signed zero/exponent spelling, reuse separation and actual compiled generated constructors. Existing exact legacy generated bytes, source AST parity, concept text snapshots and M1/M2 constructor regressions passed.

WCI3-A03 frontend evidence: decimal 0.3/0.1, off-grid decimal rounding onto a binary64 point, safe53 full interval without fractional cap, large-origin increments, outside-safe53 spacing equality rejection, oversized mantissa and exact offending spans; malformed constructed literals fail. Broader arithmetic adversarial/concurrency coverage belongs to numeric owner.

WCI3-A04/A06 frontend portion: stable-ID source preserved without loading options, explicit missing choice provider, mixed dialog ownership while text-only mapping rejects numeric fields. This is not native selection or SDL execution evidence.

WCI3-A09 frontend portion: hidden scalar and supplied-geometry SVG rejection, exact native inventory, full-entry dialog source context, no numeric-provenance bypass, truthful text/composition and atomic CLI output failure.

Remaining gates: whole candidate integration/race, real typed SDL/native gesture/validation/transaction proof, independent review and main's canonical evidence. WCI3-M2 extended text/IME is excluded, not partially claimed.

## Changed paths and SHA-256

All lane files at handoff; report excluded from its own hash. Other worker files are not included.

| File | SHA-256 |
| --- | --- |
| SDUI/docs/go-generation.md | 4fd86792c205a959d91c23cf2454edb0feff4978b1c79f1ff1edf77337cd16b9 |
| SDUI/docs/profile-0.3.md | 1a3dcb6a777f583b5371275e6f34472f79f117e2c2d301f50cc9c2bb4e91742c |
| SDUI/go/cmd/sdui/scalars_test.go | e24dab2038ea3988e35fa88cce63d989378d10bca4a1e91dbe3b0d1815f4c210 |
| SDUI/go/codegen/scalars_test.go | dfc4b363833127119bb009744e1f1e3bb2071d84f04b52f498cf17ab7c46c675 |
| SDUI/go/parser/interaction_identity.go | d4cfcaa05fb586dd2f880943ff9b662ba824951e01e63f42242b0eb36904fccd |
| SDUI/go/parser/parser.go | b29258d14de454d2a978d9d6b20bb78d9ef3820850fb5e84c5224c6320c1872e |
| SDUI/go/parser/scalars.go | 87ea2cdd1725be9f39ef50a98e475aeb26c552a4c419bd65d308e25a952ca342 |
| SDUI/go/parser/scalars_test.go | 880fad100c3d5c3bf94f81df5e9fff34bb84a4b4006d95dcf4a9e4dbf27951a5 |
| SDUI/go/parser/validate.go | d708e98b37b889e92d6fee2a552a7408a02ada6118c89c0af7ba7b2669cd2191 |
| SDUI/go/presentation/check.go | f0b903e4f89a4f7ab443fb6394e058f7cfb27e25e89212e8246d98e08b92825b |
| SDUI/go/presentation/dump.go | 4575b16a6e2ea263301146f7a93a29caa1565407368f3e095cbaa27d89cb0d05 |
| SDUI/go/presentation/markdown.go | 69df8765bf785ea4a478c74e370c6303c9908d82d3b0c2896769d90d8b3c4940 |
| SDUI/go/presentation/scalars.go | 7def441a67ed6625d20958bc82a2f054e5a32a9c284d903fdda0c40314736611 |
| SDUI/go/presentation/scalars_test.go | 84777ae59fa0edd06013e397d5a364390458d23dc4c103ef933a7b793839489c |
| SDUI/go/prototype/check.go | 9ea8e824bdccc5acd7a5783cc88e771b6305f849eab0be36850b5c359f612fdd |
| SDUI/go/prototype/scalars_test.go | e28df113499e02ba2cf2c8ff5c3b4430dbb8133f3052826df023105bf648d6fc |
| SDUI/go/svg/check.go | 076ea4fdfc544dc8be3489e58909a38e8186553a964264a59358f61b94e4ec65 |
| SDUI/go/svg/scalars_test.go | 5aff46abedd37ddf5c3c96cf980356c2df6fb760d7673dd85d81b46043dcb6bc |
| SDUI/grammar/sdui-0.3.ebnf | e52c085e3efceafb7bfa69357e5a04245f97c333a089f9db9e94e5ace40355b0 |
| WCI3-frontend-API.md | f501493fdff5242afbd47dae8e9cb56e14a06c300033501c8e32a3b1eda4dac6 |

Shared numeric dependency inspected at handoff (owned by James, not edited here):

- SDUI/go/numeric/decimal.go: `403a410451b84fa14e806cdacfa661ca3f152f03be9e3fbe764ab52f859c7f2d`.
- SDUI/go/numeric/grid.go: `b322ef66b985a3528a5fcbed89110940dc197d7355cb01b29586b531e9fc56d2`.
