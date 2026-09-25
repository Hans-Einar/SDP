# Shared project management

Profile sdp-project-management/0.1. Ledger.ndjson is the sole writable lifecycle
history. Keep documents current; Git stores revisions; the ledger stores actual
transitions. Never rewrite historical bytes or infer product completion from a
management state. Event envelope 1.0 and management payload 0.1 apply.

Management IDs are MAINT-<PROJECT>-<number>, SCRUM-<PROJECT>-<number>,
SPR-<PROJECT>-<number>, REVIEW-<PROJECT>-<number> and REFACTOR-<PROJECT>-<number>.
Use x-management:created/started/updated/completed/canceled. A created item is
planned or active; started changes planned to active; completed requires real
completion. Chain previousEventId with exact from/to states and SDP-relative
document paths. Installation events use EVT-PM-INSTALL-<operation>-<step>.

Cards retain KB-<PROJECT>-<number>, x-kanban events and payload 0.1/0.2.
Card paths are KanBan-relative. Optional SprintId/ScrumId metadata groups work.
Never duplicate lifecycle events in Traceability.
