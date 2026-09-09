# SPR-SDP-007 — Provisional vNext workflow contract

Status: active

GitHub issue: `Hans-Einar/SDP#7`

Authoritative baseline: `2cb49c02145621b099c47d05786716598e414e75` (`main`)

Delivery branch: `codex/issue-7-provisional-vnext-pilot`

## Goal

Produce the explicitly experimental SDP vNext workflow and scoped work-domain
contract required by Issue #7, validate it deterministically, obtain fresh
independent review, and stop with a concrete first-pilot handoff for HSX.

## Scope

- Define provisional Feature, Refactor, Fix, Study, Slice, optional Sprint and
  optional Iteration semantics with explicit cardinalities and state meaning.
- Define reusable scoped work-domain namespaces, stable identity, cross-scope
  relations, collision prevention, concurrent Issue Master coordination and
  future repository-split identity preservation.
- Provide pilot-only records/templates, an Issue contract, a plural per-Issue
  assignment model and proportional profiles.
- Provide deterministic validation for the experimental corpus.
- Produce an evidence-specific HSX pilot handoff without changing HSX.
- Record exact verification and fresh review evidence for the final candidate.

## Invariants

- Issue #7 and its comments are the durable assignment contract.
- Issue #5's Steering-accepted study at `main@2cb49c...` is the evidence basis.
- All vNext artifacts remain explicitly provisional and isolated from canonical
  Toolkit contracts, schemas, templates and release state.
- HSX, gh-sdp, SDP-Analyzer and PR #4 are read-only/out of scope.
- Existing historical identifiers are never rewritten or silently normalized.
- Mutable GitHub facts are observed/derived rather than copied as authored
  truth; stable authority URLs and branch/base contracts may be authored.
- Product/deliverable work is delegated to fresh Workers. Review uses a fresh
  independent Reviewer against an exact candidate.
- No unresolved Blocking, High or Medium finding may remain at completion.

## Work decomposition

- `ITR-SDP-007-001` / `SLC-SDP-007-001`: activation, evidence-bound contract
  corpus, record examples/templates, deterministic validation and HSX handoff.
- `ITR-SDP-007-002` / `SLC-SDP-007-002`: exact-candidate verification,
  independent adversarial review and required rework.
- `ITR-SDP-007-003` / `SLC-SDP-007-003`: closure evidence, traceability,
  draft-PR handoff and pilot-ready disposition.

## Non-goals

- Canonical vNext Toolkit implementation or schema freeze.
- HSX migration or any change to another repository.
- PR #4 modification, merge or contract integration.
- gh-sdp or SDP-Analyzer implementation.
- Toolkit version/release change, tag, GitHub Release or publication.

## Completion signal

`PROVISIONAL_CONTRACT_READY_FOR_PILOT` may be declared only when every Issue #7
deliverable exists, deterministic validation passes, a fresh exact-candidate
review has no unresolved Blocking/High/Medium finding, a draft PR is open and
the HSX pilot handoff names the operational evidence state and first bounded
pilot contract.
