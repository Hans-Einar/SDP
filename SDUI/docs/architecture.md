# SDUI — implemented Go architecture

Updated 2026-09-22. One active SDUI 0.2 frontend. Executable modules/commands are in the [Go area](../go/README.md); the shared design is [described in SDL](../design/README.md).

| Package | Responsibility |
| --- | --- |
| go/parser | Lexer, recursive descent, AST/spans, local rules and normalization |
| go/layout | One measured geometry; relative dimensions, rows, wrap, ratio and clipping |
| go/markdown | Bounded Goldmark content; measurement and registered Mermaid provider |
| go/svg | Static export using shared geometry, Go Regular glyphs and native-control appearance |
| go/presentation | Structural console/Markdown dumps and bounded control gallery |
| go/runtime | Session, handles, accepted/draft, typed events and atomic property batches |
| go/reload | Validated candidates and compatible state/identity preservation |
| go/host/fynehost | Native input/button, focus/keyboard and UI-thread publication |
| go/codegen | Independent typed Go constructors for Document/Root |
| go/cmd | CLI and native composition; file access belongs here, not in parser |

Parse → Normalize/Compile produces Document and expanded Instance trees. AST JSON is tagged sdui-ast/0.2. Spans are half-open UTF-8 byte ranges with one-based Unicode line/column positions. Nodes preserve groups, rows, regions and suffix formatting. Runtime copies input models; Go structs are not language-level immutable. Do not mutate models in use by a host.

Explicit widget names produce public instance paths; anonymous segments use synthetic names. Session/Path/Generation/Kind identify handles. Reload preserves compatible named instances, accepted/draft and focus; deletion/type changes invalidate old handles. Source/binding failures retain the last valid model. Source watchers publish candidates through fyne.Do.

Layout receives the host's available area; source has no pixel width/height. Children use their nearest source ancestor; `{16:9,<->}` derives height from filled width. Fonts remain logical DIP during resize. SVG and Fyne share rectangles, text measurements and clipping; native controls use host rasterization/theme.

The SDL adapter lives in SDL/go/bridge. Parser ref/callback/setHandle are data; composition registers SDL modules, Go functions and typed bridge.Plan. Action-core 0.1 provides explicitly bounded execution. An accepted Go domain transaction cannot roll back if later UI publication fails; there is no automatic replay. See the [runtime contract](runtime-contract.md).

Limits and actual trials: [G1](../go/evidence/G1.md), [G2](../go/evidence/G2.md), [G3](../go/evidence/G3.md), [G4](../../SDL/go/evidence/G4.md), [G5](../../SDL/go/evidence/G5.md). No alternative SVG/Fyne/XFMD parser, extracted Rust crate or mandatory C ABI.
