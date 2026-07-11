# Standard Document Procedure (SDP)

Status: working draft

SDP is a repository-local, document-driven method for AI-assisted software
development. It keeps mandate, research, requirements, architecture, design,
implementation, verification, review, traceability and handoff connected to the
repository rather than relying on chat memory.

Core principle:

> Design horizontally. Implement vertically.

## Repository layout

This repository is the upstream SDP toolkit. It is intentionally not itself a
project-specific `SDP/` tree.

```text
SDP repository
├── docs/                     Method documentation
├── skills/                   Reusable agent skills
├── payload/                  Files safe to install into projects
│   ├── project-root/         Files placed at project repository root
│   └── sdp-root/             Files placed inside the project's SDP folder
├── examples/                 Example project-specific SDP tree
├── scripts/                  Safe installation/update scripts
├── DraftStandardDocumentProcedure.md
└── TieredDesignAndImplementation.md
```

Project-specific Mandate, Study, Requirements, Architecture, Design, Sprints,
Refactors, Traceability and Verification remain in each project. They are not
stored as active folders at the root of this toolkit.

## Recommended installation

Clone this repository once beside your projects, not inside a consuming
project's existing `SDP/` folder:

```powershell
git clone https://github.com/Hans-Einar/SDP.git C:\Users\hanse\GIT\SDP-Toolkit
```

Then install or update the reusable files into a project:

```powershell
C:\Users\hanse\GIT\SDP-Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\Users\hanse\GIT\GrassPhenology
```

This is required to work when the target project already contains an `SDP/`
directory. The installer creates missing managed subdirectories and files, but
preserves existing project-specific content.

The installer is intentionally additive by default:

- an existing project `SDP/` folder is supported
- existing project Mandate, Study, Requirements, Architecture, Design, Sprints,
  Refactors, Verification, CodeReview and Traceability content is untouched
- existing `AGENTS.md` is preserved
- existing `SDP/AGENT-REMINDERS.md` is preserved
- toolkit-managed skills are installed under `.codex/skills/`
- toolkit-managed framework guidance is installed under `SDP/Framework/`

Use `-ForceManagedFiles` only to refresh files explicitly defined as
Toolkit-managed. It still does not overwrite project-owned SDP documents.

## Important clone behavior

`git clone <url> <destination>` requires the destination not to exist, or to be
an empty directory. Therefore this will fail when the project already has a
non-empty `SDP/` folder:

```powershell
# Do not do this
git clone https://github.com/Hans-Einar/SDP.git C:\path\to\Project\SDP
```

Even if the target were empty, cloning there would create a nested repository
with its own `.git` directory. That can be intentional with a Git submodule, but
it is not SDP's default distribution model.

The supported model is:

```text
C:\Users\hanse\GIT\
├── SDP-Toolkit\            independent upstream clone
├── GrassPhenology\
│   └── SDP\                project-owned SDP documents
├── TerrainAnalyzer\
│   └── SDP\
└── agro-crm\
    └── SDP\
```

Use one upstream clone plus the installer. For a one-time copy without Git
metadata, use `git archive`, a release ZIP, or remove the copied `.git` folder.

## Skills

The first skills are deliberately small:

- `sdp-master`
- `sdp-worker`
- `sdp-reviewer`
- `sdp-architect`
- `sdp-traceability`
- `sdp-vertical-refactor`

They contain reusable procedure, while project facts remain in the project's
own SDP documents.

## Existing draft documents

- `DraftStandardDocumentProcedure.md`
- `TieredDesignAndImplementation.md`
- `docs/How-SDP-Works.md`
- `docs/Distribution-And-Upgrades.md`
- `examples/README.md`