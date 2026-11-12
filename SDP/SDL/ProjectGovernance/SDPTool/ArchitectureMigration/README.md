# SDPTool architecture migration model

Activity0001 / PLAN-SDP-0026. This concern is the source-owned cross-entrypoint
responsibility map for restructuring existing SDPTool. It complements existing
navigation, installation, ModelGovernance and SessionWorkflow detail models; it
neither replaces their behavior contracts nor claims complete code inference.

ObservedArchitecture WORK models inspected responsibilities at d0bf6dc. TARGET
moves Session adaptation/dispatch into GovernanceApplication, action attribution
into OperationRegistry and MCP framing into McpResponseAdapter. Existing actual
processes remain SdpToolHost, ControllerProcess and McpProcess. No new runtime
or Activity capability is introduced. Metadata/operation inventory and removal
rules live in the linked Activity contract; plain SDL ownership does not encode
Go imports, every public function or prove runtime behavior.

[Activity architecture contract](../../../../Activities/Activity-0001--SDPTool-Restructuring/04-Design/Architecture.md)
and [inventory](../../../../Activities/Activity-0001--SDPTool-Restructuring/02-Study/Inventory.md)
map authored units to actual code. Both WORKs were created/checkpointed through
ModelGovernance and validate with design-core/0.6, entry System.design (includes
Intent.design and Responsibilities.design). Neither is an accepted model release.
Use `sdl check WORK--TargetArchitecture/System.design` from this directory.
See the [Blueprint references](../../../../Activities/Activity-0001--SDPTool-Restructuring/04-Design/Blueprints.md).
