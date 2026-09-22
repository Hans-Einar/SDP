# SDL-verktøy — Go-inngangen

Python-CLI/projektorer er erstattet i G5-M4 etter fryste portprøver. Det finnes
én aktiv frontend for design-core 0.5 og én viewpointprojektor, i [Go-modulen](../go/README.md).
Historisk Python-kode finnes i Git før G5-M4, ikke som fallback i arbeidstreet.

Fra SDP-roten:

```sh
go -C SystemDesignLanguage/go run ./cmd/sdl check ../../SDUI/design/architecture.design
go -C SystemDesignLanguage/go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --format static --monolithic --output ../../SDUI/design/viewpoints --renderer /absolute/mmdr
go -C SystemDesignLanguage/go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --output ../../SDUI/design/navigation --project sdui
```

Standardeksport gir navigator/overview uten detaljrendering. Statisk eksport
lager separate viewpointkataloger og en valgfri samlet fil. --renderer velger
registrert mmdr; uten flagget beholdes portable Mermaid-blokker. SVG-symbolprofilen
bruker faste aktør-/use-case-/Feature-figurer og eksplisitte kildelenker.

[Utvalg/IPC/leases](../go/README.md), [G6-bevis](../go/evidence/G6.md),
[abstraksjonsnivåer/notasjon](../../docs/SDL-Viewpoint-Levels-and-Notation.md).
Dataset/database er logiske kilder til data; database betyr ikke nødvendigvis SQL.
Mode er allokerings-/driftskontekst. Egen State-maskin er ikke vedtatt i design-core.
