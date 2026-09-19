# SDUI — skjemastyrte brukergrensesnitt

**Arbeidsprofil 0.1 · 2026-09-19 · språk-/parserprototype.** Opprettet etter eierens
bestilling. Dette er søsterprosjektet til [SystemDesignLanguage](../SystemDesignLanguage/README.md)
i Hans-Einar/SDP. Det vedtar ikke nye SDL-regler eller canonical SDP-kontrakter.

SDUI beskriver navngitte bokser, layout, widgets og symbolske koblinger til SDL.
Markdown kan være dokumentverten; språket og modellen er uavhengige av Markdown,
Mermaid, FOX og React. BoxUI er komposisjonsformen i første profil.

## Leserekkefølge

1. [Mandat og studie](Mandate-and-Study.md): intensjon, faktisk gjenbruk, avgrensning.
2. [Krav og dekning](docs/requirements.md): hva denne leveransen beviser.
3. [Språkdefinisjon](docs/language.md) og [EBNF](grammar/sdui-0.1.ebnf).
4. [Arkitektur og AST](docs/architecture.md): eiere, data og konkrete kall.
5. [Runtime-kontrakt](docs/runtime-contract.md): foreslått SDL/vert-grense, ikke implementert ABI.
6. [Videre plan og beslutninger](docs/implementation-plan.md).
7. [Testbevis](evidence/verification.md).

## Kjør parseren

Python **3.11 eller nyere**, bare standardbiblioteket. Fra denne katalogen:

```sh
PYTHONPATH=src python3 -m sdui examples/main-page.sdui
PYTHONPATH=src python3 -m sdui examples/main-page.sdui -o /tmp/main-page.ast.json
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

Parseren gir versjonert JSON-AST med UTF-8-byteområder og linje/kolonne.
Normalt valideres også den lokale språkprofilen. `--syntax-only` gir en AST uten
profilgodkjenning. Feil kommer som JSON på stderr; returkode 2 betyr språkfeil,
3 I/O-feil. `-` leser kilde fra stdin. Ingen refererte SDL-filer åpnes eller kjøres.

[Hovedeksemplet](examples/main-page.sdui) konkretiserer eierens ytre boks,
leftTop/rightTop, widgetrader, SVG og `setHandle`.
[Generert AST](examples/main-page.ast.json) kan sammenlignes direkte med kilden.
[Statisk eksempel](examples/static.sdui) trenger ingen SDL-binding.

## Faktisk status

Implementert: lexer, parser, immutable AST, lokal validering, CLI og tester.
Ikke implementert: layout/rendering av SDUI, Markdown-fence i XFMD, SDL-runtime,
modullasting, binær runtime-ABI eller flytting av BoxUI ut av Mermaid-forken.
En AST er ikke et kjørbart eller visuelt verifisert grensesnitt.

Arbeidet ligger avgrenset i `SDUI/`. Eksisterende SDL-/SDP-dokumenter er ikke
endret, og deres samtidige, lokale endringer er bevart.
