# Concept1 i SDUI og konsollpresentasjon

2026-09-21. [Kjørbar kilde](../examples/concept1-bucking.sdui),
[AST](../examples/concept1-bucking.ast.json), [dump](../examples/concept1-bucking.dump.txt).
Kildegrunnlaget er Ponsse commit 882ad7c0457a723550ea6cce5bf027230997ec2a;
[manifestet](../evidence/concept1-source-manifest.json) registrerer faktiske filhashverdier.
Det finnes andre lokale Ponsse-endringer; de undersøkte UI-/schemafilene er uendret.

## Hva som er portert

| Concept1 | SDUI |
| --- | --- |
| operator-layout.mjs: top, weight 15 | bucking/top, y=15fr; length/diameter x=1fr |
| middle, weight 45 | bucking/middle, y=45fr; selection/suggestions/currentStem 25/50/25fr |
| stemTrack, weight 40 | bucking/stemTrack, y=40fr |
| OperatorBuckingView: MetricBody og MetricSourceControls | Markdown-måleverdi og navngitte Cursor-/Δ-knapper |
| SelectionBody | Markdown med treslag/sortiment/siste knapp |
| SuggestionRows, planvalg, innstillinger | Rå Markdown-tabell og statiske knapper for Canonical/alternativ/kolonner/optimalisering |
| CurrentStemDefectsBody, OperatorActualStemPanel | Representative fakta og tydelig plassholder for produksjon/feilflyt |
| StemTrackHeaderContent/FooterContent | Grupper med forklaring, modellknapper og mål i header/footer |
| OperatorUnifiedStemSvg | Tydelig Markdown-plassholder; illustrativ Mermaid-fence tester utelating |
| OperatorApp header/nav | Separat navigation-definisjon og page-wrapper |

Dette er port av struktur og representativt innhold, ikke en full erstatning for
React-appens funksjoner. Cursor, barkkalibrering, planvalg, redigering, tabs,
tilstandsavhengig innhold og domeneoperasjoner er ikke implementert.
Eksempelverdiene er valgt for lesbar prototyping, ikke hentet fra en levende maskin.
Knappene har ingen callbacks; parseren og dumpen har ingen sideeffekter mot Ponsse.

Originalen bruker React, CSS Grid-fr og en samlet viewport-transform. SDUI beholder
vekter og struktur, men skal ikke skalere font/innhold ved vindusresize. Root med
`{16:9,<->}` skal avlede høyden fra bredden. AST-et registrerer dette;
konsolldumpen beregner ikke denne geometrien. Derfor er den ikke en skjermgeometriprøve.

## Dumpkontrakt

Dump bruker validert AST → samme normaliserte instanstre → terminalrader.
Ingen separat parser og ingen Go-/GUI-avhengighet for statisk dump.
`--entry bucking` viser selve apteringsflaten; `--entry page` tar med appskallet.
`--columns` angir terminalceller, ikke kildepiksler eller fontstørrelse.

Rå Markdown beholder markører som ##, ** og tabellstreker. Lange kildelinjer
brytes ved terminalbredden; dumpen er dermed ikke en kildefil som kan rundtrippes.
Vanlige backtick-/tilde-fences med infoord mermaid utelates med en markør;
innhold inne i andre kodegjerder forblir bokstavelig kode. Ingen diagrammotor kjøres.
Dette er en enkel fence-skanner, ikke en full Markdown-parser for nestede liste-/blockquote-containere.
SVG-referanser vises som plassholdere. Terminalkontrolltegn escapes. Ved for
liten bredde gis diagnose i stedet for overlapp eller stille tap av bokser.
Full emoji-grafemhåndtering og komplekse skriftsystemer krever senere tekstmåling.

Den statiske dumpen bruker naturlige radhøyder og horisontal fr/scale-fordeling;
den viser ikke vertikal fr-fordeling, fontstørrelser eller alle layoutregler.
G2 må levere én felles målt layoutmotor før GUI-/TUI-geometri kan sammenlignes.

## Hva gh-tree bruker

Den nye Markdown-dumpen er beskrevet nederst. Den er en statisk eksport og
forutsetter ingen av TUI-leveransene nedenfor.

Kontrollert i Hans-Einar/gh-tree på commit 97cc0d8257603766dd741b49b7d8005857b421a9:
Go 1.25.0, Bubble Tea v1.3.10 og Lip Gloss v1.1.0 som direkte avhengigheter.
Se [versjonsfestet go.mod](https://github.com/Hans-Einar/gh-tree/blob/97cc0d8257603766dd741b49b7d8005857b421a9/go.mod).
Bubble Tea er TUI-rammeverket; Lip Gloss brukes til terminalstil/layout.
Dette er gh-trees faktiske pin, ikke et råd om å fryse nye SDUI-avhengigheter til samme versjoner.

Interaktiv konsoll er utsatt etter Go/Fyne-valget 2026-09-21. Første vert er
Fyne; en eventuell senere Bubble Tea-vert bruker samme Go-modell/runtime direkte,
ikke cgo/C-ABI eller en egen parser. Historiske T1–T4 er ikke aktive milepæler;
se [G1–G5-planen](implementation-plan.md). Den statiske konsolldumpen er levert.

## Statisk Markdown-dump

`--format markdown --entry bucking` skriver et Markdown-dokument til stdout;
`-o examples/concept1-bucking.dump.md` lagrer det. Se
[generert dokument](../examples/concept1-bucking.dump.md).

Først kommer layoutoversikten som tekstgjerde, deretter innholdet som faktisk
Markdown: overskrifter, fet tekst, lister, tabeller og kodeblokker. Grupper/frames
vises med nestede sitatblokker, regioner og radetiketter. Horisontale søsken
gjengis etter hverandre i innholdsdelen, mens oversikten viser kolonneplasseringen.
Dette gir et portabelt dokument uten HTML/CSS-krav. Tabeller forutsetter en leser
med tabellutvidelsen. Knapper/input er statiske etiketter; ingen hendelser kjøres.

Mermaid-diagrammer utelates med samme fence-skanner som før. Ingen ny diagrammotor
eller SDUI-parser. Relative lenker i brukerens Markdown bevares; ved flytting av
dumpen til en annen katalog må dens innholdsressurser fortsatt være tilgjengelige.

Mermaid treemap støtter tekstetiketter på seksjoner og blader, og bladverdier
styrer areal. Full Markdown som tabeller eller flere innholdsblokker inne i en
boks er ikke en dokumentert treemap-kontrakt. Den har også egen arealfordeling,
som ikke generelt følger SDUIs eksplisitte rader/kolonner.
Se [offisiell treemap-syntaks](https://mermaid.js.org/syntax/treemap).

**Korrigering etter visuell tilbakemelding:** Markdown-rapporten ovenfor løser
ikke eierens mål om en samlet GUI-layout. En [forenklet treemap-prøve](../examples/concept1-bucking.treemap.md)
viser at den lokale Rust-rendererens vekslende akser kan gjengi akkurat
Concept1s hovedrader og kolonner når hierarkiet bygges for det. Prøven har
forkortede tekstetiketter og forhåndsrendret SVG. Den er visuelt kontrollert,
men er ingen generell SDUI-eksport eller full Markdown-i-boks-implementasjon.
[Undersøkelse og begrensninger](../evidence/treemap-probe/README.md).
