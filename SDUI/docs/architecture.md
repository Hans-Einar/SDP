# SDUI — kjørbar frontend, AST og konsolldump

Oppdatert 2026-09-21. Python 3.11+, kun standardbibliotek. Én aktiv frontend for
SDUI 0.2. Go-port og runtime følger [målarkitekturen](target-architecture.md); Go-katalogene har ennå ingen kode.

| Fil | Ansvar |
| --- | --- |
| src/sdui/lexer.py | Begrenset tokenisering, strenger og UTF-8-kildeposisjoner |
| src/sdui/parser.py | Recursive descent fra 0.2-kilde til immutable AST |
| src/sdui/ast.py | Dataclasses, spans, diagnoser og tagget JSON |
| src/sdui/formatting.py | Typede formateringsregler og kanoniske egenskaper |
| src/sdui/validate.py | Lokale navn, widgets, regioner, sykluser og symbolske koblinger |
| src/sdui/normalize.py | Begrenset ekspansjon av komponentreferanser til instanstre |
| src/sdui/dump.py | Strukturell terminalvisning fra instanstre; ingen SDUI-parser |
| src/sdui/markdown_dump.py | Statisk Markdown med layoutoversikt og renderbart innhold fra samme instanstre |
| src/sdui/prototype_widgets.py | Begrenset widgetbeskrivelse fra instanser, SVG-tegning og HTML-kontroller |
| src/sdui/prototype_html.py | Lokalt HTML-galleri med felt, trykkrespons og utskriftsverdier |
| tools/build_widget_previews.py | Avgrenset referansekomposisjon med lagret treemap-geometri; ikke generell layout |
| src/sdui/__main__.py | Fil/stdin, AST/dump/markdown, utdata og exitkoder |

## AST og modell

CLI-envelope er `astFormat: sdui-ast/0.2`, validation og document. Dataclasses
har eksplisitt type-tag; tuples blir JSON-arrays. Document har profile,
references, definitions, connections og span. Definition har name/root/span.
Node beskriver kind (frame/group/markdown/widget/use), name/role, rows,
widget/arguments, text, target, variant, layout og span. Row bevarer items.
LayoutRule bevarer navn, Literal-verdi, opprinnelig spelling og span.
Connection lagrer module/object/definition og path som sekvens av navnesegmenter.
Alle AST-noder er immutable. Kommentarer/quotevalg bevares ikke som concrete syntax tree.

Span.start/end er halvlukkede UTF-8-byteområder; line/column er ettbasert.
Komponentens span inkluderer navn, variant og formatering. En rå Markdown-literal
har kildeområde inkludert delimiter. Parseren tolker ikke innholdets diagrammer.

Instance er en separat statisk modell med path, kind, layout, rows og regions.
Referanser ekspanderes med separat identitet og kildegrense; eksplisitte grupper
bevares. Anonyme segmenter har `$r…c…`-navn og utelates i offentlige setHandle-baner.
Formatering på referanser overstyrer samme egenskap på definisjonens rot, men
motstridende aksevalg avvises. Barnets egen font er bevart; fontarv/måling utføres
ikke av denne strukturmodellen. Ingen levende handles, events eller state finnes.

## Kall og feilveier

CLI leser begrensede bytes, dekoder UTF-8 og kaller parse → validate. Validering
kaller resolve_all etter lokal kontroll og sykluskontroll, med total ekspansjonsgrense.
Offentlig normalize kaller validate og returnerer instanstrær. CLI/dump kan bruke
resolve_all etter utført validering. Dette er gjenbruk av samme frontend.

AST serialiseres til JSON. Dump velger eksplisitt inngangsframe og kaller gui_dump.
Ingen SDL-fil åpnes; ingen callback eller SVG-produsent kjøres. Feil skrives som
JSON på stderr (exit 2 for kilde/profil/dump, 3 for I/O); ingen delvis stdout.
Utdata kan ikke skrives over kildefilen. Dumpen krever validert profil.

Dumpen bruker terminalceller og ASCII-bokskanter. Rader vises horisontalt og
vertikalt, med horisontal fr/scale-fordeling og ellers lik spaltefordeling.
Bare `*box`/`*b` får synlig ramme. Header/footer omslutter body; høyden er
innholdsdrevet. Piler, vertikale vekter, ratio, font, gap, padding, wrap-/scroll-
policy og nøyaktig contentmåling er bevart i modellen, men ikke utført av dumpen.
Nye tekstlinjer i dumpen er terminalombryting av rå Markdown, ikke Markdown-rendering.
Mermaid-fences utelates uten renderer. Kontrolltegn escapes før terminalutskrift.
For liten bredde gir dump-space; bredden er 20–400 celler og samlet arbeid
begrenses av et cellebudsjett. Kombinerende tegn/CJK håndteres; full emoji-
grafemsegmentering er ikke implementert. Ingen ANSI eller interaktivitet.

Denne dumpen erstatter ikke G2s felles layoutmotor. En senere Fyne-konsument
skal bruke felles modell/runtime og målekontrakt, ikke kopiere parseren.

Markdown-eksporten gjenbruker gui_dump som inngjerdet oversikt og traverserer
samme Instance-tre for innholdsdelen. Markdown-strenger blir Markdown-blokker;
frame-/gruppestruktur blir nestede sitatblokker med region- og radetiketter.
Knapper/input/SVG blir statiske etiketter med kodeformaterte literalverdier.
Ingen HTML-widgets, callbacks eller egen layoutmotor genereres.
Vanlige Mermaid-fences filtreres som i tekstdumpen. Åpne øvrige kodegjerder
lukkes ved komponentgrensen, og oversiktens ytre gjerde velges langt nok til å
inneholde kildens backticks. Eksporten har et eget budsjett på 2 MB UTF-8.
Innhold bruker vertens Markdown-profil; tabeller krever GFM-lignende støtte.
Ingen Markdown-parser er lagt til som produksjonsavhengighet.

Den separate [widgetprøven](prototype-widgets.md) bruker parse → normalize →
prototype_widgets → SVG/HTML. Den tegner eksisterende button/input uten ny parser.
Lokalt HTML-state demonstrerer bare kontrollene; callbacks, SDL, ABI og FOX er
ikke koblet til. Referansekomposisjonens plassering erstattes av G2s felles geometri.
