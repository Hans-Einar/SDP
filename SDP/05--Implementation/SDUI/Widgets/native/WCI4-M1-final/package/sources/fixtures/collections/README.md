# WCI1 native collection fixture

This fixture uses the actual SDL action-core parser, engine, Go registration and
bridge. Tree/list data is bounded and in memory; item IDs never open filesystem
paths. The native window title remains **SDUI WCI1 Collections**. Product host APIs
do not parse these fixture commands.

From `SDL/go`, build with the pinned native Fyne dependencies already available:

```sh
go build -tags desktop -o ../../wci1-fixture-native ./examples/collections/cmd/native
```

Run that binary with `--nested` for a scrolling body around the two scrolling
collections. The fixed footer preview is a sibling outside that outer viewport.
Paths remain `page/body/nav`, `page/body/entries` and `page/footer/preview` across
normal/nested variants. `--empty` starts the tree at an unloaded root. Both lazy
branches in the normal initial tree are explicitly controlled provider barriers.

## Control channel

Write newline-delimited commands to stdin. Commands run on the UI owner goroutine.
Use `tree`, `list` or exact widget paths; `preview` also names the preview input.

| Command | Fixture operation |
| --- | --- |
| `state` | Emit runtime state, viewport and row screen rectangles/clips, action count, source/action hashes and pending provider keys. |
| `pending` | Emit sorted keys of provider requests that have actually started. |
| `complete KEY success\|error\|empty\|invalid` | Release exactly that provider barrier. Unknown outcomes reject without consuming it. |
| `cancel tree` | Revoke/cancel the tree's current runtime request; provider barrier remains available for a deliberately late reply. |
| `resize WIDTH HEIGHT` | Resize the native window and request matching host geometry. Invalid numeric arguments reject before resize. Too-small geometry must retain the last valid presentation. |
| `reorder list` | Reverse the current data order through `Host.Mutate` and runtime replacement; preserve stable IDs and compatible selection. |
| `shrink list COUNT` | Keep the first COUNT current items, testing selection invalidation and offset clamping. |
| `reset list` | Restore initial fixture data through replacement, including IDs that may have been deleted. |
| `hide tree`, `show tree` | Update visibility through the runtime property batch. |
| `disable tree`, `enable tree` | Update enabled state through the runtime property batch. |
| `reload` | Prepare and publish the same source/provider identities with a fresh candidate sequence and the retained SDL engine. |
| `fail profile\|provider\|binding\|layout\|resource\|guard\|stale` | Attempt a controlled failing reload as described below. |
| `bad-reload` | Alias for `fail layout`, retained for the first native pilot. |
| `capture PATH` | Diagnostic Fyne Canvas.Capture PNG helper; see its acceptance limitation below. |
| `close` | Revoke the bundle, release pending provider barriers and close the native window. |

The commands also accept collection paths instead of aliases. `shrink` is bounded
by current item count. Invalid arguments fail explicitly and do not become user
input events. Mutations do not simulate selection or activation.

## Failure boundaries

- `profile` uses an unsupported source header.
- `provider` supplies an invalid initial stable ID with a changed provider epoch.
- `binding` omits the required SDL module from the actual bridge binding attempt.
- `layout` assigns an outer width below native minima.
- `resource` returns an error from the existing PrepareResources hook. This is a
  resource-preparation failure injection, not proof of a particular Fyne allocator
  failing.
- `guard` fails the final publication guard after detached preparation passed.
- `stale` prepares a candidate, makes a newer preview draft in the old bundle, then
  attempts publication. The old bundle and newer draft must survive rejection.

Failure commands return an error intentionally. Inspect the emitted state and
`command-result` to compare source/model identities, focus, selection, offsets,
requests and action counts with the preceding state. A later valid `reload` should
still work. No failure command invokes the domain action.

## Barrier and JSON evidence

Stdout is newline-delimited JSON with `event` and `data`. Existing `action`, `load`,
`load-return`, `state`, `error`, `ready`, `closed` and `command` events are retained.
`command-result` correlates `commandId` (monotonic within the process), `command`
(the exact input line), `status` (`"ok"` or `"error"`), `candidateSequence`, `source`,
`modelRevision`, `actionCalls`, `bundleChanged` and optional `error`. Expected failure
injections report `status:"error"` and should retain `bundleChanged:false`.
`fail stale` deliberately changes the current preview draft while preserving the
published bundle. Both CLI close and OS close use one idempotent cleanup function;
`closed` records `sessionClosed`, empty pending barriers and action count after
revocation/cleanup and before quit. A successful complete acknowledgement
means the barrier was released; it does not mean the asynchronous result has been
accepted. Wait for the correlated `load-return` and subsequent state, or confirm
that the old request remains revoked. Never choose requests by sleep duration.

A useful native sequence is: `--nested`; select/activate a row with actual pointer
and keyboard; scroll the collection to its limit and continue to the body; check
the fixed preview sibling; focus the last row; `reorder list`; `shrink list 5`;
`reset list`; exercise `resize`, hide/disable during a pending load, explicit late
completion, recovery and each failure stage; finally close with pending replies.
Only actual native input establishes A08/A10 interaction evidence. The control
channel supplies deterministic state/lifecycle changes for A11/A14.

**Capture limitation:** main's pilot2 observed black Canvas.Capture PNGs while OS
screenshots showed the live UI. Preserve `capture` only as a diagnostic helper;
use OS ImageMagick `import` on the actual window for acceptance screenshots. Do not
claim native paint/input proof from this helper or from headless fixture tests.
Source/binary identities and the external XTest/OS screenshot logs belong to the
coordinator's exact-candidate native evidence.
