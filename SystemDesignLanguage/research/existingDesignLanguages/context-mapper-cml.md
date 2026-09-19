# Context Mapper CML

[Katalog](README.md) · Kategori: **DDD- og arkitektur-DSL** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Bounded contexts, aggregates, context maps og relasjoner fra domain-driven design. Dokumentasjonen inkluderer user requirements, stakeholders og generators.

## Modellmekanismer

Identitet: navngitte contexts/aggregates og references. Relasjoner: DDD-typer. Contracts: blant annet Published Language/Open Host Service og kontraktsgenerering som egen funksjon. Views: context-map-/PlantUML-generators. Utvidelse/integrasjon: Xtext og bibliotekbruk.

## Styrker og begrensninger — vår vurdering

**Styrke:** Spesielt relevant for brukerens spørsmål om domains og ansvar. Verktøyet har split/merge-refactoreringer, inkludert split av bounded context etter brukte use cases/user stories.

**Begrensning:** En bounded context er en begreps-/modellgrense; ikke nødvendigvis én OS-prosess. Modellens kompilering etter refactor beviser ikke at produktets atferd er bevart.

## Historikk, endring og transitions

Refactoreringene oppdaterer referanser i CML. Vedvarende Feature-ID, begrunnelse og gammel→ny ansvarsmapping må fortsatt lagres; Git alene forklarer ikke meningen med endringen.

## Illustrativt eksempel

Illustrativ bounded context; ikke et forslag om å omorganisere Ponsse. Eksemplet er ikke parser-/runtime-testet.

```text
BoundedContext Measurement {
  Aggregate Log {
    Entity Reading {
      aggregateRoot
    }
  }
}
```

## Verktøy, vedlikehold og vilkår

VS Code/Eclipse-støtte dokumentert; online IDE er merket ikke lenger støttet. Repositorymetadata: Apache-2.0, ikke arkivert, siste push oppgitt 2025-07-08. Praktisk vedlikehold og kompatibilitet må undersøkes i en pilot.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [CML språkoversikt](https://contextmapper.org/docs/language-reference/)
- [Arkitekturrefactoreringer](https://contextmapper.org/docs/architectural-refactorings/)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/ContextMapper/context-mapper-dsl)
