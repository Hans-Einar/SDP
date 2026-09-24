# SDUI repository guidelines

## Purpose and authority

Read README.md, Mandate-and-Study.md, docs/language.md, docs/architecture.md,
docs/implementation-plan.md, docs/layout-language-proposal.md,
docs/target-architecture.md and affected requirements before changes. SDUI 0.2
is the current Go profile. The Python port is complete; do not reintroduce an
old frontend or fallback. Distinguish adopted profile rules, proposals and
actual test evidence. Preserve unrelated work and follow the root's phase
branch and milestone rules.

## Architecture and files

`grammar/` owns EBNF. `go/parser` owns tokens/AST, spans, local rules and
normalization. `go/layout` owns measured geometry; `go/markdown` and `go/svg`
own content and static export. `go/runtime` owns UI state; `go/host/fynehost`
owns native controls and the UI thread. `go/presentation` provides structural
dumps. `go/codegen` generates constructors for the same model. CLI/I/O belongs
under `go/cmd`. Parser/runtime must not import GUI or renderer dependencies.
No parser opens SDL refs or executes callbacks. Keep one responsibility per
file; consider splitting files around 300 lines. Document actual calls and owners.

## Change order and verification

Requirements/profile → EBNF/AST contract → parser/validator → examples/negative
cases → verification. Version incompatible changes and port active consumers
together; versioning does not require backward compatibility. Do not reinterpret
existing expressions.

All documentation, including KanBan, must be in English, as required by the root
AGENTS.md. Go 1.26+, verified with 1.27.1; Fyne is the first interactive host.
SVG and Fyne share measured geometry. No extra C ABI, FOX or TUI host without
concrete new scope. From SDUI/go: `go test -race ./...`. AST:
`go run ./cmd/sdui ../examples/main-page.sdui -o ../examples/main-page.ast.json`.
Test source positions, limits, error paths and preserved structure/state. Native
tests use a separate display server and user area. Passing parser tests do not
prove GUI, SDL-runtime or domain behavior. Preserve frozen port fixtures and
dated machine evidence; Python automation for external processes is allowed,
but no Python language implementation.
