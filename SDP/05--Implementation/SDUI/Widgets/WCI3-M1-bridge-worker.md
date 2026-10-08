# WCI3-M1 SDL bridge / values fixture Worker handoff

## Scope and candidate

Selected WCI3-M1 under KB-SDUI-003 / PLAN-SDP-0022; role SDP Worker, sdp1.1.1,
sdp-worker2.0.0 and shared document workflow reused. Canonical original
Values-and-text / Acceptance / phase Plan and both WCI3 API memos read. Current
branch sdui/widgets-wci3, phase HEAD `a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625` plus concurrent uncommitted lane
changes. Go1.27.1 linux/amd64; pinned Fyne2.8.1. No commit/branch/module,
management/Session, SDL parser/runtime or other lane file edits by this Worker.
Main owns integration, canonical evidence, native execution and Session0010.

This is a concrete evolving-candidate handoff. Own production paths are
SDL/go/bridge and NEW SDL/go/examples/values; authorized root memos/report only.
Existing commands/panes/collections fixtures are unchanged. WCI3-M2 extended
input/text/clipboard/IME behavior is not implemented or claimed.

## Delivered boundary

- Closed ScalarResult with explicit self setHandle/typed OutputField, forbidding
  dialog decision fields, retaining optional result revision/context checks.
  One result update must accept the source proposal; a different receiver cannot
  fulfill that contract and fails preflight. Legacy TextResult remains unchanged.
- ControlBoolean, ControlNumber and ChoiceOptionID consume validated typed Commit
  captures. ControlText consumes existing basic input Commit only; the separately
  discussed future extended-input Control path remains M2. No new legacy coercion.
- Numeric preflight uses frontend NumericArguments plus shared numeric Grid,
  ValidateSDL/ExactInteger/Integer. Source min/max/step/value preserve lexemes and
  diagnostics; raw number draft is checked exactly before transport. Returned
  int64 bounds are checked before float64 conversion, then grid and exact proposal.
  No second numeric parser, tolerance, text-number escape or SDL type expansion.
- Capture/result validation checks owner/model/value/draft and select option
  generation, plus post-consumption state. Returned update includes exact value,
  draft and option guards and AcceptDraft. Errors/stale delivery cannot replay;
  typed legacy Handler does not invent a domain outcome protocol.
- SDL DialogAcceptResult rejects mixed typed owned controls, remains text-only.
  Mixed persistence explicitly binds Go InteractionHandler and retains existing
  WCI2 capture/accept/outcome/Cancel/Close semantics. No new form transaction layer.
- Actual action-core fixture: four scalar echoes, explicit child Count echo,
  basic text Load/Save, text-only SDL Accept and mixed Go Accept. Options use stable
  IDs/repeated labels, whole-set replacement, no hidden collection or loader.
  Failure hooks execute only inside actual actions/Change; stdin only arranges
  trusted conditions/barriers/lifecycle. Observe reentry invalidates the original
  automatic gesture capture without a simulated Commit.
- Native CLI uses actual host geometry, Fields/choices inspector, dialog receipts,
  guarded parent hooks and per-command JSON barriers. closed.pending reads actual
  runtime collection requests after Session.Close. Source paths/JSON/protocol are
  documented in WCI3-M1-fixture-API.md. Main requested x=fill on12 scalar/text fields
  after pilot1 showed tiny intrinsic controls; paths/domain semantics stay fixed.

## Verification and limits

Run from SDL/go unless stated otherwise:

| Command | Evidence |
| --- | --- |
| `go test ./bridge` | PASS0.152s after mixed SDL rejection regression |
| `go test ./examples/values` | PASS7.073s before width-only fixture refinement |
| `go test -race ./bridge ./examples/values` | PASS bridge2.101s / values96.932s before width-only refinement; new mixed preflight test followed with bridge plain test |
| `go test ./examples/values -run '^TestValuesActualRequestAdmissionAndTrustedMute$' -count=1` | PASS2.717s after x=fill refinement, both variants |
| `go build -tags desktop -o /tmp/wci3-values-native-pilot1 ./examples/values/cmd/native` | PASS, immutable earlier pilot artifact |
| `/tmp/wci3-values-native-pilot1 --help` | PASS, --nonmodal documented |
| `git diff --check -- SDL/go/bridge SDL/go/examples/values` (clone root) | PASS before documentation handoff |

Bridge tests prove actual typed SDL echo and no replay, both safe53 endpoints,
near-integer/off-grid/raw-invalid rejection with no action, fractional Go-grid
SDL rejection, wrong mode/owner/selector/output/self receiver, malformed/unsafe/
different/stale draft/options outputs, and mixed SDL Accept rejection. Existing
bridge M1/M2 and legacy tests remain in that package.

Fixture tests prove connected modal/nonmodal request admission, exact source paths,
programmatic mute/readOnly/strict condition commands, actual observer reentry,
error/malformed/newer-draft/options conflicts with honest domain maps, proposal
Cancel vs explicit child Commit, mixed Go Accept and text-only SDL post-domain
conflict/Cancel, and failed/successful reload retaining invalid drafts.

Main owns all native execution on :190. Main reports pilot1 Boolean8PASS; this is
attributed evolving-pilot evidence, not final acceptance. The pilot1 immutable
binary is `/tmp/wci3-values-native-pilot1`, SHA256
`78344a92a90036df0da2248ddbb07de58e07997021a4ebad4396a0bec919d4ac`. It predates the
x=fill source refinement; never treat it as proof of current geometry. No binary
overwrites or native processes by this Worker. Pilot2 awaits coordinated host/
runtime fixes plus current fixture source. Final full-candidate native proof,
independent review and canonical evidence remain Main responsibilities.

Initial connected-admission failure exposed preparation's input-only receiver
allowlist; host/preparation owner corrected it without bridge/parser weakening.
No competing runtime/frontend stubs were introduced. ReadOnly update and exact
source self acceptance were explicitly coordinated with runtime; actual APIs are
used. Future M2 ControlText capture refinement is acknowledged without code.

## Latest scoped regression / lane freeze

`go test -race ./bridge ./runtime ./codegen ./examples/commands ./examples/panes ./examples/collections`
from SDL/go passed: bridge3.592s, runtime cached, codegen5.822s, commands47.535s,
panes67.129s, collections37.258s. This includes the added mixed-SDL rejection
regression and unchanged existing fixture source files. Final owned-path diff
whitespace check passed. Own source is frozen pending concrete review/native
findings; pilot artifacts are immutable. Width-refined values request admission
passed as above; prior full values race remains identified by its earlier source.
Independent reviewer was notified; no independent approval is claimed here.

## Exact owned file SHA256 inventory

- `SDL/go/bridge/bind.go`: `aa91e2dd55cfab4ad9fb120f2fc7c774f2ceb29eb41ecff694aca9d9c7752362`
- `SDL/go/bridge/interaction.go`: `1695292b99951770d4507ce4a6ee5b15eaff97f4c1ea173325f1371bf550a6bb`
- `SDL/go/bridge/values.go`: `3de306696660d107d66fabb5b3b43319df4c43cfabb8c7a22f6509b26bf4eef5`
- `SDL/go/bridge/results.go`: `4d9559305f272f6887eedc4ef3ee82f0ddcd9a08606b0638f94ba5964ce1fd24`
- `SDL/go/bridge/scalars.go`: `ac890c9376e3acd34ec61164df0a70b00f1e189b19683c9b32e3ae4176481d80`
- `SDL/go/bridge/scalars_test.go`: `4fdc0b2484c2c73b18b652046c4bb69571c747913eb849c01fc81d58b7108bc3`
- `SDL/go/bridge/README.md`: `00246c51e9230f4f12204bf9b88f983e115af4a56745f69e4bc49a9deecf80c7`
- `SDL/go/examples/values/choices.go`: `32fc49db97b05d642a9db4e2b51286ba16a495b2539012a423712bc8d09a4966`
- `SDL/go/examples/values/cmd/native/main.go`: `fc99208b3da5791227cc51f817e76c7563beccb69135c3ab1379b18a631b99b3`
- `SDL/go/examples/values/controls.go`: `addaa856cca08bc925a39aa7e3ebe47f31ab35b7e38b1991d1e1fa45eeaeedbd`
- `SDL/go/examples/values/fixture.go`: `819e8f91ecefaa90794cdd01f58614e2d322fb15b14d60adda3f7f0baba37326`
- `SDL/go/examples/values/fixture_test.go`: `cc9b9805c8844ffdb834925b80691b5fe424bcac06a327de786797fea8b9dfc5`
- `SDL/go/examples/values/source.go`: `5a36e51129fb5b49d8e6bf02e4c94a153dc1c5da270eeba962bffcb9547660de`
- `SDL/go/examples/values/README.md`: `33074da0a4ae8372c2929f0e606c13776d9ceb873aa6b0fe4b457f2ee7a0775f`
- `WCI3-M1-fixture-API.md`: `a2e57d8991b5e51365551e5726004922497238bcbd5fb3cf2301c736201339d5`

## Independent bridge-core checkpoint

Reviewer 01a11854-523d-7c11-adf4-f622b19b68c6 reports the examined bridge-core seam
passes: exact self receiver, retained lexemes/safe53 before conversion, captured
options, pre/post Execute state/receiver guards and text-only SDL Accept.
Independent full bridge race PASS1.933s; reviewed scalars.go SHA256
ac890c9376e3acd34ec61164df0a70b00f1e189b19683c9b32e3ae4176481d80.
Reviewer reports 14 owned product-path hashes matched. This is attributed review
evidence for the bridge core, not independent approval of the full values fixture,
native candidate or whole WCI3-M1 stage. No source/binary changed for this entry.

## Approved required-empty fixture variant / bounded refreeze

Coordinator authorized only source.go, cmd/native/main.go, fixture tests, README
and memo updates for --required-empty. It changes only main Mode's source to
required=true,value="", composing with --nonmodal. FormMode remains optional alpha;
paths, actual option sets, default Source, domain/state protocol and bridge core
are unchanged. Runtime owner supplies the approved compatible already-required
absence rule; this fixture does not implement or bypass successor policy.

Actual Request/Host.Adopt tests in both modal variants show invalid editable empty
absence with zero callbacks, successful unchanged reload preserving empty/invalid
and no automatic choice, then one explicit runtime/actual SDL beta repair. The
main native pilot supplies real pointer/key input separately. Default admission
and trusted-mute test still passes. From SDL/go:
- `go test ./examples/values -run '^(TestRequiredEmptyVariantAdmissionReloadAndExplicitSDLRepair|TestValuesActualRequestAdmissionAndTrustedMute)$' -count=1`: PASS1.835s.
- `go test -race ./examples/values -run '^TestRequiredEmptyVariantAdmissionReloadAndExplicitSDLRepair$' -count=1`: PASS11.775s.
- `go build -tags desktop -o /tmp/wci3-values-required-empty-check ./examples/values/cmd/native`: PASS.
- `/tmp/wci3-values-required-empty-check --help`: PASS, both flags documented.
- Owned fixture diff whitespace check: PASS.

Separate compile-check binary SHA256
6034b1115793b4888db7766c18d29fce97722a51f88ca5d855a013faedee1df8;
not launched or substituted for any immutable native candidate. Source SHA256
5a36e51129fb5b49d8e6bf02e4c94a153dc1c5da270eeba962bffcb9547660de.
The five authorized file hashes above are refreshed; all other owned bytes remain
frozen. Main/Noether own final candidate build after their runtime/host fixes and
native required-empty evidence. No whole-stage approval is inferred.

Read-only artifact cross-check: Noether's /tmp/wci3-m1-noether-slider-fixed has the
same full6034b1115793b4888db7766c18d29fce97722a51f88ca5d855a013faedee1df8
SHA256 as the required-empty compile check. Its --help exposes both required-empty
and nonmodal. This confirms fixture inclusion, not native behavior. Both owners
were notified; no product source/binary was changed.
