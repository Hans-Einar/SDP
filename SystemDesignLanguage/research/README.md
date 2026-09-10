# P04 — Syntese og videre undersøkelse

**Dato:** 2026-09-10. **Status:** research og anbefalinger; ingen språkbeslutning.
**Oppdrag:** [#10/P04](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5618871491).
**Underlag:** [32 profiler og metode](existingDesignLanguages/README.md).

## Eierens intensjon

SDP skal gjøre Features, deres ansvar og pathways synlige gjennom endringer. Architect/Designer skal kunne foreslå et avgrenset design-delta før Workers endrer kode. Eier og Verifier skal kunne se hva som er godkjent, hva som er endret, og hvilke påstander som faktisk er verifisert. Språkvalg, historieformat og de tre foreslåtte Feature-klassene er åpne hypoteser.

Erfaringene med SharedUI-drift og store refactorløp er her **eieropplysninger fra P04**, ikke nye verifiserte funn om GrassPhenology, TerrainAnalyzer eller HSX. Ponsse UI-beskrivelsen er et mulig case fra [design discussion #72](https://github.com/Hans-Einar/ponsse/issues/72), ikke en kontrollert beskrivelse av dagens implementasjon.

## Hva kartleggingen viser

Det finnes allerede byggeklosser for mye av ideen. Et nytt tekstformat er ikke nødvendigvis den største oppgaven; det vanskelige er presis semantikk, identitet, kontroll av kode og en forståelig endringsprosess.

| Behov | Observerte mekanismer | Hva SDP fortsatt må tilføre |
| --- | --- | --- |
| Flere views fra én modell | Structurizr, LikeC4; modelleringsstandarder som UML/SysML | Hvilket utsnitt en FEAT trenger, hvilken baseline det kommer fra og hva viewet utelater |
| Ansvar og begrepsgrenser | Context Mapper contexts/aggregates og refactoreringer | Produktets stabile FEAT-identitet og verifisering av faktisk ansvar etter endring |
| Varianter av et produkt | UVL og Clafer | Skille konfigurasjons-feature fra varig kapabilitet og tidsrevisjon |
| Eksplisitte kontrakter | OpenAPI, AsyncAPI, Protobuf, Smithy | Consumer-impact, semantiske enheter, leverings-/feilregler og releasekobling |
| Maskinelle fences | OCL, CUE, Rego og formell analyse | Pålitelig observerte kodefakta og prosessen for godkjente unntak |
| Modelltransformasjoner | ATL, QVT, Epsilon | Bevaringskrav, identitetsmapping og kontroll av informasjonstap |
| Historikk/provenance | Edapt change-modeller; PROV-O revision/derivation | Godkjenningsstatus, ansvarsmigrasjon og evidens om produktet |
| Kravutveksling | ReqIF | Stakeholders/UseCases→REQ→FEAT→verifikasjon→release |

Dette er en funksjonsoversikt, ikke rangering. Begrunnelse og primærkilder finnes i de respektive profilene.

Tre funn fortjener særlig oppmerksomhet:

1. [Context Mapper](existingDesignLanguages/context-mapper-cml.md) har allerede refactoreringer som deler bounded contexts etter use cases/user stories, samt merge-operasjoner. Det er et nært forbildet for å studere ansvar gjennom endring. Korrekt CML etter operasjonen er likevel ikke bevis på bevart programatferd.
2. [Edapt](existingDesignLanguages/edapt.md) beskriver eksplisitt operasjonshistorikk og migrasjon av modellinstanser. [Epsilon](existingDesignLanguages/epsilon.md) tilfører matching, transformasjon og validering. Dette utfordrer antakelsen om at all transition-infrastruktur må oppfinnes fra grunnen.
3. [Structurizr](existingDesignLanguages/structurizr-dsl.md) og [LikeC4](existingDesignLanguages/likec4.md) viser at tekstlig modell→flere viewpoints er praktisk tilgjengelig. Ingen dokumentert mekanisme her alene beviser komplette FEAT-pathways eller release-readiness.

## Fire forskjellige former for endring

| Form | Eksempel | Må ikke forveksles med |
| --- | --- | --- |
| Runtime-transition | En handling endrer UI-/domenetilstand | Endring av arkitektur over commits |
| Produktvariant | WebUI eller DesktopUI velges | Før/etter-revisjon av samme produkt |
| Modell-/schema-evolusjon | En type eller ansvarsgrense deles | Bevis på at programmet migreres korrekt |
| Levert produktendring | En bruker får endret funksjonalitet i en bestemt build | At et issue eller en PR er lukket |

SCXML history og temporal støtte i Alloy/TLA+ løser andre problemer enn Git-historikk. Protobufs feltidentitet beskytter deler av kontraktevolusjonen, mens PROV-O kan uttrykke avledning. Vi bør kombinere disse ideene ut fra behov fremfor å kalle alt «history».

En **anbefalt undersøkelseshypotese** er stabil designobjekt-ID + eksplisitte modellrevisjoner + separat transition-beskrivelse. Git kan lagre materialet; semantikken må si om ansvar er bevart, flyttet, splittet, slått sammen, erstattet eller ikke kartlagt. Ikke all historie trenger å ligge i hvert enkelt featureobjekt.

Godkjent baseline, foreslått design og observert implementasjon bør være separate tilstander. Et generert blueprint bør kunne identifisere modellrevisjon, generatorversjon og utsnittsregler. Et grønt blueprint er ikke tilstrekkelig dersom kodeinventaret mangler en dependency.

## Utfordring av Feature-/Refactor-taksonomien

**Anbefaling til neste sammenligning:** prøv kategoriene som uavhengige dimensjoner før de gjøres til en typehierarki.

- **Kapabilitet/Feature:** hva systemet skal kunne gjøre og for hvem, med REQ og akseptkriterier.
- **Designobjekt:** en ansvarsenhet, et lag, en kontrakt, en implementeringsbinding eller en beslutning som realiserer kapabiliteter.
- **Endringsobjekt:** hvorfor og hvordan noe endres, hvilke objekter som berøres og hvilke egenskaper som skal bevares.
- **Klassifikasjon:** domain, design og architecture kan tagge ansvar/endringsomfang; en endring kan ha flere av dem.

«Distribuerbarhet» kan være et kvalitetskrav eller en kapabilitet, mens socket-valget kan være en arkitekturbeslutning med kontraktskonsekvenser. En ny renderer kan realisere et nytt brukerbehov eller bare endre implementasjonen. Det avgjøres av intensjon og observerbar atferd, ikke rollenavnet.

Refactor bør ha eksplisitte bevaringskrav. Hvis brukeratferd endres samtidig, må den funksjonelle endringen være synlig. Architecture→design→domain kan være en planleggingsretning, men gjennomføring bør følge dependencies og migrasjonsrisiko; en liten vertikal overgang kan gi bedre læring enn tre store sekvensielle omskrivinger.

«Domain» er et godt fagbegrep, men det bør ikke bindes én-til-én til prosess. Skill begrepsområdet, bounded context, logisk komponent/container og deployert prosessinstans. Namespace/package og source-folder er implementeringsbindinger til slike grenser. Et lag bør også kunne inneholde flere tydelige ansvarsenheter; lagdeling alene hindrer ikke en stor, uoversiktlig fil.

Interfaces og komposisjon kan være designregler uten å kreve klasser eller arv. [Smithy mixins](existingDesignLanguages/smithy.md) illustrerer gjenbruk i en modell, men er ikke et argument for runtime-mixins. Modellen bør angi tillatte avhengigheter og kontrakter fremfor å gjøre ett programmeringsparadigme universelt.

## Contracts, gjenbruk og release-spor

En senere modell må kunne skille kontraktsform fra kontraktsmening: type/schema, enhet, command/event-retning, ordering, levering, feil, timeout, idempotens og kompatibilitetsregler. For UI-caset må også eierskap/livstid for representations, flere presentasjoner av samme instans, abonnement/avmelding og ubehandlede events avklares. Hint/lokalisering er ikke det samme som domenets identitet. En test-renderer kan verifisere meldingsløp, men erstatter ikke alle tester av en faktisk web- eller desktop-adapter.

Bibliotekbindinger bør kunne angi package, tillatt versjon/kontrakt, implementeringsspråk, source-område og beslutning om reuse/adapt/build-new. Et modellfelt alene hindrer ikke SharedUI-drift: et uavhengig kodeinventar må kontrollere faktiske imports og eventuell duplisering. Go er en eierpreferanse, ikke et generelt SDP-krav. Gjenbruk bør begrunnes i reelle behov og kontrakter fremfor spekulative universalabstraksjoner.

Traceability må nå frem til det som faktisk er levert. En **anbefalt hypotese** for senere gh-sdp-arbeid er et release-manifest med stabil FEAT-ID, levert revisjon/variant, brukerrettet endringsbeskrivelse, kontraktsversjoner, verifikasjonsreferanser, commit og identifisert build-artefakt. PR/issue er støttelenker, ikke release-notes-authority.

Et slikt manifest må håndtere delvis levert Feature, superseded arbeid, backports, features som går over flere releases og implementert kode som ennå ikke er eksponert for brukeren. CodeReview og interne refactoroppgaver trenger ikke automatisk en brukerrettet release-note, men må fortsatt kunne spores til endringen de kontrollerer. Ingen gh-sdp-funksjonalitet implementeres i P04.

## Vesentlige åpne spørsmål

1. Hva er minste autoritative designobjektsett, og hvilke identiteter skal overleve rename/split/merge?
2. Hva er det presise kriteriet for bevart Feature-atferd, inkludert kvalitetskrav?
3. Er pathway en tillatt strukturell rute, et runtime-scenario eller begge med forskjellige relasjonstyper?
4. Hvordan registreres overgangstilstander, kompatibilitetsvinduer, rollback og ufullstendig mapping?
5. Hvordan hentes observerte kodefakta uavhengig av modellforfatter/Worker, og hvordan vises usikkerhet?
6. Hvilke kontroller er syntaktiske, semantiske, statiske kodekontroller, runtime-tester eller faglig review?
7. Hvordan holdes små prosjekter lette, og hva er maksimal kostnad for modellvedlikehold og genererte viewpoints?
8. Hvem kan endre authority, og når må en Worker stoppe fremfor å «fikse» modellen for å få grønt resultat?

Modellen bør være en begrenset beskrivelse av ansvar og kontrakter, ikke en kopi av hver klasse og kodelinje. Arkitekt-/designer-skills, Worker-skills og uavhengig verifikasjon trenger samme identifiserte oppdragsgrunnlag. Rework-budsjett, stoppkriterier og eskalering må ligge i arbeidsprosessen; en DSL kan bære disse opplysningene, men kan ikke alene stoppe agentløkker.

## Neste bounded assignment — forslag, ikke startet

**Formål:** sammenligne representasjon av én liten endring; ikke velge/fryse språk.

**Input:** en separat bekreftet beskrivelse av dagens MVP1-UI og ett foreslått delta fra #72. Bruk én output-kapabilitet og én input-action, én kontrakt og maksimalt seks ansvarsenheter. Dersom dagens design ikke er verifisert, bruk et tydelig hypotetisk case fremfor å erklære det som baseline.

**Tre kontrasterende spor:** (A) en eksisterende arkitektur-DSL som Structurizr eller LikeC4, (B) CML eller SysML som semantisk modell, (C) en liten eksplisitt JSON/CUE-modell med refererte kontrakter. Før utføring velges ett konkret verktøy i A/B og versjon/lisens avklares. Dette er forsøksspor, ikke shortlist med vedtatt vinner.

**Leveranse:** samme baseline og foreslåtte endring i tre spor; FEAT→REQ→ansvar→contract-view, før/etter-ansvarsmapping og et Worker-/Verifier-utsnitt. Registrer hvilke begreper som er native, metadata-konvensjon eller krever egen kode. Test rename og ett split/move, en brutt referanse og en forbudt dependency. Vis også hva verktøyet ikke oppdager. Sammenlign redigeringsmengde, informasjonstap, eierlesbarhet og integrasjonskostnad.

**Grenser:** ingen produktkode, ingen full Concept1→MVP1-mapping, ingen generell compiler eller skill. Maksimalt én reparasjonsrunde per spor; ved fortsatt verktøyproblem registreres hindringen i stedet for en åpen Worker/Verifier-løkke. Stopp etter sammenligning og owner-review av resultatet.

Concept1→MVP1 kan senere brukes som stresstest med eksplisitt delvis ekvivalens. At modellen virker for MVP1 er ikke bevis på universell anvendelighet; et lite prosjekt og et annerledes kontrakts-/deploymentmønster bør etter hvert prøves.
