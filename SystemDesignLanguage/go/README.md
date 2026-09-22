# SDL i Go

G4-M1 leverer parser, kildeposisjonert AST, symbol-/typekontroll, data-/wire- og
Channel-/scenariovalidering og kanonisk form for **design-core 0.5**.
Ingen struktursetning utføres. Checkpoint-kandidater og hele MVP1-korpuset er
ikke del av profilen. Den eksplisitte kjøreprofilen følger G4-M2.

Modul: `github.com/Hans-Einar/SDP/SystemDesignLanguage/go`, Go 1.26 som felles
baseline. Verifisert med Go 1.27.1. Strukturkjernen bruker bare standardbiblioteket.

Fra denne katalogen:

```sh
go test -race ./...
go run ./cmd/sdl check ../../SDUI/design/architecture.design
go run ./cmd/sdl ast ../../SDUI/design/architecture.design
go run ./cmd/sdl format ../../SDUI/design/architecture.design
```

`parser.Data` beholder Python-profilens navngitte JSON-AST-form, inkludert start-
og sluttposisjoner. 151 porttilfeller og full SDUI-modell har sammenligningsbevis.
[Faktiske bevis](evidence/G4.md), [felles plan](../../SDUI/docs/implementation-plan.md),
[språkdefinisjon](../../docs/Design-Language-Definition.md).

Python beholdes bare som midlertidig portgrunnlag frem til alle konsumenter og
viewpoint-generatoren er erstattet. Det finnes ingen Python-fallback i Go.
