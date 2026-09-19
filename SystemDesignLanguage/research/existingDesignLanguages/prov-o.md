# PROV-O med RDF/Turtle

[Katalog](README.md) · Kategori: **Provenance-ontologi (tilgrensende)** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

W3C-vokabular for entities, activities, agents og avledning. Turtle kan serialisere grafen; PROV-O er semantikken, ikke et eget arkitekturdiagram-DSL.

## Modellmekanismer

Identitet: IRIs. Relasjoner: derivation, revision, attribution og generation. Contracts: ingen domene-/API-kontrakter i seg selv. Views: graph queries og egne renderere. Utvidelse: RDF-vokabularer og kvalifiserte relasjoner.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kan knytte en modellrevisjon eller et blueprint til kilder og produksjonsaktivitet uten å kopiere alt kildeinnhold.

**Begrensning:** En provenance-påstand er ikke bevis på at avledningen er korrekt. SDP trenger fortsatt identitetsregler og en avklart policy for hvilke fakta som er pålitelige.

## Historikk, endring og transitions

wasRevisionOf og wasDerivedFrom uttrykker historiske koblinger eksplisitt. De gir ikke en kjørbar migrasjon eller forklaring av hvert ansvar som flyttes.

## Illustrativt eksempel

Illustrativ provenancegraf. Entity-navnene er revisjoner/artefakter, ikke automatisk den varige Feature-identiteten. Eksemplet er ikke parser-/runtime-testet.

```turtle
@prefix prov: <http://www.w3.org/ns/prov#> .
@prefix ex: <https://example.org/design/> .
ex:revision2 a prov:Entity ;
  prov:wasRevisionOf ex:revision1 .
ex:blueprint2 a prov:Entity ;
  prov:wasDerivedFrom ex:revision2 .
```

## Verktøy, vedlikehold og vilkår

W3C Recommendation er tilgjengelig. RDF-verktøy og validatorer velges separat med egne lisenser; ontologien alene gir ingen komplett repositoryløsning.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [W3C PROV-O](https://www.w3.org/TR/prov-o/)
