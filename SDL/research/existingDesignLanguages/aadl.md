# AADL

[Katalog](README.md) · Kategori: **Architecture Description Language** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Arkitektur for software-/hardware-systemer, særlig der prosesser, tråder, porter, deployment og analyseelementer er viktige. OSATE er et konkret modellerings- og analysemiljø.

## Modellmekanismer

Identitet: packages og navngitte classifiers/implementations. Relasjoner: komponenthierarki, forbindelser og bindings. Contracts: porter, datatyper og properties; annexes kan tilføre semantikk. Views/analyser genereres av verktøy. Utvidelse gjennom properties og annexes.

## Styrker og begrensninger — vår vurdering

**Styrke:** Mer presist om eksekveringsstruktur enn generelle boksdiagrammer. Interessant for distribuerte prosesser og ressurskrav.

**Begrensning:** Sanntids-/embedded-orienteringen kan bli tung for en liten webapp. AADL-feature betyr eksempelvis en port eller tilgang, ikke en brukerrettet SDP Feature.

## Historikk, endring og transitions

Modes beskriver driftskonfigurasjoner; de er ikke versjonshistorikk. Reorganisering av ansvar over revisjoner krever mapping og ekstern versjonering.

## Illustrativt eksempel

Illustrativ minimal komponentdefinisjon; ingen analyseegenskaper eller deploymentpåstand. Eksemplet er ikke parser-/runtime-testet.

```aadl
package Demo
public
  system Monitor
  end Monitor;
  system implementation Monitor.impl
  end Monitor.impl;
end Demo;
```

## Verktøy, vedlikehold og vilkår

OSATE-dokumentasjonen er tilgjengelig som 2.19.0. SAE-standarden og OSATE-distribusjonen har separate vilkår. Repositorymetadata ga ingen entydig SPDX-lisens; konkrete LICENSE/NOTICE og eventuelle annex-verktøy må undersøkes før bruk.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OSATE og AADL-støtte](https://osate.org/about-osate.html)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/osate/osate2)
