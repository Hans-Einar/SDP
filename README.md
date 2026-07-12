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
├── 01--Mandate/ ... 07--Implementation/   lifecycle templates
├── Sprints/ Refactors/ Fixes/ Releases/   operating records/templates
├── CodeReview/ Verification/ Traceability/
├── SDP.manifest.yaml                      Toolkit release manifest
├── RELEASE-NOTES.md                       canonical Toolkit release notes
├── Toolkit/
│   ├── skills/                            versioned reusable skills
│   ├── schemas/                           machine-readable contracts
│   ├── scripts/                           installer, build identity, validation
│   └── payload/                           Toolkit-managed project files
├── docs/                                  method and release guidance
└── examples/                              example manifests and events
```

The root lifecycle folders demonstrate the recommended project-local `SDP/`
structure. In consuming projects, populated lifecycle/release/traceability files
are project-owned and authoritative.

## Version model

Released software and SDP Toolkits use SemVer. The Toolkit that existed before
formal version metadata is documented as migration baseline `0.1.0`; this
backward-compatible capability addition targets `0.2.0`. It remains
**unreleased** and no corresponding tag or GitHub Release exists.

Sprint/Refactor, Iteration, Slice/Fix, revision and Git SHA are separate
development coordinates. See `docs/Release-And-Versioning.md` and
`docs/Development-Identity.md`.

## Install or update a project

Keep one independent clone:

```powershell
git clone https://github.com/Hans-Einar/SDP.git C:\Users\hanse\GIT\SDP
```

Preview an installation or migration:

```powershell
C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\path\to\Project `
  -Preview
```

Apply it:

```powershell
C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\path\to\Project
```

Existing non-empty `SDP/` directories are the normal migration case. The
installer creates missing project manifests and release notes, updates clearly
Toolkit-managed Framework/skills, backs up replaced managed files and never
replaces populated project-owned records. `AGENTS.md` is managed;
`AGENTS-project.md` is preserved.

Use `-ForceManagedFiles` only to restore same-version managed files. Use
`-InitializeProjectStructure` to add missing template documents. See
`docs/Installer-Migration.md`.

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
with `SDP.manifest.yaml`.

## Validation

```powershell
python -m pip install -r Toolkit\tests\requirements.txt
python Toolkit\scripts\validate_sdp.py
python -m unittest discover -s Toolkit\tests -p "test_*.py" -v
.\Toolkit\tests\Install-SDP.Tests.ps1
```

GitHub Actions runs contract tests on Linux and installer fixture tests on
Windows.

## Main documentation

- `SDP-DOCUMENT-GUIDE.md`
- `docs/How-SDP-Works.md`
- `docs/Distribution-And-Upgrades.md`
- `docs/Release-And-Versioning.md`
- `docs/Release-Lifecycle.md`
- `docs/SDP-Analyzer-Compatibility.md`
