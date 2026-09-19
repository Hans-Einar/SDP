# Structurizr DSL / C4

[Katalog](README.md) · Kategori: **Arkitekturmodell som tekst** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Definerer én arkitekturmodell med flere C4-diagrammer: systemer, containere, komponenter, deployment og dynamiske forløp. C4 er abstraksjonsmodellen; Structurizr DSL er et konkret språk.

## Modellmekanismer

Identitet: DSL-identifiers og modelelementer. Relasjoner: modellerte forbindelser. Contracts: beskrivelser/metadata og eksterne lenker, ikke full API-semantikk. Views: sterke eksplisitte utsnitt. Utvidelse: properties, tags, includes, workspace extension og plugins.

## Styrker og begrensninger — vår vurdering

**Styrke:** God kandidat til lesbare eier-viewpoints og featurefiltrerte fremstillinger over en felles modell.

**Begrensning:** FEAT→REQ, fences og semantisk kompatibilitet må spesifiseres utenfor standard C4-begrepene. DSL-navn må ikke forveksles med en stabil livsløps-ID.

## Historikk, endring og transitions

Git kan versjonere DSL. Serverdokumentasjonen har workspace-versjoner/branches, men dette er verktøystøtte; en flytting av ansvar er ikke en innebygget, verifisert Feature-transition.

## Illustrativt eksempel

Illustrativ variant av dokumentasjonens workspace/model/views-struktur. Eksemplet er ikke parser-/runtime-testet.

```text
workspace {
  model {
    owner = person "Owner"
    app = softwareSystem "Monitor"
    owner -> app "Reads measurements"
  }
  views {
    systemContext app { include * }
  }
}
```

## Verktøy, vedlikehold og vilkår

Nåværende dokumentasjon viser en samlet kommandoflate og markerer eldre Lite/CLI/cloud under end-of-life. Ikke bygg neste pilot på antakelser om gammel produktpakking. Dokumentasjonen sier at prebygget server krever lisens; øvrige kommandoer er gratis å bruke. Dette er ikke en generell lisensgodkjenning.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [DSL-eksempel](https://docs.structurizr.com/dsl/example)
- [Nåværende verktøy og bruksvilkår](https://docs.structurizr.com/)
