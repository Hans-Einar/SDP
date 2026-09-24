# Checkpoint #1 — implementert Go, runtime og dokumentnavigasjon

Oppdatert 2026-09-22 etter G1–G6. Dette er gjeldende implementasjonsstatus;
tillegg 07–09 bevarer beslutnings-/V-fasegrunnlaget. Checkpoint-kandidater blir
ikke vedtatte språkregler fordi Go-verktøyene finnes.

| Område | Levert profil og bevis |
| --- | --- |
| SDUI-frontend | SDUI 0.2, AST/spans, navn/formatering/bindinger, normalisert instanstre. [G1](../../../SDUI/go/evidence/G1.md) |
| Layout/presentasjon | Relative mål, ratio, rader/wrap, font i DIP, klipp, SVG og Fyne med native button/input; avgrenset Markdown. [G2](../../../SDUI/go/evidence/G2.md) |
| SDUI-runtime | Typede handles/events/propertybatches, accepted/draft, fokus, kompatibel modellreload og siste gyldige UI. [G3](../../../SDUI/go/evidence/G3.md) |
| SDL-frontend | Design-core 0.5 med struktur/bruksmål, data, channels/scenarioer og planfakta; 151 fryste parserprøver. [G4](../../../SDL/go/evidence/G4.md) |
| SDL-kjøring | Action-core 0.1: typed records/actions og registrerte Go-funksjoner; SDL↔SDUI-bridge, modelreload, Go-rebuild/restart. [Profil](../../../SDL/docs/profiles/SDL-Executable-Action-Profile.md) |
| Go-generering | Typede modellkonstruktører, separate håndskrevne domenefunksjoner, identiske kilde-/kompilerte hendelsesspor og state-dokumentasjon. [G5](../../../SDL/go/evidence/G5.md) |
| Dokumenter | 11 viewpoints, statiske kataloger eller lett navigator, typed utvalg, semantiske SVG-symboler, cache/IPC/leases og ekte XFMD-paneler. [G6](../../../SDL/go/evidence/G6.md) |
| Klasser | Class-core 0.1: eksplisitte medlemmer, roller, multiplisitet, aggregation/composition og kildekoblede diagrammer. [Profil](../../../SDL/docs/profiles/SDL-Class-Profile.md) |

G-fasegrensene er fullført innen de dokumenterte profilene. Ingen aktiv
Python-implementasjon av SDUI 0.2/design-core eller viewpointprojektor beholdes.
Fryste testdata og daterte bevis er bevart; Python brukes fortsatt til enkelte
native testdrivere og i andre, ikke erstattede SDP-eksperimenter.

## Åpne og kjøre

Start med [generert utviklingsplan](../../../SDUI/design/viewpoints/implementation.md),
[viewpoint-indeks](../../../SDUI/design/viewpoints/index.md),
[navigator](../../../SDUI/design/navigation/navigator.md) og
[valgt UI-tilstand](../../../SDUI/design/runtime-preview/entry.md).
Diagrammer, fakta og utskrifter er produsert av SDL-verktøyet fra validerte modeller.
Ingen manuelt tilføyde relasjoner i genererte diagrammer.

[SDL-kommandoer](../../../SDL/go/README.md) dekker check/AST,
viewpoints/view, sdl-viewsd, sdl-gen, sdl-compiled og sdl-document.
[SDUI-kommandoer](../../../SDUI/go/README.md) dekker AST/dump/Markdown/SVG og Fyne.
Linux desktop trenger X11/OpenGL/C-byggmiljø. Kjerne-/headless-kjøring trenger
ikke GUI. Go 1.26-baseline, prøvd med Go 1.27.1/Fyne 2.8.1.

XFMD-integrasjonen er [PR #38](https://github.com/Hans-Einar/xfmd/pull/38),
basert på main c245fd9 i eget worktree. Den gir separat navigator/hovedpanel,
eksplisitt vindu/panel og frigjør dokumentleases ved bytte/lukking. Native
klikk, to vinduer, ugyldige forespørsler, restart og opprydding er prøvd.
Det gamle xfmd-boxui-worktreet og Mermaid-repoene er ikke endret.
Mermaid brukes som diagrambackend; SDUI-layout ligger i eget Go-bibliotek.

## Grenser som fortsatt gjelder

- G4s EditAptCell er en merket Go-simulering. Ingen Ponsse-maskinlogikk eller automatisk MVP1-produktmigrering er levert. Det separate [MVP1-korpuset](../../../experiments/mvp1_sdl/README.md) er en kandidatøvelse.
- Design-core beskriver systemdesign, action-core kjører en liten eksplisitt handlingsprofil, class-core beskriver klassedesign. De er ikke tre alternative tolkninger av samme kilde.
- Mode er drifts-/allokeringskontekst. Egen State/state-machine-profil er ikke vedtatt. Database betyr persistent datakilde, ikke nødvendigvis SQL.
- Innebygd Markdown har en avgrenset profil. Mermaid i SDUI er prøvd for flowchart/graph med registrert backend. Ingen garanti om alle diagramtyper eller full XFMD-innholdsdekning. Scroll-layout avvises, og svg-widget er plassholder.
- Hot reload publiserer validerte modeller. Endret Go bygges og prosessen restartes; ingen native dynamisk Go-kodeutskifting. Domenestate over prosessrestart krever egen lagring.
- En utført domenetransaksjon kan ikke automatisk rulles tilbake etter UI-publiseringsfeil. Ingen automatisk event-replay. Broker beholdt lease etter usikker levering krever eksplisitt opprydding hvis leseren krasjet.
- Dokumenttjeneste/XFMD IPC er Linux-profil med privat socket og samme UID. Native SDUI bruker Fyne; FOX-widgets, TUI og C-ABI er ikke ekstra leveranser.

## Sporbarhet

G1 → G2 → G3 → G4 → G6 → G5 er branchstakken; G6 ble implementert før G5
fordi G5-M3 konsumerer viewpointprojektoren. Én commit per verifisert milepæl,
push etter fase, samlet PR mot sdp-vNow. [Historikk](../../Development-Branch-Stack.md)
og [akseptanseplan](../../../SDUI/docs/implementation-plan.md). Ingen automatisk merge.


## Separat repoavvik

Overordnet SDP-validator feiler på ni eksisterende ID-formatavvik i
Traceability/Relations.yaml (ITR-/SLC-/SPR-prefikser). Toolkit og Traceability er
uendret fra målbranchens 9ad4324. [Bevis](../../../SDL/go/evidence/sdp-validator-baseline.txt).
Dette inngår ikke i språkporten og betyr at hele repoets kontraktscheck ikke er grønt.
