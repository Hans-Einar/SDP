# Smithy 2.0

[Catalogue](README.md) · Category: **Service/data IDL** · Research: **2026-09-10**

## Purpose and abstraction level

Protocol-independent service/data models with shapes, traits and code generation. Textual IDL and machine model representation.

## Model mechanisms

Identity: namespace/shape ID. Relations: members, operation inputs/outputs and services. Contracts: structural constraints/traits and protocol bindings. Extension: custom traits/generators. Views: tooling/projections, not complete architecture models.

## Strengths and limitations — our assessment

**Strength:** First-class shape identity/explicit extensions support FEAT-linked contracts. Mixins reuse model members without dictating runtime class design.

**Limitation:** Service contracts do not cover feature responsibilities across all layers. Custom traits gain machine meaning only through supporting validators/generators.

## History, change and transitions

Service versions/projections are not design-decision history. Link change/compatibility rules and stable IDs to baselines.

## Illustrative example

Data shape, without service or selected transport. The example has not been parser/runtime tested.

```text
$version: "2"
namespace example.measurements
structure Length {
    @required
    meters: Double
}
```

## Tools, maintenance and terms

Official 2.0 documentation/repository available. Metadata: Apache-2.0, not archived. Select/pin generators/protocols separately.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Smithy IDL](https://smithy.io/2.0/spec/idl.html)
- [Smithy mixins](https://smithy.io/2.0/spec/mixins.html)
- [Official repository; metadata checked through GitHub API](https://github.com/smithy-lang/smithy)
