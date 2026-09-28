# SDP development client using Codex app-server

| Field | Value |
| --- | --- |
| id | KB-SDP-038 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-09-28T15:10:13.034937+00:00 |
| source | Owner conversation 2026-09-28: register app-server-based system development and assess SDPTool MCP |
| next_review | With KB-SDP-036/037 when selecting the bounded routine and client DesignPlan |
| tags | codex, app-server, client, MCP, supervision, blueprints, observability |
| Systems | SDPTOOL |

## Owner intent

Explore using Codex app-server to develop systems with Codex through an
SDP-aware client. The owner needs coherent project dialogue, bounded assignments
and visible work state, including agent roles, skills, MCP activity and the
position in the work tree. A right-side panel is one presentation option, not
the whole capability. Preserve the ability to use ChatGPT Pro authentication.

The owner normally talks to Steering/Project Manager. A fresh Master coordinates
each bounded task, with Worker, independent Reviewer and appropriate Verifier
children. Carry forward the role clarification in the
[RGS1 study](../../02--Requirements/RoutineGovernance/Study.md); do not make the
project-facing conversation responsible for remembering every execution detail.

This request authorizes backlog registration, not a Codex fork, client
implementation, host configuration or automatic agent delegation.

## Proposed responsibility split

Recommendation: keep both app-server and an SDPTool MCP adapter, serving different
purposes. These are proposed boundaries, not delivered interfaces.

| Component | Responsibility |
| --- | --- |
| Codex app-server | Agent conversations/threads, model execution, authentication, host approvals and activity events |
| SDPTool core | Project discovery, current process/model context, selected work and validated routine/evidence operations as those capabilities are implemented |
| SDPTool MCP adapter | Let Codex agents query that context and invoke bounded SDPTool operations; reuse the same core as CLI |
| SDP-aware client | Owner dialogue, project/task navigation, role/thread association, observation and decision presentation |

App-server is the interface our client uses to operate Codex. SDPTool MCP is the
interface Codex agents use to access SDP. Neither replaces the other. A client
may also use an appropriate SDPTool snapshot/event API directly; it need not
pretend to be an agent or route every UI read through a model/MCP call.

Do not duplicate business rules in the client or MCP adapter. Keep durable SDP
routine/task identities separate from Codex thread/session identities. Link them
explicitly so replacement sessions and retries preserve the same work history.
MCP availability alone does not enforce procedure compliance or intercept arbitrary
shell/filesystem changes; retain RGS1's measured host-integration boundaries.

## Initial investigation and candidate scope

- Compare a separate Kitty observer, a small Codex TUI extension and a custom
  app-server client. Select presentation after a bounded prototype; do not assume
  a permanent Codex fork is necessary.
- Demonstrate project-facing supervision handing a bounded task to a fresh Master,
  with explicit context/blueprint references, child work and returned evidence.
  Respect actual host nesting/concurrency limits; a role tree need not map to
  unlimited recursively spawned agents.
- Correlate observed tool/agent events with SDP-owned plan, milestone, routine and
  gate state. Codex conversational plan completion is not SDP acceptance.
- Distinguish available skills, explicit invocation/loading and declared task use;
  unknown activation stays unknown. Separate available MCP tools from active calls
  and their results. Do not label instruction compliance as directly observed.
- Establish connection/session ownership and recovery. Do not assume passive
  attachment to an arbitrary running TUI, safe multi-client control or reliable
  event replay without testing the selected app-server version.
- Reuse Codex-managed ChatGPT sign-in for subscription access. Verify the active
  auth mode and applicable usage limits. Do not silently fall back to separately
  billed API-key use. Client/server transport authentication is a separate concern.
- Keep private project context scoped to the selected project and avoid exposing
  credentials, complete transcripts or unnecessary source material in the monitor.

## Existing findings and boundaries

[Official app-server documentation](https://learn.chatgpt.com/docs/app-server)
describes ChatGPT-managed authentication as well as API-key authentication, and
agent/tool event interfaces. ChatGPT Pro can be used; app-server does not require
switching to API-key billing. Applicable subscription limits still apply; see
[pricing](https://learn.chatgpt.com/docs/pricing). This records documentation checked
on 2026-09-28, not a tested local app-server session or account configuration.

[RGS1 Evidence E11](../../02--Requirements/RoutineGovernance/Evidence.md#e11--codex-panel-follow-up)
records existing Codex inspection commands and the limits of the panel inquiry.
The SDL assignment compiler, generic SDP routine engine and client proposed here
must not be described as existing functionality.

## Relationships and next action

- [KB-SDP-036](%23036--Proposal--Request-routines-and-process-control.md) owns request
  routing, routine selection and entry guidance.
- [KB-SDP-037](%23037--Proposal--Observable-routine-state-machine.md) owns durable
  routine execution, MCP-facing process operations and observation contracts.
  This card owns the Codex client/integration experience and consumes those
  contracts; it does not create a second workflow engine or duplicate MCP project.
- [KB-SDP-031](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md)
  and planned [BP2](../../04--Design/SDPTool/Blueprints/Plan.md) own assignment context
  and blueprint design. Do not activate BP2 merely by linking it.

Next action: select a proportionate DesignPlan or bounded investigation together
with KB036/037, resolving presentation and lifecycle ownership. No new software
System or source directory is adopted by this registration.

A selected pilot should prove one small end-to-end assignment: context resolution,
supervision-to-Master handoff, permitted child execution/review, visible evidence
and a correct return disposition. Include scope escalation, stale candidate,
disconnection/reconnect, restart and missing-routine cases. Demonstrate actual
ChatGPT sign-in and telemetry without manufacturing owner approvals. A UI mock or
successful MCP connection alone does not complete the capability.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-28T15:10:13.034937+00:00 | codex; EVT-KB-SDP-000215 | Registered the requested app-server client investigation and complementary SDPTool MCP boundary | Backlog only; plan, prototype and implementation unselected |
