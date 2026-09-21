# Statisk Markdown-dump — verifikasjon 2026-09-21

**30 av 30 tester bestod**, inkludert seks nye Markdown-eksporttester.
[Testlogg](markdown-dump-tests-2026-09-21.txt),
[kode-/miljømanifest](markdown-dump-manifest.json),
[generert Concept1-dump](../examples/concept1-bucking.dump.md).
Dette er gjeldende eksportbevis etter den tidligere frontend-/tekstdumpleveransen;
eldre manifest gjelder kildeversjonen før Markdown-eksporten ble lagt til.

Fra SDUI-katalogen:

```sh
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format markdown --entry bucking -o examples/concept1-bucking.dump.md
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

Kontrollert: Markdown beholdes som renderbart innhold, Mermaid utelates,
andre kodeblokker bevares, åpne gjerder avgrenses ved widgetgrensen, backticks
i innhold/etiketter kan ikke lukke eksportens omsluttende gjerde, skjulte noder
utelates, region-/radrekkefølge og nesting bevares. CLI krever profilvalidering,
gyldig entry og tilstrekkelig oversiktsbredde; feil gir ingen delvis stdout.
Lagret Markdown er byteidentisk med ny CLI-kjøring. `git diff --check` bestod.

En separat midlertidig venv med markdown-it-py 4.2.0 rendret dokumentet til HTML
med CommonMark og tabellutvidelsen. Kontrollen fant h2-overskriften Lengde,
fet måleverdi, tabellcellen Sagtømmer, kodeetiketten Cursor AV og nestede
blockquote-elementer. Ingen button/input-elementer eller Mermaid-kodeblokker
ble generert. Dette er strukturell Markdown→HTML-verifikasjon, ikke visuell
XFMD-/nettleser-/PDF-verifikasjon. Ingen ny produksjonsavhengighet er lagt til.

Ren Markdown har ikke en generell mekanisme for SDUI-kolonner og vekter.
Eksporten viser derfor terminaloversikten først og renderbart innhold i
leserekkefølge etterpå. Ingen interaktivitet, callbackkjøring, ny parser,
rendererendring i Mermaid/XFMD, commit eller push.
