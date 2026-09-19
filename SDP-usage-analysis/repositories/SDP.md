# Hans-Einar/SDP usage study

This report studies the canonical method repository itself. Its root lifecycle
files are simultaneously examples/templates and, for release/traceability, live
Toolkit records. That dual role is material evidence rather than an assumption
that every folder represents an executed process.

## Evidence frame

- **OBSERVED — canonical baseline:** the exact default-branch study commit is
  `e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15` on `main`, committed
  `2026-07-17T12:59:57Z`. Unless explicitly marked otherwise, canonical file
  claims below refer to that tree.
- **OBSERVED — current-state checkpoint:** Issue
  [#5](https://github.com/Hans-Einar/SDP/issues/5), draft PR
  [#6](https://github.com/Hans-Einar/SDP/pull/6), and the study branch at
  `74b933db431c6b92ddfc76bffcb23d825c1e1130` are later execution evidence from
  `2026-08-20`; they are not part of the fixed activity-window inclusion proof
  or the canonical baseline.
- **OBSERVED — unmerged-evidence boundary:** draft PR
  [#4](https://github.com/Hans-Einar/SDP/pull/4) was inspected only at exact head
  `d611b8bf72aeb30d86c5ef28902469462b06a803`. It is open, draft, unmerged, has
  no GitHub review objects or conversation comments, and is not treated as
  `main` truth. Its repository records are valuable evidence of an attempted
  execution and of contracts not yet accepted into the default branch.
- **OWNER DIRECTION:** Issue #5, created `2026-08-20T21:08:52Z`, is newer than
  the July Feature-governance proposal and explicitly directs the study not to
  preserve Sprint as capability owner merely for compatibility.

## Repository facts

- **OBSERVED — identity:** `Hans-Einar/SDP` is a public, non-archived repository;
  its default branch is `main`. It was created `2026-07-11T09:59:13Z` and has no
  declared license in GitHub metadata.
- **OBSERVED — activity:** the inventory found 12 advertised refs, 108 unique
  reachable commits, and 105 commits with committer timestamps in the inclusive
  `2026-04-20T21:08:52Z`–`2026-08-20T21:08:52Z` window. The earliest qualifying
  commit is `2eef0fa98a146ecfc345ffde5041ba2e03adae6e` at
  `2026-07-11T09:59:06Z`; the latest qualifying advertised-ref evidence is
  synthetic PR #4 merge ref commit
  `adcad55319205f8aa27ceeb8fe3f87620e8c2d5d` at
  `2026-07-17T15:27:23Z`. The exact default-branch study commit remains
  `e398ebaf...`, not that pull-ref commit.
- **OBSERVED — installation/generation:** this is the upstream Toolkit authoring
  repository, not a normal installed project. At the baseline it has no root
  `AGENTS.md`, no `AGENTS-project.md`, and no `SDP/` directory. Its
  `SDP.manifest.yaml` declares Toolkit `0.2.0` as `unreleased`, Framework and
  AGENTS contracts `1.0.0`, ten versioned `1.0.0` Toolkit skills, and a migration
  compatibility range `>=0.1.0 <1.0.0`. There are no Git tags and no GitHub
  Releases.
- **OBSERVED — agent entry locations:** the actual managed entry files are
  templates at `Toolkit/payload/project-root/AGENTS.md.template` and
  `Toolkit/payload/project-root/AGENTS-project.md.template`; project-local
  reminders are templated at
  `Toolkit/payload/sdp-root/AGENT-REMINDERS.md.template`. The absence of a live
  root agent entry means the repository cannot reconstruct its own agent
  operating contract solely from the studied tree without externally supplied
  workspace instructions.
- **OBSERVED — lifecycle/template locations:** the baseline root contains
  `01--Mandate/` through `07--Implementation/`, `Instructions/`, `Sprints/`,
  `Refactors/`, `Fixes/`, `CodeReview/`, `Verification/`, `Traceability/`, and
  `Releases/`. The numbered lifecycle documents and most operating-area files
  are short `Status: template` records or README guidance.
- **OBSERVED — absent locations:** `Features/`, `Steering/`, and
  `CurrentAssignment.yaml` do not exist at the canonical baseline. Nor do
  populated Sprint, Refactor, Fix, CodeReview, or Verification records.
- **OBSERVED — live baseline state:** `Traceability/CurrentIndex.yaml`,
  `Traceability/Relations.yaml`, `Traceability/Ledger.ndjson`,
  `Releases/REL-0.2.0.yaml`, `RELEASE-NOTES.md`, and `SDP.manifest.yaml` are live
  Toolkit release state, not neutral examples. `CurrentIndex` has no active
  Sprint/Refactor/Iteration/Slice/Fix but declares `REL-0.2.0` unreleased;
  `Relations` maps that release; the Ledger contains two release events.
- **OBSERVED — skills:** the canonical managed skills are under
  `Toolkit/skills/`: Master, Worker, Reviewer, Architect, Traceability, Vertical
  Refactor, Release, Versioning, Auditor, and Verifier. A second root `skills/`
  tree contains only six older, unversioned variants and is not the installer
  source. The baseline also retains a root `payload/` copy and a root
  `scripts/Install-SDP.ps1` compatibility forwarder alongside canonical
  `Toolkit/` assets.
- **OBSERVED — automation:** `.github/workflows/validate.yml` runs Python
  contract tests on Linux and PowerShell installer fixtures on Windows for pull
  requests. `Toolkit/scripts/validate_sdp.py` validates manifests, versioned
  skills, Release/Fix records, release notes, and release Ledger events; it only
  checks that `CurrentIndex.yaml` and `Relations.yaml` parse as YAML objects at
  the baseline. `Toolkit/scripts/Install-SDP.ps1` installs managed agent/skill/
  Framework files, preserves project-owned records, and optionally initializes
  lifecycle folders.

### Relevant Issues and pull requests

- **OBSERVED — PR #1:**
  [#1](https://github.com/Hans-Einar/SDP/pull/1), branch
  `feature/sdp-toolkit-and-skills`, created `2026-07-11T10:10:33Z`, merged
  `2026-07-11T11:48:05Z` as
  `e435a5f05bc1aa1759318378794c504fc437c1e6`. It introduced the reusable
  Toolkit/method shape. GitHub records zero checks and zero reviews.
- **OBSERVED — PR #2:**
  [#2](https://github.com/Hans-Einar/SDP/pull/2), reusing
  `feature/sdp-toolkit-and-skills`, created `2026-07-11T12:32:22Z`, merged
  `2026-07-11T12:42:47Z` as
  `52c5a2b70e9517acdb6a316c84fe065a9c8abfcf`. It split Toolkit-managed
  `AGENTS.md` from project-owned `AGENTS-project.md`; GitHub again records zero
  checks and zero reviews.
- **OBSERVED — PR #3:**
  [#3](https://github.com/Hans-Einar/SDP/pull/3), branch
  `agent/sdp-release-versioning`, created `2026-07-12T08:10:32Z`, merged
  `2026-07-12T20:12:17Z` as
  `bc110bb5fd60009ba67015cf640ad6ddbfe1b04b`. Exact product head
  `e05f74ddf8fb222398fab708176b5603e80adddf` had green Linux/Windows CI and one
  owner-authored GitHub review comment describing a fresh independent Reviewer
  pass. It did not add a populated Sprint, CodeReview, or Verification record.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:**
  [#4](https://github.com/Hans-Einar/SDP/pull/4), branch
  `codex/sdp-install-contract-v1`, created `2026-07-14T00:15:25Z`, remains an
  open draft at exact head `d611b8bf...`. Its two exact-head checks passed on
  `2026-07-17`. It contains the repository's richest actual Sprint, study,
  review, verification, traceability, and handoff execution, but none of that is
  merged canonical behavior.
- **OBSERVED — direct main governance changes:**
  `docs/Feature-Governance-And-SDP-2.0.md` was added at
  `b90e4511973c990b17e9122e3baf5446a7e6f854` on `2026-07-17`; Architect skill
  Steering-interaction guidance was added at baseline commit `e398ebaf...`.
  GitHub's commit-to-PR endpoint associates neither commit with a pull request.
  Both commits had successful `contracts` and `installer` push checks, but no
  recorded independent review.
- **OBSERVED — Issue #5 and PR #6:** Issue #5 is the repository's only GitHub
  Issue. The owner comment at `2026-08-20T22:06:21Z` binds baseline
  `e398ebaf...`, branch `codex/issue-5-sdp-usage-study`, draft PR #6, and
  `SPR-SDP-005 / ITR-SDP-005-001 / SLC-SDP-005-001`. PR #6 was opened within
  seconds of the first study commits, demonstrating the newer early-draft-PR
  pattern.

## How SDP is used in this repository

### Work start

- **OBSERVED:** early canonical work started directly as a feature/agent branch
  and pull request. PRs #1–#3 have no related Issue, and the two July 17
  governance commits landed directly on `main` without a PR. Repository
  authority was therefore the branch/PR/commit itself, not a durable work
  request.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** a bounded local study and Slice
  contract appeared in commit
  `3826f1ba25077a2602a02dfd644151c608c47928` on `2026-07-13T20:56:11Z`.
  Product implementation, initial verification, and several review/remediation
  passes occurred before draft PR #4 was created on July 14. The draft PR was a
  durable CI/integration surface, but not an early assignment surface.
- **OBSERVED — current transition:** Issue #5 was written first as the durable
  contract; its active Sprint/traceability was initialized from recorded
  baseline `e398ebaf...`; the branch and draft PR were then created before broad
  repository-report work. This is the first observed workstream in this
  repository that follows Issue-authority and early-draft-PR practice.
- **INFERENCE:** work-start practice evolved in one month from PR-only/direct
  commits, through a locally defined Sprint followed by a later draft PR, to an
  Issue-first bounded assignment with early branch/PR evidence.

### Issue authority

- **OBSERVED:** no Issue existed for PRs #1–#4. PR #4's executable boundary was
  `Sprints/Sprint-001/ScrumIterations.md` and the PR body, while downstream
  acceptance arrived through the human owner and was preserved as
  `REV-SPS-001-006` rather than a GitHub Issue or review object.
- **OWNER DIRECTION:** Issue #5 makes GitHub Issues the future operational
  assignment boundary and allows one bounded Codex session to act as Master for
  one Issue.
- **OBSERVED:** the live #5 execution mirrors the Issue URL and baseline into
  `Sprints/SPR-SDP-005--recent-sdp-usage-study/sprint.md`, `Handoff.md`,
  `CurrentIndex.yaml`, `Relations.yaml`, Issue milestone comments, and PR #6.
  It initially attempted general Sprint/Slice Ledger events; the current branch
  removed them after merged validation proved that event contract unsupported.
- **INFERENCE:** Issue #5 is authoritative in current practice, but the
  repository has not yet established which mirrored fields are canonical or how
  disagreement with GitHub will be detected.

### Master, delegation, Architect, Verifier, and Reviewer

- **OBSERVED — template contract:** `AGENTS.md.template` says an undelegated
  agent acts as Master. `sdp-master` reads current coordinates, refines the
  contract, delegates one bounded implementation unit, requests fresh
  verification/review, integrates truthful state, and stops at the Slice/Fix
  boundary. `sdp-worker` implements one contract and stops; `sdp-reviewer`
  inspects the actual diff/evidence in fresh context; `sdp-architect` owns
  long-term boundaries; `sdp-verifier` checks exact-candidate evidence.
- **OBSERVED — merged execution:** PR #3 asserts a fresh Reviewer found and
  drove several remediations, and the final GitHub review comment is attached to
  exact head `e05f74dd...`. The visible GitHub actor is the repository owner, so
  GitHub proves the reviewed head and recorded disposition, not that the review
  context was independent.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** `SPS-001` records multiple fresh
  Worker/remediation and Reviewer passes. Seven review files identify exact
  candidates: reviews 001, 002, 003, and 006 required changes; 004, 005, and 007
  approved later candidates. The first review found five High and ten Medium
  issues despite passing local tests. This is strong evidence that separated,
  adversarial review improved correctness.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** Verification records are authored and
  Ledger-recorded primarily as `sdp-master`; they use exact commits, commands,
  environments, limitations, CI jobs, and supersession links. The repository
  does not demonstrate a separately identifiable Verifier agent, despite having
  a canonical Verifier skill. No execution record demonstrates an Architect
  pass; the installation study says the decision owner was the Master.
- **INFERENCE:** Worker/Reviewer separation is demonstrated most clearly in
  unmerged PR #4. Architect and Verifier separation remain principally written
  role aspirations in this repository.

### Branches and draft PRs

- **OBSERVED:** branch naming moved from `feature/...` (#1/#2), to `agent/...`
  (#3), to `codex/...` (#4/#6). There is no canonical branch naming or schema
  binding it to a stable SDP ID at the baseline.
- **OBSERVED:** PRs #1–#3 were not drafts. PR #4 remains draft behind an explicit
  downstream acceptance gate; PR #6 was opened draft at study initialization.
- **INFERENCE:** the draft state is being used as a truthful “work and evidence
  are not yet accepted” gate. The timing improved materially in Issue #5, where
  the draft PR exists from the start rather than after local completion work.

### Feature, Refactor, Fix, Sprint, Iteration, and vertical Slice

- **OBSERVED — canonical baseline:** `docs/How-SDP-Works.md` models
  `Project -> Sprint or Refactor -> Iteration -> Slice or bounded Fix`.
  `Sprints/README.md` normally permits one active Sprint, Iteration, and Slice.
  A Slice contract must state goal, why now, expected files, invariants,
  non-goals, verification, and completion signal. `TieredDesignAndImplementation.md`
  separately calls a Tier a vertical capability normally implemented by one
  Sprint or several Slices.
- **OBSERVED:** the baseline has no executed Feature, Refactor, or Fix record.
  Refactor and Fix have templates/policies and skills/schemas only. Many commits
  named `fix:` in PRs #3/#4 were remediations within larger work, not stable
  `FIX-*` workflow records.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** actual delivery used
  `Sprint-001 -> SPI-001 -> SPS-001`. One large installation-contract Slice was
  repeatedly verified, reviewed, remediated, completed, then reopened as
  revision 1 after downstream review. The append-only history retained both
  prior approvals and the newer rejection rather than rewriting them.
- **OBSERVED — current transition:** Issue #5 uses a different ID family:
  `SPR-SDP-005 -> ITR-SDP-005-001 -> SLC-SDP-005-001`, with four planned
  Iteration/Slice pairs. This shows that stable Sprint/Iteration/Slice naming is
  not yet canonical even within this repository.
- **OBSERVED — July proposal:** `SDP-PROP-001` is explicitly `Status: design
  proposal`. Its additive 1.x model keeps Sprint as implementation owner and
  places `FEAT-NNN`, a Feature backlog, and mandatory-by-default Integration
  Study before Release/Sprint scheduling. No proposed Feature/Steering folder,
  schema, backlog, or event type was implemented on `main`.
- **OWNER DIRECTION:** Issue #5 supersedes the proposal's strategic assumption:
  implementation capability must belong to Feature or Refactor; Sprint must not
  remain the conceptual owner solely because 1.x used it; a small Fix may remain
  only for genuine correction; a vertical Slice remains the bounded delivery
  unit.
- **INFERENCE:** Sprint-first delivery and the parallel Tier term are legacy
  structures needing simplification. PR #4 proves the utility of bounded Slice
  contracts and revisions, but not that Sprint should own product intent.

### Later requirements and design refinement

- **OBSERVED:** root Mandate through Implementation files are templates, not an
  executed project foundation. `docs/How-SDP-Works.md` calls them living records
  but presents a linear progression and forbids silent contradiction.
- **OBSERVED:** `InstallationContractStudy.md` in PR #4 was created after the
  Toolkit and release/versioning implementation already existed. It compared
  alternatives and decided a later installation contract inside a Sprint,
  demonstrating useful local study without rerunning the whole numbered
  lifecycle. It did not create linked Requirements, Architecture, or Design
  revisions.
- **OBSERVED:** the July proposal allows an accepted Feature Integration Study
  to update Requirements/Architecture/Design when needed, but still routes
  implementation through the old Sprint-first hierarchy.
- **OWNER DIRECTION:** Issue #5 requires the initial pass to establish a
  horizontal architectural/design skeleton while later Features may introduce
  bounded studies, requirement refinements, design refinements, and explicit
  architecture revisions.
- **INFERENCE:** actual practice already uses later local study; canonical
  lifecycle documentation has not yet explained the decision threshold between
  local design refinement and a durable architecture revision.

### Studies and convergence

- **OBSERVED:** canonical `02--Study/study.md` is one broad template with
  questions, evidence, alternatives, assumptions, risks, and recommendations.
  No canonical independent-study IDs, claim/evidence levels, convergence gate,
  or pending physical-validation state exists.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** one bounded installation-contract
  study drove one Slice. It records alternatives and an explicit decision, but
  it is not a parallel independent-study/convergence model.
- **OBSERVED — current transition:** Issue #5 and the live Sprint deliberately
  decompose study evidence into one report per repository, with synthesis
  deferred until reports exist. At the checkpoint used here, repository reports
  had begun but no completed convergence/synthesis record existed.
- **INFERENCE:** the repository demonstrates the value of bounded studies and
  delayed synthesis, but not yet a completed generic convergence gate.

### Review, verification, and handoff

- **OBSERVED — canonical contract:** `CodeReview/README.md` and the Reviewer
  skill require exact locations, actual diffs/evidence, severity, and a
  disposition. `Verification/README.md` and the Verifier skill require exact
  commands, environments, artifacts, limitations, and candidate commit.
- **OBSERVED — merged execution:** PRs #1/#2 lack recorded checks/reviews. PR #3
  has green CI and a final review comment, but no durable in-tree review or
  verification record.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** review and verification are precise,
  cumulative, adversarial, and candidate-bound. Hosted CI found failures that
  local evidence missed; later records preserve these failures and supersede
  earlier completion evidence. The final Handoff names the exact reviewed
  candidate, CI run, open draft PR, cross-repository acceptance message, next
  gate, and explicit non-authorization of merge/release.
- **OBSERVED — current transition:** Issue #5 has a minimal active Handoff with
  Issue, branch, baseline, Sprint/Iteration/Slice, next result, and stop
  boundary. Fresh independent review is planned but not yet evidence at the
  checkpoint.
- **INFERENCE:** exact-candidate evidence and a durable stop/next-gate Handoff
  materially improve recoverability. The repository does not need every
  implementation-specific adversarial test to become generic SDP; it does need
  the generic evidence shape and independence requirement.

### CurrentIndex, Relations, and Ledger

- **OBSERVED — baseline:** `CurrentIndex` is a concise current-state declaration;
  `Relations` is a manually maintained graph; `Ledger` is append-only history.
  Stable release/publication truth is represented carefully, but baseline
  validation only schemas release events and only parses the other two files as
  YAML objects.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** `Relations.yaml` grew to map one
  Sprint, one Iteration, one Slice, seven reviews, six verification records, and
  exact commits/paths/CI URLs. The Ledger grew to 22 work/review/verification
  events, including re-open and failed-CI history. This made evidence navigable
  and preserved changes of mind, but repeated the same status, path, commit, and
  relation data across Sprint, Handoff, Relations, Ledger, verification,
  reviews, and PR body.
- **OBSERVED — current contradiction and correction:** Issue #5 initially
  appended general Sprint/Slice transitions to the baseline Ledger. PR #6
  Actions run
  [32423995193](https://github.com/Hans-Einar/SDP/actions/runs/32423995193)
  failed because `release-event.schema.json` requires `releaseId` and permits
  only release event types; its Windows installer job passed. Commit
  `22b3d68bd373c592a5078fe4f3725e9921709e1b` then removed the unsupported lines
  without relabelling them as release events or changing the schema. This
  proves that current written guidance (“append work transitions”) and current
  merged executable schema disagree.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** PR #4 adds a generic Ledger envelope
  and richer CurrentIndex/Relations validation intended to resolve that class
  of problem. Because it is unmerged, the study cannot assume those contracts
  are canonical or available to Issue #5.
- **INFERENCE:** the three-part model is useful, but manual duplication and
  incomplete merged schemas make it fragile. Git/GitHub-derived facts and
  reverse links should be generated or validated rather than copied by hand.

### Steering and CurrentAssignment

- **OBSERVED:** baseline `main` has no `Steering/` directory or
  `CurrentAssignment`. Commit `e398ebaf...` only adds Architect-skill guidance
  to preserve exact Steering prompts, raw Master responses, and later
  assessments in project-owned interaction records.
- **OBSERVED — July proposal:** `SDP-PROP-001` proposes
  `Steering/CurrentAssignment.yaml`, Master Reports, Steering Decisions, and
  Blockers, but says exact schemas must be decided later. It treats GitHub as an
  operational interface while SDP remains the project-owned source of truth.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** a `gh-sdp` Steering Group disposition
  was conveyed through the repository owner, preserved truthfully in
  `REV-SPS-001-006`, and used to reopen the Slice. The record explicitly does
  not claim a GitHub review object or URL. There is no CurrentAssignment record.
- **OWNER DIRECTION:** Issue #5 places ChatGPT Steering Group above a bounded
  Issue Master, but repository/GitHub evidence—not chat—remains authoritative.
- **OBSERVED — current transition:** Issue #5's active assignment is distributed
  across the Issue, Sprint, Handoff, CurrentIndex/Relations, and PR #6; the
  release-only Ledger cannot represent it, and no separate machine-readable
  CurrentAssignment exists.
- **INFERENCE:** Steering is now real practice, but its durable assignment
  contract is not canonical. A small pointer/cache is preferable to another
  full prose authority.

### GitHub binding

- **OBSERVED:** PR #4's Sprint, Relations, Ledger, verification, reviews, and
  Handoff bind SDP IDs to exact commits, a draft PR URL, and Actions runs, but no
  Issue exists. Issue #5 binds Issue, baseline, branch, PR, and current SDP
  coordinates in prose and partial relation fields.
- **OBSERVED:** no baseline schema binds Issue number, repository identity,
  branch, PR, baseline/head, or status. No automation reconciles a closed/merged
  PR or changed branch head against `CurrentIndex`, `Relations`, or Handoff.
- **INFERENCE:** exact GitHub identities improved traceability when recorded,
  but the duplication is convention-driven and can drift. GitHub should remain
  authoritative for GitHub-native facts; SDP should store stable work identity,
  owner decisions, durable links, safety/stop boundaries, and evidence not
  derivable from GitHub.

### Automated versus conversational convention

- **OBSERVED — automated:** Toolkit/skill/manifest version agreement, Release
  and Fix schemas, release-event syntax, released-note immutability, build
  identity, installer preservation/migration, Python tests, PowerShell fixtures,
  and pull-request CI are executable at baseline.
- **OBSERVED — partly automated:** `CurrentIndex` and `Relations` must parse, but
  their semantic agreement and path graph are not validated on merged main.
  Release events are validated, while general work events are rejected.
- **OBSERVED — convention only:** Issue intake, Issue-to-work mapping, branch and
  draft-PR creation, Master/Worker delegation, fresh-context proof, Architect/
  Verifier assignment, Sprint/Slice semantics, Handoff content,
  CurrentAssignment, Steering acceptance, Feature governance, and GitHub state
  reconciliation are instructions/prose rather than executable contracts.
- **INFERENCE:** the most mature automation is release and installation
  machinery, not the product-work workflow that Issue #5 now seeks to define.

## What worked well

- **OBSERVED:** “repository records override chat memory” is consistent across
  `README.md`, `docs/How-SDP-Works.md`, Framework guidance, agent templates, and
  skills. PR #4's exact-candidate records and Issue #5's durable Issue/PR surface
  demonstrate why this helps recovery after a session boundary.
- **OBSERVED:** separating SemVer release identity from Sprint/Iteration/Slice/
  Fix/revision/Git identity prevented PRs #3/#4 from falsely claiming Toolkit
  `0.2.0` publication. No tag or GitHub Release exists, and publication fields
  remain null.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** repeated fresh review found serious
  filesystem/schema/compatibility defects that passing tests missed. Exact-head
  CI then found platform gaps. Reviews, failed evidence, remediation, and
  supersession stayed visible rather than being retrospectively cleaned up.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** the local study compared alternatives
  before changing the installer contract, and the Handoff stopped at a bounded
  downstream acceptance gate. This supported cross-repository coordination
  without modifying `gh-sdp` opportunistically.
- **OBSERVED — current transition:** Issue #5's Issue-first baseline, immediate
  draft PR, separate repository reports, and explicit non-goals provide a more
  recoverable and bounded assignment than the earlier PR-only workstreams.
- **INFERENCE:** the strongest reusable practices are exact work contracts,
  truthful state/publication, small vertical outcomes, fresh adversarial review,
  exact verification, append-only correction history, and explicit handoff/stop
  gates—not the old hierarchy that happens to contain them.

## Pain points and accidental complexity

- **OBSERVED:** canonical `main` describes a complete method but contains almost
  no executed lifecycle/Sprint/review/verification examples. The strongest
  example is isolated in unmerged PR #4, so consumers of `main` receive
  aspiration without the tested operational records that refined it.
- **OBSERVED:** the same root tree is both live Toolkit project state and an
  installation source. At baseline the installer copies live
  `Traceability/CurrentIndex.yaml` and `Relations.yaml`; optional initialization
  recursively copies root operating folders. PR #4's study identified live
  Toolkit record leakage into a consuming project as a real downstream problem.
- **OBSERVED:** duplicate `skills/`, `payload/`, and installer entry points make
  authority less obvious. The root skills are stale and incomplete relative to
  `Toolkit/skills/`; only the root installer explicitly identifies itself as a
  compatibility forwarder.
- **OBSERVED:** Sprint/Iteration/Slice IDs changed from
  `Sprint-001/SPI-001/SPS-001` to
  `SPR-SDP-005/ITR-SDP-005-001/SLC-SDP-005-001`, while the July proposal uses
  examples such as `SPR-004/ITR-004/SLC-012`. The IDs carry ceremony without a
  stable, validated format.
- **OBSERVED — PR #4, UNMERGED EVIDENCE:** one “smallest” Slice included a major
  portable installation contract, schema suite, installer refactor, project
  templates, validator modes, extensive tests, documentation, seven reviews,
  six verification records, CI remediation, and downstream conformance revision.
  The boundary was coherent, but very large for a single Slice.
- **OBSERVED:** CurrentIndex, Relations, Ledger, Sprint, Handoff, review,
  verification, PR body, and Issue comment repeat current coordinates and
  evidence. PR #6's schema failure shows that this manual state can be both
  duplicated and invalid under merged automation.
- **OBSERVED:** review independence is asserted in records but not independently
  attestable from the GitHub actor history. PR #4 has no GitHub review objects;
  PR #3's review is owner-authored. This does not invalidate the review content,
  but limits what external evidence proves.
- **OBSERVED:** the July Feature proposal and Steering skill change landed
  directly on `main` with green push checks but without a PR or recorded
  independent review, despite governing the workflow that requires review for
  material work.
- **OBSERVED:** Toolkit `0.2.0` and all versioned skills remain unreleased, while
  the installation-contract hardening needed by downstream clients remains in
  a month-old draft PR #4. Canonical and demonstrated-next contracts are split.
- **INFERENCE:** the method over-invests in manually repeating coordinates while
  under-specifying the high-value assignment boundary, Feature/Refactor
  ownership, current GitHub binding, and evolution after initial design.

## Direct contradictions and limitations

- **OBSERVED:** the studied baseline says work transitions belong in the Ledger,
  but its only Ledger schema permits release events. The live Issue #5 work
  followed the prose, failed CI, and removed the unsupported events. The
  unmerged PR #4 solution cannot be retroactively treated as canonical.
- **OBSERVED:** `docs/How-SDP-Works.md` presents Sprint/Refactor as delivery
  owner; July's proposal deliberately retains that hierarchy; Issue #5 later
  directs a Feature/Refactor-first model. These are chronological changes, not
  equivalent descriptions to normalize.
- **OBSERVED:** baseline Architect guidance prefers repository Steering
  interaction logs, but no Steering directory/schema/example exists and current
  Issue #5 state uses Issue/Sprint/Handoff/traceability instead.
- **OBSERVED:** canonical instructions require agent entry files to be read, but
  the upstream repository contains only templates. Externally supplied workspace
  instructions governed this study; that external context is not proof of what
  a fresh clone alone instructs.
- **OBSERVED:** no completed Feature or Refactor execution exists in this
  repository, and no completed multi-study convergence exists at the current
  checkpoint. Feature/Refactor and convergence recommendations therefore rely
  on owner direction, structural gaps, and other repository reports—not a
  successful local example.
- **OBSERVED:** GitHub proves commits, PR states, review objects, checks, and
  dates, but it cannot prove “fresh context” from the same owner identity. PR #4
  repository records assert fresh Worker/Reviewer passes and contain detailed
  evidence; they remain unmerged and have no GitHub review object.
- **OBSERVED:** Issue #5 and PR #6 began after the fixed activity window and
  after the canonical study commit. They are included only because the
  assignment requires current repository state and newer owner direction.
- **INFERENCE:** this repository is a high-value source for principles,
  automation boundaries, review discipline, and the transition to Issue-first
  work. It is weak evidence that the old Sprint-first hierarchy, proposed
  Feature schema, or Steering record layout has worked canonically in practice.

## Carry forward

- **RECOMMENDATION:** retain repository/GitHub evidence as authority, with the
  GitHub Issue as the bounded operational assignment and an early draft PR as
  the exact diff/CI/review surface.
- **RECOMMENDATION:** retain one Master for the bounded Issue, fresh Workers for
  implementation, a separate fresh Reviewer, optional Architect involvement
  when horizontal contracts change, and a separately checkable Verifier role or
  verification pass. Record exact candidate identities; do not rely on role
  summaries alone.
- **RECOMMENDATION:** retain the initial horizontal architecture/design skeleton
  and vertical Slice rule, but make Feature or Refactor own capability intent.
  Treat Sprint/Iteration as optional execution/timebox grouping, not mandatory
  product hierarchy.
- **RECOMMENDATION:** retain bounded studies, explicit alternatives/decisions,
  and convergence before implementation when evidence streams are independent.
  Allow these studies and requirement/design refinements inside later Features,
  long after project bootstrap.
- **RECOMMENDATION:** retain exact verification commands/environments/artifacts,
  fresh severity-based review, append-only correction/reopen events, truthful
  publication identities, and Handoff records naming stop boundary, residual
  limitations, exact head, and next authorization gate.
- **RECOMMENDATION:** retain `CurrentIndex` as a small declared current-state
  cache, `Relations` as the stable non-derivable graph, and `Ledger` as
  append-only transition history. Validate them together and generate GitHub/
  Git facts instead of copying them into every record.
- **RECOMMENDATION:** provide one canonical managed asset tree under `Toolkit/`
  and neutral project templates physically separated from this repository's
  live records. Keep only explicit, tested migration forwarders during a
  deprecation window.
- **RECOMMENDATION:** preserve the generic review lesson from PR #4—adversarial
  attempts and exact-head multi-environment CI—while leaving installer-specific
  path/link/hash safety mechanics in the installer domain rather than generic
  SDP workflow.

## Legacy / do not carry forward

- **RECOMMENDATION:** do not carry the numbered lifecycle progression forward as
  a one-time waterfall implementation decomposition. Keep the records as a
  living project foundation; replace “finish Design, then implement the original
  plan” with Feature/Refactor-local evolution through vertical Slices.
- **RECOMMENDATION:** do not carry Sprint as the conceptual owner of a capability
  or preserve both Tier and Slice as competing vertical-delivery identities.
  Feature/Refactor should own intent; Slice should be the smallest verifiable
  vertical increment; optional scheduling concepts should not add stable IDs
  unless they support a real decision.
- **RECOMMENDATION:** do not carry forward Issue-less material work, direct-main
  governance changes, or draft PRs opened only after most implementation/review.
  Issue #5's newer Issue-first and early-draft-PR execution replaces those
  patterns.
- **RECOMMENDATION:** do not treat July proposal details—mandatory Integration
  Study for every Feature, Sprint-first implementation ownership, exact
  `FeatureBacklog.yaml` fields, or full raw-chat logging—as accepted canonical
  behavior. They are proposal evidence. Re-evaluate them against Issue #5's
  newer owner direction, proportionality, privacy, and derivability.
- **RECOMMENDATION:** do not copy live repository state as installation
  templates, retain duplicate stale root skills/payload as coequal authority, or
  depend on unmerged PR #4 contracts. Neutral templates and one versioned
  Toolkit authority must replace the ambiguity after a separately reviewed
  implementation/migration issue.
- **RECOMMENDATION:** do not hand-maintain every branch/head/check/status/path in
  CurrentIndex, Relations, Ledger, Handoff, Issue, and PR prose. Let GitHub,
  `gh-sdp`, or SDP-Analyzer derive GitHub-native facts; persist only stable IDs,
  decisions, safety boundaries, evidence references, and historical transitions.
- **RECOMMENDATION:** do not infer independent review merely from a prose claim
  or reuse a prior approval after the candidate changes. Preserve PR #4's
  exact-candidate/supersession discipline and improve reviewer provenance where
  practical.
