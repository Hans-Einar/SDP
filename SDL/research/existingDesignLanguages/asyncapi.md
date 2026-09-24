# AsyncAPI 3.0

[Katalog](README.md) · Kategori: **Event-/meldingskontrakt** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Beskriver meldingsorienterte API-er med channels, messages, operations og protokollbindings.

## Modellmekanismer

Identitet: component-/operation-navn og references. Relasjoner: operation→channel→message. Contracts: payload og bindings. Views: generert dokumentasjon og kode via verktøy. Utvidelse: extensions og eksterne schemaformater.

## Styrker og begrensninger — vår vurdering

**Styrke:** Nært brukerens input/output-events og abonnenter. Kan gjøre meldingsgrenser eksplisitte uavhengig av renderer.

**Begrensning:** En send/receive-definisjon garanterer ikke bestilling, levering, backpressure eller at alle runtime-events håndteres; dette må spesifiseres og testes.

## Historikk, endring og transitions

Dokumentversjoner og meldingskompatibilitet er forskjellige. Endringer trenger consumer-impact, migrasjonsplan og historikk utenfor enkeltkontrakten.

## Illustrativt eksempel

Illustrativ minimal meldingskontrakt; ingen valgt broker eller leveringsgaranti. Eksemplet er ikke parser-/runtime-testet.

```yaml
asyncapi: 3.0.0
info:
  title: Measurements
  version: 1.0.0
channels:
  readings:
    address: measurements
    messages:
      measured:
        payload:
          type: number
operations:
  emitReading:
    action: send
    channel:
      $ref: "#/channels/readings"
```

## Verktøy, vedlikehold og vilkår

Offisiell 3.0.0-spesifikasjon er tilgjengelig. Verktøyversjoner og binding-version må pinnes i en pilot. Spesifikasjonsvilkår og generator-/runtime-lisens må vurderes separat.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [AsyncAPI 3.0.0](https://www.asyncapi.com/docs/reference/specification/v3.0.0)
