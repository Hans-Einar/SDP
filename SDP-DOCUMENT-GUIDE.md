# SDP Document Guide

The folders at this repository root demonstrate the recommended structure inside
a consuming project's `SDP/` directory. They are reference templates; populated
records in consuming projects remain project-owned.

## Lifecycle documents

| Folder | Purpose |
|---|---|
| `01--Mandate` | Purpose, outcomes, scope and non-goals |
| `02--Study` | Evidence, alternatives, assumptions and unknowns |
| `03--Requirements` | Stable functional, quality and constraint IDs |
| `04--Architecture` | Boundaries, ownership, dependencies and contracts |
| `05--DesignAnalysis` | Horizontal analysis and implementation fan-out |
| `06--Design` | Chosen detailed designs and rationale |
| `07--Implementation` | Ordered vertical delivery plan |

## Operating and release documents

| Area | Purpose |
|---|---|
| `Sprints/` or `Refactors/` | Planned delivery container |
| `Fixes/` | Strictly bounded smaller corrections |
| `Releases/` | First-class `REL-X.Y.Z` records |
| `CodeReview/` | Fresh independent review evidence |
| `Verification/` | Exact checks and artifacts |
| `Traceability/` | Current state, relations and append-only events |
| `Instructions/` | Project-specific operating rules |
| `SDP-project.manifest.yaml` | Project release and development coordinates |
| `RELEASE-NOTES.md` | Unreleased and immutable released history |

SemVer identifies public releases. Sprint/Refactor, Iteration, Slice/Fix and
revision remain separate development coordinates.

## Installation ownership

The installer refreshes only clearly Toolkit-managed `AGENTS.md`, Framework and
skill files. It creates project manifest, release notes and templates only when
missing. Existing project records are never replaced, including when
`-ForceManagedFiles` is used.
