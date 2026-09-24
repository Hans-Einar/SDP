# Epsilon language family

[Catalogue](README.md) · Category: **Model validation, transformation and generation** · Research: **2026-09-10**

## Purpose and abstraction level

EOL underpins ETL (transformation), EVL (validation), EGL (text generation), ECL (matching), EML (merge) and Flock (model migration).

## Model mechanisms

Identity/relations come from connected models; ECL describes correspondences. Contracts: EVL constraints. Views: EGL/Picto derive presentations. EMC connects model formats. Extension: rules, templates and drivers.

## Strengths and limitations — our assessment

**Strength:** Directly relevant to exploring blueprints, validation and model deltas within one ecosystem.

**Limitation:** Multiple small languages/drivers add integration costs. Automatic object matching must not be treated as verified responsibility preservation.

## History, change and transitions

Flock handles metamodel-change updates; ECL/EML assist comparison/merge. Neither automatically provides release history or code-conformance evidence.

## Illustrative example

EVL over a hypothetical Feature type; not an adopted SDP rule. The example has not been parser/runtime tested.

```text
context Feature {
  constraint HasIdentity {
    check: self.id.isDefined()
    message: "Feature needs an identity"
  }
}
```

## Tools, maintenance and terms

Eclipse lists Mature and EPL-2.0, with 2.8 dated 2025-02-20 and future-dated 2.9 (2026-10-19). The latter is not treated as released at the research date. Driver compatibility and bundle terms remain open.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Primary documentation](https://eclipse.dev/epsilon/doc/)
- [Epsilon prosjektstatus](https://projects.eclipse.org/projects/modeling.epsilon)
