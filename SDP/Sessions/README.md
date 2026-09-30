# Goal-oriented working Sessions — manual format 1

Owner request, 2026-09-30: adopt and distribute the bounded manual format under
[MAINT-SDP-0011](../Maintenance/RL1/Plan.md). The reusable distribution guide is
[Template/sdp-root/Sessions](../../Template/sdp-root/Sessions/README.md).
Existing cards/plans and ProjectManagement remain authoritative. No Session ledger
kind, automatic transcript recorder or event timeline engine is introduced.
[KB-SDP-042](../KanBan/backlog/%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md)
retains broader identity/capture integration; KB043 retains generated timelines.

## The missing overview

A Session follows one goal across turns, cards, plans and possibly replacement
agents or chat threads. A Turn is an interval of agent work ending in a returned
result or interruption; a steering prompt may join an existing execution turn.
A Step is a meaningful piece of the route and may span many turns. One turn may
advance multiple steps. Keep these identities separate.

One primary card normally starts a Session. Other cards may join it, and a card
may participate in later Sessions if its outcome needs follow-up. A Session is
not a Sprint: a Sprint selects a work package, while the Session preserves the
journey toward a goal and may link several plans or Sprints. Do not mandate all
six plan types for every task. A trivial answer does not need a Session file.

## Source of truth and states

- The Session owns goal, scope, step dependencies, proposed next step, revisions
  to the route and links to turns. It records why the direction changed.
- Cards own lifecycle/CardState; typed plans own their existing planned, active,
  completed or canceled lifecycle and milestone evidence. Session tables are
  dated projections. Preserve initial and planned-final snapshots; do not label
  the current state as an actual final outcome before disposition.
- A plan registry may show **document readiness**: planned (not created), planning,
  ready, on-going, completed, canceled or superseded. This is separate from its
  canonical lifecycle. A ready plan does not imply execution authorization.
- Proposed step display states: planned, next, on-going, waiting, blocked,
  completed, canceled, superseded. State is the first table column. Normally one
  primary next step; explicit parallel lanes may each have one. Store dependencies,
  reason/authority and completion evidence beside the status.
- Completing a turn does not complete a step. Completing a plan does not complete
  a broader card. A Session can close only with goal evidence and an explicit
  disposition for every included step/card; transferred scope names its successor.

Do not copy detailed milestones into a competing plan. Link to them. When a
Session projection disagrees with its source, show it as stale and reconcile it;
never silently change a card to make the overview look consistent.

## Proposed turn routine

Before work: recover goal/current step, classify the prompt, select existing
routine(s), check scope/authority and identify any missing procedure. During work:
record meaningful discoveries and replan affected steps rather than silently
changing the goal. At return: link evidence, update statuses, name the next step
and explain why. Record blocked/interrupted turns honestly and resume from durable
state. A new idea may be captured in backlog without derailing the active goal.

The manual pilot records loaded skills and agent-reported procedure use, not an
unobserved claim of enforced routine execution. Future entries need routine ID,
version, run ID, relevant step/check outcomes and provenance. The existing
[governance study](../02--Requirements/RoutineGovernance/Synthesis.md) and KB036/037
own routine selection and runtime enforcement; Sessions must not duplicate them.

## Roadmap and revisions

Owner refinement, 2026-09-30: put the roadmap immediately after the document
title, before metadata, goal and history. Keep a Gantt view followed by a step table. Derive the diagram from the step data
when tooling exists; in this pilot both are maintained together. Use text states
as well as colors. Preserve completed history; added/removed/reordered steps get
a dated reason and source turn. Keep original intent and revised target distinct.

Where no estimates exist, label the chart **sequence only**. The pilot uses a
synthetic 2000-01-01 anchor and one equal slot per step, not dates or durations
promised to the owner. A true calendar schedule requires actual estimates.
Mermaid is a projection, not the state store. The syntax follows the
[official Gantt reference](https://mermaid.js.org/syntax/gantt.html); no init directive
or custom click handler is needed. Diagram source must remain useful as Markdown.

## Event-derived timeline proposal

[KB-SDP-043](../KanBan/backlog/%23043--Proposal--Event-derived-session-timelines.md)
owns the detailed proposal: one planned task per step lane, turns across time,
colored work segments, point events and current/as-of position. KB042 continues
to own Session format/adoption. Distinguish a planned Task from an activity Segment;
creating a card or document can segment work on the same lane without adding a
new planned task. Turn markers and event markers have distinct roles.

Use a turn-ordinal overview or a measured seconds view with explicit turn
boundaries. Unknown timing stays unknown. A static export freezes the cursor at
its observation point; it does not claim live monitoring. ProjectManagement owns
card/plan lifecycle; routine/client capture supplies observations through a future
versioned correlation design. No unsupported Session events or parallel lifecycle
ledger are introduced by this document update. The current sequence-only diagram
remains a fallback, not a turn/time measurement or implemented generator.

## Turn capture and app-server

For a custom client that owns the conversation, app-server is the recommended
capture point: retain submitted prompts and completed user/agent message items,
correlate thread/turn/item IDs, and use turn completion status. Steering input can
belong to the same turn. This recommendation follows the
[official app-server documentation](https://learn.chatgpt.com/docs/app-server),
checked 2026-09-30; it is not a tested recorder implementation.

The adapter must preserve multiple inputs, distinguish commentary from final
answers, handle absent phase metadata, interruptions and missing final messages,
and deduplicate replay after reconnect. App-server turn completion is not SDP
acceptance. A Session ID must survive thread replacement. MCP provides explicit
SDP operations; it is not a passive transcript recorder. The existing
[version-pinned app-server study](../02--Requirements/RoutineGovernance/AppServer-Study.md)
remains the compatibility baseline, with a later implementation probe required.

Save exact owner prompts and visible final responses when observed. Link original
attachments with availability information; do not pretend a text excerpt captures
an image. Do not collect private reasoning. Record any explicit redaction instead
of silently presenting modified text as verbatim. Never extract credentials into
the record. Automatic capture in an arbitrary already-running TUI is unverified.
Until integration exists, manual entries must distinguish exact quotations from
summaries and leave unavailable host IDs/timestamps unknown. An agent cannot claim
it has captured its eventual final response through tools before sending it.

## Remaining automation work

MAINT-SDP-0011 adopts/distributes the manual format; shared DesignPlan work
under KB042 with KB036–038 still covers machine ownership,
Session IDs, routine-run/turn correlation and capture. Evaluate an extension to
the existing project-management ledger/schema; do not invent a parallel lifecycle
ledger or append unsupported Session events today. Transcript items are content,
not one management event per token. Traceability remains system design/code evidence.

The [template](Session-template.md) and
[first proposed roadmap](session-%230001--SDL_expansion.md) are the pilot artifacts.
This pilot registers links in card prose, without adding unsupported required
SessionId metadata to existing schemas. Installation inventory now includes only the reusable guide/template, with
initialize-if-missing ownership; published upgrade availability is recorded in
MAINT-SDP-0011 rather than inferred from local template presence.

[Session 0002](session-%230002--SDPTool_output.md) tracks portable SDPTool output,
including its explicitly marked late registration.
