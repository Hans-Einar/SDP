# WCI3-M1 runtime worker handoff

Date: 2026-10-08. Worker: James. Scope: runtime and the new stdlib-only numeric package; no host, parser, layout, bridge, management, Session, module or branch edits. No commits. Loaded routines: SDP 1.1.1, SDP Worker 2.0.0, shared document workflow and SDUI instructions. Coordinator retains integration, native verification, canonical records and independent gate ownership.

## Candidate and authority

Isolated clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci3`, base HEAD `a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625` plus the exact lane files below. This is an uncommitted shared candidate, not evidence for unchanged HEAD. Other lanes have concurrent, separately owned changes.

Authority: selected WCI3-M1 and original `SDP/04--Design/SDUI/Widgets/Values-and-text.md`, including reviewed numeric no-cap/endpoint-spacing and closed-dialog successor refinements. Existing WCI3 runtime/frontend memos supplied preparation; the implemented runtime README and exported declarations describe the actual seam. This work implements checkbox, slider, select and number only. No extended input, multiline, IME or later-stage implementation.

## Delivered behavior and downstream API

- `numeric.NewGrid`, `Parse`, `Tick`, `At`, `Text`, `Nearest`, `LastTick`, `ValidateSDL`, `ExactInteger` and `Integer`: original lexemes use bounded exact rational checks; typed Number round-trips to an exact reconstructed grid point. Gesture snapping is separate. The no-cap exception is restricted to exact safe53 integral min/max/step. All other grids enforce the 2^26 interval cap and strict two-direction endpoint spacing. Scanning bounds and exponent saturation precede arbitrary-precision construction. The immutable package imports only the standard library.
- Typed `Value` Number/OptionID and copied `FieldState`, `FieldTarget`, `FieldChange`, `ControlCommit`, `OptionTarget`, `ChoiceOption` APIs. Legacy input Widget.Value/Draft stays authoritative; new scalar storage is separate. `Snapshot.Fields` is detached and exact-instance-path keyed.
- `EditField`, `EditTick`, `ChooseOption`, `ObserveChanges`, `ValidateFieldWith`, `RevertField`, `CaptureCommit`; existing Dispatch/Handler accepts typed Commit with exact captures. Automatic Commit uses the original `change.Field.Target`, preventing observer reentrance from restamping an old gesture. Invalid numeric text stays visible. ReadOnly, built-in validation and checked programmatic state changes remain distinct.
- `BindChoices`, `ReplaceChoices`, `Option`, `ValidateOptionTarget` provide an exact complete inventory, explicit supplied-empty generation and stale-popup barriers. Replacement retains invalid accepted-ID diagnostics and clears invalid proposals without label/index fallback. No hidden collection/provider framework.
- Apply atomically stages mixed properties using the original batch revision baseline, independent of distinct-property order; validation feedback stages last. Typed acceptance requires captured value/draft revisions and option generation when applicable. Callback and validator reentrance cannot overwrite accepted work.
- `DialogControls` provides mixed Go Accept while DialogFields/DialogField stays text-only. Mixed writes preserve the 256-write budget, domain outcomes, replay blocks, exact surface lifetime and published-only receipts. Unbound automatic dialog Commit keeps the proposal pending; explicit child Commit survives Cancel.
- `SuccessorWithChoices` validates complete choices and retained accepted values before publication. Compatible main/inactive-page drafts survive, including invalid number text. Closed-dialog successors reset unaccepted proposals to current accepted values; earlier explicit child commits survive. Failed candidates preserve live state. Detached successors copy no closures. Compatibility Reload retains compatible handlers/validators/observers and guards reentrant validation.
- Existing M1 panes, M2 commands/surfaces, WCI1 collection lifecycle, accepted presentation tickets and 0.2/basic 0.3 input behavior remain covered by existing suites.

## Reviewer finding and correction

The independent P2 reproduction showed `[AcceptedValue, ReadOnly]` rejected the second property after candidate revisions advanced, while reversed order succeeded. Revision validation now uses `s.checkFieldUpdate` against the original session; candidate application only stages the checked values. `TestTypedBatchChecksOneBaselineForEveryPropertyOrder` asserts both accepted value and read-only state, exactly one counter increment and identical complete snapshots for both orders. The reviewer confirmed the scoped fix. Validator reentrant reproductions already passed; they remain covered rather than being reported as current defects.

## Reproducible checks

Environment: Go 1.27.1, linux/amd64. Commands run from SDUI/go unless stated otherwise.

- `go test -race ./runtime ./numeric ./layout ./preparation`: PASS; runtime 2.232s, numeric cached, layout 1.597s, preparation 1.087s on the pre-refinement frozen checkpoint.
- `go vet ./runtime ./numeric`: PASS.
- From SDL/go: `go test -race ./bridge ./examples/values/...`: PASS; bridge 3.587s, values example 92.685s.
- `go test -race -overlay=/tmp/wci3-runtime-review-l12g23h7/overlay.json ./runtime -run 'TestReviewer|TestTypedBatch' -count=1`: PASS 1.191s on the pre-refinement frozen checkpoint.
- `git diff --check -- SDUI/go/runtime SDUI/go/numeric`: PASS.
- Added 19 initial runtime tests covering typed storage, Change/Commit ordering, invalid raw drafts, gate failure, exact self acceptance, stale observer/validator captures, malformed values, property batches, option generations, inventory atomicity, mixed Accept/outcomes, child Commit/Cancel, reload/successor retention and the 256/257 budget boundary.
- Added six numeric tests covering adversarial lexical/exponent bounds, exact decimal membership versus rounded lookalikes, safe53/wide grids, endpoint representability, gesture-only clamping, exact display and concurrent immutable use.

Independent reviewer reported numeric scope PASS with an independent big.Rat oracle over 445 admitted grids (exponents -320 through 296), exact Text/Parse round-trips and hostile numeric inputs; overlay `/tmp/wci3-numeric-review-76twbmoy/overlay.json`. Runtime scoped reviewer reproductions and the Reload validator-reentrance delta also passed. Final independent runtime/numeric race with all four reviewer reproductions passed in 3.085s/1.054s; reviewer explicitly retained scoped runtime checkpoint approval after inspecting the final guards and choice-rebinding delta. These statements identify reported scoped review evidence, not whole-stage approval.

## Remaining integration ownership

No known runtime/numeric blocker at this checkpoint. Host/native behavior, actual SDL workflow proof and the final independent integration gate belong to the coordinator and their owners. The shared `SDUI/docs/runtime-contract.md` is outside this lane's current explicit write boundary; coordinator should reconcile it using the updated runtime and numeric READMEs. No M2 text work is started. Current numeric source and runtime API are frozen for downstream builds; any further fixes require a reported checkpoint delta.

## Approved required-empty select refinement

After the frozen checkpoint, main and independent reviewer approved a narrow
contract refinement: a compatible already-required select with accepted empty
absence can survive replacement as editable invalid state, including a closed
form. It does not become accepted/valid, gain a value revision, or select an option.
Newly-required blank and removed/disabled nonempty accepted values still reject.
This is a changed contract disposition, not a retroactive defect claim.

Main authorized only `field_successor.go`, related tests and this report for the
bounded M1 unfreeze. The validator now receives the predecessor Required fact;
retained empty values reject only when Required is newly introduced. All other
retained-value checks remain unchanged. No public API adjustment.

Three added regressions (22 runtime tests total across the two new test files):

- `TestRequiredEmptyChoiceSuccessorKeepsEditableAbsence`: supplied-empty and
  available-option sets, exact empty accepted/proposed state, unchanged accepted
  revision, invalid feedback/no Commit, detached predecessor preservation and
  compatibility Reload.
- `TestRequiredEmptyChoiceClosedSuccessorDiscardsOnlyProposal`: an open dialog has
  a valid but unaccepted choice; successor closes it, restores invalid empty
  baseline without accepting the proposal, leaves predecessor untouched, blocks
  reopened Accept/Commit, and remains editable for explicit repair.
- `TestRequiredEmptyExceptionDoesNotRelaxNewOrMissingChoices`: newly-required empty,
  removed accepted ID and disabled accepted ID all reject atomically.

`go test -race ./runtime ./numeric`: PASS, runtime 2.871s, numeric cached.
`git diff --check -- SDUI/go/runtime`: PASS. Independent reviewer approved this exact delta (`d2a79882` source / `753aa3b4`
tests), inspecting closed-form reset, newly-required and removed/disabled guards.
Full runtime/numeric race with four independent overlay regressions passed in
1.807s/1.024s. No confirmed runtime/numeric blocker remained on these bytes;
overall M1/native gate remains separate.
All other runtime/numeric source files remain frozen. The separate authorized
read-only `WCI3-M2-runtime-API.md` incorporates the approved refinement and newline
precedence; it is not M2 implementation.

## Exact owned file manifest

SHA-256 values include all changed/new owned source, tests and READMEs; this report is excluded from its own manifest.

| Path | SHA-256 |
| --- | --- |
| `SDUI/go/numeric/README.md` | `a4b33de5da1f39470d3df1958fa095f0ab438f8d0308e63b823bf1731e2dcc3a` |
| `SDUI/go/numeric/decimal.go` | `403a410451b84fa14e806cdacfa661ca3f152f03be9e3fbe764ab52f859c7f2d` |
| `SDUI/go/numeric/grid.go` | `b322ef66b985a3528a5fcbed89110940dc197d7355cb01b29586b531e9fc56d2` |
| `SDUI/go/numeric/grid_test.go` | `5c5f70c967a00f3f735245d17505309c45b0bbf6899b4a4e8aeef79e689281fa` |
| `SDUI/go/runtime/README.md` | `114ec56cf672458580922c30c08dd0d6aa48412a690c5eaa7d972f5098688696` |
| `SDUI/go/runtime/choices.go` | `6df86b0891c8170ccfd11be6313bfa35fd1b3ce5900550f683b6e20707cd1c61` |
| `SDUI/go/runtime/command_capture.go` | `3f9e3376c6e28bf5905c9b49b48960e25aef56ccfa32a2d3fcaac3214b9cab7a` |
| `SDUI/go/runtime/command_dispatch.go` | `aec860cb928a5bc219c104f69d85195da1157c25617af60923bce62eb2869eaf` |
| `SDUI/go/runtime/command_state.go` | `92ed544626997d1727e3d676723127617621c45547d15ce89704bc486336391c` |
| `SDUI/go/runtime/command_types.go` | `63bfa2a3673cddb1c8c3866b4f8b05c41c357d0a335d0809e849650186a446ee` |
| `SDUI/go/runtime/events.go` | `eb46696cdcbdce04f97cfb5df4a0c9689815b17c63cd2cd3f4866f77016c04f8` |
| `SDUI/go/runtime/field_commit.go` | `d68ccd9bfc05a26e14b1f6b88720a9d70c4763d875ad05c623f52d33dd492f63` |
| `SDUI/go/runtime/field_edit.go` | `d584d16554bff61f3fd4a10e615a94d320c3304884dabbb6f733b19dbd4bec18` |
| `SDUI/go/runtime/field_lifecycle_test.go` | `753aa3b47dacac7e818b3d65c54ebc655b6e81b7f6f406cc81f5d69149f5a87e` |
| `SDUI/go/runtime/field_successor.go` | `d2a79882efbcaa706d487bcfb5f811f73910d8fe12a00b5881365459add087dc` |
| `SDUI/go/runtime/field_types.go` | `88a91fa823e2272f17d5cb34bb2dff9aabb5ab34ab35b790a93545e92bb670c6` |
| `SDUI/go/runtime/field_update.go` | `637265c5c96c9f919f751487a8d9f3c4783a0e5f55830e421644e3316cb0420f` |
| `SDUI/go/runtime/fields.go` | `417b88501efcc3e19d07d7239fec4c64473ee7ba2834845c4e6f3f7762ace525` |
| `SDUI/go/runtime/fields_test.go` | `eea4b384e8b674744f9c3202cddfb790e5b15c07d05d36238264186c0701da65` |
| `SDUI/go/runtime/interactions.go` | `7b70e3d6bc6044f15dbdc0efa09d0a65865323c0abdd9a791e598c314fe16bf5` |
| `SDUI/go/runtime/properties.go` | `193d0f5366706a9b2d5a7815c1a943c5d324dee296c1d5093bf76cfb3b95f46d` |
| `SDUI/go/runtime/reload.go` | `ec670327fdb27800e97cc579a4649d035853ec8d0449e9023683cf033d7018eb` |
| `SDUI/go/runtime/session.go` | `f19c8014dbb31b270956dedb930d25ac7e5295a96b6d6b5f49e385369d82d1c9` |
| `SDUI/go/runtime/state.go` | `e2b74972e2632c9c91f300e470f9fa40ea63a988f22534de4a106363e6fe99ce` |
| `SDUI/go/runtime/successor.go` | `b980001006c9756d89d2380959c13d0365de8201378c74c689c16f659e067465` |
| `SDUI/go/runtime/surfaces.go` | `2e9f52d3c8eb3a878781a6cfb70e79d69e41c1910c80c33a581ad70091f706b2` |
| `SDUI/go/runtime/types.go` | `12aeb597dd7b3d7476b2bc18515d050cbd0db531a2c85f59fbd0fe1f5a59d363` |

Canonical Values-and-text observed at report time: SHA-256 `424c67b4a88827c7c38a41dfed49f94cfe983c3d1e5e4b49f8f3a4dd3d4dbbba`. This mutable document may include later coordinator refinements; it is not the candidate code identity.

Required-empty refinement canonical SHA-256: `461c70a63f5a3111c2ebff837b3a12c21e59cd8c28c3ede4645129321da2415f`.
