# Design-to-implementation handoff

[RuntimeDistribution WORK](../../../SDL/ProjectGovernance/SDPTool/RuntimeDistribution/WORK--SharedRuntime/System.design)
checkpoint00002 is the validated preliminary bounded model (checkpoint00001 is
preserved as the original handoff). It is not an accepted
whole-system release. The ImplementationContract records concrete path, helper,
public operation and delivery contracts.

| Modeled responsibility | Canonical implementation | Evidence boundary |
| --- | --- | --- |
| RuntimeRegistration | SDPTool/runtime_setup.go; action_runtime.go; runtime_cli.go | same handler for CLI/action, durable target, explicit conflicts and shell diagnostics |
| SharedSduiFacade | SDUI/go/preview; SDPTool/action_runtime.go | original preview/check contract reused; stale/source/readiness tests |
| LocalHelperMaterializer | SDPTool/native_host.go; native_payload_bundled.go; package.sh | private content-addressed extraction, corruption/symlink refusal, actual GUI launch |
| ExternalGh | gh-sdp/cmd/package and scripts/package.sh | verified canonical bytes, no runtime bootstrap; actual install/upgrade/offline/recovery |
| ExternalXfmd | application/sdp and application/SduiTools in XFMD | registered path, capability handshake, native reading/Activity/SDUI probes |

GitHub CLI has no post-install execution hook. One explicit setup is part of the
selected user contract; read-only commands never register themselves. The
source-based helper commands remain thin development adapters, not another release
channel. Canonical packaged engine carries the desktop helper in its own signed
bytes. Release v2.3.0 uses a static engine and an AlmaLinux10.2-native GUI payload.
