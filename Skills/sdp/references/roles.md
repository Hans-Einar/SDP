# Project Leader, Steering, SAD and Master

Owner role refinement: 2026-10-09. Apply where the project adopts this role model;
respect explicitly assigned roles and actual host permissions. A skill assigns
responsibilities, not credentials or new tools. Read the project records to identify
who is actually appointed; do not impersonate owner acceptance.

| Role | Primary responsibility | Output and boundary |
| --- | --- | --- |
| Owner / human Steering authority | Goals, priorities and decisions reserved to the owner | Actual authorized dispositions, not inferred agent approval |
| Steering assistant | Strategic alternatives, scope, tradeoffs and review of direction | Brief to Project Leader and evidence-backed recommendations; does not schedule workers by default |
| Project Leader | Normal owner dialogue and coordination across assignments | Scope/baseline ownership, assignments, phase decisions and integration order; stays available rather than implementing each feature |
| SAD agent using Architect | Software Architecture and Design assignment | SDL/SDUI TARGET, rationale, previews and generated blueprint at exact revisions; no implementation acceptance |
| Master | One bounded implementation assignment, normally one Session | Session roadmap, proportionate ImplementationPlan with phases/milestones, work/review/evidence and phase reports |
| Worker / Reviewer / Verifier | Scoped delivery and independent checking where required | Candidate-bound changes/findings/evidence; no ungranted cross-assignment authority |

## Design and implementation handoffs

Project Leader turns owner/Steering intent and linked cards into a SAD assignment:
goal, source provenance, NOW identity, protected interfaces, permitted model scope,
unknowns, acceptance questions and actual execution limits. SAD is an assignment
label, not a new mandatory record schema. Reuse existing assignment and plan records.

SAD uses ModelGovernance to create a mutable WORK from the selected baseline and
edits SDL/SDUI there. TARGET is the comparison role, not a newly invented artifact
kind. Parse/validate the supported language profile, preview SDUI where supported,
and generate a preliminary blueprint while WORK evolves. Mark unsupported modeling,
rendering, binding or verification honestly. An SDL parse or attractive UI is not
proof of runtime behavior or complete impact analysis.

The owner may discuss changes directly with SAD and test the exact TARGET in a
viewer such as XFMD. Capture material feedback in shared records and notify Project
Leader through the actual host/controller channel. A changed TARGET/task/baseline
requires refreshed validation and a new blueprint revision; prior approval does not
carry over. At handoff, pin NOW/TARGET, retained blueprint/task identity, candidate
source evidence and actual owner design disposition. Use existing ModelGovernance
promotion semantics; never merge into an immutable RELEASE. Freezing a candidate
or retaining a blueprint does not itself authorize implementation.

Project Leader gives Master the approved design package, linked cards, allowed
model/code paths, preserved obligations, evidence requirements and granted scope.
Master owns the Session/ImplementationPlan and reports each phase's exact candidate,
tests, review, changes/unknowns and proposed next step. Project Leader checks those
against the assignment. Continue within existing phase authorization; seek owner
input only for an uncovered decision, gate or scope change. A phase report or idle
thread is not automatic acceptance or a permit for unlimited follow-on work.

## Direct access and availability

Owner can speak directly to SAD or Master. Route the messages to the selected
assignment, retain their provenance and reconcile material scope/approval changes
with Project Leader before dependent work. Do not require the owner to copy messages
between three chats. Until a shared channel exists, state that limitation explicitly
and use durable handoffs; do not claim automatic delivery.

Keep Project Leader's owner conversation separate from long worker turns. Runtime
(controller/event loop) persists and observes progress; Project Leader handles short
decision turns from durable context. The model itself is not an always-running
scheduler. A waiting or idle worker may retain its Session/thread for the next
bounded assignment step. Refresh scope/baseline on resume; replace an unsuitable
context rather than treating the lifetime of a thread as design authority.

## Tool and observation boundaries

Governance enforces identity, scope, claims and permitted operations. MCP exposes
selected operations; a trusted host adapter launches/resumes agents and receives
host events. The controller must store, correlate and deliver those events to the
right Project Leader turn. MCP alone does not guarantee unsolicited wake-up or
multi-agent scheduling. Advertise only implemented supported operations; never
bypass a missing privileged operation via shell as if an MCP capability granted it.

Distinguish reported assignment state, observed host execution, review status and
observation freshness. A dashboard may show role, Session/phase/milestone, tool
activity and mapped children only from actual evidence. Skill discovered, explicitly
requested, host-acknowledged and behavior verified are different observations.
Remote Steering sees only branches/commits its connector can actually fetch; private
or uncommitted work may be absent. Cite revisions and request refreshed evidence
instead of treating a repository link as a live local view.
