# SDUI — frame-, widget- og layoutforslag

**ID:** SDUI-LAYOUT-002 · **Dato:** 2026-09-21 · **Status:** designforslag.
Dette bevarer eierretning og forslag fra 2026-09-21. **Leserveiledning
2026-09-24:** SDUI 0.2 er nå implementert i Go; Python og 0.1 er utgått.
[Språkprofil](language.md), [målekontrakt](go-layout-contract.md) og
[Go-runtime](../go/runtime/README.md) avgrenser hva som faktisk er valgt/levert.
Eksempelvis er scroll fortsatt avvist ved layout, font måles i DIP, og
[Markdown-profilen](markdown-provider.md) er avgrenset. Forslag nedenfor utvider
ikke disse kontraktene. [Implementasjonsplanen](implementation-plan.md) eier milepælene.

## 1. Minste nyttige språk

Målet er skjemaer, kontrollpaneler og dokumentflater med native input, knapper,
Markdown og diagrammer. Vi trenger forutsigbar nesting, størrelse, avstand,
justering og overflow. Vi trenger foreløpig ikke CSS, fri posisjonering,
animasjonslayout, spans over ruter eller en generell constraint solver.

| Form | Betydning i forslaget |
| --- | --- |
| `page = [...];` | Navngitt frame-definisjon; ingen ekstra rot-tilordning eller `{}` rundt definisjonen. |
| `mainBody = <...>;` | Navngitt widgetgruppedefinisjon som kan instansieres med en referanse. |
| `[...]` | Frame; kan inneholde underframes og widgetgrupper i samme innholdsliste. |
| `<...>` | Widgetgruppe; widgets, Markdown-strenger og nestede widgetgrupper. |
| `name = button("OK")` | Navngitt widget; argumentstrengen er en etikett, ikke en Markdown-widget. |
| `"# Forklaring"` som listeelement | Kortform for en Markdown-widget. Kan navngis med `intro = "..."`. |
| `]*box` | Framevariant med ramme/overskrift. `*b` er et fast alias for `box` i denne profilen. |
| `node { ... }` | Formatering/layout på komponenten foran; gjelder også Markdown-strenger og komponentreferanser. |

Eierens idé om entydige prefikser er bevart i historikken. Anbefalt første løsning
er det eksplisitte aliaset `b`, slik at installerte komponenter ikke endrer
betydningen av gammel kilde. Normalisert modell lagrer `box`. Flere varianter
krever registrerte egenskaper/måling/presentasjon; `*` er ikke dynamisk kodelasting.

En frame er en layout-/klippegrense, ikke nødvendigvis et FOX-vindu. Udekorert
frame tegner ingen kant. `header=`, valgfri `body=` og `footer=` er innholdsregioner;
vanlig umerket innhold blir body. `heading` erstattes av header, ikke et alias.
Tomme frames/grupper er lovlige plassholdere uten iboende innholdsstørrelse.
[Komposisjonsforslaget](frame-composition-proposal.md) beskriver regioner, referanser og font.

Lokale navn er unike per komponentdefinisjon; instansbanen skiller gjenbrukte grupper. Widgets med callback eller eksternt
handle må navngis; ubundne knapper/input kan være anonyme. Anonym statisk Markdown kan få kildebasert intern ID,
men ingen garanti om stabil identitet ved kildeendring. Layoutendring alene skal
ikke endre navngitte widgetidentiteter. `ref`/callback og deklarativ `setHandle`
kan videreføres som bindingsdata; parseren åpner eller kjører aldri SDL-kilder.

## 2. Rader, grupper og akser

Eierpresisering 2026-09-21: hver UI-komponent kan ha en etterfølgende `{...}`.
Blokken tilhører hele komponenten foran og kommer **før** dens `,` eller `;`.
Dette gjelder navngitte og anonyme widgets, Markdown-strenger, widgetgrupper,
frames og komponentreferanser, også som verdier for header/body/footer.
Rekkefølgen er komponent → eventuell `*variant` → eventuell `{...}` → separator.
Siste komponent kan avsluttes direkte av omsluttende `>` eller `]` etter blokken.
Forslag: én samlet formateringsblokk per komponent; egenskapene valideres mot
komponenttypen. En løs blokk etter separator har ingen komponent og avvises.

```text
mainBody = <
  "Tekst" {font=10}, ok = button("OK") {font=12};
  "Neste rad" {font=10}
> {<->};
page = [mainBody {font=12}]*b {16:9, <->};
```

På en referanse gjelder formateringen den aktuelle instansen. Forslag til
fontarv: frame-/gruppefont er default for etterkommere; eksplisitt font på et
barn overstyrer den. Referanseformatering endrer ikke originaldefinisjonen.

Komma skiller elementer i samme rad; semikolon starter neste rad. Dette gjelder
både frame-body og `<>`. Header/body/footer er regionroller, ikke elementer
i den samme kommaraden. Komma inne i widgetargumenter/layout tilhører den indre
konstruksjonen. Parseren tolker aldri `<` inne i en Markdown-streng som widgetstart.
Neste rad begynner under hele den foregående radens utstrekning, med valgt gap.
Semikolon etter en toppnivådefinisjon avslutter definisjonen; det lager ingen UI-rad.
`<>` kan nestes i flere nivåer med egen layout per gruppe. Separatorer gjelder
bare sin egen liste; en undergruppe er ett element i morens liste. Eksplisitte
grupper bevares ved normalisering, også med ett barn. Se
[gruppeeksemplet](frame-composition-proposal.md#nestede-grupper-med-egen-layout).

Normalisering lager en vertikal stabel av rader. Én rad er en horisontal gruppe.
Rader deler ikke automatisk kolonnebredder: dette er ikke et tabell-/grid-språk.
En rad med ett element senkes uten ekstra geometrisk mellomledd, slik at
elementets høyde/vekt gjelder direkte i den vertikale stabelen. Bruk navngitte
underframes for eksplisitte spor når flere kolonner skal dele samme høyde/vekt.

Eksempel med to kolonner over en fullbredde fot:

```text
sdui 0.2;
page = [
  tools = [<"## Verktøy"; run = button("Kjør")>]*box {scale-x=0.25},
  content = [<"## Resultat">]*box {x=1fr};
  status = <"Klar"> {x=fill}
] {scale=1, gap=0.01, padding=0.01};
```

I flerradsformen får en rad med flere elementer innholdshøyde. Vil man fordele
høyde mellom slike rader, pakk hver rad i en eksplisitt frame som i panel-eksemplet.
Ikke innfør skjult propagering av én tilfeldig barnewidgets `y=1fr` til raden.

## 3. Kanoniske former, relative mål og sideforhold

Eierpresisering 2026-09-20: layoutdimensjoner i kilden skal være relative til nærmeste
ancestor. Ingen width/height i piksler. `x:y` angir frame-sideforhold, og skalering
kan styre x, y eller begge; med låst sideforhold kan bare én akse styres.
Konkrete navn som `scale-x` nedenfor er forslag til skrivemåte for disse reglene.

### Mini arrows og kanonisk skrivemåte

Retningsparet angir to uordnede justeringer. `v<` og `<v` betyr begge ned/venstre,
men formatterer/eksempler bruker `v<`. Eierens kanoniske hjørner er nå
`^<`, `>^`, `v<`, `>v`. `<^` normaliseres som `^<`, og `^>` som `>^`.
Kilde-AST bevarer opprinnelig skrivemåte for diagnoser.

| Kanonisk form | Betydning |
| --- | --- |
| `^<`, `-^`, `>^` | Venstre/midt/høyre, øverst |
| `<-`, `--`, `->` | Venstre/midt/høyre, vertikalt sentrert |
| `v<`, `-v`, `>v` | Venstre/midt/høyre, nederst |
| `<->` | Bruk hele tildelte bredden; med ratio avledes høyden |
| `^\|v` | Fyll tildelt område vertikalt |
| `>-<` | Pakk til innholdets horisontale utstrekning |
| `>\|<` | Pakk til innholdets vertikale utstrekning |

Vertikale former staves `^|v` og `>|<`; tabellen escaper bare Markdown-tegnet.
Uavhengige regler i blokken kan også stå i vilkårlig rekkefølge. `{v<, >-<, >|<}`
pakker innholdet og plasserer det nedre venstre i det tildelte området.
Innover-/utoverformene er hele operatorer: `<->` og `>-<` har ulike betydninger,
selv om de bruker samme tegn. Permutasjonsregelen for retningspar gjør dem ikke
like. Det bevarer eierens forskjell mellom å strekke og å pakke.

`¤` får foreløpig ingen betydning. Det er en ledig kandidat til en senere tydelig
operasjon; vi innfører ikke et tilfeldig synonym eller en parserregel nå.

### Nærmeste ancestor og relative dimensjoner

En eksplisitt frame/widgetgruppe etablerer referanse for barna. En node bruker
sin nærmeste slike ancestors **indre tilgjengelige rektangel**, etter header/footer,
dekorasjon og padding. Automatisk opprettede rader i normalisert modell endrer
ikke denne kildereferansen. Widgetargumenter er ikke ancestors. Roten bruker
layoutområdet verten stiller til rådighet, uten vindusdekorasjon/verktøylinjer.

`scale-x=0.5` ber om halvparten av referansebredden; `scale-y=0.5` halvparten av
referansehøyden. `scale=0.5` angir begge. Uten sideforhold kan x/y være forskjellige.
Alle scaleverdier er positive, endelige forholdstall; over 1 er lovlig med
avtalt overflowpolicy. Ingen skjult prosent-/pikselkonvertering av bare tall.

Referansen er mors innholdsområde **før** søskenspor fordeles. To barn på 0.5
hver pluss et gap fyller derfor mer enn mor; det skal ikke stille korrigeres.
Bruk relative vekter når søsken skal dele plassen som er igjen etter gap.
Er referanseaksen ubestemt, rapporteres en størrelsesavhengighet som ikke kan
løses; vi hopper ikke stille til en fjernere ancestor eller til hovedvinduet.

### Aspect ratio og styrende akse

`16:9` i layoutblokken betyr ytterbredde/ytterhøyde for framen, med header, body, footer og
padding innenfor denne formen. Begge ledd må være positive endelige tall.
Normalisering lagrer ett forholdstall; `32:18` og `16:9` har samme betydning.

La referanseområdet være W × H og forholdet r = x/y:

- `{16:9, <->}` gir hele tildelte bredde B og høyde B/r; for root er B=W.
- `{16:9, ^|v}` gir hele tildelte høyde og bredde høyde*r.
- `{16:9, scale-x=0.75}` gir bredde 0.75W og høyde bredde/r.
- `{16:9, scale-y=0.5}` gir høyde 0.5H og bredde høyde*r.
- `{16:9, scale=0.75}` og samtidig scale-x/scale-y avvises, også når tallene
  tilfeldigvis ville gitt riktig form. Bare én akse skal styre.
- Med ratio avvises også samtidig `<->` og `^|v`.
- Uten eksplisitt scale eller strekk foreslås contain: største rektangel med forhold r som
  passer referanseområdet. Den begrensende aksen velges, den andre avledes.

Ingen automatisk etterjustering av eksplisitt scale for å få det til å passe;
hvis avledet størrelse går utenfor mor, gjelder overflowregelen. Min/max eller
native minimum må aldri deformere sideforholdet stille. Når sideforhold er
låst, kan den frie aksen ikke få en annen fill/content-/størrelsesregel.
Eksplisitt størrelsesvalg for slike frames kan være én scaleakse eller én
strekkakse. Auto-contain er bare for manglende aksevalg. `<->` skal aldri
nedskaleres til contain fordi vindushøyden er for liten. Blanding med fr
må defineres før støtte; ingen skjult andre styrende akse.

### Øvrig relativ layout

| Egenskap / form | Foreslått første profil |
| --- | --- |
| scale-x, scale-y, scale | Forhold til referanseområdet som definert over; ingen px. |
| x, y | `content` (default), `fill` eller positiv `Nfr`; ingen tall med absolutt lengde. |
| min-x, min-y, max-x, max-y | Relative grenser mot samme referanseakser; valgfrie, min ≤ max. |
| align-x, align-y | start/center/end; default start; kan skrives med mini arrows. |
| gap, gap-x, gap-y | Ikke-negative forholdstall mot referansebredden/-høyden. Default 0.01; aksevariant overstyrer gap. |
| padding | Relativ faktor eller `(top,right,bottom,left)`; vannrett mot referansebredde, loddrett mot referansehøyde. Default 0, eller 0.01 for box. |
| justify | start/center/end/between; styrer restplass mellom barn langs hovedaksen. |
| items | start/center/end/stretch; barnas default på tverraksen. |
| overflow-x, overflow-y | error/clip/scroll; default error. |
| wrap | none/wrap, bare widgetgruppe med én eksplisitt rad i første profil. |
| font | Absolutt tekststørrelse, eksempelvis `font=10`; ingen ancestor-/resizefaktor. Se målekontrakten nedenfor. |

Tall på gap/padding/grenser er altså relative, ikke de tidligere foreslåtte
logiske pikslene. De bruker nodens ancestor-referanse som scale; dette er et
forslag for å unngå sirkulær selvreferanse på innholdstilpassede containere.
Eksplisitt min/max valideres uten å endre native widgetminimum. `scale-x` og
`x=...` er konkurrerende størrelsesvalg og kan ikke angis samtidig; tilsvarende y.
To regler for samme egenskap avvises, heller enn at siste tekstlige regel vinner.

`fill` deler restplassen som 1fr på hovedaksen, og fyller tildelingen på tverraksen.
Fr er en **søskenvekt**, scale er en **andel av mor**. Relative scale-barn og
innholdsbarn reserveres først, deretter fordeles rest etter fr med min/max:
`size = clamp(lambda * weight, min, max)`. Dette er ikke gammel minimum+grow.
Ved max-metning går restplass til justify. Ingen skjult flex-shrink innføres.
`items` gir bare defaults når barnet ikke har eksplisitt tverraksestørrelse/justering.

## 4. Måling, skalering og responsivitet

Layoutdimensjoner bruker relative mål/størrelsespolicy. Eierpresisering
2026-09-21: `font=10` angir absolutt tekststørrelse, ikke en relativ fontskala.
Vindusresize endrer tilgjengelig layoutområde, plassering og eventuell
tekstombryting, men skalerer ikke tekst eller øvrig innhold samlet.
`scale-x`, `scale-y` og `scale` bestemmer utstrekningen til layoutområdet;
de er ikke visuelle transformasjoner av innholdet. Widgets kan få større eller
mindre tildelte rektangler gjennom layoutreglene uten at etikettfonten skaleres.
Renderer og native toolkit får konkrete rektangler og tekstmål fra verten.
Visningszoom og DPI er egne vertsegenskaper, uavhengige av vindusresize.

Før native måling implementeres må profilen angi den felles absolutte fontenheten
og konverteringen til Fyne/Markdown/SVG. Absolutt størrelse er avklart; valget
mellom eksempelvis punkter og logiske skjermenheter er ennå ikke spesifisert.
Ingen backend får velge en privat enhet. Markdownstiler kan bruke relative
overskriftsstørrelser fra den absolutte basisfonten uten avhengighet av vinduet.

Eierens presisering: ytterframen med `{16:9,<->}` skal fylle rollen til
FixedAspectViewport, med bredde bestemt av verten og ratioavledet høyde.
Det er breddefylling, ikke contain. Concept1s samlede innholdsskalering ved resize
videreføres ikke. Behold samordnet geometri for visning, native kontroller og input.

`wrap` bryter før neste widget når den ikke passer, med kilde-/tabrekkefølge
bevart. Relative scale-/innholdsstørrelser brukes til radvalg; fr/fill på
hovedaksen sammen med wrap avvises foreløpig for å unngå sirkulær fordeling.
En for stor widget får overflowdiagnose/policy, ikke automatisk usynlighet.

`clip` klipper også inputområdet. `scroll` krever endelig viewport og
vertskapabilitet. Første leveranse støtter scroll-y; x-scroll får eksplisitt
unsupported-diagnose inntil vert/eksport er testet. Scrollbar og målpass er
budsjettert. Sideforhold bevares; scroll endrer innholdsområde, ikke frameformen.

Markdownhøyde måles ved avtalt bredde. Native widgets leverer minimum/preferred
og eventuell baseline ved den valgte absolutte fontstørrelsen. Bredde/høyde-
avhengigheter med innholdstilpasset mor må oppdages; ingen uendelig målfeedback.
Baselinejustering utsettes foreløpig.

`visible=false` fjerner noden fra layout/tabrekkefølge, men beholder instansen.
`enabled=false` beholder geometri og blokkerer handling. Skjuling må ha avtalt
fokus-/draft-policy. Programmatisk oppdatering er ikke en brukerhendelse.

## 5. Markdown og grenser

Vanlige strenger bruker eksplisitt escaping. Forslag: triple doble anførselstegn
for multiline Markdown; innholdet er rått, uten dedent/interpolasjon, og slutter
ved neste triple delimiter. En bokstavelig triple delimiter må foreløpig skrives
i en vanlig escaped streng. Normaliser aldri Markdown-innrykk ved parsing.

Quoted tekst som listeelement eller header/footer-innhold er Markdown;
widgetetikett som button("OK") er fortsatt en vanlig egenskapsstreng.
Mermaid-fences inne i Markdown behandles av
innholdstjenesten, ikke SDUI-lexer. Tjenesten får kildebase, bredde, tema,
kansellering og budsjett; returnerer mål, diagnose, kildekart og visningsressurs.

Første leveranse skal gjenbruke XFMDs støttede Markdown/diagramprofiler. Eksterne
ressurser følger samme vertsregler som dokumentet. Nestede sdui-fences i en
Markdown-widget avvises foreløpig med lokal diagnose, slik at ingen utilsiktet
rekursiv UI-instans opprettes. Dette er en eksplisitt første profilgrense.
Print bruker frosset akseptert state, ingen callbacks og ingen usendte drafts.
`scroll` eksporteres som klippet viewport med tydelig overflowindikasjon i første
profil; full utvidelse/paginering av alle scrollområder er en senere funksjon.

## 6. Representative fixtures

Eksemplene under kan parses/valideres i 0.2-frontenden; geometriutførelsen gjenstår.

```text
sdui 0.2;
panel = [
  header="## Arbeidsflate",
  top = [
    length = [<"## Lengde\n12,4 m">]*b {x=1fr},
    diameter = [<"## Diameter\n32 cm">]*b {x=1fr}
  ] {y=15fr, items=stretch};
  middle = [
    selection = [<"Valg">]*b {x=25fr},
    suggestions = [<"Forslag">]*b {x=50fr},
    stem = [<"Stamme">]*b {x=25fr}
  ] {y=45fr, items=stretch};
  track = [<"Stammeprofil">]*b {y=40fr}
]*box {16:9, <->, items=stretch, gap=0.01};
```

Dette låner Concept1s seks boksidentiteter og 15/45/40, 25/50/25-fordeling.
Header/dekorasjon, minimum og innhold gjør at det ikke automatisk er identisk
med CSS Grid. Geometri-fixturen må oppgi viewport og hva som sammenlignes.

```text
sdui 0.2;
form = [
  fields = <"## Kontekst"; context = input("Navn", value="C1")>
    {x=fill};
  actions = <cancel = button("Avbryt"), ok = button("OK")>
    {v<, >-<, >|<, gap=0.01}
] {scale=1, justify=between, padding=0.01};
```

Negative prøver: doble ID-er; ukjent variant; layout uten mål; min > max;
negativ gap; x=fill sammen med `>-<`; uavsluttet Markdown-streng;
wrap kombinert med fr; x-scroll uten kapabilitet; callback på anonym widget;
native minimum som ikke får plass; rå layouttokens inne i streng bevart som tekst;
px/bare absolutt layoutdimensjon (font er et eksplisitt unntak); løs formateringsblokk
etter separator; ugyldig ratio; ratio med begge scaleakser; uløselig
ancestor-referanse; ny ekte forelder kontra syntetisk rad som ikke endrer referansen.

## 7. Bevisst utsatt

Fri koordinatposisjonering, overlapp/z-order, calc-uttrykk, media queries,
tabell-grid med spans, endring av rekkefølge uten kildeendring, dock/splittere,
layoutanimasjon, dynamisk pluginlasting og virtuelle lister. Først utvides
språket når en konkret prototype viser et behov som nesting ikke løser godt.
