# python SDP usage study

## Evidence frame and repository facts

- **OBSERVED:** This report studies private, non-archived
  `Hans-Einar/python` at the inventory-recorded default branch and exact current
  identity
  [`main@6bd5923ef6ab3041d853a66ed3b63ebbe489aa18`](https://github.com/Hans-Einar/python/tree/6bd5923ef6ab3041d853a66ed3b63ebbe489aa18),
  committed `2026-06-19T22:23:24Z`. The live GitHub branch record still agreed
  on that identity when checked on 2026-08-21. The studied repository was not
  modified.
- **OBSERVED:** In Issue #5's inclusive window,
  `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`, two commits
  qualified: `49b3c0a959b0cc5ae962c4e44dc015b683fee89a` at
  `2026-06-19T22:22:08Z` and the studied merge commit at
  `2026-06-19T22:23:24Z`. The captured repository had three advertised refs
  and 41 reachable commits.
- **OBSERVED:** The first qualifying commit, `ponsse parser update`, reorganized
  Ponsse protocol captures/tools and added analysis scripts and generated
  artifacts. The second is a direct two-parent merge of two `main` histories
  and adds further Ponsse captures while consolidating old parsers. Neither
  qualifying commit is connected to a GitHub Issue or PR.
- **OBSERVED:** No SDP installation or version marker exists. The exact default
  tree contains no `AGENTS.md`, `AGENTS-project.md`, `SDP/`, `Instructions/`,
  `Features/`, `Steering/`, `Sprints/`, `Refactors/`, `Fixes/`,
  `Traceability/`, `CodeReview/`, `Verification/`, Handoff, `CurrentIndex`,
  `Relations`, Ledger or repository-owned `.github/workflows/` path.
  `ponsse/PONSSE_PROTOCOL_ANALYSIS.md` and tool READMEs are domain-analysis
  artifacts, not SDP Study/requirements/design records.
- **OBSERVED:** Current GitHub state has zero Issues and one historical pull
  request: [PR #1](https://github.com/Hans-Einar/python/pull/1), created and
  merged on 2026-02-02, outside the study window. Its retained head branch is
  `copilot/fix-protocol-issues@55ba7cf85d06421b3129759be0e1a8ecbf813a76`;
  merge commit `d42aa59bb40ff5864de02eff265e388c9326d54d` is an ancestor of
  the studied default tip.
- **OBSERVED — historical PR evidence:** PR #1 was opened by GitHub Copilot,
  was not a draft, contained 12 commits, and supplied a detailed change summary,
  usage examples, test command and the original prompt. The owner submitted one
  native `COMMENTED` review with body `Ok` and one top-level `Ok` comment. There
  was no approval, change request, inline review comment or reported check.
- **OBSERVED:** GitHub currently exposes one active hosted/dynamic workflow,
  `Copilot coding agent` at `dynamic/copilot-swe-agent/copilot`, but the default
  study commit has zero check runs and zero status contexts, and PR #1 reports
  no checks on its head. Neither live branch is protected.
- **OBSERVED — limitation:** PR #1 is useful historical agent/PR evidence, but
  it predates the qualifying June work and cannot be treated as its assignment,
  review or verification contract. Default-tree truth and the retained
  non-default Copilot branch are therefore kept separate.

## How work starts and is controlled

- **OBSERVED — historical PR evidence:** PR #1 began from a prompt preserved in
  its body and used a Copilot-owned feature branch. It provided more durable
  scope and test context than the later direct commits, but there was no
  separate Issue and no Feature/Refactor/Fix or SDP identity.
- **OBSERVED — qualifying/default evidence:** The June work began as direct
  repository commits and a merge of `main` histories. There is no durable
  goal/why-now/scope/non-goal contract for those changes and no early draft PR.
- **OBSERVED:** No Codex Master, Worker, Architect, Verifier or Reviewer model is
  defined. PR #1 proves use of a Copilot coding agent and an owner comment; it
  does not prove bounded Master delegation, independent review, fresh context
  or exact-head verification.
- **OBSERVED:** No Feature, Refactor, Fix, Sprint, Iteration or Slice was used.
  The Ponsse protocol analysis and scripts show iterative technical study, but
  no process record splits independent Studies, defines a convergence gate, or
  relates the analysis to requirements, architecture, design or a vertical
  delivery contract.
- **OBSERVED:** There is no Handoff, `CurrentIndex`, `Relations`, Ledger,
  Steering/`CurrentAssignment`, Issue/PR-to-SDP binding, standalone review or
  verification record. The dynamic Copilot workflow supports agent execution;
  it did not supply CI evidence for the studied commit or PR.
- **OWNER DIRECTION:** Issue #5 makes a GitHub Issue the future operational
  assignment boundary and a bounded Codex session Master for that Issue. A
  detailed agent PR body is useful evidence but is not a substitute for the
  owner/Steering assignment contract.

## What worked well

- **OBSERVED:** PR #1 isolated agent work on a branch, retained the initiating
  prompt, documented intended changes and usage, and exposed the result on a
  GitHub review surface before merge.
- **OBSERVED:** The PR body named a focused test suite and documented a deliberate
  tradeoff: duplicate the decoder rather than risk breaking the existing
  snooper. This makes the implementation rationale more recoverable than the
  June commit subjects alone.
- **OBSERVED:** Later default history preserves the merged Copilot work while
  continuing the Ponsse investigation with additional captures and analysis
  tools.
- **INFERENCE:** The repository demonstrates that lightweight branch/PR agent
  work can improve recoverability even without SDP, but not that the result had
  independent or automated acceptance evidence.

## Pain points and limitations

- **OBSERVED:** The only PR is outside the activity window; the qualifying work
  reverted to direct commits and a generic merge without an Issue, PR, review
  or check surface.
- **OBSERVED:** The owner review and comment both say only `Ok`; neither records
  examined invariants, exact verification, limitations or disposition as a
  formal approval.
- **OBSERVED:** No CI/check evidence validates either the historical PR or the
  studied default commit. The hosted Copilot workflow is execution automation,
  not verification automation in the observed runs.
- **OBSERVED:** The repository mixes many unrelated Python experiments and
  generated/captured artifacts. No Feature or bounded work identity says which
  requirements, design decisions and evidence belong together.
- **INFERENCE:** Retroactively classifying the domain analysis files as SDP
  Studies would overstate the process. Their value is technical content; their
  assignment and convergence semantics were not recorded.

## Carry forward

- **RECOMMENDATION:** Preserve PR #1's useful lightweight practices: isolated
  agent branch, detailed change summary, initiating prompt/reference, explicit
  rationale, usage and focused verification command.
- **RECOMMENDATION:** Add a durable GitHub Issue before the next material agent
  change, binding an exact baseline, Feature/Refactor/Fix identity, scope,
  protocol invariants, expected areas, verification, independent review, branch,
  draft PR and stop condition.
- **RECOMMENDATION:** When exploratory protocol analysis is needed, give the
  Study a bounded claim/evidence question and an explicit decision/convergence
  point before implementation.
- **RECOMMENDATION:** Keep process adoption light for this multi-experiment
  repository: one current assignment/Handoff and stable evidence links are more
  valuable than a large manually reconstructed historical registry.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not treat an agent-generated PR body or embedded prompt
  as sufficient owner assignment authority when a durable Issue can carry the
  contract.
- **RECOMMENDATION:** Do not treat a terse `Ok` comment, unreported local test
  claim or dynamic agent workflow as independent review and exact verification.
- **RECOMMENDATION:** Do not carry forward direct default-branch work for
  material changes after the repository has demonstrated a workable branch/PR
  path.
- **RECOMMENDATION:** Do not infer SDP lifecycle use from technical files named
  analysis, README or test; retain them as domain evidence and add prospective
  process binding only when authorized.
