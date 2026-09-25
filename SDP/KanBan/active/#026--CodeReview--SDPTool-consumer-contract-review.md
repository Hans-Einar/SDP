# Consumer contract and end-to-end delivery review

| Field | Value |
| --- | --- |
| id | KB-SDP-026 |
| project | SDP |
| type | CodeReview |
| CardState | ready |
| ScrumId | SCRUM-SDP-0002 |
| SprintId | SPR-SDP-0001 |
| Systems | SDPTOOL |
| created | 2026-09-25T12:38:14Z |
| source | KB-SDP-017; owner-conversation-2026-09-25 |
| next_review | Sprint start, then each owning phase milestone |

## Delivery and acceptance — T4-M1

Review the combined producer against the owner's journeys: preview a saved
.design file, discover a project, open its plan, browse all supported concept
groups, select fresh-source diagrams and browse KanBan/SDUI entries. Depends on
KB-SDP-021 through KB-SDP-025; review is substantive, not a count of closed cards.

Publish executable examples and versioned fixtures plus a deterministic consumer
harness for XFMD. Test current-source regeneration, unsupported capabilities,
errors, cancellation/request ordering, stale results and temporary resource
lifetime. Record tested tool versions and distinguish stubbed consumer checks
from native GUI interoperability. Confirm repeated runs preserve authored sources.

At execution create a REVIEW-SDP record linked to this card/Sprint; record actual
review findings, fixes and remaining limitations with system verification evidence.
Do not invent independent review participation or owner approval. Handoff the
producer contract to XFMD KB-XFMD-014/015; XFMD code changes remain there. T4-M2
real consumer verification is retained in KB-SDP-017 until that consumer exists.
Close only when the sprint acceptance is met or explicit unfinished dispositions
are recorded; review completion alone cannot mark missing code delivered.

## Source and scope transfer

| Source | Transferred scope | Retained scope |
| --- | --- | --- |
| [KB-SDP-017](%23017--Proposal--sdptool-and-project-navigation.md) | T4-M1: Consumer contract and end-to-end delivery review | P0-M2, T4-M2 and T5 remain in the source; sibling cards own the other selected milestones |

Typed partial split: KBO-SDP-000004. The [single feature plan](../../05--Implementation/SDPTool.md)
owns phase/milestone acceptance; the [Sprint](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) owns
membership and completion. No implementation is delivered by this registration.

## Worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000129: registered with acceptance and dependency boundaries. Await Sprint start.

2026-09-25T13:08:54Z — EVT-KB-SDP-000139: Owner selects Sprint execution before the skills and installer Scrums; backlog to active/ready.
