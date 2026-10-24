# PGD1 — Pilot acceptance and requirement mapping

These are executable acceptance specifications for the implementation handoff,
not passing runtime tests. [Design](Design.md) owns behavior; cards and RGS2 own
the source requirements. Use an isolated project fixture with one real bounded
Maintenance change, existing work records, and a deliberately unrelated dirty file.

| Case | Governing source | Trigger and required observable result |
| --- | --- | --- |
| PG-A01 Direct answer | KB036; RGS2 explain versus act | Submit status/unrelated question: answer through the read route; no new card, plan or run |
| PG-A02 Existing work | KB036/038 | Submit covered change: resolve current card/plan/authority once; show routine and exact assignment; no duplicate plan or approval request |
| PG-A03 Missing or ambiguous route | KB036/037 | Remove routine, make it unreadable, incompatible or conflicting: four distinct dispositions; dependent mutation rejected; repeat does not create duplicate gap cards |
| PG-A04 Bounded execution | KB038; RGS2 role model | Fresh Master receives governing context, owned paths and return conditions; preserve unrelated edits; observed child identities and context rules match the claimed role mapping |
| PG-A05 Scope change | KB036/037 | Steer with an out-of-scope addition during work: retain original goal, record delta, stop dependent progress, reconcile in-flight effects and revise coverage |
| PG-A06 Evidence and rework | KB037 | Completion without required evidence refuses; review finding returns for rework; self-check cannot satisfy required independent review |
| PG-A07 Candidate changes | KB037; Traceability | Modify relevant source after verification: evidence remains historical but cannot satisfy the current gate; unrelated excluded file change alone must not invalidate scoped evidence |
| PG-A08 Core retry/crash | KB037 | Repeat same operation/key before and after simulated crash: same receipt/revision; changed arguments with same key refuse; orphan generation does not advance head |
| PG-A09 Record publication recovery | KB037/042 | Crash between journal/card/ledger writes: reconcile exactly once; preserve original ledger prefix; changed external preimage stops with explicit conflict |
| PG-A10 Run version and authority | KB036/037 | Upgrade routine or submit forged owner role: old run stays pinned; forged input cannot grant rights; existing authorization still works within scope |
| PG-A11 Observer and controller | KB038 | Close observer: work continues. Lose controller transport: state becomes uncertain; reconnect to mapped work without duplicate turn submission |
| PG-A12 Conversation capture | KB038/042 | Multiple owner inputs, steering, commentary, final and interruption: retain observed distinctions and IDs; no fabricated final response; replay completed items without duplication |
| PG-A13 Session continuity | KB042 | Resume on a new thread/attempt: same Session/assignment, new attempt; current next step and source revisions recovered; authored text and old entries preserved |
| PG-A14 Honest timing | KB043 input | Clock reset or disconnected interval: epoch/gap explicit; no invented active time or live cursor; static output exposes as-of position |
| PG-A15 Compatibility and auth | KB038; RGS2 | Pinned app-server and chosen MCP SDK negotiate in a real local probe; managed account path exercised without API-key fallback; unsupported capabilities refuse explicitly |
| PG-A16 Return disposition | KB037/038/042 | Completed model turn alone does not close card; exact candidate/test/review result and remaining work are recorded; owner-reserved acceptance requires actual provenance |
| PG-A17 Model boundary | Owner ProjectGovernance scope; KB049/050 | Maintenance with recorded no-design-impact rationale works without model store/blueprints; design changes require WORK under PG-A21; optional unavailable model reference stays explicit; no model artifact is mutated by process bookkeeping |
| PG-A18 Bounded protocol | Existing SDPTool strict input conventions | Duplicate/unknown fields, escaping paths, oversized frames and invalid cursors refuse without corrupting state; finite page limits enforced |

## T039 required governance cases

Owner requirement, [mandatory governance contract](Orchestration.md#t039--mandatory-governance-admission-and-lifecycle-supervision).
These are acceptance specifications, not passing tests or delivered guards.

| Case | Governing source | Trigger and required observable result |
| --- | --- | --- |
| PG-A19 Goal fit and new work | KB052; owner T039 | During an active Session submit a covered request, coherent scope revision, logically unrelated feature in the same repository, and pure status question. Reuse covered records; revise coherent scope before dependent work; explain and establish separate KB card(s), Session and selected plan for the unrelated goal before assignment, preserving current work; status creates no new hierarchy. Retry creates no duplicate records. |
| PG-A20 Assignment chain and supervision | KB052; owner T039 | Remove/mismatch a card, Session step, selected plan, revision or authority before launch/resume; dependent action refuses. SAD/Master/child assignments have exact bounded coverage, children reuse the parent chain. Steering/report/resume triggers fresh checks; scope drift cannot ride through a previously valid admission. |
| PG-A21 Design and Blueprint decisions | KB052; owner T039; PG-A17 | Prose-only task may justify no WORK. A small bug/refactor that changes behavior/contracts requires SDL WORK despite its label; absent/inadequate model blocks dependent implementation until prepared. Evaluate Blueprint need separately: required context must exist at exact revisions; justified omission never waives WORK. Changed TARGET/task invalidates handoff before continuation. |
| PG-A22 Model, code and plan coordination | KB052; owner T039; handoff 1/2 | Two disjoint Sessions integrate compatible WORKs from known ancestry; overlapping model element/contract, code path or shared plan/ledger writer is serialized or refused with a scoped reason. Include dirty work, textual-clean semantic conflict, missing ancestry and duplicate IDs. Retain sources and history; validate combined model/code, regenerate affected views and refresh stale review. Local claims never imply unobserved cross-clone exclusion. |
| PG-A23 Model release before closeout | KB052; owner T039; PG-A16 | Finish implementation with unreleased WORK: implementation progress may publish but whole-assignment completion refuses. Validate merged WORK, freeze exact CANDIDATE, obtain required actual acceptance and promote RELEASE with code evidence/predecessor. Changed candidate, competing accepted head, pending conflicts or model-only evidence cannot satisfy implementation completion. Interrupted promotion/closeout recovers without duplicate release/events; explicit owner transfer retains successor and incomplete scope. |
| PG-A24 Mandatory skill and observable checks | KB052; owner T039 | Start/resume Project Leader with missing/unreadable/incompatible required skill: dependent dispatch refuses. Native evidence identifies the supplied skill version, host acknowledgment where available and actual admission/supervision/closeout decisions. A discoverable file or claimed skill use alone cannot pass the case; skill instructions alone do not prove runtime guards. |

## Evidence ladder

1. Pure-core state/fixture tests establish guards, persistence and replay only.
2. Adapter conformance tests establish equivalent results through CLI and MCP;
   record the selected protocol/SDK versions and actual negotiation.
3. A real controller/app-server trial establishes thread/item/auth behavior and
   reconnect coverage. Schema export alone cannot satisfy this level.
4. One complete Maintenance task establishes useful owner-to-evidence behavior.
   Independent review uses a context independent of implementation, the governing
   intent and exact candidate. If host provenance is unavailable, record the gap.

For every result record candidate revision/digests, command/environment, inputs,
observed outputs, verdict and limitations. Do not label simulated protocol fixtures
as a live Codex trial or an agent's declared skill use as measured compliance.
PG-A15 requires real account/host conditions; inability to run it is an explicit
remaining integration condition, not permission to mark the entire pilot complete.

## Fixture preservation and review

Use pre/post digests for unrelated files, original ledger prefix and authored Session
content. Include negative/recovery paths alongside the successful workflow. Reuse
project management validators for publication output. The full 16-family catalog,
general software-publication routines and timeline renderer remain deferred after a
successful pilot; T039 model-release supervision is part of the selected workflow.
