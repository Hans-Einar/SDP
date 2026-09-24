# R2 — konkrete dokumentkonflikter

Inventargrunnlag: [R1s dokumentkart](../R1/Documentation-index.md). Tabellen
skiller observerte feil fra uløste språk-/prosessvalg. Bevis er eksisterende
implementasjon/bevis fra G-fasene; R2 er dokumentvedlikehold.

| ID | Dokument / observert konflikt | Grunnlag for behandling | R2-M2 |
| --- | --- | --- | --- |
| D01 | SDLs KanBan sier koden fortsatt ligger i SystemDesignLanguage | R1 flyttematrise og SDL/go | Rettet til SDL |
| D02 | SDL språkdefinisjon peker på isolert parserprototype; docs-oversikten kaller hele studies/ utforskning | Definisjonens §§1–11 + registrerte profiler (§12 er åpent arbeid) er design-core 0.5; §§13–16 er kandidater; Go-parser er aktiv | Blandet autoritet presisert; utgått parserinngang fjernet |
| D03 | Data-/Channel-profiler bruker 0.3/0.4 som nåstatus og gamle Python-testtall uten klar seksjonsgrense | design-core 0.5 i Go inkluderer V2–V4; G4/G6/G5-bevis | Aktive regler skilt fra daterte V2–V4-bevis |
| D04 | Leveranseplanprofil sier all G-status fortsatt er planned | Arkitekturmodellen og G-fasebevis; implementation-status er fortsatt bare en eksplisitt kildepåstand | V4-snapshot merket; global nåstatuspåstand rettet |
| D05 | Navigasjons-/notasjonsdesign sier G6 ikke er implementert | SDL/go/documents, query, broker, reader, G6-bevis og G7-launcher | Merket som opprinnelig design med lenker til implementasjon og begrensninger |
| D06 | SDUI runtime-kontrakt sier binding ennå er planlagt og har uavklarte regler som nå finnes i Go | SDUI/go/runtime/README, SDL action-core/bridge, G3/G4 | Levert én kort gjeldende grenseoversikt som peker til pakkekontraktene |
| D07 | SDUI layout-/komposisjonsforslag sier Python er aktiv og layout/runtime gjenstår | SDUI language, go-layout-contract, architecture og G1–G3 | Designbakgrunn beholdt; aktiv profil og avvik merket eksplisitt |
| D08 | prototype-widgets har kjørbare instruksjoner for slettet Python-verktøy | G5-M4, Go SVG-CLI og sdl-document | Utgåtte kjørekommandoer fjernet; datert artefaktbeskrivelse beholdt |
| D09 | SDUI mangler samlet docs-inngang; checkpoint omtales som gjeldende autoritet | README, aktive Go-kontrakter, implementasjonsplan; checkpoint er datert historie | Dokumentkart og tydelige aktive innganger levert |
| D10 | Portinventarets nåværende Go-sti er gammel; IR-studie kan leses som ingen runtime finnes | R1, action-core 0.1 og G4 | Aktiv sti rettet; bred IR-studie avgrenset fra levert handlingsruntime |

Videre arbeid: bred kandidatsemantikk, SDP-prosess-/malprofil og full faglig
harmonisering. R2 hevder ikke at alle formuleringer i 7000+ linjer studier er
normativt avklart. Genererte viewpoints endres bare via verktøyet ved endret kilde.

## Kontrollert resultat

D01 gjelder også KB-SDP-001s foreldede førtilstand. D07 omfatter retting av
språkprofilens påstand om at fontenheten ikke var fastsatt: Go bruker DIP.
D10 omfatter source-tree-studiens gamle 0.1-/Python-inngang; bred workspace/IR
forblir kandidat, mens avgrensede data-/Channel-/action-profiler er implementert.
SDUI har et [komplett dokumentkart](../../../SDUI/docs/README.md) over alle 17
faglige docs-filer. SDLs [inngang](../../../SDL/docs/README.md) beskriver eksplisitt
at språkdefinisjonen har både aktiv kjerne og forslag. Runtimekontrakten peker
til én kilde per pakkeansvar; opprinnelig runtimeforslag bevares i Git.

Historiske tall under V2–V4 og opprinnelige G6-designskisser er ikke oppdatert
som om de var nye prøver. R2 markerer alder/autoritet og peker til aktive innganger.
