# SDUI — frame, widget and layout proposal

**ID:** SDUI-LAYOUT-002 · **Date:** 2026-09-21 · **Status:** design proposal. Preserves owner direction/proposals from that date. **Reading guide, 2026-09-24:** Go now implements SDUI 0.2; Python/0.1 are retired. [Language](language.md), [measurement](go-layout-contract.md) and [runtime](../go/runtime/README.md) bound selected/delivered behavior. Scroll is still rejected, fonts use DIP, and [Markdown coverage](markdown-provider.md) is bounded. These proposals do not extend active contracts. The [plan](implementation-plan.md) owns milestones.

## 1. Smallest useful language

Target forms, control panels and document surfaces with native inputs/buttons, Markdown and diagrams. Need predictable nesting, sizing, spacing, alignment and overflow; initially no CSS, free positioning, layout animation, grid spans or general constraint solver.

| Form | Proposed meaning |
| --- | --- |
| `page = [...];` | Named frame definition; no extra root assignment/wrapping braces. |
| `mainBody = <...>;` | Named reusable widget group. |
| `[...]` | Frame containing subframes/groups in one content list. |
| `<...>` | Widgets, Markdown strings and nested groups. |
| `name = button("OK")` | Named widget; argument is a label, not a Markdown widget. |
| `"# Explanation"` as list item | Markdown-widget shorthand; may be named `intro = "..."`. |
| `]*box` | Decorated frame; `*b` is a fixed alias for `box`. |
| `node { ... }` | Formatting/layout of preceding component, including strings/references. |

The owner's unique-prefix idea remains historical. Recommend fixed alias b so installed components cannot change old source meaning. Normalize to box. Additional variants need registered properties/measurement/presentation; `*` is not dynamic loading.

Frames are layout/clip boundaries, not necessarily FOX windows. Undecorated frames draw no border. header/body/footer are regions; unmarked content becomes body. header replaces heading, not an alias. Empty frames/groups are placeholders without intrinsic size. [Composition proposal](frame-composition-proposal.md).

Names are unique per definition; instance paths distinguish reuse. Callbacks/external handles require widget names; unbound controls may be anonymous. Anonymous Markdown may receive source-based internal IDs without reload-stability guarantees. Layout-only changes preserve named identity. ref/callback/setHandle remain binding data; parsing never opens/executes SDL.

## 2. Rows, groups and axes

Owner clarification, 2026-09-21: formatting applies to the preceding whole component **before** its separator, for named/anonymous widgets, strings, groups, frames, references and region values. Order: component → optional variant → optional formatting → separator. Final items may close directly with > or ]. Propose one combined formatting block per component, validated by type. Detached post-separator blocks are invalid.

```text
mainBody = <
  "Text" {font=10}, ok = button("OK") {font=12};
  "Next row" {font=10}
> {<->};
page = [mainBody {font=12}]*b {16:9, <->};
```

Reference formatting affects that instance, not its definition. Proposed font inheritance: frame/group defaults, overridden explicitly by children.

Comma separates row siblings; semicolon starts below the entire preceding row plus gap. Applies to body/groups; regions are separate roles. Inner argument/layout commas belong to inner syntax. `<` within strings never starts widgets. Top-level semicolons end definitions, not UI rows. Nested groups are single parent-list items; separators are local. Preserve explicit groups even with one child. [Group example](frame-composition-proposal.md#nested-groups-with-independent-layout).

Normalization creates vertical row stacks/horizontal rows, without shared column widths; this is not a grid/table language. A one-item row adds no geometric intermediate, so its height/weight applies to the stack. Use explicit subframes for shared-height/weight tracks.

```text
sdui 0.2;
page = [
  tools = [<"## Tools"; run = button("Run")>]*box {scale-x=0.25},
  content = [<"## Result">]*box {x=1fr};
  status = <"Ready"> {x=fill}
] {scale=1, gap=0.01, padding=0.01};
```

Multicomponent rows use content height. To distribute height, wrap rows in explicit frames. Do not propagate arbitrary child y=1fr to synthetic rows.

## 3. Canonical shapes, relative dimensions and aspect ratios

Owner clarification, 2026-09-20: source dimensions reference nearest ancestors, with no pixel width/height. Ratio x:y applies to frames. Scale may drive x/y/both, but ratio permits one axis only. Names such as scale-x are proposed spellings.

### Mini arrows and canonical spelling

Direction pairs express unordered alignments: v< and <v both mean down-left, canonically v<. Corners: ^<, >^, v<, >v. Normalize <^→^< and ^>→>^; source AST retains spelling.

| Canonical form | Meaning |
| --- | --- |
| `^<`, `-^`, `>^` | Left/center/right, top |
| `<-`, `--`, `->` | Left/center/right, vertically centered |
| `v<`, `-v`, `>v` | Left/center/right, bottom |
| `<->` | Fill assigned width; ratio derives height |
| `^\|v` | Fill assigned height |
| `>-<` | Fit horizontal content extent |
| `>\|<` | Fit vertical content extent |

Vertical operators are ^|v and >|<; table escaping is Markdown only. Independent rules may reorder. `{v<, >-<, >|<}` packs content at lower left. Inward/outward operators are complete tokens: `<->` and `>-<` differ despite sharing characters; direction-pair permutation does not equate them.

¤ remains undefined, reserved as a possible future meaningful operation, not an arbitrary synonym/parser rule.

### Nearest ancestor and relative dimensions

Explicit frames/groups establish child references: their **inner available rectangle** after regions, decoration and padding. Synthetic rows do not change source ancestors; arguments are not ancestors. Root references host-provided layout area excluding window decoration/toolbars.

scale-x=0.5 requests half reference width; scale-y=0.5 half height; scale=0.5 both. Without ratio, x/y may differ. Values are positive finite ratios; above 1 is permitted with overflow policy. Bare numbers imply no hidden pixel/percentage conversion.

References use parent content area **before** sibling allocation. Two 0.5 children plus gap overflow; do not silently correct. Use weights to share remaining space after gaps. Indeterminate reference axes produce dependency diagnostics, not fallback to distant ancestors/main window.

### Aspect ratio and driving axis

16:9 means outer width/height including regions/padding. Both terms positive/finite; normalize 32:18 and 16:9 equivalently. For reference W×H and r=x/y:

- `{16:9,<->}`: assigned width B, height B/r; root B=W.
- `{16:9,^|v}`: assigned height, width height*r.
- `{16:9,scale-x=0.75}`: width 0.75W, derived height.
- `{16:9,scale-y=0.5}`: height 0.5H, derived width.
- Reject scale=0.75 or simultaneous scale-x/scale-y even if values accidentally fit the ratio.
- Reject simultaneous horizontal/vertical fill with ratio.
- Without explicit scale/fill, propose contain: largest fitting ratio rectangle, deriving one axis.

Never adjust explicit scale silently to fit; apply overflow. Min/max/native minima never silently distort ratios. The derived axis cannot have another fill/content/size policy. One scale/fill axis drives explicit ratio sizing; auto-contain only applies without a chosen axis. Width fill never becomes contain because height is insufficient. Define fr interactions before supporting them.

### Other relative layout

| Property / form | Proposed initial profile |
| --- | --- |
| scale-x, scale-y, scale | Ancestor ratios; no px. |
| x, y | content (default), fill or positive Nfr; no absolute lengths. |
| min-x/min-y/max-x/max-y | Optional relative bounds on same axes; min ≤ max. |
| align-x/align-y | start/center/end, default start; mini-arrow equivalents. |
| gap/gap-x/gap-y | Nonnegative width/height ratios, default 0.01; axis-specific overrides gap. |
| padding | Relative factor or (top,right,bottom,left); horizontal relative to width, vertical to height; default 0, box 0.01. |
| justify | start/center/end/between; main-axis leftover spacing. |
| items | start/center/end/stretch; child cross-axis defaults. |
| overflow-x/overflow-y | error/clip/scroll; default error. |
| wrap | none/wrap, initially one explicit group row. |
| font | Absolute text size such as font=10; no ancestor/resize scaling. |

Gap/padding/bounds are relative, not previously proposed logical pixels. Ancestor reference avoids circular self-reference in content-sized containers. Validate explicit min/max without changing native minima. scale-x and x policies conflict, likewise y. Duplicate property rules are errors, not last-one-wins.

fill acts as 1fr on the main axis and fills cross-axis allocation. **fr weights siblings; scale takes a parent fraction.** Reserve scaled/content children first, then distribute remainder with `size = clamp(lambda * weight, min, max)`, not old minimum+grow. Max saturation leaves space for justify. No implicit flex-shrink. items supplies defaults only without explicit child cross-axis size/alignment.

## 4. Measurement, scaling and responsiveness

Owner clarification, 2026-09-21: font=10 is absolute, not relative. Resize changes area, placement and wrapping, not overall content/text scale. Frame scale properties define layout extents, not visual content transforms. Widgets may receive larger/smaller rectangles without changing label fonts. Hosts supply concrete rectangles/text measurements; display zoom/DPI are independent.

At this proposal date shared absolute units and Fyne/Markdown/SVG conversions remained to be defined. Absolute sizing was decided; point versus logical-display units was open. Backends cannot choose private units. Heading styles may scale relative to the absolute base font, independently of windows. Current Go units are defined in the linked measurement contract.

Outer `{16:9,<->}` expresses FixedAspectViewport's role through host width/derived height, not contain. Do not copy Concept1's overall content scaling. Keep display/control/input geometry consistent.

wrap starts before nonfitting widgets while preserving source/tab order. Use relative/content sizes for row selection; initially reject main-axis fr/fill with wrap to avoid circular allocation. Oversized widgets follow overflow rather than disappear.

clip also clips input areas. scroll needs finite viewports/host support. The original proposal called for scroll-y first, x-scroll unsupported until tested, with bounded scrollbar/measurement passes. **Current implementation rejects both scroll axes**, as stated above. Preserve ratios; scrolling affects content areas, not frame shapes.

Measure Markdown height at agreed width. Native widgets provide minimum/preferred sizes and optional baselines at selected fonts. Detect content-parent width/height cycles; no unbounded measurement feedback. Baseline alignment deferred.

visible=false removes layout/tab participation but retains instances. enabled=false retains geometry while blocking actions. Define hiding/focus/draft policy. Programmatic updates are not user events.

## 5. Markdown and boundaries

Ordinary strings escape explicitly. Proposed triple double quotes preserve raw multiline content without dedent/interpolation until the next triple delimiter; literal triple delimiters initially require ordinary escaped strings. Never normalize Markdown indentation during parsing.

Quoted list/region content is Markdown; button labels remain ordinary property strings. Content services, not SDUI lexers, handle Mermaid fences. Inputs: source base, width, theme, cancellation, budget. Outputs: measurements, diagnostics, source maps, resources.

The original proposal planned XFMD's supported Markdown/diagram profiles and host resource rules. Nested sdui fences were to be rejected locally, preventing recursive UI instantiation. Printing was to use frozen accepted state, no callbacks/unsent drafts. Proposed scroll export clips viewports with visible overflow; full expansion/pagination is later work. Active Go Markdown limits supersede these unimplemented broader proposals.

## 6. Representative fixtures

At the proposal date these parsed/validated in 0.2; geometry was still pending.

```text
sdui 0.2;
panel = [
  header="## Workspace",
  top = [
    length = [<"## Length\n12.4 m">]*b {x=1fr},
    diameter = [<"## Diameter\n32 cm">]*b {x=1fr}
  ] {y=15fr, items=stretch};
  middle = [
    selection = [<"Selection">]*b {x=25fr},
    suggestions = [<"Suggestions">]*b {x=50fr},
    stem = [<"Stem">]*b {x=25fr}
  ] {y=45fr, items=stretch};
  track = [<"Stem profile">]*b {y=40fr}
]*box {16:9, <->, items=stretch, gap=0.01};
```

Borrows Concept1's six box identities and 15/45/40, 25/50/25 distribution. Headers, decoration, minima and content prevent automatic CSS Grid equivalence. Geometry fixtures must specify viewports/comparison targets.

```text
sdui 0.2;
form = [
  fields = <"## Context"; context = input("Name", value="C1")>
    {x=fill};
  actions = <cancel = button("Cancel"), ok = button("OK")>
    {v<, >-<, >|<, gap=0.01}
] {scale=1, justify=between, padding=0.01};
```

Negative cases: duplicate IDs; unknown variants; targetless layout; min>max; negative gap; fill with content packing; unterminated Markdown; wrap with fr; unsupported scroll; anonymous callbacks; nonfitting native minima; layout tokens preserved inside strings; px/bare absolute layout lengths (font excepted); detached formatting; invalid ratio/two scale axes; unresolved ancestors; real parents versus synthetic rows.

## 7. Deliberately deferred

Free coordinates, overlap/z-order, calc expressions, media queries, grid spans, source-independent reordering, docks/splitters, layout animation, dynamic plugins and virtual lists. Extend only when concrete prototypes demonstrate needs poorly served by nesting.
