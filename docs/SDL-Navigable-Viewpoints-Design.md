# SDL — navigerbare viewpoints og generering ved behov

**Dato:** 2026-09-22 · **Status:** planlagt G6, ingen implementert lenkehandler,
daemon eller XFMD-utvidelse. Eierens bestilling er navigasjon som Markdown,
små viewpoint-sider og senere to Markdown-paneler i XFMD.
[Felles faseplan](../SDUI/docs/implementation-plan.md) eier milepælene;
[SDL-modellen](../SDUI/design/architecture.design) eier ansvar og scenarioer.

## Ett grunnlag, flere dokumenter

SDL-kilde → validert modell → valgt viewpoint/projeksjon → dokumentseksjoner →
Markdown med Mermaid eller SVG. Full rapport, enkeltside og navigator bruker
samme prosjektor og kildekart. Navigatoren beskriver et utvalg; den skal ikke
skrape overskrifter fra en stor generert Markdown-fil. En framtidig innholdsmodell
med stabile seksjons-ID-er kan gjenbruke seksjoner i begge produktene.

Foreslått statisk pakke (ikke dagens utdataformat):

```text
design-docs/
  index.md
  navigator.md
  implementation.md
  viewpoints/
    index.md
    VP02/
      index.md
      architecture.md
    VP06/
      index.md
      G1FrontendPort.md
      G6NavigableDocumentation.md
    VP08/
      index.md
      SelectedViewOpened.md
  assets/
    <content-hash>.svg
  manifest.json
```

Hver side har tittel, tilbake-/overordnetlenke, kilderevisjon og lenker til
relaterte modellobjekter. Stabile ID-er, ikke overskriftens tekst, bestemmer
ankre og filnavn. Navn escapes for Markdown/URI og kan ikke bli filbaner uten
kontroll. Relativ bildesti regnes fra siden som bruker den. Manifestet eier
genererte filer; brukerens notater skal bevares ved reeksport.

Samlet viewpoints.md/printout.md kan velges som ekstra eksport. De trenger
ikke produseres ved hvert klikk. Vanlige relative lenker fungerer uten SDL-
tjeneste; det er første milepæls nyttige leveranse.

## Navigatorlenke og vertsadapter

Kandidatnotasjon, verken standard Markdown-funksjon eller implementert URI-skjema:

```markdown
[G1 — frontend](sdl-view://sdui-design/VP06?focus=G1FrontendPort&target=main&consumer=xfmd)
[Statisk G1-side](viewpoints/VP06/G1FrontendPort.md)
```

XFMDs lenkehandler oversetter første lenke til en typet ViewSelection.
Prosjekt-ID slås opp i vertens konfigurerte prosjekter. Viewpoint, fokusobjekt,
eventuell modus og støtteprofil valideres mot valgt SDL-revisjon.
Ukjent utvalg gir diagnose og beholder dokumentet som allerede vises.

`consumer=xfmd` velger en lokalt registrert adapter med fast program og
argumentliste. Lenken inneholder ikke en shell-kommando eller vilkårlig
etterkommando. Dette oppfyller ønsket om å velge hvilket program som åpner
resultatet; flere Markdown-lesere kan registreres med samme kontrakt.
Vanlige lesere bruker den statiske lenken. Endelig URI-navn og CLI fastsettes
med XFMD-implementasjonen; disse eksemplene er forslag.

Verten legger til request-ID, sesjon, konkret vindu og panelet som skal motta
resultatet. `target=main` lar navigatoren stå i sidepanelet. En lenke med
`target=navigation` kan senere bytte navigasjonsside. Historikk og scroll
holdes per panel. G6 gjelder dokumentvisning; Fyne er fortsatt første vert
for kjørbare SDUI-widgets.

## Fra klikk til dokument

1. XfmdDocumentHost fanger vindu/panel og sender SelectViewRequest.
2. DocumentBroker validerer utvalget og velger et konsistent kildesnapshot.
3. SdlViewpointGenerator bygger bare valgt dokument og nødvendige ressurser.
4. ViewArtifactStore publiserer komplett pakke og gir dokumentreferanse/lease.
5. ViewerLaunchAdapter sender referansen til den valgte leseren.
6. Leseren bekrefter åpning; broker svarer på den opprinnelige forespørselen.

De samme stegene står i SDL-scenarioet SelectedViewOpened og genereres til VP08.
InvalidViewSelectionRejected og ViewProjectionFailed viser feil før åpning.
Publiseringsfeil, bortfalt vindu og avbrudd skal også testes i G6; scenarioene
er eksempler og påstår ikke full feildekning.

Foreslått vertsgrense er `OpenDocument(requestId, entry, revision, lease,
windowId, pane)`. Kontraktfeltene i SDL beskriver argumentgrensene, men kjører
ikke verifikasjon av feltverdier eller transport. Go-typer og IPC-versjon
fastsettes ved implementasjon. Prosjektor importerer verken XFMD eller
prosesslauncher; DocumentBroker og ViewerLaunchAdapter eier koordineringen.

`xfmd --active <entry>` er ønsket CLI-adferd. Den undersøkte lokale
XFMD-inngangen (`src/application/main.cpp`) dokumenterer bare en fil/katalog,
`--help` og `--version`; ny vindus-/IPC-støtte må implementeres i XFMD.
Det aktive XFMD-vinduet må velges når forespørselen mottas. Ved klikk inne i
navigatoren brukes avsendervinduets ID. Senere fokusendring må ikke omdirigere
resultatet. Eksplisitt `--window <id> --pane main` er foreslått adapterform.
Lukket mål gir en feil til samme forespørsel; åpning av nytt vindu er en eksplisitt
vertsinnstilling. Global vindusfokusdeteksjon er ikke påkrevd for panelklikk.

## Minne, publisering og levetid

En daemon er nyttig for gjenbruk av AST, prosjektorresultater og ressurser,
men er ikke nødvendig for første enkeltforespørsel. G6-M2 kan være én CLI-prosess;
G6-M4 legger til en brukereid bakgrunnstjeneste med lokal IPC.

Et objekt i Go-minne blir ikke en vanlig fil bare ved å gi det et navn.
En kompatibel filbasert leser trenger en faktisk lesbar dokumentpakke.
Første adapter publiserer små, behovsgenererte pakker under brukerens
`$XDG_RUNTIME_DIR/sdl/views/<session>/<revision>/<selection>/`, med entry.md,
SVG-er og manifest. Dette er midlertidige runtime-filer, ikke en eksport
av alle viewpoints til arbeidsrepoet. Ren minneoverføring krever en separat
stream-/IPC-konsument og er utsatt; ingen FUSE eller pseudo-fil trengs nå.

[XDG-spesifikasjonen](https://specifications.freedesktop.org/basedir/latest/)
definerer en brukereid runtime-katalog med modus 0700 og sesjonsbundet levetid.
Den garanterer ikke at innholdet alltid ligger i RAM. Bruk ikke en global
/run/SDL-katalog. Manglende gyldig runtime-katalog gir eksplisitt fallback til
en privat midlertidig katalog, med melding om valgt plassering. Begrens både
pakkestørrelse og totalcache; store arkiver eksporteres eksplisitt et annet sted.

Publisering bruker nye immutable revisjonskataloger og atomisk ferdigstilling
før OpenDocument. Pakken omfatter også relative bilder og lenkemål som trengs
for siden. Delvise filer vises aldri. Cache-nøkkel inkluderer kilde-/importhasher,
profil, prosjektorversjon, utvalg og renderer/tema. Endret kilde invaliderer
de berørte resultatene; ukjent avhengighet gir full invalidasjon.

Leserens lease beholder pakken mens den kan laste eller reload-e ressurser;
det holder ikke å slette rett etter open-ack. Frigjøring skjer ved dokumentbytte/
lukking. Generiske lesere uten lease-støtte får pakker beholdt til eksplisitt
opprydding/sesjonslutt; ved fullt budsjett avvises nye forespørsler framfor å
slette synlige ressurser. Timeout, krasj, kvoter og opprydding testes i G6-M4.

Raskere andre klikk gjør første resultat foreldet for samme vindu/panel.
Broker kontrollerer request-ID og revisjon før åpning. Kilde-/renderfeil beholder
siste gyldige visning og gir diagnose. Oppdatering av navigator og hoveddokument
viser hvilken revisjon de bruker. Ingen SDL-domenehandling utføres ved navigasjon.

## Ansvar og akseptanse

| Grense | Eier |
| --- | --- |
| Viewpoint-/seksjonsutvalg, Markdown, navigasjon og kildekart | SDL-verktøyet / SdlViewpointGenerator |
| Revisjon, kø, avvisning av gamle resultater og cache-nøkler | DocumentBroker |
| Pakker, fullstendig publisering, leases og opprydding | ViewArtifactStore |
| Valg av registrert leser og åpneprotokoll | ViewerLaunchAdapter |
| To Markdown-paneler, vindu/panel-ID, lenkeruting og visningsbekreftelse | XFMD-arbeidet / XfmdDocumentHost |

G6-M1: statisk katalog/navigasjon; ingen døde lenker eller duplisert projeksjon.
G6-M2: utvalg ved behov, likt innhold som samme utvalg i full eksport, ugyldig
utvalg og renderfeil beholder siste visning. G6-M3: ekte XFMD-panelklikk,
flere vinduer, fokusendring og lukket mål. G6-M4: cache-hit/invalidering,
raske klikk, crash/restart, lease, ressursgrenser og opprydding.
Alle fire milepæler er planlagt.

Go-porten av den strukturelle SDL-frontenden (G4-M1) er forutsetning for G6-M1;
dagens Python-verktøy er portgrunnlag. G6 trenger ikke SDL-domenekjøring,
Go-kodegenerering eller SDUI-layout. Viewpoint-porten flyttes fra G5-M3 til
G6-M1; G5-M3 bruker så denne dokumenteksporten. Nummeret G6 er en arbeidsstrøm,
ikke et krav om at G1–G5 må ferdigstilles først.
