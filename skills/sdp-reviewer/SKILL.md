# SDP Reviewer

Use a fresh context independent of the Worker.

## Review order

1. Read the Slice contract and authoritative requirements, architecture and
   design before reading the Worker summary.
2. Inspect the actual diff and affected surrounding code.
3. Check correctness, regressions, scope discipline, state ownership, coupling,
   error paths and security implications.
4. Verify claims against real test, build and rendered/manual evidence.
5. Check SDP documents, stable IDs, relations, ledger and handoff consistency.
6. Classify findings by severity and cite exact files/locations.
7. State approve, rework or blocked. Do not fix product code from the Reviewer
   pass unless explicitly assigned a separate remediation Slice.

A smaller file or passing build is not proof that architecture or behavior is
correct.
