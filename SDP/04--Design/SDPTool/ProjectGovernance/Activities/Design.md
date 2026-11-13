# Activity management — lifecycle, documents and navigation

Design candidate, 2026-11-11. Owner selects Activity integration **before** the SDL
browsing redesign. [KB061](../../../../KanBan/completed/%23061--Proposal--Need-driven-activities.md)
and [Session0016](../../../../Sessions/session-%230016--Activity_design.md) own this
design. [Initial study](../../../../02--Requirements/ProjectGovernance/Activities/Study.md)
records the need. Proposed wire names/layouts below are not installed capabilities.
The owner later selected real manual registration under the
[bootstrap contract](../../../../ProjectManagement/Activity-records.md); runtime
Activity services remain unimplemented.

Owner refinements, 2026-11-12: README is activity-independent directory guidance;
Mandate.md owns the need/outcome; Sessions lives directly under the Activity root.
System-owned requirements/models and the shared canonical ledger boundary are
confirmed. Unrelated discussion/capture is allowed, but implementation requires
the selected Activity scope. Timeline event markers supplement lifecycle intervals.

Owner sequence refinement, Session0016 T6: Activity0001 restructures SDPTool with
no new functionality. Activity support is the next, manually created Activity;
browsing is the third and first created through SDPTool. Existing Activity model
and this feature design remain future work under on-hold KB061. This supersedes
the combined architecture-migration/Activity delivery proposal.

## 1. Purpose and selected boundary

Owner Session0016 T3 selects SDPTool-based creation of Activities and Sessions.
[Command design](Commands.md) defines the proposed facade, shared service placement
and the distinction between Session document creation and PG execution.
[ActivityWorkflow WORK](../../../../SDL/ProjectGovernance/SDPTool/ActivityWorkflow/README.md)
now expresses the bounded target in validated SDL; runtime implementation remains open.

An Activity follows one owner/customer need through study, requirements, design,
implementation and verification. It owns scope, need-level outcome, current route
and links to the work. Plans own milestones; cards own bounded obligations;
Sessions own work history; PG owns actual runs/attempts. Management Activity does
not reinterpret the existing SDL activity declaration.

Deliver Activity documents and lifecycle, producer-generated Activities navigation
replacing the Sessions tab for adopted projects, truthful event timelines, and
integration with existing Session/plan/Blueprint/PG identities. Keep the current
six supported PlanTypes. No new RefactorPlan or parallel execution engine is needed.

Activity work is now the predecessor of PLAN-SDP-0025 / KB060. This supersedes the
earlier suggestion that browsing should proceed independently. Preparing this
design is authorized; implementation readiness requires the concrete contract,
schema/model and plan work described below, not an assumed runtime feature.

## 2. Activity home and document ownership

Recommended new-record layout:

```text
SDP/
  Activities/
    README.md
    Activity-0001--Activity-Management/
      activity.yaml
      README.md
      01-Mandate/
        Mandate.md
      02-Study/
        Study.md
      03-Requirements/
        Requirements.md
      04-Design/
        Architecture.md
        Design.md
        Models.md
        Blueprints.md
      05-Implementation/
        Plan.md
      06-Verification/
        Verification.md
        Evidence/
      Sessions/
        session-0016--Activity-design.md
```

This is an illustrative future layout, not allocated identities or a command to
move current files. Preserve the owner's six activity folders without replacing
the project's existing five-phase process or adding numbered top-level phases.
Create the useful documents; do not fill every folder with empty placeholder text.
The producer may display all six phase nodes, clearly marking not started or not
applicable sections. Architecture and Blueprint references belong under Design.

Activity README.md is standard directory guidance, identical in meaning across
Activities: one short section per folder, its standard contents and relative
navigation links. It contains no Activity-specific title, need, desired outcome,
owner, status, checklist progress or route. Put the actual need and desired outcome
in 01-Mandate/Mandate.md. The generated Activity landing page assembles a summary
from the mandate, descriptor and ledger; it is not the directory README. Evolving
the standard README template does not silently overwrite project-local documents.

Activity-local studies, mandate, change requirements, plans and verification can
be authored directly here. Current system requirements and SDL source retain
their authoritative system-owned homes. Activity Requirements records the delta,
rationale and links; update the actual current requirements rather than copy them
into a competing specification. Models.md identifies the relevant WORK/reviewed
model and selected revisions. Blueprint references point to their retained
artifact/assignment identities. Do not create a second mutable model or Blueprint
store in each Activity. Evidence copies may be immutable snapshots with provenance.

Folders provide the default document category. activity.yaml supplies identity,
typed references and explicit exceptions, not a hand-maintained copy of every
discovered file. Ordinary documents in the six folders are discovered directly.
References support older documents that remain at their current paths.

### Identity and descriptor

Proposed canonical identity: ACTIVITY-<PROJECT>-NNNN. The folder number is the
project-local presentation of that identity; the title/slug may change without
changing identity. Allocate under the canonical project writer, never from a
filename guess alone. Reject duplicate identity homes and ambiguous registrations.

Proposed descriptor version: sdp-activity/0.1. Required content:
identity, project, title, owner/accountable role, mandate reference and source
card references. Optional: dependency Activity references, systems, external
project contributions and typed document/plan/model/Blueprint/Session references.
Reference entries include role, canonical identity where available, project/path,
and revision when a reviewed/immutable input is required. Planned routes may name
uncreated outputs as intentions without inventing paths or record identities.

**Do not author current lifecycle state in activity.yaml.** Ledger replay owns it.
Any displayed status/history is a derived snapshot with ledger revision/as-of
time. Keep need-specific narrative in Mandate.md and progress in its existing
plan/evidence sources, not in the standard README.

## 3. Sessions and phase progress

Use Activity-root Sessions/ as the canonical home for new local Session logs.
Study, requirements, design, planning, implementation and verification Sessions
share that home. Each linked Session records its phase or phase contributions and
selected plan milestones. The UI shows one primary Sessions node directly under
the Activity, alongside the six phase nodes. A phase overview links its relevant
Sessions without copying logs. This supersedes the earlier 05-Implementation/Sessions
proposal; no existing records are physically relocated by this design update.

Old SDP/Sessions documents stay in place and remain addressable. Link them by
identity/path; later relocation requires a checked migration preserving links and
history. A new chat, retry or compaction does not create a new Session. Instantiate
successive Sessions at meaningful work/handoff boundaries; plan future ones as
route slots. A Session has one primary Activity; unrelated findings become cards.

Track phase readiness independently of Activity lifecycle: not-started, in-progress,
ready, not-applicable. Ready means the phase's named output/disposition exists,
not that implementation or owner acceptance passed. Model phase applicability and
criteria explicitly; corrections may reopen a phase. Keep detailed milestone
state in its plan, not duplicated in the Activity descriptor. Derive phase readiness
where evidence suffices; authored applicability/waiver decisions are recorded facts.

## 4. Lifecycle and active-work rule

Canonical Activity states: planned, active, on-hold, completed, canceled.
The first four are always top-level navigation groups. Show Canceled when needed;
cancellation must not be counted as successful completion. Show legacy/unassigned
Sessions separately when present, without pretending they are Activities.

| From | Operation / next state | Required meaning |
| --- | --- | --- |
| absent | create → planned | Stable identity and need brief registered |
| planned | activate → active | Work selected; project has no other active Activity |
| active | hold → on-hold | Explicit reason, resume condition and actual execution disposition |
| on-hold | resume → active | Resume selected; active-slot and stale-input checks pass |
| active | complete → completed | Need-level acceptance/disposition and mandatory evidence recorded |
| planned / active / on-hold | cancel → canceled | Reason and disposition of linked work recorded |
| completed / canceled | reopen → planned | Explicit reason for the same need; prior history/evidence retained |

Do not reopen a record for a different need. A study may conclude no system change
is needed; completion then records that accepted outcome and not-applicable phases.
Completion does not require every broadly related card to be closed, but every
selected contribution must have a truthful disposition. Unmet mandatory scope
requires continued work or an explicit scope change; do not hide it by deferral.

Select **at most one active Activity per canonical project authority**. Activity
is active during study/design as well as implementation. A spanning Activity has
one coordinating project; other repositories contribute references instead of
independent copies competing for the same need. Parallel assignments inside the
selected Activity remain possible when authorized.

An explicit switch holds the previous Activity and activates the next as one
checked operation. No automatic switch because another folder was opened. PG must
prevent new assignments under a hold request and account for in-flight work before
claiming the switch completed. Preserve accepted output and recovery state; do not
equate ledger hold with killing a process. Pending switch/quiescence is operation
state, not a misleading second active Activity.

## 5. Activity ledger and state authority

Implement an **ActivityLedger domain** over the existing canonical
SDP/ProjectManagement/Ledger.ndjson. Each state change appends an Activity event
with date/time. This satisfies the Activity history requirement while preserving
one management-history writer/recovery boundary. Do not mirror authoritative
events into a second writable SDP/Activities ledger.

An Activity-only history/export is a filtered projection, generated on demand. If
a standalone activity-ledger file is exported, label it derived and record its
canonical ledger revision. It is never read as a competing state source.

The existing projecthistory package already provides exact-byte replay, expected
revision checks, idempotent append and cooperating transactions. It does **not**
implement Activity validation, authorization or cross-clone coordination. Add a
versioned Activity payload/domain validator; do not insert kind Activity into the
current management payload 0.2 enum or rewrite old events.

Proposed event family: x-activity:created, transitioned, updated, linked, adopted,
and corrected under the generic event envelope. Final schema versioning must be
checked against installed readers before adoption. Each lifecycle event records:

| Field | Meaning |
| --- | --- |
| eventId / operationId | Stable event and idempotent operation identities |
| subjectId / project | Activity and its canonical authority |
| occurredAt | Actual effective timestamp, RFC3339 with timezone; live writes use UTC |
| recordedAt | Time the event was recorded; distinct for historical adoption |
| actor / authority reference | Who requested/performed the transition and its actual source |
| previousEventId | Previous Activity event for causal validation |
| from / to | Exact validated lifecycle transition |
| reason / evidence references | Why state changed and the supporting disposition |
| descriptor revision | Exact relevant descriptor revision, preserving identity/path changes |

For absent historical timestamps, record unknown effective time in a dedicated
adoption payload; do not use the import timestamp as a fictitious start date.
Creation adopts explicit initial state only through the checked import path; it
does not invent earlier transitions. Metadata/link updates do not reset lifecycle
intervals. A correction references the affected event and retains its bytes; any
effective reinterpretation is explicit and revalidated against the whole history.

Replay determines group placement and active-slot ownership. Missing/invalid
history is a visible diagnostic, not an assumption that a folder is planned.
Expected project-ledger revision and per-Activity predecessor protect mutations.
History order/causal links determine replay; timestamps alone never order states.
Different branches with concurrent activation histories must be reconciled explicitly,
not concatenated by timestamp. Local file locks do not enforce global single-active
work across disconnected clones; only the selected canonical writer can publish
Activity transitions. Offline changes are proposals until checked against it.

Descriptor plus ledger mutations use the existing recovery/transaction mechanism.
An Activity switch uses one atomic domain operation containing both transitions;
replay must never expose a half-switch. Retried operations retain the same IDs and
bytes. Invalid tail, stale revisions, mismatched metadata and pending recovery
produce actionable diagnostics and no inferred success. Future CLI/MCP adapters
share this domain service and its authority checks.

## 6. Producer-generated Activities tab

```text
Activities                           [project lifecycle overview]
  Planned
    SDPTool Browsing
  Active
    Activity Management              [need, phase progress, next action, timeline]
      Mandate
      Study
      Requirements
      Design
        Architecture and design documents
        Models
        Blueprints
      Implementation
        Plans
      Verification
      Sessions
  On hold
  Completed
  Canceled                           [when nonempty]
  Unassigned sessions                [legacy records, when present]
```

The sketch illustrates future behavior; it does not register these Activities.
Fixed state/phase order comes from SDPTool. Activity order within a state uses
explicit priority if supplied, then deterministic title/identity order; display
last transition separately. Do not sort by filesystem modification time.

Activation behavior:

- Activities: project timeline and concise current/future work overview.
- State group: activities in that state, their need, phase, next action and holds.
- Activity: mandate summary, current state/phase, route, linked plans/cards,
  evidence/readiness and its event-derived timeline.
- Phase: readable section overview, direct documents, references and relevant
  work; folder expansion merely reveals children.
- Plan/Session/Blueprint/document: advertised existing reader/action in its context.

No raw IDs or complete ledger tables in ordinary pages. Human titles and readable
source links retain internal identity. An explicit History page shows event time,
transition, reason and source; machine JSON remains available for tools.

The producer publishes Activities instead of Sessions when the adopted capability
is supported. Unassigned legacy Sessions remain reachable in the new tab; old CLI
Session operations and deep links continue to work. Projects not yet adopted
retain their existing Sessions tab. Do not publish both tabs with duplicate lists
by default. XFMD's tab construction already consumes producer roots; verify its
protocol/action support before assuming this works without changes.

Navigation discovery reads descriptors/history and existing file/reference indexes;
it does not generate every document, timeline or viewpoint. Primary actions request
their content on demand. Large trees use advertised lazy expansion/paging with
stable IDs and no silent loss. Cross-project references switch through a supported
project action; inaccessible targets remain explicit links/diagnostics, not invented
local files or unsafe path escapes.

## 7. Project lifecycle timeline

One row per Activity. Creation is a point event followed by a planned interval;
activation/hold/resume create colored state intervals; completion/cancellation
are terminal markers. Reopening adds a new interval without erasing the earlier
closure. Reasons and source records are available on selection. Use text/patterns
as well as colors.

| Timeline rule | Required behavior |
| --- | --- |
| Time basis | Actual event timestamps; UTC storage, explicit display timezone |
| Current interval | Extend only to the snapshot's as-of time; label ongoing |
| Date/time precision | Retain exact timestamps; scale the visible axis without discarding events |
| Planning estimates | Optional separate forecast lane; never mix estimates into observed history |
| Unknown/imported history | Visible unknown segment/marker; no fabricated duration before first evidence |
| Clock anomalies | Surface inconsistent times; do not render negative durations or reorder causal events |
| Corrections | Default corrected interpretation with correction marker; original event history remains accessible |
| Duration | Calendar time in each state; active time is not effort or productive hours |
| Scope | This shows Activity lifecycle, not the complete project's unrecorded pre-adoption history |

### Observed event markers

Overlay point markers on the Activity row for attributed commits and ledger events.
Use distinct colors **and** shapes/legend labels, for example a blue solid tick
for a commit and an amber diamond/tick for a ledger event. Lifecycle transitions
remain identifiable as transition markers, not generic work events. Color choices
are a renderer/theme detail; their meanings must be consistent and accessible.

Markers show that an event was recorded, not how long a person or agent worked.
Keep markers visible even during planned/on-hold periods when discussion, capture
or historical imports legitimately occurred; they do not silently activate work.

- Attribution must follow explicit Activity/Session/assignment/evidence references
  or a checked commit association. Do not assign every project commit to whichever
  Activity happened to be active at its timestamp, or infer ownership from paths.
- A commit marker retains repository identity and immutable commit identity
  internally. Its label shows a readable subject and timestamp; use the recorded
  committer timestamp with its provenance, not an alleged work-start time. A
  governed completion event has its own event time; these need not match.
- Ledger markers distinguish event occurrence from later recording/import.
  Unknown effective time remains unknown; show recording time only with that label.
  Include relevant card/plan/Session-linked management events, not unrelated ledger
  traffic. Support filters for event kinds and linked work.
- Deduplicate the same commit referenced by several receipts. If a ledger event
  records that commit, correlate the markers/details rather than presenting two
  independent deliveries. A commit affecting multiple Activities may appear on
  each affected row with shared-event attribution.
- Selecting a marker shows subject/event meaning, time basis, reason and linked
  Session/card/plan/source. Link to available commit details through a supported
  action; do not assume an external viewer or expose machine IDs as reading labels.
- At coarse zoom, aggregate dense markers with a count and expandable event list;
  never silently discard events or change their timestamps. Missing Git history
  leaves a visible availability limitation, not a claim of no work.

SDPTool owns attribution, event selection and marker positions. The timeline
snapshot pins the ledger revision plus the selected repository history/evidence
inputs; later commits do not silently change a retained view. No new canonical
ledger event is needed merely to draw a marker for an already attributed commit.

SDPTool computes timeline data and an on-demand reader artifact from one ledger
snapshot. XFMD renders it, preserving filter/selection/history. A generated SVG
with exact-time data/companion event list is a viable baseline; native interactive
zoom/selection is an optional renderer capability, not a reason to recompute
lifecycle semantics in XFMD. Verify dense same-day events and print legibility.
Do not maintain a manual Gantt alongside the event-derived history.

## 8. ProjectGovernance, plans and model gates

PG resolves Activity selection in addition to the existing Session/plan/assignment
context for adopted workflows. Admission checks active Activity, current Session
membership, selected scope, prerequisite readiness and the existing design/review
rules. Operational run state does not overwrite Activity/card/plan state.

### Discussion is open; implementation is scoped

The active Activity limits implementation, not the owner's ability to discuss
other technical matters. Classify intended effects before admission:

| Request | Disposition |
| --- | --- |
| Technical question, wish, read-only investigation or discussion outside active scope | Allow; preserve the active Activity and identify unrelated context |
| Capture findings, create/update a KB card, draft a need brief or plan a future Activity | Allow explicit record-only work; do not activate or start its implementation |
| Product/code/model change, changes to accepted requirements/design, installation/publication or execution of an implementation assignment | Require selected active Activity, scope coverage and existing action authority |
| Necessary work exceeds active scope | Record the finding and return a scope-change/new-Activity route; do not silently expand or switch |
| Implementation with no selected Activity or against planned/on-hold/completed work | Refuse dispatch with a clear reason and the record/selection action needed |

Do not use capture/planning as a loophole to change an unrelated product or silently
replace accepted system requirements/design. The admission layer checks effects
and selected scope, not just words such as discussion in a prompt. A card link
alone does not authorize implementation. Apply the same boundary to owner-driven
and agent-proposed assignments through CLI/MCP/controller paths; a purely advisory
agent convention is not proof of enforcement. Recheck context at dispatch and
guarded writes, preserving pending-run/recovery behavior when Activity changes.

This boundary permits normal maintenance of the discussion record and future
planning without requiring every unrelated idea to interrupt the current Activity.
Read-only investigation may include isolated analysis/test artifacts; it must not
silently become a product change or installation under the label of investigation.

Activity completion is an explicit governed disposition using actual need-level
evidence. Completed plan/Session/agent work is necessary only where selected and
is not sufficient by itself. Blueprint and model identities/freshness remain owned
by existing services. A preliminary Blueprint can inform planning; implementation
and final model release follow the existing selected-target governance contract.
Software release/publication remains separately authorized.

Reference memberships carry contribution/whole-scope semantics for shared cards
and plans. Linking a card does not reopen or close it. Canceled Activities preserve
delivered artifacts and explicitly return remaining work to backlog/another need.
Dependencies use stable Activity references with checked cycles and explicit
satisfaction; changing a title/path cannot alter the dependency.

## 9. Adoption and compatibility

1. Specify/version Activity records, events and producer capabilities; update
   source process guidance, then distribute through the supported installer.
2. Add the read/replay service and preview an adoption mapping. Preserve original
   records and ledger bytes; unresolved membership stays Unassigned sessions.
3. Select a canonical writer and register the Activity using checked commands.
   Bootstrap this work with an explicit adoption event pointing to this design
   Session, not fictitious historical Activity events.
4. Link existing Sessions/plans/models where they are. New work may use the Activity
   folder layout. Physical migration is optional and separately verified.
5. Enable Activities navigation and PG admission for the adopted version. Old
   clients either use a documented compatible view or explain missing capability;
   they must not silently write unsupported lifecycle data.

Session discovery currently recognizes SDP/Sessions, and plan discovery has existing
path assumptions. Implementation must support Activity-local records through the
domain catalog as well as legacy paths. Adding folders alone does not implement
discovery, model binding, lifecycle or recovery. Preserve unrelated active work;
adoption does not automatically hold existing runs or choose their need ownership.

## 10. Delivery boundaries and acceptance

Proposed implementation increments, to turn into a typed plan after design/model
reconciliation:

| Increment | End-to-end outcome |
| --- | --- |
| A1 Domain and records | Create an Activity from a mandate; replay its canonical state and references with versioned schemas and compatibility checks |
| A2 Lifecycle and recovery | Activate/hold/resume/switch/close with single-active checks, idempotency, transaction recovery and truthful disposition |
| A3 Discovery and reader | See the Activities tab and six phases in XFMD; open current and legacy linked documents, Sessions/plans/Blueprints with producer-owned actions |
| A4 Timeline | Read actual lifecycle intervals and attributed commit/ledger markers, with timestamps, source provenance, dense-event handling and no effort inference |
| A5 Governance/adoption | PG allows unrelated discussion/capture but rejects out-of-scope implementation; adopt a real project with preserved history, root-level Sessions and no duplicated run state |
| A6 Integrated qualification | Independent review and actual native workflow prove the selected need, including installer/older-client behavior; explicitly dispose remaining scope |

Mandatory cases: two simultaneous activation attempts; stale descriptor/history;
interruption between descriptor and event publication; retried switch; active run
during hold; malformed/conflicting history; missing/duplicate descriptors; one
card contributing to multiple needs; no-code study outcome; canceled/reopened work;
legacy Sessions outside Activities; external unavailable references; large/dense
timeline; unknown historical time; clock reversal; changed reviewed Blueprint;
manual/old writers; and no lifecycle inference from directory names or agent success.

Also verify: standard README contains no Activity-specific narrative/state; mandate
is the sole authored need/outcome narrative; Sessions is a root sibling of phases;
unrelated ideas can become KB cards without switching Activity; planned future work
cannot dispatch; an unrelated commit during an active interval is not attributed;
one commit referenced by several events is correlated; absent Git/unknown times
are explicit; marker filters and coarse-zoom aggregation preserve all source events.

No implementation/test/installation claim is made here. Candidate decisions still
to validate: exact schema compatibility, supported renderer output and platform
writer support. This document recommends the conceptual ownership and behavior;
the next design step binds them into SDL WORK and an executable implementation plan.

## Activity0002 execution authority

Owner accepted Activity0001 and selected Activity0002. PLAN-SDP-0027 now owns
implementation and qualification. Earlier design-only/on-hold statements above
record their original scope; they no longer restrict the authorized execution.
The existing 0.1 manual records remain valid inputs. Exact implementation contracts
are refined in Activities/Activity-0002--Activity-Management/04-Design while current
requirements and SDL remain in their existing authoritative homes.
