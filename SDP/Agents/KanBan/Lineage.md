# Sammenslåing og splitting av KanBan-kort

Kontrakt 0.2, avklart 2026-09-24. Gjelder samtalekort og arbeidets opphav;
ikke automatisk kravoppfyllelse eller implementasjonsstatus.

## Arbeidsmåte

Ved utvalg til active gjennomgås beslektede kort, senere presiseringer og
motstridende forslag. Ny kunnskap kan endre hva neste avgrensede arbeid bør være.
En sammenslåing godkjenner ikke alle ideene i kildene. Skill eierbeslutninger,
anbefalinger, åpne spørsmål og dokumentert observasjon.

| Operasjon | Kilder | Resultat |
| --- | --- | --- |
| merge | Minst to eksisterende hovedkort | Ett nytt samlekort med ny ID |
| split | Ett eksisterende hovedkort | Minst to nye kort med hver sin nye ID |

Eksempel med **fiktive** kort: A + B → C; deretter C → D + E. Det er en graf med
sammenløp og forgreninger, ikke bare statusendringer langs én dokumentlinje.
En liten presisering av samme behov kan fortsatt oppdateres i eksisterende kort
med reviewed-hendelse; den trenger ingen oppkonstruert merge/split.
Mange-til-mange beskrives som to eksplisitte operasjoner via et samlekort.

1. Avgrens hva som videreføres fra hver kilde, hva som er erstattet og hvorfor,
   og hva som fortsatt er uavklart. Bevar nyttige fakta, datoer og bevis.
2. Opprett nye målkort med nye ID-er. Sett dem til active bare når deres omfang
   er valgt og autorisert; ellers backlog. Ved split kan målene få ulik status.
3. Sett en kildetabell i hvert målkort med klikkbare kilder. Sett utfall/etterfølgere
   i kildene. Metadata kan ha `sources` og `superseded_by`; ledgerens typede
   lineage-felt er autoritet for opphavsgrafen. `source` er fortsatt samtaleproveniens.
4. Flytt en fullt erstattet kilde til superseded. Ved split lenkes **alle** målene.
   En delvis overført kilde beholdes åpen med eksplisitt restarbeid, eller restarbeidet
   opprettes som et eget målkort før hele kilden kan avsluttes.
5. Registrer created per mål, og moved/reviewed per kilde som beskrevet nedenfor.
   Oppdater indekser og nåværende Markdown-lenker samlet med filene.
6. Kontroller at hele forløpet er representert og at ingen spørsmål er mistet.
   Dette er et faglig kontrollpunkt; schema kan ikke bevise tekstlig fullstendighet.

Kildene slettes aldri. Bevar tidligere utsagn med dato og marker hva som er
opphevet av hvilken beslutning. Superseded betyr at kortets videre behandling
har et annet hjem, ikke at funksjonen er levert. En Ref er ikke en opphavskilde;
bruk dens hovedkort, og oppdater lokal påvirkning separat. Historiske Ref-koblinger
skal ikke automatisk omskrives til ett tilfeldig mål ved split.

## Ledger-payload 0.2

[Skjema 0.2](ledger-payload-0.2.schema.json) beholder feltene og hendelsestypene
fra [0.1](ledger-payload.schema.json), og åpner for valgfritt `lineage`.
Envelope er fortsatt 1.0. Gamle 0.1-hendelser beholdes uendret. Nye hendelser
bruker payload 0.2; vanlig created/moved/reviewed trenger ikke lineage.

Hver deltakerhendelse i én merge/split har en **identisk** lineage-beskrivelse:

```json
{
  "operationId": "KBO-DEMO-000001",
  "kind": "merge",
  "sources": [
    {"id": "KB-DEMO-001", "scope": "Første idé", "remaining": ""},
    {"id": "KB-DEMO-002", "scope": "Senere presisering", "remaining": ""}
  ],
  "targets": [
    {"id": "KB-DEMO-003", "scope": "Samlet, avgrenset arbeid"}
  ]
}
```

`operationId` har eget fortløpende navnerom per koordinerende prosjekt:
`KBO-<PROJECT>-<nummer>`, minst seks sifre. Det er ikke en kort-ID eller event-ID.
`scope` beskriver overført omfang på kildesiden og mottatt omfang på målsiden.
`remaining` beskriver det som beholdes i kilden; tom streng betyr full erstatning.
Motstridende eller forkastet innhold begrunnes i kildetabellen og hendelsens reason.
Ved split fordeler målenes scope kildens omfang; samme del må ikke bli to parallelle
primæroppdrag ved et uhell. Skjemaet kontrollerer struktur; faglig review kontrollerer
fordeling og at restarbeidet er dekket.

### Hendelser og kontroller

- Kilder eksisterer før operasjonen. Mål er nye, unike ID-er. Kilder og mål er
  disjunkte, og ingen ID gjentas innen listene. Dermed danner nytt opphav ingen sykel.
- Hver kilde og hvert mål får nøyaktig én deltakerhendelse for operationId.
  Manglende hendelse gir **ufullstendig operasjon**, ikke ferdig sammenslåing/splitting.
- Målet får created med links til alle kilde-ID-er. Kilden får links til alle mål-ID-er.
  Disse generelle lenkene beholdes for enkle lesere; lineage avgjør relasjonstypen.
- Fullt erstattet kilde får moved til superseded. Delvis kilde har ikke-tom remaining,
  får reviewed uten status-/stiendring og må være i backlog, active eller onHold.
- Også full erstatning starter fra backlog/active/onHold. Historisk avsluttede kort
  kan siteres som grunnlag uten ny supersession. Hvis arbeidet faktisk gjenåpnes,
  registreres gjenåpningen først med egen begrunnelse.
- created-hendelsene for målene legges før kildehendelsene i samme ledger. Alle
  vanlige previousEventId-/sti-/statusregler gjelder fortsatt per kort.
- OperationId må ikke gjenbrukes med en annen definisjon eller ekstra deltakere.
  Detaljene skal være like i alle deltakerhendelser, også ved flere tavler.
- Innkommende lenker oppdateres til kildens nye plassering. Kilden leder deretter
  til etterfølgerne; ikke slett gamle identiteter fra historikken.

For én tavle commits filer, indekser og hendelser samlet. Ved flere tavler beholder
kortene sine prosjekt-ID-er; hver tavle registrerer sine egne hendelser med samme
operationId og definisjon. Koordiner alle endringene, og kontroller dem samlet.
Separate repoer har ikke en felles atomisk commit. Inntil alle deltakerhendelser
er tilgjengelige, skal en samlet leser vise operasjonen som ufullstendig.
En utilgjengelig tavle er ukjent grunnlag, ikke et tomt eller ferdig prosjekt.
Kryssrepo-transaksjoner, global låsing og automatisk ID-oppslag er ikke implementert.

## Eksempler og verifikasjonsgrense

[Eksempelledgeren](examples/lineage.ndjson) inneholder bare fiktive DEMO-kort:
full og delvis merge/split, inkludert kilder som beholder restarbeid.
Den ligger utenfor produksjonsledgeren og registreres ikke i boards.json.
[Kontrollen](examples/verify_lineage.py) prøver schema og semantiske forløp samt
negative tilfeller med manglende deltakere, sprikende operasjoner og tapte koblinger.
Det er en kontraktprøve av eksempler, ikke en generell KanBan-CLI, graf eller
validator av dokumentinnhold/filplassering i flere repoer.

Kjør fra repo-roten med Python 3 og jsonschema:

```sh
python3 SDP/Agents/KanBan/examples/verify_lineage.py
```
