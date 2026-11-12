# Restructuring obligations

[Mandate](../01-Mandate/Mandate.md) and
[PLAN-SDP-0026](../../../05--Implementation/SDPTool/ArchitectureMigration/Plan.md)
define this change's constraints and completion. Preserve existing public behavior
and exposure; document known defects separately from compatibility obligations.

The canonical [system requirements](../../../02--Requirements/SDPTool.md) retain
their source home. This page maps migration verification to those requirements;
it does not establish a competing full requirements baseline.

| Canonical requirement/contract | Preservation evidence |
| --- | --- |
| REQ-SDPTOOL-001 preview | Root preview/consumer tests and native producer integration |
| REQ-SDPTOOL-002 explicit project resolution | Root discovery/CLI routing tests and compiled baseline comparison |
| REQ-SDPTOOL-003/004 navigation, identity and freshness | Root navigation/source/deferred tests and native XFMD navigation, overview, deferred and source-link gates |
| REQ-SDPTOOL-007 installation ownership/recovery | install and nested bootstrap suites; gh-sdp forwarding suite |
| Existing governance and Session contracts (SDPTool/governance/README.md) | Governance recovery/authority suites and compiled CLI/MCP journeys |
| Existing model/Blueprint contracts (SDPTool/model/README.md; blueprints/README.md) | Model, Blueprint and assignment suites; frozen public catalogue |
| Activity0001 internal architecture constraint | Architecture negative specimens, typed dispatch/parameter/caller tests, MCP framing policy tests and independent review |

REQ-SDPTOOL-005/006 remain their previously declared later scope. No new plan/history
product feature is claimed. Candidate-specific results and limitations live in
[Verification.md](../06-Verification/Verification.md).
