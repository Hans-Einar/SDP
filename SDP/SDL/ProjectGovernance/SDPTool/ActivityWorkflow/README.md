# SDPTool Activity command concern

Owner Session0016 T3, 2026-11-12, selects SDL design for Activity and Session
creation through SDPTool. This is a bounded design concern of SDPTool/SdpToolHost,
not a separate deployed system. The sibling System, ModelGovernance and
SessionWorkflow models retain their existing authority and history.

WORK--ActivityCommands is a new preliminary ModelGovernance WORK, created with
--initial because this concern has no accepted release. It is not a fork claiming
the entire existing SessionWorkflow TARGET is implemented. Existing SessionLifecycle
WORK :00003 was inspected clean and remains unchanged. No generated viewpoints or Blueprint were produced in that initial design-only
assignment. Activity0002 now supplies the separate implementation handoff below.

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

## Activity0002 implementation handoff

PLAN-SDP-0027 selects ActivityCommands checkpoint
b3e45cab-c629-490b-9192-3ed2b4628387:00002 against the bounded absence baseline
73a3ef6f-6869-4b83-a6d8-993d2a8c4bda:00001. The baseline represents the previously
absent Activity services, not the absence of the existing SDPTool system.
Blueprint revision 830b6836603b55b150e04d6f9018f636ce8332c144e7b7d01ee91349725b25ec
passes structural constraints; runtime and code conformance remain unknown until
implementation evidence is recorded. The Activity's 04-Design/Blueprints.md and
ImplementationContract.md bind the worker scope and exact artifact locations.
No accepted model release is implied.
