# Scrum for versioned SDP installation and upgrades

| Field | Value |
| --- | --- |
| id | KB-SDP-028 |
| project | SDP |
| type | Study |
| CardState | backlog |
| ScrumId | SCRUM-SDP-0002 |
| Systems | SDP |
| created | 2026-09-25T12:38:14Z |
| source | Owner conversation 2026-09-25 |
| next_review | Separate future Scrum before the related migration or distribution change |
| tags | maintenance, scrum, installation, upgrade, versions |

## Requested outcome

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

## Required Scrum decisions and Maintenance acceptance

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

Deliver the Scrum decision and bounded Maintenance plan first. Coordinate
[KB-SDP-027](%23027--Study--Skills-review-and-project-activation.md) for skills and
[KB-SDP-018](%23018--Study--Toolkit-audit-and-organization.md) for broader Toolkit
ownership. Extend existing mechanisms and conformance tests rather than replacing
them merely because process folders changed.

ScrumId records registration provenance in Scrum-0002, not completion of the
future topic-specific Scrum. No SprintId is assigned and no Maintenance job is
started by this card.

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
