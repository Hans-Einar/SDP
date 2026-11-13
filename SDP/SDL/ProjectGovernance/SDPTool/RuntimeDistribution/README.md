# Shared runtime distribution concern

Activity4 owns this bounded model of SDPTool distribution/registration and its
external consumers. The canonical engine is the gh extension executable itself;
GitHub CLI owns replacement, SDPTool owns explicit local registration and managed
SDUI helper dispatch, XFMD consumes the registered engine offline. Libraries are
units in SdpToolHost, not separately deployed services. ExternalGh and ExternalXfmd
are collaborating callers, not bundled containers.

WORK--SharedRuntime/System.design uses design-core0.6. Validate and checkpoint
with ModelGovernance; this preliminary WORK is an implementation handoff, not an
accepted release or runtime evidence. The Activity4 ImplementationContract defines
path ownership, conflicts, signatures, exact byte identity and release ordering.
Existing system/navigation/governance concerns retain their authority.
