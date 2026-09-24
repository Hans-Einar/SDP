# SDUI — prototyping av brukergrensesnitt

[Dokumentoversikt](docs/README.md) skiller gjeldende profil, implementasjon,
designbakgrunn og daterte referanser.

[SDUIs KanBan](SDP/Agents/KanBan/README.md) følger ideer og lokale konsekvenser
av SDP-planleggingen. Kortene endrer ikke den implementerte språkprofilen.

SDUI 0.2 er implementert i Go: parser/AST, validering, normalisering, relativ
layout, SVG, strukturell konsoll-/Markdown-dump, runtime og Fyne-vert med modellreload.
SDL-kobling og Go-generering bruker de samme modellene og runtimene.
Ingen aktiv Python-frontend eller 0.1-kompatibilitetsvei beholdes.

Fra SDP-roten, med Go 1.26+ (verifisert med 1.27.1):

```sh
go -C SDUI/go test ./...
go -C SDUI/go run ./cmd/sdui ../examples/concept1-bucking.sdui --format dump --entry bucking
go -C SDUI/go run ./cmd/sdui ../examples/concept1-bucking.sdui --format svg --entry bucking -o /tmp/concept1.svg
go -C SDUI/go run -tags desktop ./cmd/sdui-fyne -entry bucking ../examples/concept1-bucking.sdui
```

Linux desktop trenger OpenGL/X11 og C-kompilator. Parser/runtime kan brukes uten
GUI. Fyne er første interaktive vert; XFMD viser generert dokumentasjon og kreves
ikke for UI-kjernen. [Go-innganger](go/README.md), [felles kjøreeksempel](../SDL/go/README.md).

[Concept1-kilden](examples/concept1-bucking.sdui) har seks hovedbokser og
representative kontroller. Det er eksempeldata, ikke portert React- eller
apteringslogikk. [AST](examples/concept1-bucking.ast.json),
[konsolldump](examples/concept1-bucking.dump.txt) og [Markdown-dump](examples/concept1-bucking.dump.md)
bevarer strukturen. [UI/state-dokumentasjon](design/runtime-preview/entry.md)
viser SVG fra felles layout og eksplisitt valgt tilstand.

`[]` er frame, `<>` nestede grupper, `*b` BoxUI-dekorasjon, `{}` formatering.
Komma fortsetter horisontalt; semikolon starter neste rad. Kildestørrelser er
relative; font er absolutt i logiske DIP. Parseren åpner aldri SDL-ref eller kjører
callbacks. En eksplisitt verts-/bridge-registrering kreves for domenekall.

- [Språk og EBNF](docs/language.md), [layoutprofil](docs/go-layout-contract.md), [Markdown-profil](docs/markdown-provider.md).
- [Arkitektur](docs/architecture.md), [runtime](go/runtime/README.md), [Go-generering](docs/go-generation.md).
- [Krav](docs/requirements.md), [milepæler](docs/implementation-plan.md), [datert checkpoint](../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md).
- [SDL-design og genererte viewpoints](design/README.md), [mandat](Mandate-and-Study.md).

Avgrensninger: scroll-layout avvises; innebygd Mermaid er verifisert for
flowchart/graph med registrert renderer. Fyne uten slik provider viser plassholder.
SVG-widgeten er en merket plassholder. FOX-widgets, TUI, full Markdown/Mermaid-
dekning og produksjonsintegrasjon med Ponsse er ikke levert. Det eldre
[HTML-kontrollgalleriet](examples/prototype-controls.html) er en statisk historisk
prøve; ny generell UI-plassering går gjennom Go-layouten.
