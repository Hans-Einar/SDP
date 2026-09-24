# TLA+ og PlusCal

[Katalog](README.md) · Kategori: **Formell atferdsspesifikasjon** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Beskriver systemers mulige tilstander og steg med matematisk presisjon. PlusCal er en algoritmenotasjon som oversettes til TLA+.

## Modellmekanismer

Identitet: moduler, konstanter, variabler og operators. Relasjoner/contracts: invariants og temporale properties. Views: utforskede traces/counterexamples gjennom verktøy, ikke ferdige arkitekturviews. Moduler og operators gir gjenbruk.

## Styrker og begrensninger — vår vurdering

**Styrke:** Relevant for køer, samtidighet, leveringsgarantier og om en overgang bevarer et uttrykt invariant.

**Begrensning:** Krever riktig abstraksjon og ekspertise. TLC-søk i et konfigurert endelig rom er ikke et generelt bevis på programkoden; modell og kode må kobles gjennom egne evidenskrav.

## Historikk, endring og transitions

Temporal logikk beskriver systematferd over tid. Refinement kan sammenligne abstraksjoner, men gir ikke automatisk FEAT-ID-historikk, Git-migrasjon eller godkjenning av et refactorløp.

## Illustrativt eksempel

Illustrativ spesifikasjon; mangler TLC-konfigurasjon og er ikke modellkontrollert her. Eksemplet er ikke parser-/runtime-testet.

```text
---- MODULE Counter ----
EXTENDS Naturals
VARIABLE n
Init == n = 0
Next == n' = n + 1
Spec == Init /\ [][Next]_n
====
```

## Verktøy, vedlikehold og vilkår

Forfatterens dokumentasjon og TLA+-verktøyrepository er tilgjengelige. Metadata for tlaplus/tlaplus: MIT, ikke arkivert. Velg eksplisitt mellom modellkontroll og eventuelt teorembevis i en pilot.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Leslie Lamports TLA+-side](https://lamport.azurewebsites.net/tla/tla.html)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/tlaplus/tlaplus)
