# WCI3-M1 host/preparation worker handoff

Status: owned product/test/document files frozen for independent review and main
integration; whole-stage acceptance remains coordinator-owned. No confirmed
remaining host/preparation blocker is known at this handoff. M2 extended input,
text/IME, public scalar SVG rendering and WCI3 final acceptance are not claimed.

## Scope and authority

Worker: Noether; exclusive clone `/tmp/sdp-sdui-widgets`, branch
`sdui/widgets-wci3`, base `a28b3ccc72b1ecc44300d4aa7c3cd679f4bc9625`.
Read/reused SDP entrypoint, sdp-worker 2.0.0, document workflow, SDUI instructions,
canonical original Values-and-text/Acceptance/phase plan, WCI3 runtime/frontend
memos and the agreed scalar layout API. The coordinator owns Sessions, management,
canonical evidence and native acceptance. No original workspace, modules, other
lane files, commits or branch changes were made.

Owned implementation is checkbox/slider/select/number only. `DocumentRequest.Choices`
is an exact normalized-path option inventory. Missing, invalid and hidden inventories
reject before native resources. Preparation keeps provider/host/widget capabilities
distinct; connected binder and HasBinding postchecks precede final native resource
preparation. Pure geometry validation still precedes application binding. The checked
connection receiver allowlist now includes the four scalar kinds and legacy input.
Documents with no select do not acquire an unrelated choices publication, preserving
legacy `.2` state revision behavior.

## Implementation and invariants

- Native Check, Slider, Select and number Entry project runtime FieldState. Ordinary
  sync retains their objects and gesture; only accepted runtime tickets publish
  native resources. Failed Entry edit preparation restores its last accepted raw
  native projection without another user Change.
- Automatic checkbox/option/number-step/slider-tap Commit uses the exact returned
  `change.Field.Target`. A reentrant observer is never authorized by recapturing a
  fresh automatic target. Slider release/key-up and explicit number Enter capture
  the current field. Programmatic updates, initialization and sync are muted.
- Slider uses Fyne's normalized visual track only. Legal uint64 ticks are applied
  through runtime EditTick and shared numeric.Grid; keyboard stepping is exact even
  at large origins. Holding a key or pointer changes proposals without Commit;
  release commits once. Escape suppresses the remainder of a canceled drag.
- Number invalid raw text stays visible with runtime validation; stepping refuses
  it. Read-only preserves focus, selection and native copy while suppressing edits,
  cut/paste/undo and Commit. Steppers do not introduce extra Tab stops.
- Fyne Select's label-index/private-popup limitations require the existing bounded
  native menu adapter. Every option captures its own exact OptionTarget, preserving
  duplicate-label IDs. Disabled/stale targets cannot activate. Ordinary sync retains
  popup focus; generation replacement closes even an originally empty opening.
- FieldMeasurer uses actual native themed minima and one fixed ellipsized feedback
  row. Diagnostic length cannot change admitted geometry. Runtime retains complete
  feedback. Actual renderer track/thumb and number subcontrol rectangles are exposed
  via `Inspect.fields`; `Inspect.choices` exposes actual option rect/clip/ID/target.
  Native controls are included in actual background omission inventory.

Number step chrome reuses the native Button renderer under a nonfocusable adapter.
The first native pilot exposed zero owner Size during renderer Refresh; synchronizing
the renderer owner's public Size fixed the actual chrome and hit behavior. A native
TapCanvas center-hit regression now checks both step buttons and verifies entry-only
Tab behavior. Main confirmed actual XTest increment/SDL call and OS painting on pilot2.
A later alleged left-label clipping finding was withdrawn after original-detail
inspection showed complete text in the same artifact; no product change was made.

## Verification

All commands below ran in the isolated clone with `GOCACHE=/tmp/sdp-wci0-gocache`.
No network/sandbox block occurred and no permission bypass was used.

| Working directory / command | Result |
| --- | --- |
| `SDUI/go`: `go test ./host/fynehost ./preparation -run '^$'` | PASS after actual API integration |
| `SDUI/go`: `go test ./host/fynehost/... ./preparation -count=1` | PASS host 6.029s, preparation 0.006s |
| `SDUI/go`: `go test ./host/fynehost ./preparation -count=1` after binder ordering | PASS host 6.116s, preparation 0.007s |
| `SDUI/go`: `go test -race ./host/fynehost/... ./preparation -count=1` | PASS exit 0, host 88.911s, preparation 1.047s; no hung process |
| `SDUI/go`: `go test -race ./host/fynehost -run TestScalar -count=1` | PASS exit 0, 27.382s; includes final large-origin/mute regression |
| `SDUI/go`: `go test -race ./preparation -count=1` after legacy no-extra-publication guard | PASS exit 0, 1.047s |
| `SDL/go`: `go test ./bridge ./examples/values -count=1` | PASS bridge 0.193s, values 5.408s |
| `SDL/go`: `go build -tags desktop -o /tmp/wci3-m1-noether-final ./examples/values/cmd/native` | PASS exit 0 |
| root: `git diff --check -- SDUI/go/host/fynehost SDUI/go/preparation` | PASS |

The complete host race run preceded the additional large-origin test and the final
preparation-only legacy no-extra-publication guard; the subsequent scoped race runs
cover those final changes. No repeated broad suite was run without a new change.
Main owns the final whole-workspace SDUI/SDL suite and exact-candidate native rerun.

Persistent host regressions cover Change-before-single-Commit, held slider key/drag,
Escape/late release, failed-resource zero Change/Commit/no speculative paint,
invalid numeric drafts and readonly copy, duplicate-label option IDs and stale opening, reentrant observer refusal with accepted nested state retained,
hidden unsupplied choices before resources, popup focus/empty generation, actual
track/thumb/step inspector objects, fixed feedback at fonts 14/20/28, native center
step hits/Tab order, safe53 large-origin step and programmatic mute. Preparation
regressions cover binder failures/no-op adapters before native resources, distinct
hidden host/provider/widget/read-only capabilities, supplied-empty inventory and
legacy state revision preservation. Independent runtime/layout tests are not
claimed as worker-authored evidence.

Main-reported provisional native successes include checkbox, choice identity,
held slider, numeric pilot2 (11 checks), modal/nonmodal mixed forms, keyboard,
typed failures/reentrant observers, legacy text bridge/dialog and mixed Accept
failure handling. Those are coordinator evidence, not this worker's own native
execution or a final full-stage approval. The original worker pilot2 binary is
unchanged at `/tmp/wci3-m1-noether-pilot2`, SHA256
`d8d726fb7fbbe4b1a4b90f1b0ca715a393e3475d3e0e91bc2a22382b77469ee9`.

## Frozen candidate identity and changed files

Final desktop artifact after slider input and visible-value corrections: `/tmp/wci3-m1-noether-value`.
SHA256: `3c1abba76adca3cdaf9d710302ebc6a16575150acfb1e6effcd85bb9b18b2ef6`.
The unchanged `/tmp/wci3-m1-noether-slider-fixed` / `6034b1115793b4888db7766c18d29fce97722a51f88ca5d855a013faedee1df8` artifact predates visible numeric feedback and is also superseded.
The unchanged `/tmp/wci3-m1-noether-final` / `0e92fdf7c1085b2a865d8d3a08e82f6f227860bb127c7d4f7aec0b6f87b0f305` artifact is superseded pre-review evidence, not the final candidate.

The following inventory covers this lane's 13 changed source/test/document files,
not other workers' changes. Its aggregate is SHA256 of sorted UTF-8 lines
`<file SHA256>  <relative path>\n`:
`3a957d5e035786e166d9d90ed619831c85862391c98108e52947a9d772f0420c`.

| Changed path | SHA256 |
| --- | --- |
| `SDUI/go/host/fynehost/README.md` | `12d892ab3d1e26c32603eae7add6a17929f6bc1395b35d148e2fce1a4374758b` |
| `SDUI/go/host/fynehost/admission/admission.go` | `ddafedeba66fdf48fc483eebfe4be968a4fa77494deead9e69e01576fb5d83aa` |
| `SDUI/go/host/fynehost/document.go` | `2b2370242904e16efc9007db3b0705c61614a4d984105e4bd6637e8df718b469` |
| `SDUI/go/host/fynehost/document_commands.go` | `25802189906d15acf87dbfe085eb8df4415debf3461706474c32bb37cd60cab6` |
| `SDUI/go/host/fynehost/document_sync.go` | `0daba7ec0a627586fc183a6b8f46ce4fcb9da2274adda42124a4627112754e0c` |
| `SDUI/go/host/fynehost/inspection.go` | `b2c604cec83433675ae03747832d7b00f5f45984422740d5e34f368186ab1e8e` |
| `SDUI/go/host/fynehost/scalar_choices.go` | `7fbe95923eef796490d8f6a87b43aa97b3b420b3deb68b643c4fb44dd5e7bb01` |
| `SDUI/go/host/fynehost/scalar_control.go` | `19049aa6b53560eabd0cb695ab85f2515a0f1fbda61328a0150ae9278cc29244` |
| `SDUI/go/host/fynehost/scalar_control_test.go` | `6367c744db1165f32b46f16cf3cd972123b87b4ec49442196bb255b10178f57d` |
| `SDUI/go/host/fynehost/view.go` | `0414e7b6baeb623d65f20fe82a19e97fabea143759d4bd23f013d11fd1973779` |
| `SDUI/go/preparation/capabilities.go` | `2ec88ca9e781f06c4b600f981a5e4fe87a94689398151a37f298a432ae829d42` |
| `SDUI/go/preparation/prepare.go` | `fbd200d3d38837021ca9d762b9769796f3321e465696a24d8faa76dc3d106771` |
| `SDUI/go/preparation/prepare_test.go` | `2b540e46a3546748ecb458c062e3bf3c02392653aa4d08baada9ef480d680a0c` |

Worker evidence file: `WCI3-M1-integration-worker.md` (excluded from its own hash
inventory). Full combined candidate inventory, integration and approval remain main
responsibility. Native clipboard readonly/focus workflows and final exact artifact
acceptance must be finished/recorded by main; headless success is not that proof.

## Bounded post-freeze slider review corrections

The coordinator accepted two independent Mendel findings and explicitly reopened
only `scalar_control.go`, `scalar_control_test.go` and this report. All other owned
file hashes remain unchanged. M2 memo work was paused; no M2 product work occurred.
The former `0e92...` artifact is pre-review evidence and must not be certified.

1. Slider Tapped previously called EditTick and then commitCurrent, recapturing a
   newer target after a reentrant Change observer. A synchronous native tap now
   marks its input scope and passes commit=true to the existing EditTick dispatch
   path. That path uses the original returned `change.Field.Target`; no post-tap
   recapture remains. Explicit later drag/key release still captures current state.
2. FocusLost previously retained the held key identity. It now cancels gesture
   bookkeeping, clears held-key/name/changed state, resets pointer dragging and
   preserves the current uncommitted proposal. Focus loss itself emits neither
   Commit nor Revert. Obsolete releases are inert; a new opposite-key gesture or
   new drag can complete normally.

Three persistent regressions cover stale automatic tap with a reentrant other-field
draft (zero callback, old accepted value, proposal and nested native draft retained),
then a fresh successful tap; lost key with late old releases and one subsequent
opposite-key Commit; lost pointer with late DragEnd and one subsequent new drag.
These supplement the existing native held-input tests; old pilot outcomes alone
could not detect these two defects.

Delta commands, all with the same /tmp Go cache:

- `SDUI/go`: `go test ./host/fynehost -run 'TestScalarSlider|TestScalarGesture|TestScalarPointer' -count=1`: PASS, exit 0, 0.601s.
- `SDUI/go`: `go test -race -overlay /tmp/wci3-host-review-jn0shb3w/overlay.json ./host/fynehost -run 'TestReviewerSlider|TestScalar' -count=1`: PASS, exit 0, 26.129s; exact independent repro plus every persistent scalar regression.
- `SDL/go`: `go build -tags desktop -o /tmp/wci3-m1-noether-slider-fixed ./examples/values/cmd/native`: PASS, exit 0.
- root: `git diff --check -- SDUI/go/host/fynehost/scalar_control.go SDUI/go/host/fynehost/scalar_control_test.go`: PASS.

Product/test files are frozen again with updated hashes in the inventory above.
Independent recheck was explicitly requested from Mendel with both file hashes;
its result remains coordinator/reviewer-owned and is not claimed here before the
reply. Main must integrate this two-file delta with its original-baseline guards,
run the affected actual pointer/key/observer proof and final whole-candidate suites.
No further known uncorrected host finding is asserted at this handoff; independent
review and exact final native acceptance remain open until their actual results.

Independent delta recheck received after the handoff: Mendel reran the exact
original `TestReviewerSlider` overlay with race checking and reported PASS in
3.349s on `scalar_control.go` SHA `c50a466d092261a8fc91d5ab0f28f935e93204fe12c39d9bf63043f7c6396a5b`.
He inspected original tap Change.Target preservation and FocusLost cancellation
and closed both reported findings at source/component level. His remaining M1
source audit and main's final native/candidate evidence are separate and pending.
The new binary's `-h` output also confirms the fixture's independently owned
`--required-empty` variant is included; the worker did not launch a native window.

## Visible slider value: final bounded inventory correction

The coordinator accepted the existing owner-card slider value-feedback obligation
and authorized the native label correction plus inspection-only `labelText`.
Changed product/test files relative to the preceding handoff are solely
`scalar_control.go`, `scalar_control_test.go`, `scalar_choices.go`; this report
contains the evidence delta. Layout, source schemas, runtime, modules and other
lanes remain untouched. M2 memo/product work remains paused.

The existing Label row now displays the shortest round-trip proposed float64
first, then ` *` when Dirty, then ` | ` and the original source label. Initial,
held-gesture, revert, accepted and programmatic read-only state use the same copied
runtime FieldState; label synchronization emits no Change or Commit. Validation
continues to use its independent Feedback row. The source label and numeric source
arguments are unchanged. `fields[path]["labelText"]` reads the actual native
`widget.Label.Text`, rather than reconstructing text from runtime for inspection.

A value-first truncated prefix alone was explicitly rejected by the coordinator.
The final minimum reserves the complete bounded number and dirty marker. For each
native font, measure full native Label strings consisting of 24 repetitions of
each numeric ASCII glyph (0123456789-+.e), followed by ` * | …`, and take their
maximum width. Shortest round-trip binary64 output is bounded by 24 characters;
this conservative reservation covers the full value, marker, separator and native
ellipsis. Max this constant with the original source-label and native slider minima.
Only the trailing source label may be ellipsized. The reservation is independent
of live proposal, accepted value, dirty/validation/read-only state and constraints;
therefore edits do not change the admitted field geometry. Initial minimum width
may rise to honestly reserve that output; row heights/shared geometry API do not
change. Existing layout rejects insufficient final allocation before publication.
Gibbs confirmed this requires no layout file/API change.

During implementation, a test caught that a single glyph's native ink width is not
an adequate substitute for complete native label advance; final measurement uses
the complete repeated strings. Another test comparison initially included ReadOnly
policy in a geometry equality assertion; it now compares rectangles independently
of the intentionally changed policy. Explicit Apply tests use runtime BatchRevision,
not event Sequence, for the already-reviewed separate batch counter. These were
corrected before freeze; earlier failing runs are not passing evidence.

New persistent tests verify actual native label/inspector text during a held gesture,
no Commit while held, unchanged label/control/feedback geometry, independently visible
validation, silent Escape, read-only programmatic update and refused user mutation.
A short source label `X` and fonts 14/20/28 cover negative extreme finite float64,
zero, small exponent and safe53 integer formatting; the full value plus marker and
ellipsis fit the measured Label region, with no value-dependent minimum. These are
native-widget/component measurements, not claims of actual OS painting.

Final delta commands (same isolated module roots and /tmp Go cache):

- `SDUI/go`: `go test ./host/fynehost -run 'TestScalarSliderVisible|TestScalarSliderLabel|TestScalarFeedback' -count=1`: PASS, exit 0, 0.397s.
- `SDUI/go`: `go test -race -overlay /tmp/wci3-host-review-jn0shb3w/overlay.json ./host/fynehost -run 'TestReviewerSlider|TestScalar' -count=1`: PASS, exit 0, 33.072s; both original reviewer regressions and every persistent scalar regression.
- `SDL/go`: `go test ./examples/values -run TestValuesActualRequestAdmissionAndTrustedMute -count=1`: PASS, exit 0, 1.942s; checks actual fixture admission with the larger honest minimum.
- `SDL/go`: `go build -tags desktop -o /tmp/wci3-m1-noether-value ./examples/values/cmd/native`: PASS, exit 0.
- root: scoped `git diff --check` for the three changed code/test files: PASS.

Mendel independently confirmed exact hashes before/after his race run (8.861s),
covering the original slider reproductions and visible-proposal/fixed-geometry/
tiny-label/scientific/font tests. He closed visible numeric feedback at source/
component level and reported no further source finding. This is scoped independent
approval, not whole-M1 acceptance. Main's queued state-barrier native `slider_reentry`
and `slider_focus`, held/revert/programmatic visible value, original-detail OS images,
combined candidate integration and final full suites remain coordinator-owned.
All product files are frozen again at the current inventory and final binary above.
