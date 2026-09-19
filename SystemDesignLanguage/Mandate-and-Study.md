# Arbeidsdokument: modellbasert og kontrollert programvareutvikling med agenter

**Dokument-ID:** SDL-MANDATE-STUDY-001  
**Revisjon:** 0.1, 2026-09-14  
**Status:** Renskrevet eierintensjon med innledende studie og eksplisitt merkede forslag. Ikke vedtatt språk, schema eller implementeringsplan.  
**Kontekst:** Hans-Einar/SDP, issue #5, #7 / PR #8, #9 og #10.  
**Målgruppe:** Eieren, arkitekter, designere og agenter som skal undersøke, videreutvikle eller bruke SDP.

## 0. Les dette først

Vi ønsker en standardisert utviklingsmodell som gjør det mulig å styre agentdrevet programvareutvikling gjennom en kompakt, maskinlesbar beskrivelse av systemets hensikt, struktur, ansvar, kontrakter og atferd. Verktøy skal kunne utlede oppgavetilpasset kontekst, endringsbeskrivelser og kontroller fra modellen. Store endringer skal skje gjennom eksplisitte revisjoner og faseoverganger, mens læring og videreutvikling fortsatt er mulig.

Det sentrale problemet er at lokale kodeendringer ofte mangler en forstått sammenheng fra brukerbehov til faktisk virkning gjennom hele systemet. Gjentatte lokale reparasjoner kan skjule manglende integrasjon, flytte ansvar og gradvis bryte arkitekturen. Modellen, arbeidsprosessen og kontrollene skal sammen gjøre slike avvik synlige og håndterbare.

Vi ønsker detalj der den styrer utviklingen, men ikke en manuelt vedlikeholdt kopi av hele implementasjonen. Maskinelle svar må oppgi hvilke egenskaper, kodeområder og kjørevarianter de faktisk dekker. Ukjent eller ukontrollert samsvar skal være synlig.

**Leserveiledning:** Del I renskriver mandatet. Del II er en innledende studie og forslag til hvordan det kan realiseres. Del III angir åpne spørsmål, mulige neste undersøkelser og regler for videre arbeid. En agent som får dette dokumentet, skal ikke behandle forslagene som allerede godkjente SDP-kontrakter.

Snarveier: [Mandat](#del-i--mandat) · [Faser og læring](#4-faser-revisjoner-og-læring) · [Foreslått endringsmodell](#11-forslag-til-fase--og-endringsmodell) · [Fra design til implementasjon](#13-fra-designbeskrivelse-til-implementasjon) · [Standarder](#15-standarder-rundt-designbeskrivelse-livsløp-og-implementasjon) · [Verifikasjonsgrenser](#17-hva-kan-verifiseres-programmatisk) · [Videre arbeid](#del-iii--videre-arbeid)

## 1. Proveniens, status og forholdet til eksisterende SDP

### 1.1 Hva som er grunnlag for dette dokumentet

Renskrivingen bygger på eierens innspill i samtalen 14. september 2026, inkludert den gjengitte Gemini-dialogen. Den bygger også videre på eksisterende SDP-arbeid; dette er ikke en ny start som opphever tidligere beslutninger.

| Spor | Kontrollert status 2026-09-14 | Betydning her |
| --- | --- | --- |
| [Issue #5](https://github.com/Hans-Einar/SDP/issues/5) / PR #6 | Studien er akseptert og merget; #5 er lukket. Akseptert merge er `2cb49c02145621b099c47d05786716598e414e75`. | Evidensgrunnlag og vedtatt overordnet retning videreføres. |
| [Issue #7](https://github.com/Hans-Einar/SDP/issues/7) / [PR #8](https://github.com/Hans-Einar/SDP/pull/8) | Åpent, foreløpig pilotarbeid. PR #8 er draft på `ea9fcf1cdd3198aeac515b89b55398282c463838`. Siste statuskommentar sier at rework mangler ny uavhengig verifikasjon/review. | Dette dokumentet erklærer ikke pilotkontrakten ferdig eller reviewfunn lukket. |
| [Issue #9](https://github.com/Hans-Einar/SDP/issues/9) | Åpent Steering- og koordineringsspor. | Samlet videre retning og samordning hører hjemme her. |
| [Issue #10](https://github.com/Hans-Einar/SDP/issues/10) / [PR #11](https://github.com/Hans-Einar/SDP/pull/11) | Modellinnspill og research. PR #11 er draft på `319ee2a43fe8a05bc3aac0860ff937bba57819ee`. | Eksisterende kartlegging av 32 språk/notasjoner er underlag, ikke språkvalg. |

Kilder til aksept og foreløpig status: [#5 Steering-aksept](https://github.com/Hans-Einar/SDP/issues/5#issuecomment-5368329242), [#7 WIP-status](https://github.com/Hans-Einar/SDP/issues/7#issuecomment-5607093722) og [#10 R04-leveranse](https://github.com/Hans-Einar/SDP/issues/10#issuecomment-5622626432). Dette er en statuskontroll av issue-/PR-informasjon, ikke en ny gjennomgang av #7-koden eller alle reviewfunnene. Issue #10s opprinnelige stopptekst må leses sammen med senere autorisasjon og R04; den beskriver ikke alene dagens leveransestatus.

Eksisterende faglig underlag:

- [Akseptert studieretning](../SDP-usage-analysis/ProposedSDPWorkflow.md).
- [Språkkatalog med 32 profiler](research/existingDesignLanguages/README.md).
- [P04-syntese, avgrensninger og sammenligningsforslag](research/README.md).

### 1.2 Hvordan utsagn skal tolkes

| Merking / plassering | Betydning |
| --- | --- |
| Mandat, Del I | Renskrevet eierintensjon. Kravene er mål for utviklingsarbeidet, ikke påstander om eksisterende verktøyfunksjonalitet. |
| Kildebelagt observasjon | Hva de oppgitte kildene faktisk beskriver, innenfor angitt undersøkelsesomfang. |
| Forslag / vurdering, Del II | Assistentens forslag til konkretisering. Krever evaluering før det blir normativ SDP. |
| Åpent spørsmål, Del III | Et spørsmål som ennå ikke er avgjort. |

Denne revisjonen dokumenterer oppdraget og utfører en innledende studie. Den endrer ikke canonical Toolkit eller godkjenner migrasjon av andre prosjekter. De nye Feature-/Functionality-/Channel-begrepene og fasereglene må samordnes med #7/#9 før de eventuelt blir canonical SDP.

# Del I — Mandat

## 2. Formål og problem som skal løses

### 2.1 Formål

Utvikle en gjenbrukbar metode og et tilhørende designspråk/verktøygrunnlag for programvareprosjekter der agenter skriver og vedlikeholder mye av koden. En ny agent skal raskt kunne forstå systemet på relevant granularitetsnivå og få presise rammer for en endring, uten å lese store mengder historiske dokumenter.

Metoden skal gi fleksibilitet til å endre arkitektur, design og implementasjon gjennom prosjektets levetid, samtidig som endringene blir eksplisitte, sporbare og kontrollerbare. Rigiditet gjelder særlig ansvar, gjeldende beslutningsgrunnlag, endringsmyndighet og kriterier for å gå videre.

### 2.2 Feilmønsteret vi vil motvirke

1. En agent endrer en lavtliggende del av systemet og antar at den ønskede funksjonen dermed virker fra ende til ende.
2. Nødvendige endringer i andre komponenter, kontrakter, konfigurasjon eller brukerflater overses.
3. Tester eller praktisk bruk avdekker mangler. Agenten reparerer lokalt uten å revurdere helheten.
4. Nye lokale løsninger flytter ansvar, dupliserer funksjonalitet eller bryter avtalte grenser.
5. Systemet kan til slutt se ut til å virke, men implementasjonen har avveket vesentlig fra tiltenkt arkitektur og design.
6. Senere endringer blir stadig vanskeligere, og eieren mister oversikten over hva som er bygget og hvorfor.

Dette er eierens beskrevne erfaring og motivasjon. Dokumentet hevder ikke å ha utført en ny empirisk analyse av alle berørte produktrepoer eller fastslått årsaken til konkrete produktfeil.

### 2.3 Sammenheng med tidligere SDP-retning

Ønsket om mindre «agile» tankegang forstås her som mindre ad hoc-prioritering og lokal improvisasjon uten tydelig systemgrunnlag. Det skal samtidig være mulig å lære, levere små vertikale endringer og revidere tidligere valg.

Dette viderefører #5s retning: Feature, Refactor eller reell Fix eier intensjonen; issue avgrenser et oppdrag; Slice er et sammenhengende, kontrollerbart resultat. Sprint er valgfri organisering. En fasebasert modell må kunne romme dette uten at en tidlig implementeringsplan blir styrende for all fremtidig utvikling.

## 3. Ønsket system- og designspråk

### 3.1 En maskinlesbar designmodell

Systembeskrivelsen skal kunne uttrykkes som faktisk, versjonerbar designkode med definert syntaks og semantikk. En interpreter eller modellkompilator skal kunne laste modellen, kontrollere den og utlede informasjon. Den trenger ikke kjøre produktets funksjonalitet.

Vi skal undersøke tre muligheter: bruke et eksisterende designspråk, utvide et eksisterende språk eller definere et eget språk. Språkvalg og implementasjonsteknologi er åpne. Knytning til etablerte standarder og metoder skal bidra til presise begreper og gjennomtenkte prosesser.

### 3.2 Struktur på flere nivåer

Modellen skal kunne beskrive:

- Systemkontekst, eksterne aktører og systemgrenser.
- Arkitekturkomponenter med ansvar, innhold og relasjoner.
- Designbestanddeler, her omtalt som *design constituents*, med det innholdet som er nødvendig for å forstå og styre implementasjonen.
- Binding til relevant source tree, kodeområder og faktiske implementasjoner.

Detaljeringsnivået skal kunne variere mellom delsystemer og utviklingsfaser. Det skal ikke automatisk kreves et manuelt designobjekt for hver klasse, metode eller kodelinje.

### 3.3 Førsteklasses designobjekter på tvers av strukturen

Eieren ønsker særlig å undersøke følgende objekter:

| Objekt | Renskrevet hensikt | Avklaring som fortsatt trengs |
| --- | --- | --- |
| Feature | Gjøre en varig kapabilitet og dens sammenheng gjennom systemet synlig. | Nøyaktig forhold til brukerbehov, REQ, kvalitetskrav og eksisterende SDP-Feature. |
| Functionality | Beskrive hvordan funksjonalitet realiseres og henger sammen gjennom systemet. | Grensen mot Feature, funksjon, scenario, pathway, sekvens og tilstandsmaskin. |
| Channel | Beskrive kommunikasjon mellom selvstendige enheter, med protokoll og kontrakt. | Logisk kontra fysisk kanal, endepunkter, retning, versjonering og støtte for flere deltakere. |

«Førsteklasses» betyr at objektene kan identifiseres, refereres til, analyseres og endres eksplisitt. De skal knytte sammen den strukturelle modellen; de trenger ikke ha hvert sitt isolerte deltre eller hver sin kopi av komponentene.

Disse tre er utgangspunkter, ikke et uttømmende eller vedtatt typehierarki. Særlig Functionality trenger en definisjon som tilfører informasjon utover Feature.

## 4. Faser, revisjoner og læring

### 4.1 Fasebasert utvikling

Utviklingen skal kunne organiseres i definerte faser, eksempelvis Concept og MVP, med identifiserte gjennomføringer som Concept1, Concept2 og MVP1. Arkitektur, design og implementasjon skal kunne knyttes til den aktuelle gjennomføringen og dens modenhet.

Store arkitektur- eller designendringer skal kreve en eksplisitt faseovergang eller ny gjennomføring innenfor samme fase. Hva som teller som «stor», og hvordan mindre endringer håndteres, må defineres i metoden.

Hver gjennomføring skal ha et identifiserbart source tree og tilhørende designgrunnlag. Den konkrete Git-/katalogorganiseringen er ikke avgjort. Sporbarhet til et tre innebærer ikke nødvendigvis et nytt repository, en permanent separat branch eller full omskriving.

### 4.2 Eksempler på tillatte utviklingsforløp

```text
Concept1
  |-- nye arkitektur-/designvalg -------------> Concept2
  |-- annen implementasjon av samme design ---> Concept2
  `-- tilstrekkelig læring -------------------> MVP1
          |
          `-- lærdom og endringsbehov bevares før MVP1-designet utformes
```

Concept2 er ikke et obligatorisk mellomsteg. En ny gjennomføring kan bevare designet og endre implementasjonen, eller endre både arkitektur, design og kode. Hva som videreføres, må fremgå.

### 4.3 Lærdom skal kunne bevares uten å beslutte neste design

Concept-fasen er prototyping og læring. Erfaringer kan vise at neste versjon bør ha andre ansvarsgrenser, kontrakter eller tekniske løsninger. Vi trenger en strukturert måte å bevare observasjonen, begrunnelsen, usikkerheten og ønsket virkning av en endring på.

Det er ikke alltid riktig å utforme MVP1s konkrete arkitektur mens Concept1 pågår. Et ønsket fremtidig målbilde skal derfor kunne uttrykkes uten å presenteres som et ferdig eller godkjent design. Issues og kommentarer kan være innspill og referanser; de bør ikke alene bære den samlede, gjeldende lærdommen.

Ved oppstart av MVP1 skal relevant lærdom brukes til å utforme et sammenhengende arkitektur- og designgrunnlag for det avtalte MVP1-omfanget. Ambisjonen om et komplett grunnlag for dette omfanget må forenes med at vi ikke skal detaljmodellere hele fremtidens implementasjon.

## 5. Utledede visninger og agentgrunnlag

Verktøyet skal kunne gi oversikt på flere nivåer og generere avgrensede utsnitt til bruk i agentinstruksjoner. Eksempler er systemoversikt, komponentinnhold, Feature-forløp, Functionality-atferd og Channel-kontrakter.

Et utsnitt må bevare nødvendig kontekst: relevante ansvar, eksterne berøringspunkter, begrensninger og kriterier for et ferdig resultat. Å gjøre vedlegget kort må ikke skjule avhengighetene som avgjør om oppgaven lykkes.

Eieren skal også kunne få en forståelig forklaring på hva som foreslås endret, hvorfor og hva som faktisk er kontrollert. Agentenes oversikt og eierens oversikt skal utledes fra samme identifiserte modellgrunnlag.

## 6. Endringsbeskrivelse og konsekvensanalyse

### 6.1 Design før implementasjon

Tiltenkte designendringer skal først beskrives i designkoden. Verktøyet skal sammenligne et identifisert utgangspunkt med et foreslått målbilde og beskrive hva som må endres. Feil og ny læring under implementering skal kunne føre til revisjon av forslaget gjennom en definert prosess.

Vi trenger å holde fra hverandre designet som er godkjent, designet vi foreslår, og det vi vet om det faktiske systemet. Et tidligere designvedtak kan ha avvik fra dagens kode; dagens kode kan også avvike fra det som faktisk er deployert.

### 6.2 Semantisk diff og rest-diff

Diffen skal forklare endringer i betydning, ansvar, kontrakter og forbindelser, ikke bare endrede tekstlinjer. Den skal kunne vises for hele systemet og filtreres gjennom Feature, Functionality, Channel eller andre designobjekter.

En rest-diff skal synliggjøre endringer som ikke er forklart gjennom de valgte objektene. Den må ikke gjemme bort arbeid som mangler tilknytning eller begrunnelse. Samme endring kan være relevant for flere objekter.

### 6.3 Konsekvenser gjennom systemet

Analysen skal hjelpe oss å finne både det som blir berørt, og det som må endres eller kontrolleres for at en ønsket effekt skal bli tilgjengelig fra ende til ende. Den skal også støtte vurdering av om sterke eller vidtrekkende avhengigheter bør reduseres gjennom endret arkitektur.

Systemet må kunne vise manglende informasjon. Fravær av en modellert avhengighet er ikke i seg selv bevis på at avhengigheten ikke finnes.

## 7. Samsvar mellom modell og faktisk programvare

Ambisjonen er å kunne få programmatiske svar på om implementasjonen følger avtalte egenskaper i designet. Vi skal undersøke flere evidenskilder:

- Strukturerte annotasjoner i kildefiler, eventuelt Doxygen eller en tilsvarende mekanisme.
- Uttrekk fra kildekode, typesystem, compiler-API eller annen statisk analyse.
- Kompilatorens mellomrepresentasjon, eksempelvis LLVM IR, der dette er relevant.
- Tester og observasjoner som kontrollerer faktisk atferd og integrasjon.

Det skal også undersøkes hvordan påstander i kommentarer eller annotasjoner kan kontrolleres mot koden de viser til. Samsvar mellom to beskrivelser alene beviser ikke samsvar med faktisk funksjon.

Å kompilere hele eller deler av designspråket til LLVM IR er en utforskende idé. Det er ikke en forutsetning for første løsning, og vi må undersøke hvilken sammenligning eller hvilket bevis det faktisk muliggjør.

## 8. Målbare mål for metoden

| ID | Ønsket resultat | Hva en senere pilot må demonstrere |
| --- | --- | --- |
| M01 | Rask systemforståelse | En ny agent kan finne gjeldende ansvar, kontrakter og Feature-forløp fra et kompakt utsnitt. |
| M02 | Kontrollerte endringer | Før/etter-modell og begrunnelse finnes før en designendrende kodeoppgave starter. |
| M03 | Bevart lærdom | Et Concept-funn kan overføres til planlegging av MVP uten et forhåndsbestemt MVP-design. |
| M04 | Synlig helhet | En endring i en Channel avdekker relevante produsenter, konsumenter og berørte Feature-forløp. |
| M05 | Fullstendig endringsregnskap | Alle oppdagede modelldeltaer fremgår i total-diff, objektvisninger eller rest-diff uten tap. |
| M06 | Etterprøvbart samsvar | Resultatene angir kontrollert egenskap, kilde-/build-identitet, evidens og dekningsbegrensning. |
| M07 | Kontrollert avvik | En agent kan ikke gjøre brudd på designet akseptabelt ved stille å endre forventningen. |
| M08 | Begrenset vedlikehold | Kostnad i modellredigering, antall håndskrevne opplysninger og kontekststørrelse måles. |
| M09 | Gjenbrukbar metode | En enkel prosjektprofil fungerer uten hele kompleksiteten til et distribuert system. |

Tallfestede terskler for tidsbruk, kontekststørrelse og vedlikehold må avtales før en pilot måles. Denne studien hevder ikke at målene allerede er oppnådd.

# Del II — Innledende studie og forslag

## 9. Vurdering: hva som må skilles før vi velger språk

**Forslag:** Del løsningen konseptuelt i fire deler:

1. **Metode:** Hvem kan beslutte hva, hvilket grunnlag kreves og når må arbeidet revurderes?
2. **Semantisk modell:** Hvilke objekter, relasjoner, regler og identiteter uttrykker designet?
3. **Representasjon og verktøy:** Hvordan redigeres, valideres, sammenlignes og presenteres modellen?
4. **Evidens:** Hvordan undersøkes kode, bygg og kjøring uavhengig av modellens påstander?

Et rikt språk kan mangle en god endringsprosess. Et skjema kan være gyldig selv om implementasjonen bryter arkitekturen. En god prosess kan på sin side bli tung hvis informasjonen må kopieres manuelt mellom mange dokumenter.

## 10. Forslag til presise modellbegreper

### 10.1 En graf med ulike relasjoner

Modellen bør kunne ha en strukturell nedbrytning, men også relasjoner som går på tvers. Et enkelt mappetre er utilstrekkelig for delte tjenester, flere konsumenter og Features som går gjennom mange komponenter.

| Begrep | Foreslått betydning |
| --- | --- |
| Component | En ansvarsbærende del av arkitekturen på et oppgitt nivå. |
| Design constituent | En modellert bestanddel av designet; ikke synonymt med klasse eller fil. Eksakt mapping mot IEEE 1016 må studeres. |
| Feature | En varig kapabilitet med formål, krav og observerbare akseptkriterier. |
| Functionality | Foreløpig kandidat: et sammenhengende atferdsansvar som realiserer deler av én eller flere Features. Bør bare bli egen type hvis piloten viser at skillet er nyttig. |
| Scenario / witness | Et konkret observerbart forløp som kan brukes til å undersøke en Feature-påstand. |
| Pathway | En definert kjede av ansvar og overganger. Strukturelt tillatt rute og faktisk runtime-forløp må ha ulik merking. |
| Channel | En logisk kommunikasjonsforbindelse med identifiserte deltakere og kontrakt. Transport og deployering bindes separat ved behov. |
| Contract | Tillatte interaksjoner, data, betydning, feil og kompatibilitetsforutsetninger. |
| Implementation binding | En kobling fra designobjekt til package, symbol, kodeområde, build-variant eller deployeringsenhet. |
| Constraint / invariant | En uttrykt regel, for eksempel forbudt avhengighet eller krav om aktuell revisjon. |

Eksempel på relasjonstyper som kan prøves: «inneholder», «realiserer», «bruker kontrakt», «sender», «mottar», «må bevares», «verifiseres av». Relasjonene må ha definert retning og betydning. En generell `depends_on`-kant alene forklarer sjelden hva en endring krever.

Én Feature kan bruke flere Functionalities; én Functionality kan støtte flere Features. Én Channel kan inngå i flere forløp. Domain, logisk komponent, lag, prosess og source-folder bør ikke automatisk være samme objekt.

### 10.2 Stabil identitet og endring over tid

**Forslag:** Skill objekt-ID fra navn, plassering, revisjon, fase og release. Ved rename beholdes identiteten når ansvaret er det samme. Ved split/merge trengs eksplisitt mapping av ansvar og bevaringskrav; navnelikhet er ikke tilstrekkelig.

Git lagrer revisjonene, mens modellen eller et separat endringsobjekt forklarer hva endringen betyr. Vi bør prøve en separat transition-beskrivelse før vi bygger all historikk inn i hvert designobjekt. Identitetsreglene må samordnes med #7s foreløpige regler om scope, reservasjon og flytting mellom repositories.

### 10.3 Deklarert dekning fremfor digital tvilling

**Forslag:** Beskriv manuelt det som uttrykker hensikt og styrer utvikling: ansvar, kontrakter, vesentlige tilstander, invariants, forbudte avhengigheter og akseptscenarier. Hent mekaniske fakta som filinventar og imports fra verktøy når det er mulig.

Hvert modellområde bør oppgi hva som er bindende, hva som er beskrivende, hva som utledes, og hva som er utenfor dekningen. En lett profil kan begynne med komponentgrenser, kritiske kontrakter og ett Feature-forløp. Ekstra detaljer innføres der feilrisiko eller endringsbehov begrunner dem.

## 11. Forslag til fase- og endringsmodell

### 11.1 Separate koordinater

Følgende opplysninger bør ikke presses inn i ett versjonsnummer:

| Koordinat | Eksempel | Hva den identifiserer |
| --- | --- | --- |
| Livsløpsfase | Concept, MVP | Formål og forventet modenhet. |
| Fasegjennomføring | Concept1, Concept2, MVP1 | En navngitt gjennomføring med omfang og inngangsgrunnlag. |
| Modellrevisjon | D17 | Et bestemt designinnhold. |
| Kilderevisjon | Git commit SHA | Et konkret source tree. |
| Build / deployment | Artefaktdigest + konfigurasjon | Hva som er bygget eller faktisk kjøres. |
| Aksept | Designgodkjent, Feature-verifisert, eierakseptert | Hvilken beslutning eller evidens som foreligger. |

Innenfor Concept1 kan det finnes flere ordinære modell- og kilderevisjoner. En beslutning om ny fasegjennomføring er en egen hendelse. MVP1 er heller ikke det samme som programvarens SemVer-release 1.0.0.

### 11.2 Læringsregister

**Foreslått minimumsinnhold i en læringspost:** stabil ID; hvor observasjonen kom fra; hva vi observerte; evidens eller uttrykkelig antakelse; konsekvens; ønsket egenskap ved en fremtidig løsning; vurderte alternativer; uavklarte spørsmål; og senere beslutning/videreføring.

Illustrativt eksempel, ikke et produktfunn:

```yaml
id: LEARN-EXAMPLE-01
origin: Concept1
observation: "UI og domenelogikk kan ikke testes uavhengig i forsøket."
evidence: "Referanse til konkret forsøk må fylles inn."
desired_outcome: "Domeneforløpet skal kunne kjøres uten UI-renderer."
options:
  - "Tydelig port mellom domene og presentasjon"
  - "Separat presentasjonsprosess"
decision: deferred
candidate_target: MVP1
```

Posten binder oss til å behandle problemet, ikke til å velge en bestemt arkitektur. Ved MVP1-planlegging får relevant lærdom en eksplisitt disposisjon: tatt inn i krav/design, undersøkes videre, utsatt eller avvist med begrunnelse. Et beslutningsnotat kan referere læringsposten uten å overskrive den opprinnelige observasjonen.

### 11.3 Fasegater

**Forslag til gater; detaljene er ikke vedtatt:**

| Gate | Minimum før videre arbeid | Mulige utfall |
| --- | --- | --- |
| Starte gjennomføring | Formål, omfang, relevant lærdom, foreløpige grenser og planlagte undersøkelser. | Start Concept2/MVP1, eller avklar manglende grunnlag. |
| Starte avgrenset implementasjon | Sammenhengende design for endringen, eksplisitte ukjente forhold, kontrakter, konsekvenser og akseptscenarier. | Implementer, eller utfør et avgrenset forsøk først. |
| Endre arkitektur/design vesentlig | Begrunnelse, før/etter-ansvar, berørte kapabiliteter, migrasjon og nye kontrollbehov. | Ny fasegjennomføring, faseovergang eller avvist/utsatt forslag. |
| Integrere resultat | Kode og modell er identifisert; relevante struktur-, kontrakt- og integrasjonskontroller foreligger. | Integrer innenfor uttrykt aksept, eller returner et avgrenset funn. |
| Avslutte gjennomføring | Resultater, lærdom, kjente avvik og uferdig arbeid har eksplisitt disposisjon. | Fortsett i ny gjennomføring, gå til neste fase, eller avslutt. |

Et «stort» designinngrep kan defineres gjennom beskyttede grenser: flyttet ansvar mellom komponenter, endret prosess-/trust-grense, brutt offentlig kontrakt, endret dataeierskap eller endret Feature-semantikk. Antall endrede linjer er et dårlig hovedkriterium.

Mindre rettelser innenfor godkjente kontrakter kan behandles som vanlige revisjoner. En utforskende spike kan starte med lettere designgrunnlag dersom hypotese, avgrensning og disponering av resultatet er tydelig. Ingen fase trenger å love et «endelig design» for resten av produktets levetid.

## 12. Forslag: tre grunnlag og tre forskjellige sammenligninger

La **B** være godkjent designbaseline, **T** foreslått måldesign og **O** observerte implementasjonsfakta for en identifisert kode-/build-/konfigurasjonsvariant. O er et faktagrunnlag med kjent dekning, ikke en automatisk fullstendig rekonstruksjon av designet.

| Sammenligning | Spørsmål |
| --- | --- |
| B → T | Hva foreslås endret i designet, og hvorfor? |
| O mot B | Hvilke observerte forhold samsvarer med eller avviker fra gjeldende design? |
| O mot T | Hva gjenstår eller avviker i implementasjonen av måldesignet? |

Den første kan være en semantisk modelldiff. De to siste trenger oversetting av kodefakta til egenskaper og relasjoner modellen kan kontrollere. De er ikke nødvendigvis vanlig diff mellom to like dokumenter.

Observasjon av kildekode er ikke observasjon av en deployering. Påstander om «systemet som kjører nå» trenger identifisert artefakt, konfigurasjon og miljø. Eksisterende avvik må registreres før en ny endring, slik at de ikke feilaktig tilskrives den nye implementasjonen.

### 12.1 Semantisk endringsregnskap

Hver oppdaget endring bør ha én identitet og kunne vises i flere utsnitt. Eksempler på operasjoner er lagt til, fjernet, flyttet ansvar, endret kontrakt, splittet, slått sammen og endret invariant. Automatisk matching må ikke late som den vet om to ulike objekter har samme hensikt; usikker mapping krever avklaring.

```text
D_total = alle registrerte semantiske endringer mellom B og T
D_objekt = unionen av endringer forklart i valgte objektvisninger
D_rest = D_total minus D_objekt
```

Objektvisningene kan overlappe. En kontraktendring trenger ikke telles som tre uavhengige endringer selv om den berører en Feature, en Functionality og en Channel. Rest-diff kan inneholde legitim infrastrukturendring eller manglende tilknytning; begge må få en forklaring. Dette regnskapet sier bare noe om oppdagede modelldeltaer. Et separat kodeinventar må finne kodeendringer som mangler modell-/oppdragskobling.

### 12.2 Konsekvenser er mer enn grafnaboer

Analysen bør skille mellom:

- Direkte endrede objekter.
- Objekter som potensielt påvirkes gjennom en definert relasjon.
- Avledede forpliktelser: noe må vurderes, tilpasses eller testes.
- Konkret avgjørelse: endres, beholdes med begrunnelse, eller er uavklart.

Ikke alle konsumenter må endres ved enhver schemaendring. Men de relevante konsumentene må vurderes. Hver påvirkningspåstand bør vise hvilken regel og relasjonskjede som førte til den. Dynamisk registrering, konfigurasjon og eksterne tjenester kan kreve andre kilder enn statiske imports.

Færre avhengigheter er heller ikke alltid bedre arkitektur. Resultatet må vurderes mot ansvar, kohesjon, ytelse, driftsbehov og migrasjonskostnad.

## 13. Fra designbeskrivelse til implementasjon

Dette er en **foreslått arbeidssekvens**, ikke en arbeidsflyt foreskrevet av IEEE 1016.

1. **Bekreft oppdrag og utgangspunkt.** Finn behov, berørte Feature-revisjoner, B, aktuell kildebaseline og kjente avvik.
2. **Undersøk horisontalt.** Følg berørte ansvar, kontrakter, konsumenter og forløp gjennom systemet. Avklar ukjent wiring eller kritiske designvalg før gjennomføring.
3. **Beskriv T og endringsbegrunnelsen.** Registrer hva som skal bevares, endres, flyttes og eventuelt migreres. Velg riktig revisjon eller fasegate.
4. **Utled forpliktelser og aksept.** Beskriv forventet ende-til-ende-resultat, feilforløp, nødvendige kontroller og kriterier for avvik.
5. **Planlegg vertikale Slices.** Rekkefølgen styres av avhengigheter og risiko. En kompatibel kontraktutvidelse eller adapter kan være første steg; alle lag må ikke omskrives samtidig.
6. **Gi agenten et identifisert oppdragsutsnitt.** Agenten får kontrakter, tillatte endringer, berørte naboer og stoppregler.
7. **Implementer og innhent evidens.** Kontroller kodefakta og relevant runtime-atferd. Nye designbehov returneres som et forslag fremfor å innarbeides stille.
8. **Integrer og kontroller Feature-forløpet.** Lokalt godkjente Slices erstatter ikke kontroll på faktisk integrasjonsbaseline og aktuell konfigurasjon.
9. **Registrer aksept og gjenværende arbeid.** Bevar forskjellen på implementert, integrert, verifisert, eierakseptert og inkludert i release. Oppdater gjeldende pekere gjennom riktig beslutning.

### 13.1 Oppdragsutsnitt til Worker og Verifier

**Foreslått minimum:** oppdrags-ID og mål; B/T-revisjoner; kildebaseline; generator-/regelversjon; berørte Feature-/Functionality-/Channel-ID-er; ansvar og kontrakter; tillatte kodeområder; delte berøringspunkter; bevaringskrav; konkrete kontrollkrav; kjente avvik; uavklarte forhold; og hvem som kan beslutte en designendring.

Utsnittet bør gi et kort sammendrag med lenker til nødvendige detaljer og vise hva som er utelatt. Det må kunne avvises som foreldet dersom modell, kontrakt eller integrasjonsgrunnlag endres. Worker og Verifier trenger samme kontraktgrunnlag, mens verifikasjonen også må kunne oppdage feil og hull i dette grunnlaget.

### 13.2 Hindre reparasjonsløkker

Hvert oppdrag bør ha et avtalt rework-budsjett og tydelige stoppkriterier. Eksempler: samme feilklasse gjentas uten ny forklaring; nødvendig arbeid krysser en beskyttet grense; ny konsument dukker opp; eller forventet Feature-effekt mangler til tross for lokale testpass.

Ved stopp skal agenten levere observasjon, hypotese og konkret behov for ny analyse eller designbeslutning. En feil i modellen kan rettes, men endring av kontrollgrunnlaget må være synlig. Verktøyet skal ikke normalisere drift ved å overskrive godkjent design med det koden tilfeldigvis gjør.

## 14. Gjennomgående eksempel: en Channel-endring som ellers kan gi tunnelsyn

**Hypotetisk case.** Dette er ikke en kontrollert beskrivelse av Ponsse, HSX eller annen produktkode.

Feature `FEAT-LENGTH` lar en operatør se aktuell målt lengde. Functionality-kandidaten `FUNC-PRESENT-LENGTH` fører målingen fra domenelogikk til presentasjon. Channel `CH-MEASUREMENT` sender måleverdi til klienten.

```text
Målekilde -> domenelogikk -> CH-MEASUREMENT -> klientadapter -> presentasjon
                                 |
                                 `-> historikklagring
```

**B:** Kanalens verdi er i millimeter; presentasjonen viser centimeter. Historikklagring bruker også millimeter.  
**T:** Kanalen skal sende meter; begge konsumentene må fortsatt oppfylle sine avtalte egenskaper. Utvendig visning skal være uendret.

En naiv lokal rettelse endrer bare domenets serialisering. JSON-typen kan fortsatt være `number`, alle imports kan være uendret, og tester som kun kontrollerer tallformat kan passere. Likevel er visning og historikk feil.

En nyttig modell og analyse skal gi:

| Utledning | Konsekvens / kontroll |
| --- | --- |
| Kanalens semantiske enhet endres | Registrer kontraktbrudd selv om syntaktisk datatype er den samme. |
| Klientadapter konsumerer kanalen | Vurder konvertering og kontraktversjon. |
| Historikklagring konsumerer kanalen | Vurder lagringsenhet, kompatibilitet og behandling av eksisterende data. |
| FEAT-LENGTH bruker forløpet | Krev scenario med kjent fysisk verdi helt frem til faktisk presentasjon. |
| Gamle og nye komponenter kan sameksistere | Beskriv eksplisitt overgang, kompatibilitetsvindu og eventuell rollback. |
| Deployment-konfigurasjon velger adapter | Bekreft at riktig implementasjon faktisk registreres og brukes. |

Eksempel på witness: En input som representerer 1 meter skal ende som 100 centimeter i presentasjonen og korrekt verdi i historikken, med avtalt presisjon. En headless test er nyttig; den må suppleres hvis den ikke bruker samme binding som den faktiske klienten.

Diffen kan vises både under Feature, Functionality og Channel, med samme endrings-ID. Endret deployeringskonfigurasjon skal inngå med forklaring eller fremgå i rest-diff. Hvis historikk-konsumenten ikke er kartlagt, må uavhengig inventar/observasjon kunne avdekke hullet; grafen alene kan ikke utlede en ukjent konsument.

Hvis det i Concept1 oppdages at enheter håndteres inkonsistent, kan dette først registreres som lærdom: «Neste gjennomføring trenger en eksplisitt enhetskontrakt.» Det bestemmer ikke på forhånd om MVP1 skal bruke meter, millimeter eller en type som bærer enheten.

## 15. Standarder rundt designbeskrivelse, livsløp og implementasjon

### 15.1 Kildegrunnlag og avgrensning

Oversikten nedenfor er kontrollert mot offentlige primærkilder 2026-09-14, hovedsakelig ISO-/IEEE-kataloger og sammendrag. Full normativ tekst er ikke gjennomgått klausul for klausul. Vi kan derfor angi dokumentert virkeområde og foreslå bruk, men ikke hevde standardkonformitet eller en komplett liste over normative referanser i IEEE 1016.

«Økosystem» brukes her som et praktisk kart over beslektede behov. Det betyr ikke at alle standardene er del av én formell IEEE 1016-familie, eller at 42010 erstatter 1016.

### 15.2 Design og arkitektur

| Referanse og observert status | Dokumentert virkeområde | Foreslått bruk i SDP |
| --- | --- | --- |
| [IEEE 1016-2009](https://standards.ieee.org/ieee/1016/4502/) — Inactive-Reserved; inaktivert 2020-03-05 | Innhold og organisering av Software Design Descriptions, for overordnet og detaljert design. Foreskriver ikke én designmetode eller ett designspråk. | Faglig grunnlag for hva modellen og avledede SDD-visninger bør kunne uttrykke. Statusen må følge referansen. |
| [ISO/IEC/IEEE 42010:2022](https://www.iso.org/standard/74393.html) — publisert | Krav til arkitekturbeskrivelser, beskrivelsesrammeverk, språk, viewpoints og model kinds; skiller arkitektur fra beskrivelsen av den. | Terminologi og struktur for concerns, modelltyper og målrettede visninger. |
| [ISO/IEC/IEEE 42020:2019](https://www.iso.org/standard/68982.html) — publisert; revisjon under utvikling | Prosesser for styring, forvaltning og utforming av arkitekturer gjennom levetiden. | Ansvar og beslutningsprosess for å utforme og revidere arkitektur. |
| [ISO/IEC/IEEE 42030:2019](https://www.iso.org/standard/73436.html) — publisert; markert for revisjon | Rammeverk for organisering og dokumentering av arkitekturevaluering. | Vurdere alternativer, kvalitet og risiko før større overgang. |

**Presisering av 1016:** Den er fortsatt relevant for idéarbeidet, men bør ikke omtales som en aktiv, oppdatert universalløsning. En SDD behøver ikke være én stor manuelt skrevet Markdown-fil; IEEE beskriver også andre medier og verktøybaserte representasjoner. [IEEE 1016-2009](https://standards.ieee.org/ieee/1016/4502/)

**Viewpoint versus view:** Arbeidsdefinisjonen vi bør bruke, er at et viewpoint angir hvilke spørsmål og konvensjoner en visning skal følge, mens et view er den konkrete fremstillingen for et bestemt modellgrunnlag. Dette må presiseres mot valgt standardutgave under videre modellarbeid. 42010 velger ikke diagramverktøy eller utviklingsmetode for oss. [ISO/IEC/IEEE 42010:2022](https://www.iso.org/standard/74393.html)

Påstanden i Gemini-utkastet om «IEEE 1016s 12 standard-viewpoints» er ikke brukt som et normativt krav her. Nøyaktig liste, terminologi og krav til valg/tilpasning må kontrolleres i full 1016-tekst før vi gjør en slik mapping. Vi skal ikke anta at alle viewpoints må materialiseres for hvert prosjekt.

### 15.3 Livsløp, krav, dokumentasjon og kontrollert endring

| Referanse og observert status | Dokumentert virkeområde | Foreslått bruk i SDP |
| --- | --- | --- |
| [ISO/IEC/IEEE 12207:2026](https://standards.ieee.org/ieee/12207/11416/) — Active; erstatter 2017-utgaven | Felles prosessrammeverk for programvarens livsløp; prosessene kan brukes samtidig, iterativt og rekursivt. | Forankre sammenhengen mellom design, realisering, integrasjon, verifikasjon og videre utvikling. Velg relevante prosesser per fase. |
| [ISO/IEC/IEEE 15288:2023](https://www.iso.org/standard/81702.html) — publisert | Systemlivsløpsprosesser; foreskriver ikke én livsløpsmodell eller metode. | Utvide konteksten der programvaren inngår i maskinvare, operatørarbeid og større systemer. |
| [ISO/IEC/IEEE 24748-1:2024](https://www.iso.org/standard/84709.html) — publisert | Veiledning om livsløpsstyring, modeller, stadier og tilpasning, i sammenheng med 12207 og 15288. | Særlig relevant for å definere Concept/MVP, inngangs-/utgangskriterier og nye gjennomføringer. |
| [ISO/IEC/IEEE 29148:2018](https://www.iso.org/standard/72089.html) — publisert; markert for revisjon | Requirements engineering gjennom livsløpet. | Skille behov og normative krav fra løsningsvalg; knytte dem til Feature og verifikasjon. |
| [ISO/IEC/IEEE 15289:2019](https://www.iso.org/standard/74909.html) — publisert; bekreftet 2025 | Formål og innhold i livsløpets informasjonsprodukter. | Bestemme hvilken informasjon vi trenger, og hvilke dokumenter som kan utledes fra modellen. |
| [ISO 10007:2017](https://www.iso.org/standard/70400.html) — publisert; markert for revisjon | Veiledning for konfigurasjonsstyring fra konsept til avvikling. | Baselines, identifiserte konfigurasjoner og kontrollert endring fremfor umerket design drift. |
| [ISO/IEC/IEEE 14764:2022](https://www.iso.org/standard/80710.html) — publisert | Programvarevedlikehold som del av livsløpet. | Undersøke hvordan endringsanalyse, gjennomføring og videre vedlikehold bør organiseres. |

Ingen av disse kildene gjør navnene Concept1, Concept2 og MVP1 til standardiserte stadier. Den konkrete fasemodellen er et SDP-valg. 12207s prosesser må heller ikke fremstilles som en obligatorisk engangssekvens. [ISO/IEC/IEEE 12207:2026](https://standards.ieee.org/ieee/12207/11416/)

Utgavene er heller ikke automatisk harmonisert på alle detaljer: 14764:2022 viser uttrykkelig til vedlikeholdsprosessen i 12207:2017. En senere SDP-standardprofil må kontrollere slike utgavekoblinger før den overfører prosess- eller klausulreferanser til 12207:2026. [ISO/IEC/IEEE 14764:2022](https://www.iso.org/standard/80710.html)

### 15.4 Verifikasjon, validering og kvalitet under implementasjon

| Referanse og observert status | Dokumentert virkeområde | Foreslått bruk i SDP |
| --- | --- | --- |
| [IEEE 1012-2024](https://standards.ieee.org/ieee/1012/7324/) — Active | Verifikasjon og validering av systemer, programvare og maskinvare. | Utforme en begrunnet V&V-strategi, med konkrete oppgaver, evidens og passende uavhengighet. |
| [ISO/IEC/IEEE 29119-2:2021](https://www.iso.org/standard/79428.html) — publisert | Prosesser for styring og gjennomføring av programvaretesting på tvers av livsløpsmodeller. | Koble design-/endringsgrunnlag til planlagte tester og faktisk utføring. |
| [ISO/IEC 25010:2023](https://www.iso.org/standard/78176.html) — publisert | Produktkvalitetsmodell i SQuaRE-serien. | Sørge for at kontrakter og aksept også omfatter relevante kvalitetsegenskaper, ikke bare happy-path-funksjon. |

**Vurdering:** 1016 hjelper med designbeskrivelsen, men er ikke alene svaret på hvordan implementasjonen angripes. Det mest relevante studiesporet for selve gjennomføringen er samspillet mellom 12207, 24748-1, konfigurasjonsstyring, V&V og testprosesser, anvendt på en konkret endring. Arbeidssekvensen i kapittel 13 er vårt forslag til en slik anvendelse, ikke et standardkrav.

**Videre standardarbeid:** Velg en liten referanseprofil og lag en sporbar mapping mellom aktuell standardutgave, relevant begrep/prosess, SDP-felt eller gate, kontrollform og eventuelle gap. Det krever tilgang til fulltekst for de standardene der vi vil hevde konkret samsvar. Unngå å gi en lang litteraturliste status som implementert metode.

## 16. Designspråk og metodikker å bygge videre på

Dette kapitlet bruker den eksisterende P04-kartleggingen som underlag. Den er datert 2026-09-10; profilene er ikke alle undersøkt på nytt i denne revisjonen. Ingen kandidat er valgt, og P04s eksempler er ikke parser-/runtime-testet.

| Byggekloss | Relevante kandidater | Vurdering for dette mandatet |
| --- | --- | --- |
| Struktur og visninger | [Structurizr DSL](research/existingDesignLanguages/structurizr-dsl.md), [LikeC4](research/existingDesignLanguages/likec4.md) | Gode forsøksspor for én modell med flere utsnitt. Undersøk behovet for ekstra Feature-, atferds- og endringssemantikk. |
| Rikere system-/designmodell | [SysML v2](research/existingDesignLanguages/sysml-v2.md), [UML](research/existingDesignLanguages/uml.md), [AADL](research/existingDesignLanguages/aadl.md), [Context Mapper](research/existingDesignLanguages/context-mapper-cml.md) | Vurder presisjon og innebygde begreper opp mot redigeringskostnad og agentbruk. |
| Kontrakter | [OpenAPI](research/existingDesignLanguages/openapi.md), [AsyncAPI](research/existingDesignLanguages/asyncapi.md), [Protobuf](research/existingDesignLanguages/protobuf.md), [Smithy](research/existingDesignLanguages/smithy.md) | Referer eksisterende kontrakter der de passer. Kompletter eksplisitt med semantikk som ikke dekkes av den valgte kontrakten. |
| Constraints og avgrenset atferd | [CUE](research/existingDesignLanguages/cue.md), [OCL](research/existingDesignLanguages/ocl.md), [SCXML](research/existingDesignLanguages/scxml.md), [TLA+](research/existingDesignLanguages/tla-plus.md), [Alloy](research/existingDesignLanguages/alloy.md) | Ulike verktøy for forskjellige spørsmål; modellkontroll er ikke automatisk kontroll av produktkoden. |
| Modellendring og provenance | [Epsilon](research/existingDesignLanguages/epsilon.md), [Edapt](research/existingDesignLanguages/edapt.md), [PROV-O](research/existingDesignLanguages/prov-o.md) | Studer matching, transformasjon, historikk og avledning før vi bygger egen transition-infrastruktur. |
| Presentasjon | [Mermaid](research/existingDesignLanguages/mermaid.md), [PlantUML](research/existingDesignLanguages/plantuml.md), [D2](research/existingDesignLanguages/d2.md) | Mulige utdataformater. Diagramnotasjon alene definerer ikke hele systemmodellens semantikk. |

C4 er en nyttig abstraksjons- og visningstilnærming, men bør ikke uten videre gjøres til metamodel for alle designbestanddeler og atferdsobjekter. JSON/YAML kan serialisere en modell, men velger ikke semantikken. Langium og Xtext er mulige verktøy for å bygge språk/editor; de bør vurderes etter at behovet for egen syntaks er dokumentert, ikke som alternativer av samme type som Mermaid.

To metodiske tillegg fortjener undersøkelse:

- **Arcadia/Capella:** Arcadia skiller blant annet behovsanalyse fra arkitekturutforming og har en verktøystøttet modelleringstilnærming. Vår vurdering er at dette er relevant for å studere sammenhengen mellom operativt behov, funksjoner og struktur. Vi har ikke gjennomført en Capella-pilot eller bekreftet dekning av SDP-transition og agentutsnitt. [Eclipse om metode og verktøy](https://www.eclipse.org/community/eclipse_newsletter/2017/december/article3.php)
- **ATAM:** SEIs metode undersøker avveininger i programvarearkitektur. Vår vurdering er at scenario- og kvalitetsbasert evaluering kan gi bedre grunnlag for faseoverganger enn antall avhengigheter alene. Full ATAM er ikke foreslått som obligatorisk seremoni for små prosjekter. [SEI: The Architecture Tradeoff Analysis Method](https://www.sei.cmu.edu/library/the-architecture-tradeoff-analysis-method/)

**Foreløpig anbefaling:** Definer hvilke spørsmål modellen skal besvare og hvilke endringer den skal representere før språkvalg. Prøv først gjenbruk eller en liten semantisk utvidelse. Eget språk fra bunnen av bør begrunnes med konkrete gap demonstrert i samme case.

## 17. Hva kan verifiseres programmatisk?

### 17.1 Kontroller på forskjellige nivåer

| Nivå | Eksempel | Hva et positivt resultat ikke beviser |
| --- | --- | --- |
| Modellform | ID-er finnes; referanser løses; felter har gyldig type. | At modellen beskriver riktig system. |
| Modellregler | En forbudt dependency finnes ikke i modellgrafen. | At koden ikke har avhengigheten. |
| Kodestruktur | Faktiske imports/typesymboler samsvarer med definerte grenser. | At riktig funksjon er koblet inn ved kjøring. |
| Kontrakt og kompatibilitet | Observerte endepunkter og meldinger følger angitt schema/versjon. | At for eksempel enhet, autorisasjon og feilhåndtering er riktig hvis dette ikke kontrolleres. |
| Integrert atferd | Et definert Feature-scenario virker på aktuell build og konfigurasjon. | Alle mulige input, tidsforløp, miljøer og feiltilstander. |
| Formell egenskap | En presist formulert invariant er bevist innenfor en formell modell og dens forutsetninger. | At vilkårlig produktkode automatisk er en korrekt realisering av modellen. |

**Forslag til resultatmodell:** `satisfied`, `violated`, `unknown`, `not_applicable`. Et kontrollresultat trenger egenskap, metode, verktøy-/regelversjon, kilde-/build-identitet, variant, evidens og dekningsomfang. Parserfeil eller manglende dekning må ikke bli et tomt «ingen avvik»-resultat. Policyen bestemmer hvilke ukjente forhold som blokkerer en bestemt gate.

«Korrekte programmatiske svar» må bety korrekte svar innenfor en definert semantikk og kontrollmodell. Verktøyet skal ikke love å avgjøre full atferdsekvivalens mellom en vilkårlig designbeskrivelse og vilkårlig programkode.

### 17.2 Doxygen og annotasjoner

Doxygen kan generere XML som andre verktøy kan bearbeide. Det gir et konkret mulig integrasjonspunkt for dokumentasjon og kodekoblinger. [Doxygen: XML output](https://www.doxygen.nl/manual/customize.html#xmloutput)

**Vurdering:** Annotasjoner er nyttige for å knytte et symbol til en stabil design-ID og uttrykke hensikt som ikke kan utledes fra kode. Det er fortsatt en deklarasjon. Hvis både modellen og kommentaren sier «bruker CH-MEASUREMENT», men koden benytter en annen vei, kan en sammenligning av tekstene gi falsk trygghet.

Undersøk derfor annotasjoner sammen med uavhengig symbol-/dependencyanalyse og relevante tester. Ikke kopier hele designbeskrivelsen inn i kodekommentarer. Doxygen-støtte og uttrekkskvalitet må prøves for det konkrete programmeringsspråket; dette er ikke en universell løsning for alle SDP-prosjekter.

### 17.3 AST og språkspesifikke analyseverktøy

Clang tilbyr blant annet AST-matchere for å finne mønstre i C/C++-kode. Dette er et dokumentert eksempel på et analysegrunnlag nær kildekodens struktur. [Clang AST Matchers](https://clang.llvm.org/docs/LibASTMatchers.html)

**Vurdering:** Start med det pilotens språkverktøy kan hente pålitelig: symboler, imports, kall, kontraktimplementasjoner og registreringer. Et syntakstre alene har ikke nødvendigvis typeoppløsning eller komplett kallgraf. Plugins, refleksjon, generert kode og konfigurasjonsstyrt wiring krever eksplisitt behandling. Analyzerens funn må angi hvilke build-varianter og kodeområder som ble undersøkt.

### 17.4 LLVM-sporet

LLVM IR er en typet, lavnivå mellomrepresentasjon med blant annet funksjoner, instruksjoner, kontrollflyt og metadata. Den kan være et grunnlag for avgrenset kodeanalyse. [LLVM Language Reference](https://llvm.org/docs/LangRef.html)

**Vurdering:** En Feature-ID, brukerhensikt eller logisk Channel følger ikke automatisk med som et forstått domeneobjekt. Mapping må bevares eller konstrueres. Optimering, valgt frontend og debug-/metadataoppsett påvirker hva vi kan gjenfinne. LLVM er heller ikke et gitt felles mellomformat for alle aktuelle språkverktøykjeder.

Å kompilere både design og implementasjon til LLVM gjør ikke meningssammenligningen automatisk. Ulike instruksjoner kan realisere samme egenskap, og samme lokalberegning kan inngå i forskjellige systemforløp. Et design med flere lovlige realiseringer trenger dessuten en presis relasjon mellom spesifikasjon og implementasjon, ikke krav om identisk IR.

Et senere LLVM-forsøk kan være verdifullt for én avgrenset egenskap med eksplisitt mapping og forutsetninger. En egen **design-IR** som normaliserer objekter, relasjoner og kontrakter er et annet konsept enn LLVM IR og kan være nyttig uten å generere maskinkode.

**Foreløpig anbefaling:** Begynn med modellkontroll, kodebindinger, struktur-/kontraktanalyse og Feature-witnesses. Behold LLVM og formelle bevis som avgrensede studiespor. Velg dem når et konkret kontrollbehov rettferdiggjør kostnaden.

# Del III — Videre arbeid

## 18. Åpne spørsmål og beslutninger som må tas

| ID | Spørsmål | Egnet undersøkelse |
| --- | --- | --- |
| Q01 | Når tilfører Functionality et eget ansvar utover Feature, pathway og scenario? | Beskriv samme lille case med og uten objekttypen. |
| Q02 | Hvilke grenser utløser Concept2/MVP2 fremfor vanlig revisjon? | Prøv lokal Fix, intern refactor, kontraktbrudd og flyttet prosessansvar. |
| Q03 | Hvor komplett må designet være før implementasjon på hvert modenhetsnivå? | Definer minimumsprofil og eksplisitt tillatt uavklart innhold. |
| Q04 | Hva skal håndskrives, hva skal utledes og hva skal bare refereres? | Mål redigeringskostnad og risiko for motstridende opplysninger. |
| Q05 | Hvordan overlever objektidentitet rename, split, merge og repo-flytting? | Samordne med #7 og prøv ansvarsmapping med ukjente deler. |
| Q06 | Hvordan oppdages kode/konsumenter som ingen har annotert? | Bruk et uavhengig kode-/konfigurasjonsinventar og bevisst utelatt konsument. |
| Q07 | Hvordan holdes oppdragsutsnitt korte uten å skjule nødvendig kontekst? | La en ny agent forklare og planlegge samme endring fra ulike utsnitt. |
| Q08 | Hvilke properties skal være maskinelt bindende? | Klassifiser hver regel etter egnet kontroll og uunngåelige ukjente forhold. |
| Q09 | Hvordan bevares lærdom uten at en uvalgt løsning blir gjeldende design? | Følg én læringspost fra Concept1 til MVP1-disposisjon. |
| Q10 | Hvilke standardbegreper adopteres, og hvilke er SDP-spesifikke? | Lag en liten standardmapping med eksplisitte avvik og fulltekstsjekk ved behov. |

## 19. Foreslått neste avgrensede undersøkelse

Dette er et forslag til senere oppdrag, ikke en allerede startet implementasjon.

**Formål:** Test om vi kan beskrive én liten endring slik at eier, Worker og Verifier ser samme sammenheng, og at verktøygrunnlaget viser minst ett bevisst avvik.

**Case:** Bruk kapittel 14 som tydelig hypotetisk fixture, eller en separat bekreftet produktendring. Begrens til én Feature, én Functionality-kandidat, én Channel og høyst seks ansvarsenheter. Ikke modellér hele Concept1 → MVP1 først.

**Tre sammenligningsspor**, i tråd med P04s forslag:

1. Ett konkret eksisterende arkitektur-DSL, for eksempel Structurizr eller LikeC4.
2. Ett rikere modellspor, for eksempel CML eller SysML v2.
3. En liten eksplisitt modell i JSON/CUE eller tilsvarende, med egen begrenset SDP-semantikk.

Før utføring må konkrete verktøy, versjoner og bruksvilkår velges. Alle spor skal bruke samme behov, B/T, forutsetninger og vurderingskriterier. Skill innebygd støtte fra metadata-konvensjoner, manuelt arbeid og ny kode.

**Forventede leveranser:** før/etter-modell; total- og objekt-diff med rest; ansvarsmapping; ett Worker- og ett Verifier-utsnitt; læringspost og fasebeslutning; samt en kort evaluering av vedlikehold, forståelighet og kontrollgrenser.

**Negative prøver:** brutt referanse; rename uten ansvarsendring; split/move med delvis bevart atferd; feil enhet med uendret datatype; en konsument utelatt fra modellen; kode som bryter en tillatt-avhengighetsregel; manglende runtime-registrering; og utsnitt basert på gammel modellrevisjon. Prøvene kan være manuelle i første sammenligning, men må merkes som det. Et spor skal ikke få kredit for en automatisk kontroll som bare er håndlaget i rapporten.

Vurderingen skal måle hvor mye som må redigeres, hva som faktisk oppdages, hvilke opplysninger som går tapt og hva som fortsatt er ukjent. Registrer kostnad ved én normal endring og én revisjon av selve modellskjemaet.

**Stopp:** Lever sammenligningen og anbefalt neste beslutning. Ikke frys språk/schema, bygg generell compiler eller migrer produktrepoer som en skjult del av undersøkelsen. P04s forslag om maksimalt én reparasjonsrunde per verktøyspor kan videreføres; verktøyproblemer dokumenteres dersom de ikke løses innenfor rammen.

## 20. Regler for agenter som viderefører arbeidsdokumentet

1. Bevar skillet mellom eierintensjon, eksisterende aksept, undersøkte fakta og egne forslag.
2. Start fra gjeldende dokumentrevisjon og relevante repo-/issue-kilder. Ikke behandle chatminne eller gamle statusformuleringer som dagens autoritet.
3. Studer et konkret spørsmål og registrer spørsmål, metode, kilder, resultat, begrensning og anbefaling.
4. Ikke gjør Functionality, fasegrenser, en DSL eller et teknologiønske normativt uten en eksplisitt beslutning.
5. Bruk daterte primærkilder og konkrete standardutgaver. Skill katalog-/sammendragslesing fra full normativ gjennomgang.
6. Bevar opprinnelig lærdom og beslutningshistorikk; bruk revisjon og supersession når en oppfatning endres.
7. Beskriv negative resultater og blindsoner. Ingen modell-, test- eller reviewrapport skal påstå bredere dekning enn evidensen gir.
8. Samordne metodeforslag med #9 og #7/#8; bevar #5s grunnlag og eksisterende reviewfunn.
9. Hold eksempler adskilt fra produktfakta. Produktbaseline må bekreftes før den modelleres som dagens system.
10. Avslutt hvert forskningsoppdrag med en konkret beslutning eller et avgrenset nytt spørsmål, slik at utforskningen ikke blir en ny endeløs iterasjonsløkke.

## 21. Revisjonslogg og kontroll av denne leveransen

| Revisjon | Dato | Endring / beslutning |
| --- | --- | --- |
| 0.1 | 2026-09-14 | Renskrevet mandat, forhold til eksisterende SDP, første standardoversikt, foreslått fase-/diff-/evidensmodell og avgrenset videre studie. Ingen ny metode-/språkaksept registrert. |

Kontrollomfang: lest lokale indeks-/syntesedokumenter og relevante issue-/PR-opplysninger; undersøkt offentlige standardkataloger og de tekniske primærkildene lenket i studien; kontrollert lokale Markdown-lenker og kodeblokker. Ingen compiler, parser eller produktintegrasjon er implementert/testet her. Ingen full standardkonformitetsvurdering eller uavhengig reviewgodkjenning påstås.
