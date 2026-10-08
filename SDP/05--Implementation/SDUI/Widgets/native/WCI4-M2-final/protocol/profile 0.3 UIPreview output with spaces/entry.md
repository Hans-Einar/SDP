# SDUI — ` Main `

Static GUI dump. Buttons and fields are text labels; no callbacks execute.

## Layout overview

Row/column structure in terminal cells. Heights follow content; this is not measured GUI geometry.

```text
SDUI GUI dump | Main | 160 columns | structural preview
Blå: [å🙂]                                                                       π prose
Blå: [å🙂]                                                                       π prose
```

## Content

Markdown is rendered as content. Nested blockquotes represent groups and frames. Horizontal siblings appear in reading order here; the overview shows placement. Mermaid diagrams are omitted.

> **Frame:** ` Main `
>
> **Row 1 · 1 component from left to right**
>
> > **Group:** ` Main/a `
> >
> > **Row 1 · 2 components from left to right**
> >
> > **Input:** ` Blå ` — ` å🙂 `
> >
> > ---
> >
> > > π prose
> >
> > ---
> >
>
> ---
>
> **Row 2 · 1 component from left to right**
>
> > **Group:** ` Main/b `
> >
> > **Row 1 · 2 components from left to right**
> >
> > **Input:** ` Blå ` — ` å🙂 `
> >
> > ---
> >
> > > π prose
> >
> > ---
> >
>
> ---
>
