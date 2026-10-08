# WCI2-M1 — tabs, split geometry and native pane interaction

Implemented, integrated, verified and independently reviewed.
Full KB-SDUI-003 inventory remains in progress.

## Candidate

Phase `sdui/widgets-wci2`, commit `403c540`, follows WCI1 handoff `53d031e`.
[Inventory](candidate-WCI2-M1.json) records all 72 changed source/test/document
files. Raw reviewed inventory SHA-256:
`eac6779750aade0678f4fbbcd9307a1236b6bbd16b6aef2690475a49a51f1757`.
Frozen native binary SHA-256:
`f1de87cefd039f84321c145d02f4b79bbd219916a30b76dfaa3a2007600239cc`.

All 72 original-workspace targets were compared with the phase baseline before
integration; none diverged. Their integrated bytes match the inventory exactly.
Unrelated governance, preview and dependency changes remain preserved. No merge,
release or publication is part of this milestone.

## Behavior and proof

Exact unreleased profile 0.3 adds named tabs/page/split composition, retains source
provenance through code generation and leaves 0.2 behavior/output unchanged.
Runtime owns page intent/activity, remembered focus/drafts, hidden offsets and
split proportions/collapse. Shared measured layout owns legal split bounds.
The host uses native AppTabs with a bounded input adapter and a divider projecting
that geometry. Pure probes cannot publish pending presentation; accepted final
tickets alone promote resources after runtime publication.

Actual SDL Page actions receive stable previous/new page IDs. Silent initial,
programmatic, fallback and reload operations invoke no action. Native errors
retain succeeded/unknown domain outcomes when the UI cannot publish the reply.
Composition/text/codegen retain pane structure; unsupported public pane SVG and
standalone native capability are explicitly diagnosed.

| Contract | Evidence |
| --- | --- |
| A01 | Profile/schema/placement/reuse/provenance and compiled constructor tests; frontend report |
| A02 | Actual pointer/keyboard pages, retained drafts/scroll/focus, eligibility fallback, disabled native agreement and failure recovery |
| A03 | Both split axes, drag/keys, positive-minimum collapse, resize, rejected/successful restore and compatible reload |
| Pane portions of A07/A08 | Real typed SDL bridge; domain error/malformed output/reentrant draft/resource conflicts; six failed reload stages; provider revocation and late-result rejection |
| Pane portions of A09 | Structural exports and explicit SVG rejection, full race suites, consumer regression |

[Native archive](native/WCI2-M1-final/README.md) retains 58 passing actual XTest
checks: panes 31, horizontal split 11, vertical split 9, lifecycle 7. All stderr
files are empty. Coordinator inspected OS screenshots, including disabled headers
and restored vertical geometry. Unit/headless checks are separate evidence.

Frozen phase checks on Go 1.27.1 linux/amd64:

- SDUI/go `go test -race ./...`: PASS; host 73.780s.
- SDL/go `go test -race ./bridge ./runtime ./codegen ./examples/panes ./examples/collections`: PASS; panes 66.304s.
- SDPTool `go test ./...`: PASS; root 11.131s.

The [pilot history](native/Pilot-WCI2.md) records header focus/disabled-selection
corrections and native input barrier corrections. Final source remained frozen
through independent review. Reviewer `01a11854-523d-7c11-adf4-f622b19b68c6`
approved all 72 files, independently rebuilt the native binary byte-identically,
checked all 58 actual native assertions and all 23 archive entries. Original full
SDUI race and full SDPTool suites passed, including preserved governance/preview
work. Their logs are in the archive.

## Remaining inventory

M2 command/button/toggle/menu/context/dialog behavior is not claimed by this
milestone. WCI3 typed values/text/IME and WCI4 providers/consumer package preparation
remain required. Native evidence uses X11/Xvfb without a window manager. No full
widget inventory, external XFMD package, merge or release acceptance is inferred.
