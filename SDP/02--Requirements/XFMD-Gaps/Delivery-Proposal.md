# Proposed delivery after the XFMD gap study

Input: [Study](Study.md), [Evidence](Evidence.md), PLAN-SDP-0009.
This is a **proposal for selecting successor plans**, not an active ImplementationPlan
or permission to modify XFMD. Phase IDs below belong to this proposal only.

## Owner scope clarification — 2026-09-29

SDUI primarily describes the concept, composition and placement of a UI using a
limited widget vocabulary. Running a basic SDUI and calling SDL remains part of
its purpose. It need not reproduce every function of a rich native application.
A tree is a possible useful addition, not an approved requirement to rebuild
XFMD in SDUI. Full editors, sophisticated dialogs, application-wide command
systems and host parity remain optional proposals requiring demonstrated need.
This clarification governs the study's successor selection; the gap inventory
records consumer observations, not a mandatory product backlog.

## Recommended selection

Select KB-SDL-005 next for its bounded System/source-set contract. Independently,
KB-SDUI-004 is a small concrete fidelity fix with no new grammar dependency.
Select the first design phase of KB-SDUI-003 before new interactive widgets.
Do not wait for full XFMD UI parity to try SDL on a real project.

| Phase / milestone | Bounded result and owner | Exit evidence / dependencies |
| --- | --- | --- |
| XGP1 / M1 | SDL Frontend: define versioned System/input membership and resolution under KB-SDL-005 | Positive multi-file system; missing/duplicate members, unresolved reference, source spans, limits and revision invalidation. No ad hoc concatenator. |
| XGP1 / M2 | SDL Viewpoints/DocumentService and SDPTool: consume one coherent source revision | Same model facts in check, navigation tree and selected Markdown generation; source change invalidates relevant output; independent consumer check. KB-SDP-020 migrates existing sources only after support. MVP1's experimental profile needs separate semantic decisions beyond file assembly. |
| XGP2 / M1 | SDUIPresentation: text fidelity under KB-SDUI-004 | Shared measured/painted glyph policy, word-aware wrapping and nested list preservation; probe fixture reproduces before/fixed after; inspected exports at narrow/normal/wide sizes. Can run independently of XGP1. |
| XGP3 / M1 | SDUIFrontend/Runtime/NativeHost, SDL bridge: capability and activation DesignPlan under KB-SDUI-003 | One explicit matrix separates parse/profile, layout, static provider and interactive host support. Typed event identities and complete preflight defined; missing capability/module/signature cannot partially activate. Current bindings remain covered. |
| XGP3 / M2 | Same systems: one collection + viewport vertical implementation | Small navigation tree with stable item identities, expand/loading/error/selection and scroll. Keyboard/pointer behavior, rejected/stale event and disposed generation verified in the selected native host, not only AST/SVG. Depends on M1, not on all SDL source-set implementation. |
| XGP4 / M1 | SDUIRuntime/NativeHost: selected panes and shared commands | Tabs/splits, command state and context cancellation share identity/focus rules; tests cover hidden vs removed pages, reload, resize and stale context target. Select only after pilot feedback. |
| XGP4 / M2 | SDUIRuntime/NativeHost/application adapter: typed values and draft forms | Boolean/number/choice semantics reuse runtime drafts and atomic batches; ranges, cancel and persistence failure proved with a controlled application fixture. No SDUI-owned business transaction. |
| XGP5 / M1 | SDUIFrontend/Presentation/NativeHost: actual shared components and editor/preview boundary design | Source-set/import decision backed by two real screens; instance identity and diagnostics. Typed host-surface/provider capabilities with revision, resource disposal and unsupported/fallback policy. Deferred until demanded by a consumer. |
| XGP5 / M2 | Selected provider/host: bounded rich-content acceptance | Individually verified diagram/image/math/text families, hit testing and anchors as selected; measured/painted agreement, lifecycle and explicit fallback. No claim of “full Markdown” from a single picture. |
| XGP6 / M1 | SDL requirements/behavior/evidence work, KB-SDL-001/006 and KB-SDP-004 | Map one real workflow and four requirement identity states; distinguish model assertion, registered action and actual product evidence. May proceed alongside earlier phases; it is not a prerequisite for every UI fix. |

Historical recommendation, superseded for the widget inventory by the
2026-10-07 refinement and execution request in KB-SDUI-003 / Session0010:
Only XGP1, XGP2 and XGP3-M1 were recommended near-term selections. XGP4–5
are optional extension candidates, not required SDUI completion scope. XGP6
retains separately owned SDL research; it is not a prerequisite for basic UI
prototyping. The tree pilot is a candidate, not a selected product commitment.
Choose BranchPolicy/CommitPolicy in each actual plan: current working branch and
milestone commits normally suffice for a small fix; retain KB-SDL-005's explicit
phase-branch commitment. No merge or release authorization is implied.

## Candidate interactive pilot

Use a small provider-supplied navigation tree with group rows, two lazy branches
and a bounded list longer than the viewport. Activation selects a stable model
item and requests a view; expansion alone does not open it. Show pending/error
state and a read-only result preview. A fixture provider suffices for SDUI host
acceptance; real SDL/SDPTool generation is the next consumer integration gate.

| Case | Required result |
| --- | --- |
| Expand / select / activate | Expansion is separate from activation; command receives stable item and model revision. Keyboard/pointer agree. |
| Loading, error and retry | Per-request identity prevents an older load from overwriting a newer subtree; error/retry is visible and bounded. |
| Scroll and resize | Needed content remains reachable, offsets clamp to new extent, hit testing follows clipping. |
| Deleted item / reload | Old selection/event cannot invoke a replacement item accidentally; compatible named state retained and removed handles revoked. |
| Missing capability or binding | Failure precedes activation with source/capability diagnostics; no half-installed handlers or working-looking inert controls. |
| Cancel / dispose | Outstanding work cannot update a disposed UI; owned resources released. Domain cancellation remains explicit and may be cooperative. |
| Unsupported static export | Explicit rejection or deliberately requested labeled preview; never claim the picture is an interactive tree. |

Native XFMD navigation integration is independently owned consumer work and can
start against current supported single-file models. It need not wait for this
SDUI pilot or use SDUI for its sidebar widgets. Its
admission tests must include dirty generated-navigation rejection, Save followed
by failed open, invalid URI while another generation is active, stale completion,
cancel and lease release. Reuse its current navigation service; an SDUI study
does not authorize a FOX bridge or replacing native XFMD widgets.

## Follow-up ownership and traceability

- [KB-SDL-005](../../KanBan/backlog/%23005--SDL--Change--System-and-source-sets.md): SDL-001.
- [KB-SDL-001](../../KanBan/backlog/%23001--SDL--Proposal--Requirements-narrative.md): SDL-002 narrative/identities.
- [KB-SDL-006](../../KanBan/backlog/%23006--SDL--Study--Executable-channel-tests-and-unit-bindings.md): SDL-003 behavior composition.
- [KB-SDP-004](../../KanBan/backlog/%23004--Proposal--Design-traceability.md): SDL-002/003 source/evidence relationships.
- [KB-SDUI-003](../../KanBan/completed/%23003--SDUI--Proposal--Capabilities-and-navigation-pilot.md): SDUI-001–010, staged rather than ten separate cards.
- [KB-SDUI-004](../../KanBan/backlog/%23004--SDUI--Bug--Text-and-Markdown-fidelity.md): SDUI-011.
- External XFMD register and MAINT-XFMD-0004: DOC-001–004 and EVID-001–002; no local resolution guarantee for external IDs.

At study completion the backlog remained unselected. KB-SDUI-003 is now selected
under [PLAN-SDP-0021](../../04--Design/SDUI/Widgets/Plan.md), with all listed
panes, commands, typed controls and basic multiline input required. Only component
source sets, full editors and broader rich-content research remain optional.
Completing PLAN-SDP-0009 means these study
results and follow-ups exist; it does not close the product gaps. Actual future
design/code evidence belongs in system-prefixed Traceability linked to its
selected plan. Management-only study events remain in ProjectManagement.
