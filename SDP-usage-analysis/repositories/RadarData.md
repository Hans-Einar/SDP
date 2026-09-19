# RadarData

## Repository facts and evidence boundary

- **OBSERVED:** Repository: `Hans-Einar/RadarData` (private); default branch:
  `main`; exact default-branch study commit:
  `ca9129209eb2eb5d10854b7043f8255a69e90953`, committed
  `2026-07-26T09:34:38Z`.
- **OBSERVED:** The fixed-window inventory found 84 qualifying reachable
  commits. The first was `f2e0f918f46bcd5da50d459e90b578906e25de27`
  at `2026-07-25T10:57:27Z` on
  `agent/add-radardata-poc-concept`; the last was synthetic pull-request merge
  ref commit `d84a5718c5ddc307efb72cddb9963149b286d64b` at
  `2026-07-26T23:59:05Z` on `refs/pull/2/merge`.
- **OBSERVED:** GitHub currently says PR
  [#2](https://github.com/Hans-Einar/RadarData/pull/2) is open and unmerged.
  Therefore its `refs/pull/2/merge` inventory object is evidence of advertised
  activity, not evidence that PR #2 was merged.
- **OBSERVED:** The default study tree contains eight files: the README, three
  concept/study documents, two radar-probe scripts and a probe workflow. It has
  no `AGENTS.md`, `AGENTS-project.md`, `SDP/`, or installed SDP manifest.
- **OBSERVED:** PR
  [#1](https://github.com/Hans-Einar/RadarData/pull/1) merged the concept and
  source-study material to `main` at the exact study commit. It has no Issue or
  PR comments, submitted reviews, review threads, or commit status contexts.
  One pull-request-triggered `Radar archive probe` workflow run
  (`30179012016`) succeeded.
- **OBSERVED:** Current Issue
  [#3](https://github.com/Hans-Einar/RadarData/issues/3) is open with one
  milestone comment. Current draft PR #2 targets the study commit, has 56
  commits and 48 changed files, and remains at exact head
  `16d7a5f7030d8d16968c2b5fb872c1f0b1bdbb76`. Its final-head
  `Validate SDP bootstrap` workflow run `30226387974` succeeded. The PR has one
  top-level comment, no submitted reviews, no inline review threads and no
  commit status contexts.
- **OBSERVED:** PR #2 is the only current SDP installation evidence. Its
  project-local profile is `radardata.sdp-feature-pilot.v1`; no manifest or
  released Toolkit identity establishes a canonical SDP version. It is
  explicitly a draft bootstrap awaiting Steering review.
- **INFERENCE:** Canonical default-branch truth is “concept and studies, no
  SDP.” The richer Feature-first material below is useful current in-flight
  evidence only and must not be reported as accepted or installed practice.

## SDP and agent locations

- **OBSERVED — default study commit:** No `AGENTS.md`, `AGENTS-project.md`,
  `AGENT-REMINDERS.md`, `SDP/`, `Instructions/`, `Features/`, `Steering/`,
  `Sprints/`, `Refactors/`, `Fixes/`, `Traceability/`, `CodeReview/`,
  `Verification/`, `Handoff.md`, `CurrentIndex.yaml`, `Relations.yaml` or
  `Ledger.ndjson` exists.
- **OBSERVED — open PR #2 head:** There is still no `AGENTS.md`,
  `AGENTS-project.md`, `AGENT-REMINDERS.md` or `SDP/Instructions/`.
  `SDP/` contains numbered Mandate through Implementation locations,
  `ProductVision/`, `UseCases/`, nine `Features/`, `Steering/`,
  `Verification/`, `CodeReview/` and `Traceability/`.
- **OBSERVED — open PR #2 head:** No tracked `Releases/`, `Sprints/`,
  `Refactors/`, `Fixes/`, local Sprint/Iteration/Slice records or `Handoff.md`
  exists. They appear in the proposed directory diagram, but Git does not
  preserve empty directories; the report does not infer their use.
- **OBSERVED — open PR #2 head:** The operating records are
  `docs/governance/RadarData-SDP-Feature-Bootstrap-and-GitHub-Workflow-v0.1.md`,
  `SDP/README.md`, `SDP/Steering/README.md`,
  `SDP/Steering/CurrentAssignment.yaml`, one Steering interaction,
  `SDP/Traceability/{CurrentIndex.yaml,Relations.yaml,Ledger.ndjson}`,
  `SDP/Verification/VER-BOOT-001.md`, and
  `SDP/CodeReview/REV-BOOT-001.md`.

## How work actually starts and is controlled

- **OBSERVED — merged default truth:** Work started as concept and bounded
  source studies on `agent/add-radardata-poc-concept`, delivered through draft
  PR #1 and merged without a GitHub Issue, formal review, SDP assignment or
  traceability records.
- **OBSERVED — open PR #2 evidence:** The branch was created and its draft PR
  opened before Issue #3: PR #2 was created at `2026-07-26T23:45:23Z`, while
  Issue #3 was created at `2026-07-26T23:55:24Z`. The Ledger calls the source
  assignment `authorized-in-chat`, and the later Issue became the durable
  review/disposition boundary for `BOOT-001`.
- **OBSERVED:** Issue #3 is contract-like: work-item ID, type, related Feature,
  exact base, branch, draft PR, goal, scope, non-goals, required verification
  and stop boundary. Its comment reports the exact final head, workflow run,
  exclusions, pending fresh review and preserved stop.
- **OBSERVED:** The governance proposal defines Issue as authorized work item
  and communication anchor, PR as implementation/evidence/review/integration
  surface, the Steering interaction as exact assignment/response/assessment,
  and the Ledger as append-only lifecycle evidence. It explicitly says one
  Issue is not necessarily one Feature.
- **OWNER DIRECTION:** Issue #5 requires GitHub Issue identity to become the
  operational assignment boundary and a bounded Codex session to act as Master.
- **RECOMMENDATION:** Carry forward RadarData's Issue fields and stop boundary,
  but create/authorize the Issue before branch work, then derive branch, PR and
  assignment state from that one contract.

## Master, delegation and roles

- **OBSERVED:** The in-flight proposal says “Codex” should open a draft PR,
  preserve the assignment, update repository/Issue/PR evidence, verify/review
  and stop. It does not define Master, Worker, Architect, Verifier or Reviewer
  responsibilities and contains no delegation records.
- **OBSERVED:** `VER-BOOT-001` is a verification artifact, but it does not name
  an independent Verifier. `REV-BOOT-001` is still “Pending fresh review,” and
  GitHub has no review submission. There is no evidence that a fresh Reviewer
  performed the review.
- **INFERENCE:** RadarData demonstrates a bounded Codex delivery session and a
  separate intended Steering decision, but not the Issue #5 Master-per-Issue
  delegation model or independent-role provenance.
- **RECOMMENDATION:** Retain the compact Issue/PR/stop mechanics while adding
  explicit Master ownership, delegated role/subject/head, and a fresh-review
  identity and outcome.

## Feature, Refactor, Fix, Sprint, Iteration and Slice

- **OBSERVED — open PR #2 evidence:** Nine proposed `FEAT-NNN` records form a
  capability backlog. Each has an adjacent `integration-study.md` and relations
  to Use Cases and the horizontal architecture. A Feature record alone does not
  authorize implementation.
- **OBSERVED:** The proposal defines Use Case as outcome, Feature as reusable
  capability, horizontal layer as responsibility boundary and vertical Slice
  as observable delivery across required layers. Eight `SLC-NNN` outcomes are
  proposed, beginning with a web-to-server walking skeleton.
- **OBSERVED:** No Feature, Sprint, Iteration or Slice has been implemented or
  activated. `CurrentIndex.active` is intentionally empty and CI checks that
  condition. No Refactor or Fix contract exists in the tree.
- **INFERENCE:** RadarData is evidence that Feature records, per-Feature
  integration studies and vertical Slice planning can be expressed clearly;
  it is not evidence that this hierarchy has survived real implementation,
  rework, extension or release.
- **OWNER DIRECTION:** Issue #5 requires implementation capability to belong to
  Feature or Refactor and asks whether Sprint/Iteration remain useful.
- **RECOMMENDATION:** Use RadarData's Feature/Use Case/architecture relations
  and observable Slice contracts as candidate schema. Do not use this dormant
  bootstrap as evidence that Sprint is unnecessary or that its Feature granularity
  is proven.

## Initial design, later refinement and Studies

- **OBSERVED — default truth:** PR #1 delivered several pre-Mandate studies.
  The bootstrap later classifies them as evidence inputs rather than accepted
  decisions and leaves formal Requirements unbaselined pending a corrected
  hourly precipitation source Study.
- **OBSERVED — open PR #2 evidence:** Draft architecture establishes horizontal
  web, application/API, durable job, scientific/domain, source-adapter and
  persistence boundaries. The proposed implementation plan cuts observable
  vertical outcomes through them.
- **OBSERVED:** Each Feature has an integration-study placeholder/proposal, and
  the governance text allows several studies and Slices for one Feature. There
  is no exercised parallel independent-Study process, synthesis artifact,
  convergence gate or post-foundation design revision.
- **INFERENCE:** Separating pre-Mandate evidence from accepted SDP decisions is
  a useful truthfulness pattern. The adjacent integration studies are useful
  local-impact prompts, but their unreviewed proposal status prevents stronger
  conclusions.
- **RECOMMENDATION:** Preserve “evidence input is not authority,” formal Study
  authorization, and horizontal-skeleton/vertical-growth language. Require a
  named synthesis/review/Steering gate when multiple studies change accepted
  requirements or design.

## Review, verification and handoff

- **OBSERVED:** The repository validator checks required paths, inactive work
  state, Use Case/Feature paths, relation paths and basic Ledger JSON/unique IDs.
  A GitHub Actions workflow runs it for relevant PR and branch changes.
- **OBSERVED — contradiction:** `VER-BOOT-001` and the last Ledger event name
  verified head `9dcad0b6...` and workflow `30226313156`; `Relations.yaml` names
  verified head `f8c5c062...` with that same workflow; final Issue/PR comments
  and current GitHub state name final head `16d7a5f...` and successful run
  `30226387974`. Final-head CI exists, but the repository's manually copied
  verification coordinates disagree.
- **OBSERVED:** `REV-BOOT-001` remains pending, no GitHub review exists, and the
  Steering assessment is pending. No `Handoff.md` exists; recovery context is
  distributed among Issue #3, PR #2, CurrentAssignment, interaction, indexes and
  comments.
- **INFERENCE:** Automated structure validation improved confidence, while
  manually repeated exact heads weakened the otherwise strong evidence model.
- **RECOMMENDATION:** Verification records should bind a subject and immutable
  tested head; live head/check state and review summaries should be generated
  from GitHub. A changed head must invalidate or supersede earlier verification.
  Keep handoff prose only for unresolved risk and next legal action.

## CurrentIndex, Relations, Ledger and Steering assignment

- **OBSERVED:** `CurrentAssignment.yaml` compactly binds interaction,
  `BOOT-001`, Issue #3, PR #2, status, branch and stop boundary. It omits exact
  base/head, review, verification and next gate.
- **OBSERVED:** `CurrentIndex.yaml` provides current pointers; `Relations.yaml`
  supplies Mandate/Product Vision/Use Case/Feature/architecture/plan/Steering/
  verification/review edges; six Ledger events preserve assignment, branch, PR,
  catalogue, Issue and verification transitions.
- **OBSERVED:** The same status and GitHub coordinates are copied across
  CurrentAssignment, CurrentIndex, Relations, Ledger, interaction, Issue, PR,
  verification and comments. The verification-head contradiction demonstrates
  concrete drift from that duplication.
- **OBSERVED:** The Issue and PR templates standardize work-item identity,
  exact base, branch, Feature IDs, scope, verification/review, traceability and
  stop boundary. Creation and synchronization remain conversational/manual;
  only structure validation is automated.
- **INFERENCE:** The compact assignment record and semantic graph are useful;
  repeating mutable GitHub/check/head state is not.
- **RECOMMENDATION:** Keep a minimal machine-readable assignment containing
  Issue/work ID, authority baseline, active branch/PR, active work unit,
  authorization/stop, delegated roles and next gate. Generate live PR/head/
  check/review status, record durable semantic relations once, and append only
  material decisions/corrections to the Ledger.

## What worked well

- **OBSERVED:** Both PRs are narrow in declared purpose, use draft review
  surfaces, state exclusions explicitly and preserve stop boundaries.
- **OBSERVED:** Default `main` remains untouched by the unaccepted bootstrap,
  making the default-versus-proposal distinction recoverable.
- **OBSERVED:** Issue #3 and PR #2 bind a stable work ID, Feature, baseline,
  branch, verification and stop condition, and final-head CI passed.
- **OBSERVED:** The bootstrap cleanly describes a horizontal architecture and
  small observable vertical outcomes without pretending that source feasibility
  or detailed Requirements are already proven.
- **INFERENCE:** RadarData's strongest generic contribution is a compact
  Issue/PR/assignment/stop loop plus automated internal-link/state validation,
  provided its proposal and acceptance status remain explicit.

## Pain points, contradictions and limitations

- **OBSERVED:** The authoritative default branch has no SDP, while the rich
  process exists only in a five-week-old open draft PR awaiting review. Current
  adoption is therefore unresolved.
- **OBSERVED:** The bootstrap began from chat/branch/PR and only later created
  its Issue, contrary to its own desired Issue-authorized lifecycle.
- **OBSERVED:** No `AGENTS.md` or `Instructions/` tells a fresh agent how to
  discover and apply the pilot, and no specialized role/delegation contract is
  present.
- **OBSERVED:** Nine Features, nine integration studies and eight Slices were
  initialized before the primary source Study or Requirements baseline. Their
  explicit “proposed” state is honest, but this creates substantial speculative
  catalogue maintenance if feasibility changes.
- **OBSERVED:** CurrentAssignment is too small for exact-head recovery, while
  the other state surfaces copy too much; the result is both missing and stale
  information.
- **OBSERVED:** Independent review and Steering acceptance remain absent.
  Therefore no claim that the pilot worked operationally is supportable.
- **OBSERVED:** GitHub review/check inspection is limited to objects currently
  exposed for PRs #1 and #2. No access limitation was encountered.

## Carry forward

- **RECOMMENDATION:** Contract-shaped Issue with stable work ID, related
  Feature, exact baseline, branch/draft PR expectation, goal, scope/non-goals,
  verification, independent review, stop and report contract.
- **RECOMMENDATION:** Early draft PR and concise milestone comments with exact
  head and CI identity.
- **RECOMMENDATION:** Compact machine-readable Steering assignment plus one
  semantic relation graph and material-event Ledger.
- **RECOMMENDATION:** Use Case -> Feature -> architecture -> proposed vertical
  Slice relations, with Feature records not authorizing implementation by
  themselves.
- **RECOMMENDATION:** Pre-Mandate evidence remains evidence until a formal
  Study/review/Steering decision promotes it.
- **RECOMMENDATION:** Horizontal responsibility skeleton, server-authoritative
  boundaries and observable end-to-end Slice outcomes.
- **RECOMMENDATION:** Deterministic structure/traceability validation and an
  explicit inactive bootstrap state.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not start an assignment in chat and create the Issue
  after the branch/PR when Issue-first authorization is available.
- **RECOMMENDATION:** Do not treat this unmerged, pending-review pilot as
  installed canonical SDP or proof that its Feature hierarchy is exercised.
- **RECOMMENDATION:** Do not duplicate mutable head, workflow and status data
  across assignment, index, relations, Ledger, verification, Issue and PR.
- **RECOMMENDATION:** Do not accept a passed structure validator as independent
  review or Steering acceptance.
- **RECOMMENDATION:** Do not pre-maintain a large speculative Feature/Slice
  catalogue beyond what remains useful after bounded feasibility evidence.
- **RECOMMENDATION:** Do not infer use of Releases, Sprints, Refactors or Fixes
  from a proposed directory diagram when no records exist.
