# Prosjektstruktur, Template og studier per fase

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-001 |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | Ved neste konsolideringsmilepæl i KB-SDP-010 |
| owner | Codex |

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Bakgrunn og eierretning

SDP skal bruke sin egen prosess: repoets `SDP/` skal eie utviklingsdokumentasjonen
for SDP-produktet. `Template/` skal etter planlegging eie malene andre prosjekter
bruker. Eksisterende nummererte rotkataloger må undersøkes før de flyttes;
malinnhold og faktisk prosjekthistorie må ikke blandes.

SDL og SDUI ønskes på sikt i egne GitHub-repoer, innkoblet i SDP. Verktøyene skal
fungere med samme underkataloger i monorepo og med separate repoer. Hvert prosjekt
har eget SDP-område, plan, ledger og KanBan. Git-mekanisme er ikke valgt.

## Prosess og nummerering — forslag som skal avklares

Behold nummererte kataloger, tettere knyttet til eksisterende A0–A5-abstraksjonslag.
Mandatet kommer fra Project Owner, som ikke nødvendigvis er utvikler. Mandate er
prosessinngang, ikke automatisk et nytt SDL-abstraksjonslag. Ikke erstatt A0–A5
eller likestill implementasjonsfaser med abstraksjonslag uten eksplisitt beslutning.

- `01--Mandate/Mandate.md`: eierens oppdrag.
- Valgfri `01--Mandate/Mandate-study.md`: forståelse, uklarheter og mulige måter å
  realisere oppdraget. Tolkning blir ikke automatisk godkjent mandat eller krav.
- Mandat og eventuell studie må samlet identifisere stakeholders.
- Egen `02--Study` foreslås fjernet som fast fase; valgfri studie per fase i stedet.
- `02--Requirements/`: valgfri Requirements-study, `01--Actors.design`,
  `02--UserStories.design`, `03--UseCases.design`. Videre nummerering er uavklart.
- Aktører/historier uttrykkes direkte i lesbar SDL; språkavklaringen eies av [KB-SDL-001 — Stakeholders, actors, user stories og lesbar SDL](../../../../SDL/SDP/Agents/KanBan/backlog/%23001--Proposal--Requirements-narrative.md).

## Neste avgrensede arbeid

R1s inventar og fysisk flytting er levert. Neste arbeid er en samlet fase-/malprofil
med studier per fase, kobling til A0–A5 og kravmodell; prosjektregistrering som
fungerer med og uten Git-grense; og en plan for eventuell SDL/SDUI-utskilling.
Dette må ta utgangspunkt i eksisterende Toolkit-installasjon og KB-SDL-001.
Submodule/subtree/annen innkobling er fortsatt uavklart; ingen nye repoer opprettes
bare fordi fysisk housekeeping er ferdig. KB-SDP-010 konsoliderer dokumentgrunnlaget.

## Akseptanse for senere omstrukturering

Godkjent migreringsplan; verifiserte lenker og bygg/installasjon; ingen dobbelt
aktiv mal- eller språkimplementasjon; samme prosjektoppslag med og uten Git-grense.
Historiske bevis beholdes som historie. Etter R1 ligger SDL-koden i `SDL/go`; `SDL/SDP/` eier språkprosjektets
prosessdokumenter. Samme prosjektoppslag på tvers av Git-grenser er fortsatt
et akseptansekriterium for kommende verktøyarbeid, ikke en levert R1-funksjon.

## Relaterte grunnlag

[docs/SDL-Viewpoint-Levels-and-Notation.md](../../../../SDL/docs/integration/SDL-Viewpoint-Levels-and-Notation.md)
[Toolkit/SDP-install.manifest.json](../../../../Toolkit/SDP-install.manifest.json)
[KB-SDP-002 — sdptool: prosjektoppslag, implementasjonsplan og viewer](../backlog/%23002--Proposal--sdptool.md)

## Aktivering 2026-09-24

Eieren ber nå om housekeeping av repo-roten og at dette kortet tas i arbeid.
R1 kartlegger eierskap, samler malene, flytter prosjektets egne driftsdokumenter
inn i SDP, samler SDL-området og etablerer dokumentinnganger/migreringskart.
Plan og avgrensninger føres i [R1-planen](../../../Maintenance/R1/Plan.md).
Full redaksjonell konsolidering av dokumentinnhold har eget kort KB-SDP-010.

## Etter R1 — utført og gjenstående

Maler er samlet i Template; prosjektets egne records ligger i SDP; SDL-kode og
språkdokumentasjon ligger nå i SDL. Checkpoint #1 er felles datert historikk i
SDP/History. Tidligere plasseringer omtalt over er forhistorien til denne migreringen.

Kortet forblir active for endelig fase-/malprofil (studier per fase og kravmodell),
prosjektregistrering som fungerer med og uten repo-grenser, og plan for separat
SDL/SDUI-repo. R1 oppretter ingen submodules eller nye GitHub-repoer.
KB-SDP-010 eier faglig dokumentkonsolidering; KB-SDP-011 eier eksisterende ID-avvik.
Disse holdes atskilt fra denne fysiske migreringen. Se R1-plan og kontrollbevis.

## Arbeidslogg og revisjoner

Denne loggen starter etter R1. Tidligere revisjoner finnes i Git, ikke som
nykonstruerte loggrader med gamle tidsstempler. [Historikk/diff](../History.md).

| Tid | Aktør / hendelse | Arbeid og resultat | Bevis / restarbeid |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000021 | R1 gjennomgått: M1 `d269bc7`, M2 `f42859e`, M3 `f722dc2` leverte fysisk organisering. | [R1-bevis](../../../Maintenance/R1/Evidence.md). Fase-/malprofil, prosjektregister og repo-utskilling gjenstår; dokumentkonsolidering følges i KB-SDP-010. |
| 2026-09-24T15:03:31Z | Codex; EVT-KB-SDP-000024 | R2 retter gjeldende neste steg og nåplassering; ingen ny migrering eller profil vedtas. | KB-SDP-010 samler dokumentgrunnlag; fase-/malprofil og prosjektregister er neste strukturarbeid. |
