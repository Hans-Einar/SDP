# weight_app_flutter SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/weight_app_flutter` at the exact
  current default-branch identity
  [`master@538577830a2dbc13f603a2986f4a00f5afcc66a1`](https://github.com/Hans-Einar/weight_app_flutter/tree/538577830a2dbc13f603a2986f4a00f5afcc66a1)
  (tree `8e9f940cc8cc0e98a9f277b185a35e999b489fc3`), committed
  `2026-07-14T07:30:41Z` as the merge commit of PR #1. GitHub's advertised
  default branch, remote `HEAD`, the local read-only commit object, and the API
  branch record all agreed on that identity when checked on 2026-08-21. The
  studied repository was not modified.
- **OBSERVED:** GitHub Issues, PRs, issue/PR comments, review comments, native
  review submissions, timelines, check runs, commit statuses, branches, tags,
  and releases were inspected through authenticated GitHub CLI/API calls on
  2026-08-21. GitHub Issues are enabled, but the repository has zero Issues.
  Its only issue-numbered objects are merged PR
  [#1](https://github.com/Hans-Einar/weight_app_flutter/pull/1) and open draft PR
  [#2](https://github.com/Hans-Einar/weight_app_flutter/pull/2).
- **OBSERVED:** Default-tree truth is kept separate from open/unmerged evidence.
  PR #2 has exact head
  [`c4c31cbae4aa1529da9cf5a4255f30ee83429066`](https://github.com/Hans-Einar/weight_app_flutter/tree/c4c31cbae4aa1529da9cf5a4255f30ee83429066)
  (tree `24d7df92639a42d324ea0ddf410c9dbaf1fce8e7`). It is an extensive, internally
  reviewed REF-001 execution record, but it is not on `master` and remains open,
  draft, and unmerged. Statements about its implementation, active coordinates,
  reviews, verification, and traceability are explicitly **open/unmerged
  evidence**, not current default-tree state.
- **OBSERVED:** Issue #5's fixed activity window is
  `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`, inclusive. The
  authoritative inventory found six advertised refs, 36 unique reachable
  commits in total, and 35 commits inside the window. The earliest qualifying
  commit was
  `4e687ff4b54b62fc071ee26a8bca980053a57087` at
  `2026-04-25T20:21:44Z` on `master`; the latest was GitHub's synthetic open-PR
  merge ref `15455be0ed03f052790807ecd541e8d32a23a6ac` at
  `2026-07-14T19:08:09Z`. The synthetic merge ref established activity only; it
  is not a reviewed or accepted repository commit.
- **OBSERVED:** No product commands were rerun for this read-only usage study.
  Verification assessments below concern committed records, exact Git/GitHub
  identities, current GitHub check state, and the repository's own documented
  limitations. PR #2 records substantial local Flutter verification, but
  GitHub reports zero check runs and zero status contexts for both PR heads.
- **INFERENCE:** The repository's `master`, Worker, and Reviewer ledger actors
  and its phrases such as “separate fresh Reviewer” are strong process records,
  but GitHub cannot independently prove context freshness or actor separation.
  Both PRs have zero native review submissions and all GitHub publication and
  merge actions use the single account `Hans-Einar`.

## Repository facts

- **OBSERVED:** Repository: `Hans-Einar/weight_app_flutter`; private,
  non-archived, not a fork; default branch `master`; primary language Dart;
  created `2023-05-15T19:41:00Z`. GitHub reported `pushed_at` as
  `2026-07-14T19:07:20Z`. There are no tags or GitHub Releases, and neither of
  the three remote branches is protected.
- **OBSERVED:** The default tree has 273 regular files, including 110 under
  `SDP/`. Root authority is
  [`AGENTS.md`](https://github.com/Hans-Einar/weight_app_flutter/blob/538577830a2dbc13f603a2986f4a00f5afcc66a1/AGENTS.md),
  with the project-owned `AGENTS-project.md`. Ten installed role/task skills
  live under `.codex/skills/sdp-*`.
- **OBSERVED:** `AGENTS-project.md` is still the generic placeholder; it does
  not record actual Flutter commands, directories, hardware limits, or protocol
  constraints. Relevant project-specific guidance instead appears in
  `SDP/AGENT-REMINDERS.md`, `SDP/Instructions/ReviewAndComplianceProcess.md`,
  Slice contracts, and verification records.
- **OBSERVED:** Current default-tree lifecycle locations are
  `SDP/01--Mandate/` through `SDP/07--Implementation/`, `SDP/Instructions/`,
  `SDP/Sprints/`, `SDP/Refactors/`, `SDP/Fixes/`, `SDP/Traceability/`,
  `SDP/CodeReview/`, `SDP/Verification/`, `SDP/Releases/`, and
  `SDP/Framework/`. There is no `SDP/Features/`, no `SDP/Steering/`, no
  `CurrentAssignment`, and no `.github/workflows/`.
- **OBSERVED:** `SDP/Sprints/` contains 11 completed Sprint folders. Each has a
  Sprint record, one Iteration containing one Slice, `implementationNotes.md`,
  and `Handoff.md`; the directory contains 45 files including its README.
  `SDP/Verification/` contains 11 historical verification records and one
  tracked hardware evidence excerpt. Default `SDP/CodeReview/` contains a full
  record only for `RVW-012`; `RVW-001` through `RVW-011` exist as IDs and prose
  in Slice/verification/ledger records, not as standalone review files.
- **OBSERVED:** `SDP/Refactors/REF-001-FieldDataIntegrity/` is the only populated
  Refactor. On `master` it is `planned`, has its own Mandate, Study,
  Requirements, Architecture, Design, Implementation plan, README, and Handoff,
  and explicitly authorizes no product-code work. `SDP/Fixes/` contains only a
  generic README; no Fix was actually used. `SDP/Features/` does not exist.
- **OBSERVED:** The repository is a hybrid of two SDP generations. Commit
  `4e687ff4...` (`New version`, 2026-04-25) added the product change plus the
  populated numbered lifecycle, 11 Sprint histories, 11 verification records,
  local Instructions, and initial traceability in one 120-file commit. Its
  `SDP/README.md` says this surface was adapted from the working SDP in
  `farmStatistics`.
- **OBSERVED:** Commit `e368e317335c311bedb6a01f6ef070296c7d584d`
  (`initialized SDP`, 2026-07-14T06:21:18Z) later installed Toolkit `0.2.0`
  without replacing those project-owned records. It added `AGENTS.md`, the ten
  skills, Framework material, release/fix/refactor surfaces, project manifest,
  and a second set of unnumbered lifecycle templates beside the populated
  numbered files.
- **OBSERVED:** `SDP/Framework/installed-toolkit.manifest.yaml` identifies
  Toolkit `0.2.0`, Framework `1.0.0`, AGENTS contract `1.0.0`, installer
  `0.2.0`, installation time `2026-07-14T06:14:11Z`, and exact source commit
  `bf20832bed618ab240cf87c17517fc31ea721311`. The default project manifest is
  still an uninitialized project-owned template: project name `TBD`, version
  `0.0.0`, no release target, and all development coordinates null.
- **OBSERVED:** The duplicated lifecycle files have explicit, different
  authority. `RVW-012` states that the populated numbered files are project
  authority and the new unnumbered files are installer templates. The latter
  all say `Status: template`; they are not a newer accepted project design.
- **OBSERVED:** PR #1, “docs(sdp): record RVW-012 and plan REF-001,” was created
  `2026-07-14T07:01:58Z`, changed 14 SDP files by 700 additions/5 deletions,
  had exact head `16276ed2f6f97249d87683168346411456948fe3`, and merged without a
  native review or check as `538577830...` at `2026-07-14T07:30:41Z`. It is
  governance-only accepted/default-tree evidence.
- **OBSERVED:** PR #2, “REF-001: enforce field data integrity and state
  ownership,” was created `2026-07-14T19:05:51Z`, more than twelve hours after
  its first branch activation commit. It has 17 commits, 63 changed files,
  9,449 additions/631 deletions, exact head `c4c31cbae...`, and remains open and
  draft. It has no comments, review comments, native reviews, check runs, or
  status contexts.

## How SDP is used in this repository

### Work start and GitHub Issue authority

- **OBSERVED — historical:** The April work was request-driven and Sprint-first.
  Sprint notes repeatedly begin with “Request,” “Bug Report,” or a user-reported
  runtime observation, then record the Master's interpretation. GitHub has no
  Issue containing those requests. All eleven Sprints and the corresponding
  application changes first appear together in commit `4e687ff4...`, so Git
  cannot reconstruct an activation branch, per-Sprint candidate, or review head
  for that work.
- **OBSERVED — current default:** `SDP/README.md` says work starts from
  `SDP/Traceability/CurrentIndex.yaml` and then the active Sprint folder. On
  `master`, `active.sprint`, `active.iteration`, and `active.slice` are all null;
  `REF-001` is only planned. PR #1 deliberately cleared the stale completed
  Sprint-011 coordinates and made owner authorization a required boundary.
- **OBSERVED — open/unmerged:** PR #2's README and implementation notes say the
  owner authorized all five REF-001 Slices, branch publication, and an unmerged
  PR on 2026-07-14. That authorization is durable only as repository prose and
  ledger entries; there is no Issue, issue comment, Steering decision record, or
  GitHub review carrying the source instruction.
- **INFERENCE:** The repository is excellent at reconstructing a contract after
  activation but weak at proving the durable request that activated it. A fresh
  agent can see what the Master says was authorized, but cannot inspect the
  owner's exact assignment or determine whether later prose paraphrased it.
- **OWNER DIRECTION:** [Hans-Einar/SDP#5](https://github.com/Hans-Einar/SDP/issues/5)
  makes a GitHub Issue the future operational assignment boundary and makes a
  bounded Codex session Master for that Issue. The no-Issue history here is a
  migration gap, not evidence to preserve conversational authority.
- **RECOMMENDATION:** Future work should start from an exact Issue/comment and a
  compact repository assignment binding before broad changes. The Issue should
  carry the full goal, why-now, scope, non-goals, baseline, invariants, safety,
  verification, review, PR, stop, and reporting contract. The repository should
  reference that authority and mirror only machine-useful assignment fields.

### Codex Master, delegation, Architect, Verifier, and Reviewer

- **OBSERVED:** Root `AGENTS.md`, `SDP/AGENT-REMINDERS.md`, and
  `AgentLifecycleAndDelegation.md` establish that a non-delegated session is
  Master. The Master identifies active coordinates, refines documents and
  traceability before product work, delegates product code to a fresh Worker,
  delegates review to a separate fresh Reviewer, inspects real verification,
  maintains notes/handoff/index/relations/ledger, and stops at the Slice/Fix
  boundary.
- **OBSERVED:** The installed Master, Worker, Reviewer, Architect, Verifier,
  Traceability, and Vertical-Refactor skills sharpen those boundaries. Worker
  instructions constrain files, scope, verification, and discoveries. Reviewer
  instructions require a fresh read, actual diff, exact findings, manifests,
  evidence, and a disposition. The Verifier requires evidence to apply to the
  exact commit. The Architect owns durable boundaries rather than product-code
  implementation.
- **OBSERVED — historical:** Sprint notes and the 79-line default Ledger record
  Worker start/completion, Reviewer completion, rework, Master verification, and
  closure. Sprint 011 shows review finding a queued-poll race, focused Worker
  rework, re-verification from 24 to 26 tests, and final review. These records
  show the role loop was practiced before the formal Toolkit skills were
  installed, although the historical review records are not exact-head-pinned.
- **OBSERVED — open/unmerged:** REF-001 makes delegation materially stronger.
  It activates one Slice at a time, commits the contract before implementation,
  uses exact candidate commits, runs fresh review passes that return
  `changes required`, assigns bounded rework, re-verifies, and uses another
  fresh review before advancing. Review found substantive defects after green
  tests: lifecycle-hidden known-red tests, stale-generation serial writes,
  cross-wrapper lost updates, delayed UI intent rollback, post-unmount
  `setState`, and duplicate persistent IDs.
- **OBSERVED:** No distinct Architect execution record exists. Architecture was
  authored inside `RVW-012`/REF-001 planning and later Master reconciliation.
  A separate Verifier identity is also not persisted for most work; the Master
  verifies and Reviewers often rerun evidence. The roles exist as skills more
  clearly than as durable assignment objects.
- **INFERENCE:** Bounded Worker and fresh Reviewer separation is this
  repository's strongest agent-process evidence. Architect/Verifier separation
  should be required based on risk rather than turned into empty ceremony, but
  exact actor/task/output identities would improve auditability.
- **RECOMMENDATION:** Canonical Issue-Master guidance should preserve the exact
  loop demonstrated by REF-001: activate one bounded contract, Worker result,
  independent exact-head review, rework by a Worker context, re-review, real
  verification, durable evidence, then stop or advance only within explicit
  Issue authority. Record a small delegation identity, owned paths, expected
  output, and reviewed commit for material Worker/Architect/Verifier/Reviewer
  passes rather than trusting actor names in a chat-derived ledger.

### Branches and draft PRs

- **OBSERVED:** The accepted governance plan used branch
  `steering/code-review-001-refactor-plan` and PR #1 against `master`. REF-001
  execution uses `refactor/ref-001-field-data-integrity` from exact baseline
  `538577830...` and open draft PR #2 against `master`.
- **OBSERVED:** REF-001's branch has progressive commits that correspond to
  activation, Worker candidates, review rework, Slice closeout/next activation,
  governance reconciliation, exact-candidate attestations, and PR publication.
  This is much more auditable than April's single 120-file commit.
- **OBSERVED:** PR #2 was opened only after all four product Slices and the
  governance candidate had passed their local review/automated gates. The
  tested governance candidate is `a9f615173...`; later commits
  `8672ed612...` and `c4c31cbae...` add review/verification/publication
  attestations. Seven SDP files differ between the tested candidate and current
  PR head, while product code and tests do not.
- **OBSERVED:** Repository prose carefully states that PR #2 is open, draft,
  unmerged and blocked on physical HEOS/Bluetooth evidence. However, branch/PR
  identity is repeated manually in README, Handoff, implementation notes,
  Ledger, PR body, and verification/review records; neither `CurrentIndex` nor
  `Relations` owns a stable GitHub binding.
- **CONTRADICTORY EVIDENCE:** The workflow values exact-head review, but the
  actual GitHub PR has no CI, native review, branch protection, or comment-based
  terminal attestation. Current PR head is evidence-only beyond the tested
  candidate and its state is established through live GitHub inspection, not a
  repository-enforced gate.
- **RECOMMENDATION:** Open the draft PR immediately after the assignment and
  activation commit so review, CI, and milestone evidence accumulate on the
  durable surface. Store one stable binding among Issue, assignment,
  Feature/Refactor, baseline, branch, and PR; query mutable draft/check/review/
  merge state from GitHub or capture timestamped snapshots rather than copying
  it across many prose files.

### Feature, Refactor, Sprint, Iteration, vertical Slice, and Fix

- **OBSERVED — historical:** Sprints are the conceptual owners of every April
  capability and correction: the MVP, Bluetooth overflow fix, reconnect/debug,
  no-data diagnostics, newline framing, firmware retest, quick save/recall,
  workflow UX, dialog table, statistics, and settings/terminal/zero. Each Sprint
  has exactly one Iteration and one Slice. Several Sprints are genuine Features,
  several are Fixes, and Sprint 006 is verification-only, but the hierarchy does
  not distinguish those semantics.
- **OBSERVED:** Historical Slices often are usable end-to-end outcomes, but
  their verticality varies. Sprint 008 crosses Bluetooth state, project
  persistence, multiple screens, shared widgets, and tests and is a broad
  vertical capability. Sprint 005 changes a narrow protocol seam. Sprint 006
  changes no product code. Thus “Slice” means bounded work unit more often than
  strictly vertical product behavior.
- **OBSERVED — current default:** `REF-001` is a first-class owner separate from
  Sprint. It contains one `REF-001-IT-001` and five planned Slices. It has its
  own local Study, Requirements, Architecture, Design, and implementation plan.
  On `master` it is planned and owner-gated.
- **OBSERVED — open/unmerged:** REF-001 Slices 2–4 are strong vertical refactor
  examples: they change behavior through state models/services/repositories/UI
  and tests while keeping the app runnable and compatibility explicit. Slice 1
  is a baseline/contract-test gate and Slice 5 is governance/closeout, so even
  the strongest current practice still overloads Slice with non-product
  enabling and closure stages.
- **OBSERVED:** `SDP/Fixes/README.md` defines a bounded correction path, but no
  Fix record exists. Small corrections became full one-Iteration/one-Slice
  Sprints. `SDP/Features/` is absent. REF-001 calls calibration and raw serial
  logging later “feature work” and records feature seeds, but no canonical
  Feature identity, folder, requirements, Issue, release inclusion, or
  completion contract exists.
- **INFERENCE:** Iteration provides little observed value here: eleven historical
  Sprints wrap one Iteration/one Slice each, and REF-001 wraps five sequential
  Slices in one Iteration completed during one day. It is an ID layer rather
  than an observed planning/timebox boundary.
- **OWNER DIRECTION:** Issue #5 requires implementation work to belong to a
  Feature or Refactor, permits a true small Fix, and retains a Slice as a small
  vertical end-to-end capability. Sprint must not remain capability owner only
  for backward compatibility.
- **RECOMMENDATION:** Preserve the completed Sprint history without
  relabelling it. Prospectively use `Feature | Refactor | Fix` as the semantic
  owner. Make Sprint an optional execution/timebox grouping and Iteration
  optional unless it records a real replanning or feedback cycle. Define a
  vertical Slice as the smallest runnable, independently verifiable behavior
  through all affected horizontal layers; treat baseline research, contract
  preparation, governance closeout, and release work as gates/tasks unless they
  themselves deliver an end-to-end operational capability.

### Initial horizontal lifecycle and later requirements/design evolution

- **OBSERVED:** The populated initial Mandate, Study, Requirements,
  Architecture, Design Analysis, and Design establish a useful horizontal
  skeleton: field workflow, Bluetooth protocol ownership, project repository,
  versioned JSON, thin UI screens, operation-mode boundary, offline use, and
  field ergonomics. Sprint 001 then implements a vertical MVP across those
  layers.
- **OBSERVED:** Product growth did not stop at the initial plan. Sprints 2–11
  introduced requirements `R-013` through `R-031`, including diagnostics,
  auto-reconnect, protocol framing, recall, wakelock, richer project UX,
  terminal ownership, and zero-offset commands. Those requirements lived in
  CurrentIndex, Relations, Sprint contracts, notes, verification, and code, but
  the numbered Requirements file on default `master` defines only `R-001`
  through `R-012` and `R-019` through `R-020`.
- **OBSERVED:** `RVW-012` correctly classified that mismatch as Medium. It also
  found stale/unrealized Architecture and Design claims. PR #1 cleared active
  coordinates and planned a refactor but left the numbered requirement gap for
  later bounded governance reconciliation.
- **OBSERVED — open/unmerged:** REF-001 demonstrates a better evolutionary
  method. It introduces local `REF-001-RQ-*` and `REF-001-QR-*`, compares
  alternatives in a local Study, changes state ownership in local Architecture
  and Design, delivers vertical Slices, then reconciles only stable project-wide
  numbered Requirements/Architecture/Design during SL-005. It explicitly does
  not line-by-line mirror the code.
- **INFERENCE:** The repository proves that later requirements and design
  refinement are necessary and productive; the failure was not evolution but
  the lack of a canonical Feature/local-change authority and a clear rule for
  when stable horizontal documents must be revised.
- **RECOMMENDATION:** Keep the initial lifecycle for a minimum stable horizontal
  skeleton, not a complete future-product plan. Let each Feature or Refactor own
  bounded local Study/requirements/design. Revise project Architecture only
  when stable boundaries, ownership, contracts, persistence/integration, or
  cross-cutting constraints change. Keep implementation detail local to the
  Slice and reconcile stable project records at acceptance rather than allowing
  `CurrentIndex` to become a shadow requirements document.

### Studies and convergence

- **OBSERVED:** The initial project uses one short broad Study. REF-001 uses one
  bounded local Study that records observed failure mechanisms, persistence and
  measurement alternatives, compatibility rules, and a recommendation. It
  converts an independent review into one coherent Refactor programme.
- **OBSERVED:** `RVW-012` itself is cross-repository: it inspects exact
  application/governance baselines plus HEOS `main` and LogParser `main`, then
  separates two later Feature seeds from the corrective Refactor. This shows a
  Steering/review pass performing bounded evidence synthesis across repository
  boundaries.
- **OBSERVED:** There is no set of parallel independent `Study` records, no
  claim/evidence-level schema, and no formal convergence gate. REF-001's five
  Slices converge through the programme plan and sequential reviews, not through
  independent Study acceptance.
- **RECOMMENDATION:** Carry forward a single bounded Study for coherent
  questions such as this persistence/measurement refactor. Permit multiple
  Studies and a convergence decision when evidence owners, safety domains, or
  uncertainties can be investigated independently. Do not import ActionCam's
  physical-safety claim levels wholesale; this repository's generic lesson is
  simply to make cross-repository assumptions, hardware limitations, and the
  decision that authorizes implementation explicit.

### Independent review, verification, and completion gates

- **OBSERVED:** The April corpus records local format/analyze/test/APK evidence,
  hardware sessions for the Bluetooth serial path, Reviewer conclusions, and
  residual risks. It truthfully records when hardware failed, when no bytes were
  received, and when later firmware returned 71 parsed weights. This is valuable
  operational evidence.
- **OBSERVED:** Those historical records do not name an exact verified commit,
  and the 11 Sprint histories first appear together with their product changes
  in `4e687ff4...`. `RVW-012` therefore refused to treat them as reproducible
  release evidence. Open PR #2 later classifies all eleven as
  `historical_unpinned`, with `verified_commit: null`, rather than falsely using
  their Git first-appearance commit as the verification identity.
- **OBSERVED — open/unmerged:** REF-001 verification is exact-commit and
  environment-bound. Each Slice records Flutter/Dart/OS, commands, test counts,
  APK build, exact commit, Reviewer disposition, and hardware limitation. The
  final automated candidate passed 152 full tests, 113 opt-in REF contracts,
  APK build/hash, diff hygiene, YAML/NDJSON parsing, relation resolution, and
  indexed-path checks.
- **OBSERVED:** Independent review repeatedly blocked apparently green
  candidates and produced actionable High/Medium findings. This materially
  improved state ownership and recovery contracts. The method also correctly
  keeps `REF-001`, its Iteration, and SL-005 active because physical
  HEOS/Bluetooth scenarios are unavailable; it does not claim completion,
  merge, release, or tag.
- **CONTRADICTORY EVIDENCE:** The installed Toolkit `0.2.0` project validator
  fails on the open branch with 843 compatibility errors: nine CurrentIndex,
  twelve Relations, and 822 legacy-ledger envelope errors. The branch records
  this as a migration warning rather than claiming validation passed. This is
  honest evidence, but it means the current installed method and the
  repository's authoritative traceability cannot fully validate each other.
- **RECOMMENDATION:** Preserve exact-head evidence, independent review that can
  block on Medium findings, changed-head re-review, artifact identity, explicit
  historical-evidence classes, and truthful unavailable/manual gates. Add CI
  and native/check-based PR evidence where practical, while retaining repository
  review records for reasoning that CI cannot express.

### Handoff and recovery after session loss

- **OBSERVED:** Every historical Sprint has a concise Handoff with its
  coordinates, resulting behavior, verification, residual risks, and next
  field action. REF-001 has a programme Handoff with exact accepted commits,
  active Slice, outstanding hardware scenarios, branch/PR identity, migration
  warning, and stop boundary. Together with implementation notes, this enables
  a new Master to resume without relying on chat memory.
- **OBSERVED:** Historical Handoffs all retain headings such as “Active sprint:
  SPR-007” even though the global CurrentIndex is now clear. They are accurate
  snapshots of their local closeout context but misleading if treated as live
  state. Eleven folders also duplicate verification status and next action that
  may later become stale.
- **RECOMMENDATION:** Retain a concise prose Handoff for context, decisions,
  residual risk, and next authorized action. Treat it as a timestamped work-item
  snapshot, not a second current-state registry. Generate active coordinates,
  exact Issue/branch/PR/check/merge state, and evidence links from a compact
  assignment/index plus live GitHub data.

### CurrentIndex, Relations, and Ledger

- **OBSERVED:** Default `CurrentIndex.yaml` is both the active pointer and a
  large registry for one Mandate, nine Use Cases, 31 requirements, six
  architecture IDs, one Refactor, eleven Sprints/Iterations/Slices, twelve
  reviews, and eleven verification records. PR #1 fixed its stale active
  coordinate, but the file still contained 17 requirement IDs without canonical
  definitions in the numbered Requirements document.
- **OBSERVED:** `Relations.yaml` explicitly connects Mandate to use cases,
  requirements to architecture, and each historical Slice to review and
  verification. PR #1 adds `RVW-012 -> REF-001` and preserved requirements.
  This makes “why does this work exist?” recoverable, but many relations depend
  on manually maintained IDs that lacked primary records.
- **OBSERVED:** The default Ledger has 79 append-only NDJSON events: Slice
  creation, Worker start/completion, rework, review, verification, closure,
  Steering review, and Refactor planning. The open branch grows it to 137 lines
  with detailed exact-Slice transitions, audit starts, migration warnings,
  hardware gates, branch push, and PR creation.
- **OBSERVED:** The Ledger provides a readable chronology, but it duplicates
  implementation notes, verification, review, CurrentIndex, and GitHub state.
  Its historical envelope predates installed Toolkit 0.2 schemas. The open
  branch correctly refuses either to rewrite immutable history or to rename
  repository-authoritative REF IDs merely to make the new validator green.
- **INFERENCE:** CurrentIndex is valuable but overburdened as both “current” and
  a complete registry. Relations keeps important justification but can mask
  missing source documents. Ledger is excellent for material recovery/review
  decisions and too detailed when every agent lifecycle transition is manually
  repeated.
- **RECOMMENDATION:** Keep a small current assignment/index and stable graph for
  `Issue <-> Feature/Refactor/Fix <-> important Study/requirement/architecture
  decisions <-> Slice <-> review/verification <-> PR/commit <-> release`.
  Generate registries, paths, GitHub status, and routine Worker/PR events when
  possible. Keep append-only events for authorization, decision, rework,
  acceptance, blocked gate, migration, correction, merge, and release rather
  than every conversational step.

### Steering assignment and GitHub binding

- **OBSERVED:** Steering exists as a role but not as a canonical surface.
  `RVW-012` identifies `Reviewer: GPT Steering Group`; the Ledger records actor
  `steering-group`; branch `steering/code-review-001-refactor-plan` and PR #1
  convert its changes-required disposition into accepted default-tree
  governance. That pass correctly chooses REF-001 remediation while leaving
  calibration and raw logging as later Features.
- **OBSERVED:** There is no `SDP/Steering/`, no `CurrentAssignment.yaml`, no
  Master report, no Steering decision object, and no GitHub Issue. Branch, PR,
  authorization, active Refactor/Slice, stop boundary, hardware gate, and
  evidence heads are spread among REF README/Handoff/notes, CurrentIndex,
  Ledger, review/verification, and PR prose.
- **INFERENCE:** `RVW-012` is strong evidence for the value of a Steering role:
  it inspects repository and cross-repository evidence, distinguishes Refactor
  from future Features, establishes severity and safety gates, and stops before
  product changes. It is weak evidence for assignment durability because the
  underlying owner request and acceptance are not independently addressable.
- **OWNER DIRECTION:** Issue #5 places a ChatGPT Steering Group above bounded
  Codex Issue-Masters while repository/GitHub evidence remains authoritative.
- **RECOMMENDATION:** Add a compact assignment record binding an immutable
  Issue/comment authority, work type and ID, exact base, branch/draft PR,
  authorized Slices, required Studies/decisions, invariants, permissions,
  verification/review gates, stop condition, and current disposition. Keep
  mutable GitHub state generated or as timestamped observations. Steering
  decisions that change authorization or acceptance should be durable and
  linked, not implied by a role name in the Ledger.

### Automated versus conversational convention

- **OBSERVED — automated locally:** Flutter format, analyze, test, APK build,
  artifact hash, diff hygiene, YAML/NDJSON parse, relation/path resolution, and
  strict-schema validation are recorded. REF-001's manifest-driven known-red
  test gate automatically enables each defect contract when its fixing Slice is
  active; this is a particularly useful connection between SDP state and
  executable behavior.
- **OBSERVED — not automated on GitHub:** There is no Actions workflow, branch
  protection, required review, check run, or status context. PR #1 merged with
  none; PR #2 remains open with none. GitHub Issue/Feature/Refactor/branch/PR
  binding and stale-state detection are not automated.
- **OBSERVED — conversational/manual:** owner authorization, Master/Worker/
  Reviewer selection, branch/PR creation, CurrentIndex/Relations/Ledger updates,
  Handoff, review identity, and advancement to the next Slice are manually
  integrated into Markdown/YAML/NDJSON.
- **RECOMMENDATION:** Canonical tooling should validate the assignment and
  Feature/Refactor graph, compare declared identity with live GitHub, classify
  legacy generations, and generate safe derived state. Project-specific Flutter
  and hardware checks should stay project-owned; the generic method should
  require truthful evidence shape and exact identity, not hard-code these
  commands.

### ChatGPT Steering Group versus a general SDP skill

- **OBSERVED:** The GPT Steering Group's useful work in `RVW-012` is governance:
  independent repository assessment, cross-repository context, severity
  disposition, selection of a Refactor, separation of later Feature seeds,
  authorization boundary, and physical-gate awareness. It does not directly
  implement product code.
- **OBSERVED:** The reusable SDP knowledge needed by every role is different:
  reading AGENTS/manifests/current state, interpreting IDs, deciding which
  records are authoritative versus templates, drafting bounded contracts,
  maintaining Feature/Refactor/Slice relations, validating evidence/handoff,
  and detecting Toolkit-generation incompatibility.
- **RECOMMENDATION:** Keep separate skills. A ChatGPT `steering-group` skill
  should inspect portfolio/repository/PR/CI/review evidence, make or escalate
  product/architecture/authorization decisions, create/refine Issues, accept or
  reject results, coordinate cross-repository gates, and stop at decision
  boundaries. A general ChatGPT `sdp` skill should explain and inspect the
  method, bootstrap/read state, help author Issue/Feature/Refactor contracts,
  understand IDs and migrations, validate readiness/handoff conceptually, and
  surface contradictions. It must not silently make Steering decisions or
  treat ChatGPT memory as authority.
- **RECOMMENDATION:** Codex role skills should gain explicit Issue authority,
  Feature support, assignment binding, and exact PR/check/head handling. The
  current Master/Worker/Reviewer/Verifier boundaries should remain; the general
  SDP skill supplies method knowledge, while project `AGENTS-project.md`
  supplies Flutter/hardware commands and constraints.

## What worked well

- **OBSERVED:** The repository evolved quickly through real field feedback.
  Eleven bounded work records cover an MVP, transport debugging, device
  verification, tractor ergonomics, project workflow, and serial tooling in
  roughly 2026-04-24 through 2026-04-25. Product development demonstrably
  continued beyond the initial Design rather than freezing inside it.
- **OBSERVED:** The initial horizontal boundaries—Bluetooth service, project
  repository, versioned export, thin screens, operation modes—were strong
  enough to support later vertical work. `RVW-012` found drift and ownership
  defects, but later Sprints did not require a wholesale rewrite of the product
  concept.
- **OBSERVED:** Slice contracts consistently include goal, why now, expected
  modules, invariants, non-goals, traceability, verification, and completion
  signal. Field-specific safety and usability constraints are repeated where
  they matter, supporting Worker autonomy without silent scope expansion.
- **OBSERVED:** Review is not ceremonial. Historical Sprint 011 and open
  REF-001 reviews find concrete asynchronous, state-ownership, and persistence
  defects after tests pass; rework is bounded and independently re-reviewed.
- **OBSERVED:** REF-001 is an exemplary evolutionary Refactor structure. It
  begins from an independent finding set, creates local Study/requirements/
  architecture/design, delivers compatibility-preserving vertical changes, and
  reconciles stable project records only after accepted implementation evidence.
- **OBSERVED:** Exact-commit verification and candid evidence classification in
  PR #2 materially improve recoverability. The branch does not invent historical
  commit identities, hardware success, Toolkit compatibility, CI, merge,
  release, or tags.
- **OBSERVED:** Known-red contracts are lifecycle-aware and become ordinary
  tests in their fixing Slice. This turns review findings into executable gates
  without hiding them permanently or requiring hardware for every defect.
- **OBSERVED:** Handoff, implementation notes, CurrentIndex, Relations, and the
  Ledger make the complex open Refactor resumable after context loss. Exact
  accepted heads and pending hardware work are explicit.
- **OBSERVED:** Cross-repository review is bounded well. HEOS and LogParser are
  evidence sources; `Version_1.0--Union`, calibration, and raw logging are not
  opportunistically pulled into REF-001.
- **OBSERVED:** The open branch stops honestly at a real physical gate. Keeping
  the draft PR open and active coordinate uncleared is better evidence practice
  than converting automated confidence into a false field-complete claim.

## Pain points and accidental complexity

- **OBSERVED:** There is no durable GitHub Issue authority. The owner request,
  approval to execute five Slices, and expected acceptance boundary cannot be
  inspected independently from the Master's own repository summary.
- **OBSERVED:** The April history compresses eleven work loops, all product
  changes, 77 initial ledger events, and 11 verification records into one
  120-file commit. The repository documents rich progressive activity that Git
  cannot independently correlate to per-Slice code/review heads.
- **OBSERVED:** Sprint, Iteration, and Slice are nearly always 1:1:1. This adds
  three IDs and four per-Sprint documents even for one bug fix or verification
  pass, while failing to distinguish Feature, Refactor, Fix, and non-product
  gate semantics.
- **OBSERVED:** New requirements were added to CurrentIndex/Relations and Slice
  prose but not the canonical numbered Requirements file. The active registry
  looked complete while its declared source document was incomplete. Review had
  to reverse-engineer and reconcile later behavior.
- **OBSERVED:** Toolkit installation produced duplicate numbered project
  authority and unnumbered template files. Review had to state which set wins;
  ordinary file discovery alone could easily select the wrong document.
- **OBSERVED:** `AGENTS-project.md` is an empty template while project-specific
  commands and constraints are spread through several SDP files. A fresh agent
  must assemble the actual Flutter/hardware operating contract indirectly.
- **OBSERVED:** Live state is manually duplicated across CurrentIndex,
  Relations, Ledger, Sprint/Refactor README, implementation notes, Handoff,
  review, verification, PR body, and GitHub. Historical Handoffs retain local
  “active” labels, and PR #2 was not opened until near completion.
- **OBSERVED:** Review evidence is committed but not GitHub-enforced. There are
  no native review submissions, CI checks, required contexts, or protected
  branches. A single GitHub identity authors, publishes, and merges.
- **OBSERVED:** Default CodeReview lacks standalone records for RVW-001 through
  RVW-011. CurrentIndex/Relations claim stable review objects whose detailed
  findings and exact inspected commits cannot be retrieved as primary files.
- **OBSERVED:** Traceability is detailed but internally generation-split. The
  repository's 137-line open ledger and REF IDs are authoritative, yet its
  installed Toolkit validator rejects them with 843 schema errors. Neither a
  broad rewrite nor permanent incompatibility is an acceptable end state.
- **OBSERVED:** `CurrentIndex` is both current-state pointer and historical
  registry. `Relations` can connect an ID that lacks a canonical definition.
  Ledger lifecycle events repeat detailed notes without supplying an external
  assignment identity.
- **OBSERVED:** Hardware gates are truthfully preserved, but the older
  verification records mix authentic device observations with unpinned
  application state. They are useful history and not release-grade proof.
- **OBSERVED:** Open PR #2's product work and reviews were completed in one
  intense day, yet the PR remains unmerged more than a month later with no
  comment/check/review activity. Repository state is recoverable, but no
  machine or Steering surface signals who owns the next physical-gate action.

## Carry forward

- **RECOMMENDATION:** Repository-authoritative bounded contracts with explicit
  goal, why now, scope, expected areas, invariants, non-goals, verification,
  review, completion, and stop conditions.
- **RECOMMENDATION:** Master-per-Issue coordination with bounded Workers,
  independent exact-head Reviewers, risk-based Architect/Verifier passes,
  evidence inspection by the Master, and no automatic scope expansion.
- **RECOMMENDATION:** Exact baseline/candidate/review/verification/PR identities,
  progressive commits, changed-head re-review, candid limitations, and explicit
  historical evidence classes.
- **RECOMMENDATION:** Feature/Refactor-local Study, requirements, architecture,
  and design refinement followed by small runnable vertical Slices and selective
  reconciliation of stable project-wide contracts.
- **RECOMMENDATION:** REF-001's state-ownership, compatibility, known-red test,
  failure-injection, exact-commit evidence, rework, and hardware-gate discipline
  as a strong generic Refactor example.
- **RECOMMENDATION:** A concise recovery Handoff plus a compact current
  assignment/index and stable relations among Issue, work owner, important
  decisions, Slice, evidence, PR/commit, and release.
- **RECOMMENDATION:** Append-only records for material authorization, decision,
  review/rework, blocked gate, acceptance, migration, merge, and release
  transitions, with routine/live fields generated where possible.
- **RECOMMENDATION:** A distinct ChatGPT Steering Group role for assignment,
  cross-repository/product/architecture decisions and acceptance, plus a
  separate general SDP method skill for inspection, authoring, validation, and
  migration support.
- **RECOMMENDATION:** Preserve project-local physical constraints—500 ms
  polling, protocol framing, HEOS/Bluetooth scenarios, field ergonomics—in
  project authority. Canonical SDP should require their explicit treatment but
  not generalize their content to unrelated repositories.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward conversational owner authorization
  or a Ledger actor label as sufficient assignment authority. Use an exact
  GitHub Issue/comment and repository assignment binding prospectively.
- **RECOMMENDATION:** Do not preserve Sprint as capability owner or mandatory
  `Sprint -> Iteration -> Slice` ceremony for every Feature, bug, verification
  pass, or governance task. Preserve existing IDs as history; change the
  prospective semantic model.
- **RECOMMENDATION:** Do not call every bounded activity a vertical Slice.
  Distinguish end-to-end product/refactor Slices from Study, baseline,
  verification, governance, migration, and release gates.
- **RECOMMENDATION:** Do not treat `CurrentIndex` as a substitute Requirements
  document or allow later requirements to exist only in indexes and Sprint
  prose. Give them Feature/Refactor-local authority and revise stable project
  requirements deliberately.
- **RECOMMENDATION:** Do not keep duplicate populated/template lifecycle files
  indefinitely without machine-readable generation and authority markers. Do
  not delete or overwrite project-owned history; migrate additively with an
  explicit mapping.
- **RECOMMENDATION:** Do not mirror live branch/PR/check/merge state across many
  prose records. Generate it from GitHub, timestamp observations, and flag
  contradictions.
- **RECOMMENDATION:** Do not require a manually authored Ledger event for every
  Worker lifecycle message, and do not discard the useful append-only history.
  Retain only decision/recovery/audit-significant events and generate routine
  activity.
- **RECOMMENDATION:** Do not infer independent review or exact verification from
  role prose, green tests on another commit, or historical file first-appearance.
  Preserve the `historical_unpinned` distinction and add exact, attributable
  evidence for new work.
- **RECOMMENDATION:** Do not treat open PR #2 as accepted default behavior,
  despite its unusually strong internal evidence. Its implementation,
  governance reconciliation, active coordinate, and Toolkit migration warning
  remain unmerged until the physical/Steering decision is made.
- **RECOMMENDATION:** Do not generalize this app's Bluetooth, tractor UI,
  persistence backend, HEOS protocol, calibration, or raw-logging rules into
  canonical SDP. Generalize the bounded-contract, cross-repository evidence,
  compatibility, exact-verification, and truthful physical-gate patterns.
