# Ecore / Eclipse Modeling Framework

[Catalogue](README.md) · Category: **Metamodel and modeling infrastructure** · Research: **2026-09-10**

## Purpose and abstraction level

Ecore defines model types, attributes and references. EMF surrounds models with serialization/code generation and other infrastructure; it is not one standalone design language.

## Model mechanisms

Identity: EPackage namespace URI/model elements; define instance-ID policy. Relations: EReference cardinality/containment. Contracts: structure; richer constraints need other languages/tools. Extension: custom metamodels. Views/editors need tooling.

## Strengths and limitations — our assessment

**Strength:** Established foundation for custom SDP semantics and model-transformation integration.

**Limitation:** A custom Ecore metamodel remains our language responsibility. JVM/Eclipse infrastructure and model persistence may exceed small-project needs.

## History, change and transitions

Ecore describes types, not their evolution automatically. Edapt/Epsilon Flock provide separate migration mechanisms; Git can add revisions.

## Illustrative example

Ecore class without attributes, not an SDP metamodel decision. The example has not been parser/runtime tested.

```xml
<ecore:EPackage xmi:version="2.0"
 xmlns:xmi="http://www.omg.org/XMI"
 xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
 xmlns:ecore="http://www.eclipse.org/emf/2002/Ecore"
 name="demo" nsURI="urn:example:demo" nsPrefix="demo">
  <eClassifiers xsi:type="ecore:EClass" name="Feature"/>
</ecore:EPackage>
```

## Tools, maintenance and terms

Eclipse project/documentation available; project page lists EPL-2.0. Check bundle NOTICE files/version compatibility before use. No toolchain installed or selected.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Eclipse EMF project](https://projects.eclipse.org/projects/modeling.emf)
