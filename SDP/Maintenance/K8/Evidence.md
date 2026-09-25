# K8-M1 — consolidation evidence

Date: 2026-09-25. Base: `29d5828`. Scope: the [plan](Plan.md).

## Delivered

Reviewed all 12 original backlog cards across three registered boards. Full
merge `KBO-SDP-000001` transfers KB-SDP-002 and KB-SDP-016 to new queued
KB-SDP-017. Both sources are superseded with history and successor links.
SDL/SDUI Refs retain their IDs, point to #017 and record their filename moves.
The other eight original cards have explicit review outcomes and remain separate.
The review matrix explains each ownership/dependency boundary.

The current proposal records the owner's native XFMD SDP tab with KanBan / SDL /
SDUI subtabs, model-derived viewpoint/use-case navigation, status/card browsing
and sdptool delegation. It retains all advanced process scope from #002 while
separating the first browsing delivery from later synthesis/history work. The
metadata filename and validity contract remain open. No source proposal was
silently promoted into an implemented schema or language rule.

## Checks

Run from the repository root:

```sh
python3 SDP/Maintenance/K8/verify_consolidation.py
python3 SDP/Agents/KanBan/examples/verify_lineage.py
python3 SDP/Maintenance/L1/verify_documents.py
git diff --check
```

The K8 verifier is a dated snapshot check, not a new KanBan CLI. It validates
schema/replay/physical placement, exactly one valid CardState per card, primary
Ref resolution, all 12 review events, all three merge participants and unchanged
historical ledger prefixes. It expects 21 cards and 96 events, 11 remaining
backlog cards (7 SDP / 3 SDL / 1 SDUI), one queued primary (#017), and #010 still
at gate-review. The implementation Traceability ledger is byte-identical to base.
The existing lineage suite covers full/partial merge and split plus 15 invalid
cases. The documentation check preserves frozen records and generated outputs
and verifies current local Markdown links/anchors.

Results: all four checks passed. The documentation check verified 105 frozen
records/ledger prefixes, 574 generated outputs, 2,171 local file links and 130
fragments. No whitespace errors were reported. Manual scope review compared both
source proposals with the transfer table and retained later milestones in #017.

## Limits and external handoff

No runtime, Go parser, tool executable, generated SDL view or XFMD code changed.
No GUI execution or new runtime tests are claimed. Source inspection confirmed
the launcher's explicit-model support and fixed project registration. XFMD's
companion remains in its own repository; the [plan](Plan.md) records its required
producer-link and owner-layout update. External bare file URLs are not checked by
the Markdown link verifier. Historical ledger paths remain historical addresses.
