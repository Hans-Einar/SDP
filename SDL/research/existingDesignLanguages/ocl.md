# Object Constraint Language (OCL) 2.4

[Catalogue](README.md) · Category: **Model constraints and queries** · Research: **2026-09-10**

## Purpose and abstraction level

Declarative model expressions: invariants, preconditions, postconditions and navigation. Complements UML and other models.

## Model mechanisms

Identity/relations derive from the host model. Contracts: precise logical conditions. Queries can select viewpoint elements; rendering is external. Machine evaluation depends on metamodel, standard library and implementation.

## Strengths and limitations — our assessment

**Strength:** Can express forbidden dependencies or missing contract links.

**Limitation:** Provides neither a system model nor source-code analysis by itself. Constraint names do not prove correct formalization.

## History, change and transitions

Not a migration language. Constraints can check before/after models if both are represented; establish their correspondence separately.

## Illustrative example

OCL assumes a hypothetical metamodel with Feature.requirements; creates no normative SDP rule. The example has not been parser/runtime tested.

```ocl
context Feature
inv HasRequirement: self.requirements->notEmpty()
```

## Tools, maintenance and terms

OMG publishes OCL 2.4. Specification and engine terms differ. Parser/dialect compatibility not tested.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [OMG OCL](https://www.omg.org/spec/OCL/)
