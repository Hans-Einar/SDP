# WCI3-M2 — extended native text and IME

Implemented, integrated, verified and independently reviewed at `69d0a3329608fdfc0dfd25417d9b051119ba4e16`.
[Candidate](candidate-WCI3-M2.json) records 224 source/test/document/dependency paths
on parent d6742742. The tested inventory is 19bc6aa20386470ab6a0b0ef2246db59fea6409ba13f8709fa777e5c998b6d6e;
only architecture/requirements delivery-status prose changed after freeze.
Coordinator, worker and independent reviewer builds produce identical native SHA256
2579df7a405f18396c34a62498d4bfd6461a6a60176f39864e6cf36c9dcbb352.
Guarded integration matched all 224 paths in the original workspace, preserving
unrelated changes and retaining pre-integration backups.

Explicit input-policy argument presence selects extended text; absent arguments
preserve basic 0.2/0.3 behavior. Widget.Value/Draft remain the single text authority.
Native Entry owns caret, selection, undo/redo and internal scroll. Exact typed
Commit/self-echo guards prevent stale or non-echo acceptance and automatic replay.
Required/CR-LF constraints, read-only copy, placeholder and form Accept/Cancel follow
the reviewed [text contract](../../../04--Design/SDUI/Widgets/Values-and-text.md).
Same displayed bytes preserve native history; changed programmatic bytes reset it.
An actual rejected native edit restores authoritative text on the same focused
Entry and may reset editing history. Failed Commit/reload/probe retain history.

[Final native archive](native/WCI3-M2-final/README.md) retains 14 fresh workflows,
92 passing checks, 5 published openings with exactly one receipt each,
empty fixture stderr and clean session/provider teardown. Actual XTest and external
X11 clipboard exercise single/multiline editing, native undo/redo, CR/LF refusal,
read-only copy, failed actions/reloads, page retention, constraints, modal/nonmodal
forms, native wrap/scrollbar and independent ancestor scrolling. Configured real
IBus/XIM proves preedit does not alter drafts or call SDL, consumed Return/Escape
stay out of control dispatch, and ordinary/explicit Return subsequently saves once.
OS captures were inspected at original resolution by coordinator and reviewer.

The selected GLFW source retains 145 pinned upstream files and notices with one
X11 filtered-key conditional patch. Exact reverse-patch verification and independent
upstream ZIP checksum review pass; both native module roots explicitly select it.
Actual SDL-root build info and IME proof are retained. SDUI-root packaged-helper IME
proof remains WCI4. Fixture stderr is clean; isolated D-Bus/portal environment
permission warnings are retained, and visible preedit panels are not claimed.

Full SDUI race suite passes (host 155.521s). SDL bridge, runtime/codegen, all non-Fyne
packages and collections/commands/panes/values/text fixture suites pass. The final
text race run finishes in 508.408s within its unchanged ten-minute limit; the earlier
600.275s timeout is preserved. Original SDPTool `go test ./...` exits zero, including
preserved concurrent governance work and cached unchanged packages. Exact commands,
flags, logs and source manifests are retained without inventing uniform environments.
Fresh bridge/values/commands/panes and 14 native runs verify the final correction;
unchanged SDUI and default text/collections/non-Fyne/SDPTool evidence is reused
with the independently approved applicability recorded in post-review-delta.json.

Independent reviewer 01a11854-523d-7c11-adf4-f622b19b68c6 approved all five scoped
lanes, integrated source and final evidence. Actual native testing found and closed
missing Primary+Shift+Z routing. Late independent integration review found command/
tab TextResult updates omitted the extended receiver draft revision; the bounded
correction preserves the captured guard and is proved by actual SDL tests and a
new native shared-command Load workflow. Prior 0078 approval was suspended; its
13-run archive remains explicitly superseded. Other failed pilots record native undo grouping,
unsupported Ctrl+Home/End and the fixture's intentional stale-publication revision;
the corrected harness follows reviewed behavior, not weakened product acceptance.
Detached Entry measurement now omits text only where independent minimum-size
checks prove equivalence; the live control retains the complete text.

X11/Xvfb evidence does not establish Wayland, window decoration or other-platform
IME behavior. Legacy standalone RuntimeView still explicitly rejects new adapters;
connected DocumentHost fixtures are the runnable route. Provider previews, matching
consumer packaging and all-family final review remain WCI4. KB004 glyph fidelity
and KB005 preview ownership remain separate. No merge, publication or full-card
completion is inferred.
