# SPR-SDP-007 Iterations and Slices

## ITR-SDP-007-001 — Provisional contract construction

Status: active

### SLC-SDP-007-001 — Pilot contract corpus

Status: active

#### Goal

Create one coherent, explicitly provisional contract corpus that is sufficient
to start a real HSX workflow pilot without modifying HSX or freezing canonical
Toolkit behavior.

#### Required outcomes

- workflow/entity semantics and cardinalities;
- scoped-domain namespace and identity semantics;
- cross-scope, concurrent-assignment and repository-split behavior;
- pilot record shapes and reusable Issue contract;
- proportional profiles;
- deterministic corpus validation; and
- concrete HSX pilot handoff.

#### Verification

- parse every machine-readable example;
- validate all required record kinds, identifiers, relations and assignments;
- exercise positive and representative negative collision/scope cases;
- run existing applicable repository validation; and
- confirm the diff cannot affect canonical Toolkit or downstream repositories.

#### Completion signal

A committed candidate is ready for fresh verification and independent review.

## ITR-SDP-007-002 — Verification, review and rework

Status: planned

### SLC-SDP-007-002 — Exact-candidate challenge

Status: planned

Run fresh deterministic verification and an independent adversarial review at
the same candidate. Resolve all Blocking, High and Medium findings and repeat
the required evidence after material head changes.

## ITR-SDP-007-003 — Pilot-ready closure

Status: planned

### SLC-SDP-007-003 — Handoff and disposition

Status: planned

Integrate final evidence, traceability and the exact draft PR/head, then stop at
`PROVISIONAL_CONTRACT_READY_FOR_PILOT` without opening or changing HSX work.
