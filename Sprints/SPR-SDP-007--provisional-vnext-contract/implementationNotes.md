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
- Opened early draft PR #8 from the dedicated delivery branch to `main` and
  bound the assignment record to its stable URL.

The root `Traceability/Ledger.ndjson` remains release-only and is intentionally
unchanged; Issue #7 work transitions must not be misrepresented as release
events before a general event contract exists.

## 2026-08-21 — Worker candidate and first adversarial gate

- A fresh Worker produced the complete provisional corpus under
  `SDP-vnext-pilot/`, including workflow/scope/concurrency contracts, seven
  typed templates, three positive examples, seventeen negative fixtures,
  deterministic validation and the concrete HSX pilot handoff.
- Master rework removed legacy Sprint/Iteration/Slice IDs from vNext
  reservations, minted durable domain UID
  `urn:uuid:9df0765b-4f22-4c5e-8fda-b5e71c96ca81`, added the resolvable
  `STU-007` dogfood record and bound the assignment to canonical reservation
  digest `sha256:50d6768eecce05e540c7674b34a78ba6b4432cabae6e079b88fb51e669bebd4d`.
- Committed and pushed immutable candidate
  `f79e3dfc18c7a1650f1f3ae66167dda9b69692b4`.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests. Exact-candidate GitHub
  Actions run `32492184311` passed.
- Fresh `VER-SDP-007-001` returned FAIL at that candidate with Blocking 0,
  High 1, Medium 3 and Low 0 after finding validator false negatives for
  concurrency cycles/empty convergence, duplicate Issue authority, domain-root
  overlap and NFKC-introduced forbidden path characters.
- Fresh `REV-SDP-007-001` returned changes required with Blocking 0, High 2,
  Medium 3, Low 0 and Note 4. In addition to the verifier findings, it required
  complete Slice/assignment/owner/path cardinality, reusable digest-bound
  concurrent reservation records, evidence-qualified accepted state,
  structured legacy/issued/reserved identity no-reuse, and current recovery
  coordinates.
- Completed `ITR-SDP-007-001` / `SLC-SDP-007-001` and activated
  `ITR-SDP-007-002` / `SLC-SDP-007-002` for required rework and new evidence.

No finding is waived. The candidate is not pilot-ready.

## 2026-08-21 — High/Medium rework candidate

- A fresh rework Worker resolved every implementation finding from
  `VER-SDP-007-001` and `REV-SDP-007-001` under `SDP-vnext-pilot/` only.
- Added bidirectional Slice/owner/Issue/assignment/path authorization checks;
  a reusable embedded-or-external canonical reservation-set contract with
  unique Issues, dependency DAG, complete ordered merge and terminal
  convergence; common evidence-qualified accepted state; structured
  prospective/legacy identity inventory with provenance and no-reuse; recursive
  domain-root collision rules; and NFKC-before-semantics path validation.
- Repaired the scoped example with three real assignments and one digest-bound
  convergence set; repaired the move example to avoid invented delivery state
  while preserving a provenance-backed historical `DBG-RF-001` identity.
- Expanded deterministic validation to eight templates, three positive
  scenarios, fifty-three exact negative fixtures and forty-eight diagnostic
  codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed reworked candidate
  `b0dd81cabf88a7e7dec977e9d58db2001d82f937`.

The candidate now requires fresh verification and a new independent review.
No earlier review or verification disposition is reused.
