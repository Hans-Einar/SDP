# Checkpoint #1 — generert G1–G5-design

Dato: 2026-09-22. V2–V4 er levert for gjennomgang før Go-implementasjon.
**Start med [den genererte implementasjonsrapporten](../../SDUI/design/viewpoints/implementation.md).**
Den kommer fra [SDL-kilden](../../SDUI/design/architecture.design), via SDLs
parser, validator og viewpoint-verktøy. Faser, ansvar, avhengigheter og
scenariofigurer er ikke håndskrevet inn i rapporten.

## Hva som kan gjennomgås

Fem planlagte faser, 18 milepæler og eksplisitte milepælkoblinger til alle 94
Functionality-er beskriver frontendport, layout/presentasjon, UI-runtime/reload,
SDL-runtime/binding og native Go-generering. Alle har status **planned**.
Rapporten viser logiske eiere og kilde-ID-er for ansvar og avhengigheter.

[Den samlede utskriften](../../SDUI/design/viewpoints/printout.md) inneholder
alle elleve viewpoints: bruksmål/bidrag, arkitektur, ansvar, porter, modi,
arbeidsplan, Feature-allokering, Channel-sekvenser, data/kontrakter, packet og
faktaregister. [Mermaid-utgaven](../../SDUI/design/viewpoints/viewpoints.md)
og de enkelte diagramkildene genereres samtidig.

Ti scenarioer illustrerer kompilering, statisk/interaktiv presentasjon,
ubundet Go-handling, akseptert/avvist binding, UI-/SDL-reload og native bygg.
47 [MessageSet-oppføringer](../../SDUI/design/viewpoints/message-sets.json)
avledes fra Channel, modus, tillatte meldinger og deltakere.

## Leveranse og verifikasjon

| Fase | Levert verktøyfunksjon |
| --- | --- |
| V2 | Typede data-/feltkontrakter, persistent Database-begrep, Datagram-projeksjon og eksplisitt packet-layout |
| V3 | Channel-roller, tillatte meldinger, ordnede scenarioer, request/resultat-korrelasjon og avledet MessageSet |
| V4 | Sporbar G1–G5-plan i SDL, generisk implementasjonsrapport og samlet designutskrift |

Aktiv SDL-profil er **design-core 0.5**. SDUI er fortsatt **0.2**; numrene
tilhører forskjellige språk. Erstattede SDL-profiler har ingen aktiv fallback.
Checkpointets øvrige kandidater og MVP1-korpus er ikke dermed vedtatt grammatikk.

[Maskinrapporten](../../SystemDesignLanguage/tools/verification.json) registrerer
61 SDL-parsertester, 24 verktøytester og 36 SDUI-tester: **121 bestått**.
Modellen har 368 deklarasjoner og 1106 fakta og gir 142 SVG-diagrammer.
Verifikasjonen kontrollerer kildeposisjoner, diagramfakta, SVG-etiketter,
bitområder, sekvensrekkefølge/korrelasjon, hasher og byte-identisk reeksport.
Fasegraf, parsersekvens og bindingssekvens er visuelt stikkprøvekontrollert.
Ingen fysisk utskrift eller paginert PDF er testet.

Fra repoets rot kan hele kontrollen gjentas:

```sh
python3 SystemDesignLanguage/tools/verify_design.py --phase V4 --renderer /home/warloc/git/mermaid-rs-renderer/target/debug/mmdr
```

Rendererstien velger et eksisterende program. Leveransen endrer verken
mermaid-rs-renderer eller XFMD og legger ingen SDUI-parser i dem.

## Grenser før G-fasene

Det finnes fortsatt ingen Go-parser eller Go-runtime i denne leveransen.
Scenarioene er validerte designbaner, ikke utført domenelogikk. Tokens, AST,
normalisert modell og prepared frame er opake bytesfelt ved scenarioets grenser;
interne Go-typer, funksjonssignaturer, atomisk publisering, state-migrering,
widgetlivstid og font-/måleenheter må konkretiseres og testes i G-fasene.

VP07 rapporterer 15 manglende allokeringer i Feature-/modusutsnitt. Modellen er
derfor ikke en fullstendig deploymentplan, selv om alle ansvar har milepælkobling.
Packet-eksemplet er en eksplisitt prototypeprofil, ikke en vedtatt runtime-ABI
eller et P1000-/StanForD-format. Database betyr persistent datatilgang ved behov,
uten krav om SQL.

Arbeidet følger [branchstakken](../Development-Branch-Stack.md), med tre commits
per V2–V4-fase og push ved faseavslutning. Samlet PR har `sdp-vNow` som base;
merge og oppstart av G-fasene inngår ikke i denne leveransen.
