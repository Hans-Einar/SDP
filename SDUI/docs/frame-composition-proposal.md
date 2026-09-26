# SDUI — frames with header, body, footer and reuse

**Date:** 2026-09-21 · **Status:** owner input and concrete design proposals. Supplements the [layout proposal](layout-language-proposal.md). **Status clarification, 2026-09-24:** Go frontend, normalization, layout and runtime are delivered. [Language](language.md), [measurement](go-layout-contract.md) and [runtime](../go/runtime/README.md) govern execution; proposals below are not another active profile.

## 1. Width-driven aspect surface

```text
page = [ <"Content"> ] {16:9, <->};
```

Root fills host width W and derives height W×9/16. Available height does not reduce width. Insufficient height uses agreed overflow/scroll rather than contain/ratio changes. This expresses FixedAspectViewport's role in SDUI.

Horizontal fill drives x; ratio drives y. Vertical fill may drive y/derive x; both conflict with ratio. Subframes reference assigned areas in nearest parents. scale-x=0.75 selects an explicit ancestor fraction; fill uses allocated space. Root fill and scale-x=1 have equal width targets.

Resize changes geometry, not content scale. Fonts stay absolute; do not copy Concept1's fixed-surface overall scaling.

## 2. Regions as content values

```text
page = [
  header = [ ],
  <"Body">,
  footer = [ ]
]*b {16:9, <->};
```

Recommend header/body/footer as structural frame roles, not arbitrary widget IDs. Regions may contain frames, groups, Markdown or references. String formatting belongs to Markdown widgets. Roles also apply to undecorated frames; *b selects BoxUI presentation.

Proposed rules:

- At most one region of each kind; absent regions consume no space.
- Unmarked content becomes body. Optional explicit body cannot mix with implicit body content.
- Regions stack header/body/footer independently of source order, not as ordinary comma siblings. Body retains comma/semicolon rows.
- Measure header/footer at frame width; body receives remainder after decoration/padding/gaps. Ratio covers the entire frame, not body plus external regions.
- Oversized headers/footers follow overflow without ratio distortion.
- Regions define semantic content areas: body children reference remaining body area; synthetic rows do not change references.
- Empty frames are valid placeholders; box does not invent headers/footers.

Replace heading with header rather than maintain two mechanisms. `header="## Heading"` uses ordinary Markdown, supporting levels/multiline content without new widget types.

## 3. Named groups and instances

### Nested groups with independent layout

Owner clarification, 2026-09-21: nested <> groups are individual parent-list components with independent layout.

```text
mainBody = <
  <"group1", "with two UI components"> {v<, >-<, >|<},
  <"group2"; <"subgroup", "with two UI components"> {>-<}>
    {>^, >-<, >|<}
> {<->};
```

Outer comma places groups horizontally; group1's comma places its texts horizontally; group2's semicolon places its subgroup below its label without affecting the outer row. Separators are list-local.

Formatting after > belongs to that group. Group alignment places it within parent allocation; internal rules govern children. Preserve explicit groups in AST/normalization even with one child, retaining formatting/ancestor references. Groups need no name to receive layout.

### Reusing named groups

Owner example with proposed top-level semicolon separators:

```text
sdui 0.2;
mainBody = <"main body widgets", button1 = button("OK")>;
page = [
  header = "##Markdown Heading" {<-},
  mainBody,
  footer = <
    "footerText left center" {<-, font=10},
    "footerText right center" {->, font=12}
  >
]*b {16:9, <->};
```

mainBody is a reusable definition instantiated by reference, not a string/callback/shared live FOX object. `body=mainBody` makes its role explicit.

Recommend definition-local name uniqueness and full runtime identity including instance path (page/mainBody/button1) plus generation. Two uses need distinct names, such as left/right=mainBody, avoiding internal button1 collision. This refines the former flat namespace rule. Reject ambiguous refs, duplicate instance names and cycles. G1 Go supports forward refs/static expansion/full setHandle paths; Python/0.1 are removed. This does not automatically bind SDL runtime.

Only instantiation establishes UI lifetime. Library definitions do not mount as extra pages. Host APIs select entry frames without depending on the name page or source order.

## 4. Footer alignment and fonts

<- and -> mean left-center/right-center. Right alignment needs free space. Proposed single-row footer defaults: full width, justify=between, placing two content-sized texts at opposite edges. Explicit rules may override. Insufficient width follows overflow, not an undeclared layout switch.

Owner clarification, 2026-09-21: font=10/12 are absolute and retained on resize. Fonts affect measurement/region height; width changes can rewrap text/change height. At proposal time, shared native units were still to be specified; absolute sizing was settled. Proposed inheritance: parent base font, explicit child override. Current Go units are DIP.

Formatting follows any variant and precedes comma/semicolon, including strings, calls, groups and references. Comma advances horizontally; semicolon starts below the previous row. Regions follow section 2. [Layout examples](layout-language-proposal.md).

Ordinary Markdown usually requires a space in `## Markdown Heading`. The owner's `##Markdown Heading` is preserved above; SDUI must not silently repair it. Pass content to the provider unchanged.

## 5. G1–G4 verification cases

Canonical/equivalent corners; ratio with only x or y fill and rejection of both; preserved width under height overflow; regions inside ratio; equivalent explicit/implicit bodies separately and conflict when mixed; multiline headers; right footer reaching the edge; stable fonts/correct rewrap; all-component suffix formatting and detached-block rejection; three independent group levels/local separators; preserved single-child groups/ancestor references; separate instance state; cycles; unbound callbacks; new source/instance generations and stale-event rejection.
