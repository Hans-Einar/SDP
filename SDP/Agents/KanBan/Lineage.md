# Merging and splitting KanBan cards

Contract 0.2, clarified 2026-09-24. This records conversation cards and work
provenance, not automatic requirements satisfaction or implementation status.

## Workflow

When selecting work for active, review related cards, later clarifications and
conflicting proposals. New knowledge can change the next bounded delivery.
Merging does not approve every source idea: distinguish owner decisions,
recommendations, open questions and documented observations.

| Operation | Sources | Result |
| --- | --- | --- |
| merge | At least two existing primary cards | One new consolidated card with a new ID |
| split | One existing primary card | At least two new cards, each with a new ID |

**Fictional** example: A + B → C, then C → D + E. This is a graph with convergence
and branching, not merely status changes along one document's history. Small
clarifications can update the same card with a reviewed event; do not invent
merge/split operations. Represent many-to-many changes as two explicit operations
through an intermediate consolidated card.

1. Define what transfers from each source, what is superseded and why, and what
   remains unresolved. Preserve useful facts, dates and evidence.
2. Create target cards with new IDs. Use active only for selected, authorized
   scope; otherwise use backlog. Split targets can have different statuses.
3. Add a linked source table to each target, and outcomes/successors to sources.
   Metadata may contain `sources` and `superseded_by`; typed ledger lineage is
   authoritative for the provenance graph. `source` remains conversation provenance.
4. Move fully replaced sources to superseded. Link **all** split targets. Keep
   partially transferred sources open with explicit remaining work, or create
   another target for that work before closing the whole source.
5. Record created for targets and moved/reviewed for sources as specified below.
   Update indexes and current Markdown links together with the files.
6. Check that the full transition is represented and no questions were lost.
   This requires substantive review: schema validation cannot prove completeness.

Never delete sources. Preserve dated statements and identify decisions that
supersede them. Superseded means further handling has another home, not that a
feature is delivered. A Ref is not a lineage source: use its primary card and
update local impact separately. Do not redirect historical Ref links to an
arbitrary target after a split.

## Ledger payload 0.2

[Schema 0.2](ledger-payload-0.2.schema.json) retains the fields/event types from
[0.1](ledger-payload.schema.json) and adds optional `lineage`. The envelope remains
1.0. Preserve old 0.1 events unchanged. New events use payload 0.2; ordinary
created/moved/reviewed events do not require lineage.

Every participant event in a merge/split has an **identical** lineage description:

```json
{
  "operationId": "KBO-DEMO-000001",
  "kind": "merge",
  "sources": [
    {"id": "KB-DEMO-001", "scope": "Initial idea", "remaining": ""},
    {"id": "KB-DEMO-002", "scope": "Later clarification", "remaining": ""}
  ],
  "targets": [
    {"id": "KB-DEMO-003", "scope": "Consolidated, bounded work"}
  ]
}
```

`operationId` has a sequential namespace per coordinating project:
`KBO-<PROJECT>-<number>`, with at least six digits. It is neither a card ID nor an
event ID. `scope` describes transferred scope for sources and received scope for
targets. `remaining` describes work retained in a source; an empty string means
full replacement. Explain conflicting/discarded content in the source table and
event reason. Split targets divide source scope; do not accidentally create two
parallel primary assignments for the same work. Schema checks structure; review
checks allocation and coverage of remaining work.

### Events and checks

- Sources exist before the operation. Targets have new, unique IDs. Source and
  target sets are disjoint, with no repeated IDs. New lineage therefore adds no cycle.
- Each source and target has exactly one participant event per operationId.
  A missing event means an **incomplete operation**, not a completed merge/split.
- Target created events link to all source IDs. Sources link to all target IDs.
  Keep these ordinary links for simple readers; lineage defines the relation type.
- A fully replaced source moves to superseded. A partial source has nonempty
  remaining, receives reviewed without status/path changes, and remains in
  backlog, active or onHold.
- Full replacement also starts from backlog/active/onHold. Cite closed historical
  cards without superseding them again. If work reopens, first record reopening
  with its own rationale.
- Append target created events before source events in the same ledger. Existing
  per-card previousEventId, path and status rules still apply.
- Never reuse an operationId for a different definition or extra participants.
  Details must match across every participant event, including different boards.
- Update incoming links to each source's new location. Sources then lead to their
  successors; do not erase historical identities.

Within one board, commit files, indexes and events together. Across boards,
cards keep project IDs and each board records its events with the same operationId
and definition. Coordinate and validate all changes together. Separate repositories
have no shared atomic commit. Until all participant events are available, a combined
reader must show an incomplete operation. An inaccessible board is unknown evidence,
not an empty or completed project. Cross-repository transactions, global locking
and automatic ID lookup are not implemented.

## Examples and verification limits

The [example ledger](examples/lineage.ndjson) contains only fictional DEMO cards:
full/partial merge and split, including sources with remaining work. It is outside
production ledgers and is not registered in boards.json. The
[check](examples/verify_lineage.py) tests schemas, semantic histories and negative
cases such as missing participants, inconsistent operations and lost links.
It tests example contracts; it is not a general KanBan CLI, graph or validator of
multi-repository card content/location.

Run from the repository root with Python 3 and jsonschema:

```sh
python3 SDP/Agents/KanBan/examples/verify_lineage.py
```
