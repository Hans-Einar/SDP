# sdptool: prosjektoppslag, implementasjonsplan og viewer

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-002 |
| project | SDP |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Eierens ønskede bruk

```sh
sdptool ~/git/XFMD generate ip
sdptool generate ip
sdptool view ip
```

Valgfri prosjektsti kommer før kommandoen. Uten sti brukes `.`. Undersøk først
om valgt katalog selv er et gyldig SDP-område, ellers om den har et gyldig `SDP/`.
Ikke gjett foreldreprosjekt eller skift til underprosjekt uten eksplisitt regel.
Et versjonert prosjektmerke foreslås for gjenkjenning; format er ikke vedtatt.
`ip` er foreslått kortform for `implementation-plan`.

`view ip` skal åpne planens hovedside og navigator i konfigurert SDP/SDL-viewer,
for eksempel XFMD, med kilde/prosjekt/verktøy riktig registrert. Gjenbruk dagens
kildebaserte dokumentgenerering, viewer-adapter og ferdigbygde programmer.
Generer detaljer fra gjeldende kilder ved behov; ingen obligatorisk full eksport,
kompilering ved oppstart eller daemon. Aktuell navigator kan være et snapshot;
endringer i modellens struktur må få en eksplisitt oppdateringsmekanisme.

## Ansvar og planforslag

SDL er språket, SDP er prosessen. SDP tools skal bruke SDL-/SDUI-biblioteker,
ikke duplisere parser eller runtime. Prosess-CLI skal fungere likt for monorepo
og underprosjekter i separate repoer. Katalog-/Git-migrering eies av [KB-SDP-001 — Prosjektstruktur, Template og studier per fase](../active/%23001--Proposal--Project-structure.md).

Støtt både å foreslå vertikale slices fra SDL-modell og å validere manuelt eller
agentforeslåtte slices. En slice skal gi en avgrenset, verifiserbar evne gjennom
relevante lag/containere, med eksplisitte forutsetninger og akseptanse.
Verktøyet kan kontrollere dekning og avhengigheter; forretningsprioritet og nytte
må komme fra prosjektet. Foreslåtte planer blir ikke automatisk vedtatt.
`generate ip` må bevare godkjent plan og vise hva som er foreslått endret.
Planformat, beslutningsprosess og hvordan mål/rammer oppgis, gjenstår å avklare.

Bruk modell, vedtatt plan og Traceability til å vise roadmap og gradvis vekst i
bruksmål, features og functionality. Kontrakten for status og bevis eies av
[KB-SDP-004 — Traceability mellom SDL-design, slices, kode og bevis](%23004--Proposal--Design-traceability.md). Svake `links` må aldri telle som implementasjonsbevis.

## Neste leveranse og akseptanse

Avtal kommandokontrakt og prosjektoppslag før implementasjon. Første vertikale
prøve bruker én eksisterende SDL-modell, en liten slice-plan og ledgeren; viser
forslag/valideringsdiagnoser og kildekoblet fremdrift uten oppdiktet status.
`view ip` skal virke fra repoets rot og eget SDP-område samt med eksplisitt sti.
Ukjent prosjekt/ugyldig kilde får diagnose; kilde/vedtatt plan overskrives ikke.

## Tillegg: kortets revisjoner og diff

Eieren ønsker 2026-09-24 å kunne følge også innholdsendringer mens kort er active.
Fremtidig `sdptool` bør tilby historikk og diff ved stabil kort-ID, både committed
revisjon mot revisjon og arbeidsutkast mot siste commit. Kommandonavn er ikke vedtatt.
Gjenbruk Git; ikke lag separat repository/versjonsmotor for hvert kort.
Vis status-/reviewhendelser sammen med Git-revisjoner, og følg eksplisitt lineage
til kilde-/målkort ved merge/split. Ikke bruk rename-heuristikk som kortidentitet.
Prosjektregister må løse samme oppgave innen monorepo og etter repo-utskilling.
Akseptanse: korrekt historisk sti/innhold og diff ved flytting, samme-status-revisjon,
merge/split og eksplisitt diagnose ved manglende historikk. [Manuell arbeidsmåte](../History.md).

## Nåstatus for verktøyet

Dette er planlegging, ikke en eksisterende `sdptool`-kommando. `sdl-design` og
SDLs Go-CLI er eksisterende gjenbruksgrunnlag. [SystemDesignLanguage/scripts/README.md](../../../../SDL/scripts/README.md)

## Arbeidslogg og revisjoner

| Tid | Aktør / hendelse | Behandling | Restarbeid |
| --- | --- | --- | --- |
| 2026-09-24T14:52:35Z | Codex; EVT-KB-SDP-000019 | Registrerer integrert korthistorikk/diff som fremtidig sdptool-funksjon; manuell Git-arbeidsmåte leveres i K4. | Avtale vertikal verktøyleveranse før implementasjon. |
