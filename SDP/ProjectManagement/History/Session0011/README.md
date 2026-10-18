# Session0011 integration history reconciliation

Owner authorization, 2026-10-18: resolve the competing branches, integrate to main,
then complete Actions/JSON Session0011. This import makes previously branch-local
identities unambiguous without changing the original records in their source commits.

The canonical ledger preserves the complete published release-lineage ledger as a
byte-identical prefix. Exact source ledgers from release 3475367, blueprint f64a5ba
and PG de4b795 are archived here. import.json records full commits, SHA-256 hashes,
source line numbers and original-to-canonical event and subject mappings.

- Published Runnable SDUI programs keeps KB-SDP-051; Actions becomes KB-SDP-053.
- Published RSP2 keeps MAINT-SDP-0015; PLR1 role guidance becomes MAINT-SDP-0016.
- Identical records imported by several branches are deduplicated by full content,
  never solely by event ID. Distinct records receive new monotonic event IDs.
- Imported projections use the import timestamp and current canonical predecessor;
  the original event timestamp and complete payload remain in the archive. Reviewed
  events cannot undo a later active state. The import reason names its source.
- PG owns the current active KB038/052 records and Session0007; the stale backlog
  copy of KB052 is removed. Session0008 remains paused, with owner pilot review open.
- Traceability combines unique original events, rejects differing records sharing
  an ID, and preserves original bytes. No implementation evidence is invented here.

This is a one-time branch-history import, not permission to edit canonical history
on normal updates. Historical evidence references still identify archived source
records; current action/role documents use reconciled IDs. Original signed releases
and their evidence are unchanged. These archives are not writable authorities.

Run verify.py here plus the standard ProjectManagement validator after integration.
