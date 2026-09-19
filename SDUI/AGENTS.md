# SDUI Repository Guidelines

## Formål og autoritet

Les README.md, Mandate-and-Study.md, docs/language.md, docs/architecture.md og
berørte krav før endringer. SDUI 0.1 er en autorisert parserprototype, ikke en
vedtatt utvidelse av SDL eller et ferdig UI-produkt. Skill eierkrav, valgte
profilregler, forslag og testbevis. Bevar parallelt arbeid i SystemDesignLanguage
og resten av SDP; endre ikke deres autoritet eller branch som sideeffekt.

## Arkitektur og filer

`grammar/` eier språkets EBNF. `src/sdui/lexer.py` lager tokens, `parser.py` bygger
syntakstreet, `ast.py` eier immutable datatyper, og `validate.py` sjekker lokal
semantikk. `__main__.py` er CLI/I/O-adapter. Parser og validator skal verken åpne
SDL-referanser, kjøre callbacks eller importere en GUI-/rendereravhengighet.
Hold én tydelig rolle per fil; vurder oppdeling rundt 300 linjer.

## Endringsrekkefølge

Oppdater krav/profil → EBNF og AST-kontrakt → parser/validator → eksempler og
negative tester → dokumentert verifikasjon. Endre grammatikk og implementasjon
samlet. Behold konkrete kall og kildeeiere i arkitekturdokumentet. Versjoner
inkompatible språk- eller AST-endringer; ikke omtolk eksisterende uttrykk stille.
Symbolreferanser og setHandle er data i denne profilen, aldri Python eval/exec.

## Stil og kontroll

Norsk prosjektprosa, engelske kodeidentifikatorer. Python 3.11+, fire mellomrom,
standardbibliotek, dataclasses for AST og eksplisitte SduiError-diagnoser.
Fra SDUI: `PYTHONPATH=src python3 -m unittest discover -s tests -v`.
Regenerer eksempel-AST med `PYTHONPATH=src python3 -m sdui
examples/main-page.sdui -o examples/main-page.ast.json` når formatet endres.
Test kildeposisjoner, feilveier, grenser og faktisk bevart struktur. Oppgi reelle
resultater og kjente hull. Ikke hev SDL-kompatibilitet, rendering eller runtime-
verifikasjon ut fra en grønn parser-test. Ikke inkluder andres ucommittede arbeid
i en commit/push; en publisering må ha tydelig avgrenset innhold og historikk.
