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

## PGL2-C1 preparation contract — 2026-10-10

Owner T015 continues KB052. Select one runnable preparation increment before adding
launch authority: trusted local controller prepares a SAD or Master assignment,
inspects its pinned inputs, and refuses stale inputs or overlapping local scope.
Prepared assignments are not runs, attempts or host threads. They share the existing
checkout store, lock and atomic generation publisher; no second management ledger
or model/blueprint database is created.

SAD requires an active DesignPlan or ArchitecturePlan; Master requires an active
ImplementationPlan. Reuse canonical card/plan/Session/authority source checks and
management lineage. Keep the older Maintenance route unchanged. The API pins the
complete coverage revision and references a controller-selected authority source;
it does not interpret natural-language authorization or prove an owner disposition.
Worker capabilities cannot prepare assignments. A caller-supplied role is only the
assignment label; the trusted adapter determines the caller capability.

A SAD package pins NOW. A Master package pins NOW, TARGET, task, blueprint and a
design-disposition record. Each is an existing regular project file with SHA-256,
read through the core's confined/symlink-refusing file reader. These are exact file
references, not new ModelGovernance identities or an assertion that those services
validated or approved them. Provider-aware revision resolution and authenticated
owner design disposition remain PGL2-C2 prerequisites; inspection explicitly says
`dispatch: unavailable` and `designAcceptance: not-assessed`. A payload saying
approved cannot grant launch. The prepared package also retains bounded outcome,
writable/excluded paths and return criteria.

Prepare requires explicit checkout binding, expected current generation, a fresh
operation UUID and matching coverage/file digests. Under the existing store lock,
reject parent/child/equal writable-path overlaps with any prepared assignment or
nonterminal legacy run in this checkout. Paths are literal, not globs; excluded
subpaths do not subtract from this conservative overlap check. No global worktree
reservation is claimed. Other checkout preflight/claims remain external prerequisites.
Validate sources again before publishing; files are not locked against outside
writers, so inspection and future dispatch must always recheck them.

Persist immutable prepared identity/revision 1 plus the idempotency receipt in the
same generation. Same operation/content retry returns the original prepared record,
including after an uncertain committed save; conflicting content refuses. It does
not return a fresh eligibility claim. Recheck through inspect. The operation-ID
namespace is shared with legacy receipts; collisions refuse in both directions.
Inspect lists assignments in stable ID order and reports changed/missing sources.
No TTL expires a scope claim. Retirement/revision/continuation belongs to PGL2-C2;
this first increment never releases scope merely because a source is stale.

The first preparation explicitly upgrades that store to operational profile
`sdp-project-governance-state/0.2`. New readers retain 0.1 support. Old binaries must
refuse the new generation rather than drop prepared records; use a separate PGL2
trial store and never upgrade the owner's original C6 lab implicitly. Application
Maintenance wire profile and worker MCP stay unchanged. No stable distribution or
online host is installed. Role skill adoption from PLR1 is deliberately sequenced
with the integrator-owned distribution boundary before role launch, not copied over
release manifests during this preparation slice.

Commands: `sdp-controller assignment context|prepare|inspect --project ROOT --state
DIR`, JSON input via `--input FILE` for context/prepare, optional `--assignment UUID`
for inspect. Context reports canonical coverage plus the operational generation;
prepare persists only after the checks above. All responses identify their narrow
preparation boundary. No model/account access occurs. A reproducible separate
local trial must show SAD and Master preparation, retry after restart, conflicting
scope refusal and a changed TARGET producing stale inspection.
