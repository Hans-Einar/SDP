# HEOS SDP usage study

## Evidence frame and limitations

- **OBSERVED:** This report studies `Hans-Einar/HEOS` at the inventory's exact default identity `main@c41995ecc503684f4bb5200ae8f8da3665d57e1c` (tree `16d1094fd3db4600b3475333bce332c58817aaf1`), committed `2026-04-25T20:30:33Z`. The studied tree was read only.
- **OBSERVED:** GitHub state was re-inspected on `2026-08-21`. The current default head remains the studied commit. The current non-default branch `Version_1.0--Union@8d4bddc04bace27d695f0faa0c320bba1b323390` (tree `400ec6368a61fc6bcc37784296a1840d19e689f6`) contains an unmerged SDP installation and is reported separately.
- **OBSERVED:** In the fixed window `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`, four commits are reachable from advertised refs: `Version_2.0@9b5cd82315dd6a06edf2c08a99cbcedc8f48fdcf` at `2026-04-24T13:15:39Z`, default `main@c41995e...` at `2026-04-25T20:30:33Z`, then non-default SDP initialization commits `b804135a906e64aca32346e87fb6f1ff44047cae` and `8d4bddc...` on `2026-07-14T06:19:09Z` and `2026-07-14T06:32:00Z`.
- **OBSERVED:** No build was run for this study. The non-default SDP record itself says PlatformIO verification was blocked because `pio` was unavailable.
- **LIMITATION:** Default-tree truth contains no SDP evidence. Conclusions about SDP practice are necessarily limited to the current non-default branch and must not be presented as accepted `main` behavior.

## Repository facts and current GitHub state

- **OBSERVED:** `HEOS` is private, unarchived, and uses `main`. The default tree has 60 blobs/21 directories and contains source, libraries, simulator files and PlatformIO configuration.
- **OBSERVED:** At the studied default commit there is no `AGENTS.md`, `AGENTS-project.md`, `SDP/`, `Instructions/`, `Features/`, `Steering/`, `Sprints/`, `Refactors/`, `Fixes/`, `Traceability/`, review, verification, or `CurrentAssignment` path.
- **OBSERVED:** GitHub has no Issues and no open PR. Five PRs, [#1](https://github.com/Hans-Einar/HEOS/pull/1) through [#5](https://github.com/Hans-Einar/HEOS/pull/5), were merged into non-default `Version_2.0` on `2025-06-15`, outside the study window.
- **OBSERVED:** Those five PRs were small Codex-linked build/platform fixes. They were opened ready-for-review, merged within minutes, and have no PR comments, formal reviews, inline review comments or checks. PRs #1, #2 and #5 explicitly recorded failing build commands before merge; PRs #3 and #4 recorded passing PlatformIO commands in prose only.
- **OBSERVED:** The repository has no GitHub Actions runs, check runs or commit statuses for the studied default or non-default SDP head.

## Non-default SDP installation

- **OBSERVED — NON-DEFAULT EVIDENCE:** Commit `b804135...` added a full older SDP seed, root `Study.md`, phase documents, Sprint 001, Traceability and a CodeReview backlog directly to branch `Version_1.0--Union`; commit `8d4bddc...` then added current Toolkit assets, `.codex/skills`, root agent files, new template phase documents, manifests, `Fixes/`, `Refactors/`, `Releases/` and `Verification/`.
- **OBSERVED — NON-DEFAULT EVIDENCE:** `SDP/Framework/installed-toolkit.manifest.yaml` identifies Toolkit `0.2.0`, framework `1.0.0`, agent contract `1.0.0`, installer `0.2.0`, installation time `2026-07-14T06:25:22Z`, and source commit `bf20832bed618ab240cf87c17517fc31ea721311`.
- **OBSERVED — NON-DEFAULT EVIDENCE:** The same tree contains both substantive legacy documents such as `SDP/01--Mandate/01--Mandate.md` and new empty templates such as `SDP/01--Mandate/mandate.md`; it has singular `SDP/Refactor/README.md` and plural `SDP/Refactors/README.md`; root `Study.md` duplicates the substantive `SDP/02--Study/02--Study.md`.
- **OBSERVED — NON-DEFAULT EVIDENCE:** `AGENTS.md` makes non-delegated sessions Master and names Worker, Reviewer, Architect and task skills. `AGENTS-project.md` is only the unfilled installer guidance template.
- **OBSERVED — NON-DEFAULT EVIDENCE:** `Features/`, `Steering/` and `CurrentAssignment` remain absent. `Fixes/`, `Refactors/`, `Verification/` and new phase files contain templates/READMEs only; there is no executed Feature, Refactor, Fix or verification record.

## How SDP is used

### Work start, authority, branch and PR practice

- **OBSERVED — NON-DEFAULT EVIDENCE:** The durable SDP work contract is branch-local `M-001 / S-001 / SP-001 / IT-001`, created after a user-requested tagged-union experiment. The root Study says it is a basis for a future SDP folder and cites workstation paths to `farmStatistics` and `HSX`.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Sprint 001 is `Union Value System`; IT-001 initializes SDP and prepares a next candidate implementation slice, IT-002. Product-code implementation is explicitly `Not started` and the Handoff says to open IT-002 before coding.
- **OBSERVED:** There is no GitHub Issue, PR or draft PR for `Version_1.0--Union`, no SDP-to-GitHub identity binding, and no branch protection/check evidence. The two SDP commits were pushed directly to the non-default branch.
- **INFERENCE:** This is a planning/bootstrap experiment, not evidence that the proposed union implementation or SDP workflow was completed or accepted.
- **OWNER DIRECTION:** Issue #5 requires future SDP assignments to use a durable GitHub Issue, exact baseline, bounded Master, early draft PR and explicit stop condition.
- **RECOMMENDATION:** Treat the non-default seed as migration evidence only. Any continuation should begin with a new Issue that reconciles the branch against current `main`, selects one canonical document set and records exact branch/PR authority.

### Master, delegation and review roles

- **OBSERVED — NON-DEFAULT EVIDENCE:** `AGENTS.md` and `AgentLifecycleAndDelegation.md` define Master, bounded Worker and separate Reviewer roles. Product code should be delegated, and a Reviewer must be separate from the Worker.
- **OBSERVED — NON-DEFAULT EVIDENCE:** No Worker implementation, Reviewer result, Architect pass or Verifier pass occurred. `CodeReview/00--ReviewBacklog.md` contains three open review questions (one High, two Medium), and verification is blocked/unstarted.
- **INFERENCE:** Role separation is installed convention, not demonstrated HEOS practice.
- **RECOMMENDATION:** Retain the role boundary if work resumes, but do not generate role/evidence IDs before an actual Slice candidate exists.

### Feature, Refactor, Fix, Sprint, Iteration and Slice

- **OBSERVED — NON-DEFAULT EVIDENCE:** There is no Feature identity. The work is framed as one broad Sprint owning ten implementation items, with IT-001 for initialization and IT-002 proposed for core union files.
- **OBSERVED — NON-DEFAULT EVIDENCE:** IT-002 is a reasonably bounded vertical foundation Slice candidate: two core files, no Weight migration, no compatibility wrappers, no heap allocation, compile/static review evidence. It was not activated.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Refactor semantics are contradictory: the mandate is a complete replacement/refactor, but the work is modeled as a Sprint; both `Refactor/` and `Refactors/` only contain README scaffolding. No Fix exists.
- **INFERENCE:** The intent would fit a Refactor owner better than a Sprint owner. Sprint and Iteration add hierarchy but do not clarify capability identity.
- **RECOMMENDATION:** Model any renewed union-system effort as one Refactor bound to an Issue, implemented through small vertical Slices. Sprint/Iteration should be optional scheduling/grouping rather than the authority owner.

### Requirements, architecture, design and Study evolution

- **OBSERVED — NON-DEFAULT EVIDENCE:** A broad Study was expanded into 14 requirements, nine architecture decisions, six design-analysis items, ten design decisions and ten implementation items before product implementation began.
- **OBSERVED — NON-DEFAULT EVIDENCE:** The Study records a meaningful direction change: compatibility wrappers were first considered and then explicitly rejected in favor of a clean tagged-union replacement. It also links HEOS decisions to HSX architecture.
- **OBSERVED — NON-DEFAULT EVIDENCE:** No later implementation learning, Feature-local requirement/design refinement or convergence gate exists. All horizontal documents were seeded together.
- **INFERENCE:** The Study's alternative rejection is useful, but generating the entire design/traceability graph before the first core Slice produced waterfall-like breadth and a large untested state surface.
- **RECOMMENDATION:** Keep the architectural skeleton and recorded alternative, then allow Slice-local requirements/design refinement from implementation evidence rather than completing all planned design detail up front.

### Review, verification and handoff

- **OBSERVED — NON-DEFAULT EVIDENCE:** The Handoff clearly names branch, active Sprint/Iteration, next files, invariants and the PlatformIO blocker. This is sufficient to recover the planning state without chat.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Verification `V-001` is truthfully marked blocked; the Sprint is active and product implementation not started. The CodeReview backlog is prospective rather than falsely closed.
- **CONTRADICTORY EVIDENCE:** The branch was created in July, but substantive documents and seven Ledger events are dated `2026-04-25`; the installed Toolkit timestamp is July. The Git history proves the records were committed in July, not that the recorded April transitions occurred then.
- **CONTRADICTORY EVIDENCE:** New `SDP-project.manifest.yaml` has project `TBD` and all development IDs null, while legacy `CurrentIndex.yaml` says project HEOS, branch `Version_1.0--Union`, active `SP-001 / IT-001`. The installer did not reconcile old and new current-state formats.
- **RECOMMENDATION:** Preserve truthful blocked evidence and concise handoff, but require migration tooling to reconcile legacy state into the current manifest without creating duplicate canonical documents or backdated events.

### CurrentIndex, Relations, Ledger and automation

- **OBSERVED — NON-DEFAULT EVIDENCE:** Legacy `CurrentIndex` registers 14 requirements, nine architecture items, six design-analysis items, ten design items, ten implementation items, Sprint/Iteration, blocked verification and open reviews. `Relations` maps their graph; the Ledger has seven manually written events.
- **OBSERVED — NON-DEFAULT EVIDENCE:** The records contain no Issue/PR/commit edge and no CurrentAssignment. GitHub-native state is not derived or reconciled.
- **OBSERVED — NON-DEFAULT EVIDENCE:** No Steering record or actor exists. The
  branch-local Handoff and Sprint state substitute for Steering assignment and
  disposition, without a durable owner decision surface.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Toolkit assets and manifests provide installed-version automation, but there is no repository workflow, validation run or product verification proving the installation/state is valid.
- **INFERENCE:** The graph is detailed relative to zero implemented product Slices and duplicates state between legacy CurrentIndex and the new manifest.
- **RECOMMENDATION:** Keep only stable, non-derivable relationships. Generate Git/Issue/PR status and avoid pre-registering every planned implementation item in three files.

## What worked well

- **OBSERVED — NON-DEFAULT EVIDENCE:** The Study clearly records the rejected wrapper approach, target architecture, non-goals and verification blocker rather than hiding uncertainty.
- **OBSERVED — NON-DEFAULT EVIDENCE:** The Handoff narrows the next step to core union types and prevents an oversized simultaneous Weight migration.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Installed agent guidance establishes bounded Worker and independent Reviewer roles before product implementation.
- **INFERENCE:** These records could support safe restart if first reconciled with current GitHub/default-tree truth.

## Pain points, limitations and contradictions

- **OBSERVED:** Canonical `main` has no SDP at all; the only SDP is unmerged on a non-default branch with no PR or Issue.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Legacy substantive files coexist with empty Toolkit templates, two refactor folder spellings and duplicate Study content, leaving document authority ambiguous.
- **OBSERVED — NON-DEFAULT EVIDENCE:** CurrentIndex and the new project manifest contradict each other about active development state.
- **OBSERVED — NON-DEFAULT EVIDENCE:** Seven backdated Ledger events and a large traceability graph were committed months after their timestamps, limiting provenance.
- **OBSERVED:** Earlier HEOS PRs demonstrate quick task/PR execution but no review/check discipline; some merged with explicitly failing tests.
- **INFERENCE:** HEOS is stronger evidence of installation/migration hazards than of successful SDP adoption.

## Carry forward

- **RECOMMENDATION:** Carry forward the clear Study decision, non-goals, concise Handoff, truthful blocked verification and small proposed core Slice.
- **RECOMMENDATION:** Carry forward repository authority plus bounded Master/Worker/fresh Reviewer roles only when an actual Issue-bound Slice begins.
- **RECOMMENDATION:** Use HEOS as a migration test case for reconciling legacy CurrentIndex/Relations/Ledger with a versioned Toolkit manifest and current GitHub state.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not carry forward duplicate legacy/template phase documents, singular/plural Refactor scaffolds, unfilled project manifests, or multiple competing current-state authorities.
- **RECOMMENDATION:** Do not carry forward backdated synthetic events, broad preimplementation ID graphs, or a Sprint-first label for work whose intent is a Refactor.
- **RECOMMENDATION:** Do not treat a non-default branch install as repository adoption without an Issue, draft PR, review, verification and accepted merge.
- **RECOMMENDATION:** Do not carry forward quick merges with failing tests or no review/check evidence. The newer Issue-first and exact-evidence direction supersedes that practice.
