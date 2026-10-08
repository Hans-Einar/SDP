# WCI3-M2 bridge worker handoff

Status: lane implementation complete and owned files frozen; scoped tests passed
on the current host candidate. Main owns final
independent review, integration/native evidence and candidate selection. No native
window was launched by this worker. No merge/release or further milestone implied.

## Scope and candidate

Assigned Worker lane only: SDL bridge plus new `examples/text`, fixture API and this
handoff. No modules, runtime/frontend/layout/host, legacy fixtures or management
records changed by this worker. No commits, branch changes or codegen changes.
Branch `sdui/widgets-wci3`; selected M1 implementation base
`ea49991f2065679c93e39fe02a3d458fe72803df`; current coordinator record commit
`d6742742f0968661d15698518df688ffb1d28946`. Working tree contains concurrent lane
changes. These twelve hashes identify this lane, not a whole-product commit.
Loaded/reused SDP entrypoint, Worker 2.0.0 and document workflow; original canonical
Values-and-text including reviewed native-edit exception and actual M2 API memos.

## Delivered behavior

Explicit source-policy FieldState.Input alone selects extended text behavior.
TextResult self-receiver preflight applies regardless of request selector. Typed
ControlText reads String/Control.ValueRevision, forbids numeric raw/option capture.
Extended handlers recheck full FieldTarget after actual SDL Execute and require
byte-exact echo; updates check value/draft revisions and accept the captured draft.
No normalization, replay or rollback claim. Basic input/nonself Commit and legacy
Load remain available; Load can update an extended readOnly receiver with exact
draft guards. Text-only DialogAccept remains closed; no heterogeneous SDL mapping.

Actual SDL fixture provides extended single/required/multiline text, inactive tab,
split, independent outer scroll, readOnly preview/Load, composed text Accept and
child Commit, long Unicode and one-shot actual publication rejection. stdin only
arranges conditions; OS keys/pointers/clipboard/IME belong to Main. API memo fixes
all paths, aliases, text slots, action/domain counters and JSON command barriers.
`--nonmodal` and `--required-empty` are source variants with truthful source hashes.
No product seam departures or competing APIs introduced.

## Verification

Environment: Linux amd64, Go 1.27.1. Commands below run from `SDL/go`.

- `go test ./bridge`: PASS 0.164s. New typed text self/selector preflight, exact
  Unicode/empty/CRLF echo, errors/malformed/non-echo/newer-draft/state conflict,
  no replay, basic nonself and Load-to-extended-readOnly tests; prior tests retained.
- `go test ./examples/text`: PASS 124.591s on pre-required-empty fixture. Actual
  Request admission, modal/nonmodal paths, trusted mute/invalid controls, truthful
  domain failures, one-shot resource gate, failed reload preservation, text dialog
  Cancel/child Commit/Accept, and source tightening. This is component evidence.
- `go test ./examples/text -run '^TestActualRequestPathsAndTrustedConditions/modal$' -v -timeout 45s`:
  PASS 40.431s (Go regexp also selected nonmodal). No OS/native claim.
- `go test -race ./bridge ./examples/text ./examples/values ./examples/commands ./examples/panes`:
  Old host candidate: bridge PASS 3.612s, values PASS 204.330s, commands PASS
  51.342s, panes PASS 88.445s; text TIMED OUT at 600.275s (default 10m).
  At timeout TestPublicationDiagnosticOneShotAndFailedReload had run 13s;
  stack was drawing Fyne SVG background during initial Host.Adopt. Earlier
  text cases consumed the remaining budget. This is a failed run, not a waiver
  or a proof of native dispatch failure. Concurrent native timeout logs retained.
- `go test -race ./examples/text -count=1 -v`: PASS **508.408s** on Noether's
  independently reviewed detached measurement optimization, within the unchanged
  default 10m limit and with unchanged fixture bytes. All six tests passed:
  actual Request/condition paths 210.08s; five echo/domain outcomes 102.75s;
  one-shot publication rejection/failed reload 55.11s; dialog Cancel/child Commit/
  captured Accept 81.30s; source tightening 24.16s; required-empty initialization/
  unchanged reload/explicit SDL repair 33.84s. No race finding reported.
  Transcript: `/tmp/wci3-m2-text-race-final.log`.
  Transcript SHA256: `4b07ef369d24ec2c0ca1c99f04bb92e664e13929b54e291e3856a61bacc2c162`.
  The prior timeout above remains a failed result; this is a distinct current-host
  run, not a relabeling of its evidence.
- `git diff --check -- SDL/go/bridge SDL/go/examples/text`: PASS.
- `go build -tags desktop -o /tmp/wci3-text-worker-pilot1 ./examples/text/cmd/native`:
  PASS; SHA256 `d4cdefe4ca6cf7a9ad57bd336019dbb35b9a7ed5389c35365648fef1a7b8844d`.
- `go build -tags desktop -o /tmp/wci3-text-worker-pilot2 ./examples/text/cmd/native`:
  PASS; SHA256 `6595a4e8a9737ff499fdc460a89a2177a6ee70b466b1a1bcffa78735bade60c8`.
  Pilot2 includes required-empty and narrower diagnostic gate. Both binaries are
  immutable worker artifacts; dependencies/other lanes were evolving at build time.
  Main's separate pilot1 SHA256
  `0cff7b7bcd6a9e259edbb7b6c5d5fc3c822166750e9cf2c00b3bffd268c30015`
  is different; no source/binary equivalence claimed.

## Remaining evidence

Frozen scoped bridge/fixture review APPROVED by reviewer thread
`01a11854-523d-7c11-adf4-f622b19b68c6`: all twelve manifest entries matched;
actual SDL Request/Plans, condition-only native channel, failure/domain accounting,
source tightening/required-empty/dialog cases and typed self/exact-echo barriers
were inspected. Independent full bridge `-race` PASS 2.066s. Reviewer audited the
exact 508.408s fixture transcript/hash; did not rerun that costly suite or claim
OS proof. No concrete blocker. Main owns the final integrated manifest, native
archive and overall M2 gate.
Final exact-product build/manifest and native history/retention/wheel/IME/clipboard/
focus proof remain with Main. Host inspector exposes actual Entry observations,
not an invented editor/history state. The worker tests use a Fyne test app and
runtime events; they cannot prove physical input behavior. No M1 evidence modified.

## Pilot investigation (read-only)

Main preserved 15-second startup/interaction timeouts; these are not product
acceptance or confirmed failures. In pasteguard pilot1, the first state already
contains accepted 1100×850 split/outer geometry. Pinned Fyne 2.8.1 SetContent resizes
to max(current canvas size, content minimum) before this fixture's explicit Resize;
DocumentHost minimum is 1×1, so a transient small request may correctly fail the
pane minimum. No minimum guard changed and no causal timeout claim is made.

In worker-pilot2 readonly-startup45 raw tail, actual Load eventually executed once.
Both runtime Value/Draft (revision 2, clean) and native Entry text equal exact
`Blåbær 日本語 🙂 é`, with ReadOnly true. Teardown reports closed/zero pending.
The harness timed out before that action; this is useful diagnosis, not a native
PASS receipt. Main owns subsequent reruns on its optimized native candidate. No fixture
code changed in response to these observations. Worker pilot2 predates Noether's
subsequent native Primary+Shift+Z adapter fix and is not a final product candidate.

## Final lane readiness

All twelve owned hashes below revalidated after the final race run. Manifest file
`/tmp/wci3-m2-bridge-worker.sha256` has SHA256
`dc7cc4bdd7a919efc39f0c368e0bd25c89d20ec3bced149100a036db40f7e43f`.
No remaining bridge/fixture implementation blocker is known. Main's separately
reported native successes and harness corrections are not promoted here into a
final exact-product receipt. No worker test/build processes remain active.

## Owned-file SHA256 freeze

```text
4866f47630c01bc1ebd8e31b509ec8a1ebf03be516c46103a5c402309795004a  SDL/go/bridge/README.md
7a7d5cf98ec3bba68b6cfd18a552e0883d01ef3a932159017b43015de8272572  SDL/go/bridge/results.go
24b43064b4ffcddcf3a4097e75fc8772a57f03229deaa664f0c7790bcd79d04e  SDL/go/bridge/scalars.go
5e2f213cff287bf851c0422ed4daf217542e80ce56ecd1530e3714904ae140c5  SDL/go/bridge/values.go
7fbd925a8e21da1e3f00ddef51f0ecbb4d35e8dc026c9b8a00fc8cdbbca4a522  SDL/go/bridge/text_test.go
1c85721d3f8fc24ec9790b164083b3b1ba33fc6983e03b684d09bd7dddc1da1a  SDL/go/examples/text/README.md
1dcec36567926d4e1bab498a595dabf4b1e088be2a8163be2954e83067e6dc20  SDL/go/examples/text/source.go
defd5b50158bfa67a4dee8cda27f8171dd29279f7eb6afeccc2bfb8d605e05b3  SDL/go/examples/text/fixture.go
851f4653babcc80d07ec5abdb03a7bdf26f673b710678d174c8f9613c7b873c3  SDL/go/examples/text/controls.go
9a295c5fd04a839569cdb10ccf35b8da3b791a1e86b388e8a70dfa5297e1642b  SDL/go/examples/text/fixture_test.go
9bfc2eab722591c08a28d28171b56e5f45e472a5e3e2c50175780a1ec5dc5585  SDL/go/examples/text/cmd/native/main.go
bf808eda35217fb11ab7a08c16b650438e771174c5bade702b52e9344e7fd727  WCI3-M2-fixture-API.md
```

## Authorized remaining SDL non-Fyne regressions

Completed after frozen lane review, with no product writes or repeats of bridge,
values, commands, panes or text. Durable commands/logs/exits, dependency inventory,
source byte comparison and hashes: `/tmp/wci3-m2-sdl-nonfyne-regressions/README.md`
and `results.json` / `SHA256SUMS` in that directory.

`go test -race -count=1 -p 1 ./runtime ./codegen`: exit 0; runtime 1.048s,
codegen 2.919s. Remaining 20 non-Fyne packages: exit 0. In total 13 packages with
tests passed and 9 no-test packages compiled. All SDL Go and module bytes remained
unchanged. Exact package lists and log SHA256 are in the durable report.

No non-Fyne default-build package gap remains. `examples/collections` is the sole
other standard package not covered by this follow-up or Main's five listed frozen
passes; it depends on Fyne and needs matching current evidence or a separate
bounded run by Main. Eight desktop-only packages are inventoried separately for
Main's compile/native gate; no desktop or OS claims are added here.

## Subsequently authorized collections gap closure

`go test -race ./examples/collections/... -count=1` completed with exit 0:
collections PASS 18.031s (command elapsed 19.540s). No other expensive package was
repeated. Source and module bytes unchanged. Durable exact command/environment,
log, exit, source hashes and checksum manifest are in
`/tmp/wci3-m2-collections-regression/`. Log SHA256:
`b74c15003aa164f7655cc95f289dcb608ad1bb27f4711f446c541b3c0c716e9e`.
This closes the default-build SDL package gap identified in the preceding section;
desktop build/native evidence remains separate. The earlier non-Fyne artifact set
is retained unchanged. Gibbs receives both artifact sets for Main's archive.
