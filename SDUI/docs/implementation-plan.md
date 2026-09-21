# SDL/SDUI — implementasjonsplan i Go

**Anbefalt prioritering 2026-09-22:** [SDL-viewpoint-løpet V0–V4](../../docs/checkpoint%231/08-SDL-Viewpoints-and-Implementation-Status.md)
konkretiserer modell-/språkhull før mer UI-runtime. V0 og V1 er levert som SDL-verktøy med strukturelle visninger, UseCase/Feature og
modusallokering. V2-data og V3-Channel-scenarioer må avklares
før full viewpoint-dekning. G1–G5 beholdes som Go-plan, ikke markert ferdig av diagrammer.

**ID:** SDUI-PLAN-003 · **Revisjon:** 2026-09-21.
Erstatter PLAN-002s P0–P6-løp for Rust/C-ABI/FOX. Nye milepæler bruker G-prefiks;
henvisninger til P-faser i eldre bevis gjelder historien, ikke aktive leveranser.
[Målarkitektur](target-architecture.md) og
[checkpoint](../../docs/checkpoint%231/07-SDUI-0.2-and-Go-Direction.md).

Omfang: SDL/SDUI-parser og runtime i Go, felles SDUI-layout, SVG-eksport, første
Fyne-vert, modellreload og senere Go-generering. Bare G0 er levert i denne runden.
Python SDUI 0.2, strukturparseren design-core og eksemplene er eksisterende portgrunnlag.

Designgrunnlag 2026-09-22: [felles SDL-strukturmodell](../design/README.md) beskriver
ansvar og avhengigheter for G1–G5 og passerer eksisterende parser. Dette er
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
| G1-M1 | Modulstruktur, Go-baseline og CLI; parser/AST med kildeposisjoner for SDUI 0.2, uten GUI |
| G1-M2 | Validator/normalisering, relative regler, frame-regioner, instansbaner og kildegrenser; porterte positive/negative Python-tilfeller |
| G1-M3 | Concept1 og bindingseksempel gir kontrollert AST/diagnoser; konsoll-/Markdown-dump og prototypeeksport har en portert erstatning før Python-veien fjernes |

Krav R01–R09, R13–R15, R19, R21–R24. Ingen SDL-kilde åpnes av parseren.
Ingen Go- eller Python-fallback for SDUI 0.1. Parallell Python brukes bare som
midlertidig portorakel; overgangens fullføring avhenger også av G2s eksport.

## G2 — én layoutmodell, SVG og første Fyne-vindu

Avhenger av G1-M2. Eier: SDUI-implementasjonen.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G2-M1 | Målekontrakt og fontenhet; relative akser, ratio, rader/grupper, header/body/footer, gap/padding og eksplisitte overflowgrenser |
| G2-M2 | Generell SVG-eksport fra målt modell; liten frame + knapp/input og deretter Concept1 uten håndplassert fixture; geometri- og visuell kontroll |
| G2-M3 | Fyne-vindu med samme geometri, knapper/input, Tab/fokus og registrert lokal Go-funksjon; ingen SDL-avhengighet |
| G2-M4 | Avgrenset Markdown-provider og separate Mermaid-ressurser; dokumentert støtteprofil og negative tilfeller, ikke påstått full diagramdekning |

Krav R10/R11/R14/R15/R17/R20/R21/R24. Begynn med liten vertikal prøve i M1–M3;
full Concept1/Markdown følger. Mål oppstart, resize og hukommelse før påstander
om «lettvekts». Fyne håndterer widgetinteraksjon; kjerne og SVG skal fungere uten GUI.
SVG-eksport gjengir valgt tilstand og utfører aldri callbacks.

## G3 — UI-runtime og SDUI hot reload

Avhenger av G1-M2/G2-M3. Eier: SDUI-implementasjonen.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G3-M1 | Typede hendelser/oppdateringer, stabile widgetinstanser og revisjoner; ubundet knapp gir eksplisitt status |
| G3-M2 | Filendring → parse/valider → publisering; ugyldig kilde beholder siste gyldige UI og viser kildediagnose |
| G3-M3 | Bevar kompatibel verdi/draft/fokus; typebytte/sletting håndteres, stale events og callback etter teardown avvises; ingen gjentatte domenehandlinger ved reload |

Krav R12/R16/R18/R25. UI-state og domene-state holdes adskilt. Avtal eventgrense
og UI-tråd før asynkronisering; ingen kompleks plugin-/prosessmekanisme i denne fasen.

## G4 — avgrenset SDL-parser/runtime og kobling til SDUI

SDL-parserarbeid kan starte ved siden av G1. Integrasjon avhenger av G3-M1.
Eier: SDL-implementasjonen for semantikk/runtime; SDUI eier UI-siden av porten.

| Milepæl | Leveranse og akseptanse |
| --- | --- |
| G4-M1 | Port design-core-struktur til Go med dokumentert grammatikk og tester; checkpoint/MVP1-kandidater blir ikke automatisk støttet |
| G4-M2 | Definer én kjørbar profil for navngitt handling, typed input/resultat og binding til registrert Go-funksjon; negative og manglende bindinger avvises |
| G4-M3 | SDUI-knapp/input → SDL-handling → Go-funksjon → UI-oppdatering; kildekart og samme kontrakt med eksplisitt simulert domene |
| G4-M4 | Reload av SDL-modell med siste gyldige versjon, tilstandsregel og håndtering av pågående hendelser; Go-endring bygges/restartes |

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

## Avgrensning og gjenbruk

Fyne er første interaktive vert. FOX/XFMD-integrasjon, C-ABI og Bubble Tea er
utsatt uten aktiv implementasjonsleveranse. Eksisterende renderer-/worktree-kode
kan gi ideer, tester og egnede algoritmer med dokumentert proveniens; det kreves
ingen uttrekkscrate eller merge til Mermaid. Fullt Markdown/Mermaid-innhold må
prøves mot en avtalt profil. Ingen flere renderere bygges bare for å holde valg åpne.

SDL- og SDUI-katalogene har hver sin README for kodeansvar. Denne planen er felles;
ikke opprett konkurrerende faseplaner under begge. Første kodeleveranse er G1-M1/M2
og en liten G2-M1–M3-prøve. Ingen commits, push, installasjon eller worktree-sletting
inngår i katalog-/checkpointoppdateringen.
