# Recommended follow-up implementation Issues

These are proposed work contracts in dependency order. Issue #5 implements none
of them. The order deliberately pilots the sparsely observed owner-directed
model before schemas and downstream writers are frozen.

## 1. Approve a provisional workflow pilot contract

**Goal:** Convert Steering-selected parts of `ProposedSDPWorkflow.md` into an
explicitly experimental, non-Toolkit pilot contract and list unresolved choices.

Decide provisionally:

- Feature/Refactor/Fix/Slice identities and minimum records;
- Feature/Issue/Slice cardinality and reopen rules;
- optional Sprint/Iteration treatment;
- per-Issue assignment records and concurrent Issue-Master coordination;
- architecture-revision threshold;
- declared/observed/accepted state and proportional profiles.

Do not call the pilot schema canonical or publish/migrate the Toolkit.

**Depends on:** Steering acceptance of Issue #5's study boundary.

## 2. Run non-destructive real workflow pilots

**Goal:** Exercise the provisional semantics manually/project-locally before
locking machine schemas or write-capable tooling.

Use separate Issues/branches/draft PRs for at least:

- one Feature with later local requirement/design refinement;
- one behavior-preserving Refactor;
- one genuine small Fix;
- two concurrent Issue Masters from a common base with disjoint ownership and a
  convergence/merge-order gate; and
- one minimal/no-SDP repository profile.

Pilots add isolated experimental records only to their work branches, obtain
fresh review and report friction/cardinality/status/migration discoveries. They
do not migrate historical records or release a Toolkit.

**Depends on:** Issue 1.

## 3. Finalize semantics, schemas, neutral templates and fixtures

**Goal:** Revise the provisional contract using pilot evidence, then implement
the accepted machine contracts.

Deliver:

- Feature/Refactor/Fix/Slice and per-Issue assignment schemas/templates;
- concurrent assignment/index, dependency/conflict and merge semantics;
- general Ledger event envelope/vocabulary;
- semantic Relations/status/timestamp contracts;
- neutral templates separated from live Toolkit records;
- legacy/profile capability declarations;
- positive/negative fixtures including pilot failures and branch contradictions.

**Depends on:** Issues 1–2.

## 4. Update canonical documentation and Codex skills

**Goal:** Align How-SDP-Works, document guide, AGENTS templates and versioned
Codex skills with the accepted Issue/Feature/Refactor model.

Include deprecation guidance, proportional profiles, Issue templates,
concurrency/merge rules, exact-head review and traceability simplification. Do
not publish a Toolkit release yet.

**Depends on:** Issue 3.

## 5. Create separate ChatGPT `sdp` and update `steering-group` skills

**Goal:** Implement the role boundary in `SkillsAndRoles.md`.

Verify with scenarios for repository inspection, Issue authoring, concurrent
assignment, evidence acceptance, architecture escalation, migration ambiguity
and stop boundaries.

**Depends on:** Issue 4.

## 6. Add multi-generation read support to SDP-Analyzer

**Goal:** Read accepted schemas and observed legacy profiles before migrations
write them.

Implement capability discovery, Markdown/work-record coverage, preserved
unknown fields, declared/observed/accepted projections, concurrent assignment
and contradiction findings, plus representative fixtures. Keep Analyzer
read-only.

**Depends on:** Issue 3; may run in parallel with Issues 4–5.

## 7. Implement gh-sdp status and deterministic migration planning

**Goal:** Add read-only status plus no-write migration plan/result contracts for
the new method and legacy profiles.

Status reports Toolkit/source/profile, all current assignments/work owners,
GitHub observations, accepted evidence, release inclusion and contradictions.
Planning covers ownership/collisions/backups/preconditions/rollback with zero
mutation.

**Depends on:** Issue 3; may run in parallel with Issues 4–6 using the same
accepted fixtures rather than depending on Analyzer implementation.

## 8. Implement gh-sdp additive install/update/migrate apply

**Goal:** Safely install/update new managed skills/schemas/templates and apply an
approved deterministic migration plan.

Must preserve project-owned history, remain idempotent, recheck assumptions,
support rollback and never promote non-default evidence without authority.

**Depends on:** Issues 4 and 7.

## 9. Pilot the migration/write path

**Goal:** Exercise actual status/plan/apply behavior in representative real
repositories after the early workflow pilots shaped the contract.

Recommended profiles:

- no-SDP bootstrap: tplink or LogParser;
- dense legacy Sprint: LogClassifier or GrassPhenology;
- Feature-first: Lyndata or TerrainAnalyzer;
- hybrid Toolkit/project authority: weight_app_flutter;
- branch-only/stale-default: HEOS or HSX;
- nested SDP roots/safety: ActionCam.

Each migration uses its own Issue, preview, branch/draft PR, backup/rollback
evidence, exact validation and fresh review.

**Depends on:** Issues 6 and 8.

## 10. Add Analyzer GitHub connector, graph/path views and staleness UI

**Goal:** Visualize/validate the accepted Issue-to-release graph from repository
facts plus timestamped GitHub snapshots.

Implement entity/status lists and path queries before derived graph/Gantt views.
Show contradictions and unknowns; do not auto-repair.

**Depends on:** Issues 6–7 and migration fixtures from Issue 9.

## 11. Reconcile duplicate Toolkit authorities and release a compatible version

**Goal:** After both pilot waves, remove/deprecate duplicate root skill/payload/
template authorities, finalize migration notes and prepare the appropriate
Toolkit SemVer release under the existing two-phase publication gate.

This Issue must explicitly coordinate—but not silently merge or replace—the
separate installation-contract workstream in PR #4.

**Depends on:** Issues 4, 8 and 9 plus explicit Steering publication
authorization.

## 12. Evaluate write-assisted repair separately

**Goal:** Only after read-only detection is trusted, study reviewed repair
proposals for stale assignment/index/relation state.

No Analyzer or status command should mutate project records before this separate
authorization and rollback design exists.

**Depends on:** Issue 10 plus real contradiction/pilot evidence.
