# UML 2.5.1

[Catalogue](README.md) · Category: **Model/design language** · Research: **2026-09-10**

## Purpose and abstraction level

General structure/behavior: components, classes, interfaces, activities, sequences and state machines. OMG publishes metamodel/XMI artifacts; UML is more than class diagrams.

## Model mechanisms

Identity: model elements/XMI IDs. Relations: typed dependencies, association and realization. Contracts: interfaces/constraints, often OCL. Viewpoints: multiple diagram types. Extension: profiles/stereotypes; machine interchange through XMI.

## Strengths and limitations — our assessment

**Strength:** Broad semantics and multiple perspectives on one model. Interfaces need not require object-oriented implementation.

**Limitation:** Large models/tool-specific profiles complicate overview/interchange. A diagram alone does not ensure cross-view consistency.

## History, change and transitions

State machines describe execution, not model revisions. Baselines, semantic diffs and split/merge tracking need repository/model tools and explicit identity rules.

## Illustrative example

Pseudonotation for UML elements, not standardized textual UML grammar. The example has not been parser/runtime tested.

```text
Class Measurement
  property value : Real
  property unit : String
Component Domain realizes interface MeasurementSource
Component Representation uses interface MeasurementSource
```

## Tools, maintenance and terms

OMG lists 2.5.1 (December 2017). Normative specification available; this does not imply every tool supports it fully. Assess specification terms and model-tool licenses separately; no editor selected.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Primary documentation](https://www.omg.org/spec/UML/)
