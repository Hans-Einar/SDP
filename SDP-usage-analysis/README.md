# Recent SDP Usage Analysis

This directory is the evidence and analysis surface for
`SPR-SDP-005` / `ITR-SDP-005-001` / `SLC-SDP-005-001`, implementing the study
contract in [Issue #5](https://github.com/Hans-Einar/SDP/issues/5). It is study
material, not a replacement for the canonical SDP Toolkit or its instructions.

## Evidence labels

- **OWNER DIRECTION** — a constraint or decision stated by the repository owner
  in Issue #5 or its owner-authored comments.
- **OBSERVED** — a fact read from GitHub, an advertised Git ref, a reachable Git
  commit, or repository content at an identified commit.
- **INFERENCE** — an interpretation derived from observed evidence.
- **RECOMMENDATION** — a proposed future practice; it is not current canonical
  behavior until separately reviewed and accepted.

Unlabelled prose in this file describes the study structure or method. Reports
must label substantive evidence and conclusions with this vocabulary.

## Fixed inventory contract

- **INFERENCE — fixed study method:** the inclusive activity window is anchored
  to Issue #5's creation time and is
  `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`.
- **OWNER DIRECTION:** every repository owned by `Hans-Einar` and visible to the
  authenticated GitHub account is considered, including private, archived, and
  fork repositories.
- **INFERENCE — fixed study method:** the owner's “actual commits” criterion is
  made deterministic by requiring at least one commit whose *committer*
  timestamp is inside the window and which is reachable from an advertised Git
  ref. Repository creation, `updated_at`, `pushed_at`, and author timestamps are
  discovery inputs, not final inclusion evidence.
- **OWNER DIRECTION:** the exact default-branch tip used as the repository study
  commit is recorded separately from the commit or commits proving activity.
- **OBSERVED:** the resulting considered, included, and excluded sets, exact
  commit evidence, commands, limitations, and one-to-one report manifest are in
  [RepositoryInventory.md](RepositoryInventory.md).

## Completed repository corpus

The study contains one file beneath `repositories/` for each of the 17 in-scope
repositories. The inventory is the authoritative report-path manifest. All 17
reports have been Master-reviewed and pass `validate_analysis.py`; synthesis is
the active work boundary.

Each repository report must cover:

1. repository name, default branch, exact studied commit, and relevant activity
   dates;
2. SDP installation and apparent version/generation, relevant GitHub Issues and
   pull requests, and the locations or absence of `AGENTS.md`,
   `AGENTS-project.md`, `SDP/`, `Instructions/`, `Features/`, `Steering/`,
   `Sprints/`, `Refactors/`, `Traceability/`, review, and verification records;
3. actual work-start and delegation practice, including the authoritative role
   (if any) of a GitHub Issue, Master/Worker/Architect/Verifier/Reviewer roles,
   branch/draft-PR creation, and issue/PR/SDP coordinate binding;
4. actual use of Feature, Refactor, Fix, Sprint, Iteration, Slice, Studies and
   convergence, independent review, handoff, `CurrentIndex`, `Relations`,
   `Ledger`, Steering, and `CurrentAssignment`;
5. what worked well and what introduced pain or accidental complexity; and
6. explicit `Carry forward` and `Legacy / do not carry forward` sections, with
   every substantive statement labelled `OBSERVED`, `OWNER DIRECTION`,
   `INFERENCE`, or `RECOMMENDATION` as applicable.

Folder presence alone is not evidence that a process was used. Reports must
distinguish copied Toolkit material from project-local practice, current or
accepted practice from abandoned experiments, and automated state from
conversational convention.

## Synthesis and Steering reading order

Read these first:

1. [ProposedSDPWorkflow.md](ProposedSDPWorkflow.md) — the Master-owned future
   hierarchy, semantics, Issue/Master/Steering model and traceability decision;
2. [CrossRepositoryPatterns.md](CrossRepositoryPatterns.md) — observed patterns,
   counterevidence and constraints on the recommendation;
3. [MigrationImpact.md](MigrationImpact.md) — profile-aware migration plus
   SDP-Analyzer and gh-sdp implications;
4. [LegacyAndDeprecation.md](LegacyAndDeprecation.md) — what remains durable,
   what is deprecated and the removal gates; and
5. [EvidenceMatrix.md](EvidenceMatrix.md) — exact 17-row comparative evidence.

Supporting decisions:

- [SkillsAndRoles.md](SkillsAndRoles.md) — ChatGPT `steering-group` versus `sdp`
  and Codex skill updates;
- [FollowUpIssues.md](FollowUpIssues.md) — dependency-ordered implementation
  assignments; and
- [RepositoryInventory.md](RepositoryInventory.md) plus `repositories/` — exact
  scope and repository evidence.

Run the deterministic corpus check with:

```powershell
python SDP-usage-analysis/validate_analysis.py
```

## Study boundary

- **OWNER DIRECTION:** studied repositories are read-only.
- **OWNER DIRECTION:** this study does not change canonical SDP documentation,
  Toolkit schemas/templates/installers, `gh-sdp`, `SDP-Analyzer`, release state,
  or any consuming repository.
- **OWNER DIRECTION:** pull request #4 is unmerged evidence only and is not this
  study's base.
- **OWNER DIRECTION:** the work stops at a recommendation ready for Steering
  Group review; implementation and migration require later assignments.
