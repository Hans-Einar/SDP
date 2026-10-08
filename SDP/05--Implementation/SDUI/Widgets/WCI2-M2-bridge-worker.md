# WCI2-M2 SDL bridge and commands fixture Worker handoff

## Boundary and candidate

Selected by coordinator under KB-SDUI-003 / PLAN-SDP-0022 WCI2-M2. Role: SDP
Worker (sdp 1.1.1, sdp-worker 2.0.0, shared document-workflow reused). Canonical
original Panes-and-commands.md governs; runtime/frontend/layout M2 memos supplied
the concrete seams. No material contract departure. No WCI3 implementation.

Worktree `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci2`, observed HEAD
`a4f2c22435d91fe07935b8b9d5fcf46fafc6e36b` plus uncommitted multi-owner changes. Go1.27.1 linux/amd64,
pinned Fyne2.8.1. Only bridge, NEW commands fixture and authorized root handoffs
were written by this Worker. No commits, branch changes, PM/Session edits, or
runtime/frontend/host/layout/panes/collections edits. Main owns those records,
integration and native acceptance. Other lanes were active during these checks;
this report identifies this lane's exact files, not a frozen whole-repository
candidate. Coordinator reconciled fancyfs v0.0.1 module metadata outside this lane.

## Delivered behavior

- Canonical command/dialog callback owners use InteractionHandler; deduplicated
  handles and cleared promoted Widget binding prevent double registration. Basic
  buttons retain the old Handler path; M1 tab selectors and legacy values remain.
- Typed CommandContextItemID, CommandChecked and DialogFieldValue/FieldPath with
  owner/type/path preflight. Runtime/frontend resolve owned input identities once
  per copied owner plan, then Accept reads captured String fields. No private
  lexical resolver, live input reread, select coercion or WCI3 mapping.
- Closed TextResult/default and DialogAcceptResult. Actual engine signatures and
  returned records are checked. Accept requires boolean decision/text message,
  forbids all receiver/revision fields, returns zero extra updates, and checks
  captured field count <=256. Runtime retains combined-batch/revision ownership.
  False is rejected; malformed/Execute errors unknown; true is succeeded before
  post-domain state checks. Replay blocking and Cancel/Close receipts stay runtime
  responsibilities; no reconciliation or retry framework was added.
- Real canonical action-core Run/Toggle/Inspect/Save registrations. Three visible
  footer receivers respect the existing one-connection-per-widget invariant:
  Run preview, Toggle flagPreview, Inspect itemPreview. Draft-conflict changes
  precisely that action's receiver; Save changes captured Name after domain save.
- Native CLI title SDUI WCI2 Commands, 1100x750; --nonmodal changes only Settings
  modal source argument. Standard icon resource, host inspector geometry passthrough,
  actual OnDialogResult observer and guarded NativeParentHidden/Closed hooks.
  Stdin controls set conditions/barriers/lifecycle only. Exact case Run/Flag map to
  page/run and page/flag for enable/disable. No fake command/Accept/Cancel/input
  dispatch. Main supplies actual native input.

## Reproducible checks

Commands below ran from `/tmp/sdp-sdui-widgets/SDL/go` unless stated otherwise.

| Command | Result |
| --- | --- |
| `go test ./bridge ./examples/commands` | PASS bridge0.152s / commands1.405s |
| `go test -race ./bridge ./examples/commands` | PASS bridge2.117s / commands16.051s, final owned bytes |
| `go test -race ./bridge ./runtime ./codegen ./examples/panes ./examples/collections` | PASS bridge2.668s, runtime cached, codegen3.365s, panes51.507s, collections27.073s; preceding added bridge regressions, production bridge unchanged |
| `go build -tags desktop -o /tmp/wci2-commands-native ./examples/commands/cmd/native` | PASS, fresh artifact after projection corrections |
| `/tmp/wci2-commands-native --help` | PASS; documents nonmodal flag |
| `git diff --check -- SDL/go/bridge SDL/go/examples/commands` (clone root) | PASS |

Tests cover one actual shared SDL handler via toolbar/menu/key; typed proposed
checked values; captured context vs current selection and stale context refusal;
Accept true/false/error/malformed/draft/resource conflict and edit-then-Cancel
preserving outcome; strict preflight/no partial handlers; owned-field isolation;
caller plan immutability; per-action receiver conflict; implicit vs legacy route;
trusted enable/disable and invalid control refusal; failed/successful reload draft
retention; provider late delivery, one-shot barriers and close cleanup.

Initial tests correctly exposed invalid duplicate receiver connections and SDL
canonical ordering in the fixture; both source issues were corrected without
parser changes. Runtime owner corrected illegal item/separator/refbutton snapshot
projection discovered by connected fixture preparation; fresh fixture checks pass.

## Limits and next step

Coordinator reports native commands pilot started from an immutable copy of
artifact 1b213abf: zero preflight SDL calls and one toolbar Run passed; remaining
checks ongoing under /tmp/wci2-m2-pilot1-commands on :190. This is preliminary
coordinator-reported evidence, not final-candidate acceptance. Native binary
overwrites are held during that pilot.

No native process was launched by this Worker. Build and in-process Fyne/runtime
checks do not prove actual keyboard/pointer/menu placement, modal focus blocking,
nonmodal parent lifetime, OS close protocol, IME or WM chrome. Main has isolated
:190 and frozen harness protocol; it owns actual-input acceptance and final whole
candidate identity. Host/native work remains concurrent; binary hash does not
prove later edits. Independent review is separate, not claimed here. Session0010
and plan/card/evidence reconciliation are handed to Main per exclusive delegation.

Native artifact: `/tmp/wci2-commands-native`
SHA256 `1b213abfb32a825c7d86a844111fb1e958e4fd2bb8d2ac941b7af454242077ed`.

## Exact owned file SHA256 inventory (after dynamic-parent fixture addition)

- `SDL/go/bridge/bind.go`: `44259d5e8e51868f3dd63d09b52b83257bf3523db9e87479371feaaeab9378ae`
- `SDL/go/bridge/interaction.go`: `7ff0fdbe268062ff00cc6d03f097563a314d5135b6f777573f9502860db39d6b`
- `SDL/go/bridge/values.go`: `33ce0c15614ce5979b5aea9be07d4c599a4bfe54557355b466a72c27ebe48539`
- `SDL/go/bridge/results.go`: `09162e9e4d33e17dbb6390653fb9dbebc1da2f21909f744ba6494448e5f8d018`
- `SDL/go/bridge/commands_test.go`: `42ac49b26fabe730de8f4238ad2088b888f275f76e62092f679e7294e05f0533`
- `SDL/go/bridge/README.md`: `37b05beaea2e6ae9cafc67ac3c3708d0395d14ddea3a659849cde8efa3b870aa`
- `SDL/go/examples/commands/cmd/native/main.go`: `e1ce093c53397d1ae2df85c2c8afe64e55ecb74f7def97365a0765358f8a303b`
- `SDL/go/examples/commands/controls.go`: `27278fa57c0e019dbf520a1f763cc6fdb2b504282b3a84ea0b8e179330bf07ba`
- `SDL/go/examples/commands/fixture.go`: `98768ceaa295e8b93c80e758391e70790d3257288a1a2be2bcba93d6576b440a`
- `SDL/go/examples/commands/fixture_test.go`: `3007da6f8ce3deaa5a12208ed38b6582b773e10ceeae971a2aa647cd62c46fb6`
- `SDL/go/examples/commands/providers.go`: `14a0b4444c3b3342eed50cb5d0733164af6a1c30cd08b48ddff9ebb6ce10fc4c`
- `SDL/go/examples/commands/source.go`: `d23b373149de45b082603673a9e975e72ed01a0d846fa4c4b055266b18883fb2`
- `SDL/go/examples/commands/README.md`: `9b9fe0e83cfb729cf7e1ef322063c76bb650262108a2d119b159026433dfa4ff`
- `WCI2-M2-fixture-API.md`: `6c7a1b39da9a172afce557aae0966005be23c9fa549e9601707a35d3803eb042`

## Independent scoped review follow-up

Reviewer 01a11854-523d-7c11-adf4-f622b19b68c6 reports all 14 owned-file SHA256
entries match the reviewed files and independently ran
`go test -race ./bridge ./examples/commands`: PASS bridge1.920s / commands15.891s.
Review covered per-receiver conflicts, promoted/legacy routing and condition-only
fixture controls/CLI; no remaining blocker in this lane. This is attributed
reviewer evidence, not a new Worker run or whole-M2/native acceptance. Product
files and binary remain frozen pending coordinator findings.

## Bounded A06 fixture-only addition and refreeze

Coordinator authorized the concrete missing dialog context-menu acceptance setup.
Settings now owns Mark (unchecked local toggle, no callback) and Text actions
context menu on Name, with exact paths page/settings/mark,
page/settings/editMenu and page/settings/editMenu/markItem. Source references use
settings/mark and settings/form/body/name because the containing definition scope
is page; no new scope/resolver or host behavior was introduced. Auxiliary
menu/command declarations do not consume ordinary dialog layout tracks.

Only commands source.go, fixture_test.go, README.md and the fixture API memo were
changed; this report's inventory is refreshed. Existing bridge files, control
protocol, domain callbacks and native binary were untouched. Prior independent
review/race and binary1b213abf evidence predate these four changed files.

Validation from SDL/go after this delta:
- `go test ./examples/commands ./bridge`: PASS commands1.588s / bridge0.142s.
- `go test -race ./examples/commands -run '^TestDialogContextMenuActualRequestAdmission$' -count=1`:
  PASS3.726s. Actual connected Request/Host.Adopt succeeds in modal and nonmodal
  variants, resolver identities target the owned Name input/local toggle, and no
  SDL handler runs during admission.
- Clone-root `git diff --check -- SDL/go/examples/commands`: PASS.

Source refrozen. Main/host own the next binary build and actual secondary-click
Name, Escape(menu only), Escape(draft revert), Escape(clean Cancel) sequence.
No native ordering pass is claimed by admission tests or the previous binary.

## Bounded dynamic-parent fixture addition and refreeze

Coordinator requested an actual native regression setup for the independently
found dynamic-parent ordering defect. Added root-sibling modal page/x, titled
SDUI WCI2 Child, scale-x/scale-y .4, ordinary input page/x/childInput and explicit
Close button page/x/childClose. Settings owns local effect command
page/settings/openChild targeting x, exposed only through the existing Name
context menu at page/settings/editMenu/childItem (label Open child). The suggested
name more conflicts with the existing bar submenu in this definition; childItem
preserves existing paths and the strict unique-name rule. Settings still has
exactly five input/button controls. No SDL action, stdin protocol, product bridge,
host code or binary change was made by this Worker.

Actual connected Request/Host.Adopt admission tests verify both variants, including
--nonmodal retaining Child modal, root-sibling target and dialog ownership,
ordinary input/Close identities, unchanged five Settings controls, and zero
preflight domain calls. Validation from SDL/go:
- `go test ./examples/commands ./bridge`: PASS commands1.832s / bridge0.158s.
- `go test -race ./examples/commands -run '^TestRootSiblingChildActualRequestAdmission$' -count=1`:
  PASS3.757s.
- Clone-root `git diff --check -- SDL/go/examples/commands`: PASS.

Source SHA256 d23b373149de45b082603673a9e975e72ed01a0d846fa4c4b055266b18883fb2;
source/tests/docs/API memo refrozen and inventory refreshed. Earlier independent
review/native artifact evidence predates this addition. Noether owns the surface
ordering fix and next separate binary; Main owns actual child pointer Close,
Settings canvas/.4 scale proof and parent-hide cascade. Admission alone makes no
claim about those native geometry/lifetime outcomes. Both owners were notified.
