# Manual Activity bootstrap evidence

Recorded 2026-11-12 for ACTIVITY-SDP-0001 / SESSION-SDP-0017. Records only;
no SDPTool functionality or product verification is implied.

## Checked outcomes

- Exactly one Activity directory and identity; all six phase folders and root
  Sessions exist. README is generic directory guidance; need/scope are in Mandate.
- Descriptor validates against activity-descriptor-0.1.schema.json, contains no
  lifecycle state, and all typed references/mandate resolve to actual files.
- Every ledger envelope validates against the existing generic schema. Both new
  Activity payloads validate against activity-payload-0.1.schema.json.
- Explicit replay checks created → planned → active, predecessor/state agreement,
  distinct event/operation identities, timestamps, descriptor hashes and at most
  one active Activity. Final replay: ACTIVITY-SDP-0001 is active.
- Created: 2026-11-12T18:18:37.309814+00:00 (`EVT-ACTIVITY-SDP-000001`).
- Activated: 2026-11-12T18:18:37.310129+00:00 (`EVT-ACTIVITY-SDP-000002`).
- Descriptor SHA-256: `437dd70088035e41c7243a9c97219a592bc00c90137d93c36a3f1a3798f97543`.
- Prior ledger prefix: 513530 bytes, SHA-256
  `07b2d3571f897ce6840c31fa9fb6666ac8e2354f53f8aedb46004b9747822664`; unchanged byte-for-byte. HEAD ledger prefix also
  remains intact. Seven current events appended: three card updates, two plan
  updates and two Activity lifecycle events.
- Existing validator's validate_current on the explicitly filtered non-Activity
  event set passes: 76 cards, 50 management records, four lineage operations.
  Current plan links were corrected by append-only updates; board index now uses
  lifecycle placement while KB060's own CardState stays queued.

## Limitations and reproducibility

The unmodified whole-ledger validator reports `unknown event type` for x-activity.
It is not claimed to pass. Activity schemas and the semantic assertions above were
checked with Python/jsonschema/PyYAML during manual registration; this is not an
installed automatic Activity validator. The current registry and descriptor can
be rechecked against the published schemas, and the two events replayed in their
recorded order with the stated invariants.

No runtime discovery, XFMD Activities tab, automatic scope gate, concurrent-writer
or crash-recovery guarantee was implemented. This manual setup adds no Go code,
SDL changes or generated views. Legacy sessions stay at their original paths.
Activity 2 and 3 do not yet exist; the directory index records intentions only.
M0 inventory remains next and PLAN-SDP-0026 remains planned.
