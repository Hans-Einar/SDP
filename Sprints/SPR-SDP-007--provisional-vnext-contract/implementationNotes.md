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

## 2026-08-21 — Fourth adversarial gate

- Froze full candidate `60978c306f7ce3c09033603abf3300839b2a251d` for the
  reserved fourth evidence pass; exact-head Actions run `32505098158` passed.
- Fresh `VER-SDP-007-004` independently replayed all 157 negative fixtures and
  prior findings, exercised a complete revision-4 history control and returned
  PASS with Blocking 0, High 0, Medium 0, Low 0 and Note 6.
- Separate fresh `REV-SDP-007-004` confirmed all prior findings resolved but
  returned changes required with Blocking 0, High 2, Medium 1, Low 1 and Note
  4. New findings require terminal Feature/Refactor/Fix closure over accepted
  represented Slices/assignments and accepted decision prerequisites; rejection
  of missing schema/kind/experimental downgrade; globally normalized semantic
  edges without circular or cross-surface duplicates; and stable diagnostics
  rather than an encoding exception for surrogate input.

The passing verifier evidence remains exact-candidate evidence; it does not
override the independent review. No reviewer finding is waived.

## 2026-08-21 — Fifth candidate after terminal/marker rework

- A new fresh Worker resolved every `REV-SDP-007-004` finding and the verifier
  satisfiability Note under `SDP-vnext-pilot/` only.
- Added terminal Feature/Refactor/sliced-Fix closure over accepted represented
  Slices/assignments, accepted Study decision and gating dependency rules,
  unconditional exact v0 markers, one globally normalized semantic-edge set,
  self/cross-surface duplicate rejection, stable invalid-Unicode diagnostics,
  and owned/shared/prohibited path disjointness.
- Dogfood assignment revision 4 retains exact revision-1/2/3 snapshots,
  reserves exact REV/VER-005 paths and binds reservation digest
  `sha256:715c4f18e244447411e51bd780807cdef9586a56f21a539cc1aa05fae4cb1e48`.
- Deterministic coverage now contains eight templates, four positive examples,
  six positive mutation controls, two hundred sixteen exact negative fixtures
  and one hundred nineteen diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `48a7ee7f0d0684a96e402caaaf49ed856134dc23`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-005` gate. Earlier evidence remains historical.

## 2026-08-21 — Fifth adversarial gate

- Froze full candidate `15a476dd76bc80de5573ab2abb65af8c963a9c0e` for
  fresh `REV/VER-SDP-007-005`; exact-head Actions run `32509812239` passed.
- Fresh `VER-SDP-007-005` confirmed every earlier finding resolved, then
  returned FAIL with Blocking 0, High 1, Medium 1, Low 0 and Note 5 because
  immutable allocation authority incorrectly blocked later Issue assignments
  and a delivered sliced Fix could omit `affectedWork`.
- Separate fresh `REV-SDP-007-005` confirmed the full prior regression set,
  then returned changes required with Blocking 0, High 2, Medium 2, Low 0 and
  Note 5. It additionally proved sequential accepted Slices at distinct
  candidates cannot be retained under one Issue, terminal standalone-Fix and
  aggregate Issue closure are incomplete, multi-hop `depends_on` cycles pass,
  and non-finite JSON scalars can be hashed.

No fifth-pass finding is waived. The next rework must preserve allocation and
execution identities separately and keep earlier Slice evidence immutable.

## 2026-08-21 — Sixth candidate after multi-Issue lifecycle rework

- A new fresh Worker resolved every `REV/VER-SDP-007-005` finding and the two
  remaining satisfiability Notes within the experimental pilot area.
- Separated immutable prospective `allocationIssue` from later execution
  authorities; added complete `authorizedSlices` history and `0..1`
  `activeSlices`; preserved independently accepted Slice candidates across
  sequential/reopened work; closed every terminal owner/Issue/Fix target;
  enforced relation-type DAGs; rejected non-finite JSON; and required active
  implementation host/write-surface satisfiability.
- Dogfood assignment revision 5 retains exact revision-1 through revision-4
  snapshots, reserves exact REV/VER-006 paths and binds reservation digest
  `sha256:145665468058a57f0278b5c45471d9fbaee180546e20fd3e9189317fa67f8144`.
- Deterministic coverage now contains eight templates, four positive examples,
  fourteen positive mutation controls, two hundred thirty-nine exact negative
  fixtures and one hundred thirty-two diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `cbfa241db51e040f3b3ae8c7480298861de82100`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-006` gate. Earlier evidence remains historical.

## 2026-08-21 — Sixth adversarial gate

- Froze full candidate `dfdea4fe259a9e342d651254656bcfac5363d0b9`
  for fresh `REV/VER-SDP-007-006`; exact-head Actions run `32514457819`
  passed.
- Fresh `VER-SDP-007-006` confirmed all named sixth-candidate lifecycle,
  terminal, DAG, strict-JSON and satisfiability controls, then returned FAIL
  with Blocking 0, High 1, Medium 0, Low 0 and Note 5 because preserved
  historical assignments from different reservation episodes were still
  treated as one concurrent set.
- Separate fresh `REV-SDP-007-006` returned changes required with Blocking 0,
  High 2, Medium 4, Low 1 and Note 6. It additionally proved revision history
  can replace the primary work and erase accepted authorized Slices; reverse
  active-Slice projection is incomplete; an active Fix target can lack any
  accepted represented outcome; top-level relation kinds bypass endpoint-kind
  rules; duplicate JSON keys are not rejected; and invalid branch Unicode can
  raise rather than diagnose.

No sixth-pass finding is waived. The next rework must distinguish historical
reservation episodes from current concurrency and preserve revision semantics
from actual historical assignment/reservation bytes.

## 2026-08-21 — Seventh candidate after reservation-epoch rework

- A new fresh Worker resolved every `REV/VER-SDP-007-006` finding, including
  the Low Unicode robustness item, under the experimental pilot area only.
- Grouped concurrency/base/DAG/order/path checks by reservation epoch while
  preserving accepted historical assignments; rebuilt revision comparison from
  exact historical assignment and reservation bytes with immutable Issue,
  source, workRef and accepted authorized-Slice evidence; closed reverse active
  projection and active Fix target history; enforced one top-level/embedded
  relation kind matrix; rejected duplicate JSON keys recursively; and stabilized
  invalid branch Unicode diagnostics.
- Dogfood assignment revision 6 retains exact revision-1 through revision-5
  snapshots, reserves exact REV/VER-007 paths and binds reservation digest
  `sha256:08cbabb798c0cc95e31052110ff2a910d7ee91e52f370165398c194aea9b552d`.
- Deterministic coverage now contains eight templates, four positive examples,
  twenty-one positive mutation controls, two hundred sixty-three exact negative
  fixtures, two raw duplicate-key parse fixtures and one hundred thirty-eight
  diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `8726e577ac4f6ec7ca1fb81688305d60359e3db0`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-007` gate. Earlier evidence remains historical.

## 2026-08-21 — Seventh adversarial gate

- Froze full candidate `0655002cbe5e14543b007525cdc3e3bad82b6816`
  for fresh `REV/VER-SDP-007-007`; exact-head Actions run `32519142555`
  passed.
- Fresh `VER-SDP-007-007` independently replayed all 263 negatives, 21 positive
  controls, raw parse fixtures and the complete r1-r6 history, then returned
  PASS with Blocking 0, High 0, Medium 0, Low 0 and Note 6.
- Separate fresh `REV-SDP-007-007` confirmed all prior findings resolved but
  returned changes required with Blocking 0, High 1, Medium 2, Low 0 and Note
  5. New findings require exact historical reservation-row equality with the
  historical assignment projection, consistent `(reservation ID,digest)` epoch
  identity throughout resolution/materialization/rebinding, and full structural
  validation of typed unreferenced reservation sets.

The passing verifier record remains exact historical evidence and does not
override the independent review. No reviewer finding is waived.

## 2026-08-21 — Eighth candidate after epoch-identity rework

- A new fresh Worker resolved every `REV-SDP-007-007` finding under the
  experimental pilot area only.
- Historical reconstruction now requires the exact reservation row/base/order/
  convergence projection to equal its exact historical assignment even under
  coordinated byte/digest/pointer rewrites; epoch identity is consistently
  `(reservation ID,digest)` across lookup/grouping/materialization/rebinding;
  and every typed bound or unbound reservation set receives full semantic
  validation.
- Dogfood assignment revision 7 retains exact revision-1 through revision-6
  snapshots, reserves exact REV/VER-008 paths and binds reservation digest
  `sha256:ce963d187746fce12ce14e54ad98b39e2ab13463a5ec8c4f15bd6d1a4d335872`.
- Deterministic coverage now contains eight templates, four positive examples,
  twenty-five positive controls, two hundred eighty exact negative fixtures,
  two raw parse fixtures and one hundred forty diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `36205b14450d59c64a1c071a2d0890eb5fc620b6`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-008` gate. Earlier evidence remains historical.

## 2026-08-21 — Eighth adversarial gate

- Froze full candidate `133cfaee9cb194b6181ac1e8fa9e1b474f79c00f`
  for fresh `REV/VER-SDP-007-008`; exact-head Actions run `32522761761`
  passed.
- Fresh `VER-SDP-007-008` independently replayed all 280 negatives, 25
  positives, raw fixtures and r1-r7 history, then returned PASS with Blocking
  0, High 0, Medium 0, Low 0 and Note 6.
- Separate fresh `REV-SDP-007-008` confirmed prior findings resolved but
  returned changes required with Blocking 0, High 2, Medium 1, Low 0 and Note
  6. New findings require historical source-candidate ancestry/order proof,
  repository/domain/inventory/current-collision coherence for unbound
  preparatory sets, and exact revision-scoped early compatibility rules.

The passing verifier record remains exact historical evidence and does not
override review. No reviewer finding is waived.

## 2026-08-21 — Ninth candidate after ancestry/preparatory rework

- A new fresh Worker resolved every `REV-SDP-007-008` finding under the
  experimental pilot area only.
- History source candidates are now strict chronological ancestors of the
  validated HEAD; unbound preparatory sets receive repository/domain/record/
  inventory/current-collision coherence checks; and early compatibility
  normalization is exact and revision-scoped to r1 and r2-r4 only.
- Dogfood assignment revision 8 retains exact revision-1 through revision-7
  snapshots, reserves exact REV/VER-009 paths and binds reservation digest
  `sha256:888daed562aaafb05530e48f5ebb11b344d7005a25ce3890d3e83db301d6ea70`.
- Deterministic coverage now includes eight templates, four examples, twenty-six
  positive mutations, two hundred eighty-eight standard negatives, two raw
  parse fixtures, eight compatibility controls and six disposable-Git ancestry
  controls spanning one hundred forty-three diagnostics.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `f7f31a626430749a1527021705834fa3c8c984b0`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-009` gate. Earlier evidence remains historical.

## 2026-08-21 — Ninth adversarial gate

- Froze full candidate `a6a23d50e1bb807c9db350758aad1737c01f7846`
  for fresh `REV/VER-SDP-007-009`; exact-head Actions run `32526591987`
  passed.
- Fresh `VER-SDP-007-009` independently replayed the full stack, r1-r8
  ancestry/history, preparatory controls and neighboring probes, then returned
  PASS with Blocking 0, High 0, Medium 0, Low 0 and Note 6.
- Separate fresh `REV-SDP-007-009` confirmed all earlier findings resolved but
  returned changes required with Blocking 0, High 1, Medium 1, Low 0 and Note
  6. A reserved ID may still be absent from structured inventory, permitting
  terminal-to-preparatory reuse, and a preparatory authorized Slice may claim a
  non-null accepted candidate without represented acceptance evidence.

The passing verifier remains exact historical evidence; it does not override
the independent review. Neither finding is waived.

## 2026-08-21 — Tenth candidate after durable-reservation rework

- A new fresh Worker resolved both `REV-SDP-007-009` findings under the
  experimental pilot area only.
- Added a durable pre-materialization `reserved` inventory state, mandatory
  inventory/allocator resolution for every bound or preparatory reservation,
  cross-epoch no-reuse, promotion to materialized `prospective`, and null-only
  preparatory authorized-Slice candidates until bound acceptance evidence
  exists.
- Dogfood assignment revision 9 retains exact revision-1 through revision-8
  snapshots, reserves exact REV/VER-010 paths and binds reservation digest
  `sha256:6bdb62d86c07aad489fde8be141eb90a83df3ea56de2b43d59e59b2a1273a30e`;
  it correctly does not re-add already allocated `STU-007` to reserved IDs.
- Deterministic coverage now contains eight templates, four examples,
  twenty-eight positive controls, two hundred ninety-nine negatives, two raw
  parse fixtures, compatibility and temporary-Git controls spanning one
  hundred forty-eight diagnostic codes.
- Master validation passed the pilot validator, Issue #5 analysis validator,
  Toolkit validator and all nine Toolkit Python tests.
- Committed and pushed technical candidate
  `f7eefcb20d1a3f051550e9caeca8672b032c0775`.

The exact remote head containing this recovery update now requires the reserved
fresh `REV/VER-SDP-007-010` gate. Earlier evidence remains historical.
