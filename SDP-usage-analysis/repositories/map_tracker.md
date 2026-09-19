# map_tracker SDP usage study

## Evidence frame and repository facts

- **OBSERVED:** This report studies private, non-archived
  `Hans-Einar/map_tracker` at the inventory-recorded default branch and exact
  current identity
  [`master@22e03b49a27a7e888ee9dbdb97874e2038b43fb4`](https://github.com/Hans-Einar/map_tracker/tree/22e03b49a27a7e888ee9dbdb97874e2038b43fb4),
  committed `2025-12-31T22:14:55Z`. The live GitHub branch record still agreed
  on that identity when checked on 2026-08-21. The studied repository was not
  modified.
- **OBSERVED:** The default commit predates Issue #5's inclusive activity
  window, `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`. Four
  qualifying commits exist on non-default refs: first
  `3367e01aec45057815cfcba6707f3ff71ae06f91` at
  `2026-05-09T13:41:34Z` and last
  `19cb24d9ccf7d3d1d1ea8c5654de1e683e7574d5` at
  `2026-05-10T09:56:19Z`, both reachable from `V2`.
- **OBSERVED:** Current GitHub has three unprotected branches:
  `master@22e03b49...`,
  `V2@19cb24d9ccf7d3d1d1ea8c5654de1e683e7574d5`, and
  `V2_IOS@1883dc2c6f1d90babb3b21693ed79fdc18d66329`. It has zero
  Issues, zero pull requests, zero comments/native reviews, zero Actions
  workflows, zero check runs and zero commit-status contexts.
- **OBSERVED — default truth:** `master` contains no SDP installation. It has
  no `AGENTS.md`, `AGENTS-project.md`, `SDP/`, `Instructions/`, `Features/`,
  `Steering/`, `Sprints/`, `Refactors/`, `Fixes/`, `Traceability/`, standalone
  review/verification record, Handoff, `CurrentIndex`, `Relations`, Ledger or
  `.github/workflows/`. Its `codereview/00--Study.md`, first committed in
  2025-12, is an initial non-prescriptive product/code study, not proof of a
  review workflow or an installed SDP generation.
- **OBSERVED — non-default `V2` evidence:** Exact current `V2` contains a
  20-file, 2,553-line `SDP/` surface: `Instructions/README.md`, four Sprint
  folders with Sprint, Iteration/Slice, implementation-note and Handoff records,
  plus `Traceability/CurrentIndex.yaml`, `Relations.yaml` and `Ledger.ndjson`.
  It has no AGENTS files, manifest/version marker, Features, Steering,
  Refactors, Fixes, numbered initial lifecycle, separate CodeReview or separate
  Verification folder.
- **OBSERVED:** `SDP/Instructions/README.md` calls this a bootstrap and points
  agents to the external absolute path
  `C:\Users\hanse\GIT\farmStatistics\SDP\Instructions`; no installed Toolkit
  version or source commit is determinable. The repository-local operating
  records are therefore project-owned, pre-template/ad hoc evidence rather than
  a versioned canonical installation.
- **OBSERVED — limitation:** The substantial SDP practice is not merged to the
  advertised default branch and has no PR. It is current non-default branch
  evidence, not accepted default-tree truth. Live GitHub state supplies no
  issue, review, CI or merge decision that resolves that ambiguity.

## How SDP is used on `V2`

### Work start and Issue authority

- **OBSERVED — non-default evidence:** Work starts from user-reported/requested
  needs recorded by the Master in Slice prose and Ledger notes: runtime error
  visibility, Gardskart recovery, smoother GPS prediction, curved prediction,
  and Android/iOS build separation. There is no inspectable GitHub Issue or
  Steering decision carrying the original request.
- **OBSERVED:** Sprint 003 was first recorded as
  `planned_pending_branch`, then activated on `V2` after the user requested a
  shared location-provider boundary. The record names intended files,
  invariants, non-goals, verification and a macOS stop limitation, but does not
  bind an exact Issue/comment authority or baseline.
- **INFERENCE:** A fresh agent can reconstruct the Master's interpretation of
  each request but cannot compare it with the owner's exact wording or prove
  that later scope changes remained authorized.
- **OWNER DIRECTION:** Issue #5 makes a GitHub Issue the future operational
  assignment boundary. The request-derived practice here is useful contract
  evidence but should not remain the sole authority mechanism.

### Master, delegation and review roles

- **OBSERVED — non-default evidence:** The Ledger has 36 NDJSON events with
  explicit `master`, `worker` and `reviewer` roles. Slice contracts require code
  implementation by a Worker, separate independent Reviewer results, focused
  rework, Master verification and traceability/Handoff updates.
- **OBSERVED:** The loop was substantive. Sprint 001's third Slice required
  three review passes to resolve four findings, including a retry race. Sprint
  003 review found and re-reviewed an iOS status overwrite. Sprint 004 review
  triggered hardening for invalid fixes and per-segment speed. Reviewers did
  not merely restate green tests.
- **OBSERVED:** No Architect pass is recorded. Verification IDs exist, but the
  Ledger records verification under role `master`; no distinct Verifier actor or
  standalone verification artifact is present. There are no repository-local
  AGENTS role rules.
- **INFERENCE:** The documents strongly support practiced Worker/Reviewer
  separation and Master integration, but GitHub cannot prove fresh contexts,
  exact agent identity or the exact commit each review inspected.

### Branches and pull requests

- **OBSERVED:** Product and SDP work accumulated directly on long-lived
  non-default `V2`; `V2_IOS` was a separate source branch. Sprint 003 says the
  needed iOS scaffold was copied without merging `V2_IOS`.
- **OBSERVED:** Commit `19cb24d...` has subject `merged in IOS version`, but is
  a single-parent commit and `V2_IOS@1883dc2...` is not its ancestor. Git and
  Sprint prose therefore support copy/integration, not a Git merge. This is a
  concrete branch-language ambiguity.
- **OBSERVED:** No draft PR, native review, check or branch protection exists.
  No branch/PR identity is recorded in a Steering assignment or traceability
  relation, and `V2` remains unmerged more than three months after the studied
  activity.
- **RECOMMENDATION:** Future bounded work should open an early draft PR from an
  exact Issue/baseline and use exact candidate/review/verification heads. Branch
  names and commit subjects should describe Git operations truthfully.

### Feature, Refactor, Fix, Sprint, Iteration and Slice

- **OBSERVED — non-default evidence:** Sprint is the capability owner. Sprint
  001 contains three Iterations with one Slice each; Sprints 002, 003 and 004
  each contain one Iteration and one Slice. There is no Feature, Refactor or Fix
  identity even where the work is semantically a capability, architectural
  partition or correction.
- **OBSERVED:** Slice contracts consistently define goal, why now, expected
  files, invariants, non-goals, traceability IDs, verification and completion.
  Several are genuinely vertical: the Gardskart recovery Slice crosses typed
  domain failure, app coordination, settings UI and tests; the iOS partition
  crosses build entrypoints, provider boundaries, UI exclusion, platform
  scaffold and verification.
- **OBSERVED:** The 1:1:1 pattern in Sprints 002-004 provides little evidence
  that both Sprint and Iteration are necessary. Review rework occurs inside the
  same Iteration rather than being represented as a new feedback cycle.
- **OWNER DIRECTION:** Issue #5 requires future implementation ownership by
  Feature or Refactor, permits a genuine small Fix, and retains Slice as a small
  vertical end-to-end capability. Sprint must not remain capability owner only
  because this older branch used it.
- **RECOMMENDATION:** Preserve these IDs as branch history. Prospectively place
  the demonstrated vertical contracts under Feature/Refactor/Fix, make Sprint
  an optional execution grouping, and use Iteration only when a real replanning
  or feedback boundary exists.

### Studies and later requirement/design evolution

- **OBSERVED — default truth:** `codereview/00--Study.md` documents current
  structure, data flow, dependencies, evolution constraints and open domain
  questions. It is one broad, non-prescriptive study with no claim/evidence
  levels, independent sub-studies or convergence gate.
- **OBSERVED — non-default evidence:** `V2` has no canonical Requirements,
  Architecture or Design folders. Instead, Sprint contracts and CurrentIndex
  introduce use-case (`UC-*`), requirement (`R-*`) and design (`D-*`) records as
  later work needs them. Sprint 003 adds a platform service boundary and Sprint
  004 adds a rolling turn-rate design long after the default-tree initial study.
- **INFERENCE:** The branch demonstrates useful feature-local requirement and
  design evolution rather than a frozen one-time plan. It does not demonstrate
  deliberate reconciliation with an initial horizontal architecture because no
  authoritative initial SDP architecture/design exists.
- **RECOMMENDATION:** Carry forward later local refinement, but anchor it under
  Feature/Refactor authority and require explicit architecture revision when a
  stable system boundary changes. Use multiple Studies/convergence only when
  uncertainty warrants them, not for every small change.

### Review, verification and Handoff

- **OBSERVED — non-default evidence:** Review notes are embedded in
  `ScrumIterations.md`; implementation and verification evidence is repeated in
  `implementationNotes.md` and Handoffs. No standalone review/verification
  files or GitHub reviews exist.
- **OBSERVED:** Verification is candid and risk-aware. Focused Flutter tests,
  format, analyze, diff checks, an Android APK build and a project-specific iOS
  import-graph verifier are recorded. Full-suite `SerialGPIO` teardown failure
  is repeatedly retained as a known warning, and iOS no-codesign verification
  remains pending macOS/Xcode rather than being claimed complete.
- **OBSERVED:** Each Sprint has a concise Handoff naming authoritative records,
  completed behavior, evidence, residual warnings and likely next work. Sprint
  003 records exactly what remains for macOS and that `V2_IOS` can be removed
  only after verification.
- **INFERENCE:** The Handoff/evidence practice materially improves recovery
  after session loss. Its main weakness is identity: commands and findings are
  not pinned to exact candidate/review/verification commits and are not backed
  by GitHub checks.

### CurrentIndex, Relations, Ledger, Steering and automation

- **OBSERVED — non-default evidence:** `CurrentIndex.yaml` is a 538-line,
  72-item registry rather than a compact live pointer. It contains use cases,
  requirements, designs, four Sprints, six Iterations/Slices, review IDs,
  verification IDs and five resolved findings. `Relations.yaml` has 280 lines
  connecting use case to requirement, requirement to design, design to Slice,
  and Slice to review/verification/finding.
- **OBSERVED:** The 36-event Ledger records Slice planning/start, Worker
  completion, Reviewer result, corrective rework, re-review and Master
  verification. It provides a readable chronology but repeats details already
  in Iteration notes, implementation notes, Handoff and the registry.
- **OBSERVED:** There is no `Steering/`, `CurrentAssignment`, Steering actor,
  GitHub Issue/PR binding, release binding or compact current assignment. The
  final Handoffs say work is complete and the next priority must be chosen, but
  no machine-readable current pointer represents that stop state.
- **OBSERVED:** Automation is project-local and manually invoked: Flutter
  tests/analyze/format, build, diff checks and the iOS import-graph verifier.
  GitHub supplies no CI checks, required review or stale-state detection.
  CurrentIndex, Relations, Ledger and Handoff were hand-maintained.
- **INFERENCE:** Stable requirement/design-to-Slice and review/verification
  relations are useful; maintaining 72 registry objects plus 36 prose-rich
  events for six Slices is disproportionately expensive and still fails to
  answer the live Issue/branch/PR/acceptance question.
- **RECOMMENDATION:** Keep a compact assignment/index and stable graph for
  `Issue <-> Feature/Refactor/Fix <-> important requirement/design decision <->
  Slice <-> review/verification <-> PR/commit`. Generate paths and live GitHub
  state; reserve append-only events for authorization, decision, rework,
  acceptance, blocked evidence and merge/release transitions.

## What worked well

- **OBSERVED:** Slice contracts are concrete enough for autonomous Workers:
  files, invariants, non-goals, tests and completion signals are consistently
  explicit.
- **OBSERVED:** Independent review materially improves quality and can require
  repeated bounded rework even after focused tests pass.
- **OBSERVED:** The branch allows new requirements and designs to emerge with
  later capabilities, including a major platform boundary, without pretending
  they were known in the initial study.
- **OBSERVED:** Verification records distinguish focused success, known
  unrelated suite failure, platform limitations and unperformed physical/macOS
  evidence. This is stronger than a blanket green/complete claim.
- **OBSERVED:** Handoffs, relations and the event chronology make the complex
  `V2` work understandable without chat memory.
- **OBSERVED:** Sprint 003 copied selected iOS work rather than silently merging
  an entire divergent branch, and its import-graph verifier converted a key
  architecture invariant into executable evidence.

## Pain points and accidental complexity

- **OBSERVED:** The strongest SDP evidence exists only on `V2`, not default,
  and no Issue/PR/Steering decision explains its acceptance or next disposition.
- **OBSERVED:** The bootstrap depends on an absolute path in another local
  repository and has no version/source manifest or local AGENTS role contract.
  Another machine cannot reproduce the operating instructions from this
  repository alone.
- **OBSERVED:** Sprint remains the semantic owner and three of four Sprints add
  an Iteration layer around a single Slice. Feature/Refactor/Fix distinctions
  are absent.
- **OBSERVED:** Review and verification are detailed but not exact-head-pinned,
  independently attributable on GitHub, or enforced by CI/branch protection.
- **OBSERVED:** `CurrentIndex` is a large historical registry, while Handoffs
  carry the practical next-state narrative. The repository has much machine-
  readable state but no compact current assignment.
- **OBSERVED:** Relations and Ledger improve recovery but duplicate Sprint
  prose and still omit GitHub identities. Thirty-six events for six Slices are
  manually maintained.
- **OBSERVED:** Commit `3367e01...` first adds Sprints 001-003 process records
  together with product changes after some recorded event dates. Git cannot
  independently correlate each documented Worker/review pass with a distinct
  candidate commit.
- **OBSERVED:** The subject `merged in IOS version` contradicts the actual
  single-parent/copy integration shape, illustrating branch/commit ambiguity.

## Carry forward

- **RECOMMENDATION:** Preserve the strong Slice contract shape, bounded Worker
  implementation, independent Reviewer/rework loop, Master evidence inspection
  and candid verification limitations.
- **RECOMMENDATION:** Preserve later Feature/Refactor-local requirements and
  design evolution, with explicit escalation when a stable horizontal
  architecture boundary changes.
- **RECOMMENDATION:** Preserve concise Handoffs with decisions, residual risk,
  unavailable evidence and the next authorized action.
- **RECOMMENDATION:** Preserve stable relations that answer why a requirement or
  design exists and which Slice/review/verification proves it, while reducing
  registry and event duplication.
- **RECOMMENDATION:** Promote generic executable invariants such as import-graph
  validation as a pattern; keep Flutter/iOS/Android commands and thresholds
  project-local.
- **RECOMMENDATION:** Prospectively bind the work to a GitHub Issue, compact
  assignment record, exact baseline, branch, early draft PR and exact reviewed/
  verified head so Steering can accept or reject it from repository evidence.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not treat a long-lived non-default branch plus internal
  prose as accepted project state. Preserve `V2` as evidence until the owner
  decides migration/merge, but require explicit Issue/PR/Steering disposition
  for future work.
- **RECOMMENDATION:** Do not carry forward absolute external instruction paths,
  unversioned installs or missing local AGENTS authority as a reproducible SDP
  bootstrap.
- **RECOMMENDATION:** Do not preserve mandatory
  `Sprint -> Iteration -> Slice` ownership when Feature/Refactor/Fix and one
  vertical Slice describe the work more truthfully.
- **RECOMMENDATION:** Do not treat actor labels in a manually written Ledger as
  proof of fresh independent review or exact-head verification.
- **RECOMMENDATION:** Do not make CurrentIndex both a complete historical
  registry and the current assignment surface, and do not repeat live branch/
  acceptance state across Handoff, notes and events when it can be generated.
- **RECOMMENDATION:** Do not infer a canonical multi-study/convergence model
  from `codereview/00--Study.md`; it is one broad project study without those
  controls.
