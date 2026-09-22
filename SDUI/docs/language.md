# SDUI 0.2 — implementert kildeprofil

Oppdatert 2026-09-22. Go-frontenden parser SDUI 0.2 og bygger AST.
Runtime/Fyne/SVG er implementert separat; parseren utfører aldri domenekall.
[EBNF](../grammar/sdui-0.2.ebnf) og denne profilen erstatter gammel 0.1-syntaks.
[Layoutforslaget](layout-language-proposal.md) beskriver også framtidig geometri;
akseptert formatering er ikke det samme som implementert layoututførelse.

## Kilde og struktur

Dokumentet begynner med eksakt `sdui 0.2;`, fulgt av eventuelle `ref:`-deklarasjoner,
en eller flere navngitte komponentdefinisjoner og eventuelle `setHandle`-koblinger.
Ingen gammel profil/fallback. Definisjoner kan referere framover til andre
komponentdefinisjoner. En konsument velger inngangsframe; definisjoner monteres
ikke automatisk. CLI krever `--entry` hvis flere frames kan velges.

```text
sdui 0.2;
mainBody = <<"Gruppe 1", "Mer tekst"> {v<}, <button("OK")>>;
page = [header="## Aptering", body=mainBody, footer="Statisk prototype"]*b {16:9, <->, font=10};
```

`[]` er en frame. `<>` er en widgetgruppe som kan nestes; eksplisitte grupper
bevares også med ett barn. Frames kan inneholde frames, grupper, Markdown,
widgetkall og referanser. Widgetgrupper kan inneholde alle disse unntatt frames.
Tom frame og tom gruppe er gyldige. Strenger som komponenter er Markdown.
`button("OK")` har vanlig etiketttekst; `text(...)` fra gammel profil er fjernet.

Komma fortsetter raden horisontalt; semikolon begynner neste rad under hele
foregående rad. Separatorer tilhører sin egen liste. Avsluttende semikolon er
valgfritt inne i en container, obligatorisk etter en toppnivådefinisjon.
Trailing comma og tomme rader avvises. Formatering `{...}` følger komponenten,
etter eventuell `*box`/`*b`, før separatoren. En løs formateringsblokk avvises.

I frames er `header`, `body`, `footer` regionroller, høyst én av hver. Eksplisitt
body kan ikke blandes med umerket bodyinnhold. Regioner trekkes ut ved
normalisering; komposisjonen er header, body, footer uavhengig av kilderekkefølge.
I widgetgrupper er slike tilordningsnavn vanlige navn.

## Formatering og navn

Lokalt validerte egenskaper: scale/scale-x/scale-y, x/y=content/fill/Nfr,
min-x/min-y/max-x/max-y, gap/gap-x/gap-y, padding eller fireleddet paddingtuple,
align-x/align-y, justify, items, overflow-x/overflow-y, wrap, font, enabled, visible.
Tillatte verdier står i [layoutforslaget](layout-language-proposal.md).
Ukjente egenskaper, duplikater og konkurrerende størrelsesregler avvises.
Bare relative layoutdimensjoner; `font` er positiv absolutt tekststørrelse.
Den native fontenheten er ikke fastsatt. Dumpen skalerer ikke terminalens font.

Kanoniske hjørner er ^<, >^, v<, >v; omvendt retningspar er ekvivalent.
Kortformer senkes til egenskaper i normalisert modell, mens AST bevarer stavingen.
`<->` og `>-<` er forskjellige hele operatorer. Ratio x:y gjelder bare frames;
én scaleakse eller én fillakse kan styre. To akser, scale på begge, content/fr
med ratio eller ikke-positive forhold avvises. Ingen geometriløsning hevdes.
Wrap gjelder én grupperad uten horisontal fill/fr på gruppen eller dens barn.

Navn er ASCII `[A-Za-z_][A-Za-z0-9_]*`, case-sensitive; reserverte ord er
sdui/ref/true/false/null/setHandle. Egenskapsnavn kan i tillegg ha bindestrek.
Lokale navn er unike i hver definisjon. Gjenbruk får separat instansbane:
`left=mainBody` og `right=mainBody` gir ulike widgetbaner. Ualiasert referanse
bruker definisjonsnavnet. Sykler og tvetydige instansbaner avvises.
Anonyme grupper får interne kildeposisjonsbaserte banesegmenter; offentlige
setHandle-baner utelater disse. Runtime-reglene for kompatibel identitet og reload står i Go-runtimekontrakten.

## Widgets og symbolske koblinger

| Kall | Første argument | Øvrige argumenter |
| --- | --- | --- |
| button | label: streng | callback: symbolsk referanse |
| input | text: streng | value: streng, callback: symbolsk referanse |
| svg | source: symbolsk referanse | label: streng |

Første argument er obligatorisk, posisjonelt eller navngitt; deretter bare
navngitte argumenter. Duplikater og feil typer avvises. Anonyme ubundne widgets
er gyldige. Callback krever widgetnavn; eksterne handles adresserer navngitte widgets.
SVG-kilden kjøres aldri av parseren eller dumpen; dumpen viser en plassholder.

```text
ref: domain "missing-on-purpose.sdl";
# inne i en komponent:
# field=input("Navn", callback=domain.object.@changed)
# etter definisjonene:
# domain.object.setHandle(page.form.field);
```

Referansen har formen module.object.@member. Aliaset må finnes; ekstern fil,
objekt og metode undersøkes ikke. setHandle må treffe en widgetinstans gjennom
definisjon og navngitte foreldre. Én kobling per SDL-objekt og widget. Ingen
vilkårlige metodekall, imports, interpolering eller eval.

## Strenger, posisjoner og grenser

UTF-8 med Unicode-skalarposisjoner og halvlukkede UTF-8-byteområder. Linje/kolonne
starter på 1; LF starter ny linje, også i CRLF. `#` starter kommentar utenfor streng.
Enkle/doble strenger støtter n/r/t, backslash, anførselstegn og fire hexsifre etter
u. Rå kontrolltegn, NUL og surrogate escapes avvises. Triple doble anførselstegn
bevarer multiline Markdown uten dedent, escaping eller interpolering.
Tall følger JSON-lignende syntaks og endelig binary64; ingen garanti over 2^53.

Grenser: 256 KiB kilde, 50000 tokens, 64 syntaktiske/ekspanderte nivåer,
2048 kildekomponenter, 32 widgetargumenter, 32 formateringsregler per komponent,
8192 ekspanderte komponenter totalt over dokumentets definisjoner.
Feil er strukturert med code/message/span; ingen recovery eller delvis CLI-output.
`--syntax-only` bevarer også ukjente widgets/formateringsnavn uten profilgodkjenning.
EOF kreves. Ressursgrenser er vern mot ubegrenset arbeid, ikke hard sanntid.
