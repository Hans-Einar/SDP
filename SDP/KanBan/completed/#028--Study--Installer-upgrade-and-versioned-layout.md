# Versioned SDP installation and upgrades

| Field | Value |
| --- | --- |
| id | KB-SDP-028 |
| project | SDP |
| type | Study |
| CardState | completed |
| ScrumId | SCRUM-SDP-0002 |
| Systems | SDP |
| created | 2026-09-25T12:38:14Z |
| source | Owner conversation 2026-09-25 |
| next_review | Closed; implementation delivered by MAINT-SDP-0003; live rollout separately selected |
| tags | maintenance, scrum, installation, upgrade, versions |

## Original requested outcome

The owner requests a future Scrum resulting in Maintenance work on SDP install
and upgrade scripts. The desired workflow builds a versioned declarative layout
and applies it programmatically to a project's SDP area, including existing
manually initialized projects such as XFMD. The conversation's isolated “SDL
folder” wording is interpreted as the SDP process tree; language sources are
project-owned and are not moved by implication.

## Existing mechanisms to reuse

| Authority | Existing responsibility |
| --- | --- |
| [SDP.manifest.yaml](../../../SDP.manifest.yaml) | Toolkit release and capability facts |
| [SDP-install.manifest.json](../../../Toolkit/SDP-install.manifest.json) | Explicit copied/generated inventory, ownership and install/migration policies |
| [Project manifest](../../../Toolkit/docs/Project-Manifest.md) | Project-owned release/work facts and relative pointer to installed Toolkit facts |
| [Installed manifest schema](../../../Toolkit/schemas/installed-toolkit-manifest.schema.json) | Generated Toolkit, Framework, AGENTS, installer, skill and capability versions |
| [Install contract](../../../Toolkit/docs/Installation-Contract.md) | Planning/apply semantics, managed versus project-owned files |
| [Conformance scenarios](../../../Toolkit/conformance/install-v1/README.md) | Portable plans, upgrades, preservation and failure cases |

The supported entry is Toolkit/scripts/Install-SDP.ps1; New-SDPBuildIdentity.py
records build identity, not a completed new layout compiler. The local five-phase
process and common management ledger have advanced beyond distributed templates.
Do not invent another competing version file or claim current installers already
migrate this profile. JSON is already used for inventory; choose any YAML source
or build-to-JSON step only with a concrete authority and reproducibility reason.

## Required planning decisions and Maintenance acceptance

Define the versioned process/layout profile and deterministic build/validation
of the configuration consumed by installers. Reconcile existing schemas and
capabilities with KB-SDP-022 discovery and KB-SDP-014 KanBan distribution. Keep
project product version, process/profile version, Toolkit release and authored
language versions separate. Make installed version/capabilities machine-readable
through an agreed existing manifest or justified compatible extension.

Plan ordered upgrade steps with explicit source/target profiles, dry-run diff,
preconditions, backups/recovery, project-owned data preservation and repeat-run
idempotence. Handle unsupported schemas, downgrade attempts, local modifications,
interrupted apply, path hazards and old/absent manifests before mutation. Unknown
or manual XFMD installations need inspected baseline/adoption; never invent an
old version or overwrite project content to make detection succeed.

A successful upgrade must emit a Maintenance entry describing actual old/new
versions (or explicit unknown baseline), configuration/source identity, changes,
preserved/conflicting files and verification, linked to project-management history.
Define failure/rollback records without claiming a failed upgrade succeeded.
Traceability receives system changes only, not duplicate management transitions.
Test clean install, supported upgrade, current local profile, repeat install,
failed/partial recovery and a fixture of manually adopted XFMD. A real XFMD
migration requires its own selected work; this card does not perform one.

The later owner decision below selects a direct Maintenance plan without Scrum. Coordinate
[KB-SDP-027](../completed/%23027--Study--Skills-review-and-project-activation.md) for skills and
[KB-SDP-018](../backlog/%23018--Study--Toolkit-audit-and-organization.md) for broader Toolkit
ownership. Extend existing mechanisms and conformance tests rather than replacing
them merely because process folders changed.

ScrumId preserves original registration provenance in Scrum-0002. The owner
subsequently selected direct Maintenance planning; no new Scrum or Sprint is
created. MAINT-SDP-0003 owns execution separately from this Study.

## Worklog

2026-09-25T12:38:14Z — EVT-KB-SDP-000132: registered; future Scrum and Maintenance execution remain unstarted.

## Sprint-0001 handoff

Sprint-0001 now implements local SDP/navigation.json bindings, including an
explicit default model, model/source lists and optional project-manifest pointer.
It intentionally does not duplicate installed Toolkit versions. The facade reads
installation facts as declared and distinguishes navigation eligibility from full
installation validation. Include binding creation/preservation, default model,
local board profile support and prebuilt Go tool distribution in the future
upgrade/configuration build design. See Toolkit/SDPTool/Contract.md; no installer
or XFMD migration was performed by this Sprint.

2026-09-25T13:42:06Z — EVT-KB-SDP-000153: future Scrum remains unstarted.

## Scrum-0003 skills handoff

[MAINT-SDP-0002](../../Maintenance/SK1/Plan-and-Evidence.md) is now planned from
completed skills Study #027. It selects one root Skills/ source and a tested
Codex .agents/skills symlink adapter. The current installer still copies ten
legacy roles into .codex/skills, while the new router sdp is excluded by the old
validator's skill-ID pattern. Native name/description and distribution version
metadata must coexist. Root SKILL.md alone did not register a skill in the
installed Codex CLI 0.156.1 probe. File discovery does not prove agent behavior.

This card's Scrum is recommended next, before the skills job changes distribution.
Set the shared ordering/compatibility contract for canonical source migration,
host adapters, manifest facts, generated AGENTS and repeat installation. Keep
SK1's mechanical skill-consumer migration bounded; this card owns broad process
profile upgrades and manual-project adoption. It remains an unexecuted Scrum,
not a delivered installer merely because the dependency is now documented.

2026-09-25T14:18:18Z — EVT-KB-SDP-000157: handoff recorded; card remains backlog.

## Completed SK1 handoff

The owner selected skills Maintenance before this Scrum. [MAINT-SDP-0002](../../Maintenance/SK1/Plan-and-Evidence.md)
is now complete: root Skills/ is canonical, .agents/skills exposes this repository,
and generated .codex/skills copies retain the existing installed destination.
Codex CLI 0.156.1 discovers both tested adapters; do not assume that all hosts do.
All thirteen roles and shared references are explicit installation inventory.
Native metadata v2 coexists with legacy installed metadata validation; the bare
sdp router is supported. The Toolkit remains unreleased 0.2.0; AGENTS contract
is 2.0.0, ten replacement roles are 2.0.0 and three new roles 1.0.0.

A same-Toolkit-version skill contract change fails before writes unless using
ForceManagedFiles; reviewed force plans preserve backups/project-owned content.
Nineteen normative conformance scenarios and the faithful old-installation probe
cover this boundary. Reuse these mechanisms; do not reimplement or silently bypass
them when designing broader process-profile migration. No real consuming project
was upgraded. Versioned layout builds, manual adoption and old/new-profile
Maintenance entries remain this card's future Scrum and implementation scope.

2026-09-25T14:48:15Z — EVT-KB-SDP-000161.

## Direct Maintenance selection — 2026-09-25

EVT-KB-SDP-000163: selected and in-progress for planning. The owner chooses
KBCard → Maintenance because this is one coherent card; the earlier suggestion
of another Scrum is superseded. This changes the planning route, not the scope.

## Planning outcome

[MAINT-SDP-0003](../../Maintenance/IU1/Plan.md) contains the selected Maintenance
plan: three phases, six milestones, current authority map, migration/preservation
rules, report/recovery semantics and acceptance scenarios. This completes the
Study's planning outcome only. Implementation and any actual consuming-project
upgrade remain planned; no new Scrum, Sprint or wrapper card has been created.

2026-09-25T15:03:56Z — EVT-KB-SDP-000164: planning complete; Maintenance execution not started.

## Maintenance completion

2026-09-25T16:02:42Z — EVT-KB-SDP-000167: [MAINT-SDP-0003](../../Maintenance/IU1/Plan.md)
is complete across IU1–IU3. The explicit profile/configuration drives read-only
plans, journaled apply/resume, truthful installed facts and target Maintenance
history. Updated SDPTool consumes both supported fact schemas.
[Verification](../../Maintenance/IU1/Evidence.md) includes independent review,
forced interruption and a disposable copy of the actual XFMD baseline.
No live consumer was upgraded. Earlier worklog statements describe their dated
planning state, not the current implementation outcome.

The later typed-plan and Planning-skill direction is tracked separately in
[KB-SDP-029](../active/%23029--Proposal--Typed-plans-and-planning-skill.md).
