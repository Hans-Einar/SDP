# SDUI — user interface prototyping

The [documentation map](docs/README.md) distinguishes current profiles, implementation, design background and dated references. [SDUI KanBan](SDP/Agents/KanBan/README.md) tracks ideas and local effects of SDP planning; cards do not change the implemented language profile.

SDUI 0.2 is implemented in Go: parser/AST, validation, normalization, relative layout, SVG, structural console/Markdown dumps, runtime and Fyne host with model reload. SDL bindings and Go generation share these models/runtimes. No active Python frontend or 0.1 compatibility path remains.

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

Limits: scroll layout is rejected; embedded Mermaid is verified for flowchart/graph with a registered renderer. Fyne without that provider displays placeholders. The SVG widget is a labeled placeholder. FOX widgets, TUI, full Markdown/Mermaid coverage and Ponsse production integration are not delivered. The older [HTML control gallery](examples/prototype-controls.html) is a static historical trial; new general UI placement uses Go layout.
