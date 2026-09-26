# TLA+ and PlusCal

[Catalogue](README.md) · Category: **Formal behavior specification** · Research: **2026-09-10**

## Purpose and abstraction level

Mathematically precise descriptions of possible system states/steps. PlusCal is algorithm notation translated into TLA+.

## Model mechanisms

Identity: modules, constants, variables and operators. Relations/contracts: invariants and temporal properties. Views: tool-explored traces/counterexamples, not ready-made architecture views. Reuse through modules/operators.

## Strengths and limitations — our assessment

**Strength:** Relevant to queues, concurrency, delivery guarantees and invariant preservation through transitions.

**Limitation:** Requires appropriate abstraction/expertise. TLC searches finite configured spaces, not general program proofs; connect models/code through evidence requirements.

## History, change and transitions

Temporal logic describes behavior over time. Refinement compares abstractions but does not automatically supply FEAT-ID history, Git migration or refactoring approval.

## Illustrative example

Specification without TLC configuration; not model-checked here. The example has not been parser/runtime tested.

```text
---- MODULE Counter ----
EXTENDS Naturals
VARIABLE n
Init == n = 0
Next == n' = n + 1
Spec == Init /\ [][Next]_n
====
```

## Tools, maintenance and terms

Author documentation/TLA+ tools repository available. tlaplus/tlaplus metadata: MIT, not archived. Explicitly choose model checking versus theorem proving in a pilot.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Leslie Lamports TLA+-side](https://lamport.azurewebsites.net/tla/tla.html)
- [Official repository; metadata checked through GitHub API](https://github.com/tlaplus/tlaplus)
