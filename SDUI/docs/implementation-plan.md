# SDL/SDUI — implementasjonsplan i Go

**Status 2026-09-22:** SDL-løpet V0–V4 er levert og pushet som fasebrancher.
Eieren har autorisert implementasjon av alle G-faser i én sammenhengende økt,
med fasebrancher, milepælcommits og push ved faseslutt. G1–G4 er levert;
øvrige milepæler står som planlagt inntil deres akseptanse er verifisert.
[Faktiske implementasjonsbevis](../go/evidence/G1.md) holdes atskilt fra
[den genererte designplanen](../design/viewpoints/implementation.md).
Designmodellens planned-status beskriver målstrukturen inntil samlet modelloppdatering.

**ID:** SDUI-PLAN-003 · **Revisjon:** 2026-09-22.
Erstatter PLAN-002s P0–P6-løp for Rust/C-ABI/FOX. Nye milepæler bruker G-prefiks;
henvisninger til P-faser i eldre bevis gjelder historien, ikke aktive leveranser.
[Målarkitektur](target-architecture.md) og
[checkpoint](../../docs/checkpoint%231/07-SDUI-0.2-and-Go-Direction.md).

Omfang: SDL/SDUI-parser og runtime i Go, felles SDUI-layout, SVG-eksport, første
Fyne-vert, modellreload og senere Go-generering. G0/G1–G4 er levert; G2 er under implementasjon.
Python SDUI 0.2, strukturparseren design-core og eksemplene er eksisterende portgrunnlag.

Designgrunnlag 2026-09-22: [felles SDL-strukturmodell](../design/README.md) beskriver
ansvar og avhengigheter for G1–G6 og passerer eksisterende parser. Dette er
designdekning, ikke fullførte implementasjonsmilepæler.

## G0 — oppdatert grunnlag og kataloger

| Milepæl | Akseptanse | Status |
| --- | --- | --- |
| G0-M1 | Checkpoint beskriver SDUI 0.2, faktisk implementasjon, Go-retning og historiske avvik | Levert |
| G0-M2 | Aktive arkitektur-/plan-/handoff-instrukser peker samme vei | Levert |
| G0-M3 | SDL og SDUI har egne Go-områder med parser/runtime og dokumentert ansvar | Levert; ingen Go-kode eller avhengigheter |

## G1 — kjørbar SDUI-frontend i Go

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G1-M1 | **Levert:** Go-modul, syntax-only CLI og parser/AST med kildeposisjoner; alle eksempel-AST-er samsvarer med Python-fixturene |
| G1-M2 | **Levert:** Validator/normalisering, relative regler, frame-regioner, instansbaner og kildegrenser; porterte positive/negative Python-tilfeller |
| G1-M3 | **Levert:** Concept1 AST/diagnoser og identiske konsoll-/Markdown-dumper; SVG/HTML-kontrollgalleri uten callbacks. Generell geometri følger G2 før Python-fixtureplassering fjernes |

Krav R01–R09, R13–R15, R19, R21–R24. Ingen SDL-kilde åpnes av parseren.
Ingen Go- eller Python-fallback for SDUI 0.1. Parallell Python brukes bare som
midlertidig portorakel; overgangens fullføring avhenger også av G2s eksport.

## G2 — én layoutmodell, SVG og første Fyne-vindu

Avhenger av G1-M2. Eier: SDUI-implementasjonen.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G2-M1 | **Levert:** Målekontrakt og fontenhet; relative akser, ratio, rader/grupper, header/body/footer, gap/padding og eksplisitte overflowgrenser |
| G2-M2 | **Levert:** Generell SVG-eksport fra målt modell; liten frame + knapp/input og deretter Concept1 uten håndplassert fixture; geometri- og visuell kontroll |
| G2-M3 | **Levert:** Fyne-vindu med samme geometri, knapper/input, Tab/fokus og registrert lokal Go-funksjon; ingen SDL-avhengighet |
| G2-M4 | **Levert:** Avgrenset Markdown-provider og separate Mermaid-ressurser; dokumentert støtteprofil og negative tilfeller, ikke påstått full diagramdekning |

Krav R10/R11/R14/R15/R17/R20/R21/R24. Begynn med liten vertikal prøve i M1–M3;
full Concept1/Markdown følger. Mål oppstart, resize og hukommelse før påstander
om «lettvekts». Fyne håndterer widgetinteraksjon; kjerne og SVG skal fungere uten GUI.
SVG-eksport gjengir valgt tilstand og utfører aldri callbacks.

## G3 — UI-runtime og SDUI hot reload

Avhenger av G1-M2/G2-M3. Eier: SDUI-implementasjonen.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G3-M1 | **Levert:** Typede hendelser/oppdateringer, stabile widgetinstanser og revisjoner; ubundet knapp gir eksplisitt status |
| G3-M2 | **Levert:** Filendring → parse/valider → publisering; ugyldig kilde beholder siste gyldige UI og viser kildediagnose |
| G3-M3 | **Levert:** Bevar kompatibel verdi/draft/fokus; typebytte/sletting håndteres, stale events og callback etter teardown avvises; ingen gjentatte domenehandlinger ved reload |

Krav R12/R16/R18/R25. UI-state og domene-state holdes adskilt. Avtal eventgrense
og UI-tråd før asynkronisering; ingen kompleks plugin-/prosessmekanisme i denne fasen.

## G4 — avgrenset SDL-parser/runtime og kobling til SDUI

SDL-parserarbeid kan starte ved siden av G1. Integrasjon avhenger av G3-M1.
Eier: SDL-implementasjonen for semantikk/runtime; SDUI eier UI-siden av porten.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G4-M1 | **Levert:** Port design-core-struktur til Go med dokumentert grammatikk og tester; checkpoint/MVP1-kandidater blir ikke automatisk støttet |
| G4-M2 | **Levert:** Definer én kjørbar profil for navngitt handling, typed input/resultat og binding til registrert Go-funksjon; negative og manglende bindinger avvises |
| G4-M3 | **Levert:** SDUI-knapp/input → SDL-handling → Go-funksjon → UI-oppdatering; kildekart og samme kontrakt med eksplisitt simulert domene |
| G4-M4 | **Levert:** Reload av SDL-modell med siste gyldige versjon, tilstandsregel og håndtering av pågående hendelser; Go-endring bygges/restartes |

R12/R18/R25. Bruk et avgrenset EditAptCell-scenario fra MVP1 som referanse etter
at enkel binding virker. Hele 66-fils korpuset er ikke et parserakseptansemål ennå.
Ingen maskin-/domenealgoritmer utledes fra struktur alene. Produktkoden i Ponsse
endres ikke automatisk som del av språkimplementasjonen.

## G5 — Go-generering og samlet dokumentasjon

Avhenger av avklart G4-profil og felles runtime. Eier: SDL/SDUI sammen.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G5-M1 | Generert Go oppretter samme modeller/bindinger; bygger sammen med separate håndskrevne domenefunksjoner |
| G5-M2 | Filbasert utviklingsmodus og generert program gir samme hendelsesspor/tilstand for avtalt profil; ufullstendig semantikk gir diagnose |
| G5-M3 | Reproduserbar SVG/Markdown-dokumentasjon fra valgt UI/state; kilde-/verktøyversjon og renderbevis |
| G5-M4 | Port fullført: gamle aktive Python-innganger/fixtureplassering fjernet eller erstattet, lenker/kommandoer oppdatert; én språkimplementasjon per profil |

Krav R19/R26. Go-byggcache gjenbrukes ved restart; dynamisk maskinkodeutskifting
eller separate workerprosesser er ikke nødvendig akseptanse for hot reload.

## G6 — navigerbare dokumenter og generering ved behov

**Status: under implementasjon**, med [designkontrakt og XFMD-handoff](../../docs/SDL-Navigable-Viewpoints-Design.md).
Eier: SDL for projeksjon/publisering; XFMD for dokumentpaneler og lenkeruting.
Dette er dokumentvisning, adskilt fra G2s Fyne-vert for interaktive SDUI-widgets.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G6-M1 | **Levert:** Alternative eksportformer: navigator/overview uten detaljdiagrammer, eller statisk pakke; A0–A5-kataloger, typeinventar og stabile lenker/ankre fra samme modell; port eksisterende SDL-projektor; samleeksport valgfri; lenker/bilder og determinisme kontrollert |
| G6-M2 | Typet utvalg av relasjoner, retning, dybde, nivå og mode ved klikk/CLI, revisjon og publisering av bare valgt dokument med ressurser; samme innhold som tilsvarende full eksport; feil beholder siste visning |
| G6-M3 | XFMD med navigasjons-/hovedpanel, registrert leseradapter og eksplisitt vindu/panel; klikk, fokusbytte, flere vinduer og lukket mål testet |
| G6-M4 | Valgfri Go-bakgrunnstjeneste med lokal IPC, cache/invalidering, leser-lease, request-rekkefølge, kvoter og opprydding; ingen døde bilder ved dokumentbytte/reload |
| G6-M5 | Fast symbol-/pilprofil med UML der semantikken stemmer; aktørfigurer og use-case-ellipser; rendererprøver kontrollerer faktiske figurer/markører, ikke bare exitkode |
| G6-M6 | Senere eksplisitt klasse-/relasjonsprofil med multiplisitet og aggregation/composition; språk/validator før kildekoblede klassediagrammer, ingen automatisk oversettelse fra contains |

G6-M1 avhenger av G4-M1s strukturelle frontendport, ikke SDL-runtime eller G5s
Go-generering. Viewpoint-port og kildekart flyttes fra G5-M3 til G6-M1;
G5-M3 blir konsument av denne eksporten. G6-M2 → M3 → M4 følger hverandre. M5 avhenger av M1 og kan utvikles
ved siden av vertsarbeidet; M6 følger M5 og krever avklart klassekontrakt.
Den nåværende Python-generatoren er portgrunnlaget. URI-notasjon og XFMD-flagg i
designet er forslag og skal verifiseres med XFMD-implementasjonen.
G6-D1/D2 leverer design og parsede planer/scenarioer, ikke implementerte
G6-M1–M6. [Nivåer og notasjon](../../docs/SDL-Viewpoint-Levels-and-Notation.md)
presiserer eksportformene, A0–A5, Mode/State og semantisk diagramprofil.

## Avgrensning og gjenbruk

Fyne er første interaktive SDUI-vert. XFMD-dokumentnavigasjon er planlagt i G6;
FOX-baserte SDUI-widgets, C-ABI og Bubble Tea er fortsatt utsatt. Eksisterende renderer-/worktree-kode
kan gi ideer, tester og egnede algoritmer med dokumentert proveniens; det kreves
ingen uttrekkscrate eller merge til Mermaid. Fullt Markdown/Mermaid-innhold må
prøves mot en avtalt profil. Ingen flere renderere bygges bare for å holde valg åpne.

SDL- og SDUI-katalogene har hver sin README for kodeansvar. Denne planen er felles;
ikke opprett konkurrerende faseplaner under begge. Første kodeleveranse er G1-M1/M2
og en liten G2-M1–M3-prøve etter designgjennomgangen. Git følger
[én branch per fase og commit per milepæl](../../docs/Development-Branch-Stack.md);
push er autorisert etter hver fullført fase.
