# Epsilon-språkfamilien

[Katalog](README.md) · Kategori: **Modellvalidering, transformasjon og generering** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

EOL er uttrykksgrunnlag for blant annet ETL (transformasjon), EVL (validering), EGL (tekstgenerering), ECL (matching), EML (merge) og Flock (modellmigrasjon).

## Modellmekanismer

Identitet/relasjoner: fra tilknyttede modeller; ECL kan beskrive korrespondanser. Contracts: EVL-constraints. Views: EGL/Picto kan avlede fremstillinger. EMC kobler flere modellformater. Utvidelse: regler, templates og drivers.

## Styrker og begrensninger — vår vurdering

**Styrke:** En direkte relevant verktøyfamilie for å undersøke blueprints, validering og modell-delta i samme økosystem.

**Begrensning:** Flere små språk og drivers gir integrasjonskostnad. En automatisk match mellom objekter må ikke få autoritet som verifisert bevart ansvar.

## Historikk, endring og transitions

Flock håndterer oppdatering ved metamodelendring; ECL/EML kan hjelpe sammenligning/merge. Dette er ikke automatisk releasehistorikk eller bevis på kodekonformitet.

## Illustrativt eksempel

Illustrativ EVL over en hypotetisk Feature-type; ingen vedtatt SDP-regel. Eksemplet er ikke parser-/runtime-testet.

```text
context Feature {
  constraint HasIdentity {
    check: self.id.isDefined()
    message: "Feature needs an identity"
  }
}
```

## Verktøy, vedlikehold og vilkår

Eclipse-siden angir Mature og EPL-2.0. Den lister 2.8 fra 2025-02-20 og en fremtidsdatert 2.9 (2026-10-19); sistnevnte behandles ikke som utgitt per researchedato. Driverkompatibilitet og konkrete bundle-vilkår gjenstår.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Epsilon språk og verktøy](https://eclipse.dev/epsilon/doc/)
- [Epsilon prosjektstatus](https://projects.eclipse.org/projects/modeling.epsilon)
