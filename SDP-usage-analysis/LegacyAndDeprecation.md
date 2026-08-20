# Legacy and deprecation analysis

## Classification rule

- **RECOMMENDATION:** Mark a practice legacy only when a newer contract replaces
  its purpose or the corpus shows that it causes ambiguity, drift or
  disproportionate cost. Age alone is not a reason.
- **RECOMMENDATION:** Preserve historical records and IDs. Deprecation changes
  prospective authoring/validation; it does not rewrite repositories into a
  fictional new history.
- **RECOMMENDATION:** Distinguish `legacy-supported`, `deprecated-for-new-work`
  and `unsupported/contradictory`. Analyzer and migration tooling must continue
  to read supported legacy profiles during a documented window.

## Preserve as durable foundations

The following are not legacy:

- repository/GitHub evidence over chat memory;
- living Mandate, Study, Requirements, Architecture, Design Analysis and Design;
- horizontal responsibility/ownership/contract boundaries;
- small vertical delivery;
- explicit invariants, non-goals, verification, review and stop conditions;
- exact candidate identities and truthful limitations;
- fresh adversarial review and changed-head re-review;
- append-only correction/supersession history;
- SemVer release identity separated from work identity;
- truthful two-phase release publication.

## Deprecated for new work

### One-time waterfall implementation ownership

- **OBSERVED:** Canonical `How-SDP-Works`, older repositories and the July
  proposal route later delivery primarily through the original Implementation/
  Sprint hierarchy. Multiple repositories then add later requirements/design in
  Sprint prose or misclassify new capabilities as Refactors.
- **RECOMMENDATION:** Keep numbered foundation documents, but deprecate the
  assumption that all future work must fit the initial implementation plan.
  Feature/Refactor-local evolution replaces it.

### Sprint or Tier as capability owner

- **OBSERVED:** Many repositories use Sprint as the capability, Study or Fix
  owner; SDP-Analyzer also uses Tier. Recent TerrainAnalyzer already makes
  Feature primary, and GrassPhenology/SharedUI make Refactor primary.
- **RECOMMENDATION:** Deprecate Sprint/Tier as the required semantic owner.
  Retain them as optional roadmap/timebox/grouping metadata.

### Mandatory Iteration

- **OBSERVED:** Six reports have one Iteration for every observed Sprint, and
  many others repeatedly use 1:1 wrappers.
- **RECOMMENDATION:** Deprecate mandatory Iteration IDs/folders. Retain Iteration
  only for a real multi-Slice learning/replanning cycle.

### “Slice” for every bounded activity

- **OBSERVED:** Activation, baseline, Study, tooling, verification,
  publication and governance-closeout records are called Slices; other Slices
  become large programmes.
- **RECOMMENDATION:** Reserve vertical Slice for an end-to-end implementable
  outcome. Use bounded task/gate terminology for other work.

### Permanent project Master

- **OBSERVED:** Older AGENTS contracts assume the top-level session owns the
  active project lifecycle. Recent Issue practice repeatedly re-anchors fresh
  Masters to bounded assignments.
- **RECOMMENDATION:** Deprecate reliance on one long-lived Master context. The
  repository and Steering Group preserve continuity; the Codex Master owns one
  Issue.

### Issue-less or late-PR material work

- **OBSERVED:** Several repositories have direct-main aggregate commits,
  conversational assignments or draft PRs opened only after substantive work.
  The evidence is harder to reconstruct and review.
- **RECOMMENDATION:** Deprecate Issue-less material implementation and late PR
  creation where GitHub is available. Preserve proportional exceptions for
  emergency response and tiny fixes with immediate durable reconciliation.

### Manually mirrored GitHub state

- **OBSERVED:** CurrentAssignment examples and many Handoffs/indexes copy Issue,
  PR, head, check and merge state; all exercised examples drift.
- **RECOMMENDATION:** Deprecate hand-authored mutable GitHub status fields.
  Keep stable URLs/identities and generate timestamped observations.

### CurrentIndex as historical registry

- **OBSERVED:** Dense projects place tens or hundreds of Requirements, Slices,
  reviews and verification objects in `CurrentIndex`, obscuring its current
  pointer purpose.
- **RECOMMENDATION:** Deprecate registry use. CurrentIndex holds only current
  pointers/profile/release target; Analyzer derives registries from records and
  relations.

### Hand-maintained reverse/containment relations

- **OBSERVED:** Large Relations files repeat paths and edges that can be derived
  and are frequent review-finding sources.
- **RECOMMENDATION:** Deprecate authored reverse edges, path indexes and routine
  containment. Keep semantic, non-derivable relations.

### Ledger as chat/activity transcript

- **OBSERVED:** High-volume ledgers repeat agent lifecycle messages, GitHub
  activity and status; imported/backdated/future-dated evidence appears. The
  merged SDP schema, conversely, accepts release events only.
- **RECOMMENDATION:** Replace both extremes with a versioned general envelope
  and a small high-value event vocabulary. Git/GitHub observations are derived.

### Review provenance by role label alone

- **OBSERVED:** Fresh review is substantive but normally shares the
  `Hans-Einar` GitHub identity, and some accepted state contradicts current
  findings.
- **RECOMMENDATION:** Deprecate treating `actor: reviewer`, owner comments,
  Master self-check or green CI as proof of independent review. Require an
  attestation/session identity, exact head and current disposition.

### Raw chat prompt/response as primary authority

- **OBSERVED:** Some Steering interaction experiments preserve exact prompts
  and raw responses but still lack compact Issue/branch/PR state.
- **RECOMMENDATION:** Keep raw interaction logs optional for audit/safety. The
  Issue and authored decisions are authoritative; do not require copying all
  chat transport or sensitive material into every repository.

### Version-bound Fix identity

- **OBSERVED:** Existing `FIX-<target-version>-NNN` couples work identity to a
  release target that can change; concrete Fix execution is not established on
  default in this corpus.
- **RECOMMENDATION:** Use stable `FIX-NNN` prospectively and relate it to one or
  more target Releases. Preserve legacy IDs unchanged.

## Deprecate after migration support exists

### Duplicate Toolkit and project-owned authorities

- **OBSERVED:** SDP, HEOS and weight_app_flutter show duplicate root/Toolkit
  assets or populated numbered documents beside empty installer templates.
- **RECOMMENDATION:** Converge on one versioned Toolkit-managed asset tree and
  explicit project-owned records. Keep tested forwarders/aliases only through a
  defined compatibility period.

### Unversioned project-local schemas/status vocabularies

- **OBSERVED:** ActionCam, HSX, map_tracker and other projects invent useful but
  incompatible assignments, ledgers, IDs and statuses.
- **RECOMMENDATION:** Add capability/profile declarations and canonical schemas;
  preserve unknown fields and project extensions, but do not guess their
  semantics.

### Historical “active” prose

- **OBSERVED:** Completed Sprint/Refactor documents often retain local “active”
  wording that is accurate as a historical snapshot and misleading as live
  state.
- **RECOMMENDATION:** Treat historical records as immutable snapshots and render
  current state from current pointers/GitHub. Migration may add a generated
  banner/index; it must not rewrite history silently.

## Keep project-local, not canonical

- ActionCam device selectors, temperatures, MTD/storage/ADB and physical media
  claim taxonomy;
- HSX debugger domain/ABI/address/epoch and anti-monolith decomposition;
- weight_app_flutter Bluetooth, HEOS device, tractor UI and hardware thresholds;
- Lyndata QGIS/Frost credentials, UALF and scale thresholds;
- GrassPhenology weather/Sentinel/OpenLayers/rendered workflow rules;
- gh-sdp archive/path/link/locking/journaling mechanics;
- repository-specific test commands, fixtures and platform matrices.

Canonical SDP requires their authority, evidence, verification, recovery and
stop shape; it does not standardize their domain contents.

## Removal criteria

A deprecated contract may be removed from the Toolkit only after:

1. a replacement schema/template/instruction is released;
2. Analyzer can identify and read the legacy profile;
3. gh-sdp can preview an additive migration with no project-owned overwrite;
4. representative legacy/default/non-default fixtures pass;
5. the deprecation window and rollback guidance are documented;
6. at least two real project pilots complete without loss of authority/history;
7. Steering explicitly authorizes removal in a later implementation Issue.
