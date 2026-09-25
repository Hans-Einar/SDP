# SDPTool — common project tooling

This is the owner-selected home for `sdptool`, the project-aware entry point for
SDP operations and selected SDL/SDUI services. Created on 2026-09-25. The directory
and [implementation plan](../../SDP/05--Implementation/SDPTool.md) are delivered; an executable,
discovery schema and navigator API are not implemented here yet.

The active feature assignment is
[KB-SDP-017](../../SDP/Agents/KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md).
The feature is designed in the [shared SDP process](../../SDP/README.md), using
its five numbered phases and common Traceability/KanBan. Do not create a nested
SDP directory here. Source code stays here; the single phase/milestone plan lives
in SDP/05--Implementation.

SDPTool owns project discovery/configuration, command coordination, navigation
inventory and viewer adapters. SDL and SDUI retain their parsers, runtimes and
projection/export implementations. Prefer reuse of existing APIs/commands.

Intended usage:

```sh
sdptool ~/git/XFMD generate ip
sdptool generate ip
sdptool view ip
```

An optional first path selects the project; otherwise use the current directory.
Recognition checks the selected SDP area, then its SDP child. Define validity
against the reconciled manifest contract, not directory existence alone. Do not
infer a parent project. Keep monorepo and separately checked-out projects equivalent.

For the existing tools, URI scheme and temporary-file lifecycle, read
[Current navigation and design preview](Navigation-and-Design-Preview.md).

## Navigation responsibility

The service must describe all supported groupable concepts and relationships as
tree nodes, not just UseCases. ViewPoint/phase groups, typed collections, individual
objects and supported relationship selections come from validated model facts.
Use stable identities and explicit selectable targets. Repeated relationships
need references/lazy expansion with cycle handling rather than pretending the
model itself is a tree. Produce documents on selection, without pre-rendering the
entire model. Exact wire format and command names remain to be specified.

The native XFMD SDP tab, its KanBan / SDL / SDUI subtabs, widgets and interaction
belong to [KB-XFMD-014](../../../xfmd-sdl-navigation/SDP/Agents/KanBan/backlog/%23014--Proposal--SDP-sidebar-and-generated-navigation.md)
in XFMD's own `Agents/KanBan`. No XFMD implementation is done from this workstream.
XFMD does not adopt SDP's project-development process or gain an SDP directory.
SDPTool provides the reusable producer services; XFMD supplies a native consumer.

The broader [Toolkit audit](../../SDP/Agents/KanBan/backlog/%23018--Study--Toolkit-audit-and-organization.md)
is separate. Creating this directory changes neither installation inventory nor
published Toolkit versions. Existing installers, schemas and language tools keep
their current responsibilities until an explicit migration is delivered.
