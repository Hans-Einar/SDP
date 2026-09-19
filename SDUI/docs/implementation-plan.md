# SDUI — leveranseplan og beslutninger

**ID:** SDUI-PLAN-001 · Avgrenset bestilling 2026-09-19.

## 1. Denne leveransen

| Trinn | Resultat | Status |
| --- | --- | --- |
| P0: grunnlag | Les faktisk SDL/Concept1/BoxUI, avklar autoritet og arbeidsmappe | Utført |
| P1: språk | EBNF, navn/scopes, widgets/rader, referanser og AST-kontrakt | Implementert arbeidsprofil |
| P2: prototype | Lexer, parser, lokal validator, CLI, positive/negative tester | Implementert; se testbevis |
| P3: overlevering | Praktisk kilde, generert AST, dokumentert dekning og begrensninger | Utført |

Dette er ikke en installasjon eller et versjonsbytte av XFMD. All ny kode og
prosa ligger i SDUI; pågående SDL-arbeid og repositoryets Toolkit er bevart.

## 2. Videre arbeid

1. Gjennomgå den konkrete syntaksen og størrelsessemantikken mot et skjermbilde/
   skjema fra Concept1. Velg testbare kriterier for «samme skjema».
2. Definer normalisert UI-modell og Concept1-adapter. Ikke oversett weight til grow
   uten eksplisitt bevaringsregel og geometriprøver.
3. Trekk BoxUI-kjernen ut av Mermaid-forken med uendret SVG/kontrollkart først.
   Bevar fungerende XFMD-input, kildeidentitet og PDF som regresjoner.
4. Koble SDUI til kjernen og XFMDs Markdown-fence. Vis ytre/indre bokser og målte
   widgetrader. Test lys/mørk, zoom, resizing, tekst, input og PDF visuelt.
5. Implementer runtime-grensen når SDL-objekter, signaturer og kjøresemantikk er
   fastlagt. Bruk en uavhengig simulert adapter til kontrakttester underveis.
6. Vurder direkte Cairo/Pango-presentasjon og felles SDL-IR som separate endringer.

Ikke gjenimplementer tekstredigering eller Mermaid-diagrammer i SDUI-parseren.
Ikke krev at alle framtidige widgets må ha egne grammatikker: Widget-noden er
generisk, mens den versjonerte profilen avgjør hvilke kinds/properties som støttes.

## 3. Beslutningslogg

| ID | Beslutning / begrunnelse | Status |
| --- | --- | --- |
| SDUI-D01 | Egen SDUI-katalog ved faktisk SDL-utviklingskatalog; ingen endring av SDL-authority | Utført |
| SDUI-D02 | ASCII-navn, Unicode-strenger, eksakt profilheader og fail-fast diagnoser | Valgt prototypeprofil |
| SDUI-D03 | `name =` er identitet; ikke separat handle-property | Valgt prototypeprofil |
| SDUI-D04 | Semikolon bryter widgetrad, eksplisitt axis organiserer bokser | Valgt prototypeprofil |
| SDUI-D05 | Modulpath quotes; callbacks er module.object.@member; setHandle er deklarasjon | Valgt prototypeprofil |
| SDUI-D06 | AST uten defaultnormalisering; separat lokal validator | Implementert |
| SDUI-D07 | Python-standardbibliotek som referanseprototype | Implementert, ikke produksjonsspråkvedtak |
| SDUI-D08 | Kjøring, visuell layout, runtime-ABI og Mermaid-uttrekk er neste leveranser | Avgrensning |

## 4. Bevaringskrav

Syntaksutvidelser må bevare eller versjonere AST og tidligere eksempler.
Framtidig rendererbytte må bevare identiteter, hendelsesbindinger og kildekartlegging.
Kjørbar SDUI må ikke gjøre eksisterende Markdown-import eller PDF til en skjult
oppstart av domenehandlinger. Kandidat-SDL-begreper blir ikke vedtatt fordi en
syntetisk callback kan demonstreres.
