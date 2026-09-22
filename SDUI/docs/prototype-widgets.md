# SDUI — begrensede prototype-widgets

**Historisk kontrollgalleri.** Fixturebygger og midlertidige CLI-formater er
fjernet i G5-M4. Bruk [felles Go-layout/SVG](../go/README.md) og
[state-dokumenteksport](../../SystemDesignLanguage/go/README.md) for ny UI.

Levert 2026-09-21 for R24: `button` og `input` kan tegnes som SVG i
[Markdown-prøven](../examples/concept1-bucking.widgets.md). En separat
[HTML-demo](../examples/prototype-controls.html) viser redigering, mus-/tastaturtrykk,
fokus, deaktivert knapp og utskrift av gjeldende feltverdier.

Dette bruker eksisterende SDUI 0.2-syntaks. `prototype_widgets` leser normaliserte
instanser fra samme parser og bevarer identitet, label, value og enabled.
Skjulte grener utelates; disabled arves gjennom foreldre. Callback-referanser
overføres ikke til demoen. Ingen SDL-fil åpnes og ingen SVG-produsent kjøres.

## To presentasjonsformer

| Format | Mulighet | Begrensning |
| --- | --- | --- |
| SVG-bilde i Markdown | Tegnede knapper/felt, normal/trykket/fokus/deaktivert | Statisk bilde; ingen redigering eller klikk |
| Lokal HTML i nettleser | Native HTML-input, lokal knappetrykkrespons, utskrift | Separat widgetgalleri, ingen lagring eller domenefunksjon |

CommonMark har ingen standard for skjema-widgets. [Rå HTML](https://spec.commonmark.org/0.31.2/#raw-html)
kan inngå i kilden, men leseren bestemmer hva som tillates og kjøres.
SVG brukt som bilde har også [begrensninger på interaktivitet og skript](https://developer.mozilla.org/en-US/docs/Web/SVG/Guides/SVG_as_an_image).
Dette er derfor ikke en ny Markdown-utvidelse som andre lesere må implementere.

Kontroll av lokal XFMD-kode: `src/interpreter/ModelBuilder.cpp` behandler rå HTML
som tekst, ikke som DOM-kontroller. `src/application/media/ImageDecoder.cpp`
kopierer et statisk pixbuf til Cairo; README oppgir første bilde for GIF.
SVG-bilder passer dagens dokumentvei, mens skjema og animasjon krever en annen
vertsfunksjon. HTML-demoens CSS/JS kjører i nettleseren, ikke i XFMD.

## Bygg og avgrensning

Fra SDUI:

```sh
python3 tools/build_widget_previews.py
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

Kun standardbiblioteket kreves for generering. Verktøyet leser
`concept1-bucking.sdui`, `prototype-controls.sdui` og tidligere registrert
treemap-geometri. SVG-widgets og HTML-kontroller bruker samme widgetbeskrivelse.
Boksregionene beholder treemap-prøvens rad- og kolonnefordeling; overflødige
gruppeoverskrifter fjernes for å gi kontrollene plass.
Lengde/Diameter har Cursor AV og Δ i overskriftsraden; Diameter har også
O/B · U/B og BarkNOR-kalibrering. Stammeforløp har taperNOR og Mixed i samme
rad som tittelen. SDUI-kilden plasserer disse kontrollene i `header`-regionene.
Selve tegningen av overskriftskontroller skjer i SDUI-prøvens SVG-komposisjon;
det er ikke en ny funksjon i Mermaids treemap-renderer.

Komposisjonen er manuelt plassert innen referanseboksene. Den gjengir utvalgte
kontroller og håndskrevne eksempeldata, ikke alle innholdsblokker i AST-et.
Den tolker ikke generell SDUI-layout, Markdown-innhold eller fontarv.
SVG-koordinater og demofont er presentasjonsvalg; de innfører ikke pikselmål i
SDUI-språket. Trykk/fokus i SVG er illustrasjonstilstand, ikke nye språkattributter.

HTML-demoen har ingen nettverksavhengigheter eller lagring. Inndata kopieres som
tekst til separate utskriftsblokker ved redigering og før utskrift. Lange verdier
kan da brytes uten inputfeltets horisontale klipping. Trykk-/fokusutseende fjernes
ved utskrift. Redusert bevegelse deaktiverer CSS-overganger.

## Videreføring

Dette er portgrunnlag for [Go-planens G1–G3](implementation-plan.md). G2 leverer
felles målt geometri til Fyne og SVG og erstatter fixtureplasseringen. G3/G4
kobler widgetidentiteter og egenskaper til runtime og SDL via typede Go-porter.
FOX er ikke første backend. HTML-galleriets lokale trykkteller er fortsatt bare
en demonstrasjon, ikke en SDL/SDUI-runtime.

[Verifikasjon](../evidence/prototype-widgets/README.md)
