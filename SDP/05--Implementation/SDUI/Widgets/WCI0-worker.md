# WCI0-M1 worker evidence

## Candidate and scope

Isolated clone: `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci0`.
Baseline: `00b105ab4f7c4b7150750d69990564a04c30b035` plus the uncommitted
implementation files inventoried below. Baseline includes the supplied KB005
snapshot. Supplied untracked design/implementation records were read and left
unchanged. No original-workspace edits, management edits, commits or publication.

Loaded routine: `Skills/sdp-worker/SKILL.md` and its document workflow reference.
Read the supplied Design, Preparation, Acceptance and implementation plan contracts
and applicable repository/SDUI instructions. Independent implementation review is
still pending; this report is worker evidence, not approval.

## Implementation

- Pure `preparation` package normalizes the selected document internally, checks
  exact dimension/identifier/major capabilities including hidden nodes, preserves
  diagnostic instance path/source span, and validates measured layout before
  binding an unmounted detached runtime session.
- Concrete `fynehost/admission` facts separate native button/input support from
  Markdown and symbolic SVG placeholder provider support. Neither scroll nor
  real vector execution is advertised. Source profile remains 0.2; no grammar edits.
- The frontend identifier is `sdui/0.2` with admission capability major 1; source
  profile identity and capability contract version remain distinct.
- Prototype mode reports unbound declarations without invoking binding. Connected
  mode requires an actual adapter, verifies callback handlers were installed and
  selected connections are matched. Binding adapters are trusted synchronous
  code responsible for loaded module/signature validation without action execution.
- Failed candidates close their detached session. Successful candidates belong to
  the caller; `Admit` checks source identity and live model revision, and rejects
  closed candidates. The caller owns source I/O and the final revision recheck.
- Prototype readiness uses `markdown.Prepare(root, nil)` and measured provider
  geometry. Native runtime rejects unsupported normalized models before native
  construction and installs the same capability gate for subsequent model checks.
- SDL fixture uses the existing real `bridge.Bind` with loaded action-core Echo,
  input/button callbacks and visible input receiver. Missing module, mismatched
  typed input and invalid result reject without domain calls or live-session edits.
  Explicit later dispatch proves prepared and retained live handlers work.

## Verification commands and results

Environment: `go version go1.27.1 linux/amd64`.
All Go commands use `GOCACHE=/tmp/sdp-wci0-gocache`.

- From `SDUI/go`: `go test ./preparation ./prototype` — PASS, including actual
  failed geometry, hidden unsupported Markdown, exact capability mismatch,
  unsupported nodes/widgets, provider/host separation, adapter failure disposal,
  no-op adapter rejection, selected connections, source/model stale and closed
  admission, and detached callback/document correspondence.
- From `SDL/go`: `go test ./bridge` — PASS. Initial new fixture failed because
  its two explicit post-preparation dispatches reused SDL command sequence 1;
  corrected the second dispatch to sequence 2. This was a fixture error, not an
  action executed during preparation.
- From `SDUI/go`: `go test ./preparation ./prototype ./host/fynehost ./runtime`
  — PASS. `go test -v -timeout 90s ./host/fynehost -run
  TestAdmissionRejectsBeforeNativeConstruction` — PASS for unknown widget,
  hidden unknown node, scroll and closed session.
- From `SDUI/go`: `go test -race ./preparation ./prototype ./runtime` — PASS.
- From `SDUI/go`: final `go test -race ./...` — PASS on the final inventoried
  source, including Fyne, preparation, prototype, parser, layout, Markdown,
  runtime, reload, code generation and export packages. Optional external renderer
  acceptance remains environment-gated by the existing `SDUI_MMDR` test.
- From `SDL/go`: extra `go test ./runtime ./codegen` — runtime PASS; codegen
  FAIL in `TestGeneratedConstructors`, line 33: generated constructor differs
  from frontend. The failure compares action program, UI document and normalized
  root before application creation. `git diff --name-only 00b105a --
  SDL/go/codegen SDL/go/parser SDL/go/examples SDUI/go/parser SDUI/go/codegen`
  produced no output: all compared constructor/frontend inputs are unchanged
  baseline files. No generated-model repair was made within this assignment.
- `git diff --check` — PASS at final checkpoint.

An initial test-file creation command used repository-relative paths from the Go
module directory, so those writes failed without changing files. The files were
then written from the clone root. No failed write touched the original workspace.

## Limits and handoff

This implements detached preparation and legacy native admission only. It does
not deliver atomic native bundle publication, collection events, new widgets,
source 0.3, OS interaction/IME acceptance, consumer packaging or release evidence.
Full native publication remains WCI1. Fyne tests are headless adapter tests.
Nil-renderer Markdown retains the existing labelled diagram placeholder baseline;
no real external renderer capability is claimed. No network/sandbox failure or bypass occurred. The unchanged generated-model
comparison failure above remains for coordinator disposition.

The core accepts a validated document and trusted pure layout/binding adapters;
it does not sandbox arbitrary Go functions. Final source hash verification and
closing abandoned candidates belong to the composition caller.

## Changed implementation files

- `SDUI/go/preparation/capabilities.go`
- `SDUI/go/preparation/prepare.go`
- `SDUI/go/preparation/prepare_test.go`
- `SDUI/go/host/fynehost/admission/admission.go`
- `SDUI/go/host/fynehost/runtime.go`
- `SDUI/go/host/fynehost/runtime_test.go`
- `SDUI/go/prototype/check.go`
- `SDUI/go/prototype/check_test.go`
- `SDUI/go/runtime/session.go`
- `SDL/go/bridge/preparation_test.go`
- `WCI0-worker.md` (this report)

## Tested source inventory (SHA-256)

```text
2e5628efc4227021c8fbcf43e92378be018a8417aad785cf660e77fb6e4379df  SDUI/go/preparation/capabilities.go
c66360269397ac2be14798726f66617633157e9e687823f807bbba22a10a9415  SDUI/go/preparation/prepare.go
51c20982dae64cc91dfb8bdaf9fb7508f0c7d96a50fe966d22f2c4da3a8ea9ae  SDUI/go/preparation/prepare_test.go
545daf06a476882df7aa9e5f46fe635e0a2e49a29bd1af066fee81ed13278997  SDUI/go/host/fynehost/admission/admission.go
0dcf6db2f750d1eec2288edba2c1bdef19c406eadf66184dbbef6ee19f5cc3bc  SDUI/go/host/fynehost/runtime.go
04bd46453b60e37a57a6eae6c137a66402aba8ed2309c76d8454ccc312ad99b8  SDUI/go/host/fynehost/runtime_test.go
8d362a2e69cda1bf6225aad2c069b0582af215faad72adfb1efe6866259aa919  SDUI/go/prototype/check.go
3428a2bda630194bfbdf32aa53110667012c4dca5b6345e250f635fcc6d3f8c1  SDUI/go/prototype/check_test.go
905b2b9e9bace48cded2ff36d79dfe6dfac646703279d6347d1e651506b5c76a  SDUI/go/runtime/session.go
4bea3be4ead0af041fab7224ce82bd58d6a0be27a7731cc2e39ed420cbf6dd78  SDL/go/bridge/preparation_test.go
```
