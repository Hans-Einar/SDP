# JSON Schema 2020-12

[Katalog](README.md) · Kategori: **Datavalidering og schema** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Beskriver hvilke JSON-instansdata som er gyldige, samt annotations. Kan definere formen på et lite SDP-manifest uten å kreve et nytt parser-språk.

## Modellmekanismer

Identitet: $id, anchors og references for schemaressurser. Relasjoner: $ref/$dynamicRef; ikke automatisk foreign-key-kontroll mellom instansobjekter. Contracts: typer og strukturelle constraints. Utvidelse: vocabularies. Views/grafspørringer er eksterne.

## Styrker og begrensninger — vår vurdering

**Styrke:** Mange generiske verktøy kan lese en JSON-basert modell. God inngang til lettvekts validering av obligatoriske felt.

**Begrensning:** Validering av FEAT-ID-referanser, sykluser og laggrenser krever vanligvis ekstra logikk. YAML/JSON er serialisering; semantikken kommer ikke gratis.

## Historikk, endring og transitions

$schema angir dialekt, $id identifiserer schemaressursen. Ingen av dem er en innebygd historikk for modellinstanser eller transition-migrasjon.

## Illustrativt eksempel

Illustrativ ID-formkontroll, ikke vedtatt SDP-schema eller kontroll av ID-unikhet. Eksemplet er ikke parser-/runtime-testet.

```json
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "type": "object",
  "required": ["id"],
  "properties": {"id": {"type": "string", "pattern": "^FEAT-[0-9]+$"}}
}
```

## Verktøy, vedlikehold og vilkår

2020-12-spesifikasjonen og vocabularies er tilgjengelige. Kontroller validatorens dialektstøtte. Spesifikasjonsvilkår og valgt validatorlisens er forskjellige; ingen validator er valgt.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [JSON Schema 2020-12](https://json-schema.org/draft/2020-12)
