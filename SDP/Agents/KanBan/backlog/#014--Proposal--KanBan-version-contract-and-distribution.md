# Define the KanBan version contract and reusable distribution

| Field | Value |
| --- | --- |
| id | KB-SDP-014 |
| project | SDP |
| CardState | backlog |
| type | Proposal |
| created | 2026-09-24T17:40:43Z |
| source | Owner conversation, 2026-09-24 Europe/Oslo: KanBan version contracts, distribution paths and XFMD synchronization |
| next_review | Before the next KanBan contract/template change or KanBan extraction; coordinate with KB-SDP-001 and the XFMD consumer |
| tags | process, kanban, compatibility, templates |

## Need and owner direction

XFMD borrowed the conversation-board workflow while SDP was being restructured.
Its card format must remain able to support functionality defined by the shared
KanBan system. The owner requests an explicit backlog item for proper version
contracts and future organization, plus a linked XFMD synchronization item.

The owner identifies these eventual locations relative to the SDP-vNow root:

- `KanBan/`: the shared KanBan definition and its supporting material.
- `Template/sdp-root/Agents/KanBan/`: templates distributed to consuming projects.
- `SDP/Agents/KanBan/`: the current SDP project's own operational board, where
  this request is registered. Distinguish real project cards/history from reusable
  definitions/templates; determine the exact migration boundary when selected.

This records intended direction. It does not move directories, introduce a
released contract version or authorize implementation of the migration.

## Current evidence and compatibility gap

Inspection on 2026-09-24 used SDP-vNow branch
`sdp/phase-l1-english-documentation`, commit
`cbb0dc6c26e8b99f8b01626669958fff500779cf`.
The [workflow](../README.md) already defines visible metadata, stable identities,
Refs and lifecycle; [Lineage](../Lineage.md) defines typed merge/split semantics;
[History](../History.md) and the [template](../Card-template.md) define active
worklogs and Git-backed revisions/diffs. Ledger payload schemas 0.1 and 0.2 and
board metadata 0.1 exist locally. These schema labels are not a complete released
KanBan compatibility contract for cards, workflows, tools and distributed templates.
The two proposed distribution directories did not exist at inspection time.

XFMD currently uses its own visible-metadata template and payload 0.1. A consumer
needs an identifiable supported baseline and upgrade rules, not just periodic
copying of whichever README happens to be newest.

## Proposal, scope and open questions

Agent recommendations for the selected work:

- Define what is versioned together or separately: card metadata/body conventions,
  board metadata, ledger payloads, worklogs, Ref resolution, lineage and templates.
  Choose where a consumer declares its supported version/profile and how tools
  handle older or unknown versions. No field names or version scheme are chosen here.
- Publish capability/compatibility rules, release/change notes, validation criteria
  and an upgrade path. Preserve existing IDs, historical event bytes and Git history;
  distinguish migrating current cards from interpreting older records.
- Define which features are required versus optional for a lightweight standalone
  consumer such as XFMD. Shared KanBan compatibility must not require adopting the
  rest of SDP or creating an SDP directory in XFMD.
- Plan extraction into `KanBan/` and distribution through the requested Template
  location. Identify canonical ownership and verify templates against it to avoid
  independent competing definitions. Determine installer/update responsibilities.
- Document how links, board registration and older consumer references survive
  relocation. Decide manual versus automated conformance/update checks separately.

## Related-card review

[KB-SDP-001](../completed/%23001--Proposal--Project-structure.md) retains broader
project/template structure and extraction work; this card owns the bounded KanBan
contract and distribution need. [KB-SDP-010](../active/%23010--Proposal--Document-consolidation.md)
retains editorial consolidation. Completed K2/K3/K4 cards remain evidence for
metadata, lineage and history; they do not prove a released compatibility contract.
[KB-SDP-003](%23003--Idea--KanBan-graph.md) retains the future graph, and
[KB-SDP-002](%23002--Proposal--sdptool.md) retains tooling implementation.
No existing card is superseded, moved or reduced by this registration.

Consumer follow-up: [KB-XFMD-012](../../../../../xfmd-sdl-navigation/Agents/KanBan/backlog/%23012--Ref--SDP--014--KanBan-format-sync.md) is the direct XFMD Ref. It owns the local gap
review and later adaptation, without duplicating the shared contract decision.

## Next action and completion criteria

At the next contract/template or extraction planning round, select a bounded
contract/delivery scope and coordinate with XFMD before changing the published
format. Completion requires an agreed version/compatibility policy, explicit
canonical/distributed ownership and paths, migration guidance and documented
checks against representative existing and upgraded boards. Planned relocation
must be delivered or explicitly dispositioned, not silently omitted on closure.

## Outcome and references

Registered in backlog; no version contract, migration or consumer update delivered.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-24T17:40:43Z | Codex; EVT-KB-SDP-000057 | Captured owner direction, inspected current local schemas/history and reviewed overlap with structure, consolidation and tooling cards. | Shared version contract and distribution remain to be selected; KB-XFMD-012 tracks XFMD adaptation. |

## Concurrent registration integration

The owner authorized inclusion with K5 after the other session finished. The
uncommitted event originally numbered EVT-KB-SDP-000041 was reassigned to
EVT-KB-SDP-000057 and appended after committed events. Registration time, card
identity, content and the XFMD Ref are preserved; no committed ledger bytes
were reordered. The local CardState addition does not deliver this version contract.
