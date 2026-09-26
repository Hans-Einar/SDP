# Survey of existing design languages

Research date: **2026-09-10**. Scope: **32 profiles** of languages, families, notations and explicitly adjacent formats/frameworks. Broad and targeted, not an exhaustive inventory of every language in use.

See the [synthesis](../README.md) for SDP needs, distinctions, questions and next work.

## Coverage

| Profile | Category |
| --- | --- |
| [UML 2.5.1](uml.md) | Model/design language |
| [SysML v2](sysml-v2.md) | Systems modeling language |
| [ArchiMate](archimate.md) | Architecture model and graphical notation |
| [AADL](aadl.md) | Architecture Description Language |
| [Structurizr DSL / C4](structurizr-dsl.md) | Architecture model as text |
| [LikeC4](likec4.md) | Extensible architecture DSL |
| [PlantUML](plantuml.md) | Diagram notation |
| [Mermaid](mermaid.md) | Diagram notation |
| [D2](d2.md) | Diagram DSL |
| [Context Mapper CML](context-mapper-cml.md) | DDD and architecture DSL |
| [Universal Variability Language (UVL)](uvl.md) | Feature/variability model |
| [Clafer](clafer.md) | Structural/feature modeling with constraints |
| [BPMN 2.0.2](bpmn.md) | Process model and notation |
| [SCXML 1.0](scxml.md) | State-machine/behavior language |
| [TLA+ and PlusCal](tla-plus.md) | Formal behavior specification |
| [Alloy 6](alloy.md) | Relational formal modeling language |
| [Object Constraint Language (OCL) 2.4](ocl.md) | Model constraints and queries |
| [OpenAPI](openapi.md) | HTTP API contract |
| [AsyncAPI 3.0](asyncapi.md) | Event/message contract |
| [Protocol Buffers](protobuf.md) | Data/service IDL |
| [Smithy 2.0](smithy.md) | Service/data IDL |
| [JSON Schema 2020-12](json-schema.md) | Data validation and schemas |
| [CUE](cue.md) | Constraint/configuration language |
| [DBML](dbml.md) | Database model/diagram DSL |
| [Ecore / Eclipse Modeling Framework](ecore.md) | Metamodel and modeling infrastructure |
| [ATL](atl.md) | Model-to-model transformation |
| [QVT 1.3](qvt.md) | Query/View/Transformation standard |
| [Epsilon language family](epsilon.md) | Model validation, transformation and generation |
| [Edapt](edapt.md) | Model history/migration framework (adjacent) |
| [Rego / Open Policy Agent](rego.md) | Policy language |
| [ReqIF 1.2](reqif.md) | Requirements interchange (adjacent format) |
| [PROV-O with RDF/Turtle](prov-o.md) | Provenance ontology (adjacent) |

## Methodology and evidence limits

Primary sources are standards bodies, developer documentation and official repositories. Profiles distinguish documented purposes/mechanisms from our assessment of strengths, limits and SDP relevance. Negative lifecycle-support findings concern inspected documentation/standards, not every plugin/research extension.

Examples are **our own illustrations**, not copied normative examples. None was parser/runtime tested. Fragments/pseudonotation are labeled, not implementation-ready fixtures. Full Ponsse representation belongs to the next assignment.

Available documentation, unarchived repositories and recent pushes are distinct activity signals; none alone proves production maturity, adoption or sustained maintenance. Listed versions are observed references, not universal integration support. Future-dated releases are not treated as released.

GitHub license names are metadata signals, not full-file legal reviews. NOASSERTION means no unambiguous metadata result. Standards, parsers, generators, runtimes and hosted services may have different terms. Resolve explicitly open terms for actual tools before integration.

## Significant coverage gaps

ArchiMate was partially checked: Archi's manual was available; normative Open Group pages could not be retrieved. This is not a full normative transition assessment.

Not separately profiled: SysML v1, Capella/Arcadia, Modelica, SDL, Acme/Wright, DMN, Event-B, Z, SHACL, GraphQL, Avro, Thrift, RAML and TypeSpec. Xtext, MPS and Sirius are possible language/editor tools, not selected. Gaps prevent claiming any candidate is best among all languages. Future comparisons may add a candidate for concrete unmet needs.

YAML/JSON serialize; they are not semantic SDP models. C4 is covered with Structurizr. EMF/Edapt, ReqIF and PROV-O address bounded supplementary needs, not complete alternative system-design languages.

## Verification performed

Checked expected sections, examples, sources, unique filenames and internal targets. Web sources were read for research; no complete crawler or installation of 32 tools ran. The original assignment's diff was limited to research. No independent verification approval or product-test success claimed. Translation preserves these dated claims without renewing them.
