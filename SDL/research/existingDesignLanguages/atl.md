# ATL

[Catalogue](README.md) · Category: **Model-to-model transformation** · Research: **2026-09-10**

## Purpose and abstraction level

Rules produce a target model from a source model with defined metamodels. Eclipse supplies editing/debugging tools.

## Model mechanisms

Identity: source/target elements and transformation links; copy or map stable business IDs explicitly. Relations/contracts depend on metamodels. Views: target models can feed renderers. Extension: helpers/rules and integration.

## Strengths and limitations — our assessment

**Strength:** A concrete precedent for compiling one model into multiple blueprints.

**Limitation:** A correct transformation does not prove correct product migration. Specify losslessness/direction; transformations are not automatically reversible.

## History, change and transitions

Transformation rules can describe old→new model structures. Add revision storage, decision provenance and explanations of split/merge.

## Illustrative example

ATL with hypothetical Source/Target metamodels; not executable without them. The example has not been parser/runtime tested.

```text
module CopyFeature;
create OUT : Target from IN : Source;
rule FeatureToFeature {
  from s : Source!Feature
  to t : Target!Feature (id <- s.id)
}
```

## Tools, maintenance and terms

The Eclipse project page lists Mature, EPL-2.0 and latest listed release 4.12.0 dated 2025-05-19. Runtime/metamodel compatibility was not tested; bundle NOTICE files may add terms.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Eclipse ATL project](https://eclipse.dev/atl/)
- [Eclipse ATL project](https://projects.eclipse.org/projects/modeling.atl)
