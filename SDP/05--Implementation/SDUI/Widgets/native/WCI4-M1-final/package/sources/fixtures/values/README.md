# WCI3-M1 actual SDL values fixture

This bounded fixture exercises checkbox, slider, number and stable-ID select with
real action-core registrations. Existing basic text Load/Save and text-only SDL
Accept are compatibility cases; extended input/multiline/clipboard/IME remain M2.
The separate mixed form uses an explicit Go Accept adapter. This README is not a
native acceptance or release record.

From SDL/go:

```sh
go test ./bridge ./examples/values
go test -race ./bridge ./examples/values
go build -tags desktop -o /tmp/wci3-values-native ./examples/values/cmd/native
/tmp/wci3-values-native
/tmp/wci3-values-native --nonmodal
/tmp/wci3-values-native --required-empty
```

Main window **SDUI WCI3 Values**, initially1100x850. --nonmodal changes both form
dialogs to standard windows titled **SDUI WCI3 Mixed** and **SDUI WCI3 Text**.
The coordinator supplies actual pointer/keyboard/window inputs; no fixture control
simulates them. Full source paths and exact NDJSON/state schema are frozen in
`../../../../WCI3-M1-fixture-API.md` relative to this folder.

## Actual domain work

AcceptFlag/AcceptLevel/AcceptMode/AcceptCount/ChildCount each take and echo one
Value, respectively boolean/integer/text-ID/integer/integer. Each has its own
explicit self receiver, typed Commit selector and ScalarResult. Shared numeric
validation preserves source/raw provenance and rejects unsafe SDL conversion.
Load returns supplied in-memory text; SaveText commits existing basic input.
TextSave uses WCI2 DialogFieldValue and DialogAcceptResult without a hidden receiver.

Mixed form flag/level/mode/count/note are proposals until Go Accept. ChildCount's
explicit SDL Commit is independent persistence and survives Cancel. Domain maps
and counts record actual handler work; failed UI delivery after successful domain
work never claims rollback or automatic replay. False Accept commits no domain
state. Domain errors/malformed text results and succeeded conflicts follow the
existing runtime outcome/replay/terminal-result protocol.

## Conditions and barriers

Stdin commands produce one increasing commandId/status result. `state` exposes
actual host geometry and Fields, actionCalls, changeCounts keyed full field path,
domain, source/actions hashes and pending actual runtime collection requests.
`closed.pending` is inspected after Session.Close. Options have no hidden collection
or provider loader. Use observed command-result and published state as barriers.
Action logs precede domain work; the later state reflects the actual outcome.

| Command | Effect only on trusted conditions/lifecycle |
| --- | --- |
| `action NAME success\|error\|malformed\|draft-conflict\|options-conflict` | Arm next named actual SDL call; options-conflict only AcceptMode. |
| `accept mixed\|text true\|false\|error\|draft-conflict\|options-conflict\|malformed` | Arm actual form acceptance; options-conflict mixed only, malformed text only. |
| `observe TARGET draft-conflict` | Arm one actual Change observer to publish a newer edit and invalidate the original automatic gesture capture. |
| `options Mode\|FormMode initial\|removed\|disabled\|reordered` | Replace complete options and advance generation, including repeated labels. |
| `set TARGET VALUE` | Checked silent programmatic AcceptedValue. Does not overwrite a newer dirty proposal; numeric raw value uses shared exact grid parsing. Choice literal empty means no selection. |
| `enable\|disable\|readonly\|writable TARGET` | Checked silent condition mutation. |
| `reload`, `fail profile\|binding\|resource\|guard\|stale` | Complete successor or specified failure preserving live state. |
| `parent-hide`, `parent-show`, `close`, `resize WIDTH HEIGHT` | Explicit native parent lifecycle/size through host hooks. |

Exact target aliases: Flag, Level, Mode, Count, FormFlag, FormLevel, FormMode,
FormCount, ChildCount. Action names: AcceptFlag, AcceptLevel, AcceptMode,
AcceptCount, ChildCount, Load, SaveText. No raw draft, invoke, Commit, Accept,
Cancel, key or pointer command exists. Outcome/observer injections are one-shot.

Draft-conflict changes the actual source/receiver inside the actual handler.
options-conflict replaces the set after capture; a retained ID does not revive
its old generation. Removed/disabled accepted IDs remain diagnostic; repair with
real selection or restore valid options and revert before compatible reload.
Observe reentrance clears its arm before its accepted newer edit, avoiding loops;
two Changes occur but the original stale automatic Commit cannot invoke SDL.

Tests exercise connected request admission, typed actual actions, programmatic
mute, strict command arguments, reentrance, error/malformed/newer-draft/options
conflicts, mixed proposal/child Commit/Cancel, mixed Go and text-only SDL Accept,
and failed/successful reload. They do not establish native gesture counts, popup
geometry, focus, read-only pointer behavior or OS lifecycle. Main owns that proof.

## Required-empty startup variant

`--required-empty` changes only main Mode to `required=true,value=""`; it can be
combined with --nonmodal. FormMode, option sets, paths, domain behavior and stdin
protocol remain unchanged. The main choice starts invalid and editable with no
automatic selection or SDL action. Compatible unchanged reload retains that
already-required empty accepted absence and its validation. Main's native variant
then explicitly chooses beta to produce one real AcceptMode SDL repair. Newly
required blank values and invalid nonempty accepted IDs still require runtime
rejection under the refined contract; this fixture adds no alternate reload rule.
