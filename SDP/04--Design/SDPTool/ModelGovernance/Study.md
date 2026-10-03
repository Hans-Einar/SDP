# ModelGovernance — study and decision baseline

2026-10-03. Study milestone MG1-M1; planning evidence, not runtime evidence.
[Plan](Plan.md), [Design](Design.md), [Session0006](../../../Sessions/session-%230006--Model_governance.md).

## Problem and outcome

Agents can damage a working model or change its scope without preserving a usable
baseline. The owner needs independent SDL/SDUI model workspaces, local recovery,
reviewable candidates and accepted releases without requiring project Git or
remembered exports. Keep lineage portable as ordinary files, but do not build a
complete version-control system. Blueprint impact analysis is a separate feature.

## Evidence and provenance

[Session0005](../../../Sessions/session-%230005--Blueprint_model_history.md) and
[original KB048](../../../KanBan/superseded/%23048--Proposal--Versioned-design-reviews-and-blueprint-diffs.md)
record the conversation and superseded alternatives. Latest owner corrections take
precedence over earlier recommendations. Current Go command entrypoints are
SDL/go/cmd and SDPTool/cmd; SDPTool reuses language libraries. No model-history,
commit/restore or semantic blueprint command is currently implemented by this work.
The prior disposable Git probe proved only Git storage behavior.

## Owner-selected boundaries

- Work without Git; parent Git may transport normal files but is not an operation prerequisite.
- Mutable WORK; optional PROPOSAL; immutable CANDIDATE and RELEASE. Normal route
  is WORK -> CANDIDATE -> RELEASE, not mandatory proposal creation.
- WORK starts from a release or explicit empty initial state. Integration may
  create a new WORK from other WORKs with recorded ancestry and a known base.
- No UUID in WORK names. Full identity is YAML metadata; four-character UUID
  suffixes may distinguish proposal/candidate display names, with collision checks.
- Commit messages and local changed-file history support recovery during work.
  Promotion may discard restoration payloads but retains messages and lineage.
- Preliminary previews may read WORK without locking or making a persistent snapshot.
- Commands should read naturally: create work:NAME from ..., merge ... into ... .
  Context-relative pull remains unresolved, not a second implemented workflow.
- One designated integration owner creates candidates/releases. This is workflow
  responsibility, not a promise of distributed locking across independent clones.

## Alternatives assessed

| Alternative | Disposition / reason |
| --- | --- |
| Nested Git with export/bundle | Owner rejected extra export dependency |
| Project Git as mandatory model history | Owner rejected repository/commit prerequisite |
| Copy SVN/Fossil database with project Git | Does not inherently combine independent histories safely |
| Permanent content-addressed model repository | Beyond selected scope; defer |
| WORK-local baseline and changed-file copies | Selected direction; requires explicit deletion and recovery rules |
| Full semantic blueprint engine in this feature | Separate KB050; consumes model revisions later |

## Derived requirements and unresolved design

A starting snapshot plus changed-file after-images needs deletion records. Merge
must retain exact inputs and a common base. Whole-state rollback is feasible with
checkpoints even after overlapping edits; selective undo across merged branches is
outside initial scope. Recursive .merge copies need bounded traversal and identity
checks. A clean text merge is not proof of valid model contracts.

Resolve in MG2: on-disk schema/version, crash recovery/publication, hash rules,
local numbering after restore, source identity and repeated merge handling, portable
path constraints, discovery filtering, initial-release bootstrap, and code evidence
references without imposing Git. [Design](Design.md) provides the initial proposals.

## Conclusion and next step

Proceed with bounded ModelGovernance design, then a scratch-directory storage and
merge proof before production CLI implementation. No requirement for a server,
remote PR system, global VCS, text-delta engine or selective cherry-pick. Retain
SDL/SDUI parser ownership and separate blueprint selection/impact responsibilities.
