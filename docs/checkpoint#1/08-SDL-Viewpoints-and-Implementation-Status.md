# Checkpoint #1 — SDL-status og genererte viewpoints

Dato: 2026-09-22. Gjennomgang av dokumentasjon, eksisterende parser, MVP1-korpus
og SDL-modellen av SDL/SDUI. Dette er dagens inngang for prioritering av SDL-arbeid.
SDUI 0.2 og Go/Fyne-retningen i [tillegg 07](07-SDUI-0.2-and-Go-Direction.md) beholdes.

## Konklusjon

Checkpoint #1 er langt mer omfattende enn den implementerte SDL-profilen.
Den inneholder begreper, eksempler og kandidater, ikke en ferdig kompilator for
hele beskrivelsen. Vi bør fullføre en avgrenset **SDL-viewpoint-leveranse før mer
SDUI-runtimearbeid**. Det krever ikke at all kjørbar SDL-semantikk implementeres først.

Viewpoints skal genereres av SDL-verktøyet fra validerte modellfakta. En agent skal
ikke tegne et plausibelt diagram og presentere det som verktøyets resultat.
Visningsvalg får filtrere og ordne fakta, men ikke opprette nye arkitekturfakta.

## Implementert mot dokumentert

| Område | Checkpoint/studier | Faktisk implementasjon |
| --- | --- | --- |
| Unit, Container, Functionality, Capability, Interface, Activity, Mode | Strukturelt grunnlag | Python design-core 0.3 parser/AST/validator/formatter |
| contains, owns, realizes, provides, consumes, requires in mode, refines | Typede strukturrelasjoner | Implementert med navn, eierskap, sykluser og kildeposisjoner |
| Actor, Use Case, Feature og bidrag | 02/04 og V1s avgrensede språkregler | Implementert i 0.2: pursues, supports, contributes-to; direkte og mange-til-mange-bidrag |
| Functionality-allokering | V1, eksplisitt kjørekontekst | allocated-to Container in mode Mode; høyst én Container per ansvar/modus, uendret logisk eier |
| System og flerfilskilder | Studier og MVP1-korpus | Ingen komplett resolver/import-/deployment-/instansmodell i design-core |
| Channel, deltakere, Contract og MessageSet | Kandidater; MessageSet skal avledes | Ikke implementert som generell validerbar profil |
| Dataset, Datagram, Database og feltkontrakter | [V2-profil](../SDL-Data-Contract-Profile.md) | Implementert: kilde/holder, closed/open, varianter, typed presence og fast wirelayout |
| Scenario-steg, protokollrekkefølge og tilstand | Kandidater og MVP1-eksempler | Ingen generell parser/runtime for disse |
| Viewpoint-katalog | Studie nevner containerkart, avhengigheter, realisering, scenario og endringspåvirkning | Ingen tidligere ferdig katalog; konkretisert nedenfor |
| Markdown/Mermaid-generering | Foreslått modell/diagram/kildekart-pakke | Avgrenset SDL-verktøykommando opprettet i denne leveransen |
| Go-parser/runtime, binding til Go | Valgt videre retning | Kataloger, ikke implementert kode |

Kildene er [språkdefinisjonen](../Design-Language-Definition.md),
[parseren](../../experiments/design_core/design_core.py),
[blueprint-studien](../SDL-Source-Tree-and-Compilation-Study.md),
[MVP1-korpusets grenser](../../experiments/mvp1_sdl/README.md) og
[datastudien](../SDL-Datasets-Datagrams-and-Data-Contracts.md).
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
| VP06 | Hvordan detaljeres en aktivitet? Refinement-graf | Activity refines Activity | Genereres; ingen tidsrekkefølge |
| VP07 | Hvordan ligger en Feature over arkitekturen? Markert bidragslag | Feature-bidrag, Functionality-eier og eksplisitt Container-allokering per modus | Genereres; manglende allokering rapporteres, ingen modus arves |
| VP08 | Hvordan kommuniserer deltakerne over Channel? Sekvens | Deltakere/roller, kontrakter, meldinger, scenario/protokollsteg, korrelasjon, alternativer | Mangelrapport |
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

[SDL-verktøyet](../../SystemDesignLanguage/tools/README.md) har kommandoen
`viewpoints`. Den bruker eksisterende SDL-parser, validerer kilden og lager
Markdown med Mermaid, enkeltstående .mmd-filer og et maskinlesbart kildekart.
Med valgt renderer lager den også SVG og en samlet printout.md med bildeinnbygging.
Dette er en SDL-verktøyfunksjon, ikke håndtegnede diagrammer for SDUI-eksemplet.

Prøveinput er [SDL-modellen av SDL/SDUI](../../SDUI/design/architecture.design).
[Generert Markdown](../../SDUI/design/viewpoints/viewpoints.md) og
[rendret utskrift](../../SDUI/design/viewpoints/printout.md) viser alle tilgjengelige
viewpoints og mangelstatus for resten. 227 deklarasjoner og 509 fakta kommer
fra kilden. Actor/UseCase/Feature og utvalgt Functionality-allokering er skrevet
som SDL-fakta i modellen, ikke rekonstruert av generatoren fra prosa eller navn.
De 6 bruksmålene, 6 Features og 2 Actors dekker inspeksjon, UI-prøving, domene-
binding, modellreload, dokumentasjon og native bygg. `SdlViewpointGenerator`
er også beskrevet som Unit med ansvar for projeksjon, Markdown og kildekobling.

VP07 viser 15 uspesifiserte allokeringer i modusutsnittene. Et delt ansvar kan
ha plassering i en modus der andre bidrag til samme Feature ennå ikke er
allokert. Det er en synlig grense for modellen, ikke bevis på at hele Feature-en
kan kjøres i den modusen. Funksjonelt bidrag er ennå ikke modusbetinget.
Disse utsnittene er ikke en komplett deploymentplan for alle 92 ansvar.

Verifikasjon: 73 SVG-diagrammer rendret, 17 verktøytester, 49 SDL-parsertester og
36 SDUI-tester består. Alle diagramnoder har modell-/kildekobling, alle nodenavn
finnes i SVG-en, og gjentatt eksport gir samme artefakter. Visuell kontroll er
stikkprøver, ikke fysisk print/PDF. [Rapport](../../SystemDesignLanguage/tools/verification.json).

## SDL-milepæler og neste leveranse

| Milepæl | Akseptanse |
| --- | --- |
| V0 — strukturell generator, levert | Dagens profil → valgbare viewpoints, deterministisk Markdown/Mermaid, kildekart, renderbevis og eksplisitte hull |
| V1 — mål og bidrag, levert | design-core 0.2, Actor/UseCase/Feature, direkte/indirekte bidrag og modusallokering; positive/negative tester; VP01/VP07 generert fra portert SDL/SDUI-modell |
| V2 — data og kontrakter, levert | Definer Dataset/Datagram/Database og kontraktinnhold, opprinnelse/projeksjon og ved-behov-persistens; VP09, og VP10 kun for valgt encoding |
| V3 — Channel og scenario | Typede deltakere/roller, kontraktreferanser og deklarert meldingsrekkefølge; VP08 med request/resultat og negativ kontraktkontroll |
| V4 — samlet SDUI-design | Berik SDL-kilden med de nye fakta og få alle bestilte viewpoints generert; ingen faktatillegg i generatoren |

V1 er levert som språk-/modellarbeid før en stor runtime. Neste avgrensede
leveranse er V3. Avklar grammatikk og semantikk på én SDUI-witness, oppdater definisjon/tester samlet, og port til Go i ett løp.
Ikke bygg en stor ny Python-runtime for å få diagrammer. Go-parserporten og disse
viewpoint-behovene må planlegges sammen; MVP1-korpusets kandidatsyntaks er ikke en
ferdig grammatikk som kan aktiveres med en permissiv fallback.

SDUI har et godt stoppunkt: 0.2-profil, kjørbar Python-parser, AST/dumper,
SVG-widgetprøve og planlagt Go/Fyne-runtime. Dagens UI-eksempler trenger ikke
omskrives for denne SDL-runden. Neste UI-runtimefase bør vente på målrettet
SDL-avklaring; eksisterende bevis og portgrunnlag beholdes.
