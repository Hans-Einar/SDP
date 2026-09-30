# Capabilities and a bounded navigation pilot

| Field | Value |
| --- | --- |
| id | KB-SDUI-003 |
| project | SDUI |
| type | Proposal |
| CardState | backlog |
| Systems | SDUI |
| created | 2026-09-29T16:53:54.251554+00:00 |
| source | PLAN-SDP-0009; KB-SDP-041; external XFMD gap register |
| next_review | At the next SDL/SDUI work selection |

## Owner scope clarification — 2026-09-29

SDUI primarily describes the concept, composition and placement of a UI using a
limited widget vocabulary. Running a basic SDUI and calling SDL remains part of
its purpose. It need not reproduce every function of a rich native application.
A tree is a possible useful addition, not an approved requirement to rebuild
XFMD in SDUI. Full editors, sophisticated dialogs, application-wide command
systems and host parity remain optional proposals requiring demonstrated need.
This clarification governs the study's successor selection; the gap inventory
records consumer observations, not a mandatory product backlog.

## Need and bounded next step

[PLAN-SDP-0009 study](../../02--Requirements/XFMD-Gaps/Study.md) evaluates
GAP-XFMD-SDUI-001–010. Own the follow-up here rather than creating ten loosely
coupled widget cards. First select a DesignPlan for required capabilities,
activation preflight and typed interaction state. Then evaluate whether a simple collection/viewport pilot is warranted by the
basic prototype use case; tree support remains a candidate. Registration does not authorize implementation or adopt syntax.

Reuse SDL/go/bridge typed binding validation and SDUI/go/runtime drafts, atomic
updates, identity and reload protections. Keep frontend I/O-free. A missing SDL
file passing parsing is intentional; missing supplied runtime module/signature
must fail before activation. Scroll already parses but layout rejects it.

## Staged scope and acceptance

Follow [XGP3–XGP5](../../02--Requirements/XFMD-Gaps/Delivery-Proposal.md):

- First: separate frontend, layout, provider and native-host capabilities; no
  partial activation. Typed stable item identity plus model generation; stale,
  deleted, canceled or disposed targets cannot mutate a replacement UI.
- Pilot: bounded tree/list, expansion/loading/error, selection and scrolling,
  verified with keyboard/pointer and reload in the chosen actual host. Group
  rows are not selectable; expanding alone does not navigate.
- Optional, only after a concrete basic-UI need: panes, commands/transient surfaces, typed values;
  draft transactions remain application-owned. Keep relative dimensions.
- Optional research, not required SDUI parity: component source sets, editor/preview host surfaces and broader
  resource profiles. Require a real consumer and original-source diagnostics;
  do not create an unrestricted native-object escape hatch.

New profile work must address parser/AST/normalization, diagnostics, formatting
where supported, codegen, presentation, runtime and actual host verification.
Pictures/AST do not prove interaction. Unsupported exports either reject or use
an explicit labeled fallback selected by the consumer.

## Boundaries and dependencies

[KB-SDUI-004](%23004--SDUI--Bug--Text-and-Markdown-fidelity.md) independently owns
text fidelity. [KB-SDL-005](../active/%23005--SDL--Change--System-and-source-sets.md) owns SDL
source sets, not SDUI component imports. The first SDUI fixture-driven pilot need
not wait for every SDL feature. Real generated-view integration needs a coherent
producer revision and existing SDPTool delegation.

External XFMD owns its C++ document/open/save/lease workflow and any native
integration. The external discovery register is at
`SDP/02--Requirements/XFMD-Modeling/Gaps.md` in its repository; no local checkout
is guaranteed. No FOX bridge, XFMD rewrite, new Rust renderer or Fyne editor is
selected. [KB040](%23040--Proposal--KanBan-TUI.md) and Ponsse are potential reuse
inputs, not permission to add another host now.

## Completion boundary

A design-only successor must explicitly hand remaining implementation to a
selected plan/card; this proposal cannot be closed as implemented from design
or a static export. Retain a disposition for every SDUI-001–010 gap; an explicit
out-of-scope decision satisfies disposition without requiring implementation. Native workflow evidence and independent review
are needed for claims about delivered interactive behavior.

## Worklog

- 2026-09-29T16:53:54.251554+00:00 — EVT-KB-SDUI-000012: Registered from completed producer research; remains backlog.

- 2026-09-29T17:55:11.132172+00:00 — EVT-KB-SDUI-000014: Owner clarifies limited conceptual/layout UI and basic SDL-connected execution. Rich application parity is not required; tree is a candidate. Update proposal scope; keep backlog, no implementation selected.
