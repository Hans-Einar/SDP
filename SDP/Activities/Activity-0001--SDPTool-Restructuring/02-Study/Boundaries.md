# M0 controller, MCP and external boundaries

Read-only inventory by independent bounded agent, baseline d0bf6dc. No host effects
were executed. External revisions: gh-sdp b177e143; XFMD 8893a32 (unrelated dirty
RPM/process changes excluded).

## Controller: 51 accepted command paths

| Family | Exact operations | Disposition |
| --- | --- | --- |
| Host/control | probe, execute, run, continue, status, recover, attach, provenance | Retain governancehost/governancecontroller typed calls |
| session | list, status, resume, select, next, rework, complete, accept, changes-requested, questions, answer | Retain governance.Service calls |
| overview | overview | Retain |
| inbox | inbox, inbox acknowledge | Retain |
| attention | list, request, reply, withdraw | Retain |
| coordinate | coordinate | Retain |
| leader | inspect, message, run | Retain |
| evidence | capture, inspect, attach, rework | Retain |
| review | prepare, run, inspect, supersede | Retain |
| publication | prepare, inspect, apply, abandon | Retain |
| assignment | context, prepare, inspect, revise, retire, activate, design-review, design-decide, design-check | Retain |

Authority: SDPTool/cmd/sdp-controller/main.go and named family files. execute is
accepted although omitted in the usage text; do not remove it based on help alone.
Embedded coordination controls: refresh, quit, acknowledge, reply, message.
Live-run stdin controls: detach, status, approval, steer, interrupt, scope-change.
Durable Session host requests: launch, continue, interrupt, scope-change.

session_dispatcher.go owns at most two child processes, their logs/exit tracking.
Its self-exec is retained as a true lifetime/lease boundary: each child independently
claims durable requests. It does not parse CLI output back into domain data.
No internal CLI round-trip was found in governancecontroller.

## MCP: 28 unique names across six launcher configurations

| Configuration | Exact exposed tools | Disposition |
| --- | --- | --- |
| Unbound inspection | context_read, route_evaluate, run_read | Retain; run_read refuses without binding |
| Run only | context_read, route_evaluate, run_read | Retain bound selection |
| Run/attempt/epoch | prior three plus attention_read, attention_request, events_read, run_transition | Retain; report/fail/scope-change only |
| Coordinator | coordination_read, report_acknowledge | Retain identity-bound acknowledgment |
| Session controller | session_read, session_definition, session_successor, session_rebind_context, session_design_context, session_rebind, session_resume, session_select, session_assess, session_start, session_progress, session_verify, session_review, session_dispose, session_owner_request, session_advance | Retain exact instance/skill guards and owner-request-only restriction |
| Session with host requests | prior sixteen plus session_host_read, session_host_request, session_host_cancel | Retain; queues requests, never dispatches hosts |

Authority: governancemcp/server.go, session.go, session_host.go; launcher
cmd/sdp-mcp/main.go. All call typed domain methods. Consolidate only duplicated
transport retry/error/response framing into a shared wrapper with explicit policy;
keep schema-specific limit messages, 256 KiB limit, strict decoding, at most 100
busy attempts/20 ms delays and cancellation of waiting. Stale failures never retry.
Capability registration is a deliberate separate exposure contract, not a copy of
the unrestricted local CLI registry. Existing instance/skill checks stay at calls.

## External process/library boundaries

- gh-sdp forwards all argv/streams/exit codes to the selected verified engine.
  Bootstrap resolves signed artifacts and probes version/protocol in a child.
  These are permanent external protocols; retain them. Other bootstrap callers:
  install/cli.go, install/plan.go, install/history.go, blueprintstate/state.go and
  tools/profile/sign.go all use direct library calls.
- XFMD actually invokes actions, discover, tree --model-only, select and the
  sdui-preview fallback, with --json. Authority: src/application/sdp/SdpDiscovery.cpp,
  SdpNavigation.cpp, SdpToolJob.cpp and SdpProtocol.cpp. The catalogue currently
  detects deferred capabilities; XFMD is not a general action executor.
- Preserve XFMD_SDP_TOOL/bundled/gh-sdp selection, old-producer arguments fallback,
  cancellation/process-group cleanup, revision/output checks and dirty-file refusal.
- Direct SDUI preview/check/Fyne invocation remains owned by that collaborating
  system. Renderer/host process protocols are external boundaries, not domain bridges.

## Evidence obligations

Controller tests cover typed requests/terminal parsing. Governance compiled tests
exercise Session CLI/MCP transfer, authority, gates and host recovery. MCP server,
coordination and contention tests exercise launch limits and busy behavior.
Bootstrap is a nested Go module and needs its own tests. gh-sdp has real delegation
and integration tests. Native XFMD gates are SdpNavigationGuiTest,
SdpOverviewGuiTest, SdpDeferredGuiTest, SdlSourceGuiTest and SduiWorkflowGuiTest with
the selected producer. Source inspection alone is not native acceptance.
