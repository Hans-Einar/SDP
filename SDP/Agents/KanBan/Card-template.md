# Kort tittel

| Felt | Verdi |
| --- | --- |
| id | KB-PROJECT-NNN |
| project | PROJECT |
| type | Idea |
| created | REPLACE_WITH_RFC3339_TIMESTAMP |
| source | REPLACE_WITH_CONVERSATION_OR_ISSUE_REFERENCE |
| next_review | YYYY-MM-DD |

Dette er en mal, ikke et registrert kort. Opprett fil og created-hendelse samlet.
For Ref: legg til raden `primary` med verdien `KB-OTHER-NNN` og Markdown-lenke til hovedkortet.
Valgfrie tema-tags: raden `tags` med verdien `process, tooling`. Status følger katalog og ledger.

## Behov og kilde

Hva prøver vi å løse? Skill eierens bestilling/beslutning fra anbefalinger.

## Forslag, omfang og åpne spørsmål

Hva er kjent, hva er ikke bestemt, og hvem/hvilke prosjekter berøres?

## Neste behandling og ferdigkriterier

Hva må avgjøres eller leveres? Ved active: ansvar, avgrensning og verifikasjon.
Ved onHold: grunn, gjenopptakingsbetingelse og neste vurdering.

## Utfall og referanser

Oppdater før avslutning. Lenk til beslutning, plan/slice, bevis eller etterfølger.
Ikke påstå produktimplementasjon bare fordi forslaget er vurdert eller planlagt.

## Arbeidslogg og revisjoner

Oppdater også mens kortet er active. Én rad per vesentlig behandling; behold
eldre rader og korriger med et tillegg. Gjeldende omfang og neste steg står over.
Ikke fyll inn egen commithash før commit; hendelses-ID kobler til Git-revisjonen.
Se [arbeidsmåten for historikk og diff](History.md).

| Tid (RFC3339) | Aktør / hendelse | Arbeid, funn eller beslutning | Bevis / restarbeid |
| --- | --- | --- | --- |
| Faktisk registreringstid | Aktør; EVT-KB-PROJECT-NNNNNN | Skill resultat fra forslag og eierbeslutning | Lenke og konkret neste steg |

## Opphav ved sammenslåing eller splitting

Bruk denne delen bare ved merge/split. Oppgi operationId og alle kilde-/mål-ID-er
med klikkbare lenker. Valgfrie metadatarader `sources` og `superseded_by` viser
ID-ene; ledgerens lineage-felt beskriver hele operasjonen.

| Kilde-ID og lenke | Videreført til mål-ID | Overført omfang | Erstattet utsagn og grunn | Restarbeid og hjem |
| --- | --- | --- | --- | --- |
| Fyll ut én rad per relevant kilde/mål | Nytt mål | Hva videreføres | Beslutning, dato og autoritet | I kilden eller navngitt målkort |

Se [opphavskontrakten](Lineage.md). Ikke lukk en kilde mens restarbeid mangler et hjem.
