# LogClassifier SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/LogClassifier` at the inventory's exact default identity `main@a38f9bb061f8ab77e25b801c06915c754d746e74` (tree `da17c2a9b42f07e47deaaf8ff58b26ef884b7885`), committed `2026-07-11T13:55:20Z`. The tree was read only.
- **OBSERVED:** GitHub state was re-inspected on `2026-08-21`; the default head is unchanged. There are no other advertised refs.
- **OBSERVED:** Both repository commits fall inside the fixed window: initial `a4811d6b6719ad70c8e0de6df5813b805b6430ba` at `2026-07-11T09:37:50Z` and the studied head at `2026-07-11T13:55:20Z`.
- **OBSERVED:** The initial commit already contains SDP events and work records dated `2026-07-03` through `2026-07-10`; the second commit adds Sprints 005-007 and their same-day July 11 records. Git cannot independently corroborate the many claimed Worker/Reviewer/Master candidates before those two aggregate commits.
- **OBSERVED:** No code/tests were run for this read-only study. Verification results are committed claims from the studied tree.
- **LIMITATION:** SDP generation/version is indeterminate. There is no manifest, installer record, agent entry file, local instruction set or versioned skill bundle.

## Repository facts and SDP surfaces

- **OBSERVED:** `LogClassifier` is private, unarchived and uses `main`. It contains Python detector/annotation tools, large source/generated image artifacts, CSV/JSON results, tracked `__pycache__`, and an `SDP/` folder.
- **OBSERVED:** The default tree has no `AGENTS.md`, `AGENTS-project.md`, `.codex/skills`, `SDP/Instructions/`, `Features/`, `Steering/`, `Refactors/`, `Fixes/`, `CodeReview/`, `Verification/`, `CurrentAssignment`, README or manifest.
- **OBSERVED:** SDP content consists of substantive Mandate, Study, Requirements, Architecture, Design and Implementation documents; seven Sprint folders; and `Traceability/{CurrentIndex.yaml,Relations.yaml,Ledger.ndjson}`.
- **OBSERVED:** Each Sprint contains a Sprint record, `ScrumIterations.md`, `implementationNotes.md` and `Handoff.md`. Review and verification have stable IDs but are embedded in those files and traceability rather than separate review/verification records.
- **OBSERVED:** `CurrentIndex.yaml` is a 56-item list covering requirements/design decisions, seven Sprints, seven Iterations, seven Slices, seven reviews and seven verifications. It has no active-coordinate block. `Relations.yaml` has 190 lines; `Ledger.ndjson` has 28 parseable events.
- **OBSERVED:** GitHub has no Issues, PRs, PR comments, reviews, checks, status contexts, workflow runs or branch-based evidence for this repository.

## How SDP is used

### Work start, authority and GitHub binding

- **OBSERVED:** The first Ledger event says the Master found no local SDP, created a minimal Sprint/traceability surface from the workstation reference `C:\Users\hanse\GIT\farmStatistics\SDP`, and opened `SPR-001 / IT-001-001 / SL-001-001-001` before implementation.
- **OBSERVED:** Subsequent work starts by creating one Sprint/Iteration/Slice contract containing goal, why now, expected files, invariants, non-goals, traceability IDs, verification and completion signal. The Ledger then records Worker, Reviewer/rework when needed, Master verification and closure.
- **OBSERVED:** No GitHub Issue is authoritative. No branch or draft PR was created; all durable repository history was pushed in two aggregate commits to `main`.
- **INFERENCE:** The local Slice contract—not GitHub—was the operational authority, while the actual prompts/owner decisions that initiated Sprints are not durably represented.
- **OWNER DIRECTION:** Issue #5 requires the GitHub Issue to become the bounded assignment contract and an early draft PR to become the diff/CI/review surface.
- **RECOMMENDATION:** Keep the concise local Slice contract, but bind it to an Issue, exact baseline, branch and draft PR. Do not reconstruct several completed Sprints into one later default-branch commit.

### Master, Worker, Reviewer, Architect and Verifier

- **OBSERVED:** Ledger records use explicit `master`, `worker` and `reviewer` roles. Handoffs repeatedly claim a separate Reviewer and Master verification, and state no agents remain open.
- **OBSERVED:** Sprint 006 shows the loop's value: Reviewer event `EVT-023` found a blocking persistence defect; Worker rework `EVT-024` fixed it; follow-up review and Master verification closed the Slice at `EVT-025`.
- **OBSERVED:** No Architect role is evidenced. Verification is performed/integrated by the Master, not a separately identified Verifier.
- **LIMITATION:** The same GitHub owner authored the aggregate commits, no GitHub review objects exist, and root agent instructions are absent. Role freshness and independence are therefore internal claims, not externally attributable evidence.
- **RECOMMENDATION:** Carry forward bounded Worker plus separate Reviewer and exact Master acceptance, but record candidate/provenance in a way that survives Git import and can be checked against PR state.

### Feature, Refactor, Fix, Sprint, Iteration and vertical Slice

- **OBSERVED:** There is no Feature, Refactor or Fix identity. Every change—from detector redesign and a complete GUI to a one-line wheel-direction sign swap—is modeled as a new Sprint with one Iteration and one Slice.
- **OBSERVED:** Slices are genuinely vertical. Sprint 003 delivers the complete manual annotation path from image/CSV load through navigation, ellipse placement, selection/delete and persistence. Sprint 006 delivers editing/growth-center behavior with compatibility, rework and persistence verification. Sprint 007 is a tiny behavior Fix but still receives seven stable requirement/design/Sprint/Iteration/Slice/review/verification IDs.
- **INFERENCE:** Vertical Slice discipline is useful; mandatory Sprint/Iteration wrappers and seven IDs for the one-line Sprint 007 change are disproportionate.
- **RECOMMENDATION:** Model owner-visible capability as Feature, structural change as Refactor and small regression/behavior correction as Fix. Use one or more vertical Slices as needed; make Sprint/Iteration optional grouping.

### Later requirements/design and Studies/convergence

- **OBSERVED:** Sprints 001-002 predate the formal Mandate/Study/Architecture set and focus on detector experiments. Sprint 003 adds a broad manual-GUI Mandate, GUI-backend Study, requirements, architecture, design and implementation plan after evidence showed automatic detection was insufficient.
- **OBSERVED:** Later requirements/design decisions `R-011` through `R-014` and `D-004` through `D-007` are added inside Sprint iteration records rather than extending the central Requirements/Design documents.
- **OBSERVED:** The PyQt Study compares three alternatives and makes a bounded decision. There is no independent parallel Study set, convergence gate, claim/evidence level or explicit architecture revision history.
- **INFERENCE:** The repository proves that requirements/design can evolve after initial implementation, but the evolution becomes fragmented when each later decision lives only in its Sprint record.
- **RECOMMENDATION:** Carry forward Feature-local Study and later refinement. Update or explicitly revise the relevant durable Requirements/Design/Architecture record when later behavior changes their truth; do not force every small Fix into a broad phase pass.

### Review, verification and handoff

- **OBSERVED:** Review and verification are practical and proportional to the local GUI: compile, self-test, PyQt import, offscreen behavior smokes, traceability parsing and artifact/CSV readback. Sprints explicitly retain limitations such as no manual visual pass and Alt-wheel platform behavior.
- **OBSERVED:** Handoffs list current objective, authoritative records, outcome, commands, IDs and residual risk. They are sufficient to understand each closed Slice without chat.
- **CONTRADICTORY EVIDENCE:** Sprint 006 and Sprint 007 Handoffs/implementation notes say the Slices are closed, but their Sprint status blocks still say `Active iteration` and `Active slice` are the closed IDs. Manual duplicate status was not fully reconciled.
- **OBSERVED:** There are no separate immutable review/verification files or PR/check identities; evidence is repeated across Iteration, implementation notes, Handoff, CurrentIndex and Ledger.
- **RECOMMENDATION:** Preserve exact behavioral tests, adverse findings, rework and limitations. Generate closure coordinates and provide one stable review/verification evidence object instead of repeating results in four documents.

### CurrentIndex, Relations, Ledger and automation

- **OBSERVED:** CurrentIndex provides a complete registry, Relations connects requirement/design/Slice/review/verification, and Ledger preserves chronological role events including rework.
- **OBSERVED:** The Ledger is compact enough to be useful, but many facts are manually duplicated in all four Sprint files plus traceability. CurrentIndex has no explicit `active: null` coordinates, so current state must be inferred from item statuses/Handoffs.
- **OBSERVED:** Compile/self-test/offscreen smoke and YAML/NDJSON parsing were automated locally. Issue intake, agent delegation, independence, branch/PR management, active-state agreement and GitHub reconciliation were conversational/manual conventions.
- **OBSERVED:** Generated outputs, bytecode caches and large experiment artifacts are committed alongside source and process records. The initial commit contains 1.78 million added lines, making review and evidence isolation difficult.
- **RECOMMENDATION:** Keep a small current pointer, stable non-derivable graph edges and append-only transitions. Generate registries/reverse links and separate reproducible evidence from bulk generated artifacts.

## What worked well

- **OBSERVED:** The one-Slice contracts kept scope explicit and supported autonomous Worker execution.
- **OBSERVED:** Independent review found a real cancel/persistence defect in Sprint 006 and forced rework before closure.
- **OBSERVED:** Offscreen GUI smokes made interaction regressions testable even without a manual display session, while records truthfully retained the manual-validation limitation.
- **OBSERVED:** Handoffs and chronological Ledger events preserve enough context to understand the seven changes after session loss.
- **OBSERVED:** The manual-GUI Study and vertical Sprint 003 show effective product pivoting after detector evidence changed the direction.

## Pain points, accidental complexity and contradictions

- **OBSERVED:** Seven one-Slice Sprints and 56 current-index items over-model small changes; Sprint 007's one-line sign change receives seven IDs and four Sprint files.
- **OBSERVED:** No AGENTS/Instructions file defines current role rules, and no Issue/PR/check evidence binds the repository claims to externally visible candidates.
- **OBSERVED:** Closed status is duplicated inconsistently across Sprint, Handoff, CurrentIndex and Ledger.
- **OBSERVED:** Requirements/design evolution is split between central documents and later Sprint-local IDs, weakening the durable horizontal model.
- **OBSERVED:** Two aggregate commits obscure the candidate-by-candidate history claimed by 28 detailed events and make formal review of massive generated artifacts impractical.
- **INFERENCE:** LogClassifier demonstrates useful vertical execution discipline but also the cost of applying the full Sprint/Iteration/ID ceremony to every small adjustment without GitHub integration.

## Carry forward

- **RECOMMENDATION:** Carry forward explicit vertical Slice contracts, bounded Worker ownership, separate Reviewer, rework-before-closure and exact behavioral verification with limitations.
- **RECOMMENDATION:** Carry forward concise Handoff and append-only role/outcome events, especially adverse findings and corrections.
- **RECOMMENDATION:** Carry forward later Feature-local Study/requirements/design evolution after real product evidence changes direction.
- **RECOMMENDATION:** Add Issue/baseline/branch/draft-PR binding and one authoritative current assignment so the evidence can be checked outside prose.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward mandatory Sprint and Iteration wrappers or full seven-ID graphs for one-line Fixes. Use proportional Fix records and only create IDs that support a real traceability decision.
- **RECOMMENDATION:** Do not carry forward repeated manual closure state across Sprint, Handoff, CurrentIndex, Relations and Ledger; derive registry and current-state views.
- **RECOMMENDATION:** Do not carry forward Issue-less/direct-main aggregate delivery or treat role labels in a later import commit as independently proven review provenance.
- **RECOMMENDATION:** Do not carry forward tracked caches and large generated experiment outputs as the primary verification surface. Keep reproducible commands, selected artifacts/hashes and storage policy instead.
