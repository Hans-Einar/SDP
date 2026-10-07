# ProjectGovernance app-server adapter — PGI2

The package owns a local stdio app-server process, selected cwd and bounded JSON
transport. New connections require **codex-cli 0.161.0** (measured in PGI3-C2). Typed request/result projections
were derived from that executable's generated default schemas; unknown optional
host fields are not persisted. Malformed JSON, duplicate keys, invalid UTF-8 and
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
Host requests are refused by this PGI2 adapter, not automatically approved.
ReceivedAt is a local receive timestamp; no host occurrence time or active effort
is inferred. Epoch changes retain a gap even when history establishes current state.

Controller code binds an observed idle thread to an existing run through
governance.RecordHost, then uses Binding.Request to prepare durable operation IDs.
Persist each request intent before submission and retry that exact request after
an uncertain write. The core preserves immutable receipts, attempt mappings,
completed-item deduplication, conflicting replay refusals and paginated reads.
Read-only observer closure does not mutate execution state. Lost transport must be
recorded as disconnected/unknown by the owning controller. The PGI3 controller candidate owns that loop; live workflow verification is open.

PrepareRecovery reads the already mapped thread and projects its recorded turn;
missing/uncertain history refuses. ReconcileRequest records observed state and the
new epoch. Replay allowlisted completed items with fresh receipt operations,
refreshing Binding.Result after each accepted operation. In-progress history is
not treated as completed item content. Recovery can resolve a lost launch/steer acknowledgment only when exactly one
recorded userMessage.clientId and its exact text match the durable outbound input
in the already bound thread. Persisted correlated observations precede resolution.
Missing, duplicate, changed or unmapped input stays uncertain. An in-progress
turn cannot supply completed-message evidence. No recovery path resends a model request,
changes a Session identity or treats turn completion as SDP acceptance.

## Evidence boundaries

Real negotiation proves initialization, account-mode inspection and MCP tool
discovery. Recorded fixtures prove capture, loss, replay, bounded input and recovery
guards. They do not prove reconnect to an executing model turn, native child
provenance, full controller-loop recovery or the complete owner workflow. Those
remain PGI3–PGI5. Automatic project-journal publication remains PGI4. These commands
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
