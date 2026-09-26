# Context Mapper CML

[Catalogue](README.md) · Category: **DDD and architecture DSL** · Research: **2026-09-10**

## Purpose and abstraction level

Bounded contexts, aggregates, context maps and domain-driven-design relations. Documentation includes user requirements, stakeholders and generators.

## Model mechanisms

Identity: named contexts/aggregates and references. Relations: DDD types. Contracts include Published Language/Open Host Service and separate contract generation. Views: context-map/PlantUML generators. Extension/integration: Xtext and library use.

## Strengths and limitations — our assessment

**Strength:** Especially relevant to domain/responsibility questions. Provides split/merge refactorings, including bounded-context splitting by used use cases/user stories.

**Limitation:** A bounded context is a conceptual/model boundary, not necessarily an OS process. Compiling after refactoring does not prove product behavior preserved.

## History, change and transitions

Refactorings update CML references. Persistent Feature IDs, rationale and old→new responsibility maps still need storage; Git alone does not explain meaning.

## Illustrative example

Bounded context, not a proposal to reorganize Ponsse. The example has not been parser/runtime tested.

```text
BoundedContext Measurement {
  Aggregate Log {
    Entity Reading {
      aggregateRoot
    }
  }
}
```

## Tools, maintenance and terms

VS Code/Eclipse support documented; online IDE labeled unsupported. Repository metadata: Apache-2.0, not archived, last reported push 2025-07-08. Assess practical maintenance/compatibility in a pilot.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Context Mapper documentation](https://contextmapper.org/docs/language-reference/)
- [Context Mapper documentation](https://contextmapper.org/docs/architectural-refactorings/)
- [Official repository; metadata checked through GitHub API](https://github.com/ContextMapper/context-mapper-dsl)
