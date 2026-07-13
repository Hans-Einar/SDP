---
skillId: sdp-traceability
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.traceability.update
- sdp.traceability.release-events
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Traceability

Use this skill when updating SDP traceability.

1. Read the active work and linked requirements, designs, release target, reviews
   and verification records.
2. Update `CurrentIndex.yaml` to reflect only actual current state.
3. Update `Relations.yaml` with explicit requirement, design, Sprint/Refactor,
   Iteration, Slice/Fix, review, verification and release links.
4. Append immutable events to `Ledger.ndjson`; never rewrite history.
5. Validate YAML and every NDJSON line against applicable schemas.
6. Require real tag, commit and GitHub Release identities for publication events.
7. Report stale, missing or contradictory links to the Master.

Do not invent IDs, evidence, completion state or publication events.
