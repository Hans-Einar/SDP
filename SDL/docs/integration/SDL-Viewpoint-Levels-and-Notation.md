# SDL — viewpoints etter abstraksjonsnivå og fast notasjon

**Status 2026-09-22:** designpresisering G6-D2. Eksportvalg, menystruktur og
notasjon nedenfor er planlagt; dagens generator lager fortsatt hovedsakelig
merkede flowcharts. Ingen nye SDL-typer eller UML-relasjoner vedtas av dette
dokumentet. [G6-planen](../../../SDUI/docs/implementation-plan.md) og
[navigasjonsdesignet](SDL-Navigable-Viewpoints-Design.md) hører sammen med dette.

## To likeverdige eksportformer

**Navigator/overview** er anbefalt utviklingsvisning. Generer bare navigasjon,
små faseoversikter og navne-/typekataloger. List UseCase, Feature, Functionality,
Capability, Activity, Mode, Container, Unit, Interface og Channel fra den
validerte modellen, med antall og stabile identiteter. Detaljdiagrammer,
bidragskart, tilbydere, avhengigheter og valgte scenarioer bygges først ved klikk.
Ingen detalj-SVG eller skjult full rendering kreves for å lage navigatoren.

**Statisk pakke** materialiserer et eksplisitt utvalg eller alle viewpoints for
deling, arkivering og lesere uten lenkehandler. Samlet rapport er et ekstra valg.
Begge bruker samme prosjektor og samme utvalgsbeskrivelse. Dette er alternative
leveranseformer, ikke et krav om først å generere hundrevis av diagrammer.

Navigatoren kan være dynamisk eller eksportert med relative fallback-lenker.
En ennå ikke generert side må ha en handlingslenke, ikke en død fil-lenke.
En generisk leser får forklarende oversikt og tilbud om statisk eksport.
Klikket angir viewpoint, fokusobjekt, relasjonstyper, retning, dybde, mode og
notasjonsprofil. Serveren avgrenser utvalg og arbeid; dette er ikke fri kode.
Kilderevisjon og utvalgsnøkkel følger både navigator og generert side.

## Abstraksjonslag og dokumentstruktur

Checkpoint [02, avsnitt 2 og 6](../../../SDP/History/checkpoint-1/02-SDL-Model-and-Abstraction-Levels.md)
definerer A0–A5 som organiserende nivåer, med kandidatprofilene Intent,
Architecture, Detailed design og Execution binding. Disse er ikke parser-headere.
Vi beholder denne inndelingen framfor å innføre konkurrerende nivåer:

| Dokumentområde | Nivå | Oversikt og relevante visninger |
| --- | --- | --- |
| Krav og hensikt | A0 behov/forpliktelser, A1 funksjonell hensikt | Actors, bruksmål, Features, Functionality, Capability, overordnet Activity, krav og sporbarhet |
| Arkitektur | A2 system/container, A3 interne Units/lag | Container-/Unit-kart, Channels, porter, tilbydere, bidrag og allokering |
| Detaljdesign | A4 atferd/kontrakter | Scenariosekvenser, Activity-detaljer, data-/meldingskontrakter, packet; senere klassestruktur og state-maskiner |
| Realisering og bevis | A5 | Bindinger til kode, generering, verifikasjon, utførte spor og deployment |

Foreslått minste navigatorpakke:

```text
viewpoints/
  index.md
  navigator.md
  requirements/
    index.md
    inventory.md
  architecture/
    index.md
    inventory.md
  design/
    index.md
    inventory.md
  realization/
    index.md
    inventory.md
```

Katalogene er navigasjonsplasseringer, ikke kopier av SDL-objekter. Samme
Functionality eller kontrakt kan være relevant på flere nivåer. Modell-ID og
kilde er den samme. Katalogene bygges fra et eksplisitt viewpoint-register
med nivåer, spørsmål, nødvendige typer/relasjoner, tillatte utvalg og notasjonsprofil.
Objekttypen alene bestemmer ikke abstraksjonsnivået; en Activity kan være A1 eller A4.
Manglende nivåklassifisering vises som uspesifisert, ikke gjettet fra objektnavn.

Foreløpig registerkart: VP01 primært A0/A1; VP02 A2/A3; VP03 og VP07 kobler A1
til A2/A3; VP04/VP05 A2/A3; VP08 A1 eller A4 etter eksplisitt scenarioomfang;
VP09 A1/A4; VP10 A4. VP11 er kilde-/sporbarhetsvisning på tvers av nivåene.
VP06s arbeidsplan G1–G6 hører til utviklingsarbeidet og får eget menypunkt;
den skal ikke forveksles med systemets abstraksjonslag eller runtime-Activities.
G-faser beskriver vår implementasjon, A-nivåer beskriver systemdesignet.

Requirement, System, Class og State er foreløpig ikke egne deklarasjonstyper
i design-core 0.5. Oversikten må skille **ikke støttet av profilen**, **ikke
modellert** og **ingen treff i utvalget**. En tom liste betyr ikke at systemet
mangler krav eller tilstander. Planlagt eller dokumentert innhold merkes tydelig.

## Mode og State

**Mode** er i aktiv SDL en navngitt drifts-/anvendelseskontekst som avgrenser
avhengigheter, allokering og Channel-deltakelse. UiPreview, BoundExecution og
DocumentBrowsing er eksempler. Aktiviseringsregler, gjensidig utelukkelse og
mode-hierarki er ikke implementert.

**State** er en situasjon eller tilstand hos et bestemt system-/aktivitetsobjekt,
for eksempel Ready, Processing eller Suspended. Eieren og livstiden må være
definert. Overganger krever hendelser/betingelser og virkning; navn alene er
ingen state-maskin. Checkpointets State/Activity-semantikk er fortsatt delvis
kandidatgrunnlag, ikke parserstøtte.

En Activity kan gå fra Ready til Processing og tilbake mens samme Mode gjelder.
Samme State kan være relevant i flere Modes. En modell kan senere gjøre en Mode
avhengig av et state-predikat, eller gi en Mode en eksplisitt under-state-maskin.
Ingen slik relasjon utledes automatisk. **State er altså ikke generelt et barn
av Mode, og ordene er ikke synonymer.** Navigatoren viser dem separat når
språkprofilen kan uttrykke begge.

## Fast visuelt vokabular

Bruk UML der semantikken passer og et eksplisitt SDL-symbolsett for SDL-spesifikke
begreper. Tabellen er en foreslått visningsprofil, ikke nye språkregler.
Alle symboler har typeetikett og navn; farge alene skal aldri skille typer.
Detaljnivå endrer innhold/compartments, ikke objektets gjenkjennelige typeikon.

| Begrep | Fast visuelt kjennetegn |
| --- | --- |
| Actor | UML-aktørfigur med navn; ekstern systemaktør kan merkes som sådan |
| UseCase | Ellipse med bruksmålets navn |
| Requirement | Dokumentform med ID og kort forpliktelse; planlagt type |
| Feature | Rektangel med fane og tydelig Feature-merke |
| Functionality | Avrundet ansvarskort med vertikal sidemarkør og Functionality-merke |
| Capability | Kort med fast evne-ikon og Capability-merke, forskjellig fra use-case-ellipsen |
| Activity | Avrundet aktivitetsform med Activity-merke; detaljvisning kan åpne flyt |
| Mode | Kontekstfane/ramme med Mode-merke; ingen start-/sluttmarkør |
| State | UML-state-form inne i navngitt tilstandsmaskin; planlagt type |
| System / Container / Unit | Nestbare grenser med respektive faste typeikoner; Container er ikke automatisk en UML deployment-node |
| Interface | Merket kontraktport; UML lollipop/socket bare når tilbud/krav er uttrykkelig modellert |
| Channel | Navngitt forbindelseslinje/knutepunkt med endepunkter, roller og kontrakt |
| Class | UML-klasse med compartments for egenskaper/operasjoner; bare ved eksplisitt klassemodell |

Bruksmålvisningen skal få aktører og ellipser. Feature-/Functionality-sporbarhet
kan ligge i et separat tilknyttet diagram; alt behøver ikke presses inn i en ren
UML use-case-visning. Ingen systemgrense eller include/extend-relasjon oppdiktes
fra plassering eller supports-fakta.

## Piler: UML der betydningen stemmer

[OMG UML 2.5.1](https://www.omg.org/spec/UML/2.5.1), avsnitt 7.7 om
avhengigheter/realisering og 9.5 om Property/aggregation, er semantisk referanse.
[Mermaids klassediagramprofil](https://mermaid.js.org/syntax/classDiagram.html)
har notasjon for dependency, realization, aggregation og composition; renderer-
syntaks alene er ikke bevis for at en SDL-relasjon har den samme betydningen.

| SDL-faktum / senere relasjon | Visningsregel |
| --- | --- |
| Unit consumes Interface | UML-inspirert stiplet dependency med åpen pil mot grensesnittet; behold teksten consumes. Ingen faktisk melding eller obligatorisk avhengighet utledes. |
| Participant uses Channel as sender/receiver | Vis Channel-deltakelse og rolle. Dette er ikke UML Usage bare fordi ordet uses forekommer. |
| Functionality realizes Capability | Dagens betydning er **bidrar til realisering**, ikke full oppfyllelse. Behold merket SDL-pil realizes; ikke gi den UMLs hule realiseringstrekant uten en sterkere eksplisitt kontrakt. |
| Unit provides Capability | Merket tilbyr-relasjon; Capability omdøpes ikke til Interface eller lollipop. |
| supports / contributes-to / addresses | Egne merkede SDL-sporbarhetsrelasjoner; ingen implisitt arv, include eller extend. |
| contains / owns | Logisk inndeling/ansvar; ingen UML composition-diamant uten modellert helhet/del- og livstidskontrakt. |
| Senere eksplisitt UML realization | Stiplet linje og hul trekant mot spesifikasjonen som realiseres. |
| Senere aggregation / composition | Hul/fylt diamant ved helheten; retning, roller, multiplisitet og eierskapsbetydning må være eksplisitt. |

Etikettene uses/consumes/realizes skal beholdes også når pilen har en kjent form.
En legende angir hva som er SDL og hva som er UML. UML-inspirert er ikke det samme
som en formelt definert UML-profil. Først bestemmes forholdets betydning, deretter
velges pil; ikke omvendt.

## Rendererprofil og senere klassediagrammer

Mermaids offisielle dokumentasjon beskriver nå
[usecase-beta](https://mermaid.js.org/syntax/usecase.html), med aktørfigurer og
ellipser. Det beviser ikke støtte i den lokale Rust-rendereren. Lokal prøve med
den eksisterende mmdr-binæren ga exit 0, men tegnet syntaksord som noder i stedet
for et use-case-diagram. Den er ikke godkjent for denne notasjonen.
En framtidig kapabilitetsprøve må kontrollere AST/figurer og pilmarkører, ikke bare
exitkode eller at en SVG-fil finnes. Dette oppdraget endrer ikke renderer-koden.

G6-M5 skal etablere en versjonert symbol-/pilprofil, riktig use-case-visning og
testede backend-kapabiliteter. Manglende støtte gir en tydelig merket alternativ
visning eller en diagnose; en flowchart skal ikke kalles et UML use-case-diagram.
Full-/enkeltvisning skal bruke samme symboler og projeksjon.

G6-M6 planlegger klassediagrammer senere. Før rendering må SDL uttrykke klassene,
attributter/operasjoner, assosiasjonsender/roller, multiplisitet og forskjellen på
aggregation og composition. Unit-contains blir ikke automatisk klassestruktur.
Språkprofil, validering, negative tester og kildekart må leveres samlet. State-
maskiner og include/extend krever tilsvarende eksplisitte framtidige profiler;
de legges ikke til automatisk gjennom et notasjonsvalg.

Akseptanse: navigator-only gir ingen detaljdiagrammer; utvalg matcher tilsvarende
del av statisk eksport; alle viste symboler beholder modelltype/kilde-ID; relasjoner
har korrekt retning og markør; sort/hvitt og nestede grenser er lesbare; unsupported
diagramtype kan ikke passere bare fordi renderer returnerte exit 0.

SDL-modellen skiller Feature-en NavigableDesignDocumentation (navigasjon og
dokumentlevering) fra TypedDesignInspection (nivåkataloger, relasjonsutvalg og
semantisk notasjon). Begge støtter BrowseDesignViews; sistnevnte støtter også
InspectModels. Dette gir egne kildekoblede bidragsvisninger for to ulike behov.

Verifikasjonsgrense: første samlede VP07-prøve med alle G6-ansvar i én Feature
traff generatorens 60-sekunders renderergrense. Den endelige modellen skiller de
to behovene ovenfor; en vellykket eksport av disse mindre visningene opphever
ikke renderergrensen. G6-M2 skal også ha utvalgs-/diagramgrenser og diagnoser
for tunge utsnitt. Ingen generator-timeout eller ekstern renderer er endret her.
