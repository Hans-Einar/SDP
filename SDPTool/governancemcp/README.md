# Local ProjectGovernance MCP adapter — PGI2

Build from SDPTool: `go build -o /tmp/sdp-mcp ./cmd/sdp-mcp`.
The process uses local stdio only; stdout contains protocol frames. Launch with
explicit `--project ROOT --state PRIVATE_DIR`. These are trusted local launcher
arguments, not model-supplied tool inputs. Inspection does not create state.

The official MCP Go SDK is pinned at **v1.8.0**, restricted to protocol
**2025-11-25**, the intersection exercised with Codex **0.160.1**. This is a
deliberate compatibility boundary, not a claim that the newest protocol was tested.
The SDK owns negotiation/framing; the shared governance core owns decisions.

## Tools and scope

| Tool | Input | Behavior |
| --- | --- | --- |
| context_read | Core Selection | Read coverage and source identities; no activation |
| route_evaluate | Core RouteRequest | Evaluate caller-proposed categories; no language classifier or new run |
| run_read | Empty object | Read the launcher-selected run; cannot select another run |
| events_read | cursor, limit (1–256) | Read a page for the launch-bound attempt; cursor pins run revision |
| run_transition | Core TransitionRequest | Only report, fail and scope-change for the assigned run/attempt |
| attention_read | Optional questionId | Read questions, attributed replies and applicability for the assigned run |
| attention_request | Core AttentionRequest | Ask a durable question; blocking questions prevent advancement |

Without `--run`, only the first three tools are registered; run_read refuses.
With `--run UUID`, context/route selection must match that run and run_read works.
The four remaining tools, including read-only attention/events, additionally require the trusted controller's
`--attempt UUID --epoch UUID` launch binding. Each mutation verifies the latest
persisted attempt and epoch. A disconnected or superseded binding refuses, even
if its payload names the right run. Argument text cannot grant controller rights.
The fully bound worker exposes seven tools. Native child/reviewer attribution
remains unverified; the separate controller owns bounded launches.

No tools expose checkout binding, run creation, resume, completion, owner decisions,
host-observation writes, shell execution or arbitrary project publication. A local
same-user process can bypass this API using the filesystem; the pilot does not
claim a security boundary against arbitrary same-user shell access.

Application schemas are shared from governance/schema.json. Handlers preserve
raw arguments through strict core decoding, rejecting duplicate keys, case aliases,
nulls and unknown fields. Domain refusals are MCP tool errors with structured
codes, not successful empty results. Frames are capped at 1 MiB plus 64 KiB of
transport overhead; application payloads stay capped at 1 MiB. SDK protocol errors
remain separate. Results are capped at 256 KiB before protocol wrapping; run
results summarize attempts and observation pages enforce a 240 KiB content limit.
No HTTP transport, credentials or global host configuration is
created by this adapter.

Tests cover core/MCP result parity, retry receipts, cross-assignment and controller
action refusal, strict inputs, and a compiled stdio server negotiated by the SDK
client. The live controller probe separately verifies Codex tool discovery.

PGI3-C3 keeps this authority unchanged while the controller handles explicit host
approval. Brief checkout-lock contention retries the unchanged tool request and
operation ID (at most 100 attempts, 20 ms delays), respecting MCP cancellation.
Stale revisions, changed arguments and semantic refusals are not automatically
retried. Host approval never grants controller-only actions.

## Current adoption boundary

A static project-bound MCP connection supports inspection; it cannot create or
launch work. A fully bound worker connection is issued by the trusted controller
for a specific run, attempt and epoch. Do not copy those dynamic bindings into a
permanent Project Leader configuration or add controller operations to this worker
surface. Owner/controller decisions and worker reports have different authority.

The [state responsibility map](../../SDP/04--Design/SDPTool/ProjectGovernance/State-Boundaries.md)
explains those distinctions. The current
[MCP adoption study](../../SDP/Maintenance/WorkstationReadiness/MCP-Adoption-Study.md)
records actual workstation capability, configuration and host compatibility limits.
The documentation correction adds no tools or runtime behavior.
