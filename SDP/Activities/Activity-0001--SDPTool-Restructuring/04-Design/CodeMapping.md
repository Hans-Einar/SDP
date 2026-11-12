# Candidate architecture mapping

Product candidate: 527cc9008bdf8f567f8b201b665b2960865518f1. This authored mapping
supplements the pinned Blueprint; it does not modify generated facts or turn the
bundle's unknown CODE_MAPPING result into machine proof. Each source path below
is relative to SDPTool/. [Verification](../06-Verification/Verification.md) owns
execution results, source hashes and remaining acceptance gates.

| TARGET unit / responsibility | Candidate code | Disposition and verification |
| --- | --- | --- |
| CommandLineAdapter / DecodeCommandArguments | cli.go, governance_cli.go, model_cli.go, program_run.go | One root routing table; governance definitions supply dispatch/input kind; argument/wire parity, root and compiled tests |
| JsonInvocationAdapter / DecodeActionRequest | action_cli.go, action_invoke.go, action_schema.go | Strict request decoding retained; caller policy resolved from registered action; schema/authority tests |
| OperationRegistry / DescribePublicActions, CheckActionAttribution | action_registry.go, action_blueprints.go, action_sessions.go | Attribution and availability accompany definitions; 28-action frozen catalogue, policy tests, no action-ID allowlist in Invoke |
| GovernanceApplication / AdaptSessionInvocation, DispatchGovernanceOperation | governanceapp/operations.go; governance_service.go resolves project/private store | Single typed handler and input definition per governance operation; Session adapters map their wire inputs; no argv/output bridge |
| GovernanceService / ManageGovernedWork, ManageSessionExecution | governance/ | Retained domain authorization, transactions and recovery; full suite and compiled CLI/MCP journeys |
| ProjectFacade / DiscoverProject, BuildNavigation, SelectDocument | discovery.go, navigation.go, deferred_navigation.go, selection.go, preview.go, sdui.go, viewer.go | Retained typed facade and producer contracts; root and native XFMD checks |
| ModelGovernance / ManageModelArtifacts | model/ | Retained artifact/revision state owner; model recovery/compatibility suite |
| BlueprintService / GenerateBlueprint, AssessBlueprint, ManageBlueprintAssignments | blueprints/, blueprintstate/ | Retained generation/assessment and bounded assignment-store owner; package suites and compiled/root adapters |
| ProjectHistory / PersistCanonicalHistory | projecthistory/ | Shared canonical history publication, locking and recovery retained; package and race tests plus governance recovery |
| ProgramRunner / RunSelectedProgram | program_execution.go, program_process_unix.go, program_process_other.go | Flag-free typed execution; existing process groups, input/output, cwd and exit behavior; root lifecycle and revision/exit tests |
| InstallService / PreviewInstallation, ApplyInstallation, RecoverInstallation | install/, installation.go, bootstrap/ | Retained planner/apply/recovery and bootstrap selection boundaries; install, nested bootstrap and gh-sdp tests |
| ReleaseLogService / ProjectReleaseNotes | release_log.go | Retained release-note projection; root tests; no publication performed |
| ControllerAdapter / DecodeControllerCommand | cmd/sdp-controller/ | Retained typed calls and bounded CLI controls; command/controller and compiled governance tests |
| HostCoordinator / CoordinateHostWork, LaunchIsolatedSessionController | governancecontroller/, governancehost/ | Retained host protocol, child leases/logs and explicit concurrency; controller/host and recovery tests |
| McpAdapter / ExposeBoundCapabilities | cmd/sdp-mcp/, governancemcp/server.go, session.go, session_host.go | Retained six capability configurations and instance/skill guards; MCP, compiled and authority tests |
| McpResponseAdapter / FrameGeneralMcpResponse, FrameSessionMcpResponse | governancemcp/response.go | One response/retry implementation with schema-specific policy; framing/limit/stale/busy/cancellation tests |

SdpToolHost, ControllerProcess and McpProcess remain actual process boundaries.
SDL and SDUI remain collaborating language/runtime systems. The host's root Go
package remains its composition and public facade; file splitting is not a claim
of another runtime container. Retained presentation and codexcompat packages are
adapter utilities, not additional state owners.

## Removed authority and temporary paths

- Session actions no longer construct governance argv, invoke a CLI handler,
  capture stdout or reparse JSON results. CLI and actions execute governanceapp
  operations directly; the ordinary external CLI remains supported.
- Governance operation names, handler dispatch and request classifications share
  governanceapp definitions. Root command recognition and routing share cliCommands.
- Invoke no longer has a separate action-ID attribution list. Blueprint policy
  and platform availability are registered beside its handlers.
- Ordinary and Session MCP wrappers share response framing/retry policy. Capability
  grants remain intentionally separate from unrestricted local CLI exposure.
- Program flags and process execution have separate ownership. No temporary
  adapter was introduced; there is no deferred removal obligation.

The architecture test's deliberately violating specimens cover forbidden domain
imports, Session CLI/JSON bridges, duplicate caller-ID policy and MCP framing.
Behavioral tests and independent review cover semantics those syntax checks cannot
prove. The manual Activity records remain a separate pilot; no automatic Activity
creation, lifecycle admission or browsing support is implied by this mapping.
