# Concept1 — reference for relative surfaces

Checked on 2026-09-20 against local Ponsse HEAD 882ad7c. Inspected files had no reported local changes; unrelated changes preserved. No GUI execution. Paths below are relative to /home/warloc/git/ponsse.

## Actual implementation

| Concept1 source | Observation |
| --- | --- |
| shared/ui-box/model.mjs | layoutTrackTemplate maps positive weights to CSS fr. |
| shared/ui-box/UILayout.jsx | Groups use row/column CSS Grid; IDs bind separate React content. aspectRatio becomes a CSS variable. |
| apps/operator-ui/src/ui/operator-layout.mjs | Ratio 16/9, vertical weights 15/45/40, middle 25/50/25. |
| shared/ui-box/fixed-aspect-model.mjs | Internal 1600×1200/1600×900 surfaces for 4:3/16:9; centered fit uses min(viewportWidth/designWidth, viewportHeight/designHeight). |
| shared/ui-box/FixedAspectViewport.jsx | Watches main-window resize; one overall design-surface scale transform; geometry-authority=outer-only. |
| shared/ui-box/UIAspectRegion.jsx | scaling-authority=none; no nested scaler. |
| shared/ui-box/ui-box.css | Surfaces fill 100%; gap/padding/header/control sizes still use CSS px. |
| apps/operator-ui/src/ui/operator-ui.css | Overrides ui-layout-frame aspect-ratio to auto; outer viewport matters. |

This provides relative track distribution and overall internal-surface scaling, not arbitrary ancestor-relative scaling on each SDUI frame or internally pixel-free layout.

## Transfer to SDUI

Reuse weights, aspect ratios, host-owned viewports and consistent image/control/input transforms. Children should not measure the main window independently. Do not copy fixed 1600-based dimensions, root-only scaling or CSS px into source language.

Scale-x/scale-y/scale size frames relative to nearest source ancestors. Owner clarification, 2026-09-21: outer `{16:9,<->}` expresses FixedAspectViewport's role, filling width/deriving height even when the window is too short. This differs from Concept1's min-based contain fit. SDUI also retains absolute fonts and recomputes layout/wrapping; do not copy overall content scaling. Resolve shared native font units before measurement implementation.

## G2 calculation examples

Normalized reference units, not source pixel dimensions.

| Reference | Properties | Expected frame extent |
| --- | --- | --- |
| 1 × 1 | 16:9, scale-x=0.8 | 0.8 × 0.45 |
| 1 × 1 | 16:9, scale-y=0.5 | 8/9 × 0.5 |
| 1 × 0.75 | 16:9 without scale, proposed contain | 1 × 0.5625 |
| 1 × 0.25 | 16:9, <-> | 1 × 0.5625: height overflow, preserved width |
| 1 × 1 | 4:3, scale-x=0.6 | 0.6 × 0.45 |
| Parent 0.8 × 0.45 | Child without ratio, scale=0.5 | 0.4 × 0.225; frame scale does not scale content |
| 1 × 1 | 16:9, scale=0.8 | Reject: both axes explicitly driven |
| 1 × 1 | 16:9, scale-y=1 | 16/9 × 1: overflow, no silent distortion |

These are test oracles, not results from an implemented layout engine at the study date. Synthetic rows do not change ancestor references; explicit new parents do.
