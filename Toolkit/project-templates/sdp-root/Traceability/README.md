# Traceability

Project traceability contains:

- `CurrentIndex.yaml` — actual current release and development coordinates
- `Relations.yaml` — stable links among requirements, design, work, review,
  verification, migrations and releases
- `Ledger.ndjson` — append-only events, one valid JSON object per non-empty line

An empty Ledger is valid. Canonical events use the generic Ledger envelope;
`release-*` events also use the stricter release-event contract. Project-defined
event types use `x-<namespace>:<event-name>` and still satisfy the generic
envelope. Record only transitions and evidence that actually exist.
