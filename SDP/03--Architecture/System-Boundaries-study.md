# System boundaries and shared process ownership — study

Date: 2026-09-25. Status: owner-selected system boundaries; SDL syntax implementation remains pending. Primary tracking:
[KB-SDL-004](../../SDL/SDP/Agents/KanBan/active/%23004--Change--Language-source-organization.md).

## Owner clarification

Maintaining separate full SDP process trees for SDL and SDUI appears unnecessary.
The owner prefers bringing their development documentation into root SDP and asks
how languages, tools, systems and containers should relate. This replaces the
working assumption that the migration should distribute the shared model across
two independently maintained process trees. The existing three boards and their
append-only histories have not yet been moved or combined.

## Adopted system boundaries

Use one project process area in root SDP for the work we coordinate together.
A process project may document one or more software systems; it is not itself a
C4 software-system boundary. Keep the five numbered process homes and organize
artifacts inside them by responsibility. Preserve SDL, SDUI and Toolkit/SDPTool
as implementation source locations.

The owner selects **SDL**, **SDUI** and **SDPTool** as three separate software
systems within this shared process project. Each system can contain multiple
executable tools/containers and library components. XFMD is a fourth collaborating
system with its own SDP area in its own repository. No enclosing software system
called SDP Development Tools is introduced.

This supersedes the earlier agent recommendation to begin with one integrated
system boundary. It is an explicit product/design decision, not an inference
from repository or directory names. Cross-system contracts and library reuse
remain visible even when the three systems share a Git checkout and process area.

| Concept | Meaning in this project |
| --- | --- |
| SDP process/project | Mandate, requirements, decisions, plans and evidence; may cover several systems |
| SDP ecosystem / landscape | A selected view of cooperating systems; no mandatory new runtime parent |
| Software system | A coherent user-facing purpose and responsibility boundary |
| Container | An executable application/service or data store within that boundary |
| Component / logical Unit | Responsibility implemented within a container; reusable libraries may serve several containers |
| Source module/folder | Code or authoring organization; does not establish runtime containment |

Candidate containers include the sdl CLI, sdui CLI, sdui-fyne interactive host,
optional sdl-viewsd broker and planned sdptool CLI. They must be reconciled with
existing abstract CommandLineHost/FyneHost/ViewServiceHost symbols before changes;
this table is not a new validated allocation. Additional Go commands are assessed
by runtime responsibility, not mechanically converted into containers. Parser,
runtime and layout libraries retain a single library design even where multiple
applications link them. Their code can execute in more than one container.

XFMD is an independently useful document-viewer system collaborating with this
toolchain. Its preview/navigation contracts cross that system boundary. It is
currently represented by XfmdDocumentHost in the existing model; changing its type
requires a deliberate model/profile migration. The owner authorizes initializing its SDP area and migrating KanBan now.
Full XFMD adoption/design is assigned to another agent through KB-XFMD-017;
no XFMD product code changes are part of the bootstrap.

## C4 interpretation and sources

C4 distinguishes software systems, containers, components and code. A container
is a runtime application or data-store boundary; an ordinary library/module is
usually not a container. Therefore a language/tool family can group several
containers, but that group should not be modeled as a container containing other
containers in a strict C4 decomposition.

- [Software systems](https://c4model.com/abstractions/software-system)
- [Containers and libraries](https://c4model.com/abstractions/container)
- [System landscapes](https://c4model.com/diagrams/system-landscape)

The system landscape supplies the ecosystem view. It need not introduce an
Ecosystem declaration or make the SDP methodology an enclosing software system.

## SDL language implications

System is already proposed in the
[source-tree study](../../SDL/docs/studies/SDL-Source-Tree-and-Compilation-Study.md),
which describes one explicitly declared System per compilation. Current
**design-core 0.5 does not implement System**. Implement it only through an
explicit profile change with ownership, cross-system references, validation,
source provenance and context/landscape projections defined together.

A shared process workspace with several system entries is compatible with one
system per compilation. A future landscape may reference their stable identities;
it must not merge private namespaces or imply System-contains-System semantics.
Keep process phase, model abstraction, system membership and file placement as
separate properties. A0–A5 are viewpoint concerns rather than folder-derived types.

The new uncommitted sourceinput package is only an unfinished input-loader draft.
It has not been integrated, tested, installed or claimed as source-set support.
The source split, tool integration and profile decisions remain outstanding.

## Selected authoring organization

Use root SDP/01--Mandate through 05--Implementation. Within Requirements, group
SDL, SDUI and SDPTool needs where useful. Architecture owns the landscape/system
context, container topology and connections. Detailed Design contains a directory
per actual container, plus SharedLibraries/SDL, SharedLibraries/SDUI and Contracts
for facts that should not be duplicated in each host. Implementation has one
coordinated plan with explicit product responsibilities and existing evidence.
Group architecture and detailed design by SDL, SDUI and SDPTool within the
shared numbered phases; detailed design then groups by actual container and
shared library/contract responsibilities. Do not create a separate SDP process
tree under each system or equate the folder name with a language declaration.

## Next bounded delivery

1. Apply the owner's three-system decision and shared root process ownership.
   The current monolithic model and its type assignments still need migration.
2. Update the migration plan/card ownership for the common process area. Preserve
   stable card IDs, historical ledger bytes and reference paths when consolidating
   boards; do not silently renumber SDL/SDUI history into SDP.
3. Finish explicit multi-file input/provenance independently of new System syntax,
   then split the shared model into the agreed root phase homes without losing facts.
4. Treat adoption of System and any changed runtime boundaries as a separate,
   explicitly tested semantic migration rather than part of a lossless file move.

No documentation relocation, active model replacement, new System syntax or
native rendering integration is delivered by this study.

## XFMD process bootstrap and handoff

The owner requested the separate XFMD SDP bootstrap in the same conversation.
Its board now lives at SDP/Agents/KanBan in the xfmd-sdl-navigation worktree.
KB-XFMD-016 records bootstrap evidence; KB-XFMD-017 assigns subsequent adoption
and valid .design models to another agent. This does not implement System in SDL.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/README.md

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Agents/KanBan/backlog/%23017--Proposal--Adopt-SDP-and-model-XFMD.md
