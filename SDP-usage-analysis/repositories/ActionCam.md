# ActionCam

Study performed on 2026-08-21. Evidence labels in this report have the meanings
required by `Hans-Einar/SDP#5`: **OBSERVED** is direct repository/GitHub
evidence, **OWNER DIRECTION** is the future constraint in Issue #5,
**INFERENCE** is an interpretation, and **RECOMMENDATION** is a proposed SDP
practice.

## Repository facts

- **OBSERVED — Repository and baseline.** `Hans-Einar/ActionCam` is a private,
  non-archived repository. Its default branch is `main`. The exact default-tree
  commit studied is `861a22d4b2336e1eabd24051b914c8d24219ee01`, committed
  `2026-08-20T18:32:55Z` by `Hans-Einar` and created by GitHub's `web-flow` as
  the merge commit for PR #60. This report rechecked the remote default head,
  rather than treating the PR branch as current repository state.
- **OBSERVED — Activity.** The first commit reachable from the studied default
  head is `0d483beb29d5d6df201adc651a8c7ba83dc83d47`, committed
  `2026-08-12T22:37:08Z` (`Initialize ActionCam reverse-engineering repository`).
  The studied history contains 180 reachable commits through
  `2026-08-20T18:32:55Z`; therefore ActionCam has actual commits inside Issue
  #5's `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z` window.
- **OBSERVED — Issue authority.** Issue
  [#59](https://github.com/Hans-Einar/ActionCam/issues/59), authored by GitHub
  user `Hans-Einar`, was open from `2026-08-18T18:17:00Z` to
  `2026-08-20T18:33:44Z`. It fixed the initial baseline at
  `main@36b2a8aa17786d06df98de84dafdacfd05ef86d8`, stated the research goal,
  paths, phases, physical-safety boundaries, acceptance milestones and Codex
  branch/PR workflow. Its three comments, also by `Hans-Einar`, created the
  branch/PR boundary, narrowed the merge criterion, split follow-up work, and
  recorded the landed result.
- **OBSERVED — Branch and PR.** Historical execution occurred on
  `work/59-standalone-media-stack` in draft PR
  [#60](https://github.com/Hans-Einar/ActionCam/pull/60), opened
  `2026-08-18T18:28:03Z` directly from the stated base. The PR contains 35
  commits from `75e79f21f5eb47bdaaf8aff0ab9371d228011cc1` through exact reviewed
  head `ec7a516681e7f620a7906f7ca6927ab5ff0d0478`, 357 changed files,
  31,734 additions and 100 deletions. `Hans-Einar` marked it ready and merged
  it at `2026-08-20T18:32:56Z`. The PR branch is historical evidence; its tree
  is now contained in the studied default branch through merge commit
  `861a22d4b2336e1eabd24051b914c8d24219ee01`.
- **OBSERVED — PR discussion, review and checks.** PR #60 has 13 checkpoint/
  handoff comments, all under GitHub identity `Hans-Einar`, dated
  `2026-08-18T19:46:52Z` through `2026-08-20T17:44:34Z`. It has one submitted
  GitHub review, also by `Hans-Einar`, state `COMMENTED`, at commit
  `84c7a70dfe726d7145c862c97f7b80ee17e6d86d`; that review directs a final
  merge-closure audit. There are no inline review comments, no GitHub check
  runs, no commit status contexts, and an empty PR check rollup at exact head
  `ec7a516...`. Verification was preserved in repository evidence and the
  final PR comment, not GitHub CI.
- **OBSERVED — Related GitHub work.** Issue #29 and unmerged draft PR #58 were
  the separate `sdvcam` production track; PR #58 was closed unmerged on
  `2026-08-20T12:59:15Z`. Owner direction split unfinished product work into
  open Issues [#61](https://github.com/Hans-Einar/ActionCam/issues/61)
  (image quality/calibration), [#62](https://github.com/Hans-Einar/ActionCam/issues/62)
  (recording/stills), [#63](https://github.com/Hans-Einar/ActionCam/issues/63)
  (adaptive streaming/network/reliability), and
  [#64](https://github.com/Hans-Einar/ActionCam/issues/64) (on-camera
  calibration/UI). PR #60 explicitly excludes those scopes.
- **OBSERVED — SDP installation.** There is no repository-root `SDP/`,
  `SDP.manifest.yaml`, `AGENTS.md`, or `AGENTS-project.md` at the studied tree,
  so a repository-wide installed SDP generation/version is not determinable.
  Instead, the repository has task-local SDP islands at `camera/mr500/SDP/`,
  `reverse/front-display/SDP/`, and `reverse/standalone-media/SDP/`. The last
  one says explicitly in `reverse/standalone-media/SDP/README.md` that it is a
  research-first extension of the recent TerrainAnalyzer process. It is a
  project-local evolution, not proof of a released Toolkit version.
- **OBSERVED — Standalone-media locations.** The durable local records are:
  `reverse/standalone-media/SDP/01--Mandate/01--Mandate.md`,
  `reverse/standalone-media/SDP/02--Studies/`,
  `reverse/standalone-media/SDP/03--Requirements/`,
  `reverse/standalone-media/SDP/Instructions/ResearchProcess.md`,
  `reverse/standalone-media/SDP/Sprints/Sprint-001-*` through
  `Sprint-006-*`, `reverse/standalone-media/SDP/Steering/CurrentAssignment.yaml`,
  and `reverse/standalone-media/SDP/Traceability/{CurrentIndex.yaml,Relations.yaml,Ledger.ndjson}`.
  Handoffs and verification results live in each Sprint folder. Component
  machine evidence includes
  `reverse/standalone-media/evidence/encoded-stream-handoff/verification.json`
  and `reverse/standalone-media/evidence/live-rate-control/verification.json`.
  The wider repository claim policy is `reverse/VERIFICATION.md`. There are no
  standalone-media `Features/`, `Refactors/`, `CodeReview/`, or
  `Verification/` directories and no formal Architecture/Design phase folders.

## How SDP is used in this repository

### Work start and Issue authority

- **OBSERVED —** Work started from Issue #59, not from an SDP Sprint. The seed
  workspace commit and draft PR appeared at `18:27–18:28Z`; the local research
  SDP was bootstrapped in commit `fe8c460aa393144978611dc22400fd0881193389`
  at `2026-08-18T18:55:41Z`. `ResearchProcess.md` makes the authority order
  explicit: Issue #59 body/comments, later owner direction on the Issue/PR,
  local SDP, repository evidence/rules, then pinned references.
- **INFERENCE —** This is strong evidence for Issue-as-assignment and SDP as
  the repository-resident execution/refinement surface, rather than requiring
  the Sprint record to originate the capability.
- **OWNER DIRECTION —** Issue #5 requires GitHub Issues to become the
  operational assignment boundary and a bounded Codex session to act as
  Master. ActionCam supports that direction; it does not justify retaining a
  Sprint-first owner.

### Codex Master, delegation, Architect, Verifier and Reviewer

- **OBSERVED —** Issue #59 says it is intended for the OpticStudio Codex agent.
  The local ledger contains 109 events with role `master` and five with role
  `owner`; it contains no `worker`, `architect`, `verifier`, or `reviewer`
  actor. Every `SM-VERIFY-*` completion is recorded by `master`. Process
  architecture was handled as `SM-STUDY-012`, not as a separately attributable
  Architect pass.
- **OBSERVED —** `ResearchProcess.md` correctly says substantial target-running
  code should receive separate review, unavailable review must remain pending,
  and self-check must not be called independent review. Sprint 1, 4 and 5
  handoffs say independent review remains pending. The ledger's physical
  authorization event also says it remains pending.
- **OBSERVED —** Physical execution nevertheless proceeded. The final ledger
  event `merge_closure_review_completed` is role `master`; final verification
  is also role `master`. On GitHub, the closure direction, final handoff, merge
  and every checkpoint are attributable only to `Hans-Einar`.
- **INFERENCE —** A focused review did useful work—it found and closed a real
  short-datagram panic blocker—but the durable evidence cannot establish a
  fresh, independent Reviewer or Verifier identity. Nor does it show Worker
  delegation. This repository demonstrates the Master pattern, but not the
  complete role-separation pattern requested for future canonical SDP.
- **RECOMMENDATION —** Canonical issue execution should record delegated role,
  agent/run identity, exact reviewed head, finding IDs and disposition. A
  review requested by the owner and a master-authored self-review must remain
  distinguishable from independent review.

### Branch and draft PR

- **OBSERVED —** Issue #59 and its first comment required a new branch and
  draft PR from an exact `main` commit, isolated from PR #58. The branch/PR were
  created within eleven minutes of the Issue and received progressive commits
  plus checkpoint comments. Final scope, test evidence and exact head were
  placed in the PR before merge; no force-push was reported.
- **RECOMMENDATION —** Carry forward exact base, one issue-bound branch, early
  draft PR, progressive evidence, exact-head review, and explicit parallel-PR
  isolation. Generate current PR state from GitHub instead of copying it by
  hand.

### Feature, Refactor, Sprint, Iteration and vertical Slice

- **OBSERVED —** There is no Feature or Refactor record. The hierarchy is
  `SM-MANDATE-059 -> SM-REQ-* / SM-STUDY-* -> six SM-SPRINTs -> one
  SM-ITERATION per Sprint -> 27 SM-SLICEs -> six SM-VERIFY records`.
  All six Sprints and Iterations were created and completed within roughly two
  days, so they are execution checkpoints, not calendar timeboxes.
- **OBSERVED —** Sprint 4 calls itself the “first post-VENC vertical slice,”
  while its Iteration contains four objects also named Slices. Many Slices are
  activation or verification/publication steps, not vertical end-to-end
  capabilities. Conversely, `SM-SLICE-026` grew to include H.264 VBR, H.265,
  ACM2, MPEG-TS/SRT, ABR control and a 12-hour run after later owner direction.
- **INFERENCE —** ActionCam preserves useful bounded contracts, but it uses
  “Slice” at two levels and sometimes lets one Slice become a substantial
  workstream. The one-Iteration-per-Sprint structure adds identity without an
  observable grouping benefit.
- **OWNER DIRECTION —** Future work must belong to a Feature or Refactor, with
  Fix retained only for genuine corrections. Do not preserve ActionCam's
  Sprint ownership merely because this repository used it.
- **RECOMMENDATION —** Treat the Issue as assignment to one Feature/Refactor
  (or a bounded research item), make a Slice one small vertical deliverable
  with explicit entry/completion evidence, keep Sprint optional as a temporary
  grouping, and omit Iteration where it is only 1:1 with Sprint. The final
  short-datagram correction is evidence that a small Fix can stay inside the
  active assignment rather than requiring a new hierarchy.

### Later requirements and design evolution

- **OBSERVED —** Requirements were added as evidence changed:
  `03--Requirements/02--HostSimulationPreparation.md`,
  `03--EncodedStreamHandoff.md`, `04--SeamlessLiveRateControl.md`, and
  `05--HardwareHevcLiveContribution.md` postdate the initial checkpoint.
  Process architecture was refined in `SM-STUDY-012` and
  `reverse/standalone-media/process-architecture.md`. Later owner direction
  introduced H.265, production uplink integration and continuous operation.
- **INFERENCE —** This is direct evidence that requirements and design can
  evolve after an initial horizontal study, without pretending they were all
  known up front. The absence of a formal feature-local design record also
  shows why the future method needs an explicit home for such refinement.

### Independent/multiple Studies and convergence

- **OBSERVED —** `ResearchProcess.md` gives every `SM-STUDY-*` bounded
  questions, inputs, dependencies, allowed actions, deliverables, status and
  promotion gates. Initial Studies 001–006 progressed independently; later
  Studies 012–017 were introduced as requirements evolved. Status is scoped,
  and target/reference/simulation claims are deliberately separated. The key
  rule is: design horizontally across independent studies; implement
  vertically only after a convergence gate supports a bounded end-to-end
  slice.
- **OBSERVED —** The method corrected important hypotheses rather than hiding
  them: runtime changed the initial IMX286 assumption to `imx386_mipi`; host
  simulation was never promoted to hardware proof; an inventory-process
  kernel oops invalidated a session and forced cold recovery; incomplete image
  quality remained open after encoder success.
- **RECOMMENDATION —** Carry forward multi-study contracts, scoped dependency
  blocking, claim levels and named convergence gates. Canonical statuses
  should be a controlled vocabulary, and work moved to a follow-up Issue
  should be marked transferred/superseded with the target Issue identity.
- **RECOMMENDATION —** Keep exact camera temperatures, device selectors,
  storage paths, MTD prohibitions and physical experiment commands local to
  ActionCam. Canonical SDP should carry only the generic authorization,
  recovery, evidence, bounded-action and stop-condition pattern.

### Review, verification and handoff

- **OBSERVED —** Each Sprint ends with a verification/publication Slice and a
  Handoff containing current objective, authority, exact next step, claim
  boundary, recovery and stops. PR checkpoint comments name commits, tests and
  remaining physical gates. The final audit found one Blocking and two Low
  findings, fixed all three, then repeated host/unit/race/vet/static/shell/
  reproducible-build/QEMU/evidence checks at exact head `ec7a516...`.
- **OBSERVED —** GitHub itself contains no automated checks, and the final
  review/verification cannot be independently attributed beyond `Hans-Einar`
  and ledger role `master`.
- **INFERENCE —** Handoff quality strongly improved recoverability and claim
  discipline. Acceptance confidence would be stronger if the same evidence
  were produced by CI and a separately identified Reviewer/Verifier.

### CurrentIndex, Relations and Ledger

- **OBSERVED —** At the studied tree, `CurrentIndex.yaml` contains 77 stable
  IDs. `Relations.yaml` has 152 relations: 65 `addresses`, 26 `advances`, 33
  `contains`, two `depends_on`, six `prepares`, 14 `requires`, and six
  `verifies`. `Ledger.ndjson` has 114 append-only events, including start/
  completion, owner direction, recovery, physical observation, evidence
  correction, review and verification.
- **OBSERVED —** The graph connects Mandate, Requirements, Studies, Sprints,
  Iterations, Slices and verification, but GitHub Issue, branch, PR, commit and
  follow-up Issue identities are not first-class relation nodes. They are
  repeated as fields in CurrentAssignment, Sprint prose and ledger payloads.
- **INFERENCE —** The ledger was valuable where it withdrew a false residual-
  process claim and recorded later correction. The volume of manually
  synchronized IDs/relations for a two-day assignment is also material
  ceremony.
- **RECOMMENDATION —** Retain stable IDs for Feature/Refactor, material Study,
  Requirement/decision, Slice, verification/review and release. Generate
  issue/branch/PR/commit state from GitHub and derive routine containment/
  publication events. Do not require IDs for every activation and publication
  action unless they carry an independently useful contract.

### Steering/CurrentAssignment

- **OBSERVED —** `Steering/CurrentAssignment.yaml` is the strongest generic
  experiment in this repository. It binds assignment `ISSUE-059-STANDALONE-MEDIA`
  to repository, Issue, base/branch/PR, Mandate/Sprint/Iteration/Slice, Studies,
  verification, permissions, stop boundary, recovery, architecture direction
  and physical gates. Its schema name is project-local
  `actioncam-standalone-media-assignment-v1`; no schema/validator is checked in.
- **INFERENCE —** The generic core is valuable for session recovery and
  machine-readable authority. The embedded camera-run telemetry and copied
  GitHub state made the file expensive to keep current.
- **RECOMMENDATION —** Canonicalize a much smaller CurrentAssignment:
  Issue/work-item identity, exact baseline, branch/PR, active Slice, delegated
  roles, allowed actions, stop boundary, expected verification/review and next
  handoff. Link operational evidence; do not embed live process/run state.

### GitHub binding and automated versus conventional behavior

- **OBSERVED —** Binding is explicit but distributed: Issue/PR prose,
  `SM-MANDATE-059`, every Sprint authority block, CurrentAssignment and the
  first ledger event all repeat Issue #59, PR #60 and branch/base identities.
  Progressive PR comments were manual milestones under `Hans-Einar`.
- **OBSERVED —** Component verification is substantially automated through
  unit/race/vet/static-analysis, reproducible build, QEMU, evidence JSON and
  parse/link/relation checks. Synchronizing GitHub state, CurrentAssignment,
  CurrentIndex, Relations, Ledger, Sprint status and README handoff pointers is
  a conversational/manual convention. GitHub Actions did not enforce it.
- **RECOMMENDATION —** Let `gh-sdp`/SDP tooling create the binding once and
  refresh GitHub-derived fields. CI should publish exact-head validation.
  Humans/agents should author only decisions, scope, safety, evidence meaning
  and exceptions.

## What worked well

- **OBSERVED — Bounded evidence language.** Verified/Observed/Inference/
  Unknown/Reference separation prevented public V316 behavior and host mocks
  from becoming false target claims.
- **OBSERVED — Safe convergence.** Independent studies permitted useful host
  architecture and simulation while hardware facts were blocked, then required
  named target gates before promotion.
- **OBSERVED — Recoverability.** Exact base/head identities, progressive
  commits, early draft PR, checkpoint comments, per-Sprint handoffs and the
  append-only ledger made a dense two-day effort reconstructable.
- **OBSERVED — Truthful correction.** The process retained failures, withdrew
  an incorrect “no residual process” claim, corrected the sensor identity and
  recorded the real pre-merge blocker rather than smoothing contradictions.
- **OBSERVED — Scope steering.** Owner direction stopped PR #60 from absorbing
  all remaining product work and created Issues #61–#64 with explicit
  boundaries.
- **INFERENCE —** The strongest reusable pattern is not the camera-specific
  safety detail; it is the combination of issue authority, machine-readable
  action limits, independent studies, convergence, exact evidence levels and
  durable handoff.

## Pain points / accidental complexity

- **OBSERVED — Oversized assignment.** One Issue/PR accumulated six Sprints,
  six one-to-one Iterations, 27 Slices, 35 commits, 357 files and more than
  31,000 added lines in about two days. Master-per-Issue alone does not make an
  Issue bounded; late follow-up splitting was necessary.
- **OBSERVED — Ambiguous hierarchy.** Sprint remained capability owner, every
  Sprint had exactly one Iteration, and “vertical slice” described both Sprint
  4 as a whole and its subordinate Slice records. Activation/publication Slices
  inflated the hierarchy.
- **OBSERVED — Duplicated mutable state.** GitHub, CurrentAssignment,
  CurrentIndex, Studies, Sprint files, README, Handoff and Ledger repeat status
  and coordinates. Several copies disagree at the current default head.
- **OBSERVED — Review provenance gap.** The documented rule requires separate
  review, but durable identities show master verification and owner-direction
  review only. No CI or independent GitHub review supplies a second provenance
  path.
- **OBSERVED — Project-local schema/status drift.** CurrentAssignment has no
  checked schema. Study statuses such as `exact_target_hevc_ingest_verified`
  and `exact_target_live_mutation_verified` are used even though
  `ResearchProcess.md` defines a different controlled status list.
- **OBSERVED — Discoverability gap.** There is no root AGENTS contract or
  repository-wide SDP manifest; a new agent must discover multiple nested SDP
  islands and determine which one applies.
- **INFERENCE —** The large amount of precise physical evidence is justified
  by ActionCam risk. The same density in canonical workflow state would be
  bureaucracy; evidence should be linked, not duplicated into assignment and
  traceability records.

## Direct contradictions and limitations

- **OBSERVED — Stale assignment.** At current `main@861a22d4...`,
  `Steering/CurrentAssignment.yaml` says `status: active`, phase
  `exact_camera_hevc_continuous_observation`, PR state `draft_open`, active
  `SM-SLICE-026`, and verification pending. GitHub says Issue #59 is closed and
  PR #60 is merged; `CurrentIndex.yaml`, the Ledger and final Handoff say
  Slices 026/027, `SM-VERIFY-006`, Iteration 006 and Sprint 006 are complete.
- **OBSERVED — Stale navigation.** `reverse/standalone-media/SDP/README.md`
  calls Sprint 3's Handoff the “current handoff” and later calls Sprint 6 the
  current physical continuation.
- **OBSERVED — Study-state contradictions.** `SM-STUDY-015` says
  `physical_validation_pending`, while CurrentIndex says
  `exact_target_hevc_ingest_verified`. Planned `SM-STUDY-010` still describes
  MPEG-TS/SRT integration although later Study 017 and Sprint 6 delivered it;
  planned `SM-STUDY-011` still asks coexistence-versus-replacement although
  owner direction and Sprint 6 selected replacement. Mandate and all
  Requirements remain `active` after Issue #59 was declared complete, without
  explicit transfer relations to Issues #61–#64.
- **OBSERVED — Review contradiction.** Independent review was a stated
  prerequisite and recorded pending at physical authorization, but target work
  proceeded. Final review evidence is useful but attributed to `master` and
  `Hans-Einar`, not an independently identifiable reviewer.
- **LIMITATION —** All commits, comments, review and merge actions are exposed
  through the single GitHub identity `Hans-Einar`. Repository evidence cannot
  identify which Codex session or human performed a pass behind that account.
- **LIMITATION —** This study inspected the complete relevant repository tree,
  Issue #59 timeline/comments, PR #60 comments/review/commits/check/check-status
  surfaces and committed evidence at the exact default head. It did not rerun
  physical camera experiments, and raw/private camera evidence intentionally
  retained only on microSD is not independently inspectable here.
- **LIMITATION —** The repository is private; links require appropriate access.
  No absent GitHub check can be inferred from the rich self-reported test list.

## Carry forward

- **RECOMMENDATION —** GitHub Issue as the durable assignment boundary, with
  exact baseline, bounded goal/non-goals, authorization, stop condition and
  reporting contract.
- **RECOMMENDATION —** One issue-bound Codex Master, early branch/draft PR,
  progressive evidence comments and an exact-head final handoff.
- **RECOMMENDATION —** Independent bounded Studies with scoped dependencies,
  controlled claim levels and explicit convergence before vertical
  implementation.
- **RECOMMENDATION —** A small machine-readable CurrentAssignment that binds
  Issue, Feature/Refactor/research item, branch/PR, active Slice, role
  delegation, allowed actions, verification/review and stop boundary.
- **RECOMMENDATION —** Repository-resident handoff and append-only correction
  history, especially for safety, evidence promotion and withdrawn claims.
- **RECOMMENDATION —** Fresh independent Reviewer/Verifier passes with exact
  head and durable actor identity; automate deterministic checks as CI.
- **RECOMMENDATION —** Feature-local requirements/studies/design refinement
  after the initial horizontal skeleton. Follow-up Issues #61–#64 demonstrate
  why later capability work must not be forced into the original design.
- **RECOMMENDATION —** Preserve ActionCam's generic safety pattern—explicit
  authority, bounded experiment, recovery and stop gates—while keeping its
  camera, thermal, storage, ADB and firmware specifics project-local.

## Legacy / do not carry forward

- **RECOMMENDATION —** Do not carry forward Sprint as the conceptual owner of
  a capability, or a mandatory Iteration that is always 1:1 with its Sprint.
- **RECOMMENDATION —** Do not call activation and publication bookkeeping
  “vertical Slices”; reserve Slice for a small end-to-end outcome.
- **RECOMMENDATION —** Do not carry forward one enormous Issue/PR merely
  because a Master can maintain many internal IDs. Split bounded Features or
  research decisions earlier.
- **RECOMMENDATION —** Do not manually duplicate GitHub state, active
  coordinates and operational telemetry across CurrentAssignment, CurrentIndex,
  README, Sprint, Handoff and Ledger.
- **RECOMMENDATION —** Do not retain project-invented status vocabularies or a
  schema-less ActionCam assignment document as canonical format; extract and
  validate the generic subset.
- **RECOMMENDATION —** Do not treat master self-verification or an owner
  steering review as independent review.
- **RECOMMENDATION —** Do not canonicalize ActionCam's device-specific safety
  mechanisms, evidence volume, physical temperature thresholds, removable-SD
  layout, ADB selectors, MTD rules or media-stack claim taxonomy. They are
  justified local controls, not general software-process requirements.
