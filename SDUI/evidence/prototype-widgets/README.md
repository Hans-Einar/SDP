# Prototype-widgets — verifikasjon 2026-09-21

R24 er levert som SVG-referansekomposisjon og lokalt HTML-widgetgalleri.
Ingen generell layout, FOX, SDL-runtime eller ABI er implementert i denne prøven.
Oppdatert samme dag: kontrollene for Lengde, Diameter og Stammeforløp er flyttet
til overskriftsradene; Diameter viser nå også Δ. Kilde, AST og begge strukturelle
dumper er regenerert. SVG-en er rendret med librsvg og inspisert på nytt uten
overlapp mellom titler og kontroller. Manifestet og testloggen gjelder denne utgaven;
nettleserbeviset gjelder det uendrede HTML-galleriet.

## Kontroller

- 36/36 Python-tester består; [rå logg](tests.txt). Seks nye tester dekker modell-
  og kildegjenbruk, utelatte callbacks, synlighet/disabled-arv, SVG-tilstander,
  escaping, utskriftsblokker og de genererte SVG-artefaktene.
- Nettleserprøve med Playwright og lokal Google Chrome: redigering, synlig musetrykk,
  tastaturaktivering, deaktivert knapp, aktuelle utskriftsverdier, linjebryting,
  smal visning, tilbakestilling ved reload, ingen eksterne forespørsler og ingen
  JavaScript-feil. [Maskinresultat](browser-verification.json).
- [Skjerm med knapp nedtrykket](browser-pressed.png) og
  [utskriftsvisning med utfylte verdier](print-preview.png) er visuelt inspisert.
  Utskriftskontrollen emulerer nettleserens print-CSS. Fysisk utskrift,
  paginert PDF og XFMD GUI er ikke verifisert her.
- Begge SVG-er er rendret med lokal librsvg og visuelt inspisert: kontroller og
  innhold er synlige uten overlapp. SVG-en er et bilde; ingen klikk eller redigering.
- Generering gir identiske bytes ved gjentakelse. [SHA-256-manifest](manifest.json)
  identifiserer kilde og artefakter fra denne prøven.

## Reproduksjon

Fra SDUI, standardbibliotek for produksjon og enhetstester:

```sh
python3 tools/build_widget_previews.py
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

Den separate nettleserprøven trenger Python-pakken `playwright` og
`/usr/bin/google-chrome`; dette er verifikasjonsverktøy, ikke produksjonsavhengigheter:

```sh
python3 evidence/prototype-widgets/verify_browser.py
```

Den åpner det lokale HTML-galleriet og regenererer de to skjermbildene og
nettleserrapporten. Ingen server eller eksterne nettressurser kreves av demoen.

Se [bibliotekets omfang](../../docs/prototype-widgets.md) for begrensninger,
XFMD-kildekontroll og videreføring til felles geometri/FOX.
