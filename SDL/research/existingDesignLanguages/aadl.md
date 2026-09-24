# AADL

[Catalogue](README.md) · Category: **Architecture Description Language** · Research: **2026-09-10**

## Purpose and abstraction level

Architecture of software/hardware systems, especially where processes, threads, ports, deployment and analysis matter. OSATE is a concrete modeling/analysis environment.

## Model mechanisms

Identity: packages and named classifiers/implementations. Relations: component hierarchy, connections and bindings. Contracts: ports, data types and properties; annexes may add semantics. Tools generate views/analyses. Extension through properties and annexes.

## Strengths and limitations — our assessment

**Strength:** More precise about execution structure than generic box diagrams; relevant to distributed processes and resource requirements.

**Limitation:** Real-time/embedded orientation may be heavy for a small web app. An AADL feature can mean a port or access, not a user-facing SDP Feature.

## History, change and transitions

Modes describe operating configurations, not version history. Reorganizing responsibilities across revisions requires mapping and external versioning.

## Illustrative example

Minimal component definition, without analysis properties or deployment claims. The example has not been parser/runtime tested.

```aadl
package Demo
public
  system Monitor
  end Monitor;
  system implementation Monitor.impl
  end Monitor.impl;
end Demo;
```

## Tools, maintenance and terms

OSATE documentation is available as 2.19.0. The SAE standard and OSATE distribution have separate terms. Repository metadata provided no unambiguous SPDX license; inspect actual LICENSE/NOTICE files and annex tools before use.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [OSATE and AADL support](https://osate.org/about-osate.html)
- [Official repository; metadata checked through GitHub API](https://github.com/osate/osate2)
