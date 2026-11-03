# PGF2-C1c — Durable leader-directed Session host execution

T053 connects explicitly granted Session MCP operations to trusted host execution.
The owner authorized continued completion of Session0007; this increment implements
its existing HostActionGateway contract. The overall Session remains active because
integrated SDL work, in-flight disposition, full recovery/capacity and final acceptance
are still required by PLAN-SDP-0018.

## Delivered behavior

A Session host request pins its project/checkout, registered Session revision,
committed run/revision, action and named Project Leader skill. CLI, JSON invocation
and explicitly granted MCP use the same core. `session_host_read` supplies current
request inputs from one snapshot. Ordinary coordinator, worker and Session MCP
launchers receive no additional rights. The explicit host-request adapter adds
three tools to the existing fifteen Session tools.

Operational profile 0.23 retains queued, claimed, acknowledged, uncertain and canceled
requests. The controller holds a request lease across claim, host effect and settlement.
The request UUID becomes the durable outbound host input identity. Settlement derives
acknowledgement from that exact input; a missing response never triggers another send.
A claimed request without an input can be canceled under lease. Requests with an input
require existing run recovery. Pending requests block Session disposition, advance
and rebind, and cannot be bypassed by another launch input. Historical command proofs
remain attached to retained Session executions.

`execute --host-request` launches the committed assignment and selected skills, or
continues the existing attention/reply workflow. Active Session controllers consume
interrupt/scope-change on their own host connection. Scope change stops the routine
before host interruption. A separate process cannot pretend to own that live epoch.

The Project Leader launcher can grant up to eight registered Sessions explicitly.
Each has one scoped MCP server; configuration and mandatory selected leader skill
are retained. Before a leader submission the host must expose exactly the selected
server/tool inventory. A changed grant refuses retained-dispatch reuse. With explicit
`--execute-session-requests`, the launcher owns up to two separate worker controllers,
with private per-request logs and per-run leases. Owner input and leader turns retain
their separate durable queue. Leader filesystem access remains read-only. Host delivery,
observed work, verification, independent review and owner acceptance stay distinct.

## Verification boundary

Focused service and compiled adapter tests cover exact retries, canceled identities,
stale Session/run/skill/canonical bases, command lease exclusion, privileged caller and
MCP grant isolation, cross-adapter replay and read parity, lost launch response, no
blind resend, and scope-change interruption. Inventory fixtures reject extra/missing
servers/tools, duplicate servers and partial pages. Leader fixtures require the named
selected skill and registered, unique Session grants. The invocation trial found the
explicit attribution allow-list missing the new operations and the existing rebind;
the final adapter admits those controller mutations while preserving its read-only
caller restrictions.

A first native fixture setup refused before dispatch because the helper substituted
synthetic leader prose without required frontmatter. Its corrected fixture pins
canonical skills before assessment. An initial leader launch refused missing skill
discovery before preparing a turn; adding the explicit project discovery link made
that enabled skill available. These were fixture corrections, not silently discarded
passing trials. The retained earlier setup/launch results are identified separately.

Native Codex 0.166.0 executed a committed no-design task through an explicit durable
request. A subsequent native Project Leader loaded its actual skill, inspected
KB/Session/plan/assessment/scope, requested one launch through MCP, and returned
without claiming completion. The trusted launcher started a separate Master/Worker
controller. External file hashes showed only `target.txt` changed, to the requested
38 bytes. The agent reported to verification; no review or owner acceptance was
invented. A fresh trial on the frozen final candidate is recorded in the accompanying
candidate/results, including requested-versus-observed skill evidence.

The current WORK remains SessionLifecycle :00002. The scoped diagnostic Blueprint
passes constraints; its observed proposal baseline is not an accepted NOW. No target
model release or whole-Session acceptance is inferred from this increment.

## Remaining scope

This is one native leader/worker chain. It does not establish the planned two-Master
native concurrency/conflicting-third matrix, full new-goal creation/admission,
complete attention/stop/restart supervision, SDL SAD/WORK/Blueprint/owner/model-release
integration, lossless archive rotation, multi-hour operation, independent PG source
review, actual final owner trial, remote-Steering disposition or the final combined PR.
Finite store bounds remain explicit. The owner need not authorize continued work again.

## Candidate and recorded checks

[Candidate manifest](evidence-T053/host-candidate.json) pins baseline ebc5184,
37 changed product files, the three final binaries, Blueprint and actual log hashes.
The frozen full SDPTool suite passes (governance: 406.942 seconds); focused race
checks pass (governance: 27.323 seconds), and full vet passes. Observed native request,
LeaderView and LeaderContext validate against their shared schemas. Project-management
and Toolkit validators pass before final journal/ledger updates and are rerun below.

[Final native result](evidence-T053/leader-final-result.json),
[observed skill/tool calls](evidence-T053/leader-final-skill-proof.json) and
[compatibility/replay](evidence-T053/compatibility-replay.json) state their exact limits.
Lab commit 20cb4bc7e7667d00fdc0baab710dc7e50210f266 retains all four new trial fixtures;
all 634 earlier tracked lab paths are unchanged. Private runtime state and binaries
remain under `/home/warloc/.cache/sdp-workstation/T053/`; Git fixture records alone
are not a backup of in-flight host state. No merge or software publication occurred.

Final record validation passes after journal/ledger updates: 73 cards, 47 management
records, 4 lineage operations and 724 events. Toolkit validation and diff checks pass.
