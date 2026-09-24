# ArchiMate

[Catalogue](README.md) · Category: **Architecture model and graphical notation** · Research: **2026-09-10**

## Purpose and abstraction level

Enterprise architecture across business, application and technology perspectives. Relevant to showing owner goals/responsibilities above class/function level.

## Model mechanisms

Identity/relations: model elements with typed connections. Viewpoints select model elements for different stakeholders; Archi reuses elements across views. Contracts: high-level service/interface descriptions, not wire schemas. Properties/interchange depend on tools/standards.

## Strengths and limitations — our assessment

**Strength:** Owner-oriented overview and reuse across views; distinguishes models from diagrams.

**Limitation:** Too coarse alone to enforce import boundaries and precise input/output contracts. Graphical editing may raise the entry barrier for a text-based agent workflow.

## History, change and transitions

Implementation/migration is an ArchiMate area, but this research did not verify precise transition semantics against normative Open Group text. Archi undo/redo is editor history, not persistent Feature provenance.

## Illustrative example

Element list for a graphical language, not ArchiMate DSL syntax. The example has not been parser/runtime tested.

```text
Application Component: Domain service
Application Interface: Measurement input
Application Service: Provide measurements
View: component, interface and service in one selection
```

## Tools, maintenance and terms

The inspected Archi manual is 5.10.0. Open Group specification pages could not be retrieved; normative version/licensing assessment remains open. Archi tooling and the ArchiMate standard have different terms.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Archi user guide](https://www.archimatetool.com/downloads/archi/Archi%20User%20Guide.pdf)
