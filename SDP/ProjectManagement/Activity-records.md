# Activity records — preserved manual bootstrap contract 0.1

Current implementation: Activity0002 now reads this unchanged contract and writes
versioned0.2 records through SDPTool. See ../../SDPTool/Activities.md and
../Activities/Activity-0002--Activity-Management/06-Verification/Verification.md.
The integrated Python/Go validators support the full mixed history. The historical
bootstrap limitations below describe the earlier manual creation, not the current
qualified candidate; schemas/event bytes remain unchanged.

Owner instruction, 2026-11-12: manually create the first Activity and the same
canonical history that later tooling will produce. This selects a local record
contract before automated support. These records are real, not sample fixtures.
No SDPTool commands, installer capabilities or runtime enforcement are added.

## Identity, paths and schemas

Descriptor: [activity-descriptor-0.1.schema.json](activity-descriptor-0.1.schema.json),
YAML with schema sdp-activity/0.1, identity ACTIVITY-PROJECT-NNNN. The mandate path
is relative to the Activity root; typed reference paths are relative to SDP/.
All paths must resolve to actual records, stay inside their declared root and
reject symlinks/ambiguity for later automated mutation. No state field is allowed.
A prior-discussion reference is context, not retrospective Session membership.

Events use the unchanged generic envelope 1.0 and
[activity-payload-0.1.schema.json](activity-payload-0.1.schema.json). Event IDs are
EVT-ACTIVITY-PROJECT-NNNNNN, subjectId is the Activity identity. This is a new domain
in ProjectManagement/Ledger.ndjson, not kind Activity in management payload 0.2.

Initial supported event semantics:

- x-activity:created: absent to planned, null predecessor, descriptor and need
  already exist; links identify selected records.
- x-activity:transitioned: exact predecessor and previous state; permit planned to
  active, active to on-hold/completed/canceled, on-hold to active/canceled, planned
  to canceled, and completed/canceled to planned for explicit reopening.
- At most one active Activity per canonical project authority after each operation.
  Creation never implicitly activates. Switching and historical adoption need
  additional designed operation contracts; this bootstrap does neither.

occurredAt is the actual effective time, recordedAt is recording time; both are
RFC3339 UTC for this live creation. Do not backdate to prior study work. Actor
identifies the recorder, authority points to the actual owner instruction recorded
in the mandate. operationId is a stable idempotency identity, separate from eventId.
Identical retries must not append duplicates. A changed operation under an existing
identity conflicts. previousEventId chains this Activity's history; event order
is authoritative, not sorting timestamps.

Descriptor path/hash pins the bytes for the event. Future descriptor updates must
retain the referenced version and append the applicable update event; this initial
contract does not authorize untracked descriptor edits. No ongoing state is cached
in the descriptor. Session0017 membership is part of the initial descriptor; future
membership changes require the separately designed linked/update operation.

## Compatibility and manual procedure

Future Activity tooling must read these 0.1 records or perform an explicit,
byte-preserving adoption/version migration. Do not recreate or renumber them, infer
older times, or call them automatically created. Extend unsupported operations
through explicit contract versions rather than silently reinterpreting this data.

For manual creation, resolve current authority and IDs, validate files and references,
inspect existing Activity history, prepare descriptor/documents, then append created
and activated events after confirming the current ledger bytes are unchanged.
Verify schemas, causal replay, descriptor hashes and the single-active rule.
This is a manually checked operation, not proof of crash-atomic multi-file writes,
cooperating concurrent-writer locking or installed automated lifecycle validation.

The existing project management Python validator rejects x-activity events; it also
had unresolved Session links in the latest PLAN0025/0026 entries before this work.
New management update events correct those current projections without changing
historical bytes. Its legacy-domain checks can be run on the filtered legacy event
set, while Activity events receive explicit schema/replay checks. This separation
is a bootstrap verification method, not a full unmodified-validator pass.

The future Activity implementation must integrate domain-aware validation/readers,
including the existing Python validator. Existing SDPTool/PG consumers are not
claimed Activity-aware. No managed capability flag is enabled by this manual pilot.
