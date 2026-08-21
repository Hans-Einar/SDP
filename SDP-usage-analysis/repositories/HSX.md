# HSX SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies merged/default truth at the exact current
  default-branch identity `main@e374da88f4dd470bad2d8ec6a2a14f1ce367e40e`
  (tree `fa16ce642f437e1fe8f9be0398a9424f6d389932`, committed
  `2025-11-09T17:55:00Z`). A fresh mirror and `git ls-remote --symref` on
  2026-08-21 both resolved `HEAD -> refs/heads/main` to that commit.
- **OBSERVED:** Issue #5's inventory separately captured 82 advertised refs,
  544 reachable commits and 104 qualifying commits in the fixed inclusive
  window `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`. The first was
  `de84aafd732ec67f53f322481f75d77a17940a0a` at
  `2026-08-19T15:31:33Z`; the last was
  `86aac69ac89f20d8ccce0fe83fd15af0daed8d9d` at
  `2026-08-20T20:55:36Z`. Neither is reachable from `main`; both are
  non-default branch evidence.
- **OBSERVED:** Advertised refs are mutable. The live 2026-08-21 mirror had 83
  refs and 545 unique reachable commits, while reproducing the same 104
  qualifying commits by enumerating every reachable commit before filtering
  its UTC committer timestamp. Later live tips included
  `codex/dbg-rf-002-003@69a54aeb3394d3cd4792bce620748e15bab69f1f`
  (`2026-08-20T21:44:53Z`) and
  `codex/dbg-rf-004@82154c614a31284723bf3e6a337c5bedfb8aba5d`
  (`2026-08-20T22:23:26Z`), both after the frozen window end. They are useful
  current-work evidence but did not add to the 104-commit inclusion count.
- **OBSERVED:** The relevant recent branch chain starts from
  `Implementation/vscode@a1daa1c62605c44ac67e58e2b71320006f73cdd9`,
  not `main`. The important advertised checkpoints are
  `sdp/debugger-gap-analysis@e5a50ab45acdcb515ccd3602ce99487bd668cdfd`,
  `codex/dbg-rf-001@977da1c19117e805a341a6055d881eefaff58ddf`,
  `codex/dbg-da-001@c0003d070c6840f0d55487d73b227880cfd27494`,
  `codex/dbg-st-006@bf92c9be6cf81a7cb704778dafe88e55fee2e235`,
  `codex/dbg-rf-002-003@69a54aeb3394d3cd4792bce620748e15bab69f1f`
  and `codex/dbg-rf-004@82154c614a31284723bf3e6a337c5bedfb8aba5d`.
  No one of these is merged to the default branch.
- **OBSERVED:** GitHub Issues, comments, PR metadata, formal review objects and
  check rollups were inspected on 2026-08-21. The latest practice is visible
  through one owner identity, `Hans-Einar`; repository evidence cannot prove
  which human or agent session acted behind that account.
- **OBSERVED:** No tests were rerun for this read-only usage study. Test counts
  below are claims in committed verification/review records and owner Issue
  comments, not newly reproduced results.
- **INFERENCE:** “Repository is active” and “default branch is current” are
  different predicates. HSX satisfies the first because actual advertised-ref
  commits fall inside the window. It fails the second because `main` predates
  the whole evolved SDP process by more than nine months.
- **RECOMMENDATION:** Canonical studies should always record both identities:
  the default-branch commit as merged/canonical tree truth, and exact
  non-default heads as explicitly unmerged activity evidence. Active-repository
  inclusion must not promote branch-local practice to accepted repository-wide
  truth.

## Repository facts

- **OBSERVED:** Repository: [`Hans-Einar/HSX`](https://github.com/Hans-Einar/HSX),
  public, unarchived, non-fork; default branch `main`.
- **OBSERVED:** Default study commit
  `e374da88f4dd470bad2d8ec6a2a14f1ce367e40e` is merge commit for PR
  [#34](https://github.com/Hans-Einar/HSX/pull/34). It contains no `SDP/`
  directory. Its process surface is the older phase-oriented
  `main/02--Study/`, `main/03--Architecture/`, `main/04--Design/`,
  `main/05--Implementation/` plus top-level `refactor/` packages.
- **OBSERVED:** At `main`, root authority is lowercase `agents.md`; additional
  scopes include `main/AGENTS.md`,
  `main/05--Implementation/01--GapAnalysis/AGENTS.md` and per-module files such
  as
  `main/05--Implementation/01--GapAnalysis/09--Debugger/AGENTS.md`.
  `main/AGENTS.md` still says Design is active and Implementation has not
  begun, contradicting the extensive implementation history already present.
- **OBSERVED:** The default tree's older implementation method uses one grand
  GapAnalysis plan, eleven module-level Study/ImplementationPlan/
  ImplementationNotes/`04--git.md` packages, design/implementation/integration/
  comprehensive review gates, and a top-level `refactor/#0_template/` sequence
  of context -> objectives -> analysis -> requirements -> design -> DoD ->
  implementation -> Git summary. It has no CurrentIndex, Relations, Ledger,
  Steering assignment, stable Feature IDs, Sprint/Iteration/Slice hierarchy or
  formal Verification directory.
- **OBSERVED:** SDP is present only on recent non-default branches. At
  `82154c...`, the top-level SDP tracks are `SDP/HSX/`, `SDP/Debugger/`,
  `SDP/AVR/` and `SDP/Shared/`. The active evidence surfaces include:
  `SDP/Shared/Process.md`; numbered Study, Requirements, Architecture,
  DesignAnalysis and Design documents; `SDP/Debugger/CodeReview/`;
  `SDP/Debugger/GapAnalysis/`; singular `SDP/Debugger/Refactor/`;
  `SDP/Debugger/Sprints/`; `SDP/Debugger/Verification/`; per-track
  `Traceability/CurrentIndex.yaml`, `Issues.yaml`, `Relations.yaml` and
  `Ledger.ndjson`; and Handoffs under the active Sprint and HSX Study.
- **OBSERVED:** The new tracks are selective rather than a complete copied
  lifecycle. `SDP/Debugger/` has `02--Study` through `06--Design` plus
  CodeReview, GapAnalysis, Refactor, Sprints, Traceability and Verification;
  `SDP/HSX/` has Study, Requirements, Architecture and Design plus Reviews,
  Traceability and Verification; `SDP/AVR/` currently has Traceability only.
  There is no new-track Mandate or Implementation phase directory; legacy
  foundation/provenance remains under `main/`, and product code remains in the
  ordinary source tree.
- **OBSERVED:** There is no `SDP/Features/`, `SDP/Steering/`,
  `CurrentAssignment`, `Instructions/`, `AGENTS-project.md`, `SDP.manifest.yaml`
  or project manifest on the recent branches. Apparent SDP generation/version
  is therefore indeterminate. The first multi-track process commit is
  `de84aafd...`; this is a project-local evolution, not evidence of an installed
  released Toolkit version.
- **OBSERVED:** Root `AGENTS.md` on the recent branches combines the old HSX
  build/product constraints with new Master/Worker/Reviewer/Issue/Feature/
  Refactor/Slice rules. PR #49 initially exposed a Windows case collision
  between `agents.md` and `AGENTS.md`; Issue #36 comment
  [`5345020713`](https://github.com/Hans-Einar/HSX/issues/36#issuecomment-5345020713)
  blocked activation until the files were consolidated into one Windows-safe
  path.
- **OBSERVED:** Relevant GitHub work is Issues
  [#36](https://github.com/Hans-Einar/HSX/issues/36) (program acceptance),
  [#37](https://github.com/Hans-Einar/HSX/issues/37) (`DBG-RF-001`),
  [#38](https://github.com/Hans-Einar/HSX/issues/38) (`DBG-DA-001` and the
  continuing Steering gate), [#39](https://github.com/Hans-Einar/HSX/issues/39)
  through [#46](https://github.com/Hans-Einar/HSX/issues/46)
  (`DBG-RF-002..009`), [#47](https://github.com/Hans-Einar/HSX/issues/47)
  (`HSX-ST-001`) and [#48](https://github.com/Hans-Einar/HSX/issues/48)
  (`AVR-ST-001`, reserved/blocked).
- **OBSERVED:** Draft PR
  [#49](https://github.com/Hans-Einar/HSX/pull/49), created
  `2026-08-19T15:44:19Z`, targets non-default `Implementation/vscode` from
  `sdp/debugger-gap-analysis`; current head is `e5a50ab...`. Its GitHub review
  list, comments and check rollup are empty. There is no PR for the later
  `codex/dbg-rf-*` delivery branches.
- **OBSERVED:** At live branch head `82154c...`, Debugger Ledger has 114 NDJSON
  records and HSX Ledger has 24. Those records use project-local `event` fields
  and no checked-in event schema/validator was found. No `.github/workflows/`
  exists on either `main` or the latest branch.

## How SDP is used

### Work start and GitHub Issue authority

- **OBSERVED:** The old default tree starts work from phase checklists, design
  documents, grand/module implementation plans and TODO lists. GitHub identity
  is not a first-class coordinate in that tree.
- **OBSERVED:** Recent work reverses that order. Root `AGENTS.md` and
  `SDP/Shared/Process.md` require a Master to read the track CurrentIndex, then
  the active GitHub Issue and referenced Study/Requirements/Design/Refactor/
  Slice contracts before product work. Issues are explicitly the durable
  Steering Group/Codex interface; chat is not reconstructable project state.
- **OBSERVED:** Issue #36 accepted the initial process/package only after an
  exact-head audit found missing provenance, missing Slice contracts and the
  Windows AGENTS case collision. Issue #38 then became the operational parent
  for architecture, portable-contract Studies, the RF-002/RF-003 wave,
  interface refreeze, publication and RF-004 authorization. Issue #47 was the
  separate HSX Study/contract gate.
- **OBSERVED:** Later owner comments amend earlier Issue bodies. Issue #42's
  `DBG-RF-005` body names RF-002/RF-003 dependencies, while comment
  [`5362515750`](https://github.com/Hans-Einar/HSX/issues/42#issuecomment-5362515750)
  adds the accepted RF-004 typed-interface dependency and warns that the old
  body alone cannot authorize work.
- **INFERENCE:** HSX shows a strong Issue-as-authority model, but not a clean
  one-Issue/one-work-domain mapping. Issue #38 owns multiple Studies,
  architectural decisions, Refactors, Iterations, branches and Steering
  stops; child Issues #39/#40 remained open with no comments even after
  branch-local records and Issue #38 declared RF-002/RF-003 complete.
- **OWNER DIRECTION:** Issue #5 establishes GitHub Issue as the future
  operational assignment boundary and a bounded Codex session as Master for
  that Issue.
- **RECOMMENDATION:** Carry forward Issue authority and exact comment-level
  amendments, but require one explicit assignment record to identify the
  Issue's Feature/Refactor/research domain, exact base, branch/PR and stop. A
  parent Steering Issue may coordinate child work, but each completed child
  Issue must be reconciled rather than left visibly open and stale.

### Master, delegation, Architect, Reviewer and Verifier

- **OBSERVED:** Root `AGENTS.md` says every non-spawned session is Master. The
  Master owns current-track discovery, requirements/design/dependencies,
  traceability, Handoff, worker decomposition, separate review, verification
  evidence, exact-head sign-off and Steering reporting. It prohibits the
  Master from substantial product-code implementation while allowing
  documentation-only SDP work.
- **OBSERVED:** Workers receive exact requirement/design/Feature-or-Refactor/
  Slice IDs, owned files, non-goals, verification and concurrent ownership
  constraints. Issue #38 comment
  [`5356801432`](https://github.com/Hans-Einar/HSX/issues/38#issuecomment-5356801432)
  records disjoint RF-002/RF-003 workers and a read-only shared contract.
  When they found an impossible generation handshake, both stopped before a
  product commit and returned the contract to Steering.
- **OBSERVED:** A separate Architect role is not defined. Architecture is
  produced through bounded Study/design workers, Master synthesis and a fresh
  architecture Reviewer (`DBG-RVW-001-002-001..003`). A separate Verifier role
  is implied by formal `DBG-VER-*`/`HSX-VER-*` records but is not given its own
  root role contract comparable to Worker/Reviewer.
- **OBSERVED:** Review history is substantive: `DBG-DA-001` needed two REWORK
  reviews before PASS; portable HSX contracts needed five REWORK reviews
  before `HSX-RVW-001-001-006` PASS; first-wave implementation review found
  dropped terminal results, saturated-inbox loss, false healthy reopen and
  failed replacement-open continuity before final PASS.
- **LIMITATION:** All commits, Issues and comments use `Hans-Einar`, and PR #49
  has no formal GitHub review submission. Fresh/independent roles are stated in
  committed records, but external actor independence is not independently
  attributable from GitHub.
- **RECOMMENDATION:** Preserve Master prohibition on direct product work,
  bounded Worker ownership, fresh Reviewer and formal verification. Add a
  durable agent/run attestation and clarify Architect/Verifier separation for
  architecture-changing or high-risk work without requiring ceremonial roles
  on trivial Slices.

### Branches, draft PR and durable publication

- **OBSERVED:** PR #49 was opened early for the initial process baseline, but
  against `Implementation/vscode`, which was recorded as 84 commits ahead of
  and zero behind `main`. Even merging PR #49 would not make its SDP current on
  the default branch.
- **OBSERVED:** Subsequent work used a chain of `codex/dbg-*` branches. Review,
  verification and sign-off often happened locally before publication. Twice,
  Steering blocked acceptance because reported exact commits were not
  GitHub-resolvable: the portable package in Issues #47/#38 and the RF-002/
  RF-003 package in Issue #38. Both were pushed and reverified from fresh
  clones before Steering accepted them.
- **OBSERVED:** The first-wave publication chain is remote-resolvable on
  `codex/dbg-rf-002-003`: RF-002 signed product `a0640203...`, integration
  `860a98a6...`, RF-003/combined product `1e479536...`, parent sign-off
  `f52447ab...`, decision head `b0a9e816...`, final publication head
  `69a54aeb...`. Issue #38 comment
  [`5362252462`](https://github.com/Hans-Einar/HSX/issues/38#issuecomment-5362252462)
  records fresh-checkout `121 passed`, clean traceability and no later
  non-SDP product diff.
- **INFERENCE:** Exact-head remote publication is a valuable evidence gate, but
  branch publication without a PR/merge path produces several competing
  “current” trees. No GitHub CI or review surface protects the later delivery
  branches.
- **RECOMMENDATION:** Keep early draft PR and remote-resolvable/fresh-checkout
  gates, but bind every implementation assignment to one PR targeting the
  declared integration branch and ultimately reconcile it to the default
  branch. A branch-local accepted package must remain “accepted in-flight,”
  not “repository current,” until merged.

### Feature, Refactor, Sprint, Iteration, Slice and Fix

- **OBSERVED:** `SDP/Shared/Process.md` defines stable Feature (`*-FEAT-###`)
  and Refactor (`*-RF-###`) peers. Feature owns intended capability; Refactor
  owns remediation/structural change. Product work outside either requires an
  explicit Master exception.
- **OBSERVED:** No actual Feature record or `Features/` directory exists in
  HSX's evolved branch. The observed work is entirely Study/GapAnalysis/
  Refactor-led. Therefore HSX supplies evidence for Refactor semantics, not
  evidence that its written Feature contract has been exercised.
- **OBSERVED:** `DBG-SPR-001` is a long stabilization programme, not a calendar
  timebox. It contains Iteration 001 (RF-001), 002 (DesignAnalysis), 003
  (portable HSX contracts), 004 (RF-002/RF-003/integration) and active 005
  (RF-004). Capability and structure are owned by Refactors; Sprint/Iteration
  group execution stages.
- **OBSERVED:** RF-001 used one narrowly bounded, end-to-end production DAP
  Slice. RF-002/RF-003 used two domain-scoped foundation Slices plus an early
  integration Slice. Planned RF-004 uses six sequential domain Slices—legacy
  oracle, identities/address, artifact index, source resolver, recipes/stack
  and epoch inspection—with vertical convergence in the later Slice and
  parent gate. This matches the project's explicit exception allowing a
  coherent domain-scoped Refactor Slice when later vertical integration is
  named.
- **OBSERVED:** No stable Fix class is defined in `SDP/Shared/Process.md`.
  Corrective product commits (`fix(debugger): ...`) stay inside the active
  Slice/Refactor review-rework loop.
- **OWNER DIRECTION:** Issue #5 requires Feature and Refactor as first-class
  implementation owners, retains Fix only for genuine corrections and asks
  whether Sprint/Iteration should remain.
- **RECOMMENDATION:** Carry forward Feature/Refactor ownership and Refactor
  fan-out from CodeReview/GapAnalysis. Make Sprint optional execution grouping
  and Iteration an actual decision/learning cycle, not mandatory hierarchy.
  Preserve vertical Slice as the default, while allowing HSX-style local
  foundation Slices only when a frozen boundary, local verification and a
  named later integration Slice make them coherent. A Fix should be a small
  correction attached to the active Feature/Refactor/Slice unless it needs an
  independent Issue for risk or release reasons.

### Later requirements, design and first-class Studies

- **OBSERVED:** HSX directly demonstrates post-foundation evolution. The
  legacy default tree supplied broad horizontal design; recent work did not
  treat it as immutable. `DBG-ST-002..005` separately studied controller/
  concurrency, inspection/resources, frontend/package and legacy reuse. The
  Master synthesized `DBG-A-001..008` and `DBG-D-001..010`; a fresh
  architecture review challenged them before Steering froze the baseline.
- **OBSERVED:** Cross-track uncertainty then created Debugger Study
  `DBG-ST-006` and HSX coordinator `HSX-ST-001`, which spawned
  `HSX-ST-002..008`. Their synthesis created later portable requirements
  `HSX-R-001..036`, architecture `HSX-A-001..005` and design
  `HSX-D-001..005` long after the initial project design. Six independent
  review attempts and verification preceded Steering freeze in Issue #47.
- **OBSERVED:** Implementation discovery did not silently rewrite design. The
  RF-002/RF-003 generation-handshake contradiction stopped workers and returned
  to Issue #38. Steering superseded `dbg.controller-gateway/1` with reviewed
  `1.1` reserve -> authoritative success -> promote semantics before workers
  resumed.
- **INFERENCE:** This is strong evidence for an initial horizontal skeleton
  plus later Feature/Refactor-local Study/requirements/design refinement. It
  also shows a convergence gate does not need to be a single file: issue
  decision, reviewed synthesis, frozen contracts and traceability can jointly
  form the gate.
- **RECOMMENDATION:** Canonical SDP should support one broad foundation Study,
  multiple bounded Studies during later work, Master synthesis, fresh design
  review and explicit Steering acceptance. Technical Studies may inform new
  requirements/design; only an explicit relation/decision freezes them as
  implementation authority.

### Review, verification and handoff

- **OBSERVED:** Reviews are exact-head Markdown records under
  `SDP/Debugger/Sprints/001--Debugger_Stabilization/Reviews/` and
  `SDP/HSX/Reviews/`. Verification is separately recorded under per-track
  `Verification/`, followed by exact-content/parent Master sign-off. Blocking,
  High and Medium findings force rework and a fresh review.
- **OBSERVED:** Verification combines focused test matrices, broad regression
  results, YAML/NDJSON/Markdown parsing, scope/protected-path diffs, ancestry,
  object connectivity and clean-worktree checks. Known broad-suite failures
  and absent Linux evidence remain explicit rather than being silently called
  PASS.
- **OBSERVED:** Handoffs are detailed recovery surfaces:
  `SDP/Debugger/Sprints/001--Debugger_Stabilization/Handoff.md` and
  `SDP/HSX/02--Study/Handoff.md` record authority, exact heads, evidence,
  blockers, dirty-worktree isolation and the next legal action.
- **OBSERVED:** Review repeatedly found stale Handoff/CurrentIndex/review-stage
  statements. Correcting current-state drift consumed several review attempts,
  even when the technical contracts had already passed.
- **RECOMMENDATION:** Preserve exact-head review, separate verification,
  changed-head re-review, explicit limitations and concise Handoff. Generate
  coordinates/status/check identities and historical review lists; reserve
  hand-authored Handoff prose for risks, rationale and next boundary.

### CurrentIndex, Relations, Ledger and Steering state

- **OBSERVED:** Per-track CurrentIndex is both current pointer and historical
  dossier. At `69a54a...`, the Debugger CurrentIndex embeds six HSX review
  attempts, many RF-002/RF-003 review attempts, verification heads, Issue
  comment IDs, publication history, frozen interfaces, every Slice and its
  corrective sign-off. The RF-004 branch extends it again. This is much more
  than a minimal current index.
- **OBSERVED:** `Relations.yaml` supplies useful cross-track graph edges such
  as Study -> Requirement/Architecture/Design and HSX Design -> Debugger
  Design/Iteration. `Issues.yaml` binds IDs to Issue numbers and statuses.
  Ledgers preserve failed reviews, corrections, publication and Steering
  transitions rather than rewriting history.
- **OBSERVED:** There is no `Steering/CurrentAssignment` file. The effective
  assignment is distributed across Issue #38 comments, track CurrentIndex,
  Issues/Relations/Ledger, Sprint/Refactor/Slice records and Handoff.
- **OBSERVED — contradiction:** On `codex/dbg-rf-002-003@69a54a...`, HSX
  CurrentIndex and Handoff still name `89cb74a10ce36d8b0f4d0cc60332d3070c2c635f`
  as the fresh-checkout verified remote head. Issue #47's later publication
  closeout and Steering acceptance, plus the HSX Ledger, name final coordination
  head `bf92c9be6cf81a7cb704778dafe88e55fee2e235` and say a fresh checkout of that
  head passed. Even within one unmerged branch, “current” surfaces disagree on
  the accepted remote evidence head.
- **OBSERVED — contradiction:** The latest Debugger Ledger contains event
  timestamps `2026-08-21T06:00:00+02:00` and `06:01:00+02:00` inside commit
  `82154c...`, whose committer timestamp is
  `2026-08-21T00:23:26+02:00`. The event times are over five hours later than
  the commit that already contains them. Ledger timestamps therefore cannot
  be treated as verified chronology without validation.
- **INFERENCE:** Stable IDs and the append-only correction model materially
  improved recovery and review. Manual duplication of current and historical
  state caused much of the recorded review churn and can even create
  impossible time evidence.
- **RECOMMENDATION:** Keep a small CurrentAssignment/CurrentIndex containing
  Issue/work-domain, exact integration baseline, branch/PR, active Slice,
  authorization/stop, delegated roles and next gate. Keep Relations for
  durable semantic edges and Ledger for material transitions/corrections.
  Generate Issue/PR/head/check status, containment edges, review-attempt
  summaries and timestamps from Git/GitHub where possible; reject future-dated
  events.

### GitHub binding and automation

- **OBSERVED:** GitHub binding is explicit through Issue numbers/comment IDs,
  branches and full SHAs in CurrentIndex, Issues, Relations, Ledger, contracts
  and Handoff. The durable Steering loop successfully caught an invalid
  abbreviated SHA, unpublished review heads and stale dependency wording.
- **OBSERVED:** Product/test verification is automated by pytest and scripted
  parsing/diff checks, but workflow orchestration is conversational/manual.
  There are no GitHub Actions workflows, PR #49 has no checks, later branches
  have no PR review surface, and no checked schema validates CurrentIndex,
  Relations or Ledger.
- **RECOMMENDATION:** Tooling should create Issue/branch/draft-PR binding once,
  derive live GitHub state, validate exact heads/schemas/timestamps and publish
  CI evidence. Agents should author decisions, exceptions, contracts and
  evidence meaning rather than repeatedly copying mutable coordinates.

## What worked well

- **OBSERVED:** Exact-head discipline caught three classes of false confidence:
  missing provenance, local-but-unpublished evidence and an internally
  impossible generation contract.
- **OBSERVED:** Bounded multi-Study work plus Master synthesis produced portable
  HSX contracts without allowing legacy Python behavior to redefine the target.
- **OBSERVED:** Refactor fan-out and anti-monolith responsibilities converted
  one legacy debugger review into independently owned controller, gateway,
  inspection, resources, lifecycle, DAP, VS Code/package and cross-platform
  workstreams.
- **OBSERVED:** Fresh review was allowed to return repeated REWORK decisions.
  Material defects and state drift were preserved as history, then checked
  again rather than normalized away.
- **OBSERVED:** Workers stopped at scope/design boundaries. The shared-interface
  conflict went back to Steering; later RF-004 contracts preserve disjoint
  services and keep Executive/VM/frontend/AVR paths explicitly out of scope.
- **OBSERVED:** Fresh remote checkout and no-product-drift publication checks
  make dense local agent work recoverable after context/worktree loss.
- **INFERENCE:** HSX's strongest reusable contribution is the combination of
  Issue-level Steering decisions, first-class Studies, Refactor/interface
  contracts, stop/refreeze behavior, exact-head review/verification and remote
  reconstruction.

## Pain points / accidental complexity

- **OBSERVED:** Default `main` is stale and lacks the entire evolved SDP. The
  most mature process and signed product work live across several unmerged
  branches; no single canonical tree reconstructs current practice.
- **OBSERVED:** PR discipline is incomplete. PR #49 is draft against a
  non-default base with no formal review/checks, while later product branches
  have no PR. “Published” and “Steering accepted” therefore do not imply
  reviewed integration or default-branch adoption.
- **OBSERVED:** Issue #38 became a broad programme/Steering ledger. Child
  Refactor Issues #39/#40 stayed open and unupdated after their work was
  declared complete elsewhere, so GitHub's visible work state is contradictory.
- **OBSERVED:** Feature is specified but unused; Sprint 001 spans the entire
  debugger programme; Iterations are sequential stage containers. The written
  hierarchy is more flexible than old SDP, but it has not exercised Feature
  ownership and does not demonstrate Sprint as a useful timebox.
- **OBSERVED:** CurrentIndex, Issues, Relations, two Ledgers, Sprint, Refactor,
  Slice, review, verification, Handoff and Issue comments duplicate state.
  Multiple reviews spent effort correcting stale status surfaces rather than
  technical defects; HSX CurrentIndex/Handoff and Ledger/Issue evidence still
  disagree on the final verified publication head.
- **OBSERVED:** Review/verification IDs proliferate: portable contracts needed
  six review IDs; RF-002/RF-003 parent closure retains many trace-only attempts.
  Retaining history is useful, but copying every attempt into CurrentIndex
  obscures the actual current gate.
- **OBSERVED:** No durable reviewer/session identity beyond role labels and
  `Hans-Einar` exists, and no GitHub review/check provides a second provenance
  route.
- **OBSERVED:** Ledger time evidence is not trustworthy at the latest branch
  head because records are future-dated relative to the commit containing
  them. No schema or CI catches this.
- **INFERENCE:** HSX demonstrates that strong evidence discipline can coexist
  with excessive manual bookkeeping. The next canonical model should preserve
  stops, contracts and evidence while deriving mutable GitHub/graph/history
  views automatically.

## Direct contradictions and current-state rules

- **OBSERVED:** `main@e374da88...` is the canonical default tree and contains
  no SDP, while branch `82154c...` claims active `DBG-RF-004`, Iteration 005 and
  frozen Slices. Both facts are true only when branch scope is retained.
- **OBSERVED:** `main/AGENTS.md` says implementation has not begun; recent
  branches record completed product Refactors, hundreds of tests and Steering
  acceptance. The former is merged legacy truth, not current operational
  truth.
- **OBSERVED:** PR #49 is the only open PR and stops at RF-001 activation head,
  while descendant branches contain the later architecture, portable contracts
  and RF-002/RF-003 work. PR state alone does not describe all active work.
- **OBSERVED:** Issue #38/branch records say RF-002/RF-003 complete; their own
  Issues #39/#40 remain open without completion comments. GitHub Issue state
  and branch-local SDP disagree.
- **RECOMMENDATION:** Use an explicit evidence precedence model:
  1. default branch = merged/canonical repository truth;
  2. exact Issue/branch/PR head = bounded in-flight assignment truth;
  3. Steering Issue comment = durable authorization/decision that must be
     mirrored into that assignment branch before dependent work;
  4. only merge/reconciliation promotes in-flight practice to default truth.
  Inventory activity may examine all advertised refs, but canonical method
  adoption and “current repository state” must report merge status and branch
  coordinate, never silently choose the newest commit.

## Carry forward

- **RECOMMENDATION:** Issue/Issue-comment authority with exact baseline,
  Feature/Refactor/research identity, invariants, non-goals, verification,
  independent review, stop and reporting contract.
- **RECOMMENDATION:** Bounded Codex Master-per-Issue, fresh product Workers,
  independent exact-head Reviewers, formal Verifiers and Master sign-off.
- **RECOMMENDATION:** Feature/Refactor as peers above execution, with the HSX
  CodeReview -> GapAnalysis -> independent Refactor fan-out pattern for
  structural remediation.
- **RECOMMENDATION:** First-class numbered Studies at any lifecycle point,
  parallel bounded Study domains, Master synthesis, fresh design review and a
  named Steering convergence/freeze gate.
- **RECOMMENDATION:** Initial horizontal skeleton plus later work-local
  requirements/architecture/design refinement; implementation contradictions
  return to a controlled design refreeze.
- **RECOMMENDATION:** Vertical Slices by default, with an explicit exception for
  coherent domain-scoped Refactor foundations when interfaces are frozen and a
  later integration Slice is named.
- **RECOMMENDATION:** Remote-resolvable exact heads, fresh-checkout
  reconstruction, no-product-drift publication checks, truthful limitations
  and append-only correction history.
- **RECOMMENDATION:** A small machine-readable assignment/current pointer and a
  durable semantic relation graph, with live GitHub and repetitive history
  generated automatically.
- **RECOMMENDATION:** Canonical study tooling must expose both default-branch
  age and qualifying non-default activity so stale-default active repositories
  remain in scope without misrepresenting unmerged practice.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward the default tree's one-time
  Study -> Architecture -> Design -> grand implementation/module-plan model as
  the only route for later capability work.
- **RECOMMENDATION:** Do not treat the newest advertised branch as canonical
  merely because `main` is stale, or call a branch-local Steering acceptance
  merged repository truth.
- **RECOMMENDATION:** Do not allow accepted work to accumulate indefinitely on
  non-default branches without a PR/integration/default reconciliation plan.
- **RECOMMENDATION:** Do not use one broad Steering Issue as the only status
  surface while child implementation Issues remain open and stale.
- **RECOMMENDATION:** Do not preserve Sprint as mandatory capability owner or
  Iteration as a compulsory container; in HSX the Refactor owns the structural
  outcome and Sprint/Iteration merely group stages.
- **RECOMMENDATION:** Do not claim HSX proves the Feature model works in
  practice; Feature exists only as process text here.
- **RECOMMENDATION:** Do not duplicate complete review history, GitHub state,
  branch heads and publication state across CurrentIndex, Issues, Relations,
  Ledger, Handoff and prose.
- **RECOMMENDATION:** Do not accept manually future-dated ledger events,
  schema-less machine state, Master self-attestation or a single GitHub account
  as sufficient independent provenance.
- **RECOMMENDATION:** Do not canonicalize HSX's debugger-specific ID volume,
  anti-monolith domains, ABI/address/epoch contracts or exact test matrices.
  Preserve the generic boundary/evidence pattern; keep subsystem mechanics
  project-local.
