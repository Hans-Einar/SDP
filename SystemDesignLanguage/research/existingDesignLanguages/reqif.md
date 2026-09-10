# ReqIF 1.2

[Katalog](README.md) · Kategori: **Kravutveksling (tilgrensende format)** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Utveksler krav med typer, attributter, hierarkier og relasjoner mellom kravverktøy. Relevant for stakeholders→UseCases→REQ-sporet, ikke en komplett systemarkitektur.

## Modellmekanismer

Identitet: IDENTIFIER for elementer. Relasjoner: SpecRelations og strukturer. Contracts: krav kan beskrive dem; ingen API-eksekveringssemantikk. Utvidelse: typer/attributter og tool extensions. Maskinell XML-utveksling; views avhenger av kravverktøy.

## Styrker og begrensninger — vår vurdering

**Styrke:** Et eksisterende alternativ til å finne opp et nytt interchange-format for krav og identitet.

**Begrensning:** Et importert krav dokumenterer ikke at en Feature oppfyller kravet. Round-trip kan være verktøyavhengig og må testes.

## Historikk, endring og transitions

Metadata som LAST-CHANGE er ikke et revisjonsarkiv. Baselines, endringsgodkjenning og kobling til FEAT/releases må etableres separat.

## Illustrativt eksempel

Illustrativt XML-fragment, ikke et komplett schema-validert ReqIF-dokument. Eksemplet er ikke parser-/runtime-testet.

```xml
<SPEC-OBJECT IDENTIFIER="REQ-1"
 LONG-NAME="Display measured length"
 LAST-CHANGE="2026-09-10T00:00:00Z"/>
<!-- Fragment: type og øvrig ReqIF-dokument er utelatt. -->
```

## Verktøy, vedlikehold og vilkår

OMG publiserer ReqIF 1.2. Standardvilkår er separate fra importør-/eksportørlisenser. Ingen konkret round-trip eller editor er evaluert.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OMG ReqIF](https://www.omg.org/spec/ReqIF/)
