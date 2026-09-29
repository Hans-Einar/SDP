# Goal-oriented Sessions, roadmaps and turn history

| Field | Value |
| --- | --- |
| id | KB-SDP-042 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| Systems | SDPTOOL |
| created | 2026-09-29T22:21:29.329185+00:00 |
| source | Owner Session proposal, 2026-09-30 |
| next_review | Before adopting/distributing the Session format or selecting routine/client implementation |
| tags | sessions, roadmap, planning, routines, turn-history |

## Need and owner intent

Make the route toward a goal visible across many turns: what has happened, what
is happening, what comes next and why. Usually start from one primary card and
include other cards as needed. Preserve direction changes rather than relying
on chat memory. Use SDP/Sessions/session-#NNNN--Topic.md with Goal, card snapshots,
plan registry, Gantt/step table and per-turn prompts, responses and routine use.
Cards link back to Sessions that handled or plan to handle them.

## Reviewable document pilot

- [Session guide](../../Sessions/README.md): distinctions, authority and capture boundaries.
- [Template](../../Sessions/Session-template.md): required sections and state projections.
- [SDL expansion roadmap](../../Sessions/session-%230001--SDL_expansion.md): concrete
  proposed sequence, existing study input, uncreated plan placeholders and manual
  first-turn entry. It does not activate SDL implementation.

The entry card owns the goal; a Session connects its route across plans/turns.
Existing typed plans retain execution detail and evidence. Document readiness
(planned/planning/ready/on-going/completed) is a separate view from adopted plan
lifecycle (planned/active/completed/canceled). Do not introduce conflicting state
machines in Markdown. A future projection should derive card/plan state and Gantt
from authoritative data; this first pilot is manually maintained and labeled.

## Ownership and bounded adoption

KB042 owns the Session format, identity, roadmap/plan relationships and adoption.
[KB036](%23036--Proposal--Request-routines-and-process-control.md) still owns
request classification and routine selection; [KB037](%23037--Proposal--Observable-routine-state-machine.md)
owns execution state/monitoring; [KB038](%23038--Proposal--Codex-app-server-development-client.md)
owns app-server capture/client integration. Do not create a second routine engine
or a new client project here.

Next select a proportionate MaintenancePlan for manual adoption, template/skill
changes and distribution impact. A shared DesignPlan with those existing cards
must settle programmatic ownership, schema changes and replay/correlation before
implementing automation. This registration and draft do not adopt Session as a
new management kind or PlanType. Existing ledger cannot yet accept Session events;
record this proposal/review using supported KanBan events. Traceability stays clean
of management-only events. No installer/template payload changes yet.

## Acceptance and edge cases

1. Start a goal from a primary card; include several cards/plans without copying
   them. Distinguish initial, planned-final, current and actual-final dispositions.
2. One visible proposed-next step, with dependency/authority; multiple explicit
   parallel lanes allowed. Completed turn does not imply completed work. Interrupted,
   blocked, canceled and superseded steps have reasons and resume/successor links.
3. Replan after a discovery while retaining original intent and route-change history.
   Preserve membership snapshots when scope splits, merges or moves to another Session.
4. Resume the same Session after replacing the agent/chat. Session IDs are independent
   of app-server thread/turn IDs and optional Sprint membership.
5. Capture exact visible prompts and final responses when observed, including steering
   inputs; distinguish manual summaries and unavailable attachments/IDs. Reconnect
   must deduplicate items and must not manufacture a final response for an interrupted turn.
6. Record routine identity/version/run, checks and provenance. Available skills,
   declared skill use, observed tool activity and validated routine transitions are
   different evidence. Missing routine must remain visible and route to existing KB036.
7. Render the same roadmap as Gantt and table with status labels; no fictitious
   deadline when only sequence is known. Existing card/plan changes expose stale
   snapshots instead of silently overwriting history.
8. Close with goal evidence and explicit card dispositions; new ideas outside the
   goal remain backlog. Lightweight use must not require six plans or a Sprint per task.

## App-server assessment

For a client that owns conversation submission, app-server is the recommended
capture point. The official contract supplies user/agent message items and
turn completion, and steering may append to an existing turn. Correlation,
replay, final-response identification and retention remain client work. MCP
alone does not observe all prompts. Current TUI passive capture is unverified.

[Official app-server documentation](https://learn.chatgpt.com/docs/app-server),
checked 2026-09-30; [local versioned study](../../02--Requirements/RoutineGovernance/AppServer-Study.md).
No live protocol/model/account operation was executed for this registration.

## Source prompt

Verbatim visible owner text, with Markdown blockquote prefixes and trailing
whitespace normalized for storage. This is manually captured, not a host transcript.

> ok so how do you now keep track of each step that is next and next and next? we did use plan documents for requirement plan, architecture plan, design plan and so on. but what I'm really missing now is a visualisation of all the steps we plan ahead during a work session like this. maybe we need to implement yet another structure that we call a session. a session can be "whatever we do", so everything we work on is a working session. but I usually use the word session when I ask you to do "all steps in one session" maybe we should call that a "turn". cause this is more or less turn based. I give you a prompt, and you do something in your turn and comes back with a summary of your turn.  a session is when we work in multiple turns toward a goal. so maybe now we are able to describe it all...
>
> when we start to work on something, we usually start from a kb card. that kb cards describes something we want to get done. but there might be several kb cards involved to actually realizing a goal. and on the road to realizing that goal, we often have to have many turns back and forth, and sometimes the direction changes slightly depending on what we discover on the road to realizing the goal. when we are on that road towards the goal, we can call it a session. so what we need then is a sessionPlan that can describe the steps we need to take along the road to reach our goal. that session plan can also be viewed as a roadmap in a mermaid gant schema / diagram. I suggest that we create a folder in projects SDP folder called Sessions. in that we create a document for each session. session-#0001--SDL_expansion.md  for example. in that document there will be  several chapter defined in a template that the session have to fill in.  first is the Goal.  next is a list with link to KB cards that are affected. the table could have a column that shows the state of the kb card at the beginning of the session. the planned end state and the actual end state. if it was canceled, superseeded, completed and so on.
>
> the next chapter shows a table with link to all  *Plan documents" to be made during the session. the states of the documents should be shown in a column.  for example planned, planning, ready, on-going, completed.
>
> then comes the actual session plan with a gant diagram first and then a table for each step. the first column can show the state of the steps, like planned, on-going, completed and next. when you are working on a step it will show up as on-going, and when you have completed the step, it shows up as completed. then the next step will have state next step.
>
> then comes the chapter that is actively filled out during the session which should add some lines for each turn. it would be great to save here every prompt I give you, and every summary you come back with for each turn.
>
> is it easiest to acheive that last one if we go through the app server? then we also enter a section for which routines that has been usesd for each turn. I'm too tired now to write more, but I think you undersand the point.
> one KB cards should ideally be the entry point of one session. the kb card should be updated with a link to the session that took care of it.  i think the idea to add sesssions and routine checks now completes our work process with sdp.

## Worklog

- 2026-09-29T22:21:29.329185+00:00 — EVT-KB-SDP-000242: Proposal and provisional guide/template/roadmap recorded; remains backlog.

- 2026-09-29T23:38:37.439968+00:00 — EVT-KB-SDP-000247: Correct Session Gantt after owner-reported XFMD parser rejection: remove axisFormat and todayMarker. Standalone renderer acceptance did not prove XFMD-path compatibility; card remains backlog.

## External Mermaid compatibility handoff — 2026-09-30

[Prepared XFMD card](../../02--Requirements/XFMD-Gaps/XFMD-Mermaid-Compatibility-Card-Draft.md)
covers broader Mermaid gaps and the actual parser/model/renderer boundary.
**Registration in XFMD is pending**: current permissions allow reading but not
writing that repository. No KB-XFMD ID was reserved or lifecycle event invented.
The next write-capable XFMD session must allocate the card and update its index
and ledger together. This is external follow-up, not a new SDP renderer scope.

2026-09-29T23:44:34.061920+00:00 — EVT-KB-SDP-000248: Prepare owner-requested external XFMD Mermaid compatibility card covering all families/construct categories and Gantt/gitGraph directive regressions. XFMD registration remains pending because this session cannot write that repository; no external ID reserved.

## External handoff registered — 2026-09-30

The earlier pending handoff is resolved: external **KB-XFMD-020** is registered
in XFMD with index and ledger. [Handoff pointer](../../02--Requirements/XFMD-Gaps/XFMD-Mermaid-Compatibility-Card-Draft.md).
The external card owns follow-up; no application change selected here.

- 2026-09-29T23:48:33.074200+00:00 — EVT-KB-SDP-000249: Owner restored write access: register external KB-XFMD-020 with its board index and ledger, replace local draft with a handoff pointer; Session proposal remains backlog. No external reference resolution required.
