# IPD — design SDPTool installation and upgrade in SDL

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0002 |
| project | SDP |
| state | completed |
| PlanType | DesignPlan |
| BranchPolicy | current |
| CommitPolicy | phase |
| Systems | SDPTOOL, SDP |
| source | KB-SDP-033; owner instruction 2026-09-27 to activate and create the installation DesignPlan |

## Outcome and authority

Extend the canonical SDL design of SDPTool to describe installation and upgrade,
including a thin gh-sdp client. Produce parser-checked design, reproducible generated
Markdown viewpoints and concrete contracts sufficient to plan Go implementation.
The owner selected this design work; this plan does not authorize live installation,
release publication or merging. It does not mark the Go engine implemented.

Primary: [KB-SDP-033](../../../KanBan/completed/%23033--Study--XFMD-SDP-adoption-and-SDL-pilot.md).
Governing requirement: [REQ-SDPTOOL-007](../../../02--Requirements/SDPTool.md).
Reuse the existing [architecture](../../../03--Architecture/SDPTool.md),
[SDL model](../../../03--Architecture/SDPTool.design),
[producer contract](../../../../SDPTool/Contract.md) and
[feature implementation plan](../../../05--Implementation/SDPTool.md).
This is the installation design work record, not a replacement feature plan.

## Scope and inspected baseline — IPD-0-M1

Baseline f10a276d176f4254e07239b10aacbee72cca747e, following
[ST1](../../../Maintenance/ST1/Plan.md). Source is root SDPTool, a Go module.
The working tree contains unrelated untracked SDL/go/sourceinput; preserve it.

| Area | Observed implementation | Design obligation |
| --- | --- | --- |
| SDPTool | Discovery, tree/selection, preview, viewer coordination; installation facts/journal reader | Extend the existing command entry point, not another Go product |
| SDL model | design-core 0.5 actors/use cases, one SDPTool process, units/interfaces and delivery status for preview/navigation | Add installation responsibility and planned delivery; preserve existing facts |
| Installer | Retained PowerShell engine with deterministic plans, ownership, backups, journal/resume and installed facts | Map reusable contracts to Go responsibilities; separate observed behavior from selected guarantees |
| Distribution | Toolkit 0.2.0 declared unreleased; no working gh-sdp CLI established by this work | Define compatible SDPTool binary/artifact identities and thin-client delegation without inventing a published release |
| XFMD | Prior preflight observed a manual five-phase bootstrap with unknown installed version | Refresh exact worktree facts before a future adoption trial; the dated snapshot is not a standard release |

Inputs: [process installation contract](../../../../Toolkit/docs/Process-Installation.md),
[profile source](https://github.com/Hans-Einar/SDP/blob/56919a19edcbd6e672c8a26661d87d16828be4f8/Toolkit/profiles/five-phase.json),
[profile artifact](https://github.com/Hans-Einar/SDP/blob/56919a19edcbd6e672c8a26661d87d16828be4f8/Toolkit/profiles/five-phase.artifact.json),
[conformance](../../../../Toolkit/conformance/install-v2/README.md), and the dated
observations in KB-SDP-033. Reconcile the old gh-sdp Study against these inputs;
its previous recommendation of a client-owned apply engine is superseded.

## Selected boundaries

SDPTool owns the installation engine and common command interface. gh-sdp locates
and invokes a compatible SDPTool; it must not independently implement migration
policy. Execution is local and must not require PowerShell on Linux/macOS.
SDL/SDUI retain language semantics. Required legacy schemas, profiles and fixtures
remain until their replacements have evidence; no blanket Toolkit deletion.

The intended wrapper workflow is `gh sdp upgrade --manifest xfmd-upgrade.yaml`.
Its exact flags and the corresponding direct sdptool interface must be specified,
including explicit project selection and the default current-directory behavior.
The custom adoption manifest describes a verified manual baseline and migration
mapping; it is not the authoritative inventory for an SDP release.

Scope includes clean install, known-version upgrade, manual-baseline adoption,
read-only preview, apply, interruption/resume and no-change repetition. It includes
preserving project documents, KanBan and historical ledger bytes, managed-file
edit handling, root instructions/skills, installed facts and Maintenance reporting.

The general XFMD SDL-modeling pilot, blueprint generation, native XFMD changes,
new SDL grammar, production installation code and release publishing are outside
this DesignPlan. KB-SDP-033 keeps the remaining adoption/pilot scope; completing
this plan must not silently close that card as a delivered live upgrade.

## Phases and milestones

| Phase | Milestone | Observable acceptance | State |
| --- | --- | --- | --- |
| IPD-0 — selection and baseline | IPD-0-M1 | Activate KB-SDP-033, link this typed plan, inspect current responsibilities and record scope/Git policy | Completed |
| IPD-1 — SDL model | IPD-1-M1 | Model installation/upgrade actors, use cases, capabilities, functionality, process boundaries, units and interfaces; allocate responsibilities to SDPTool versus thin gh-sdp; mark all new delivery activities planned; run SDL check/AST and preserve existing navigation/preview facts | Completed |
| IPD-2 — detailed contracts | IPD-2-M1 | Specify direct CLI/wrapper protocol, release inventory versus installed receipt versus adoption manifest, deterministic plan/apply binding, ownership/migration/error/recovery behavior and compatibility; trace each contract to model elements and acceptance cases | Completed |
| IPD-3 — generated review and implementation handoff | IPD-3-M1 | Generate selected Markdown viewpoints from validated SDL; record source/tool identities and coverage/gaps; walk the XFMD adoption scenario and failure cases; write a bounded ImplementationPlan with phases/milestones and remaining decisions | Completed |

## Model and contract method

Extend the single canonical SDPTool.design using syntax supported by the current
SDL Go parser. Do not maintain a competing hand-authored diagram or copy of the
model. The earlier document location spans abstraction levels; express requirement
intent, architecture and unit details without forcing every concept into a container.
Model gh-sdp as its actual collaborating process responsibility; record system-level
meaning in prose if the current grammar cannot represent it faithfully.

For supported Channel/contract/data/scenario constructs, model actual designed
exchanges and validate them. Never invent keywords to fill a diagram. If a required
sequence or data relationship cannot be expressed/projected, record the exact gap
and its impact instead of claiming full SDL coverage. Detailed schemas and protocol
examples may remain linked authored contracts; generated views only reflect facts
present in the model.

Resolve in IPD-2:

- Request/response and error contracts, protocol/engine compatibility, process
  invocation, exit status, cancellation and ownership of generated artifacts.
- Release selection and provenance verification, binary retrieval/cache/offline
  behavior, schema/capability versions and accurate installation receipts.
- Ownership classification and old/current/new comparison; explicit renames and
  relocation prerequisites; handling missing/modified managed files and unknown
  baselines without overwriting project-owned work.
- Exact input-bound plan review/apply; locks, before/after hashes, backup and journal
  finalization. Preserve the distinction between forward resume, per-file atomic
  replacement and whole-operation rollback; do not promise the latter accidentally.
- Project instruction/skill migration, append-only history preservation, generated
  Maintenance entries and how existing SDPTool discovery sees incomplete/completed
  operations. Local manifests cannot authorize arbitrary executable migration code.
- A narrow shared set of canonical fixtures usable by the Go engine and the retained
  installer, with explicit differences and conditions for retiring legacy execution.

## Generated review and verification

Use the repository's SDL toolkit to parse, validate, emit AST and generate viewpoints.
Do not write Mermaid or modify generated Markdown by hand. First verify the existing
catalog's actual capabilities, then select views covering use cases/features,
architecture, allocation and designed exchanges/data. Mark unsupported projections
as gaps. Keep a compact index and links to generated views; do not pre-render an
unbounded set. Retain model revision, generator identity and output hashes so a
second run can reproduce the result.

Scenario review must cover at least: empty-project install; known installed version
upgrade; manual XFMD adoption with unknown origin; changed input after preview;
managed versus project-owned edits; collision/unsafe path/unsupported schema;
interrupted apply and resume; repeated no-change operation; incompatible wrapper
and SDPTool; unavailable or invalid release artifact. These are design acceptance
cases now, not executed Go-installer tests. Refresh the live XFMD inventory read-only
before deriving a disposable fixture; no live mutation is part of this plan.

At IPD-3, map these cases to implementation slices and executable checks. Include
an end-to-end disposable XFMD trial before a separately authorized live adoption.
Record review findings and unresolved decisions honestly; generated diagrams and
parser success are not runtime correctness or owner acceptance.

## Git policy and delivery

Continue on sdp/study-xfmd-install-adoption. Use one phase commit per completed
phase, with its milestone ID and delivery evidence. Push completed phases under
existing authorization; do not merge or publish releases. This small design plan
does not change the historical feature plan's stacked implementation branches.
The later ImplementationPlan must explicitly select its own Git policy.

## IPD-0 handoff — historical

IPD-0 inspected the current model, architecture, requirement, producer/installer
contracts and management state at f10a276 plus this planning diff. This delivery
only creates the plan and activates its card; no installation SDL model additions,
generated installation views or new runtime behavior are claimed.

IPD-0 verification: project-management replay passes (41 cards, 13 management
records, 280 events); the documentation check preserves 105 frozen records/prefixes
and 574 generated outputs and validates local links. `git diff --check` passes.
These checks validate this planning delivery, not installation model/runtime
behavior. Next execution milestone is IPD-1-M1.
CardState is ready after the plan handoff and becomes in-progress when modeling
begins. The plan stays active until its actual design deliverables are complete.

IPD-1-M1: [Evidence](Evidence.md) records the checked model and preserved baseline.
Next: IPD-2 contracts and modeled exchanges. CardState is in-progress.

IPD-2-M1: [Contract](Contract.md) and [Scenarios](Scenarios.md) define execution,
identity/trust, preservation and recovery. Twelve channels and five scenarios
validate in SDL; detailed schemas and runtime tests remain implementation work.

## IPD-3 delivery and review

All selected design phases have delivered their artifacts. [Generated review](review/index.md)
is reproduced through the SDL CLI by [review.py](review.py), with source/tool hashes,
semantic negative checks and repeatability evidence in [Evidence](Evidence.md).
The [working contract](Contract.md) and [scenario review](Scenarios.md) state remaining
runtime/schema limits. [PLAN-SDP-0003](../../../05--Implementation/SDPTool/Installation/Plan.md)
is a planned ImplementationPlan, not an executed Go implementation.

Owner review can now assess concrete choices: preview-only default and explicit
saved-plan apply, descriptor trust/key distribution, cached executable delegation,
forward-only recovery and the proposed implementation phases. KB-SDP-033 stays
active at gate-review for this actual deliverable, with the later adoption/SDL
pilot explicitly outstanding. No independent approval or owner acceptance is inferred.
