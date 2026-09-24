# Konsolider dokumenter og skill aktiv kontrakt fra historikk

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-010 |
| project | SDP |
| type | Proposal |
| created | 2026-09-24T13:54:42Z |
| source | owner-conversation-2026-09-24 |
| next_review | Før neste fase-/malprofilleveranse i KB-SDP-001 |
| owner | Codex |

## Bestilling og omfang

Kartlegg og konsolider dokumenter i docs, checkpoint #1, SDL og SDUI. Det finnes
mange overlappende studier og planer, og enkelte eldre dokumenter omtaler levert
Go-funksjonalitet som uimplementert. Finn ett aktivt inngangsdokument per ansvar,
og skill språkprofil, prosess, implementasjonsstatus, forslag og daterte bevis.

Checkpoint #1 inneholder mest SDL/SDUI, men også SDP-prosess og felles beslutninger.
Det skal ikke flyttes til SDL og presenteres som en ren, gjeldende språkspesifikasjon.
Bevar opprinnelige bevis/fingeravtrykk og registrer hvor materialet hører hjemme.

## Sammenheng og neste behandling

[KB-SDP-001](../active/%23001--Proposal--Project-structure.md) eier fysisk struktur,
stimigrering og inventar. Dette kortet eier redaksjonell konsolidering av innhold;
ikke slå kortene sammen bare fordi de berører samme filer. Bruk R1-inventaret
som start, og vurder mindre språk-/prosessleveranser med lineage ved behov.

## Akseptanse

Alle kildedokumenter har eier og status. Motstridende påstander er gjennomgått
mot gjeldende profil/kode/bevis. Erstattede dokumenter peker til riktig etterfølger
eller arkiveres eksplisitt. Aktiv dokumentasjon og genererte viewpoints har
forskjellig autoritet; genererte filer bygges av verktøyet, ikke håndredigeres.
Ingen historisk test eller kandidat påstås å være bevis for nyere implementasjon.

## Aktiv leveranse R2

[Plan](../../../Maintenance/R2/Plan.md) avgrenser første runde til innganger,
autoritet og konkrete feil om hva som er implementert. [Funn og oppfølging](../../../Maintenance/R2/Findings.md)
viser hver rettelse og grunnlaget. Omnummerering/installasjonsprofil og
prosjektregister tilhører fortsatt KB-SDP-001. Ingen nye språkregler innføres.

## Arbeidslogg og revisjoner

| Tid | Aktør / hendelse | Arbeid og funn | Bevis / neste steg |
| --- | --- | --- | --- |
| 2026-09-24T14:55:42Z | Codex; EVT-KB-SDP-000023 | R2-M1: aktivert etter R1. Funnet foreldet SDL-sti, planlagt G4/G6 i leverte profiler og aktive Python-instrukser. | [Funn](../../../Maintenance/R2/Findings.md); R2-M2 retter og kontrollerer disse. |
| 2026-09-24T15:03:31Z | Codex; EVT-KB-SDP-000025 | R2-M2: kartlagte statuskonflikter behandlet, SDUI-dokumentkart opprettet og runtimeforslag erstattet som aktiv kontrakt av pakkekontrakter. | [Kontroller](../../../Maintenance/R2/Evidence.md); kandidatsemantikk og SDP-prosess-/malprofil gjenstår. |

## Neste behandling etter R2

R2s avgrensede statusrunde er ferdig. Kortet er fortsatt active: SDL-kjernens
aktive og foreslåtte deler deler ennå én stor kildefil, og prosessforslagene må
harmoniseres med vedtatt fase-/malprofil. Neste avgrensning er å behandle denne
profilen sammen med KB-SDP-001 og KB-SDL-001, og så konsolidere de berørte
dokumentene. R2 vedtar ikke nye SDL-keywords, A0–A5-nivåer eller installasjonsnavn.
