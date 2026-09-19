# QVT 1.3

[Katalog](README.md) · Kategori: **Query/View/Transformation-standard** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

OMG-standardfamilie for operasjoner på modeller: Relations, Core og Operational Mappings. Kan beskrive konsistens og transformasjoner mellom metamodeller.

## Modellmekanismer

Identitet/relasjoner: modelelementer i deltagende modeller. Contracts: constraints og transformasjonsbetingelser. Views: avledede modeller; visuell rendering er separat. Maskinell utføring krever en implementasjon av valgt del av standarden.

## Styrker og begrensninger — vår vurdering

**Styrke:** Relevant når før/etter-korrespondanser skal uttrykkes eksplisitt og ikke bare tegnes.

**Begrensning:** QVT-delene har ulik operasjonell betydning og verktøystøtte. En deklarert relation er ikke automatisk en entydig, tapsfri toveis migrasjon.

## Historikk, endring og transitions

QVT kan uttrykke en mapping; et varig revisjonsarkiv og audit av hvilket regelsett som ble brukt må håndteres utenfor selve transformasjonen.

## Illustrativt eksempel

Illustrativ QVT Relations-skisse med hypotetiske metamodeller; implementasjonsdialekt og navneoppløsning må valideres før kjøring. Eksemplet er ikke parser-/runtime-testet.

```text
transformation Preserve(source : Old, target : New) {
  top relation PreserveIdentity {
    fid : String;
    checkonly domain source a : Feature { id = fid };
    enforce domain target b : Feature { id = fid };
  }
}
```

## Verktøy, vedlikehold og vilkår

OMG publiserer 1.3. Standardens vilkår og QVTo/QVTr-implementasjonenes lisenser/støtte må undersøkes hver for seg. Ingen motor er valgt.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OMG QVT](https://www.omg.org/spec/QVT/)
