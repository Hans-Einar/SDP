# DBML

[Katalog](README.md) · Kategori: **Databasemodell-/diagram-DSL** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Tekstlig beskrivelse av databaser: tabeller, kolonner, indekser og relasjoner. Brukes blant annet av dbdiagram.

## Modellmekanismer

Identitet: schema-/tabell-/kolonnenavn. Relasjoner: refs og keys. Contracts: databaseform, ikke komplette tjenestecontracts. Views: ER-diagram via verktøy. Gjenbruk/annotations og SQL-import/eksport avhenger av støttede constructs.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kompakt og lesbar modell for et avgrenset data-viewpoint. Et godt eksempel på verdien av et lite domeneorientert språk.

**Begrensning:** En database-relasjon er ikke en feature-pathway. SQL-dialekter og import/eksport kan gi tap; egnethet må vurderes på den faktiske databasen.

## Historikk, endring og transitions

Tekstversjonering er mulig; dokumentert struktur er ikke i seg selv en kjørbar, tapsfri migrasjonsplan. Rename og datatransformasjoner må håndteres separat.

## Illustrativt eksempel

Illustrativ datarelasjon; ikke forslag til Ponsse-database. Eksemplet er ikke parser-/runtime-testet.

```text
Table readings {
  id integer [pk]
  length_m double
  session_id integer
}
Table sessions {
  id integer [pk]
}
Ref: readings.session_id > sessions.id
```

## Verktøy, vedlikehold og vilkår

Offisiell språkguide og repository er tilgjengelige. DBML-kode: Apache-2.0 ifølge metadata. Det gir ikke automatisk samme vilkår for den hostede dbdiagram-tjenesten.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [DBML språkguide](https://dbml.dbdiagram.io/docs/)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/holistics/dbml)
