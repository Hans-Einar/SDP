# Standard Document Procedure (SDP)

Status: working draft  
Toolkit-Version: 0.2.0 (unreleased)

Start with the [documentation navigator](docs/README.md), the
[SDL language/tooling overview](SDL/docs/README.md), or [SDUI](SDUI/README.md).
The repository's own work lives under [SDP](SDP/README.md); reusable project seeds
live under [Template](Template/README.md). The R1 reorganization preserves the
installed Toolkit process while its next SDL-based phase profile is designed.
The dated [checkpoint #1](SDP/History/checkpoint-1/README.md) remains a shared
SDP/SDL/SDUI history entry, not the sole source of current implementation status.

SDP is a repository-local, document-driven method for AI-assisted software
development. It connects mandate, research, requirements, architecture, design,
implementation, verification, review, traceability, releases and handoff to the
repository instead of relying on chat memory.

Core principle:

> Design horizontally. Implement vertically.

## Repository layout

```text
SDP repository
|-- Template/                 neutral project-owned templates
|-- Toolkit/                  installer, schemas, managed payload and skills
|-- SDP/                      this project's KanBan, records, studies and history
|-- SDL/     SDL (moving to SDL in R1-M3)
|-- SDL/                      SDL project KanBan
|-- SDUI/                     UI language and project documentation
|-- SDP.manifest.yaml         Toolkit release/capability manifest
|-- RELEASE-NOTES.md          Toolkit release notes
|-- docs/                     documentation navigator and SDP process
`-- examples/                 neutral Toolkit contract examples
```

Project records under `SDP/` are not template inputs. Installation sources are
explicitly listed in `Toolkit/SDP-install.manifest.json`. Root-level numbered
seed duplicates were consolidated into `Template/`; installed destinations have
not changed. The [R1 plan](SDP/Maintenance/R1/Plan.md) and [document map](SDP/Maintenance/R1/Documentation-index.md) record the migration and remaining work.

## Version model

Released software and SDP Toolkits use SemVer. The Toolkit that existed before
formal version metadata is documented as migration baseline `0.1.0`; this
backward-compatible capability addition targets `0.2.0`. It remains
**unreleased** and no corresponding tag or GitHub Release exists.

Sprint/Refactor, Iteration, Slice/Fix, revision and Git SHA are separate
development coordinates. See `Toolkit/docs/Release-And-Versioning.md` and
`Toolkit/docs/Development-Identity.md`.

## Install or update a project

Use an independent clone or an extracted GitHub source archive. Let `$SdpSource`
name its root, the directory containing `Toolkit/SDP-install.manifest.json`:

```powershell
$SdpSource = 'C:\path\to\extracted-or-cloned-SDP'
$Project = 'C:\path\to\Project'
```

Produce the deterministic, mutation-free JSON plan:

```powershell
& "$SdpSource\Toolkit\scripts\Install-SDP.ps1" `
  -ProjectRoot $Project `
  -PlanJson
```

Use `-Preview` for human-readable mutation-free output, or apply the plan:

```powershell
& "$SdpSource\Toolkit\scripts\Install-SDP.ps1" `
  -ProjectRoot $Project
```

Existing non-empty `SDP/` directories are the normal migration case. Managed
files refresh according to the installation contract and are backed up before
replacement. Project-owned files are created only when missing and are never
replaced, including under `-ForceManagedFiles`.

Install paths are portable, case-collision checked and physically confined to
link-free source and target roots. Plans are schema- and semantics-validated;
apply executes the exact planned source/target pairs. Equivalent supported YAML
does not cause churn, and AGENTS migration conflicts use deterministic
content-hash destinations.

Installed `sourceCommit` provenance is null for archives without trustworthy Git
metadata. In a dirty clone, a non-null value identifies the available `HEAD`
baseline; it does not attest that installed bytes equal that commit.

`-InitializeProjectStructure` adds only missing neutral lifecycle and operating
seeds. It never copies this repository's active Sprint, `REL-0.2.0`, release
notes, Ledger history, review or verification evidence. See
`Toolkit/docs/Installation-Contract.md` and `Toolkit/docs/Installer-Migration.md`.

## Build identity

Generate framework-neutral build metadata:

```powershell
python Toolkit\scripts\New-SDPBuildIdentity.py `
  --project-root C:\path\to\Project `
  --output C:\path\to\Project\src\generated\sdp-build.json
```

A Vite app may import the JSON, while non-Vite systems may embed or expose the
same schema by another build step. Unreleased builds are visibly marked `-dev`.

## Canonical skills

- `sdp-master`, `sdp-worker`, `sdp-reviewer`, `sdp-architect`
- `sdp-traceability`, `sdp-vertical-refactor`
- `sdp-release`, `sdp-versioning`, `sdp-auditor`, `sdp-verifier`

Every skill has machine-readable YAML front matter and its version must agree
with `SDP.manifest.yaml` and the generated installed facts.

## Validation

```powershell
python -m pip install -r Toolkit\tests\requirements.txt
python Toolkit\scripts\validate_sdp.py --mode toolkit
python Toolkit\scripts\validate_sdp.py --mode project --project-root C:\path\to\Project
python -m unittest discover -s Toolkit\tests -p "test_*.py" -v
.\Toolkit\tests\Install-SDP.Tests.ps1
```

Toolkit mode is the backward-compatible default. Project mode validates an
installed consuming project without assuming the Toolkit repository layout.
The regression suite includes an offline fixture pinned to the supported
`gh-sdp` commit and normalized tree identity.
See `Toolkit/docs/Validation.md` for exact boundaries and limitations.

## Main documentation

- `SDP/History/legacy-bootstrap/SDP-DOCUMENT-GUIDE.md`
- `docs/process/How-SDP-Works.md`
- `Toolkit/docs/Installation-Contract.md`
- `Toolkit/docs/Installer-Migration.md`
- `Toolkit/docs/Validation.md`
- `docs/process/Feature-Governance-And-SDP-2.0.md`
- `Toolkit/docs/Distribution-And-Upgrades.md`
- `Toolkit/docs/Release-And-Versioning.md`
- `Toolkit/docs/Release-Lifecycle.md`
- `Toolkit/docs/SDP-Analyzer-Compatibility.md`

## Exploratory studies

The [design-language study](SDL/docs/studies/Design-Language-Definition.md)
preserves vocabulary, grammar and type proposals. The studies below supply
rationale and scenarios; use the [active-profile overview](SDL/docs/README.md)
for implemented language boundaries.

- [Features, Functionality, Containers and Channels](SDL/docs/studies/Feature-Functionality-and-Channel-Study.md)
  — owner-model clarification and comparison with established architecture and
  contract languages; research recommendations, not an adopted Toolkit contract.
- [MVP1 design evolution and the proposed SDP skills](docs/process/MVP1-Design-Evolution-and-SDP-Skills.md)
  — chronological design evidence, the owner's selected September 10 direction,
  inherited obligations, and practical application of the candidate skills.
- [Scenarios, state and implementation traceability](SDL/docs/studies/Scenarios-State-and-Implementation-Traceability.md)
  — capability realization, behavior/state terminology, and navigation from
  scenario flows to units, contracts and source-level change impact.
- [Vocabulary and grammar exploration](SDL/docs/studies/Vocabulary-and-Grammar-Exploration.md)
  — typed nouns and verbs, qualifiers, conditions, modality, temporal expressions
  and an initial structural grammar fragment.
- [MVP1 design-language example](SDL/docs/studies/MVP1-Design-Language-Example.md)
  — system overview, detailed UI responsibilities, scenario/state descriptions
  and inspected implementation bindings expressed through the proposed vocabulary.
- [Design-language conformance scenarios](SDL/docs/studies/Design-Language-Conformance-Scenarios.md)
  — scenario-driven vocabulary development, canonical sentence candidates,
  rejected ambiguities and the validation needed to prevent dialects.
- [SDL source tree and compilation](SDL/docs/studies/SDL-Source-Tree-and-Compilation-Study.md)
  — proposed multi-file authoring structure, name resolution, container boundary
  exports, traceable blueprints and staged language/compiler development.
- [Executable SDL IR and runtime](SDL/docs/studies/SDL-Executable-IR-and-Runtime-Study.md)
  — proposed execution profile, concurrency semantics, ports/adapters, external
  renderer contract and bounded experiments for running a design.
- [Whole-system MVP1 SDL exercise](experiments/mvp1_sdl/README.md)
  — multi-file system model, selected UI layers, 20 scenario graphs and a pinned
  inventory of 332 system/local obligations; executable gaps remain explicit.
- [MVP1 exercise findings and local language extensions](SDL/docs/studies/MVP1-SDL-Exercise-Findings-and-Extensions.md)
  — candidate constructs recorded before use, with modeling findings and limits;
  isolated from the general language definition and the current parser.
- [Datasets, Datagrams and data contracts](SDL/docs/studies/SDL-Datasets-Datagrams-and-Data-Contracts.md)
  — proposed datasets, optional databases and Dataset-linked Datagram families;
  contract-defined variants, source identity and the SDL-to-IR completion boundary.
- [ControlSets, layer boundaries and data access](SDL/docs/studies/SDL-ControlSets-Layer-Boundaries-and-Data-Access.md)
  — HEOS/HSX and external-model evidence; Commands/Values, generated MessageSets,
  internal Channels, queryable data versus live arrivals and explicit durability.
