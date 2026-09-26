# Project discovery and SDL/SDUI extraction — review proposal

**Owner review update:** KB-SDP-001 is accepted and closed. This document is the
accepted planning deliverable; the original proposal/review wording below records
its R3 submission. It is not an installed implementation contract. Outstanding
metadata choices are handled with KB-SDP-002/016; repository-linking choices remain
explicit prerequisites before extraction. KB-SDP-010 review can proceed separately.

Status: **proposed**, R3-M1, 2026-09-24. Completes the design package for
[KB-SDP-001](../../SDP/KanBan/completed/%23001--Proposal--Project-structure.md).
No sdptool implementation or repository extraction is delivered here.

## Actual baseline

The Toolkit's project manifest schema 1.0 owns installed release/work state;
[its contract](../../Toolkit/docs/Project-Manifest.md) remains authoritative.
The local SDP, SDL/SDP and SDUI/SDP process areas currently contain no installed
SDP-project.manifest.yaml. Their existence is not proof of a conforming installed
project. The [board registry](../../SDP/ProjectManagement/History/SDP/boards.json) registers boards,
not complete projects. Avoid silently treating it as sdptool discovery.

SDL/go still declares module github.com/Hans-Einar/SDP/SystemDesignLanguage/go;
SDUI/go declares github.com/Hans-Einar/SDP/SDUI/go. SDL has a relative replace to
../../SDUI/go. The common design source is SDUI/design/architecture.design.
These are concrete extraction dependencies, not reasons to duplicate parsers/models.

## Proposed recognition and registry

Keep project identity separate from Git roots and from installed release identity.
Propose an explicit SDP/project.json discovery marker; its schema/version are a
new contract to implement, not an addition to the existing release manifest.
A marker records stable project ID, repository root relative to its SDP area,
process profile, selected model entry and registered child-project locations.
Keep viewer executable registration in host configuration, not arbitrary Markdown
links. Existing release manifests continue to own release/work coordinates.

Illustrative marker fields (unimplemented proposal, not an installed config):

```json
{
  "schemaVersion": "proposal/1",
  "projectId": "SDP",
  "repositoryRoot": "..",
  "processProfile": "proposed-phases-1",
  "model": "../SDUI/design/architecture.design",
  "projects": {
    "SDL": "../SDL/SDP",
    "SDUI": "../SDUI/SDP"
  }
}
```

Marker paths are relative to the marker directory, independent of cwd. Resolve
and validate registered paths explicitly; no recursive scan of the home directory,
implicit Git discovery, network cloning or automatic selection of a child project.
Duplicate IDs/conflicting registrations fail. Shared models must have one declared
owner and explicit cross-project references. A missing sibling repository produces
a useful diagnostic; it never falls back to stale output or a different project.

## Command resolution

Accept `sdptool [PROJECT_PATH] COMMAND ...`, including generate ip and view ip.
Without PROJECT_PATH use cwd. Recognition order:

1. If the selected directory has a valid supported project marker, it is the SDP area.
2. Otherwise inspect only its SDP child for a supported marker.
3. An invalid/unknown marker is an error, not permission to select another project.
4. Neither found: explain the missing marker; do not search parent directories.

A conforming installed release manifest without a new marker is recognized only
through an explicitly versioned legacy adapter; it must report missing model/viewer
configuration. A random directory named SDP is not sufficient. Migration can create
markers only after project identity/configuration is explicit, including the three
current local process areas. This step is absent today and belongs to implementation.

| Selected directory / condition | Expected result |
| --- | --- |
| Repository root with valid SDP/project.json | That project's SDP area |
| The same SDP directory | Exactly the same project/model/plan |
| Explicit SDL directory or SDL/SDP | SDL, whether monorepo or nested Git repository |
| A deeper source directory | Not found; no implicit parent search |
| Unknown/invalid marker plus a different nested marker | Error, not silent fallback |
| Duplicate project ID or unavailable registered dependency | Explicit diagnostic |
| Paths with spaces; cwd unrelated to source | Same resolved model and viewer target |

These are acceptance cases, not tests of a delivered resolver.

## First bounded sdptool delivery

After review, separate KB-SDP-002 into tested milestones: marker/resolver with the
matrix above; explicit project/model/viewer configuration; view ip using the existing
Go document projector and registered viewer; checked source/provenance/error handling.
Reuse prebuilt tools, on-demand selection and explicit XFMD window/panel delivery.
Do not compile during browsing, shell-execute arbitrary URI commands or require a
daemon. Keep approved plans immutable unless an explicit proposal is accepted.

Automatic slice proposals, validation of user-proposed slices, roadmap overlays
and Git-backed card history are later bounded deliveries. Their evidence/status
rules depend on KB-SDP-004; weak links and source assertions are not implementation
proof. Do not claim generate ip is implemented by rebranding sdl-design.

## Repository extraction plan

| Milestone | Deliverable / prerequisite |
| --- | --- |
| E1 — inventory and owner decision | Choose repository names/hosting, license/provenance, public module paths and Git link mechanism. Map source, tests, docs, board IDs/ledgers and shared design ownership. |
| E2 — prepare boundaries | Remove build dependence on parent checkout through explicit module/workspace configuration; preserve a monorepo development path. Move no code without passing parser/runtime/renderer tests. |
| E3 — extract with provenance | Preserve history or document the exact source commit/file map; keep each project's SDP/Agents/KanBan and stable IDs. Update cross-project registry, Refs and source links together. |
| E4 — verify consumers | Test standalone SDL and SDUI plus SDP integration, code generation, source/view launch, examples and installed package paths; no hidden sibling checkout or duplicate parser. |
| E5 — switch entry points | Update supported build/distribution paths and remove replaced active paths only after both layouts pass. Preserve historical records and branch review boundaries. |

Recommendation for the eventual umbrella checkout: pinned Git submodules, because
SDL/SDUI are intended to have independent repositories/releases. Subtree is an
alternative if ordinary clones must contain all code and independent publication
is secondary. This recommendation is **not adopted**. Tool resolution is marker/path
based under either mechanism, and `.git` must never be a required project boundary.
No repositories, submodules or remote publications are created by this proposal.

## Review gate

Approve or revise the discovery marker/resolution contract and extraction direction.
Repository names, hosting and actual Git integration must be chosen before E3.
The phase/template proposal can be approved independently; neither decision adds
SDL keywords. [Numbered phase proposal](Project-Phase-Profile-Proposal.md).
