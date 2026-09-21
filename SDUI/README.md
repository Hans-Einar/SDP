# SDUI — prototyping av brukergrensesnitt

**2026-09-21: SDUI 0.2 parser, AST, lokal normalisering og konsolldump er kjørbare.**
Frames `[]`, nestede widgetgrupper `<>`, Markdown-strenger og formatering `{}`
er implementert i Python-prototypen. Gammel 0.1-syntaks, grammatikk og `text()`-
widget er erstattet; eksempler og tester er portert. Ingen kompatibilitetsmodus.

[Concept1-eksemplet](examples/concept1-bucking.sdui) beskriver apteringsflaten
med seks hovedbokser og representative kontroller. Verdiene er statiske eksempeldata;
React-komponenter og apteringslogikk er ikke automatisk oversatt eller kjørt.
[AST](examples/concept1-bucking.ast.json) og [GUI-dump](examples/concept1-bucking.dump.txt)
er generert fra samme kilde. Se [kartlegging og TUI-retning](docs/concept1-console.md).
[Markdown-dumpen](examples/concept1-bucking.dump.md) gir samme statiske oversikt
og renderbart Markdown-innhold, uten interaktive widgets eller Mermaid-avhengighet.
Den viser innholdet separat og oppfyller ikke målet om en samlet grafisk GUI-dump.
En [forenklet treemap-prøve](examples/concept1-bucking.treemap.md) viser boksene
og kort tekst samlet, som et forhåndsrendret SVG-bilde i Markdown.
En [utvidet widgetprøve](examples/concept1-bucking.widgets.md) tegner knapper og
inputfelt i samme boksinndeling. [HTML-galleriet](examples/prototype-controls.html)
lar deg fylle ut felter, prøve knappetrykk og skrive ut aktuelle verdier i nettleseren.
Dette er avgrensede prototypevisninger, ikke den kommende felles layoutmotoren.
Se [widgetbibliotek og avgrensninger](docs/prototype-widgets.md).

## Kjør

Python 3.11+, kun standardbiblioteket. Fra SDUI-katalogen:

```sh
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format dump --entry bucking --columns 160
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format dump --entry page --columns 200
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format markdown --entry bucking -o examples/concept1-bucking.dump.md
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

`-o path` lagrer resultatet; `-` som kilde leser stdin. Standardformat er JSON-AST.
`--syntax-only` hopper over profilvalidering og kan ikke brukes med dump.
Feil kommer som JSON på stderr: exit 2 for språk/dump, 3 for I/O.
Parseren åpner eller kjører aldri SDL-referanser.

Dumpen viser rå Markdown, utelater Mermaid-fences og viser knapper/input som
tekstkontroller. Den viser rader, grupper, bokser og horisontale vekter;
vertikale vekter, ratio og native tekstmåling gjenstår i den felles layoutmotoren.
Høydene i dumpen følger innholdet. Dette er en strukturell konsollforhåndsvisning.
Markdown-formatet legger deretter til innhold med overskrifter, lister, tabeller
og kodeblokker. Nesting vises med sitatblokker og radrekkefølge med etiketter;
innholdsdelen gjengir ikke GUI-kolonner. Tabeller krever en leser med tabellstøtte.
Se [verifikasjon av Markdown-dump](evidence/markdown-dump-2026-09-21.md).

## Dokumentasjon og videre arbeid

- [Implementert språk](docs/language.md) og [EBNF](grammar/sdui-0.2.ebnf).
- [Kjørbar arkitektur](docs/architecture.md), [krav](docs/requirements.md) og [testbevis](evidence/frontend-console-2026-09-21.md).
- [Mandat](Mandate-and-Study.md), [layoutforslag](docs/layout-language-proposal.md) og [gruppe-/regiondesign](docs/frame-composition-proposal.md).
- [Implementasjonsplan](docs/implementation-plan.md), [Go-målarkitektur](docs/target-architecture.md) og [runtime-kontrakt](docs/runtime-contract.md).
- [Egen SDL-designmodell](design/README.md), [worktree-kart](docs/renderer-extraction-and-language-direction.md) og [XFMD-handoff](docs/handoff-xfmd-sdui.md).

Neste leveranser er parser/runtime i Go, felles layout, SVG-eksport, Fyne-visning
og modellreload; senere SDL-kobling og Go-generering. [Go-området](go/README.md)
har nå parser/runtime-kataloger, men ingen Go-kode. Fyne er første vert;
FOX/XFMD og interaktiv TUI er utsatt. Se
[checkpointets oppdatering](../docs/checkpoint%231/07-SDUI-0.2-and-Go-Direction.md).
Python er kjørbart portgrunnlag til Go erstatter de aktive funksjonene.
Ingen Fyne-integrasjon, SDL-runtime eller generell layoutmotor er implementert.

[Hovedeksemplet med symbolske bindinger](examples/main-page.sdui) og
[statisk eksempel](examples/static.sdui) bruker også bare den nye profilen.
[0.1-verifikasjonen](evidence/verification.md) er historisk og beskriver utgått kode.
