# Checkpoint #1 — SDL status and generated viewpoints

**Updated implementation status: [supplement 11](11-Go-Implementation-and-Navigation.md).** This retains dated design/V-phase evidence; old Python commands are historical.

Date: 2026-09-22. Documentation/parser/MVP1/SDL-model review; counts below are V4. See [supplement 10](10-SDL-Viewpoint-Navigation.md) for later G6 design and supplement 11 for implementation. [Supplement 07](07-SDUI-0.2-and-Go-Direction.md) retains SDUI 0.2/Go/Fyne direction.

## Conclusion

Checkpoint #1 exceeds implemented SDL: concepts/examples/candidates, not a compiler for its entire vocabulary. V0–V4 deliver bounded **SDL viewpoints**. [Supplement 09](09-SDL-Generated-Go-Design-Review.md) introduces generated G1–G5 review before Go implementation.

SDL tools must generate viewpoints from validated facts. Agents must not draw plausible diagrams and present them as tool output. Views may filter/order facts, not invent architecture.

## Implemented versus documented

| Area | Checkpoint/studies | Implementation at V4 |
| --- | --- | --- |
| Unit, Container, Functionality, Capability, Interface, Activity, Mode | Structural basis | Python design-core 0.5 parser/AST/validation/formatting |
| contains, owns, realizes, provides, consumes, requires in mode, refines | Typed structure | Names, ownership, cycles, source positions |
| Actor, UseCase, Feature, contributions | 02/04 and bounded V1 | 0.2 pursues/supports/contributes-to; direct/many-to-many contributions |
| Functionality allocation | Explicit V1 context | allocated-to Container in mode Mode; at most one Container per responsibility/mode, unchanged logical owner |
| System/multiple files | Studies/MVP1 | No complete resolver/import/deployment/instance model |
| Channel/participants/Contract/MessageSet | [V3](../../../SDL/docs/profiles/SDL-Channel-Scenario-Profile.md) | Typed roles, permits, modes; derived MessageSet |
| Dataset/Datagram/Database/fields | [V2](../../../SDL/docs/profiles/SDL-Data-Contract-Profile.md) | Source/holder, completeness, variants, typed presence, fixed wire layout |
| Scenario steps/request/result | Explicit V3 steps | Order/correlation validated, no runtime/state machine |
| Viewpoint catalogue | Container/dependency/realization/scenario/impact ideas | First concrete catalogue below |
| Markdown/Mermaid | Proposed model/diagram/source package | Bounded tool command in this delivery |
| Go parser/runtime/bindings | Selected direction | Directories, no code yet |

Sources: [definition](../../../SDL/docs/studies/Design-Language-Definition.md), [parser entry point](../../../SDL/go/README.md), [blueprint study](../../../SDL/docs/studies/SDL-Source-Tree-and-Compilation-Study.md), [MVP1 limits](../../../experiments/mvp1_sdl/README.md), [data study](../../../SDL/docs/studies/SDL-Datasets-Datagrams-and-Data-Contracts.md). MVP1 inventory audit checks coverage, not complete parsing.

## Viewpoint catalogue

Bounded recommendation for the owner's printout, not adoption of every construct. Generator IDs identify viewpoints; new entries must declare required facts.

| ID | Question / diagram | Required facts | Tool status |
| --- | --- | --- | --- |
| VP01 | Who wants what? Use cases and contributions | Actor/UseCase, pursues/supports/contributes-to; model scope, no invented System | Generated |
| VP02 | Architecture decomposition | Declarations/contains; allocation separately in VP07 | Generated |
| VP03 | Capability realization and ownership | owns/realizes/provides | Generated |
| VP04 | Used interfaces | consumes; later providers/bindings/contracts | Table generated; no inferred consumer/provider binding |
| VP05 | Requirements in selected mode | requires Interface in mode Mode | Generated |
| VP06 | Phases/milestones | refines/addresses/delivers/depends-on/illustrates/status | Coverage/dependencies generated, no scheduler |
| VP07 | Feature across architecture | Contributions/owners/explicit per-mode allocation | Gaps reported, no mode inheritance |
| VP08 | Channel communication sequence | Participants/roles/contracts/messages/ordered steps/correlation/named alternatives | Generated |
| VP09 | Data origins/retrieval | Dataset, Datagram source/projection, contracts, optional Database, owners | Generated |
| VP10 | Datagram variant encoding | Explicit order/width/offset/variant | Generated for explicit Encoding |
| VP11 | Evidence and gaps | Declarations/facts/source positions/profile limits | Generated |

Feature overlays filter/decorate the same architecture without moving owners/inventing Containers. Preserve many-to-many relations and direct Functionality→UseCase contributions. Use verified notation; labeled goal/traceability flowcharts may substitute for unsupported use-case syntax.

Channel contracts permit collaboration, not necessarily one sequence. VP08 selects declared protocol/scenario order, never sorted relation text masquerading as time. [Mermaid sequences](https://mermaid.js.org/syntax/sequenceDiagram.html) draw supplied order; SDL supplies semantics.

[Packets](https://mermaid.js.org/syntax/packet.html) require bit positions/fields. Logical Datagrams need not be binary. VP09 shows Dataset/Database links from VP10; never invent bit widths for JSON/files/unspecified encodings. Checkpoint 03's 128-bit example was an explicit notation experiment, not an adopted wire contract.

## Database — owner clarification

Database is **an umbrella term for where persistent data can be retrieved when needed**, not necessarily SQL, tables, servers or separate processes. Files, object stores or persistent services may realize it. Dataset defines logical data; Datagram its agreed transfer/projection. Not every Dataset needs persistence.

Clarifies CP1-D11: persistence belongs to Database meaning; duration, availability, ownership and operations still need contracts. Temporary caches are not automatically Databases. Read older retained/queryable wording accordingly. This clarification alone adds no parser support.

## Executable delivery at this snapshot

[SDL tools](../../../SDL/tools/README.md) provide viewpoints using the existing parser/validator, Markdown/Mermaid, separate .mmd files and source maps. Registered renderers also produce SVG/combined printout.md. This is tool functionality, not hand-drawn SDUI-specific diagrams.

[Input](../../../SDUI/design/architecture.design), [generated Markdown/rendered report](../../../SDUI/design/viewpoints/viewpoints.md): 368 declarations, 1106 facts, all available views and explicit gaps. Actor/UseCase/Feature/allocation are authored SDL facts, not reconstructed prose. Six use cases, six Features and two Actors cover inspection, UI trials, domain binding, reload, documentation and native builds. SdlViewpointGenerator is modeled with projection/Markdown/source responsibilities.

VP07 reports 15 unspecified per-mode allocations. Shared responsibilities may be allocated where other Feature contributions are not. This does not prove whole-Feature execution in that mode; contributions are not yet mode-conditioned. These are not full deployment plans for 94 responsibilities.

Verification: 142 SVGs; 24 tool, 61 SDL parser and 36 SDUI tests pass. Diagram nodes link to sources and appear in SVG; repeated exports identical. Visual spot checks, not physical print/PDF. [Report](../../../SDL/tools/verification.json).

## SDL milestones and next delivery

| Milestone | Acceptance |
| --- | --- |
| V0 — delivered | Selectable structural views, deterministic Markdown/Mermaid, source maps, rendering evidence, explicit gaps |
| V1 — delivered | 0.2 goals/contributions/mode allocation; positive/negative tests; generated VP01/VP07 from ported model |
| V2 — delivered | Data/contracts/source/projection/persistence; VP09 and explicit-encoding VP10 |
| V3 — delivered | Typed participation/contracts/ordered messages; request/result VP08 and negative validation |
| V4 — delivered | Enriched source and all requested generated views; no invented generator facts |

After V0–V4, owner scope was review of [generated G1–G5](../../../SDUI/design/viewpoints/implementation.md) before implementation. Five phases/18 milestones cover 94 responsibilities; all were planned. Do not build a large Python runtime for diagrams; coordinate Go port/view needs. MVP1 candidates are not ready grammar for permissive fallback.

SDUI's stopping point then was 0.2 Python frontend/AST/dumps, SVG trial and planned Go/Fyne runtime. No UI-source rewrite needed for this SDL round. Await targeted SDL clarification before runtime work, retaining port evidence.
