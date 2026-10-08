# WCI2-M1 panes fixture

This fixture uses the real SDL action-core parser, registered Go Page action,
engine and bridge. It exercises only tabs/page/split, existing String inputs and
WCI1 tree providers. The window title is **SDUI WCI2 Panes**. No M2 commands,
menus/dialogs or WCI3 values are enabled. IDs and provider data are in memory;
provider IDs never open filesystem paths.

From `SDL/go`:

```sh
go test -race ./bridge ./examples/panes ./examples/collections
go build -tags desktop -o /tmp/wci2-panes-native ./examples/panes/cmd/native
```

The default is horizontal. Run `/tmp/wci2-panes-native --vertical` to replace only
that split axis with vertical. Source paths, minima, proportions, action bindings,
control commands and JSON protocol are unchanged; the source hash and snapshot
axis identify the selected variant. Compatible reload retains that selected source.
`--help` documents the flag without creating a native window. For vertical native
acceptance, main uses Up/Down and vertical divider dragging instead of horizontal
axis input. Building this variant is not evidence that those native checks passed.

Main runs the binary with an isolated X11 display/XDG area. Actual XTest pointer
and keyboard input supplies user interaction. Headless tests and control commands
are integration checks, not native-input evidence. Use OS screenshots: inherited
Canvas.Capture limitations make that helper unsuitable for acceptance.

## Source and observation

The outer split (horizontal by default, vertical with `--vertical`) has positive
minima .25/.20 and initial proportion
.65. Its first child contains Overview/Notes tabs; Overview owns an editable input
and a scrolling tree with 64 stable rows plus a lazy group. Notes and the second
split child have separate editable inputs. A fixed footer shows the actual SDL
result `PreviousId -> PageId`. Preparing, reloading, same-page activation and
programmatic selection must not invoke Page.

| Identity | Normalized path |
| --- | --- |
| Split | `page/body` |
| Tabs | `page/body/workspace` |
| Overview input | `page/body/workspace/overview/edit` |
| Overview tree | `page/body/workspace/overview/nav` |
| Notes input | `page/body/workspace/notes/notesEdit` |
| Side input | `page/body/aside/note` |
| SDL result receiver | `page/footer/preview` |

`state` includes the host inspector's snapshot, widgets/drafts, focus, measured
viewports/rows and pane geometry, plus action count, provider keys, candidate
sequence and source/action hashes. Native automation obtains live coordinates
from inspection; do not hard-code pane chrome sizes.

Stdout is newline-delimited JSON `{event,data}`. `action` carries stable `page`,
`previous`, actual SDL sequence, call count and injected outcome; it is emitted
from the registered Go action, never a control-channel simulation. `ready` marks
initial native setup. `command-result` correlates command ID/text, status,
candidate/model/source identities, action count and whether the bundle changed.
`closed` follows idempotent host/provider cleanup. `error` is diagnostic; consult
state and outcome to distinguish domain failure from UI publication conflict.

## Fixture control channel

Send newline-delimited commands on stdin. Commands execute on the UI owner
goroutine. There is deliberately no user `activate`, synthetic pointer, keyboard,
or direct callback command. `select`, `ratio`, `collapse` and `restore` below are
trusted application operations for silent-state/regression checks; use actual
input for corresponding native acceptance.

| Command | Effect |
| --- | --- |
| `state`, `pending` | Emit diagnostic state or sorted provider keys that have actually started. |
| `complete KEY success\|error\|empty\|invalid` | Release exactly one provider barrier; invalid outcome/key does not consume it. |
| `cancel` | Revoke the tree request; its provider barrier remains available for deliberate late completion. |
| `action success\|error\|invalid\|draft-conflict` | Configure only the next real Page action. Does not invoke it. |
| `resource-error` | Fail the next final resource preparation once. Arm immediately before native tab input; do not interpose another state mutation. |
| `select overview\|notes` | Silent programmatic page selection. |
| `ratio NUMBER` | Strict finite programmatic split proportion; measured out-of-bounds requests reject. |
| `collapse first\|second`, `restore` | Programmatic explicit collapse/restore with retained expanded ratio. |
| `hide TARGET`, `show TARGET`, `disable TARGET`, `enable TARGET` | Change declaration intent via runtime Apply. TARGET is `workspace` (the entire tabs owner at `page/body/workspace`), `overview` or `notes` (one page). |
| `resize WIDTH HEIGHT` | Request positive finite window/host size at most 32768 per axis. Rejected geometry keeps the last valid presentation. |
| `reload` | Compatible reload preserving the SDL engine and provider identity. |
| `fail profile\|binding\|layout\|resource\|guard\|stale` | Attempt the controlled failing reload described below. |
| `close` | Close once, revoke live state and release outstanding provider barriers. |

Use `disable workspace` / `enable workspace` to test globally disabled tabs,
then supply actual pointer/keyboard input through the native harness. These
commands do not activate a page or call SDL; they retain the existing
`command-result` acknowledgement. `hide workspace` / `show workspace` target the
same whole composition, rather than altering individual page intent.

An action `error` returns a domain error; `invalid` returns a malformed typed
record. Both leave the old page/receiver and report uncertain domain outcome.
`draft-conflict` makes a newer accepted runtime draft mutation during Execute,
then returns a successful SDL record: UI result application must reject while
preserving that draft and the succeeded domain annotation. `resource-error`
exercises successful domain execution followed by final presentation failure.
Consumed event sequences never replay. Inspect retained page/focus/drafts and
call counts after each failure; future fresh tab gestures remain possible.

Failure reloads use unsupported profile, missing SDL module in actual bridge
preflight, too-small layout, resource-hook failure, final guard failure, or a
newer live draft after detached preparation. Each should preserve the live bundle
and invoke zero Page actions. `fail stale` intentionally changes the old preview
draft; all other preservation comparisons should use unchanged live state.

Provider `load` logs the exact request/key; `complete` acknowledges only release,
not accepted publication. Wait for its matching `load-return` and subsequent
state. Cancellation/hide retains the barrier to prove rejection of late replies.
Never choose a request by sleeping or assume a released result is already visible.

## Required native sequence

Edit Overview, scroll its tree, switch by pointer and header keys, return and
verify draft/scroll retention. Confirm previous/new IDs and one real SDL call per
changed page, zero for same-page/programmatic/fallback/reload. Start the lazy
provider, hide its page, deliver the obsolete result, reveal and verify no
automatic restart. Exercise current explicit recovery with a new request key.

Drag and keyboard-adjust the divider; collapse with positive minima, resize while
collapsed, restore with current measured bounds, and reload in expanded/collapsed
states. Check actual focus and agreement between visible geometry and snapshot.
Inject Page/receiver/resource/reload failures; speculative page selection must
never become visible. Successful independent draft changes must survive failure.
Finish with pending work and verify cleanup. Main owns actual XTest/OS captures,
exact binary/candidate identity and independent acceptance; fixture presence or a
successful build does not claim those passes.
