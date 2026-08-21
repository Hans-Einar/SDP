# SPR-SDP-007 Implementation Notes

## 2026-08-21 — Issue activation

- Read the complete Issue #7 body and confirmed it has no comments at activation.
- Read the complete Issue #5 body, all four comments and Steering acceptance.
- Fast-forwarded local `main` to the merged study baseline
  `2cb49c02145621b099c47d05786716598e414e75`.
- Read the required study synthesis, migration, legacy, roles/follow-up and HSX
  evidence documents in `SDP-usage-analysis/`.
- Created delivery branch `codex/issue-7-provisional-vnext-pilot` from that exact
  baseline.
- Activated `SPR-SDP-007` / `ITR-SDP-007-001` / `SLC-SDP-007-001`.
- Established the experimental assignment at
  `SDP-vnext-pilot/Steering/Assignments/ISSUE-007.yaml` with common base,
  owned/shared paths, reserved pilot identities, stale-base checks and explicit
  downstream/PR #4 boundaries.

The root `Traceability/Ledger.ndjson` remains release-only and is intentionally
unchanged; Issue #7 work transitions must not be misrepresented as release
events before a general event contract exists.
