# CodexClient — project supervision and bounded task execution

## Boundary and implementation state

This is a **proposed** terminal client system, separate from Codex itself.
`CodexClientHost` owns one app-server controller and its read-only observation
projection. No client source/binary exists in the studied repository. Candidate
binary name `sdp-codex` and client libraries are provisional; choosing a TUI toolkit
or implementing a Codex fork is outside this model delivery.

The source basis is the [app-server study](../../../02--Requirements/RoutineGovernance/AppServer-Study.md)
and [synthesis](../../../02--Requirements/RoutineGovernance/Synthesis.md). Codex
0.158.0 schema generation was inspected; runtime behavior, authentication,
native-child context and reconnect behavior were not exercised. Pin and test
actual compatibility before using the candidate model as an executable contract.

## Responsibility and execution

The owner works with a project-facing Steering/PM agent. A fresh bounded Master
owns decomposition, permitted Worker/Reviewer delegation and integrated return
results. The client presents, launches and maps execution requested through this
role model; it does not become an autonomous project manager or replace Master
coordination. Native child delegation is the desired mapping. Client-created
independent threads remain an alternative to evaluate, not a silently adopted
replacement or proof of native parentage.

`ExecutionController` checks managed authentication and starts the selected fresh
Master. `AssignmentMapping` binds an execution attempt to durable SDP work.
`ObserverProjection` distinguishes observed activity from verified delivery and
marks disconnected observations stale. `Reconciling` scopes uncertainty recovery;
no blind resubmission of turns/approvals after an ambiguous response. One controller
owns input/approvals; an observer reads its redacted projection and cannot mutate
work or send model turns. These are target behavior, not implemented enforcement.

`ExecutionRequest` and `ExecutionOutcome` are **normalized client-domain records**.
They are not the raw `thread/start` request/response: the exact versioned app-server
adapter must translate supported fields, retain SDP assignment/attempt mapping
locally, and omit unsupported fields. `Role` is a recorded responsibility, never
an authorization credential. The context read precedes the fresh thread request
in the example; independent context freshness and review must still be tested.
`RunObservation` carries SDP revision/state/coverage separately from Codex status.

Managed ChatGPT authentication is the study recommendation; missing or mismatched
mode should block model submission without silent API-key fallback. This model
neither reads an account nor proves entitlement. Credentials must not enter
observer logs. Skill availability, explicit supply, observed invocation and
claimed compliance remain distinct. A completed turn is not an SDP acceptance.

`ExternalCodexAppServer` and `ExternalSdpToolCore` are unowned endpoint stubs. Codex
is maintained outside this repo; [SDPTool](../SDPTool/README.md) owns workflow state.
MCP is the agent-to-SDP path, while app-server is client-to-Codex; the model does
not collapse them into one protocol or assume passive attachment to any TUI.

## Model and validation boundary

[System.design](System.design) is the independently parsed entry for **design-core
0.5**. Its declarations belong to this file only; directory names are catalog
metadata, not language namespaces or imports. Activities distinguish implemented
source responsibilities from planned delivery. A valid model or rendered scenario
is not runtime evidence, owner acceptance or authorization to implement a proposal.

From the repository root, with a built SDL CLI:

```sh
sdl check SDP/SDL/ProjectGovernance/CodexClient/System.design
sdl ast SDP/SDL/ProjectGovernance/CodexClient/System.design
sdl viewpoints SDP/SDL/ProjectGovernance/CodexClient/System.design --format static --output /tmp/sdp-codexclient-views
```

Build the CLI from `SDL/go` with `go build -o /desired/bin/sdl ./cmd/sdl`.
The authoring verification used the repository CLI with Go 1.27.1; the enclosing
[Ecosystems plan](../../../03--Architecture/Ecosystems/Plan.md) records consolidated
candidate hashes and export evidence. Generated output belongs in temporary or
explicitly marked derived directories, never in this authored source folder.
The [ecosystem index](../README.md) identifies cross-system dependencies.
Navigation registration must select this entry explicitly; folder placement alone
does not make it available to a viewer.
