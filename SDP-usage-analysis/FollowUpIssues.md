# Recommended follow-up implementation Issues

These are proposed work contracts in dependency order. Issue #5 implements none
of them.

## 1. Accept the canonical work semantics

**Goal:** Turn the Steering-approved parts of `ProposedSDPWorkflow.md` into a
normative method specification.

Decide:

- `FEAT-NNN`, `REF-NNN`, `FIX-NNN`, `SLC-NNN` semantics/statuses;
- Feature/Issue/Slice cardinality and reopen rules;
- Sprint optionality and Iteration removal/optional use;
- architecture-revision threshold;
- Issue/Master/Steering responsibilities;
- declared/observed/accepted state and proportional profiles.

**Depends on:** Steering acceptance of Issue #5 study.

## 2. Define versioned schemas, neutral templates and conformance fixtures

**Goal:** Implement the accepted machine contracts without changing consumers
implicitly.

Deliver:

- Feature/Refactor/Fix/Slice and compact CurrentAssignment schemas/templates;
- general Ledger event envelope/vocabulary;
- semantic Relations/status/timestamp contracts;
- neutral project templates separated from live Toolkit records;
- legacy/profile capability declarations;
- positive/negative fixtures including branch-only and contradiction cases.

**Depends on:** Issue 1.

## 3. Update canonical documentation and Codex skills

**Goal:** Align How-SDP-Works, document guide, AGENTS templates and versioned
Codex skills with the accepted Issue/Feature/Refactor model.

Include deprecation guidance, proportional profiles, Issue templates, exact-head
review and traceability simplification. Do not publish a Toolkit release yet.

**Depends on:** Issues 1–2.

## 4. Create separate ChatGPT `sdp` and update `steering-group` skills

**Goal:** Implement the role boundary in `SkillsAndRoles.md`.

Verify with scenarios for repository inspection, Issue authoring, evidence
acceptance, architecture escalation, migration ambiguity and stop boundaries.

**Depends on:** Issues 1–3.

## 5. Add multi-generation read support to SDP-Analyzer

**Goal:** Read accepted schemas and observed legacy profiles before migrations
write them.

Implement capability discovery, Markdown/work-record coverage, preserved unknown
fields, declared/observed/accepted projections, contradiction findings and
representative fixtures. Keep Analyzer read-only.

**Depends on:** Issue 2; may run in parallel with Issues 3–4 after schemas freeze.

## 6. Implement gh-sdp status and deterministic migration planning

**Goal:** Add read-only status plus no-write migration plan/result contracts for
the new method and legacy profiles.

Status reports Toolkit/source/profile, assignment/work owner, GitHub observations,
accepted evidence, release inclusion and contradictions. Planning covers
ownership/collisions/backups/preconditions/rollback with zero mutation.

**Depends on:** Issues 2 and 5's compatibility profile outputs.

## 7. Implement gh-sdp additive install/update/migrate apply

**Goal:** Safely install/update new managed skills/schemas/templates and apply an
approved deterministic migration plan.

Must preserve project-owned history, remain idempotent, recheck assumptions,
support rollback and never promote non-default evidence without authority.

**Depends on:** Issues 3 and 6.

## 8. Pilot the workflow and migrations

**Goal:** Exercise the accepted contracts in representative real repositories.

Recommended pilots:

- no-SDP minimal bootstrap: tplink or LogParser;
- dense legacy Sprint project: LogClassifier or GrassPhenology;
- Feature-first: Lyndata or TerrainAnalyzer;
- hybrid Toolkit/project authority: weight_app_flutter;
- branch-only/stale-default: HEOS or HSX;
- nested SDP roots/safety profile: ActionCam.

Each pilot uses its own Issue, branch/draft PR, fresh review and exact evidence.

**Depends on:** Issues 3, 5 and 7.

## 9. Add Analyzer GitHub connector, graph/path views and staleness UI

**Goal:** Visualize/validate the accepted Issue-to-release graph from repository
facts plus timestamped GitHub snapshots.

Implement entity/status lists and path queries before derived graph/Gantt views.
Show contradictions and unknowns; do not auto-repair.

**Depends on:** Issues 5–6; use pilot fixtures from Issue 8.

## 10. Reconcile duplicate Toolkit authorities and release a compatible version

**Goal:** After pilots, remove/deprecate duplicate root skill/payload/template
authorities, finalize migration notes and prepare the appropriate Toolkit SemVer
release under the existing two-phase publication gate.

This Issue must explicitly coordinate—but not silently merge or replace—the
separate installation-contract workstream in PR #4.

**Depends on:** Issues 3–8 and explicit Steering publication authorization.

## 11. Evaluate write-assisted repair separately

**Goal:** Only after read-only detection is trusted, study reviewed repair
proposals for stale assignment/index/relation state.

No Analyzer or status command should mutate project records before this separate
authorization and rollback design exists.

**Depends on:** Issue 9 plus real contradiction/pilot evidence.
