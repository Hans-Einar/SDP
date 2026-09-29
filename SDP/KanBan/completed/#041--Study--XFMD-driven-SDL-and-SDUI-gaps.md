# Study XFMD-driven SDL and SDUI modeling gaps

| Field | Value |
| --- | --- |
| id | KB-SDP-041 |
| project | SDP |
| type | Study |
| CardState | completed |
| Systems | SDL, SDUI, SDPTOOL |
| PlanId | PLAN-SDP-0009 |
| created | 2026-09-29T11:39:17.906434+00:00 |
| source | Owner request 2026-09-29; external XFMD modeling gap register |
| tags | language, modeling, host-capabilities, consumer-feedback |
| next_review | When selecting the next SDL/SDUI study; before finalizing related extension contracts |

## Need and authority

The owner requests a producer-side card referencing XFMD's discovered modeling
gaps. The next selected task should study **whether and how** each gap belongs in
SDL/SDUI. Registration does not adopt proposed keywords/widgets or authorize
language, runtime, renderer or XFMD application changes. This is one cross-language
study card on the shared board; existing primary cards retain their owned scope.

## External source and provenance

External project: XFMD, worktree `xfmd-sdl-navigation`. Source-relative path:
`SDP/02--Requirements/XFMD-Modeling/Gaps.md`.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/02--Requirements/XFMD-Modeling/Gaps.md

The external report belongs to XFMD and remains its maintained discovery register.
Its linked baseline/probes and MAINT-XFMD-0004 are external evidence, not local
SDP records and not dependencies that board validation must resolve. The local
path is a reading convenience; other checkouts need not have that worktree.
Do not copy the report into SDP as a competing authority or invent a public link.

| Observation | Value |
| --- | --- |
| Read at | 2026-09-29T11:39:17.906434+00:00 |
| Observed XFMD HEAD | 41fa4948c5377edd792ba3f503e6b92b0fd4717a |
| File state at inspection | Untracked; observed HEAD does not identify this file's bytes |
| Report SHA256 | 6bb59190e2a35ccd8034137009e2fbf044c3793680eab2fbf24ea754ba169f6b |
| Report-declared producer baseline | 3b4cff695ecab36a0e897052e3f096bbeb9f3ea0 |

Refresh the external report and pin the studied bytes when the study starts.
Preserve original GAP-XFMD IDs in the result table as its contents evolve. Claims
below summarize the external report; this intake did not rerun its producer or
native interaction probes. New ecosystem modeling after its baseline is source
organization, not evidence that the reported parser/runtime gaps are resolved.

## Gap inventory and existing ownership

| Source IDs | Study concern | Existing route / distinction |
| --- | --- | --- |
| GAP-XFMD-SDL-001 | System identity, explicit multi-file membership, resolution and original-source diagnostics | [KB-SDL-005](../backlog/%23005--SDL--Change--System-and-source-sets.md); reuse its contract work, do not create a second loader |
| GAP-XFMD-SDL-002 | Requirement identities, normative narrative, amendments and acceptance/source links | [KB-SDL-001](../backlog/%23001--SDL--Proposal--Requirements-narrative.md); preserve external hyphenated identities without assuming new keyword spelling |
| GAP-XFMD-SDL-003 | Detailed state, failure/cancellation invariants and source/evidence bindings | [KB-SDL-006](../backlog/%23006--SDL--Study--Executable-channel-tests-and-unit-bindings.md) and [KB-SDP-004](../backlog/%23004--Proposal--Design-traceability.md); first assess existing structural/action/class/scenario profiles |
| GAP-XFMD-SDUI-001 | Trees, lists, stable item identity and lazy collections | Collection semantics and host interaction, not merely drawing rows |
| GAP-XFMD-SDUI-002 | Tabs, splitters, retained pane proportions and focus | Distinguish composition from interactive pane state |
| GAP-XFMD-SDUI-003 | Scrolling and source-linked viewports | Reported syntax already parses; identify layout/host capability and synchronization responsibility |
| GAP-XFMD-SDUI-004 | Menus, dialogs, contextual targets and focus ownership | Host transient surfaces, shared commands and explicit cancellation |
| GAP-XFMD-SDUI-005 | Multiline editor and rich preview surfaces | Evaluate host-surface contracts; do not absorb XFMD editing/domain behavior into the language by default |
| GAP-XFMD-SDUI-006 | Numeric, boolean and choice controls; transactional forms | Reuse existing draft/revert and atomic update behavior; application owns persistence/transactions |
| GAP-XFMD-SDUI-007 | Icons, toggles, command state, keyboard and accessibility | Shared semantic command contract rather than host constants or arbitrary styles |
| GAP-XFMD-SDUI-008 | Callback signatures, binding resolution and host capability checks | Separate linking/activation checks from an I/O-free parser; preserve revision/handle protections |
| GAP-XFMD-SDUI-009 | Cross-file UI reuse and parser/layout/host profile negotiation | SDL source sets and SDUI component imports are related but distinct contracts |
| GAP-XFMD-SDUI-010 | Markdown, images, math, Mermaid families and rich resource behavior | Provider/renderer capability and explicit fallback, not automatic grammar expansion |
| GAP-XFMD-SDUI-011 | Glyph coverage, word wrapping, nested lists and font fallback | Diagnose implementation/provider fidelity before proposing syntax; visual/host evidence required |

GAP-XFMD-DOC-001–004 and GAP-XFMD-EVID-001–002 remain XFMD-owned documentation,
identity, semantic coverage and verification concerns. Include their implications
in the study's evidence assessment, but do not convert them into SDL syntax
requirements or claim this card closes them. In particular, preserve actual
Save-then-failed-open semantics and ordering of invalid-URI rejection/cancellation;
a failed operation does not universally imply full state rollback.

## Study questions and expected result

Create a proportionate StudyPlan when selected, using the adopted RequirementPlan
record type. Link its result back here. For each of the 14 SDL/SDUI gap IDs, record:

1. The concrete XFMD workflow/invariant and baseline evidence; rerun a minimal
   relevant probe against the selected producer candidate, identifying untested
   behavior separately from observed failure.
2. Whether existing syntax/profiles/APIs already express it, whether documentation
   or an implementation fix suffices, or whether it needs grammar/model semantics,
   binding validation, runtime state, layout/presentation or a host/provider API.
3. A recommended disposition: existing support, clarify, fix, propose extension,
   defer or retain application ownership. A finding need not lead to language growth.
4. Alternatives, reuse and costs against SDUI's deliberately bounded prototyping
   purpose. Compare a generic host surface with built-in widgets where appropriate;
   distinguish reusable needs shared with KanBanTUI/Ponsse from XFMD-specific ones.
5. Validation cases, consumer migration and the owning card/plan. Any selected
   language change must update Go parsing/model checks, diagnostics, formatting,
   projection/codegen where affected, and independent consumer/host evidence.

Use one small XFMD navigation workflow as the first integration candidate. Cover
success, rejection/cancel, stale identity after reload and resource lifetime as
appropriate. Static pictures and AST acceptance alone do not prove interactive
semantics, complete Markdown fidelity or working C++/Go/native bindings.

## Sequencing and completion criteria

Coordinate System/source-set decisions with KB-SDL-005 first; this report is a
consumer input, not a requirement to finish every SDUI widget before source sets.
Evaluate capability/binding contracts before selecting larger widget additions.
Then propose a small implementation order based on the studied dependencies.
[The ecosystem catalog](../../SDL/Catalog.md) supplies current tool ownership;
[KB040](../backlog/%23040--Proposal--KanBan-TUI.md) is a potential second consumer, not proof
that its needs or implementation are already established.

The study completes when all 14 language/UI entries have evidence-qualified
recommendations and explicit follow-up ownership, the six XFMD-local entries
have a boundary disposition, and a proposed phased delivery/verification plan
is available for the chosen subset. No implementation commitment follows from
registering or completing research. The owner authorized study execution on 2026-09-29 under PLAN-SDP-0009; implementation remains unselected.

## Worklog

- 2026-09-29T11:39:17.906434+00:00 — EVT-KB-SDP-000237: Read the external gap register, recorded its
  untracked-file provenance and all 14 SDL/SDUI gap IDs, linked existing producer
  cards, and registered this study in backlog. No producer probes rerun, no XFMD
  writes and no language/runtime capability adopted.

## Selected study

[PLAN-SDP-0009](../../02--Requirements/XFMD-Gaps/StudyPlan.md) now owns execution.

- 2026-09-29T16:40:18.297106+00:00 — Study activated; CardState in-progress.

## Outcome — XGS3-M1

[Study](../../02--Requirements/XFMD-Gaps/Study.md),
[evidence](../../02--Requirements/XFMD-Gaps/Evidence.md) and
[delivery proposal](../../02--Requirements/XFMD-Gaps/Delivery-Proposal.md)
complete the authorized evaluation of all 14 language/UI and six external
concerns. The study recommends existing SDL source-set work, independent text
fidelity fixes and a capability-first collection/viewport pilot. Existing typed
bridge and runtime state protections are retained. No product gap is closed by
research and no implementation or owner design acceptance is inferred.

Existing KB-SDL-001/005/006 and KB-SDP-004 retain ownership. New
[KB-SDUI-003](../backlog/%23003--SDUI--Proposal--Capabilities-and-navigation-pilot.md)
and [KB-SDUI-004](../backlog/%23004--SDUI--Bug--Text-and-Markdown-fidelity.md)
retain all SDUI follow-up in backlog. XFMD-local work remains external.

Actual capability probes and ten selected Go packages support bounded findings;
external Mermaid test skipped, no native host/XFMD or independent review claimed.
Git commits/pushes could not be made under the read-only .git mount.

- 2026-09-29T16:53:54.251554+00:00 — EVT-KB-SDP-000240: Study completed; follow-up remains backlog.

## Owner follow-up — 2026-09-29

The owner confirms limited conceptual/layout prototypes and basic SDL-connected
execution as SDUI's purpose; recreating a rich XFMD interface is not required.
The study/proposal and KB-SDUI-003 distinguish optional extensions from native
XFMD navigation responsibilities. Research remains complete. Git access is now
restored and the owner requested the consolidated study recovery commit.

- 2026-09-29T17:55:11.132172+00:00 — EVT-KB-SDP-000241: Record owner scope clarification and authorized Git recovery after study-time read-only restriction. Research remains completed; distinguish existing producer services from external XFMD integration.
