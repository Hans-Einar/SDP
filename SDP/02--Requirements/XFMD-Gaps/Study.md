# XFMD-driven SDL and SDUI gaps — study result

Study: [PLAN-SDP-0009](StudyPlan.md), 2026-09-29. Status: research delivered;
the recommendations below are **not adopted language rules or authorized implementation**.
[Evidence](Evidence.md) identifies the candidate, reproducible probes and limits.
[Delivery proposal](Delivery-Proposal.md) turns the findings into bounded successor work.

## Owner scope clarification — 2026-09-29

SDUI primarily describes the concept, composition and placement of a UI using a
limited widget vocabulary. Running a basic SDUI and calling SDL remains part of
its purpose. It need not reproduce every function of a rich native application.
A tree is a possible useful addition, not an approved requirement to rebuild
XFMD in SDUI. Full editors, sophisticated dialogs, application-wide command
systems and host parity remain optional proposals requiring demonstrated need.
This clarification governs the study's successor selection; the gap inventory
records consumer observations, not a mandatory product backlog.

## Recommendation

Do not make complete XFMD UI parity a prerequisite for using SDL. Continue its
supported structural models and existing navigation service while delivering
SDL System/source sets through KB-SDL-005. For SDUI, correct text fidelity and
define host capability checks before adding one interactive collection/viewport
workflow. Defer a full editor and rich preview to explicit host/provider contracts.

The register contains real limitations, but they do not all require more syntax:

- SDL already has structural, action, class and channel/scenario facilities.
  It lacks their complete composition into the requested whole-system model and
  evidence mapping. An executable Go action is not proof of XFMD C++ behavior.
- The SDL–SDUI bridge already validates typed binding plans before installing
  handlers. Extend this path; do not create a competing binding engine.
- SDUI already has drafts, atomic property updates, stale-event rejection and
  reload identity handling. Missing controls must use those mechanisms.
- Scroll syntax already parses. Its layout/host implementation is missing.
- Missing glyphs, character-based wrapping and flattened lists are presentation
  defects/limitations, not reasons to grow the grammar.

The source register is external and remains open there. This producer study does
not close XFMD acceptance or rewrite its gap dispositions.

## Evidence and interpretation

The study built tracked Go sources at `a274265` in a temporary archive, excluding
the unrelated untracked `SDL/go/sourceinput` draft. All 21 external capability
probes reproduced their expected outcomes. Five additional probes distinguished
local component reuse, unsupported import syntax, rejected image/HTML content,
and a successful SVG containing a diagram placeholder rather than a diagram.

After restoring missing tracked fixtures to the archive, ten selected Go packages
passed: SDL parser/runtime/bridge/reload and SDUI parser/runtime/reload/layout/
Markdown/codegen. The configured external Mermaid renderer test was skipped.
These checks are bounded regression evidence, not native-host or XFMD acceptance.

Two SVG exports were rasterized and visually inspected at widths 220 and 420.
They reproduce missing-symbol boxes and flat nested lists; the narrower export
splits “ordinary” into “or” and “dinary”. Direct probes identify glyph index zero
for ↻, ↺, ✓ and 界 without a rendering error; → and ø have real glyphs. The
presentation model stores no list depth. These findings apply to this embedded
font/exporter, not every font or XFMD's separate text rendering path.

## SDL dispositions

| Source ID | Current support and missing responsibility | Recommendation and owner | Required proof |
| --- | --- | --- | --- |
| GAP-XFMD-SDL-001 | design-core 0.5 accepts Container; System/include are rejected. One checked file is not a composed system. Directory organization cannot supply language identity. | Propose the existing KB-SDL-005 System/source-set contract, owned by SDL Frontend and its tool consumers. Keep parsing I/O-free; explicit input resolution belongs outside it. | Declared membership, cross-file references, duplicate/missing identities, deterministic results, original spans, bounded path/cycle handling and revisions covering all inputs; run through navigation as well as checker. |
| GAP-XFMD-SDL-002 | Actors, UseCases, Features and Functionality do not carry all normative external requirement text/status/amendments. A proposed Requirement declaration is rejected; external hyphenated IDs need explicit mapping. | Clarify requirements/narrative through KB-SDL-001 and evidence relations through KB-SDP-004. Retain external identity plus source/candidate references; do not derive compliance from prose or renumber requirements. | Current, amended, future and reserved identities survive round-trip and projection. Narrative, design satisfaction, implementation and evidence status remain distinct. |
| GAP-XFMD-SDL-003 | Structural Mode is not a state-machine state. Typed actions call registered Go handlers; channel/class/scenario facilities do not automatically bind to real XFMD functions, ordering or resource lifetime. | Map each invariant to an existing profile before extending it. KB-SDL-006 owns executable channel/unit binding research; KB-SDP-004 owns source/evidence relationships. XFMD retains its actual state transitions and tests. | Navigation success, cancellation, stale generation, dirty rejection and lease release each name the responsible unit, trigger, expected outcome and evidence candidate. Model checks and real-code tests are reported separately. |

For SDL-003, the current action engine validates registration, inputs and results;
it serializes calls and consumes request sequence identities even when a handler
fails. Context cancellation is cooperative. Neither an error nor cancellation
rolls back arbitrary Go effects. The class profile describes relationships; it is
not an object runtime. Use structural channels for collaboration and action tests
for explicitly bound behavior until a selected composition contract connects them.

### Navigation invariant mapping

This is the requested decomposition of SDL-003, not a new state-machine grammar.

| Invariant | Expressible now | What remains outside that proof |
| --- | --- | --- |
| Successful selection produces the requested view | Structural units/channels and a scenario can describe participants/messages; an explicitly registered action can validate a request/result shape | No automatic binding from those declarations to XFMD's generation process or preview delivery; the real workflow needs candidate-pinned tests |
| Cancellation prevents obsolete delivery | A typed action may carry an explicit cancel/outcome value; the Go handler receives context | Context alone is cooperative, not an ordering or rollback proof. XFMD owns child cancellation and delivery admission; cancellation ordering needs its actual tests |
| Old result cannot replace the current generation | Action runtime sequence/revision checks and SDUI event/reload protections cover their own boundaries | XFMD generation IDs and its renderer/result lease are different identities. A source/evidence mapping must show the explicit translation and stale-result rejection |
| Dirty document rejects unsolicited generated navigation | Model the application responsibility and interaction structurally; a registered action can implement a deliberate admission predicate | No general SDL precondition language or structural Mode declaration proves this predicate in XFMD. DocumentViews checks dirty state; no native verification run here |
| Save followed by failed open preserves the correct state | Scenario can distinguish save success from later read failure and describe both outcomes | A universal failed-open rollback rule is false. DocumentCoordinator saves before reading; domain state/evidence must represent the already completed save separately |
| Invalid URI does not cancel a valid in-flight request | Model rejection as a distinct scenario before the cancellation interaction | Structural facts do not enforce execution ordering. Source shows validation before cancel; a regression test must establish the behavior |
| Resources are released after cancel/replacement/disposal | Describe ownership and release interaction; handler/host code can perform release | No automatic resource-lifetime checker follows from those edges. Lease identity, asynchronous completion and exactly-once/effective-idempotent release policy need an explicit application contract and tests |

Consequently, first add explicit model-to-code/test references under KB-SDP-004
and investigate executable contract composition under KB-SDL-006. Do not claim
every invariant requires a new keyword: several are application tests whose
identity, scope and outcome need to be attached to the design. The missing generic
facility is trustworthy composition/evidence binding, not necessarily a second
programming language inside SDL.

## SDUI dispositions

The IDs below preserve the external vocabulary without adopting its suggested
widget names. **KB-SDUI-003** owns staged capability/interaction follow-up for
001–010; **KB-SDUI-004** owns the independently reproducible text-fidelity work.

| Source ID | Existing support / gap | Disposition and minimal boundary | Acceptance focus |
| --- | --- | --- | --- |
| GAP-XFMD-SDUI-001 | tree/list spellings rejected; frames/Markdown can only illustrate collections. | Propose a bounded collection model: stable item IDs, selection, expansion and explicit loading/error states. Host projects rows; application supplies items, sorting and operations. Begin with bounded data and explicit limits, not an unproven virtualization promise. | Loading/error/retry, deleted and stale items, keyboard/pointer activation, nonselectable group labels; no open caused by expansion alone. |
| GAP-XFMD-SDUI-002 | Relative frames exist; tabs/splitter spellings rejected and no corresponding interaction state. | Defer until collection pilot. Add stable page identity, selected page and bounded relative split proportions. Distinguish hidden state from removed identity; do not introduce pixel dimensions. | Resize/collapse/restore and reload retain compatible proportions, focus and selection; unavailable page has an explicit fallback. |
| GAP-XFMD-SDUI-003 | overflow-y=scroll passes AST; layout returns unsupported-scroll. Fyne's clipping container currently disables scrolling. | Implement a negotiated viewport capability in layout/runtime/host; no grammar change is needed merely to recognize scroll. Logical extent/offset belongs at the host boundary; source-anchor synchronization belongs to XFMD. | Long content, resize/clamping, keyboard/wheel, clipping/hit tests and feedback-loop suppression. Static export explicitly states its viewport policy. |
| GAP-XFMD-SDUI-004 | menu/dialog and onContext rejected; no transient-surface contract. | Propose host-owned transient surfaces and shared commands after the pilot. Target identity, modality, invocation source and explicit cancel/dismiss outcomes are necessary. | Dismiss does not activate, change history or commit drafts; stale target refused; focus restored predictably. |
| GAP-XFMD-SDUI-005 | input has drafts; textarea/editor semantics absent. | Defer a full built-in editor. Prefer a typed host-surface extension for existing document/preview implementations. Revision, selection, edit intents and disposal must be bounded; no unrestricted native pointer or opaque callback escape hatch. | Unicode/IME, undo, multiline edits, stale preview, anchors and dirty admission verified in the actual host/application. |
| GAP-XFMD-SDUI-006 | slider/checkbox/select rejected; runtime already supplies draft/revert and atomic updates. | Propose typed number/boolean/choice state after collection acceptance. Application owns form transaction/persistence; SDUI owns visible draft/value and feedback. | Range/step validation, pointer/keyboard equivalence, update atomicity, dirty conflict, cancel and failed persistence; programmatic updates emit no user action. |
| GAP-XFMD-SDUI-007 | icon/checked rejected; basic button state is insufficient for shared commands. | Propose stable semantic commands with enabled/checked state, accessible name and invocation identity; host maps shortcuts and resource presentation. Avoid FOX/Fyne constants in source. | Same command across toolbar/menu/keyboard; disabled action refused, icon-only label available, focus and accessibility inspected at host level. |
| GAP-XFMD-SDUI-008 | Missing ref file parses intentionally. Existing bridge validates module/action signature, input sources and explicit output handle before installing handlers. Runtime guards revisions and stale handles. | Extend existing explicit linking/activation with host capability preflight and new typed event/collection payloads. Do not make the parser load files. Current bridge output is a text field into an input, not a generic host binding. | Missing module, wrong signature, invalid handle and unsupported capability cause no partial activation. Reload revokes old events/handlers; host disposal prevents later mutation. |
| GAP-XFMD-SDUI-009 | Same-document definition reuse passes; proposed import syntax rejected. ref names a domain module, not a component library. | Defer cross-file UI source sets until the first real shared component requires them; coordinate source identity with SDL-001 without conflating the languages. Separate content/profile requirements from host/provider availability. | Two instances of a shared toolbar have distinct paths and common command meaning; missing/cyclic/duplicate definitions and incompatible capabilities fail with original source diagnostics. |
| GAP-XFMD-SDUI-010 | Bounded Markdown provider rejects images/HTML, flattens inline styling/links and limits registered Mermaid adapter to flowchart/graph. svg is symbolic. No-renderer output can succeed with a placeholder. | Preserve the eventual rich-Markdown goal but expose honest content profiles and explicit fallback now. Extend provider/resource interfaces per verified family; no parser-side fetching or new renderer fork implied. XFMD may retain its richer preview host. | Content measured and painted by the same provider; real rendering distinguished from placeholder. Resource lifetime, limits, links/selection and representative image/math/diagram families independently verified. |
| GAP-XFMD-SDUI-011 | Direct probes and inspected exports reproduce glyph loss, rune-based word splits and discarded list depth. | Fix SDUIPresentation under KB-SDUI-004: shared measurement/painting fallback, word wrapping with a defined long-word fallback, preserved list structure. No widget syntax needed. | Supported and missing glyphs, combining sequences, long words and nested lists at several widths; deterministic output, explicit missing-glyph behavior, host/export consistency evidence. |

## Binding and ownership design recommendation

Preserve this division of responsibility:

1. **Frontend:** validate syntax, local identities and profile rules; retain spans.
2. **Explicit input/link preparation:** resolve supplied source sets/modules and
   validate signatures. Reuse `SDL/go/bridge`; extend its restricted payloads only
   through a selected contract. Profile acceptance and host availability differ.
3. **Activation preflight:** compare required widget/event/viewport/content
   capabilities with the concrete host/provider; either reject coherently or use
   an explicitly requested static fallback. Do not silently render an inert tree
   as if it were interactive.
4. **SDUI runtime:** own UI instance identity, state revisions, event admission,
   drafts and atomic updates. Application owns business transactions and I/O.
5. **Presentation/provider:** own measurement, layout and paint/resource results.
   **Native host:** own native controls, UI thread, focus and resource disposal.
6. **SDL action runtime/application adapter:** invoke deliberately registered
   domain operations. XFMD owns document admission, open/save and generated-view
   lifecycle; its native C++ interface remains a separate integration decision.

This extends existing architecture, not a new service or mandatory ABI. Fyne is
the current SDUI host. A FOX bridge, TUI host, C ABI or new Rust renderer requires
a concrete separately selected consumer contract. Modeling XFMD does not require
porting XFMD to SDUI or replacing its current navigator implementation.

## Alternatives and costs

| Approach | Benefit | Cost / failure mode | Recommendation |
| --- | --- | --- | --- |
| Keep only static button/input compositions | Immediately useful design pictures with existing tooling | Cannot test tree/pane/focus behavior; placeholders can mislead | Retain as explicitly static baseline, insufficient for the interactive pilot |
| Implement every XFMD widget as core SDUI syntax now | Potentially expressive application mockups | Large state/host/verification surface before proven consumer value; duplicates mature editor behavior | Reject as the next delivery |
| Bounded collection/viewport plus typed host extensions | Reuses runtime and supports a meaningful navigator; room for existing editor/preview hosts | Requires precise identity/capability and host tests; extension cannot become an untyped escape hatch | Recommended staged route |
| Treat arbitrary Markdown/HTML or raw native objects as all-purpose widgets | Appears to avoid profile work | Loses predictable layout, lifecycle and portability; embeds another application/runtime | Reject as a substitute for the contract |

Ponsse needs reliable typed controls and text fidelity; the proposed KanBanTUI
could share collection identities and command meaning. These are potential reuse
opportunities, not proof that a graphical tree should determine a future TUI API.
Use actual second-consumer examples when selecting the profile.

## XFMD-local boundaries

| Source ID | Disposition |
| --- | --- |
| GAP-XFMD-DOC-001 | XFMD must reconcile semantic ownership across its 281-file inventory. The 144 files without literal blueprint filenames are a search result, not 144 proven defects. No SDL declaration per source file. |
| GAP-XFMD-DOC-002 | XFMD reconciles installed process/toolkit metadata, product release and development build identities separately. Producer study does not rewrite receipts or infer version equivalence. |
| GAP-XFMD-DOC-003 | XFMD translates maintained prose and incorporates authoritative amendments while retaining IDs and frozen evidence. Language extensions do not replace this work. |
| GAP-XFMD-DOC-004 | XFMD retains its adopted linked evidence records until it selects a machine-readable Traceability contract. Do not impose the producer's ledger schema through modeling. |
| GAP-XFMD-EVID-001 | External model coverage, native workflow acceptance and independent review remain pending. This study adds producer evidence only; it does not make the external register accepted. |
| GAP-XFMD-EVID-002 | Source inspection confirms requestOpen resolves unsaved data before reading its target, and save updates the baseline. Save followed by failed read can leave the original document clean. DocumentViews rejects malformed/unregistered URI before canceling active generation. XFMD must test these actual invariants; blanket rollback is wrong. |

No external worktree was modified or native XFMD test executed. These six entries
remain owned by XFMD's existing MAINT-XFMD-0004 and gap register.

## Adoption and unresolved choices

The proposed contracts still need owner selection: how far the first collection
profile goes, which host demonstrates it, whether an unsupported static export
rejects or deliberately shows a labeled preview, and the supported text repertoire
and fallback policy. Do not choose widget spelling before those semantics.

Every selected extension must cover source/AST/validation, canonical formatting
where provided, normalization, codegen, projections, runtime and relevant hosts.
Port maintained fixtures/consumers when semantics change; retain frozen historical
evidence without maintaining an unnecessary legacy execution path. Version by the
actual public compatibility impact, not by this study's card number.

No current requirements/profile file was changed: all proposed obligations are
recorded here and in successor cards. That preserves the distinction between a
research recommendation and an adopted language contract.

## What XFMD still needs

For native SDP navigation, XFMD consumes the existing SDPTool operations:
`discover`, `tree`, `select`, `preview` and the registered SDUI preview service.
See the [implemented producer contract](../../../SDPTool/Contract.md). These
services already cover a bounded saved, single-file model; they do not require
new SDUI widgets. External KB-XFMD-014/015 remain in backlog at this inspection;
older proposal paragraphs in those cards must be read against the current
producer contract, not mistaken for missing producer implementation.

The remaining responsibilities are distinct:

- **Producer language work:** composed SDL System/source sets, richer requirement
  identity/narrative and model-to-code/evidence bindings are covered by the
  existing follow-ups. Whole-system composition still needs that work; a native
  sidebar can first use supported independent models.
- **XFMD integration:** native tabs/tree population, selection/refresh, invoking
  the producer, showing diagnostics and previews, dirty-document admission,
  stale-response rejection and generated-resource lifetime stay in XFMD. The
  producer study will not implement them. Unsaved-buffer preview is a separate
  future source-snapshot contract; the current facade reads saved sources.
- **XFMD modeling/verification:** reconcile actual requirement and source
  ownership, amendments, documentation and real workflow tests. Language tools
  cannot create evidence that the application's behavior was tested.
- **Optional SDUI expressiveness:** a rich XFMD reproduction would still lack
  full editors, menus/dialogs, advanced panes, resource fidelity and other widget
  semantics. Under the owner's clarification this is not a blocker or a promise
  that the remaining SDUI work will implement those features.

No newly established mandatory producer blocker for the first native navigation
slice was found beyond the existing supported-profile/single-file limits. This
is a contract/source assessment, not a new end-to-end XFMD acceptance run.
