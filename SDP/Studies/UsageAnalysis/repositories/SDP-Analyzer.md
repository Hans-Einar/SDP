# SDP-Analyzer SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/SDP-Analyzer` at the exact
  current default-branch identity
  [`main@632991a878100e8cd8c4efbb7d724edb3694d98a`](https://github.com/Hans-Einar/SDP-Analyzer/commit/632991a878100e8cd8c4efbb7d724edb3694d98a),
  committed `2026-07-17T12:59:20Z`. GitHub's advertised default-branch head
  and the local `origin/main` object agreed. Files were read from that Git
  object; the studied repository was not checked out, edited, committed, or
  pushed.
- **OBSERVED:** Live GitHub repository, Issue, pull-request, review, branch,
  commit-status, check-run, workflow, release, tag, and head-commit-comment
  state was inspected on `2026-08-21`. GitHub returned zero Issues, zero pull
  requests, zero Actions workflows, zero check runs/status contexts on the
  studied head, zero releases, zero tags, and no head-commit comments. The only
  advertised branch was unprotected `main`.
- **OBSERVED:** The local working checkout was at the older clean product
  commit `48388cb40faedc4391bc42aa1758ce6ef1cbf75e`, two commits behind
  `origin/main`. The two newer default-tree commits add only
  `SDP/Steering/README.md` and
  `SDP/Steering/Interactions/STR-2026-001--feature-governance-and-slc-009.md`.
  This report does not mistake the local checkout's older view for current
  default-tree state.
- **OBSERVED:** No product tests were run for this read-only study. Statements
  about implemented behavior come from current source and permanent tests;
  command results come from committed `VER-*`/`REV-*` records and are identified
  as recorded evidence rather than independently rerun evidence.
- **LIMITATION / OBSERVED:** Fresh Worker, Reviewer, Master, and supervising
  Architect context is asserted in committed SDP records, but Git/GitHub cannot
  independently identify those sessions. All 30 commits are authored and
  committed by GitHub user `Hans-Einar`; there are no PR reviews or review
  threads to corroborate distinct identities.
- **LIMITATION / OBSERVED:** The installed SDP generation is not versioned by a
  manifest. There is no `SDP.manifest.yaml`, installed-toolkit manifest, or
  project manifest. The repository self-identifies its implemented Analyzer
  input as `sdp-toolkit-structured-core-v1`; that is an Analyzer compatibility
  profile, not proof of the installed Toolkit release.

## Repository facts

- **OBSERVED:** Repository: [`Hans-Einar/SDP-Analyzer`](https://github.com/Hans-Einar/SDP-Analyzer)
  (public); default branch: `main`; study commit:
  `632991a878100e8cd8c4efbb7d724edb3694d98a`.
- **OBSERVED:** All 30 commits reachable from `main` fall inside Issue #5's
  fixed window (`2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`). The
  first is `122407e4f0cea8c32305d2fc3e7910d66e30422f` at
  `2026-07-11T13:11:10Z`; the latest is the studied head at
  `2026-07-17T12:59:20Z`. All are single-parent or root commits; there are no
  merge commits.
- **OBSERVED:** SDP is extensively installed at root `AGENTS.md`,
  `AGENTS-project.md`, `.codex/skills/`, and `SDP/`. The lifecycle is
  `SDP/01--Mandate/` through `SDP/07--Implementation/`. Operating material is
  under `SDP/Instructions/`, `SDP/Sprints/`, `SDP/Traceability/`,
  `SDP/CodeReview/`, `SDP/Verification/`, and `SDP/Steering/`.
- **OBSERVED:** `SDP/Features/` does not exist. `SDP/Refactors/` contains only
  the Toolkit-style README and no active or historical Refactor record.
  `SDP/Fixes/`, repository-local release records, and a machine-readable
  `CurrentAssignment` do not exist.
- **OBSERVED:** `SDP/Sprints/SPR-001/` contains completed Tier 1 contracts,
  notes, and handoff. `SDP/Sprints/SPR-002/` contains active Tier 2 contracts,
  notes, and handoff. Review records are `REV-SLC-001` through `REV-SLC-008`
  plus `REV-TIER-001`; verification records are `VER-SLC-001` through
  `VER-SLC-008`, `VER-TIER-001`, and `verification-plan.md`.
- **OBSERVED:** `SDP/Traceability/CurrentIndex.yaml` declares project
  `SDP-ANALYZER` in `active-development`, Tier `TIER-002`, active Sprint
  `SPR-002`, active Iteration `ITR-002`, and Slice `SLC-008`, with no active
  Refactor. The Slice is explicitly completed and intentionally retained as
  the stop boundary. `SLC-009` and `SLC-010` are planned.
- **OBSERVED:** `Relations.yaml` is a 300-line registry for lifecycle documents,
  three Tiers, two Sprints, two Iterations, Slices `SLC-001` through
  `SLC-010`, review records, and verification records. It has no Feature,
  GitHub Issue, branch, PR, commit, Fix, or release section.
- **OBSERVED:** `Ledger.ndjson` contains 75 parseable append-only records. Its
  last event is `EVT-2026-07-14-075`, completing `SLC-008`; Ledger timestamps
  are dates such as `2026-07-14`, not precise instants. The Steering commits of
  `2026-07-17` have no corresponding Ledger event.
- **OBSERVED:** There are no GitHub Issues or PRs to list as actual Analyzer
  work authorities. There are therefore also no GitHub Issue comments, PR
  comments, formal PR reviews, or PR checks. The current head's combined
  status is `pending` solely because it has zero status contexts, not because a
  check is running.
- **OBSERVED:** `main` has no branch protection. GitHub exposes no Actions
  workflow. Verification was executed and recorded locally in SDP records,
  rather than enforced by GitHub CI.
- **OBSERVED:** The current application is Vite/React/TypeScript `0.1.0` with
  `SharedUI@0.1.0` and `yaml@2.9.0`. Its implemented product boundary is a
  deterministic read-only structured-core Analyzer plus a browser-directory
  source adapter that is not wired into the UI.

## How SDP is used

### Work start and authority

- **OBSERVED:** Work started repository-first, not Issue-first. The first 20
  commits on `2026-07-11` established the SDP installation, then successively
  authored Mandate, one broad Study, Requirements, Architecture, Design
  Analysis, Design, Implementation plan, Sprint/Iteration/Slice contracts,
  verification plan, CurrentIndex, Relations, Ledger, handoff, and a
  repository-specific SharedUI policy before the first product-bearing commit.
- **OBSERVED:** The project Mandate describes a supervising ChatGPT
  conversation as responsible for long-range product/architecture oversight,
  while Codex acts as repository-local Master. That supervision was a
  conversational convention during implemented work: no GitHub Issue records
  it.
- **OBSERVED:** The first durable Steering assignment appears only at the
  current head in `STR-2026-001`. It records an exact prompt, expected base
  `48388cb40faedc4391bc42aa1758ce6ef1cbf75e`, target role, Tier/Sprint/
  Iteration/next-Slice coordinates, and strict stop boundaries. Its raw Master
  response and Steering assessment remain `pending`.
- **INFERENCE:** The repository evolved from chat-supervised, repository-local
  execution toward durable Steering interaction records, but it did not reach
  Issue-authoritative assignment or a machine-readable CurrentAssignment
  before work stopped.
- **OWNER DIRECTION:** Issue #5 requires GitHub Issues to become the future
  operational assignment boundary and permits a bounded Codex session to act
  as Master for one Issue. The Analyzer repository's older non-Issue start is
  evidence to migrate from, not a reason to reject that direction.
- **RECOMMENDATION:** Future Analyzer work should begin from a GitHub Issue and
  mirror only durable machine-readable coordinates into a Steering assignment
  record. A copied raw prompt may remain supporting evidence, but should not be
  the sole assignment identity or the only place branch/PR/stop conditions are
  expressed.

### Master role and delegation

- **OBSERVED:** Root `AGENTS.md` makes any non-delegated agent the Master. The
  prescribed read order is root instructions, project instructions, reminders,
  Framework, the matching role skill, then active SDP/traceability.
- **OBSERVED:** The installed Master skill requires the Master to identify
  active Sprint or Refactor, Iteration, and Slice; refine SDP before
  implementation; delegate bounded product-code work to a Worker; delegate
  review to a fresh Reviewer; inspect real evidence; update notes,
  verification, Relations, CurrentIndex, Ledger, and handoff; then stop at the
  Slice boundary.
- **OBSERVED:** Worker, Reviewer, Architect, traceability, and vertical-refactor
  skills are separate. The Worker implements one assigned Slice; Reviewer
  inspects independently and does not fix; Architect changes long-term
  contracts and does not implement product code.
- **OBSERVED:** There is no installed `sdp-verifier` skill and root `AGENTS.md`
  does not name Verifier as a separate role. The Master is required to inspect
  real verification evidence, while `VER-*` records describe the verification
  function. A separately identifiable Verifier session is therefore not part of
  this repository's demonstrated role model.
- **OBSERVED:** Committed review records document real changes-required cycles.
  `REV-SLC-002` required source-list failure correction; `REV-SLC-003` found
  two parser defects; `REV-SLC-005` found missing regressions and incomplete
  SDP007 provenance; `REV-SLC-006` went through multiple correction/re-review
  cycles, including reopening a previously completed Slice. The append-only
  Ledger preserves these transitions.
- **OBSERVED:** `SLC-007` was blocked at commit
  `4e2018a45614f76a3ed207272a86badc366d1fac` when Tier acceptance contradicted
  Markdown and verification semantics. A supervising Architect decision added
  `DEC-STU-015` and `DEC-STU-016`, unblocked it, and completed Tier 1 at
  `afb1d97cf537d80cb3ff7b84b17a807a591b104e`.
- **LIMITATION / INFERENCE:** The correction history supports the value of a
  separate Reviewer perspective, but fresh identity and Worker/Reviewer
  separation remain process attestations. All product, review, verification,
  and closure files ultimately enter the same linear `main` history under
  `Hans-Einar`.
- **RECOMMENDATION:** Carry the role separation forward, but bind each Worker,
  Verifier, and Reviewer pass to an exact branch/PR head and record the
  attestation/session identity or equivalent evidence. A changed head should
  require a fresh review disposition rather than relying only on prose that an
  uncommitted tree was rechecked.

### Branches, commits, draft PRs, reviews, and checks

- **OBSERVED:** Actual Analyzer delivery did not use branches or PRs. The 30
  commits form one linear `main` history, including all product Slices,
  verification/review records, and Steering documentation.
- **OBSERVED:** `AGENTS-project.md` says no pull requests or commits unless the
  human explicitly requests them. `STR-2026-001` likewise says not to create a
  PR and not to commit or push without explicit human instruction. This is a
  project-local authorization boundary, not an early-draft-PR workflow.
- **OBSERVED:** Product completion commits combine product changes, tests,
  review/verification documents, traceability, and handoff. For example,
  `48388cb40faedc4391bc42aa1758ce6ef1cbf75e` adds the browser adapter and all
  `SLC-008` evidence in one commit.
- **OBSERVED:** `REV-SLC-008` reviewed accepted Tier 1 commit `afb1d97...` plus
  an uncommitted tree. The later immutable product/evidence commit is
  `48388cb...`; no PR or check suite establishes that exact commit as the
  separately reviewed candidate.
- **INFERENCE:** The repository achieved short feedback loops, but direct-main
  integration reduced the durable independence and exact-head evidence that a
  draft PR would provide.
- **OWNER DIRECTION:** Issue #5 mandates branch plus early draft PR for the
  future Issue-driven operating model.
- **RECOMMENDATION:** Use one Issue-bound branch and early draft PR. Record base,
  current head, reviewed head, check-run identities, merge commit, and any
  post-merge reconciliation separately. Do not treat a local uncommitted-tree
  review as equivalent to an immutable reviewed PR head.

### Feature, Refactor, Fix, Sprint, Iteration, and Slice

- **OBSERVED:** Implemented hierarchy is
  `Project -> Tier -> Sprint -> Iteration -> Slice`. `SPR-001/ITR-001` contains
  `SLC-001` through `SLC-007`; `SPR-002/ITR-002` contains `SLC-008` through
  planned `SLC-010`. There is one Iteration per Sprint.
- **OBSERVED:** Tiers are long-range capability owners: structured-core
  analysis, local acquisition, Markdown coverage, automation, traceability
  exploration/graph, and repair assistance. Sprints are execution containers
  for a Tier, rather than time-boxes independent of the capability.
- **OBSERVED:** Each Slice contract is unusually complete: goal, why now,
  linked requirements/architecture/design, required files/modules, invariants,
  explicit non-goals, tests, verification, independent review, discoveries
  policy, completion signal, and stop condition.
- **OBSERVED:** Slices preserve a buildable path and are bounded, but Tier 1's
  sequence is substantially component/layer staged: shell, discovery, parsers,
  normalization, rule engine, UI, then integrated acceptance. It becomes an
  end-to-end user capability only after several Slices. `SLC-008` is more
  vertical in testing the browser adapter through the existing full analysis
  pipeline, while intentionally withholding UI.
- **INFERENCE:** This shows the benefit of small gates but also the risk of
  calling every bounded layer increment a vertical Slice. Some Analyzer Slices
  are architecture-enabling increments, not independently user-visible
  end-to-end capabilities.
- **OBSERVED:** No Feature record is implemented. The committed pending
  `STR-2026-001` prompt instructs a future Master to allocate a provisional
  `Project Governance Visualization` Feature and Integration Study before
  implementing `SLC-009`. It explicitly says the upstream SDP Feature contract
  is not accepted and must not be invented locally as final schema.
- **OBSERVED:** No Refactor or Fix was executed. Although the normalized code
  accepts a `refactors` relation section, current Relations has none and the UI
  offers no Refactor view.
- **OWNER DIRECTION:** Future SDP makes Feature and Refactor the implementation
  owners; a Fix is reserved for a genuine bounded correction. Sprint must not
  remain capability owner merely because this repository used Tiers/Sprints.
- **RECOMMENDATION:** Reframe future Analyzer capability delivery as
  `Feature/Refactor -> optional execution grouping -> vertical Slice(s)`. Keep
  Tier only as a product roadmap/maturity label if useful, not the traceability
  owner. Distinguish a user-capability Slice from a prerequisite/enabler Slice
  explicitly rather than stretching “vertical” to cover both.

### Initial design and later requirements/design

- **OBSERVED:** The repository is a strong example of initial horizontal
  design. `ARC-001`, `DAN-001`, and `DES-001` establish acquisition, parsing,
  profile/compatibility, normalized domain, validation, findings, application,
  React, SharedUI, report/export, fixture/test, security, provenance, and state
  ownership boundaries before product work.
- **OBSERVED:** Product code follows those boundaries: `src/core` owns pure
  source/discovery/parsing/normalization/validation models;
  `src/application` orchestrates; `src/adapters` supplies fixture/browser
  sources; `src/ui` owns React/SharedUI. The Tier 2 browser adapter reuses the
  existing `ProjectSource` path without putting browser types into core.
- **OBSERVED:** Planning is not completely frozen. The `SLC-007` block caused
  Study, Requirements, Architecture, Design Analysis, Design, and
  Implementation documents to be amended with structured-core and verification
  qualification decisions. This is truthful later refinement in response to
  implementation evidence.
- **OBSERVED:** The pending Steering prompt proposes a Feature-local Integration
  Study for governance visualization, declared-vs-observed lifecycle, Feature
  compatibility, and future Slice fan-out after the initial architecture was
  already accepted.
- **INFERENCE:** Analyzer demonstrates both sides of the desired model: a useful
  initial horizontal skeleton and evidence-driven later design refinement.
  What it lacks is a canonical Feature owner for that later refinement.
- **RECOMMENDATION:** Preserve the horizontal analysis-core boundaries. Let a
  later Feature own new requirements, an Integration Study, and design
  refinements. Require an architecture revision only when dependency direction,
  source authority, normalized facts, or runtime boundaries change; keep
  ordinary view/query design local to the Feature/Slice.

### Studies and decisions

- **OBSERVED:** `STU-001` is one broad initial Study with 16 stable decisions.
  It covers profile, formats, parsers, acquisition, security, provenance,
  normalized model, findings, determinism, validation, verification,
  graphing, SharedUI, packaging, structured-core boundary, and qualification.
- **OBSERVED:** It was amended, not replaced, when later implementation exposed
  a contradiction. Earlier decisions remain visible and `DEC-STU-015`
  explicitly states which Tier 1 implications it narrows/supersedes.
- **OBSERVED:** There are no multiple independent Study records, parallel study
  assignments, evidence/claim levels, or convergence gate. `STR-2026-001`
  proposes one future Feature-local Integration Study, but it remains an
  unexecuted prompt.
- **RECOMMENDATION:** Carry forward stable decision IDs and explicit
  supersession/narrowing. Support one coherent Study by default, Feature-local
  Studies later, and multiple independent Studies/convergence only where the
  uncertainty can actually be split. Do not impose ActionCam-style multi-study
  machinery on every Analyzer Feature.

### Independent verification and review

- **OBSERVED:** Verification is Slice-specific and concrete. Records include
  exact commands, environment, test counts, focused and full suites, build
  output, dependency identities, scope scans, limitations, residual risk, and
  traceability checks. `VER-SLC-008` records 23 files/183 tests, 34 focused
  browser-adapter tests, a 2,070-module build, and exact dependency versions.
- **OBSERVED:** Review is not a ceremonial approval. Earlier dispositions and
  failed/incorrect evidence harness attempts are preserved, and the Ledger
  records review-required, reopened, corrected, re-reviewed, and completed
  transitions. `SLC-006` is a particularly strong recovery example.
- **OBSERVED:** The Analyzer's own `SDP007` rule embodies a useful evidence
  rule: a completed Slice qualifies only when an explicit `verification`
  relation resolves to a verification entity with exact `outcome: passed` and
  a non-empty `check` or `command`. Plan/review/summary prose does not qualify.
- **LIMITATION / OBSERVED:** The rule does not read verification Markdown or
  execute the command. The repository has no CI checks and no formal GitHub
  review. Review records frequently describe an uncommitted tree, so exact
  immutable review identity remains weaker than their content quality.
- **RECOMMENDATION:** Preserve separate verification facts and review
  dispositions, real correction history, and no pre-claiming. Add immutable
  branch/PR head, reviewer attestation, check-run URL/ID, and post-change
  re-review. Analyzer should report “recorded passed evidence,” not “the code is
  correct,” unless an explicitly trusted verifier has rerun it.

### Handoff and recoverability

- **OBSERVED:** Each Sprint has `Handoff.md`. The current handoff states accepted
  Tier 1 baseline, current completed Slice, policies, planned work, unallocated
  requirements, evidence, and the exact no-continue boundary. It is compact
  enough for re-anchoring compared with the much larger contracts.
- **OBSERVED:** `STR-2026-001` begins with a detailed re-anchor list and requires
  verifying committed state before doing anything. This directly addresses
  context/session loss.
- **PAIN / OBSERVED:** Handoff, CurrentIndex, Relations, Sprint status, README,
  and Ledger each repeat overlapping state manually. The Steering record then
  adds a second assignment/re-anchor surface. Their agreement is maintained by
  discipline, not generation.
- **RECOMMENDATION:** Keep a short prose handoff for discoveries, risks, and next
  decision. Generate coordinates, Issue/branch/PR/head/check/review/status from
  machine-readable sources so re-anchoring does not depend on reconciling six
  manual copies.

### CurrentIndex, Relations, Ledger, and Steering assignment

- **OBSERVED:** The three-file structured core provides distinct useful claims:
  CurrentIndex is declared present state; Relations is the explicit entity/link
  registry; Ledger is append-only claimed history. The code deliberately does
  not infer missing entities from relations or Ledger prose.
- **OBSERVED:** The current state exposes a semantic ambiguity: the YAML key is
  `active.slice`, but its value `SLC-008` is completed and retained as a stop
  boundary while Tier/Sprint/Iteration remain active. This is documented and
  therefore not simply a mistake, but a consumer cannot distinguish “current
  active implementation” from “last completed accepted boundary” from the
  field alone.
- **OBSERVED:** `STR-2026-001` gives a durable human-readable assignment but no
  `CurrentAssignment.yaml`. Its response/assessment are pending; no Ledger
  event connects it to the 75-event lifecycle; no Issue/branch/PR identity is
  available.
- **OBSERVED:** Relations repeats long requirement, architecture, decision,
  verification, and review lists on Tiers, Sprints, Slices, verification, and
  reviews. The manual graph is rich but expensive to maintain.
- **INFERENCE:** CurrentIndex is valuable as a fast pointer, but it should not
  also be an ambiguous lifecycle conclusion. Relations should own stable graph
  edges once; current GitHub/runtime facts and derived reverse edges should be
  generated.
- **RECOMMENDATION:** Separate at least `currentAssignment`, `activeExecution`,
  and `lastAcceptedBoundary`. A machine-readable assignment should bind Issue,
  Feature/Refactor/Fix, base, branch, draft PR, Master role, stop condition,
  and expected verification/review without copying every issue paragraph.
  Generate reverse/duplicated links and GitHub state; keep the Ledger append-only
  for transitions that are not already immutable GitHub events.

### GitHub binding and automation

- **OBSERVED:** No Analyzer SDP ID is bound to a GitHub Issue, branch, PR,
  check, tag, or release. There is no actual GitHub automation. The only GitHub
  repository identity in the planning chain is prose such as
  `Repository: Hans-Einar/SDP-Analyzer`.
- **OBSERVED:** Product automation is strong and deterministic at the local
  test level. Architecture boundary tests scan for prohibited imports and
  graph/write-back behavior; parser and normalization tests exercise malformed
  inputs; UI tests exercise loading, failure, filters, provenance, and stale
  selection clearing.
- **INFERENCE:** Analyzer distinguishes deterministic product automation from
  workflow automation well. The former is implemented; the latter is absent.
- **RECOMMENDATION:** Future `gh-sdp` or a GitHub integration should populate
  Issue/PR/head/check facts and run Analyzer validation in CI. Analyzer should
  consume a versioned, timestamped GitHub evidence snapshot or an explicit live
  connector; it should never silently treat stale cached GitHub data as current.

## Implemented Analyzer behavior relevant to future SDP formats

### Default-tree current behavior

- **OBSERVED:** Discovery recognizes exact core paths
  `SDP/Traceability/CurrentIndex.yaml`, `Relations.yaml`, and `Ledger.ndjson`;
  it detects the seven lifecycle directories plus Sprints, Verification,
  CodeReview, and Traceability by canonical path prefix. It does not discover
  Features, Steering, Refactors, Fixes, Releases, or GitHub bindings as standard
  directories.
- **OBSERVED:** CurrentIndex normalization exposes project `id`, `name`,
  `status`, `tier`, and active `sprint`, `refactor`, `iteration`, and `slice`.
  There is no active Feature, Fix, Issue, assignment, PR, or release field.
- **OBSERVED:** Relations parsing preserves arbitrary top-level YAML, but
  normalization only accepts `documents`, `tiers`, `sprints`, `iterations`,
  `slices`, `reviews`, `verification`, and `refactors`. A future `features`,
  `issues`, `pull_requests`, `commits`, `fixes`, or `releases` section would be
  diagnosed as unsupported and skipped, not normalized as graph entities.
- **OBSERVED:** Recognized entity kinds are mandate, study, requirement,
  architecture-decision, design-decision, tier, sprint, iteration, slice,
  verification, review, and unknown. The TypeScript kind is extensible, but no
  supported section currently creates Feature, Issue, PR, commit, Fix, or
  release entities.
- **OBSERVED:** Supported relation fields are `derives_from`, `decisions`,
  `tier`, `sprint`, `iteration`, `slice`, `slices`, `requirements`,
  `architecture`, `study_decisions`, `design`, `verification_plan`,
  `verification`, and `review`. Authority, branch, PR, commit, supersession,
  release inclusion, and removal are not normalized.
- **OBSERVED:** Ledger normalization preserves ordered events, event type,
  subject ID, timestamp string, payload, and source. It does not reconstruct an
  observed lifecycle or compare that lifecycle to declared entity status.
- **OBSERVED:** Rules `SDP001` through `SDP008` cover core source availability,
  duplicate structured definitions, dangling relation endpoints, malformed
  Ledger lines, unresolved active declarations, explicit contradictory active
  hierarchy, completed-Slice verification qualification, and profile support.
  There is no staleness, orphaned Slice, Issue mismatch, PR/head mismatch,
  declared-vs-observed status, release inclusion, review freshness, or complete
  Requirement-to-verification path rule.
- **OBSERVED:** The UI offers bundled clean/broken fixture selection, project and
  declared-active summary, diagnostics, findings filters, and detailed
  provenance. It has no live folder picker, entity table/navigation, relation
  exploration, graph/timeline, current-assignment view, GitHub state, or release
  view. The browser directory adapter exists but is not connected to the UI.
- **OBSERVED:** Markdown content is not read. Therefore current Analyzer cannot
  discover IDs/statuses/relations found only in Feature, Study, Slice,
  verification, review, handoff, Steering, or CurrentAssignment Markdown.

### Aspirational default-tree documentation

- **OBSERVED:** The Mandate aspires to navigate requirements and decisions
  through implementation/verification; detect stale/forgotten work; explore
  verification by Slice, requirement, commit, or date; compare current index,
  relation graph, and Ledger; and later support graphing, reports, CLI/CI, and
  skills.
- **OBSERVED:** Requirements already reserve `REQ-V-009` for stale detection
  using explicit status/timestamp policy, `REQ-UI-005` for relation navigation,
  and `REQ-UI-006` for a derived graph. `REQ-D-003` reserves Markdown coverage;
  `REQ-C-005` reserves versioned machine reports.
- **OBSERVED:** `DAN-001` orders later work as TIER-003 Markdown/document
  coverage, TIER-004 reports/CLI/CI, TIER-005 relation navigation/derived graph,
  and TIER-006 repair. `REQ-V-009` and some T2-tagged requirements remain
  explicitly unallocated before Tier 2 closure.
- **OBSERVED:** `STR-2026-001` goes further by proposing an Analyzer 1.x Product
  Governance Visualization Feature: Release/Feature/Sprint/Iteration/Slice/
  Refactor/Fix listings, declared and Ledger-observed status shown separately,
  a Ledger-timestamp Gantt, item detail with relations/events/findings/
  verification/review/provenance, optional future Feature records, and backward
  compatibility for repositories with no Features. It defers raw Steering
  interaction browsing to Analyzer 2.0.
- **INFERENCE:** These are credible road-map inputs because they are committed,
  but they are not accepted requirements or implemented behavior. The Steering
  interaction remains `prompt-issued` with pending response and assessment.

### Open or unmerged evidence

- **OBSERVED:** There is no open or unmerged Analyzer Issue/PR evidence. No
  branch exists besides `main` in the live GitHub repository.
- **OBSERVED:** The only “open” evidence is repository-local planning state:
  active Tier/Sprint/Iteration, completed retained `SLC-008`, planned SLC-009/
  SLC-010, and pending `STR-2026-001` response/assessment. These records are on
  default `main`, not an unmerged branch.
- **RECOMMENDATION:** Cross-repository synthesis should not describe
  `Project Governance Visualization` as an accepted Feature or claim Feature
  support exists. It is a useful owner/Steering proposal awaiting execution and
  upstream SDP schema decisions.

## Future Analyzer graph, status, staleness, and legacy-format assessment

### Issue-to-release graph

- **RECOMMENDATION:** The future Analyzer graph should support the candidate
  path `Issue <-> Feature/Refactor <-> Study/requirements/decisions <-> Slice
  <-> PR/commit <-> verification/review <-> release`; Issue #5 asks the study
  to evaluate this shape rather than mandating its exact schema.
- **RECOMMENDATION:** Extend the serializable normalized domain, not a UI graph
  library, with first-class entities for GitHub Issue/Issue comment authority,
  Feature, Refactor, Fix, Study, Requirement, architecture/design decision,
  optional execution group/Sprint, Iteration when retained, Slice, branch, PR,
  commit/tree, verification run/record, review/disposition, release record, Git
  tag, and GitHub Release.
- **RECOMMENDATION:** Add explicit directed relations such as `authorizes`,
  `implements`, `refines`, `depends_on`, `supersedes`, `removes`, `executed_by`,
  `delivered_by`, `head_at`, `merged_as`, `verified_by`, `reviewed_by`,
  `includes`, and `published_as`. Do not overload generic `derives_from` or
  infer Issue/PR identity from matching strings.
- **RECOMMENDATION:** Preserve source class and observation time on every fact:
  repository declaration, append-only repository event, Git object, live or
  captured GitHub fact, generated derivation, or Markdown extraction. A graph
  edge must retain provenance and compatibility-profile ownership.
- **RECOMMENDATION:** Add path queries that can answer Requirement -> Feature/
  Refactor -> Slice -> exact PR head/commit -> qualifying verification -> fresh
  review -> release. Missing links should yield explainable findings, not
  invented intermediate nodes.
- **RECOMMENDATION:** Keep graph/timeline views derived from this model, in line
  with existing `ADR-006`/`REQ-UI-006`. No graph library object should become
  source of truth.

### Declared, observed, and accepted status

- **OBSERVED:** The current Analyzer displays CurrentIndex declarations and
  validates only a few explicit contradictions; it does not derive observed
  state from Ledger or GitHub.
- **RECOMMENDATION:** Model at least three separate status projections:
  `declared` from CurrentIndex/entity records, `observed` from Ledger/Git/GitHub,
  and `accepted` from qualifying verification plus independent review/Steering
  disposition. Never silently choose one when they disagree.
- **RECOMMENDATION:** Represent contradiction findings with both sources and a
  deterministic precedence policy for display only. Examples include an active
  assignment whose Issue is closed, a planned PR already merged, a completed
  Slice without qualifying evidence, a PR head changed after review, or a
  release record whose tag/Release does not exist.
- **RECOMMENDATION:** Distinguish `active implementation`, `blocked`, `on hold`,
  `last accepted boundary`, `awaiting Steering`, `merged`, `released`,
  `canceled`, `removed`, and `superseded`. The current completed `SLC-008` under
  `active.slice` is the concrete reason a single `active` field is insufficient.
- **RECOMMENDATION:** Treat GitHub state as observed external evidence, not as a
  replacement for repository product/architecture decisions. Cache/live reads
  require a snapshot timestamp, repository identity, and API limitations.

### Staleness and contradiction detection

- **OBSERVED:** Existing requirements correctly defer staleness until an
  explicit timestamp/status policy is supplied and give validation a fixed
  `analysisTime`. The current Ledger's date-only timestamps and missing
  post-`2026-07-14` Steering events demonstrate why ambient “older than N days”
  guesses would be unreliable.
- **RECOMMENDATION:** Define staleness per entity/work type and state, not one
  universal threshold. Examples: active Issue with no assignment/PR activity;
  assignment base behind default branch; draft PR with no update; completed
  Slice still declared active; review older than current head; verification
  older than changed inputs; handoff/CurrentIndex older than a Ledger/GitHub
  transition; release record disagreeing with tag/Release.
- **RECOMMENDATION:** When timestamps are missing, date-only, from different
  clocks, or outside a declared policy, emit `unknown`, not warning/error. Keep
  deterministic analysis by supplying both analysis time and policy explicitly.
- **RECOMMENDATION:** Allow a generated status command to refresh GitHub facts
  and derive staleness, while Analyzer itself stays read-only. Never rewrite
  project files automatically in the analysis path.

### Legacy and compatibility profiles

- **OBSERVED:** Current Analyzer intentionally supports only one
  structured-core profile and explicitly reports partial/unknown compatibility.
  The repository itself lacks a manifest declaring the installed SDP version.
- **RECOMMENDATION:** Introduce versioned capability profiles rather than one
  binary “old/new” schema: legacy phase/Sprint-first without Features; Sprint/
  Refactor structured core; project-local Feature additions; Feature/Refactor-
  first with Issue/CurrentAssignment bindings; and any safety-specific
  extensions. Detect capabilities from a manifest when available, and use
  conservative adapters/heuristics only when it is absent.
- **RECOMMENDATION:** Preserve raw unknown fields/sections and continue analysis
  of supported neighbors. A legacy repository without Features should not be an
  error merely for being legacy; show its profile and migration opportunities.
  Conversely, do not normalize an unknown `features:` mapping as canonical
  Feature semantics without a profile/schema.
- **RECOMMENDATION:** TIER-003 Markdown parsing must precede claims of complete
  legacy support, because recent repositories place important Feature,
  Steering, review, verification, and handoff facts outside the three current
  structured files.
- **RECOMMENDATION:** Build representative fixtures from multiple real SDP
  generations, including contradictions and partial migrations. Do not use the
  Analyzer repository's self-similar traceability alone as proof of ecosystem
  compatibility.

### Suggested implementation order after Issue #5 review

- **RECOMMENDATION:** First settle canonical Feature/Refactor/Fix, Issue binding,
  assignment, relation, status, timestamp, review, verification, and release
  schemas upstream in SDP. Analyzer must not invent them ahead of the method.
- **RECOMMENDATION:** Next add version/capability discovery plus Markdown/work
  document coverage and representative compatibility fixtures. This is the
  minimum evidence foundation for trustworthy cross-generation analysis.
- **RECOMMENDATION:** Then add lifecycle projection and declared-vs-observed
  contradiction rules, followed by Issue/PR/commit/check/release acquisition
  through a versioned snapshot/connector boundary.
- **RECOMMENDATION:** Add entity/relation/status listings and path queries before
  a graph or Gantt. Finally derive graph/timeline views and staleness UI from the
  tested projection. Keep repair/write-back separately mandated.

## What worked well

- **OBSERVED:** Repository-first authority, explicit read order, and compact
  handoffs make a fresh session able to reconstruct state without chat memory.
- **OBSERVED:** The initial horizontal architecture was detailed enough that
  later browser acquisition reused the same source/discovery/analysis pipeline
  without UI/core leakage.
- **OBSERVED:** Small Slice contracts, explicit non-goals, stop boundaries, and
  “record discoveries rather than expand scope” supported bounded autonomy.
- **OBSERVED:** Review materially improved code and evidence. Changes-required,
  reopening, correction, and re-review were preserved rather than rewritten as
  first-pass success.
- **OBSERVED:** CurrentIndex, Relations, and append-only Ledger separate present
  declaration, graph, and history. The blocked/unblocked SLC-007 and reopened
  SLC-006 sequences show the recovery value of immutable transitions.
- **OBSERVED:** Exact evidence standards are strong: deterministic parsers,
  provenance, explicit unknown/partial states, no invented entities, no opaque
  health score, and verification that requires a real structured relation and
  passed check description.
- **OBSERVED:** Read-only/security boundaries are disciplined. Source access is
  adapter-based; analyzed code and recorded commands are never executed;
  browser traversal is canonical/root-confined; partial acquisition is not
  misreported as missing data.
- **OBSERVED:** The committed Steering proposal already recognizes two crucial
  future principles: declared and Ledger-observed state should coexist, and
  optional Feature support must remain backward compatible with repositories
  that have no Features.
- **INFERENCE:** These qualities form a strong technical foundation for the
  future Issue-to-release graph, even though the current schema and operating
  workflow do not yet represent it.

## Pain points and accidental complexity

- **OBSERVED:** A six-day, 30-commit project accumulated 16 Study decisions, 59
  requirements, three Tiers, two Sprints, ten Slice IDs, nine review IDs, nine
  verification IDs, and 75 Ledger events. The evidence is rich, but the ratio
  of coordination records to product age is high.
- **OBSERVED:** Status and links are duplicated across CurrentIndex, Relations,
  Sprint contracts, implementation notes, handoff, README, verification,
  review, Ledger, and now Steering interactions. None is generated from GitHub
  or from one canonical assignment object.
- **OBSERVED:** `active.slice` can mean a completed retained boundary. This is
  documented but semantically surprising and makes staleness/current-assignment
  analysis ambiguous.
- **OBSERVED:** Ledger dates lack precise instants, and its history stops before
  later Steering commits. A future timeline or staleness rule cannot safely
  infer exact sequencing/duration from these records alone.
- **OBSERVED:** Direct-main development provides no early PR surface, immutable
  reviewed head, branch identity, CI enforcement, or GitHub acceptance trail.
  All formal role separation is internal documentation.
- **OBSERVED:** Relations repeats very large requirement/decision lists on
  several entities. Reverse edges and coverage summaries that could be derived
  are hand-maintained.
- **OBSERVED:** The lifecycle initially decomposed a complete planned product
  into Tiers and component-oriented Slices. Later Feature governance is a
  pending add-on rather than a current canonical owner.
- **OBSERVED:** The raw Steering interaction contract requires exact prompt and
  raw response copies. This preserves context but can duplicate lengthy chat
  transport into the repository while still lacking a small machine-readable
  Issue/branch/PR assignment.
- **OBSERVED:** The repository cannot declare which SDP generation it uses. The
  Analyzer profile is hard-coded in source/UI and does not yet mediate real
  installed manifests or profile adapters.
- **OBSERVED:** The current UI calls itself fixture mode and cannot analyze the
  repository selected by the user despite the accepted browser adapter. It
  cannot yet expose the graph/status/staleness conclusions central to its
  Mandate.
- **INFERENCE:** The principal accidental complexity is not stable IDs
  themselves; it is manual duplication of facts that Git, GitHub, schemas, or
  derived queries could own once.

## Carry forward

- **RECOMMENDATION:** Carry forward repository authority, re-anchoring, bounded
  Master/Worker/Reviewer/Architect roles, and the stop-at-Slice boundary.
- **RECOMMENDATION:** Carry forward the initial horizontal skeleton: source
  acquisition, compatibility/profile, parsing, normalized serializable domain,
  validation, application/query, and derived UI boundaries.
- **RECOMMENDATION:** Carry forward small contracts with goal, why now,
  authority, linked decisions/requirements/design, invariants, non-goals,
  verification, independent review, discoveries, completion, and stop
  conditions.
- **RECOMMENDATION:** Carry forward explicit partial/unknown compatibility,
  preservation of neighboring evidence, deterministic analysis time, stable
  rule IDs/fingerprints, and source provenance on every finding.
- **RECOMMENDATION:** Carry forward append-only correction history and the rule
  that verification/review are evidence objects, not status prose.
- **RECOMMENDATION:** Carry forward concise handoff for human context, while
  generating current coordinates and GitHub facts.
- **RECOMMENDATION:** Carry forward derived graph views rather than a graph
  library as canonical model, and keep Analyzer read-only even when a separate
  tool can propose or apply reviewed repairs.
- **RECOMMENDATION:** Carry forward backward-compatible optional Feature support
  and declared-vs-observed status separation from the pending governance
  proposal, but only after upstream SDP accepts the schema.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward a one-time full lifecycle/Tier plan as
  the permanent owner of all later capability work. Replace it with initial
  horizontal foundation plus Feature/Refactor/Fix-owned evolution.
- **RECOMMENDATION:** Do not carry forward Sprint or Tier as the mandatory
  capability owner. Retain them only as optional roadmap/execution grouping
  where they add value.
- **RECOMMENDATION:** Do not call a component-only increment a vertical Slice
  without qualifying it as an enabler. The canonical Slice should normally
  produce a bounded end-to-end capability through the necessary layers.
- **RECOMMENDATION:** Do not carry forward direct-main implementation, local-only
  uncommitted-tree review, or verification without an immutable branch/PR head
  when GitHub is available.
- **RECOMMENDATION:** Do not carry forward chat/Steering prompt text as the
  operational authority in place of a GitHub Issue. Keep it as supporting
  context or audit evidence.
- **RECOMMENDATION:** Do not carry forward ambiguous `active` fields that may
  point to completed work. Split active assignment, active execution, and last
  accepted boundary.
- **RECOMMENDATION:** Do not carry forward manually duplicated GitHub state,
  reverse relations, coverage summaries, branch/PR/commit identities, or
  release facts. Generate them and retain only genuinely authored decisions.
- **RECOMMENDATION:** Do not carry forward day-only timestamps for lifecycle
  projection where exact time matters, nor infer false precision from legacy
  dates.
- **RECOMMENDATION:** Do not carry forward an unversioned one-profile assumption
  or treat an absent Feature folder as an error. Use explicit capability
  profiles and conservative legacy adapters.
- **RECOMMENDATION:** Do not promote every project-local Steering transcript
  field or the proposed Feature lifecycle list into canonical SDP before Issue
  #5's synthesis and Steering review settle the method.
