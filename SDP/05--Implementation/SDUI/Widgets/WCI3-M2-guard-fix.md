# WCI3-M2 late interaction receiver correction

Coordinator implementation under the loaded SDP Worker role; independent review
by 01a11854-523d-7c11-adf4-f622b19b68c6. This supplements the five lane reports:
the earlier bridge report describes its frozen pre-correction bytes, not this
coordinator-owned delta. Whole-stage correction/archive review is approved; see Evidence-WCI3-M2.

The integration worker noticed interactionHandler built TextResult updates without
ExpectedDraftRevision. An independent overlay reproduced both command and tab
callbacks binding successfully, executing SDL once, then returning UI conflict with
domain succeeded because an extended receiver requires its actual draft revision.
Changing the fixture to a basic receiver would hide a supported-workflow defect.

The correction mirrors the existing legacy-handler extended-input guard: after
unchanged post-Execute handle/value/draft checks, include the originally captured
receiver.DraftRevision only for extended input. No revision check is relaxed; newer
drafts still reject, domain outcomes stay truthful and calls are never auto-replayed.
New retained integration tests cover read-only/writable/legacy command receivers,
tab success/conflict and no automatic replay. Targeted race tests pass in 1.148s.
An independent actual same-event redispatch overlay additionally proves rejection
with domain-not-called and an unchanged call count; retained supplemental evidence
is distinct from the committed test's no-automatic-replay assertion.

The native CLI adds an optional false-default --command-load flag. It changes only
the fixture source to declare page/loadText and routes the existing Load button
through that shared command, with the same SDL Plan and read-only receiver. The
independent reviewer removed only the import, flag and guarded block and obtained
byte-identical previous CLI source. No default fixture, SDUI, runtime or dependency
byte changed. The new native command_load variant verifies the actual command
presentation before testing Load, copy, read-only behavior and silent updates.

Main and independent reviewer native builds match SHA256
2579df7a405f18396c34a62498d4bfd6461a6a60176f39864e6cf36c9dcbb352.
The 224-path manifest is 19bc6aa20386470ab6a0b0ef2246db59fea6409ba13f8709fa777e5c998b6d6e,
with exactly these three paths different from the previous 222-path candidate:

| Path | SHA256 |
| --- | --- |
| SDL/go/bridge/interaction.go | 4c06ad5758b3107bb5c0c688c550214017d2adb45742ef4471403179a893f0f5 |
| SDL/go/bridge/text_interaction_test.go | 389e498190027342df2e3dbbe47586692cb43c8aaaebe29e3b9a2683a4d5b3c0 |
| SDL/go/examples/text/cmd/native/main.go | 295085c0c1488ca4fdc1c51dc525569830da726541df4654128f1ee35c3c99f8 |

The remaining 221 paths match exactly. Main guarded the original-workspace
integration against prior hashes and retained backups. The final native matrix
and affected bridge/values/commands/panes race run are fresh. Unchanged full SDUI,
default text fixture, collections/non-Fyne and SDPTool evidence is reused under the
reviewer's explicit applicability assessment, recorded in post-review-delta.json.
No blind repeat of the eight-minute unchanged text fixture is required. The earlier
0078 archive and suspended approval remain explicit; this report does not relabel
their results as the corrected candidate.
