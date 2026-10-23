# SDUI — user interface prototyping

The [documentation map](docs/README.md) distinguishes current profiles, implementation, design background and dated references. [SDUI KanBan](../SDP/KanBan/README.md) tracks ideas and local effects of SDP planning; cards do not change the implemented language profile.

SDUI 0.2 and the selected 0.3 development profile are implemented in Go: parser/AST, validation, normalization, relative layout, SVG, structural console/Markdown dumps, runtime and Fyne host with model reload. SDL bindings and Go generation share these models/runtimes. No active Python frontend or 0.1 compatibility path remains.

From the SDP root, with Go 1.26+ (verified using 1.27.1):

```sh
go -C SDUI/go test ./...
go -C SDUI/go run ./cmd/sdui ../examples/concept1-bucking.sdui --format dump --entry bucking
go -C SDUI/go run ./cmd/sdui ../examples/concept1-bucking.sdui --format svg --entry bucking -o /tmp/concept1.svg
go -C SDUI/go run -tags desktop ./cmd/sdui-fyne -entry bucking ../examples/concept1-bucking.sdui
```

Linux desktop requires OpenGL/X11 and a C compiler. Parser/runtime work without a GUI. Fyne is the first interactive host; XFMD displays generated documentation and is not required by the UI core. [Go entry points](go/README.md); [shared execution example](../SDL/go/README.md).

The [Concept1 source](examples/concept1-bucking.sdui) has six main boxes and representative controls. These are example data, not ported React/bucking logic. The [AST](examples/concept1-bucking.ast.json), [console dump](examples/concept1-bucking.dump.txt) and [Markdown dump](examples/concept1-bucking.dump.md) preserve structure. [UI/state documentation](design/runtime-preview/entry.md) shows SVG from shared layout and explicitly selected state. Norwegian UI labels in this example are intentional localized sample data.

`[]` is a frame; `<>` nested groups; `*b` BoxUI decoration; `{}` formatting. Comma continues horizontally; semicolon starts a new row. Source dimensions are relative; fonts are absolute logical DIP. The parser never opens SDL refs or runs callbacks. Domain calls require explicit host/bridge registration.

- [Language/EBNF](docs/language.md), [layout profile](docs/go-layout-contract.md), [Markdown profile](docs/markdown-provider.md).
- [Architecture](docs/architecture.md), [runtime](go/runtime/README.md), [Go generation](docs/go-generation.md).
- [Requirements](docs/requirements.md), [milestones](docs/implementation-plan.md), [dated checkpoint](../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md).
- [SDL design/generated viewpoints](design/README.md), [mandate](Mandate-and-Study.md).

Source 0.3 adds collections and shared scroll viewports, tabs/splits, commands and
menus, composed dialogs, typed scalar fields, native text editing and explicit
prepared SVG/Markdown previews. Connected applications use DocumentHost with their
own providers and typed SDL registrations; the staged examples under SDL/go/examples
show the runnable routes. The standalone sdui-fyne helper retains its basic
RuntimeView adapters and reports unsupported new families explicitly. Parsing or
structural documentation does not imply connected readiness.

Explicit previews support a closed supplied SVG shape subset and bounded Markdown
prose. Missing or unsupported content follows each declaration's label/reject
policy. Native diagrams embedded in Markdown remain unavailable and receive an
explicit per-diagram fallback; legacy SVG placeholders remain unchanged. Public
SVG export rejects explicit resource previews without replacing existing output.
Native preview acceptance uses FYNE_THEME=light; dark/system-theme contrast and
Linux OS screen-reader delivery are not verified. Full Unicode/glyph fidelity is
tracked separately in KB-SDUI-004.

[Widget plan and evidence](../SDP/05--Implementation/SDUI/Widgets/Plan.md) and the
[acceptance matrix](../SDP/04--Design/SDUI/Widgets/Acceptance.md) record the actual
milestone state. Prepared helper payloads do not imply an installed XFMD upgrade.
FOX widgets, TUI, a full editor/browser, arbitrary Markdown/Mermaid rendering and
Ponsse production integration are not delivered. The older
[HTML control gallery](examples/prototype-controls.html) is a historical trial;
current UI placement uses shared Go layout.
