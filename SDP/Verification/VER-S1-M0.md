# VER-S1-M0 — system decisions and XFMD process handoff

Date: 2026-09-25. Phase branch: sdp/phase-s1-language-source-organization.
Base: 1a3f5e8. Scope: adopted design direction and coordinated process bootstrap.

The owner selects SDL, SDUI and SDPTool as three software systems sharing root
SDP, with XFMD as a separate collaborating system and separate SDP project.
[System study](../03--Architecture/System-Boundaries-study.md) records the decision
and superseded single-system recommendation. Existing SDL/SDUI board locations
remain until the planned consolidation; no source split has run. System syntax
is not implemented by this delivery. Sourceinput Go files remain unfinished,
untested drafts outside the milestone commit.

## Verification

The three SDP boards replay 25 cards/113 events with schema/placement agreement
and preserved prefixes from 1a3f5e8. Documentation checks retain 105 frozen
records/ledger prefixes and 574 generated outputs; 2,252 local links/130 anchors pass.
Incoming XFMD card links use its new SDP/Agents/KanBan location.

XFMD's SDP1-M1, commit fadf90c on maintenance/phase-sdp1-process-bootstrap, initializes five phase entry points and
migrates/aligns its board to the pinned SDP baseline. It preserves all 29 original
board files and historical ledger/non-Markdown bytes. Its 17 cards/69 events
pass schema/replay/state/index checks; lineage cases and 15 negative examples,
37-blueprint/70-requirement checks and 279 implemented-callee checks pass.
No XFMD product code or generated .design files are claimed. Its next-agent card
is KB-XFMD-017; KB-XFMD-012 retains future version-release alignment.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Maintenance/SDP1/Plan-and-Evidence.md

## Remaining work

S1-M1/M2 still own explicit source inputs and source organization. Adopt System
through a versioned semantic change; do not silently reinterpret existing Unit or
Container declarations. Full XFMD design adoption belongs to its next agent.
