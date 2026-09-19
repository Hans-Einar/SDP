# SDUI — arkitektur, AST og plumbing

**ID:** SDUI-ARCH-001. Implementert parserprofil 0.1; runtime/rendering er framtidig.

## 1. Lag og eiere

| Fil | Eneansvar | Konsument |
| --- | --- | --- |
| src/sdui/lexer.py | Begrenset tokenisering og kildeposisjoner | Parser |
| src/sdui/parser.py | EBNF-struktur → immutable syntax tree | CLI, framtidig språkadapter |
| src/sdui/ast.py | Dataclasses, Span, strukturert feil og JSON-serialisering | Parser, validator, CLI |
| src/sdui/validate.py | Lokal profil, navnerom, egenskaper og referansemål | CLI, framtidig kompilator |
| src/sdui/__main__.py | Fil/stdin, JSON-output, exitkoder | Utvikler/verktøy |

Ingen globale registre, GUI-importer, nettverk eller SDL-I/O i kjernen. Prosjektet
har ingen produksjonsavhengighet til Mermaid. Dataklassen Reference er en symbolsk
verdi; ingen Python-funksjonsreferanse lagres i AST.

## 2. AST-format

CLI-envelope: `astFormat: sdui-ast/0.1`, `validation: local-profile | syntax-only`,
`document: Document`. Alle dataclasses serialiseres med eksplisitt `type`.
Tupler blir JSON-arrays med bevart kilderekkefølge. Ingen skjulte defaults legges
inn i syntax tree; senere normalisering blir en annen representasjon.

| Node | Bevarte felt |
| --- | --- |
| Document | profile, references, definitions, connections, span |
| ModuleRef | alias, path, span |
| Definition | name, root, span |
| Box | name (null for anonym), properties, children, content, span |
| Property | name, Literal, span |
| Content / Row | eksplisitte rader / widgets, span |
| Widget | name, kind, arguments, span |
| Argument | name (null for posisjonelt), Literal eller Reference, span |
| Literal | kind, value, span |
| Reference | module, object, member, span |
| Connection | module, object, definition, widget, span |

`Span.start/end` er halvlukket UTF-8-byteområde. `line/column` er ettbasert ved
start. Box-span begynner på `[`, mens navnet er lagret separat; Widget-span
begynner på navnet. Row-span inkluderer ikke separatoren etter raden.
Document-span går fra header til EOF, inkludert etterfølgende kommentarer.
Dette er ikke et concrete syntax tree: kommentarer, quotevalg og tallstavemåte
bevares ikke. Omskriving av kilden med identisk formattering er ikke støttet.

## 3. Semantikk og framtidig normalisering

Validering muterer ikke AST. Navneregistring bruker definisjonens scope, mens
modulaliaser og definisjonsnavn er dokumentomfattende. Lokal kontroll kan fastslå
at et widgetmål finnes, men ikke at en ekstern SDL-metode eksisterer eller matcher.

Neste representasjon bør være en rendererfri UI-modell med eksplisitte grupper,
standardverdier, widgetegenskaper, identiteter og bindinger. SDUI og Concept1-
adapter skal kunne produsere denne. AST og normalisert UI-modell har forskjellige
formål og skal ikke blandes sammen.

## 4. Feil og ressursgrenser

Fail-fast: én SduiError med code/message/span, ingen delvis vellykket CLI-output.
Ingen feilrecovery for editor i denne versjonen. Maks 256 KiB kilde, 50000 tokens,
64 nivåer med bokser, 2048 bokser/widgets og 32 argumenter per widget. Lexer bygger
posisjonskart én gang; parseren bruker begrenset lookahead. Ingen hard tidsgrense
hevdes. Filinnlesing er begrenset før dekoding; API-kilden finnes allerede i minnet.

## 5. Plumbing

| Inngang | Kalt symbol | Fil | Data/resultat og feilvei |
| --- | --- | --- | --- |
| python -m sdui | main | src/sdui/__main__.py | bounded bytes → UTF-8; encoding/I/O-diagnose |
| CLI / bibliotekskonsument | parse | src/sdui/parser.py | source → Document eller SduiError |
| Parser.__init__ | lex | src/sdui/lexer.py | str → tokens med spans; byte/token/escape-grenser |
| Parser.parse | box / content / widget / value | src/sdui/parser.py | strukturelle noder; depth/node/argument-grenser |
| CLI | validate | src/sdui/validate.py | lokale symboler/profil; ingen ekstern resolusjon |
| CLI | to_data | src/sdui/ast.py | tagget AST → JSON-envelope; ingen kodekjøring |

## 6. Teknisk valg

Recursive descent er valgt fordi grammatikken er liten og feilposisjoner må være
lette å følge. EBNF er ikke automatisk kompilert til parserkode. Samsvar kontrolleres
med konkrete positive/negative eksempler og manuell grammatikkgjennomgang; testene
beviser ikke formell ekvivalens. Ingen parsergenerator eller runtime installeres.

Python-prototypen er en referanse og et verktøy for språkiterasjon. Rust-/C++-
implementasjon, eventuell felles SDL-IR og stabil binær ABI er senere beslutninger.
