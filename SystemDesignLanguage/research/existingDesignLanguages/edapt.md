# Edapt

[Katalog](README.md) · Kategori: **Modellhistorikk-/migrasjonsrammeverk (tilgrensende)** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Lagrer endringer mellom Ecore-versjoner som en historie med migrasjonsinstruksjoner; runtime kan migrere eldre modellinstanser.

## Modellmekanismer

Identitet/relasjoner: Ecore-/change-modeller. Contracts: metamodel og migrasjonsoperasjoner. Views: inspeksjon av operasjoner/effekter gjennom verktøy. Utvidelse: nye operasjoner og migrasjonsinstruksjoner. Ikke et generelt arkitektur-DSL.

## Styrker og begrensninger — vår vurdering

**Styrke:** Det tydeligste undersøkte eksemplet på eksplisitt operasjonshistorikk knyttet til modellendring.

**Begrensning:** Migrasjon av modellinstanser er forskjellig fra refactorering av et distribuert produkt. En ny Ecore-struktur sier ikke at alle FEAT-kontrakter fortsatt er oppfylt.

## Historikk, endring og transitions

Historikk er uttrykkelig en del av rammeverkets change-modeller, ikke bare Git. Designansvar, eiergodkjenning og implementeringsbevis må fortsatt modelleres separat.

## Illustrativt eksempel

Illustrativ operasjonsbeskrivelse, ikke Edapt API-/filsyntaks. Eksemplet er ikke parser-/runtime-testet.

```text
Ecore revision A -> revision B
record operation: rename classifier Reading to Measurement
attach migration instructions
migrate stored instances from A to B
```

## Verktøy, vedlikehold og vilkår

Eclipse-siden angir Mature, EPL-2.0 og siste oppførte release 1.5.0 fra 2022-05-20. Det er ikke bevis på opphørt aktivitet, men nyere kompatibilitet er et konkret uavklart punkt.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Edapt prosjekt, historie og lisens](https://projects.eclipse.org/projects/modeling.edapt)
