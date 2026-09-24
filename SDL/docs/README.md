# SDL — documentation map and authority

SDL code now lives in SDL/go. R1 moved documentation by responsibility; this is not a new language version. The [Go README](../go/README.md) and active parser profiles define executable scope. [SDL KanBan](../SDP/Agents/KanBan/README.md) contains proposals that are not yet language rules.

## Active profiles and boundaries

Design-core 0.5, action-core 0.1 and class-core 0.1 are separate bounded profiles.

- [Structural core](studies/Design-Language-Definition.md): sections 1–11 and registered additions; section 12 onwards contains open work/candidates.
- [Data/wire contracts](profiles/SDL-Data-Contract-Profile.md).
- [Channels and scenarios](profiles/SDL-Channel-Scenario-Profile.md).
- [Delivery-plan facts](profiles/SDL-Delivery-Plan-Profile.md).
- [Executable actions](profiles/SDL-Executable-Action-Profile.md) and [grammar](../grammar/action-core-0.1.ebnf).
- [Classes](profiles/SDL-Class-Profile.md) and [grammar](../grammar/class-core-0.1.ebnf).
- [Symbol profile](profiles/SDL-Symbol-Profile.md).

Profiles distinguish active rules from dated V2–V4 delivery evidence. Old test counts and planned statuses apply to the named milestones. Current [Go commands](../go/README.md) and the [implementation plan](../../SDUI/docs/implementation-plan.md) lead to later G-phase evidence; a model assertion alone is not verification.

## Integration, studies and history

The [Go port inventory](integration/SDL-Go-Port-Inventory.md), [navigation design](integration/SDL-Navigable-Viewpoints-Design.md) and [levels/notation](integration/SDL-Viewpoint-Levels-and-Notation.md) record reuse and original design foundations. G6 is implemented within bounded profiles; dated design prose cannot independently extend them. [G6 evidence](../go/evidence/G6.md) and the [launcher](../scripts/README.md) identify delivered scope.

[Design Language Definition](studies/Design-Language-Definition.md) contains both active core and explicitly bounded proposals; the studies/ directory alone does not determine authority. Other studies explore concepts beyond active profiles. [Research](../research/README.md) compares other design languages. The [MVP1 corpus](../../experiments/mvp1_sdl/README.md) is a candidate exercise. [Checkpoint #1](../../SDP/History/checkpoint-1/README.md) is a shared dated SDP/SDL/SDUI snapshot, not a unified active SDL specification.

The [document map](../../SDP/Maintenance/R1/Documentation-index.md) locates all former docs files.
