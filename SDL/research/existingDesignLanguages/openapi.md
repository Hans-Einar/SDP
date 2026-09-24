# OpenAPI

[Catalogue](README.md) · Category: **HTTP API contract** · Research: **2026-09-10**

## Purpose and abstraction level

Machine-readable HTTP operations, parameters, responses, schemas and security requirements.

## Model mechanisms

Identity: paths, operationId and component references. Relations: $ref/API links. Contracts: HTTP/data schemas. Views: tool-generated documentation/clients. Extension: specification extensions; parser/generator/validator support varies.

## Strengths and limitations — our assessment

**Strength:** Can provide an authoritative contract along a feature pathway between backend and client.

**Limitation:** Does not describe all domain behavior, UI layers or reuse policies. Schema validation does not prove backward-compatible user experience.

## History, change and transitions

info.version declares an API-document version, not complete history. Changes require Git, compatibility analysis and consumer mappings.

## Illustrative example

3.1.0 example, deliberately not a demonstration of latest syntax. Payload omitted. The example has not been parser/runtime tested.

```yaml
openapi: 3.1.0
info:
  title: Measurement API
  version: 1.0.0
paths:
  /measurements:
    get:
      operationId: listMeasurements
      responses:
        "200":
          description: Measurements available
```

## Tools, maintenance and terms

The latest page showed 3.2.0 (19 September 2025). Check tool support per version. Specification states Apache-2.0; generators/services may have different terms.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Primary documentation](https://spec.openapis.org/oas/v3.2.0.html)
