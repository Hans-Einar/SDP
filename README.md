# Standard Document Procedure (SDP)

Status: working draft  
Toolkit-Version: 0.2.0 (unreleased)

SDP is a repository-local, document-driven method for AI-assisted software
development. It connects mandate, research, requirements, architecture, design,
implementation, verification, review, traceability, releases and handoff to the
repository instead of relying on chat memory.

Core principle:

> Design horizontally. Implement vertically.

## Repository layout

```text
SDP repository
|-- 01--Mandate/ ... 07--Implementation/   this repository's SDP records
|-- Sprints/ Refactors/ Fixes/ Releases/   live operating and release records
|-- CodeReview/ Verification/ Traceability/
|-- SDP.manifest.yaml                      Toolkit release manifest
|-- RELEASE-NOTES.md                       canonical Toolkit release notes
|-- Toolkit/
|   |-- SDP-install.manifest.json          canonical installation contract
|   |-- project-templates/                 neutral project-owned seeds
|   |-- payload/                           Toolkit-managed copied files
|   |-- skills/                            versioned reusable skills
|   |-- schemas/                           machine-readable contracts
|   |-- scripts/                           installer, build identity, validation
|   `-- tests/                             deterministic contract fixtures
|-- docs/                                  method and compatibility guidance
`-- examples/                              neutral contract examples
```

The root lifecycle and traceability folders are this repository's live SDP
instance. They are not installation templates. A conforming installer copies
only entries explicitly listed in `Toolkit/SDP-install.manifest.json`; neutral
project seeds come only from `Toolkit/project-templates/`.

## Version model

Released software and SDP Toolkits use SemVer. The Toolkit that existed before
formal version metadata is documented as migration baseline `0.1.0`; this
backward-compatible capability addition targets `0.2.0`. It remains
**unreleased** and no corresponding tag or GitHub Release exists.

Sprint/Refactor, Iteration, Slice/Fix, revision and Git SHA are separate
development coordinates. See `docs/Release-And-Versioning.md` and
`docs/Development-Identity.md`.

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
`docs/Installation-Contract.md` and `docs/Installer-Migration.md`.

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
See `docs/Validation.md` for exact boundaries and limitations.

## Main documentation

- `SDP-DOCUMENT-GUIDE.md`
- `docs/How-SDP-Works.md`
- `docs/Installation-Contract.md`
- `docs/Installer-Migration.md`
- `docs/Validation.md`
- `docs/Feature-Governance-And-SDP-2.0.md`
- `docs/Distribution-And-Upgrades.md`
- `docs/Release-And-Versioning.md`
- `docs/Release-Lifecycle.md`
- `docs/SDP-Analyzer-Compatibility.md`
