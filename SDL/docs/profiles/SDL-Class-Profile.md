# SDL class-core 0.1

Eksplisitt G6-M6-profil for klassedesign; ingen design-core-fakta omtolkes.
Felles SDL-lexer/kildeposisjoner brukes. Alle deklarasjoner kommer før fakta,
sortert etter navn; fakta sorteres etter sin kanoniske tekst.

`class` deklarerer en klasse. `attribute Name as Type` og
`operation Name returns Type` angir egenskaper og parameterløse operasjonssignaturer.
Type er text/integer/boolean eller navnet på en deklarert klasse. Dette er
signaturer, ikke implementerte funksjoner eller en kjørbar objektruntime.

`association` deklarerer et forhold med to eksplisitte klasser, rollenavn og
multiplisitet ved begge ender. `multiplicity 0 to many` betyr 0..*, `1 to 1`
betyr 1. Begge tall må være 0–65536; many er eneste åpne øvre grense.
Endene uttrykker antall objekter på den enden per objekt på motsatt ende.

Hver association må ha nøyaktig én `links` og én `ownership`:

- none: vanlig assosiasjon, ingen eierskaps-/livstidsbetydning.
- aggregation: første ende er helheten; delt del, uten påstått kaskadesletting.
- composition: første ende er helheten; delen har høyst én helhet samtidig,
  og inngår i helhetens livstid. Første endes øvre multiplisitet må derfor være 1.

Profilen krever også asyklisk composition mellom klassetypene. Dette er en
streng, avgrenset modellregel; den påstår ikke å dekke alle UML-modeller.
Roller må være forskjellige og unike for en klasses assosiasjonsender.
Egenskaper/operasjoner har unike navn innen hvert compartment. Maks 128 klasser,
256 assosiasjoner og 64 members per klasse. Udeklarerte typer, feil intervaller,
gjentatte/ufullstendige relasjoner og ubegrunnet composition avvises.

Mermaid classDiagram bruker hul diamant for aggregation og fylt diamant for
composition ved første ende. Vanlig association bruker strek. Begge rollenavn
og begge multiplisiteter følger diagrammet. Ingen contains/owns/allocation
oversettes hit. Klasseprofiler blir ikke automatisk del av action-core-runtime.
