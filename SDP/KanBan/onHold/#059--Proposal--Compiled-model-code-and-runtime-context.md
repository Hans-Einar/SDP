# Model-linked code and runtime context for agents

| Field | Value |
| --- | --- |
| id | KB-SDP-059 |
| project | SDP |
| type | Proposal |
| CardState | onHold |
| created | 2026-10-23T22:24:47.632367+00:00 |
| source | Owner compiled-context/runtime-state question during Session0007 T029 |
| Systems | SDL, SDPTOOL, ProjectGovernance |
| next_review | After initial SDL-only use, only if the owner explicitly reselects this extension |

## Need and proposal

Explore prebuilt agent context from SDL/Blueprints, compiler-derived implementation
facts and mapped runtime observations/state. The client injects a compact package;
MCP supplies pinned detail expansion. Consider trace stepping and, separately,
explicitly defined design simulation. Keep modeled, code-derived and observed facts
and their revisions distinct.

[Initial feasibility study](../../02--Requirements/BlueprintContext/Compiled-Context-Study.md)
records current foundations, limitations and proposed increments. Coordinate existing
KB050 Blueprint context, KB004 model/code mapping and KB-SDL-006 executable bindings;
this card owns the additional cross-source context/runtime proposal, not their
already assigned implementations. No alternative model/management authority.

## Next action and acceptance

Discuss the proposed smallest path at Blueprint review; select a bounded plan only
if the owner chooses it. No new parser/runtime, trace collector, global state snapshot,
MCP API or Ponsse/model execution is authorized by this registration. Compatibility
PGC1 remains the current implementation.

## Worklog

2026-10-23T22:24:47.632367+00:00: EVT-KB-SDP-000393 registers the owner question and source-backed initial
study; backlog state. No implemented context compiler or automatic preload claim.

## Owner deferral — 2026-10-23T23:23:59.095012+00:00

Owner restricts current context/Blueprint work to SDL and prioritizes finishing
the existing system for actual use. Code-derived context, runtime snapshots/traces
and new simulation semantics are deferred; they are not prerequisites or selected
implementation. Restart only after SDL-based use and explicit owner selection.
Preserve this proposal for later learning. CardState onHold; EVT-KB-SDP-000398.
