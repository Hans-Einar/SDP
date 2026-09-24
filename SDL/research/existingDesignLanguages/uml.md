# UML 2.5.1

[Katalog](README.md) · Kategori: **Modell-/designspråk** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Generell struktur og atferd: komponenter, klasser, interfaces, aktiviteter, sekvenser og tilstandsmaskiner. OMG publiserer metamodel og XMI-artefakter; UML er mer enn klassediagrammer.

## Modellmekanismer

Identitet: modelelementer/XMI-ID. Relasjoner: typede dependencies, association og realization. Contracts: interfaces og constraints, gjerne OCL. Viewpoints: flere diagramtyper. Utvidelse: profiles/stereotypes; maskinell utveksling via XMI.

## Styrker og begrensninger — vår vurdering

**Styrke:** Bred semantikk og flere perspektiver på samme modell. Interfaces kan beskrives uten å kreve objektorientert implementasjon.

**Begrensning:** Store modeller og verktøyspesifikke profiler kan gjøre oversikten og utvekslingen krevende. Et diagram alene sikrer ikke konsistens mellom views.

## Historikk, endring og transitions

Tilstandsmaskiner beskriver kjøreatferd, ikke modellrevisjoner. Baselines, semantiske differanser og split/merge-sporing trenger repository-/modellverktøy og eksplisitte identitetsregler.

## Illustrativt eksempel

Pseudonotasjon for UML-elementer, ikke en standardisert tekstlig UML-grammatikk. Eksemplet er ikke parser-/runtime-testet.

```text
Class Measurement
  property value : Real
  property unit : String
Component Domain realizes interface MeasurementSource
Component Representation uses interface MeasurementSource
```

## Verktøy, vedlikehold og vilkår

OMG-listingen angir 2.5.1 (desember 2017). Normativ spesifikasjon er tilgjengelig; dette sier ikke at alle verktøy støtter hele standarden. Standarddokumentets bruksvilkår og modellverktøyets lisens må vurderes separat; ingen konkret editor er valgt.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OMG UML spesifikasjon og maskinlesbare artefakter](https://www.omg.org/spec/UML/)
