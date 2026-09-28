# From requests to governed work: routines, MCP and observable execution

| Field | Value |
| --- | --- |
| Study | RGS1 |
| Date | 2026-09-28 |
| Plan | PLAN-SDP-0006 |
| Source cards | KB-SDP-036, KB-SDP-037 |
| Status | Research recommendation; production design and adoption not selected |
| Scope | SDP, XFMD, Ponsse, HSX, agro-crm, TerrainAnalyzer, GrassPhenology, processor emulator/debugger |

The missing capability is a durable connection between a request, the applicable
procedure, its authority, the selected work and evidence of what actually
happened. A documentation MCP helps retrieve instructions. A governed workflow
also needs executable transition rules, explicit uncertainty and recovery across
sessions. SDL blueprints address the worker's design context; they cannot alone
decide whether the right work was selected or whether a release was authorized.

**Recommendation:** extend the Go SDPTool with a small routine service used by
CLI and MCP, strengthen the existing SDP entry skill, and add a read-only terminal
observer. Begin with one complete maintenance workflow and its missing-routine
fallback. Grow a versioned base catalog through observed project use. Do not
build a new portfolio orchestrator or migrate every project before the pilot.

This is a study, not implementation or an adopted replacement process. The
[plan](Plan.md) records authority. [Evidence](Evidence.md) separates observed facts
from project assertions and owner reports. [Routine catalog](Routine-Catalog.md)
is the candidate coverage map; [scenarios](Scenarios.md) challenge the proposal
and define future acceptance cases.

## 1. What the portfolio teaches us

The owner reports repeated implementation drift and difficulty recovering project
context after interruptions. The sampled repositories contain substantial
instructions and careful reviews already. Their existence has not made current
authority, compatible candidates or the next permitted step easy to recover.
This suggests an operational gap rather than simply a shortage of documentation.

| Evidence | Observed pattern | Consequence for the proposal |
| --- | --- | --- |
| E01, E02 | SDP installation, project-specific rules and actual application adoption are different facts | Resolve installed profile and current project authority separately |
| E03, E04 | Multiple projects/branches have distinct handoffs and stop conditions | Bind work to project, checkout and assignment; never equate default branch with current work |
| E05 | A local update can affect other consumers and derived views | Analyze impact and invariants before assigning the implementation |
| E06, E08 | Stored PR status contradicts current GitHub observations | Derive external state, preserve historical decisions, expose reconciliation work |
| E07 | A review can complete while its recommended fixes remain unselected | Completion is scoped; a finding does not grant new implementation authority |
| E08 | Fake-runtime tests and compatible real-runtime evidence answer different questions | Track evidence level and paired provider/consumer revisions |
| E09 | Prior studies already contain useful workflow ideas | Reuse intent and lessons without reviving superseded Issue-first authority |
| E10 | MCP and host hooks expose useful but bounded interfaces | Measure actual integration coverage; do not promise universal enforcement |

These findings support the proposed controls. They do not prove that this design
would have prevented every historical failure. In particular, no process engine
can infer all unmodeled dependencies or establish correct product behavior from
an agent's completion message.

## 2. Cover the work before the work package

The process must support a conversation becoming useful work without forcing
every message through a long ceremony. A request may ask a question, correct an
active scope, select a decision or introduce an idea that should remain in
backlog. These are different outcomes. Most messages should reuse current work.

The full lifecycle is a graph with loops, optional paths and explicit stops:

1. **Recover context:** resolve the project and installed capabilities, selected
   checkout, active work, governing decisions and existing authorization.
2. **Understand intent:** distinguish explanation, observation, idea, defect,
   desired outcome, constraint and instruction. Record material ambiguity.
3. **Relate it to existing work:** find overlapping cards, known decisions,
   current plans and abandoned/superseded alternatives. Preserve the original
   objective when a new message merely steers it.
4. **Assess project fit:** identify users, expected value, affected capabilities,
   dependencies, unknowns, ownership and integration consequences. Compare an
   isolated fix with a broader design problem when evidence warrants it.
5. **Select work:** prioritize against active commitments, readiness and available
   review capacity. A Scrum or Sprint is optional. Capture an idea without
   activating it; make deferral and cancellation valid results.
6. **Study or decide where needed:** investigate uncertain causes, model gaps and
   alternatives. Record actual owner decisions, not an agent's imagined approval.
7. **Plan proportionately:** reuse or revise an existing typed plan; select a new
   one only when needed. Define outcome, boundaries, phases/milestones, evidence,
   branch/commit policy and stop conditions before implementation.
8. **Prepare the assignment:** combine selected intent with current design facts,
   observed code context, NOW/TARGET difference, invariants and unknown impact.
9. **Execute and verify:** follow the selected milestone, capture emerging scope,
   preserve unrelated work and attach candidate-specific evidence.
10. **Review and integrate:** distinguish self-check, independent review, owner
    acceptance and integration authorization. Rework invalidates affected evidence.
11. **Release or install when selected:** verify the exact candidate, publish only
    within existing authorization, reconcile actual remote/install results.
12. **Close and learn:** reconcile records, retain unresolved work, capture a
    procedure gap when justified, and make resumption understandable tomorrow.

Project creation/adoption, reverse engineering, incident recovery, backlog
consolidation, design-only work and decommissioning also fit this lifecycle.
It does not imply that every request executes all twelve steps.

### Assignment readiness is an explicit boundary

Before a worker changes a model or product, the coordinator should be able to
answer: What outcome is selected? Which card and plan cover it? What is current,
what should change, what must remain true, and how will we know? Which surrounding
units, consumers, channels and data contracts may be affected? What is unknown?
Who can decide an exception, and what must stop the worker?

The ready check needs:

- Selected outcome, constraints and scope; a project/worktree/candidate identity.
- Card, plan and milestone references, owned paths and explicit non-goals.
- Current authority and applicable procedure version; existing authorization.
- NOW and TARGET context, with added, changed **and removed** relationships.
- Integration invariants, compatible external revisions where relevant and a
  stated frontier beyond which impact has not been established.
- Verification and review expectations at the correct level; stop/recovery rules.

[BP2](../../04--Design/SDPTool/Blueprints/Plan.md) remains the separate blueprint
contract/pilot. A routine run references its bundle and revision; it does not
invent a competing bundle schema. Initially use manually bounded context plus
facts supported by the current SDL parser. The unsupported MVP1 experiment must
be labeled partial input, not presented as a complete executable system model.
A source tag links code and design; it does not prove their behavioral agreement.

## 3. Request routing and the missing-routine path

The agent interprets intent; a deterministic service checks declared conditions.
Neither replaces the other. Keyword similarity can suggest a routine but cannot
grant authority. Match on intended action, affected artifacts, current work and
project constraints, not just words such as “fix” or “release”.

Use a small set of action categories and orthogonal facets: system, lifecycle
stage, artifact kind, external effect, domain restrictions and active-work relation.
For example, “update the install template and release” combines Maintenance with
Release and Install facets; it should not create three competing assignments.
A read-only status question during implementation stays a question.

Every substantive route should expose a compact explanation: **request category,
selected routine/version, linked work, next permitted step and unresolved gap**.
For trivial questions the route can be `no-project-work`; do not create cards,
plans or journals full of everyday conversation merely to make a dashboard busy.

| Resolution | Meaning | Permitted response |
| --- | --- | --- |
| Matched | Applicable current procedure with satisfied entry conditions | Reuse/select covered work and proceed |
| No procedure needed | Factual or unrelated request with no governed project action | Answer directly; no backlog item |
| Missing | A project action needs procedure coverage that does not exist | Enter bounded gap handling; record why |
| Ambiguous | Several routes materially disagree | Gather discriminating evidence; ask only the material unresolved question |
| Unavailable | A known procedure or service cannot be read | Report unavailable, not missing; use a declared fallback if permitted |
| Incompatible | Installed schema/capability cannot execute this procedure | Identify the mismatch; do not silently run a newer definition |
| Conflicting/outdated | Sources disagree or current facts invalidate assumptions | Reconcile authority; block only dependent work |
| Deviation discovered | Work has left the selected scope or repeated exceptions expose a gap | Reassess affected steps and capture a linked gap |

### How a gap becomes a procedure

The following is a proposed meta-routine, not authority for automatic process edits:

1. Capture the request, intended effect, catalog revision, candidates considered,
   unmet conditions and current work. Store a concise rationale, not unrestricted
   transcripts, secrets or model reasoning traces.
2. Search existing procedure-gap cards and decisions. Add evidence to an existing
   card when this is the same gap; do not create one card per repeated prompt.
3. Continue permitted read-only investigation and independent authorized work.
   Hold the dependent mutation if its scope or authority cannot be established.
   Say exactly what is missing and which rule requires the hold.
4. Recommend a bounded disposition: extend a routine, add a project specialization,
   correct metadata, use a recorded exception, or decide that no routine is needed.
5. A selected plan authors the candidate. Trial it against positive and negative
   cases, review its obligations, then adopt/distribute a version explicitly.
6. Preserve which version a trial used. Retire or merge overlapping procedures;
   keep supersession links. Never turn repeated agent behavior into normative
   policy just because it happened several times.

Signals include repeated manual rescue, recurrent missing prerequisites, repeated
exceptions, near-identical work with inconsistent sequencing, recurring release
reconciliation and incompatible installed rules. These are reasons to investigate,
not automatic proof that another routine is useful. An infrequent action can need
a routine when its effect is consequential; frequency alone is insufficient.

## 4. Authority, state and evidence

Separate these dimensions instead of inventing another universal “status” field:

| Dimension | Authority | Example |
| --- | --- | --- |
| Intended work and permission | Owner instructions and adopted project decisions | Implement a milestone; release remains unselected |
| Management lifecycle | Existing KanBan/plan records and ProjectManagement ledger | Card in progress; plan active |
| Operational execution | Routine instance and its event history | Verification step waiting for a tool result |
| System conformance | Candidate-specific evidence and Traceability | Contract test passed on a particular candidate |
| External observation | Versioned observation with source/time | PR merged; installation receipt present |
| Agent activity | Session/lease/last observation | Worker connected, disconnected or unknown |

A merged PR does not prove owner acceptance. A live worker does not imply useful
progress. A disconnected worker is not automatically failed. A completed study
does not complete its capability card. Current management records remain the
place for material lifecycle dispositions, not every low-level tool event.

Proposed hierarchy: project → request → routine instance → step → child assignment
or nested routine. Multiple requests can steer one instance; one request may
create related bounded instances. Maintain stable identities and causal links,
not one global state per repository. Support parallel children only when the
host and assignment authorize delegation; selecting a routine cannot override
that rule or let a worker impersonate its independent reviewer.

Candidate step states: pending, ready, running, waiting, blocked, failed, completed,
skipped and canceled. A waiting reason distinguishes input, review and dependency.
These are proposed operational states, **not new CardState values**. “Skipped”
records permitted authority/reason and never masquerades as passed verification.

A transition identifies instance, expected revision, step, requested event,
idempotency key, actor role and evidence/authorization references. The service
checks prerequisites, role and evidence applicability. Reject stale concurrent
updates and conflicting replay; return the existing result for an identical retry.
A worker cannot satisfy a required owner decision with its own statement.

Evidence must identify its candidate and inputs: commits alone are insufficient
for dirty trees, generated payloads or paired repositories. Record relevant file
hashes, invocation, result, environment and evidence level. Invalidate affected
checks when inputs change; retain unrelated valid evidence with an explicit basis.
Unknown dependency coverage requires conservative invalidation, not invented
precision. Preserve historical results instead of rewriting them as current.

### Recovery and ownership of records

Recommend a local durable service/core first. Evaluate an event journal plus
transactional state store, such as SQLite, during design. No database choice is
adopted here. The operational journal is not a second KanBan/Traceability ledger:
material transitions project into the existing records with common operation IDs.
Use an outbox/reconciliation mechanism or equivalent recoverable transaction so a
crash between state persistence and file updates cannot manufacture completion.

A run pins its project identity and checkout identity separately. Two worktrees
share a logical project but may have incompatible candidates. Concurrent workers
need bounded ownership and revision checks; stale leases become unknown/recoverable,
not silently reallocated. Restart recovers the persisted run, then refreshes
observations before allowing affected work. Preserve explicit pause/cancel authority.

Sensitive inputs remain local and access-controlled by default. Store minimum
necessary request/evidence metadata; expose summaries to the observer. A service
using a network transport needs a separately verified authentication and access
boundary. Files or tool responses are data, not instructions allowed to override
the installed procedure.

## 5. Interfaces: knowledge, operations, host entry and observation

The [OpenAI Docs MCP](https://developers.openai.com/learn/docs-mcp) demonstrates
read-only documentation discovery. For SDP, expose version-resolved authority,
procedures, selected work and model context as resources. Search should return
provenance, installed/profile version and current-versus-historical status.
Do not serve “latest upstream” instructions as though they govern an older project.

A second interface exposes bounded operations through the same SDPTool core used
by CLI. Illustrative API responsibilities, **not implemented command names**:

| Operation family | Responsibility |
| --- | --- |
| resolve context / read authority | Resolve explicit project/worktree, installed capabilities and active work |
| inspect routes | Return candidate procedures, applicability basis and unmet conditions |
| open or resume run | Link authorized work and pin the procedure/context revision |
| check / apply transition | Validate preconditions, authority, revision and evidence before changing state |
| record observation / gap | Attach a fact or procedure gap without treating it as an approval |
| snapshot / events after cursor | Support recovery, read-only monitoring and replay |

Use the [official Go SDK](https://go.sdk.modelcontextprotocol.io/) as a candidate
adapter dependency, not a separate implementation of business rules. Existing
SDPTool discovery, navigation and install libraries should remain reusable. Do not
wrap every shell command as a new MCP tool. MCP transport success is not task
success, and tool descriptions are not a security boundary.

The entry skill provides semantic interpretation and routing guidance. A host
adapter supplies request entry/resumption and, where supported, mutation checks.
MCP alone does not observe every chat turn. The documented
[Codex hooks](https://learn.chatgpt.com/docs/hooks) justify a bounded integration
experiment, not a claim that this session already has enforcement. Test the
actual host version, tool surfaces and bypasses before naming anything “guarded”.

Describe coverage honestly:

- **Documented:** instructions and procedure definitions exist.
- **Instrumented:** actual requests/actions produce attributable events.
- **Guarded for tested operations:** specified operations are rejected when their
  prerequisites fail, on named host/tool versions.
- **Isolated execution:** stronger controls require an explicitly designed
  execution boundary; arbitrary shell effects cannot be inferred from a regex.

The core can reject invalid calls through itself. An agent with unrestricted
filesystem/shell access may bypass it. A read-only docs MCP, skill or visualization
does not remove that limitation. Host integration should fail visibly and apply a
declared fallback for the affected action, without freezing unrelated work or
silently treating service absence as “all checks passed”.

### The owner's monitor

Start with a second Kitty terminal showing a stable textual tree and current
snapshot, refreshed by polling if sufficient. Show project, request, route/version,
card/plan/milestone, current step, next permitted action, waiting reason and evidence.
Allow drill-down into nested routines/worker assignments and back via breadcrumbs.
Separate live mode from historical replay. Keep labels/icons as well as colors.

A later stable flowchart may highlight actual nodes and edges, with parallel work
and stale observations visible. Reuse rendering components when justified; do not
make a new GUI, Mermaid enhancement or XFMD application change a pilot prerequisite.
Closing the monitor changes no work state. Reconnecting resumes after an event
cursor or reloads a snapshot. A checklist is a projection of state/evidence, never
an animation driven solely by agent prose. Do not display invented percentages.

## 6. Distribution and project specialization

Many routines should ship with SDP. Recommend a versioned base catalog alongside
canonical skills, inventoried by the existing signed release payload. The exact
source/installed paths are a DesignPlan decision; no new root hierarchy is adopted
here. Reuse the Go installation/upgrade machinery rather than introduce another
installer or PowerShell entry point.

Each definition needs stable ID, version/content hash, schema version, purpose,
match/exclusion conditions, required capabilities, authority references,
prerequisites, steps/transitions, evidence, exception/stop rules and supersession.
Procedure version, catalog version, SDP distribution version and product version
are separate identities. A project receipt resolves the exact installed catalog.

Project-owned overlays can add domain obligations: Ponsse offline evidence,
XFMD viewer/application tests, emulator/DAP compatibility or data-migration checks.
An explicit extends/override relation identifies the base version. Refinements
must not silently weaken mandatory base conditions; conflicts require a recorded
disposition. A signed payload proves package integrity, not the truth or safety
of every project-specific instruction.

An upgrade preserves project-owned overlays and detects incompatibilities before
mutation. Existing runs keep their pinned procedure; migration or restart is
explicit. Old definitions needed for history remain resolvable. Rollback of files
is distinct from undoing external effects such as a release. Adoption may be
staged; installed capability does not automatically mean a project has replaced
its former working method.

External project/card references remain valid as **unverified external references**.
They do not require all repositories checked out or a fragile global worktree map.
When a procedure needs actual peer evidence, acquire a pinned snapshot/contract
explicitly or report it unavailable. An external reference by itself cannot
satisfy that gate. Portfolio summaries are projections of accessible projects,
with missing/stale coverage shown, not a new central source of authority.

## 7. Requirement set and evaluation

These are proposed requirements, not implemented capabilities. Each has concrete
challenge cases in [Scenarios.md](Scenarios.md). A test pass must name the actual
implementation/host/candidate; the present walkthrough only evaluates the design.

| ID | Requirement and acceptance intent | Cases |
| --- | --- | --- |
| R01 | Resolve project, checkout, profile and current authority; show conflicts and dirty inputs | S02, S04, S06, S10 |
| R02 | Route every governed request and scope change, with a no-work path and no duplicate plan/card per message | S01, S03, S15 |
| R03 | Distinguish missing, unavailable, ambiguous and incompatible routines; deduplicate gap handling | S05, S11, S12 |
| R04 | Preserve owner intent/authorization and hold only dependent work when a prerequisite is missing | S01, S04, S07, S15 |
| R05 | Check project fit, consumers, invariants and unknown impact before worker assignment | S05, S08, S09, S16 |
| R06 | Connect assignment context to BP2 and honest SDL coverage; no generated certainty from missing models | S05, S09 |
| R07 | Keep management, operational, evidence and external facts distinct | S06, S07, S09, S13 |
| R08 | Reject unauthorized, stale or insufficient-evidence transitions; preserve historical results | S08, S09, S13, S14 |
| R09 | Recover after restart with idempotency, concurrency checks and reconciled material records | S12, S13 |
| R10 | Distribute versioned base routines and preserve compatible project refinements and pinned runs | S02, S14 |
| R11 | CLI and MCP share transition rules; measure host capture/guard coverage and limitations | S11, S12 |
| R12 | Read-only monitoring supports reconnect, nested/parallel work, waiting reasons and historical replay | S12, S13 |
| R13 | Review procedure gaps and repeated exceptions without automatic policy invention | S11, S16 |
| R14 | Preserve role boundaries and distinguish self-check, independent review and owner acceptance | S04, S07, S09, S13 |
| R15 | Treat external references as unverified until required evidence is supplied; minimize private data exposure | S09, S10 |
| R16 | Keep process overhead proportionate and measure owner rescue, duplicate records and false routing | S03, S15, S16 |

For the pilot, measure route disagreements, missed scope changes, duplicate gap
cards, repeated owner reminders, recovery success and time/records added per task.
Record false-positive blocking as seriously as missed gates. Compare matched tasks
with a baseline where feasible; never claim causal drift reduction from a single
successful demonstration. Critical acceptance cases include wrong candidate,
missing real integration evidence, interrupted record updates and a bypass attempt.

## 8. Alternatives and staged next work

| Alternative | Benefit | Limitation | Assessment |
| --- | --- | --- | --- |
| Docs MCP plus stronger skills only | Small, useful retrieval improvement | No durable run or transition validation; compliance still relies on the agent | Useful component, insufficient for the owner's full goal |
| New central agent/portfolio orchestrator | Potentially stronger execution ownership | Large migration, new authority, difficult host integration and premature abstraction | Defer |
| SDPTool core, small MCP adapter, entry integration and observer | Reuses current process/tooling; can be introduced one workflow at a time | Enforcement remains bounded; requires deliberate evidence and recovery design | Recommended |

The next deliverable should be a bounded **DesignPlan**, selected under KB036/037,
for routine identity, request resolution, state/evidence and adapter boundaries.
The following phases are a proposed roadmap, not an activated implementation plan
or permission to merge/release. Select Git policy in that plan.

| Phase | Milestones | Exit evidence |
| --- | --- | --- |
| RG1 — contract and entry | One Maintenance routine plus no-work and gap paths; define state/authority/version contracts; select actual host adapter | Reviewed positive/negative fixtures, clear enforcement scope, no competing PM ledger |
| RG2 — durable vertical workflow | SDPTool core + CLI + minimal MCP read/transition interface; persisted run and record reconciliation | Same decisions through CLI/MCP; rejected invalid transitions; crash/retry recovery |
| RG3 — visible pilot | Read-only terminal view; one sandboxed SDP maintenance-to-release rehearsal with worker context | Reconnect/replay and evidence drill-down; no real publication; owner can recover context without chat |
| RG4 — blueprint and second project | Coordinate with selected BP2 work; trial one XFMD workflow and later one bounded Ponsse MVP1 workflow | Impact/unknown frontier visible; paired contracts where needed; project overlay survives upgrade |
| RG5 — distribution and refinement | Release catalog/adapters via existing installer after authorization; evaluate gaps and retire overlap | Fresh install/upgrade, incompatible overlay case, version-pinned historical run and measured owner burden |

Do not wait for complete SDL runtime, all MVP1 language support, a graphical TUI,
full portfolio onboarding or automatic code generation. Do not silently turn this
study into the implementation of any of them. Existing roadmap cards retain their
scope; concrete new findings can be added there without multiplying wrapper work.

## 9. Decisions and limitations to carry forward

The study recommends the responsibility split and a small pilot. Selection of the
next plan, exact routine schema/store, host guard coverage, project overlay policy
and acceptance thresholds remain design work. Any actual external publication or
project migration requires its applicable authorization; the study grants none.

The process should make it possible to tell the owner: “This changes the active
scope; these existing rules apply; I can continue this part now; this particular
part needs a decision.” It should retrieve that basis itself and preserve the
answer. The owner should not have to remember a card number or police every step.

There is no guarantee that a state machine makes an agent understand the whole
system. Its useful role is narrower and testable: make missing context, authority,
evidence and procedure coverage explicit before an unsupported transition occurs,
and retain enough history that tomorrow's session can resume responsibly.
