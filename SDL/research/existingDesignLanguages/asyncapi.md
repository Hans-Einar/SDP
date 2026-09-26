# AsyncAPI 3.0

[Catalogue](README.md) · Category: **Event/message contract** · Research: **2026-09-10**

## Purpose and abstraction level

Message-oriented APIs with channels, messages, operations and protocol bindings.

## Model mechanisms

Identity: component/operation names and references. Relations: operation→channel→message. Contracts: payloads and bindings. Views: tool-generated documentation/code. Extension: extensions and external schema formats.

## Strengths and limitations — our assessment

**Strength:** Close to the owner’s input/output events and subscribers; makes messaging boundaries explicit independently of renderers.

**Limitation:** A send/receive definition does not guarantee ordering, delivery, backpressure or handling of every runtime event; specify and test these separately.

## History, change and transitions

Document versions and message compatibility differ. Changes need consumer-impact analysis, migration plans and history beyond one contract.

## Illustrative example

Minimal message contract without a selected broker or delivery guarantee. The example has not been parser/runtime tested.

```yaml
asyncapi: 3.0.0
info:
  title: Measurements
  version: 1.0.0
channels:
  readings:
    address: measurements
    messages:
      measured:
        payload:
          type: number
operations:
  emitReading:
    action: send
    channel:
      $ref: "#/channels/readings"
```

## Tools, maintenance and terms

Official 3.0.0 specification is available. Pin tool/binding versions in a pilot. Assess specification terms separately from generator/runtime licenses.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [AsyncAPI 3.0.0](https://www.asyncapi.com/docs/reference/specification/v3.0.0)
