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

## 2026-08-21 — Second adversarial gate

- Frozen full reworked/recovery candidate
  `a1c64f4dadb663e3b3b7c03c1f8759b94c7f2be3` for a new fresh evidence pass.
- Exact-candidate GitHub Actions run `32496386287` passed its contracts and
  installer jobs.
- Fresh `VER-SDP-007-002` independently confirmed every earlier verifier and
  reviewer regression resolved, then returned FAIL with Blocking 0, High 0,
  Medium 2, Low 0 and Note 3 for malformed/empty authority/recovery/owner/
  cardinality values plus self-conflict and duplicate graph edges.
- Fresh `REV-SDP-007-002` independently confirmed the first-pass resolution and
  returned changes required with Blocking 0, High 2, Medium 5, Low 0 and Note
  4. It additionally required canonical GitHub Issue/repository/PR identity,
  cross-registry collision safety, minimum lifecycle/value/Fix rules, portable
  identity source/provenance and truthful no-Issue legacy support, qualified
  exact-candidate accepted evidence, Study authority consistency, and an
  incremented dogfood assignment revision with a durable refreeze chain.

No second-pass finding is waived. `ITR-SDP-007-002` / `SLC-SDP-007-002`
remain active for the final required rework and new exact-candidate evidence.

## 2026-08-21 — Final rework candidate

- A new fresh Worker resolved every second-pass implementation finding within
  `SDP-vnext-pilot/` and left canonical Toolkit/downstream/release surfaces
  unchanged.
- Added canonical GitHub repository/Issue/comment/PR parsing and repository
  coherence; one-registry-per-repository collision safety; supported pilot
  state/minimum value and zero-Slice Fix rules; portable record-source binding;
  truthful nullable legacy authority; qualified exact-candidate acceptance
  evidence; Study/issued authority consistency; self/duplicate/overlapping edge
  rejection; and dogfood assignment revision/refreeze history.
- The real dogfood assignment is now revision 2, links the revision-1 snapshot
  at `f79e3df...`, reserves exact REV/VER-003 paths and binds current reservation
  digest `sha256:f5805cecd8f72394d4502eed8dbb4f1ec899f890faed3dd5f4e368e1a4f448bb`.
- Deterministic coverage now contains eight templates, three positive scenarios,
  ninety-seven exact negative fixtures and seventy-nine diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical rework candidate
  `730424f5d6078356353da57c8234ef4237b4eae8`.

The full branch head containing this recovery update must now receive new fresh
verification and an independent `REV/VER-SDP-007-003` pass. Earlier failed
dispositions remain historical only.

## 2026-08-21 — Third adversarial gate

- Froze full candidate `0de8a8957b3212404007160d90da47445d7b4e7b` for
  fresh `REV/VER-SDP-007-003` evidence; exact-head Actions run `32500882771`
  passed.
- Fresh `VER-SDP-007-003` independently confirmed all prior findings resolved,
  then returned FAIL with Blocking 0, High 0, Medium 1, Low 1 and Note 4 for a
  reusable revision-2+ history omission/fabrication gap and same-row duplicate
  ID reservation.
- Fresh `REV-SDP-007-003` independently confirmed all prior findings resolved,
  then returned changes required with Blocking 0, High 2, Medium 3, Low 2 and
  Note 4. The new findings require accepted assignment/Slice/Study/prerequisite/
  conflict state coherence; qualified embedded Feature/Refactor/Fix/Study
  semantic edges; complete multi-refreeze history with generic local
  reconstruction; complete portable path/shared/boundary validation; unique
  prospective record sources; valid move/key hints; and exact pilot schema
  markers.

No third-pass finding is waived. `ITR-SDP-007-002` / `SLC-SDP-007-002`
remain active for bounded rework and a new exact-candidate gate.

## 2026-08-21 — Fourth candidate after state/history/graph rework

- A new fresh Worker resolved every `REV/VER-SDP-007-003` finding under the
  experimental pilot area, including the two Low adjacency findings.
- Added cross-record assignment/Slice/Study/Fix/prerequisite/conflict state
  coherence; qualified embedded Feature/Refactor/Study/Fix semantic edges;
  complete generic repository-aware multi-refreeze history; conservative
  portable path/shared/boundary validation; normalized prospective source
  uniqueness/materialization; move registry/key-hint checks; exact schema
  markers; and same-row reservation uniqueness.
- Dogfood assignment revision 3 preserves revision-1 at `f79e3df...`, revision-2
  at `0de8a895...`, reserves exact REV/VER-004 paths and binds reservation
  digest `sha256:c6b9245e236c8843ed4fc76fa8c0e540f4afadb7bb75ee33dfc1f3ff0502623d`.
- Deterministic coverage now contains eight templates, four positive examples,
  one positive mutation control, one hundred fifty-seven exact negative
  fixtures and one hundred three diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `390541cfc480e920d15aa15366f6078e8c020835`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-004` gate. No earlier disposition is reused.
