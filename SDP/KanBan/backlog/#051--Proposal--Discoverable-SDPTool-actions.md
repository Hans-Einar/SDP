# Discoverable SDPTool actions for native clients

| Field | Value |
| --- | --- |
| id | KB-SDP-051 |
| project | SDP |
| type | Proposal |
| CardState | queued |
| created | 2026-10-09T09:44:29.345109+00:00 |
| source | Owner Session0008 T013, 2026-10-09 |
| next_review | Session0011 S1 contract and implementation-plan selection |

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

## T014 — JSON request direction (2026-10-09)

Owner confirms subprocess JSON responses are sufficient and proposes JSON inputs
as the programmatic API. Recommendation: retain readable CLI commands and add one
versioned JSON invocation adapter over the same typed services/operation registry.
This is a design proposal, not a delivered command or implementation selection.

Illustrative syntax: `sdptool invoke --request -` (one JSON object on stdin), with
`--request FILE` for reproducible saved requests. `gh sdp` forwards the invocation
and stdin. Prefer stdin over a large inline JSON argument: no shell quoting or
argument-length dependency. Start with one request/response per process; a daemon,
HTTP server, batch protocol or persistent JSON-RPC stream is not needed for this
bounded outcome. Revisit transports only for measured needs.

A proposed envelope separates protocol version, correlation ID, stable action ID,
context (such as project path) and typed parameters. The action catalogue describes
those exact inputs. JSON mode returns one structured success/error envelope on
stdout; diagnostics go to stderr, and exit codes remain meaningful. Require explicit
rules for unsupported versions/actions, duplicate/unknown fields, bounded input,
path resolution, cancellation and output errors. A correlation ID is not automatic
mutation deduplication; preserve existing operation IDs, expected revisions and
trusted authority boundaries. Request data must not grant caller authority.

Existing blueprint assignment apply already reads a JSON request file, but general
Run currently accepts only argv and output writers; there is no generic JSON input
adapter yet. XFMD currently opens child stdin from /dev/null and rejects nonzero
exit status before decoding stdout. External KB-XFMD-030 must therefore cover a
bounded stdin writer and preservation of structured failure results, rather than
assuming JSON input works unchanged. gh-sdp already forwards os.Stdin.

Use this shared invocation/catalogue boundary for eventual MCP mapping without
claiming that a JSON subprocess protocol itself implements MCP. The first complete
slice should exercise catalogue, one read action and one revision-bound mutation
through both human CLI and JSON input, with equivalent validation and results.
CardState remains backlog; next is the bounded design/implementation plan.

## Queue — Session0011 selected (2026-10-09)

Owner selects [Session0011](../../Sessions/session-%230011--SDPTool_actions_and_JSON_API.md)
for focused delivery, pausing predecessor Session0008. CardState backlog → queued:
next is contract refinement and a bounded ImplementationPlan, before execution
activation. Existing CLI/API evidence and the consumer handoff are available; no
new daemon, native XFMD implementation or release is selected by this transition.
Session0008 resumes at BP2-A owner pilot review after this goal or owner redirection.

## Queue prerequisite — Session0011 T002

Owner requires green concurrent-work preflight before implementation. The
[check is RED](../../Sessions/evidence/0011-concurrent-work/README.md): governance
and runnable-program work overlap shared SDPTool files/contracts, and a distinct
runnable-program card in the primary worktree also declares KB-SDP-051. Refer to
this card by subject/path plus branch until identities are reconciled. Remain queued;
no unilateral renumbering or execution activation. Next: coordinate ownership and
baseline/identity reconciliation, then repeat the preflight before S1 execution.


## T007 — bounded contract and implementation route (2026-10-10)

Owner resumes Session0011 and has handed KB052 to PG Codex. Its owning worktree
now has that card active. The [refreshed preflight](../../Sessions/evidence/0011-concurrent-work/refresh-T007.md)
remains RED for shared product changes. This card stays queued; registration is not
execution activation and its colliding identity is not silently renamed.

[Contract proposal](../../04--Design/SDPTool/Actions/Contract.md) defines catalogue,
strict JSON invocation, service ownership, local principal attribution, errors,
paths and retry/cancellation boundaries. [PLAN-SDP-0024](../../05--Implementation/SDPTool/Actions/Plan.md)
is planned: ACT0 integration, ACT1 schema/SDL contract, ACT2 read workflow, ACT3
blueprint operations, ACT4 actual consumer/wrapper tests and ACT5 handoff. Existing
CLI operations remain compatible. No new command is implemented by these documents.
