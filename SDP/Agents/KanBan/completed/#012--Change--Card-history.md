# K4-M1: levende kort, revisjonslogg og Git-diff

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-012 |
| project | SDP |
| type | Change |
| created | 2026-09-24T14:52:35Z |
| source | owner-conversation-2026-09-24 |
| owner | Codex |

## Bestilling og avgrensning

Eieren ønsker at kortene dokumenterer arbeidet mens de er active, og at vi kan
se revisjoner og diff gjennom livsløpet. Eksisterende Git gjenbrukes. Maintenance
beholdes for faseplaner/inventar/bevis med tilbakekobling til det aktive kortet.
Ingen separat Git-motor, automatisk snapshotting eller viewer implementeres her.

## Plan K4

Branch `sdp/phase-k4-card-history` fra R1 `f722dc2`. Én milepæl K4-M1:
arbeidsmåte og kortmal, synlig logg på aktivt kort, prøvbar historikk/diff og
registrerte akseptansekriterier for fremtidig verktøy/graf. KB-SDP-010 fortsetter
som neste dokumentoppdrag; ingen andre backlogforslag implementeres ved registrering.

## Resultat og bevis

[History](../History.md) definerer manuell arbeidsmåte. [Kortmal](../Card-template.md)
har arbeidslogg; KB-SDP-001 oppsummerer levert arbeid og restarbeid. Historikk/diff
føres videre i KB-SDP-002/003. Dokumentene i SDL/SDUI bruker samme felles arbeidsmåte.
Ingen payload-schema eller gamle ledgerhendelser er endret.

Git-eksemplene er kjørt mot KB-SDP-001: `--follow` viser K1, K2 og R1; eksakt
blob-diff fra `bb3728c` backlog til `f722dc2` active viser 32 tillegg / 11 slettinger.
Oppslag på eventId 000015 identifiserer `d269bc7`. K3s lineage-prøver består med
full/delvis merge/split, gammel payload og 15 negative tilfeller.
Schema/replay og plassering består for 16 kort / 26 hendelser på tre tavler.
Ledgerprefikser er bytebevart mot R1 `f722dc2`; 1900 lokale Markdown-filmål
består R1-kontrollen, sammen med 574 genererte artefakter og migreringsgrensene.
`git diff --check` består.
Dette beviser dokumentarbeidsmåten; ingen integrert historikkvisning er levert.

## Arbeidslogg og revisjoner

| Tid | Aktør / hendelse | Arbeid og resultat | Restarbeid |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000018 | Registrerer og avgrenser K4 fra eierens bestilling. | Manuell historikk og aktiv arbeidslogg. |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000022 | K4-M1 levert som dokumentert, kontrollert arbeidsmåte. | Verktøy/graf eies fortsatt av KB-SDP-002/003. |

Kortet har én committed leveranserevisjon; created/moved i samme milepæl er to
behandlingshendelser, ikke to separate Git-snapshots av innholdet.
