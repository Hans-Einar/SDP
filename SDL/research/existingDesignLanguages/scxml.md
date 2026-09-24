# SCXML 1.0

[Catalogue](README.md) · Category: **State-machine/behavior language** · Research: **2026-09-10**

## Purpose and abstraction level

XML descriptions of event-driven state machines, including hierarchical/parallel states and transitions.

## Model mechanisms

Identity: state IDs/event names. Relations: transitions/nesting. Contracts: event-driven behavior; agree payload/schema separately. SCXML processors execute models. Extension: data models/supported execution mechanisms; external viewpoints.

## Strengths and limitations — our assessment

**Strength:** Precisely describes input events that can change component state.

**Limitation:** Correct statecharts provide neither architecture boundaries, FEAT traceability nor automatic tests of actual renderers.

## History, change and transitions

The history element remembers previous runtime substates, explicitly not model history or source-code migration.

## Illustrative example

Event-driven startup. The example has not been parser/runtime tested.

```xml
<scxml xmlns="http://www.w3.org/2005/07/scxml" version="1.0" initial="idle">
  <state id="idle">
    <transition event="start" target="running"/>
  </state>
  <state id="running"/>
</scxml>
```

## Tools, maintenance and terms

W3C Recommendation from 2015 available. Assess processor support, data model and distribution license separately; no runtime selected.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [W3C SCXML Recommendation](https://www.w3.org/TR/scxml/)
