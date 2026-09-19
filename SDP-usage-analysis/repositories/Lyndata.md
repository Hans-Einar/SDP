# Lyndata SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies private repository `Hans-Einar/Lyndata` at
  the exact current default-branch identity
  `main@132a72f2c706ee6de32dae0f90491f5ac11d1aa3`, tree
  `7b28ef79b15199cb804394c25ce549b1d44d0da4`, committed
  `2026-08-05T16:21:20Z`. The local `origin/main` and GitHub's current `main`
  API both resolved to that identity on `2026-08-21`; the studied repository
  was not changed.
- **OBSERVED:** GitHub Issues, pull requests, comments, timelines, branch heads,
  checks and review APIs were inspected on `2026-08-21`. The current default
  tree is distinguished below from historical accepted evidence and from open,
  unmerged draft PR [#16](https://github.com/Hans-Einar/Lyndata/pull/16) at
  `519439053c78b4160ee3174adfc13292b6b6368a`, tree
  `74f6f33c536f755edc2980c18756dfa1bf12ca5f`, committed
  `2026-08-20T21:40:15Z`.
- **OBSERVED:** Issue #5's inventory found 165 qualifying commits reachable
  from advertised Lyndata refs in the inclusive window
  `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`. The first was initial
  commit `2afe7d4ef36ef79155d27c9fdb519b223a916acb` at
  `2026-07-30T16:01:23Z`; the latest qualifying commit was PR #16 handoff
  `9ad04a848e5a2d2219c0f709db91abcdd1ad3fae` at
  `2026-08-07T11:49:16Z`. The current PR #16 head is after the fixed window and
  is included only as current open evidence.
- **OBSERVED:** No product tests or QGIS experiments were rerun for this
  read-only process study. Test, package, runtime and independence statements
  below are repository/GitHub claims tied to the recorded commits and Actions
  runs.
- **LIMITATION:** Every Issue comment, PR action and merge is exposed through
  GitHub identity `Hans-Einar`. Repository records name Master, Worker,
  Verifier, Reviewer and Steering passes, but GitHub cannot independently
  prove which human, ChatGPT conversation or Codex context performed them.
  All eight PRs have zero native GitHub review submissions and zero inline
  review comments; “independent review” is evidenced by committed Markdown,
  exact commit identities and owner-authored reports rather than a separate
  GitHub reviewer identity.
- **LIMITATION:** PR #16 remains draft/open/unmerged. Its records are valuable
  current-practice evidence, but they are not the authoritative default tree
  and must not be treated as accepted project state.

## Repository facts

- **OBSERVED:** Repository: `Hans-Einar/Lyndata`; visibility: private; created
  `2026-07-30T15:54:29Z`; default branch: `main`; not archived and not a fork.
- **OBSERVED:** SDP is extensively installed under `SDP/`. Root authority files
  are `AGENTS.md` and `AGENTS-project.md`; role/task skills are installed under
  `.codex/skills/`. There is no `SDP/Instructions/` directory. Lifecycle
  guidance is distributed among the root agent files, `SDP/README.md`,
  `SDP/AGENT-REMINDERS.md`, `SDP/Framework/README.md`, the installed skills,
  phase documents and work records.
- **OBSERVED:** The installed generation is an explicitly unreleased SDP
  Toolkit `0.2.0` working draft, installed `2026-07-30T16:27:13Z` from exact
  source commit `e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15` and recorded in
  `SDP/Framework/installed-toolkit.manifest.yaml`. `SDP/SDP-project.manifest.yaml`
  declares profile `experimental-project-local-feature-first`; the project
  correctly does not claim the Feature-first additions are supported upstream.
- **OBSERVED:** The current default tree contains the horizontal lifecycle
  documents `SDP/01--Mandate/` through `SDP/07--Implementation/`, seven Feature
  records plus `FeatureBacklog.yaml`, one accepted Integration Study
  (`STU-001`), five Sprint folders, Steering/CurrentAssignment and seven
  interaction records, three Use Cases, `Traceability/CurrentIndex.yaml`,
  `Relations.yaml` and `Ledger.ndjson`, seven principal review records and nine
  verification/evidence records.
- **OBSERVED:** `SDP/Refactors/` and `SDP/Fixes/` are absent, and no Refactor or
  Fix identity was used. The project has no Release identity; its manifest
  remains version `0.0.0`, unreleased. Handoff records live inside each
  `SDP/Sprints/Sprint-*` folder.
- **OBSERVED:** Default-tree `Ledger.ndjson` has 162 parseable, unique events
  from `2026-07-30T16:23:51Z` through the post-merge SLC-007 extension
  authorization at `2026-08-05T15:36:21Z`. Open PR #16 appends 39 more for 201
  events.
- **OBSERVED:** Relevant GitHub work is:
  - open bootstrap Issue [#1](https://github.com/Hans-Einar/Lyndata/issues/1)
    with merged PR [#2](https://github.com/Hans-Einar/Lyndata/pull/2);
  - open accepted-Study Issue [#3](https://github.com/Hans-Einar/Lyndata/issues/3)
    with merged PR [#4](https://github.com/Hans-Einar/Lyndata/pull/4);
  - closed SLC-001 Issue [#5](https://github.com/Hans-Einar/Lyndata/issues/5)
    with merged PR [#6](https://github.com/Hans-Einar/Lyndata/pull/6);
  - closed compatibility Issue [#7](https://github.com/Hans-Einar/Lyndata/issues/7)
    with merged PR [#8](https://github.com/Hans-Einar/Lyndata/pull/8);
  - closed credential Issue [#9](https://github.com/Hans-Einar/Lyndata/issues/9)
    with merged PR [#10](https://github.com/Hans-Einar/Lyndata/pull/10);
  - reopened FEAT-007 Issue [#11](https://github.com/Hans-Einar/Lyndata/issues/11),
    accepted/merged PR [#13](https://github.com/Hans-Einar/Lyndata/pull/13),
    and current draft PR #16;
  - prepared inactive FEAT-008 Issue
    [#12](https://github.com/Hans-Einar/Lyndata/issues/12) and FEAT-009 Issue
    [#17](https://github.com/Hans-Einar/Lyndata/issues/17); and
  - closed SLC-007 tooling Issue
    [#14](https://github.com/Hans-Einar/Lyndata/issues/14) with merged PR
    [#15](https://github.com/Hans-Einar/Lyndata/pull/15).
- **OBSERVED:** Final-head Actions checks shown by GitHub passed for PRs #2,
  #4, #6, #8, #10, #13 and #15. PR #16's current head has successful
  deterministic/package and PyQGIS 4 jobs from run `32420572667`; its ordinary
  `build-exact-pr` check was skipped by path/event rules, while owner comment
  `5362256799` separately records successful forced artifact run `32420697536`.

## How SDP is used in this repository

### Work start and GitHub Issue authority

- **OBSERVED:** Work began with Issue #1, not with an active Sprint. Its body
  defined exact base, stable IDs, intended Feature catalogue, architecture,
  authorized files, credential safety, verification, independent review,
  branch/draft-PR rules and a hard stop before product work. BOOT-001 installed
  the Toolkit and created proposed Features/Study/Slices while deliberately
  keeping every delivery coordinate inactive.
- **OBSERVED:** Subsequent Issues are complete operational contracts. They
  identify exact accepted baseline, goal, scope/non-goals, invariants,
  authorization limits, expected branch/draft PR, verification matrix,
  independent-review requirements, reporting format and stop condition. The
  Master first reconciles GitHub facts that supersede stale repository state,
  selects non-duplicate SDP coordinates, posts a work-start comment, updates
  SDP, branches, and opens a draft PR early.
- **OBSERVED:** Later owner comments materially revise authority. Examples are
  the SLC-001 compatibility continuation and live-request gate in Issue #5,
  SLC-005 rework direction in Issue #9, Issue #14's owner closure despite a
  deferred live-PR smoke, and Issue #11's two owner-acceptance corrections.
  Repository records explicitly say GitHub live state outranks stale copied
  state until reconciled.
- **OBSERVED:** Issue titles do not use one consistent canonical work type:
  #1 is `BOOT-001`, #3 a Study, #5 a Slice, #7 compatibility, #9 auth, #11 a
  Feature and #14 DEVX. This flexibility helped work start, but stable SDP
  ownership must be inferred from the Issue body and relations.
- **OWNER DIRECTION:** Issue #5 makes GitHub Issue the future operational
  assignment boundary and permits one bounded Codex session to act as Master.
  Lyndata is strong direct evidence for that direction.
- **RECOMMENDATION:** Carry forward the complete Issue contract and exact
  comment amendments. Require a small machine-readable binding to one primary
  Feature, Refactor, Fix, Study or bootstrap assignment; do not infer the work
  type from the Issue title.

### Master-per-Issue, delegation and role boundaries

- **OBSERVED:** `AGENTS.md` makes an undelegated session Master. The installed
  Master skill requires contract reconciliation, bounded Worker delegation,
  separate verification and fresh independent Review, evidence inspection,
  traceability/handoff update, and a stop at the Slice/Fix boundary. Worker,
  Reviewer, Verifier and Architect skills define separate scopes.
- **OBSERVED:** Issue comments repeatedly say “You are the Codex Master” for
  that Issue. Sprint contracts say product code is changed by delegated
  Workers and reviewed by a separate fresh Reviewer. Implementation, failed
  verification, remediation, re-verification and Review integration usually
  appear as distinct commits.
- **OBSERVED:** Architect practice is clearest in `STU-001`: the Study calls
  itself an Architect candidate and separates documented fact, probe
  observation, decision, assumption and unknown. Later Feature architecture
  refinement was generally reconciled by the Master before Worker product
  implementation; separately attributable Architect passes are not evident.
- **OBSERVED:** Fresh Verifier/Reviewer passes found material defects: SLC-005
  needed five verification attempts; merged SLC-006 needed five; PR #16
  revision 2 needed five verification and three Review attempts. Review
  findings were not treated as equivalent to owner acceptance: owner testing
  rejected already verified/reviewed PR #16 revision 1 and supplied a new
  architecture/state-ownership contract.
- **INFERENCE:** The role separation improved product quality and recoverability,
  but “freshness” is a process attestation rather than externally attributable
  identity because all GitHub actions use `Hans-Einar` and no native GitHub
  review exists.
- **RECOMMENDATION:** Preserve Master-per-Issue, bounded Worker, evidence-owning
  Verifier and fresh Reviewer. Record an agent/run attestation and exact
  product/evidence heads in a compact review object. Use Architect only when a
  Study or Feature changes architecture/state ownership, as Issue #11
  correction 2 did.

### Branches, draft PRs and exact identity

- **OBSERVED:** The normal accepted sequence was exact `main` base -> named
  `codex/*` branch -> SDP activation commit -> early draft PR using `Refs #N`
  rather than `Closes` -> progressive Worker/Verifier/Reviewer commits ->
  exact-head Issue report -> Steering acceptance -> ready/merge -> additive
  post-merge reconciliation. PRs #2, #4, #6, #8, #10, #13 and #15 follow this
  pattern.
- **OBSERVED:** Records carefully distinguish product head, verification-record
  head, reviewed candidate, review-record head, handoff head and merge commit.
  Terminal CI/artifact facts are sometimes posted externally after the final
  commit to avoid an evidence commit invalidating its own exact-head claim.
- **OBSERVED:** PR #16 shows one Issue may retain an open draft PR through
  multiple owner revisions. It starts from current `main@132a72f...`, preserves
  revision 1 history, and carries revision 2 at `519439053...`; its 51-file,
  8,491-addition delta is unmerged.
- **CONTRADICTORY EVIDENCE:** SLC-007's post-merge exact-main extension was
  authorized as a local direct-`main` commit with “no remote push,” but commits
  `1ad5115...` through current remote `main@132a72f...` were pushed. That path
  bypassed the otherwise consistent Issue-bound branch/draft-PR acceptance
  surface.
- **RECOMMENDATION:** Preserve the normal early-draft-PR model. Treat direct
  default-branch work as an exceptional Fix requiring explicit GitHub-bound
  review and immediate state reconciliation; an Issue comment should not
  silently waive the normal exact-head review surface.

### Feature, Refactor, Sprint, Iteration and vertical Slice

- **OBSERVED:** BOOT-001 created `FEAT-001` through `FEAT-006` before product
  implementation. Features describe a stable problem/outcome, scope,
  non-goals, constraints, open questions, Study link, backlog/Relations links
  and current disposition. `FEAT-007` was later created by Issue #11 and used
  in delivered product work.
- **OBSERVED:** Feature ownership is not consistently 1:1 with Issue or Slice.
  SLC-001/Issue #5 delivered a vertical MVP spanning FEAT-001 through FEAT-005;
  SLC-004 also spans those five; SLC-005/Issue #9 belongs to FEAT-004;
  SLC-006/Issue #11 belongs to FEAT-007; SLC-007/Issue #14 is a tooling Slice
  under FEAT-005. BOOT-001 related six Features without delivering any.
- **OBSERVED:** FEAT-007 demonstrates reopen/extension: PR #13 was accepted and
  merged, Issue #11 closed, the same Issue reopened, and PR #16 reused
  `SPR-004 / IT-004-001 / SLC-006` with explicit revisions 1 and 2. Prepared
  Issues #12 and #17 reserve FEAT-008/009 only as owner planning and do not
  activate implementation.
- **OBSERVED:** Each delivered Issue created exactly one Sprint, one Iteration
  and one active Slice: `SPR-001/IT-001-001/SLC-001` through
  `SPR-005/IT-005-001/SLC-007`. The Sprint and Iteration usually add no grouping
  beyond the Issue/Slice. Proposed SLC-002 and SLC-003 were skipped but kept
  reserved, forcing later non-duplicate SLC-004 onward.
- **OBSERVED:** SLC-001 is genuinely vertical: UI/time/space/auth/network/parser/
  domain/layer/package behavior crosses the horizontal architecture to deliver
  one user outcome. SLC-004, SLC-005 and SLC-006 are similarly bounded
  cross-layer outcomes. SLC-007 is delivery tooling and therefore an enabler,
  not a user-facing vertical product capability.
- **OBSERVED:** No Refactor or Fix record was used. Even a behavior-changing
  acceptance correction reused SLC-006 revisions, while tooling defects were
  remediated inside SLC-007.
- **INFERENCE:** Lyndata proves Feature-first value, but also shows why Feature,
  Issue and Slice must not be forced 1:1. Feature is the durable capability;
  Issue is one bounded authorization; Slice is one implementable end-to-end
  increment. A reopened Feature may need a new Slice rather than reactivating
  completed Sprint/Iteration/Slice state indefinitely.
- **RECOMMENDATION:** Make Feature/Refactor the capability owner and Issue the
  assignment owner. Keep Sprint optional and omit Iteration when it is 1:1.
  Reserve Slice for a bounded vertical product/enabler outcome with explicit
  entry and completion evidence. Add a new Slice for material post-acceptance
  behavior; use `revision` only for bounded rework of an unaccepted candidate.

### Initial horizontal design and later Feature-local evolution

- **OBSERVED:** BOOT-001 established one horizontal `REQSET-001`, `ARC-001`,
  `DAN-001`, `DES-001` and `IMPL-001` before product code. The architecture
  separated UI, controller, query service, Frost adapter, HTTP/auth boundary,
  normalized records and QGIS layer adapter. SLC-001 then implemented a
  complete vertical path through those layers.
- **OBSERVED:** Later Issues changed Requirements, Architecture, Design Analysis
  and Design after initial delivery. SLC-005 introduced QGIS Authentication
  Manager ownership; SLC-006 added scalable models and selection; SLC-007 added
  a tooling control plane. Issue #11 correction 2 explicitly superseded the
  earlier “QGIS layer selection is the single source of truth” decision with
  persistent selected membership plus independent per-table row buffers.
- **OBSERVED:** PR #16 correctly reconciled those architecture/design records
  before Worker implementation. However, the changes are written back into the
  same global phase files; Feature-local decision deltas exist mainly in the
  Issue, SLC record and Git history, not in a dedicated Feature-local
  requirements/design folder.
- **INFERENCE:** Lyndata strongly supports the owner-directed horizontal
  skeleton plus later vertical Feature growth. It also demonstrates the audit
  weakness of mutating one global design: current prose describes the latest
  state well, but the reason a prior accepted rule was superseded is scattered
  across GitHub comments, interaction history and Git.
- **RECOMMENDATION:** Keep the initial horizontal skeleton. Put later Study,
  requirement and design deltas under the Feature/Refactor, and explicitly link
  any system-wide architecture revision/supersession into the canonical
  architecture graph.

### Studies and convergence

- **OBSERVED:** Lyndata used one broad `STU-001` Integration Study shared by
  FEAT-001 through FEAT-005 and later FEAT-007. Issue #3 authorized research
  and a narrowly bounded live probe but no product implementation. The Study
  preserved an initial blocked probe, owner direction, one sanitized successful
  observation, exact-head Verification, independent Review, Steering
  acceptance and merge before Issue #5 activated product delivery.
- **OBSERVED:** `STU-001` distinguished documented facts, probe observations,
  decisions, assumptions and unresolved unknowns. This prevented undocumented
  Frost/UALF semantics and project safety caps from being represented as
  provider guarantees.
- **OBSERVED:** The project did not use multiple independent Studies or a
  formal convergence manifest. `ISTUDY-FEAT-006` remains a placeholder despite
  the governance pilot being extensively exercised. Later FEAT-007 semantics
  and revisions used Issue-bounded source study and design reconciliation
  without a new Study ID.
- **RECOMMENDATION:** Carry forward the evidence taxonomy, separate Study
  authorization, blocked/resumed evidence and Steering acceptance gate. Support
  one broad or one Feature-local Study by default; use multiple independent
  Studies/convergence only when uncertainties are genuinely separable.

### Verification, review, acceptance and handoff

- **OBSERVED:** Verification and Review are first-class committed records under
  `SDP/Verification/` and `SDP/CodeReview/`. They state exact candidates,
  environments, commands, package hashes, Actions runs, findings and
  disposition. Failed attempts and findings remain in later records rather
  than being erased.
- **OBSERVED:** The loop caught substantive defects. SLC-001 required product
  remediation for a real empty-200 observation and then a separate history
  repair after one Issue comment was overwritten. SLC-005 exposed credential
  lifecycle and process findings. SLC-006 exposed repeated 100,000-row
  performance/selection failures. PR #16 exposed map gesture, state ownership,
  halo and virtual selection defects.
- **OBSERVED:** Each Sprint has `Handoff.md`, usually with exact base, product,
  verification/review heads, package identity, Issue report and stop boundary.
  These records make failed/reworked candidates reconstructable after a new
  context starts.
- **CONTRADICTORY EVIDENCE:** Handoffs and assurance records can themselves be
  stale. Default-tree Sprint 002 `Handoff.md` still says “awaiting-steering”
  after merge; the SLC-001 interaction still says in-progress; several Feature
  records contain “Study is in progress and not accepted” despite `STU-001`
  being accepted. Current PR #16's `SDP/README.md` says fresh Worker
  implementation is next even though its CurrentAssignment and GitHub report
  say revision 2 is implemented, verified and approved.
- **RECOMMENDATION:** Preserve exact-head verification, changed-head re-review,
  owner acceptance after review and concise handoff. Generate mutable status,
  Issue/PR/check state and accepted head from authoritative sources; keep prose
  for rationale, risks, limitations and next decision.

### CurrentIndex, Relations, Ledger and CurrentAssignment

- **OBSERVED:** `CurrentIndex.yaml` combines active pointer, GitHub snapshot,
  evidence registry and detailed delivery record. `Relations.yaml` links
  Mandate, Use Cases, Features, Study, requirements/design, Sprints,
  Iterations, Slices, assignments, findings, verification and review. The
  append-only Ledger preserves transition order and explicit corrections.
- **OBSERVED:** This traceability was useful when SLC-001 restored an overwritten
  historical Issue comment without rewriting 81 earlier events. Exact IDs also
  make five rapid delivery assignments and many failed candidates
  reconstructable.
- **OBSERVED:** `Steering/CurrentAssignment.yaml` is a strong recovery contract:
  it binds assignment, Issue, branch, PR, exact base, Toolkit, predecessor,
  coordinates, authorization flags, compatibility, scope/non-goals,
  verification/review and stop boundary.
- **CONTRADICTORY EVIDENCE:** Manual duplication failed at the current default
  tip. GitHub says Issue #14 is closed completed and PR #15 merged, while
  `SDP/README.md`, project manifest, CurrentAssignment, CurrentIndex, Relations,
  Sprint 005 and implementation plan still call SLC-007 active. The default
  Ledger stops at extension authorization and does not record the extension's
  implementation, verification, review, owner closure or later findings.
- **CONTRADICTORY EVIDENCE:** More seriously, current default
  `SDP/CodeReview/REV-SLC-007.md` says the post-merge exact-main extension is
  `changes required` and records unresolved Medium findings
  `FND-SLC-007-004` through `-006` plus Low `-007`. `CurrentIndex.yaml`
  nevertheless labels `REV-SLC-007` approved, and Relations does not register
  findings 004-007. Owner closed Issue #14 at `2026-08-05T15:52:32Z` when
  `main` was `1ad5115...`; verification `d056890...`, the changes-required
  Review `07c380f...`, and current tip `132a72f...` followed. No fresh
  verification/re-review resolving those findings exists on default `main`.
- **OBSERVED:** PR #16 repairs the primary assignment snapshot and appends 39
  events, but it is unmerged and declares SLC-007 closed while carrying the
  unresolved default-tree review record unchanged. It therefore cannot be
  treated as resolution of the SLC-007 review contradiction.
- **INFERENCE:** Stable relations and append-only corrections added real value;
  162 events in six days, a large CurrentIndex, a large Relations graph and a
  large CurrentAssignment also created more mutable truth copies than agents
  could keep coherent.
- **RECOMMENDATION:** Retain stable IDs/relations for Feature/Refactor, material
  Study/requirement/decision, Slice, verification/review/finding and Release.
  Reduce CurrentAssignment to authority, baseline, work item, active Slice,
  permissions, roles, required evidence and stop. Generate Issue/PR/branch/
  check/current-head state and routine lifecycle events from GitHub/Git; make
  stale or contradictory review state a validator failure.

### Steering, GitHub binding and automation

- **OBSERVED:** `SDP/Steering/README.md` defines Issue -> CurrentAssignment/
  interaction -> Master/Worker -> exact-head Verification/fresh Review -> PR
  evidence/Master report -> Steering disposition. Interaction records preserve
  stable references and chronology without copying the full Issue body.
- **OBSERVED:** Binding is explicit but repeated across Issue, PR, CurrentAssignment,
  CurrentIndex, Relations, Feature, Slice, Sprint, interaction, Handoff and
  Ledger. GitHub numbers are correctly treated as dynamic while `FEAT-*`,
  `STU-*`, `SPR-*`, `IT-*` and `SLC-*` remain stable.
- **OBSERVED:** Automation is strong for deterministic behavior: validators,
  YAML/NDJSON/path checks, Python tests/syntax, package double-build/inspection,
  PyQGIS 4 checking, CI, exact-PR artifact generation and local synchronization.
  It is weak for governance synchronization: Issue/PR state, assignment state,
  findings and handoff status are manually copied and the validator did not
  reject the current GitHub/default-tree contradictions.
- **INFERENCE:** The practical Steering role is real—Issues contain owner
  activation, rework, acceptance and closure decisions—but repository evidence
  cannot distinguish ChatGPT Steering from the human owner behind the shared
  GitHub identity. Chat remains transport; Issue/repository records are the
  durable evidence.
- **RECOMMENDATION:** Canonical SDP should define Steering responsibilities and
  require durable Issue decisions, without naming ChatGPT as authority. Tooling
  should create one GitHub/SDP binding and continuously compare, rather than
  manually replicate, its mutable fields.

## What worked well

- **OBSERVED — Feature-first bootstrap:** Features existed before delivery and
  later Issues could add FEAT-007 and prepared FEAT-008/009 without pretending
  the initial design contained every future capability.
- **OBSERVED — Horizontal then vertical:** The initial architecture provided
  stable UI/application/source/QGIS boundaries; SLC-001 delivered a complete
  vertical user outcome and later Slices evolved it without collapsing those
  boundaries.
- **OBSERVED — Issue-bounded autonomy:** Exact baseline, authorization flags,
  non-goals, branch/PR expectations, verification and stop conditions let a
  new Master act autonomously while remaining bounded.
- **OBSERVED — Truthful evidence:** `STU-001`'s fact/observation/decision/
  assumption/unknown taxonomy, sanitized probe record and explicit provider
  limitations prevented unsupported claims.
- **OBSERVED — Review quality:** Multiple fresh passes found real functional,
  security, scale and traceability defects; failed attempts were retained and
  owner testing remained a separate acceptance gate.
- **OBSERVED — Recoverability:** Exact commit/package/run identities,
  interaction chronology, handoffs and append-only events make a very dense
  week of work reconstructable after context loss.
- **OBSERVED — Safe experimentation:** The repository labels Feature-first
  fields as project-local and Analyzer-limited rather than falsely claiming a
  stable Toolkit schema.

## Pain points / accidental complexity

- **OBSERVED — Duplicated mutable truth:** GitHub, README, two manifests,
  CurrentAssignment, CurrentIndex, Relations, Feature backlog, Sprint, Slice,
  interaction, review/verification, Handoff and Ledger repeat current state.
  Default `main` proves that these copies can contradict both GitHub and each
  other.
- **OBSERVED — Redundant hierarchy:** Five Sprints each contain one Iteration
  and one Slice. Sprint and Iteration provide naming/ceremony more often than
  scheduling or grouping value.
- **OBSERVED — ID gaps and reuse pressure:** Proposed inactive SLC-002/003 were
  retained, so later work jumped to SLC-004. Material owner revisions then
  reused completed SLC-006 and its completed Sprint/Iteration with revision
  numbers, making “active,” “completed” and “accepted” coexist.
- **OBSERVED — Feature completion ambiguity:** FEAT-001 through FEAT-005 remain
  `in-progress` after several accepted Slices; FEAT-006 remains proposed even
  though its governance pilot is the repository's dominant process. There is
  no canonical completion/reopen/release rule for a Feature.
- **OBSERVED — Global design mutation:** Later Feature decisions update the
  global phase files. The current architecture is clear, but Feature-local
  rationale and supersession require Issue/Git-history reconstruction.
- **OBSERVED — Heavy event bookkeeping:** 34 Ledger events existed at bootstrap,
  85 by SLC-001 handoff, 162 on current default main and 201 on PR #16. Many
  are manually copied activation, report and status transitions which GitHub
  or tooling could derive.
- **OBSERVED — Review provenance gap:** Independent review is substantive but
  not represented by a distinct GitHub reviewer. Current SLC-007 review state
  also demonstrates that committed findings can be bypassed by owner closure
  and stale indexes.
- **OBSERVED — Direct-main exception:** The SLC-007 extension diverged from the
  successful branch/draft-PR pattern, was remotely pushed despite a no-push
  boundary, and ended with unresolved current review findings.
- **INFERENCE:** The process's strongest quality controls are Issue scope,
  exact-head evidence and independent challenge. Its weakest part is the
  volume of hand-synchronized status metadata around those controls.

## Direct contradictions and limitations

- **OBSERVED:** GitHub Issue #14 is closed completed; default SDP says SLC-007
  active. Default Review says post-merge changes required; CurrentIndex says
  approved. Relations/Ledger omit the later findings and closure.
- **OBSERVED:** GitHub Issues #1 and #3 remain open despite their PRs being
  merged and repository records calling BOOT-001 owner-accepted and STU-001
  accepted. Open state may intentionally preserve discussion, but it prevents
  Issue state alone from meaning assignment activity.
- **OBSERVED:** Several historical “current” documents retain pre-Steering
  status after accepted merge. PR #16 repairs much of this only on an open
  branch and still leaves some stale top-level prose.
- **OBSERVED:** Issue #11 correction 2 supersedes an accepted architecture rule.
  The new implementation/review evidence is unmerged, so current default
  product practice remains the PR #13 model even though current owner direction
  requires the PR #16 model.
- **LIMITATION:** Private GitHub links require authorized access. Local ignored
  QGIS profiles/builds and owner-observed runtime facts cannot be independently
  inspected from GitHub after the fact.
- **LIMITATION:** Successful Actions checks prove the automated jobs at their
  exact heads, not the correctness of owner-only QGIS observations or role
  independence. No native PR review corroborates the committed Review records.

## Carry forward

- **RECOMMENDATION:** GitHub Issue as the durable operational assignment, with
  exact baseline, primary Feature/Refactor/Fix/Study identity, scope/non-goals,
  invariants, authority flags, verification, independent review, branch/draft
  PR, stop condition and reporting contract.
- **RECOMMENDATION:** One bounded Codex Master per Issue, with repository-first
  reconciliation, bounded Workers, architecture pass when needed, exact-head
  Verifier, fresh Reviewer and a stop for Steering/owner acceptance.
- **RECOMMENDATION:** Stable Feature records as capability owners, including
  later Feature-local Study/requirements/design refinement and explicit
  completion/reopen/release inclusion semantics.
- **RECOMMENDATION:** Initial horizontal architecture/design skeleton followed
  by small vertical Slices. Lyndata's SLC-001 is a strong example of a Slice
  that crosses the required layers to produce one usable outcome.
- **RECOMMENDATION:** Early issue-bound branch and draft PR, `Refs` rather than
  automatic close, exact base/product/review/merge identities, and external
  terminal evidence where necessary to avoid self-referential commits.
- **RECOMMENDATION:** `STU-001`'s evidence taxonomy, separate Study authority,
  explicit unknowns, sanitized operational evidence and Steering convergence
  before implementation.
- **RECOMMENDATION:** Exact-head deterministic CI plus real-environment evidence
  where justified, failed-attempt preservation, changed-head re-review, and
  owner testing as a distinct acceptance layer.
- **RECOMMENDATION:** A concise generated handoff/current-assignment view and a
  smaller stable traceability graph, with append-only corrections for material
  decisions and findings.
- **RECOMMENDATION:** Validator/Analyzer checks that compare live GitHub state,
  current review findings and repository declarations, surfacing rather than
  normalizing contradictions.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward Sprint as capability owner or a
  mandatory Iteration that is always 1:1 with Sprint and Slice. Keep Sprint an
  optional execution/timebox/group only.
- **RECOMMENDATION:** Do not reserve proposed Slice numbers indefinitely or
  reactivate a completed Slice/Sprint/Iteration for material new owner behavior.
  Create a new Slice; reserve revision for rework before acceptance.
- **RECOMMENDATION:** Do not require one GitHub Issue to map 1:1 to one Feature
  or vice versa. Require one primary assignment identity and explicit links for
  cross-cutting Slices.
- **RECOMMENDATION:** Do not manually copy Issue/PR/branch/check/current-head
  state across ten documents or log every routine transition as a hand-authored
  Ledger event. Generate those facts and retain only decisions/exceptions.
- **RECOMMENDATION:** Do not retain one mutable global requirements/design file
  as the sole home for later Feature decisions. Preserve Feature-local deltas
  and explicit architecture supersession.
- **RECOMMENDATION:** Do not treat committed “approved” labels, green CI,
  Master self-check or owner closure as a substitute for a current independent
  Review with no unresolved material finding. The SLC-007 contradiction is the
  concrete counterexample.
- **RECOMMENDATION:** Do not carry forward direct-default-branch product/tooling
  changes that bypass the issue-bound draft-PR review model.
- **RECOMMENDATION:** Do not canonicalize Lyndata's QGIS/Frost-specific
  credential rules, provider limits, Windows runtime matrix, UALF semantics or
  100,000-row UI thresholds. They are justified project-local contracts; the
  reusable SDP part is their authority, evidence, verification and stop-gate
  structure.
