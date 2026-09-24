# SDUI — composition clarification, 2026-09-21

Documentation only; no parser/runtime/renderer change.

The owner's canonical corners ^</>^/v</>v replace the preceding proposal. Ratio with
horizontal stretch uses full width and derived height, not contain-fit. The composition
proposal preserves header/body/footer and named mainBody references. Later clarification:
font=10/12 is absolute text size; resize does not scale content. The native font unit
remained open. Component formatting precedes separators; comma continues horizontally,
semicolon starts a row. R21 and the then-current P1/P2/P5 plan record remaining parser,
measurement and GUI checks. Nested `<>` groups have independent layout; a three-level
example demonstrates local separators. R13 and P1/P2 require explicit single-child groups
to survive normalization.

Checks: relative Markdown links, balanced fences, no pixel dimensions in new examples,
and W=1, H=0.25, ratio=16/9 → 1×0.5625 with height overflow. This is documentation/arithmetic,
not layout implementation. No GUI/PDF run, commit/push or external repository changes.
