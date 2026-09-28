# RGS2 — Codex app-server and SDP development client

| Field | Value |
| --- | --- |
| Milestone | RGS2-C-M1 |
| Plan | [PLAN-SDP-0007](StudyPlan.md) |
| Card | [KB-SDP-038](../../KanBan/active/%23038--Proposal--Codex-app-server-development-client.md) |
| Research date | 2026-09-28 UTC |
| Status | Research delivery; client implementation and deployment remain unselected |
| Candidate | Installed codex-cli 0.158.0, Linux x64; generated protocol schemas |

A small app-server client is a viable **design candidate** for the owner's
supervisor → fresh task Master → Worker/Reviewer workflow. Prefer a terminal
client owning its own app-server process and explicit thread mapping. A Kitty
observer remains useful for durable SDP state, but the study does not establish
passive attachment to an arbitrary existing Codex TUI. No GUI, Codex fork or
runtime client was implemented.

## Evidence and compatibility boundary

This report extends [E11](Evidence.md#e11--codex-panel-follow-up), using official
OpenAI documentation and measured schema generation. The
[evidence manifest](Evidence/RGS2-AppServer-evidence.json) records commands,
exit codes, source URLs/hashes, repository HEAD, launcher/native binary hashes,
and generated-file hashes. The
[schema extract](Evidence/RGS2-AppServer-schema-extract.json) retains selected
structures and complete method-name inventories, including experimental-only
field differences.

Four isolated CLI commands succeeded: version, app-server help, normal schema
generation and generation with `--experimental`. They produced 314 and 440 JSON
files respectively. Generation used a temporary working directory and replacement
HOME/CODEX_HOME/XDG_CONFIG_HOME with no inherited authentication environment.
The CLI warned that the disposable CODEX_HOME did not exist; generation still
succeeded. Preliminary version/help commands also ran in the ordinary environment;
they were not account/configuration reads. Native binary SHA-256:
`167c0148a849d2444f1b5a7fb5f8bb2de1de5ae13a2a504b833fc765980f5cd9`.

**Coverage:** one installed version/platform, two generated surfaces, zero runtime
protocol probes, zero model turns, zero account inspections and zero active-session
attachments. No accounts, configuration or product files were changed. Schema
presence proves exported structure, not working execution, entitlement or stability.
No source commit for the packaged binary was established. Official pages were
retrieved on the research date; their mutable URLs are not release contracts.

The [official app-server page](https://learn.chatgpt.com/docs/app-server) identifies
app-server as a rich-client integration and recommends version-specific generated
schemas. Comparison found material drift: its collaboration example says
`collabToolCall`, while 0.158.0 exports `collabAgentToolCall`; its detached-review
example omits the schema's deprecation. It also warns that paginated-history
creation/recovery is not supported yet. Pin the pilot to legacy history and the
installed schema; do not infer readiness from an exported method alone.

## Fresh tasks, contexts and identities

The following are **0.158.0 schema observations**, followed by design implications.

| Need | Exported contract | SDP implication |
| --- | --- | --- |
| Fresh bounded task | `thread/start` accepts cwd, instructions, model/provider and sandbox/approval options | Start a new Master with a deliberate assignment package; do not copy the supervisor conversation by default |
| Continue same execution | `thread/resume(threadId)`; description says a running thread is rejoined | Resume only the mapped execution after recovering durable work and authority |
| Branch prior history | `thread/fork` copies a source; optional `lastTurnId` cuts history inclusively and cannot reference an in-progress turn | A new thread ID does not establish a fresh or independent context |
| Progress and steering | `turn/start`, `turn/steer`, `turn/interrupt` with thread/turn identity | Scope correction, cancellation and replacement require explicit assignment handling |
| Review | `review/start`; detached delivery deprecated in its schema | Create a fresh thread, then inline review if the built-in reviewer fits; SDP review still needs its own role/acceptance contract |

Thread responses carry `id`, `sessionId`, optional `parentThreadId`, `forkedFromId`,
`agentRole`, `agentNickname`, source, runtime status, cwd and captured Git metadata.
`sessionId` groups a session tree; only native subagents receive `parentThreadId`.
These are not PLAN/card/routine-run identifiers. `originator` records original
creation, not current ownership. Captured Git metadata is not a continually
verified candidate. Preserve an SDP mapping of project/worktree, assignment ID
and revision, role, execution attempt, thread/session IDs and returned evidence.
A replacement thread changes the attempt, not the durable assignment.

The owner wants the Master to use Worker/Reviewer children. The Master retains
responsibility for choosing, decomposing and delegating work within its authorized
assignment. The client presents, launches and maps requested execution; it does
not become an autonomous project supervisor or task coordinator. Two technical
mappings need evaluation without changing that responsibility:

1. **Native delegation from the task Master — desired mapping:** observe
   collaboration items and child identities while the Master uses Codex's native
   delegation. The schema reports sender, receiver thread IDs, requested model/
   effort and last-known agent states. It does not prove which history a child
   received. Freshness, nesting and independent review require a separate versioned
   behavioral check of the actual spawning surface.
2. **Client-created independent threads — fallback to evaluate:** the client can
   launch Worker/Reviewer execution requested by the Master with bounded inputs
   and relay handoffs/results. SDP owns the responsibility mapping; do not label
   these native child threads when the server reports no parent. This provides
   explicit fresh-context construction but needs execution scheduling, escalation
   transport and return handling in the client. A fresh client-created reviewer
   is an alternative if native context isolation proves insufficient, not an
   owner-adopted replacement orchestration policy.

A fresh conversation can still load project instructions and shared files. Review
independence therefore means a separately selected reviewer context, governing
intent and exact candidate, not merely another nickname or model. Never treat
`approvalsReviewer: auto_review` as SDP independent product review: that schema
field concerns permission/risk approval.

The [configuration reference](https://learn.chatgpt.com/docs/config-file/config-reference)
defines `agents.max_concurrent_threads_per_session` excluding the primary thread;
`agents.max_threads` is its legacy alias and an unset value chooses a default.
The [subagent guide](https://learn.chatgpt.com/docs/agent-configuration/subagents)
describes delegated work and returned summaries. Neither establishes this pilot's
effective capacity or a portable nesting limit. This research host separately
allows four concurrent agents including its root; that session instruction is not
an app-server guarantee. Serialize work when capacity is exhausted, reserve review
capacity and record refused/deferred spawns. Do not increase independent threads
to evade a host/account limit. Depth and context-inheritance behavior remain unmeasured.

## Authentication and usage

The normal generated `LoginAccountParams` includes `apiKey`, managed `chatgpt`
and `chatgptDeviceCode` alternatives. It also exposes `chatgptAuthTokens` with an
explicit internal-use-only prohibition. Export in the normal bundle is therefore
not sufficient permission to use every variant. Select managed ChatGPT login for
a later pilot; do not build token extraction or external-token injection.

`GetAccountResponse` distinguishes ChatGPT from API-key accounts and provides
`planType`; the enum contains `pro`, `prolite`, `promax` and `unknown`, among others.
The [pricing page](https://learn.chatgpt.com/docs/pricing) includes Codex with
ChatGPT Pro and distinguishes API-key usage billed at API pricing. This establishes
a documented route, not the owner's actual entitlement or available quota.
No sign-in or account-state assertion is made here.

Proposed client behavior: show the observed authentication mode and returned plan,
require the selected ChatGPT mode before a model turn, and stop on missing/mismatched
auth rather than silently choosing an API key/provider. Do not infer a commercial
plan from a loosely matched enum name. The schema's rate-limit response supports
multiple limit buckets, percentages, reset timestamps, credit information and
nullable `ordinaryUsageAllowed`; its description forbids inferring recovery from
percentages/reset time when that permission is unavailable. Display unknowns,
source timestamps and token usage separately from quota or monetary spend.

Installed help advertises STDIO, Unix sockets and WebSocket transports, with
capability-token or signed-bearer-token WebSocket authentication. Those credentials
protect access to the local server; they do not buy model usage or replace ChatGPT
login. The first design should use a child-process STDIO connection, avoiding a
network listener. Keep credentials and account identifiers out of observer logs;
store only the minimum mode/plan/usage facts needed for the display.

## What a panel can truthfully display

The following inventory comes from the generated protocol, not a working monitor.

| Display | Evidence available | Limit |
| --- | --- | --- |
| Agent work tree | Thread identities/source plus `collabAgentToolCall` sender/receivers and states | Native execution tree and SDP responsibility tree must remain distinguishable |
| Turn progress | Turn/item lifecycle, message deltas, errors, plan updates | Model completion does not establish milestone delivery or acceptance |
| MCP activity | `mcpToolCall` server/tool/arguments/status and optional result/error/duration | Redact content; tool success alone does not prove a durable SDP transition |
| Model plan | `turn/plan/updated` includes threadId, turnId and step/status list | Agent plan is not the authoritative SDP typed plan |
| Skills | `skills/list`, `skills/changed`, explicit `UserInput` skill name/path; response instruction-source paths | Available, selected, source-observed and agent-declared are separate facts |
| Approvals | Server requests for command/file/permission approval, input/elicitation; `serverRequest/resolved` | Permission response is not owner scope, integration or release approval |

No generic skill-activation/deactivation notification appears in the inspected
notification inventory. Instruction-source paths do not establish how every skill
was invoked or complied with. Show “available”, “explicitly supplied”, or “agent
reports loaded” with provenance; avoid an unsupported “currently active skills”
light. The SDP pane should join durable routine/assignment/evidence state from
SDPTool with these observations, retaining their different authorities.

## Connections, recovery and ownership

[Official app-server documentation](https://learn.chatgpt.com/docs/app-server)
describes per-connection initialization/subscriptions, stored reads without resume,
and remote TUI connection to an explicitly started server. Those features do not
establish passive observation of an arbitrary existing TUI. This study did not
attach to one.

The inspected initialization response lacks a negotiated protocol-version/capability
matrix: pin CLI identity and keep a tested method allowlist. Its capability request
has experimental opt-in and notification suppression. Resume parameters can alter
instructions/configuration; a second “observer” must not casually resume with
stale overrides. Help also advertises daemon/proxy entrypoints, but their existence
does not prove safe observer routing or multi-client approval ownership.

No global event-sequence/replay cursor was found in the inspected notification
envelope. History pagination cursors are not evidence of notification replay.
The study leaves six runtime questions open: subscription fan-out, approval-request
routing between clients, competing input arbitration, reconnect event gaps,
process-crash recovery, and survival of pending approvals. A separate client
must not claim exactly-once recovery from these schemas. A second read-only
app-server observer connection is unverified and requires a future probe; it is
not part of the recommended first slice.

Recommended initial rule: one controller owns each server/execution; the observer
reads a redacted projection produced by that controller. On disconnect mark live
state stale. Recover stored thread/items and durable SDP revision, then reconcile
uncertain submissions before allowing another mutation. Never blindly resend a
turn or approval after a lost response. A local journal should map request IDs,
execution IDs and candidate/revision to outcomes without claiming that JSON-RPC
request IDs provide mutation idempotency. The follow-on design must specify
pending-approval recovery and detect concurrent ownership before expanding to
multiple controlling clients.

## UI alternatives and bounded next slice

| Alternative | Value | Cost and boundary |
| --- | --- | --- |
| Separate Kitty pane | Fast read-only view of SDP work, evidence and staleness; independent of Codex rendering | No verified live TUI attach; terminal arrangement is host-specific, while plain text output can travel |
| Small Codex TUI patch | Could render SDP status within existing interaction | Source integration, upstream acceptance, release tracking and regression maintenance; no patch/source build verified here |
| App-server terminal client | Owns dialogue, fresh-task handoff and telemetry consistently | Must implement input, approvals, events and recovery; strongest fit without requiring a GUI |

These cost comparisons are design judgments. Recommend a **single-controller
terminal pilot** over STDIO with a read-only SDP pane. Start with one selected
maintenance assignment: supervisor dialogue creates a bounded Master thread;
the Master decomposes and delegates to Worker/Reviewer children, selecting a
reviewer context with intent and exact candidate, then returns evidence and
unresolved findings to supervision. The client launches/maps requested execution
and displays results. Test native child context/nesting behavior first; retain
client-created fresh execution, particularly review, as a technical fallback for
the next design to assess explicitly. Do not silently adopt a different coordinator.
Exercise the shared pilot paths: inspect/triage, bounded Maintenance,
resume/reconcile and procedure-gap handling. Test request capture and
host-to-assignment binding; neither this client nor MCP establishes universal
process enforcement. Governance and transition rules remain in SDPTool, not
in screen widgets or agent prose.

The successor DesignPlan should require five acceptance groups, all **unexecuted**:

1. Fresh Master/Reviewer inputs contain the assignment and required authority,
   without silently copying the supervisor/worker transcript; IDs survive replacement.
2. ChatGPT mode is visible; missing/expired/wrong-mode authentication blocks model
   submission without API fallback, and no secret enters the projection.
3. Tool/plan/skill observations preserve provenance and uncertainty; a successful
   turn cannot mark SDP work delivered without its evidence/gate contract.
4. Interruption, dropped responses and restart reconcile state without duplicate
   work, stale approvals or invented completion; unsupported methods fail clearly.
5. Read-only observation cannot send turns, approve tools or mutate SDP; protocol
   version mismatch and exhausted capacity are visible with a recoverable next step.

Owner disposition is needed to select this successor and its authentication/deployment
scope. Routine engineering choices include rendering library and JSON framing.
Multi-controller sharing, remote transport, a permanent TUI patch and native-child
freshness remain explicit later decisions. The study does not activate them.
