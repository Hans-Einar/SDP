# Prototype widgets — verification, 2026-09-21

R24 delivered an SVG reference composition and local HTML gallery, not general layout,
FOX, SDL runtime or ABI. Same-day update moved Length/Diameter/Stem-track controls into
heading rows and added Diameter Δ. Source, AST and both structural dumps were regenerated.
Librsvg visual review found no title/control overlap. Manifest/log describe this revision;
browser evidence covers the unchanged HTML gallery. Localized labels remain example data.

## Checks

- 36/36 Python tests passed; [log](tests.txt). Six new tests cover model/source reuse,
  omitted callbacks, visibility/disabled inheritance, SVG states, escaping, print blocks
  and generated SVG artifacts.
- Playwright/local Chrome tested editing, visible mouse press, keyboard activation,
  disabled controls, current print values, wrapping, narrow view, reload reset, no
  external requests and no JS errors. [Results](browser-verification.json).
- [Pressed button](browser-pressed.png) and [filled print preview](print-preview.png)
  were visually checked. Print CSS was emulated; physical printing, paginated PDF and
  XFMD GUI were not verified.
- Both SVGs were rasterized with local librsvg and checked for visible, nonoverlapping
  controls/content. SVG is a static image, with no clicks/editing.
- Repeated generation was byte-identical; [SHA-256 manifest](manifest.json) identifies
  original source/artifacts.

## Historical reproduction

These commands describe the retired builder at the recorded revision, from SDUI:

```sh
python3 tools/build_widget_previews.py
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

The separate browser probe requires playwright and /usr/bin/google-chrome as verification
tools, not production dependencies:

```sh
python3 evidence/prototype-widgets/verify_browser.py
```

It opens the local gallery and regenerates two screenshots/report. The demo requires no
server or external resources. [Scope and continuation](../../docs/prototype-widgets.md)
records limits, XFMD source inspection and the later shared geometry direction.
