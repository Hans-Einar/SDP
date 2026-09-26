# Concept1 in SDUI and console presentation

2026-09-21. [Source](../examples/concept1-bucking.sdui),
[AST](../examples/concept1-bucking.ast.json), [dump](../examples/concept1-bucking.dump.txt).
Ponsse source commit: 882ad7c0457a723550ea6cce5bf027230997ec2a; the
[manifest](../evidence/concept1-source-manifest.json) records file hashes.
Other local Ponsse changes existed; inspected UI/schema files were unchanged.

## Ported content

| Concept1 | SDUI |
| --- | --- |
| operator-layout.mjs: top, weight 15 | bucking/top, y=15fr; length/diameter x=1fr |
| middle, weight 45 | bucking/middle, y=45fr; selection/suggestions/currentStem 25/50/25fr |
| stemTrack, weight 40 | bucking/stemTrack, y=40fr |
| MetricBody and MetricSourceControls | Markdown measurement and named Cursor/Δ buttons |
| SelectionBody | Markdown species/assortment/last button |
| SuggestionRows, plan selection, settings | Raw Markdown table and static Canonical/alternative/columns/optimization buttons |
| CurrentStemDefectsBody, OperatorActualStemPanel | Representative facts and explicit production/defect placeholders |
| StemTrackHeaderContent/FooterContent | Legend groups, model buttons and measurements in header/footer |
| OperatorUnifiedStemSvg | Explicit Markdown placeholder; illustrative Mermaid fence tests omission |
| OperatorApp header/nav | Separate navigation definition and page wrapper |

This ports structure and representative content, not the full React application's
behavior. Cursor, bark calibration, plan selection, editing, tabs, conditional content
and domain operations are not implemented by this example. Values are prototype data,
not live machine readings. Buttons have no callbacks; parsing/dumping has no Ponsse
side effects. The localized UI text is intentional example content.

The original uses React, CSS Grid-fr and a viewport transform. SDUI preserves weights
and structure while keeping fonts/content unscaled on window resize. Root `{16:9,<->}`
derives height from width. The AST records this; console output does not compute that
geometry and is not screen-geometry evidence.

## Dump contract

Validated AST → shared normalized instance tree → terminal rows. The active static
dump is Go, without a separate parser or GUI dependency. `--entry bucking` shows the
bucking surface; page includes its shell. `--columns` measures terminal cells, not
source pixels or font size.

Raw Markdown retains ##, ** and table separators. Long lines wrap, so dumps are not
roundtrippable source. Ordinary backtick/tilde fences with info word mermaid become
omission markers; content inside other fences remains literal. No diagram engine runs.
The fence scanner is not a full Markdown parser for nested list/blockquote containers.
SVG references become placeholders; terminal control characters are escaped. Insufficient
width produces a diagnostic instead of overlap or silently missing boxes. Full emoji
graphemes/complex scripts require later measurement.

Static dumps use natural row heights and horizontal fr/scale allocation, not vertical
fr allocation, font sizes or every layout rule. The historical G2 dependency was a shared
measurement engine before GUI/TUI geometry could be compared; see current G2 evidence.

## What gh-tree uses

The static Markdown export below does not depend on a TUI implementation. Inspection
of Hans-Einar/gh-tree at 97cc0d8257603766dd741b49b7d8005857b421a9 found Go 1.25.0,
Bubble Tea v1.3.10 and Lip Gloss v1.1.0 as direct dependencies:
[pinned go.mod](https://github.com/Hans-Einar/gh-tree/blob/97cc0d8257603766dd741b49b7d8005857b421a9/go.mod).
Bubble Tea supplies the TUI; Lip Gloss supplies terminal styling/layout. These are
observed pins, not recommendations to freeze new dependencies to those versions.

Interactive console work was deferred after the 2026-09-21 Go/Fyne choice. Fyne is
first; a later Bubble Tea host would share the Go model/runtime, without cgo/C ABI or
another parser. Historical T1–T4 are inactive; see [plan](implementation-plan.md).
The static console dump is delivered.

## Static Markdown dump

`--format markdown --entry bucking` writes Markdown to stdout; `-o` saves it.
See [generated document](../examples/concept1-bucking.dump.md). A fenced text layout
comes first, followed by actual Markdown headings, emphasis, lists, tables and code.
Nested blockquotes, regions and row labels represent groups/frames. Horizontal siblings
appear sequentially in the content section; the overview shows their columns. This
is portable without HTML/CSS; tables need a reader extension. Buttons/inputs are labels,
not interactive controls. The same fence scanner omits Mermaid. Relative content links
are preserved; moved exports still need access to those resources.

Mermaid treemap supports section/leaf labels and area weights. Full Markdown tables or
multiple content blocks inside a box are not a documented treemap contract; its area
allocation does not generally follow explicit SDUI rows/columns.
[Official syntax](https://mermaid.js.org/syntax/treemap).

**Correction after visual feedback:** this report does not meet the owner's single
GUI-layout goal. A [treemap probe](../examples/concept1-bucking.treemap.md) shows that
this Rust renderer's alternating axes can reproduce Concept1's main rows/columns with
a suitable hierarchy. It uses shortened labels and prerendered SVG, visually checked.
It is not a general SDUI export or full Markdown-in-box implementation.
[Probe and limits](../evidence/treemap-probe/README.md).
