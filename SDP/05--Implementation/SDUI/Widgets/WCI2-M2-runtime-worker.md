# WCI2-M2 runtime worker evidence

Date: 2026-10-08. Worker James; bounded coordinator-selected M2 runtime lane.
Worktree /tmp/sdp-sdui-widgets, branch sdui/widgets-wci2. Observed HEAD a4f2c22435d91fe07935b8b9d5fcf46fafc6e36b.
This is HEAD plus dirty lane files below, not an unchanged-commit test claim.
Canonical original Panes-and-commands substantive approval 9a9c8710; latest read
source hash c8af54b42a5c3094eb34dcda294899936c9358f6896e81ee15d78a5dc496f0d7.
Reused SDP entrypoint, Worker and document workflow; main owns Session/PM records.
No commits, branch changes, host/parser/layout/bridge/module/management edits.
Other exclusive owners' concurrent changes are preserved and excluded below.

## Delivery

Implemented concrete canonical command/presentation/menu/surface state and detached
accessors/capture. Frontend ResolveInteractions and ResolveDialogField own source
identity; no duplicate resolver. Basic buttons remain legacy unless explicitly
opted in. Shared button/menu/key routes bind one canonical owner; local effects
invoke only their selected action. Checked/exclusive state is atomic and silent
under programmatic Apply. Native key capture requires the active canvas.

Menus use stored MenuScope identity independent of current revision for cleanup;
invocation separately checks current revision and exact context generations. Old
selection cleanup cannot dismiss a replacement. Success dismisses in the command
publication; failed invocation dismisses without replay or restamping. Runtime
contains no GUI input-scope object or native dependency.

Dialog open generations, actual opener/parent, dynamic modal child routing,
remembered/first-eligible focus fallback, published acknowledgment, terminal
results and bounded synchronous acceptance diagnostics are implemented. Accept
captures owned hidden-page fields, excludes nested dialogs and validates UTF-8,
exact revisions, reserved changes and the combined 256-write limit. Unknown or
succeeded UI-conflict blocks replay while retaining edit/Revert/Cancel/Close.
Reentrant close/reload finalizes the exact queued receipt after callback unwind.
No generic transaction/reconciliation/ResolveAccept framework was added.

M1 pure probes and final tickets remain the publication authority. Failed geometry
or resource preparation retains live surfaces/drafts/provider requests; diagnostic
sequence/outcome changes promote no pending bundle. Real native lifetime revocation
cannot be vetoed by a fallible gate. Successor starts all surfaces/menus closed,
preserves generation/sequence and compatible checked state, and discards only
unaccepted dialog drafts while retaining accepted values and ordinary main/page
drafts. Ordinary 0.2 and WCI1 collection/provider behavior remains covered.

Updated runtime README and shared runtime-contract to actual M2 APIs, typed bridge
selectors/result mode and clear native-evidence boundaries. Frozen API memo is
WCI2-M2-runtime-API.md. No WCI3 implementation or future family claim.

## Verification

30 substantive M2 test functions plus subtests, alongside the complete existing
runtime suite. Cases include shared routes/basic opt-in, exclusivity, exact old/new
menu cleanup, clicked versus selected item and copied callback payload, state and
resource conflicts, false/unknown/succeeded acceptance, malformed outcome pairs,
256+1 returned updates and 257-field preflight, hidden/nested/UTF-8 fields,
published-only results, reentrant close/reload, modal dynamic children, parent hide
with pending provider and rejected geometry, forced teardown, copied snapshots,
successor reload draft semantics, remembered focus/inactive keys, tab replies to
commands with reserved activation owners excluded, and frontend-valid projection.

Final commands run after substantive fixes:

- SDUI/go: `go test -race ./runtime ./layout ./preparation` PASS
  (runtime 1.651s, layout 1.405s, preparation 1.064s after exact native revoke addition).
- SDUI/go: `go vet ./runtime` PASS.
- SDL/go: `go test ./bridge` PASS (cached).
- SDL/go: `go test ./examples/commands -run TestTrustedEnabledConditionsAndStrictNames`
  PASS 1.013s, modal and nonmodal fixture admission.
- `git diff --check -- SDUI/go/runtime SDUI/docs/runtime-contract.md` PASS.

Earlier runs exposed reserved-pane regression from broadened callback capture and
invalid projected menu/ref-button schemas; both corrected with retained regression
tests. Earlier bridge setup failed on missing fancyfs sums; owning lanes corrected
the dependency, then bridge and fixture commands passed. A mistaken test path
SDUI/go/bridge was corrected to the actual SDL/go/bridge package; it supplied no
bridge evidence. No unexecuted test is represented as passing.

Independent reviewer separately reported fresh overlay/full runtime race PASS
1.737s on properties.go hash 86c41d03e30d4fe76d6c177d9ab26fe5827cdca6eb9ce02d0cbf367f55bb8b48,
closing five reported runtime regressions. This is a received interim result, not
this worker's independent approval or whole-M2/native certification.

## Remaining boundary

Runtime lane ready for integrated verification. A coordinated host handoff added
RevokeSurface(exact target, lifecycle reason) after the initial freeze: real native
OnClosed bypass now revokes self plus descendants without geometry veto. Its new
regression proves failed ordinary close, forced exactly-once result, repeat cleanup
and stale callback protection after reopening. Documentation/API are updated;
latest runtime/layout/preparation race and vet above include this addition.
The reviewer subsequently reported scoped approval of this additive operation on
surfaces.go f3fd9b8ce69c78f69bec7a884bb67ce5ca8874f70d724f196ae36d3503b79c51,
with independent full runtime race PASS 1.726s, including the new native-loss
regression. Host call-site review and final native proof remain pending host freeze.
 Main/host own actual native menus,
modal/nonmodal canvas lifecycle, OS focus and SDL end-to-end evidence on the final
candidate. No outstanding material contract departure known. Keep the files frozen
except coordinated fixes found by integration/review; no M3/WCI3 work selected.

## Lane file SHA-256 inventory

```text
cd18db2f684869123fafe2d2cfeedc78e21f27b04a1086a4325830521c9a8eb6  SDUI/docs/runtime-contract.md
105331f975d114f91b6c23beceeb499368541c6b11751c49225c635d541343c8  SDUI/go/runtime/README.md
38a99715e2d06ee73911d5d51450341cf3065e6baf0004677cf7cc2f78f91988  SDUI/go/runtime/collection_mutation.go
ba177e5752e6c1e9f80284555b190a80d23f3ebd75d4972ff05b3ea1fc85d6f3  SDUI/go/runtime/command_capture.go
f211ebe0bd3d667a5fe5c23b2a5420983520a5ecd9d915604b490af0e147ea17  SDUI/go/runtime/command_context_test.go
3ec820279918c34e9a89d941afb3b849a2f6f68fe6d47e805cca720ec993b020  SDUI/go/runtime/command_dispatch.go
95ddcaa429f24d19ed3a4f8efc77c6e790ab746c895807639ab3a6782aef39ee  SDUI/go/runtime/command_init.go
d378b6ed25c15a25a9caa85f3f8bae97084a04473b6d74937c186f19dfdd0493  SDUI/go/runtime/command_state.go
5dafe4ae075c6249f77ef616fb0f8f6887cf3b5d574b081bf401d1cbc8726e29  SDUI/go/runtime/command_successor.go
8c4fc527e09c4862c5a71a6933383f0a2a172c2f702be2113c38a7c2ff3d9877  SDUI/go/runtime/command_types.go
6bfe7e52d32e0fdca54d48ed94f2295cdfa398893dd5039a35066dfabe6c03df  SDUI/go/runtime/commands_test.go
2aa860e5eb750fbed9c1d4cf501f5e6c98830c84956b0dca017785655fe19e2a  SDUI/go/runtime/dialogs_test.go
0494d3322a5c0e68ddbfa886bbeb400cd1ac66e303eff3106f75a354ad25ff14  SDUI/go/runtime/events.go
0a8d9ab719c5dde028d0e5e3ff597ecbc422df9b84f9b81a419a64ecd4f6647c  SDUI/go/runtime/interactions.go
c2c14c799970f3f10feacb58a9f465f79a502b9df4450dbd195de557a7d02fd9  SDUI/go/runtime/pane_activity.go
49de3192892ae4489ad4156a9845ff9528ebfb453fd2eaacf83f7b36c321d267  SDUI/go/runtime/panes.go
86c41d03e30d4fe76d6c177d9ab26fe5827cdca6eb9ce02d0cbf367f55bb8b48  SDUI/go/runtime/properties.go
4754d5284ecf65c47987765545634c5897bcb73d63e58166f350f02df1af6b5f  SDUI/go/runtime/reload.go
7839c530cd33f3b870b0db9e966f7eb60125937d53e8815be469a9847eecaadd  SDUI/go/runtime/session.go
843bb05e0b1f77a59683308ebd8848bc268e668a1697814f1e755feacd13024d  SDUI/go/runtime/state.go
13524cb02a4dbd0160b0eaf1ca41576749a73a44c1353398fca2cedcec6550b1  SDUI/go/runtime/successor.go
d0532aabfa490d745e58a81f4493fdb55dd2f6fab59b148d05a7fc175ec3e263  SDUI/go/runtime/surface_lifecycle_test.go
f3fd9b8ce69c78f69bec7a884bb67ce5ca8874f70d724f196ae36d3503b79c51  SDUI/go/runtime/surfaces.go
c3503d705525326b65fd07972bb16ed96446c85dcaae52a39b4d616da79e8ca7  SDUI/go/runtime/types.go
```
