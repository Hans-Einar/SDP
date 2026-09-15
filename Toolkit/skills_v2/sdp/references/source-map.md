# Source map and adoption decisions

Inspected 2026-09-15. This is a bounded provenance map, not an assertion that
every historical comment was reviewed. Links identify source material; the
installed consuming-project contract remains the operational authority.

## Working baseline

PR #4 branch at `d611b8bf72aeb30d86c5ef28902469462b06a803` is the draft's
baseline. All 16 Markdown documents immediately under `docs/` were read,
alongside root README, document guide, DraftStandardDocumentProcedure,
TieredDesignAndImplementation and Instructions/README. Relevant Issue bodies
and issue-thread comments reviewed: #5, #7, #9, #10 and PR #4/#8. This does not
claim an exhaustive review of PR inline comments or all root lifecycle records.

| Source | Status and useful content | Candidate treatment |
|---|---|---|
| [How SDP Works](https://github.com/Hans-Einar/SDP/blob/d611b8bf72aeb30d86c5ef28902469462b06a803/docs/How-SDP-Works.md) | Baseline lifecycle, living design and existing work hierarchy | Explicit document read/update obligations; preserve installed conventions |
| [Installation Contract](https://github.com/Hans-Einar/SDP/blob/d611b8bf72aeb30d86c5ef28902469462b06a803/docs/Installation-Contract.md) | Explicit inventory, ownership and neutral seeds | No implicit skill installation or copying repository live state |
| [Validation](https://github.com/Hans-Einar/SDP/blob/d611b8bf72aeb30d86c5ef28902469462b06a803/docs/Validation.md) | Structural validation does not prove Markdown meaning or acceptance | Add substantive documentary review; retain schema checks |
| [Feature Governance proposal](https://github.com/Hans-Einar/SDP/blob/d611b8bf72aeb30d86c5ef28902469462b06a803/docs/Feature-Governance-And-SDP-2.0.md) | SDP-PROP-001, additive 1.x proposal plus future 2.0 design | Carry forward intent/integration analysis; do not claim proposed schemas are installed |
| [Tier proposal](https://github.com/Hans-Einar/SDP/blob/d611b8bf72aeb30d86c5ef28902469462b06a803/TieredDesignAndImplementation.md) | Working draft: vertical capability across layers | Clarify vocabulary; defer canonical adoption and renaming |
| [Issue #5](https://github.com/Hans-Einar/SDP/issues/5) and [PR #6](https://github.com/Hans-Einar/SDP/pull/6) | Accepted cross-project study and direction | Preserve durable intent, bounded work and living foundation |
| [Issue #7](https://github.com/Hans-Einar/SDP/issues/7) and [PR #8](https://github.com/Hans-Einar/SDP/pull/8) | Experimental contract with unresolved rework | Borrow principles; no schema/identity migration |
| [Issue #9](https://github.com/Hans-Einar/SDP/issues/9) | vNext handoff; Feature contract and integrated acceptance questions | Distinguish Slice evidence from capability acceptance without new role hierarchy |
| [Issue #10 P03](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5618854397) | Owner discussion of durable Features, Use Cases, requirements and domain ownership | Preserve rationale and cross-layer impact; no mandatory new object graph |
| [Issue #10 P04](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5618871491) | Hypothesis about Domain/Design/Architecture Features requiring study | Not an accepted UserFeature/SystemFeature taxonomy; defer to vNext |

## Conflicts and deliberate limits

- The root document guide describes root folders as reference templates, while
  the PR #4 installation contract distinguishes live repository records from
  neutral `Toolkit/project-templates`. Use the installation contract for
  distribution; do not repeat the older description in new installation rules.
- Accepted study direction makes Sprint optional and Iteration nonmandatory;
  the installed vNow contracts still use them. A study's acceptance does not
  constitute a downstream installation/migration. Keep existing obligations
  until a scoped compatibility change adopts the new model.
- The pilot's SkillsAndRoles recommends a method `sdp` skill distinct from
  supervisory ChatGPT Steering, and avoids a new Codex steering-group role.
  This draft's `sdp-steering` is an evaluation candidate for explicitly scoped
  governance assistance across hosts, not an accepted historical role change.
  The router cannot appoint itself owner or provide a human disposition.
- No accepted precise Feature-versus-Functionality or user/system entity
  taxonomy was established by this bounded review. The document-workflow guide
  uses behavior descriptions and actor views without new canonical types.
- PR #11's system-description-language implementation direction remains
  outside this vNow draft. Reading its originating questions does not adopt it.

## Promotion is separate work

The shared document guide and role links are a bundle dependency. Promotion
must list every required reference and skill explicitly in the installation
manifest; v1 does not expand directories automatically. Reconcile native
discovery, generated AGENTS paths, metadata, validators and compatibility
fixtures together. Test installation and actual skill selection in a fresh
target-host session. This candidate collection has not performed that work.
