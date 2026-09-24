# P04 — synthesis and further investigation

**Date:** 2026-09-10. **Status:** research/recommendations, no language decision. **Assignment:** [#10/P04](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5618871491). **Basis:** [32 profiles and methodology](existingDesignLanguages/README.md).

## Owner intent

SDP should expose Features, responsibilities and pathways through change. Architects/Designers propose bounded design deltas before Workers change code. Owners/Verifiers see approvals, changes and verified assertions. Language selection, history formats and three proposed Feature classes remain hypotheses.

SharedUI drift/large refactors are **owner reports from P04**, not newly verified GrassPhenology, TerrainAnalyzer or HSX findings. Ponsse UI is a possible case from [discussion #72](https://github.com/Hans-Einar/ponsse/issues/72), not verified current implementation.

## Survey findings

Existing building blocks cover much of the idea. A new text format may be easier than precise semantics, identity, code checking and understandable change processes.

| Need | Observed mechanisms | SDP still needs |
| --- | --- | --- |
| Multiple views/one model | Structurizr, LikeC4, UML/SysML | FEAT selections, baseline identity and explicit omissions |
| Responsibility/concept boundaries | Context Mapper contexts/aggregates/refactoring | Stable FEAT identity and actual responsibility verification |
| Product variants | UVL, Clafer | Distinguish configuration choices, persistent capabilities and revisions |
| Explicit contracts | OpenAPI, AsyncAPI, Protobuf, Smithy | Consumer impact, units, delivery/errors and releases |
| Machine-checkable boundaries | OCL, CUE, Rego, formal analysis | Reliable observed code facts and approved-exception process |
| Model transformations | ATL, QVT, Epsilon | Preservation, identity mapping and loss checks |
| History/provenance | Edapt change models, PROV-O revision/derivation | Approval, responsibility migration and product evidence |
| Requirements interchange | ReqIF | Stakeholder/UseCase→REQ→FEAT→verification→release |

This is a capability map, not ranking; profiles provide rationale/primary sources.

1. [Context Mapper](existingDesignLanguages/context-mapper-cml.md) already splits bounded contexts by use cases/stories and merges contexts. Useful precedent for responsibility changes; valid resulting CML does not prove preserved program behavior.
2. [Edapt](existingDesignLanguages/edapt.md) provides explicit operation history/model migration; [Epsilon](existingDesignLanguages/epsilon.md) adds matching/transformation/validation. Not all transition infrastructure needs reinvention.
3. [Structurizr](existingDesignLanguages/structurizr-dsl.md) and [LikeC4](existingDesignLanguages/likec4.md) demonstrate practical textual model→viewpoints. No mechanism alone proves complete FEAT pathways/release readiness.

## Four kinds of change

| Kind | Example | Not equivalent to |
| --- | --- | --- |
| Runtime transition | Action changes UI/domain state | Architecture evolution across commits |
| Product variant | Select WebUI/DesktopUI | Before/after revision |
| Model/schema evolution | Split type/responsibility | Correct product migration evidence |
| Delivered product change | User receives changed behavior in identified build | Closed issue/PR |

SCXML history and Alloy/TLA+ temporal support solve different problems from Git. Protobuf field identity supports contract evolution; PROV-O expresses derivation. Combine by need rather than call everything history.

**Recommended hypothesis:** stable design IDs + explicit model revisions + separate transition descriptions. Git stores data; semantics explain preserved/moved/split/merged/replaced/unmapped responsibility. Avoid embedding all history in every Feature.

Separate approved baselines, proposals and observed implementations. Blueprints identify model revisions, generator versions and selection rules. Passing blueprint checks are insufficient if code inventories omit dependencies.

## Challenge the Feature/Refactor taxonomy

**Recommendation:** trial independent dimensions before a type hierarchy:

- Capability/Feature: what the system does, for whom, with REQ/acceptance.
- Design object: responsibility, layer, contract, implementation binding or decision realizing capabilities.
- Change object: why/how changes occur, affected objects and preservation requirements.
- Classification: domain/design/architecture tags may overlap in one change.

Distributability may be quality/capability; sockets may be architecture decisions with contract effects. New renderers may satisfy new needs or merely replace implementation. Intent/observable behavior decides, not role names.

Refactors require preservation criteria; simultaneous behavior changes must be explicit. Architecture→design→domain can guide planning, but implementation follows dependencies/migration risk. Small vertical transitions may teach more than three large rewrites.

Domain is not one-to-one with processes. Distinguish conceptual domain, bounded context, logical component/Container and deployed process instance. Namespaces/packages/folders bind implementations. Layers can contain multiple responsibilities; layering alone does not prevent large tangled files.

Interfaces/composition need not imply classes/inheritance. [Smithy mixins](existingDesignLanguages/smithy.md) illustrate model reuse, not a case for runtime mixins. Specify contracts/dependencies rather than universal programming paradigms.

## Contracts, reuse and release traceability

Separate contract shape from meaning: types/schemas, units, command/event direction, ordering, delivery, failures, timeouts, idempotency and compatibility. UI needs representation ownership/lifetime, multiple presentations per instance, subscribe/unsubscribe and unhandled-event rules. Hints/localization differ from domain identity. Test renderers check messaging but do not replace all real web/desktop adapter tests.

Library bindings should identify packages, permitted versions/contracts, implementation language, source areas and reuse/adapt/build-new decisions. Model fields alone cannot prevent SharedUI drift: independently inspect imports/duplication. Go is an owner preference, not universal SDP policy. Justify reuse by actual needs/contracts, not speculative abstractions.

Traceability must reach actual deliveries. A **future gh-sdp hypothesis** is a release manifest with stable FEAT ID, revision/variant, user-facing changes, contract versions, verification references, commit and build artifact. PRs/issues support it; they do not own release-note truth.

Handle partial Features, superseded work, backports, multi-release Features and implemented-but-unexposed code. Reviews/internal refactors need traceability without automatic user-facing release notes. P04 implements no gh-sdp feature.

## Significant open questions

1. Minimum authoritative objects and identities surviving rename/split/merge?
2. Precise preserved-Feature behavior, including quality?
3. Pathway as allowed route, runtime scenario or distinct types for both?
4. Transition states, compatibility windows, rollback and partial mappings?
5. Independent observed code facts and visible uncertainty?
6. Which checks are syntax, semantics, static analysis, runtime or substantive review?
7. Lightweight small projects and bounded model/view maintenance cost?
8. Change authority and Worker stop conditions instead of changing models to pass?

Models should bound responsibilities/contracts, not copy every class/line. Architect/Designer/Worker skills and independent verification share identified assignments. Rework budgets, stops and escalation belong to process; DSLs can carry them but cannot alone prevent agent loops.

## Next bounded assignment — proposed, not started

**Purpose:** compare one small change representation without selecting/freezing language.

**Input:** separately confirmed current MVP1 UI and one #72 delta, with one output capability, one input action, one contract and at most six responsibility units. If current design is unverified, use a labeled hypothetical case.

**Tracks:** A: Structurizr/LikeC4-style architecture DSL; B: CML/SysML semantic model; C: small JSON/CUE model with referenced contracts. Select concrete tools/versions/terms first. These are experiments, not an approved shortlist winner.

**Delivery:** same baseline/delta in all tracks; FEAT→REQ→responsibility→contract view, before/after mappings and Worker/Verifier selections. Identify native concepts, metadata conventions and custom code. Trial rename, split/move, broken references and forbidden dependencies, including undetected failures. Compare edits, information loss, owner readability and integration cost.

**Limits:** no product code, full Concept1→MVP1 map, general compiler or skill. At most one repair round per track; record unresolved tool issues rather than endless Worker/Verifier loops. Stop after comparison and owner review.

Later, Concept1→MVP1 can stress-test explicit partial equivalence. MVP1 success does not prove universal suitability; eventually trial a small project and different contract/deployment pattern.
