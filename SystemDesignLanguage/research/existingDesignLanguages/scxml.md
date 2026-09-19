# SCXML 1.0

[Katalog](README.md) · Kategori: **Tilstandsmaskin-/atferdsspråk** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

XML-basert beskrivelse av eventdrevne tilstandsmaskiner, inkludert hierarkiske/parallelle tilstander og transitions.

## Modellmekanismer

Identitet: state-ID-er og eventnavn. Relasjoner: transitions og nesting. Contracts: hendelsesstyrt atferd; payload/schema trenger separat avtale. Maskinell kjøring i en SCXML-prosessor. Utvidelse via datamodeller og støttede eksekveringsmekanismer; viewpoints er eksterne.

## Styrker og begrensninger — vår vurdering

**Styrke:** Presis beskrivelse av hvilke inputhendelser som kan endre en komponenttilstand.

**Begrensning:** Et korrekt statechart gir ikke arkitekturgrenser, FEAT-sporing eller automatisk test av en faktisk renderer.

## Historikk, endring og transitions

Elementet history husker en tidligere runtime-substate. Det er uttrykkelig ikke modellhistorikk eller kildekodemigrasjon.

## Illustrativt eksempel

Illustrativ eventdrevet oppstart. Eksemplet er ikke parser-/runtime-testet.

```xml
<scxml xmlns="http://www.w3.org/2005/07/scxml" version="1.0" initial="idle">
  <state id="idle">
    <transition event="start" target="running"/>
  </state>
  <state id="running"/>
</scxml>
```

## Verktøy, vedlikehold og vilkår

W3C Recommendation fra 2015 er tilgjengelig. Prosessorstøtte, datamodell og distribusjonslisens må undersøkes separat; ingen runtime er valgt.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [W3C SCXML Recommendation](https://www.w3.org/TR/scxml/)
