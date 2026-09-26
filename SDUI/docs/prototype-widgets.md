# SDUI — limited prototype widgets

**Historical control gallery.** G5-M4 removed its fixture builder and temporary CLI
formats. Use [shared Go layout/SVG](../go/README.md) and
[state document export](../../SDL/go/README.md) for new UI.

Delivered for R24 on 2026-09-21: button/input SVG in the
[Markdown example](../examples/concept1-bucking.widgets.md). The separate
[HTML demo](../examples/prototype-controls.html) supports editing, mouse/keyboard press,
focus, disabled buttons and printing current values.

The gallery used SDUI 0.2. prototype_widgets consumed normalized instances from the
same parser, retaining identity, label, value and enabled. Hidden branches are omitted;
disabled is inherited. Callbacks are not transferred; no SDL file opens or SVG producer
runs. This description records the historical implementation.

## Two presentation formats

| Format | Capability | Limit |
| --- | --- | --- |
| SVG image in Markdown | Drawn buttons/fields, normal/pressed/focused/disabled | Static; no editing/clicks |
| Local HTML in browser | Native input, local button response, printing | Separate gallery; no storage/domain function |

CommonMark has no standard form widgets. [Raw HTML](https://spec.commonmark.org/0.31.2/#raw-html)
may appear in source, but readers control execution. SVG images also have
[interaction/script limits](https://developer.mozilla.org/en-US/docs/Web/SVG/Guides/SVG_as_an_image).
This is not a new Markdown extension for other readers.

Local XFMD inspection found ModelBuilder.cpp treats raw HTML as text, not DOM controls;
ImageDecoder.cpp copies a static pixbuf to Cairo, and README specifies GIF's first
frame. SVG fits its document path; forms/animation need another host capability.
HTML CSS/JS runs in a browser, not XFMD.

## Current export and historical scope

From SDUI, export with the active Go CLI:

```sh
go -C go run ./cmd/sdui ../examples/concept1-bucking.sdui --format svg --entry bucking -o /tmp/concept1.svg
```

The Python gallery builder is removed. Old commands remain in Git/dated evidence,
not as another supported build path. On September 21 it read concept1-bucking.sdui,
prototype-controls.sdui and recorded treemap geometry. SVG/HTML shared widget descriptions.
Regions retained reference rows/columns, with redundant group headings removed for space.
Length/Diameter have Cursor AV and Δ in their heading row; Diameter also has O/B · U/B
and BarkNOR calibration. Stem track has taperNOR/Mixed beside its title. Source places
these in header regions; SDUI's SVG composition draws them. This adds no Mermaid treemap
heading-control capability.

Placement was manual inside reference boxes, rendering selected controls/handwritten
data rather than every AST block. It did not implement general layout, Markdown or
font inheritance. SVG coordinates/demo fonts introduce no pixel dimensions into SDUI.
Pressed/focused SVG states are illustrations, not language attributes.

HTML has no network/storage dependency. Editing and pre-print handling copy values as
text into print blocks, allowing long values to wrap without input clipping. Print
removes pressed/focus styles; reduced motion disables CSS transitions.

## Continuation

This supplied port material for [G1–G3](implementation-plan.md). G2 replaced manual
placement with shared Fyne/SVG geometry; G3/G4 connect identities/properties to runtime
and SDL through typed Go ports. FOX is not first. The HTML press counter remains a
local demonstration, not an SDL/SDUI runtime.
[Verification](../evidence/prototype-widgets/README.md).
