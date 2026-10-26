# Session execution and checked KanBan lifecycle — T041

Owner Session0007 T041 accepts the T040 direction and selects its concretization.
This is the bounded design for PLAN-SDP-0018 / KB052, linked to KB038 and KB037.
It specifies the target behavior; runtime now includes the bounded C1, C2a and
C2b1/C2b2/C2c1 subsets described below. Full lifecycle delivery remains open.
[Orchestration](Orchestration.md) retains role/event behavior and T039 mandatory
admission/model rules. [Acceptance](Acceptance.md) owns the test case register.
The [SDL concern](../../../SDL/ProjectGovernance/SDPTool/SessionWorkflow/README.md)
is a mutable WORK, not a released model or an executable workflow implementation.

## 1. Selected boundary and rationale

Extend the existing Go governance core and controller with a versioned Session
execution binding and checked lifecycle effects. Keep CLI and role-bound MCP as
adapters to the same command/revision rules. Keep existing operational generations,
receipts and canonical management history; do not introduce a second engine,
database, general-purpose workflow language or required daemon for this increment.
The existing controller can host the scheduler; closing an observer is independent.
A stopped controller requires explicit recovery/start of that controller, not a
promise that agents or scheduling survive termination of every owning process.

Rejected alternatives: agent-only reminders cannot enforce required record effects;
executing arbitrary Markdown makes edits ambiguous and unsafe to resume; a second
Session engine duplicates run/attempt/recovery ownership. One enormous state enum
mixing card status, host liveness, review and waiting conditions loses their distinct
meaning. Use a fixed bounded stage graph plus orthogonal holds and existing attempts.

This increment covers one selected card-derived Session plan, multiple ordered or
explicitly independent step assignments, design/no-design routing, implementation,
review/rework, model release and disposition. Other routine families and cross-machine
scheduling remain separate. Model release here is distinct from a software release.

## 2. Authoritative state and identifiers

| Meaning | Owner | Binding rule |
| --- | --- | --- |
| Goal, rationale, roadmap step IDs | Authored Session Markdown | No second authored milestone plan; retain journal and human prose |
| Outcome/milestone/acceptance and Git policy | Existing selected plan | Exact plan ID/path and reviewed revision; changes require reconciliation |
| Card state, placement and work history | Existing card metadata plus canonical management ledger | One writer protocol; one CardState; no private competing card truth |
| Executable selection | Versioned Session execution binding | Explicit references and dependencies, not a copy of prose or executable shell |
| Current stage, holds, receipts, pending effects and instance revisions | Existing PG operational store, extended by a new supported profile | One nonterminal execution instance per Session/binding revision in the selected checkout authority |
| Assignment/run and host attempts | Existing PG run/attempt/controller structures | Session instance references existing runs; retry never means a new Session |
| Models, Blueprint assignment and evidence | Existing ModelGovernance/Blueprint/verification owners | Capture exact identities and applicability; PG does not reimplement their lifecycle |

Select a project-owned JSON binding at
`SDP/Sessions/execution/<SESSION-ID>.json`, linked from its Session. This location
and `sdp-session-execution/0.1` are target design, not an already installed schema.
Validation must be implemented before runtime adoption; old manual Sessions remain
readable and are not silently auto-enrolled. Keep the binding in Git; operational
instances remain in the existing private checkout store. Markdown alone is not a
backup of host/effect state.

Required binding fields: schema, Session ID/path, selected plan ID/path, workflow
ID/version/digest and a nonempty set of selected steps. Each step references its
existing Session step and plan milestone, source KB IDs with coverage `whole` or
`contribution`, dependency step IDs, bounded scope/protected paths and required
review/acceptance policy. A single-step plan uses one entry. Resolve record paths
from canonical identities at adoption; do not let moving a card create a new ID.
Reject missing/duplicate step IDs, cycles, self-dependencies, unsupported versions,
contradictory membership and unknown fields. No shell text or arbitrary next-state
expression is executable input. Adding/removing steps is a recorded replan.

Instance identity includes project/checkout, binding digest, pinned routine digest,
Session identity and revision. Each selected step owns stage, design-impact rationale,
Blueprint requirement, exact model/task/evidence refs, optional run ID, role/scope
claims and holds. Business stages are admission, design, implementation, verification,
review, model-release, disposition and done. Holds are named prerequisites/conflicts
with reason, evidence and clearing operation. Host running/interrupted/unknown is
read through attempts; it never silently changes the business stage or card folder.
Instance lifecycle is active, on-hold, completed, canceled or superseded. Pause/cancel
requires an actual disposition; ordinary waiting is an active instance with a hold.

## 3. Commands and transition/effect table

Every mutation supplies operation ID, instance ID (except adoption), expected
instance revision and reason/evidence references. Trusted adapters supply principal
and capabilities, never a caller-supplied role string. Same operation/content returns
the original receipt; changed content under the same ID refuses. Reads return current
revision, next permitted actions, unmet conditions and observation freshness.

Names below are shared service operation concepts, not shipped command spellings.

| Operation / source | Checks | Result and required canonical effects |
| --- | --- | --- |
| adopt-session / absent | KB -> Session -> plan chain; execution authority; unique instance; supported binding/routine; fresh preflight | Create active instance, steps at admission. Record binding/adoption in Session; no host launch or fake card start |
| select-step / admission | Dependencies satisfied; whole/contribution coverage; authorization; allocate bounded writer/scope claims; plan selected | Activate selected plan/card if needed; backlog/queued -> active/ready; preserve an already active card. Record selection/history, then permit assignment preparation |
| assess-design / admission | Attributed rationale tied to task/plan revision; required skill and design criteria | Design-impact yes -> design; no -> implementation with explicit no-model disposition. Record separate Blueprint requirement; unknown impact holds |
| start-design / design | Ready records; bounded SAD authority; claims; baseline or explicit initial bootstrap selected | Publish owned work as in-progress, then dispatch SAD through controller. SAD creates/updates SDL WORK; missing model is preparation work, never a code bypass |
| accept-design / design | Valid WORK, exact NOW/TARGET/task; required Blueprint and review; actual required design acceptance; current scope | Pin implementation handoff -> implementation. Changed model/task stales prior evidence; no automatic code launch |
| start-implementation / implementation | Current handoff or justified no-design route; Master scope/skills; dependencies/claims; no holds | Record owned assignment/in-progress before host dispatch. Use existing run/attempt identities; pending dispatch is visible separately from observed host running |
| submit-result / implementation | Bound agent report with exact candidate and selected milestone; preserve unrelated work | Record progress/evidence references -> verification. Card stays in-progress; reported success is not closure |
| verify-result / verification | Required actual checks against candidate, coverage and provenance | Passing checks -> review; failure -> implementation with rework reason and retained failed evidence |
| review-result / review | Independent review where required, exact candidate, current design/task; actual owner acceptance only if policy requires | Failed review -> appropriate implementation/design rework. Passed/current review -> model-release if design changed, otherwise disposition. A concrete owner decision creates a waiting-owner hold and gate-review; automated review alone does not |
| release-model / model-release | Reconcile WORK/current model head and code; combined validation; frozen exact CANDIDATE; integration-owner acceptance and verified evidence | Call existing model service, record RELEASE identity/predecessor/code evidence -> disposition. No model-only proof of implementation; pending/conflicted release holds full completion |
| dispose-step / disposition | Result publication committed; review/evidence current; model release satisfied or explicit scoped transfer; no pending host/effects | Mark this contribution done and update milestone/Session worklog. Complete/move a card only if its whole outcome is satisfied and every contribution has a disposition; otherwise keep it active with remaining work |
| complete-session / all selected steps disposed | Goal achieved; selected plan/card obligations reconciled; transfers named; no unresolved claims/effects | Close satisfied plan/Session, release remaining claims, preserve evidence. A Session completing one contribution cannot close a broader card |
| replan / nonterminal | Authorized changed scope; identify affected dependencies and in-flight effects | Hold affected continuation, reconcile/stop affected host work, refresh binding/handoffs and impacted stages; preserve unaffected completed work and historical evidence |
| defer / nonterminal | Explicit authorized postponement, restart condition and affected membership | Instance/selected work on-hold as appropriate, card onHold only if its remaining work is deferred; no inference from crash or idle time |
| cancel or transfer / nonterminal | Actual disposition, remaining scope and successor for transfers; reconcile in-flight work | Preserve receipts, update only affected card/plan outcomes. Terminal instances remain historical; reopening creates a successor instance with fresh authority |

Progress updates are mandatory at meaningful phase/report boundaries even when no
folder changes. Worklog includes result, evidence, remaining work and next action.
One step may loop through implementation/verification/review for several milestones;
new phases cannot exceed the selected plan. New independent goals get separate
records and execution instances under T039. A gate-review card returns to in-progress
on requested rework; an accepted gate followed by more execution also resumes the
active state rather than pretending whole-card completion.

## 4. KanBan effects and crash-safe publication

Extend the existing publication service into typed lifecycle intents; do not build a
parallel writer. An intent binds operation/instance revisions, source preimages,
card/plan identities, semantic transition, exact after-images, event IDs and expected
postconditions. Changes may include metadata/worklog, source removal and destination
creation for a move, index/current links, Session/plan projections and ledger append.
Historical ledger bytes and frozen evidence references are never rewritten.

Protocol: validate/claim -> persist prepared intent -> apply checked file effects ->
validate canonical postconditions -> mark committed and advance instance -> dispatch
any dependent host action using its own durable intent. A transient effects-pending
hold prevents further dependent transitions. Multiple renames are not an atomic
filesystem transaction; coordinated readers report pending publication or retain
the last committed view until reconciliation. Only after postconditions pass may
PG advertise the new lifecycle. Direct filesystem readers may see intermediate files.

Use one project writer coordination protocol shared with existing projecthistory
and Blueprint writers. Establish lock order (project transaction claim, canonical
history lock, then short operational-store lock) and never hold these while waiting
on agents or owner input. Allocate/check IDs under the canonical writer boundary.
No long-held whole-project execution lease: scope claims serialize only conflicting
assignments while record writes serialize briefly. Migrate all cooperating writers
before claiming they are safe together; arbitrary shell editors remain outside it.

Resume checks each effect: exact preimage -> apply; exact intended postimage ->
already applied; neither -> conflict, preserve external bytes and stop dependent
work. For moves, check both paths and identity; never delete a different destination
or overwrite an edited source. Do not launch until canonical effects are committed.
For model promotion, save target label/candidate/expected predecessor and inspect
existing model receipts/artifact identity after uncertainty; if equivalence cannot
be established, require reconciliation instead of blindly releasing again.

## 5. Input changes without losing authority

Replace blanket global-ledger invalidation for the new Session profile with explicit
coverage reconciliation; leave legacy run semantics unchanged until explicit adoption.
Always validate complete canonical history integrity before accepting a new snapshot.
Track relevant card/plan/Blueprint predecessor events, selected binding revision,
shared authority documents, model/task/candidate digests and scope claims.

- Own committed effects: advance coverage to exactly recorded postimages/receipts.
- Unrelated append-only valid history: preserve semantic pins and record refreshed
  observation after confirming no selected subject, shared authority or claim changed.
- Unexplained edits to selected documents, membership, task, model or authority:
  hold and classify the delta; never blindly accept the new digest. Even an apparently
  editorial manual edit requires explicit reviewed rebind before advancing.
- Replan: new binding revision with impact analysis, authority and refreshed affected
  evidence. Preserve prior snapshots; do not translate an old review into new approval.

Old and new runtime profiles must refuse unsupported writes. Resume stays pinned to
its original routine version/digest; upgrade is an explicit validated migration or a
successor instance. Resume cannot silently install the newest workflow definition.

## 6. CLI, MCP, scheduler and observations

Select the existing trusted local controller as the host-action owner. Owner CLI and
an explicitly privileged leader adapter use the shared Session operations. The old
worker MCP remains restricted to its assignment reports/questions; no role argument
upgrades it. Leader may adopt/select/continue within granted scope, but cannot author
owner acceptance. Owner CLI can record actual owner decisions and request recovery;
that still cannot bypass missing tests/model release/conflicts.

Proposed user verbs are session status, inspect, resume, answer and watch; exact CLI
and MCP names/schemas will be registered through the existing action/discovery
contracts during implementation. Inspection includes Session/step/card/plan, role/run,
business stage, card state, host observation, pending effects, evidence/model release,
next permitted actions and blocking reasons. Tool events show their observed source,
status and cursor; unavailable children/activity remain unknown. Do not expose private
reasoning or claim to see uninstrumented external tools.

A bounded scheduler selects eligible steps after effects commit, serializes decisions
for each leader thread, respects configured concurrency and emits durable notices.
Model completion only queues evaluation; it is not unconditional next-step dispatch.
If the leader is absent, owner CLI can drive an already authorized deterministic
continuation. A missing semantic scope/design decision remains waiting for an actual
leader/owner assessment; the scheduler does not invent that judgment.

Resume attaches to the existing instance: inspect pinned profile/binding and current
coverage, reconcile effects, recheck claims, locate current attempts, then return
ready/waiting/conflict/uncertain or perform the authorized next action. Existing live
work is reattached/observed, never duplicated. Lost acknowledgments reconcile using
existing dispatch identity and host evidence. Unknown execution holds dispatch until
resolved. Stale owner/leader requests fail expected-revision checks. Observation
clients may close independently. Cross-machine relocation needs a separate portable
handoff of operational state and host reconciliation; Git files alone are insufficient.

## 7. Bounded implementation sequence and acceptance

Keep PLAN-SDP-0018 and its PGF checkpoints. PGF1 must first implement the disposition
contract with typed canonical effects; PGF2 adds Session binding/admission/scheduling
and shared adapter commands; PGF3 adds required model integration/release gates;
PGF5 exercises recovery across them; PGF7 verifies one combined native candidate.
Deliver each as a runnable vertical increment with retained compatibility. Do not
claim whole-session readiness from isolated schema tests or a diagnostic Blueprint.

The first lab scenario is a small design-affecting feature with one primary card,
Session and selected plan, plus an unrelated dirty file. Run selection, SDL WORK/
Blueprint handoff, Master implementation, deliberately failed review, correction,
passing review, model release and card/Session disposition. Interrupt between card
move and ledger write, and lose a host-launch acknowledgment. Recover via owner CLI
while the leader is absent, then reconnect the leader to the same instance.

PG-A25–30 additionally cover a plain no-design task, shared-card partial contribution,
manual edits, unrelated ledger append, routine version change, concurrent conflicting
Sessions and all adapter authority refusals. Prove no duplicate launch, lifecycle
event or model release. This design is ready for bounded implementation planning;
model parsing and preview evidence establish structural validity only. Independent
review, native acceptance and final model release remain explicit obligations.

## T042 implementation boundary — PGF1-C1

C1 implements binding/context, pinned durable adoption, inspection and registration
recovery through the core, CLI and six discoverable JSON actions. Wire contract is
`sdp-session-execution/0.1`; persisted profile is 0.13. Adoption prepares one exact
Session-journal effect; resume commits/reconciles it. All selected steps remain at
admission with explicit activation holds. This implements the foundation of the
transition table, not its later select/start/advance/disposition behavior.

The first increment deliberately retains conservative whole-ledger coverage and
one instance per Session in the selected store. C2 must deliver the scoped
reconciliation and lifecycle effects described above. One-file registration uses
the existing project writer lock and exact preimage/postimage recovery. It does not
claim a global transaction or prevent uncooperative external edits. A pending
publication from another store blocks registration; an external edit after a
registration write blocks recovery rather than being overwritten.

The authority document is pinned evidence supplied by a trusted controller, not
natural-language permission interpretation. JSON invocation host authority is
checked separately from request fields. Existing restricted leader MCP remains
unchanged; common privileged leader commands and dispatch are still PGF2 work.
See the [implementation evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C1.md)
for candidate identity, tests and explicit outstanding gates.


## T043 implementation boundary — PGF1-C2a

The first C2 checkpoint implements select-step and durable canonical effects, with
one dependency-free selection per instance. C2's remaining start/progress/disposition
and dependency completion are explicit follow-ons. Selection stays at admission;
there is no no-design shortcut to execution or completion. Existing WORK :00002
covers SessionCoordinator, CanonicalEffectPublisher and RunRepository ownership.

Selected backlog/queued cards move to active/ready; active ready/in-progress cards
retain state. Planned plans become active. A typed intent pins all canonical inputs,
exact destination/source effects, index/current links, journals and allocated events.
Prepare increments instance revision but does not apply records; resume applies and
validates them before committing selection. Whole/contribution coverage is retained.

Current PG writers and projecthistory/Blueprint history now share the canonical
ledger lock and durable transaction reservation. Lock order is project coordination,
canonical history, then operational store, with no host wait. Legacy PG publication
keeps its existing private pending marker for recovery compatibility and also uses
the canonical reservation. Direct manual writers and older independent processes
remain outside this contract; no universal write interception is claimed.

Move references are bounded to selected cards, plan, Session and board index. Other
references or unsupported syntax refuse before effects, avoiding unapproved edits
to historical records. A selected scope claim serializes conflicting Session/card/
plan or writer/reader selections across current stores. Claim release and broader
legacy-run admission remain future gates. A new store cannot register a Session
already selected elsewhere. Claims do not expire merely because an agent exits.

After committed selection, exact historical-prefix plus unrelated validated ledger
appends are reconciled; selected-document/subject changes and new Blueprint events
hold. Prepared effects retain exact all-input preconditions and cannot be silently
rebased. [Evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2a.md)
records bounded verification and remaining C2/PGF2–9 obligations.


## T044 implementation boundary — PGF1-C2b1

The first C2b checkpoint implements assess-design through the shared Session core,
CLI and trusted-controller JSON action. An explicit seven-domain assessment and
separate Blueprint decision are tied to the selected binding/step, instance revision,
one to eight evidence pins and required Project Leader plus Change Analysis/Architect
skill snapshots. This refines section 3's attributed assessment contract; it does
not claim semantic inference from prose, skill consumption or independent approval.

Changed impact routes to design and requires SDL WORK/accepted handoff. Explicit
model-file/source writable scope cannot be classified as no-design. Unknown impact
or unknown Blueprint applicability prevents no-design admission. All-unchanged impact
permits the implementation business stage; required Blueprint remains a hold. Every
route retains a start/dispatch hold until the remaining integration is delivered.
The existing WORK :00002 covers RequireDesignWork, AdvanceWorkflowStage and checked
effect ownership; no semantic model rewrite or release is selected here.

A prepared assessment has no canonical effects. Resume publishes the exact attributed
JSON record, Session/plan/card journals and history, then commits stage visibility.
CardState remains ready/in-progress. Operational profile 0.15 stores bounded
assessment history (eight records), while preserving the selection intent and all
prior assessment artifacts. A new assessment can resolve uncertainty against the
unchanged pinned basis. Changed evidence or procedures require explicit reconciliation,
not an automatic new authority. Existing checked writer locks, reservations and
pre/postimage recovery apply. Old resume receipts cannot release newer pending work.

C2b2 remains responsible for start/progress and exact Session-to-existing-run binding.
Legacy Maintenance and approved-Blueprint run contexts have their own validation;
relabelling them would bypass those checks. No second execution engine or fabricated
run report is introduced. C2c owns verified disposition/dependencies/claim release;
PGF2/3 retain privileged leader and model handoff/release gates.


## T045 implementation boundary — PGF1-C2b2

Start/progress now implement the bounded no-design path through the existing
run/attempt/controller, not a second execution engine. An immutable start records
exact Session/step/assessment identity, outcome and Master/Worker skill snapshots.
Resume commits all selected cards as in-progress and creates one ready run atomically
with operational start commitment, after canonical effects. Explicit controller launch
happens later; committed intent does not prove host execution. The full binding keeps
whole/contribution coverage for MaintenancePlan and ImplementationPlan alike.

The run input revision is stable across its own canonical journals, while current
Session coverage independently validates pinned inputs and unrelated ledger appends.
Existing run scope, report capability, host intent/recovery and evidence ownership
apply. Current legacy admission/dispatch checks persistent claims optimistically;
manual writers/older binaries and universal cross-store arbitration remain outside
coverage. Writable scope cannot overlap pinned governance records or required skills.

Progress pins the exact assigned report, run revision, evidence artifact and candidate.
Every effect rechecks that basis. Unknown host outcomes or outstanding attention/input
hold progress; committed progress enters verification and keeps cards in-progress.
It does not accept the implementation or complete a contribution/card. Existing
single-card result publication and further advance/relaunch/rework cannot bypass
future Session disposition. Stopping operations remain available.

Operational profile 0.16 and the two trusted-controller CLI/JSON actions reuse the
existing effect publisher and receipts; old readers refuse. The current bound is
one start and one progress for one selected step. Rework/dependency progression,
whole/contribution disposition and claim release/replan remain C2c. Changed/unknown
design and required/unknown Blueprint cannot use this start path; PGF3 must implement
WORK/model handoff and release gates. PGF2 retains privileged leader orchestration.
[Evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2b2.md)
records fault, candidate drift, strict wire, multi-card and actual-controller protocol
tests. No native model execution, independent review or release is asserted.


## T046 implementation boundary — PGF1-C2c1

The verify-result/review-result prerequisites now bind the Session stage to existing
candidate evidence and observed review ownership. Verification consumes an exact
passing test receipt for the progressed report/candidate. It cannot advance on a
claim, failed/inconclusive result or contradictory current failure. Canonical records
commit before a required review can start. Existing review preparation carries the
full Session binding and requires the verified receipt among its selected passing
tests. Fresh-thread provenance remains owned by the existing review engine.

Review intake pins the exact terminal passing review revision, its tests and complete
payload/binding. Current candidate/artifact/skill checks and absence of conflicting
review/evidence are rechecked at each canonical effect. The shared effect protocol
publishes journals for all selected cards, the plan and Session, immutable result
artifacts and management history. Prepared effects never advance the visible stage;
external drift holds recovery and preserves bytes. Operational profile 0.17 retains
these snapshots without rewriting historical start/progress artifacts.

Committed verification routes to review where required; explicit not-required policy
routes to disposition. Committed review intake routes to disposition while preserving
any owner-acceptance hold. These stages do not mark a contribution/card completed or
release claims. C2c2 remains responsible for owner decisions, failed-result rework,
whole/contribution disposition, dependencies and claim release/replan. This is a
bounded implementation of the existing table, not a waiver of those obligations.
WORK :00002 remains TARGET; model review/integration/release is still required.
[Evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c1.md)
records actual scope and limits; synthetic review observations are not independent
PG review or native whole-workflow acceptance.


## T047 implementation boundary — PGF1-C2c2a

`dispose-work` now has a bounded no-design, owner-acceptance-not-required path.
Exact committed verification/review, current candidate and run revision, and absence
of outstanding host input/attention are prerequisites. The controller submits all
ordered card outcomes: whole completion or a completed contribution with explicit
remaining work. Another bound step retaining a whole card's scope blocks its closeout.
Required owner acceptance refuses until the trusted decision adapter is implemented.

Preparation retains only named selected preimage bodies and a deterministic effect
plan. Retained-state validation recomputes those effects. Resume uses the existing
publisher to journal selected records, move eligible whole cards, update metadata/
selected links/index, append management history and separate system evidence, and
remove source cards/claim last. The pending canonical reservation remains until
operational commitment; current legacy claim checks honor that reservation too.
External edits hold recovery without overwriting bytes. Same-operation retry cannot
repeat committed effects. This does not arbitrate manual writers or old executables.

Profile 0.18 retains the disposition. Only commitment completes the Run and selected
step. Plan/Session stay active; owner decisions, failed-result rework, dependencies/
replan and explicit overall closure remain later C2c2 work. Full workflow delivery,
SDL handoff/model acceptance/release and privileged leader integration remain open.
The existing WORK :00002 target already describes this responsibility; no model
source is relabeled accepted or released by implementing this subset.
[Evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c2a.md)
records the actual candidate, fault matrix and compiled lab trial.


## T048 implementation boundary — PGF1-C2c2b

Concrete owner gates now use the shared effect publisher. Controller authority may
request an exact-result review; only the trusted local owner adapter can decide it.
The JSON action propagates the actual host principal rather than silently replacing
it with controller authority. Committed request/decision moves all selected cards
through gate-review/in-progress. Passing tests/review do not substitute for acceptance.

Continuation archives each selected execution with its own assessment, run, evidence,
owner decisions and disposition. Dependency-checked next-step selection and rework
from exact failed tests/reviews or owner rejection require fresh assessment and work.
A completed prior contribution can satisfy a later step's dependency without losing
its evidence. The failed run is retired only after canonical effects commit. History
is flat and bounded; this is not the sustained-use archive/rollover solution.

An initially activated plan readies all bound execution cards, satisfying the existing
plan contract; only the selected step acquires a scope claim. Overall closure is a
separate operation: every bound step and whole-card obligation must be satisfied.
Remaining contribution scope still refuses without implemented successor transfer.
Authoring a changed binding still requires the later explicit rebind/replan path.
These retained boundaries are not waivers of the original transition table.

Profiles 0.19/0.20 retain deterministic effects and flat execution history. Current
readers accept prior profiles; mutators preserve newer profile requirements. Same
candidate checks, exact pre/postimage recovery, shared reservation and immutable
receipts apply. CLI/JSON are exposed; privileged leader MCP and native integration
remain later checkpoints. Existing WORK :00002 covers these responsibilities and
stays unreleased TARGET. [Evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF1-C2c2b.md)
records passing full race, focused recovery/tamper, wire and compiled lab checks.


## T048 common adapter boundary — PGF2-C1a

A separate `sdp-mcp --session-controller` launcher capability selects one registered
instance and a named, pinned project-local Project Leader skill. Twelve tools expose
the supported core operations and receipts. Ordinary coordinator/worker grants do
not expand. Instance switching and owner decisions refuse; adoption and host dispatch
are absent. Each call checks the skill snapshot; actual native loading/application
remains separately unverified. This implements ExternalLeaderMcp's shared-command
boundary without making MCP the scheduler. [Evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGF2-C1a.md)
records final protocol, race and durable replay checks. WORK :00002 remains TARGET.


### T048 owner host entry through the current Session

PGF2-C1b resolves the current committed Session assignment for existing owner host
operations. Historical executions remain inspectable by run ID but are never selected
implicitly. Dispatch still requires ready/current canonical coverage and the exact
start-pinned skills; it uses existing host leases, durable intents and uncertain-send
reconciliation. Missing or pending start effects refuse launch. No Session transition
or default coordinator privilege is added. This implements the existing WORK's
HostActionGateway boundary; it does not provide autonomous leader scheduling.


### T048 terminal projection and current-attempt replies

PGF4-C1 exposes existing Session operations through a terminal adapter with unambiguous
canonical-ID selection, guarded automatic record commitment and explicit recovery.
It owns no new runtime state. Exact owner gates and current attention questions remain
bound to the original result/basis; EOF is not input. An already committed record
resume is read-only. Agent execution remains a separate host operation.

A native continuation revealed that comparing only the old question basis to the new
attempt falsely marked its selected reply stale. The attention projection now retains
current applicability only for an exact reply referenced by the latest claimed or
acknowledged launch, mapped to the current gap-free attempt of the same running
assignment with current sources. Later input/scope/source changes, uncertainty or
terminal routine state cannot inherit that exception. Admission, retry identity and
semantic-consumption claims remain unchanged. The initial failed trial and corrected
native trial are retained in PGF4-C1 evidence.
