# Go measurement contract — G2

G2-M1 implements one measurement/placement model in `go/layout`. SVG and native hosts receive `Box` values containing rectangle, clip, font, enabled state and instance path. Measurement runs no callbacks. Sources contain no pixel dimensions.

`font` uses absolute logical display units (DIP). `font=10` means ten units in both SVG viewBox and host. Physical unit size follows host DPI, not window size. The reference measurer embeds Go Regular through OpenType at 72 DPI; Markdown headings may use relative font sizes. SVG must use the same font or disclose substitution.

Rows do not change the ancestor for relative dimensions. Fr tracks distribute remaining space using `clamp(lambda * weight,min,max)`; minima are not added on top of weights. Scaled/content-sized siblings do not shrink to conceal overflow. Reserve header/footer before body; regions use available width unless another x policy is selected. `items=stretch` lets a single row component fill width and equalizes component heights in multicolumn rows.

A content-sized ancestor with relatively sized children on the same axis is rejected as `layout-dependency`. Width-driven ratio derives height without requiring a height-defined ancestor. Min/max never distort ratios. Padding/gap refer to the node's ancestor, as specified in the language proposal.

Overflow `error` produces diagnostics; `clip` also clips hit testing. Scroll requires host-owned scroll state and currently reports `unsupported-scroll`; it is not silently treated as clip. Measurement is limited to 200000 recursive operations and viewports up to 32768 logical units per axis. Parser source/depth/expansion limits apply first.

Tests cover fr versus scale, minima/maxima, resize-independent fonts, ratio versus contain, region reference areas, clipped hit testing and the actual Concept1 model's six boxes at 1920 × 1200. Other viewport sizes may reject content that does not fit; this is expected without explicit clip.
