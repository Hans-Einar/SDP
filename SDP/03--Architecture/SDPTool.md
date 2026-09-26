# SDPTool — architecture and SDL model

Feature: **SdpTool**. Design record: **DES-SDPTOOL-001**. Status: initial design.
[Requirements](../02--Requirements/SDPTool.md) · [Detailed design](../04--Design/SDPTool.md) ·
[Model](SDPTool.design) · [Implementation plan](../05--Implementation/SDPTool.md).

## Architecture and reuse

SdpTool is an offered feature; SdpToolProcess is its proposed command-line runtime
boundary. The feature is not a directory or container. Functionality belongs to
focused units: PreviewCoordinator, ProjectContext, NavigationInventory,
ViewerBridge, PlanCoordinator and CardHistory. They contribute to SdpTool and
its use cases; named interfaces describe collaboration, not completed wire schemas.

SdlProjectionPort delegates to existing Go parser/viewpoint/documents services.
Those services already reuse mmdr and SDL-specific Go SVG symbols. KanBanReadPort
and SduiDocumentPort retain their existing semantic owners. No new language parser,
renderer repository, mandatory daemon or startup compilation is required.

ViewerDeliveryPort is the producer/consumer boundary. Native XFMD tabs and direct
.design source/preview behavior belong to KB-XFMD-014 and KB-XFMD-015. SDPTool does
not own FOX widgets. Preview/configuration protocols remain separately specified. The installation
extension adds explicit Channels and scenarios, with no invented packet encoding.
It has one local namespace and imports no external design model.


## Model boundaries

This self-contained model includes A0/A1 goal/feature facts, A2/A3 ownership and
planned delivery activities for several viewpoints. It is stored here because
architecture integrates these facts; it is not a claim that every fact is A2/A3.
Current SDL has no cross-file imports. Do not maintain parallel model copies
in Requirements, Design and Implementation. Installation A4 exchanges are modeled; preview A4 details retain their own scope.

Activity implementation-status values in the SDL model distinguish delivered and
planned work; they are not proof of native XFMD acceptance. SDPTool interfaces describe required collaboration
boundaries, not implemented signatures. Native XFMD work remains external; no
functionality in this model owns its widget implementation.

## Installation ownership — GIP delivery, 2026-09-27

SDPTool is the common Go command entry point at root SDPTool/. It owns Go
installation planning, validation, migration, backup/recovery and installed facts.
gh-sdp is a thin distribution/invocation client; it must not duplicate migration
policy. No separate Toolkit Go runtime is selected. The legacy PowerShell engine
and reusable contracts remain in Toolkit until a verified Go migration replaces
it. The install package owns policy and execution; installation.go keeps the
compatible facts/journal discovery adapter. The separate stdlib bootstrap module
is shared by the engine and thin client; it has no project migration policy.

[REQ-SDPTOOL-007](../02--Requirements/SDPTool.md) records the target outcome.
The canonical SDL model now allocates installation work to InstallationCoordinator,
InstallationBaselineInspector, InstallationReleaseResolver, InstallationPlanner,
InstallationExecutor, InstallationJournal and InstallationRecorder within
SdpToolProcess. GhSdpProcess contains only GhSdpLauncher; ReleaseRepositoryProcess
is an external artifact service. The five installation/client/adoption activities now carry implemented status
based on GIP evidence. IPD-2 adds twelve logical channels, five validated scenarios and file-based data
holders. [Installation contract](../04--Design/SDPTool/Installation/Contract.md)
and [scenario review](../04--Design/SDPTool/Installation/Scenarios.md) define the
design beyond SDL payload shapes. Runtime verification is recorded separately in
[PLAN-SDP-0003 evidence](../05--Implementation/SDPTool/Installation/Evidence.md).

[PLAN-SDP-0002](../04--Design/SDPTool/Installation/Plan.md) owns this design delivery;
[IPD evidence](../04--Design/SDPTool/Installation/Evidence.md) records parser checks.
[KB-SDP-033](../KanBan/active/%23033--Study--XFMD-SDP-adoption-and-SDL-pilot.md)
retains the subsequent implementation/adoption workflow.

Review the [generated installation viewpoints](../04--Design/SDPTool/Installation/review/index.md)
and [planned Go implementation](../05--Implementation/SDPTool/Installation/Plan.md).
