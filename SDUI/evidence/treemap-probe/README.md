# Treemap as static UI sketch — 2026-09-21

The owner wanted one box layout with content, matching the text dump. The earlier
ASCII-plus-blockquote Markdown report did not meet that visual goal.

Local renderer inspection found horizontal allocation at even depth and vertical at odd
depth. A root with three groups and their leaves gives the required pattern. Parents sum
child weights: 7.5+7.5 / 11.25+22.5+11.25 / 40 produces rows 15/45/40.

The existing local mmdr 0.3.1 binary was used without build/source changes. verification.json
records its hash/checks. Repository HEAD was afab5e9; this does not prove the binary's build
source. mmdr consumed view.mmd/config.json and produced SVG/layout.json; rsvg-convert produced
a visually inspected PNG. Assertions covered row weights, middle columns, equal top columns
and nonempty labels in all six leaves. The SVG is displayed by the corresponding example
Markdown document.

Long top labels disappeared because measured blocks did not fit. Shortening them produced
six visible labeled boxes. Leaf weights are also printed; metadata/padding mean inner box
heights are not exactly 15/45/40 of the root.

This is a manually bounded probe, not general AST→treemap export or full Markdown in boxes.
Other Mermaid engines need not place identically. XFMD documents local SVG support, but no
XFMD GUI was run here; Markdown uses prerendered SVG. No runtime, interactive controls,
parser rules, commit or push were added.
