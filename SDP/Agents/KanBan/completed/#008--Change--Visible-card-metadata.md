# K2-M1: synlig metadata i KanBan-kort

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-008 |
| project | SDP |
| type | Change |
| created | 2026-09-23T22:22:08Z |
| source | owner-conversation-2026-09-24 |

## Behov og omfang

Eieren rapporterer at YAML-frontmatter skjules i VS Code og tegnes som overskrifter
av XFMD. Metadata skal kunne leses som en vanlig Markdown-tabell.

## Implementasjonsplan K2

Ansvar: Codex. Fasebranch `sdp/phase-k2-readable-metadata` fra K1 (`bb3728c`).
K2-M1 konverterer metadata i alle eksisterende kort og malen til én synlig tabell
rett under tittelen. Bevar feltnavn, verdier, innhold, identiteter og eksisterende
ledgerhistorikk. Oppdater formatbeskrivelsen; ingen viewer-/parserendring.

## Verifikasjon og utfall

Kontroller at metadata og øvrig kortinnhold er bevart, at ingen kort bruker
frontmatter, og at tabellstrukturen, lenkene og ledgerforløpet er gyldige.
**K2-M1 levert 2026-09-24.** Tolv kort og malen har synlige metadatatabeller.
Sammenligning før/etter konvertering bekreftet identiske feltnavn/verdier i alle
13 dokumenter og uendret innhold i alle 12 kort. Malinstruksen er oppdatert til
rader i tabellen. Etterpå er dette kortets utfall ferdigstilt.

Kontroll av tre tavler, 21 statuskataloger, 12 kort, fire Ref-kort og 101 lokale
Markdown-lenker bestod, inkludert JSON-schema og replay av ledger mot filplassering.
Etter avslutning har ledgerne samlet 14 hendelser. `git diff --check` bestod.
Dette er en dokument-/strukturkontroll; native visning i VS Code og XFMD er ikke
prøvekjørt. Det er ingen endring i viewernes kode eller Markdown-parser.
