# OpenAPI

[Katalog](README.md) · Kategori: **HTTP API-kontrakt** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Maskinlesbar beskrivelse av HTTP-operasjoner, parametere, responses, schemaer og sikkerhetskrav.

## Modellmekanismer

Identitet: paths, operationId og component-referanser. Relasjoner: $ref og API-links. Contracts: HTTP og dataskjema. Views: dokumentasjon/klienter via verktøy. Utvidelse: specification extensions; parser/generator/validator-støtte varierer.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kan være en autoritativ kontrakt på en feature-pathway mellom backend og klient.

**Begrensning:** Beskriver ikke hele domenets atferd, UI-lag eller gjenbrukspolicy. Et schema som validerer er ikke et bevis på bakoverkompatibel brukeropplevelse.

## Historikk, endring og transitions

info.version er API-dokumentets oppgitte versjon, ikke en komplett historikk. Git, kompatibilitetsanalyse og consumer-mapping trengs ved endringer.

## Illustrativt eksempel

Illustrativ 3.1.0-beskrivelse, bevisst ikke et forsøk på å demonstrere nyeste syntaks. Payload utelatt. Eksemplet er ikke parser-/runtime-testet.

```yaml
openapi: 3.1.0
info:
  title: Measurement API
  version: 1.0.0
paths:
  /measurements:
    get:
      operationId: listMeasurements
      responses:
        "200":
          description: Measurements available
```

## Verktøy, vedlikehold og vilkår

Latest-siden viste 3.2.0 (19. september 2025). Verktøystøtte må sjekkes per versjon. Spesifikasjonen angir Apache-2.0; generators og tjenester kan ha andre vilkår.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OpenAPI 3.2.0 spesifikasjon](https://spec.openapis.org/oas/v3.2.0.html)
