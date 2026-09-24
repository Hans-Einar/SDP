# R2 — konsolider dokumentautoritet og implementasjonsstatus

KB-SDP-010, første avgrensede leveranse etter fysisk housekeeping i R1.
Branch `sdp/phase-r2-document-consolidation` fra K4 `1713778`.

## Mål og grenser

En leser skal finne gjeldende språkprofil, kjørbare innganger, bevis og åpne
forslag uten å forveksle dem. R1s inventar av 50 docs-filer gjenbrukes. R2
supplerer med SDUI-dokumentenes eierskap og gjennomgår konkrete motstridende
nåstatuspåstander mot eksisterende kode, profiler og milepælbevis.

Ingen ny syntaks/runtime, generering av viewpoints for hånd, nye repoer eller
omnummerering av maler. Ingen endring i historiske testlogger, hashmanifest,
fryste testdata eller checkpoint-snapshot. Gamle kandidatstudier kan beholdes
med eksplisitt avgrensning; uløste faglige spørsmål skal fortsatt være synlige.

## Milepæler

| ID | Leveranse | Verifikasjon | Status |
| --- | --- | --- | --- |
| R2-M1 | Aktiver 010, kartlegg konkrete konflikter og avgrens autoritet | Kort/ledger/lenker, kildehenvisninger til profil/kode/bevis | Levert |
| R2-M2 | Konsolider leserveier og rett de kartlagte statuskonfliktene | Dokumentdiff, profil-/CLI-stikkprøver, lokale lenker og bevarte historiske/genererte bytes | Levert |

[Funn](Findings.md) er behandlingstabellen; [bevis](Evidence.md) registrerer
utførte kontroller. Aktivt kort har arbeidslogg med ledgerhendelser. Hver milepæl
får egen commit; fasen pushes etter M2. Kortet lukkes bare når hele dets avtalte
omfang er håndtert, ikke fordi R2s avgrensede statusrunde er levert.

## Neste faglige arbeid etter R2

Skille gjenværende kandidatsemantikk fra aktive deler av den blandede SDL-
språkdefinisjonen, samordne SDP-prosessforslag og ferdigstille fase-/malprofil.
Dette må koordineres med KB-SDP-001 og KB-SDL-001; dokumentredigering skal ikke
stille vedta nye språk- eller installasjonsregler. KB-SDP-011s gamle ID-feil er
separat arbeid. Integrert KanBan-historikk/graf ligger i KB-SDP-002/003.
