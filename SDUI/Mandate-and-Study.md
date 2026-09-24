# SDUI — mandat og innledende studie

**Gjeldende retning 2026-09-21:** SDL/SDUI-parser og runtime utvikles videre i Go,
med Fyne som første interaktive vert og SVG-eksport fra felles modell/layout.
Dette erstatter eldre språk-/vertsteknologivalg nedenfor; avsnittene er mandatets
historie. [Checkpoint tillegg 07](../SDP/History/checkpoint-1/07-SDUI-0.2-and-Go-Direction.md)
og [PLAN-003](docs/implementation-plan.md) eier dagens leveranser.

**Implementasjon 2026-09-22:** Eierens Concept1-bestilling er konkretisert som
[SDUI-kilde, AST og konsolldump](docs/concept1-console.md). Go-frontenden, layout/SVG, Fyne, runtime og modellreload er levert.
[Checkpoint tillegg 11](../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md)
avgrenser SDL-kjøring, generering og dokumentnavigasjon. Eldre avsnitt
nedenfor bevarer mandatets historie, ikke en ekstra aktiv parserprofil.

**Presisering 2026-09-21:** Kanoniske hjørner er ^< og >^ øverst, v< og >v
nederst. Ytterframe `{16:9,<->}` bruker hele layoutbredden og avleder høyden;
den fyller rollen til FixedAspectViewport uten høydebegrenset contain-fallback.
Header/footer kan være frames eller widgetinnhold; `body=` kan være eksplisitt.
Eierens mainBody/page-eksempel innfører navngitte gjenbrukbare widgetgrupper.
`font=10/12` er absolutte tekststørrelser; vindusresize skalerer ikke innholdet.
Hver UI-komponent kan ha formatering `{...}` før separatoren. Komma fortsetter
horisontalt, semikolon starter en ny rad under den foregående. Detaljer står i
[komposisjonsforslaget](docs/frame-composition-proposal.md).

**Ytterligere eierpresisering 2026-09-20:** Frame-sideforhold uttrykkes som x:y.
Scale er relativ til nærmeste ancestor, for root til vertens layoutområde.
Med ratio styres bare én akse; ellers kan x, y eller begge skaleres. Ingen
pikselbredde/-høyde i ny kilde. Mini arrows har kanonisk shape uten at rekkefølgen
på hjørneparets retninger endrer betydningen; ned/venstre skrives v<.
Tegnet ¤ er en idé uten fastlagt betydning. Se det oppdaterte layoutforslaget.

**Tillegg 2026-09-20:** Eieren autoriserer videreutvikling av libsdui,
SDUI-runtime og layout/presentasjon og ber om phases/milestones. Generelle frames,
widgetlister, Markdown-innhold og layoututtrykk videreutvikles i
[layoutforslaget](docs/layout-language-proposal.md). Ingen støtte for gammel
SDUI 0.1/BoxUI-kontrakt kreves: port eksempler og fjern erstattede kjøreveier.
[Implementasjonsplanen](docs/implementation-plan.md) er gjeldende videreplan.
Resten av dokumentet bevarer den opprinnelige parserleveransens grunnlag.

**ID:** SDUI-MANDATE-001 · **Revisjon:** 0.1 · **Dato:** 2026-09-19.
Eierens bestilling autoriserer dokumentasjon, språkdefinisjon i EBNF og en parser
som genererer AST. SDUI er et arbeidsnavn; ingen navnekollisjonsundersøkelse eller
endelig navneregistrering er gjort.

## 1. Eierens hensikt

Et UI skal kunne beskrives som en stor ytre boks med navngitte underbokser og
widgetinnhold, slik Ponsse/Concept1s BoxUI er organisert. Plasseringen skal være
forutsigbar og skjemastyrt. Et lett språk inne i Markdown skal kunne beskrive
layout og koble widgets til objekter/funksjoner som er definert i en SDL-fil.

`input1_boxui` er widgeten. `sdlFile.input1_sdl` er SDL-objektet. Callback går fra
brukerhendelsen til SDL-objektet; `setHandle(BoxUIDefinition.input1_boxui)` gir
SDL-objektet en logisk referanse for senere oppdateringer. Dette er to retninger,
ikke en rå FOX-peker eller en callback som tegner direkte på skjermen.

Eieren ønsker et avgrenset eget språk og støtter retningen om å skille det fra
Mermaid. Denne leveransen flytter ikke allerede implementert kode eller endrer
XFMD-installasjonen. Et senere integrasjonsløp skal gjenbruke eksisterende arbeid.

## 2. Undersøkt grunnlag

| Kilde | Faktisk observasjon | Konsekvens |
| --- | --- | --- |
| Ponsse `882ad7c`, Concept1/shared/ui-box/model.mjs og UILayout.jsx | group/axis/weight/box-id; CSS-grid fordeler spor. Innhold leveres separat som React-komponenter. | Bevar skillet mellom plassering og innhold. Fullt UI er ikke allerede ett portabelt deklarativt skjema. |
| Concept1/apps/operator-ui/src/ui/operator-layout.mjs | Nestede grupper med 15/45/40 vekt og stabile boks-ID-er. | Realistisk framtidig kompatibilitetsfixture. |
| Rendererfork `61a85b6`, src/boxui/layout.rs, svg.rs, model.rs | BoxUI har selvstendig row/column-plassering, tekstcallback og SVG; bruker ikke treemap-layout. | Kjernen kan undersøkes for uttrekk; det finnes ingen nødvendig treemap-avhengighet. |
| XFMD `ffb98e8`, BoxUiAbi, BoxUiFrame, BoxUiSession | Parse/prepare/free-ABI, typed modell, native kontroller, lokale simulerte deltakere. Ingen SDL-runtime. | Gjenbruk identitet/revisjon/utkast; ikke kall dagens renderer-ABI en SDL-ABI. |
| SDL-arbeidsdokumentet i søsterkatalogen | Utforskende designmodell. Modellinterpreter trenger ikke kjøre produktfunksjonalitet. | Ikke forutsett en ferdig kjørbar SDL-runtime eller vedta dens semantikk her. |

SDL-arbeidsdokumentet har lokale, ucommittede endringer. Vi behandler det som
arbeidsgrunnlag, ikke som en frosset godkjent kontrakt. Ingen nye eksterne
biblioteker er valgt i denne parserleveransen.

## 3. Avgrensning og beslutning

Velg én avgrenset kildeprofil, ett eksplisitt AST-format og separat lokal
validering. Python-standardbibliotek er valgt for en liten, inspiserbar
referanseprototype. Dette vedtar ikke implementasjonsspråk for den framtidige
UI-kjernen eller vertsadapteren.

Prototypen støtter bokser, rader med text/button/input/svg, modulreferanser og
deklarativ setHandle. Den kjører ingen av disse koblingene. Heller ikke syntaktisk
korrekte referanser er bevis på at SDL-objektet eller medlemmet finnes.

Et eget språk gir kortere uttrykk enn wire-JSON, men krever tydelig grammatikk,
diagnoser, versjonering og tooling. UI-modellen skal senere være felles inngang
for SDUI og relevante eksisterende skjemaer; det er en plan, ikke en ferdig adapter.

## 4. Åpne spørsmål

- Endelig navn og profilnavnerom; behold SDUI som arbeidstittel.
- Faktiske SDL-funksjonssignaturer, objektopprettelse, asynkronitet og feilmodell.
- Første mål for UI-kjernen etter uttrekk: Rust eller annen implementasjon.
- Detaljert størrelse-/overflowprofil og hvor tett Concept1-geometrien skal samsvare.
- Om senere SDUI normaliseres til en delt IR med SDL. Det behøver ikke endre kildeprofilen.

Disse spørsmålene blokkerer ikke parsing til et kildetro AST. De blokkerer påstander
om kjørbar oppkobling og kompatibel UI-rendering.
