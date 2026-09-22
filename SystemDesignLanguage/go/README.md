# SDL i Go

G4-M1 leverer parser, kildeposisjonert AST, symbol-/typekontroll, data-/wire- og
Channel-/scenariovalidering og kanonisk form for **design-core 0.5**.
Ingen struktursetning utføres. Checkpoint-kandidater og hele MVP1-korpuset er
ikke del av profilen. G4-M2 leverer også den eksplisitte kjøreprofilen
[action-core 0.1](../../docs/SDL-Executable-Action-Profile.md), parser og runtime
med typede records og registrerte Go-funksjoner.

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

Kjøreprofilkontroll: `go run ./cmd/sdl action-check examples/echo.sdl`.
`runtime.New` krever eksplisitt signaturregistrering; ingen kilde kjøres som Go.

G4-M3 har en typet SDUI-port i `bridge/`. Enkel Echo og en eksplisitt simulert
EditAptCell går gjennom begge runtimene. `go run ./cmd/sdl-simulate` skriver
korrelert hendelsesspor; `go run -tags desktop ./cmd/sdl-demo` viser den native
Fyne-prototypen. `examples/simulation` er separat håndskrevet Go-domenelogikk.
