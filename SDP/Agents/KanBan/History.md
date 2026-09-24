# Arbeidslogg, revisjoner og diff for KanBan-kort

Kortet er et levende arbeidsdokument også i `active`. Behold én fil med stabil
ID gjennom livsløpet. Git eier innholdsrevisjonene; KanBan-ledgeren eier
behandlingshendelsene; Traceability eier implementasjon og verifikasjon.
Dette er K4s manuelle arbeidsmåte, ikke et nytt versjonskontrollsystem eller en
implementert `sdptool history`-kommando.

## Mens kortet er aktivt

Kortets øverste deler viser gjeldende behov, avtalt omfang, ansvar, neste steg
og ferdigkriterier. Legg til en kort arbeidslogg med tidspunkt, aktør, hendelses-ID,
arbeid/resultat og bevis. Oppdater ved vesentlig funn, omfangsendring, beslutning,
verifikasjon og milepæl, også uten statusendring. Loggen er ikke en chattranskripsjon.

- Skill observert funn, agentforslag og eierbeslutning. Skriv hvem som vedtok hva.
- Bruk `x-kanban:reviewed` ved en behandling uten flytting, med samme sti/status
  og lenker til resultatet. Flere nært sammenhengende småendringer kan samles i én
  revisjon; rene stave-/lenkerettelser trenger ikke en egen ledgerhendelse.
- Før ny kunnskap inn i gjeldende tekst. Forklar hva som erstattes i en ny loggrad;
  ikke fjern tidligere begrunnelse fra loggen eller omskriv gammel ledger.
- Oppgi eksplisitt restarbeid ved hver milepæl. Et ferdig delarbeid lukker ikke
  hele kortet. Behold opprinnelig oppdrag og skille mellom levert og planlagt.
- Commit kort, relevante dokumenter, indeks og ledger samlet ved milepælen.
  Bruk både milepæl-ID og kort-ID i committeksten. Mellom milepæler er endringer
  lokale utkast; Git kan ikke vise mellomversjoner som aldri ble committet.

Arbeidslogg er et menneskelesbart sammendrag; gamle logger korrigeres med tillegg.
Ikke legg nye snapshotkopier av kortet i en revisjonsmappe. Ikke skriv commitens
egen hash inn i filen som inngår i commiten. Ledgerens `commit: null` er fortsatt
gyldig: Git-commiten som introduserer den unike event-ID-en identifiserer revisjonen.
Henvis gjerne til tidligere, allerede eksisterende commits i senere loggrader.

## Maintenance og de aktive kortene

`SDP/Maintenance/<fase>/` kan eie en sammenhengende faseplan, inventar og større
kontrollbevis. Kortet peker dit og summerer fremdrift, neste steg og restarbeid.
Det skal være mulig å forstå kortets situasjon uten å lete i tilfeldig plasserte
arbeidsnotater. Ikke kopier hele planen eller rå testutskrifter inn i kortet.
Ved avslutning bevares både kort og bevis. Maintenance er en leveranseadresse,
ikke en ny obligatorisk SDP-fase eller et nytt abstraksjonslag.

## Se historikken med Git nå

Kjør fra repoets rot. Dette konkrete eksemplet gjelder KB-SDP-001:

```sh
git log --follow --date=iso-strict --format='%h %ad %s' -- 'SDP/Agents/KanBan/active/#001--Proposal--Project-structure.md'
git log --follow -p -- 'SDP/Agents/KanBan/active/#001--Proposal--Project-structure.md'
git diff -- 'SDP/Agents/KanBan/active/#001--Proposal--Project-structure.md'
git diff --cached -- 'SDP/Agents/KanBan/active/#001--Proposal--Project-structure.md'
```

De to siste viser henholdsvis ustagede og stagede endringer. For historiske
revisjoner, finn commit og daværende sti; de behøver ikke være dagens sti:

```sh
git show 'bb3728c:SDP/Agents/KanBan/backlog/#001--Proposal--Project-structure.md'
git diff 'bb3728c:SDP/Agents/KanBan/backlog/#001--Proposal--Project-structure.md' 'f722dc2:SDP/Agents/KanBan/active/#001--Proposal--Project-structure.md'
git log --format='%h %s' -G '"eventId"[[:space:]]*:[[:space:]]*"EVT-KB-SDP-000015"' -- SDP/Agents/KanBan/Ledger.ndjson
```

`--follow` følger én fils antatte omdøpinger og er nyttig for vanlig flytting,
men Git gjetter renames fra likhet. Det er ikke autoritet for kortidentitet og
følger ikke semantisk merge/split. Ved tvil, slå opp ID-ens historiske stier i
ledgeren og sammenlign de eksakte `commit:sti`-blobene som over. Finn også eldre
ledgerstier via repoets migreringskart dersom selve tavlen er flyttet.

## Flere kilder, repoer og ufullstendig historikk

Merge/split følger [Lineage](Lineage.md): nye målkort har egne revisjoner og
bevarte kildekort har sine. En livsløpsvisning må vise dette som flere grener,
ikke late som alt var én fil. Ved delvis overføring fortsetter kildekortets logg.

Et fremtidig verktøy må slå opp prosjektets Git-repo fra prosjektregisteret,
deretter stabil kort-ID, ledgersti og revisjon. Bruk repoidentitet sammen med
commithash; en hash alene er ikke en kryssrepoadresse. Ref-kort peker til
hovedkortets historie og beholder sin egen lokale behandling.

Ucommittede utkast, manglende repo, grunn klone, manglende Git-objekt og avbrutt
historikk skal vises som nettopp det. Ikke rekonstruer innhold fra mtime eller
hevde at en manglende revisjon var tom. En eksport av bare Markdown inneholder
arbeidslogg, men ikke full Git-historikk. Rebase/squash kan endre commitidentiteter;
eierens fase-/milepælcommits skal bevares, mens kort- og event-ID-er er stabile.

Integrert historikk/diff følges i [sdptool-kortet](backlog/%23002--Proposal--sdptool.md)
og [grafidéen](backlog/%23003--Idea--KanBan-graph.md). Det er fortsatt backlog.
