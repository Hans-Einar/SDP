# DBML

[Catalogue](README.md) · Category: **Database model/diagram DSL** · Research: **2026-09-10**

## Purpose and abstraction level

Textual database descriptions: tables, columns, indexes and relations. Used by dbdiagram, among others.

## Model mechanisms

Identity: schema/table/column names. Relations: refs/keys. Contracts: database shape, not complete service contracts. Views: tool-generated ER diagrams. Reuse/annotations and SQL import/export depend on supported constructs.

## Strengths and limitations — our assessment

**Strength:** Compact, readable model for a bounded data viewpoint; illustrates the value of a small domain-oriented language.

**Limitation:** Database relations are not feature pathways. SQL dialect/import/export differences may lose information; assess the actual database.

## History, change and transitions

Text versioning is possible; documented structure alone is not an executable, lossless migration plan. Handle renames/data transformations separately.

## Illustrative example

Data relation, not a proposed Ponsse database. The example has not been parser/runtime tested.

```text
Table readings {
  id integer [pk]
  length_m double
  session_id integer
}
Table sessions {
  id integer [pk]
}
Ref: readings.session_id > sessions.id
```

## Tools, maintenance and terms

Official guide/repository available. Metadata lists DBML code as Apache-2.0; hosted dbdiagram services do not automatically share those terms.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [DBML language guide](https://dbml.dbdiagram.io/docs/)
- [Official repository; metadata checked through GitHub API](https://github.com/holistics/dbml)
