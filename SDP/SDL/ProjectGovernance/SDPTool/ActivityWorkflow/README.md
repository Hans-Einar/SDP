# SDPTool Activity command concern

Owner Session0016 T3, 2026-11-12, selects SDL design for Activity and Session
creation through SDPTool. This is a bounded design concern of SDPTool/SdpToolHost,
not a separate deployed system. The sibling System, ModelGovernance and
SessionWorkflow models retain their existing authority and history.

WORK--ActivityCommands is a new preliminary ModelGovernance WORK, created with
--initial because this concern has no accepted release. It is not a fork claiming
the entire existing SessionWorkflow TARGET is implemented. Existing SessionLifecycle
WORK :00003 was inspected clean and remains unchanged. No generated viewpoints
or Blueprint are produced in this assignment.

System.design composes Features, Responsibilities, Contracts and Scenarios.
The model owns new Activity application responsibilities, canonical record creation,
Activity-scoped Session document creation, state checks, navigation/timeline
projections and the boundary to existing PG admission. Creation does not dispatch
a run. Management Activity records are modeled as data; ActivityCommandsDelivery
is an SDL delivery activity marked planned, not a runtime Activity instance.

[Command design](../../../../04--Design/SDPTool/ProjectGovernance/Activities/Commands.md)
owns proposed CLI/actions, payload constraints, code placement and error behavior.
[Overall design](../../../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md)
owns lifecycle, document layout, timeline and adoption rules. SDL shape/scenario
validation does not prove authorization, concurrency, atomic recovery or runtime
implementation. Message text fields are bounded structural placeholders, not
the complete future JSON schema.

Use built repository tools for source-only verification:

```sh
sdl check SDP/SDL/ProjectGovernance/SDPTool/ActivityWorkflow/WORK--ActivityCommands/System.design
sdptool SDP/SDL/ProjectGovernance/SDPTool/ActivityWorkflow model status work:ActivityCommands --json
```

Model metadata/checkpoints are created by ModelGovernance, never edited manually.
Views remain on-demand in XFMD. No CANDIDATE, accepted RELEASE or product delivery
is implied by a clean WORK checkpoint.
