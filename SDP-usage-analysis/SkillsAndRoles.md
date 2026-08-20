# Skills and role recommendations

## Decision

- **RECOMMENDATION:** Maintain two separate ChatGPT skills:
  `steering-group` for repository governance and `sdp` for reusable method
  knowledge/operations.
- **RECOMMENDATION:** Update existing Codex SDP role skills for Issue,
  Feature/Refactor/Fix, exact GitHub binding and declared/observed/accepted
  state. Do not create a new skill for every document or duplicate the Reviewer
  workflow under `sdp-code-review`.
- **RECOMMENDATION:** Repository `AGENTS-project.md` remains the authority for
  product/domain commands, platform constraints, safety and local exceptions.

## ChatGPT `steering-group` skill

Purpose: govern one repository or a coordinated repository set above bounded
Codex Issue assignments.

Responsibilities:

- discover repository/default/current evidence and active SDP state;
- inspect product direction, architecture and cross-repository dependencies;
- decide the next bounded Feature/Refactor/Fix/Study/release assignment;
- create/refine high-quality GitHub Issues and explicit comment amendments;
- inspect the real draft PR, exact head, CI, verification, independent review,
  traceability and limitations returned by the Issue Master;
- accept, require changes, block, split, cancel or open follow-up Issues;
- coordinate release readiness and dependency order;
- stop at owner/product/architecture decision boundaries.

It does not:

- implement product code;
- treat chat memory as repository truth;
- accept a Master summary without evidence inspection;
- silently change project state or expand the Issue;
- replace the human owner's authority.

Durable outputs are Issues/comments, repository decisions and inspected GitHub
identities. Raw conversation can be retained when required for audit/safety but
is not the primary assignment contract.

## ChatGPT `sdp` skill

Purpose: provide method knowledge and reusable, repository-neutral SDP
operations without governing a product.

Responsibilities:

- explain the current SDP model and evidence vocabulary;
- inspect/bootstrap/profile an SDP installation;
- identify authoritative versus template/legacy/project-local records;
- help author Feature/Refactor/Fix, Study, Slice and Issue contracts;
- explain stable IDs, status projections, relations and migration;
- check conceptual readiness, Handoff and evidence completeness;
- surface contradictions and suggest proportionate next records;
- cooperate with the Steering skill and downstream Codex/gh-sdp/Analyzer tools.

It does not:

- select product priority, architecture direction or release scope;
- accept/reject a Codex result;
- authorize implementation;
- infer missing evidence or silently migrate files.

## Why the skills remain separate

- **OBSERVED:** weight_app_flutter's Steering review performs product/
  architecture disposition and separates Refactor from later Features, while
  reusable method knowledge is needed by every repository role.
- **INFERENCE:** Combining them would let generic method assistance accidentally
  assume governance authority. Separating them also allows `sdp` to serve users
  who need explanation/inspection without booting a supervisory conversation.

## Codex role skill changes

### `sdp-master`

Add:

- Issue/comment authority and exact baseline reconciliation;
- Feature/Refactor/Fix primary owner and current Slice;
- early draft PR and declared/observed/accepted state checks;
- bounded multi-Worker/Architect/Verifier/Reviewer delegation;
- current-head review invalidation and Issue stop/split policy;
- compact assignment/Handoff and GitHub contradiction handling.

Change the stop rule from only Slice/Fix boundary to the assigned Issue boundary;
the Master may execute several explicitly authorized sequential Slices but may
not start another Issue.

### `sdp-worker`

Add Issue/work-owner/Slice identity, exact owned areas/shared touchpoints,
concurrency rules and discovery escalation. Preserve “implement one bounded
contract and stop.”

### `sdp-reviewer`

Require:

- exact base/head/tree and current PR/check observations;
- reviewer/session attestation independent of implementation;
- validation of Feature/Refactor/Issue scope and architecture deltas;
- review of declared/observed/accepted contradictions;
- mandatory changed-head re-review after material fixes;
- explicit Blocking/High/Medium/Low/Note disposition.

### `sdp-verifier`

Add qualifying-evidence semantics: subject, candidate, environment, command,
artifact/run, outcome and limitation. It must distinguish recorded evidence from
rerun evidence and external/manual gates from automated checks.

### `sdp-architect`

Add the architecture-revision threshold, Feature/Refactor-local design deltas,
parallel Study convergence/refreeze and Issue/Steering escalation. Preserve its
non-implementation boundary.

### `sdp-traceability`

Replace blanket registry/event authoring with:

- compact current pointers;
- non-derivable semantic relations;
- material general events under a versioned envelope;
- generated Git/GitHub facts and contradiction detection;
- declared/observed/accepted status separation;
- compatibility-profile-aware preservation.

### `sdp-vertical-refactor`

Bind every Refactor to a stable `REF-NNN` and Issue assignment, require baseline
and compatibility, and express ordered vertical Slices plus explicit temporary
adapter removal conditions. Preserve behavior-first evidence from
GrassPhenology and weight_app_flutter without copying their domain mechanics.

### Release/versioning/auditor skills

- Keep release identity separate from Feature/Refactor/Slice.
- Add Release inclusion relations for accepted work.
- Auditor validates GitHub contradictions, current-head evidence, profiles and
  stale assignment state.
- Preserve truthful two-phase publication and existing SemVer rules.

## Missing or deferred Codex skills

- **RECOMMENDATION:** Do not add `sdp-code-review`; `sdp-reviewer` owns it.
- **RECOMMENDATION:** Do not add a generic `sdp-planner` until a distinct
  repeatable workflow exists beyond Master/Architect/Steering responsibilities.
- **RECOMMENDATION:** Do not add a separate Codex `steering-group` role. The
  supervisory ChatGPT skill operates above repository-local Codex assignments;
  Codex Master records and responds to Steering decisions.
- **RECOMMENDATION:** Consider a narrowly executable `sdp-assignment` helper
  only after the CurrentAssignment/Issue schema is accepted and multiple tools
  need the same create/validate operation. Until then, integrate it into Master,
  traceability and gh-sdp rather than inventing another prose skill.

## Skill packaging and versioning

- One canonical managed skill tree under `Toolkit/skills/`.
- Stable skill IDs and independent SemVer remain.
- Feature/Issue changes are compatibility-significant and must update capability
  metadata, manifests, fixtures and installed-version rules.
- Project-local instructions extend canonical roles without changing their
  evidence, independence, ownership or publication invariants.
- A Toolkit release and migration plan are required before calling the new
  skill behavior canonical for consuming repositories.
