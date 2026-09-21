# SDUI Repository Guidelines

## Formål og autoritet

Les README.md, Mandate-and-Study.md, docs/language.md, docs/architecture.md,
docs/implementation-plan.md, docs/layout-language-proposal.md,
docs/target-architecture.md og berørte krav før endringer. SDUI 0.2 er dagens
kjørbare Python-profil; 0.1-syntaks og grammatikk er fjernet. Eierens beslutning
2026-09-20 er å porte eksempler og fjerne legacy-veier, ikke bevare to dialekter.
SDUI er ikke en vedtatt utvidelse av SDL eller et ferdig UI-produkt. Skill eierkrav, valgte
profilregler, forslag og testbevis. Bevar parallelt arbeid i SystemDesignLanguage
og resten av SDP; endre ikke deres autoritet eller branch som sideeffekt.

## Arkitektur og filer

`grammar/` eier språkets EBNF. `src/sdui/lexer.py` lager tokens, `parser.py` bygger
syntakstreet, `ast.py` eier immutable datatyper, og `validate.py` sjekker lokal
semantikk. `formatting.py` validerer layoutregler, `normalize.py` lager statiske
instanstrær og `dump.py` viser strukturen i terminalceller uten egen parser.
`markdown_dump.py` eksporterer statisk Markdown fra samme modell.
`prototype_widgets.py` tegner avgrensede button/input-prøver fra modellen;
`prototype_html.py` lager et lokalt widgetgalleri. Fixtureplassering i
`tools/build_widget_previews.py` er ikke en generell layoutmotor eller runtime.
`__main__.py` er CLI/I/O-adapter. Parser og validator skal verken åpne
SDL-referanser, kjøre callbacks eller importere en GUI-/rendereravhengighet.
Hold én tydelig rolle per fil; vurder oppdeling rundt 300 linjer.

## Endringsrekkefølge

Oppdater krav/profil → EBNF og AST-kontrakt → parser/validator → eksempler og
negative tester → dokumentert verifikasjon. Endre grammatikk og implementasjon
samlet. Behold konkrete kall og kildeeiere i arkitekturdokumentet. Versjoner
inkompatible språk- eller AST-endringer; port lokale eksempler/tester samlet og
fjern erstattet kode. Versjonering krever ikke bakoverkompatibilitet. Ikke omtolk
eksisterende uttrykk stille. Ikke opprett en gammel parser-/ABI-fallback.
Symbolreferanser og setHandle er data i denne profilen, aldri Python eval/exec.

## Stil og kontroll

Norsk prosjektprosa, engelske kodeidentifikatorer. Dagens Python-baseline bruker
Python 3.11+, standardbibliotek, dataclasses og eksplisitte SduiError-diagnoser.
Valgt videre implementasjon er Go; `go/parser` og `go/runtime` er foreløpig
bare kataloger. Fyne er første interaktive vert, SVG er statisk eksport fra
felles geometri. Ikke innfør Rust/C-ABI eller en ekstra parser i verten.
Python-frontenden beholdes til Go-porten og dens konsumenter er verifisert;
deretter fjernes erstattede aktive veier. Ingen Go-runtime er implementert.
Checkpoint #1 tillegg 07 og PLAN-003 erstatter eldre FOX-/Rust-planer.
Fra SDUI: `PYTHONPATH=src python3 -m unittest discover -s tests -v`.
Regenerer eksempel-AST med `PYTHONPATH=src python3 -m sdui
examples/main-page.sdui -o examples/main-page.ast.json` når formatet endres.
Test kildeposisjoner, feilveier, grenser og faktisk bevart struktur. Oppgi reelle
resultater og kjente hull. Ikke hev SDL-kompatibilitet, rendering eller runtime-
verifikasjon ut fra en grønn parser-test. Ikke inkluder andres ucommittede arbeid
i en commit/push; en publisering må ha tydelig avgrenset innhold og historikk.
