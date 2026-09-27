# Request classification and mandatory routine selection

| Field | Value |
| --- | --- |
| id | KB-SDP-036 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-09-27T23:51:37.811184+00:00 |
| source | Owner conversation 2026-09-28: mandatory routines for every request and scope change |
| next_review | Before selecting the next code or process implementation after RP3 |
| tags | workflow, governance, skills, request-routing, scope-control |
| Systems | SDP, SDPTOOL |

## Need and owner direction

Established routines become ineffective when agents and owners must remember to
invoke them. Enthusiasm, interrupted attention, context changes and scope drift
can bypass existing KanBan, planning and review work or cause another competing
routine to be invented. The process must carry this memory rather than depend on
the owner recalling the correct procedure in each conversation.

The owner requests a mandatory entry routine for every request: classify intent,
look for an applicable established routine and identify the required path before
acting. Apply this again when a new message changes scope during ongoing work.
Changes to SDL/SDUI models and product code must be covered by KanBan and a plan.
Process/template/installer cleanup must have proportionate Maintenance planning.
The owner asks the agent to be more consistent about these constraints, including
briefly explaining the selected category and routine when action is involved.

This card captures the direction and proposed acceptance criteria. It does not
claim a new skill, enforced hook or released workflow already exists.

## Existing authority and the concrete recent example

[The SDP entrypoint](../../../Skills/sdp/SKILL.md) already routes role selection;
[KanBan](../README.md), [Plans](../../ProjectManagement/Plans.md) and the role
skills already define much of the work. Prefer strengthening this shared entry
and referring to existing procedures over adding a competing workflow hierarchy.

The recent template consolidation was recorded under
[MAINT-SDP-0009 / TS1](../../Maintenance/TS1/Plan.md). Go-only removal and release
were recorded in [MAINT-SDP-0010 / RP3](../../Maintenance/RP3/Plan.md), including
the owner's scope and version corrections. The gap is not an absence of these
plans: the request-to-routine/plan connection needs to be visible, repeatable and
checked when scope changes. Do not rewrite this history as undocumented work.

## Proposed first category set

These are candidates for study and trial, not newly adopted process rules.
Categories may combine; the resulting obligations must compose without creating
one new plan or card for every message.

| Request category | Candidate existing route |
| --- | --- |
| Explanation or status | Read current facts and answer; no mutation or unnecessary plan |
| Idea or new requirement | Find related cards, capture owner intent and reconcile overlap |
| Study or unresolved decision | Existing study/architecture route; bound unknowns and deliver a decision basis |
| Planning or reprioritization | Planning skill; select/revise the linked plan, optional Scrum/Sprint |
| Implementation or defect correction | Applicable card and plan before model/code mutation; analysis where the cause is unknown |
| Maintenance, migration or cleanup | MaintenancePlan covering affected files, preservation and verification |
| Review or verification | Appropriate independent review/evidence route; findings return to tracked work |
| Integration, release, installation or upgrade | Existing release/install routine, exact candidate and authorization boundaries |
| Continuation, correction, pause or cancellation | Reuse the active assignment; check scope, evidence and state changes |

## Proposed routine contract

1. Classify the request using current project authority and active work. Several
   categories or a non-project/out-of-scope result are valid. Do not force every
   question into an implementation workflow.
2. Locate the authoritative matching routine and current version/status. Reuse a
   matching active card/plan when it covers the work; avoid duplicate proposals.
3. State the route briefly for substantive work: category, routine, card/plan,
   next permitted step and any material gap. Do not require the owner to remember
   or supply these identifiers. Routine selection is not another approval prompt.
4. Before mutation, check prerequisites and record coverage. For code/model work,
   select the card and plan first. For a scope change, assess whether the current
   plan can be revised or separate work is needed; do not silently continue under
   stale scope or abandon the previous objective.
5. If no routine fits, identify the gap and register it. Do not invent a durable
   replacement routine implicitly. Clarify only the material unresolved decision;
   continue independent already-authorized work where possible.
6. Verify the applicable completion conditions and update existing records before
   saying done. Reclassify at milestones and new substantive input, not only at
   session startup. Reuse loaded unchanged instructions where practical.

Keep a routine record small: identity, purpose, matching criteria, required
inputs, current authority, allowed steps, completion evidence, exception handling
and supersession. Link existing skills/plans instead of duplicating them. Decide
where this registry belongs during planning; no new directory is assumed here.

## Activation and enforcement questions

- Determine whether to extend the mandatory SDP entrypoint or introduce an
  sdp-routine skill beneath it. A discoverable SKILL.md alone is not proof that
  an agent actually loaded or followed it.
- Separate host/session instruction activation, agent instruction compliance and
  programmatic checks. Establish what each supported host can actually enforce;
  do not promise that prose can intercept every request or prevent every write.
- Consider a later SDPTool/MCP preflight that returns the selected routine,
  linked work and missing prerequisites. A machine-checkable transition can
  strengthen controls but does not guarantee that arbitrary shell/file writes
  elsewhere are intercepted. Treat this as a separately bounded implementation.
- Cover compound requests, ambiguous categories, existing authorization, scope
  changes, missing routines, resumption after compaction and competing versions.
- Preserve proportionate work: direct factual questions should stay easy; a
  matching active plan should not be recreated, and explicit owner authorization
  should not trigger repeated permission requests.

## Next action and completion criteria

Select a proportionate MaintenancePlan to inventory current routes and design a
single request-entry contract, then trial it before changing distributed skills
and templates. This is proposed next work, not started implementation.

Acceptance should demonstrate representative requests from the table, including
"also remove PowerShell" during release preparation, template restructuring,
"continue", a pure status question and an unrelated factual request. Evidence
must show when routing was performed, which authority/card/plan was reused or
revised, what gaps blocked dependent mutations and how scope/closeout was kept
current. A declarative instruction or file-presence check alone is insufficient.

Related: [KB-SDP-029: typed plans](../completed/%23029--Proposal--Typed-plans-and-planning-skill.md),
[KB-SDP-027: skill activation](../completed/%23027--Study--Skills-review-and-project-activation.md),
[KB-SDP-031: assignment context](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md).
No reopening of completed cards or MCP implementation is implied by this link.

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-27T23:51:37.811184+00:00 | codex; EVT-KB-SDP-000207 | Captured owner direction and proposed routine/category/activation contract during RP3 closeout | Registration only; select the MaintenancePlan before implementation |
