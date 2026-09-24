# SDP — KanBan

## Kortoversikt

| ID | Type | Status | Dokument |
| --- | --- | --- | --- |
| KB-SDP-001 | Proposal | backlog | [Prosjektstruktur, Template og studier per fase](backlog/%23001--Proposal--Project-structure.md) |
| KB-SDP-002 | Proposal | backlog | [sdptool: prosjektoppslag, implementasjonsplan og viewer](backlog/%23002--Proposal--sdptool.md) |
| KB-SDP-003 | Idea | backlog | [KanBan-graf med tidsakse og trinnvis detaljering](backlog/%23003--Idea--KanBan-graph.md) |
| KB-SDP-004 | Proposal | backlog | [Traceability mellom SDL-design, slices, kode og bevis](backlog/%23004--Proposal--Design-traceability.md) |
| KB-SDP-005 | Ref | backlog | [SDL: links through i SDP-viewpoints](backlog/%23005--Ref--SDL--002--Links-through.md) |
| KB-SDP-006 | Ref | backlog | [SDL: kravmodeller i SDP-prosessen](backlog/%23006--Ref--SDL--001--Requirements-narrative.md) |
| KB-SDP-007 | Change | completed | [K1-M1: etablere KanBan og bevare samtalens forslag](completed/%23007--Change--KanBan-foundation.md) |
| KB-SDP-008 | Change | completed | [K2-M1: synlig metadata i KanBan-kort](completed/%23008--Change--Visible-card-metadata.md) |
| KB-SDP-009 | Change | completed | [K3-M1: sammenslåing og splitting](completed/%23009--Change--Card-merge-and-split.md) |

Indeksen vedlikeholdes sammen med flytting; ledgeren eier hendelseshistorikken.

## Formål og autoritet

Dette er prosjektets inngang for ideer, spørsmål og endringsønsker fra samtaler.
Opprett kort ved registrering, også når innholdet foreløpig er uavklart. Å skrive
et kort er ikke å vedta en språkregel eller autorisere hele implementasjonen.
Bevar forskjellen mellom eierens beslutning, agentens anbefaling og åpne valg.
Ikke kopier allerede leverte G-faser inn som nye ønsker.

`SDP/Agents/KanBan/` er valgt plassering, også i underprosjekter. Ingen ekstra
`SDP/KanBan/` opprettes. [Prosjektregisteret](boards.json) peker til de tre
lokale tavlene; [SDL](../../../SDL/SDP/Agents/KanBan/README.md) og
[SDUI](../../../SDUI/SDP/Agents/KanBan/README.md) har egne ID-serier og ledgere.
Registeret er en lokal kontrakt 0.1, ikke implementert sdptool-prosjektoppdagelse.
Underrepoene er ennå ikke skilt ut. Relative lenker fungerer i dagens katalogtre;
ved flytting/utskilling må register, policy og lenker migreres og verifiseres.

## Typer/tags

Første tag er kortets `type` og del av filnavnet. Valgfrie `tags` i metadatatabellen
kan angi tema som `tooling`, `language`, `process` eller `visualization`.
Type beskriver innhold, ikke prioritet eller arbeidsstatus.

| Type | Bruk |
| --- | --- |
| Idea | Mulighet som skal vurderes; ingen ferdig løsning kreves |
| Proposal | Konkret forslag med hensikt, åpne valg og akseptansekriterier |
| Question | Avklaring som trenger et dokumentert svar |
| Study | Avgrenset undersøkelse, med spørsmål og forventet resultat |
| Change | Bestilt endring med kjent omfang |
| Bug | Observert avvik med forventet atferd og reproduksjon |
| Decision | Beslutning med beslutter, dato, grunnlag og konsekvenser; tag alene er ikke godkjenning |
| Ref | Lokal påvirkning og lenke til ett hovedkort i et annet prosjekt |

Bruk `UserStory` og andre modellbegreper i SDL når profilen støtter dem;
KanBan-tags er ikke nye SDL-keywords. Unngå nye typer uten et eget behov.

## Identitet og dokumenter

Filnavn: `#003--Idea--KanBan-graph.md`. Referanseeksempel:
`#001--Ref--SDP--002--sdptool.md`. Tall er fortløpende per prosjekt på tvers av
status og type, minst tre sifre. Ikke gjenbruk slettede/avsluttede nummer.
Stabil ID er `KB-<PROJECT>-<nummer>`, eksempelvis `KB-SDP-003`; navn og plassering
kan endres uten at ID gjør det. Sjekk alle mapper og ledger før nummer tildeles.
Parallelle registreringer må løse ID-/eventkollisjoner før commit, ikke overskrive.

[Mal](Card-template.md) angir minste innhold. Metadata skrives som en vanlig
Markdown-tabell rett under tittelen, med kolonnene `Felt` og `Verdi`. Behold
feltnavnene `id`, `project`, `type`, `created`, `source` og eventuell `next_review`,
`primary` og `tags`. Verdiene skal være synlige i Markdown-visere med tabellstøtte.
Ikke legg en ekstra kopi i YAML-frontmatter; det gir ulik visning og to kilder
som kan komme ut av takt. Ledgeren beholder JSON-formatet og hendelseshistorikken.
Hvert hovedkort eier ett sammenhengende
behov. Ref har eget nummer/status, `primary` med hovedkortets stabile ID, klikkbar
Markdown-lenke og lokal påvirkning. Ref peker direkte til hovedkort, ikke en Ref-kjede.
Et hovedkort kan registreres i hvilken som helst tavle; velg nærmeste faglige eier
når det er praktisk. Ikke dupliser behovet for å oppnå perfekt plassering.

Lokal ferdigbehandling av Ref lukker ikke hovedkortet. Ved behov opprettes egne
implementasjonsslices med tilbakekobling til hovedkortet. Gi kortets ID i eventuelle
GitHub-issues/PR-er; GitHub-status er ikke automatisk lokal KanBan-status.

## Statuskataloger

| Katalog | Betydning og krav |
| --- | --- |
| backlog | Registrert og venter på prioritering eller avklaring; ha konkret neste vurdering |
| active | Avtalt, avgrenset arbeid pågår; oppgi omfang, ansvar og ferdigkriterier |
| onHold | Beholdes, men blokkert/utsatt; oppgi grunn, utløsende betingelse og vurderingsdato |
| completed | Kortets avtalte utfall er oppnådd, med lenke til beslutning/leveranse/bevis |
| canceled | Arbeid som var aktuelt eller besluttet, er aktivt avbrutt; begrunn valget |
| superseded | Fullt erstattet, slått sammen eller splittet; alle etterfølger-ID-er og lenker kreves |
| irrelevant | Vurdert som utenfor behov/omfang eller ikke lenger relevant; begrunn vurderingen |

Katalognavnene er case-sensitive; bruk `superseded`, ikke `superseeded`.
En Idea/Question/Study kan bli completed når den avtalte vurderingen er levert.
Det betyr ikke at produktfunksjonen er implementert. Ved godkjent viderearbeid:
lenk til ny plan/slice; implementasjonsstatus hentes fra Traceability. For et kort
som faktisk lover implementasjon, krever completed dokumentert implementasjon
og avtalt verifikasjon. Ikke lukk det bare fordi en plan finnes.

Alle statuser kan gjenåpnes med begrunnelse; historikken beholdes. Ikke slett gamle
kort for å få tom backlog. Duplikater flyttes til superseded og peker på hovedkortet.
Sammenslåing og splitting følger [opphavskontrakten](Lineage.md): nye målkort,
bevarte kilder og eksplisitt restarbeid. Bare fullt erstattede kilder avsluttes;
delvise kilder beholdes åpne. Vurder beslektede backlogkort før valg til active.

## Arbeidsrytme og omfang

Ved oppstart av arbeid: les aktuell tavle og berørte Ref-kort, velg et avgrenset
kort og flytt det til active med eksplisitt omfang. Ved nye funn som endrer oppdraget:
registrer eller oppdater et kort og lenk fra pågående arbeid før du bytter retning.
Registrering gir ikke i seg selv autorisasjon til å implementere det nye omfanget.

Ved milepælslutt: oppdater utfall, ledger og referanser. Før neste fase og ved
avtalt vurderingsdato: gjennomgå backlog/onHold og velg konkret neste arbeid,
utsatt vurdering, canceled, superseded eller irrelevant. Standard første
vurderingsfrist for disse nyregistrerte kortene er 2026-09-30; det er ikke en
leveransefrist. Nye frister settes ut fra faktisk prosjektbehov. Aldring alene
skal ikke automatisk slette eller avvise en idé. Målet er få, tydelige og behandlede
kort, ikke et permanent arkiv av uavklarte ønsker i backlog.

## Flytting og lenker — manuell arbeidsflyt i første versjon

1. Les gjeldende kort og siste hendelse. Noter ID, nåstatus og gammel sti.
2. Dokumenter grunn/utfall, eventuelt etterfølger, plan og bevis i kortet.
3. Flytt hele filen til ny statuskatalog; behold ID og normalt filnavn.
4. Append én hendelse i `Ledger.ndjson` med reell tid, aktør, gammel/ny status og sti.
5. Oppdater tavleindeks og alle innkommende Markdown-lenker i registrerte tavler.
   Søk på stabil ID og gammelt filnavn; prosentkod `#` som `%23` i lenkemålet.
   Ikke skriv om gamle ledgerhendelser eller deres historiske stier.
6. Kontroller lenker og at replay av ledger stemmer med fysisk plassering, og
   commit dokument/flytting/ledger/indekser samlet. Ved avvik: stopp og avklar;
   ikke gjett status fra mtime. Det er ingen automatisk flyttekommando ennå.

Varig dokumentreferanse er ID; Markdown-lenken er nåværende lokal åpneadresse.
Flytting innen samme tavle beholder relativ dybde. For flytting mellom prosjekter
beholdes opprinnelig hovedkort/ID foreløpig og nytt prosjekt får Ref; en generell
transfer-/ID-migreringskontrakt er ikke implementert.

## Ledgerkontrakt — payload 0.1 og 0.2

Hver tavle har [board.json](board.json) og append-only `Ledger.ndjson`, én JSON-
hendelse per linje. Gjenbruk eksisterende
[SDP event-envelope](../../../Toolkit/schemas/ledger-event.schema.json) med
`schemaVersion: "1.0"`. Nye hendelser bruker [payload 0.2](ledger-payload-0.2.schema.json);
historiske hendelser beholder [payload 0.1](ledger-payload.schema.json).
Dette endrer ikke Toolkit-skjemaet eller implementasjonsledgeren i Traceability.

- `eventId`: `EVT-KB-<PROJECT>-<løpenummer>`, unik og fortløpende innen tavlen.
- `eventType`: `x-kanban:created`, `x-kanban:moved` eller `x-kanban:reviewed`.
- `subjectId`: stabil kort-ID. `occurredAt`: faktisk RFC3339-tid med tidssone,
  her UTC. `actor`: hvem som registrerte. `commit` er null ved registrering;
  Git-commiten som introduserer linjen dokumenterer hendelsen. Ikke skriv om
  historikk bare for å sette hendelsens egen commithash etter commit.
- Payload: `schemaVersion`, `projectId`, `previousEventId` (for samme kort),
  `from`, `to`, `fromPath`, `toPath`, `reason`, `links`. Versjon 0.2 tillater også
  `lineage` med operationId, merge/split, kilder, mål og overført/gjenstående omfang.
  Se [opphavskontrakten](Lineage.md) for deltakerhendelser og fullstendighetskrav.
- Stier er bokstavelige UTF-8-stier relativt til tavlen, uten `..` eller absolutt
  prefiks. De er historiske hendelsesdata, ikke URL-er. `links` inneholder stabile
  kort-/slice-ID-er eller bevis-/beslutningsstier relativt til tavlen.
- created: første hendelse, tidligere event/status/sti er null. Registrering kan
  starte i backlog eller active; eldre forløp skal ikke oppdiktes.
- moved: `previousEventId` treffer siste hendelse; from/fromPath treffer forrige
  tilstand; ny sti/status beskriver faktisk flytting eller navneendring.
- reviewed: samme status og sti, med resultat/begrunnelse og eventuell ny frist
  i kortet. Ingen skjult flytting. Korrigering av gjeldende tilstand registreres
  som en ny begrunnet hendelse; tidligere feil kan forklares, ikke slettes.
- completed og superseded krever minst én utfalls-/etterfølgerreferanse i `links`
  ved overgangen. For superseded må alle etterfølgere være navngitte, eksisterende hovedkort.

Replay følger linjerekkefølge og forrige-hendelse-kjeden, ikke bare klokkeslett.
Dette gir entydig forløp også når flere hendelser har samme tidsstempel. Ved
Git-merge må nye, ikke-integrerte event-ID-kollisjoner løses før integrasjon.
Siste hendelse må samsvare med nøyaktig én fil for kortet i oppgitt katalog.

Skjemaene validerer envelope/payload. Kjederegler, referansemål, statusoverganger,
filplassering og lenker må også kontrolleres; ingen generell validator eller
interaktiv graf er levert her. K1s konkrete kontrollresultat dokumenteres i
[leveransekortet](completed/%23007--Change--KanBan-foundation.md).
