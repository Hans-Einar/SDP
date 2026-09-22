# SDL/SDUI-parser og runtime beskrevet med SDL

**Start med [generert G1–G5-implementasjonsplan](viewpoints/implementation.md).**
Den viser planlagte milepæler, eiere, avhengigheter og eksempelforløp fra SDL.

[SDL-verktøyets genererte viewpoints](viewpoints/viewpoints.md) og
[samlet rendret utskrift](viewpoints/printout.md) er avledet fra denne kilden.
Regenerering og utvalg står i [verktøydokumentasjonen](../../SystemDesignLanguage/tools/README.md).
Use Case/Feature og allokering genereres fra SDL. Datakart, packet og kontraktkontrollerte Channel-sekvenser genereres også;
ingen diagramfakta tegnes inn manuelt.

Oppdatert 2026-09-22. [architecture.design](architecture.design) er den felles
målstrukturen for **begge språkene**. Den erstatter den tidligere korte
SDUI/vertsmodellen på samme sted. Ingen kopi opprettes under SDL-katalogen.

Modellen er skrevet i implementert `design-core 0.5` og passerer den eksisterende
SDL-parseren. Dette profilnummeret gjelder SDL-struktur, ikke utgått SDUI 0.1.
Den beskriver planlagt Go-kode; kjørbar Go-parser/runtime finnes ennå ikke.

## Åpne modellen

- [SDL-kilde](architecture.design): autoritativt eierskap, struktur og grensebruk.
- [Generert ansvarsoversikt](architecture.catalog.md): Units, underenheter, ansvar og porter fra AST.
- [Generert AST med symboltabell](architecture.ast.json): kildeposisjoner og typede noder.
- [Valideringsrapport](architecture.validation.json): resultat, tellinger og SHA-256 for kilder/verktøy/artefakter.

Fra SDP-roten:

```sh
python3 experiments/design_core/design_core.py check SDUI/design/architecture.design
python3 experiments/design_core/design_core.py ast SDUI/design/architecture.design
python3 SDUI/tools/export_design.py
```

De to første kommandoene er den eksisterende SDL-CLI-en. Den siste bruker samme
`check`, `canonicalize`, `symbol_table` og `to_json` for å regenerere AST,
ansvarsoversikt og rapport. Én frontend brukes; V1–V4 er definert i design-core 0.5, uten gammel fallback.
Ugyldig modell stopper eksport før eksisterende artefakter erstattes.

Verifisert kilde 2026-09-22: 368 deklarasjoner og 1106 fakta. 94 Functionality-er
har eksplisitt milepælkobling. De 35 Activities inkluderer fem G-faser og 18
milepæler; 17 Channels, 35 Messages og ti scenarioer beskriver samarbeid.
Syntaks, typer, eierskap, kontrakter, korrelasjon og kanonisk form kontrolleres.
Tallene viser modellomfang. G1–G5 har eksplisitt status planned.

## Bruksmål og arkitekturbidrag

`DesignAuthor` og `DesignReviewer` er modellens eksterne roller. Bruksmålene
kobles til Features gjennom supports og til Functionality gjennom contributes-to.
`TraceViewpointFacts` bidrar direkte til `InspectModels`, uten en ekstra Feature.
`StructuralModelInspection` støtter både inspeksjon og redigering; delte ansvar
beholder samme identitet og eier i alle utsnitt.

Allokeringen plasserer konkrete ansvar i `CommandLineHost` eller `FyneHost` i
eksplisitte modi. Dette flytter ikke bibliotekene eller deres eierskap. De nye
allokeringene er designpåstander, ikke bevis på at Fyne-/Go-koden kjører.
VP07s 15 hull viser ufullstendige modusutsnitt, og modellen har ikke en komplett
runtimeallokering. En delt Functionality gjør ikke alle dens Features kjørbare
i alle modi hvor dette ene ansvaret har en plassering.

`SdlViewpointGenerator` under `SdlLibrary` eier `ProjectSdlViewpoints`,
`ExportViewpointMarkdown` og `TraceViewpointFacts`. Python-verktøyet er dagens
implementasjon av disse avgrensede ansvarene; Go-porten gjenstår.

## Parserdesign

`SdlLibrary` og `SduiLibrary` er logiske bibliotekgrenser. Hver frontend inneholder
sin egen lexer, parser, validator og normalisering. Dette er to forskjellige
kildespråk, ikke to implementasjoner av samme grammatikk. Frontend koordinerer
kjeden; de enkelte Units eier hver sin Functionality.

| Trinn | SDL | SDUI |
| --- | --- | --- |
| Tokens/AST | SdlLexer, SdlParser | SduiLexer, SduiParser |
| Navn/profil | SdlValidator: symboler, typede strukturregler og profil | SduiValidator: navn, widgetargumenter, relativ formatering og symbolske bindinger |
| Normalisering | SdlNormalizer: validerte designfakta og kildekart | SduiNormalizer: definisjonsekspansjon, header/body/footer og kildekart |
| Kjørbarhet | SdlExecutionGate kontrollerer den framtidige kjøreprofilen | SduiRuntime kan vise en ubundet, validert UI-modell |

`SourceLoader` leser begrensede kilder og identifiserer revisjonen; parseren får
et allerede innlest kildesnapshot. SDUI `ref` åpner ikke filer under parsing.
`DiagnosticReporter` er ansvarlig for kilde-/bindingsdiagnoser. Ressursgrenser
skal kontrolleres i kjeden, ikke bare før innlesing.

`stateless`/`deterministic` på rene frontendoperasjoner er designkrav med kilden
og profil som eksplisitte input. Det utelukker ikke midlertidige parserobjekter;
det betyr ingen beholdt sesjonstilstand mellom kall. Dette er ikke en påstand
om deterministisk domeneatferd, nettverk eller fontmåling.

## Runtime og domenekobling

`SdlRuntime` eier opprettelse/lukking av SDL-instans. `SdlStateStore` forvalter
modellert domenetilstand og snapshots. `SdlFunctionRegistry` registrerer de
håndskrevne Go-funksjonene og kontrollerer signaturer. `SdlDispatcher` validerer
handlingens input, utfører registrert funksjon, korrelerer resultat og avbryter
ventende handlinger. `GoDomainImplementation` er en ekstern logisk Unit som eier
produktalgoritmen, ikke en ny påstått prosess eller maskinimplementasjon.

`SduiRuntime` eier UI-instansen. `SduiInstanceStore` forvalter identiteter og
generations; `SduiPropertyStore` validerer/publiserer egenskapsbatcher og sporer
draft; `SduiDispatcher` validerer og ruter events og avviser gamle events.
`UiStateReconciler` matcher kompatible widgets og bevarer eller nullstiller UI-state.
Fyne eier den konkrete tekstmarkøren/fokusmekanismen, mens runtime beholder den
logiske widgetidentiteten og avtalt verdi/draft. Ingen toolkitpeker er et SDL-handle.

`SdlUiBindingAdapter` ligger utenfor begge kjerner. Den løser symbolske callbacks,
kobler typed widgethandles, ruter kall og publiserer domeneoppdateringer til UI.
Den kan kobles fra uten å gjøre en statisk SDUI-modell ugyldig. `DomainBindingPort`
er grensen UI ser; `SdlExecutionPort` er grensen til SDL-kjøringen.

Planlagt hendelsesforløp; de eksplisitte eksempelbanene i VP08 validerer
meldingsrekkefølge og kontrakter, mens tilstandsvirkningene nedenfor gjenstår:

1. Fyne sender aktivering/commit med widgetidentitet, generation og revisjon.
2. SDUI sjekker aktuell widget, enabled, binding og verdi før dispatch.
3. Adapteren finner SDL-handlingen; SDL kontrollerer input og Go-signatur.
4. Registrert Go-funksjon utfører domenearbeidet og gir et korrelert resultat.
5. SDL publiserer akseptert domeneendring; adapteren lager UI-oppdatering.
6. SDUI validerer hele batchen og publiserer nytt snapshot/layout til verten.

Et programmatisk value-sett utløser ikke automatisk samme input-callback.
Utestående eller sene resultater må kontrolleres mot instansens generation.
Dette konkretiserer [runtime-kontrakten](../docs/runtime-contract.md); eksakte
signaturer, feiltyper, leveringsregler og skjemaer må fortsatt defineres i G3/G4.

## Hot reload og Go-generering

`SourceWatcher` samler filendringer. `ReloadCoordinator` bygger kandidatmodeller
via begge frontender ved behov og publiserer først etter vellykket validering.
`UiStateReconciler` håndterer UI-kompatibilitet; `DomainStateMigrator` vurderer
SDL-state mot eksplisitt migrerings-/resetpolicy. Adapteren bytter bindinger ved
samme avtalte generationsgrense. Gammel modell beholdes ved feil.

Kun SDUI trenger å parses på nytt ved en isolert UI-endring. Endret SDL-binding
kan kreve ny bindingsvalidering av UI selv om SDUI-kilden er uendret. Dette er en
avhengighetsregel som må implementeres; vi har ikke en inkrementell kompilator nå.
Ingen reload skal gjenta en allerede utført domenehandling. Publisering må
samordnes med pågående hendelser, kansellering, fokus/draft og ressurser.

`GoCodeGenerator` lager modellkonstruktører og bindingsregistrering. Kjøremodellen
skal være den samme ved filbasert og generert programstart. `GoBuildRunner` bygger
og restarter når selve Go-koden endres. Den bytter ikke maskinkode i prosessen.
Genererte kilder skal aldri overskrive håndskrevne domenefunksjoner.

## Layout, innhold og vert

`SduiLayout` måler, løser ancestor-relative dimensjoner, fordeler geometri,
beregner klipping og lager et felles prepared frame. `SduiPresentation` eksporterer
SVG, konsoll og Markdown-dokumentasjon fra dette grunnlaget. Eksport kjører ikke
domenehandlinger. `MarkdownProvider`, `DiagramProvider` og `ResourceStore` er
separate innholdstjenester. Mermaid er en ekstern diagrammotor, ingen språkgren i SDUI.

`FyneBackend` oppretter/oppdaterer konkrete widgets, håndterer fokus/tekstinntasting
og frigjør native ressurser. `FyneHost` setter sammen interaktiv kjøring og eier
UI-trådspublisering. `CommandLineHost` er alternativet for eksport uten vindu.
Disse er Containers i modellen; bibliotekene er Units. `contains` uttrykker
logisk inndeling, ikke at alle bibliotekene er plassert i hver prosess.

Fyne og SVG bruker samme målte geometri. Fyne-importer hører til vert/backend,
ikke parser/runtime. Forskjellige tekst-/SVG-renderere må fortsatt verifiseres
visuelt; felles modell alene beviser ikke identisk rendering.

## Porter og modi

`consumes` uttrykker bruk, ikke et faktisk kall eller krav i alle modi.
`provides` gjelder Capability, **ikke Interface**, i dagens parser. Koblingen
mellom portnavn, konkret tilbyder og funksjonssignatur er derfor en kontrakt som
forklares her og senere må formaliseres; vi forfalsker ikke den relasjonen med
feil objekttype i SDL-kilden.

| Portgruppe | Planlagt tilbyder og kontraktgrense |
| --- | --- |
| SourceInputPort, FileChangePort | Vertens fil-/watcher-adapter; kildebase, begrensninger og feil |
| SourceSnapshotPort, DiagnosticPort | SourceLoader og DiagnosticReporter; kildeidentitet/revisjon og posisjonerte diagnoser |
| SdlFrontendPort, SduiFrontendPort | Respektiv frontend; kilde → validert modell/diagnoser |
| SdlModelPort, SduiModelPort | Frontendnormalisering eller genererte konstruktører; immutable modell med kildekart |
| ExecutionProfilePort | Avtalt SDL-profil/kjørbarhetskontroll; støttede operasjoner og komplett binding |
| DomainFunctionPort, DomainStatePort | Registrerte Go-implementasjoner og SdlStateStore; signatur, stateeier, resultat og revisjon |
| SdlExecutionPort, DomainBindingPort | SdlRuntime og SdlUiBindingAdapter; typed handlinger, korrelasjon og livstid |
| UiSessionPort, UiStatePort, UiSnapshotPort | SDUI-runtime/stores; logisk state, hendelser og immutable snapshots |
| SdlReloadPort, UiReloadPort, BindingReloadPort, ReloadPort | Respektiv runtime, adapter og koordinator; forberedelse, kompatibilitet, publisering og avbrudd |
| MeasurementPort, ContentProviderPort | Vertens måletjeneste og MarkdownProvider; bredde→høyde, fontenhet og begrenset arbeid |
| DiagramPort, DiagramEnginePort, ResourcePort | DiagramProvider, eksisterende diagrammotor og ResourceStore; inert innhold, ressursidentitet og release |
| PreparedFramePort, WidgetBackendPort | SduiLayout og FyneBackend; samordnet geometri/klipp, widgetlivstid og events |
| GeneratedArtifactPort, BuildToolPort, ExportSinkPort | GoCodeGenerator, vertens byggeverktøy og filutdata; eide filer og eksplisitt utførelse |

Modes er uavhengige navngitte kontekster, uten implisitt arv eller kjørbar
aktiveringsbetingelse. `SourceInspection` er strukturell kontroll; `UiPreview`
er ubundet UI; `BoundExecution` krever SDL-/domenekobling; `LiveEditing` er UI-reload;
`BoundLiveEditing` beskriver nødvendige porter for samordnet SDL/SDUI-reload;
`StaticExport` er dokumentasjon; `NativeBuild` er Go-bygg; `RichDocument` er rik
innholdsbehandling. En vert kan senere aktivere flere kontekster etter avtalt policy.
Aktivitetenes `refines` er detaljering, ikke sekvens, tilstandsmaskin eller scheduler.

## Språkgrense og sporbarhet

Modellen dekker ansvar for alle delene i [G1–G5-planen](../docs/implementation-plan.md).
Parseren kontrollerer typede recordfelt, eksplisitte scenario-steg og korrelasjon
innen den avgrensede design-core 0.5-profilen. Den kontrollerer ikke Go-signaturer,
kall i implementert kode, atomisitet, state-maskiner, ressursbudsjetter, trådregler
eller samsvar mellom Go-kode og modellen. V2–V4 utvidet grammatikk og tester
samlet; kjørbar semantikk må få presise profiler og kontrakter i G-fasene.

| Modellområde | Plan / fremtidig kodeområde |
| --- | --- |
| SduiFrontend og barn | G1; SDUI/go/parser |
| SduiRuntime og barn | G3; SDUI/go/runtime |
| SdlFrontend og barn | G4-M1; SystemDesignLanguage/go/parser |
| SdlRuntime og barn, SdlUiBindingAdapter | G4-M2–M4; SDL-runtime og vertens kobling |
| SduiLayout, innhold, Fyne, SduiPresentation | G2; kataloger opprettes ved implementasjon |
| ReloadCoordinator, SourceWatcher | G3/G4; delt utviklingsvert, ikke duplisert i parserne |
| GoCodeGenerator, GoBuildRunner | G5; kataloger opprettes ved implementasjon |

Disse kodeplasseringene er manuell sporbarhet. SDL-fakta gir ikke eksisterende
Go-pakker eller filbindinger. [Målarkitekturen](../docs/target-architecture.md)
eier prinsippene; denne modellen eier den detaljerte ansvarsfordelingen.

## V2 — data og valgt packet-prøve

`DesignSourceDocuments` ligger hos `DesignSourceArchive`, en planlagt persistent
kildebase eid av SourceLoader. `UiSessionState` er transient hos SduiInstanceStore.
Kontraktene beskriver påkrevd kilde/revisjon og UI-generation med optional draft;
optional draft er ikke automatisk en tom streng.

`ProjectUiGeneration` projiserer UI-tilstanden til familien `UiGenerationNotices`.
`UiGenerationWire` er et eksplisitt valgt **prototypeformat** med 16-bit versjon
og 64-bit generation, big-endian/MSB-first. Dette beskriver en dokumenterbar
wireprøve, ikke en vedtatt transport mellom Go-bibliotekene, SDL-runtime-ABI,
P1000- eller StanForD-format. Native Go-kall trenger ikke serialiseres slik.

## V3 — meldinger og scenarioer

V3 innførte sju logiske Channels og fire eksplisitte scenarioer som dekker akseptert/avvist
UI-handling og modellreload. BoundActionAccepted går FyneBackend → UI-dispatch →
bindingsadapter → SDL-dispatch → registrert Go-domeneansvar, med korrelert retur
for hvert kall. BoundActionRejected avsluttes ved UI-valideringen uten domene-kall.
Dette er planlagt samarbeid, ikke observerte kjøreresultater.

Reload-scenarioene skiller publisert modell fra avvist kandidat. Bare den
aksepterte banen inneholder UiGenerationNotices. Kontraktene har optional draft,
generation og diagnose der fravær er meningsfullt; ingen kjørende tilstandsmigrering
hevdes. [MessageSet](viewpoints/message-sets.json) er avledet av SDL-verktøyet.

## V4 — samlet gjennomgangsgrunnlag

Fase-/milepælplanen er nå også SDL. `addresses` kobler milepæler til Functionality,
`delivers` kobler faser til Features, `depends-on` angir eksplisitte forutsetninger,
og `illustrates` kobler scenarioer til arbeidsaktiviteten. implementation.md lages
av SDL-verktøyet fra disse faktaene, uten en innebygd G1–G5-plan.

De ti scenarioene viser frontendens lexer/parser/validering/normalisering,
felles layout for Fyne og SVG, lokal Go-handling uten SDL, akseptert/avvist
SDL-binding, UI-reload, SDL-reload og native bygg. Meldingskontraktenes bytesfelt
for tokens/AST/modell/scene er bevisst opake artefaktgrenser; komplette Go-structs,
instanslivstid, atomisk reload/state-migrering og dimensjons-/fontenhet skal
realiseres og testes i G-fasene. Ingen datatransformasjon eller tilstandsendring
utledes bare fordi to meldinger følger hverandre i et scenario.
