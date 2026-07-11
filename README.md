# Standard Document Procedure (SDP)

Status: working draft

SDP is a repository-local, document-driven method for AI-assisted software
development. It keeps mandate, research, requirements, architecture, design,
implementation, verification, review, traceability and handoff connected to the
repository rather than relying on chat memory.

Core principle:

> Design horizontally. Implement vertically.

## Repository layout

This repository has two purposes:

1. show the recommended project-local `SDP/` structure directly at repository root
2. provide reusable skills and installation tooling under `Toolkit/`

```text
SDP repository
├── 01--Mandate/               Template and reference document
├── 02--Study/
├── 03--Requirements/
├── 04--Architecture/
├── 05--DesignAnalysis/
├── 06--Design/
├── 07--Implementation/
├── Sprints/
├── Refactors/
├── CodeReview/
├── Verification/
├── Traceability/
├── Instructions/
├── SDP-DOCUMENT-GUIDE.md      Detailed use of each document area
├── Toolkit/
│   ├── skills/                Reusable Codex skills
│   ├── scripts/               Safe installer/update scripts
│   └── payload/               Toolkit-managed project files
├── docs/                      Extended method documentation
└── examples/                  Larger examples
```

The root lifecycle folders are templates and reference material. In consuming
projects, the corresponding folders inside that project's `SDP/` directory are
project-owned and authoritative.

## Recommended installation

Clone this repository once at the existing shared location:

```powershell
git clone https://github.com/Hans-Einar/SDP.git C:\Users\hanse\GIT\SDP
```

If `C:\Users\hanse\GIT\SDP` already contains this clone, use `git pull` instead.

Install reusable Toolkit files into a project:

```powershell
C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\Users\hanse\GIT\GrassPhenology
```

The target project may already contain a non-empty `SDP/` directory. This is a
hard requirement and the normal migration case. Existing project documents are
preserved.

To also add only missing standard template documents and folders:

```powershell
C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\Users\hanse\GIT\GrassPhenology `
  -InitializeProjectStructure
```

## Managed and project-owned agent instructions

`AGENTS.md` is Toolkit-managed and is refreshed on every install. Consuming
projects should not edit it.

Repository-specific instructions belong in `AGENTS-project.md`, which is
project-owned and is never overwritten by the installer. The managed
`AGENTS.md` requires agents to read `AGENTS-project.md` when it exists.

When migrating an existing project:

- an old, differing `AGENTS.md` is copied to `AGENTS-project.md` when that file
  does not yet exist
- when `AGENTS-project.md` already exists, the old `AGENTS.md` is kept as a
  timestamped migration backup
- the canonical managed `AGENTS.md` is then installed

The installer otherwise remains additive:

- existing project Mandate, Study, Requirements, Architecture, Design, Sprints,
  Refactors, Verification, CodeReview, Instructions and Traceability are untouched
- existing `AGENTS-project.md` is preserved
- existing `SDP/AGENT-REMINDERS.md` is preserved
- Toolkit-managed skills are installed under `.codex/skills/`
- Toolkit-managed Framework guidance is installed under `SDP/Framework/`
- standard project templates are copied only when the destination is missing

Use `-ForceManagedFiles` to refresh Toolkit-managed Framework and skill files.
It does not replace project-owned SDP documents or `AGENTS-project.md`.

## Clone behavior

Do not clone this repository into a consuming project's existing `SDP/` folder:

```powershell
# Wrong
git clone https://github.com/Hans-Einar/SDP.git C:\path\to\Project\SDP
```

`git clone` requires the destination to be absent or empty and creates its own
`.git` directory. Cloning there would either fail against an existing non-empty
folder or create a nested Git repository.

Supported layout:

```text
C:\Users\hanse\GIT\
├── SDP\                     independent upstream clone
├── GrassPhenology\
│   └── SDP\                 project-owned documents
├── TerrainAnalyzer\
│   └── SDP\
└── agro-crm\
    └── SDP\
```

## Skills

Canonical skills live under `Toolkit/skills/` and are installed to
`.codex/skills/` in each consuming project:

- `sdp-master`
- `sdp-worker`
- `sdp-reviewer`
- `sdp-architect`
- `sdp-traceability`
- `sdp-vertical-refactor`

Project facts remain in project-specific SDP documents; skills contain reusable
procedure.

## Main documentation

- `SDP-DOCUMENT-GUIDE.md`
- `DraftStandardDocumentProcedure.md`
- `TieredDesignAndImplementation.md`
- `docs/How-SDP-Works.md`
- `docs/Distribution-And-Upgrades.md`
- `examples/README.md`