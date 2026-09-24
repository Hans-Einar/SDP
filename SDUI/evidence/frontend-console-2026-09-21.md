# SDUI 0.2 — Concept1, AST and console dump

2026-09-21: **24/24 tests passed**.
[Log](frontend-console-tests-2026-09-21.txt), [code/environment manifest](frontend-console-manifest.json),
[Concept1/gh-tree source basis](concept1-source-manifest.json).
Historical Python evidence; the active frontend has since moved to Go.

## Verified

Commands then run from SDUI:

```sh
PYTHONPATH=src python3 -m unittest discover -s tests -v
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format dump --entry bucking --columns 160
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format dump --entry page --columns 200
```

- Frames, nested groups, local rows, regions and formatting in parser/AST.
- Unambiguous instances, forward references, cycles and separate expansion budget.
- UTF-8/CRLF/multiline positions and raw string preservation.
- Typed widgets, named callbacks, anonymous unbound buttons and declarative setHandle.
- Canonical arrows, ratio/axis conflicts, invalid absolute dimensions/font rejection.
- Local separators, retained single-child groups, formatting placement and stray block rejection.
- Source/token/depth/node/argument limits, all minimal-fixture truncations and 500 deterministic noise strings.
- Visible dump rows/columns, raw Markdown, omitted Mermaid, escaped control characters and CJK/combining cell widths.
- CLI entry, no partial failure output, JSON diagnostics, UTF-8/I/O/source protection.
- Six Concept1 box IDs and original 15/45/40 and 25/50/25 weights.

Both stored ASTs and the Concept1 dump matched fresh CLI output byte-for-byte. The page
wrapper ran at 200 columns; 160-column console inspection showed boxes/raw tables without
Mermaid execution. Local links, fences and git diff --check were checked.

## Scope and cleanup

The 0.1 parser, Box/Content AST, grammar, text widget and active examples were replaced,
without fallback. Historical logs/manifests retain their old-code meaning. This was a
Python frontend/static structure, not native libsdui, C ABI or runtime. No new package
dependencies. Natural row heights did not verify vertical fr, 16:9 geometry, arrows,
font measurement, full Markdown or FOX. Values/SVG/defect/production content are explicit
examples/placeholders. No SDL module, Ponsse operation or callback executes.

GitHub API inspection confirmed Go, Bubble Tea and Lip Gloss in gh-tree. Interactive TUI
was planned then, not implemented. No GUI/PDF run, external changes, branch switch,
commit or push occurred in that delivery.
