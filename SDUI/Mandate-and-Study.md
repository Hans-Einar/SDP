# SDUI — mandate and initial study

**Current direction, 2026-09-21:** develop SDL/SDUI parsers/runtimes in Go, Fyne as first interactive host, SVG export from shared models/layout. This supersedes earlier language/host choices below, retained as mandate history. [Checkpoint 07](../SDP/History/checkpoint-1/07-SDUI-0.2-and-Go-Direction.md) and [PLAN-003](docs/implementation-plan.md) own current deliveries.

**Implementation, 2026-09-22:** the Concept1 request now has [SDUI source, AST and console dump](docs/concept1-console.md). Go frontend, layout/SVG, Fyne, runtime and reload are delivered. [Checkpoint 11](../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md) bounds SDL execution, generation and navigation. Earlier sections preserve history, not another active parser profile.

**Clarification, 2026-09-21:** canonical corners: ^<, >^ at top; v<, >v at bottom. Outer `{16:9,<->}` fills width/derives height, expressing FixedAspectViewport without height-constrained contain fallback. Header/footer may hold frames/widgets; body= is optional. Named reusable groups come from the owner's mainBody/page example. font=10/12 are absolute; resize does not scale content. Formatting follows every component before separators; comma continues horizontally, semicolon begins below the preceding row. [Composition proposal](docs/frame-composition-proposal.md).

**Owner clarification, 2026-09-20:** frame ratio x:y; scale relative to nearest ancestor/root host area. Ratio allows one driven axis; otherwise x/y/both may scale. No source pixel width/height. Canonical mini-arrow shapes preserve direction-pair equivalence; down-left is v<. ¤ has no defined meaning. See layout proposal.

**Addition, 2026-09-20:** owner authorizes libsdui, runtime and layout/presentation development with phases/milestones. [Layout proposal](docs/layout-language-proposal.md) develops general frames, widget lists, Markdown and layout expressions. No old SDUI 0.1/BoxUI compatibility required: port examples/remove replaced paths. [Implementation plan](docs/implementation-plan.md) is current. The remainder preserves the initial parser foundation.

**ID:** SDUI-MANDATE-001 · **Revision:** 0.1 · **Date:** 2026-09-19. Owner request authorizes documentation, EBNF and an AST-generating parser. SDUI is a working name; no collision investigation/final registration performed.

## 1. Owner intent

Describe UI as a large outer box with named nested boxes/widget content, like Ponsse/Concept1 BoxUI. Placement should be predictable/schema-driven. A small language in Markdown describes layout and binds widgets to SDL-defined objects/functions.

`input1_boxui` is the widget; `sdlFile.input1_sdl` the SDL object. Callbacks travel from user events to SDL objects; `setHandle(BoxUIDefinition.input1_boxui)` supplies a logical reference for later updates. These are two directions, not raw FOX pointers or callbacks drawing directly on screen.

The owner supports a bounded independent language, separate from Mermaid. This initial delivery moves no existing implementation or XFMD installation. Later integration should reuse existing work.

## 2. Inspected foundations

| Source | Observation | Consequence |
| --- | --- | --- |
| Ponsse `882ad7c`, Concept1/shared/ui-box/model.mjs and UILayout.jsx | group/axis/weight/box-id, CSS Grid tracks; separate React content. | Preserve placement/content separation; full UI is not already one portable declarative schema. |
| Concept1/apps/operator-ui/src/ui/operator-layout.mjs | Nested 15/45/40 groups and stable box IDs. | Realistic future compatibility fixture. |
| Renderer fork `61a85b6`, src/boxui/layout.rs, svg.rs, model.rs | Independent row/column placement, text callbacks and SVG, without treemap layout. | Investigate extraction; treemap is not required. |
| XFMD `ffb98e8`, BoxUiAbi, BoxUiFrame, BoxUiSession | Parse/prepare/free ABI, typed model, native controls, local simulated participants; no SDL runtime. | Reuse identity/revision/draft concepts; renderer ABI is not SDL ABI. |
| Sibling SDL working document | Exploratory design model; interpreters need not execute product functionality. | Do not assume completed executable SDL or adopt its semantics here. |

The SDL working document had uncommitted local changes; treat as working evidence, not frozen approved contract. No new external libraries selected for this parser delivery.

## 3. Scope and decision

Select one bounded source profile, explicit AST format and separate local validation. Python standard library was chosen for a small inspectable reference prototype, not as a future UI-core/host technology decision.

The prototype supports boxes, text/button/input/svg rows, module references and declarative setHandle, executing none. Syntactically valid refs do not prove object/member existence.

A language is shorter than wire JSON but needs grammar, diagnostics, versioning and tools. A future shared model for SDUI/existing schemas remains a plan, not a completed adapter.

## 4. Open questions at the initial delivery

- Final name/profile namespace; retain working name SDUI.
- SDL signatures, object creation, asynchronous behavior and error model.
- First extracted UI-core target: Rust or another implementation.
- Detailed sizing/overflow and desired Concept1 geometry fidelity.
- Possible shared normalized IR with SDL, without necessarily changing source profile.

These did not block source-faithful AST parsing; they blocked executable-binding and compatible-rendering claims. Later direction/status is recorded above.
