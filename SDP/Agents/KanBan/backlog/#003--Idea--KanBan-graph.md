# KanBan-graf med tidsakse og trinnvis detaljering

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-003 |
| project | SDP |
| type | Idea |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Idé fra eieren

Vis hver statuskatalog som en vertikal kolonne. Tiden går nedover, i utgangspunktet
én rad per uke. Hvert dokument får et datapunkt når det opprettes eller flyttes,
og en sammenhengende sti gjennom statuskolonnene viser forløpet frem til nå.
Et klikk på en uke utvider dagene; videre klikk fokuserer på timene der det faktisk
var aktivitet. Punkter skal kunne åpne dokumentet og forklaringen til hendelsen.

## Grunnlag og avgrensning

Les append-only KanBan-ledger, ikke filenes mtime eller antatt Git-committid.
Bevar reelle tidsstempler og deterministisk rekkefølge ved samtidige hendelser.
Skille hendelser fra nåværende plassering og vise gjenåpning/tilbakeflytting.
Vis Ref som referanse til hovedkort, ikke en ekstra implementert leveranse.
Filtrering på prosjekt, type/tag og tidsrom kan gjøre grafen lesbar.

Interaktiv visning, tidszoom og valg av vert er åpne designvalg. Markdown/SVG kan
være en statisk eksport, men statisk Markdown alene oppfyller ikke klikkbar tidszoom.
Ikke bygg en ny renderer eller velg Fyne/XFMD/SDUI før et avgrenset forsøk er avtalt.

## Mulig første forsøk

Lag en statisk ukevisning fra en liten validert ledger med flyttinger og gjenåpning.
Kontroller status, stier, tidspunkt og lenker. Avklar deretter en interaktiv vert,
tidssone/ukegrenser, håndtering av tette hendelser og skjulte uker uten aktivitet.
Ingen graf er implementert i KanBan-grunnlaget.

## Presisering 2026-09-24: sammenløp og forgrening

Grafen skal kunne vise merge (flere kilder → ett mål) og split (én kilde → flere
mål) fra ledgerens typede lineage-felt. Ikke gjett disse forbindelsene fra fritekst
eller vanlige links. Delvis overføring viser at kilden fortsatt har restarbeid.
OperationId binder deltakerhendelsene sammen; en ufullstendig operasjon eller
utilgjengelig prosjekttavle merkes som ukjent/ufullstendig, ikke ferdigbehandlet.
Kildenes historikk og tidspunkter beholdes. [Kontrakt](../Lineage.md).
