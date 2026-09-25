# Define SDP discovery, version and viewer capability metadata

| Field | Value |
| --- | --- |
| id | KB-SDP-016 |
| project | SDP |
| type | Proposal |
| CardState | superseded |
| created | 2026-09-24T20:52:36+00:00 |
| source | Owner conversation in XFMD, 2026-09-24: SDP sidebar and version/capability marker |
| next_review | Continue in KB-SDP-017 |
| tags | discovery, compatibility, SDL, SDUI, XFMD |

## Consolidation outcome — 2026-09-25

Full scope transferred to [KB-SDP-017](../active/%23017--Proposal--sdptool-and-project-navigation.md)
under `KBO-SDP-000001`. No work remains assigned to this source. The proposal
and queue statements below are historical; the successor owns current planning.
Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDP-000074. Superseded is not implemented.


## Owner request and consumer

The owner wants a version/capability file inside a project's SDP directory so XFMD
can recognize a valid SDP project, offer an SDP sidebar tab and configure generated
SDL navigation without requiring all settings as launch arguments. Requested facts
include the SDP version, SDL and SDUI language versions used, whether Agents/KanBan
exists/is supported, SDL navigation support, and inputs needed to open/generate
project views. Coordinate the contract with the agent working on SDP-vNow.

The owner first named `SDP_Version.yaml`, then `SDL_version.yaml`. Naming, case and
format remain open; do not create both or present either as an adopted schema.

Consumer-side companion: **KB-XFMD-014**, owning XFMD discovery integration, tab
layout and viewer lifecycle (the owner has further UI ideas).

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Agents/KanBan/backlog/%23014--Proposal--SDP-sidebar-and-generated-navigation.md

## Existing work and overlap review

Inspected SDP-vNow at afd9edb and XFMD at 9b93edb on 2026-09-24:

- [KB-SDP-002](%23002--Proposal--sdptool.md) is queued for project discovery and
  `view ip`. It already proposes source-based on-demand generation with prebuilt
  tools, correct viewer/source/project registration and explicit navigator refresh.
  This new card captures the specific metadata contract and XFMD consumer, not a
  competing sdptool implementation or replacement of that queued work.
- [KB-SDP-001 / R3 proposal](../../../docs/process/Project-Discovery-and-Extraction-Proposal.md)
  proposes `SDP/project.json`: schema, project ID, repository root, process profile,
  model and child-project registration. Recognition checks the selected SDP area
  then its SDP child, with no implicit parent search; this remains a proposal.
- [SDP.manifest.yaml](../../../SDP.manifest.yaml) already declares toolkit,
  framework, agent/skill versions and toolkit capabilities. It is distribution
  metadata, not automatic evidence of a valid local project.
- [SDP-project.manifest.yaml template](../../../Template/sdp-root/SDP-project.manifest.yaml)
  and [project-manifest contract](../../../Toolkit/docs/Project-Manifest.md)
  already own installed project/release/work state and project capabilities.
  Reconcile reuse, extension or an explicitly separate discovery descriptor before
  introducing another version source. Do not replace release semantics silently.
- [KB-SDP-014](../backlog/%23014--Proposal--KanBan-version-contract-and-distribution.md)
  owns KanBan compatibility/distribution. A boolean folder check is not proof of a
  compatible board contract. Coordinate capability/version declarations there.

Keep existing cards and queue state unchanged. A later consolidation may transfer
this scope into KB-SDP-002, but must preserve source/successor links and ledger
history; recording this related proposal does not perform that consolidation.

## Proposed contract questions, not a schema

Separate descriptor schema version, SDP process/profile version, product release
version and the SDL/SDUI language/profile versions in which the model is authored.
Define compatibility ranges and optional/unknown capability behavior. A process
area need not contain SDUI or KanBan; distinguish absent, declared and validated.

Candidate declarations: stable project ID; SDP/model roots; model entry; supported
navigation/viewpoint identifiers; initial/main view and navigator projection or
path; KanBan relative location/contract; SDL navigation protocol/profile; required
SDL tool/renderer capabilities. Resolve paths from an explicit base, not launch cwd.
Agree CLI override precedence and how models are shared across subprojects.

Executable paths belong to deliberate host/tool registration under the R3 proposal;
a project capability is not blanket permission to run an arbitrary repository
command. Runtime window/client/request/lease IDs must remain separate from durable
project metadata. Decide optional broker configuration without requiring a daemon
for direct on-demand generation.

## Existing XFMD navigation boundary

XFMD's DocumentViewConfig currently accepts navigator/tool/source/project/renderer/
window and optional broker. Clicking sdl-view:// for the registered project invokes
the SDL tool to generate the selected Markdown/resource bundle, then delivers it
to main or navigation in the originating window. A standalone navigation panel
already exists; XFMD owns its planned relocation/integration with the SDP tab.
SDL owns projection and generated resources. Ordinary Markdown display does not
make XFMD a BoxUI/SDUI widget host.

The XFMD companion maps --navigator, --sdl-tool, --sdl-source, --project,
--renderer and --window-id to candidate configuration/runtime sources, and records
--window/--pane/--client/--request delivery plus the optional broker/lease mode.
Reuse these contracts; do not design a second viewer command protocol from scratch.

## Next action and completion criteria

Bring this to the SDP-vNow owner/agent's discovery review alongside the XFMD card.
First deliver a decision on naming/location, relationship to existing manifests
and R3 marker, schema/compatibility, capability semantics, path bases, host tool
resolution, and initial view/navigator generation/refresh. Include example valid,
minimal, unsupported and invalid project descriptions and a consumer mapping.

Only then select bounded schema/validator/discovery integration work with tests
for repo-root versus SDP-area resolution, unrelated cwd, nested/missing projects,
optional capabilities, unsupported versions and on-demand XFMD view delivery.
No YAML file, new schema, sdptool command, migration or generated view is delivered
by this capture. No other agent has been contacted in this registration turn.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-24T20:52:36+00:00 | Codex; EVT-KB-SDP-000068 | Registered owner request and linked XFMD consumer; reviewed overlap with R3, manifests and queued sdptool scope. | Joint contract/name review before implementation; CardState remains backlog. |
