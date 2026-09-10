# Protocol Buffers

[Katalog](README.md) · Kategori: **Data-/service-IDL** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Typede meldinger med binær serialisering og kodegenerering på flere implementeringsspråk. Servicebeskrivelser kan brukes av RPC-verktøy.

## Modellmekanismer

Identitet: package, message/type names og felt-numre. Relasjoner: typereferanser og services. Contracts: datatyper/wire-format; semantiske enheter må avtales. Views: eksterne dokumentasjonsverktøy. Utvidelse: options og generatorplugins.

## Styrker og begrensninger — vår vurdering

**Styrke:** Feltidentiteter og dokumenterte evolusjonsregler gir et konkret eksempel på kontrakter som kan endres kontrollert.

**Begrensning:** Wire-kompatibilitet er ikke det samme som meningsbevaring. Å endre meter til centimeter kan være en alvorlig feil uten at typen endres.

## Historikk, endring og transitions

Felt-numre skal ikke gjenbrukes etter sletting; reserved kan beskytte fjernede felt. Dette støtter kontraktevolusjon, men lagrer ikke hele ansvarshistorikken.

## Illustrativt eksempel

Illustrativ proto3-melding. Meter er valgt som eksplisitt konvensjon i feltnavnet; ingen universell SDP-unitregel er bestemt. Eksemplet er ikke parser-/runtime-testet.

```proto
syntax = "proto3";
package measurements.v1;
message Length {
  double value_m = 1;
  reserved 2;
}
```

## Verktøy, vedlikehold og vilkår

Proto3-guiden er tilgjengelig; andre editions finnes og proto3 er ikke brukt som påstand om nyeste språkvariant. Repositoryet er ikke arkivert; metadata ga NOASSERTION, så konkret LICENSE/NOTICE og runtime må kontrolleres.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Proto3 språk- og evolusjonsguide](https://protobuf.dev/programming-guides/proto3/)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/protocolbuffers/protobuf)
