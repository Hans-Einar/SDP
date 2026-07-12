# Traceability

Project traceability normally contains:

- `CurrentIndex.yaml` — actual current release and development coordinates
- `Relations.yaml` — links among requirements, design, work, review, verification,
  migrations and releases
- `Ledger.ndjson` — append-only event history, one valid JSON object per line

Release events conform to `Toolkit/schemas/release-event.schema.json`. Only append
state transitions that occurred. In particular, do not record
`release-tag-created` or `release-published` before the real tag or GitHub Release
exists. Corrections are new events; historical lines are never rewritten.

Release IDs use `REL-X.Y.Z`. Fix IDs use `FIX-X.Y.Z-NNN` or, for emergency work,
`HOTFIX-X.Y.Z-NNN`.
