# SDUI — krav og sporbarhet

Oppdatert 2026-09-22. Gjeldende fasebevis er [G1–G3](../go/README.md),
[G4–G6](../../SDL/go/README.md). Tabellen følger Go-porten;
de navngitte Python-baselinetilfellene er fryst som testdata, ikke aktive tester.

| ID | Krav | Eier / bevis | Status |
| --- | --- | --- | --- |
| SDUI-R01 | Navngitte nestede frames og grupper | parser.Node; nested_groups_and_local_rows_preserved | Implementert struktur; tidligere Box-profil erstattet |
| SDUI-R02 | Horisontale elementer og eksplisitte rader | parser.rows; formatting_binds_before_separator | Implementert/testet |
| SDUI-R03 | Symbolske modul-/medlemsreferanser | ast.Reference; ported_example_bindings_and_instance_paths | Implementert; ingen ekstern resolusjon |
| SDUI-R04 | Deklarativ setHandle | ast.Connection, validate; local_symbol_errors | Implementert lokalt |
| SDUI-R05 | Entydige navne-/instansbaner | normalize; references_forward_reuse_and_cycles | Implementert; runtime-handles i G3 |
| SDUI-R06 | Versjonert EBNF/AST med kildeposisjoner | raw_markdown_and_utf8_spans; version_is_exact_no_legacy | Implementert; ingen formell EBNF-ekvivalens |
| SDUI-R07 | Ugyldig profilsemantikk avvises | formatting; shapes_ratio_and_relative_dimensions, widget_contracts | Implementert lokale regler |
| SDUI-R08 | Begrenset arbeid uten kjøring av kilde | bounded_work, truncations_and_malformed_inputs_are_structured | Implementert; ingen hard sanntid |
| SDUI-R09 | Standalone AST uten GUI-/Mermaid-avhengighet | CLI; cli_source_protection_and_exit_codes | Implementert/testet |
| SDUI-R10 | Samme normaliserte UI-modell fra SDUI og Concept1-layout | Concept1-referanse og Go-normalisering/layout | Verifisert referansestruktur; generell React-adapter utenfor profilen |
| SDUI-R11 | Synlig ytre boks, tittel, vekter og widgetlayout i Fyne og SVG/Markdown | Go layout/svg/fynehost, G2 | Implementert/verifisert |
| SDUI-R12 | SDL-objekt kan motta handle og oppdatere text/value | Go runtime/bridge, G3/G4 | Implementert med action-core og typede Go-funksjoner |

Testnavnene i tabellen er forkortet for lesbarhet; filen inneholder fulle navn.
[Historisk 0.1-verifikasjon](../evidence/verification.md) registrerer den utgåtte prototypens kjøring.
Ingen semantisk validator påstår at en ekstern SDL-fil eller funksjon finnes.

## Leveransekrav og fasebevis

| ID | Krav | Fase / bevis |
| --- | --- | --- |
| SDUI-R13 | Generelle frames, header/body/footer, nestede grupper med egen layout og gjenbrukbare grupper med instansbaner | G1/G2: grammatikk, AST og bevart gruppestruktur; lokale separatorer og egen layout på flere nivåer |
| SDUI-R14 | Markdown-streng blir innholdswidget med korrekt kildekart | G1/G2: multiline/escaping, bredde→høyde, innholdsprovider |
| SDUI-R15 | Ancestor-relative størrelser, ratio x:y, én styrende scaleakse med ratio; relative grenser/avstander og kanoniske pilformer | G1/G2: avvis px og to ratioakser; normalisering, geometriorakler og overflow |
| SDUI-R16 | Typed runtime/backend-porter med stabile instanser og atomiske oppdateringer | G3: Go-konsument, fake backend og ekte Fyne-prøver |
| SDUI-R17 | SVG, native kontroller og Markdown bruker samordnet geometri/klipp/print | G2: kontrollkart, transforms og visuell PDF/GUI-verifikasjon |
| SDUI-R18 | SDL kan kobles over eksplisitt adapterport; simulering merkes | G3/G4: toveis korrelasjon, revisjonskonflikt, teardown; ekte SDL-runtime separat |
| SDUI-R19 | Én ny aktiv frontend/API-vei; port og fjern legacy | G1-M3/G5-M4: target-/importinventar, unsupported-version, ingen fallback |
| SDUI-R20 | Markdown i frame gjenbruker vertens avtalte profil og diagramdekning | G2/G3: profilmatrix, lenker/merking, scroll og ressurser |
| SDUI-R21 | Formatering etter hver UI-komponent før separator; absolutt font uten innholdsskalering ved resize | G1/G2/G3: suffix på alle komponenttyper, radstruktur, stabil font og ny tekstombryting ved resize |

R13/R14/R15/R21 er verifisert gjennom frontend, målt layout og native prøver.
R19 er fullført i G5-M4. R20 er avgrenset av markdown-provider.md: flowchart/graph,
ikke alle Mermaid-typer; scroll avvises. XFMDs gamle BoxUI-worktree er bevart som
annet arbeid; den nye dokumentintegrasjonen er en separat konsument i G6.

R10 konkretiseres til samme referansestruktur/geometriregler i første leveranse;
generell Concept1/React-adapter er ikke nødvendig. R11 er revidert 2026-09-21:
Fyne er første interaktive vert, SVG i Markdown er dokumentasjonsvisningen. R12 har ekte SDL action-core-runtime/bridge i G4. Det håndskrevne apteringsdomenet
er fortsatt en tydelig merket simulering.

Eierpresisering 2026-09-20: R15 tillater ikke pikselbredde/-høyde i kilden.
Root bruker vertens layoutområde; barn nærmeste ancestor. Retningsparets
tegnrekkefølge endrer ikke justering, men formattereren bruker én shape.
Forslaget gjør også øvrige lengder relative. Presisering 2026-09-21:
`{16:9,<->}` fyller bredden og avleder høyden, uten contain-fallback.
Header/footer ligger innen ratio. Fontstørrelse er absolutt; resize skalerer ikke
innholdet. Felles fontenhet er logiske DIP med Go Regular-måling.
Formatering `{...}` står etter komponenten og før separator; komma fortsetter
horisontalt, semikolon starter neste rad under foregående rad.
Kanoniske hjørner er ^<, >^, v<, >v; høyre/midt skrives ->.

SDUI-R22: statisk GUI-dump fra felles modell med rå Markdown og utelatte Mermaid-fences.
Implementert/testet i go/presentation. Interaktiv konsoll er planlagt i [TUI-retningen](concept1-console.md).

SDUI-R23: statisk Markdown-dump fra samme modell, med layoutoversikt og renderbare
Markdown-innholdsblokker. Bevar tabeller/lister/kode, marker widgets som statiske,
og utelat Mermaid som i tekstdumpen. Ingen interaktivitet eller ny parser.
Markdownens innholdsdel viser leserekkefølge og nesting, ikke eksakt GUI-geometri.

SDUI-R24: et avgrenset SDUI-prototypebibliotek tegner button/input som statiske
SVG-widgets i Markdown. Samme widgetbeskrivelse kan demonstreres i lokal HTML
med redigering, trykk-/fokusutseende og utskrift av aktuelle feltverdier.
Ingen SDL-kall, nettverk, FOX-avhengighet eller ny kildeparser. Prototypeplassering
og demonstrasjonstilstand må skilles fra den framtidige felles layout/runtime.
Levert som avgrenset prøve; [omfang](prototype-widgets.md) og
[test-/nettleserbevis](../evidence/prototype-widgets/README.md).

SDUI-R25: kildeendringer kan parses/valideres og publiseres uten å lukke UI-vinduet.
Bevar siste gyldige modell ved feil og kompatibel verdi/fokus ved reload; avvis
stale hendelser. SDL-state krever migrerings-/resetregel. Go-funksjonsendringer
bygges/restartes. Implementert/verifisert G3/G4.

SDUI-R26: senere Go-generering bruker samme modell/runtime som utviklingsmodus
og skiller generert kode fra håndskrevne domenefunksjoner. Ukjent eller ufullstendig
SDL-kjøresemantikk gir diagnose. Implementert/verifisert G5.
