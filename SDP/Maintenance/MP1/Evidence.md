# MP1 integration evidence

## MP1-R-M2 — historical Traceability compatibility

Candidate: 1419bb0 plus the MP1-R-M2 diff. All 38 baseline validator findings
are resolved by preserving legitimate historical work IDs, accepting the adopted
system-prefixed verification IDs in the verification category, and repairing a
misplaced YAML block. The nine SDPTool verification objects had been inserted
between SPS-001's verification key and its values, incorrectly moving its review
ownership to SDPTOOL-VER-P0-M1. They now live under verification; SPS-001 again
owns its original six verification and seven review links. No historical ID,
append-only ledger prefix, source evidence or review disposition is rewritten.

CurrentIndex and Relations accept the same historical SPR/ITR/SLC spellings.
These are bounded identifier forms, not adoption of Issue #7 workflow semantics.
Malformed IDs, verification-as-slice, dangling references and missing reverse
links remain invalid. The regression tests exercise these distinctions against
the actual repository graph and both schema surfaces.

- Full Toolkit validation against origin/main: pass, previously 38 findings.
- Toolkit unittest discovery: 106 tests, pass; 19 optional PowerShell integration
  tests skipped because that host was not supplied in this run. Those tests are
  required separately for MP1-R-M3/V-M1 and are not counted as passed here.
- Dedicated compatibility regression: two test groups pass.

## Remaining gates

MP1-R-M3 installation conformance, exact-candidate verification, independent
review and remote CI remain outstanding. Local compatibility checks do not
establish main integration readiness.
