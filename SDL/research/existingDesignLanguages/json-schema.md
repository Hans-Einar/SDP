# JSON Schema 2020-12

[Catalogue](README.md) · Category: **Data validation and schemas** · Research: **2026-09-10**

## Purpose and abstraction level

Defines valid JSON instance data and annotations. Can specify a small SDP manifest without requiring a new parser language.

## Model mechanisms

Identity: $id, anchors and references for schema resources. Relations: $ref/$dynamicRef, not automatic foreign-key checking between instances. Contracts: types/structural constraints. Extension: vocabularies. Views/graph queries are external.

## Strengths and limitations — our assessment

**Strength:** Many generic tools can read JSON-based models; lightweight entry point for required-field validation.

**Limitation:** FEAT-ID references, cycles and layer boundaries generally require additional logic. YAML/JSON provide serialization, not semantics for free.

## History, change and transitions

$schema names a dialect; $id identifies a schema resource. Neither provides built-in instance history or transition migration.

## Illustrative example

ID-shape validation, not an adopted SDP schema or ID-uniqueness check. The example has not been parser/runtime tested.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["id"],
  "properties": {"id": {"type": "string", "pattern": "^FEAT-[0-9]+$"}}
}
```

## Tools, maintenance and terms

2020-12 specification/vocabularies available. Check validator dialect support. Specification terms and validator licenses differ; no validator selected.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12)
