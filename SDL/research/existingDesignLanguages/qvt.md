# QVT 1.3

[Catalogue](README.md) · Category: **Query/View/Transformation standard** · Research: **2026-09-10**

## Purpose and abstraction level

OMG family for model operations: Relations, Core and Operational Mappings. Describes consistency/transformations between metamodels.

## Model mechanisms

Identity/relations: elements in participating models. Contracts: constraints/transformation conditions. Views: derived models; visual rendering separate. Execution requires an implementation of the selected standard portion.

## Strengths and limitations — our assessment

**Strength:** Relevant when before/after correspondences must be explicit rather than merely drawn.

**Limitation:** QVT variants differ in operational meaning/tool support. Declared relations do not automatically provide unambiguous, lossless bidirectional migration.

## History, change and transitions

Can express mappings; persistent revision archives and audits of applied rule sets remain outside transformations.

## Illustrative example

QVT Relations sketch with hypothetical metamodels; validate implementation dialect/name resolution before execution. The example has not been parser/runtime tested.

```text
transformation Preserve(source : Old, target : New) {
  top relation PreserveIdentity {
    fid : String;
    checkonly domain source a : Feature { id = fid };
    enforce domain target b : Feature { id = fid };
  }
}
```

## Tools, maintenance and terms

OMG publishes 1.3. Assess standard terms and QVTo/QVTr implementation licenses/support separately. No engine selected.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [OMG QVT](https://www.omg.org/spec/QVT/)
