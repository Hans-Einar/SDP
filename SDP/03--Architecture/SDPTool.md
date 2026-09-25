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
not own FOX widgets. The model intentionally does not invent Channels, scenarios
or packet contracts before the preview/configuration interfaces are specified.
It has one local namespace and imports no external design model.


## Model boundaries

This self-contained model includes A0/A1 goal/feature facts, A2/A3 ownership and
planned delivery activities for several viewpoints. It is stored here because
architecture integrates these facts; it is not a claim that every fact is A2/A3.
Current SDL has no cross-file imports. Do not maintain parallel model copies
in Requirements, Design and Implementation. Detailed A4 protocols remain open.

Every Activity remains planned. SDPTool interfaces describe required collaboration
boundaries, not implemented signatures. Native XFMD work remains external; no
functionality in this model owns its widget implementation.
