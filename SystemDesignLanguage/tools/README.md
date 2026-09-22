# SDL-verktøy — genererte viewpoints

Dette er SDLs avgrensede verktøyinngang. Det er ikke installasjonsverktøyet i
repoets `Toolkit/` (Standard Document Procedure). Implementert i Python med
standardbibliotek og **gjenbruk av eksisterende design-core-parser**.
Ingen ny parser, Go-runtime eller håndskrevne SDUI-diagrammer.

Dette er et kjørbart verktøyinkrement på dagens frontend. Ved Go-porteringen
portes samme viewpoint-regler og tester; Python-verktøyet erstattes samlet.
Det innfører ingen egen runtime eller permanent parser ved siden av Go.

Fra SDP-roten, generer alle viewpoints som dagens profil kan beskrive:

```sh
python3 SystemDesignLanguage/tools/sdl.py viewpoints SDUI/design/architecture.design --output SDUI/design/viewpoints
```

Velg ett eller flere viewpoints:

```sh
python3 SystemDesignLanguage/tools/sdl.py viewpoints SDUI/design/architecture.design --output /tmp/sdl-architecture --viewpoint VP02 --viewpoint VP03
```

Generer også faktiske SVG-bilder og samlet Markdown-utskrift med lokal renderer:

```sh
python3 SystemDesignLanguage/tools/sdl.py viewpoints SDUI/design/architecture.design --output SDUI/design/viewpoints --renderer /home/warloc/git/mermaid-rs-renderer/target/debug/mmdr
```

`--renderer` er eksplisitt sti til et tilgjengelig mmdr-program, ikke et krav om
at SDL-parseren linker mot Mermaid. En annen konsument kan rendre de genererte
Mermaid-blokkene selv. Rendererens eksakte binærhash og versjon registreres.
SVG-generering beviser ikke interaktivitet eller fysisk utskrift/PDF-paginering.

| Artefakt | Innhold |
| --- | --- |
| viewpoints.md | Valgte viewpoints med Mermaid-fences, tabeller og mangelstatus |
| diagrams/*.mmd | De samme diagramkildene som separate filer |
| diagrams/*.svg | Resultat fra valgt renderer, når `--renderer` er oppgitt |
| printout.md | Valgte viewpoints med SVG-bilder, når renderer er brukt |
| implementation.md | Faser, milepæler, ansvar, eiere og scenarioer fra VP06; ingen innebygd G-plan |
| message-sets.json | Avledede meldingssett per Channel og modus fra VP08 |
| manifest.json | Kildehash, verktøyhash, rendereridentitet, fakta, node-/kantkoblinger og kildeposisjoner |

Se [viewpoint-katalog og språkstatus](../../docs/checkpoint%231/08-SDL-Viewpoints-and-Implementation-Status.md).
VP01/02/03/05/06/07/08/09/10 gir diagrammer, VP04 tabell, VP11 hele faktaregisteret.
VP08 følger deklarerte, validerte scenario-steg. VP10 krever en eksplisitt Encoding. VP01/07 viser også hull i modellen. Capability omdøpes ikke til Feature;
Interface-bruk omgjøres ikke til Channel-sekvens; bitbredder gjettes ikke.
Alle f-ID-er er lokale til rapportens kildehash; modellnavn og kildeposisjoner
følger med slik at endringer ikke forveksles med vedvarende faktumidentitet.

Testkommando:

```sh
python3 -m unittest discover -s SystemDesignLanguage/tools -p 'test_*.py' -v
```

Testene kontrollerer faktisk modellprojeksjon, kildekobling, endret input,
viewpoint-utvalg, urealiserte ansvar, ukjent profil og repeterbar CLI-eksport.

Eksporten bygges i en midlertidig katalog før publisering. Parser-/rendererfeil
beholder siste rapport. Et nytt utvalg fjerner utgåtte genererte diagrammer;
andre notatfiler bevares. Ikke legg SDL-kilden i utdataområdet. Uten `--renderer`
fjernes tidligere generert printout/SVG så gammel rendering ikke vises som ny.
Rendererresultater kan gjenbrukes når både Mermaid-kilde, renderer og SVG-hash stemmer.

V4-snapshot, verifisert 2026-09-22: 24/24 verktøytester, 61/61 SDL-parsertester og 36/36
SDUI-tester. Prøvemodellen gir 142 SVG-diagrammer; alle kartlagte nodenavn finnes
i SVG-teksten. Parser-, runtime-, arkitektur- og modusvisninger er visuelt
stikkprøvekontrollert. Ingen fysisk utskrift eller paginert PDF er testet.
[Maskinrapport](verification.json) og [generert printout](../../SDUI/design/viewpoints/printout.md).

Leveransehistorikk: V1 innførte **design-core 0.2**, som erstatter aktiv 0.1. `goal_views.py` projiserer
UseCase-sporbarhet og Feature-utsnitt; den inneholder ingen navn/fakta fra SDUI-
eksemplet. VP01 er en merket flowchart, ikke en påstand om formell UML-støtte.
VP07 viser én tegning per Feature og eksplisitt allokeringsmodus. Alle bidrag
vises, også de uten plassering i valgt modus. Ingen contains-kant eller Mode
arves som allokering. Modus følger tegningstittel og kildefaktum i manifestet.
Modellhull finnes både i Markdown og `model_gaps` i manifestet; de avviser ikke
en ellers gyldig delmodell. Et utsnitt beviser ikke fullstendig Feature-kjørbarhet.

```sh
python3 SystemDesignLanguage/tools/sdl.py viewpoints SDUI/design/architecture.design --output /tmp/sdl-goals --viewpoint VP01 --viewpoint VP07
```

V2 legger til data-/kontraktkart og packet fra design-core 0.3. Den generelle
verifikasjonskommandoen kjører testene, regenererer via CLI og kontrollerer
kildekobling, SVG-etiketter, eksakte bitområder og repeterbar eksport:

```sh
python3 SystemDesignLanguage/tools/verify_design.py --phase V4 --renderer /home/warloc/git/mermaid-rs-renderer/target/debug/mmdr
```

V3 innførte design-core 0.4. message-sets.json er et generert katalogprodukt;
scenario-piler har kilde-ID-er for steg, deltakelse og governing kontrakt.

Aktiv profil er **design-core 0.5** (V4); eldre profiler er erstattet.
`plan_views.py` lager både VP06 og implementasjonsrapporten fra validerte fakta.
[Start med G1–G5-rapporten](../../SDUI/design/viewpoints/implementation.md).
V4s fasegraf, parsersekvens og bindingssekvens er visuelt stikkprøvekontrollert.

G6 er senere lagt til som **planlagt design** i den felles SDL-kilden.
[Checkpoint 10](../../docs/checkpoint%231/10-SDL-Viewpoint-Navigation.md)
viser nye modelltall. Katalogbaserte sider, navigasjonslenkehandler og daemon
er ennå ikke implementert av dette verktøyet; eksportformatet er uendret.

G6-D2s [nivå- og notasjonsprofil](../../docs/SDL-Viewpoint-Levels-and-Notation.md)
planlegger navigator-only, A0–A5 og typede symboler. Dagens flowchart-utdata er
fortsatt gjeldende verktøyimplementasjon; nye use-case-/klasseregler er ikke aktivert.
