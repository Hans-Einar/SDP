# Rego / Open Policy Agent

[Katalog](README.md) · Kategori: **Policyspråk** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Deklarative regler som evaluerer strukturert input og data. Kan brukes til maskinelle designfences når en pålitelig modell eller kodeinventar finnes.

## Modellmekanismer

Identitet/relasjoner: felter/referanser i input, ikke et innebygget arkitekturregister. Contracts: policybeslutninger. Views: resultater/avvik; rendering er ekstern. Gjenbruk: packages/regler og data. Maskinell evaluering med OPA.

## Styrker og begrensninger — vår vurdering

**Styrke:** Et separat kontrollag kan avvise avtalebrudd som forbudte dependencies uten å tvinge alle designobjekter inn i et nytt språk.

**Begrensning:** Policy kan bare kontrollere fakta den mottar. En Worker som leverer feil inventar kan få en grønn modellkontroll uten at koden er riktig.

## Historikk, endring og transitions

Git og policy-bundles kan versjonere regler/data. Ingen innebygd FEAT-transitionmodell. Policyversjon og observert code-SHA må kobles til evidensen.

## Illustrativt eksempel

Illustrativ Rego v1-policy over et hypotetisk dependency-inventar; ingen scanner eller håndheving implementert. Eksemplet er ikke parser-/runtime-testet.

```rego
package sdp.fences
import rego.v1

deny contains "Domain must not import Renderer" if {
  some edge in input.dependencies
  edge.from == "Domain"
  edge.to == "Renderer"
}
```

## Verktøy, vedlikehold og vilkår

Offisiell policyguide og OPA-repository er tilgjengelige. Metadata: Apache-2.0, ikke arkivert. Evalueringsversjon og inputkontrakt må pinnes.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Rego policy language](https://www.openpolicyagent.org/docs/policy-language)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/open-policy-agent/opa)
