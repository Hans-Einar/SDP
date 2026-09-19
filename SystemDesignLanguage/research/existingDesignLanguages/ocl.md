# Object Constraint Language (OCL) 2.4

[Katalog](README.md) · Kategori: **Modellconstraints og queries** · Research: **2026-09-10**

## Formål og abstraksjonsnivå

Deklarative uttrykk over modeller: invariants, preconditions, postconditions og navigering. Komplementerer blant annet UML.

## Modellmekanismer

Identitet/relasjoner: arves fra vertsmodellen. Contracts: presise logiske betingelser. Viewpoints: queries kan velge elementer, men rendering er eksternt. Maskinell evaluering avhenger av metamodel, standardbibliotek og implementasjon.

## Styrker og begrensninger — vår vurdering

**Styrke:** Et mulig middel for å uttrykke at en modell ikke tillater bestemte dependencies eller manglende kontraktskoblinger.

**Begrensning:** Gir ikke selv en systemmodel eller kildekodeanalyse. Navnet på en constraint er ikke bevis på at den tilsvarende egenskapen er riktig formalisert.

## Historikk, endring og transitions

OCL er ikke et migrasjonsspråk. Constraints kan sjekke før/etter-modeller dersom begge representeres; koblingen mellom dem må etableres separat.

## Illustrativt eksempel

Illustrativ OCL som forutsetter en hypotetisk metamodel med Feature.requirements; dette oppretter ingen normativ SDP-regel. Eksemplet er ikke parser-/runtime-testet.

```ocl
context Feature
inv HasRequirement: self.requirements->notEmpty()
```

## Verktøy, vedlikehold og vilkår

OMG publiserer OCL 2.4. Spesifikasjon og konkret OCL-motor har separate vilkår. Parser-/dialektkompatibilitet er ikke testet.

## Primærkilder

Alle kilder kontrollert 2026-09-10; se katalogens metode for evidens- og lisensbegrensninger.

- [OMG OCL](https://www.omg.org/spec/OCL/)
