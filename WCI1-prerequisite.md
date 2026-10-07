# WCI1 prerequisite: active generated-model provenance repair

## Authority and candidate

Authorized bounded repair for KB-SDUI-003, EVT-KB-SDUI-000021, following the
accepted WCI0 baseline failure investigation. Work performed only in
`/tmp/sdp-sdui-widgets`, on `sdui/widgets-wci0`.
Candidate: HEAD `fac09f263637ab233c8ee5890e29f809774c5146` plus the three
implementation/artifact files listed below. No commit or branch switch performed.
The supplied untracked design and implementation records were not changed.

Applied the previously loaded SDP Worker routine and document workflow; read root
instructions and checked subtree instruction discovery. SDUI instructions loaded
in this assignment chain remain applicable to its codegen regression. This report
is the only worker evidence update; Master owns management/Session reconciliation.

## Root cause and repair

KB005 added normalized `Instance.Declaration` and `Uses`. The active checked-in
UI constructor predated those fields, leaving root Declaration empty while the
frontend produces `page`. The reflection-based generator already supports all
exported fields, including UseSite slices. No generator or parser change was
needed.

Regenerated the active bundle through `sdl-gen`, without editing generated Go or
manifest bytes manually. `ui_gen.go` now includes Declaration and Uses fields;
`manifest.json` contains its regenerated output hash. `actions_gen.go` remains
byte-identical. Original source bytes, source revision and generator version remain
unchanged. Existing localized example text was preserved as source content; no
translation or source change was part of this bounded repair.

Added `SDUI/go/codegen/provenance_test.go`. The regression generates constructors
for two instances of an alias that reuses a leaf definition, compiles and executes
them in a temporary Go module against the local parser, and checks:

- Complete Document and normalized Root equality with the frontend.
- Root declaration and both entries of the reused definition chain.
- Original declaration, intermediate use and distinct outer use source spans.
- Separate instance paths and independent provenance slices across siblings and
  repeated constructor calls.

This tests executable constructor behavior, not merely presence of field names in
emitted text. Temporary generated test artifacts are managed by `testing.T.TempDir`.
No frozen fixtures or historical pinned evidence were modified.

## Commands and actual results

Go environment: `go1.27.1 linux/amd64`; cache `GOCACHE=/tmp/sdp-wci0-gocache`.

From `SDL/go`:

```sh
GOCACHE=/tmp/sdp-wci0-gocache go run ./cmd/sdl-gen -actions examples/edit-apt-cell.sdl -ui examples/edit-apt-cell.sdui -entry page -package generatedmodel -output examples/generatedmodel
GOCACHE=/tmp/sdp-wci0-gocache go test ./codegen ./bridge
```

Both commands PASS. SDL codegen includes TestGeneratedConstructors, exact bundle
comparison, and compiled-executable parity; bridge passes (cached, unchanged inputs).

From `SDUI/go`:

```sh
GOCACHE=/tmp/sdp-wci0-gocache go test ./codegen -count=1 -v
```

PASS: TestLiteralBoundaries and TestGeneratedReuseProvenance, including the child
Go compilation/execution. `git diff --check` also PASS.
No network or sandbox block occurred. No full SDUI race rerun was needed for this
artifact/test-only prerequisite; the earlier WCI0 race result is separate evidence.
No WCI1 widget implementation or native publication claim is made.

## Changed files

- `SDL/go/examples/generatedmodel/ui_gen.go` — regenerated active UI model.
- `SDL/go/examples/generatedmodel/manifest.json` — regenerated output hash.
- `SDUI/go/codegen/provenance_test.go` — executable reuse/provenance regression.
- `WCI1-prerequisite.md` — this report.

The reported baseline TestGeneratedConstructors failure is resolved. No remaining
repair issue was observed; review and integration remain Master's responsibility.

## Candidate file hashes (SHA-256)

```text
41a2319e75725f4e7758ce3a5697b125cecbd55bb157ced8524107cb2c7ed524  SDL/go/examples/generatedmodel/ui_gen.go
c5ae3847d50582fc6fd5b83321a2adf68d473bc29bdf8283dbf7b31f1947a3eb  SDL/go/examples/generatedmodel/manifest.json
ca5f894c4a386ceb78dc59be54bd3c4b76978b7379954a601805282538f38696  SDUI/go/codegen/provenance_test.go
```
