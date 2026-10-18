# WCI3-M1 native pilot history

Evolving candidates only; none of these entries is final whole-stage acceptance.
Main uses the actual SDL values fixture, Xvfb :190 1600x1200 and XTest. Condition
commands set provider/domain/reload circumstances; gestures execute through native
controls. The harness and inspector preserve actual source paths/typed state.

## Pilot 1 — 78344a92

Binary SHA-256 `78344a92a90036df0da2248ddbb07de58e07997021a4ebad4396a0bec919d4ac`,
HEAD a28b3cc plus evolving worker changes. Exact native binary retained under
`/tmp/wci3-values-native-pilot1`; source was not frozen for final evidence.

- Boolean workflow: eight checks passed, including real click/Space, programmatic
  mute, readonly/disabled refusal and native teardown.
- Choice workflow: seven checks passed. Equal display labels select their distinct
  IDs; disabled choice/Escape invoke no action; options replacement revokes an
  open popup; removed accepted ID stays diagnosed until explicit repair.
- Numeric workflow: invalid intermediate/out-of-range/off-grid raw drafts remain
  visible without SDL; Tab does not commit; Enter accepts exact integer 42. Native
  increment click then produces no action and times out. This is an open host
  routing/geometry issue, not a passed workflow.
- Actual OS inspection found intrinsic fixture widths too narrow for useful
  demonstration: tiny slider, truncated choice label and narrow numeric text.
  Fixture owner adds explicit fill widths, without changing identity or domain
  semantics. Number step chrome also appears wrong; host owner investigates the
  failed click with actual inspector geometry and raw input records.

Raw pilot logs and OS screenshots: `/tmp/wci3-m1-pilot1`. Final archive will retain
relevant failure evidence and independently verified reruns; no failed pilot is
relabeled as final acceptance.

- Slider workflow: six checks passed, with actual held mouse/key input separating
  Change from one release Commit, and Escape preventing a later drag-release Commit.
- Typed action failures: eight checks passed for errors, malformed replies, newer
  numeric drafts and replaced option generations after real domain execution.
- Reentrant observer: two checks passed; an accepted newer Change prevents the
  original automatic gesture from acquiring a fresh target and calling SDL.
- Legacy text alongside typed fields: four checks passed for basic text Commit,
  explicit Load receiver and real text-only SDL dialog Accept. No M2 text extensions
  or IME proof is implied.
- Mixed Go Accept failure workflow: 13 checks passed for unknown persistence
  errors and succeeded-domain/newer-draft or option-generation conflicts. The open
  form retains the outcome, blocks a second Accept, remains editable and preserves
  the attempt in its Cancel receipt without rolling back persistence.
- Host diagnosis of number step failure: delegated Button renderer had size zero
  while the outer adapter inspector rectangle was positive. Owner added actual
  Resize synchronization and a native canvas tap regression; OS rerun still required.

- Lifecycle workflow: nine checks passed. Five candidate failures preserve both
  main and open-form drafts; successful reload retains the main draft and closes
  the form with its unaccepted proposal reset. Parent hide closes without persistence.

## Pilot 2 — d8d726fb

Binary SHA-256 `d8d726fb7fbbe4b1a4b90f1b0ca715a393e3475d3e0e91bc2a22382b77469ee9`,
retained `/tmp/wci3-m1-pilot2-native`. Explicit fill widths and step-button Resize
correction are present. Numeric workflow passes all 11 checks, including actual
increment 42→43, invalid-draft step refusal and Escape restoring accepted 43.
Mixed form passes eight checks, including explicit child SDL commit surviving
Cancel and one successful Go Accept atomically persisting captured mixed values.

OS inspection confirms real step-button chrome now appears. An initial reduced
image presentation appeared to clip unrelated left-hand text. Independent inspection
and coordinator reinspection of the same PNG at original detail show every label
and value fully visible; no artifact/product defect is supported. The tentative
visual finding was withdrawn without a product change.

- Separate nonmodal mixed form repeats all eight form checks successfully.
- Keyboard workflow: seven checks pass for duplicate-label Home/Down/Return,
  Space reopening after focus restoration, source-order Tab/Shift-Tab and numeric
  arrow Commit. Initial harness attempt misread Inspector.focused as a Handle;
  it is the documented path string. Correcting the assertion requires no product
  change; the rerun passes.
- Readonly workflow: six checks pass for number typing/Enter/arrows/step refusal,
  actual native copy/paste into a different text field, silent trusted replacement,
  blocked choice popup and unchanged slider under pointer/key input.

## Pre-review 0e92fdf7 and corrected 6034b111

Candidate 0e92fdf7 (88-file inventory 37cdefe3) passed full SDUI race (host133.379s),
affected SDL race including values105.019s, and original full SDPTool tests, all
exit zero. It is NOT final acceptance: independent review found a slider tap
recapturing after reentrant Change and held-key state surviving focus loss. The
first actual tap reproduction initially read a nested intermediate snapshot too
early. A queued UI-owner state command exposes the real later SDL call and accepted
80, confirming the bug. All WCI3 changed-state waits now use this completion barrier
and reassert the condition; the initial false-positive tap run is excluded.

Bounded fixes preserve the original automatic-tap target and clear interrupted
held-gesture bookkeeping on focus loss. Independent reviewer regressions pass.
6034b111 repeats the actual reentrant tap with zero SDL calls and unchanged accepted
25 after the completion barrier. Final native/suite evidence remains pending.

A separately approved contract refinement retains the invalid editable baseline
of an already-required accepted-empty choice across unchanged reload. Newly required
blank or ineligible nonempty accepted values still reject. The fixture has an
explicit --required-empty startup variant for actual native verification.

Review also found the owner's slider value-feedback obligation was not yet met:
thumb position alone is insufficient. A bounded native label/value addition is
in progress, preserving independent validation feedback and measured geometry.
No final M1 completion is inferred from these passing intermediate checks.

- Corrected 6034b111 native held-key focus-loss sequence passes four checks;
  returning to the slider permits the next opposite complete gesture to commit once.
- Required-empty variant passes four checks: invalid empty initial state executes
  nothing, unchanged reload retains invalid emptiness and accepted revision, and
  an explicit real beta selection commits through SDL and clears validation.

- Actual 0e92 held-key focus-loss reproduction also fails at the expected missing
  next Commit after correlated focus transitions. An earlier attempt had inspected
  focus before the native Tab was delivered; the harness now waits for each actual
  focus transition rather than relying on immediate state-command scheduling.

## Final frozen candidate — 3c1abba7

Coordinator, host worker and independent reviewer builds match SHA-256
`3c1abba76adca3cdaf9d710302ebc6a16575150acfb1e6effcd85bb9b18b2ef6`.
The 88-file tested manifest is `cbd842e334b646f75a61374e063e85d3ec7075da0bf0fff4549708fe5e9deef6`.
Final SDUI race, affected SDL race/fixtures and original SDPTool suites exit zero.
The original modal forms run passes eight behavior/teardown checks but emits
Fyne Preferences load EOF at app/preferences.go:105, matching the prior observed
signature. Preserve that run separately; one fresh isolated-configuration repeat
passed all eight checks with empty stderr. No dependency correction or proven causal explanation is claimed.
The coordinator's preliminary all-stderr-empty summary was premature and corrected
before final acceptance. Remaining final audit is recorded in WCI3-M1-final.
