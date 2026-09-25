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

## MP1-R-M3 — installation readiness repair

The v1 reference mismatch is a stale conformance authority after PL1's authorized
Planning skill distribution. Generated candidates were written outside the repo
and structurally compared with every existing expected outcome before acceptance.
Exactly 12 of 19 outcomes change, each by three Planning skill entries only:
create/missing-target for empty profiles, unchanged/content-matches for populated
profiles. Removing these entries and normalizing sequence numbers yields exact
structural equality with every previous outcome. Seven blocked/fatal outcomes
remain byte-identical. [Audit](ConformanceAudit.json) records before/after hashes.
No installer behavior or error expectation was weakened to obtain a pass.

The Linux workflow now checks out full history for the upgrade test, which reads
a specific earlier committed five-phase artifact. The test now names its full
40-character commit rather than an abbreviation. It still installs the old
artifact and checks history preservation and current Planning capability on
upgrade. The old source is not replaced by current output.

Local candidate generation and plan/apply conformance: all 19 scenarios pass with
PowerShell /tmp/sk1-pwsh/pwsh. Normal comparison, profile installation, interruption
matrix and remote Windows/Linux gates are recorded below when actually complete.
