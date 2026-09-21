# SDUI ↔ SDL-runtime — foreslått kontrakt

Status 2026-09-21: Go er valgt implementasjon; typede Go-grensesnitt er første port.
Ingen C-ABI eller FOX-avhengighet kreves. Dette erstatter tidligere teknologivalg.

Kontraktgrunnlag: grunnlag for den nye kontraktleveransen i
[målarkitekturen](target-architecture.md) og G3/G4 i [planen](implementation-plan.md).
Identitet, livstid og hendelsesregler kan gjenbrukes, men gammel kilde-/wireform
er ikke et kompatibilitetskrav. Ingen runtime er implementert her.

**ID:** SDUI-RUNTIME-001 · **Status:** designforslag, ikke implementert runtime/API.
Parseren produserer bare Reference og Connection. Det finnes ingen SDL-loader,
callbackutfører eller generell widgetoppdatering i denne leveransen.

## 1. Oppkobling og to retninger

```text
ref: sdlFile "some_SDL_file.sdl";
# callback=sdlFile.input1_sdl.@callback
sdlFile.input1_sdl.setHandle(BoxUIDefinition.top.rightTop.input1_boxui);
```

`input1_sdl` eies av SDL-modulen, `input1_boxui` av UI-instansen. Aliaset er en
modulreferanse, ikke import av en Markdown-fil. Callback er en symbolsk
medlemsreferanse, ikke en minneadresse i AST.

Framtidig rekkefølge:

1. Parse/valider SDUI og opprett en UI-instans med widgetregister.
2. Vertsadapter løser modulkilder relativt til dokumentets avtalte base og policy.
3. SDL-adapter løser objekter/medlemmer og kontrollerer typer/signaturer.
4. Utfør deklarerte setHandle-oppkoblinger og registrer callbacks.
5. Publiser en fullstendig ramme; aktiver handlinger når bindingene er klare.

Manglende runtime eller binding skal gi eksplisitt unbound/diagnose, ikke en
påstand om at koden er kjørt. Initialisering skal være adskilt fra vanlig rendering
og PDF-eksport. Parseren skal aldri få ansvaret for modulopprettelse eller kjøring.

## 2. Logisk widgetreferanse

Foreslått identitet: session, definitionInstance, widgetInstancePath, generation.
Statisk 0.2-oppslag bruker definisjon og navngitte komponentforeldre; anonyme
grupper gir ikke offentlige banesegmenter. Navn gir
ikke alene en gyldig runtime-referanse: én definisjon kan senere ha flere instanser.

Handle er typet og vertseid. Det overlever kompatibel omplassering; sletting,
dokumentlukking og inkompatibel ny instans invaliderer det. Det er aldri en
Fyne-widgetpeker, FOX-peker eller DOM-node. setHandle må kontrollere
forventet widgetkapabilitet. Profilens lokale validering kan bare kontrollere at
målet er en deklarert widget, ikke SDL-objektets faktiske type.

## 3. Brukerhendelser og UI-oppdateringer

| Retning | Foreslått payload | Regel |
| --- | --- | --- |
| UI → runtime | instans/generation, widget, binding, event-ID, context/value revision, activate/commit og typed verdi | Kontroller aktuell binding, type og enabled før dispatch |
| Runtime → UI | instans/generation, batchrevision, typede property-endringer | Valider hele batchen og publiser atomisk på vertens UI-tråd |
| Runtime → UI | accepted/rejected/pending/unknown + samme command-ID | Resultat er ikke automatisk en ny domeneobservasjon |

Input skiller label (`text`), akseptert verdi (`value`) og brukerens draft. Det
framtidige SDL-API-et kan tilby `input1_sdl.text(...)` og `.value(...)`, men må
oversette dette til validerte UI-/verdibindinger med avtalt eierskap. Programmatisk
oppdatering skal ikke automatisk kalle samme input-callback igjen.

Enter sender draft; Escape gjenoppretter akseptert verdi. En ekstern oppdatering
mens feltet er skittent skal ikke stille overskrive brukerens arbeid. Revisjon og
konfliktpolitikk må defineres før kjørbar integrasjon. Egenskapsoppdateringer som
endrer tekstmål må utløse ny layout, men bevare kompatibel fokus/draft.

`svg`-produsenten publiserer en ressurs med identitet, revisjon, størrelse og inert
SVG-innhold. Den kjøres ved en avtalt hendelse/oppdatering, aldri som vilkårlig
callback fra malerens tegneoperasjon. Budsjett, avbrudd og SVG-validering må bevares.

## 4. API, ABI og felles IR

En logisk kontrakt trengs selv om SDUI og SDL senere deler runtime/IR. En binær
ABI trengs bare der faktisk implementasjon krever det, for eksempel en native
C++/Rust-grense. Ikke innfør en ekstra binærprotokoll uten en konkret konsument.

Foreløpig anbefaling: behold språkenes syntaks adskilt og normaliser SDUI til en
UI-modell. Undersøk felles IR når SDL har eksplisitt kjørbar semantikk. Dagens
XFMD parse/prepare/free-ABI er en renderergrense og er ikke denne runtime-kontrakten.

## 5. Avklaringer før implementasjon

Signaturer for callbacks og setHandle; eierskap for domeneverdier; synkron/asynkron
utførelse; thread-affinity; cancellation/teardown; feil og duplikater; modulreload;
handle-revokering; skjemaversjoner; ressursbudsjetter; hvordan en prototype blir
uttrykkelig aktivert. Det er ikke parserens jobb å gjette disse egenskapene.
