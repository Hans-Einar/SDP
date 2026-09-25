# K1-M1: establish KanBan and preserve conversation proposals

| Field | Value |
| --- | --- |
| id | KB-SDP-007 |
| CardState | completed |
| ScrumId | SCRUM-SDP-0001 |
| project | SDP |
| type | Change |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed discussion time. The directory and ledger record lifecycle status.

## Assignment and scope

Owner request of 2026-09-23: establish the process/directories, propose tags, and capture recent SDP/SDL/SDUI discussions in primary and reference cards. This is K1-M1; do not implement the other cards' proposals now.

## Implementation plan K1

One phase branch `sdp/phase-k1-kanban` from G7-M1 (`d03eb78`); one milestone:

1. Establish seven status directories in each of three projects, stable IDs, indexes, card template, ledger format and move/closure workflow.
2. Register six primary proposals and four reference cards; identify sources, questions, ownership, next handling and criteria without adopting a new language profile.
3. Check JSON/schema, event histories, locations, identities and new local Markdown links. Record results, complete the milestone, commit and push.

No existing process directory, implementation ledger, Go code, GitHub issue or generated SDL viewpoint is moved/changed. No separate repository creation or merge.

## Acceptance and evidence

**K1-M1 delivered locally on 2026-09-23.** Three boards, seven status directories per board, six primary proposals, four Ref cards and this delivery card were registered. Tags, cadence, ID/reference rules, card template and versioned ledger payload were documented. Root agent instructions and documentation entry points link to the process.

Python 3/jsonschema checks against the existing Toolkit envelope and KanBan payload passed: all JSON files readable, all 11 cards with unique/correct IDs and types, all Refs directly linked to existing primary cards. Per-card event replay checked previous events, before/after status, paths and physical files. All 99 local Markdown links in boards and affected entry points were checked. `git diff --check` passed. Checks also passed after closing this card; 12 events include its active → completed transition. This was a one-off delivery check, not a general KanBan tool.

No Go code, old process directory, Traceability ledger or Toolkit schema changed. No new language profile, sdptool command, graph or repository extraction was implemented; these remain backlog proposals. Commit/push closes the phase in Git; this file's milestone commit is traceable with `git log --follow`.
