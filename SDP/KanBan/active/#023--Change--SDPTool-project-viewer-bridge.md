# Project resolver and configured viewer bridge

| Field | Value |
| --- | --- |
| id | KB-SDP-023 |
| project | SDP |
| type | Change |
| CardState | ready |
| ScrumId | SCRUM-SDP-0002 |
| SprintId | SPR-SDP-0001 |
| Systems | SDPTOOL |
| created | 2026-09-25T12:38:14Z |
| source | KB-SDP-017; owner-conversation-2026-09-25 |
| next_review | Sprint start, then each owning phase milestone |

## Delivery and acceptance — T2-M1/M2

Implement the T1 contract in the Go SDPTool CLI at Toolkit/SDPTool. Support
an explicit project path before the command and current-directory selection when
omitted. Resolve either the selected valid SDP area or its SDP child. Implement
view ip with the project's plan/main page and navigator via configured prebuilt
tools; preserve existing direct language CLIs.

Depends on KB-SDP-021/022. Acceptance: test repo root, SDP area, explicit path
from unrelated cwd, missing/invalid/unsupported metadata, nested projects,
optional capabilities, unavailable executables and paths with spaces. Verify
viewer arguments with a deterministic consumer harness; distinguish this from
real GUI evidence. Preserve authored plans/sources, report process errors, and
assign temporary bundle cleanup and request/window targeting explicitly. No
required daemon, startup build or hardcoded SDUI/repository default. A missing
approved plan is a diagnostic, not permission to invent or overwrite one.

## Source and scope transfer

| Source | Transferred scope | Retained scope |
| --- | --- | --- |
| [KB-SDP-017](%23017--Proposal--sdptool-and-project-navigation.md) | T2-M1/M2: Project resolver and configured viewer bridge | P0-M2, T4-M2 and T5 remain in the source; sibling cards own the other selected milestones |

Typed partial split: KBO-SDP-000004. The [single feature plan](../../05--Implementation/SDPTool.md)
owns phase/milestone acceptance; the [Sprint](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) owns
membership and completion. No implementation is delivered by this registration.

## Worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000126: registered with acceptance and dependency boundaries. Await Sprint start.

2026-09-25T13:08:53Z — EVT-KB-SDP-000136: Owner selects Sprint execution before the skills and installer Scrums; backlog to active/ready.
