# SDL tools — Go entry point

G5-M4 replaced Python CLI/projectors after frozen port comparisons. One design-core 0.5 frontend and one viewpoint projector remain in the [Go module](../go/README.md). Historical Python code is in Git before G5-M4, not a working-tree fallback.

From the SDP root:

```sh
go -C SDL/go run ./cmd/sdl check ../../SDUI/design/architecture.design
go -C SDL/go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --format static --monolithic --output ../../SDUI/design/viewpoints --renderer /absolute/mmdr
go -C SDL/go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --output ../../SDUI/design/navigation --project sdui
```

Default export creates navigator/overview without detailed rendering. Static export creates viewpoint directories and an optional combined file. --renderer selects registered mmdr; otherwise portable Mermaid blocks remain. SVG uses consistent actor/use-case/Feature symbols and explicit source links.

[Selection/IPC/leases](../go/README.md), [G6 evidence](../go/evidence/G6.md), [abstraction levels/notation](../docs/integration/SDL-Viewpoint-Levels-and-Notation.md). Dataset/database are logical data sources; database does not necessarily mean SQL. Mode is allocation/operating context. A separate State machine is not adopted in design-core.
