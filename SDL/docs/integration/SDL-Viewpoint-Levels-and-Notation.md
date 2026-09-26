# SDL — viewpoints by abstraction level and consistent notation

**Design basis:** G6-D2, 2026-09-22. **Status clarification, 2026-09-24:** navigator/static export, symbol profile and class-core have since been delivered in G6. [Go commands](../../go/README.md), [symbol profile 1](../profiles/SDL-Symbol-Profile.md), [class-core 0.1](../profiles/SDL-Class-Profile.md) and [evidence](../../go/evidence/G6.md) define implemented scope. The remainder preserves rationale/proposals, not support claims for every proposed type/UML relation. See the [G6 plan](../../../SDUI/docs/implementation-plan.md) and [navigation design](SDL-Navigable-Viewpoints-Design.md).

## Two equivalent export forms

**Navigator/overview** is recommended for development. Generate navigation, small phase overviews and name/type catalogues from validated UseCase, Feature, Functionality, Capability, Activity, Mode, Container, Unit, Interface and Channel objects, with counts/stable identities. Build detailed diagrams, contributions, providers, dependencies and selected scenarios on click. Navigator creation requires neither detailed SVG nor hidden full rendering.

**Static packages** materialize explicit selections/all viewpoints for sharing, archives and readers without handlers. Combined reports are optional. Both forms share projectors/selection descriptions; neither requires first rendering hundreds of diagrams.

Navigators may be dynamic or exported with relative fallbacks. Ungenerated pages need action links, not broken file links. Generic readers receive explanatory overviews/static-export options. Clicks specify viewpoint, focus, relations, direction, depth, mode and notation profile. Servers bound selection/work; requests are not arbitrary code. Navigator/pages identify source revisions/selection keys.

## Abstraction levels and document structure

Checkpoint [02, sections 2 and 6](../../../SDP/History/checkpoint-1/02-SDL-Model-and-Abstraction-Levels.md) defines A0–A5 with candidate Intent, Architecture, Detailed design and Execution binding profiles, not parser headers. Retain these levels:

| Document area | Level | Overview / relevant views |
| --- | --- | --- |
| Requirements and intent | A0 needs/obligations, A1 functional intent | Actors, use cases, Features, Functionality, Capability, high-level Activity, requirements, traceability |
| Architecture | A2 system/Container, A3 internal Units/layers | Container/Unit maps, Channels, ports, providers, contributions, allocation |
| Detailed design | A4 behavior/contracts | Sequences, Activity details, data/message contracts, packets; later classes/state machines |
| Realization and evidence | A5 | Code bindings, generation, verification, executed traces, deployment |

Proposed minimum navigator package:

```text
viewpoints/
  index.md
  navigator.md
  requirements/
    index.md
    inventory.md
  architecture/
    index.md
    inventory.md
  design/
    index.md
    inventory.md
  realization/
    index.md
    inventory.md
```

Directories locate navigation, not duplicate SDL objects. Functionality/contracts may matter at multiple levels with the same ID/source. Build from explicit viewpoint registries specifying levels, questions, required types/relations, permitted selections and notation. Type alone does not determine level: Activity may be A1 or A4. Show missing classification as unspecified, not inferred from names.

Provisional registry: VP01 primarily A0/A1; VP02 A2/A3; VP03/VP07 connect A1 to A2/A3; VP04/VP05 A2/A3; VP08 A1 or A4 by explicit scenario scope; VP09 A1/A4; VP10 A4; VP11 cross-level source/traceability. VP06's G1–G6 work plan has its own menu entry; implementation phases are not system abstraction levels/runtime Activities. G phases describe our work; A levels describe system design.

Requirement, System, Class and State are not declaration types in design-core 0.5. Distinguish **unsupported profile**, **not modeled** and **no selection matches**. Empty inventories do not prove absent requirements/states. Label planned/documented content clearly.

## Mode and State

Active SDL **Mode** names operating/application contexts bounding dependencies, allocation and Channel participation, such as UiPreview, BoundExecution and DocumentBrowsing. Activation, mutual exclusion and hierarchies are unimplemented.

**State** describes a particular system/activity object's condition, such as Ready, Processing or Suspended, with defined ownership/lifetime. Transitions require events/conditions/effects; names alone do not form state machines. Checkpoint State/Activity semantics remain partly candidate material, not parser support.

Activity may move Ready→Processing→Ready within one Mode. States may matter in several Modes. Future models might condition Modes on state predicates or define Mode-local state machines, but nothing is inferred. **State is not generally a child of Mode; the terms are not synonyms.** Display them separately when profiles support both.

## Consistent visual vocabulary

Use UML when semantics match and explicit SDL symbols for its own concepts. This proposed display profile introduces no language rules. Every symbol carries type/name; color alone never distinguishes types. Detail changes contents/compartments, not recognizable type icons.

| Concept | Visual identity |
| --- | --- |
| Actor | Named UML actor figure; external system actors may be labeled accordingly |
| UseCase | Named ellipse |
| Requirement | Document shape with ID/short obligation; planned type |
| Feature | Tabbed rectangle with Feature label |
| Functionality | Rounded responsibility card with vertical side marker/type label |
| Capability | Consistent capability icon/type label, distinct from UseCase ellipse |
| Activity | Rounded activity shape/type label; details may expose flow |
| Mode | Context tab/frame/type label, no start/end marker |
| State | UML state inside a named state machine; planned type |
| System / Container / Unit | Nestable boundaries with distinct type icons; Container is not automatically UML deployment node |
| Interface | Labeled contract port; lollipop/socket only for explicit provided/required modeling |
| Channel | Named connection/junction with endpoints, roles, contract |
| Class | UML compartments for attributes/operations; explicit class models only |

Use-case views need actors/ellipses. Feature/Functionality traceability may be a linked separate diagram. Infer neither System boundaries nor include/extend from placement/supports facts.

## Arrows: UML where semantics match

[OMG UML 2.5.1](https://www.omg.org/spec/UML/2.5.1), sections 7.7 (dependency/realization) and 9.5 (Property/aggregation), is the semantic reference. [Mermaid class diagrams](https://mermaid.js.org/syntax/classDiagram.html) support dependency, realization, aggregation and composition notation; renderer syntax does not prove matching SDL semantics.

| SDL fact / future relation | Display rule |
| --- | --- |
| Unit consumes Interface | UML-inspired dashed dependency/open arrow toward Interface; retain consumes label. Infer no actual message or mandatory dependency. |
| Participant uses Channel as sender/receiver | Channel participation/role, not UML Usage merely because it says uses. |
| Functionality realizes Capability | Means **contributes to realization**, not complete satisfaction. Retain labeled SDL realizes arrow; no hollow UML triangle without stronger explicit contract. |
| Unit provides Capability | Labeled provision; do not rename Capability to Interface/lollipop. |
| supports / contributes-to / addresses | Labeled SDL traceability; no implicit inheritance/include/extend. |
| contains / owns | Logical decomposition/responsibility; no composition diamond without whole/part/lifetime contract. |
| Future explicit UML realization | Dashed line/hollow triangle toward realized specification. |
| Future aggregation / composition | Open/filled diamond at whole; explicit direction, roles, multiplicity and ownership. |

Keep uses/consumes/realizes labels even with familiar arrows. Legends distinguish SDL/UML. UML-inspired is not a formal UML profile. Establish meaning before selecting arrows.

## Renderer profile and later class diagrams

Official Mermaid documentation describes [usecase-beta](https://mermaid.js.org/syntax/usecase.html) with actors/ellipses. This does not prove local Rust-backend support. A local mmdr probe returned 0 but drew syntax words as nodes, not a use-case diagram. It is not approved for that notation. Future capability tests must inspect AST/shapes/markers, not merely exit code/SVG existence. This assignment changes no renderer code.

G6-M5 establishes versioned symbols/arrows, correct use-case views and tested capabilities. Missing support produces labeled alternatives/diagnostics; do not call flowcharts UML use-case diagrams. Full/selected views share symbols/projections.

G6-M6 plans class diagrams after explicit classes, attributes/operations, association ends/roles, multiplicity and aggregation/composition semantics. Unit-contains does not imply class structure. Deliver profiles, validation, negative tests and source maps together. State machines/include/extend similarly need future explicit profiles.

Acceptance: navigator-only renders no details; selections match static-export subsets; symbols retain type/source ID; relations have correct directions/markers; monochrome/nested boundaries remain readable; unsupported diagrams cannot pass on exit code alone.

The model separates NavigableDesignDocumentation (navigation/delivery) from TypedDesignInspection (level catalogues, relation selection, semantic notation). Both support BrowseDesignViews; the latter also supports InspectModels, producing distinct source-linked contribution views.

Verification limit: an initial VP07 with all G6 responsibilities in one Feature hit the 60-second rendering limit. The final model separates the two needs; successful smaller exports do not remove that limit. G6-M2 also needs selection/diagram bounds and heavy-selection diagnostics. No generator timeout or external renderer changed here.
