# Event-derived Session roadmaps and turn timelines

| Field | Value |
| --- | --- |
| id | KB-SDP-043 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| Systems | SDPTOOL |
| created | 2026-09-30T10:36:36.093682+00:00 |
| source | Owner Session timeline proposal, 2026-09-30 |
| next_review | With KB-SDP-042 format adoption and KB-SDP-037/038 observability/client design |
| tags | sessions, timeline, events, lineage, turns, visualization |

## Goal and ownership

Generate the Session roadmap/Gantt and table from recorded steps, tasks, turns and
activity/lifecycle events. The owner wants an immediate view of what happened,
what was added, how many turns a task consumed and what is happening now.

[KB042](%23042--Proposal--Goal-oriented-sessions-and-roadmaps.md) owns the Session
format/adoption. This card owns timeline projection and its event requirements;
[KB037](%23037--Proposal--Observable-routine-state-machine.md) owns routine execution
observability, [KB038](%23038--Proposal--Codex-app-server-development-client.md)
owns client capture, and SDPTool owns process operations. Avoid a second lifecycle
engine or parallel authoritative history. [Session 0001](../../Sessions/session-%230001--SDL_expansion.md)
is the originating pilot; timeline implementation is not part of its SDL delivery.

## Owner's vocabulary and visual behavior

- Step: one stable horizontal lane, arranged vertically in roadmap order.
- Task: the planned unit of work for that step; one planned task per step. The
  current activity is a segment of that task, not a new planned task per tool call.
- Turn: an observed interaction/execution interval. A task can span turns, and
  one turn can touch several steps. Preserve stable local and optional host IDs.
- Segment: interval on the same lane showing the activity performed, with links
  to its task, turn, event and affected artifact. Split the visual bar when focus
  changes; do not erase the earlier segment or turn the split into card lineage.
- Event: point occurrence, such as card created, CardState changed to in-progress,
  document created, task resumed or completed. Display an event marker and label.

Place SessionRoadmap (Gantt then table) immediately after the document title,
before goal/metadata/history. Show current turn and next step near it. This order
is applied to the local manual pilot/template now; automatic generation is proposed.

Owner's marker sketch, requiring visual prototyping rather than literal Mermaid:

| Marker | Intended meaning |
| --- | --- |
| Upward marker below the lane, `^` with a vertical stem | Turn boundary; label T001, T002, etc. |
| Downward marker above the lane, `v` with a vertical stem | Event occurrence within the turn |
| `#` with card ID | KanBan card activity/event; link to the card |
| `@` with filename | Document activity/event; link to the artifact |
| `<\|` | Suggested segment-end marker |
| Vertical line through lanes | Current turn/current recorded position |

Use red for KanBan activity and green for document activity, with a distinct
ordinary task color between boundaries. An event may start a new colored segment
on the same lane. Preserve labels/icons/patterns so color is not the sole meaning
and red does not silently mean failure. Card activation is an event; card work
may be an interval. File creation is an event; writing/editing may be an interval.
Do not infer a long work interval solely from the existence of one event.

## Axis and timing

The owner first proposed turn numbers on the x-axis, then refined the idea to
seconds with clearly marked turns and events at offsets from each turn's start.
Study two projections of the same records:

1. Turn overview: ordinal turn bands for counting turns per task; within-turn
   placement is schematic unless measured timings are available. Equal-width
   turns must not be labeled equal elapsed durations.
2. Detailed elapsed-time view: seconds, turn-start boundaries and event offsets.
   Preserve actual timestamps/turn boundaries. Decide explicitly whether idle
   gaps remain in wall-clock scale or appear as marked axis breaks. Never silently
   label concatenated active intervals as wall-clock time.

A live current-turn line follows observed state; a static export freezes it at
its recorded as-of event/time. It is not a calendar todayMarker. Completed Sessions
have a last-recorded marker, not a fabricated running turn. Count distinct turns
that touched a task, not the difference between first/last turn numbers. Elapsed
seconds are not proof of active effort; waiting/interruption should remain visible.
The current manual journal has no reliable turn-start/subturn timestamps. Keep
those unknown; ledger write time is not retroactive execution time.

## Record and lineage design to investigate

ProjectManagement owns actual card/plan lifecycle events. Project history already
has timestamps and IDs, but not the complete task/turn/activity correlation needed
for this chart. Link to its existing events rather than repeating state changes in
a second ledger. Routine/client observations may need versioned telemetry records
and a deliberate schema extension before replay can drive this projection.

Candidate information: sessionId, stepId, taskId, turnId, segmentId/eventId, kind,
observedAt, turn-relative monotonic offset, activity start/end, subject reference,
source event ID, previous/causal/superseding reference and capture provenance.
These are design inputs, not an adopted schema. Distinguish recorder reception time
from action occurrence; retain unknown timestamps and manual summaries honestly.
Handle duplicate/replayed/out-of-order observations and clock resets/reconnects.
No private reasoning or unobserved tool execution may be invented for the chart.

Preserve route revisions: when scope is added/replanned, retain the old task intent
and a successor/revision relationship while each step still has one current planned
task. KanBan merge/split lineage remains its existing separate meaning. Link it
when relevant; a red segment is not itself a card merge/split operation.

## Renderer and delivery boundary

First select a supported representation and test the complete consumer path.
The existing pilot has observed XFMD rejections for axisFormat/todayMarker, and
standalone renderer acceptance did not prove native viewer support. Do not promise
that Mermaid Gantt provides all proposed glyphs, same-lane segmentation, coloring
and live cursors. Evaluate Mermaid's verified subset versus generated SVG embedded
in Markdown or a later native/client timeline. No new renderer is selected here.
External KB-XFMD-020 owns broader Mermaid compatibility; no XFMD changes authorized.
Markdown/table fallback must retain meaningful events and links.

## Acceptance for a future plan

- Replay recorded inputs twice: same rows, segment IDs, events, labels and links;
  no invented timestamps or lifecycle transitions.
- Show one task across several turns with nested card/document activities on its
  original step lane; preserve one planned task and all earlier segments.
- Show card creation and activation, document creation/edit, pause/resume and
  task completion at their actual recorded positions. Explain missing timings.
- Count distinct contributing turns correctly; maintain visible turn boundaries,
  a live/as-of current marker, event cursors and accessible colors/text.
- Replanning, canceled work, concurrent activity, interruptions and reconnect replay
  cannot overwrite history or falsely show a task completed.
- Gantt/timeline and table are projections of the same records with revision/as-of
  metadata; stale generated views can be identified and regenerated.
- Verify native consumer/export behavior, not just source syntax. Follow artifact
  links after a card moves using stable identity resolution or explicit unresolved
  state; do not assume the old filename remains valid forever.

## Worklog and limits

Captured in Session T009; EVT-KB-SDP-000251. Current delivery registers the proposal
and reorders the manual Session. No timer, event collector, generated timeline,
new management schema, app-server client or Mermaid extension is implemented.
