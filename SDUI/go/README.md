# SDUI — implementasjonsområde i Go

Opprettet 2026-09-21. **Kun katalogstruktur; ingen Go-parser eller runtime er implementert.**
Omfanget er SDUI 0.2 med frames, grupper, Markdown, formatering og symbolske bindinger.

[Felles SDL-design for parser/runtime](../design/README.md) beskriver detaljert
ansvarsfordeling og valideres av eksisterende SDL-parser.

| Katalog | Ansvar |
| --- | --- |
| parser/ | Kilde, AST, kildeposisjoner, diagnoser og lokal validering; ingen utførelse eller GUI-import |
| runtime/ | Instanser, typed tilstand/hendelser og eksplisitte porter; ingen parserkopi eller GUI-import |

Den felles [implementasjonsplanen](../../SDUI/docs/implementation-plan.md) og
[målarkitekturen](../../SDUI/docs/target-architecture.md) eier fasevalg og avhengigheter.
[Checkpoint #1](../../docs/checkpoint%231/07-SDUI-0.2-and-Go-Direction.md) beskriver
status og språkgrenser. Eksisterende [Python-grunnlag](../src/sdui) beholdes til
porten er verifisert. SDL design-core 0.1 er ikke SDUI 0.1.

Modulnavn, Go-versjon, eventuell go.work og felles portpakker fastsettes ved første
kodeleveranse. Ingen tomme API-er eller falske go test-resultater opprettes nå.
De tomme parser/runtime-katalogene spores med .gitkeep og erstattes av kode.
Fyne skal bare ligge i en senere vert; kjernen må kunne testes uten vindussystem.
