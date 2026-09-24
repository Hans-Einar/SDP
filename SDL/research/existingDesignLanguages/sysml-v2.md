# SysML v2

[Catalogue](README.md) · Category: **Systems modeling language** · Research: **2026-09-10**

## Purpose and abstraction level

System structure, behavior, requirements, connections and analyses across software/physical systems. V2 has textual/graphical notation and builds on KerML; it is not just another UML profile.

## Model mechanisms

Identity: model elements/qualified names. Relations include parts, ports, connections and requirements links. Contracts: defined interfaces/constraints. Views/viewpoints and libraries belong to the modeling domain. Assess API/repository support as separate specifications/tools.

## Strengths and limitations — our assessment

**Strength:** Links requirements, structure and behavior more closely than pure diagram DSLs; textual format is relevant to agent work.

**Limitation:** Large semantic surface/learning cost. SysML/KerML feature must not automatically mean SDP FEAT.

## History, change and transitions

Distinguish language models from repository commits/branches/API support. Before/after models do not automatically prove migrations. Agree persistent FEAT identity/responsibility mappings.

## Illustrative example

Small textual structure, not a Ponsse mapping or requirements pilot. The example has not been parser/runtime tested.

```sysml
package Demo {
  part def Sensor;
  part def Monitor {
    part sensor : Sensor;
  }
}
```

## Tools, maintenance and terms

OMG has a formal 2.0 page with normative documents. Systems-Modeling provides pilot implementation/examples; pilot support does not prove full conformance. Check OMG document terms and selected implementation LICENSE files separately.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [OMG SysML 2.0](https://www.omg.org/spec/SysML/2.0/)
- [Official repository; metadata checked through GitHub API](https://github.com/Systems-Modeling/SysML-v2-Release)
