---
id: KB-SDUI-001
project: SDUI
type: Ref
created: 2026-09-23T18:25:13Z
source: owner-conversation-2026-09-23
next_review: 2026-09-30
primary: KB-SDP-002
---

# SDUI som underprosjekt og bibliotek for SDP tools

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

Hoveddokument: [KB-SDP-002 — sdptool: prosjektoppslag, implementasjonsplan og viewer](../../../../../SDP/Agents/KanBan/backlog/%23002--Proposal--sdptool.md)

## Lokal påvirkning

SDUI skal ha eget SDP-område og KanBan, og kunne brukes som bibliotek/underprosjekt
med samme struktur i monorepo og senere separat repo. `sdptool` må kunne velge
SDUI-prosjektet eksplisitt. [KB-SDP-001 — Prosjektstruktur, Template og studier per fase](../../../../../SDP/Agents/KanBan/backlog/%23001--Proposal--Project-structure.md) eier selve strukturmigreringen.

Bevar nåværende Go/Fyne/SVG-retning. Dette er ikke en bestilling av FOX-vert,
ny parser, ny KanBan-GUI eller full Markdown/Mermaid-støtte. Grafidéen
[KB-SDP-003 — KanBan-graf med tidsakse og trinnvis detaljering](../../../../../SDP/Agents/KanBan/backlog/%23003--Idea--KanBan-graph.md) avgjør ikke valg av SDUI som vert.

## Behandling

Avklar lokal prosjektkonfigurasjon, modell-/plankilder og biblioteksgrenser før
migrering. Hold kodeplanens G-faser og denne idebacklogen tydelig atskilt.
[SDUI/docs/implementation-plan.md](../../../../docs/implementation-plan.md)
