# SPR-SDP-007 Iterations and Slices

## ITR-SDP-007-001 — Provisional contract construction

Status: complete

### SLC-SDP-007-001 — Pilot contract corpus

Status: complete

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

Completed candidate:
`f79e3dfc18c7a1650f1f3ae66167dda9b69692b4`.

## ITR-SDP-007-002 — Verification, review and rework

Status: active

### SLC-SDP-007-002 — Exact-candidate challenge

Status: active

Run fresh deterministic verification and an independent adversarial review at
the same candidate. Resolve all Blocking, High and Medium findings and repeat
the required evidence after material head changes.

The first exact-candidate pass at `f79e3df...` returned rework:

- `VER-SDP-007-001`: FAIL, Blocking 0 / High 1 / Medium 3 / Low 0;
- `REV-SDP-007-001`: changes required, Blocking 0 / High 2 / Medium 3 /
  Low 0 / Note 4.

The active work is to resolve every High and Medium finding, then use fresh
verification and a new independent Reviewer at the reworked candidate.

## ITR-SDP-007-003 — Pilot-ready closure

Status: planned

### SLC-SDP-007-003 — Handoff and disposition

Status: planned

Integrate final evidence, traceability and the exact draft PR/head, then stop at
`PROVISIONAL_CONTRACT_READY_FOR_PILOT` without opening or changing HSX work.
