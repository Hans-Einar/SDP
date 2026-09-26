# MAINT-SDP-0004 evidence

## PL1-A-M1 — planning contract

Baseline 28bf156 plus the PL1-A document diff. The owner authorized both planning
and execution. One branch with phase commits was selected explicitly; historical
large-plan commitments remain unchanged. The contract uses existing Maintenance
identity, optional Sprint grouping, and one authoritative plan text. Payload and
profile 0.2 are selected here; their executable support is PL1-B work.

Root instructions, card workflow and project-management guidance now point to
that contract. KB029 is active/in-progress because its plan executes, not merely
because a plan was drafted. Frozen ledgers are retained byte-for-byte. No live
consumer is changed by this phase. The next phase must deliver schema/readers,
skill and distribution before claiming the new profile usable.

## PL1-B-M1 — implementation and distribution

Candidate: c2d8edf plus the typed planning implementation diff. Canonical skill
count is fourteen; sdp-planning 1.0.0 and router 1.1.0 are distributed with their
portable references. AGENTS contract 2.1.0 adds per-plan policy. Toolkit remains
unreleased 0.2.0. The explicit process artifact selects management profile 0.2,
retains layout 0.1, and adds sdp.planning.v1. Old schemas/history remain readable;
no existing consumer is upgraded automatically.

Local and consumer validators check typed-plan identity/type, state, declared Git
policy and Sprint membership. Payload 0.2 separates plan and direct-card arrays.
Explicit removal can leave an empty Sprint for closure without finishing deferred
plans; creation/start still require selected work. Tests cover six plan types,
legacy Maintenance adoption, immutable types, missing references/policies,
reciprocal tags and valid/invalid Sprint starts and closures.

Six local management tests pass. The [installed history test](planning-history-tests.txt)
passes negative cases and a positive empty-Sprint deferral case. The
[upgrade/consumer test](upgrade-and-consumer.txt) passes old-profile installation
→ explicit managed refresh → profile 0.2, preserving original history/project
notes, blocking downgrade, repeating without changes and exposing installed
facts/KanBan through the prebuilt SDPTool. This upgrade run preceded the final
consumer history hardening; the installed history test and final engine checks
cover that correction.

Both skill metadata/reference groups, both reproducible artifact/negative-input
groups and skill-creator quick validation pass. [Native catalog evidence](catalog.json)
observed Codex CLI 0.156.1: fourteen enabled unique skills at project root, nested
cwd, separate worktree and actual install-v1 destination; an independent nested
repository did not inherit the parent collection. Catalog evidence establishes
discovery, not full instruction loading. The later deferred-Sprint wording fix
does not change catalog names/paths; hashes identify the measured guidance.

Four [independent forward trials](report.md) produced actual plan/disposition
artifacts under trials/. These are qualitative planning outcomes, not execution
of the fictional jobs. The [independent implementation review](implementation-review.md)
required and verified three corrections: deferred-Sprint closure, consumer plan
type immutability and ordinary empty Markdown metadata cells. All are resolved;
review approval is scoped and does not impersonate owner/release acceptance.

Go SDPTool race tests and vet pass. The complete legacy PowerShell fixture suite
passes with the new inventory. Full Toolkit Python run: 104 tests, one known
repository-baseline failure, 19 dependency skips; PowerShell/consumer tests were
run separately with explicit executables. The Toolkit validator reports exactly
the same 38 SK1/KB-SDP-011 findings, with no additions/removals.

PL1-C owns final integrated evidence/record checks and closure. No live XFMD
rollout, full fault-matrix rerun, remote CI pass or release publication is claimed.

## PL1-C-M1 — final verification and closeout

Implementation candidate: 16dacf2. All eight files/artifacts identified in
[final independent review hashes](implementation-review-final-source-hashes.json)
match this commit. The remaining C-phase diff contains current-document/index,
evidence and lifecycle closure updates only. The shared Sprint entrypoint was
reconciled with the new plan-grouping contract, retaining the completed historical
Sprint's actual state.

[Final engine checks](final-engine-tests.txt) passed three integrated groups:
clean apply/no-change repetition, manual history-prefix preservation and actual
installed SDPTool workflow. The final installed malformed-history suite includes
the corrected empty-cell path and passes. Recovery/journal implementation was not
redesigned; this work does not claim a new exhaustive interruption matrix. The
prior IU3 evidence remains attached to its own candidate.

Reproduction on Linux: set SDP_TEST_PWSH to PowerShell 7.4+ and SDP_TEST_TOOL to
the prebuilt SDPTool; run Toolkit/tests/test_process_install.py groups named in
the retained logs. Run local management tests with unittest discover under
SDP/ProjectManagement, skill/profile tests under Toolkit/tests, Go test -race ./...
and go vet ./... under Toolkit/SDPTool. The catalog probe is
SDP/Maintenance/SK1/verify_catalog.py with the PowerShell executable argument.
Host evidence used PowerShell 7.6.6, Go 1.27.1 and Codex CLI 0.156.1.

Management, lineage, historical-prefix, generated-output and local-link checks
pass. Full Toolkit validation retains exactly 38 pre-existing findings under
KB-SDP-011; this is not a green repository-wide validation claim. Windows profile
execution remains experimental. No live consumer migration or release occurred.

Backlog review: KB-SDP-014 still owns standalone KanBan distribution; KB-SDP-018
still owns the broader Toolkit audit. Neither is closed by typed planning.
KB-SDP-029's selected plan, skill, validation and distribution are delivered.

Final closure checks: 36 cards, nine management records, three lineage operations
and 247 events; six management tests and 15 lineage negative cases pass. Document
checks preserve 105 frozen records/ledger prefixes and 574 generated outputs,
resolving 2,445 local links and 130 fragments. git diff --check passes.
