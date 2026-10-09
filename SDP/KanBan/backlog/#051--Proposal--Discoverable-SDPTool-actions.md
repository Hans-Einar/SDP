# Discoverable SDPTool actions for native clients

| Field | Value |
| --- | --- |
| id | KB-SDP-051 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-10-09T09:44:29.345109+00:00 |
| source | Owner Session0008 T013, 2026-10-09 |
| next_review | Before selecting the XFMD dynamic-action integration plan |

## Outcome and current evidence

Expose supported operations programmatically so XFMD can build an SDPTool menu
and configurable toolbar as tools evolve. Existing Go services and CLI JSON form
the execution foundation. Project navigation/capabilities describe content; they
do not currently describe all callable actions or their parameters.

[Integration handoff](../../../SDPTool/XFMD-Blueprint-Integration.md) records exact
current commands, source inspection, protocol boundaries and branch testing.
Related local cards: KB-SDP-017 and KB-SDP-050. External dependency:
`external:KB-XFMD-030`, dynamic menus/preferences; its checkout is not required.

## Proposed scope for a bounded design plan

- Versioned action catalogue from the same operation definitions used by execution:
  stable IDs, label/group/description, icon key, typed input/output contracts,
  defaults/context bindings, mutation class and supported protocol versions.
- Advertised support separate from contextual availability and caller authority.
  Revalidate actual selection/revisions, prerequisites and permissions at invocation.
- Typed invocation shared by CLI and adapters; retain existing command compatibility.
  Choose the exact catalogue/invocation command shape during design. No help scraping
  or arbitrary shell snippets as an integration API.
- Deterministic catalogue, validation, unknown parameter/action errors, refresh when
  the selected runtime changes, and explicit unsupported-schema behavior.
- First vertical workflow: blueprint generation/retention, listing, assessment and
  revision-bound lifecycle. Mutations preserve canonical history and retry gates.
- Consumer fixture tests: older producer, new action with generic icon, unsupported
  input type, stale selection, rejected mutation and stable personal preferences.
- Document repeatable branch testing: direct executable for XFMD, then a signed
  test descriptor through gh-sdp with isolated cache/config. Track wrapper and engine
  identities separately. Alpha/beta publication remains an explicit release action.

## Ownership and limits

SDPTool owns action metadata and execution validation; XFMD owns menu placement,
Preferences/ToolBar, icons and personal config. Generic forms support only a defined
input vocabulary; new interaction semantics may need more than a new icon.
Do not add native XFMD code or silently introduce HTTP/MCP authentication here.
Coordinate eventual MCP exposure with existing governance work, without duplicating it.

## Worklog

2026-10-09: registered from owner request; no implementation selected. Blueprint
pilot KB050 remains gate-review independently of this proposal.
