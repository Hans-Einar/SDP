# BPMN 2.0.2

[Katalog](README.md) · Kategori: **Prosessmodell og notasjon** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Forretningsprosesser, oppgaver, hendelser, gateways og samarbeid mellom deltakere. Kan beskrive UseCase-forløp og menneskelige godkjenninger.

## Modellmekanismer

Identitet: element-ID-er i XML. Relasjoner: sequence/message flows. Contracts: melding-/deltakergrenser, men payload og software-importregler krever mer. Views: prosess/choreography/collaboration. Extension elements og maskinlesbare XSD/CMOF; eksekvering er motoravhengig.

## Styrker og begrensninger — vår vurdering

**Styrke:** God til å gjøre ansvar og beslutningspunkter synlige for eier og andre stakeholders.

**Begrensning:** Prosessflyt er ikke et arkitekturkart eller en feature-pathway gjennom kode. Ikke alle gyldige diagrammer er eksekverbare.

## Historikk, endring og transitions

Runtime-prosessinstanser og motorens prosessversjoner er ikke modellens generelle overgangssemantikk. Flytting av Feature-ansvar trenger eksterne modellkoblinger.

## Illustrativt eksempel

Illustrativ prosess uten diagram-layout eller motorkonfigurasjon. Eksemplet er ikke parser-/runtime-testet.

```xml
<definitions xmlns="http://www.omg.org/spec/BPMN/20100524/MODEL"
 targetNamespace="urn:example:sdp">
  <process id="review" isExecutable="false">
    <startEvent id="start"/>
    <task id="inspect" name="Inspect design"/>
    <sequenceFlow id="f1" sourceRef="start" targetRef="inspect"/>
  </process>
</definitions>
```

## Verktøy, vedlikehold og vilkår

OMG lister 2.0.2 fra januar 2014 med normative schemafiler. Standardvilkår og valgt editor/runtime-lisens er separate; ingen motor er evaluert.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OMG BPMN](https://www.omg.org/spec/BPMN/)
