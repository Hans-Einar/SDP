# SPR-SDP-005 — Recent SDP Usage Study

Status: complete — awaiting Steering Group review

GitHub issue: `Hans-Einar/SDP#5`

Authoritative baseline: `e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15` (`main`)

Study branch: `codex/issue-5-sdp-usage-study`

## Goal

Produce the evidence-backed, cross-repository SDP usage study required by Issue
#5 and stop at a recommendation ready for Steering Group review.

## Scope

- Establish the exact four-month repository inventory from actual commit
  evidence.
- Create one evidence-backed report for every in-scope repository.
- Synthesize observed practices into a proposed future SDP workflow.
- Cover Feature, Refactor, Fix, Sprint, Iteration, Slice, GitHub Issue,
  Master-per-Issue, Steering Group, CurrentAssignment, traceability, migration,
  legacy, SDP-Analyzer, gh-sdp and reusable skill implications.
- Validate the corpus and obtain fresh independent review.

## Invariants

- Issue #5 is the durable work contract.
- The study baseline is current authoritative `main` at the commit above.
- Studied repositories are read-only.
- Pull request #4 is unmerged evidence only and must not be modified, merged or
  used as this branch's base.
- Evidence is labelled `OBSERVED`, `OWNER DIRECTION`, `INFERENCE` or
  `RECOMMENDATION` as applicable.
- Every repository report records its exact studied commit.
- The Master owns inventory decisions and synthesis; Workers provide bounded
  evidence reports; a fresh Reviewer independently challenges the result.
- No canonical SDP, Toolkit, gh-sdp or SDP-Analyzer implementation change is in
  scope.

## Non-goals

- Implementing the recommended future workflow.
- Changing Toolkit schemas, templates, installer behavior or release state.
- Modifying studied repositories.
- Closing Issue #5 or merging the study pull request.

## Completion signal

`STUDY_READY_FOR_STEERING_REVIEW` may be declared only when the complete corpus
exists, deterministic validation passes and independent review has no unresolved
Blocking, High or Medium finding.
