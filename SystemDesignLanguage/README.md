# SystemDesignLanguage

SDL-verktøyet er implementert i Go og genererer Markdown/Mermaid/SVG-viewpoints
fra validerte modellfakta. G1–G6 er levert innen de eksplisitte profilgrensene.
[Checkpoint #1 — gjeldende implementasjon](../docs/checkpoint%231/11-Go-Implementation-and-Navigation.md)
skiller dette fra bredere kandidater og framtidig Ponsse-produksjonslogikk.

- [Go-verktøy og kjørekommandoer](go/README.md): check/AST, viewpoints, dynamiske utvalg, dokumenttjeneste, runtime og Go-generering.
- [SDUI-systemets SDL-modell](../SDUI/design/README.md) og [generert utviklingsplan](../SDUI/design/viewpoints/implementation.md).
- [Felles faseplan](../SDUI/docs/implementation-plan.md) og [Git-sporbarhet](../docs/Development-Branch-Stack.md).

Design-core 0.5 beskriver struktur, bruksmål/features, allokering, datakontrakter,
Channels, scenarioer og planfakta. Action-core 0.1 kjører typede handlinger via
registrerte håndskrevne Go-funksjoner. Class-core 0.1 beskriver eksplisitte klasser,
roller, multiplisitet og aggregation/composition; den er ikke objektruntime.
Ingen gammel Python-design-core eller viewportprojektor er aktiv.

[MVP1-korpuset](../experiments/mvp1_sdl/README.md) er fortsatt en separat
66-fils kandidatøvelse, ikke kjørbar design-core eller ferdig maskinlogikk.
EditAptCell-eksemplet er en merket simulering, ikke et Ponsse-backend.
[Mandat og studie](Mandate-and-Study.md), [research](research/README.md) og
oppdrag [SDP #10](https://github.com/Hans-Einar/SDP/issues/10) beholdes med sin
opprinnelige autoritet; dette endrer ikke canonical Toolkit/SDP-kontrakter.
