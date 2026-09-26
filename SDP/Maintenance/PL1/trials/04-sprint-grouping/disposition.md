# Proposed membership changes
- Create the planned SPR-VIEW-0001 record with a separate Plans metadata row.
- Add SprintId SPR-VIEW-0001 to the one metadata table in each existing plan and append actual membership-update events, preserving state planned.
- Retain KB-VIEW-003's completed state, original closure evidence and successor links. No wrapper card, new planning Study, source-card reopening or duplicate plan is necessary.
- Append the real transitions to ProjectManagement/Ledger.ndjson, not Traceability. Never insert plan IDs into the legacy members array.

This is a proposed record set in an isolated evaluation. No repository or live ledger transitions were performed.
