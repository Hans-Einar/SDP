# ProjectGovernance app-server adapter — PGI2

The default transport owns a local stdio app-server process, selected cwd and bounded JSON
transport. PGI3-C4 optionally connects to a separately running private socket host. New connections require **stable codex-cli 0.164.0 or later** within the bounded
PG wire profile. The actual executable and initialize identity must agree; socket
connections retain the observed server version. Later stable versions are eligible,
not automatically workflow-tested. Missing required operations/shapes refuse
explicitly; unknown optional fields are not persisted. Malformed JSON, duplicate keys, invalid UTF-8 and
frames over 1 MiB refuse. The PGI3 [owner controller candidate](../governancecontroller/README.md) adds
explicit work-thread/start/steer/interrupt requests. Sign-in and account-mode
fallback remain unsupported. A new connection has a new epoch.

Build from SDPTool: `go build -o /tmp/sdp-controller ./cmd/sdp-controller`.

```sh
/tmp/sdp-controller probe --project /absolute/project
/tmp/sdp-controller probe --project /absolute/project --mcp-binary /tmp/sdp-mcp --state /private/governance-state
```

The first probe checks version, initializes once, sends initialized, and reads only
account mode with token refresh disabled. The second creates an ephemeral idle
read-only thread with a thread-local MCP configuration, then inspects only that
server's tools. Neither sends turn/start, edits persistent host configuration,
changes account mode or authorizes an agent. An unavailable/non-ChatGPT account is
reported explicitly; the probe cannot switch to API-key billing. Raw account
responses, email, credentials and host stderr are not retained.

## Capture and reconciliation

Visible capture retains completed user and assistant messages with thread/turn/item
identities and commentary/final/unknown phase. It retains only status metadata for
selected completed tool/agent items. Reasoning, authentication notifications, raw
tool output and streaming deltas are discarded. Supported local attachment paths
are retained with unknown availability; remote image/audio payloads are marked
unavailable rather than copying embedded or credential-bearing URLs. This is a
bounded capture profile, not full attachment transport.

The queue holds at most 256 observations. Overflow or malformed captured events
close the connection with an observation gap; they are never silently truncated.
Host requests default to refusal. Explicit PGI3 ask mode exposes only the measured
SDP MCP approval form and accepts a correlated owner answer; it never grants
approval automatically.
ReceivedAt is a local receive timestamp; no host occurrence time or active effort
is inferred. Epoch changes retain a gap even when history establishes current state.

Controller code binds an observed idle thread to an existing run through
governance.RecordHost, then uses Binding.Request to prepare durable operation IDs.
Persist each request intent before submission and retry that exact request after
an uncertain write. The core preserves immutable receipts, attempt mappings,
completed-item deduplication, conflicting replay refusals and paginated reads.
Read-only observer closure does not mutate execution state. Detected transport loss is recorded as disconnected/unknown by the owning controller.
An abrupt controller kill can leave an as-of running snapshot; reattachment must
measure the actual mapped turn before continuing observation.

PrepareRecovery reads the already mapped thread and projects its recorded turn;
missing/uncertain history refuses. ReconcileRequest records observed state and the
new epoch. Replay allowlisted completed items with fresh receipt operations,
refreshing Binding.Result after each accepted operation. In-progress history is
not treated as completed item content. Recovery can resolve a lost launch/steer acknowledgment only when exactly one
recorded userMessage.clientId and its exact text match the durable outbound input
in the already bound thread. Persisted correlated observations precede resolution.
Missing, duplicate, changed or unmapped input stays uncertain. An in-progress
turn cannot supply completed-message evidence. No recovery path resends a model request,
changes a Session identity or treats turn completion as SDP acceptance. Recovery
and reattachment require the connected host version to match the recorded attempt;
a software upgrade cannot silently relabel existing work. Legacy 0.160.1/0.161.0
records remain readable through the core. See [PGC1](../../SDP/05--Implementation/SDPTool/ProjectGovernance/Codex-Compatibility.md).

## Evidence boundaries

Real negotiation proves initialization, account-mode inspection and MCP tool
discovery. Recorded fixtures prove capture, loss, replay, bounded input and recovery
guards. PGI3-C4 separately verifies actual executing-turn reattachment after graceful
detach and abrupt controller loss. Native child/skill provenance, independent
review and the complete owner-to-evidence workflow remain PGI3–PGI5. Automatic project-journal publication remains PGI4. These commands
are repository build targets, not a published installer/release inventory.

Protocol references: [OpenAI app-server](https://learn.chatgpt.com/docs/app-server),
[MCP SDK v1.8.0](https://github.com/modelcontextprotocol/go-sdk/tree/v1.8.0).
Local generated schema identities and actual probe output are retained in PGI2
implementation evidence; online documentation does not override that measured pin.


PGI3-C2 adds real managed-account execution, steering, observer closure,
interruption, replacement-attempt continuation and a deliberately dropped
turn/start response. The last case resolves from recorded history without a second
submission. Host approval refusals are preserved; the live run_transition call
required approval under the selected never policy. This is still not independent
review or evidence acceptance. Historical 0.160.1 stores remain readable, but a
new 0.161.0 connection cannot silently relabel an old bound attempt's host version.


PGI3-C3 supports the measured mcpServer/elicitation/request form for the assignment's
run_transition tool, retaining a bounded message and exact transition arguments.
Transport-owned request IDs distinguish JSON integers and strings. Pending IDs are
single-use, bound to thread/turn/epoch, and invalidated on resolution or turn end.
Up to 16 pending requests share the existing bounded observation queue. Replies
omit persistence metadata. Only serverRequest/resolved establishes observed closure.
Actual schema export includes granular approval policy, but the host gates it on
experimentalApi; explicit ask uses the stable on-request/user path instead.

PGI3-C4 adds an optional Linux private Unix WebSocket connection using
coder/websocket v1.8.15. It keeps JSON frame limits, strict wire validation and
visible-event filtering. The initialize user agent must report the measured host
version. Socket identity combines requested path, socket inode/creation metadata,
peer user/process/start time and boot identity; it is a same-host continuity check,
not a credential against arbitrary same-user processes. Close only closes the
connection. The separately launched server remains explicitly owned by its caller.

PrepareAttachment preflights the covered run and unresolved inputs/approvals before
reading and resubscribing to the one mapped active thread. thread/resume receives
only threadId. No turn/start, permission override or lost-request replay is sent.
Missing/full-history coverage and a completion race refuse explicitly. A new
connection epoch retains a visible gap even after terminal-history replay.

PGI3-C5 preserves the optional thread/start instructionSources and native
session/parent/role/source values, separately from the controller-assigned role.
InspectSkills uses the pinned stable skills/list shape (`cwds`, `forceReload`).
It snapshots at most 64 project-local skills, hashes readable files up to 64 KiB
through a confined project-root open, and counts external entries without reading
their contents. Invalid, ambiguous or excessive inventory is reported unavailable
by the controller. A discovery result never establishes skill activation.


Explicit skill dispatch uses the measured 0.161.0 text-plus-skill input array and
stable clientUserMessageId. Completed notifications and full thread history preserve
skill name/path/order. The adapter fingerprints only the lossless supported subset,
normalizing absent/empty text_elements. Unknown fields/content and nonempty elements
cannot prove correlation. Exact acceptance establishes input identity, not consumed
skill bytes or procedural compliance. No native child dispatch is introduced.
