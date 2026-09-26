# ReqIF 1.2

[Catalogue](README.md) · Category: **Requirements interchange (adjacent format)** · Research: **2026-09-10**

## Purpose and abstraction level

Exchanges typed requirements, attributes, hierarchies and relations between requirements tools. Relevant to stakeholder→UseCase→REQ traceability, not complete system architecture.

## Model mechanisms

Identity: element IDENTIFIER. Relations: SpecRelations/structures. Contracts can be described, without API execution semantics. Extension: types/attributes/tool extensions. Machine XML interchange; views depend on requirements tools.

## Strengths and limitations — our assessment

**Strength:** Existing alternative to inventing requirements/identity interchange formats.

**Limitation:** Imported requirements do not prove Feature satisfaction. Tool-dependent round trips need testing.

## History, change and transitions

LAST-CHANGE metadata is not a revision archive. Establish baselines, change approval and FEAT/release links separately.

## Illustrative example

XML fragment, not a complete schema-validated ReqIF document. The example has not been parser/runtime tested.

```xml
<SPEC-OBJECT IDENTIFIER="REQ-1"
 LONG-NAME="Display measured length"
 LAST-CHANGE="2026-09-10T00:00:00Z"/>
<!-- Fragment: type and the rest of the ReqIF document are omitted. -->
```

## Tools, maintenance and terms

OMG publishes ReqIF 1.2. Standard terms differ from importer/exporter licenses. No specific round trip/editor evaluated.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [OMG ReqIF](https://www.omg.org/spec/ReqIF/)
