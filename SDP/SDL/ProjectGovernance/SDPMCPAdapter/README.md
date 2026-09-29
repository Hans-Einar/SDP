# SDPMCPAdapter — governed agent access

## Boundary and implementation state

This proposed system provides the MCP-facing entry into SDPTool. There is **no
adapter implementation or binary** in the inspected repository. `SdpMcpHost` is a
proposed separate executable boundary (candidate binary `sdp-mcp`); a reusable
adapter library may sit beside it. The shared SDPTool application service remains
the only owner of routine policy and accepted transitions. Local Go library calls
are a valid implementation of `CoreCalls`; this model does not require network
traffic or a second SDP daemon.

The source basis is the [MCP study](../../../02--Requirements/RoutineGovernance/MCP-Study.md)
and [synthesis](../../../02--Requirements/RoutineGovernance/Synthesis.md).
All `AdapterDelivery` functions are planned. Local STDIO is a recommendation;
protocol and SDK selection require a measured compatibility intersection with the
actual Codex host. The pinned study is evidence about documents and source, not
successful negotiation or runtime enforcement.

## Functions and public contract

`ProtocolAdapter` exposes context resources and validates incoming tool envelopes.
`CoreBridge` passes bound operations to SDPTool. `ResultProjection` returns source,
revision and coverage facts without converting tool success into accepted work.
Passive resource reads do not start runs, mutate cards, generate preview files or
launch viewers. Those side effects require explicit operations.

`ToolOperation` and `ToolOutcome` are **proposed conceptual application envelopes**,
not replacements for MCP JSON-RPC or committed SDK schemas. They capture selected
project/checkout/run, revision, idempotency, action and evidence. Strict JSON shape,
field domains, trusted caller binding, authorization and resource limits remain
DesignPlan work. They cannot be established by the design-core field checker.

`CoreOutcomeReturned` and `CoreRefusalPreserved` show correlated outer and inner
requests. The core is responsible for the disposition; the adapter preserves it.
Internal handoffs between adapter units are intentionally abstracted. These sample
paths do not imply atomicity, exactly-once delivery, or permission granted by a
caller-supplied role. Operational cursor/replay semantics belong to SDPTool, not
MCP request/session IDs. Unsupported protocol, unavailable core and stale context
must remain explicit failures, not empty successful resources.

`ExternalMcpCaller` represents the agent host; `ExternalSdpToolCore` represents the
[SDPTool system](../SDPTool/README.md). They are unowned unit stubs solely for this
standalone model's boundary sequence. Names do not link across entry files. The
future shared contract must have one authority when source-set linking exists;
these local envelopes must not be treated as independently versioned wire APIs.

## Model and validation boundary

[System.design](System.design) is the independently parsed entry for **design-core
0.5**. Its declarations belong to this file only; directory names are catalog
metadata, not language namespaces or imports. Activities distinguish implemented
source responsibilities from planned delivery. A valid model or rendered scenario
is not runtime evidence, owner acceptance or authorization to implement a proposal.

From the repository root, with a built SDL CLI:

```sh
sdl check SDP/SDL/ProjectGovernance/SDPMCPAdapter/System.design
sdl ast SDP/SDL/ProjectGovernance/SDPMCPAdapter/System.design
sdl viewpoints SDP/SDL/ProjectGovernance/SDPMCPAdapter/System.design --format static --output /tmp/sdp-sdpmcpadapter-views
```

Build the CLI from `SDL/go` with `go build -o /desired/bin/sdl ./cmd/sdl`.
The authoring verification used the repository CLI with Go 1.27.1; the enclosing
[Ecosystems plan](../../../03--Architecture/Ecosystems/Plan.md) records consolidated
candidate hashes and export evidence. Generated output belongs in temporary or
explicitly marked derived directories, never in this authored source folder.
The [ecosystem index](../README.md) identifies cross-system dependencies.
Navigation registration must select this entry explicitly; folder placement alone
does not make it available to a viewer.
