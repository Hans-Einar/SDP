# Alloy 6

[Catalogue](README.md) · Category: **Relational formal modeling language** · Research: **2026-09-10**

## Purpose and abstraction level

Structures, relations and constraints for solver analysis. Alloy 6 includes mutable elements and temporal logic.

## Model mechanisms

Identity: signatures/atoms within analysis. Relations: typed fields and logic. Contracts: facts, predicates and assertions. Views: instances and counterexamples. Extension: modules; machine analysis within explicit scope.

## Strengths and limitations — our assessment

**Strength:** Can find counterexamples to assumptions such as a layered dependency graph or a mapping covering all required responsibilities.

**Limitation:** Solver atoms are not persistent FEAT IDs. Absence of counterexamples must be interpreted within analysis scope, time bounds and assumptions.

## History, change and transitions

Temporal support can model transitions but is not built-in version control. Preserved-property proofs require those properties to be modeled.

## Illustrative example

Unexecuted cycle-check model; four is the analysis scope, not an architecture rule. The example has not been parser/runtime tested.

```alloy
sig Layer { depends: set Layer }
fact Acyclic { no l: Layer | l in l.^depends }
run {} for 4
```

## Tools, maintenance and terms

Official Alloy 6 notes and Analyzer are available. Repository metadata returned NOASSERTION for license; inspect LICENSE/NOTICE before distribution.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Alloy 6](https://alloytools.org/alloy6.html)
- [Official repository; metadata checked through GitHub API](https://github.com/alloytools/org.alloytools.alloy)
