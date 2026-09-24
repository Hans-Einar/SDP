# FOX/XFMD — historical handoff superseded by Go/Fyne

Updated 2026-09-22. The former FOX widget handoff is superseded. G6 documentation
navigation was implemented separately in [XFMD PR #38](https://github.com/Hans-Einar/xfmd/pull/38);
see [G6 evidence](../../SDL/go/evidence/G6.md). The owner selected independent SDL/SDUI
Go libraries, Fyne as first interactive host, and static SVG export.
[Architecture](target-architecture.md) and [PLAN-003](implementation-plan.md) replace
P4–P6/Rust/C-ABI instructions.

XFMD can display generated SVG in Markdown without an SDUI parser or SDL runtime.
Interactive SDUI inside XFMD is a possible later consumer requiring its own agreed
deliverable/port. Existing FoxBoxUiOverlay/FoxBoxUiInput, local simulation, diagrams
and Cairo/PDF code remain reference material in their repositories. G6 includes no
cleanup or merge of the old BoxUI worktree. Historical pins/reuse sources are in the
[worktree map](renderer-extraction-and-language-direction.md).

G6 started from verified clean XFMD main c245fd9 in xfmd-sdl-navigation. At this
handoff phases 049/050 were pushed and PR #38 was not merged. Do not start a FOX
backend, private parser or C ABI to follow superseded instructions. The first Go
deliveries are under SDUI/go and SDL/go.
