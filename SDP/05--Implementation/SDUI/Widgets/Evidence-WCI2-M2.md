# WCI2-M2 — shared commands, native menus and composed surfaces

Implemented, integrated, verified and independently reviewed. Full KB-SDUI-003
remains in progress for WCI3/WCI4.

## Candidate and integration

Phase commit `0fc15c8f85b09ed334ba50b2f28f537f85d2cab2`, parent `a4f2c224`, contains
113 source/test/document files. [Candidate](candidate-WCI2-M2.json) identifies every
file. Tested inventory SHA-256 `34d3dda34ff2274aab46c66959390a4e9a3812b53c5e3ceba6ca36068bcd8710`
is retained separately in the native archive. Only architecture/requirements status
prose changed after approval; no product code changed. Frozen native binary SHA-256
`2c1db6d173c67ff873da4a7fdadb2a32722d3153705523622328e0d509be7806` was independently
rebuilt byte-identically. All 113 files match the safely integrated original tree;
all 168 SDUI Go/module files match across trees. Unrelated governance and KB005
changes remain preserved. Four module files add only pinned fancyfs v0.0.1, no upgrade.

## Behavior and proof

Unreleased 0.3 adds explicit shared command/button/toggle identity, keyboard and
native menu/context invocation, composed modal/nonmodal dialogs and exact published
opening results. Legacy basic button behavior remains intact. Source Close/Cancel
retains user sequence and modality; native owner loss uses exact lifecycle revocation,
parent-closed and sequence zero. Runtime owns drafts, accepted state, outcome and
replay protection; native resources publish only from accepted presentation tickets.
Actual SDL Run/Toggle/Inspect/Save handlers establish typed binding and domain effects.
Public structural/text/codegen output preserves identity; unsupported SVG declares
its boundary, and native background preparation uses the actual per-canvas inventory.

[Native archive](native/WCI2-M2-final/README.md) retains 118 passing assertions in
12 completed native workflows and one terminal result for each of 26 published
openings. All completed-run stderr is empty. The dynamic root-sibling modal dialog
uses its actual nonmodal parent canvas and size; closing that parent drains both
openings. Source Close/Cancel, nested Escape, focus restoration, disabled/stale command
refusal, menu dismissal, five reload failures and unknown/succeeded action failures
have distinct observable assertions. Coordinator inspected actual OS captures.

Final clean-exit suites (Go 1.27.1 linux/amd64):

- SDUI/go: `go test -mod=readonly -race ./...`, PASS, host 109.795s.
- SDL/go: `go test -mod=readonly -race ./bridge ./runtime ./codegen ./examples/commands ./examples/panes ./examples/collections`, PASS.
- Original SDPTool: `go test -mod=readonly ./...`, PASS, including concurrent preserved work.

Independent reviewer `01a11854-523d-7c11-adf4-f622b19b68c6` approved the full bounded
milestone: all 113 files, byte-identical build, 70 archive entries, 118 assertions,
26 exact receipts, full suite logs and inspected native failure state. Review found
and closed actual dynamic-parent ordering and native parent-close blockers. Seven
additional independent keyboard regressions passed. Earlier 58 pane checks remain
valid for the unchanged pane behavior on their separately identified predecessor
binary; they are not mislabeled as final-binary runs.

## Limits and remaining work

Predecessor 22539389 is explicitly not final acceptance. Its preferences-load EOF
and clean same-binary repeat are preserved; no causal claim or dependency patch is
made. A final multi-run process interruption (143) is preserved but excluded from
pass counts; the remaining three scenarios completed individually. Startup minimum
rejection and injected diagnostic events are expected and retained. X11/Xvfb without
window manager does not prove Wayland or decoration-click behavior.

Backlog review retains glyph fidelity as KB-SDUI-004 and combined-preview dependency
as separate KB-SDUI-005 gate-review. WCI3 values/native text and WCI4 provider/package
preparation remain required. No merge, publication, release or full-card completion.
