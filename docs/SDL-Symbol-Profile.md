# SDL symbolprofil 1

G6-M5: kildeidentitet/type følger hver figur; relasjon og kildefaktum følger
hver kant. Mermaid-rs-renderer eier fortsatt plassering og ruting. SDL-verktøyets
presentasjonsadapter leser dens eksisterende `--dumpLayout`-kontrakt og tegner
faste SVG-figurer fra samme node-/kantgeometri. Ingen parser eller layoutmotor
kopieres inn. Backendens node-ID-er, kanter og geometri kontrolleres før tegning.

Actor: menneskefigur med hode, kropp, armer og ben. UseCase: ellipse.
Feature: kort med fane. Functionality: avrundet kort med sidemarkør.
Capability: sekskant. Activity: avrundet aktivitetskort. Mode: stiplet
kontekstramme. Container: dobbel grense. Unit: enkel grense. Interface: merket
portkort. Database: sylindersymbol for persistent datakilde. Øvrige data-/
kontraktobjekter bruker kort med eksplisitt typeetikett. Farge er ikke typekoden.

`consumes` bruker stiplet dependency med åpen pil mot Interface.
Andre strukturelle relasjoner har merkede SDL-piler. `realizes` betyr fortsatt
bidrag til realisering og får ingen UML-hul trekant. `contains`/`owns` blir aldri
composition. `uses` i Channel-kontrakten er deltakelse/rolle, ikke UML Usage.
Ingen System-grense eller include/extend-relasjon utledes.

SVG er symbolprofilens autoritative bilde. Mermaid-kilden er tilgjengelig som
portabel, merket strukturvisning; den hevder ikke samme figurer i alle lesere.
Sequence og packet bruker eksisterende, prøvde backend-former. Native
`usecase-beta` er ikke valgt: den lokale backendens tidligere prøve ga feil
figurer til tross for exitkode 0. Klassediagram krever egen eksplisitt profil.

Grunnlag: [Mermaid flowchart](https://mermaid.js.org/syntax/flowchart.html),
[Mermaid classDiagram](https://mermaid.js.org/syntax/classDiagram) og
[prosjektets semantiske notasjonsregler](SDL-Viewpoint-Levels-and-Notation.md).
