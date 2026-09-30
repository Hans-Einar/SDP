# Text and Markdown fidelity across SDUI exports

| Field | Value |
| --- | --- |
| id | KB-SDUI-004 |
| project | SDUI |
| type | Bug |
| CardState | backlog |
| Systems | SDUI |
| created | 2026-09-29T16:53:54.251554+00:00 |
| source | PLAN-SDP-0009; KB-SDP-041; external XFMD gap register |
| next_review | At the next SDL/SDUI work selection |

## Observation and reproduction

[Study evidence](../../02--Requirements/XFMD-Gaps/Evidence.md) isolates
GAP-XFMD-SDUI-011 on producer a274265. The embedded Go Regular font maps ↻, ↺,
✓ and 界 to glyph zero, while svg.Text returns no error. layout.Lines breaks
ordinary prose at rune boundaries; Markdown's Block discards nested list depth.
Inspected exports at widths 220 and 420 show missing-glyph boxes and flat lists;
the narrow export also splits an ordinary word.

Use [the fixture](../../02--Requirements/XFMD-Gaps/evidence/typography.sdui),
[glyph probe](../../02--Requirements/XFMD-Gaps/evidence/glyph_probe.go) and
[reproduction commands](../../02--Requirements/XFMD-Gaps/Evidence.md#reproduction).
This identifies an implementation/provider fidelity problem, not new grammar.
The source report's full-coverage expectation exceeds the current bounded
provider; choose an explicit repertoire/fallback policy rather than promise all
Unicode implicitly.

## Proposed fix and acceptance

Select a small ImplementationPlan following XGP2. Define one measured/painted
font fallback or explicit unsupported-glyph outcome; preserve word boundaries
where possible with a defined long-word fallback. Carry list nesting into layout
and painting, including continuation indentation. Avoid a second measurement
path that produces geometry different from export/native rendering.

Acceptance: supported and missing glyphs, combining sequences, whitespace,
long words, narrow/normal/wide widths, nested and multiline lists, deterministic
geometry and inspected SVG exports. Verify relevant host consistency separately;
static evidence alone cannot establish Fyne or XFMD Pango behavior. Existing
Concept1 geometry and parser fixtures must remain valid or have an explicitly
explained contract change. No pixel-based UI dimension rule is introduced.

## Scope and authority

SDUIPresentation owns the implementation; use SDUIRuntime/NativeHost only if a
selected shared text contract requires them. [KB-SDUI-003](%23003--SDUI--Proposal--Capabilities-and-navigation-pilot.md)
owns new interaction/resource capabilities. No XFMD changes or legacy parser
fallback. Registration records observed limitations and recommended work;
implementation, merge and release remain unselected.

## Worklog

- 2026-09-29T16:53:54.251554+00:00 — EVT-KB-SDUI-000013: Registered from completed producer research; remains backlog.
