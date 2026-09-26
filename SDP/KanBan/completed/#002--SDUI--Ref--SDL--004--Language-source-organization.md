# KB-SDUI-002 — Split the shared design into phase-owned sources

| Field | Value |
| --- | --- |
| id | KB-SDUI-002 |
| project | SDUI |
| type | Ref |
| created | 2026-09-25T10:32:43Z |
| source | owner-conversation-2026-09-25 |
| CardState | completed |
| ScrumId | SCRUM-SDP-0001 |
| primary | KB-SDL-004 |

## Local impact

Primary: [KB-SDL-004](../superseded/%23004--SDL--Change--Language-source-organization.md). Adopt the five phase homes and split SDUI/design/architecture.design. Keep SDUI-specific internals under this project and shared integration in one canonical SDL-owned source set. Retain dated exports as historical evidence.

## Worklog

- 2026-09-25: in-progress; coordinated migration selected with the primary card.

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

## Scrum-0001 routing outcome

Local consequence is recorded in [KB-SDP-020](../backlog/%23020--Change--Shared-design-source-organization.md). The former separate-board
routing is no longer needed. This closes only this Ref, not language, tool or
model implementation. Historical identity and original primary remain recorded.
