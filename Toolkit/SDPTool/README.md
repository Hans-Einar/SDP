# SDPTool — common project tooling

This is the owner-selected home for `sdptool`, the project-aware entry point for
SDP operations and selected SDL/SDUI services. Created on 2026-09-25. The directory
and [implementation plan](../../SDP/05--Implementation/SDPTool.md) are delivered.
The Go executable supports saved design preview, project discovery and configured
plan viewing, typed navigation data and SDUI structural Markdown previews.

The active feature assignment is
[KB-SDP-017](../../SDP/KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md).
The feature is designed in the [shared SDP process](../../SDP/README.md), using
its five numbered phases and common Traceability/KanBan. Do not create a nested
SDP directory here. Source code stays here; the single phase/milestone plan lives
in SDP/05--Implementation.

SDPTool owns project discovery/configuration, command coordination, navigation
inventory and viewer adapters. SDL and SDUI retain their parsers, runtimes and
projection/export implementations. Prefer reuse of existing APIs/commands.

Current project operations:

```sh
sdptool /path/to/project discover
sdptool tree
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

The broader [Toolkit audit](../../SDP/KanBan/backlog/%23018--Study--Toolkit-audit-and-organization.md)
is separate. Creating this directory changes neither installation inventory nor
published Toolkit versions. Existing installers, schemas and language tools keep
their current responsibilities until an explicit migration is delivered.

## Implemented saved-file preview — P0-M1

The Go facade now implements the [producer contract](Contract.md). Build it once
from this directory, then use the prebuilt executable for document requests:

```sh
go build -o /desired/bin/sdptool ./cmd/sdptool
sdptool preview ../../SDP/03--Architecture/SDPTool.design --output /tmp/design-preview
```

The local module replacements locate existing SDL/SDUI libraries in this checkout;
no Python parser or second language implementation is introduced. Tests run with
`go test -race ./...`. Discovery and view ip use the navigation registration described in the contract;
direct preview works without project metadata.

## Project selection and viewer

```sh
sdptool /path/to/project discover
sdptool /path/to/project view ip --model sdptool --viewer /path/to/xfmd --sdl-tool /path/to/sdl
```

The repository registers two SDL model entries and explicitly defaults to sdptool;
select --model sdl-sdui for the shared language design.
The viewer command waits until that window exits to keep generated navigation
resources alive. Host options can also use SDP_XFMD, SDP_SDL_TOOL and SDP_MMDR.
`generate ip` remains later scope and never runs implicitly while viewing.

## Navigation services

```sh
sdptool /path/to/project tree --model sdptool
sdptool /path/to/project select --model sdptool --uri 'sdl-view://sdp-vnow/VP02?diagram=VP02-roots&target=main&consumer=xfmd' --revision SOURCE_HASH --output /tmp/selected-view
sdptool /path/to/project sdui-preview --model concept1 --entry page --output /tmp/ui-document
```

Use the revision/targets returned by tree, not a hardcoded hash. KanBan nodes carry
current file paths, CardState and optional Scrum/Sprint grouping; SDL tree nodes
cover every catalog viewpoint and supported kind. SDUI preview currently delegates
the structural Markdown exporter; it does not imply interactive controls.

[Consumer examples and executable harness](Consumer-Examples.md) document the
producer boundary available to XFMD and other hosts.
