# Avklar eldre Traceability-ID-er mot Toolkit-kontrakten

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-011 |
| project | SDP |
| type | Bug |
| created | 2026-09-24T14:06:02Z |
| source | R1-baselinekontroll, videreføring av kjent CI-avvik |
| next_review | Ved neste prioritering etter R1 |

## Observasjon

RepositoryValidation og contracts-CI feiler på ni eldre Issue #5-ID-er:
SPR-SDP-005, fire ITR-SDP-005-* og fire SLC-SDP-005-*. Schema krever Sprint-, SPI-
og SPS-format. R1s baseline og etterkontroll gir identisk feilsett.
[Baseline](../../../Maintenance/R1/toolkit-baseline.txt).

## Avgrensning og akseptanse

Avklar om dette er støttet ekstern ID-form eller data som skal migreres, og lag
én konsekvent løsning med alias-/sporbarhetsregler. Ikke omskriv append-only
ledgerhistorikk eller gamle bevis for å få grønn CI. Følg alle innkommende
referanser i records, bevis, issue og schema før endring. Full Toolkit-validator
og unittest-suite skal bestå uten unntak som skjuler feil. Dette er ikke fikset
som del av fysisk katalogopprydding.
