# TerrainAnalyzer SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/TerrainAnalyzer` at the exact current default-branch commit `main@77d71dc9bd248986409103e0e1fa1c37a3dfa21f`, committed and merged on `2026-08-19T17:53:15Z` (`2026-08-19T19:53:15+02:00`). Its tree is `afeb713995f481cddd1cafcbcff62463a8499add`, identical to accepted PR #21 head `8163c74d6643d8f65fb1f4d9523b3935573bb474`. Files were read from that Git object; the studied checkout was not changed.
- **OBSERVED:** GitHub Issue and PR state was inspected on 2026-08-21. Default-tree evidence is distinguished from later Issue comments and open, non-default PRs #22 and #23.
- **OBSERVED:** The GitHub connector returns Issue-comment identities and bodies but omitted their `created_at` values. Where exact comment time is unavailable, this report gives the stable comment URL/ID and uses the enclosing Issue/PR update date or commit date without inventing a comment timestamp.
- **OBSERVED:** No test suite was run for this read-only study. Verification claims below are claims recorded in committed verification files, GitHub Issues/PRs, and Actions identities. GitHub returned no combined status or workflow run for merge commit `77d71dc...`; the workflow is PR-event-only, and the accepted PR-head runs are recorded separately.
- **INFERENCE:** “Fresh” or “independent” agent identity is strongly documented but cannot always be independently reconstructed. Most reviewer identities are role strings in the Ledger; PR #21 has no formal GitHub review submission or review thread. PR #19 does have one formal `COMMENTED` review at `7355b2d...` and two resolved inline threads, which corroborates part of the repository review history.
- **LIMITATION:** SDP generation/version is not determinable. There is no `SDP.manifest.yaml` or `SDP-project.manifest.yaml`. `SDP/README.md` says the initial shape was adapted from `weight_app_flutter/SDP` and `tsm_locations_core/SDP`, and the local instructions call themselves a working draft.

## Repository facts

- **OBSERVED:** Repository: `Hans-Einar/TerrainAnalyzer` (private, visible to the authenticated owner); default branch: `main`.
- **OBSERVED:** Study commit: `77d71dc9bd248986409103e0e1fa1c37a3dfa21f`; commit tree: `afeb713995f481cddd1cafcbcff62463a8499add`; parents: pre-preparation `main@51a323887d97acefc80f1a255ad17c526ff4bf50` and PR #21 head `8163c74d6643d8f65fb1f4d9523b3935573bb474`.
- **OBSERVED:** Within Issue #5's fixed window (`2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`), 255 commits are reachable from the studied default head. The first is `564fdca00e4f01e2984dc7023a8d4dd1d8ee4fa3` (“Study done”) at `2026-06-21T19:12:05Z`; the latest is the studied merge at `2026-08-19T17:53:15Z`.
- **OBSERVED:** SDP is extensively installed at `SDP/`. Root authority is `AGENTS.md`; there is no `AGENTS-project.md`. Agent lifecycle rules are in `SDP/AGENT-REMINDERS.md` and `SDP/Instructions/AgentLifecycleAndDelegation.md`, not a separate lifecycle folder.
- **OBSERVED:** Major surfaces are `SDP/01--Mandate/` through `SDP/07--Implementation/`, `SDP/Features/`, `SDP/Steering/`, `SDP/Instructions/`, `SDP/Sprints/`, `SDP/Traceability/`, `SDP/CodeReview/`, and `SDP/Verification/`. There is no `SDP/Refactors/`; refactor history is represented by `SDP/07--Implementation/ImplementationRefactorPlan.md`, Sprint folders such as `SDP/Sprints/Sprint-022-ProjectModeSourceCreationRefactor/`, and open GitHub Issue #10.
- **OBSERVED:** `SDP/Features/` contains `FEAT-001` through `FEAT-007` and `FeatureBacklog.yaml`. `SDP/Sprints/` contains Sprints 001-032 (with no 033/034 records on the studied tree). `CurrentIndex.yaml` registers 107 requirements, 32 Sprints, 33 Iterations, 69 Slices, 124 implementation IDs, 81 review IDs, and 69 verification IDs.
- **OBSERVED:** `SDP/Traceability/Ledger.ndjson` contains 814 parseable events from `2026-06-20T01:14:04+02:00` through `2026-08-12T23:55:46+02:00`. It stops before PR #21's 2026-08-19 merge.
- **OBSERVED:** Important GitHub work includes Study Issue [#1](https://github.com/Hans-Einar/TerrainAnalyzer/issues/1), implementation Issues [#2](https://github.com/Hans-Einar/TerrainAnalyzer/issues/2), [#5](https://github.com/Hans-Einar/TerrainAnalyzer/issues/5), [#7](https://github.com/Hans-Einar/TerrainAnalyzer/issues/7), [#8](https://github.com/Hans-Einar/TerrainAnalyzer/issues/8), [#12](https://github.com/Hans-Einar/TerrainAnalyzer/issues/12), and [#13](https://github.com/Hans-Einar/TerrainAnalyzer/issues/13), planned Refactor [#10](https://github.com/Hans-Einar/TerrainAnalyzer/issues/10), planned Study [#16](https://github.com/Hans-Einar/TerrainAnalyzer/issues/16), and reconciliation Issue [#17](https://github.com/Hans-Einar/TerrainAnalyzer/issues/17).
- **OBSERVED:** Merged delivery PRs include [#3](https://github.com/Hans-Einar/TerrainAnalyzer/pull/3), [#6](https://github.com/Hans-Einar/TerrainAnalyzer/pull/6), [#9](https://github.com/Hans-Einar/TerrainAnalyzer/pull/9), [#11](https://github.com/Hans-Einar/TerrainAnalyzer/pull/11), [#18](https://github.com/Hans-Einar/TerrainAnalyzer/pull/18), [#19](https://github.com/Hans-Einar/TerrainAnalyzer/pull/19), reconciliation PR [#20](https://github.com/Hans-Einar/TerrainAnalyzer/pull/20), and preparation PR [#21](https://github.com/Hans-Einar/TerrainAnalyzer/pull/21). PR #21 was created `2026-08-12T18:33:21Z` and merged `2026-08-19T17:53:15Z`.
- **OBSERVED:** Open draft PR [#22](https://github.com/Hans-Einar/TerrainAnalyzer/pull/22) (Issue #15 Diff) and [#23](https://github.com/Hans-Einar/TerrainAnalyzer/pull/23) (Issue #14 Mesh) start from exact `77d71dc...` but are not default-branch evidence. They were created after the studied merge and remained open/draft on 2026-08-21.

## How SDP is used

### Work start and Issue authority

- **OBSERVED:** Older June-July work followed the phase/Tier plan in `02--Study`, requirements studies, horizontal designs, and one Tier-oriented Sprint after another. The newer Issue-driven practice begins clearly with Issue #1/#2 on 2026-08-05.
- **OBSERVED:** The newer start sequence is: owner/Steering Issue or authoritative Issue comment; exact `main` baseline; stale-state reconciliation; new Feature/Sprint/Iteration/Slice and evidence IDs; branch; SDP-only activation commit; early draft PR; then product work. Issue #2 comment [`5195965237`](https://github.com/Hans-Einar/TerrainAnalyzer/issues/2#issuecomment-5195965237) records the exact base and selected coordinates before code, and comment [`5196098449`](https://github.com/Hans-Einar/TerrainAnalyzer/issues/2#issuecomment-5196098449) records the early draft PR while product work is still unstarted.
- **OBSERVED:** GitHub Issues are the live authority, not just links. Issue bodies specify baseline, files to read, invariants, non-goals, allowed paths, verification, review, stop/merge authority, and reporting. Later owner comments can amend scope or scheduling: Issue #13 gained Slices 7-8 through comments, and Issue #12 comment [`5257061940`](https://github.com/Hans-Einar/TerrainAnalyzer/issues/12#issuecomment-5257061940) replaced sequential #14→#15 scheduling with a same-base parallel preparation gate.
- **OBSERVED:** Issue-to-Feature is not strictly 1:1. Issue #12 owns `FEAT-005` and later comment-scoped `FEAT-007`; Issue #1 remains an open parent Study while separate Issues authorize child Features. Closed Issue #13 later received an owner clarification and bounded authority extension in comment [`5359420561`](https://github.com/Hans-Einar/TerrainAnalyzer/issues/13#issuecomment-5359420561), after the studied default commit.
- **OWNER DIRECTION:** The repository's Issue practice matches Issue #5's desired durable operational boundary: repository state is authoritative, while Issue comments are authorized Steering inputs that must be reconciled into repository records.
- **RECOMMENDATION:** Canonical SDP should allow `Feature.authority` to be an Issue or exact Issue comment and allow one Issue to govern multiple Features only when each bounded assignment has its own durable Feature identity. Require a post-merge reconciliation event or generated snapshot so Issue and repository state do not diverge.

### Codex Master, delegation, Architect, Verifier, and Reviewer

- **OBSERVED:** `AGENTS.md` makes every non-spawned top-level session the Master. The Master reads `CurrentIndex`, updates the Slice contract and traceability first, delegates product code to fresh Workers, delegates review to a different fresh Reviewer, decides acceptance from verification evidence, and updates Sprint/traceability/handoff. It explicitly prohibits the Master from direct product-code implementation.
- **OBSERVED:** Recent Sprint 030-032 contracts repeatedly say each product Slice gets a fresh Worker and separate fresh Reviewer; the Master owns coordination, SDP integration, exact evidence, PR state, remediation decisions, and Issue reporting. Ledger events record worker/reviewer/verifier role labels and exact candidates.
- **OBSERVED:** A distinct Verifier exists in newer Issues and evidence (`VER-*`), but the canonical instruction is mixed: `AGENTS.md` says the Master uses verification evidence and `SoftwareDevelopmentProcess.md` says “run verification”; Issue #2 explicitly asked for a Verifier. Verification records sometimes say “formal verifier-and-master.”
- **OBSERVED:** No repository instruction assigns a distinct Architect role. Architecture/design changes are owned through Study/Feature/Slice contracts and Master integration. Separate architecture documents (`A-*`, `DD-*`) are strong, but separate fresh Architect-agent evidence was not found.
- **CONTRADICTORY EVIDENCE:** Older Sprint 022 says the Master integrated frontend work and fixed verifier findings directly, and calls fresh agents “verifiers” while using them for review. This conflicts with the stricter current root `AGENTS.md`. It is historical evolution, not evidence that direct Master implementation is still accepted.
- **INFERENCE:** The role model became stricter over time: early Master implementation/review blending gave way to Worker/Reviewer separation and exact-commit re-review, while Architect and Verifier boundaries remain less formal than Worker/Reviewer.
- **RECOMMENDATION:** Preserve Master-per-Issue, Worker, independent Reviewer, and evidence-owning Verifier. Define when Verifier may be the Master and when safety/complexity requires a separate context. Add Architect only for architecture-changing Slices; do not require the role ceremonially for every Feature.

### Branches and draft PRs

- **OBSERVED:** Newer work records exact base, branch, PR, candidate head/tree, CI run, review result, merge authorization, and stop boundary. Draft PRs are opened early and kept unmerged until Steering accepts evidence. Accepted PRs are later marked ready/merged under explicit owner authority.
- **OBSERVED:** PR #3 started with an SDP activation commit. PR #18 started foundation work from reconciled `main`; PR #19 started its production consumer only after PR #18 merged and a separate owner assignment. PR #20 was a short-lived governance-only reconciliation. PR #21 established a same-base parallel child topology.
- **OBSERVED:** Older PR #6 was stacked on PR #3 (`base=codex/issue-2-nn2000-meshbuilder`), whereas the later preparation contract explicitly reserves sibling branches from one exact merge commit. This is a concrete evolution away from fragile stacked delivery when independent work can be made disjoint.
- **OBSERVED:** Exact reviewed PR head is distinguished from merge commit. Sprint 029/031 handoffs record first/second parents and explicitly refuse to claim the merge commit itself was reviewed.
- **RECOMMENDATION:** Keep early draft PR, exact base/head, and reviewed-head/merge-identity separation. Prefer sibling branches after a designed convergence/preparation gate; allow stacked PRs only when the dependency is real and retest/review after rebasing onto `main`.

### Feature, Refactor, Sprint, Iteration, and vertical Slice

- **OBSERVED:** `SDP/Features/README.md` defines Features as owner-visible capabilities connecting GitHub Steering authority to requirements, design, Sprint, verification, and review. Stable IDs are `FEAT-###`; the Issue holds live assignment/discussion and `feature.md` holds durable scope, authority, non-goals, and delivery links.
- **OBSERVED:** `FeatureBacklog.yaml` maps each Feature to authority, parent Study, Feature path, one Sprint, and a terminal Slice. Feature records link requirements, architecture, design, Sprint/Iteration/Slices, verification, and review.
- **OBSERVED:** In current practice the Feature owns the capability and the Sprint is its execution vehicle. `FEAT-004` uses `SPR-029`, two Iterations (one superseded), and seven final Slices. `FEAT-005`, `FEAT-006`, and `FEAT-007` each use one Sprint/Iteration and 4-8 sequential Slices.
- **OBSERVED:** Slices are vertical and gated, not merely file batches. Recent Slice contracts state goal, why now, expected paths, invariants, non-goals, IDs, verification, review, and completion. They progress from activation/contracts through bounded product capability to integrated exact-head closure.
- **OBSERVED:** `FEAT-007` is called a Feature even though its outcome is behavior-preserving parallelization preparation rather than an owner-visible product capability. Its `work_type` is `parallelization_preparation` in `CurrentAssignment.yaml`.
- **OBSERVED:** Refactor has no first-class folder or current traceability type. Historical refactors are Sprint-owned via `ImplementationRefactorPlan.md` and Sprint 022. Planned Issue #10 is clearly called a Refactor but has no repository-local Refactor ID/record yet. `ImplementationRefactorPlan.md` says no active follow-up after Sprint 024, which is stale relative to planned Issue #10.
- **RECOMMENDATION:** Make Feature and Refactor peers above execution. Treat `FEAT-007`-like preparation as a Refactor/Enabler (or explicit preparation work type), not force it into owner-visible Feature semantics. Retain Sprint only as an optional execution/timebox/group inside Feature/Refactor; retain Slice as the mandatory bounded vertical gate.

### Later requirements and design refinement

- **OBSERVED:** TerrainAnalyzer provides strong evidence against a frozen one-time design. The initial requirements stop at `R-063`; Issue-driven work adds `R-064` through `R-107`. Later Features add/supersede architecture `A-002` through `A-006` and detailed designs `DD-013` through `DD-018`.
- **OBSERVED:** `FEAT-004` superseded ephemeral MeshBuilder/TerrainDiff designs with durable Geometry→SurfaceMesh→DifferenceSurface (`A-005 / DD-016`). `FEAT-005` added feature-local Study, architecture, design, and requirements (`A-006 / DD-017 / R-085..090`). `FEAT-006` refined both `DD-016` and `DD-017`; `FEAT-007` added `DD-018` and `R-103..107` without reopening accepted product behavior.
- **OBSERVED:** `Relations.yaml` explicitly uses `refines_architecture`, `refines_design`, `preserves`, `supersedes`, and `delivered_by`, allowing later work to revise the horizontal skeleton truthfully.
- **RECOMMENDATION:** Carry forward explicit feature-local Study/requirements/design refinement. Require an architecture revision when system ownership/boundaries change; use design refinement for bounded local contracts; do not rewrite older accepted evidence as if it knew the later decision.

### Studies and convergence

- **OBSERVED:** The original `SDP/02--Study/02--Study.md` is one broad repository study. Several bounded pre-design studies live unusually under `SDP/03--Requirements/`, then `PreDesignCloseoutStudy.md` acts as a sequential convergence/phase gate.
- **OBSERVED:** Issue #1 is a durable, still-open Study that refines findings through comments and authorizes separate implementation Issues one at a time. It explicitly does not itself authorize production code.
- **OBSERVED:** `SDP/02--Study/ModeUIStudy.md` is a later Feature-local Study with evidence, alternatives, decision, ownership boundary, risks, extraction criteria, and a gate into `FEAT-005`. It explicitly says implementation authority comes from Issue #12/Sprint 030, not the Study.
- **OBSERVED:** There is no standardized independent `STUDY-*` identity, parallel-study manifest, claim/evidence level, or formal multi-study convergence record. TerrainAnalyzer demonstrates bounded studies and sequential gates, not ActionCam-style independent-study convergence.
- **RECOMMENDATION:** Carry forward broad initial Study plus Feature-local Study. Add optional independent Studies and convergence gates only when uncertainty is separable; do not impose that machinery on the coherent ModeUI-style single Study.

### Review and verification

- **OBSERVED:** Review quality is high in substance. `SDP/CodeReview/013--Sprint_031_Geometry_ModeUI/01-review.md` and `014--Sprint_032_Parallel_Mesh_Diff_Preparation/01-review.md` record successive B/H/M/L/Note findings, remediation, exact commit/tree, reproducer checks, CI runs, and mandatory changed-head re-review. Green CI does not waive manual findings.
- **OBSERVED:** Verification is separate and concrete. `VER-032-001-002.md` and `VER-032-001-003.md` give exact accepted commits, tests/builds/browser repetitions, deterministic hashes, scope scans, review history, and final CI. `VER-032-001-004.md` correctly says its record-head confirmation remains external after the closeout commit.
- **OBSERVED:** Issue/PR discussion sometimes holds the last immutable acceptance evidence so the reviewed commit is not changed by an evidence-only commit. This avoids self-invalidating exact-head records but splits final truth between repository and GitHub.
- **CONTRADICTORY EVIDENCE:** PR #21 has no formal GitHub reviews, threads, or PR comments despite committed claims of fresh exact-head B0/H0/M0 review. PR #19 has only one formal GitHub `COMMENTED` review at investigation head `7355b2d...`; final reviews mostly live in SDP/Issue records. Therefore “independent” means process separation by convention, not necessarily a distinct GitHub reviewer identity.
- **RECOMMENDATION:** Preserve exact-head CI plus fresh changed-head review and the rule that B/H/M blocks acceptance. Give each review a stable evidence object containing reviewer/session identity or attestation, base/head, commands/probes, findings, and disposition. Automate copying immutable GitHub check identities into a generated acceptance snapshot after merge.

### Handoff

- **OBSERVED:** `Handoff.md` is designed as the short recovery surface for current status, coordinates, files, evidence, risks, and next action. Recent Sprint 029-032 handoffs identify exact heads/CI/review and stop boundaries well enough for a fresh session.
- **OBSERVED:** Handoff works particularly well around dependency transitions: Sprint 030 points to Sprint 031 only after separate authorization; Sprint 031 distinguishes product merge from governance reconciliation; Sprint 032 blocks child work until Steering.
- **CONTRADICTORY EVIDENCE:** Handoffs are manually maintained and can freeze at pre-transition state. Sprint 022's Sprint and Iteration are complete while its `Handoff.md` still says active and “no product-code implementation has started.” Sprint 032's default-tree handoff still says PR #21 is draft/unmerged after it became the studied merge.
- **RECOMMENDATION:** Keep a concise handoff, but generate coordinates, Issue/PR state, exact head/checks, and completed statuses from authoritative records. Reserve prose for risks, context, and next decision.

### CurrentIndex, Relations, and Ledger

- **OBSERVED:** `CurrentIndex.yaml` is a large registry as well as a current pointer. `Relations.yaml` is the explicit graph. `Ledger.ndjson` is append-only operational history. Together they make later requirements/design/slices/review/verification and supersession reconstructable.
- **OBSERVED:** Current default-tree `active` points to `FEAT-007 / SPR-032 / IT-032-001 / SL-032-001-004`, while every referenced record has status `complete`. This represented “complete pending Steering” before merge but is semantically ambiguous as an active pointer.
- **OBSERVED:** Relations connect Feature authority and Study to requirements, architecture, designs, Sprints, Slices, implementation, verification, review, PRs, supersession, and post-merge reconciliation. This is the richest part of the local SDP model and a good Analyzer input.
- **OBSERVED:** The Ledger contains 814 events and a very large actor vocabulary. Common event counts include 89 `worker_completed`, 65 `rework_requested`, 64 `verification_completed`, 58 `slice_created`, 45 `review_completed`, and 39 `slice_closed`. Line 695 has a known timestamp inversion (`17:02:09` after lines at `17:03:00` and `17:03:30`) and is retained as immutable inherited evidence.
- **CONTRADICTORY EVIDENCE:** The default merge on 2026-08-19 contains a Ledger whose last event is 2026-08-12 and says PR #21 remains draft/open/unmerged. The append-only rule prevented correction in place, but no additive post-merge event was made.
- **RECOMMENDATION:** Retain stable graph IDs for Feature/Refactor, Requirement/Architecture/Design decision, Slice, verification/review, PR/commit, and release. Generate registry/status/PR/check fields. Reduce manual implementation IDs and verbose lifecycle events unless they support a decision, audit, or recovery. Allow correction events that point to the superseded assertion rather than leaving consumers to infer truth from GitHub.

### Steering and CurrentAssignment

- **OBSERVED:** `SDP/Steering/README.md` defines Issues as live assignment/discussion, draft PRs as diff/CI, and SDP as durable contract/state/evidence/handoff. `CurrentAssignment.yaml` is a strong machine-readable bridge.
- **OBSERVED:** The assignment binds repository, Issue and exact authority comments, active/base branch, exact base, PR, Feature/Sprint/Iteration/Slice/Implementation/Verification/Review, authorization flags, prerequisites, current and later scope, invariants, forbidden work, stop boundary, exact completed commit/tree/CI/review, residual findings, and publication checkpoints.
- **OBSERVED:** Issue #12 comments [`5345960799`](https://github.com/Hans-Einar/TerrainAnalyzer/issues/12#issuecomment-5345960799) and [`5345965181`](https://github.com/Hans-Einar/TerrainAnalyzer/issues/12#issuecomment-5345965181) are the actual Steering acceptance and merge disposition: PR #21 was accepted at `8163c74...`, merged, and child Issues #14/#15 were authorized from common base `77d71dc...`.
- **CONTRADICTORY EVIDENCE:** The studied `CurrentAssignment.yaml` still says `complete_pending_steering`, `pull_request_state: draft_open`, and forbids Issue #14/#15 branches. The studied Feature, Sprint, Handoff, CurrentIndex, and Ledger repeat that stale state even though the merge commit itself proves otherwise.
- **INFERENCE:** Merging a PR whose tree intentionally describes its own pre-merge stop state guarantees temporary repository staleness unless an authorized reconciliation follows. Earlier lifecycle transitions used Issue #17 or PR #20 for exactly this reason; PR #21 had no such reconciliation before child work began.
- **RECOMMENDATION:** Canonical `CurrentAssignment` should separate immutable assignment contract from generated live disposition. On merge, a bot/tool should create or propose the additive reconciliation and clear/advance the active assignment. A Steering chat remains supervisory; the Issue, PR, commit, checks, and reconciled repository snapshot remain authoritative.

### GitHub binding

- **OBSERVED:** Binding is explicit but manually duplicated across Issue, Feature, Sprint, CurrentAssignment, Relations, Ledger, PR description, verification, review, handoff, and FeatureBacklog. Exact comment authority is used where Issue body alone is too broad.
- **OBSERVED:** Branch naming evolved from `codex/issue-*` to `agent/issue-*`; both bind the Issue number. PR descriptions report Feature/Slice identities, boundaries, exact heads, tests, review, and stop state.
- **OBSERVED:** Issue closure and Feature completion are distinct. Issue #12 remained open after `FEAT-005`; Issues #1/#2/#4/#5 remain open even after their old Features were superseded. Issue #13 was closed only after a dedicated reconciliation.
- **RECOMMENDATION:** Store the Issue/comment authority, Feature/Refactor ID, base branch/SHA, branch, PR, accepted head, merge commit, check run, review, and release link once in a canonical binding record; render the other documents from it. Support `superseded`/`reconciled_by` for open historical Issues so open state is not mistaken for active work.

### Automated versus conversational convention

- **OBSERVED — automated:** `.github/workflows/pull-request-exact-head.yml` triggers on PRs to `main`, checks out and asserts the exact PR head, verifies fixture hashes, parses two traceability YAML files and every Ledger line, runs full backend/frontend tests, typecheck, production build, Playwright production browser tests, and diff whitespace checks. It uses PR-specific concurrency cancellation.
- **OBSERVED — automated but narrow:** The workflow checks YAML/NDJSON syntax, not duplicate keys, graph semantics, current-assignment truth, Issue/PR bindings, link existence, actor identity, or chronology. `scripts/verify-tier10-traceability.py` performs semantic checks only for the older Tier 10/Sprint 019 contract.
- **OBSERVED — convention/manual:** selecting IDs, authoring Issues, creating branches/draft PRs, keeping PR draft, assigning fresh agents, enforcing ownership boundaries, updating Feature/CurrentAssignment/CurrentIndex/Relations/Ledger/Handoff, validating strict duplicate-key YAML/paths, publishing/deploying, deciding merge/Issue closure, and post-merge reconciliation are manual/agent conventions.
- **OBSERVED:** The merge commit itself has no combined status/workflow run. The method deliberately accepts an exact PR head and separately records the merge identity.
- **RECOMMENDATION:** Automate schema/semantic validation, GitHub binding/status refresh, next-ID allocation, append-only correction/link checks, and post-merge reconciliation. Keep scope, architecture, safety, findings, and Steering decisions human/agent-reviewed.

## What worked well

- **OBSERVED:** Repository-first state plus exact coordinates allowed repeated Codex sessions to recover and continue complex work without chat memory.
- **OBSERVED:** Early SDP activation and draft PR publication bounded work before code and exposed stale pointers, dependency topology, and missing CI early.
- **OBSERVED:** The Issue→Feature→Slice model enabled major later product growth and architecture replacement without pretending it belonged to the original Tier plan.
- **OBSERVED:** Strong horizontal ownership (`A-005/A-006`, `DD-016..018`) plus vertical Slice gates prevented Geometry, Mesh, Diff, ModeUI, and SharedUI responsibilities from drifting together.
- **OBSERVED:** Exact-head CI and adversarial independent review repeatedly caught defects that green tests missed. Records explicitly show CI did not waive High/Medium findings.
- **OBSERVED:** Supersession is mostly honest. Earlier ephemeral Mesh/TerrainDiff work remains traceable while `FEAT-004` and later designs become current; merge commits are not falsely labelled reviewed heads.
- **OBSERVED:** Verification evidence is unusually reproducible: exact SHA/tree/run IDs, counts, browser repeat/retry policy, deterministic hashes, fixtures, and residual Notes are preserved.
- **OBSERVED:** `CurrentAssignment.yaml` is an excellent bounded authorization/stop contract and a strong seed for canonical Steering integration.
- **OBSERVED:** The preparation Slice demonstrated safe parallelization through explicit shared/disjoint ownership and reserved namespaces before child branches.

## Pain points and accidental complexity

- **OBSERVED:** State is duplicated across too many manually edited surfaces. The PR #21 merge made Feature, Sprint, Handoff, assignment, index, and Ledger stale simultaneously even though the Issue was current.
- **OBSERVED:** Post-merge truth sometimes requires a dedicated governance Issue/PR (Issue #17, PR #20). This protects exact-head evidence but adds a second delivery cycle and can be skipped.
- **OBSERVED:** The registry has become very large: 107 requirements, 124 implementation IDs, 81 review IDs, 69 verification IDs, and 814 Ledger events in under two months. Many actor/event spellings are near-duplicates, limiting machine analysis.
- **OBSERVED:** `CurrentIndex` is both registry and active pointer; “active” pointing to complete records encodes pending Steering implicitly rather than through a clear lifecycle state.
- **OBSERVED:** Review records are comprehensive but enormous and mix plans, failed passes, remediations, final evidence, and post-merge notes in one file. This improves auditability but makes current disposition expensive to discover.
- **OBSERVED:** The exact-head/evidence-only-commit problem creates multiple valid identities (product candidate, records head, accepted final head, merge commit, reconciliation head) that must be handled explicitly.
- **OBSERVED:** Feature semantics are stretched by `FEAT-007` preparation. Refactor semantics are underdeveloped: no `Refactors/`, no Refactor ID, and a stale implementation refactor plan despite open Issue #10.
- **OBSERVED:** Older stale handoffs and historical “active” paragraphs remain inside superseded Feature records. They are truthful snapshots but can mislead readers and analyzers without strong current/superseded rendering.
- **OBSERVED:** Formal GitHub review evidence is inconsistent. Local role separation may be genuine, but generic Ledger actors do not prove fresh context or independence.
- **OBSERVED:** The exact-head workflow has accumulated project-specific fixture and legacy MeshBuilder checks; it is valuable but not a generic SDP validator.

## Carry forward

- **RECOMMENDATION:** Stable `FEAT-###` plus first-class Refactor identity, each linked to an Issue/comment authority and one or more vertical Slices.
- **RECOMMENDATION:** Issue-assigned Codex Master with explicit scope/authority/stop boundaries, fresh bounded Workers, independent Reviewers, and separate Verifier when warranted.
- **RECOMMENDATION:** Early SDP-only activation, exact baseline, branch, draft PR, and explicit acceptance/merge/closure authority.
- **RECOMMENDATION:** Horizontal initial architecture plus Feature-local Study/requirement/design refinement and explicit supersession/preservation relations.
- **RECOMMENDATION:** Exact-head CI, exact-commit review, changed-head re-review after B/H/M fixes, and reviewed-head versus merge-commit separation.
- **RECOMMENDATION:** Machine-readable `CurrentAssignment` fields for Issue/comment, base/branch/PR, Feature/Refactor/Slice, authorization, forbidden scope, invariants, stop boundary, accepted head/check/review, and residual findings.
- **RECOMMENDATION:** Concise handoff plus reproducible verification records and append-only decision/audit events.
- **RECOMMENDATION:** Optional convergence/preparation Slice before parallel sibling Features when shared touchpoints/IDs would otherwise collide.
- **RECOMMENDATION:** Truthful supersession instead of rewriting immutable history; additive post-merge reconciliation, preferably generated.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward the original one-time phase/Tier decomposition as the owner of all future capability. Retain its useful horizontal skeleton, but current Features/Refactors should own later work.
- **RECOMMENDATION:** Do not require Sprint as the conceptual capability owner. TerrainAnalyzer's best recent evidence already makes Feature primary and Sprint an execution container.
- **RECOMMENDATION:** Do not use Feature for behavior-preserving preparation merely because no Refactor/Enabler type exists.
- **RECOMMENDATION:** Do not keep manually mirrored live GitHub state in six repository files. Generate it and make stale-state detection a gate.
- **RECOMMENDATION:** Do not preserve the current unbounded ID/event proliferation. Keep stable IDs only where they support authority, architecture, delivery, evidence, release, or recovery.
- **RECOMMENDATION:** Do not treat generic Ledger actor labels as sufficient proof of independent review; add an auditable review attestation.
- **RECOMMENDATION:** Do not retain stacked PRs as the default topology once a shared preparation gate can create clean sibling bases.
- **RECOMMENDATION:** Do not carry project-specific NN2000, CRS, mesh, browser-count, Apache, or exact fixture rules into canonical SDP. They belong in TerrainAnalyzer requirements, verification profiles, and project-local agent guidance.
- **RECOMMENDATION:** Do not infer current state from historical paragraphs or folder names. `status`, supersession, GitHub disposition, accepted head, merge commit, and reconciliation must be resolved together.
