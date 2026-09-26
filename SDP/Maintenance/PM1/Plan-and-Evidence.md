# MAINT-SDP-0001 — shared board and project management

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0001 |
| project | SDP |
| state | completed |
| source | SCRUM-SDP-0001; owner-conversation-2026-09-25 |

Base b2fb900; branch sdp/phase-pm1-project-management. Owner: Project Owner;
executor: Codex. Selected directly by [Scrum-0001](../../Agents/Scrum/Scrum--%230001--Backlog-and-project-management.md),
without an extra KanBan wrapper. This is project-management maintenance, not
SDL System implementation or the outstanding model migration.

| Milestone | Acceptance | State |
| --- | --- | --- |
| PM1-M1 | Consolidate boards and histories, review all cards, preserve identities/lineage, define optional grouping and one management ledger with validation | Delivered |
| PM1-M2 | Extend shell listing for Sprint/Scrum grouping/filtering, verify workflow/negative cases/links and hand off next work | Delivered |

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

## PM1-M2 result and completion

The installed kanban command supports --group-by state|sprint|scrum, --sprint ID
and --scrum ID. Filters match exact identities and combine with AND; unassigned
cards remain visible under (none). State-only behavior, bounded directory scope,
OSC 8 terminal links and plain redirected output remain intact. Optional tags do
not create/move cards or start a Sprint. Every reviewed card now references
SCRUM-SDP-0001; there is no actual Sprint to pretend has started.

Verification: 10 CLI tests passed, including grouping/filter combinations,
missing matches, malformed/duplicate metadata, first-table scope, hostile filename
encoding, installer backup/idempotence and symlink protection. Four management
test groups and the lineage examples/15 negatives still pass. Installed into
/home/warloc/bin; the previous command was preserved by the installer backup.
Current ledger: 27 cards, two management records, two lineage operations, 164
combined events after closure. Eleven cards remain open, six routing Refs are
complete, and the split's two successor scopes remain visible.

Frozen evidence checks retain 105 records/prefixes and 574 generated outputs;
all current local links/anchors pass. The old three board ledgers remain exact
113-event archive/import bytes. No new management-only event was added to
Traceability: its single new record covers the actual project-tooling code change
and verification, referring to this maintenance and its completion event.

No SDL System/source-set behavior, model migration, native graph/sidebar, released
KanBan package or installer-template migration is claimed. KB-SDL-005 is queued
next; KB-SDP-020 depends on it; #010 remains on owner review. XFMD keeps its pinned
profile, with producer-card links updated and the future contract impact recorded
in its existing #012 Ref. No XFMD product code is changed.
