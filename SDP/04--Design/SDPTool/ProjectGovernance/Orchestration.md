# Project Leader and SAD orchestration — KB052

Owner direction: Session0007 T014, 2026-10-10, incorporating the
[revision-bound handoff](Project-Leader-Handoff.md). This refines the earlier
[PGD1 design](Design.md); it does not claim its new operations already exist.
[PLAN-SDP-0018](../../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
retains the implementation branch and milestone policy.

## Outcome and current gap

The owner normally talks to Project Leader, can also address SAD/Master directly,
and sees consequential feedback reflected once in shared work. Project Leader
commissions bounded work and stays available while workers execute. SAD prepares
ModelGovernance WORK/TARGET and generated blueprint evidence; Master owns one
implementation Session and its phases. Steering provides strategic, revision-bound
advice. Role guidance does not create runtime privileges or owner acceptance.

The owner reports having tried the lab and needing more familiarity. This is use
feedback, not rejection or acceptance. The old lab demonstrates one Maintenance
Master and cannot demonstrate the new coordination model. Keep its actual state and
candidate intact. Do not require further old-lab acceptance before selected KB052 work.

Current runtime has one checkout-local run store, run/attempt/host identities,
controller-only launch, restricted worker MCP, replay receipts, exact skill inputs
and correlated recovery. Both controller leases are checkout-wide. Master context,
coverage and routine selection are Maintenance-only; the role validator accepts
only Master. Simply changing a prompt or launching multiple current controllers
would not implement the proposed system. No native child mapping is proven.

## Selected architecture and alternatives

Reuse the Go governance services and existing app-server controller. Extend their
operational state and adapters in bounded increments; preserve existing Maintenance
clients and their retry/authority contracts. Extend the dedicated terminal client
rather than selecting an unverified stock-TUI plugin or maintaining a Codex fork.
The first view is a standalone read-only command; a later event-loop client owns
conversations and panels. No XFMD changes or general KB040 KanBan TUI are selected.

Use logical assignment identity independent of run/attempt/thread. Explicit
controller-issued capabilities distinguish Project Leader operations from worker
reports. A payload role string grants nothing. Future dispatch uses expected
assignment/bundle revisions and stable operation identity; uncertain sends must be
reconciled against the original attempt before any new launch. Existing worker MCP
keeps its five restricted operations until a separately tested change is selected.

The controller owns durable inbox delivery and per-thread turn serialization.
Project Leader uses short turns over that state; it is not the event loop. Prioritize
owner input, coalesce reports by stable identities, retain them across restarts,
and revalidate delayed reports against newer steering before dispatch. Host turn
completion, reported phase completion, reviewed evidence and owner acceptance
are distinct. Delivery exactly once means idempotent processing, not a promise
that transport never duplicates a message.

Design-bound implementation packages reference existing ModelGovernance and
Blueprint services, exact NOW/TARGET/task/blueprint revisions, allowed/protected
paths, source evidence and actual design disposition. A changed dependency makes
the package stale before continuation. Do not implement another model-history or
blueprint store. General Maintenance remains usable without model prerequisites.
Integration of those separately developed services needs a refreshed preflight and
selected combined baseline before dependent code changes.

Shared cross-worktree claims/IDs need one selected authority, not branch-local
registries advertised as global exclusion. Until that capability is implemented,
use the existing manual preflight and explicit coordinator sequencing. The current
known ID collisions remain integration work; preserving both histories is required.

## PGL1: observable work before orchestration

Add `sdp-controller overview --project ROOT --state DIR`, with readable default
and `--json`. It inspects every registered run in that one store, using a single
persisted generation. No host process, registration, new turn or approval occurs.
This is a registered-work view, not a complete project board or a live scheduler.

Show run/attempt/thread identity, task, selected Session/step/card/plan, routine state,
current coverage and pinned input revision. Show the last recorded host status and
timestamp as historical observations, with gaps and current liveness unknown.
Assigned role comes from controller provenance; never infer it from a run name.
Show observed tool-completion count/last identity, explicitly requested skill count
and acceptance separately from unobserved compliance. Child mapping and owner
acceptance stay unknown/unobserved. No transcript, owner-input content or capability
credentials belong in the overview. Terminal strings must not execute control codes.

Per-run source checks occur during inspection and are not an atomic filesystem
snapshot. Unavailable/stale records must be visible; operational-store corruption
fails the operation. Missing operational state yields an empty, unbound view without
creating a directory. Sort by stable run identity. Never derive live state from an
arbitrary observation-age threshold. This is the first part of handoff scenario 5;
all multi-agent orchestration scenarios remain open.

## Verification and next contract

PGL1 tests must distinguish two runs, stale coverage, no attempt, old/reconnected
observations, completed host versus unfinished routine, and terminal-control text.
Run the compiled command against the owner's existing lab state read-only and verify
both project and state bytes remain unchanged. This tests inspection, not a new
model workflow or independent review. Freeze a separate development binary for use.

Next PGL2 specifies and implements bounded SAD/implementation assignments with
revision gates, persisted identity, explicit capabilities and an executable
prepare/inspect/refuse path. Then dispatch/inbox/continuation and direct conversation
follow, with the six handoff scenarios retained in the plan. Online Steering needs
separate selected-host evidence and cannot be proved by local skill files.
