# SystemDesignLanguage

**2026-09-22:** [SDL-verktøyet](tools/README.md) genererer nå valgte strukturelle
viewpoints fra validert SDL til Markdown/Mermaid og valgfri SVG-utskrift.
[Checkpoint-status og viewpoint-katalog](../docs/checkpoint%231/08-SDL-Viewpoints-and-Implementation-Status.md)
beskriver design-core 0.5 med bruksmål, allokering, datakontrakter, Channels,
scenarioer og implementasjonsplan. V2–V4 er levert som verktøyfunksjoner. Diagrammene genereres av verktøyet, ikke manuelt per eksempel.

[Generert G1–G5-design](../SDUI/design/viewpoints/implementation.md) er inngangen
for gjennomgang før Go-implementasjonen.

Utforskende modellutvikling for SDP. Første leveranse er [researchkatalogen](research/existingDesignLanguages/README.md), med [syntese og neste avgrensede oppgave](research/README.md).

Les [arbeidsdokumentet med renskrevet mandat og innledende studie](Mandate-and-Study.md) for eierens samlede hensikt, faseoverganger, design-diff, kodekontroll og standardreferanser. Dokumentet er utarbeidet fra nye eierinnspill 2026-09-14 og skiller mandat fra forslag; det vedtar ikke språk eller canonical SDP-kontrakter.

Oppdrag: [SDP #10, P04](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5618871491), autorisert av eieren med «Ok utfør P04» 2026-09-10. [#9](https://github.com/Hans-Einar/SDP/issues/9) er fortsatt Steering-sporet. Denne katalogen endrer ikke accepted authority, canonical Toolkit eller det foreløpige #7/#8-arbeidet.

Status 2026-09-21: research og kandidatsemantikk er fortsatt ikke et normativt
SDP-schema. [Checkpoint #1](../docs/checkpoint%231/07-SDUI-0.2-and-Go-Direction.md)
er oppdatert med SDUI 0.2 og valgt Go-retning. [Go-området](go/README.md) har
kataloger for parser og runtime, men ingen implementasjon ennå.
Eksisterende [design-core-parser](../experiments/design_core/README.md) er
strukturell Python-kode. [MVP1-korpuset](../experiments/mvp1_sdl/README.md)
er en separat kandidatprofil, ikke et kjørbart system.
Felles videreplan står i [SDUI-planen](../SDUI/docs/implementation-plan.md).
