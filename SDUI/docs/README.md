# SDUI — dokumentasjon og autoritet

Start med [SDUI README](../README.md) for dagens kjørbare omfang og
[Go-kommandoene](../go/README.md) for bruk. Bare SDUI 0.2 i Go er aktiv.

## Gjeldende profil og implementasjon

| Dokument | Eier hva? |
| --- | --- |
| [language](language.md) | Språkprofil 0.2, AST og lokale valideringsregler |
| [architecture](architecture.md) | Go-pakker, livsløp og ansvarsgrenser |
| [go-layout-contract](go-layout-contract.md) | Måling, relativ layout, DIP-font, klipp og grenser |
| [markdown-provider](markdown-provider.md) | Faktisk Markdown-/Mermaid-dekning; ingen full Markdown-påstand |
| [runtime-contract](runtime-contract.md) | SDUI/SDL-grensen og lenker til detaljerte pakkekontrakter |
| [go-generation](go-generation.md) | Genererte modeller/bindinger; samme runtime som kildekjøring |
| [requirements](requirements.md) | Krav, leveransegrenser og referanser til bevis |
| [implementation-plan](implementation-plan.md) | Felles SDL/SDUI-faser og milepæler med bevis |
| [target-architecture](target-architecture.md) | Valgt Go/Fyne-retning, gjenbruk og grenser mot andre repoer |
| [concept1-console](concept1-console.md) | Concept1-kilde, AST og strukturelle dumpformater |

Språk, layout og runtime har hvert sitt ansvar. En bredere idé i en studie
utvider ikke parseren eller runtime-kontrakten. Ved et reelt avvik mellom kode
og kontrakt registreres feilen; den skal ikke skjules ved å kalle all kode autoritativ.

## Designbakgrunn og daterte referanser

| Dokument | Hvordan det skal leses |
| --- | --- |
| [layout-language-proposal](layout-language-proposal.md) | Opprinnelig forslag; aktive valg avgrenses av language og målekontrakten |
| [frame-composition-proposal](frame-composition-proposal.md) | Eierinnspill og anbefalinger; ikke en parallell språkprofil |
| [concept1-relative-layout-study](concept1-relative-layout-study.md) | Datert undersøkelse av Ponsse 882ad7c; ingen påstand om dagens Ponsse |
| [renderer-extraction-and-language-direction](renderer-extraction-and-language-direction.md) | Datert worktree-/gjenbruksinventar, med senere Go-retning |
| [handoff-mermaid-extraction](handoff-mermaid-extraction.md) | Historisk koordinering; ikke bestilling av ny uttrekksrenderer |
| [handoff-xfmd-sdui](handoff-xfmd-sdui.md) | Avløst FOX-handoff og referanse til separat dokumentnavigasjon |
| [prototype-widgets](prototype-widgets.md) | Statisk SVG/HTML-galleri fra før felles Go-layout; ikke ny runtime |

Alle 17 faglige dokumenter i denne katalogen er fordelt over. Genererte
[viewpoints](../design/README.md) kommer fra verktøyet og modellens fakta.
[Checkpoint #1](../../SDP/History/checkpoint-1/README.md) og G-fasebevis er daterte
snapshots; de oppdateres ikke for å ligne dagens kontrakter. Åpne oppgaver går
til [KanBan](../SDP/Agents/KanBan/README.md); redaksjonell konsolidering følges i
[KB-SDP-010](../../SDP/Agents/KanBan/active/%23010--Proposal--Document-consolidation.md).
