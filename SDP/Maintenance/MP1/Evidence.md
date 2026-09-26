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

## Gates remaining after MP1-R-M2 (historical)

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

## MP1-V-M1 — candidate checks and independent review

Implementation candidate: 0ee89e1ce8695752a93e2cca4c99d6020b0de430.
A detached worktree excludes the unrelated sourceinput draft. With Go 1.27.1,
`go test -race ./...` and `go vet ./...` pass in SDL/go, SDUI/go and
Toolkit/SDPTool. The optional real-renderer tests were additionally run with
MMDR/SDUI_MMDR pointing to the existing mermaid-rs-renderer target/debug/mmdr:
SDL documents (including semantic/class backend checks) and SDUI markdown pass.
Maintained parser/runtime tests are Go; old Python parser implementations have
been removed. No new desktop GUI acceptance or feature readiness is claimed.

Exact-checkout Toolkit discovery passes 106 tests (19 host tests deferred to the
separate host runs); six management test groups pass. The document check passes
105 frozen records/prefixes, 574 generated outputs, 2,467 local links and 130
fragments when checked out beside the existing XFMD repository. Its first /tmp
checkout correctly failed on an existing external sibling link; the repository
layout requirement was restored without changing or suppressing link checks.

Independent reviewer /root/mp1_review reviewed 0803470 through 0ee89e1 and approved
these bounded changes, with one low documentation finding: the Traceability
README still described KB-SDP-011 as unresolved. That current assertion is now
corrected. The reviewer independently confirmed exact SPS-001 restoration,
preserved ledger prefixes, ID and relation rejection controls, all conformance
audit hashes/deltas, and a passing real PowerShell old-profile upgrade (107 s).
This is not a new review of the entire earlier development stack. The reviewer
explicitly retained remote Windows/Linux and full interruption-matrix gates.

Remote check run 36197129935 was started for 0ee89e1; no result is claimed here.
The next documentation-only candidate must also receive green remote checks
before its SHA is selected for integration. Publication and version release
remain outside scope.

## MP1-R-M3 / MP1-V-M1 — completed remote gates

Selected candidate: 0cf22458002809c58748f65354edbadcc70a8279; tree
f19bb3f99ee01eec618d03f6c640142a2c0ae924. Independent reviewer /root/mp1_review
also approved its documentation-only delta after the low finding was corrected.
[Run 36197403423](https://github.com/Hans-Einar/SDP/actions/runs/36197403423)
completed all three jobs successfully for that exact head:

- contracts: full Toolkit validator and 106 tests, followed by v1 package checks;
- Windows installer: PowerShell fixture suite and all 19 v1 reference scenarios;
- Linux process profile: two build tests, all 19 installation/upgrade/consumer
  tests, and all 69 journal steps at backup/write/journal boundaries plus
  prepare/completion exits, with forward resume, preserved history and
  exactly-once finalization.

The separate local PowerShell integration run completed all 19 tests in 521.539 s.
The exact final documentation candidate passed management, document and Toolkit
checks in a detached checkout. The Go implementation is unchanged from 0ee89e1;
its passing race/vet and real backend checks retain that candidate identity.
All 574 generated outputs and their design-source revision bindings match
(538 viewpoints, 14 navigation outputs and 22 runtime-preview outputs).

## MP1-I-M1 — staging integration

Owner resumed the MergePlan on 2026-09-26 after the authentication interruption.
[PR #35](https://github.com/Hans-Einar/SDP/pull/35) was merged at
2026-09-26T08:00:43Z using the exact-head guard and a merge commit:
8039700198cecf8dbd27d6394773955b3849178e. Its parents are the original staging
head 9ad432407004080dd7f4f0ab06d107523f4316fd and verified candidate
0cf22458002809c58748f65354edbadcc70a8279. Its tree is exactly
f19bb3f99ee01eec618d03f6c640142a2c0ae924. No conflict resolution or code changes
were needed. The 34 Issue #7-specific commits remain outside this ancestry.

GitHub automatically closed the included staging PRs; PR #4 remains open against
main, and pilot PR #8 remains separate. [PR #36](https://github.com/Hans-Einar/SDP/pull/36)
now proposes this exact staging result to main. Its own checks are pending at this
recording point. The old archive tag remains unchanged; main is not yet advanced.

## MP1-I-M2 — main integration

[Run 36228514036](https://github.com/Hans-Einar/SDP/actions/runs/36228514036)
passed contracts, Windows installer/conformance and Linux process-profile tests
on exact PR #36 head 8039700198cecf8dbd27d6394773955b3849178e. The Windows job
took 12m21s and Linux 12m7s; the complete fault matrix passed. Before merging,
remote main still equaled the original baseline and merge-tree equaled the
verified staging tree. No checks were bypassed.

PR #36 merged at 2026-09-26T08:13:57Z as
2f0256f7e822ef2965e8f3d287095f1d90eac46c. Its parents are former main
2cb49c02145621b099c47d05786716598e414e75 and staging
8039700198cecf8dbd27d6394773955b3849178e. Its tree is
f19bb3f99ee01eec618d03f6c640142a2c0ae924, exactly the verified candidate.
All selected milestone commits are ancestors. The 34 Issue #7-specific commits
remain outside main; draft PR #8 and its branch remain untouched. GitHub marked
included PR #4 merged at 08:13:59Z; PRs #12–#34 were already reconciled by staging.

Remote annotated old remains bf420cceb773210757376d0813415bf09ee8bdba,
peeling to 2cb49c02145621b099c47d05786716598e414e75. No phase branch was deleted
or rewritten. The local main worktree was not reset or changed; this observation
is about remote main. The untracked SDL/go/sourceinput draft remains excluded.

## MP1-C-M1 — administrative closeout boundary

Current AGENTS and development entry points now target main and preserve earlier
staging history as dated records. KB-SDP-031/032/033 remain backlog studies.
No Toolkit version, published release or live XFMD installation is claimed.
The earlier staging administration commit c9f64a7 was independently approved
with no findings; its evidence and KB-SDP-030 closure join this follow-up delivery.
The closeout is documentation/history only and is delivered through a separate
reviewed PR to main. Its GitHub merge record will own the final merge identity;
the already observed main promotion above remains the stable evidence baseline.
