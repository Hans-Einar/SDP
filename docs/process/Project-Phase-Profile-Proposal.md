# Project phases, templates and per-phase studies — review proposal

**Owner review update:** KB-SDP-001 is accepted and closed. This document is the
accepted planning deliverable; the original proposal/review wording below records
its R3 submission. It is not an installed implementation contract. Outstanding
metadata choices are handled with KB-SDP-017; repository-linking choices remain
explicit prerequisites before extraction. KB-SDP-010 review can proceed separately.

**Local adoption, 2026-09-25:** the owner explicitly directs first real use of
these five numbered phases in SDP-vNow/SDP, starting with SDPTool. The local
[project entry](../../SDP/README.md) is initialized accordingly. Installer/template
migration remains separate; local adoption does not select discovery metadata.

Original submission: **proposed for owner review**, R3-M1, 2026-09-24. Primary card:
[KB-SDP-001](../../SDP/KanBan/completed/%23001--Proposal--Project-structure.md).
This is the single proposed replacement for the numbered directory profile.
The current [installation manifest](../../Toolkit/SDP-install.manifest.json)
continues to distribute 01–07; this document does not change installed destinations.

## Decisions already supplied by the owner

Each product has its own SDP area, including SDP itself and eventual SDL/SDUI
repositories. Template contains reusable sources, not live project records.
Keep numbered phases, relate them to abstraction levels and use studies where
needed within each phase. The Project Owner supplies the mandate and need not be
a developer. Mandate and any study together identify stakeholders. Actors may
represent people or systems. Requirement narratives should become SDL when its
profile supports them. Monorepo/subrepository layout should not alter tool behavior.

## Proposed profile for approval

| Directory | Responsibility and typical files | SDL abstraction relationship |
| --- | --- | --- |
| 01--Mandate | Mandate.md; optional Mandate-study.md; stakeholder identification | Process input and authorization, not a new A-level |
| 02--Requirements | Requirements-study.md when needed; 01--Actors.design, 02--UserStories.design, 03--UseCases.design; requirement/acceptance records | A0 needs/obligations and A1 functional intent |
| 03--Architecture | Architecture.design and optional Architecture-study.md; containers, Units, ownership, Channels, allocation | A2 system/containers and A3 internal Units/layers |
| 04--Design | Contract/scenario/data/class models, design decisions and optional Design-study.md | A4 detailed behavior/contracts |
| 05--Implementation | Approved implementation plan/roadmap and realization references; source code remains in its language/project directories | A5 realization/evidence; vertical slices can cross every earlier level |

These are process homes, not five new language abstraction levels. A0–A5 keep
[their existing meaning](../../SDL/docs/integration/SDL-Viewpoint-Levels-and-Notation.md).
A document can contain facts viewed at several levels. A phase number does not
assign a semantic level to each object, and phases need not be completed once in
strict waterfall order. Changes revisit affected earlier decisions explicitly.

Unnumbered Agents/KanBan, Traceability, Verification, CodeReview, Releases,
Sprints, Refactors, Fixes and Instructions keep their distinct roles. Maintenance
holds repository maintenance plans/evidence; History holds dated records. They
are not additional abstraction levels. Do not use a new Maintenance note as a
substitute for a primary backlog card.

## Studies and authority

A study is optional where the input is already clear. It records questions,
observations/sources, considered alternatives, proposed decisions and unresolved
items. A study's interpretation does not override an approved mandate, requirement
or profile. The responsible owner approves decisions separately; material changes
link back to the relevant card and affected design. Design analysis remains a
practice within the appropriate phase, not a mandatory separate numbered directory.

Use plain English Markdown for stakeholder descriptions, actor kind and stories
until the required SDL syntax is adopted. Current actor/usecase/feature facts may
already use design-core 0.5. Do not create apparently executable .design files with
unsupported stakeholder, userstory, comment, links, include or shortcut syntax.
[KB-SDL-001](../../SDP/KanBan/backlog/%23001--SDL--Proposal--Requirements-narrative.md)
owns that language work; existing candidates remain examples, not alternate syntax.
The proposed Actors/UserStories filenames become active templates only with a
supported profile and parser-tested examples.

## Migration from the currently distributed profile

| Current template directory | Proposed destination | Handling |
| --- | --- | --- |
| 01--Mandate | 01--Mandate | Keep mandate ownership; add optional study guidance |
| 02--Study | Relevant phase's study document, or History | Classify each document by content; never blindly merge unlike studies |
| 03--Requirements | 02--Requirements | Preserve IDs, acceptance and trace links |
| 04--Architecture | 03--Architecture | Preserve models, sources and ownership |
| 05--DesignAnalysis | Relevant phase study/decision, commonly 04--Design | Inventory mixed content and split deliberately with provenance |
| 06--Design | 04--Design | Detect name collisions; no silent overwrite |
| 07--Implementation | 05--Implementation | Preserve approved plans, release/slice/evidence references |

R1 already separated neutral templates from live records. Apply this future
migration to template definitions first, then generate a per-project migration
plan. Project-owned files remain protected even under force. Update manifest
sources/destinations, schema/capability/version contracts, documentation, agents
and tests together; do not reinterpret an installed 1.x tree just because new
templates use different numbering. Capture hashes and old→new paths. Update local
links/entry points; retain historical event paths and source fingerprints. A
collision or ambiguous study owner requires an explicit disposition, not deletion.

Acceptance of the later implementation: dry-run plan, neutral new installation,
upgrade with modified project documents preserved, stable IDs/ledgers, working
links, parser-tested supported examples and no second active template copy.
This proposal does not assert those installer changes/tests already exist.

## Review gate and follow-up

Requested owner decision: approve this five-directory profile, per-phase study
policy and migration approach, or identify changes. Approval authorizes preparing
a separately bounded installer/template implementation; it is not evidence that
such implementation has happened. Existing distributed 01–07 stays authoritative
until that versioned change passes its acceptance tests.

[Project discovery and extraction proposal](Project-Discovery-and-Extraction-Proposal.md)
completes KB-SDP-001's remaining design package. The proposed next implementation
planning item is [KB-SDP-017](../../SDP/KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md),
limited first to discovery and view-ip integration, rather than all slice synthesis.

## Local management placement update — 2026-09-25

PM1 adopts one operational SDP/KanBan board and SDP/ProjectManagement ledger for
SDP, SDL and SDUI, with Scrum reviews under SDP/Agents/Scrum. This supersedes
separate language-board placement for this project. It does not change frozen
consumer versions or installer/template destinations. The
[current profile](../../SDP/ProjectManagement/README.md) owns that local contract.
