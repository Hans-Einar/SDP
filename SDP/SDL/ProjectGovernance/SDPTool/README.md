# SDPTool — project facade and shared governance core

## Boundary and implementation state

SDPTool owns project discovery, command coordination, navigation and installation,
and is the proposed owner of routine decisions and durable run state. It is one
software system in ProjectGovernance. The `SdpToolHost` container represents its
command/service execution boundary; new routine logic is units/libraries within
this responsibility, not a new independently deployed workflow engine.

The existing Go command is `sdptool`, built from
[SDPTool/cmd/sdptool](../../../../SDPTool/cmd/sdptool/main.go). Its reusable Go
facade and [installation package](../../../../SDPTool/install/Records.md) exist.
The [existing architecture](../../../03--Architecture/SDPTool.design) remains the
authority for detailed navigation/preview/install functions, channels and scenarios.
This new entry owns the **proposed routine boundary overview**; its two coarse
`ExistingReadFacade` and `ExistingInstallFacade` activities summarize source
responsibilities, not a second detailed installation model or a claim of complete
SDPTool functionality. Shared Go libraries are not separately deployed containers.

| Model concern | Evidence and status |
| --- | --- |
| Registered project inventory | Implemented facade: [project.go](../../../../SDPTool/project.go), [navigation.go](../../../../SDPTool/navigation.go), [kanban.go](../../../../SDPTool/kanban.go) |
| Installation command coordination | Implemented facade: [installation.go](../../../../SDPTool/installation.go), [install/cli.go](../../../../SDPTool/install/cli.go); full guarantees/evidence remain in the existing design and installation plan |
| Routine matching and checked transitions | Proposed; no generic routine service in the studied candidate |
| Durable run/revision recovery | Proposed; existing installer journal is precedent, not this implementation |
| Checkout authority resolution | Proposed extension beyond current path/project discovery |

## Functions, contracts and state

`RoutineCoordinator` resolves applicability and checks transitions. `ContextFacade`
resolves project/checkout authority. `RunRepository` owns operational snapshots;
`OperationalRunStore` means persistent retrievable data, not a SQL database or a
new competing management ledger. Run state must refer to ProjectManagement and
Traceability records rather than become a second authority for either.

`RoutineRequest` identifies checkout, run, expected revision, action, operation and
evidence. `RoutineOutcome` reports revision, disposition, event cursor and coverage.
These are conceptual required fields for a future application service, not a
published wire schema. Text fields do not enforce revision/hash domains or authorize
callers. `TrustedAssignmentBinding` is supplied through a trusted launch/auth
context; a role string in the payload cannot confer owner/reviewer authority.

The accepted and refused example scenarios expose the two result families. Detailed
recovery, idempotency, stale evidence checks and multi-file publication must be
specified and tested before implementation is called guarded. Modes scope modeled
responsibilities; they are not a state transition engine. Route/gap/resume behavior
follows the [four-path synthesis](../../../02--Requirements/RoutineGovernance/Synthesis.md).

`ExternalCommandAdapter` is an unowned caller boundary standing for CLI/MCP/client
invocations, not a container bundled into SDPTool. SDL/SDUI own parsing, execution
and rendering. XFMD and gh-sdp are external consumers. The candidate application
API may be linked into adapters without turning every library into a service.

## Model and validation boundary

[System.design](System.design) is the independently parsed entry for **design-core
0.5**. Its declarations belong to this file only; directory names are catalog
metadata, not language namespaces or imports. Activities distinguish implemented
source responsibilities from planned delivery. A valid model or rendered scenario
is not runtime evidence, owner acceptance or authorization to implement a proposal.

From the repository root, with a built SDL CLI:

```sh
sdl check SDP/SDL/ProjectGovernance/SDPTool/System.design
sdl ast SDP/SDL/ProjectGovernance/SDPTool/System.design
sdl viewpoints SDP/SDL/ProjectGovernance/SDPTool/System.design --format static --output /tmp/sdp-sdptool-views
```

Build the CLI from `SDL/go` with `go build -o /desired/bin/sdl ./cmd/sdl`.
The authoring verification used the repository CLI with Go 1.27.1; the enclosing
[Ecosystems plan](../../../03--Architecture/Ecosystems/Plan.md) records consolidated
candidate hashes and export evidence. Generated output belongs in temporary or
explicitly marked derived directories, never in this authored source folder.
The [ecosystem index](../README.md) identifies cross-system dependencies.
Current discovery derives navigation from supported source files; no manual
navigation registry is required. Viewer support remains a separate consumer contract.


## Blueprint preview delivery — BPI2-M1

BlueprintFacade captures ModelGovernance source views and delegates structural
analysis/rendering to SDL/go/blueprint, then uses the shared document publisher.
GenerateBlueprintPreview is a diagnostic operation; BlueprintPreviewDelivery does
not claim catalogue discovery, assignment states or code conformance. The libraries
are dependencies of SdpToolHost, not new deployed containers. Source:
SDPTool/blueprints/generate.go and SDL/go/blueprint. Evidence:
05--Implementation/SDPTool/Blueprints/Evidence-BPI2-M1.md.
