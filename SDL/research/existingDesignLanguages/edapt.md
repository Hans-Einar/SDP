# Edapt

[Catalogue](README.md) · Category: **Model history/migration framework (adjacent)** · Research: **2026-09-10**

## Purpose and abstraction level

Records changes between Ecore versions as histories with migration instructions; runtime can migrate older model instances.

## Model mechanisms

Identity/relations: Ecore/change models. Contracts: metamodel and migration operations. Views: tool inspection of operations/effects. Extension: new operations/instructions. Not a general architecture DSL.

## Strengths and limitations — our assessment

**Strength:** The clearest investigated example of explicit operation history associated with model changes.

**Limitation:** Model-instance migration differs from distributed-product refactoring. New Ecore structures do not prove every FEAT contract remains satisfied.

## History, change and transitions

History explicitly belongs to framework change models, not just Git. Model design responsibility, owner approval and implementation evidence separately.

## Illustrative example

Operation description, not Edapt API/file syntax. The example has not been parser/runtime tested.

```text
Ecore revision A -> revision B
record operation: rename classifier Reading to Measurement
attach migration instructions
migrate stored instances from A to B
```

## Tools, maintenance and terms

Eclipse lists Mature, EPL-2.0 and latest listed release 1.5.0 dated 2022-05-20. This does not prove inactivity, but newer compatibility is unresolved.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Primary documentation](https://projects.eclipse.org/projects/modeling.edapt)
