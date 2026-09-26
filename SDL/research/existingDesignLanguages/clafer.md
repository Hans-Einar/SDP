# Clafer

[Catalogue](README.md) · Category: **Structural/feature modeling with constraints** · Research: **2026-09-10**

## Purpose and abstraction level

Combines lightweight structural modeling, cardinalities, references and feature variability. Solver support explores possible instances.

## Model mechanisms

Identity: declared clafers/references. Relations: nesting, cardinality, references and inheritance. Contracts: logical constraints. Views: generated instances/analyses, not a complete SDP viewpoint package. Machine processing through Clafer tools.

## Strengths and limitations — our assessment

**Strength:** Expresses structural alternatives and can expose impossible combinations before implementation.

**Limitation:** Solver scope and model abstraction limit findings. Class-like syntax does not require implementation through inheritance.

## History, change and transitions

Structure/variability modeling does not provide general lifecycle history. Assess temporal support in particular extensions separately; this profile does not rely on it.

## Illustrative example

Grouping of alternatives; syntax/solver configuration needs checking in the selected implementation. The example has not been parser/runtime tested.

```text
Monitor
  xor ui
    Web
    Desktop
```

## Tools, maintenance and terms

Official language site/compiler repository available. Metadata: MIT, not archived, pushed July 2026. This signals activity, not user-base size or production maturity.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Official Clafer site](https://www.clafer.org/)
- [Official repository; metadata checked through GitHub API](https://github.com/gsdlab/clafer)
