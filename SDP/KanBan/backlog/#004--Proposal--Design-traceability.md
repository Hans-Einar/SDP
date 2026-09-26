# Traceability between SDL design, slices, code and evidence

| Field | Value |
| --- | --- |
| id | KB-SDP-004 |
| CardState | backlog |
| ScrumId | SCRUM-SDP-0001 |
| Systems | SDL, SDUI, SDPTOOL |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

## Need

Design and implementation must actively use the Traceability ledger. SDP tools should link SDL design and the approved implementation plan to actual deliveries, showing current status and system growth across roadmap milestones. KanBan history does not replace this evidence.

## Contract to investigate

Inventory and reuse/extend existing envelopes, IDs and relation contracts; do not create a competing implementation ledger. Link stable model-object identity and model revision to slice/milestone identity, code/commit references, verification results and evidence. Resolve invalidation when design/evidence changes, partial delivery, missing coverage and cross-project dependencies.

Keep proposed, planned, implemented and verified distinct. A passing test or one delivered Functionality does not complete a whole Feature. Expose unknown/stale evidence; do not infer status from prose or weak `links`. Reports must show aggregation rules and scope.

## Next work and acceptance

Compare current Ledger.ndjson, Current-Index and Relations with these needs. Define a versioned contract and migration before making changes. Test one model with two slices, only one implemented/verified; generated roadmap/progress diagrams should show the difference and link to sources/evidence. [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md) is a consumer.

[Traceability ledger](../../Traceability/Ledger.ndjson)
[Ledger envelope schema](../../../Toolkit/schemas/ledger-event.schema.json)
[Traceability relations](../../Traceability/Relations.yaml)

## Backlog review — 2026-09-25

Retained as the independent evidence/status contract. KB-SDP-017 consumes it for advanced plan/roadmap reporting; ordinary navigation need not wait for evidence aggregation.

Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDP-000076. CardState remains backlog.

## Scrum-0001 review

Keep model-to-code evidence and aggregation separate from management lifecycle. Include system-prefixed new IDs and management-event references; never infer implementation from sprint/card completion.

EVT-KB-SDP-000117; next review at the next selection or relevant dependency delivery.

## Owner clarification — SDL code tags and paired model verification

2026-09-26T10:12:02Z: EVT-KB-SDP-000183. The owner intends SDL tags in implementation code to bind
model identities to actual implementation locations and support code/model checks.
Design missing/stale/duplicate/ambiguous tag diagnostics and exact model/code
revision references. A tag is a declared mapping, not proof of behavior. Assignment
bundles must support NOW/baseline tests before work and TARGET/result tests after
work, retaining unchanged-neighbor obligations. Coordinate PLAN-SDP-0001 and
KB-SDL-006; no tag grammar or executable mapping is adopted by this clarification.
