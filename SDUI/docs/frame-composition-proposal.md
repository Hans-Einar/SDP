# SDUI — frames med header, body, footer og gjenbruk

**Dato:** 2026-09-21 · **Status:** eierinnspill med konkretiserende designforslag.
Dette supplerer [layoutforslaget](layout-language-proposal.md) og bevarer
begrunnelsen fra 2026-09-21. **Statusavklaring 2026-09-24:** Go-frontend,
normalisering, layout og runtime er levert. [Språkprofil](language.md),
[målekontrakt](go-layout-contract.md) og [runtime](../go/runtime/README.md) gjelder
for kjøring. Anbefalingene nedenfor er ikke en ekstra aktiv implementasjonsprofil.

## 1. Fast forholdsflate bestemt av bredden

```text
page = [ <"Innhold"> ] {16:9, <->};
```

Rotframe bruker hele bredden W i vertens layoutområde. Høyden er W × 9/16.
Tilgjengelig vindushøyde endrer ikke denne bredden. For liten høyde håndteres
av avtalt overflow/scroll, ikke ved å bytte til contain eller endre forholdet.
Dette er eierens SDUI-uttrykk for rollen til FixedAspectViewport.

`<->` bestemmer x-aksen, ratio bestemmer y. Tilsvarende kan `^|v` bestemme y
og la ratio bestemme x. Begge strekkoperatorene med ratio er en konflikt.
En underframe bruker sitt tildelte tilgjengelige område i nærmeste mor.
`scale-x=0.75` velger en eksplisitt andel av ancestorbredden; `<->` fyller sitt
tildelte område. På rot er `<->` det samme breddemålet som scale-x=1.

Rollen som fast forholdsflate er dermed avklart. Resize endrer framegeometrien,
men skalerer ikke innholdet. Tekst har absolutt fontstørrelse. Concept1s samlede
skalering av en fast designflate skal ikke videreføres som SDUI-resizeatferd.

## 2. Regioner som egne innholdsverdier

```text
page = [
  header = [ ],
  <"Body">,
  footer = [ ]
]*b {16:9, <->};
```

Anbefaling: `header`, `body` og `footer` er strukturelle regionroller i en frame.
De er ikke en alternativ måte å gi en vilkårlig widget ID-en header/footer på.
En region kan inneholde en frame, widgetgruppe, Markdown-streng eller referanse
til en slik definisjon. Layoutsuffix på en streng gjelder Markdown-widgeten.
Rolleformen gjelder også en udekorert frame; `*b` velger BoxUI-presentasjonen.

Anbefalte regler:

- Høyst én header, body og footer per frame. Manglende region tar ingen plass.
- Vanlig umerket innhold blir body. Eksplisitt `body=...` er valgfritt, men skal
  ikke blandes med ytterligere implisitt body i samme frame.
- Header/body/footer ligger vertikalt uavhengig av rekkefølgen på rollefeltene.
  De er ikke tre vanlige søsken i samme kommarad. Inne i body gjelder fortsatt
  komma/semikolon-reglene for rader.
- Header/footer måles etter innhold ved framebredden. Body får resten etter
  dekorasjon, padding og avstand. Hele ytterframen, inkludert alle tre regionene,
  har 16:9; header/footer legges ikke utenpå ratiohøyden.
- For stor header/footer gir overflowdiagnose/policy. Ratio endres ikke.
- Regionene gir semantiske innholdsområder: body-barnas relative mål bruker det
  resterende bodyområdet. En syntetisk rad endrer fortsatt ikke målreferansen.
- En tom frame er en lovlig plassholder. Ingen skjult standardheader/footer
  opprettes bare fordi varianten er box.

Tidligere `heading='...'` erstattes av `header=...`; vi beholder ikke to aktive
overskriftsmekanismer. `header="## Overskrift"` bruker vanlig Markdown-rendering,
så headingnivå, flere linjer og annet innhold kan beskrives uten ny widgettype.

## 3. Navngitte grupper og instanser

### Nestede grupper med egen layout

Eierpresisering 2026-09-21: `<>` kan nestes i flere nivåer. Hver eksplisitte
gruppe er én komponent i den omsluttende listen og kan ha egne layoutregler.

```text
mainBody = <
  <"gruppe1", "med to UI-komponenter"> {v<, >-<, >|<},
  <"gruppe2"; <"undergruppe", "med to UI-komponenter"> {>-<}>
    {>^, >-<, >|<}
> {<->};
```

Det ytre kommaet plasserer gruppe1 og gruppe2 horisontalt etter hverandre.
Kommaet inne i gruppe1 plasserer dens to tekster horisontalt. Semikolonet inne
i gruppe2 legger undergruppen under teksten «gruppe2», uten å starte en ny rad
i den ytre gruppen. Separatorene virker bare i listen de tilhører.

Formateringsblokken etter `>` tilhører akkurat den avsluttede gruppen.
Gruppens justering plasserer den samlet i området den får av forelderen;
gruppens innvendige layoutregler styrer dens egne barn. Nestede grupper
bevares i AST og normalisert modell, også med ett barn: de må ikke flates bort
slik at egen formatering eller ancestor-referansen for barna går tapt.
Grupper kan være anonyme eller navngitte; navn kreves ikke for å gi dem layout.

### Gjenbruk av navngitte grupper

Eierens eksempel, med avsluttende semikolon som foreslått dokumentseparator:

```text
sdui 0.2;

mainBody = <"main body widgets", button1 = button("OK")>;

page = [
  header = "##Markdown Heading" {<-},
  mainBody,
  footer = <
    "footerText left center" {<-, font=10},
    "footerText right center" {->, font=12}
  >
]*b {16:9, <->};
```

`mainBody` er en gjenbrukbar definisjon; referansen i page instansierer innholdet.
Den er ikke en tekststreng, callback eller et delt levende FOX-objekt.
Man kan skrive `body=mainBody` for å gjøre regionen eksplisitt.

Anbefaling: lokale navn er unike per komponentdefinisjon. Full widgetidentitet
inkluderer instansbanen, for eksempel page/mainBody/button1, og runtimegeneration.
To anvendelser av samme komponent må ha forskjellige instansnavn, for eksempel
left=mainBody og right=mainBody. Da kolliderer ikke den interne button1.
Dette presiserer den tidligere regelen om ett flatt navnerom per UI-definisjon.
Uklar referanse, duplisert instansnavn og rekursive definisjonssykluser avvises.
G1s Go-profil støtter framoverreferanser og statisk instansekspansjon med full
setHandle-bane; se språkprofilen. Python-portgrunnlaget og gammel 0.1-parser
er fjernet. Dette er ikke automatisk oppkobling til en SDL-runtime.

Kun instansiering oppretter UI-livstid. En bibliotekdefinisjon skal ikke monteres
som en ekstra side automatisk. Verts-API-et velger hvilken frame-definisjon som
er inngang, her page. Ingen avhengighet av navnet page eller filrekkefølgen.

## 4. Footerjustering og font

`<-` og `->` betyr henholdsvis venstre/midten og høyre/midten. For at den høyre
teksten faktisk skal nå høyre kant, trenger gruppen restplass å justere innenfor.
Anbefalt footer-default er full bredde og justify=between for én rad; med to
innholdstilpassede tekster gir dette ønsket plassering uten ekstra kildeegenskaper.
Eksplisitt layout kan overstyre defaulten. Overlapp ved for liten bredde gir
avtalt overflow, ikke automatisk bytte til en annen uoppgitt layout.

Eierpresisering 2026-09-21: `font=10` og `font=12` er absolutte tekststørrelser.
De beholdes ved vindusresize. Font påvirker måling, header/footerhøyde og
innholdsrom; ny tilgjengelig bredde kan endre tekstombryting og dermed høyde.
Den felles native fontenheten må fastsettes før måleimplementasjonen; dette
gjenåpner ikke beslutningen om absolutt størrelse. Forslaget til arv er at en
gruppe/frame gir basisfont til barna, og barnets eksplisitte font overstyrer den.

Formatering er tillatt etter hver UI-komponent og hører til komponenten foran.
`{...}` kommer etter eventuell variant og før `,` eller `;`. Komma plasserer
neste komponent horisontalt etter den foregående; semikolon starter en ny rad
under foregående rad. Dette gjelder innholdslistene; regionroller følger
forslaget i del 2. Se [layoutforslaget](layout-language-proposal.md) for
eksempler med formatering på strenger, widgetkall, grupper og referanser.

Merk at vanlig Markdown normalt krever mellomrom i `## Markdown Heading`.
Eierens tekst `##Markdown Heading` er bevart ovenfor; SDUI skal ikke stille
reparere Markdown-kilden. Egenskapens innhold gis videre til innholdsprovideren.

## 5. Kontrolltilfeller for G1–G4

Kanoniske hjørner ^</>^/v</>v og ekvivalente alternative retningstavemåter;
ratio med bare x-strekk og bare y-strekk; avvis begge med ratio; behold bredde
ved for liten vindushøyde; header/footer innen ratio; eksplisitt og implisitt body
gir samme struktur hver for seg; konflikt mellom dem avvises; multiline header;
footerens høyre tekst når høyre kant; uendret font ved resize og korrekt ny
tekstombryting; formatering på alle komponenttyper før separator og avvisning
av løs formateringsblokk etter separator;
tre gruppenivåer med uavhengig layout; lokal virkning av komma/semikolon;
bevar eksplisitt gruppe med ett barn og dens ancestor-referanse;
to komponentinstanser med separat state; referansesyklus; ubundet callback;
ny kilde/instancegeneration og avvisning av events fra gammel instans.
