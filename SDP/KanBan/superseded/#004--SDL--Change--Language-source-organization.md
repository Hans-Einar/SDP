# KB-SDL-004 — Organize language development by process phase

| Field | Value |
| --- | --- |
| id | KB-SDL-004 |
| project | SDL |
| type | Change |
| created | 2026-09-25T10:32:43Z |
| source | owner-conversation-2026-09-25 |
| CardState | superseded |

## Selected scope and acceptance

Owner: Project Owner; implementation: Codex. Apply the accepted five-phase
profile to SDL and SDUI. Split the mixed shared architecture model by requirements,
structural architecture, container/library detail and implementation planning.
Keep stable symbols and all existing semantic facts. Phase placement does not
assign A0–A5 levels. Container detail gets container directories; shared libraries
and contracts keep one source instead of copies in every host.

A bounded explicit source-set input is necessary for parsing/projection across
files. It belongs in Go tool input, with per-file ASTs, positioned diagnostics,
explicit membership and revision tracking across all files. This does not adopt
System/import/namespace syntax from the historical workspace study. Keep one
combined SDL/SDUI system model distributed by responsibility unless the owner
selects independent models.

Acceptance: phased project entries, no maintained monolithic architecture source,
lossless canonical composition, working CLI/nav/selected projection, original
source provenance, invalid-member and cross-file diagnostic tests, preserved
historical evidence and working links. XFMD code and installer templates are
outside this delivery.

Plan: [S1](../../Maintenance/S1/Plan.md).

## Worklog

- 2026-09-25: selected and started on the owner's instruction. Located the original
  source-tree study; inspected current parser, exporter, broker and consumers.
  No explicit multi-file source loader currently exists.

## Owner clarification and current handoff — 2026-09-25

The owner now prefers a common root SDP process area for SDL/SDUI/SDPTool and
asks for system/container boundary guidance before migration. This supersedes
the earlier separate-process-tree working assumption. See the
[system-boundary study](../../03--Architecture/System-Boundaries-study.md). System boundaries in that study are
recommendations, not adopted language rules.

CardState: in-progress → ready while the architecture discussion is resolved.
The source loader is an uncommitted, untested draft; integration and migration
are unfinished. No model files, historical exports or boards have been moved.
Next: revise the destination plan for shared root ownership, then finish the
bounded lossless split. Preserve card identities and append-only ledger histories
if boards are later consolidated.

## Owner decision — three systems, shared process

The owner now selects SDL, SDUI and SDPTool as separate Systems documented under
root SDP. This supersedes the initial combined-system recommendation and the
separate language process-tree assumption. XFMD owns its own SDP area; bootstrap
and the next-agent adoption card are delivered separately. System declarations
remain a planned SDL profile extension. CardState is ready for the revised
source/profile migration; no architecture decision remains pending on this point.

## Scrum-0001 disposition

Fully split by KBO-SDP-000003 into [KB-SDL-005](../backlog/%23005--SDL--Change--System-and-source-sets.md)
and [KB-SDP-020](../backlog/%23020--Change--Shared-design-source-organization.md).
The first owns semantics/tool support; the second owns the model/source migration.
No unfinished work is closed as implemented. The uncommitted loader draft remains
input for the first successor, not milestone evidence.
