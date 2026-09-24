# SDUI — ` page `

Static GUI dump. Buttons and fields are text labels; no callbacks execute.

## Layout overview

Row/column structure in terminal cells. Heights follow content; this is not measured GUI geometry.

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

## Content

Markdown is rendered as content. Nested blockquotes represent groups and frames. Horizontal siblings appear in reading order here; the overview shows placement. Mermaid diagrams are omitted.

> **BoxUI-frame:** ` page `
>
> **header**
>
> > ## EditAptCell · simulert domene
>
> **Row 1 · 1 component from left to right**
>
> > Celle LengthA1 · eksempelverdi 400. Prøven godtar heltall 100–1000.
> > Dette er ikke Ponsse-domenevalidering eller maskinkontroll.
>
> ---
>
> **Row 2 · 2 components from left to right**
>
> **Input:** ` Lengdeverdi ` — ` 440 `
>
> ---
>
> **Button:** ` Bruk verdi `
>
> ---
>
> **Row 3 · 1 component from left to right**
>
> > Enter eller knappen sender draft og forventet domenerevisjon. Escape forkaster lokal redigering.
>
> ---
>
