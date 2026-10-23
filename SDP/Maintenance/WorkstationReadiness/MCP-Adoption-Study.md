# ProjectGovernance through MCP — adoption study

Session0013 T002 / MAINT-SDP-0018 WR2-M4, 23 October 2026.
The owner defers MVP1 and confirms an existing Ponsse pilot model. This study
addresses use of the broader SDP workflow; it does not propose another pilot model.
[Execution evidence](Execution-Evidence.md) pins the candidate, environment and checks.

## Answer: what can be used now?

**Project-bound inspection can be used in a controlled development setup now.**
Actual Codex 0.164.0 discovers the compiled MCP server and successfully calls
context_read and route_evaluate against this SDP project without submitting a
model turn or creating operational state. Context has no coverage conditions;
a status route yields no-work. An unbound run_read refuses as designed.

**Complete ProjectGovernance is not ready for ordinary autonomous use.** The worker
MCP adapter is not a Project Leader control API. It cannot start/resume workers,
record owner acceptance, consume reports as a leader or publish canonical outcomes.
The current controller also refuses this workstation's Codex version before launch.
Those are specific remaining boundaries, not a reason to discard implemented PG.

## Measured implementation progress

| Capability | Existing delivery and current observation | Adoption consequence |
| --- | --- | --- |
| Coverage/routing and durable run core | PGI1 delivered; current core/adapter tests and direct MCP calls exercise existing records | Useful inspection and state primitives already exist. Natural-language classification is not supplied by route_evaluate. |
| Worker MCP transport | PGI2 adapter, SDK v1.8.0 / protocol 2025-11-25; real compiled stdio calls verified here | Local stdio works; no HTTP/remote deployment is implemented by this adapter. |
| Bounded execution/recovery | PGI3 has prior native evidence through Codex 0.161.0, steering/recovery/skill inputs; native child/review mapping remains open | Current 0.164.0 controller use is blocked by exact version pin. Earlier trials are historical evidence, not current compatibility certification. |
| Work inspection and packages | PGL1 overview; PGL2 preparation, revision/retirement, live provider review, attributed decision and activation delivered | Services/CLI exist. They are not automatically accessible as leader MCP tools. |
| Reporting and attention | PGL3-C2 reports/read-only inbox, C3 questions/replies, C4 guarded continuation delivered | Worker can report/ask under a controller-issued capability. Inbox consumption is explicitly not-tracked; responsive leader remains unfinished. |
| Whole owner-to-accepted-outcome workflow | PGI4/5, later PGL, SAD dispatch/scope closure, native parallel/capacity and independent-review gates remain | Do not claim automatic checked delivery, canonical Session publication or unattended multi-hour operation. |

This is a capability assessment rather than a percentage-complete estimate:
substantial primitives are implemented, but the missing integration gates determine
whether a particular workflow is usable.

## Inspection coverage is intentionally narrow

The successful condition-free context read selected WR2 while its MaintenancePlan
and card were active. Unbound context_read uses the legacy Maintenance coverage
path; it is not a whole-board reader or a generic ImplementationPlan inspector.
An additional [actual check](evidence/WR2-implementation-context.json) against PG's
active ImplementationPlan returns `required plan identity prefix: MAINT-` and
`active MaintenancePlan required`. That is a coverage limitation, not a reason to
rename the PG plan or create a duplicate Maintenance wrapper. Bound activated
implementation runs use the separate run-context path.

After this bounded WR2 plan closes, its old selection must not be presented as an
active assignment. Generic status-only routing still reads without activation;
select current eligible work when evaluating Maintenance coverage. Broader project
status/overview and prepared assignment inspection currently belong to controller
services/CLI and are part of the proposed leader MCP integration.

## Three distinct MCP uses

| Mode | Launch binding / tools | Current status |
| --- | --- | --- |
| Inspection | Explicit project + private state path; context_read, route_evaluate and unbound run_read registered | Verified through direct stdio and actual Codex read calls. Hide unbound run_read in a practical static config because it cannot succeed without a run. |
| Assigned worker | Trusted controller selects run + attempt + epoch; seven tools: context_read, route_evaluate, run_read, events_read, run_transition, attention_read, attention_request | Real stdio fixture verifies the surface, reports, retry and refusals. The host binding in this trial is scripted; it is not a native agent task. |
| Project Leader control | Would invoke trusted preparation/activation/dispatch, consume reports, coordinate phases and retain owner decisions | Current worker MCP does not implement this surface. Reuse existing services under explicitly designed controller authority; do not widen worker privileges or use model-supplied role strings. |

Static run IDs and copied attempt/epoch values are not a general leader integration.
Controller-issued worker capabilities must remain scoped and refreshed by the real
execution lifecycle. A run reaching verifying proves reporting; acceptance and
canonical card/Session publication are separate obligations.

## Workstation prerequisites and actual probes

A private Go 1.27.1 toolchain was fetched from the official Go distribution with
its published SHA-256 verified. Current sdptool, sdp-controller and sdp-mcp were
built from the integrated candidate with development identity 4174c68. The signed
gh-sdp engine remains 2.2.0; no stable installation was replaced.

The controller probe returns `incompatible_profile: require codex-cli 0.161.0`.
Installed CLI is 0.164.0. We did not remove that guard or relabel 0.164.0 as 0.161.0.
A separate read-only app-server probe, outside the pinned PG controller, creates an
ephemeral idle thread and calls only discovery/context/status tools. It succeeds
and creates no PG state. That proves inspection compatibility, not turn submission,
approval/recovery compatibility or complete controller compatibility.

Current Codex default schemas were exported locally. Before upgrading the controller,
compare the exact methods/fields it uses, including tool discovery, skill inputs,
turn/client correlation, approval requests and recovery. Then run a fresh isolated
native task and interruption/recovery trial on the selected supported host. Updating
a version string or relying on this inspection probe alone would be inadequate.

## Concrete inspection configuration

The adjacent [example](mcp-readonly.example.toml) points to the measured private
candidate and this project. It is a reviewable configuration, not an applied global
setting. The state path is private and inspection does not create/bind it. The
example exposes only the two useful unbound inspection tools and requires startup
to succeed. A new connection must negotiate before use; existing chat tool presence
must not be inferred from editing a file.

Official OpenAI documentation supports local stdio MCP configuration with command,
args and optional cwd, plus enabled_tools and startup/tool timeouts. It also
describes CLI registration and connection inspection. Use the host's actual
configuration mechanism; this study does not install a plugin or modify global
settings. [Official MCP configuration](https://learn.chatgpt.com/docs/extend/mcp?surface=cli).

The app-server exposes a handshake and thread-local configuration for integration
clients. This study's no-turn probe uses actual 0.164.0 schema/transport responses,
rather than treating documentation as compatibility evidence for PG's pinned
adapter. [Official app-server documentation](https://learn.chatgpt.com/docs/app-server).

A practical first request after connecting is: read the selected project's current
coverage, evaluate a status-only route and explain missing coverage without creating
work. Use known current record paths; the MCP server does not select work or infer
execution permission from free text. Dynamic worker access requires the separate
controller path and a real run.

## Recommended next implementation sequence

These are bounded recommendations within the existing PG capability owners,
not new implementation authorization by this study.

| Increment | Work and owner | Acceptance evidence |
| --- | --- | --- |
| MCP1 Inspection packaging | Maintain a pinned local candidate and selected-project config; align source/binary identity and documented tool surface | A fresh intended client sees the two useful inspection tools, reads real records, handles invalid/missing selection and writes no project or operational state. The mechanism is verified here; persistent user setup remains a deployment choice. |
| MCP2 Current host compatibility | PLAN-SDP-0018 owner: compare 0.164.0 schemas, adapt supported protocol deliberately and retain old-attempt provenance | Real account/host handshake plus bounded native task, reporting, attention, interruption/reconnect, lost-response handling and protected-file checks on exact binaries. No unconditional version-pin removal. |
| MCP3 Leader access and report processing | KB052: choose the trusted leader/controller integration over existing services, retain worker restrictions and exact-revision authority | Design and test leader inspection/preparation/dispatch/report-consumption/continuation; delayed or duplicate reports cannot launch new work after changed steering. MCP calls alone supply no wake-up loop. |
| MCP4 Evidence/disposition integration | Existing PGI4/PGI5 and retained PGL gates: review/rework, recovery and canonical outcome publication | Candidate-bound evidence, actual independent review and truthful card/plan/Session transitions; crash/conflict reconciliation preserves history. Host completion is never accepted delivery. |
| MCP5 Sustained adoption | Selected project and workload after the preceding gates | Measured capacity/recovery and meaningful owner use; native parallel evidence only if parallel operation is selected. No remote hub or Ponsse application change implied. |

MCP2 is the first blocker for controlled execution on **this** workstation. MCP3
and MCP4 are additional blockers for using MCP as a normal Project Leader interface.
MCP1 inspection can be introduced separately while those are completed.

## Disposition and remaining uncertainty

Do not reopen or replace the existing Ponsse pilot model. Keep MVP1 deferred as
requested. This consolidation round does not finish the PG implementation plan,
release new binaries or accept earlier paused pilots. Independent review and native
execution on Codex 0.164.0 remain unverified. The next useful PG delivery is a
versioned current-host compatibility increment, followed by the leader/report and
evidence gates; inspection-only use can start from the measured candidate meanwhile.
