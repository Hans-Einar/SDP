# LikeC4

[Catalogue](README.md) · Category: **Extensible architecture DSL** · Research: **2026-09-10**

## Purpose and abstraction level

Textual architecture models with custom element types/relations. Multiple source files combine into one model; views select displayed elements.

## Model mechanisms

Identity: qualified element references. Relations: explicit connections. Contracts can be referenced/modeled, but arrows do not prove API conformance. Viewpoints: static/dynamic views and predicates. Extension: custom kinds, metadata, generators/API.

## Strengths and limitations — our assessment

**Strength:** Custom kinds can express SDP concepts without imposing original C4 levels. Programmatic model access is relevant to Analyzer.

**Limitation:** Naming an element feature does not supply requirements/release/transition semantics. Flexibility requires a validated profile.

## History, change and transitions

Cross-file model extension is composition, not history. Git and separate before/after mappings must preserve identity through rename/split/merge.

## Illustrative example

Architecture with a custom element type. The example has not been parser/runtime tested.

```text
specification {
  element service
}
model {
  source = service 'Source'
  sink = service 'Sink'
  source -> sink 'Measurements'
}
views {
  view overview { include * }
}
```

## Tools, maintenance and terms

Documentation shows CLI, editor support, API, generators and releases. Repository metadata: not archived, MIT. Pin exact version/API stability in a pilot.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [DSL introduksjon](https://likec4.dev/dsl/intro/)
- [Official repository; metadata checked through GitHub API](https://github.com/likec4/likec4)
