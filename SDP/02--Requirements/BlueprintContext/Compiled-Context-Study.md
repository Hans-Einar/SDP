# Compiled agent context from design, code and runtime observations

Initial feasibility study, 24 October 2026. Owner discussion during Session0007
T029; primary proposal KB059. This records a possible direction, not selected
implementation, an adopted trace schema or a new executable SDL profile.

## Goal and recommendation

Prepare a compact, revision-bound context package before an agent starts a task,
using reviewed SDL design, code-derived structure and selected runtime evidence.
Give the agent useful orientation immediately, with precise MCP expansion handles.
Keep those three sources distinct: intended/modelled design, observed implementation
structure and observations from a particular run. Link them where a mapping is
established and report mismatches/unknowns instead of silently reconciling them.

This is feasible in stages. Reuse Blueprint generation, ModelGovernance identities,
PG's launch payload and existing SDL/action/runtime services. Do not start another
parallel model authority or replace approved TARGET from reverse-engineered code.

## What can be compiled or captured?

| Source | Useful derived artifact | Boundary |
| --- | --- | --- |
| Reviewed SDL NOW/TARGET and authored task | Purpose/responsibility orientation, relevant contracts, semantic changes, preserved obligations, pinned design views | A model can omit implementation details or be stale; explicit correspondence evidence is still needed. |
| Code/compiler/type information | Packages, symbols, types/interfaces, imports, source locations, selected dependencies/calls and build configuration | Static declarations/possible calls are not a complete runtime flow; callbacks, dynamic dispatch and external components can remain unresolved. |
| Instrumented runtime/test execution | Correlated requests/events, selected component states/revisions, before/after transitions, outputs/errors and timing observations | One run demonstrates one observed path, not every possible execution or design intention. |
| Explicit mapping between sources | SDL element ↔ code symbol ↔ observed event/state | Bind mapping to model/code/run revisions. Declared tags locate things; they do not prove correct behavior. |

Prefer relevant domain/component state to raw memory or machine-code dumps. Debug
symbols and compiler maps can help locate implementation, but purpose, ownership
and desired behavior require model/authored meaning. Keep exact source provenance
and label any generated narrative interpretation separately from compiler facts.

## Context delivery through MCP and PG

MCP exposes data/tools; merely publishing a resource or tool does not ensure its
content is present in the model's first context. The PG/client handoff should
resolve the selected task and exact revisions, retrieve/build the small initial
package, and explicitly include that package in the assignment input. Expose
read-only detail retrieval for the remaining pinned model/code/trace views. This
reduces discovery work; initial text still consumes context tokens.

The current PG controller already constructs a bounded launch payload with owner
input, assignment, coverage and optional approved-design/attention context and sends
it through StartTurn/StartSkillTurn. That is an injection point, not a finished
Blueprint context compiler. Current launch payload has a 64 KiB bound; do not
embed entire model/code/trace inventories or silently truncate required context.
Current worker MCP has seven governance tools and does not expose a complete
Blueprint/code/trace resource surface.

Proposed initial packet: task and authority, short system orientation, selected
NOW/TARGET identities and change, protected boundaries, build/code correspondence,
important unknowns, available evidence and exact detail handles. Retrieved evidence
remains data and cannot widen the assignment or establish owner approval.
Client input/handshake support is documented in the
[official app-server guide](https://learn.chatgpt.com/docs/app-server);
local controller source defines the actual currently implemented payload.

## Logical state views and stepping

A useful observed view could highlight: request accepted by component A, message
sent over a mapped Channel, component B transition from revision N to N+1, returned
result and a subsequent UI update. A cursor moves through recorded correlated
observations while the corresponding model elements and selected state are shown.
This is trace replay/inspection, not execution of every possible path or automatic
recovery of unstored state.

Across concurrent processes, a wall-clock timestamp alone does not define a
consistent whole-system state. Capture component identity, per-component order,
request/causal correlation and observation gaps. A combined point-in-time view must
state its consistency/coverage; do not call independently sampled states an atomic
global snapshot. Deterministic re-execution needs additional control over inputs,
clocks, external effects and recorded state, beyond an ordinary debugger snapshot.

Stepping a **design simulation** is different: it requires defined transition,
input/output, ordering/error and state semantics. Structural design-core facts and
Channel/scenario sequences are not automatically executable. A simulated TARGET
trace must remain distinguishable from observations of current production code.

## Existing foundations and gaps

- [Blueprint contract/selection](../../04--Design/SDPTool/Blueprints/Contract.md)
  supplies NOW/TARGET, affected context and unknowns. The
  [owner-fit assessment](../../04--Design/SDPTool/Blueprints/Owner-Intent-Assessment.md)
  proposes a layered entry and progressive delivery.
- [KB004](../../KanBan/backlog/%23004--Proposal--Design-traceability.md) owns
  code-tag/mapping identity and evidence. It has not delivered a general automatic
  code-to-model conformance verifier.
- [Action-core 0.1](../../../SDL/docs/profiles/SDL-Executable-Action-Profile.md)
  already binds typed Actions to registered Go functions. The runtime exposes
  invocation/action/model/sequence identity; Go handlers own domain state. This
  is a useful first seam for explicit mapped events/state projection, not general
  interpretation or single-stepping of arbitrary Go or SDL Functionality.
- [KB-SDL-006](../../KanBan/backlog/%23006--SDL--Study--Executable-channel-tests-and-unit-bindings.md)
  owns unselected research into Channel/Unit handler mapping and behavioral checks.
  Structural scenario diagrams currently describe authored example paths, not
  instrumented execution.
- Current resource/budget/provenance and PG recovery services should be reused.
  A new context/trace feature must not silently expand worker control privileges
  or trigger application execution while browsing.

## Smallest proposed route

1. Layer one existing Blueprint into a useful initial packet and explicit pinned
   detail handles. Verify that a new agent can identify task, context, invariants
   and unknowns. No runtime capture required for this first useful increment.
2. Add one language's compiler/type-derived code map for the selected scope and
   build configuration. Review stable model/symbol mapping and expose drift/gaps.
3. Observe one bounded existing typed action/Channel path with selected state and
   event correlation. Render/step that recorded path through the mapped design;
   compare expected versus observed transitions without inventing a full system map.
4. Consider design simulation or broader distributed state only after required
   semantics and actual owner use justify them. Reuse the separate executable-profile
   study rather than treating source formatting as executable behavior.

The suggested acceptance is reduced exploratory source reading and a correctly
explained/implemented task, with known gaps, not elimination of every source read.
Full-system snapshots, universal call graphs and a general execution engine are
not prerequisites for a first useful package.

## Disposition

Register as KB059 backlog and coordinate KB050/KB004/KB-SDL-006. The owner asks
whether this is possible; no code/runtime instrumentation, new SDL grammar, snapshot
collector or model/Ponsse execution is selected. Continue the already authorized
PG Codex compatibility increment. Revisit this study at Blueprint owner review.
