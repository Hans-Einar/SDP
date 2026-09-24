# ArchiMate

[Katalog](README.md) · Kategori: **Arkitekturmodell og grafisk notasjon** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Enterprise-arkitektur med forretnings-, applikasjons- og teknologiperspektiver. Relevant når SDP må vise eierens mål og ansvar på et høyere nivå enn klasser og funksjoner.

## Modellmekanismer

Identitet/relasjoner: modelelementer med typede forbindelser. Viewpoints: utvalg av modellen for ulike interessenter; Archi gjenbruker elementer mellom views. Contracts: høynivå service-/grensesnittbeskrivelser, ikke wire-schema. Egenskaper og utveksling er verktøy-/standardavhengige.

## Styrker og begrensninger — vår vurdering

**Styrke:** Eierorientert helhetsbilde og gjenbruk av samme element i flere views. Skiller modell og diagram.

**Begrensning:** For grovt alene til å håndheve importgrenser og presise input/output-kontrakter. Grafisk redigering kan gi større terskel for en tekstbasert agentløype.

## Historikk, endring og transitions

Implementasjon/migrasjon er et relevant ArchiMate-område, men presis transition-semantikk er ikke verifisert mot The Open Groups normative tekst i denne runden. Archi undo/redo er editorhistorikk, ikke varig Feature-provenance.

## Illustrativt eksempel

Illustrativ elementliste for et grafisk språk; ikke ArchiMate-DSL-syntaks. Eksemplet er ikke parser-/runtime-testet.

```text
Application Component: Domain service
Application Interface: Measurement input
Application Service: Provide measurements
View: komponenten, grensesnittet og tjenesten som ett utsnitt
```

## Verktøy, vedlikehold og vilkår

Archi-manualen som ble lest er versjon 5.10.0. The Open Groups spesifikasjonssider lot seg ikke hente; normativ versjons-/lisensvurdering står åpen. Archi som verktøy og ArchiMate som standard har forskjellige bruksvilkår.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Archi-utviklernes brukerhåndbok](https://www.archimatetool.com/downloads/archi/Archi%20User%20Guide.pdf)
