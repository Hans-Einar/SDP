# SDUI — ` bucking `

Statisk GUI-dump. Knapper og felt er tekstetiketter; ingen callbacks kjøres.

## Layoutoversikt

Rad-/kolonnestruktur i terminalceller. Høydene følger innholdet; dette er ikke målt GUI-geometri.

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

## Innhold

Markdown gjengis som innhold. Nestede sitatblokker viser grupper og frames. Horisontale søsken står i leserekkefølge her; plasseringen vises i oversikten. Mermaid-diagrammer er utelatt.

> **Frame:** ` bucking `
>
> **Rad 1 · 1 komponent fra venstre mot høyre**
>
> > **Frame:** ` bucking/top `
> >
> > **Rad 1 · 2 komponenter fra venstre mot høyre**
> >
> > > **BoxUI-frame:** ` bucking/top/length `
> > >
> > > **header**
> > >
> > > > **Gruppe:** ` bucking/top/length/header `
> > > >
> > > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > > >
> > > > > ## Lengde
> > > >
> > > > ---
> > > >
> > > > > **Gruppe:** ` bucking/top/length/header/$r0c1 `
> > > > >
> > > > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > > > >
> > > > > **Knapp:** ` Cursor AV `
> > > > >
> > > > > ---
> > > > >
> > > > > **Knapp:** ` Δ `
> > > > >
> > > > > ---
> > > > >
> > > >
> > > > ---
> > > >
> > >
> > > **Rad 1 · 1 komponent fra venstre mot høyre**
> > >
> > > > **412,0 cm**
> > >
> > > ---
> > >
> > > **Rad 2 · 1 komponent fra venstre mot høyre**
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
> > > > **Gruppe:** ` bucking/top/diameter/header `
> > > >
> > > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > > >
> > > > > ## Diameter
> > > >
> > > > ---
> > > >
> > > > > **Gruppe:** ` bucking/top/diameter/header/$r0c1 `
> > > > >
> > > > > **Rad 1 · 4 komponenter fra venstre mot høyre**
> > > > >
> > > > > **Knapp:** ` Cursor AV `
> > > > >
> > > > > ---
> > > > >
> > > > > **Knapp:** ` Δ `
> > > > >
> > > > > ---
> > > > >
> > > > > **Knapp:** ` O/B · U/B `
> > > > >
> > > > > ---
> > > > >
> > > > > **Knapp:** ` BarkNOR-kalibrering `
> > > > >
> > > > > ---
> > > > >
> > > >
> > > > ---
> > > >
> > >
> > > **Rad 1 · 1 komponent fra venstre mot høyre**
> > >
> > > > **28,4 cm**
> > >
> > > ---
> > >
> > > **Rad 2 · 1 komponent fra venstre mot høyre**
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
> **Rad 2 · 1 komponent fra venstre mot høyre**
>
> > **Frame:** ` bucking/middle `
> >
> > **Rad 1 · 3 komponenter fra venstre mot høyre**
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
> > > > **Gruppe:** ` bucking/middle/suggestions/header `
> > > >
> > > > **Rad 1 · 1 komponent fra venstre mot høyre**
> > > >
> > > > > ## Apteringsforslag
> > > >
> > > > ---
> > > >
> > > > **Rad 2 · 1 komponent fra venstre mot høyre**
> > > >
> > > > > **Gruppe:** ` bucking/middle/suggestions/header/$r1c0 `
> > > > >
> > > > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > > > >
> > > > > **Knapp:** ` Optimalisering `
> > > > >
> > > > > ---
> > > > >
> > > > > **Knapp:** ` Kolonner `
> > > > >
> > > > > ---
> > > > >
> > > >
> > > > ---
> > > >
> > >
> > > **Rad 1 · 1 komponent fra venstre mot høyre**
> > >
> > > > Neste kapp · eksempeldata
> > > > | Nr | Sortiment | Lengde | Topp-Ø |
> > > > | -- | --------- | ------ | ------ |
> > > > | 1  | Sagtømmer | 430 cm | 26 cm  |
> > > > | 2  | Massevirke| 310 cm | 18 cm  |
> > >
> > > ---
> > >
> > > **Rad 2 · 1 komponent fra venstre mot høyre**
> > >
> > > > **Gruppe:** ` bucking/middle/suggestions/$r1c0 `
> > > >
> > > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > > >
> > > > **Knapp:** ` Canonical `
> > > >
> > > > ---
> > > >
> > > > **Knapp:** ` Alternativ 1 `
> > > >
> > > > ---
> > > >
> > >
> > > ---
> > >
> > > **Rad 3 · 1 komponent fra venstre mot høyre**
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
> > > **Rad 1 · 1 komponent fra venstre mot høyre**
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
> > > **Rad 2 · 1 komponent fra venstre mot høyre**
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
> **Rad 3 · 1 komponent fra venstre mot høyre**
>
> > **BoxUI-frame:** ` bucking/stemTrack `
> >
> > **header**
> >
> > > **Gruppe:** ` bucking/stemTrack/header `
> > >
> > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > >
> > > > ## Stammeforløp
> > >
> > > ---
> > >
> > > > **Gruppe:** ` bucking/stemTrack/header/$r0c1 `
> > > >
> > > > **Rad 1 · 2 komponenter fra venstre mot høyre**
> > > >
> > > > **Knapp:** ` taperNOR `
> > > >
> > > > ---
> > > >
> > > > **Knapp:** ` Mixed `
> > > >
> > > > ---
> > > >
> > >
> > > ---
> > >
> > > **Rad 2 · 1 komponent fra venstre mot høyre**
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
> > > **Gruppe:** ` bucking/stemTrack/footer `
> > >
> > > **Rad 1 · 3 komponenter fra venstre mot høyre**
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
> > > **Rad 2 · 3 komponenter fra venstre mot høyre**
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
