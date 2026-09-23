---
id: KB-SDP-007
project: SDP
type: Change
created: 2026-09-23T18:25:13Z
source: owner-conversation-2026-09-23
---

# K1-M1: etablere KanBan og bevare samtalens forslag

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Oppdrag og avgrensning

Eierens bestilling 2026-09-23: etabler prosessen og katalogene, foreslå tags, og
registrer nyere samtale om SDP/SDL/SDUI med hovedkort og referanser. Dette er
leveransen for K1-M1; forslagene i de øvrige kortene skal ikke implementeres nå.

## Implementasjonsplan K1

Én fasebranch `sdp/phase-k1-kanban` fra G7-M1 (`d03eb78`). Én milepæl:

1. Etabler sju statuskataloger i hvert av tre prosjekter, stabile ID-er, indeks,
   kortmal, ledger-format og behandling ved flytting/avslutning.
2. Registrer seks hovedforslag og fire referansekort; merk kilde, åpne spørsmål,
   eierskap, neste behandling og kriterier uten å vedta ny språkprofil.
3. Kontroller JSON/schema, hendelsesforløp, plassering, identiteter og alle nye
   lokale Markdown-lenker. Dokumenter resultat, fullfør milepælen, commit og push.

Ingen gammel katalog, implementasjonsledger, Go-kode, GitHub-issue eller generert
SDL-viewpoint skal flyttes/endres. Ingen separat repo-opprettelse eller merge.

## Akseptanse og bevis

**K1-M1 levert lokalt 2026-09-23.** Tre tavler, sju statuskataloger per tavle,
seks hovedforslag, fire Ref-kort og dette leveransekortet er registrert.
Tags, arbeidsrytme, ID-/referanseregler, kortmal og versjonert ledger-payload er
beskrevet. Rotens agentinstruks og dokumentinnganger peker til prosessen.

Kontroll utført med Python 3 / jsonschema mot eksisterende Toolkit-envelope og
KanBan-payload: alle JSON-filer kan leses, alle 11 kort har unik og riktig ID/type,
alle Ref peker direkte på eksisterende hovedkort. Hendelsene er spilt gjennom
per kort med kontroll av forrige-hendelse, før/etter-status, sti og fysisk fil.
Alle 99 lokale Markdown-lenker i tavlene og berørte dokumentinnganger ble kontrollert.
`git diff --check` bestod. Kontrollene bestod også etter avslutningsflyttingen;
ledgeren har 12 hendelser, inkludert active → completed for dette kortet.
Dette er en engangskontroll av leveransen, ikke et levert generelt KanBan-verktøy.

Ingen Go-kode, gammel prosesskatalog, Traceability-ledger eller Toolkit-skjema er
endret. Ingen ny språkprofil, sdptool-kommando, graf eller repo-utskilling er
implementert. Disse temaene står eksplisitt som backlog. Commit/push avslutter
fasen i Git; denne filens milepælcommit er sporbar med `git log --follow`.
