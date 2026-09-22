# SDUI — ` page `

Statisk GUI-dump. Knapper og felt er tekstetiketter; ingen callbacks kjøres.

## Layoutoversikt

Rad-/kolonnestruktur i terminalceller. Høydene følger innholdet; dette er ikke målt GUI-geometri.

```text
SDUI GUI dump | page | 120 columns | structural preview
+-page-----------------------------------------------------------------------------------------------------------------+
|## EditAptCell · simulert domene                                                                                      |
|Celle LengthA1 · eksempelverdi 400. Prøven godtar heltall 100–1000.                                                   |
|Dette er ikke Ponsse-domenevalidering eller maskinkontroll.                                                           |
|Lengdeverdi: [440]                                          [ Bruk verdi ]                                            |
|Enter eller knappen sender draft og forventet domenerevisjon. Escape forkaster lokal redigering.                      |
+----------------------------------------------------------------------------------------------------------------------+
```

## Innhold

Markdown gjengis som innhold. Nestede sitatblokker viser grupper og frames. Horisontale søsken står i leserekkefølge her; plasseringen vises i oversikten. Mermaid-diagrammer er utelatt.

> **BoxUI-frame:** ` page `
>
> **header**
>
> > ## EditAptCell · simulert domene
>
> **Rad 1 · 1 komponent fra venstre mot høyre**
>
> > Celle LengthA1 · eksempelverdi 400. Prøven godtar heltall 100–1000.
> > Dette er ikke Ponsse-domenevalidering eller maskinkontroll.
>
> ---
>
> **Rad 2 · 2 komponenter fra venstre mot høyre**
>
> **Inndata:** ` Lengdeverdi ` — ` 440 `
>
> ---
>
> **Knapp:** ` Bruk verdi `
>
> ---
>
> **Rad 3 · 1 komponent fra venstre mot høyre**
>
> > Enter eller knappen sender draft og forventet domenerevisjon. Escape forkaster lokal redigering.
>
> ---
>
