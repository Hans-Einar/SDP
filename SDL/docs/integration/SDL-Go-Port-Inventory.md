# Go-port — aktive innganger etter G5-M4

| Erstattet aktiv vei | Gjeldende konsument / bevis |
| --- | --- |
| SDUI/src/sdui: lexer, parser, AST, validate, formatting, normalize | SDUI/go/parser; 93 fryste porttilfeller og alle eksempel-AST-er |
| SDUI Python CLI, dump, markdown_dump | SDUI/go/cmd/sdui og presentation; identiske Concept1-dumper og feilveier |
| Python prototype_widgets/html og tools/build_widget_previews | Generell go/layout + go/svg; go/host/fynehost for interaksjon; state-eksport via sdl-document |
| Midlertidig Go prototype-svg/prototype-html | Fjernet ved portslutt; --format svg bruker felles geometri. Eksisterende HTML-galleri beholdes kun som merket historisk artefakt |
| experiments/design_core/*.py | SystemDesignLanguage/go/parser; 151 fryste parser-/diagnose-/kanoniseringstilfeller |
| SystemDesignLanguage/tools Python-projectorer/CLI | go/viewpoint, documents og cmd/sdl; 170 fryste diagramprojeksjoner og alle 11 viewpoints |
| SDUI/tools/export_design.py | sdl check/ast/viewpoints; ansvarsoversikt i VP02, faktaregister i VP11 og felles manifest |
| Python capture_port_cases | Fjernet etter frysing; original referansekilde ligger sammen med prosjektorfixturen |

Ingen runtime faller tilbake til Python. Nye parser-/renderer-/CLI-prøver er Go;
Python-scriptet for native broker/FOX-prøver er testautomatisering, ikke SDL/SDUI-
implementasjon. Historiske daterte logger/hashmanifest er bevis for gamle commits,
ikke oppdaterte byggeinstrukser. Python-koden kan finnes igjen i Git før G5-M4.

Ikke erstattet og ikke slettet: Toolkit/SDP-verktøy, MVP1-kandidatøvelsen og andre
uavhengige språkeksperimenter. De er ikke alternative aktive frontender for
SDUI 0.2 eller design-core 0.5. XFMD BoxUI og Mermaid-worktrees tilhører separat
arbeid; G6 bruker en ny, isolert XFMD-fasebranch og uendret Mermaid-backend.

Go-generering bygger typed modeller/bindingsdata; parsere, runtime, layout og
projektorer er fortsatt én implementasjon per profil. Detaljer og kommandoer:
[checkpoint](../../../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md),
[SDUI](../../../SDUI/go/README.md), [SDL](../../go/README.md).
