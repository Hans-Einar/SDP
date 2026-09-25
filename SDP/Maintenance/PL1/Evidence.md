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
