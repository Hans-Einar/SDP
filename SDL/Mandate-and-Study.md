# Working document: model-based, controlled software development with agents

**Document ID:** SDL-MANDATE-STUDY-001  
**Revision:** 0.1, 2026-09-14  
**Status:** Edited owner intent with an initial study and explicitly labeled proposals. Not an adopted language, schema or implementation plan.  
**Context:** Hans-Einar/SDP, issues #5, #7 / PR #8, #9 and #10.  
**Audience:** Owner, architects, designers and agents investigating, developing or using SDP.

## 0. Read this first

We want a standardized development model for directing agent-driven software development through compact, machine-readable descriptions of system purpose, structure, responsibilities, contracts and behavior. Tools should derive task-specific context, change descriptions and checks. Major changes should follow explicit revisions and phase transitions while allowing learning and further development.

Local code changes often lack an understood connection from user needs to actual effects across the system. Repeated local repairs can conceal missing integration, move responsibilities and gradually violate architecture. The model, process and checks should make these deviations visible and manageable.

We want detail where it directs development, not a manually maintained copy of the whole implementation. Machine-generated answers must identify covered properties, code areas and execution variants. Unknown or unchecked conformance must remain visible.

**Reading guide:** Part I restates the mandate. Part II provides an initial study and implementation proposals. Part III lists questions, possible investigations and continuation rules. Agents must not treat proposals as approved SDP contracts.

Shortcuts: [Mandate](#part-i--mandate) · [Phases and learning](#4-phases-revisions-and-learning) · [Change model](#11-proposed-phase-and-change-model) · [Implementation](#13-from-design-description-to-implementation) · [Standards](#15-standards-for-design-lifecycle-and-implementation) · [Verification limits](#17-what-can-be-verified-programmatically) · [Further work](#part-iii--further-work)

## 1. Provenance, status and existing SDP

### 1.1 Sources

Based on owner input from 14 September 2026, including the quoted Gemini dialogue, and existing SDP work. This does not restart the project or revoke earlier decisions.

| Track | Status checked 2026-09-14 | Meaning here |
| --- | --- | --- |
| [Issue #5](https://github.com/Hans-Einar/SDP/issues/5) / PR #6 | Study accepted and merged; #5 closed. Accepted merge: `2cb49c02145621b099c47d05786716598e414e75`. | Preserve evidence and accepted overall direction. |
| [Issue #7](https://github.com/Hans-Einar/SDP/issues/7) / [PR #8](https://github.com/Hans-Einar/SDP/pull/8) | Open preliminary pilot. Draft PR #8 at `ea9fcf1cdd3198aeac515b89b55398282c463838`. Latest status says rework lacks renewed independent verification/review. | Does not declare pilot contracts complete or findings closed. |
| [Issue #9](https://github.com/Hans-Einar/SDP/issues/9) | Open Steering/coordination track. | Home for overall direction and coordination. |
| [Issue #10](https://github.com/Hans-Einar/SDP/issues/10) / [PR #11](https://github.com/Hans-Einar/SDP/pull/11) | Model input/research. Draft PR #11 at `319ee2a43fe8a05bc3aac0860ff937bba57819ee`. | Existing 32-language survey is evidence, not language selection. |

Acceptance/status sources: [#5 Steering acceptance](https://github.com/Hans-Einar/SDP/issues/5#issuecomment-5368329242), [#7 WIP status](https://github.com/Hans-Einar/SDP/issues/7#issuecomment-5607093722), [#10 R04 delivery](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5622626432). These are issue/PR status checks, not renewed review of #7 code or every finding. Read #10's original stop instruction with later authorization and R04; it alone does not establish delivery status.

Existing foundations:

- [Accepted study direction](../SDP/Studies/UsageAnalysis/ProposedSDPWorkflow.md).
- [Catalogue of 32 languages](research/existingDesignLanguages/README.md).
- [P04 synthesis, limits and proposed comparison](research/README.md).

### 1.2 Interpreting statements

| Label / location | Meaning |
| --- | --- |
| Mandate, Part I | Restated owner intent: development goals, not existing tool capabilities. |
| Sourced observation | What cited sources describe within the stated investigation scope. |
| Proposal/assessment, Part II | Assistant proposals requiring evaluation before becoming normative SDP. |
| Open question, Part III | Undecided matter. |

This revision documents the assignment and initial study. It changes no canonical Toolkit contract and approves no project migration. Coordinate Feature/Functionality/Channel concepts and phase rules with #7/#9 before potential adoption.

# Part I — Mandate

## 2. Purpose and problem

### 2.1 Purpose

Develop a reusable method and associated design-language/tool foundation for projects where agents write and maintain substantial code. A new agent should quickly understand the system at relevant granularity and receive precise change boundaries without reading extensive history.

Allow architecture, design and implementation to evolve while keeping changes explicit, traceable and checkable. Rigor primarily concerns responsibilities, current decision evidence, change authority and progression criteria.

### 2.2 Failure pattern to counter

1. An agent changes a low-level subsystem and assumes the desired function now works end to end.
2. Necessary changes to other components, contracts, configuration or user interfaces are missed.
3. Tests/use expose gaps; local repairs proceed without reconsidering the system.
4. Repairs move responsibility, duplicate functionality or violate boundaries.
5. The system eventually appears functional while implementation substantially diverges from intended architecture/design.
6. Later changes become harder; the owner loses track of what was built and why.

This is the owner's reported experience/motivation, not a new empirical audit of affected products or a proven explanation of specific defects.

### 2.3 Earlier SDP direction

The desire for less “agile” thinking means less ad hoc prioritization and local improvisation without a system basis. Learning, small vertical deliveries and revisiting decisions remain possible.

Continue #5's direction: Feature, Refactor or genuine Fix owns intent; an issue bounds an assignment; a Slice provides a coherent, checkable result. Sprint organization is optional. A phase model must accommodate this without letting early implementation plans dictate all future development.

## 3. Desired system/design language

### 3.1 Machine-readable model

Express system descriptions as versionable design code with defined syntax/semantics. An interpreter/model compiler should load, validate and derive information; it need not execute product functionality.

Investigate using an existing language, extending one or defining our own. Language and implementation technology remain open. Established standards/methods should help make concepts and processes precise.

### 3.2 Structure at multiple levels

Describe system context, external actors and boundaries; architecture components with responsibilities, contents and relations; *design constituents* with information needed to guide implementation; and bindings to source trees, code areas and actual implementations.

Detail may vary by subsystem and phase. Do not automatically require a manually modeled object for every class, method or code line.

### 3.3 First-class cross-cutting objects

| Object | Owner intent | Remaining clarification |
| --- | --- | --- |
| Feature | Expose persistent capability and its path through the system. | Relationship to user needs, REQ, quality requirements and existing SDP Feature. |
| Functionality | Describe realization and connections of functionality across the system. | Boundary against Feature, function, scenario, pathway, sequence and state machine. |
| Channel | Describe communication between independent units, with protocol/contract. | Logical/physical distinction, endpoints, direction, versioning and multiple participants. |

First-class objects can be identified, referenced, analyzed and explicitly changed. They connect the structural model without requiring isolated subtrees or duplicate components. These starting points are not an exhaustive/adopted type hierarchy. Functionality especially needs a definition adding value beyond Feature.

## 4. Phases, revisions and learning

### 4.1 Phase-based development

Support phases such as Concept/MVP with named iterations Concept1, Concept2 and MVP1. Link architecture, design and implementation to iteration scope/maturity.

Major architecture/design changes require explicit phase transitions or new iterations within a phase. Define “major” and handling of smaller changes. Each iteration needs an identifiable source tree/design basis. Git/directory organization remains open: tree traceability does not necessarily require a new repository, permanent branch or rewrite.

### 4.2 Example paths

```text
Concept1
  |-- new architecture/design decisions ------> Concept2
  |-- different implementation, same design --> Concept2
  `-- sufficient learning --------------------> MVP1
          |
          `-- preserve lessons/change needs before designing MVP1
```

Concept2 is optional. A new iteration may retain design and replace implementation, or change architecture, design and code. Identify what carries forward.

### 4.3 Preserve learning without deciding the next design

Concept work prototypes and learns. Experience may suggest new responsibility boundaries, contracts or technologies. Preserve observations, rationale, uncertainty and desired effects structurally.

Do not necessarily design MVP1 architecture during Concept1. Express future desired outcomes without presenting them as approved designs. Issues/comments provide input/references but should not alone carry consolidated current learning.

At MVP1 planning, use relevant learning to produce a coherent architecture/design for agreed scope. Reconcile completeness for that scope with avoiding detailed models of all future implementation.

## 5. Derived views and agent context

Tools should provide multiple-level overviews and bounded selections for agent instructions: system structure, component contents, Feature paths, Functionality behavior and Channel contracts.

Selections must preserve responsibilities, external touchpoints, constraints and completion criteria. Brevity must not hide dependencies essential to success. Owners need understandable explanations of proposed changes, reasons and actual checks. Owner/agent views derive from the same identified model basis.

## 6. Change descriptions and impact analysis

### 6.1 Design before implementation

Describe intended changes in design code first. Compare an identified baseline with a proposed target and explain necessary work. Defects/new learning may revise proposals through a defined process.

Keep approved design, proposed design and observed system facts distinct. Approved designs may differ from code; code may differ from deployed systems.

### 6.2 Semantic and residual diff

Explain changed meaning, responsibilities, contracts and connections, not merely lines. Show whole-system diffs and filters by Feature, Functionality, Channel or other objects. Residual diffs expose changes unexplained by selected objects, including unlinked/unjustified work. One change may matter to multiple objects.

### 6.3 System-wide consequences

Find affected elements and changes/checks necessary to expose desired effects end to end. Support reconsidering architecture where dependencies are strong or widespread. Expose missing information: an absent modeled dependency does not prove no dependency exists.

## 7. Conformance to actual software

Investigate programmatic checks of agreed design properties through source annotations (possibly Doxygen), code/type/compiler-API analysis, intermediate representations such as LLVM IR where useful, and tests/observations of behavior/integration.

Check comments/annotations against referenced code. Agreement between descriptions alone does not prove functionality. Compiling some/all design language to LLVM IR is exploratory, not a first-solution prerequisite; establish what comparison/proof it enables.

## 8. Measurable goals

| ID | Desired outcome | Future pilot demonstration |
| --- | --- | --- |
| M01 | Rapid understanding | New agent finds current responsibilities, contracts and Feature paths in compact context. |
| M02 | Controlled change | Before/after model and rationale exist before design-changing code work. |
| M03 | Preserved learning | Concept findings inform MVP planning without predetermined MVP design. |
| M04 | Visible whole | Channel changes reveal producers, consumers and affected Feature paths. |
| M05 | Complete change accounting | Detected deltas appear without loss in total/object/residual diffs. |
| M06 | Reproducible conformance | Results identify property, source/build, evidence and coverage limits. |
| M07 | Controlled deviation | Agents cannot silently change expectations to excuse violations. |
| M08 | Bounded maintenance | Measure model editing, handwritten information and context size. |
| M09 | Reusable method | Lightweight profiles work without distributed-system complexity. |

Agree numerical time/context/maintenance thresholds before pilot measurement. This study does not claim these goals achieved.

# Part II — Initial study and proposals

## 9. Separate concerns before choosing a language

**Proposal:** distinguish method (decision authority, evidence and reconsideration), semantic model (objects, relations, rules, identity), representation/tools (editing, validation, comparison, presentation), and evidence (independent code/build/runtime investigation).

Rich languages may lack good change processes. Valid schemas may coexist with architecture violations. Good processes become expensive when information is manually copied across documents.

## 10. Proposed precise model concepts

### 10.1 Graph with different relations

Allow structural decomposition and cross-cutting relations. Folder trees cannot adequately represent shared services, multiple consumers and Features spanning components.

| Concept | Proposed meaning |
| --- | --- |
| Component | Responsibility-bearing architectural part at an identified level. |
| Design constituent | Modeled design part, not synonymous with class/file; study exact IEEE 1016 mapping. |
| Feature | Persistent capability with purpose, requirements and observable acceptance criteria. |
| Functionality | Provisional candidate: coherent behavioral responsibility realizing parts of Features; retain as a type only if useful in a pilot. |
| Scenario / witness | Observable concrete path for examining Feature assertions. |
| Pathway | Defined responsibility/transition chain; distinguish permitted structural routes from actual runtime paths. |
| Channel | Logical communication with identified participants/contract; bind transport/deployment separately. |
| Contract | Permitted interactions, data, meaning, errors and compatibility assumptions. |
| Implementation binding | Design-object link to package, symbol, code area, build variant or deployment unit. |
| Constraint / invariant | Expressed rule, such as forbidden dependencies or required current revision. |

Candidate relations: contains, realizes, uses contract, sends, receives, must preserve, verified by. Define direction/meaning; generic `depends_on` alone rarely explains change obligations.

Features may use multiple Functionalities and vice versa; Channels may participate in multiple paths. Domain, logical component, layer, process and source folder are not automatically one object.

### 10.2 Stable identity over time

**Proposal:** separate IDs from names, locations, revisions, phases and releases. Renames preserve identity when responsibility is unchanged. Split/merge needs explicit responsibility/preservation mappings, not name similarity.

Git stores revisions; models/change objects explain meaning. Trial separate transition descriptions before embedding all history in every object. Coordinate identity rules with #7's provisional scope/reservation/cross-repository rules.

### 10.3 Declared coverage instead of a digital twin

**Proposal:** manually describe intent-bearing information: responsibilities, contracts, important states/invariants, forbidden dependencies and acceptance scenarios. Obtain mechanical facts such as files/imports through tools.

Each model area identifies binding/descriptive/derived/out-of-scope information. Lightweight profiles may start with component boundaries, critical contracts and one Feature path. Add detail where risk/change needs justify it.

## 11. Proposed phase and change model

### 11.1 Separate coordinates

| Coordinate | Example | Identifies |
| --- | --- | --- |
| Lifecycle phase | Concept, MVP | Purpose and expected maturity. |
| Phase iteration | Concept1, Concept2, MVP1 | Named execution with scope and input basis. |
| Model revision | D17 | Particular design content. |
| Source revision | Git commit SHA | Concrete source tree. |
| Build / deployment | Artifact digest + configuration | Built or actually running system. |
| Acceptance | Design approved, Feature verified, owner accepted | Available decisions/evidence. |

Do not compress these into one version number. Concept1 may contain ordinary model/source revisions; a new iteration is a separate decision. MVP1 is not SemVer 1.0.0.

### 11.2 Learning register

**Proposed minimum:** stable ID, observation origin/content, evidence or explicit assumption, consequence, desired future property, alternatives, questions and subsequent decision/disposition.

Illustrative example, not a product finding:

```yaml
id: LEARN-EXAMPLE-01
origin: Concept1
observation: "UI and domain logic cannot be tested independently in this trial."
evidence: "Insert a reference to the actual trial."
desired_outcome: "Run the domain path without a UI renderer."
options:
  - "Explicit port between domain and presentation"
  - "Separate presentation process"
decision: deferred
candidate_target: MVP1
```

The entry commits us to handling the problem, not choosing architecture. During MVP1 planning, explicitly classify learning as incorporated into requirements/design, further investigation, deferred or rejected with reasons. Decision notes reference entries without overwriting observations.

### 11.3 Phase gates

**Proposed, not adopted:**

| Gate | Minimum basis | Outcomes |
| --- | --- | --- |
| Start iteration | Purpose, scope, learning, provisional boundaries, planned investigations. | Start Concept2/MVP1 or resolve missing basis. |
| Start bounded implementation | Coherent change design, explicit unknowns, contracts, consequences, acceptance scenarios. | Implement or run a bounded experiment first. |
| Substantial architecture/design change | Rationale, before/after responsibilities, capabilities, migration, new checks. | New iteration/phase, rejection or deferral. |
| Integrate result | Identified code/model and relevant structural, contract and integration checks. | Integrate within stated acceptance or return a bounded finding. |
| End iteration | Explicit disposition of outcomes, lessons, deviations and unfinished work. | New iteration, next phase or closure. |

Define major changes by protected boundaries: moved component responsibility, process/trust boundaries, broken public contracts, data ownership or Feature semantics. Changed-line count is a poor primary criterion.

Small fixes within approved contracts may be ordinary revisions. Exploratory spikes may use lighter design if hypotheses, limits and result disposition are explicit. No phase promises a final lifetime design.

## 12. Three bases and three comparisons

**B** is approved design baseline, **T** proposed target, **O** observed implementation facts for identified code/build/configuration. O has known coverage, not automatically complete reconstructed design.

| Comparison | Question |
| --- | --- |
| B → T | What design changes are proposed, and why? |
| O versus B | Which observations conform to/deviate from approved design? |
| O versus T | What remains or differs in target implementation? |

B→T may be semantic model diff; the others translate code facts into checkable properties/relations rather than ordinary document diff. Source observation is not deployment observation. Claims about running systems require identified artifacts/configuration/environment. Record preexisting deviations before changes to avoid misattribution.

### 12.1 Semantic change accounting

Each detected change has one identity and may appear in multiple views: addition, removal, moved responsibility, changed contract, split, merge, changed invariant. Automatic matching must expose uncertain identity mappings rather than assume common intent.

```text
D_total = all recorded semantic changes between B and T
D_object = union of changes explained by selected object views
D_rest = D_total minus D_object
```

Views may overlap; one contract change affecting Feature/Functionality/Channel remains one change. Explain residual infrastructure changes and missing links. Accounting covers detected model deltas only; independent code inventory must find code changes without model/assignment links.

### 12.2 Impact exceeds graph adjacency

Distinguish directly changed objects, potentially affected objects through defined relations, derived assessment/adaptation/test obligations, and concrete decisions to change/retain/unresolved.

Not every schema change requires every consumer to change, but relevant consumers need assessment. Impact claims should show their rules/relation chains. Dynamic registration, configuration and external services may require evidence beyond static imports. Fewer dependencies are not always better: assess responsibility, cohesion, performance, operations and migration cost.

## 13. From design description to implementation

**Proposed workflow, not prescribed by IEEE 1016:**

1. Confirm assignment/baseline: need, Feature revisions, B, source baseline, known deviations.
2. Investigate horizontally: follow responsibilities, contracts, consumers and paths; resolve unknown wiring/critical choices.
3. Describe T/rationale: preserve/change/move/migrate decisions and appropriate revision/gate.
4. Derive obligations/acceptance: end-to-end result, failure paths, checks and deviation criteria.
5. Plan vertical Slices by dependencies/risk; a compatible contract extension or adapter may come first without rewriting every layer.
6. Give agents identified context: contracts, permitted changes, affected neighbors, stop rules.
7. Implement/gather evidence: check code facts/runtime; return new design needs as proposals, not silent changes.
8. Integrate/check Feature paths: local Slice acceptance does not replace checks on actual integration baseline/configuration.
9. Record acceptance/remaining work: distinguish implemented, integrated, verified, owner-accepted and released; update current pointers through proper decisions.

### 13.1 Worker and Verifier context

**Proposed minimum:** assignment ID/goal, B/T revisions, source baseline, generator/rule versions, affected Feature/Functionality/Channel IDs, responsibilities/contracts, permitted code areas, shared touchpoints, preservation/check requirements, known deviations, unknowns and design-change authority.

Provide concise summaries with necessary detail links and visible omissions. Reject stale context after model/contract/integration changes. Worker and Verifier share contract evidence; verification must also find flaws/gaps in that evidence.

### 13.2 Prevent repair loops

Agree rework budgets and stop criteria: repeated failure classes without new explanation, crossing protected boundaries, new consumers or missing Feature effects despite local test passes.

On stopping, return observation, hypothesis and concrete analysis/decision needs. Models may be corrected, but changing verification expectations must remain visible. Tools must not normalize drift by overwriting approved designs with incidental code behavior.

## 14. Worked example: a Channel change that can cause tunnel vision

**Hypothetical**, not a verified description of Ponsse, HSX or other product code.

`FEAT-LENGTH` lets operators see current measured length. Candidate `FUNC-PRESENT-LENGTH` carries measurements from domain to presentation. `CH-MEASUREMENT` sends values to clients.

```text
Measurement source -> domain -> CH-MEASUREMENT -> client adapter -> presentation
                                     |
                                     `-> history storage
```

**B:** Channel values and history storage use millimeters; presentation shows centimeters.  
**T:** Channel sends meters; both consumers retain agreed properties, with unchanged external display.

A local serialization-only fix can retain JSON `number`, imports and passing format tests while breaking display/history.

| Derivation | Consequence / check |
| --- | --- |
| Semantic unit changes | Record contract break despite unchanged syntax/type. |
| Client adapter consumes Channel | Assess conversion/contract version. |
| History consumes Channel | Assess storage units, compatibility, existing data. |
| FEAT-LENGTH uses path | Require known physical-value scenario through actual display. |
| Old/new components may coexist | Explicit transition, compatibility window and possible rollback. |
| Deployment selects adapter | Verify correct registration and actual use. |

Witness: input representing 1 meter becomes 100 centimeters on screen and correct history at agreed precision. Supplement headless tests when they do not use actual client bindings.

Show the same change ID under Feature, Functionality and Channel. Explain deployment changes or expose them in residual diff. Independent inventory/observation must reveal omitted history consumers; graphs cannot infer unknown consumers.

Concept1 may first record inconsistent units as learning: “The next iteration needs explicit unit contracts.” This does not predetermine meters, millimeters or unit-bearing types in MVP1.

## 15. Standards for design, lifecycle and implementation

### 15.1 Sources and limits

Checked against public primary sources on 2026-09-14, mainly ISO/IEEE catalogues and summaries. Full normative texts were not reviewed clause by clause. We can state documented scope/proposed use, not standard conformance or a complete IEEE 1016 normative-reference list.

“Ecosystem” is a practical map of related needs, not a claim that all standards form one IEEE 1016 family or that 42010 replaces 1016.

### 15.2 Design and architecture

| Reference and observed status | Documented scope | Proposed SDP use |
| --- | --- | --- |
| [IEEE 1016-2009](https://standards.ieee.org/ieee/1016/4502/) — Inactive-Reserved; inactivated 2020-03-05 | Software Design Description content/organization at high and detailed levels; does not prescribe a design method/language. | Basis for model/derived SDD information; retain status with citation. |
| [ISO/IEC/IEEE 42010:2022](https://www.iso.org/standard/74393.html) — published | Requirements for architecture descriptions, frameworks, languages, viewpoints and model kinds; distinguishes architecture from its description. | Terminology/structure for concerns, model kinds and targeted views. |
| [ISO/IEC/IEEE 42020:2019](https://www.iso.org/standard/68982.html) — published; revision developing | Architecture governance, management and development through lifecycle. | Responsibility/decision processes for architecture design/revision. |
| [ISO/IEC/IEEE 42030:2019](https://www.iso.org/standard/73436.html) — published; marked for revision | Framework for organizing/documenting architecture evaluation. | Alternatives, quality and risk before major transitions. |

**1016 clarification:** still relevant to exploration, not an active, current universal solution. SDDs need not be large handwritten Markdown files; IEEE also describes other media/tool-based representations. [IEEE 1016-2009](https://standards.ieee.org/ieee/1016/4502/)

**Viewpoint versus view:** proposed working distinction: viewpoints define questions/conventions; views are concrete presentations for identified model evidence. Refine this against the selected edition. 42010 does not choose diagram tools/development methods. [ISO/IEC/IEEE 42010:2022](https://www.iso.org/standard/74393.html)

The Gemini draft's “IEEE 1016's 12 standard viewpoints” claim is not normative here. Verify exact lists, terminology and selection/adaptation requirements in full text before mapping. Do not assume every viewpoint must be materialized for every project.

### 15.3 Lifecycle, requirements, documentation and controlled change

| Reference and observed status | Documented scope | Proposed SDP use |
| --- | --- | --- |
| [ISO/IEC/IEEE 12207:2026](https://standards.ieee.org/ieee/12207/11416/) — Active; replaces 2017 | Software lifecycle process framework; processes may be concurrent, iterative and recursive. | Connect design, realization, integration, verification and evolution; select relevant processes per phase. |
| [ISO/IEC/IEEE 15288:2023](https://www.iso.org/standard/81702.html) — published | System lifecycle processes; prescribes no single lifecycle model/method. | Include hardware, operator work and larger systems. |
| [ISO/IEC/IEEE 24748-1:2024](https://www.iso.org/standard/84709.html) — published | Lifecycle management, models, stages and tailoring with 12207/15288. | Concept/MVP, entry/exit criteria and new iterations. |
| [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html) — published; marked for revision | Lifecycle requirements engineering. | Distinguish needs/normative requirements from solutions; link Features/verification. |
| [ISO/IEC/IEEE 15289:2019](https://www.iso.org/standard/74909.html) — published; confirmed 2025 | Lifecycle information-item purposes/content. | Identify required information and model-derived documents. |
| [ISO 10007:2017](https://www.iso.org/standard/70400.html) — published; marked for revision | Configuration management from concept to disposal. | Baselines, identified configurations, controlled change instead of unlabeled drift. |
| [ISO/IEC/IEEE 14764:2022](https://www.iso.org/standard/80710.html) — published | Software maintenance within lifecycle. | Organize change analysis, execution and maintenance. |

These sources do not standardize Concept1/Concept2/MVP1 names; SDP chooses its phases. Do not present 12207 processes as a mandatory once-only sequence. [ISO/IEC/IEEE 12207:2026](https://standards.ieee.org/ieee/12207/11416/)

Editions are not automatically harmonized: 14764:2022 explicitly references 12207:2017 maintenance. Check edition relationships before transferring process/clause references to 12207:2026. [ISO/IEC/IEEE 14764:2022](https://www.iso.org/standard/80710.html)

### 15.4 Verification, validation and quality during implementation

| Reference and observed status | Documented scope | Proposed SDP use |
| --- | --- | --- |
| [IEEE 1012-2024](https://standards.ieee.org/ieee/1012/7324/) — Active | System/software/hardware V&V. | Justified V&V strategy with tasks, evidence and appropriate independence. |
| [ISO/IEC/IEEE 29119-2:2021](https://www.iso.org/standard/79428.html) — published | Software test management/execution across lifecycle models. | Link design/change basis to planned tests and actual execution. |
| [ISO/IEC 25010:2023](https://www.iso.org/standard/78176.html) — published | SQuaRE product quality model. | Cover quality in contracts/acceptance beyond happy paths. |

**Assessment:** 1016 supports design descriptions but does not alone explain implementation. Study the interaction of 12207, 24748-1, configuration management, V&V and testing on concrete changes. Section 13 proposes an application, not a standard requirement.

**Further standards work:** choose a small reference profile and trace standard editions/concepts/processes to SDP fields/gates, checks and gaps. Specific conformance claims require full texts. A bibliography is not an implemented method.

## 16. Languages and methods to build on

Uses P04 research dated 2026-09-10; profiles were not all rechecked for this revision. No candidate selected; P04 examples remain untested.

| Building block | Candidates | Assessment for this mandate |
| --- | --- | --- |
| Structure/views | [Structurizr DSL](research/existingDesignLanguages/structurizr-dsl.md), [LikeC4](research/existingDesignLanguages/likec4.md) | Trial one model/multiple views; investigate extra Feature, behavior and change semantics. |
| Richer system/design models | [SysML v2](research/existingDesignLanguages/sysml-v2.md), [UML](research/existingDesignLanguages/uml.md), [AADL](research/existingDesignLanguages/aadl.md), [Context Mapper](research/existingDesignLanguages/context-mapper-cml.md) | Balance precision/concepts with editing cost and agent use. |
| Contracts | [OpenAPI](research/existingDesignLanguages/openapi.md), [AsyncAPI](research/existingDesignLanguages/asyncapi.md), [Protobuf](research/existingDesignLanguages/protobuf.md), [Smithy](research/existingDesignLanguages/smithy.md) | Reference suitable existing contracts; explicitly add uncovered semantics. |
| Constraints/bounded behavior | [CUE](research/existingDesignLanguages/cue.md), [OCL](research/existingDesignLanguages/ocl.md), [SCXML](research/existingDesignLanguages/scxml.md), [TLA+](research/existingDesignLanguages/tla-plus.md), [Alloy](research/existingDesignLanguages/alloy.md) | Different questions need different tools; model checking is not product-code checking. |
| Model change/provenance | [Epsilon](research/existingDesignLanguages/epsilon.md), [Edapt](research/existingDesignLanguages/edapt.md), [PROV-O](research/existingDesignLanguages/prov-o.md) | Study matching, transformation, history and derivation before custom transition infrastructure. |
| Presentation | [Mermaid](research/existingDesignLanguages/mermaid.md), [PlantUML](research/existingDesignLanguages/plantuml.md), [D2](research/existingDesignLanguages/d2.md) | Possible outputs; diagram notation does not define full model semantics. |

C4 helps abstraction/views but need not become the metamodel for every constituent/behavior object. JSON/YAML serialize without choosing semantics. Langium/Xtext build languages/editors; assess after demonstrating syntax needs, not as Mermaid-equivalent alternatives.

Further methodological candidates:

- **Arcadia/Capella:** separates needs analysis from architecture design with model tooling. Relevant to connecting operational needs, functions and structure. No pilot or verification of SDP transitions/agent-context coverage. [Eclipse method/tool article](https://www.eclipse.org/community/eclipse_newsletter/2017/december/article3.php)
- **ATAM:** SEI architecture tradeoff analysis. Scenario/quality evaluation may justify phase transitions better than dependency counts alone. Full ATAM is not proposed as mandatory ceremony for small projects. [SEI ATAM](https://www.sei.cmu.edu/library/the-architecture-tradeoff-analysis-method/)

**Provisional recommendation:** define questions/changes before language selection. Trial reuse or small semantic extensions first; justify a new language with concrete gaps in the same case.

## 17. What can be verified programmatically?

### 17.1 Different checking levels

| Level | Example | Passing does not prove |
| --- | --- | --- |
| Model shape | IDs exist, references resolve, fields have valid types. | Correct system description. |
| Model rules | No forbidden dependency in model graph. | No such code dependency. |
| Code structure | Imports/type symbols match boundaries. | Correct runtime wiring. |
| Contract/compatibility | Observed endpoints/messages follow schema/version. | Unchecked units, authorization or error handling. |
| Integrated behavior | Defined Feature scenario works on identified build/configuration. | All inputs, timing, environments/failures. |
| Formal property | Precise invariant proven within formal assumptions. | Arbitrary code correctly realizes the model. |

**Proposed outcomes:** `satisfied`, `violated`, `unknown`, `not_applicable`. Each result identifies property, method, tool/rule version, source/build, variant, evidence and coverage. Parse failures/missing coverage must not become empty “no violations” results. Policies decide which unknowns block particular gates.

Correct programmatic answers mean correctness within defined semantics/checking models, not deciding complete behavioral equivalence between arbitrary designs and code.

### 17.2 Doxygen and annotations

Doxygen can generate XML for downstream processing: a possible documentation/code-link integration point. [Doxygen XML](https://www.doxygen.nl/manual/customize.html#xmloutput)

**Assessment:** annotations link symbols to stable design IDs and express non-inferable intent, but remain declarations. Matching model/comment claims of CH-MEASUREMENT use can falsely reassure when code uses another route. Combine annotations with independent symbol/dependency analysis and tests. Do not copy whole designs into comments. Trial Doxygen support/extraction quality for the actual language; it is not universal.

### 17.3 AST and language-specific analysis

Clang AST matchers find C/C++ patterns near source structure. [Clang AST Matchers](https://clang.llvm.org/docs/LibASTMatchers.html)

**Assessment:** start with reliably extractable symbols, imports, calls, contract implementations and registrations. Syntax trees alone may lack type resolution/complete call graphs. Explicitly handle plugins, reflection, generated code and configuration-driven wiring. Report investigated build variants/code areas.

### 17.4 LLVM

LLVM IR is typed low-level representation with functions, instructions, control flow and metadata, potentially useful for bounded analysis. [LLVM Language Reference](https://llvm.org/docs/LangRef.html)

**Assessment:** Feature IDs, user intent and logical Channels do not automatically survive as understood domain objects. Preserve/construct mappings. Optimization, frontend and debug/metadata configuration affect recoverability. LLVM is not a universal common format across relevant toolchains.

Compiling design and implementation to LLVM does not automate semantic comparison. Different instructions may realize the same property; identical local computations may occur in different system paths. Multiple valid realizations need precise specification/implementation relations, not identical IR.

A later bounded LLVM trial may be valuable with explicit mapping/assumptions. A normalized **design IR** for objects, relations and contracts is a different concept, useful without machine-code generation.

**Provisional recommendation:** start with model checks, code bindings, structural/contract analysis and Feature witnesses. Keep LLVM/formal proof as bounded studies justified by concrete needs.

# Part III — Further work

## 18. Open questions and decisions

| ID | Question | Investigation |
| --- | --- | --- |
| Q01 | When does Functionality add responsibility beyond Feature/pathway/scenario? | Model the same case with/without it. |
| Q02 | Which boundaries trigger Concept2/MVP2 instead of ordinary revision? | Trial local Fix, internal refactor, contract break, process-responsibility move. |
| Q03 | Required design completeness at each maturity? | Define minimum profile and permitted unknowns. |
| Q04 | What is handwritten, derived or referenced? | Measure editing cost/conflict risk. |
| Q05 | Identity through rename/split/merge/repository moves? | Coordinate #7; trial mappings with unknown portions. |
| Q06 | Find unannotated code/consumers? | Independent code/configuration inventory with deliberately omitted consumer. |
| Q07 | Keep agent context short without hiding essentials? | New agent explains/plans same change from different selections. |
| Q08 | Which properties bind machine checks? | Classify rules by check method/inevitable unknowns. |
| Q09 | Preserve learning without adopting unselected solutions? | Follow Concept1 learning to MVP1 disposition. |
| Q10 | Which standard concepts versus SDP-specific ones? | Small standards mapping with explicit differences/full-text checks as needed. |

## 19. Proposed next bounded investigation

Future assignment proposal, not implementation already underway.

**Purpose:** describe one small change so owner, Worker and Verifier see the same relationships and tools reveal at least one deliberate deviation.

**Case:** section 14 as an explicitly hypothetical fixture, or a separately confirmed product change. One Feature, one Functionality candidate, one Channel, at most six responsibility units. Do not model all Concept1→MVP1 first.

**Three comparison tracks**, following P04:

1. Concrete architecture DSL, such as Structurizr/LikeC4.
2. Richer model approach, such as CML/SysML v2.
3. Small explicit JSON/CUE-style model with bounded SDP semantics.

Select tools, versions and terms before execution. Use identical needs, B/T, assumptions and assessment criteria. Distinguish built-in support, metadata conventions, manual work and new code.

**Deliverables:** before/after models; total/object/residual diffs; responsibility map; Worker/Verifier context selections; learning entry and phase decision; concise evaluation of maintenance, readability and checking limits.

**Negative cases:** broken reference; rename without changed responsibility; split/move with partial behavior preservation; changed unit with unchanged data type; omitted consumer; code violating dependency rules; missing runtime registration; stale model context. First comparisons may be manual if labeled; do not credit report-only manual checks as automation.

Measure editing effort, detected issues, lost information and remaining unknowns. Record costs of a normal change and model-schema revision.

**Stop:** deliver comparison and recommended next decision. Do not secretly freeze language/schema, build a general compiler or migrate product repositories. P04's maximum one repair round per tool track may carry forward; document unresolved tool issues within that boundary.

## 20. Continuation rules for agents

1. Preserve distinctions between owner intent, accepted decisions, observations and proposals.
2. Start from current document revision/repository/issue evidence, not chat memory or old status wording.
3. Investigate specific questions; record question, method, sources, results, limits and recommendation.
4. Do not make Functionality, phase boundaries, DSLs or technology preferences normative without decisions.
5. Use dated primary sources/specific standard editions; distinguish summaries from full normative review.
6. Preserve learning/decision history; use revisions/supersession for changed understanding.
7. Record negative results/blind spots; never claim more coverage than evidence supports.
8. Coordinate method proposals with #9 and #7/#8; preserve #5 foundations and review findings.
9. Separate examples from product facts; confirm product baselines before modeling current systems.
10. End research assignments with concrete decisions or bounded questions, avoiding endless exploration loops.

## 21. Revision log and delivery verification

| Revision | Date | Change / decision |
| --- | --- | --- |
| 0.1 | 2026-09-14 | Restated mandate, relationship to existing SDP, initial standards survey, proposed phase/diff/evidence model and bounded follow-up study. No new method/language acceptance. |

Verification scope: local index/synthesis documents and relevant issue/PR information read; public standards catalogues and linked technical primary sources examined; local Markdown links/code blocks checked. No compiler, parser or product integration implemented/tested here. No full standards-conformance assessment or independent review approval claimed.
