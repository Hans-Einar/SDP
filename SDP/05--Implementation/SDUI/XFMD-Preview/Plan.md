# SDUI combined preview and host preflight — XFMD 0.7 work unit

Execution selected by owner XFMD Session 0004 T002, 2026-10-07, under
PLAN-XFMD-0010. This is the producer work unit of that authoritative plan, not a
second lifecycle plan. KB-SDUI-005 owns producer intent. Keep the existing branch
and milestone policy; preserve concurrent governance changes. No upstream release
or gh extension replacement is selected. Exact tested source is captured with
XFMD release evidence before local installation.

Implement presentation.Combined using normalized containment and parser provenance,
while retaining Markdown's existing text output. Add instance declaration/use
metadata without changing grammar. CLI preview/check owns files and hashes;
parser remains free of I/O. The standalone Fyne host shares revision/preflight
checks and explicitly labels unbound SDL callbacks. XFMD ships these binaries
privately; SDPTool UIPreview calls the same combined presentation library.

Verify deterministic grouping, hidden/reused instances, escaped labels, UTF-8
spans, runtime/layout validation, stale revision rejection, missing entry and
native Fyne launch/close. Existing language/runtime tests remain required.
Embedding is not selected: no established FOX/Fyne parent bridge.

External Session:
file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230004--SDUI-workflow.md

Final notation: Mermaid mindmap overview for each frame/group, region and ordered
row branches, with nested-container details in the same page. Native nested
flowchart rectangles overlapped labels, so that representation was rejected by
visual inspection. No source nodes are omitted from provenance/text ownership.

## Delivery — 2026-10-07

Implemented producer projection/provenance, revision/readiness CLI and labelled
local Fyne host. SDUI race tests, SDPTool facade tests and actual XFMD native/PDF/
Fyne acceptance passed. Exact source/binary hashes and package installation are
in XFMD U4-001. Helpers ship privately with XFMD0.7; no gh extension or upstream
release was installed. KB005 awaits review; KB003/004 remain separately owned.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDUI/evidence/U4-001/README.md
