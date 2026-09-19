# ATL

[Katalog](README.md) · Kategori: **Model-to-model-transformasjon** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Regler som lager en målmodell fra en kildemodell med definerte metamodeller. Eclipse tilbyr editor-/debugging-verktøy.

## Modellmekanismer

Identitet: kilde-/målelementer og transformasjonskoblinger; stabil business-ID må kopieres eller mappes eksplisitt. Relasjoner/contracts: avhenger av metamodellene. Views: en target-model kan mate en renderer. Utvidelse: helpers/regler og integrasjon.

## Styrker og begrensninger — vår vurdering

**Styrke:** Konkret forløper for tanken om å kompilere én modell til flere blueprints.

**Begrensning:** Korrekt transformasjon er ikke korrekt produktmigrasjon. Tapsfrihet og retningsvalg må spesifiseres; en transformasjon er ikke automatisk reversibel.

## Historikk, endring og transitions

Transformasjonsregler kan beskrive gammel→ny modellstruktur. Revisjonslagring, beslutningsprovenance og forklaring av split/merge må legges til.

## Illustrativt eksempel

Illustrativ ATL med hypotetiske Source/Target-metamodeller; ingen kjørbar bundle uten dem. Eksemplet er ikke parser-/runtime-testet.

```text
module CopyFeature;
create OUT : Target from IN : Source;
rule FeatureToFeature {
  from s : Source!Feature
  to t : Target!Feature (id <- s.id)
}
```

## Verktøy, vedlikehold og vilkår

Eclipse-prosjektsiden angir Mature, EPL-2.0 og siste oppførte release 4.12.0 fra 2025-05-19. Runtime-/metamodelkompatibilitet er ikke testet; konkrete bundle-NOTICE-filer kan inneholde ytterligere vilkår.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Eclipse ATL](https://eclipse.dev/atl/)
- [Eclipse prosjektstatus og lisens](https://projects.eclipse.org/projects/modeling.atl)
