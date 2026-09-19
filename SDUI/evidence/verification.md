# SDUI 0.1 — verifikasjon 2026-09-19

Status: avgrenset parserbevis. Ingen renderer, SDL-kjøring eller XFMD-integrasjon
hevdes. Filhashene i [source-manifest.json](source-manifest.json) identifiserer
prototypen uavhengig av lokale endringer i resten av repositoryet.

## Utført

Fra SDP-worktreet:

```sh
PYTHONPATH=SDUI/src python3 -m sdui SDUI/examples/main-page.sdui \
  -o SDUI/examples/main-page.ast.json
PYTHONPATH=SDUI/src python3 -m unittest discover -s SDUI/tests -v
```

**15 av 15 tester bestod**; [full testlogg](parser-tests.txt). Kjørt med miljøet i
manifestet. Tester bruker standardbibliotek og midlertidige filer, ingen nettverk
eller refererte SDL-filer.

Kontrollert:

- Ytre boks, to underbokser i gruppe, SVG-boks, widgetrekkefølge og eksplisitte rader.
- Callbackmål og setHandle som forskjellige AST-noder med nøyaktig identitet.
- Unicode/escapes/kommentarer, CRLF og UTF-8-byteområder.
- Eksakt profilversjon, duplikatnavn, ukjent modul/definisjon/widget og feil property/type.
- Ingen vilkårlig kode, feil @-plassering, trailing commas eller tomme widgetrader.
- Byte-, token-, dybde-, node- og argumentgrenser.
- Alle trunkeringer av minimalfixture og 300 deterministiske støystrenger gir
  kontrollert parserresultat/diagnose, ikke uventet Python-feil.
- Reproduserbar JSON-AST mot det committbare eksempelartefaktet.
- CLI stdout/stderr, returkode 2/3, ugyldig UTF-8, manglende fil og vern mot å
  skrive output over kildefilen.

EBNF og parser ble gjennomgått mot samme konstruksjoner. EBNF blir ikke automatisk
kompilert; dette er ikke et formelt bevis på full språkekvivalens. Delte syntaks-
og semantikkbegrensninger er dokumentert i språkbeskrivelsen.

## Begrensninger

Kun første diagnose, ingen editor-recovery eller formatterer. Ingen lagring av
kommentarer i AST, ingen persistent handle-runtime og ingen ekstern symbolkontroll.
Ingen geometri eller visuell lesbarhet er testet. Ingen støtte for å laste en hel
Markdown-fil eller eksisterende BoxUI-JSON gjennom denne CLI-en. Python 3.11 er
profilens valgte minimum; faktisk kjørt Python-versjon står i manifestet. Ikke hev testing på andre versjoner enn registrert miljø.

## Bevaring og leveranse

Bare `SDUI/` er lagt til i SDL-utviklingsworktreet. Eksisterende modifisert
SystemDesignLanguage/README.md og utrackede SDL-arbeidsdokumenter er urørt.
Ingen XFMD-/Mermaid-kode, installasjon, branchbytte eller ekstern publisering
inngår i denne leveransen. Arbeidet er lokalt og klart for gjennomgang.
