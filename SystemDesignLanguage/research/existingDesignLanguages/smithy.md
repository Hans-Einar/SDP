# Smithy 2.0

[Katalog](README.md) · Kategori: **Service-/data-IDL** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Protokolluavhengige service- og datamodeller med shapes, traits og kodegenerering. Har tekstlig IDL og en maskinell modellrepresentasjon.

## Modellmekanismer

Identitet: namespace og shape-ID. Relasjoner: members, operation inputs/outputs og services. Contracts: strukturelle constraints/traits og protokollspesifikke bindings. Utvidelse: custom traits og generators. Views: tooling/projections, ikke en komplett arkitekturmodell.

## Styrker og begrensninger — vår vurdering

**Styrke:** Førsteklasses shape-identitet og eksplisitte utvidelser er relevante for kontrakter som kan knyttes til FEAT. Mixins viser gjenbruk av modellmedlemmer uten å bestemme klassedesignet i runtime.

**Begrensning:** Servicekontrakter dekker ikke feature-ansvar gjennom alle lag. En custom trait får bare maskinell betydning dersom validator/generator håndterer den.

## Historikk, endring og transitions

Serviceversjon og projections må ikke forveksles med en historikk over designbeslutninger. Endrings-/kompatibilitetsregler og stable IDs må kobles til baselines.

## Illustrativt eksempel

Illustrativ datashape; ingen service eller valgt transport. Eksemplet er ikke parser-/runtime-testet.

```text
$version: "2"
namespace example.measurements
structure Length {
    @required
    meters: Double
}
```

## Verktøy, vedlikehold og vilkår

Offisiell 2.0-dokumentasjon og repository er tilgjengelige. Metadata: Apache-2.0, ikke arkivert. Generatorer/protokoller må velges og versjonspinnes separat.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Smithy IDL](https://smithy.io/2.0/spec/idl.html)
- [Smithy mixins](https://smithy.io/2.0/spec/mixins.html)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/smithy-lang/smithy)
