# SDUI 0.2 — Concept1, AST og konsolldump

2026-09-21. **24 av 24 tester bestod.**
[Full logg](frontend-console-tests-2026-09-21.txt),
[kode-/miljømanifest](frontend-console-manifest.json),
[Concept1-/gh-tree-kildegrunnlag](concept1-source-manifest.json).

## Verifisert

Fra SDUI-katalogen:

```sh
PYTHONPATH=src python3 -m unittest discover -s tests -v
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format dump --entry bucking --columns 160
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format dump --entry page --columns 200
```

- Parser/AST for frames, nestede grupper, lokale rader, regioner og formatering.
- Entydige referanseinstanser, framoverreferanser, sykluskontroll, separat ekspansjonsbudsjett.
- Kildeposisjoner for UTF-8/CRLF/multiline Markdown og bevaring av rått strenginnhold.
- Typede widgets, navngitte callbacks, anonyme ubundne knapper og deklarativ setHandle.
- Kanoniske pilpar, ratio/aksekonflikter, avvisning av absolutte layoutdimensjoner og ugyldig font.
- Lokal separatorvirkning, bevart ettbarnsgruppe, formatering før separator og avvisning av løs blokk.
- Begrenset kilde/token/dybde/noder/argumenter; alle trunkeringer av minimalfixture og 500 deterministiske støystrenger.
- Dump med synlige rader/kolonner, rå Markdown, utelatte Mermaid-fences, escaped terminalkontrolltegn og CJK/combining-cellebredder.
- CLI-entryvalg, ingen delvis feiloutput, JSON-diagnoser, UTF-8-/I/O-/kildevern.
- Seks Concept1-boks-ID-er og kildevekter 15/45/40 og 25/50/25.

Begge lagrede AST-filer (main-page og concept1-bucking) og Concept1-dumpen
ble sammenlignet byte-for-byte med ny CLI-kjøring; alle var identiske.
Page-wrapperen ble også kjørt ved 200 kolonner. Dumpen ved 160 kolonner ble
inspisert i konsoll; rader/bokser og rå tabellkode er synlige uten Mermaid-kjøring.
Gjeldende Markdown-lenker, kodegjerder og `git diff --check` er kontrollert.

## Avgrensning og opprydding

Gammel 0.1-parserlogikk, Box/Content-AST, grammatikk, text-widget og aktive
kildeeksempler er erstattet. Ingen versjonsfallback. Historiske testlogger og
manifest er bevart og tydelig merket som historiske; de gjelder utgått kode.

Dette er Python-frontenden og en statisk strukturmodell, ikke native libsdui,
C-ABI eller runtime. Ingen nye pakkeavhengigheter. Konsolldumpen har naturlige
radhøyder; vertikale fr-vekter, 16:9-geometri, justeringspiler, fontmåling,
full Markdown-rendering og FOX-kontroller er ikke utført/verifisert.
Concept1-verdiene og SVG-/feil-/produksjonsinnhold er markerte eksempeldata/plassholdere.
Ingen SDL-modul, Ponsse-operasjon eller widgetcallback kjøres.

gh-tree er kontrollert via GitHub API: Go, Bubble Tea og Lip Gloss er faktiske
imports/avhengigheter. Interaktiv Go-TUI er planlagt, ikke implementert.
Ingen GUI/PDF-kjøring, endring i eksterne repoer, branchbytte, commit eller push.
