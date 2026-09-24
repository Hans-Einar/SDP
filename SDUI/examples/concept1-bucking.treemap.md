# Bucking — static treemap probe

![Bucking surface: two top boxes, three middle boxes and stem track below](concept1-bucking.treemap.svg)

This probe displays boxes and text together. The local Mermaid Rust renderer produced
SVG so Markdown readers can display the same layout without recomputing it.

Rows use 15/45/40; middle columns 25/50/25. Headings, padding and gaps occupy group space.
Printed numbers 7.5, 11.25, 22.5 and 40 are area weights, not bucking data.

Content is deliberately shortened, with intentional localized UI labels. This handcrafted
Concept1 feasibility probe is not general SDUI export. Labels are plain multiline text;
tables, full Markdown and native widgets are not implemented by the probe.

[Treemap source](../evidence/treemap-probe/view.mmd) ·
[Renderer configuration](../evidence/treemap-probe/config.json)
