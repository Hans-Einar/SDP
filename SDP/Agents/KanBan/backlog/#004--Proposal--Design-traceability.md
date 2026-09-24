# Traceability mellom SDL-design, slices, kode og bevis

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-004 |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Behov

Design- og implementasjonsløpet må bruke Traceability-ledger aktivt. SDP tools
skal kunne koble SDL-design og vedtatt implementasjonsplan til faktisk leveranse,
og vise både nåstatus og hvordan systemet vokser gjennom roadmapens milepæler.
KanBan-historikk erstatter ikke denne dokumentasjonen.

## Kontrakt som må utredes

Gjenbruk og utvid eksisterende envelope, ID-er og relasjonskontrakter etter
inventar; ikke opprett konkurrerende implementasjonsledger. Modellobjektets
stabile ID og modellrevisjon må kobles til slice-/milepæl-ID, kode-/commitreferanse,
verifikasjonsresultat og bevis. Avklar også ugyldiggjøring når design eller bevis
endres, delvis leveranse, manglende dekning og kryssprosjektavhengigheter.

Hold foreslått, planlagt, implementert og verifisert atskilt. En grønn test eller
én levert Functionality innebærer ikke at en hel Feature er ferdig. Ukjent eller
utdatert grunnlag skal være synlig; status må ikke utledes av fritekst eller svake
`links`. Aggregasjonsregler og avgrensning må vises i rapporten.

## Neste arbeid og akseptanse

Kartlegg dagens Ledger.ndjson, Current-Index og Relations mot behovene. Definer
versjonert kontrakt og migrering før endring. Prøv én modell med to slices hvor
bare én er implementert/verifisert; generert roadmap/fremdriftsdiagram skal vise
forskjellen og lenke til kilden og bevisene. [KB-SDP-002 — sdptool: prosjektoppslag, implementasjonsplan og viewer](%23002--Proposal--sdptool.md) er konsument.

[Traceability/Ledger.ndjson](../../../Traceability/Ledger.ndjson)
[Toolkit/schemas/ledger-event.schema.json](../../../../Toolkit/schemas/ledger-event.schema.json)
[Traceability/Relations.yaml](../../../Traceability/Relations.yaml)
