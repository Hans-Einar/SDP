# SDUI i Go

G1-M1/M2: SDUI 0.2 lexer/parser og kildeposisjonert AST er implementert uten GUI eller
I/O i parserpakken. Den taggede JSON-AST-en samsvarer med de eksisterende
Python-fixturene. CLI-en validerer lokal profil som standard. Normalisering bevarer regioner,
gjenbruk og instansbaner; 93 porttilfeller sammenlignes med Python-grunnlaget.

Modul: `github.com/Hans-Einar/SDP/SDUI/go`, språkbaseline Go 1.25.
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
