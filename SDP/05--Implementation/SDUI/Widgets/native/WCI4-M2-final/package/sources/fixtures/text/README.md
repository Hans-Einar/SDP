# WCI3-M2 actual SDL text editing fixture

Run from SDL/go:

```sh
go test ./bridge ./examples/text
go test -race ./bridge ./examples/text
go build -tags desktop -o /tmp/wci3-text-native ./examples/text/cmd/native
/tmp/wci3-text-native
/tmp/wci3-text-native --nonmodal
/tmp/wci3-text-native --required-empty
```

Main owns native execution, display/IME/clipboard configuration, screenshots and
receipts. A compiled binary or a Fyne test app is not OS evidence. No binary is
stored in the repository. Preserve each pilot binary until its evidence is retired.

The title is `SDUI WCI3 Text Editing`, initially 1100×850. The Editor tab contains
extended single-line, required and multiline inputs; the Other tab retains a
separate draft. The multiline editor contains deterministic long Unicode text
inside a separately scrolling outer pane. A split separates a readOnly preview.
Load really executes SDL and updates that preview. Save actions really execute
SDL and echo only the exact captured text into their self receiver. The composed
Text Form has ordinary extended text inputs, real SDL Accept, explicit Cancel /
Close and a child Commit whose domain work survives later form cancellation.

All fields use explicit opt-in arguments and fill their available width. The
native inspector reports actual field/control geometry and Entry observations.
The root `WCI3-M2-fixture-API.md` is the path and JSON protocol contract.

stdin only arranges trusted conditions or reads state. Commands cannot invoke
callbacks, edit drafts, synthesize keys/pointers, operate clipboard or fake IME.
`action` and `accept` arm the next real invocation outcome, then reset to success.
`set` uses checked programmatic Apply. `set ALIAS same` accepts the exact current
draft programmatically without firing Change/Commit. Other slots retain normal
dirty-draft admission. `constraint` attempts a source successor and retains the
previous source and live bundle on failure. `fail` deliberately fails successor
preparation/publication. Every command emits its exact line in `command-result`;
`state` emits the actual host inspector plus counters.

`reject-edit ALIAS` arms a diagnostic publication gate. The next prospective edit
with changed proposal/draft revision and unchanged accepted value fails before
publication; focus/resize or a model successor do not consume it. Only an actual
native edit proves the approved history-reset exception: the same focused Entry
must restore current authoritative draft, emit no Commit, and may clear its native
undo/redo/caret/selection/scroll state. Component tests prove gate rejection and
one-shot consumption, not native history or IME behavior.

Action counts count actual SDL handler entry. Domain counters separately record
successful Save/Accept work; a later non-echo or draft conflict does not roll it
back or replay it. Errors/malformed results in this fixture occur before domain
writes. Load is a read and does not increment domain commit counts. Closing emits
actual sessionClosed and pending request state, plus final counters. Prior M1
fixtures/evidence remain unchanged.

`--required-empty` starts only main Required with empty accepted text while still
required. It exposes initially invalid editable absence and unchanged-reload
retention until explicit user repair. FormName and defaults remain unchanged;
sourceSHA256 always reflects the actual selected source.
