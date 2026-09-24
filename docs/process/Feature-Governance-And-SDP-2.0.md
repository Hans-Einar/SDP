# Feature Governance for SDP 1.x and Product-Governance Vision for SDP 2.0

Status: design proposal  
Proposal ID: `SDP-PROP-001`  
Applies to: SDP 1.x additive evolution and future SDP 2.0 design  
Date: 2026-07-17

## 1. Purpose

SDP already provides a strong repository-local method for turning an accepted mandate into verified implementation through Study, Requirements, Architecture, Design, Sprints, Iterations, Slices, Review, Verification and append-only Traceability.

The current method is weaker at product governance after initial planning. New ideas often appear while another implementation sequence is active. Without a governed intake path, teams and agents may:

- interrupt an active Slice with unplanned work;
- forget valuable ideas that were discussed only in chat;
- add Features without studying their architectural impact;
- create one oversized Sprint containing an entire release;
- confuse product backlog decisions with implementation decisions;
- trigger hidden refactors after implementation has already begun;
- lose Steering Group and Master-agent communication when chat history disappears.

This proposal defines:

1. a deliberately small, backward-compatible **SDP 1.x Feature extension** that can be added without redesigning the current method; and
2. a broader **SDP 2.0 product-governance model** to be studied and designed before implementation.

The guiding rule remains:

> Design horizontally. Implement vertically.

The new companion rule is:

> Capture ideas immediately. Integrate them deliberately. Implement them only after scheduling.

---

## 2. Non-goals

This proposal does not immediately:

- replace GitHub Issues, Projects or Pull Requests;
- require a database or hosted project-management service;
- require OpenAI API usage;
- change existing Sprint, Iteration or Slice semantics in installed SDP 1.x projects;
- force every small project to adopt every future SDP 2.0 governance layer;
- turn process compliance into an end in itself;
- require active implementation work to stop while a new Feature is documented;
- authorize a Feature merely because its document exists.

SDP must remain adaptable. Projects may use a smaller subset when appropriate, provided the repository remains truthful about what is and is not being used.

---

# Part I — Additive SDP 1.x Feature Extension

## 3. Why SDP 1.x needs only a small addition

The existing SDP 1.x implementation model is retained:

```text
Sprint or Refactor
    └── Iteration
        └── Slice or bounded Fix
```

A Sprint is a bounded delivery objective that may contain several Iterations and Slices. An Iteration is a planning/review cycle inside the Sprint. A Slice is the smallest independently implementable and reviewable vertical capability.

The Feature extension adds product intake **before** Sprint assignment without redefining the implementation hierarchy.

```text
New idea
    ↓
Feature proposal
    ↓
Integration Study
    ↓
Backlog decision
    ↓
Release scheduling
    ↓
Sprint assignment
    ↓
Existing Sprint → Iteration → Slice flow
```

This lets a project record a new Feature while another Sprint is active without changing `CurrentIndex.yaml` or expanding the active Slice.

---

## 4. Proposed SDP 1.x directory

Projects adopting the extension add:

```text
SDP/
└── Features/
    ├── FeatureBacklog.yaml
    ├── FEAT-001--short-name/
    │   ├── feature.md
    │   └── integration-study.md
    └── FEAT-002--short-name/
        ├── feature.md
        └── integration-study.md
```

The extension is additive. Existing lifecycle folders and active work remain unchanged.

### 4.1 Feature identity

Feature IDs use:

```text
FEAT-NNN
```

Examples:

- `FEAT-001`
- `FEAT-017`
- `FEAT-204`

IDs are stable and never reused. A rejected or withdrawn Feature retains its ID and history.

### 4.2 Feature folder naming

Recommended folder form:

```text
FEAT-017--browser-directory-selection
```

The ID is authoritative. The slug is descriptive and may be corrected without changing identity, provided traceability paths are updated explicitly.

---

## 5. Feature document contract

Each Feature has one `feature.md`.

Minimum contract:

```markdown
# FEAT-017 — Browser directory selection

Status: proposed
Feature ID: FEAT-017
Created: 2026-07-17
Proposed by: human | steering | master | user-feedback | issue

## Problem

What problem or opportunity exists?

## Intended outcome

What should become possible if the Feature succeeds?

## Users and stakeholders

Who benefits or is affected?

## Value

Why is this worth doing?

## Initial scope

What appears to belong in the Feature?

## Initial non-goals

What explicitly does not belong?

## Constraints

Known security, compatibility, regulatory, performance, cost or technology constraints.

## Related evidence

Links or stable IDs for Issues, Mandate items, Requirements, prior decisions, user reports or external evidence.

## Open questions

Questions that Integration Study must resolve.

## Current disposition

proposed | studying | backlog | scheduled | active | delivered | rejected | withdrawn | superseded
```

### 5.1 Feature status semantics

| Status | Meaning |
|---|---|
| `proposed` | Recorded but not yet accepted for study |
| `studying` | Integration Study is active |
| `backlog` | Accepted as valuable, not scheduled |
| `scheduled` | Assigned to a planned Release and usually a future Sprint |
| `active` | At least one authorized implementation Slice is active |
| `delivered` | Acceptance evidence exists in a released or explicitly delivered increment |
| `rejected` | Considered and intentionally declined |
| `withdrawn` | Proposer or owner removed it before delivery |
| `superseded` | Replaced by another stable Feature ID |

Creating a Feature document does not authorize implementation.

---

## 6. Integration Study contract

Every Feature must receive an Integration Study before it is scheduled, unless the project explicitly records a justified lightweight exception.

The Integration Study answers how the Feature fits the current system, not whether the idea merely sounds attractive.

Minimum contract:

```markdown
# Integration Study — FEAT-017

Status: draft | accepted | rejected | blocked
Study ID: ISTUDY-FEAT-017
Feature: FEAT-017

## Existing system context

Relevant architecture, design, active Releases, Sprints, interfaces and constraints.

## Horizontal impact analysis

- acquisition/input
- domain model
- storage/data
- application/use-case layer
- UI/SharedUI
- security/privacy
- compatibility/migration
- observability/operations
- verification/testing
- documentation/skills/tooling

Mark unaffected layers explicitly rather than omitting them ambiguously.

## Existing contracts reused

Which accepted boundaries can carry the Feature without change?

## Required contract changes

Which Requirements, Architecture or Design records must change?

## Alternatives considered

Meaningful options and tradeoffs.

## Dependency and sequencing analysis

Prerequisites, blockers and interactions with active or planned work.

## Refactor assessment

none | local | bounded-refactor | architectural-refactor

If a Refactor is required, identify or propose its stable REF ID and explain whether it must precede Feature delivery.

## Release recommendation

- backlog only
- target Release
- urgent exception
- reject
- split into multiple Features

## Proposed vertical delivery

Candidate Sprint and Slice boundaries. These remain proposals until scheduled.

## Risks and mitigations

Technical and delivery risks.

## Decision

Accepted, rejected, blocked or returned for revision, with rationale.
```

### 6.1 Integration Study outcome

An accepted study must end in one of these actionable outcomes:

1. **Backlog:** valuable but unscheduled.
2. **Schedule:** assign to a Release and later Sprint.
3. **Split:** create smaller Features and supersede or narrow the original.
4. **Refactor first:** create a bounded Refactor dependency.
5. **Reject:** preserve the evidence and rationale.
6. **Blocked:** identify the unresolved decision or evidence required.

---

## 7. Feature backlog contract

`SDP/Features/FeatureBacklog.yaml` is a concise index, not the full source of truth.

Recommended structure:

```yaml
schema_version: sdp-feature-backlog-v1
updated: 2026-07-17

features:
  FEAT-017:
    path: SDP/Features/FEAT-017--browser-directory-selection/feature.md
    integration_study: SDP/Features/FEAT-017--browser-directory-selection/integration-study.md
    status: backlog
    priority: normal
    target_release: null
    target_sprint: null
    depends_on: []
    supersedes: []

  FEAT-018:
    path: SDP/Features/FEAT-018--steering-records/feature.md
    integration_study: SDP/Features/FEAT-018--steering-records/integration-study.md
    status: scheduled
    priority: high
    target_release: REL-0.3.0
    target_sprint: SPR-004
    depends_on: []
    supersedes: []
```

Allowed priority values should remain deliberately simple:

```text
critical | high | normal | low
```

Priority does not authorize interruption of active work. Urgent work must still use an explicit exception, Fix or revised Sprint decision.

---

## 8. Release and Sprint assignment

Features are selected into Releases. Sprints deliver bounded subsets of a Release.

```text
Feature backlog
    ↓ selection
Release plan
    ↓ decomposition
Sprint
    ↓ execution
Iteration
    ↓ bounded implementation
Slice
```

A Release may contain multiple Sprints. A Sprint should not silently become the entire Release unless the Release is genuinely small and the repository records that decision.

Recommended relation chain:

```text
FEAT-017
  → REL-0.3.0
  → SPR-004
  → ITR-004
  → SLC-012
  → VER-SLC-012
  → REV-SLC-012
```

A Feature can require more than one Sprint or Slice. Delivery status must derive from explicit relations and evidence, not from the presence of a single completed Slice.

---

## 9. Steering records for SDP 1.x

The Feature extension should be compatible with repository-local Steering Group records.

Recommended additive structure:

```text
SDP/
└── Steering/
    ├── CurrentAssignment.yaml
    ├── MasterReports/
    │   └── MREP-SLC-012.md
    ├── Decisions/
    │   └── SDEC-SLC-012.md
    └── Blockers/
        └── BLK-SLC-012-001.md
```

This proposal defines the conceptual contracts. Exact schemas should be locked in a separate Feature/Integration Study before Toolkit distribution.

### 9.1 Master Report minimum semantics

A Master Report should identify:

- stable report ID;
- project;
- Sprint, Iteration and Slice;
- analyzed commit or tree identity;
- completion or blocked status;
- changed files;
- verification records;
- review records;
- discoveries and limitations;
- CurrentIndex state;
- explicit stop boundary;
- status `awaiting-steering-review` until Steering acts.

Suggested ID:

```text
MREP-SLC-012
```

### 9.2 Steering Decision minimum semantics

A Steering Decision should identify:

- stable decision ID;
- reviewed Master Report;
- reviewed commit;
- disposition: accepted, changes-required, blocked or rejected;
- rationale;
- authorized next work, if any;
- planning documents changed or required;
- human acknowledgement when required.

Suggested ID:

```text
SDEC-SLC-012
```

### 9.3 Current Assignment

`CurrentAssignment.yaml` may provide a concise machine-readable handoff from Steering to the repository-local Master.

It must not replace the complete Slice contract. It points to it.

Example:

```yaml
schema_version: sdp-steering-assignment-v1
assignment_id: ASGN-SLC-012
status: active
project: SDP-Analyzer
sprint: SPR-004
iteration: ITR-004
slice: SLC-012
slice_contract: SDP/Sprints/SPR-004/ScrumIterations.md
steering_decision: SDEC-SLC-011
base_commit: abcdef1234567890
constraints:
  create_pull_request: false
  commit_without_human_request: false
```

---

## 10. Ledger event flow

The existing append-only `SDP/Traceability/Ledger.ndjson` remains the historical event stream.

The Feature extension adds event types; it does not create a competing ledger.

Core flow:

```text
Feature proposed
    ↓
Integration Study accepted
    ↓
Feature scheduled
    ↓
Sprint assignment
    ↓
Slice completed
    ↓
Master Report submitted
    ↓
Steering Decision recorded
    ↓
Release completed
```

Recommended event types:

```text
feature_proposed
feature_study_started
feature_study_accepted
feature_study_rejected
feature_backlogged
feature_scheduled
feature_unscheduled
feature_split
feature_superseded
feature_implementation_started
feature_delivered
release_feature_added
release_feature_removed
sprint_feature_assigned
slice_feature_assigned
master_report_submitted
steering_review_started
steering_decision_recorded
steering_changes_required
steering_blocked
steering_accepted
release_completed
```

### 10.1 Event minimum fields

Every new event should follow the existing Ledger contract and include at least:

```json
{
  "event_id": "EVT-2026-07-17-001",
  "timestamp": "2026-07-17T12:00:00+02:00",
  "type": "feature_proposed",
  "subject_id": "FEAT-017",
  "status": "proposed",
  "source": "SDP/Features/FEAT-017--browser-directory-selection/feature.md"
}
```

When applicable, add stable related IDs rather than embedding prose-only references:

```json
{
  "event_id": "EVT-2026-07-17-008",
  "timestamp": "2026-07-17T15:30:00+02:00",
  "type": "steering_accepted",
  "subject_id": "SLC-012",
  "master_report_id": "MREP-SLC-012",
  "steering_decision_id": "SDEC-SLC-012",
  "reviewed_commit": "abcdef1234567890",
  "status": "accepted",
  "source": "SDP/Steering/Decisions/SDEC-SLC-012.md"
}
```

### 10.2 Ledger rules

- Append only.
- Never rewrite a historical event because a later decision changed.
- Corrections are new events referencing the superseded event or decision.
- Document creation, status transitions and Steering actions should have events.
- The Ledger is history, not the sole current-state model.
- `CurrentIndex.yaml`, FeatureBacklog and Release records hold current declared state.
- Relations connect durable identities.
- No event may claim verification, review, publication or Steering acceptance that does not exist.

---

## 11. Relations additions

`Relations.yaml` should eventually support explicit sections or records for:

```yaml
features:
  FEAT-017:
    status: scheduled
    integration_study: ISTUDY-FEAT-017
    release: REL-0.3.0
    sprints: [SPR-004]
    slices: [SLC-012, SLC-013]

integration_studies:
  ISTUDY-FEAT-017:
    feature: FEAT-017
    disposition: accepted
    architecture: [ARC-001]
    design: [DES-001]

master_reports:
  MREP-SLC-012:
    slice: SLC-012
    commit: abcdef1234567890
    verification: [VER-SLC-012]
    review: [REV-SLC-012]
    status: awaiting-steering-review

steering_decisions:
  SDEC-SLC-012:
    master_report: MREP-SLC-012
    reviewed_commit: abcdef1234567890
    disposition: accepted
    authorizes: SLC-013
```

Exact machine-readable schemas must be versioned before Toolkit enforcement.

---

## 12. Minimal adoption rules for SDP 1.x

A project may adopt the Feature extension safely when it follows these rules:

1. Do not change active Sprint/Iteration/Slice merely to record a Feature.
2. Create a stable Feature ID and Feature document.
3. Append `feature_proposed` to Ledger.
4. Add the Feature to FeatureBacklog.
5. Complete Integration Study before scheduling.
6. Update Requirements/Architecture/Design only when the accepted Integration Study requires it.
7. Create a Refactor record when integration requires significant structural change.
8. Assign accepted Features to a Release before assigning implementation Slices.
9. Keep implementation authorization in normal Sprint/Slice contracts.
10. Record Master Report and Steering Decision before authorizing the next boundary when Steering review is required.
11. Never silently convert an idea discussed in chat into active implementation scope.
12. Allow lightweight adaptation for small projects, but record any skipped governance step and rationale.

---

# Part II — SDP 2.0 Product-Governance Vision

## 13. Motivation

SDP 2.0 should treat product governance and implementation governance as connected but distinct systems.

The proposed high-level progression is:

```text
Mandate
    │
    ▼
Product Vision
    │
    ▼
Use Cases
    │
    ▼
Features
    │
    ▼
Integration Study
    │
    ▼
Release Planning
    │
    ▼
Sprint
    │
    ▼
Iteration
    │
    ▼
Slice
```

This progression should not be interpreted as a mandatory waterfall. It defines durable information layers and decision boundaries. Feedback may move backward, but changes must be explicit and traceable.

---

## 14. Proposed SDP 2.0 layers

### 14.1 Mandate

Defines why the project exists, constraints, authority, success conditions and fundamental boundaries.

The Mandate should remain comparatively stable.

### 14.2 Product Vision

Defines the intended product direction beyond the initial implementation:

- users and stakeholders;
- value proposition;
- long-term outcomes;
- product principles;
- strategic exclusions;
- expected evolution;
- success measures.

The Product Vision prevents each new Feature from redefining the product independently.

Suggested ID:

```text
VISION-001
```

### 14.3 Use Cases

Use Cases describe stable user or system goals, not implementation Features.

Examples:

```text
UC-001 — Analyze an SDP project
UC-002 — Inspect traceability disagreement
UC-003 — Review a Master Report
UC-004 — Select a local repository
UC-005 — Prepare a machine-readable report
```

Use Cases should be fewer and more stable than Features.

Each Use Case should identify:

- actor;
- trigger;
- preconditions;
- normal flow;
- alternate/error flows;
- desired outcome;
- security/privacy constraints;
- related success evidence.

### 14.4 Features

Features are planned product capabilities that serve one or more Use Cases.

A Feature may:

- add a new capability;
- improve an existing capability;
- address a product-level defect;
- introduce an integration;
- add a governance or operational capability.

Features are not implementation tasks.

### 14.5 Integration Study

Integration Study is the decision gate between product intent and delivery scheduling.

It evaluates how a Feature fits the existing horizontal system and whether Architecture, Design, data contracts, migration or Refactor work is required.

### 14.6 Release Planning

Release Planning selects accepted Features into a coherent public or internal delivery target.

A Release Plan should define:

- release ID/version target;
- selected Features;
- excluded/deferred Features;
- compatibility and migration goals;
- quality/release gates;
- expected Sprints;
- dependency order;
- risk budget;
- release acceptance criteria.

### 14.7 Sprint

A Sprint is a bounded delivery objective inside a Release.

A Sprint may include multiple Slices and one or more Iterations. It should produce a coherent increment and should not become a dumping ground for every Feature in the Release.

### 14.8 Iteration

An Iteration is a controlled learning and correction cycle within a Sprint.

Recommended semantics:

- one Sprint may have one or more Iterations;
- each Iteration selects or refines bounded Slices;
- discoveries may modify later Slices without silently expanding the active Slice;
- an Iteration ends with an explicit planning/review checkpoint.

The Iteration is not a retry count attached to one Slice. A Slice may receive bounded corrections/revisions while remaining the same Slice, but a material redesign belongs in a new Iteration or Refactor decision.

### 14.9 Slice

A Slice remains the smallest independently implementable, verifiable and reviewable vertical capability.

A Slice contract should continue to define:

- goal;
- why now;
- requirements and design references;
- expected modules/files;
- invariants;
- non-goals;
- verification;
- completion signal;
- discoveries policy;
- stop condition.

---

## 15. Proposed SDP 2.0 repository structure

A possible structure for study—not yet an implementation contract—is:

```text
SDP/
├── 01--Mandate/
├── 02--ProductVision/
├── 03--UseCases/
├── 04--Study/
├── 05--Requirements/
├── 06--Architecture/
├── 07--DesignAnalysis/
├── 08--Design/
├── 09--Implementation/
├── Features/
├── IntegrationStudies/
├── Releases/
├── Sprints/
├── Refactors/
├── Fixes/
├── Steering/
├── CodeReview/
├── Verification/
└── Traceability/
```

Alternative numbering and migration paths must be studied carefully. Renumbering existing lifecycle folders may create unnecessary migration cost. SDP 2.0 may instead retain current folder numbers and add Product Vision/Use Cases as unnumbered governance folders.

No folder structure should be adopted merely because it looks complete.

---

## 16. Steering Group in SDP 2.0

Steering becomes a first-class repository-local role and evidence boundary.

Responsibilities:

- product direction;
- Feature disposition;
- architecture integrity;
- release scope;
- acceptance of Master Reports;
- authorization of next work boundaries;
- resolution of blockers requiring product or architecture decisions;
- consistency across multiple repositories and dependent Toolkits.

Steering does not replace the repository-local Master. The Master executes accepted local contracts. Steering decides what should happen next when the answer is not already contained in the repository.

Proposed flow:

```text
Steering Assignment
    ↓
Repository-local Master
    ↓
Workers
    ↓
Verification
    ↓
Independent Review
    ↓
Master Report
    ↓
Steering Decision
    ↓
Next authorization
```

The repository, not chat history, carries the durable handoff.

---

## 17. GitHub integration philosophy

GitHub should be used as an operational interface, while SDP remains the project-owned source of truth.

Possible mapping:

| GitHub | SDP |
|---|---|
| Issue | References a `FEAT`, defect, blocker or task ID |
| Project board item | View of Feature/Release/Sprint state |
| Branch | Implementation workspace for an authorized Slice/Sprint |
| Pull Request | Review and integration surface |
| Check run | Verification evidence source |
| Release | Published identity linked from SDP release records |

SDP must not duplicate every GitHub field. It should preserve the durable decisions and evidence required when GitHub UI state changes or is unavailable.

Whether each Sprint should use one PR depends on project scale and workflow. SDP 2.0 should support, not mandate, policies such as:

- one PR per Slice;
- one integration PR per Sprint;
- direct local work with human-controlled commits;
- hybrid approaches for small projects.

The project must state its chosen policy in `AGENTS-project.md` or an equivalent project governance record.

---

## 18. Database and central dashboard

A future self-hosted dashboard may aggregate multiple project repositories.

Recommended architecture:

```text
Project repositories
    ↓ Git fetch/webhook/indexing
Central read model
    ↓
Dashboard and notifications
```

The central database should initially be a derived cache/read model. Project repositories remain authoritative for:

- Feature records;
- Integration Studies;
- Release plans;
- active work;
- Master Reports;
- Steering Decisions;
- verification/review evidence;
- append-only Ledger history.

This avoids two competing sources of truth.

A later API-enabled agent loop may consume the same contracts. Manual ChatGPT/Codex workflows can use them without API cost by reading and writing Git records through normal connectors and local Git.

---

## 19. Complexity controls

SDP 2.0 must explicitly defend against process bloat.

### 19.1 Progressive adoption

Small projects may use:

```text
Mandate → Feature → Integration Study → Sprint → Slice
```

Larger projects may use the complete Product Vision, Use Case, Release and Steering model.

### 19.2 Required versus optional records

Every schema should mark fields and records as:

- required for correctness;
- required only above a project complexity threshold;
- optional guidance.

### 19.3 No duplicate planning documents

A Feature document explains product intent. A Slice contract explains bounded implementation. They must not contain competing copies of the same detailed design.

### 19.4 Explicit exceptions

Projects may skip a step when the repository records:

- what was skipped;
- why it was safe;
- who accepted the exception;
- what evidence still applies.

### 19.5 Automation should reduce burden

Templates, validators, Analyzer rules and dashboards should remove clerical work rather than add it.

---

## 20. Migration strategy

### Phase 1 — SDP 1.x additive Feature extension

- add `SDP/Features/` templates;
- add FeatureBacklog schema;
- add Integration Study template;
- add optional Steering record templates;
- extend Ledger event schema/examples;
- extend Relations schema additively;
- update skills to recognize proposed/backlogged Features without activating them;
- teach SDP-Analyzer to report Feature/Steering consistency when document coverage exists.

### Phase 2 — Pilot

Pilot the extension in a few active repositories:

- SDP;
- gh-sdp;
- SDP-Analyzer;
- one application repository.

Collect friction, missing fields, redundant ceremony and migration issues.

### Phase 3 — SDP 2.0 Study

Run a dedicated Study covering:

- Product Vision and Use Case contracts;
- Feature/Integration Study lifecycle;
- release portfolio planning;
- Sprint/Iteration semantics;
- Steering contracts;
- GitHub integration;
- schemas and compatibility;
- installer migration;
- skill updates;
- Analyzer support;
- central dashboard/read model.

### Phase 4 — SDP 2.0 design and compatibility plan

Only after the Study:

- fix schemas;
- define directory structure;
- define migration from 1.x;
- define minimum and full profiles;
- define SemVer/toolkit release impact;
- define deprecation periods;
- define cross-version Analyzer behavior.

### Phase 5 — Implementation

Implement through normal Releases, Sprints, Iterations and Slices. SDP 2.0 must use the governance model it introduces.

---

## 21. Proposed acceptance criteria for the SDP 1.x extension

The additive extension is ready for Toolkit release when:

1. Feature and Integration Study templates are stable.
2. FeatureBacklog has a versioned schema.
3. Ledger event additions have schema/examples and append-only tests.
4. Relations additions are backward compatible.
5. Installer adds missing templates without replacing populated project-owned records.
6. Master/Architect/Traceability skills understand Feature states and do not activate backlog Features.
7. A Feature can be recorded during an active Sprint without changing CurrentIndex.
8. Accepted Integration Study can schedule a Feature into Release/Sprint traceability.
9. Refactor dependency can be represented explicitly.
10. Master Report and Steering Decision contracts are versioned or explicitly deferred.
11. At least two real project pilots succeed.
12. Documentation distinguishes SDP 1.x Feature extension from future SDP 2.0.

---

## 22. Proposed next work

The immediate next step should be one repository-local Feature in the SDP repository:

```text
FEAT-SDP-001 — Additive Feature governance for SDP 1.x
```

Its Integration Study should decide:

- exact folder names;
- schemas;
- templates;
- Steering contract scope;
- Ledger event additions;
- Relations compatibility;
- installer behavior;
- skill changes;
- validation/tests;
- pilot repositories;
- Toolkit release target.

A separate Feature should govern the larger redesign:

```text
FEAT-SDP-002 — SDP 2.0 product-governance architecture
```

That Feature should remain backlog/study work until the 1.x extension has been piloted.

---

## 23. Decision summary

### SDP 1.x

Add only the minimum governance needed to capture and safely schedule new Features:

```text
Feature
    ↓
Integration Study
    ↓
Backlog / Release scheduling
    ↓
Existing Sprint → Iteration → Slice implementation
    ↓
Master Report
    ↓
Steering Decision
```

Use the existing Ledger for the complete historical flow:

```text
Feature proposed
    ↓
Integration Study accepted
    ↓
Feature scheduled
    ↓
Sprint assignment
    ↓
Slice completed
    ↓
Master Report
    ↓
Steering Decision
    ↓
Release completed
```

### SDP 2.0

Study and design a broader, adaptable product-governance model:

```text
Mandate
    ↓
Product Vision
    ↓
Use Cases
    ↓
Features
    ↓
Integration Study
    ↓
Release Planning
    ↓
Sprint
    ↓
Iteration
    ↓
Slice
```

The 1.x extension solves the immediate coordination problem without prematurely locking the full 2.0 architecture.
