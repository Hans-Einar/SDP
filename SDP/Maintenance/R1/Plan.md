# R1 — organiser SDP-repoet

Eier aktiverte KB-SDP-001 den 2026-09-24 og ba om housekeeping nå. Branch:
`sdp/phase-r1-repository-organization`, fra K3 `431e47e`.

## Mål og avgrensning

Ett fysisk hjem for maler, prosjektets egne records og hvert språk. Behold
semantikken i gjeldende installasjon, språk og historiske bevis. Redaksjonell
sammenskriving/konfliktavklaring ligger i KB-SDP-010. Ikke opprett GitHub-repoer,
submodules eller ny sdptool-implementasjon som del av katalogflyttingen.

## Målstruktur

| Område | Ansvar |
| --- | --- |
| Template/ | Eneste kilde for nøytrale project-root/sdp-root-maler |
| Toolkit/ | Installer, schema, verifikasjon, administrert payload og skills; egne kontraktdokumenter i docs/ |
| SDP/ | Utviklingen av SDP selv: KanBan, Traceability, Sprints, review/verifikasjon, studies og vedlikehold |
| SDL/ | Eksisterende SDL-kode/grammatikk, dokumenter, scripts og eget SDP/KanBan |
| SDUI/ | Eksisterende SDUI-kode og dokumentasjon, eget SDP/KanBan |
| docs/ | Kort felles inngang og SDP-prosessdokumenter |
| SDP/History/ | Merket checkpoint og tidligere bootstrap/prosessutkast |
| experiments/ | Navngitte kandidatøvelser; ikke produksjonsprofiler |
| examples/ | Toolkit-kontrakteksempler; fortsatt manifeststyrte |

Checkpoint #1 er et felles datert diskusjons-/implementasjonssnapshot, ikke ren
SDL-spesifikasjon. Det samles i SDP/History/checkpoint-1 og lenkes fra SDL/SDUI.
Aktive språkprofiler legges hos SDL. Historiske source-index/hashmanifest bevares;
flyttemappen dokumenterer ny adresse uten å forfalske gammelt bevis.

Nummererte rotfiler er tomme maler, like de installerte malene bortsett fra siste
linjeskift. De dedupliseres mot Template. Eksisterende malers fase-/destinasjonsnavn
beholdes i denne migreringen: omnummerering og fase-studier må utformes samlet med
kravprofilen i KB-SDL-001. Ingen nye, parserugyldige .design-maler distribueres.
Dette er gjenværende arbeid i KB-SDP-001 etter fysisk organisering.

## Milepæler

| ID | Leveranse | Kontroll | Status |
| --- | --- | --- | --- |
| R1-M1 | Inventar, eierskap, flyttematrise; aktiver 001, registrer 010 | Alle flyttekilder finnes; malduplikater sammenlignet; baseline registrert | Levert |
| R1-M2 | Template og egne prosjektrecords; arkiver gammel bootstrap | Toolkit-regresjoner, konformans, baselinevalidator, sporbarhet og lenker | Levert |
| R1-M3 | SDL samlet; docs fordelt; oversikter og fungerende innganger | Go race-tester, CLI/launcher, genererte manifest uendret, lenker | Planlagt |

## Migreringsregler

[Flyttematrise](Migration-map.json) registrerer hver kilde, mål, milepæl og gammel
hash. [Dokumentinventar](Documents.json) klassifiserer alle Markdown-filer i den
opprinnelige docs-katalogen; klassifisering er eierskap/leseveiledning, ikke ny
språkautoritet. En menneskelesbar oversikt følger som Documentation-index.md.

Reparer nåværende Markdown-lenker og eksekverbare stier. Behold gamle ledgerlinjer,
daterte hashmanifest, portfixturer og genererte viewpoints; generert innhold må
fortsatt komme fra generatoren. Eksisterende Go-modulidentitet beholdes under
katalogflytting; endelig modul-/repoidentitet velges ved separat utskilling.
Installerens source-adresser flyttes, men destinations og oppførsel beholdes.
Konformansforventningene endres bare for de deklarerte source-adressene.

Kjent førtilstand: Toolkit-validatoren har ni ID-formatavvik fra eldre Issue #5-
records. De skal ikke skjules, omnummereres tilfeldig eller blandes med nye feil.
PowerShell er ikke installert lokalt; Windows-installerens native prøver kjøres
av CI. Lokale strukturelle tester kan ikke omtales som native PowerShell-bevis.

Fullført milepæl får egen commit. Fasens hode pushes etter leveranse, uten merge.
KB-SDP-001 forblir active så lenge fase-/malprofil og endelig repo-utskilling
fortsatt krever avklaring; fysisk housekeeping er ikke hele språk-/prosessdesignet.
