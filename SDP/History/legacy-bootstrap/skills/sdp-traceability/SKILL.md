# SDP Traceability

Use this skill when creating or updating stable project links.

## Procedure

1. Read the project's ID conventions and current index.
2. Update `CurrentIndex.yaml` to reflect only genuinely active work.
3. Update `Relations.yaml` so requirements, designs, Tiers, Sprints, Iterations,
   Slices, reviews and verification evidence remain connected.
4. Append events to `Ledger.ndjson`; never rewrite historical ledger entries to
   hide a correction.
5. Validate YAML and each NDJSON line with a parser.
6. Check referenced paths and IDs exist and statuses agree across documents.
7. Record corrections as new ledger events.

Traceability is evidence navigation, not ceremonial bookkeeping. Avoid links
that do not help explain why work exists, what it changes or how it was verified.
