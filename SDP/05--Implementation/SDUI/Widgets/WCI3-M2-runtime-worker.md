# WCI3-M2 runtime worker handoff

Date: 2026-10-08. Worker: James. Candidate: clone `/tmp/sdp-sdui-widgets`, branch
`sdui/widgets-wci3`, base HEAD `ea49991f2065679c93e39fe02a3d458fe72803df` plus the
owned file manifest below. This is uncommitted M2 work, not a test of unchanged M1
HEAD. Prior M1 evidence/manifests were not edited. No commits, branch changes,
module/dependency edits, numeric changes, management or other product-lane writes.

Authority: selected WCI3-M2, ORIGINAL Values-and-text including explicit-new-argument
opt-in, exact-empty/CRLF successor rules and reviewed native editing refinement.
Reused SDP 1.1.1, Worker 2.0.0, shared document workflow and SDUI instructions.
Main owns Session/PM/traceability, canonical shared contracts, dependencies and OS
native proof. Frontend, layout, host and bridge retain their disjoint scopes.

## Delivered runtime behavior

- Actual exported `InputState{Multiline bool; Placeholder string}` plus
  `FieldState.Input *InputState`, nil for legacy/non-input. Frontend's actual
  `parser.InputOptions` is the sole source-policy resolver. Explicit false/empty
  new arguments opt in; profile/handlers/event shapes do not. All snapshots,
  validators and observers receive detached metadata/text projections.
- Widget.Value/Draft and their existing revisions are the ONLY accepted/draft text
  authority. Private input field entries store policy/feedback, never a second
  accepted/proposed/raw text value. Omitted placeholder follows current label;
  explicitly supplied placeholder, even empty, remains fixed.
- Extended EditField and Draft share one checked edit/publication/Change path.
  UTF-8/32768-byte admission rejects atomically; required/CRLF invalid drafts remain
  visible with validation and cannot Commit. Initial text admission returns the
  exact original value span/path. Required uses non-whitespace without trimming;
  only CR/LF counts as single-line breaks. Legacy Draft/event behavior remains.
- Existing typed CaptureCommit/Dispatch now operate on opted-in input metadata:
  String carries exact text, Control contains ValueRevision only (no RawDraft or
  Option). Nil-Control cannot bypass extended guards. Required/readOnly/activity,
  model/state/value/draft, observer/validator/handler reentrance and exact source
  acceptance are guarded. Missing/non-echo acceptance never claims success; failed
  callback delivery does not replay domain execution.
- Checked Apply writes Widget text directly, requires original batch value/draft
  revisions, validates the whole candidate and remains silent. ReadOnly permits
  checked programmatic Load/updates and focus, while edits/Commit/user Revert reject.
  Revert restores accepted text silently with validation; lifecycle reset handles
  hidden/readOnly controls internally. The M1 order-independent batch guard remains.
- Text-only DialogFields/DialogField remain unchanged. DialogControls and extended
  DraftField metadata support mixed Go Accept and exact text-only SDL captures,
  preserving 256 combined writes, Domain/AcceptBlocked, child Commit persistence,
  reentrant conflict barriers and published-only receipts.
- Successor retains compatible main/page drafts, including inactive invalid drafts.
  Retained accepted values violating new constraints reject. Already-required exact
  empty absence remains invalid/editable with unchanged accepted revision; newly
  required blank and whitespace-only invalid accepted text do not gain an exception.
  Multiline-to-single-line conversion rejects CR/LF in accepted text or a surviving
  main/page draft; closed-dialog draft reset happens before that check. An unchanged
  single-line invalid draft can remain visibly invalid. Detached successor copies
  no closures; compatibility Reload retains eligible observers/validators.
- Existing presentation probe/accepted-ticket publication is reused. No native
  pointers, text-scroll authority, acceptance-origin marker or history store added.
  The reviewed same-byte/history and rejected-native-edit rollback behavior remains
  host-owned; runtime state stays atomic on rejected edits and failed Commit gates.

## Coordination and documentation

Actual InputState export was sent before downstream code. Dalton's InputOptions
implementation is consumed directly; Gibbs agreed to copied metadata and effective
placeholder checks; Noether owns latest-authoritative-draft native rollback; Lorentz
uses typed ControlText, explicit self TextResult and exact post-Execute echo.
No competing stubs. Runtime README describes actual APIs and the native boundary.
Main owns the shared runtime-contract/architecture/requirements updates.

## Reproducible evidence

Environment: Go 1.27.1, linux/amd64; commands below run in SDUI/go unless stated.

- `go test -race ./runtime ./numeric ./layout ./preparation`: PASS on the manifest
  checkpoint: runtime 1.926s, numeric cached, layout 1.532s, preparation 1.081s.
- `go vet ./runtime`: PASS.
- Clone root `git diff --check -- SDUI/go/runtime`: PASS.
- Fourteen new runtime tests across input_fields_test.go, input_lifecycle_test.go
  and input_publication_test.go: opt-in/sole storage/detached projection/placeholder;
  exact Unicode/CRLF Change and Commit/legacy envelopes; validation/readOnly/atomic
  Apply; observer/validator reentrance; exact source acceptance and replay rejection;
  initial source diagnostics; inactive-page draft/form validation; required-empty,
  tightened required and newline successor cases; closed-dialog discard/child Commit;
  mixed Accept detached payload/receipt; batch baseline/silent Revert; policy changes
  and closure retention; failed Commit accepted-ticket preservation; post-domain
  validator conflict preserving newer draft and AcceptBlocked.

Independent reviewer scoped-approved this exact implementation after inspecting
all implementation paths and surrounding projection/dispatch/Accept/close/successor
callers. Independent full runtime race plus two additional overlay regressions
passed in 1.866s. Overlay: `/tmp/wci3-m2-runtime-review-4orde46x/overlay.json`.
The additional tests prove Label+AcceptedValue property-order consistency against
one Widget text store, and 256 mixed legacy/extended/scalar captured fields plus an
extra handler write rejecting atomically with DomainSucceeded/AcceptBlocked.
No confirmed runtime blocker was found. This is scoped runtime approval, not an
overall M2/native/bridge gate. Existing M1/scalar/numeric tests also pass, but component
race success does not establish actual clipboard, caret, undo, scroll, IME, native
rejection rollback or whole-stage acceptance. Those remain coordinator/host proof.

## Status and remaining work

Runtime lane complete and source/API frozen at the reviewed manifest; no known
runtime blocker. Actual integrated/native evidence remains with main and the other
owners. Any necessary correction will be reported as a new checkpoint delta. No future stage, release or merge is implied.

## Owned changed paths and SHA-256

| Path | SHA-256 |
| --- | --- |
| `SDUI/go/runtime/README.md` | `50565e384ed091c0685280b9ad2b8b843db82a92a4aa52acdbb6b35608b935c1` |
| `SDUI/go/runtime/command_capture.go` | `946ca8d643b46ea1c278564fc13acdcfcb0b36e49c8781b70aad0f65f18a3daa` |
| `SDUI/go/runtime/field_edit.go` | `ff41a0fc6a9cf2afb3f9c8943602d1ee913277837bcb91bc2fc0d5efda8c5043` |
| `SDUI/go/runtime/field_successor.go` | `8e1a427f47625b948a0f53f52c18a33bc8a43db6b9b60cbb2d42515be13dcd02` |
| `SDUI/go/runtime/field_types.go` | `5c571c969a27bcc815609a42e9667d0d0525b906c51e42a3d1c370e77e751027` |
| `SDUI/go/runtime/field_update.go` | `0ef434a43576300e6d4cf85e15506e6ded77eaa707663342418e0489dec906ce` |
| `SDUI/go/runtime/fields.go` | `c07eece21bb63bf9f9ecef052d93a72bcfcdee78fd8621b866d53ff14d3550a0` |
| `SDUI/go/runtime/input_fields.go` | `a364cea7333dc0264c4ec5050c30d9eb7336992e817c74bb8c403c04b3b8980f` |
| `SDUI/go/runtime/input_fields_test.go` | `b4e4db839bfdb7dcba63fba7cf5c10c4d390f82a63c424a395b4658075db953e` |
| `SDUI/go/runtime/input_lifecycle_test.go` | `ae63c036e196dfced5dd84d9141eee692a03eaa88112fbbe1d57eb8c664864de` |
| `SDUI/go/runtime/input_publication_test.go` | `4f042bf3bf64fd749703f02b090bd7196aefd6bc34778fa5a6e2ed12ca4becab` |
| `SDUI/go/runtime/session.go` | `e7fb4a480d215ad2b0e5f6a0e5185f65c0a524efc0e11ae04217efb0e46578e7` |

Canonical original observed SHA-256: `8199c275fecff3f6d03835540a3e586dd418ff70470f641590d2d75ace8f51be`.
