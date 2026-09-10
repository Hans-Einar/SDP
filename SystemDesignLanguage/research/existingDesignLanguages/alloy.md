# Alloy 6

[Katalog](README.md) · Kategori: **Relasjonelt formelt modelleringsspråk** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Strukturer, relasjoner og constraints som en solver kan undersøke. Alloy 6 inkluderer mutable elementer og temporal logikk.

## Modellmekanismer

Identitet: signatures/atoms innen analysen. Relasjoner: typede felt og logikk. Contracts: facts, predicates og assertions. Views: instanser og moteksempler. Utvidelse: moduler; maskinell analyse av eksplisitt scope.

## Styrker og begrensninger — vår vurdering

**Styrke:** Kan finne moteksempler til antakelser som at et dependency-graph er lagdelt eller at en mapping dekker alle nødvendige ansvar.

**Begrensning:** Solver-atomer er ikke i seg selv varige FEAT-ID-er. Et funn uten moteksempel må tolkes med analysens scope, tidsgrenser og modellforutsetninger.

## Historikk, endring og transitions

Temporal støtte kan modellere en overgang, men er ikke innebygget versionskontroll. Bevis på bevarte egenskaper krever at de faktisk er uttrykt i modellen.

## Illustrativt eksempel

Illustrativ sykluskontrollmodell; ikke kjørt. Fire er analysescope, ikke en arkitekturregel. Eksemplet er ikke parser-/runtime-testet.

```alloy
sig Layer { depends: set Layer }
fact Acyclic { no l: Layer | l in l.^depends }
run {} for 4
```

## Verktøy, vedlikehold og vilkår

Offisielle Alloy 6-notater og Analyzer er tilgjengelige. Repositorymetadata ga NOASSERTION for lisens; konkrete LICENSE-/NOTICE-filer må kontrolleres før distribusjon.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [Alloy 6](https://alloytools.org/alloy6.html)
- [Offisielt repository; metadata kontrollert via GitHub API](https://github.com/alloytools/org.alloytools.alloy)
