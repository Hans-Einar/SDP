# Bounded Markdown provider — G2-M4

The Go provider parses Markdown once before measurement. Measurement and SVG share blocks, table cells and diagram dimensions. SVG/Fyne CLIs use it; raw Markdown in console dumps remains a deliberately separate export form.

Supported presentation profile:

- Six heading levels, paragraphs, line breaks and simple list items.
- Inline strong/emphasis, code and links retain readable content without Markdown markers. This initial profile does not render separate font weights, italics or link navigation and does not claim full CommonMark presentation.
- Fenced/indented code as literal text; GFM tables with cell borders, wrapping and identical SVG/layout measurements.
- Mermaid fences become separately identified diagram resources. Without a registered renderer they show explicit placeholders, not silent omission.

HTML/images are rejected. URLs are never loaded automatically. Each Markdown widget is limited to 32 KiB, 256 blocks and eight diagrams. External resources, animation, CSS and DOM scripting are outside the profile.

`--mermaid-renderer /absolute/path/to/mmdr` registers a local SVG-export renderer. Source goes through stdin, never a shell. Initial verified capability: `flowchart`/`graph`; other diagram types produce profile errors. Maximum source: 12 KiB; renderer timeout: five seconds; maximum output: 4 MiB. SVG requires a finite viewBox and no active/external elements or references. Initialization directives are forbidden. The parser imports no renderer.

`--resources DIR` writes content-named `.svg` files separately; the combined SVG also embeds data resources and works alone. Existing resources of the same name must have identical bytes. Directory-transaction publication, cache revisions and leases follow in G6; M4 is not a document daemon.

In this phase Fyne CLI uses the provider without an external diagram renderer, so diagrams show explicit placeholders. This avoids claiming support for every Mermaid SVG feature in Fyne's rasterizer without evidence. External flowchart resources are verified in SVG export through librsvg.
