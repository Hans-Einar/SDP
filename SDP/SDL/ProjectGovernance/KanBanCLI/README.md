# KanBanCLI — read-only shell listing

## Boundary and implementation state

This system is the existing [kanban.sh](../../../../Toolkit/scripts/cli/kanban.sh)
read-only listing tool. It runs as a Bash process, installed as `kanban`; there is
no separate Go binary or public Go library for this system today. Catalog placement
does not classify this maintained CLI as an obsolete installer or authorize its
replacement. A future compiled successor must preserve or explicitly revise the
observed command contract. All `ShellListingDelivery` functions describe existing
source responsibilities, not a newly executed acceptance suite.

## Actual command contract

```sh
kanban status /path/to/SDP/KanBan
kanban state /path/to/SDP/KanBan/active --group-by sprint
kanban status --group-by scrum --scrum SCRUM-SDP-0001
```

The optional directory defaults to the current directory. A board covers itself
and its known status children; a status directory covers itself only. It does
not discover arbitrary parent projects. Grouping supports state/sprint/scrum and
filters accept typed Sprint/Scrum IDs. Only the first visible metadata table is
read. Cards without CardState are ignored; malformed/duplicate state or grouping
metadata is an error. No matching cards prints `no KanBan cards found` and exits 1;
argument/metadata errors exit 2. This tool does not validate the complete card
ledger or installed process profile.

`TerminalOutput` emits printable sorted paths. Its current terminal mode uses
OSC 8 file links, while a pipe or dumb terminal uses plain text. This describes
the existing product behavior, not a request to print control sequences in agent
responses. The model intentionally uses interfaces rather than invented network
messages for shell/filesystem/stdout operations.

The tool does not move cards, change CardState, append management history, run a
Sprint or approve evidence. It is distinct from [SDPTool's validated board
service](../SDPTool/README.md) and the proposed [KanBanTUI](../KanBanTUI/README.md).
Do not silently infer the stronger SDPTool ledger guarantees from this lightweight
script. KanBan files and ProjectManagement history remain project-owned records.

## Model and validation boundary

[System.design](System.design) is the independently parsed entry for **design-core
0.5**. Its declarations belong to this file only; directory names are catalog
metadata, not language namespaces or imports. Activities distinguish implemented
source responsibilities from planned delivery. A valid model or rendered scenario
is not runtime evidence, owner acceptance or authorization to implement a proposal.

From the repository root, with a built SDL CLI:

```sh
sdl check SDP/SDL/ProjectGovernance/KanBanCLI/System.design
sdl ast SDP/SDL/ProjectGovernance/KanBanCLI/System.design
sdl viewpoints SDP/SDL/ProjectGovernance/KanBanCLI/System.design --format static --output /tmp/sdp-kanbancli-views
```

Build the CLI from `SDL/go` with `go build -o /desired/bin/sdl ./cmd/sdl`.
The authoring verification used the repository CLI with Go 1.27.1; the enclosing
[Ecosystems plan](../../../03--Architecture/Ecosystems/Plan.md) records consolidated
candidate hashes and export evidence. Generated output belongs in temporary or
explicitly marked derived directories, never in this authored source folder.
The [ecosystem index](../README.md) identifies cross-system dependencies.
Navigation registration must select this entry explicitly; folder placement alone
does not make it available to a viewer.
