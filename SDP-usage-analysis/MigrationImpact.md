# Migration and downstream impact

Status: recommendation only; no repository migration or tooling change is
implemented by Issue #5.

## Migration principles

- **RECOMMENDATION:** Migration is additive, previewable, deterministic and
  profile-aware. Existing project-owned records, IDs and Git history are never
  overwritten merely to match the new layout.
- **RECOMMENDATION:** Do not retroactively invent Features, Issues, reviews,
  verification or Ledger events. Historical Sprint/Iteration/Slice records stay
  historical; new work adopts the new semantic owner prospectively.
- **RECOMMENDATION:** Default-branch truth and non-default/open evidence remain
  distinct. Migration does not silently promote an active branch or draft PR to
  canonical project state.
- **RECOMMENDATION:** Every plan identifies exact source/default commit,
  installed profile/capabilities, proposed file operations, collisions,
  preserved paths, validation and rollback/backup evidence.
- **RECOMMENDATION:** Mutable GitHub state is observed at a timestamp and never
  written into project decisions without explicit reconciliation authority.

## Compatibility profiles observed in the corpus

| Profile | Representative repositories | Migration posture |
|---|---|---|
| No SDP | LogParser, python, tplink | Minimal prospective bootstrap; no historical reconstruction. |
| Legacy phase/Sprint, no versioned manifest | LogClassifier, SharedUI, GrassPhenology | Preserve records/IDs; add profile marker, Issue/work-owner templates and derived-state validation. |
| SDP only on non-default/open evidence | HEOS, HSX, map_tracker, RadarData | Report branch evidence; require owner Issue/merge decision before default migration. |
| Hybrid Toolkit plus older project authority | weight_app_flutter | Resolve duplicate template/authority classes without overwriting populated documents. |
| Project-local Feature/Steering pilots | Lyndata, TerrainAnalyzer | Map capabilities/fields to accepted schemas; preserve extensions and contradictions. |
| Nested/task-local SDP islands | ActionCam | Register applicable roots/scopes; do not force a single repository-wide safety profile. |
| Canonical Toolkit authoring | SDP | Separate neutral installation templates from live Toolkit records and remove duplicate authority after compatibility gates. |
| Downstream method consumer | gh-sdp, SDP-Analyzer | Consume versioned contracts; do not independently redefine them. |

## Proposed migration phases

### Phase 0 — inspect and freeze a plan

1. resolve repository/default/current branch identities;
2. discover AGENTS, manifests, SDP roots and capability/profile markers;
3. classify populated, template, managed, project-owned, unknown and conflicting
   files;
4. parse CurrentIndex/Relations/Ledger without assuming a current schema;
5. compare declared state with Git/GitHub observations;
6. emit a no-write plan and explicit limitations.

### Phase 1 — add compatibility identity

- Add/upgrade a project manifest with schema version, installed Toolkit source,
  declared compatibility profile and recognized capabilities.
- Add no Feature/Refactor/assignment state unless actual work exists.
- Preserve unknown/project-local sections and record their owning profile.

### Phase 2 — add future-work surfaces

- Add neutral `Features/`, `Refactors/`, `Fixes/` and compact Steering assignment
  templates only when missing.
- Update managed AGENTS/skills to Issue-Master and Feature/Refactor semantics.
- Leave populated lifecycle, Sprint, Handoff, review, verification and
  traceability records untouched.
- Mark old Sprint-first work as `legacy-supported`; do not relabel it.

### Phase 3 — reconcile current state

- Require an owner/Steering Issue when default state, non-default work and
  GitHub disagree.
- Create one prospective current assignment for authorized new work.
- Clear or advance pointers only from real Issue/merge/acceptance evidence.
- Add correction/supersession events under the versioned general event contract;
  never backdate fabricated transitions.

### Phase 4 — validate and pilot

- Run schema, path, relation, chronology, Git/GitHub contradiction and project
  verification checks.
- Pilot at least one no-SDP repository, one dense legacy Sprint repository, one
  Feature-first repository, one branch-only repository and one nested SDP case.
- Preserve rollback/backup and exact before/after evidence.

### Phase 5 — deprecation window and removal

- Analyzer reads both generations and reports migration opportunities.
- gh-sdp warns rather than mutates unsupported/ambiguous cases.
- Only a later Steering-authorized release may remove deprecated Toolkit
  aliases/templates after pilot and compatibility gates pass.

## SDP-Analyzer implications

### Upstream-first dependency

- **RECOMMENDATION:** SDP-Analyzer must not invent final Feature, Refactor, Fix,
  assignment, status or event schemas before SDP accepts them. It consumes
  versioned contracts and exposes unknown/partial compatibility truthfully.

### Required normalized entities

- GitHub Issue and exact authority comment;
- Feature, Refactor, Fix and Study;
- Requirement and architecture/design decision;
- optional Sprint/Iteration group and Slice;
- branch, PR, commit/tree and check run;
- verification, review, finding and Steering disposition;
- Release, tag and GitHub Release.

Every fact/edge records provenance: repository declaration, append-only event,
Git object, timestamped GitHub observation, Markdown extraction or derivation.

### Required relations and path queries

Support `authorizes`, `implements`, `refines`, `supersedes`, `depends_on`,
`executed_by`, `head_at`, `merged_as`, `verified_by`, `reviewed_by`,
`accepted_by`, `included_in` and `published_as`.

The primary completeness query is:

```text
Requirement / decision
  -> Feature | Refactor | Fix
  -> Slice
  -> exact PR head / commit
  -> qualifying verification
  -> current independent review
  -> Steering disposition
  -> Release inclusion / publication
```

Missing or contradictory links produce explainable findings, not invented
entities or a single opaque health score.

### Declared/observed/accepted projection

- Declared state comes from repository records.
- Observed state comes from Ledger, Git and timestamped GitHub snapshots.
- Accepted state requires qualifying verification/review/Steering evidence.

Analyzer shows all projections and contradictions such as:

- assignment says draft/open after PR merge;
- closed Issue but repository says active;
- completed Slice without qualifying verification/current-head review;
- reviewed head differs from current PR head;
- Handoff/index older than a merge/decision;
- release record without tag/GitHub Release;
- branch-only work presented as default truth;
- event timestamp impossible relative to its containing commit.

### Compatibility and staleness

- Use versioned capability profiles rather than a binary old/new schema.
- An absent Feature folder in a legacy profile is not an error.
- Preserve unknown fields and analyze supported neighbors.
- Parse Markdown/work records before claiming ecosystem coverage; the current
  Analyzer's structured-core-only parser cannot see many Feature, Steering,
  review, verification and handoff facts.
- Staleness rules are entity/state/profile-specific and accept an explicit
  analysis time. Missing or date-only evidence yields `unknown`, not false
  precision.

### Recommended Analyzer implementation order

1. accepted schemas/capability profiles and representative fixtures;
2. manifest/profile discovery and Markdown/work-record coverage;
3. lifecycle projection plus declared/observed/accepted contradictions;
4. versioned GitHub snapshot connector;
5. entity/relation/status listings and path queries;
6. derived graph/timeline/Gantt and release views;
7. repair proposals only under a separate write-authorized Issue.

Analyzer remains read-only in the analysis path.

## gh-sdp implications

### Contract ownership

- **RECOMMENDATION:** Canonical SDP owns schemas, neutral templates, ownership
  classes, migration-plan/result contracts and conformance fixtures. gh-sdp is a
  portable validating client; it must not embed a second interpretation.

### Future install/update payload

After the canonical contracts are accepted, installation may add:

- neutral Feature/Refactor/Fix records/templates;
- compact `Steering/CurrentAssignment.yaml` template;
- Issue-oriented Instructions and managed AGENTS/skills;
- general event schema and semantic Relations vocabulary;
- profile/capability manifest fields;
- validator/configuration hooks.

It must not seed a live Feature, assignment, Sprint, release or Toolkit's own
CurrentIndex into a consumer.

### `gh sdp status`

Report separately:

- client version/release;
- installed Toolkit release/source, payload/plan digest and schema capabilities;
- project profile and migration warnings;
- declared Feature/Refactor/Fix, current assignment and Slice;
- observed Issue/branch/PR/head/check/review/merge facts with timestamp;
- accepted evidence/Steering state;
- target Release and inclusion/publication identities;
- contradictions, unsupported fields and offline/unknown facts.

Status is read-only. It does not overwrite repository declarations because
GitHub differs.

### Preview/apply migration

- A preview is deterministic, portable, ordered and contains no volatile
  absolute paths or future timestamps.
- It lists ownership class, precondition/hash, action, backup, collision policy,
  validation and rollback for every mutation.
- Apply rechecks all assumptions and stops with zero mutation on unsupported
  schema, ambiguous authority, path/link hazard or changed precondition.
- Repeating the same plan is idempotent.
- Historical/non-default records are never merged, renamed or promoted without
  explicit owner authority.

### GitHub-oriented operations

Future commands may scaffold an Issue contract, create the stable assignment
binding, open an early draft PR, capture a timestamped GitHub snapshot and
validate declared/observed state. These must require explicit caller authority
for external writes and preserve Issue #5's bounded workflow.

### gh-sdp product-local requirements

Archive, path, link, locking, journal, recovery, authentication and
cross-platform filesystem rules remain gh-sdp engineering requirements. They
are not generic SDP workflow fields, though their evidence and safety gates use
the canonical assignment/review model.

## Migration verification matrix

At minimum, fixtures must cover:

1. empty/no-SDP bootstrap;
2. minimal existing SDP;
3. pre-manifest Sprint-first project;
4. legacy CurrentIndex/Relations/Ledger envelope;
5. duplicate numbered authority plus empty Toolkit templates;
6. singular/plural path collision;
7. project-local Feature/CurrentAssignment extension;
8. nested/multiple SDP roots;
9. stale default with active non-default refs;
10. open PR synthetic merge ref versus actual unmerged state;
11. malformed/unsupported/future schema;
12. backdated/future-dated event and GitHub contradiction;
13. same-version/different-source unreleased Toolkit;
14. idempotent repeat and changed-precondition race;
15. backup/rollback and no-mutation failure paths.

## Migration success condition

A migration succeeds only when the exact post-apply tree validates, project
verification remains truthful, no project-owned history is lost, default versus
non-default authority is unchanged without explicit decision, and a fresh
Analyzer/Reviewer can explain every transformed or preserved fact.
