# SysML v2

[Katalog](README.md) · Kategori: **Systemmodelleringsspråk** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Systemstruktur, behavior, requirements, forbindelser og analyser på tvers av software og fysiske systemer. V2 har tekstlig og grafisk notasjon og bygger på KerML; det er ikke bare en ny UML-profile.

## Modellmekanismer

Identitet: modelelementer og kvalifiserte navn. Relasjoner: blant annet parts, ports, connections og requirement-sammenhenger. Contracts: definerte grensesnitt og constraints. Views/viewpoints og biblioteker er del av modelleringsområdet. API/repository-støtte må vurderes som egne spesifikasjoner/verktøy.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kan knytte krav, struktur og atferd tettere sammen enn rene diagram-DSL-er. Tekstformatet er interessant for agentarbeid.

**Begrensning:** Stor semantisk overflate og læringskostnad. SysML/KerML-begrepet feature må ikke automatisk tolkes som SDP FEAT.

## Historikk, endring og transitions

Skill språkets modeller fra repositoryets commits/branches og API-støtte. En før/etter-modell er ikke automatisk en verifisert migrasjon. Langlivet FEAT-identitet og ansvarsflytting trenger en avtalt mapping.

## Illustrativt eksempel

Illustrativ, liten tekstlig struktur; ingen Ponsse-mapping eller requirements-pilot. Eksemplet er ikke parser-/runtime-testet.

```sysml
package Demo {
  part def Sensor;
  part def Monitor {
    part sensor : Sensor;
  }
}
```

## Verktøy, vedlikehold og vilkår

OMG har en formell 2.0-side med publiserte normative dokumenter. Pilotimplementasjon og eksempler finnes i Systems-Modeling-repositoryet. Pilotstøtte er ikke bevis på full standardkonformitet. OMG-dokumentvilkår og LICENSE-filer i den valgte implementasjonen må kontrolleres hver for seg.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OMG SysML 2.0](https://www.omg.org/spec/SysML/2.0/)
- [Offisiell pilot og eksempler](https://github.com/Systems-Modeling/SysML-v2-Release)
