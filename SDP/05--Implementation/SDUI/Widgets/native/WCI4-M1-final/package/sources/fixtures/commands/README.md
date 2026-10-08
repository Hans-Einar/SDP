# WCI2-M2 real SDL commands fixture

This bounded integration fixture exercises shared commands and ordinary composed
modal/nonmodal dialogs with actual action-core Go registrations. It preserves the
existing panes/collections fixtures. It is not a release or native-pass record.

From `SDL/go`:

```sh
go test ./bridge ./examples/commands
go test -race ./bridge ./examples/commands
go build -tags desktop -o /tmp/wci2-commands-native ./examples/commands/cmd/native
/tmp/wci2-commands-native
/tmp/wci2-commands-native --nonmodal
```

The main title is **SDUI WCI2 Commands**, initially 1100x750. Settings is modal by
default; `--nonmodal` changes only its source `modal` argument and uses the host's
standard Fyne window, titled **SDUI WCI2 Settings**. Main supplies real native
pointer/keyboard/window-close inputs. Do not infer native focus, OS menu placement,
IME or window-manager behavior from in-process tests.

## Paths, actions and inspection

The frozen harness contract is `../../../../WCI2-M2-fixture-API.md` from this
folder. Main paths are `page/run`, `page/flag`, `page/view/header/toolbar`,
`page/view/header/actions`, `page/itemMenu`, `page/view/body/items` and
`page/settings`. Dialog inputs are `page/settings/form/body/name` and `note`;
terminal buttons are under `page/settings/form/footer`. Settings also owns
`page/settings/editMenu` (Text actions), targeted at Name, with the local Mark
toggle `page/settings/mark`. Secondary-click Name, then Escape closes only the
menu; subsequent Escapes revert its dirty draft and cancel the clean dialog.
This ordering requires actual native-input evidence. The same menu also has
Open child (`page/settings/editMenu/childItem`), through settings/openChild to the
root-sibling modal `page/x` titled SDUI WCI2 Child. Its ordinary childInput and
childClose use .4 source scale on each axis. This adds no Settings tab stops.
The shorter child path deliberately exercises dynamic parent-canvas ordering
when opened from nonmodal Settings. Native parent scaling/close cascade needs
separate actual-input proof. Run uses Primary+R,
Flag Primary+K, Settings Primary+D. Run's symbolic `fixture-run` icon is supplied
through `DocumentRequest.Icons` using Fyne's standard MediaPlayIcon.

Run/Toggle/Inspect are real typed SDL actions returning text to separate visible
footer inputs preview/flagPreview/itemPreview respectively (each x=1fr).
Save consumes captured Name/Note text through resolved DialogFieldValue paths and
returns boolean Accepted plus text Message, without an output receiver. Save's
in-memory domain values/commit count are distinct from runtime accepted inputs.
Single-line Enter commits locally; it does not invoke Save. Cancel/Close never
invoke Save. Domain success followed by UI conflict is retained truthfully.

Stdout is NDJSON `{event,data}`: `ready`, `state`, `action`, `dialog-result`,
`pending`, `load`, `command`, `command-result`, `error`, `closed`. Every stdin
command receives an increasing commandId and ok/error result. Use that barrier,
provider keys and actual published state rather than delays. Inspector controls
are flat records with `canvas`/`title`; surface/menu geometry comes directly from
host Inspect. Coordinates for nonmodal content are local to its native canvas.
The fixture does not calculate native menu/control positions.

## Stdin controls

These set trusted test conditions, release providers or request parent lifecycle;
none invokes/simulates a user command, dialog action, key, pointer or input edit.

| Command | Effect |
| --- | --- |
| `state`, `pending` | Inspect state or started provider keys. |
| `enable Run`, `disable Run` | Set Enabled on exactly `page/run`. |
| `enable Flag`, `disable Flag` | Set Enabled on exactly `page/flag`. Names are case-sensitive. |
| `action Run\|Toggle\|Inspect success\|error\|malformed\|draft-conflict` | Set next actual SDL action outcome. |
| `accept true\|false\|error\|malformed\|draft-conflict\|resource-conflict` | Set next actual Save outcome; does not invoke Accept. |
| `complete KEY success\|error\|empty\|invalid` | Release exactly one already-started lazy load. |
| `replace-items`, `cancel-load` | Replace collection generation or revoke pending request, preserving a releasable late reply. |
| `reload` | Publish compatible successor; published old dialogs close with reload result. |
| `fail profile\|binding\|resource\|guard\|stale` | Attempt specified failed successor without invoking SDL. |
| `parent-hide`, `parent-show` | Notify host before native Hide; Show does not reopen dialogs. |
| `resize WIDTH HEIGHT` | Resize parent; finite positive logical dimensions <=32768. |
| `close` | Notify host, dispose session/providers, close parent and emit closed once. |

Outcome injection is one-shot. False changes no domain values and permits a fresh
Accept. Error/malformed leave domain outcome unknown; failed publication after a
successful Save preserves domain success, blocks replay, and permits actual edits
then Cancel/Close. Draft-conflict edits an owned captured Name draft inside Save;
resource-conflict arms failure inside Save so earlier focus guards cannot consume
it. Command draft-conflict edits its own Run/Toggle/Inspect receiver after capture. These
hooks exercise existing guards; they are not application reconciliation APIs.

Provider barriers ignore cancellation deliberately so native tests can deliver a
late response and verify product rejection. Closing releases every fixture waiter.
