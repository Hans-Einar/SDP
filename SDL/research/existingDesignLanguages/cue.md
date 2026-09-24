# CUE

[Katalog](README.md) · Kategori: **Constraint-/konfigurasjonsspråk** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Kombinerer data og constraints gjennom unification. Kan validere og generere konfigurasjon og strukturert data.

## Modellmekanismer

Identitet: fields, definitions og packages; varige design-ID-er må modelleres. Relasjoner: references og constraints. Contracts: datakrav. Utvidelse/gjenbruk: definitions, imports og komposisjon. Views: eksport og egne generators.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kan samle standardkrav og prosjektspesifikke begrensninger uten å kopiere hele schemaet. Go-implementasjonen er interessant, men SDP trenger ikke kreve Go.

**Begrensning:** Ufullstendige verdier og unification har læringskostnad. Gyldig CUE-data beviser ikke at observert kode følger modellen.

## Historikk, endring og transitions

Git lagrer revisjoner. Schema-/datamigrasjon og bevart identitet må spesifiseres separat; unification er ikke en før/etter-transitionmekanisme.

## Illustrativt eksempel

Illustrativ lokal constraint og instans; ikke normativt SDP-format. Eksemplet er ikke parser-/runtime-testet.

```cue
#Feature: {
  id: string & =~"^FEAT-[0-9]+$"
  owner: string
}
feature: #Feature & {
  id: "FEAT-1"
  owner: "Measurement"
}
```

## Verktøy, vedlikehold og vilkår

Offisiell tour og repository er tilgjengelige. Metadata: Apache-2.0, ikke arkivert. Verktøyets versjon og eksportadferd må pinnes i en eventuell pilot.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [CUE constraints](https://cuelang.org/docs/tour/basics/constraints/)
- [CUE tour](https://cuelang.org/docs/tour/)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/cue-lang/cue)
