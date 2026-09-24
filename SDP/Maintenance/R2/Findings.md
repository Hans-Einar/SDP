# R2 — konkrete dokumentkonflikter

Inventargrunnlag: [R1s dokumentkart](../R1/Documentation-index.md). Tabellen
skiller observerte feil fra uløste språk-/prosessvalg. Bevis er eksisterende
implementasjon/bevis fra G-fasene; R2 er dokumentvedlikehold.

| ID | Dokument / observert konflikt | Grunnlag for behandling | R2-M2 |
| --- | --- | --- | --- |
| D01 | SDLs KanBan sier koden fortsatt ligger i SystemDesignLanguage | R1 flyttematrise og SDL/go | Rettes til SDL |
| D02 | SDL språkdefinisjon peker på isolert parserprototype; docs-oversikten kaller hele studies/ utforskning | Definisjonens §§1–12 + registrerte profiler er design-core 0.5; §§13–16 er kandidater; Go-parser er aktiv | Presiser blandet autoritet og fjern utgått parserinngang |
| D03 | Data-/Channel-profiler bruker 0.3/0.4 som nåstatus og gamle Python-testtall uten klar seksjonsgrense | design-core 0.5 i Go inkluderer V2–V4; G4/G6/G5-bevis | Skill aktive regler fra daterte V2–V4-bevis |
| D04 | Leveranseplanprofil sier all G-status fortsatt er planned | Arkitekturmodellen og G-fasebevis; implementation-status er fortsatt bare en eksplisitt kildepåstand | Merk V4-snapshot, fjern global nåstatuspåstand |
| D05 | Navigasjons-/notasjonsdesign sier G6 ikke er implementert | SDL/go/documents, query, broker, reader, G6-bevis og G7-launcher | Leses som opprinnelig design med lenker til implementasjon og begrensninger |
| D06 | SDUI runtime-kontrakt sier binding ennå er planlagt og har uavklarte regler som nå finnes i Go | SDUI/go/runtime/README, SDL action-core/bridge, G3/G4 | Én kort gjeldende grenseoversikt som peker til pakkekontraktene |
| D07 | SDUI layout-/komposisjonsforslag sier Python er aktiv og layout/runtime gjenstår | SDUI language, go-layout-contract, architecture og G1–G3 | Behold designbakgrunn, merk vedtatt profil og avvik eksplisitt |
| D08 | prototype-widgets har kjørbare instruksjoner for slettet Python-verktøy | G5-M4, Go SVG-CLI og sdl-document | Fjern utgåtte kjørekommandoer; behold datert artefaktbeskrivelse |
| D09 | SDUI mangler samlet docs-inngang; checkpoint omtales som gjeldende autoritet | README, aktive Go-kontrakter, implementasjonsplan; checkpoint er datert historie | Dokumentkart og tydelige aktive innganger |
| D10 | Portinventarets nåværende Go-sti er gammel; IR-studie kan leses som ingen runtime finnes | R1, action-core 0.1 og G4 | Rett aktiv sti og avgrens bred IR-studie fra levert handlingsruntime |

Videre arbeid: bred kandidatsemantikk, SDP-prosess-/malprofil og full faglig
harmonisering. R2 hevder ikke at alle formuleringer i 7000+ linjer studier er
normativt avklart. Genererte viewpoints endres bare via verktøyet ved endret kilde.
