# SDUI i Go

G1-M1–M3: SDUI 0.2 lexer/parser og kildeposisjonert AST er implementert uten GUI eller
I/O i parserpakken. Den taggede JSON-AST-en samsvarer med de eksisterende
Python-fixturene. CLI-en validerer lokal profil som standard. Normalisering bevarer regioner,
gjenbruk og instansbaner; 93 porttilfeller sammenlignes med Python-grunnlaget.

Modul: `github.com/Hans-Einar/SDP/SDUI/go`, språkbaseline Go 1.26 (krevd av x/text 0.42.0).
Verifisert med lokal Go 1.27.1 fra go.dev, kontrollert mot publisert SHA-256.
Installerte byggeverktøy ligger utenfor repoet.

Fra denne katalogen:

```sh
go test ./...
go run ./cmd/sdui ../examples/concept1-bucking.sdui
```

Parseren åpner ikke symbolske SDL-referanser eller kjører callbacks.
Python-frontenden beholdes som portorakel til alle konsumenter er erstattet.
[Faseplan](../docs/implementation-plan.md), [bevis](evidence/G1.md).

CLI tilbyr `--format ast|dump|markdown|prototype-svg|prototype-html`,
`--entry`, `--columns`, `--syntax-only` og `-o`. Flag kan stå før eller etter
kilden. Eksport til fil publiseres atomisk etter validering og overskriver
ikke kildefilen. Konsoll/Markdown samsvarer byte-for-byte med Concept1-fixturene.
Prototype-SVG/HTML er kontrollgallerier, ikke generell SDUI-layout; G2 gir den.

G2 leverer felles geometri (`layout`), SVG (`svg`), native Fyne-adapter
(`host/fynehost`) og avgrenset Markdown (`markdown`). [Bevis](evidence/G2.md).
Generell eksport bruker `--format svg --width 1920 --height 1200` og valgfritt
`--mermaid-renderer /sti/til/mmdr --resources DIR`. [Profil](../docs/markdown-provider.md).

Native prøve: `go run -tags desktop ./cmd/sdui-fyne -entry bucking
../examples/concept1-bucking.sdui`. CLI-flaggene må stå før filnavnet i den
native Go-flag-inngangen. Native bygg trenger OpenGL/X11/C-kompilator på Linux.
Vanlige pakkeprøver bruker Fynes minnedriver uten skjermserver.
