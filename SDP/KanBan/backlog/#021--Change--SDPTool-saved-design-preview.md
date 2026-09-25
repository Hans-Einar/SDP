# Saved-file design preview

| Field | Value |
| --- | --- |
| id | KB-SDP-021 |
| project | SDP |
| type | Change |
| CardState | queued |
| ScrumId | SCRUM-SDP-0002 |
| SprintId | SPR-SDP-0001 |
| Systems | SDPTOOL |
| created | 2026-09-25T12:38:14Z |
| source | KB-SDP-017; owner-conversation-2026-09-25 |
| next_review | Sprint start, then each owning phase milestone |

## Delivery and acceptance — P0-M1

Deliver a small Go sdptool operation for a saved .design file using existing SDL
Go parsing/projection/document services and the existing renderer boundary.
Define command arguments, supported profiles and default relevant views before
implementation. Standalone files work without an installed SDP project. Return
Markdown, Mermaid/SVG resources, manifest/provenance and actionable diagnostics
through one documented result contract; report unsupported views explicitly.

Acceptance: a supported model produces source-derived diagrams; invalid syntax,
missing inputs/tools, unsupported profiles and output failures produce useful
errors. Test paths with spaces, resource ownership/cleanup and unchanged source
bytes. Show a current-source edit changes generated output. Do not generate all
views by default or compile tools while opening a document. Native XFMD preview
and unsaved buffers are not included. Reuse the pipeline guide; do not create a
second parser or rewrite rendering in Rust.

## Source and scope transfer

| Source | Transferred scope | Retained scope |
| --- | --- | --- |
| [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md) | P0-M1: Saved-file design preview | P0-M2, T4-M2 and T5 remain in the source; sibling cards own the other selected milestones |

Typed partial split: KBO-SDP-000004. The [single feature plan](../../05--Implementation/SDPTool.md)
owns phase/milestone acceptance; the [Sprint](../../Sprints/Sprint--%230001--SDPTool-preview-and-navigation.md) owns
membership and completion. No implementation is delivered by this registration.

## Queue

Proposed first execution after completed MAINT-SDP-0001 and Scrum-0002. A visible
saved-file preview gives an early useful outcome without blocking on project
metadata or System/source-set language work. The Sprint is planned, not started.

## Worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000124: registered with acceptance and dependency boundaries. Await Sprint start.
