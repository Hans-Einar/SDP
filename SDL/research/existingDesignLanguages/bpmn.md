# BPMN 2.0.2

[Catalogue](README.md) · Category: **Process model and notation** · Research: **2026-09-10**

## Purpose and abstraction level

Business processes, tasks, events, gateways and participant collaboration. Can describe UseCase paths and human approvals.

## Model mechanisms

Identity: XML element IDs. Relations: sequence/message flows. Contracts: messaging/participant boundaries; payloads and software import rules require more. Views: process/choreography/collaboration. Extension elements and machine-readable XSD/CMOF; execution depends on engines.

## Strengths and limitations — our assessment

**Strength:** Makes responsibilities and decision points visible to owners/stakeholders.

**Limitation:** Process flow is neither an architecture map nor a feature pathway through code. Valid diagrams are not necessarily executable.

## History, change and transitions

Runtime process instances and engine process versions are not general model-transition semantics. Moving Feature responsibilities needs external model links.

## Illustrative example

Process without diagram layout or engine configuration. The example has not been parser/runtime tested.

```xml
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
 targetNamespace="urn:example:sdp">
  <process id="review" isExecutable="false">
    <startEvent id="start"/>
    <task id="inspect" name="Inspect design"/>
    <sequenceFlow id="f1" sourceRef="start" targetRef="inspect"/>
  </process>
</definitions>
```

## Tools, maintenance and terms

OMG lists 2.0.2 from January 2014 with normative schemas. Standard terms and chosen editor/runtime licenses are separate; no engine was evaluated.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [OMG BPMN](https://www.omg.org/spec/BPMN/)
