# Structurizr DSL / C4

[Catalogue](README.md) · Category: **Architecture model as text** · Research: **2026-09-10**

## Purpose and abstraction level

One architecture model with multiple C4 diagrams: systems, containers, components, deployment and dynamic paths. C4 is the abstraction model; Structurizr DSL is a concrete language.

## Model mechanisms

Identity: DSL identifiers/model elements. Relations: modeled connections. Contracts: descriptions/metadata/external links, not full API semantics. Views: explicit selections. Extension: properties, tags, includes, workspace extensions and plugins.

## Strengths and limitations — our assessment

**Strength:** Candidate for readable owner viewpoints and feature-filtered views over one model.

**Limitation:** Define FEAT→REQ, fences and semantic compatibility beyond standard C4 concepts. DSL names are not stable lifecycle IDs.

## History, change and transitions

Git versions DSL. Server documentation includes workspace versions/branches, but these are tool features; responsibility moves are not built-in verified Feature transitions.

## Illustrative example

Variant of the documented workspace/model/views structure. The example has not been parser/runtime tested.

```text
workspace {
  model {
    owner = person "Owner"
    app = softwareSystem "Monitor"
    owner -> app "Reads measurements"
  }
  views {
    systemContext app { include * }
  }
}
```

## Tools, maintenance and terms

Current documentation shows a unified command surface and marks older Lite/CLI/cloud as end-of-life. Do not base the next pilot on old packaging assumptions. Documentation says prebuilt server requires a license; other commands are free to use. This is not general license approval.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [DSL-eksempel](https://docs.structurizr.com/dsl/example)
- [Primary documentation](https://docs.structurizr.com/)
