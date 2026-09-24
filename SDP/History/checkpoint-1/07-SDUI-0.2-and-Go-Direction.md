# Checkpoint #1 — SDUI 0.2 og felles Go-retning

**Implementasjonsstatus er oppdatert i [tillegg 11](11-Go-Implementation-and-Navigation.md).**
Nedenfor beholdes det daterte design-/V-fasegrunnlaget; gamle Python-kommandoer
er historiske og erstattet av Go-inngangene.

Oppdatert 2026-09-21 etter eierens valg i samtalen. Dette tillegget korrigerer
implementasjonsretningen fra 18.–20. september; checkpointnummeret er uendret.
SDL omtales nå som **SystemDesignLanguage**. Eldre «System Description Language»
er historisk navnebruk, ikke et annet språk. SDUI er et selvstendig UI-designspråk.
Ingen av disse navnene endrer Toolkit-kontraktene for Standard Document Procedure.

**Tillegg 2026-09-22:** [Parser-/runtimedesignet](../../../SDUI/design/README.md) er nå
beskrevet samlet i SDL med AST og valideringsrapport fra eksisterende parser.
Dette konkretiserer Go-retningen strukturelt; Go-implementasjonen gjenstår.

## Gjeldende beslutninger

| ID | Valg | Konsekvens |
| --- | --- | --- |
| CP1-D18 | SDUI 0.2 med frames, grupper, Markdown og relativ layout | Python-frontenden er portert; ingen aktiv 0.1-kompatibilitet |
| CP1-D19 | Go for videre SDL/SDUI-parser og runtime | Erstatter den planlagte Rust/C-ABI-kjernen; Python er portgrunnlag til erstatningen virker |
| CP1-D20 | Selvstendig UI med Fyne som første interaktive vert | FOX/XFMD er ikke nødvendig for kjernen; ingen parallell FOX-/TUI-implementasjon nå |
| CP1-D21 | Felles UI-modell/layout for interaktiv visning og SVG-eksport | Markdown kan inkludere genererte UI-bilder; det gjør ikke bildet interaktivt |
| CP1-D22 | Hot reload av validert SDL/SDUI-modell | Bevar siste gyldige modell; kompatibel tilstand kan videreføres, inkompatible endringer trenger reset/migrering |
| CP1-D23 | Senere SDL/SDUI→Go-generering med håndskrevne domenefunksjoner | Generert og håndskrevet kode holdes adskilt; samme runtime brukes i utvikling og ferdig program |

D19–D23 er valgt utviklingsretning, ikke påstand om implementert Go-funksjonalitet.
SDLs kjørbare delmengde, eksakte Go-grensesnitt og reload-regler må fortsatt defineres.
Eldre kandidater om ControlSet, Function/Functionality og Activity blir ikke
vedtatte språkregler som følge av teknologivalget.

## Hva finnes faktisk?

| Område | Nåstatus |
| --- | --- |
| SDL | Python `design-core 0.5`: struktur, data, Channels, scenarioer og planfakta; generert G1–G5-design i tillegg 09 |
| MVP1 | 66-fils forfatterøvelse i kandidatprofil; ikke kjørbar med design-core |
| SDUI | Python-parser for 0.2, AST med kildeposisjoner, validering og statisk normalisering |
| Konsoll/Markdown | Genererte strukturdumper, rå Markdown i konsoll og separat renderbart innhold i Markdown |
| SVG | Tegnede button/input og avgrenset apteringskomposisjon; plassering fra lagret treemap-geometri og fixturekode |
| HTML | Lokalt widgetgalleri med inndata, knappetrykk og utskriftsvisning; ingen SDL-kjøring |
| Go | Kataloger for parser/runtime er opprettet; ingen Go-kode, moduler, Fyne-integrasjon eller kjørbar runtime ennå |

Den eksisterende SVG-prøven er **ikke en generell SDUI-layoutmotor**. Den tegner
utvalgte kildekontroller og eksempeldata. En generell motor må måle og plassere
vilkårlig støttet innhold, inkludert formatering og klipping.

[SDUI-testbeviset](../../../SDUI/evidence/prototype-widgets/README.md) dokumenterer
36 beståtte tester og nettleserprøven. Dette er ikke bevis for SDL-runtime,
Fyne, generell Markdown/Mermaid-komposisjon eller full systemkjøring.
Eldre hashmanifest og testlogger gjelder sine daterte baselines.

## SDUI 0.2 på én side

- `[]` er frame; `*box` eller aliaset `*b` velger boksdekorasjon.
- `<>` grupperer widgets/Markdown; gruppene kan nestes og få egen formatering.
- `header=`, valgfri `body=` og `footer=` kan inneholde sammensatte komponenter.
- Komma fortsetter horisontalt; semikolon starter neste rad. `{...}` står etter
  komponenten og før separatoren.
- Strenger som innhold er Markdown. `button`, `input` og `svg` er widgetkall;
  widgetetiketter er vanlig tekst. `svg`-referanser kjøres ikke av parseren.
- Dimensjoner er relative til nærmeste kildeancestor; root bruker vertens område.
  Ingen pikselbredde/-høyde i kilden. Fontstørrelse er absolutt.
- `16:9` låser frame-ratio. Én akse kan styre skalering/fylling med ratio.
  `{16:9,<->}` fyller bredden og avleder høyden; ingen automatisk contain-fallback.
- Hjørner skrives kanonisk `^<`, `>^`, `v<`, `>v`; omvendt tegnrekkefølge har
  samme betydning. `¤` har ingen definert betydning.
- `ref`, callback og `setHandle` er symbolske AST-data, ikke utførte kall.

```text
sdui 0.2;
mainBody = <name=input("Navn", value="Ola"), ok=button("OK")>;
page = [header="## Prototype", body=mainBody, footer="Uten domenekobling"]*b {16:9,<->,font=12};
```

Detaljer vedlikeholdes i [språkprofilen](../../../SDUI/docs/language.md),
[EBNF](../../../SDUI/grammar/sdui-0.2.ebnf) og
[layoutforslaget](../../../SDUI/docs/layout-language-proposal.md).
Parsing av en layoutregel er ikke bevis for utført geometri.

## Implementasjon og hot reload

SDUI-parseren lager modell; runtime håndterer identitet, tilstand og hendelser;
layout måler/plasserer; Fyne presenterer og SVG-eksport lager dokumentasjonsbilder.
SDL har egen parser/modell og avgrenset runtime. Koblingen går gjennom typede
Go-grensesnitt og registrerte domenefunksjoner, ikke gjennom en obligatorisk C-ABI.
Bibliotekene må kunne testes uten GUI. Fyne må ikke importeres av parser/runtime.

SDUI-endring → parse/valider → ny modell/layout → publisering. Feil i kilde beholder
siste gyldige visning. Felt/fokus kan bevares ved kompatibel identitet og type.
SDL-modellendringer publiseres ved avtalte hendelsesgrenser; tilstandsmigrering og
utestående operasjoner må håndteres eksplisitt. Reload skal ikke gjenta domenekall.

Endret Go-kode bygges og restartes i første løsning. Go-byggcache er ikke
utskifting av maskinkode i en kjørende prosess. `plugin` er ikke reloadfundamentet.
Se [Go-byggcache](https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching) og
[plugin-kontrakten](https://pkg.go.dev/plugin). Senere prosessutskifting krever
et eget målt behov; ingen slik mekanisme er opprettet nå.

Kodegenerering skal først produsere Go som oppretter den samme validerte modellen
og de samme koblingene som utviklingsmodus bruker. Domenefunksjoner leveres som
Go-implementasjoner bak avtalte kontrakter. Ufullstendig SDL-semantikk skal gi en
forklarlig diagnose, ikke tomme «vellykkede» funksjoner eller gjettet maskinlogikk.

## Les videre og neste leveranse

- [Gjeldende målarkitektur](../../../SDUI/docs/target-architecture.md) og
  [felles implementasjonsplan](../../../SDUI/docs/implementation-plan.md).
- [SDUI Go-område](../../../SDUI/go/README.md) og
  [SDL Go-område](../../../SDL/go/README.md).
- [Apteringskilde](../../../SDUI/examples/concept1-bucking.sdui) og
  [SVG-visning i Markdown](../../../SDUI/examples/concept1-bucking.widgets.md).
- [SDLs kjørbare strukturprofil](../../../experiments/design_core/README.md).
- [MVP1-oversikt](../../../experiments/mvp1_sdl/README.md),
  [System.design](../../../experiments/mvp1_sdl/SDL/MVP1/System.design) og
  [EditAptCell](../../../experiments/mvp1_sdl/SDL/MVP1/Scenarios/EditAptCell.design).

Første kjørbare Go-leveranse er en SDUI-frame med knapp/input, felles layout,
SVG-eksport, Fyne-visning og SDUI-reload mot en registrert Go-funksjon. Deretter
kommer én eksplisitt SDL-kobling før videre kjørbar grammatikk/kodegenerering.
Denne oppdateringen leverer bare dokumentasjon, strukturmodell og kataloger.
