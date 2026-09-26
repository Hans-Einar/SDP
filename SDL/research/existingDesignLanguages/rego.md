# Rego / Open Policy Agent

[Catalogue](README.md) · Category: **Policy language** · Research: **2026-09-10**

## Purpose and abstraction level

Declarative rules evaluate structured input/data. Can enforce machine-checkable design boundaries when a reliable model/code inventory exists.

## Model mechanisms

Identity/relations: input fields/references, not a built-in architecture registry. Contracts: policy decisions. Views: results/violations; external rendering. Reuse: packages/rules/data. Evaluation through OPA.

## Strengths and limitations — our assessment

**Strength:** A separate checking layer can reject forbidden dependencies without forcing every design object into another language.

**Limitation:** Policies only check supplied facts. A Worker supplying incorrect inventory may obtain passing model checks despite incorrect code.

## History, change and transitions

Git/policy bundles version rules/data. No built-in FEAT transition model; connect policy versions and observed code SHA to evidence.

## Illustrative example

Rego v1 over hypothetical dependencies; no scanner/enforcement implemented. The example has not been parser/runtime tested.

```rego
package sdp.fences
import rego.v1

deny contains "Domain must not import Renderer" if {
  some edge in input.dependencies
  edge.from == "Domain"
  edge.to == "Renderer"
}
```

## Tools, maintenance and terms

Official guide/OPA repository available. Metadata: Apache-2.0, not archived. Pin evaluation version/input contract.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Rego policy language](https://www.openpolicyagent.org/docs/policy-language)
- [Official repository; metadata checked through GitHub API](https://github.com/open-policy-agent/opa)
