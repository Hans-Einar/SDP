# Model-derived navigation and selected generation

| Field | Value |
| --- | --- |
| id | KB-SDP-024 |
| project | SDP |
| type | Change |
| CardState | completed |
| ScrumId | SCRUM-SDP-0002 |
| SprintId | SPR-SDP-0001 |
| Systems | SDPTOOL |
| created | 2026-09-25T12:38:14Z |
| source | KB-SDP-017; owner-conversation-2026-09-25 |
| next_review | Sprint start, then each owning phase milestone |

## Delivery and acceptance — T3-M1/M2

Expose versioned tree data backed by existing SDL catalog/model facts: phase,
viewpoint, typed collection, object and supported relationship selections. Cover
all groupable kinds the selected profile supports, not only UseCase. Keep stable
IDs distinct from display labels and paths; mark absent/unsupported/empty views.

Depends on KB-SDP-021/022/023. Acceptance: fixtures cover non-UseCase collections,
variable depth, shared objects, cyclic relationships and bounded expansion.
Define source revision identity, refresh after structural edits and stale-result
rejection. Generate only overview/inventory initially and selected detail on
demand; full static export is explicit. Both Markdown navigation and a native
consumer use the same semantic selection facts and sdl-view:// boundary.

Do not infer facts from generated Markdown or duplicate SDL projection logic.
Single-file supported models suffice for this sprint; advertise future System,
source-set, Stakeholder/UserStory features only after their owning SDL delivery.
Record any small required SDL-owned API extension in SDL Traceability.

## Source and scope transfer

| Source | Transferred scope | Retained scope |
| --- | --- | --- |
| [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md) | T3-M1/M2: Model-derived navigation and selected generation | P0-M2, T4-M2 and T5 remain in the source; sibling cards own the other selected milestones |

Typed partial split: KBO-SDP-000004. The [single feature plan](../../05--Implementation/SDPTool.md)
owns phase/milestone acceptance; the [Sprint](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) owns
membership and completion. No implementation is delivered by this registration.

## Worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000127: registered with acceptance and dependency boundaries. Await Sprint start.

2026-09-25T13:08:53Z — EVT-KB-SDP-000137: Owner selects Sprint execution before the skills and installer Scrums; backlog to active/ready.

2026-09-25T13:27:42Z — EVT-KB-SDP-000146: T3-M1: expose SDL-owned catalog, typed collections and object identity as a versioned navigation graph.

2026-09-25T13:30:23Z — EVT-KB-SDP-000147: T3-M1/M2 delivered: all-catalog typed model navigation, finite relationship references and current-source selected generation with revision checks.
