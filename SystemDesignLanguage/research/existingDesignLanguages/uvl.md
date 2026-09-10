# Universal Variability Language (UVL)

[Katalog](README.md) · Kategori: **Feature-/variabilitetsmodell** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Beskriver programvareproduktlinjer: tillatte kombinasjoner av features, grupper, attributter og constraints.

## Modellmekanismer

Identitet: navngitte features med hierarki/referanser. Relasjoner: obligatorisk, valgfritt, alternative grupper og constraints. Contracts: konfigurasjonsbetingelser, ikke tjeneste-API. Views: featuretre/analyser via verktøy. Maskinell parsing og analyse; språkets støttede nivå må avklares per verktøy.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kan undersøke om varianter som WebUI og DesktopUI kan velges konsistent med øvrige kapabiliteter.

**Begrensning:** Feature betyr et konfigurasjonsvalg. Det er ikke automatisk den varige REQ-koblede SDP-kapabiliteten eller en implementeringspathway.

## Historikk, endring og transitions

En produktvariant er ikke en tidsrevisjon. Evolusjon av featuremodellen, identitet ved omdøping og migrasjon av lagrede konfigurasjoner krever ekstra regler og versjonering.

## Illustrativt eksempel

Illustrativ UVL-variantstruktur. Ingen påstand om at Ponsse skal ha akkurat én UI-variant. Eksemplet er ikke parser-/runtime-testet.

```text
features
    Monitor
        alternative
            WebUI
            DesktopUI
```

## Verktøy, vedlikehold og vilkår

Offisiell språkside og parser finnes. Parserrepositoryet er ikke arkivert og metadata oppgir LGPL-3.0. Lisensen gjelder parseren, ikke automatisk alle UVL-verktøy.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [UVL språk og økosystem](https://universal-variability-language.github.io/)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/Universal-Variability-Language/uvl-parser)
