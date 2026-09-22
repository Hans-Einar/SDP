# SDUI Repository Guidelines

## Formål og autoritet

Les README.md, Mandate-and-Study.md, docs/language.md, docs/architecture.md,
docs/implementation-plan.md, docs/layout-language-proposal.md,
docs/target-architecture.md og berørte krav før endringer. SDUI 0.2 er dagens
Go-profil. Python-porten er fullført; ingen gammel frontend/fallback skal
reintroduseres. Skill vedtatte profilregler, forslag og faktisk testbevis.
Bevar annet lokalt arbeid og følg rotens fasebranch-/milepælregler.

## Arkitektur og filer

`grammar/` eier EBNF. `go/parser` eier tokens/AST, spans, lokale regler og
normalisering. `go/layout` eier målt geometri; `go/markdown` og `go/svg` innhold
og statisk eksport. `go/runtime` eier UI-tilstand; `go/host/fynehost` eier native
kontroller og UI-tråd. `go/presentation` gir strukturdumper. `go/codegen`
genererer konstruktører for samme modell. CLI/I/O ligger under `go/cmd`.
Parser/runtime importerer ikke en GUI-/rendereravhengighet. Ingen parser åpner
SDL-ref eller kjører callbacks. Hold én rolle per fil; vurder oppdeling rundt
300 linjer. Dokumenter faktiske kall og kildeeiere.

## Endringsrekkefølge og kontroll

Krav/profil → EBNF/AST-kontrakt → parser/validator → eksempler/negative prøver →
verifikasjon. Versjoner inkompatible endringer og port aktive konsumenter samlet;
versjonering krever ikke bakoverkompatibilitet. Ikke omtolk eksisterende uttrykk.

Norsk prosjektprosa og engelske kodeidentifikatorer. Go 1.26+, verifisert med
1.27.1; Fyne er første interaktive vert. SVG og Fyne deler målt geometri.
Ingen ekstra C-ABI, FOX- eller TUI-vert uten konkret nytt omfang. Fra SDUI/go:
`go test -race ./...`. AST: `go run ./cmd/sdui ../examples/main-page.sdui
-o ../examples/main-page.ast.json`. Test kildeposisjoner, grenser, feilveier og
bevart struktur/tilstand. Native prøver bruker egen skjermserver og eget bruker-
område. Grønne parserprøver er ikke bevis for GUI, SDL-runtime eller domene.
Fryste portfixturer og daterte historiske bevis beholdes; Python-testautomatisering
for eksterne prosesser er tillatt, men ingen Python-språkimplementasjon.
