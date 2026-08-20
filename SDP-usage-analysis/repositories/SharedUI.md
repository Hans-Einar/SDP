# SharedUI

## Repository facts

- **OBSERVED:** Repository: `Hans-Einar/SharedUI` (private); default branch:
  `main`; exact study commit:
  `8cb61651107ff4fa23ad283504dc2f4ffc6e5e11`, committed
  `2026-06-29T13:24:43Z`.
- **OBSERVED:** The fixed-window inventory found two qualifying commits:
  `b0247bf042e2e239b33889e92ee38117147f3cbd` at
  `2026-06-28T09:36:34Z` and the study commit at
  `2026-06-29T13:24:43Z`, both on `main`. The first changed 107 files with
  15,563 insertions and 873 deletions; the second changed one documentation file.
- **OBSERVED:** The current default branch still equals the inventory study
  commit. GitHub has no pull requests. Current Issue
  [#1](https://github.com/Hans-Einar/SharedUI/issues/1), created after the study
  commit on `2026-08-07`, is open with no comments and no implementation branch
  or PR.
- **OBSERVED:** With no PR, there are no SharedUI PR comments, submitted
  reviews, review threads or PR checks to inspect. The default tree has no
  `.github/workflows/`.
- **OBSERVED:** SDP is installed as a project-local older generation. Its
  `SoftwareDevelopmentProcess.md` identifies “Working draft v0.1” dated
  `2026-04-05`; no SDP manifest, installer metadata or released Toolkit
  identity determines a canonical installed version. `package.json` version
  `0.1.0` and `SharedUI-0.1.0.tgz` are product-package identities, not SDP
  versions.

## SDP and agent locations

- **OBSERVED:** Root `AGENTS.md` defines the authority chain: project Mandate,
  `SDP/AGENT-REMINDERS.md`,
  `SDP/Instructions/SoftwareDevelopmentProcess.md`, then the active Sprint or
  Refactor. There is no `AGENTS-project.md`.
- **OBSERVED:** All operating instructions are in `SDP/Instructions/`:
  `README.md`, `SoftwareDevelopmentProcess.md`, `HandoffProcess.md` and
  `HandoffTemplate.md`.
- **OBSERVED:** Foundation records use numbered phase folders
  `SDP/01--Mandate` through `SDP/07--Implementation`. Six implementation Sprint
  folders are nested under `07--Implementation`; `SDP/Sprints/README.md`
  reserves the top-level Sprint location but says it is unused.
- **OBSERVED:** Structural work uses singular `SDP/Refactor/` with two numbered
  tracks. Refactor 001 contains Mandate, Study, Design, implementation plan,
  notes and `Handoff.md`. Refactor 002 contains the same core records plus a
  nested Sprint/Iteration, but no handoff.
- **OBSERVED:** Review evidence is under `SDP/CodeReview/`, with one repository-
  wide DRY/SRP review and a README. No `SDP/Features/`, `Steering/`,
  `Traceability/`, `Verification/`, `Fixes/`, `CurrentAssignment`,
  `CurrentIndex.yaml`, `Relations.yaml` or `Ledger.ndjson` exists.

## How work starts and who is authoritative

- **OBSERVED:** `AGENTS.md` tells an agent to read the Mandate, reminders,
  process and active Sprint/Refactor. The process says to identify the active
  authoritative folder, restate scope/non-goals/output, update SDP before broad
  code changes, implement, verify locally, and record notes only after verified
  completion.
- **OBSERVED:** No repository record binds completed work to a GitHub Issue,
  branch, PR or exact commit. No Issue authorized the six historical Sprints or
  either Refactor.
- **OBSERVED:** Current Issue #1 is a later, high-quality Study/Architecture
  request. It states purpose, owner direction, contradictions, cross-project
  discovery, questions, required outputs, guardrails, downstream relationships
  and independent-review requirement. It has not yet produced repository work.
- **INFERENCE:** SharedUI demonstrates repository-first work recovery but not an
  Issue operational boundary. Issue #1 is evidence of later intended practice,
  not evidence that the current SDP executed that practice.
- **OWNER DIRECTION:** Issue #5 requires GitHub Issues to become the operational
  assignment boundary while repository records remain authoritative.
- **RECOMMENDATION:** Combine SharedUI's concise local authority chain with an
  Issue-first assignment and exact branch/PR/head evidence; do not replace the
  repository-local architecture and execution contracts with Issue prose.

## Master, delegation and roles

- **OBSERVED:** Instructions refer only to generic “agents” and sessions. No
  Master, Worker, Architect, Verifier or Reviewer contract exists, and no
  delegation record names a bounded subject or exact head.
- **OBSERVED:** Refactor 001's handoff records “Previous agent role” and “Next
  agent role,” but these are sequential context-transfer labels, not evidence of
  Master/Worker delegation or independent review.
- **OBSERVED:** Current Issue #1 requires an independent architecture review,
  but no branch, review artifact or GitHub review exists yet.
- **INFERENCE:** Durable agent-facing instructions and local handoff worked
  without role specialization, but independent provenance and parallel task
  ownership cannot be reconstructed.
- **RECOMMENDATION:** Preserve the authority-chain and local handoff strengths;
  add bounded Master-per-Issue, delegated Worker/Architect/Verifier/Reviewer
  identities, exact subject/head and independent dispositions.

## Feature, Refactor, Fix, Sprint, Iteration and Slice

- **OBSERVED:** There are no first-class Feature records. “Feature” appears only
  descriptively in process/non-goal prose.
- **OBSERVED:** The initial implementation strategy explicitly pairs horizontal
  design with vertical implementation slices. Six sequential Sprints are named
  as capability Slices; each Sprint has high-level goal/scope/dependencies/
  phases/exit criteria, a `ScrumIterations.md`, and implementation notes.
- **OBSERVED:** Sprints 001–004 each have one completed Iteration, Sprint 005 has
  two, and Sprint 006 has four. Iterations are detailed task/meeting surfaces,
  while Sprint documents remain high-level. Verification commands are recorded
  in the implementation notes.
- **OBSERVED:** Refactor is a first-class later-work owner. The process explicitly
  says it captures later architecture or scope shifts when earlier assumptions
  stop fitting. Refactor 001 was promoted from a repository-wide CodeReview and
  executes ordered structural slices; Refactor 002 adds its own later Study,
  Design, plan and nested Sprint/Iteration.
- **OBSERVED:** No first-class Fix records or stable Fix contract exists. Newly
  found fixes may be listed in an iteration, but there is no distinct durable
  Fix hierarchy.
- **INFERENCE:** SharedUI strongly supports Refactor as a capability owner and
  vertical Slice as an implementation unit. It also shows Sprint and Slice were
  conflated: the Sprint folder itself owns the named Slice, while Iteration is a
  task cycle beneath it.
- **OWNER DIRECTION:** Issue #5 says capability should belong to Feature or
  Refactor rather than Sprint and asks whether Sprint/Iteration remain useful.
- **RECOMMENDATION:** Carry forward Refactor's self-contained
  Mandate/Study/Design/plan/evidence and the vertical increment contract. Make
  Feature/Refactor the outcome owner; retain Sprint only as optional grouping or
  timebox and Iteration only when a genuine planning/learning cycle adds value.

## Initial design and later requirements/design/Studies

- **OBSERVED:** The initial flow is specification-first and linear: Mandate,
  one broad Study, Requirements, Architecture, Design Analysis, authoritative
  Design and an ordered six-Sprint implementation plan.
- **OBSERVED:** Architecture/Design establish stable schema, renderer/state,
  component/token, data/async and auth boundaries. The implementation strategy
  requires each Sprint to cut through the minimum necessary layers and leave a
  coherent verifiable increment, with explicit exceptions for risk-reduction
  work leaning toward one layer.
- **OBSERVED:** Later work can revise assumptions through a Refactor-local
  Mandate, Study and Design. Refactor 002 used that mechanism to correct a
  consumer-composition mismatch long after the initial design, then delivered
  a schema-driven Agro CRM path through a nested vertical Sprint.
- **OBSERVED:** Studies are one broad foundation Study or one Study per
  Refactor. There is no observed split into parallel independent studies,
  synthesis artifact, evidence-level scheme or convergence gate.
- **INFERENCE:** SharedUI demonstrates the desired horizontal-skeleton/vertical-
  growth principle and later design refinement, but later capability growth is
  forced through Refactor because Feature has no canonical place.
- **RECOMMENDATION:** Preserve horizontal responsibility boundaries and allow
  Feature- or Refactor-local Study/requirements/design refinement. Add parallel
  studies and a named convergence/review gate only when uncertainty domains are
  genuinely independent.

## Review, verification and handoff

- **OBSERVED:** The persistent CodeReview identifies concrete DRY/SRP problems,
  explains why they are review findings rather than an approved Refactor, and
  proposes tracks. Refactor 001 then promotes that evidence into an active
  structural programme with scope, non-goals, ordered work and verification.
- **OBSERVED:** The review does not identify a reviewer, reviewed commit,
  severity disposition or fresh-context provenance. GitHub supplies no second
  review surface.
- **OBSERVED:** Verification is embedded in Sprint/Refactor implementation notes
  as local commands and smoke-test outcomes. There are no separate verification
  records, immutable verified heads, CI runs or machine-readable results.
- **OBSERVED:** Refactor 001's local `Handoff.md` is a strong recovery surface:
  authority, non-goals, structure, traps, done/not-done, exact next step, open
  questions, verification and worktree notes. Other Sprints and Refactor 002
  have no handoff.
- **RECOMMENDATION:** Preserve promotion of review findings into separately
  authorized Refactors and local work-owned handoff. Add exact-head verification,
  independent reviewer identity/disposition and changed-head invalidation.

## Traceability, Steering, GitHub binding and automation

- **OBSERVED:** There is no CurrentIndex, Relations, Ledger, Steering assignment
  or stable machine-readable ID graph. Relationships are Markdown paths and
  narrative dependencies.
- **OBSERVED:** There is no historical branch/PR/Issue binding. Current Issue #1
  links downstream repositories and authority documents in prose only.
- **OBSERVED:** Product verification is automated through npm scripts and tests,
  while process state, completion, handoff and review are conversational/manual.
  No GitHub Actions workflow publishes the recorded checks.
- **INFERENCE:** Low schema ceremony made the repository approachable, but exact
  assignment/review/release reconstruction is impossible. Markdown status had
  no automated contradiction detection.
- **RECOMMENDATION:** Add only the minimal assignment and durable semantic edges
  needed for Issue -> Feature/Refactor -> Slice -> verification/review/release;
  derive live GitHub and containment state rather than copying it into prose.

## What worked well

- **OBSERVED:** Root agent guidance is concise, points to repository authority,
  and enforces stable architectural boundaries and a default composition path.
- **OBSERVED:** Horizontal design paired with verifiable vertical Sprints
  produced working schema, renderer, components, data/auth integrations and
  examples rather than isolated layer-only completion.
- **OBSERVED:** The CodeReview-to-Refactor promotion separated diagnosis from
  authorization and gave structural remediation its own scope and non-goals.
- **OBSERVED:** Refactor-local Study and Design enabled later architectural
  correction instead of pretending all consumer requirements were known during
  the initial foundation.
- **OBSERVED:** Local implementation notes and the Refactor 001 handoff make
  substantial work recoverable without chat history.
- **INFERENCE:** SharedUI's strongest reusable contribution is a clear local
  authority chain plus horizontal boundaries, vertical increments and
  self-contained later Refactors.

## Pain points, contradictions and limitations

- **OBSERVED — state contradiction:** Root `AGENTS.md`, Refactor 001 Mandate and
  plan all say Refactor 001 is active. Refactor 002's Mandate remains “Draft,”
  its plan remains “Planned,” and its `Active Slice` says broad implementation
  must not continue; nevertheless its nested Sprint/Iteration and extensive
  implementation notes say all planned phases and multiple follow-ups completed.
- **OBSERVED:** There is no single machine-readable active pointer, so a fresh
  agent cannot resolve that contradiction deterministically.
- **OBSERVED:** Completion and verification are self-reported without exact
  commits, CI or independent reviews. A 107-file, 15,563-line commit carried much
  of the recent state, reducing review and traceability granularity.
- **OBSERVED:** Sprint is both execution container and capability/Slice owner;
  Feature has no durable meaning; Fix has no contract.
- **OBSERVED:** Initial phase/Sprint documentation is extensive, but later
  capability discovery has only Refactor as a legal home. This can misclassify
  genuinely new functionality as structural correction.
- **OBSERVED:** Issue #1 is broad and cross-repository but has no assignment
  coordinates, branch/draft-PR contract, stop reporting state or milestone
  comments yet.
- **OBSERVED:** No repository or GitHub access limitation was encountered. The
  absence of PR evidence limits conclusions about review/check practice.

## Carry forward

- **RECOMMENDATION:** Concise root `AGENTS.md` that points to repository-local
  authority and the active work contract rather than duplicating product detail.
- **RECOMMENDATION:** Stable horizontal architecture/design responsibilities and
  small vertical implementation outcomes.
- **RECOMMENDATION:** Self-contained Refactor with local Study, Design, plan,
  evidence and handoff, promoted from review rather than started as incidental
  cleanup.
- **RECOMMENDATION:** Locality-based handoff with done/not-done, exact next legal
  step, open risks, verification and worktree state.
- **RECOMMENDATION:** Implementation notes updated after verification, augmented
  in future with immutable head/check identities.
- **RECOMMENDATION:** Current Issue #1's evidence-before-universalization,
  explicit guardrails and independent architecture review requirement.

## Legacy / do not carry forward

- **RECOMMENDATION:** Do not preserve the initial one-time phase plan as the
  only legal route for new capabilities after foundation.
- **RECOMMENDATION:** Do not use Sprint as capability owner or encode “Slice”
  only in a Sprint folder name.
- **RECOMMENDATION:** Do not make Refactor the catch-all for both structural
  remediation and genuinely new product capability.
- **RECOMMENDATION:** Do not keep simultaneous “active,” “draft,” “planned” and
  “completed” state in unvalidated prose.
- **RECOMMENDATION:** Do not treat self-reported local commands as immutable
  verification or the persistent CodeReview as independently proven review.
- **RECOMMENDATION:** Do not carry complete work in giant commits without an
  Issue/PR/head evidence chain when the future Issue-driven model is available.
- **RECOMMENDATION:** Do not add manual CurrentIndex/Relations/Ledger ceremony
  merely to imitate newer repositories; add stable state only where it enables
  assignment recovery, semantic traceability or contradiction detection.
