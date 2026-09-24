# CUE

[Catalogue](README.md) · Category: **Constraint/configuration language** · Research: **2026-09-10**

## Purpose and abstraction level

Combines data and constraints through unification; validates/generates configuration and structured data.

## Model mechanisms

Identity: fields, definitions and packages; model persistent design IDs explicitly. Relations: references/constraints. Contracts: data requirements. Extension/reuse: definitions, imports and composition. Views: export/custom generators.

## Strengths and limitations — our assessment

**Strength:** Combines standard requirements and project constraints without copying entire schemas. Go implementation is interesting, but SDP need not require Go.

**Limitation:** Incomplete values and unification require learning. Valid CUE data does not prove observed code follows the model.

## History, change and transitions

Git stores revisions. Define schema/data migration and identity preservation separately; unification is not a before/after transition mechanism.

## Illustrative example

Local constraint and instance, not normative SDP format. The example has not been parser/runtime tested.

```cue
#Feature: {
  id: string & =~"^FEAT-[0-9]+$"
  owner: string
}
feature: #Feature & {
  id: "FEAT-1"
  owner: "Measurement"
}
```

## Tools, maintenance and terms

Official tour/repository available. Metadata: Apache-2.0, not archived. Pin tool version/export behavior in any pilot.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [CUE documentation](https://cuelang.org/docs/tour/basics/constraints/)
- [CUE documentation](https://cuelang.org/docs/tour/)
- [Official repository; metadata checked through GitHub API](https://github.com/cue-lang/cue)
