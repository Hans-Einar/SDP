# PROV-O with RDF/Turtle

[Catalogue](README.md) · Category: **Provenance ontology (adjacent)** · Research: **2026-09-10**

## Purpose and abstraction level

W3C vocabulary for entities, activities, agents and derivation. Turtle serializes graphs; PROV-O supplies semantics, not an architecture-diagram DSL.

## Model mechanisms

Identity: IRIs. Relations: derivation, revision, attribution and generation. Contracts: no intrinsic domain/API contracts. Views: graph queries/custom renderers. Extension: RDF vocabularies/qualified relations.

## Strengths and limitations — our assessment

**Strength:** Links model revisions/blueprints to sources and production activities without copying all source content.

**Limitation:** Provenance assertions do not prove correct derivation. SDP still needs identity rules and a policy for trusted facts.

## History, change and transitions

wasRevisionOf and wasDerivedFrom explicitly express history. They provide neither executable migration nor explanations for every moved responsibility.

## Illustrative example

Provenance graph; entity names identify revisions/artifacts, not automatically persistent Feature identity. The example has not been parser/runtime tested.

```turtle
@prefix prov: <http://www.w3.org/ns/prov#> .
@prefix ex: <https://example.org/design/> .
ex:revision2 a prov:Entity ;
  prov:wasRevisionOf ex:revision1 .
ex:blueprint2 a prov:Entity ;
  prov:wasDerivedFrom ex:revision2 .
```

## Tools, maintenance and terms

W3C Recommendation available. Select RDF tools/validators separately under their own licenses; the ontology alone is not a complete repository solution.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [W3C PROV-O](https://www.w3.org/TR/prov-o/)
