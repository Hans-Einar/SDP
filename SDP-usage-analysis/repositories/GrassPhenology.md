# GrassPhenology SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/GrassPhenology` at the inventory's exact default-branch identity `master@8ba7567aaaf0524d86bdb4675a09232b41496f36` (tree `9af93c6a558843f40cc65968a9ec788102b305e9`), committed `2026-07-14T07:37:07Z`. Repository files were read from that Git object; the repository was not changed.
- **OBSERVED:** GitHub state was re-inspected on `2026-08-21`. The current default head is still the studied commit. Later open draft PR #2 and its branch are non-default evidence and are identified as such below.
- **OBSERVED:** The fixed activity window is `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`. The inventory found 22 qualifying commits across six advertised refs. The first is `8dd3011c240ad349b9f7f16985245bf7d02bc082` at `2026-07-11T08:09:46Z`; the latest is unmerged PR #2's synthetic merge ref `f54d2a7a6e8338a9ad83ae679914b7c366b6c1c7` at `2026-07-22T20:15:06Z`.
- **OBSERVED:** The Git history starts with one `initial` commit on `2026-07-11` that already contains SDP records dated from `2026-06-28` through early July. Later Refactor-001 work is incrementally committed, but the earlier Sprint-by-Sprint event history cannot be corroborated by separate Git commits in this repository.
- **OBSERVED:** No product tests were rerun for this read-only study. Test, review and verification results below are committed claims and artifacts, not new execution evidence.
- **LIMITATION:** The installed SDP version is not determinable. The tree has root agent instructions and copied role skills but no `SDP.manifest.yaml`, `SDP-project.manifest.yaml`, or installed-toolkit manifest.

## Repository facts and SDP surfaces

- **OBSERVED:** `GrassPhenology` is private, unarchived, and uses `master`. The studied commit message is `Document RFI-004A workspace review and prepare RFI-005`.
- **OBSERVED:** Root `AGENTS.md` makes the repository authoritative, assigns any non-delegated top-level session the Master role, and requires the active Sprint or Refactor, Iteration and Slice to be read before work. `AGENTS-project.md` exists but is empty.
- **OBSERVED:** `SDP/` contains `01--Mandate` through `08--Implementation`, `Instructions/`, `Sprints/`, `Refactors/`, `Traceability/`, `CodeReview/`, and `Verification/`. `Features/`, `Steering/`, `Fixes/`, and any `CurrentAssignment` file are absent.
- **OBSERVED:** The default tree contains 28 Sprint folders, one active refactor programme (`Refactor-001 - POC Workbench Architecture`), two CodeReview records, 40 `VER-*.md` verification records, 250 parseable Ledger events, and a 1,900-line Relations graph.
- **OBSERVED:** `CurrentIndex.yaml` has no active Sprint, Iteration, Slice, refactor iteration or refactor slice. It marks `Refactor-001` active, RFI-001 through RFI-004A complete, and `RFI-005 - Weather Workspace Extraction` planned. This agrees with the Refactor README: no product work is currently authorized.
- **OBSERVED:** `SDP/README.md` is stale: it says Sprint 018 is the latest completed work and describes Mandate/Study stage, while `CurrentIndex` registers Sprints through 028 and Refactor-001 is active. The more focused Refactor README and traceability records carry the later state.
- **OBSERVED:** There are no GitHub Issues. GitHub currently has two PRs: merged PR [#1](https://github.com/Hans-Einar/GrassPhenology/pull/1) and open draft PR [#2](https://github.com/Hans-Einar/GrassPhenology/pull/2). Neither has PR comments, formal reviews, inline review comments, status checks or check suites.
- **OBSERVED:** GitHub Actions shows only a successful Dependency Graph run on the initial commit; no CI workflow/check applies to the studied default commit or either PR head.

## How SDP is used

### Work start, authority and GitHub binding

- **OBSERVED:** The repository's first durable start evidence is the imported Ledger: `sdp_bootstrap_started`, `slice_started`, Worker, Reviewer, verification and closure events for `SPR-001 / IT-001-001 / SL-001-001-001` on `2026-06-28`. It records a Slice contract before implementation, then corrective work after live verification failed.
- **OBSERVED:** Later committed Refactor work starts by opening one bounded RFI Slice, finalizing its contract, establishing an exact baseline, delegating implementation, obtaining fresh review, rerunning verification, closing traceability, and clearing active pointers before the next RFI is activated.
- **OBSERVED:** GitHub Issues are not an authority surface here because none exist. The Sprint/Refactor documents and traceability files are the work authority.
- **OBSERVED:** PR #1 says the initial Mandate, Study, Requirements, Architecture and CR-002 documents were committed directly to `master` before a PR workflow was requested; only the subsequent Design, implementation plan and Refactor index used branch `refactor-001-poc-workbench-architecture`. It was opened ready-for-review, not draft, and merged seven minutes later without GitHub review or checks.
- **OBSERVED — NON-DEFAULT EVIDENCE:** PR #2 was opened early as a draft from `agent/add-radar-precipitation-wishlist` on `2026-07-22`. Its only change is `features-wishlist.md`; it explicitly says no product Slice has begun. The file is not a canonical `Features/FEAT-*` record and is still unmerged.
- **INFERENCE:** Repository-side work contracts are strong, but branch/PR/Issue identities are not bound into `CurrentIndex`, `Relations` or `Ledger`. GitHub coordination is therefore episodic rather than a canonical SDP edge.
- **OWNER DIRECTION:** Issue #5 requires future SDP work to treat a high-quality GitHub Issue as the bounded assignment and to use an early draft PR as the durable execution/review surface.
- **RECOMMENDATION:** Preserve the repository contract quality, but bind each future Feature/Refactor assignment to its Issue, baseline, branch and draft PR. Do not use PR prose or chat request alone as the work authority.

### Master, Worker, Reviewer, Architect and Verifier

- **OBSERVED:** `AGENTS.md`, `AGENT-REMINDERS.md`, and local Instructions make the Master responsible for planning, active coordinates, traceability, verification evidence and handoff; product code belongs to a bounded Worker; review belongs to a separate fresh Reviewer.
- **OBSERVED:** Refactor ledger evidence uses explicit Worker and Reviewer passes. RFI-003 review found one High stale-target defect; RFI-004 review found one Medium ordered-update semantic regression; RFI-003A and RFI-004A documentation reviews found Medium issues. Work was corrected and re-reviewed before closure.
- **OBSERVED:** Verification is usually run and integrated by the Master, with durable `VER-*` records containing commands, screenshots, hashes, dimensions and limitations. A separately delegated Verifier role is not evidenced in the default tree even though verification is first-class.
- **OBSERVED:** An Architect role skill is installed, but Refactor architecture work is recorded as Worker/Master/documentation-review activity; no separately attributable Architect pass was found.
- **INFERENCE:** Independent Worker/Reviewer separation has real value here; Architect and Verifier are capabilities invoked by need rather than consistently distinct roles.
- **RECOMMENDATION:** Carry forward Master-per-Issue plus bounded Worker and fresh Reviewer. Use a separate Architect only when a Slice changes horizontal boundaries, and make Verifier independence proportional to risk while always preserving exact evidence.

### Sprint, Iteration and vertical Slice practice

- **OBSERVED:** The 28 Sprints are mostly one Iteration and one Slice each. Early Sprints include documentation additions, requirements/design baselines, implementation, small fixes and UI behavior changes under the same Sprint hierarchy.
- **OBSERVED:** Refactor-001 introduces a clearer peer execution surface: RFI-001/RFI-001A establish and correct the baseline; RFI-002/RFI-002A establish and correct shared contracts; RFI-003 adds a typed map-intent seam; RFI-003A converges architecture evidence; RFI-004 moves the Fields workflow vertically; RFI-004A independently studies the resulting pattern; RFI-005 is prepared but inactive.
- **OBSERVED:** Refactor Slices are vertical and behavior-preserving. They state exact goal, scope, expected files/responsibilities, invariants, non-goals, opening evidence, verification, fresh review, completion signal and hard stop. `RFI-004` moved a complete Fields workflow while retaining named transitional adapters and OpenLayers execution in its prior owner.
- **OBSERVED:** Corrective `A` Slices are opened rather than silently rewriting or reusing a closed Slice. Active pointers are cleared on closure, and prepared next work is explicitly not authorized.
- **OBSERVED:** `Feature` has no canonical meaning in the default tree. Product growth is represented by Sprints/use cases/requirements; later radar work is only a wishlist in an unmerged PR.
- **OBSERVED:** Fix is represented as Sprint 021 rather than a `Fixes/` record. Refactor is first-class locally but uses `Refactor-001` plus RFI IDs outside the Sprint hierarchy.
- **INFERENCE:** The valuable unit is the bounded vertical Slice and its evidence, not the fact that it is nested under a Sprint. One-Slice Sprints add an extra ID without adding planning value.
- **RECOMMENDATION:** Make Feature and Refactor peer intent owners. Retain the GrassPhenology RFI contract discipline for Slices; make Sprint/Iteration optional grouping or timebox metadata instead of mandatory wrappers.

### Later requirements/design and Studies/convergence

- **OBSERVED:** GrassPhenology repeatedly adds requirements and design after initial bootstrap. Sprints 003-013 expand the Study, requirements, architecture and design analysis; later implementation Sprints add requirements up through `R-112`; Refactor-001 adds its own Mandate, Study, requirements, architecture, design and implementation plan.
- **OBSERVED:** Refactor Study and living Architecture explicitly distinguish original hypothesis, implemented state, proposed state, evidence, revised conclusion, risks, alternatives and next decision. RFI-003A and RFI-004A are bounded convergence/review Slices after implementation checkpoints.
- **OBSERVED:** The Refactor's Behaviour First Rule requires baseline, move one vertical responsibility without redesign, automated/rendered verification, fresh review, closure/commit, and only then separately authorized redesign.
- **OBSERVED:** This is sequential evidence convergence, not multiple independent parallel Studies. There are no study IDs with claim/evidence levels or a formal multi-study convergence schema.
- **INFERENCE:** Feature-/Refactor-local Study and design evolution worked well. The repository does not prove that ActionCam-style parallel-study machinery is necessary for ordinary work.
- **RECOMMENDATION:** Carry forward initial horizontal skeleton plus local Study/requirements/design refinement inside later Feature/Refactor work. Require explicit architecture revision when ownership or dependency direction changes; use bounded convergence Slices when implementation evidence changes the design.

### Review, verification and handoff

- **OBSERVED:** Review and verification are substantive even though they are not GitHub-native. RFI records preserve adverse findings, rework, exact command outcomes, screenshots and known limitations. Green local build/tests do not waive behavioral review.
- **OBSERVED:** Verification records distinguish rerun evidence from inspected prior evidence. RFI-004A explicitly says the prior full product suite/rendered record was inspected, not rerun, and leaves human rendered verification pending.
- **OBSERVED:** Handoffs exist for every Sprint and record objectives, evidence, open risks and next action. Refactor state is integrated primarily in its README, iteration files, implementation notes and traceability rather than one current Refactor handoff.
- **CONTRADICTORY EVIDENCE:** The root SDP README is stale, while the active traceability and Refactor README are current. GitHub has no formal reviews/checks supporting the many fresh-review claims, so freshness is repository convention rather than externally attributable evidence.
- **RECOMMENDATION:** Keep exact-candidate, severity-based review; executable and rendered verification; explicit inspected-versus-rerun language; and concise recovery handoff. Generate stable status/PR/check fields to prevent the README/handoff drift seen here.

### CurrentIndex, Relations, Ledger and automation

- **OBSERVED:** `CurrentIndex` is both current pointer and very large registry. `Relations` links requirements, Sprints, Refactor/RFI, verification and review. The append-only Ledger preserves starts, findings, corrections, rechecks and closures.
- **OBSERVED:** There is no `Steering/` surface or `CurrentAssignment`. Work
  authority and next-step state are distributed across the Refactor records,
  CurrentIndex, Handoffs and conversational/GitHub context rather than one
  durable Steering assignment.
- **OBSERVED:** The 250-event Ledger is unusually useful around Refactor-001 because it records rejected claims and exact rework instead of presenting only a clean final story.
- **OBSERVED:** The same state is repeated in Refactor README, iteration files, implementation notes, verification, CurrentIndex, Relations and Ledger. `SDP/README.md` demonstrates that not every copy is maintained.
- **OBSERVED:** Tests, builds, YAML/NDJSON parsing, hash/dimension checks and rendered evidence are executable. Role freshness, branch/PR creation, GitHub binding, active-coordinate agreement and Reviewer identity are conversational/manual conventions; there is no repository CI enforcing them.
- **RECOMMENDATION:** Keep a small CurrentIndex for current coordinates, Relations for non-derivable stable graph edges, and Ledger for transitions/corrections. Generate reverse links, Git/PR/check state and registry views rather than duplicating them manually.

## What worked well

- **OBSERVED:** Small vertical Refactor Slices protected behavior while responsibility moved across architectural boundaries. Named temporary adapters and explicit removal conditions made transitional complexity visible.
- **OBSERVED:** Fresh review found real High/Medium behavioral and documentation defects after tests/builds passed; rework and re-review prevented premature closure.
- **OBSERVED:** Baseline screenshots, exact commands, hashes, dimensions, test counts and known-defect preservation made the work recoverable and auditable after session loss.
- **OBSERVED:** Corrective Slices and append-only Ledger events retained contradictions and changes of mind instead of rewriting history.
- **OBSERVED:** Later local Study/requirements/design evolution allowed product learning through implementation without abandoning the initial horizontal architecture discipline.
- **INFERENCE:** GrassPhenology is strong evidence for behavior-first vertical slicing, evidence-driven architecture evolution and independent review.

## Pain points, accidental complexity and contradictions

- **OBSERVED:** Twenty-eight mostly one-Slice Sprints, hundreds of registry entries, 1,900 Relations lines and 250 Ledger events impose substantial manual state for a repository with no Issue binding or CI validation.
- **OBSERVED:** `SDP/README.md`, Sprint handoffs, CurrentIndex, Refactor README and Ledger duplicate current state and can disagree.
- **OBSERVED:** Early repository history was imported in one initial Git commit, so detailed pre-import event timestamps and role transitions cannot be independently matched to incremental Git candidates.
- **OBSERVED:** GitHub review is disconnected from repository review. PR #1 merged without checks/reviews; PR #2 remains an unreviewed draft; no Issue provides authority or a durable stop contract.
- **OBSERVED:** `Feature` is only a wishlist term, Fix was implemented as a Sprint, and Refactor uses its own parallel RFI hierarchy. These semantics are locally understandable but not canonical or mutually consistent.
- **OBSERVED:** Verification evidence is strong but sometimes expensive and duplicated; repository-native screenshots and 40 verification records coexist with prose summaries and ledger copies.
- **INFERENCE:** The repository preserves high-value evidence but over-encodes scheduling/work coordinates while under-encoding the external assignment and GitHub delivery boundary.

## Carry forward

- **RECOMMENDATION:** Carry forward the behavior-first vertical Slice contract, exact opening baseline, explicit invariants/non-goals, named temporary adapters and hard stop before the next Slice.
- **RECOMMENDATION:** Carry forward bounded Worker ownership, fresh independent Reviewer, severity-based rework, exact post-rework verification and truthful inspected-versus-rerun evidence.
- **RECOMMENDATION:** Carry forward append-only corrections, corrective Slices, reproducible metrics with stated limitations, and living architecture that separates implemented from proposed structure.
- **RECOMMENDATION:** Carry forward Feature-/Refactor-local Studies and requirements/design refinement after the initial horizontal skeleton.
- **RECOMMENDATION:** Add the missing Issue/branch/draft-PR binding and a small machine-readable CurrentAssignment derived from, rather than duplicating, GitHub-native state.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward one mandatory Sprint and Iteration for every single Slice. Keep them only when they add a real timebox, multi-Slice grouping or coordination decision.
- **RECOMMENDATION:** Do not carry forward the large registry-plus-Relations-plus-Ledger duplication or hand-maintained current-state prose in multiple READMEs and Handoffs. Generate derivable views.
- **RECOMMENDATION:** Do not carry forward Issue-less material work, direct-default planning commits, or PRs without exact review/check evidence. Issue #5's Issue-first, early-draft-PR direction replaces this gap.
- **RECOMMENDATION:** Do not treat the unmerged `features-wishlist.md` as a canonical Feature model, or use Sprint as a substitute for Fix/Feature/Refactor semantics.
- **RECOMMENDATION:** Do not generalize GrassPhenology's domain-specific OpenLayers, SharedUI, weather/Sentinel or rendered-workflow rules into generic SDP. Preserve only their generic contract, evidence and migration lessons.
