---
skillId: sdp-master
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.coordinate
- sdp.release.coordinate
- sdp.traceability.coordinate
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Master

Use this skill when coordinating repository work under SDP.

## Procedure

1. Read `AGENTS.md`, `AGENTS-project.md` when present, Framework guidance, the
   project manifest, current traceability and the active work contract.
2. Identify the active Sprint or Refactor, Iteration, Slice or bounded Fix, plus
   the current and target release versions.
3. Confirm goal, scope, invariants, non-goals, expected files, verification and
   completion signal before implementation.
4. Refine contradictory SDP documents before delegating work.
5. Delegate one bounded implementation unit to a Worker.
6. Delegate verification inspection and independent review to fresh contexts.
7. Inspect actual evidence; do not accept summaries as proof.
8. Update implementation notes, release notes, relations, current index and ledger
   only for transitions that actually occurred.
9. Decide complete, rework or blocked, then stop at the Slice/Fix boundary.

## Guardrails

- Repository records override chat memory.
- SemVer identifies releases; Sprint/Iteration/Slice coordinates do not alter it.
- Record discovered work instead of implementing it opportunistically.
- Never claim a tag, GitHub Release or verification record before it exists.
- Never continue automatically into the next Slice.
