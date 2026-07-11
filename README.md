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

Clone this repository once beside your projects:

```powershell
git clone https://github.com/Hans-Einar/SDP.git C:\Users\hanse\GIT\SDP
```

Then install or update the reusable files into a project:

```powershell
C:\Users\hanse\GIT\SDP\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\Users\hanse\GIT\GrassPhenology
```

The installer copies only toolkit-managed files. By default it does not replace
an existing project-owned file. Use `-ForceManagedFiles` only for files that are
explicitly documented as toolkit-managed.

## Why not clone into an existing project SDP folder?

`git clone` creates a new working tree containing its own `.git` directory.
Cloning directly into `Project/SDP` would therefore create a nested repository.
That can be intentional when using a Git submodule, but it is not the default
SDP distribution model.

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
