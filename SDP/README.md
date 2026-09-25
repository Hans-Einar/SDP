# Developing SDP

This area gathers development work for the SDP product itself.

## Current process structure

Owner application of KB-SDP-001, 2026-09-25: this project now uses the five
numbered phase directories. SDPTool is the first product feature designed through
them; code remains in Toolkit/SDPTool, with no nested SDP project.

| Phase | Responsibility / abstraction | First feature artifact |
| --- | --- | --- |
| 01--Mandate | Owner assignment and stakeholders; process input, not a new A-level | [Mandate](01--Mandate/Mandate.md) |
| 02--Requirements | Needs, goals and functional intent, A0/A1 | [SDPTool requirements](02--Requirements/SDPTool.md) |
| 03--Architecture | Containers and internal ownership, A2/A3 | [SDPTool architecture](03--Architecture/SDPTool.md) and [SDL model](03--Architecture/SDPTool.design) |
| 04--Design | Detailed behavior and contracts, A4 | [SDPTool contract work](04--Design/SDPTool.md) |
| 05--Implementation | Phases, milestones, realization and evidence, A5 | [SDPTool implementation plan](05--Implementation/SDPTool.md) |

Studies belong within the phase that needs them. One model can support viewpoints
at multiple levels; a folder does not assign every model fact the same level.
The shared records below remain project-wide. Distributed Template/installer
profiles and automatic project discovery have not been migrated by this local setup.

## Shared project records

- [KanBan](KanBan/README.md): conversation findings, prioritization and active assignments.
- [Card worklogs and Git diff](KanBan/History.md): content history during active work.
- [R1 housekeeping](Maintenance/R1/Plan.md): migration map, document inventory and verification.
- [R2 document consolidation](Maintenance/R2/Plan.md): active contracts, dated proposals and status corrections.
- [Traceability](Traceability/README.md): existing implementation ledger and relations.
- [Sprints](Sprints/README.md), [verification](Verification/README.md) and [review](CodeReview/README.md).
- [Usage study](Studies/UsageAnalysis/README.md): dated Issue #5 material.
- [Historical bootstrap](History/legacy-bootstrap/README.md): retired root copies, not installation sources.

Templates live in Template; language implementations have their own directories.
KanBan tracks proposal handling, not implementation evidence. KB-SDP-001 is the
accepted structure decision; KB-SDP-010 retains its separate editorial review.

The remaining structure/consolidation package is in [R3](Maintenance/R3/Plan.md).
Its five-phase profile is now applied locally by the owner's instruction; current
installed templates retain their documented authority until a tested migration.
CardState distinguishes review from ongoing work.

## Shared management — PM1

[One KanBan board](KanBan/README.md) contains SDP, SDL and SDUI cards.
[ProjectManagement](ProjectManagement/README.md) owns their history together with
Scrums, Sprints, Maintenance, CodeReview and Refactor. [Scrum-0001](Agents/Scrum/Scrum--%230001--Backlog-and-project-management.md)
selects [MAINT-SDP-0001](Maintenance/PM1/Plan-and-Evidence.md), without a new wrapper
card or forced Sprint. Traceability records system changes and links to that history.
