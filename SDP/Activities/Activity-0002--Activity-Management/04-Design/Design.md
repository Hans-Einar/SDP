# Selected design and implementation handoff

[Activity design](../../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md)
and [commands](../../../04--Design/SDPTool/ProjectGovernance/Activities/Commands.md)
define selected behavior. [ActivityWorkflow SDL](../../../SDL/ProjectGovernance/SDPTool/ActivityWorkflow/README.md)
is its authoritative model home. Activity0001's architecture governs dependencies.

A0 reconciles exact wire schemas, descriptor history, transactional recovery,
in-flight guard integration, admission/adoption and consumer action contracts.
Keep current 0.1 manual records readable; add a versioned extension where new
payload semantics are required. Do not mutate checkpoint metadata by hand.
