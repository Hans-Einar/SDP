# SDUI — ` bucking `

Static GUI dump. Buttons and fields are text labels; no callbacks execute.

## Layout overview

Row/column structure in terminal cells. Heights follow content; this is not measured GUI geometry.

```text
SDUI GUI dump | bucking | 160 columns | structural preview
+-length-----------------------------------------------------------------------+ +-diameter--------------------------------------------------------------------+
|## Lengde                               [ Cursor AV ]       [ Δ ]             | |## Diameter                            [ Cursor  [ Δ ]     [ O/B · U [ BarkNO|
|**412,0 cm**                                                                  | |                                       AV ]                /B ]      R-kalibr|
|Måling · aggregat                                                             | |                                                                     ering ] |
+------------------------------------------------------------------------------+ |**28,4 cm**                                                                  |
                                                                                 |Måling · O/B · eksempelprofil                                                |
                                                                                 +-----------------------------------------------------------------------------+
+-selection----------------------------+ +-suggestions----------------------------------------------------------------+ +-currentStem--------------------------+
|## Treslag / Sortiment                | |## Apteringsforslag                                                         | |## Stammen i aggregatet               |
|**Treslag**                           | |[ Optimalisering ]                     [ Kolonner ]                         | |Aktuell stamme                        |
|Gran                                  | |Neste kapp · eksempeldata                                                   | |                                      |
|                                      | || Nr | Sortiment | Lengde | Topp-Ø |                                        | |**Låste feil: 0**                     |
|**Sortiment**                         | || -- | --------- | ------ | ------ |                                        | |Ingen låste feil registrert           |
|Sagtømmer                             | || 1  | Sagtømmer | 430 cm | 26 cm  |                                        | |                                      |
|                                      | || 2  | Massevirke| 310 cm | 18 cm  |                                        | |Nåposisjon: 412,0 cm                  |
|Siste knapp: GRAN                     | |[ Canonical ]                          [ Alternativ 1 ]                     | |Maks kvistet: 412,0 cm                |
+--------------------------------------+ |Planvalg er ubundet i denne prototypen.                                     | |Produsert / kappet: 0 cm              |
                                         +----------------------------------------------------------------------------+ |Segmentrot: 0 cm                      |
                                                                                                                        |Produksjon og feilflyt er plassholdere|
                                                                                                                        |.                                     |
                                                                                                                        +--------------------------------------+
+-stemTrack-40fr-----------------------------------------------------------------------------------------------------------------------------------------------+
|## Stammeforløp                                                                 [ taperNOR ]                            [ Mixed ]                             |
|Hele estimatet · Predikert · Målt · Plan · Rest · Kvalitet 1/2 · Underkjent · Utkast                                                                          |
|**Plassholder for OperatorUnifiedStemSvg**                                                                                                                    |
|Stammeprofil, målepunkter, kapp og manuell cursor trenger en egen visuell provider.                                                                           |
|Eksempel på Mermaid-innhold nedenfor utelates av konsolldumpen:                                                                                               |
|[Mermaid utelatt]                                                                                                                                             |
|FØLGER AGGREGAT                                      Lengde: 412,0 cm                                     Rel. lengde: 0,0 cm                                 |
|Pred. Ø O/B: 28,4 cm                                 Profilvolum O/B: 0,312 m³                            Toppsylinder: 0,260 m³                              |
+--------------------------------------------------------------------------------------------------------------------------------------------------------------+
```

## Content

Markdown is rendered as content. Nested blockquotes represent groups and frames. Horizontal siblings appear in reading order here; the overview shows placement. Mermaid diagrams are omitted.

> **Frame:** ` bucking `
>
> **Row 1 · 1 component from left to right**
>
> > **Frame:** ` bucking/top `
> >
> > **Row 1 · 2 components from left to right**
> >
> > > **BoxUI-frame:** ` bucking/top/length `
> > >
> > > **header**
> > >
> > > > **Group:** ` bucking/top/length/header `
> > > >
> > > > **Row 1 · 2 components from left to right**
> > > >
> > > > > ## Lengde
> > > >
> > > > ---
> > > >
> > > > > **Group:** ` bucking/top/length/header/$r0c1 `
> > > > >
> > > > > **Row 1 · 2 components from left to right**
> > > > >
> > > > > **Button:** ` Cursor AV `
> > > > >
> > > > > ---
> > > > >
> > > > > **Button:** ` Δ `
> > > > >
> > > > > ---
> > > > >
> > > >
> > > > ---
> > > >
> > >
> > > **Row 1 · 1 component from left to right**
> > >
> > > > **412,0 cm**
> > >
> > > ---
> > >
> > > **Row 2 · 1 component from left to right**
> > >
> > > > Måling · aggregat
> > >
> > > ---
> > >
> >
> > ---
> >
> > > **BoxUI-frame:** ` bucking/top/diameter `
> > >
> > > **header**
> > >
> > > > **Group:** ` bucking/top/diameter/header `
> > > >
> > > > **Row 1 · 2 components from left to right**
> > > >
> > > > > ## Diameter
> > > >
> > > > ---
> > > >
> > > > > **Group:** ` bucking/top/diameter/header/$r0c1 `
> > > > >
> > > > > **Row 1 · 4 components from left to right**
> > > > >
> > > > > **Button:** ` Cursor AV `
> > > > >
> > > > > ---
> > > > >
> > > > > **Button:** ` Δ `
> > > > >
> > > > > ---
> > > > >
> > > > > **Button:** ` O/B · U/B `
> > > > >
> > > > > ---
> > > > >
> > > > > **Button:** ` BarkNOR-kalibrering `
> > > > >
> > > > > ---
> > > > >
> > > >
> > > > ---
> > > >
> > >
> > > **Row 1 · 1 component from left to right**
> > >
> > > > **28,4 cm**
> > >
> > > ---
> > >
> > > **Row 2 · 1 component from left to right**
> > >
> > > > Måling · O/B · eksempelprofil
> > >
> > > ---
> > >
> >
> > ---
> >
>
> ---
>
> **Row 2 · 1 component from left to right**
>
> > **Frame:** ` bucking/middle `
> >
> > **Row 1 · 3 components from left to right**
> >
> > > **BoxUI-frame:** ` bucking/middle/selection `
> > >
> > > **header**
> > >
> > > > ## Treslag / Sortiment
> > >
> > > > **Treslag**
> > > > Gran
> > > >
> > > > **Sortiment**
> > > > Sagtømmer
> > > >
> > > > Siste knapp: GRAN
> > >
> > > ---
> > >
> >
> > ---
> >
> > > **BoxUI-frame:** ` bucking/middle/suggestions `
> > >
> > > **header**
> > >
> > > > **Group:** ` bucking/middle/suggestions/header `
> > > >
> > > > **Row 1 · 1 component from left to right**
> > > >
> > > > > ## Apteringsforslag
> > > >
> > > > ---
> > > >
> > > > **Row 2 · 1 component from left to right**
> > > >
> > > > > **Group:** ` bucking/middle/suggestions/header/$r1c0 `
> > > > >
> > > > > **Row 1 · 2 components from left to right**
> > > > >
> > > > > **Button:** ` Optimalisering `
> > > > >
> > > > > ---
> > > > >
> > > > > **Button:** ` Kolonner `
> > > > >
> > > > > ---
> > > > >
> > > >
> > > > ---
> > > >
> > >
> > > **Row 1 · 1 component from left to right**
> > >
> > > > Neste kapp · eksempeldata
> > > > | Nr | Sortiment | Lengde | Topp-Ø |
> > > > | -- | --------- | ------ | ------ |
> > > > | 1  | Sagtømmer | 430 cm | 26 cm  |
> > > > | 2  | Massevirke| 310 cm | 18 cm  |
> > >
> > > ---
> > >
> > > **Row 2 · 1 component from left to right**
> > >
> > > > **Group:** ` bucking/middle/suggestions/$r1c0 `
> > > >
> > > > **Row 1 · 2 components from left to right**
> > > >
> > > > **Button:** ` Canonical `
> > > >
> > > > ---
> > > >
> > > > **Button:** ` Alternativ 1 `
> > > >
> > > > ---
> > > >
> > >
> > > ---
> > >
> > > **Row 3 · 1 component from left to right**
> > >
> > > > Planvalg er ubundet i denne prototypen.
> > >
> > > ---
> > >
> >
> > ---
> >
> > > **BoxUI-frame:** ` bucking/middle/currentStem `
> > >
> > > **header**
> > >
> > > > ## Stammen i aggregatet
> > >
> > > **Row 1 · 1 component from left to right**
> > >
> > > > Aktuell stamme
> > > >
> > > > **Låste feil: 0**
> > > > Ingen låste feil registrert
> > > >
> > > > Nåposisjon: 412,0 cm
> > > > Maks kvistet: 412,0 cm
> > > > Produsert / kappet: 0 cm
> > > > Segmentrot: 0 cm
> > >
> > > ---
> > >
> > > **Row 2 · 1 component from left to right**
> > >
> > > > Produksjon og feilflyt er plassholdere.
> > >
> > > ---
> > >
> >
> > ---
> >
>
> ---
>
> **Row 3 · 1 component from left to right**
>
> > **BoxUI-frame:** ` bucking/stemTrack `
> >
> > **header**
> >
> > > **Group:** ` bucking/stemTrack/header `
> > >
> > > **Row 1 · 2 components from left to right**
> > >
> > > > ## Stammeforløp
> > >
> > > ---
> > >
> > > > **Group:** ` bucking/stemTrack/header/$r0c1 `
> > > >
> > > > **Row 1 · 2 components from left to right**
> > > >
> > > > **Button:** ` taperNOR `
> > > >
> > > > ---
> > > >
> > > > **Button:** ` Mixed `
> > > >
> > > > ---
> > > >
> > >
> > > ---
> > >
> > > **Row 2 · 1 component from left to right**
> > >
> > > > Hele estimatet · Predikert · Målt · Plan · Rest · Kvalitet 1/2 · Underkjent · Utkast
> > >
> > > ---
> > >
> >
> > > **Plassholder for OperatorUnifiedStemSvg**
> > > Stammeprofil, målepunkter, kapp og manuell cursor trenger en egen visuell provider.
> > > Eksempel på Mermaid-innhold nedenfor utelates av konsolldumpen:
> > > [Mermaid utelatt]
> >
> > ---
> >
> > **footer**
> >
> > > **Group:** ` bucking/stemTrack/footer `
> > >
> > > **Row 1 · 3 components from left to right**
> > >
> > > > FØLGER AGGREGAT
> > >
> > > ---
> > >
> > > > Lengde: 412,0 cm
> > >
> > > ---
> > >
> > > > Rel. lengde: 0,0 cm
> > >
> > > ---
> > >
> > > **Row 2 · 3 components from left to right**
> > >
> > > > Pred. Ø O/B: 28,4 cm
> > >
> > > ---
> > >
> > > > Profilvolum O/B: 0,312 m³
> > >
> > > ---
> > >
> > > > Toppsylinder: 0,260 m³
> > >
> > > ---
> > >
> >
>
> ---
>
