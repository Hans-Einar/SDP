# Checkpoint #1 — SDUI 0.2 and shared Go direction

**Implementation status is updated in [supplement 11](11-Go-Implementation-and-Navigation.md).** This preserves dated design/V-phase foundations; old Python commands are historical and replaced by Go entry points.

Updated 2026-09-21 after owner decisions. Corrects the 18–20 September direction without changing checkpoint number. SDL is **SystemDesignLanguage**; “System Description Language” is its historical name, not another language. SDUI is independent. Neither renames Standard Document Procedure Toolkit contracts.

**Addition, 2026-09-22:** the [combined parser/runtime design](../../../SDUI/design/README.md) is described in SDL with parser-produced AST/validation. This establishes structural Go direction; Go implementation was still pending at this snapshot.

## Decisions

| ID | Choice | Consequence |
| --- | --- | --- |
| CP1-D18 | SDUI 0.2 frames, groups, Markdown, relative layout | Python frontend ported; no active 0.1 compatibility |
| CP1-D19 | Go for subsequent parsers/runtimes | Replaces planned Rust/C-ABI core; Python provides port evidence until replacement works |
| CP1-D20 | Independent UI, first interactive host Fyne | Core needs no FOX/XFMD; no parallel FOX/TUI now |
| CP1-D21 | Shared model/layout for interaction/SVG | Markdown embeds generated UI images without making them interactive |
| CP1-D22 | Validated-model hot reload | Retain last valid model/compatible state; reset/migrate incompatible changes |
| CP1-D23 | Later Go generation with handwritten domain functions | Separate generated/handwritten code; shared development/production runtime |

D19–D23 select direction, not implemented Go capabilities at this date. Executable SDL subsets, Go interfaces and reload rules still need definition. Technology choices do not adopt ControlSet, Function/Functionality or Activity candidates.

## Actual artifacts at this snapshot

| Area | Status |
| --- | --- |
| SDL | Python design-core 0.5 structure, data, Channels, scenarios, plan facts; generated G1–G5 design in supplement 09 |
| MVP1 | 66-file candidate exercise, not design-core executable |
| SDUI | Python 0.2 parser, source-positioned AST, validation, normalization |
| Console/Markdown | Generated structure dumps, raw console Markdown, separately renderable Markdown content |
| SVG | Button/input drawings and bounded bucking composition from saved treemap geometry/fixture code |
| HTML | Local editable/clickable/printable widget gallery, no SDL execution |
| Go | Parser/runtime directories only; no code/modules/Fyne/runtime yet |

The SVG trial is **not a general layout engine**: it draws selected source controls/sample data. General engines must measure/place arbitrary supported content, including formatting/clipping.

[Prototype evidence](../../../SDUI/evidence/prototype-widgets/README.md) records 36 tests/browser trials, not SDL runtime, Fyne, full Markdown/Mermaid composition or complete system execution. Old hashes/logs concern dated baselines.

## SDUI 0.2 in one page

- [] frames; *box/*b decoration.
- <> nestable widget/Markdown groups with formatting.
- header/body/footer may contain composed components; body is optional.
- Comma continues horizontally; semicolon starts rows; formatting precedes separators.
- Content strings are Markdown; button/input/svg are calls with ordinary labels. Parsing executes no SVG refs.
- Dimensions reference nearest source ancestor/root host area; no pixel width/height; fonts absolute.
- Ratio 16:9 permits one driven scale/fill axis. Width fill derives height without contain fallback.
- Canonical corners ^<, >^, v<, >v; reversed direction pairs equivalent. ¤ undefined.
- ref/callback/setHandle are symbolic AST data, not calls.

```text
sdui 0.2;
mainBody = <name=input("Name", value="Ola"), ok=button("OK")>;
page = [header="## Prototype", body=mainBody, footer="No domain binding"]*b {16:9,<->,font=12};
```

Details: [language](../../../SDUI/docs/language.md), [EBNF](../../../SDUI/grammar/sdui-0.2.ebnf), [layout proposal](../../../SDUI/docs/layout-language-proposal.md). Parsed rules do not prove executed geometry.

## Implementation and hot reload

Parser produces models; runtime owns identity/state/events; layout measures/places; Fyne presents; SVG exports. SDL has separate parser/model/bounded runtime. Typed Go interfaces/registered domain functions connect them without mandatory C ABI. Libraries must test without GUI; parser/runtime cannot import Fyne.

UI changes parse/validate, prepare model/layout and publish. Errors retain valid views; compatible identity/type may preserve fields/focus. SDL replacement occurs at defined event boundaries with explicit state migration/outstanding operations; reload never repeats domain calls.

Initial Go changes rebuild/restart. Build caching is not live machine-code replacement; plugins are not the reload foundation. [Go caching](https://pkg.go.dev/cmd/go#hdr-Build_and_test_caching), [plugin contract](https://pkg.go.dev/plugin). Later process replacement requires measured need; none created here.

Initial generation constructs the same validated models/bindings as development. Handwritten Go provides domain functions behind contracts. Incomplete semantics must diagnose rather than generate empty successful functions or guessed machine logic.

## Further reading and next delivery

- [Target architecture](../../../SDUI/docs/target-architecture.md), [shared plan](../../../SDUI/docs/implementation-plan.md).
- [SDUI Go](../../../SDUI/go/README.md), [SDL Go](../../../SDL/go/README.md).
- [Bucking source](../../../SDUI/examples/concept1-bucking.sdui), [SVG in Markdown](../../../SDUI/examples/concept1-bucking.widgets.md).
- [Structural profile](../../../experiments/design_core/README.md).
- [MVP1](../../../experiments/mvp1_sdl/README.md), [System.design](../../../experiments/mvp1_sdl/SDL/MVP1/System.design), [EditAptCell](../../../experiments/mvp1_sdl/SDL/MVP1/Scenarios/EditAptCell.design).

First Go delivery: frame/button/input, shared layout/SVG/Fyne/reload and registered Go function. Then one explicit SDL binding before broader executable grammar/generation. This update delivers documentation, structural model and directories only.
