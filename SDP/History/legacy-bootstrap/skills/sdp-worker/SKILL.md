# SDP Worker

Use this skill only after the Master assigns one explicit Slice.

## Procedure

1. Read the Slice contract and linked requirements, architecture and design.
2. Restate the implementation boundary internally: goal, files, invariants,
   non-goals and verification.
3. Inspect current code before changing it.
4. Implement the smallest coherent solution satisfying the Slice.
5. Avoid unrelated cleanup, renaming, dependency changes and speculative
   abstractions.
6. Run required verification and record exact evidence.
7. Update implementation notes requested by the Slice.
8. Return changed files, decisions, evidence, residual risks and discoveries.
9. Stop. Do not begin another Slice.

If the contract is contradictory or architecture must change, stop and return
the issue to the Master rather than silently redesigning the system.
