# Checkpoint #1 — SDL-status og genererte viewpoints

**Implementasjonsstatus er oppdatert i [tillegg 11](11-Go-Implementation-and-Navigation.md).**
Nedenfor beholdes det daterte design-/V-fasegrunnlaget; gamle Python-kommandoer
er historiske og erstattet av Go-inngangene.

Dato: 2026-09-22. Gjennomgang av dokumentasjon, eksisterende parser, MVP1-korpus
og SDL-modellen av SDL/SDUI. Tall og verifikasjon nedenfor er V4-snapshotet.
Se [tillegg 11](10-SDL-Viewpoint-Navigation.md) for senere G6-design og dagens modell.
SDUI 0.2 og Go/Fyne-retningen i [tillegg 07](07-SDUI-0.2-and-Go-Direction.md) beholdes.

## Konklusjon

Checkpoint #1 er langt mer omfattende enn den implementerte SDL-profilen.
Den inneholder begreper, eksempler og kandidater, ikke en ferdig kompilator for
hele beskrivelsen. V0–V4 er nå levert som en avgrenset **SDL-viewpoint-leveranse**.
[Tillegg 09](09-SDL-Generated-Go-Design-Review.md) er inngangen til gjennomgang av
det genererte G1–G5-designet før Go-implementasjon.

Viewpoints skal genereres av SDL-verktøyet fra validerte modellfakta. En agent skal
ikke tegne et plausibelt diagram og presentere det som verktøyets resultat.
Visningsvalg får filtrere og ordne fakta, men ikke opprette nye arkitekturfakta.

## Implementert mot dokumentert

| Område | Checkpoint/studier | Faktisk implementasjon |
| --- | --- | --- |
| Unit, Container, Functionality, Capability, Interface, Activity, Mode | Strukturelt grunnlag | Python design-core 0.5 parser/AST/validator/formatter |
| contains, owns, realizes, provides, consumes, requires in mode, refines | Typede strukturrelasjoner | Implementert med navn, eierskap, sykluser og kildeposisjoner |
| Actor, Use Case, Feature og bidrag | 02/04 og V1s avgrensede språkregler | Implementert i 0.2: pursues, supports, contributes-to; direkte og mange-til-mange-bidrag |
| Functionality-allokering | V1, eksplisitt kjørekontekst | allocated-to Container in mode Mode; høyst én Container per ansvar/modus, uendret logisk eier |
| System og flerfilskilder | Studier og MVP1-korpus | Ingen komplett resolver/import-/deployment-/instansmodell i design-core |
| Channel, deltakere, Contract og MessageSet | [V3-profil](../../../SDL/docs/profiles/SDL-Channel-Scenario-Profile.md) | Typede roller, permits og modus; MessageSet avledes |
| Dataset, Datagram, Database og feltkontrakter | [V2-profil](../../../SDL/docs/profiles/SDL-Data-Contract-Profile.md) | Implementert: kilde/holder, closed/open, varianter, typed presence og fast wirelayout |
| Scenario-steg og request/resultat | V3s eksplisitte stegprofil | Ordning og korrelasjon valideres; ingen runtime eller tilstandsmaskin |
| Viewpoint-katalog | Studie nevner containerkart, avhengigheter, realisering, scenario og endringspåvirkning | Ingen tidligere ferdig katalog; konkretisert nedenfor |
| Markdown/Mermaid-generering | Foreslått modell/diagram/kildekart-pakke | Avgrenset SDL-verktøykommando opprettet i denne leveransen |
| Go-parser/runtime, binding til Go | Valgt videre retning | Kataloger, ikke implementert kode |

Kildene er [språkdefinisjonen](../../../SDL/docs/studies/Design-Language-Definition.md),
[parseren](../../../SDL/go/README.md),
[blueprint-studien](../../../SDL/docs/studies/SDL-Source-Tree-and-Compilation-Study.md),
[MVP1-korpusets grenser](../../../experiments/mvp1_sdl/README.md) og
[datastudien](../../../SDL/docs/studies/SDL-Datasets-Datagrams-and-Data-Contracts.md).
MVP1s inventaraudit kontrollerer korpusdekning; den er ikke en full SDL-parser.

## Viewpoint-katalog

Dette er en avgrenset anbefalt katalog for eierens ønskede utskrift. Den er ikke
påstand om at alle konstruksjonene er vedtatt grammatikk. Identitetene brukes av
SDL-generatoren; nye viewpoints må angi hvilke modellfakta de trenger.

| ID | Spørsmål og diagram | Nødvendig grunnlag | Status i dagens verktøy |
| --- | --- | --- | --- |
| VP01 | Hvem vil oppnå hva? Use Case og kobling til Feature/Functionality | Actor, UseCase, pursues/supports/contributes-to; modellen er omfanget, ingen formell System-grense | Genereres |
| VP02 | Hvordan er arkitekturen delt? Container-/Unit-kart | Deklarasjoner og contains; allokering vises separat i VP07 | Genereres |
| VP03 | Hvilke ansvar realiserer en kapabilitet, og hvem eier dem? | Unit owns Functionality, Functionality realizes Capability, Unit provides Capability | Genereres |
| VP04 | Hvilke grensesnitt brukes? Samarbeidstabell, senere kontraktgraf | consumes; senere tilbyder/binding/kontrakt | Tabell genereres; kobling mellom konsument og tilbyder utledes ikke |
| VP05 | Hva kreves i valgt modus? Avhengighetsgraf | requires Interface in mode Mode | Genereres |
| VP06 | Hvordan deles arbeidet i faser og milepæler? | refines, addresses, delivers, depends-on, illustrates og implementation-status | Genereres med ansvarsdekning og eksplisitte avhengigheter; ingen scheduler |
| VP07 | Hvordan ligger en Feature over arkitekturen? Markert bidragslag | Feature-bidrag, Functionality-eier og eksplisitt Container-allokering per modus | Genereres; manglende allokering rapporteres, ingen modus arves |
| VP08 | Hvordan kommuniserer deltakerne over Channel? Sekvens | Deltakere/roller, kontrakter, meldinger, scenario/protokollsteg, korrelasjon, separate navngitte alternative baner | Genereres |
| VP09 | Hvor kommer data fra og hvor hentes de? Datakart/ER | Dataset, Datagram-kilde/projeksjon, kontrakt, valgfri Database-lagring, eier | Genereres |
| VP10 | Hvordan er én Datagram-variant kodet? Packet | Eksplisitt kodingsprofil med feltrekkefølge, bredde/offset og variant | Genereres for eksplisitt Encoding |
| VP11 | Hva støtter visningen, og hva mangler? Faktaregister/egenskaper | Alle deklarasjoner/fakta med kildeposisjoner og profilgrenser | Genereres |

Et Feature-lag skal være et filter/dekorasjon på den samme arkitekturen, uten å
flytte eierskap eller opprette nye containere. Mange-til-mange-forhold skal bevares.
Use Cases kan ha direkte Functionality-bidrag uten en oppdiktet Feature mellom.
Diagramnotasjonen velges etter verifisert rendererprofil; flowchart kan brukes
som merket mål-/sporbarhetsvisning når vertens Use Case-notasjon ikke er støttet.

En Channel-kontrakt angir tillatt samarbeid, men bestemmer ikke nødvendigvis én
sekvens. VP08 må velge en deklarert kontraktprotokoll eller et scenario som bruker
kontrakten. Generatoren skal aldri sortere relasjonslinjer og kalle det en tidsrekkefølge.
[Mermaid sequence](https://mermaid.js.org/syntax/sequenceDiagram.html) tegner den
rekkefølgen generatoren gir den; semantikken må komme fra SDL.

[Mermaid packet](https://mermaid.js.org/syntax/packet.html) beskriver bitposisjoner
og felt. Et logisk Datagram alene har ikke nødvendigvis et slikt binært format.
Dataset/Database-koblingen vises i VP09 og lenkes fra VP10. JSON, filer eller en
ikke fastlagt encoding skal ikke få oppdiktede bitbredder. Checkpoint 03s 128-bit-
eksempel var eksplisitt oppfunnet for å prøve diagramtypen, ikke en vedtatt wirekontrakt.

## Database — eierpresisering

Database er **et samlebegrep for hvor persistent data kan hentes ved behov**.
Det er ikke synonymt med SQL, tabeller, en databaseserver eller en separat prosess.
Filer, objektlager eller en ekstern persistent datatjeneste kan være realiseringer.
Dataset beskriver den logiske datamengden; Datagram beskriver avtalt overføring
eller projeksjon av den. Ikke alle Datasets trenger persistent lagring.

Dette presiserer CP1-D11: persistent tilgang er del av Database-betydningen;
konkret varighet, tilgjengelighet, eierskap og lese-/skriveoperasjoner må fortsatt
stå i kontrakten. En bare midlertidig cache blir ikke automatisk Database.
Eldre prosa om «retained/queryable» skal leses med denne presiseringen. Ingen
ny parserstøtte hevdes ved å oppdatere definisjonen i checkpointet.

## Kjørbar leveranse nå

[SDL-verktøyet](../../../SDL/tools/README.md) har kommandoen
`viewpoints`. Den bruker eksisterende SDL-parser, validerer kilden og lager
Markdown med Mermaid, enkeltstående .mmd-filer og et maskinlesbart kildekart.
Med valgt renderer lager den også SVG og en samlet printout.md med bildeinnbygging.
Dette er en SDL-verktøyfunksjon, ikke håndtegnede diagrammer for SDUI-eksemplet.

Prøveinput er [SDL-modellen av SDL/SDUI](../../../SDUI/design/architecture.design).
[Generert Markdown](../../../SDUI/design/viewpoints/viewpoints.md) og
[rendret utskrift](../../../SDUI/design/viewpoints/viewpoints.md) viser alle tilgjengelige
viewpoints og eksplisitte modellhull. 368 deklarasjoner og 1106 fakta kommer
fra kilden. Actor/UseCase/Feature og utvalgt Functionality-allokering er skrevet
som SDL-fakta i modellen, ikke rekonstruert av generatoren fra prosa eller navn.
De 6 bruksmålene, 6 Features og 2 Actors dekker inspeksjon, UI-prøving, domene-
binding, modellreload, dokumentasjon og native bygg. `SdlViewpointGenerator`
er også beskrevet som Unit med ansvar for projeksjon, Markdown og kildekobling.

VP07 viser 15 uspesifiserte allokeringer i modusutsnittene. Et delt ansvar kan
ha plassering i en modus der andre bidrag til samme Feature ennå ikke er
allokert. Det er en synlig grense for modellen, ikke bevis på at hele Feature-en
kan kjøres i den modusen. Funksjonelt bidrag er ennå ikke modusbetinget.
Disse utsnittene er ikke en komplett deploymentplan for alle 94 ansvar.

Verifikasjon: 142 SVG-diagrammer rendret, 24 verktøytester, 61 SDL-parsertester og
36 SDUI-tester består. Alle diagramnoder har modell-/kildekobling, alle nodenavn
finnes i SVG-en, og gjentatt eksport gir samme artefakter. Visuell kontroll er
stikkprøver, ikke fysisk print/PDF. [Rapport](../../../SDL/tools/verification.json).

## SDL-milepæler og neste leveranse

| Milepæl | Akseptanse |
| --- | --- |
| V0 — strukturell generator, levert | Dagens profil → valgbare viewpoints, deterministisk Markdown/Mermaid, kildekart, renderbevis og eksplisitte hull |
| V1 — mål og bidrag, levert | design-core 0.2, Actor/UseCase/Feature, direkte/indirekte bidrag og modusallokering; positive/negative tester; VP01/VP07 generert fra portert SDL/SDUI-modell |
| V2 — data og kontrakter, levert | Definer Dataset/Datagram/Database og kontraktinnhold, opprinnelse/projeksjon og ved-behov-persistens; VP09, og VP10 kun for valgt encoding |
| V3 — Channel og scenario, levert | Typede deltakere/roller, kontraktreferanser og deklarert meldingsrekkefølge; VP08 med request/resultat og negativ kontraktkontroll |
| V4 — samlet SDUI-design, levert | Berik SDL-kilden med de nye fakta og få alle bestilte viewpoints generert; ingen faktatillegg i generatoren |

V0–V4 er levert. Eierens avgrensning er nå å gjennomgå den
[genererte G1–G5-planen](../../../SDUI/design/viewpoints/implementation.md) før
Go-fasene startes. Fem faser og 18 milepæler dekker alle 94 Functionality-er.
Alle G-faser og milepæler har status planned i SDL-kilden.
Ikke bygg en stor ny Python-runtime for å få diagrammer. Go-parserporten og disse
viewpoint-behovene må planlegges sammen; MVP1-korpusets kandidatsyntaks er ikke en
ferdig grammatikk som kan aktiveres med en permissiv fallback.

SDUI har et godt stoppunkt: 0.2-profil, kjørbar Python-parser, AST/dumper,
SVG-widgetprøve og planlagt Go/Fyne-runtime. Dagens UI-eksempler trenger ikke
omskrives for denne SDL-runden. Neste UI-runtimefase bør vente på målrettet
SDL-avklaring; eksisterende bevis og portgrunnlag beholdes.
