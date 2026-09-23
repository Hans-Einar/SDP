---
id: KB-SDP-001
project: SDP
type: Proposal
created: 2026-09-23T18:25:13Z
source: owner-conversation-2026-09-23
next_review: 2026-09-30
---

# Prosjektstruktur, Template og studier per fase

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

Lag et faktisk inventar av rotkataloger, Toolkit-maler/installasjonskontrakter,
aktive dokumenter og historiske bevis. Foreslå målstruktur og eierskap per filgruppe,
prosjektregistrering, underprosjektavhengigheter og flyttematrise for lenker,
Go-moduler, skript, CI og installasjon. Avklar submodule/subtree/annen innkobling
ut fra arbeidsflyt; ingen GitHub-repoer opprettes i denne leveransen.

## Akseptanse for senere omstrukturering

Godkjent migreringsplan; verifiserte lenker og bygg/installasjon; ingen dobbelt
aktiv mal- eller språkimplementasjon; samme prosjektoppslag med og uten Git-grense.
Historiske bevis beholdes som historie. Dagens `SDL/SDP/` er bare dokumentinngang;
SDL-koden ligger fortsatt i `SystemDesignLanguage/`.

## Relaterte grunnlag

[docs/SDL-Viewpoint-Levels-and-Notation.md](../../../../docs/SDL-Viewpoint-Levels-and-Notation.md)
[Toolkit/SDP-install.manifest.json](../../../../Toolkit/SDP-install.manifest.json)
[KB-SDP-002 — sdptool: prosjektoppslag, implementasjonsplan og viewer](%23002--Proposal--sdptool.md)
