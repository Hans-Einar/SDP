# Shared project management — profile 0.2

Owner decision: 2026-09-25. This is the adopted `sdp-project-management/0.2`
profile, not a new published version of the SDP Toolkit or KanBan distribution.

## One workflow, two histories with different purposes

[Ledger.ndjson](Ledger.ndjson) is the single writable **project-management**
history for KanBan cards, Scrum discussions, Sprints, Maintenance, CodeReview
and Refactor work. Current documents describe scope/outcome; Git owns document
revisions; the ledger records actual transitions and relationships. No parallel
writable ledger remains in KanBan or the language process areas.

[Traceability](../Traceability/README.md) records changes to system designs/code
and their verification, with references to management IDs/events. Completing a
card, sprint or review never proves a system Feature is implemented or verified.
Pure board cleanup/meeting/planning events do not belong in Traceability. Existing
historical management events there remain unchanged and are labeled pre-cutover;
no new duplicate transition is appended to both ledgers.

## Scrum and work selection

A Scrum here means the owner's local backlog-review/consolidation session; it is
not a claim to implement an external Scrum framework. One concise record contains
source IDs, decisions, dispositions, selected output and remaining scope. Record
only actual participants/inputs, not invented meeting attendance or approval.
[Scrum-0001](../Agents/Scrum/Scrum--%230001--Backlog-and-project-management.md)
is the first record; its outcome is MAINT-SDP-0001, not a forced Sprint.

Supported paths include cards → Scrum → Maintenance, Scrum → Maintenance, and
card → Maintenance. A card may be executed alone; a Scrum may select a Sprint,
Maintenance, CodeReview or Refactor. None requires an extra wrapper card. Every
selected management work item has its own stable ID, scope, completion criteria,
source/reference IDs and events in this ledger. Record unknown/unselected work
honestly; a review recommendation alone does not authorize it.

## Optional sprint grouping

A Sprint groups selected plans and/or direct cards under a goal; see [Plans](Plans.md).
The paragraphs below retain the direct-card path and historical membership rules. Use one concise record
in Sprints, listing membership/goal and exceptions; do not copy the cards there.
Before start, cards remain backlog (possibly queued). On actual Sprint start,
move the selected backlog cards to active/ready together, appending a transition
per card plus the Sprint start event in one commit. An already active card can
be associated explicitly without resetting its working state. Do not move
completed cards back to active merely to fill a sprint list.

Cards carry optional `SprintId` and `ScrumId` visible metadata. Current membership
must match the Sprint record. Reviews/reassignments append events and preserve
previous membership in history; current tags are not a history database. Existing
sprints are historical records until explicitly mapped, never relabeled retroactively.
Use CardState ready/in-progress/gate-review inside active; queued remains backlog.
A Sprint may contain multiple system areas. Cancellation/removal returns unfinished
cards to a declared backlog/onHold/active destination with reasons; it does not
silently finish or lose them. Closing a Sprint requires each member's disposition.

## IDs and ownership

| Record | New local ID convention |
| --- | --- |
| Scrum | SCRUM-SDP-0001; file Scrum--#0001--Topic.md |
| Sprint | SPR-SDP-0001; distinguish from legacy SPR-SDP-005 |
| Maintenance | MAINT-SDP-0001 |
| Standalone typed plan | PLAN-SDP-0001 |
| CodeReview | REVIEW-SDP-0001 |
| Refactor | REFACTOR-SDP-0001 |
| Management event | EVT-PM-SDP-000001 |
| KanBan card/event | Existing KB-SDP/SDL/SDUI namespaces retained; do not renumber |

SDP is the managing **project** namespace, not an enclosing software system.
Optional `Systems` metadata identifies SDL, SDUI, SDPTOOL or project-wide SDP
scope. New system Traceability identities put the system first, for example
SDL-DES-0001, SDUI-REQ-0001 and SDPTOOL-VER-0001. A cross-system delivery links
separate system records to one management item. Shared project tooling uses SDP
as an explicit project-wide namespace. Old IDs keep their meaning; KB-SDP-011
owns legacy schema compatibility. Do not claim every existing Toolkit schema
already accepts this local new-ID policy.

## Ledger and migration contract

Use the existing generic envelope 1.0. KanBan events retain x-kanban event types
and payload 0.1/0.2 interpretation; new card events use 0.2 and per-namespace
EVT-KB counters. Their paths are relative to SDP/KanBan. New management records
use x-management:created/updated/started/completed/canceled and
[payload 0.1](management-payload.schema.json) for existing untyped work or
[payload 0.2](management-payload-0.2.schema.json) for typed plans/Sprint plan groups. Management
paths are relative to SDP; each Sprint event includes its current `members` ID
array. previousEventId chains each subject, with exact
from/to states and paths. Use meaningful links to source/result record IDs.

A created management record starts planned or active; started moves planned to
active; completed/canceled records a real terminal outcome. updated keeps its
current state/path and explains scope/membership/decision changes. There is no
implicit reopening rule in 0.1; select a new work record for further work and link
the predecessor. CodeReview/Refactor cards can use those types and point to their
actual work records; don't duplicate the review report in the card.

[Import manifest](History/import.json) records byte ranges, original paths and
hashes for the three frozen original ledgers. Their complete bytes form the
initial ledger prefix in SDP, SDL, SDUI order; cross-board timestamps are not
rewritten into global chronological order. Sequence is authoritative per subject;
timestamps remain actual observations. Filenames include the legacy namespace
where needed to avoid collisions. Archive ledgers are read-only migration evidence.
Old Traceability events are not copied into management history.

Validate with `python3 SDP/ProjectManagement/validate.py` and its negative tests.
Schema validation alone cannot prove that the written work was actually done.
Do not infer activity from mtime, count duplicate archive events or run new writers
against the historical roots. The generated-state graph/history UI is still backlog.

## Typed plans (0.2)

[Plans](Plans.md) is the current planning authority. New typed work uses management
payload 0.2: kind Maintenance with planType MaintenancePlan, or kind Plan with
one of the six planType values. Plan IDs use PLAN-SDP-<four or more digits>.
Existing untyped records and all payload 0.1 event bytes remain supported.
Typed plan identity/type cannot change during its lifecycle. Sprint payload 0.2
retains members for direct cards and adds plans for plan IDs; both snapshots
are required; creation/start requires selected work. Later explicit removals may
leave an empty Sprint for truthful closure. Document Plans and each plan's
SprintId must agree with history. No new ledger or CardState is introduced.
