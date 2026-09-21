# SDL/SDUI — målarkitektur i Go

**ID:** SDUI-ARCH-003 · 2026-09-21 · Valgt retning, implementasjon gjenstår.
Erstatter ARCH-002s Rust/C-ABI og obligatoriske FOX/XFMD-løp.
[Checkpoint #1, tillegg 07](../../docs/checkpoint%231/07-SDUI-0.2-and-Go-Direction.md)
eier beslutningsoversikten; [architecture.md](architecture.md) beskriver kjørbar Python-kode.

**Detaljert SDL-modell 2026-09-22:** [parser-/runtimedesignet](../design/README.md)
dekker begge språkene og deres porter, reload, layout, vert og kodegenerering.
Modellen er validert med eksisterende design-core-parser; ingen Go-runtime kjøres.

## Eiendom og avhengigheter

| Område | Ansvar |
| --- | --- |
| SystemDesignLanguage/go/parser | SDL-kilde, AST, symboler og profilvalidering; ingen domeneutførelse |
| SystemDesignLanguage/go/runtime | Avgrenset SDL-kjøring og registrerte Go-funksjoner; ingen GUI-avhengighet |
| SDUI/go/parser | SDUI 0.2, AST, diagnoser og normalisering; én språkimplementasjon etter port |
| SDUI/go/runtime | UI-instans, identitet, egenskaper, events, bindinger og modellreload |
| Senere SDUI layout | Én målt layoutmodell for interaktiv visning og eksport |
| Senere SVG-presentasjon | Statisk dokumentasjonsbilde fra samme geometri og eksplisitt tilstand |
| Senere Fyne-vert | Vindu, widgetlivstid, fokus, inndata og publisering på UI-tråden |
| Senere Go-generator | Modell-/koblingskode; håndskrevne domenefunksjoner ligger separat |

Bare de fire parser/runtime-katalogene er opprettet. Pakke-/modulnavn og
Go/Fyne-versjon fastsettes i første kodeleveranse. Ingen go.mod, go.work eller
påstått fungerende Go-API er lagt til som tom fasade.

Parser/runtime importerer ikke Fyne, FOX, XFMD eller Mermaid. En vert setter
sammen bibliotekene; SDUI-kjernen krever ikke en konkret SDL-implementasjon for
å vise et ubundet design. SDL eier domenetilstand, SDUI eier widgetidentiteter og
UI-tilstand. Koblinger går over typede Go-grensesnitt; ingen ekstra binær ABI før
en konkret konsument krever den. Navnet libsdui er bibliotekrollen, ikke krav
om en .so-fil eller C-header.

## Måling og presentasjon

Kildens relative størrelser, ratio og ancestorreferanse beholdes fram til layout.
Root får vertens tilgjengelige område. `{16:9,<->}` avleder høyden fra bredden;
resize endrer geometri og tekstombryting uten å skalere fonten. Absolutt fontenhet
og målekontrakt fastsettes før layout bygges. Header/body/footer ligger innen ratio.

Layoutresultatet inneholder widgetidentiteter, rektangler, klipping, tekstmål og
ressurser. Fyne bruker dette til interaktiv visning; SVG-eksport bruker samme
resultat og tema. Ingen separat Fyne-layout som tolker språkreglene på nytt.
Fyne-kontroller gjenbrukes særlig for tekstinntasting, fokus og tastatur. Et helt
UI som ett SVG-bilde er ikke interaktivt uten ekstra hendelses- og treffhåndtering.

Fyne har SVG-bilder og egne widgetrenderere, men bilde-/tekst-/klippstøtten for vår
konkrete SVG-profil må prøves. Ingen garanti om pikselidentiske rendererresultater.
[SVG-bilder](https://docs.fyne.io/canvas/image/),
[widgetrenderere](https://docs.fyne.io/extend/custom-widget/).

Markdown er fortsatt ønsket innhold, inkludert Mermaid. En avgrenset innholdsport
må definere profil, tekstmåling, ressurser og diagramdekning. Fyne RichText er ikke
automatisk full XFMD-kompatibilitet. Vi bygger ikke en ny Mermaid-parser.
Første prøve krever enkel Markdown; full innholdsdekning er en egen milepæl.
Dokumentasjon kan inkludere SVG som et statisk bilde uten SDUI-runtime.

## Tilstand, binding og reload

Logisk widgetreferanse består av sesjon, instansbane og generation. Omplassering
bevarer kompatibel identitet; sletting/typebytte invaliderer gamle referanser.
UI skiller akseptert verdi fra brukerens draft. Programmatisk endring er ikke et
nytt brukerklikk. Runtime kontrollerer typer, enabled og revisjon før dispatch.
Se [runtime-kontrakten](runtime-contract.md) for videre kontraktarbeid.

Reload bygger en ny validert modell før publisering. Kilde-/bindingsfeil beholder
siste gyldige modell. Kompatibel verdi/fokus kan videreføres; inkompatible endringer
krever diagnose og eksplisitt reset/migrering. Utestående callbacks kanselleres eller
avvises med generation/revisjon; reload skal ikke gjenta domenehandlinger.

SDL/SDUI-modell kan skiftes i den kjørende Go-runtime. Endringer i Go-funksjoner
krever bygg/restart i første løsning. Prosessbytte og stateoverføring er senere
muligheter, ikke innlastning av modifiserte plugins.

## Kodegenerering og portering

Første Go-generering oppretter samme modell/bindinger som filbasert utviklingsmodus;
begge bruker samme runtime. En senere direkte kodeoversetting av atferd krever
samsvarstester mot definert kjøresemantikk. Ikke implementer to SDL-semantikker.
Genererte filer overskriver aldri håndskrevet domene-Go.

Python-frontender og relevante tester er portgrunnlag. De fjernes som aktive
implementasjoner først når Go-porten dekker deres avtalte profil og brukere er
portert. Dette gir ingen bakoverkompatibilitet for SDUI 0.1. SDL design-core 0.1
er en annen profil og må ikke slettes bare fordi den har samme versjonstall.

Tidligere FOX/XFMD-/Mermaid-arbeid er gjenbruksgrunnlag, ikke en forutsetning.
Ingen eksterne worktrees eller eksisterende produkter endres av denne retningen.
