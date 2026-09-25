# MAINT-SDP-0001 — shared board and project management

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0001 |
| project | SDP |
| state | active |
| source | SCRUM-SDP-0001; owner-conversation-2026-09-25 |

Base b2fb900; branch sdp/phase-pm1-project-management. Owner: Project Owner;
executor: Codex. Selected directly by [Scrum-0001](../../Agents/Scrum/Scrum--%230001--Backlog-and-project-management.md),
without an extra KanBan wrapper. This is project-management maintenance, not
SDL System implementation or the outstanding model migration.

| Milestone | Acceptance | State |
| --- | --- | --- |
| PM1-M1 | Consolidate boards and histories, review all cards, preserve identities/lineage, define optional grouping and one management ledger with validation | Delivered |
| PM1-M2 | Extend shell listing for Sprint/Scrum grouping/filtering, verify workflow/negative cases/links and hand off next work | Planned |

## Invariants and migration

One writable ledger under ProjectManagement. Import original histories exactly,
keep archive hashes and byte ranges, continue per-card chains. Preserve all IDs;
namespace foreign card filenames to avoid collisions. Six routing Refs close only
after their consequences have a primary home. One genuine split retains full
source history. No arbitrary merging just to reduce card count.

No product Sprint is invented for this maintenance. Optional tags are tested with
fictional fixtures. System-prefixed new Traceability IDs and management references
do not rewrite historical bytes or claim legacy Toolkit conformance. Existing
SDL/go/sourceinput drafts stay untouched and outside commits. XFMD uses its pinned
profile; only producer reference addresses may be refreshed there.

## Verification to record

Import prefix/hash/inventory, unique IDs, schema/replay/placement, current metadata
and grouping, relationship/lineage completeness and invalid transitions; all local
links, historical evidence and CLI behavior including terminal/plain output and
installer safety. No GUI/runtime acceptance is implied by these process checks.

## Outcome and remaining work

Evidence and completion are recorded after executed checks. System/input work and
shared design-source migration belong to the two successor cards. Reusable release
contracts, templates and native project navigation remain separately tracked.

## PM1-M1 result

All 25 original cards retain identities; two successor cards bring the board to
27. Six routing Refs are complete, the mixed source card has a typed full split,
and eleven independent open cards have explicit Scrum dispositions. No Sprint
was invented. Three original ledgers (113 events) form an exact preserved prefix,
with frozen hash-pinned archives. Management and card transitions share one writer.

Checks passed: schema/replay/current placement and metadata, two lineage operations,
full/partial lineage examples and 15 negative cases; four management test groups
cover the three permitted origin paths, review/refactor work, 12 corrupt histories,
sprint membership and actual-start activation, and exact import bytes. The existing
105 frozen records/ledger prefixes and 574 generated artifacts remain intact;
local documentation links/anchors pass. PM1-M2 adds CLI grouping next.
