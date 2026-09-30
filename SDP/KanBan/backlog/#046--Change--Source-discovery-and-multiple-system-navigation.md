# Source discovery and navigation across systems

| Field | Value |
| --- | --- |
| id | KB-SDP-046 |
| project | SDP |
| type | Change |
| CardState | backlog |
| created | 2026-10-01 |
| source | Owner report in Session 0002: SDL/SDUI absent despite source files in XFMD |

## Observed workflow gap

The owner expects discover/tree to recognize SDL and SDUI under SDP/SDL for both
single-system and multi-system projects. Current navigation.json has explicit
models and sdui arrays; empty arrays report absent even when files exist. The
new human-output layer does not change discovery semantics.

XFMD's registration has both arrays empty. Its source README identifies
SDP/SDL/XFMD/XfmdDesktop.design and SDUI/XfmdDesktop/Desktop.sdui as current inputs;
Navigation.design and Navigation.sdui remain historical reproduction fixtures.
A blind glob must not promote every .design/.sdui file into current authority.
model-map.json is local verification metadata, not a language entry manifest.

## Verified boundary

Temporary-copy registration of the two current XFMD inputs succeeds with the
unreleased SDPTool candidate: discover declares both capabilities; tree has 2454
nodes, SDL validated, SDUI available. The live XFMD project is unchanged.
Existing tests cover design-core/0.6 source composition. One compiled source graph
contains one System. Multiple registered models are supported, but tree chooses
one SDL model via --model/defaultModel; there is no aggregate systems tree.
Published engine 1.0.0 must not be credited with unreleased 0.6 support.

## Proposed next bounded work

Select a design/implementation plan for source-aware discovery and explicit
system/entry selection, then a unified multi-system navigation inventory. Respect
source-owned includes/contains, the flat System layout and ecosystem grouping;
distinguish roots, fragments, shared code, historical fixtures and SDUI entries.
Report unregistered/incompatible source candidates without claiming they are
absent or validated. Decide how existing explicit registrations interact with
convention-based discovery; avoid a new mandatory design-set manifest.

Acceptance should cover one and multiple Systems, nested ecosystem groups,
invalid/ambiguous candidates, source exclusions, stable IDs/revisions, default
selection, SDUI ownership and existing explicit registrations. Parsing/UI-runtime
integration and automatic native XFMD changes are not implied by this card.

## Provenance and handoff

[Session 0002](../../Sessions/session-%230002--SDPTool_output.md), T003, records the
owner report and read-only diagnosis. [KB045](../completed/%23045--Change--Portable-SDPTool-presentation.md)
owns output presentation only. This card is registered, not selected for execution.
The temporary evidence path is in the Session; canonical behavior is in the
[producer contract](../../../SDPTool/Contract.md).
