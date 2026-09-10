# LikeC4

[Katalog](README.md) · Kategori: **Utvidbar arkitektur-DSL** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Tekstlig arkitekturmodell med egne elementtyper og relasjoner. Flere kildefiler samles til en modell, og views velger hva som vises.

## Modellmekanismer

Identitet: kvalifiserte elementreferanser. Relasjoner: eksplisitte forbindelser. Contracts: kan refereres/modelleres, men API-konformitet følger ikke av en pil. Viewpoints: statiske/dynamiske views og predicates. Utvidelse: specification med egne kinds, metadata og generators/API.

## Styrker og begrensninger — vår vurdering

**Styrke:** Egne kinds kan uttrykke SDP-begreper uten å låse alt til C4s opprinnelige nivåer. Programmatisk modelltilgang er interessant for Analyzer.

**Begrensning:** Et element kalt feature får ikke automatisk krav-, release- eller transition-semantikk. Fleksibiliteten trenger en validert profil.

## Historikk, endring og transitions

Modellutvidelse mellom filer er komposisjon, ikke historikk. Git og en separat før/etter-mapping må bevare identiteter ved rename, split og merge.

## Illustrativt eksempel

Illustrativ arkitektur med en egen elementtype. Eksemplet er ikke parser-/runtime-testet.

```text
specification {
  element service
}
model {
  source = service 'Source'
  sink = service 'Sink'
  source -> sink 'Measurements'
}
views {
  view overview { include * }
}
```

## Verktøy, vedlikehold og vilkår

Dokumentasjonen viser CLI, editorstøtte, API, generators og releases. Repositorymetadata: ikke arkivert; MIT. Nøyaktig versjon og API-stabilitet må pinnes i en eventuell pilot.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [DSL introduksjon](https://likec4.dev/dsl/intro/)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/likec4/likec4)
