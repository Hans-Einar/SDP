# MAINT-SDP-0003 — versioned SDP installation and upgrades

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0003 |
| project | SDP |
| state | active |
| source | KB-SDP-028; owner decision 2026-09-25: direct card to Maintenance, without Scrum |
| Systems | SDP, SDPTOOL |

This job is selected directly from [KB-SDP-028](../../KanBan/completed/%23028--Study--Installer-upgrade-and-versioned-layout.md).
No Scrum, Sprint or wrapper card is required. The card's Study is fulfilled by
this plan; this Maintenance owns the implementation.
The owner authorized execution on 2026-09-25. Live consumer rollout remains excluded.

## User outcome

A project can initialize or update its SDP process using a versioned, reproducible
configuration. Before mutation, the tool explains its observed baseline, target
profile, changes, conflicts and preservation. After successful verification,
machine-readable installed facts and a Maintenance report describe what actually
changed, including old/new versions and an explicit unknown baseline where needed.
Repeating an already completed operation does not rewrite files or invent another
successful upgrade. SDL/SDUI design sources and project history remain project-owned.

## Inspected baseline

Repository commit 6cf74e0; preceding [skills Maintenance](../SK1/Plan-and-Evidence.md)
is complete. This planning milestone changes no installer or consuming project.

| Existing authority | Reuse and gap |
| --- | --- |
| [Toolkit facts](../../../SDP.manifest.yaml) | Unreleased 0.2.0; release/capability facts, distinct from process profile and project product version |
| [Install inventory](../../../Toolkit/SDP-install.manifest.json) | Explicit copied/generated entries, ownership and policies; not a compiler for the new process layout |
| [Templates](../../../Template/README.md) | Distributed numbered folders still include separate Study and DesignAnalysis and end at 07--Implementation |
| [Local SDP structure](../../README.md) | Five numbered phases, phase-local studies, shared board/management history; not automatically installed elsewhere |
| [Install contract](../../../Toolkit/docs/Installation-Contract.md) | Existing PowerShell engine, mutation-free plans, preservation, backups, path containment and nineteen conformance scenarios |
| [Migration guide](../../../Toolkit/docs/Installer-Migration.md) | Current missing-manifest behavior treats it as pre-versioning; does not establish an observed process version; no complete transactional rollback guarantee |
| [Installed facts schema](../../../Toolkit/schemas/installed-toolkit-manifest.schema.json) | Strict 1.0 schema has Toolkit/Framework/AGENTS/installer/skills/capabilities but no process profile or configuration identity |
| [Project manifest](../../../Toolkit/docs/Project-Manifest.md) | Project-owned identity/release/work facts and pointer to installed facts; preserve it |
| [SDPTool contract](../../../Toolkit/SDPTool/Contract.md) | navigation.json declares bindings and navigation eligibility, not installation conformance; readers currently support their declared schemas |
| [Management contract](../../ProjectManagement/README.md) | Management events and reports belong here; system design/code verification remains in Traceability |

SK1 installed all thirteen canonical skills and their references without moving
consumer destinations from .codex/skills. Preserve its native metadata contract,
version facts and same-Toolkit-version ForceManagedFiles preflight; do not restore
old role sources. The full Toolkit validator still has 38 recorded pre-existing
Traceability errors under KB-SDP-011. Compare against that baseline explicitly.

## Proposed implementation boundary

Use JSON as the initial authored layout configuration, matching the existing
inventory. Keep one authority for each fact: a versioned process/layout definition
feeds a deterministic build that emits the installation inventory and template
bundle. Release facts stay in SDP.manifest.yaml; installed facts record the exact
selected profile, capabilities and configuration identity. Do not introduce YAML
plus JSON copies that both need manual maintenance, or another competing version
registry. Exact filenames/field names and schema versions are finalized in IU1-M1.
These are planned contracts, not currently supported fields or commands.

The target layout covers 01--Mandate, 02--Requirements, 03--Architecture,
04--Design and 05--Implementation, plus the shared KanBan, ProjectManagement,
Maintenance, Traceability, Agents/Scrum and optional Sprint/review areas. Studies
belong to the affected phase. Systems and actual containers organize source
content; folder placement never assigns an abstraction level. New templates carry
neutral placeholders only, never this repository's cards, IDs, ledgers, model
sources, active work or publication claims.

Retain the current PowerShell installer as the first execution engine; extend its
explicit plan/apply contracts instead of introducing a second independently
implemented engine. Reuse Go SDPTool for project recognition and declared-fact
consumption. Any install/update facade must delegate to the same engine/artifact
and preserve machine-readable outcomes. A Go installer rewrite, GUI changes and
new parser/runtime work are outside this job. Prebuilt Go tool packaging/version
checks must be defined; ordinary documentation viewing must not build tools.

## Compatibility and ownership rules

- Separate Toolkit release, process/profile, management/KanBan schema, skill,
  language and consuming-product versions. navigation.json's current process
  label is an eligibility claim, not proof of an installed process version.
- Extend the installed-facts authority through an explicitly versioned schema
  and compatible reader rollout. Its current additionalProperties=false means
  silently adding fields to schema 1.0 is not a compatible solution. Test old
  facts with new readers and unsupported new facts with old readers.
- Supported baselines initially include a clean project, the distributed legacy
  template profile, the current local five-phase/shared-management layout and
  a manually initialized XFMD-shaped fixture. Inventory actual differences and
  publish the supported transition matrix before enabling each migration.
- Missing facts mean unknown until inspection supports an adoption classification.
  Do not infer a version solely from folder names. Report mixed/ambiguous layouts
  with conflicts and proposed mappings; never silently select a destructive path.
- Managed file replacement follows declared policy and backup rules. Project-owned
  documents, IDs, ledger bytes, model sources, navigation bindings and local
  customization are preserved. Relocation needs an explicit source/destination
  map, content preconditions and verified link updates, not a copy-overwrite loop.
- Preserve legacy Study/DesignAnalysis material in a declared archive or explicit
  destination. Semantic repartition of design documents requires its own reviewed
  mapping; do not mechanically infer content ownership from old folder numbers.
- KanBan migration must identify card/ledger profiles, preserve IDs, Ref links,
  lineage and append-only history, and avoid duplicate writable ledgers. The
  minimum version/compatibility contract needed to distribute this profile is in
  scope; standalone KanBan packaging remains owned by KB-SDP-014.

## Plan, apply, recovery and reporting

The build artifact includes schema/profile identity, source/configuration digest,
explicit source inventory, ordered migration steps and capability prerequisites.
Identical inputs must produce identical bytes; runtime timestamps and project
observations belong to the operation record, not the reproducible build identity.
Build validation rejects missing sources, duplicate/overlapping destinations,
unknown operations and unsafe paths before publishing an artifact.

Inspection produces a deterministic reviewable plan with old/target profiles,
unknown facts, source and destination hashes/preconditions, preservation/conflicts,
link updates and required capabilities. Preview is read-only. Apply must consume
or prove equivalence to that exact plan, check source/target drift, prevent
concurrent writers and check physical containment again before writing.

Use operation identity and a durable step journal for multi-file migration.
Back up bytes before replacement and record completed/failed steps. Select a
bounded resume/recovery policy in IU2-M1; never claim whole-operation atomicity
from per-file renames. A crash before final facts/report publication must leave
a detectable incomplete operation, with retry-safe finalization and no duplicate
ledger event. If full rollback is unsafe, stop with preserved data and explicit
recovery instructions rather than overwrite edits made after failure.

A successful operation writes a readable Maintenance report in the target SDP
area and links it to that project's management history. Record operation identity,
observed old version/profile or unknown, actual target versions/capabilities,
configuration digest, changed/preserved/conflicting paths, backup/recovery location,
verification results and outcome. Failed/partial/recovered operations must be
reported as such. No-change runs report no-change without another upgrade entry.

Plan the bootstrap case: a project may not yet have a management ledger or may
have an older supported one. Establish/migrate the selected history contract
before finalizing the operation's report/event; preserve old history and assign
IDs from that target, never copy MAINT-SDP-0003 into every consumer. If history
cannot be safely finalized, the operation is incomplete, not silently successful.
System/tool code changes made here receive Traceability records referring to this
Maintenance; target installation movements are management history, not proof that
its product features are implemented.

## Phases and milestones

Execution uses one stacked branch per phase and a commit for each delivered
milestone. Update this plan/evidence and review related backlog at every milestone;
push completed phases, offer a combined PR against sdp-vNow, and do not merge.

| Phase / branch | Milestone | Concrete acceptance | State |
| --- | --- | --- | --- |
| IU1 / sdp/phase-iu1-profile-build | IU1-M1 | Finalize implementable profile/configuration, ownership/transition matrix, version/read compatibility, operation/report and neutral-template contracts against existing code and fixtures | Completed |
| IU1 | IU1-M2 | Implement deterministic configuration build and validated target templates/inventory; reproducibility and negative cases pass; preserve an explicit supported legacy artifact/reader boundary | Completed |
| IU2 / sdp/phase-iu2-safe-upgrades | IU2-M1 | Implement inspection and mutation-free exact-plan generation for clean, versioned, local and manual baselines, with explicit conflicts/unknown facts and input identities | Completed |
| IU2 | IU2-M2 | Implement apply, preservation, backups, ordered migrations and crash/concurrency recovery; verify target facts and retry-safe Maintenance/history finalization | Planned |
| IU3 / sdp/phase-iu3-consumer-validation | IU3-M1 | Update installed-fact readers and supported navigation bindings; define prebuilt tool distribution; installed SDPTool discovers actual facts without invented models/viewer settings | Planned |
| IU3 | IU3-M2 | Complete end-to-end transition/rollback-resume conformance, independent review and operator guidance; disposition remaining gaps before Maintenance closure | Planned |

IU1-M1 resolves field/version identifiers, exact packaging destinations, supported
host matrix and report finalization semantics before code depends on them. This
is ordinary implementation design within the selected boundary; escalate only
material unresolved owner choices, not routine technical decisions already covered.

## Verification and completion

Test clean install, each supported old-to-target transition, current local profile,
manual adoption, ambiguous/mixed baseline, repeat/no-change, unsupported schema,
downgrade, customization/conflicts, missing tools and read-only preview. Exercise
source/target drift, path escapes, links/reparse points, collision/case behavior,
interruption at each mutation/finalization boundary and concurrent attempts.
Verify recovery, exactly-once operation/event finalization, byte-preserved history,
IDs and working relative links. Test target reader behavior on both old/new facts.

Use disposable fixtures, including a provenance-recorded XFMD-shaped snapshot;
inspect the actual XFMD baseline when preparing that fixture. Do not perform a
live XFMD upgrade or change its application from this plan. Its agent can use the
completed tool/report contract through its own selected work. Any live consumer
rollout must identify the concrete project and supported baseline separately.

Completion means the built configuration drives a usable, tested install/update
workflow, installed versions/capabilities are truthful, history and user content
survive, and failures are recoverable and accurately reported. It does not mean
all historical documents have been semantically reorganized, every host is
supported or a Toolkit release has been published.

## Dependencies and remaining scope

Reuse SK1 and the existing install-v1 conformance suite. Coordinate
[KB-SDP-014](../../KanBan/backlog/%23014--Proposal--KanBan-version-contract-and-distribution.md)
for minimum distributable management compatibility; keep its independent-consumer
packaging/publication scope open. [KB-SDP-018](../../KanBan/backlog/%23018--Study--Toolkit-audit-and-organization.md)
retains the broad Toolkit audit. KB-SDL-005/KB-SDP-020 own language/source-set
semantics and shared model splitting; installer work must not bypass them.
No graphical history, implementation-plan synthesis or native XFMD integration
is absorbed from other cards. Existing #010 owner review remains non-blocking.

## PM4-M1 planning verification

This planning delivery was checked against the current inventory, templates,
strict installed/project schemas, SDPTool navigation contract and SK1 results.
No new process fields, commands or recovery guarantees are described as already
implemented. The owner-selected direct card → Maintenance route is supported by
the existing management contract; no methodology change is needed.

Checks pass: 35 cards, eight management records, three lineage operations and
231 events; four management test groups and 15 lineage negative cases. Existing
ledger bytes remain an exact prefix and Traceability is unchanged. Documentation
verification preserves 105 frozen records/prefixes and 574 generated artifacts
and resolves 2,397 local links and 130 fragments. git diff --check passes.
These checks verify planning/history integrity, not future installation behavior.

## IU1-M1 execution decision

The [versioned installation contract](../../../Toolkit/docs/Process-Installation.md)
fixes the artifact, transition matrix, explicit 2.0 facts/plan boundary and
forward-resume recovery policy. Linux PowerShell 7 is the initial verified host;
Windows remains experimental until tested. XFMD baseline cf11709e was inspected
read-only. KB-SDP-014/018 retain standalone distribution and wider audit work.
