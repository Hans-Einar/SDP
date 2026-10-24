# ProjectGovernance core — PGI1/PGI2

This package implements the first slice of
[PLAN-SDP-0018](../../SDP/05--Implementation/SDPTool/ProjectGovernance/Plan.md):
read-only work inspection, one pinned Maintenance routine and durable local run
transitions. Linux supports mutation; other platforms explicitly refuse writes.
The CLI and [MCP adapter](../governancemcp/README.md) use the same core.
The [app-server adapter](../governancehost/README.md) supplies protocol inspection
and bounded observed attempt recovery. The [owner controller candidate](../governancecontroller/README.md) adds durable
input and bounded submission; full PGI3 verification remains open. Bounded result
publication and Session registration are described below; final lifecycle disposition
remains a later increment.

Current integrated-state responsibilities and the distinction between legacy
Maintenance routing, prepared design packages and execution runs are summarized
in the [state responsibility map](../../SDP/04--Design/SDPTool/ProjectGovernance/State-Boundaries.md).
For the measured workstation and MCP adoption limits, see the
[adoption study](../../SDP/Maintenance/WorkstationReadiness/MCP-Adoption-Study.md).

## Inspect and select work

```sh
sdptool /path/to/project governance help
sdptool /path/to/project governance context --input selection.json --json
sdptool /path/to/project governance route --input route.json --json
```

`context` takes a Selection object, for example:

```json
{
  "card": "SDP/KanBan/active/card.md",
  "plan": "SDP/Maintenance/Fix/Plan.md",
  "session": "SDP/Sessions/session.md",
  "authority": "SDP/Maintenance/Fix/Authority.md",
  "step": "S1"
}
```

Use actual existing records. The selected card must be ready/in-progress and
reference an active MaintenancePlan with a `MAINT-` identity; the Session must be
active and contain the named roadmap step. Selected card/plan paths and lifecycle
must agree with their canonical management event chains. This is bounded coverage
inspection, not full validation of every project record. Use the existing project
management validator for whole-repository validation. Missing references appear
as conditions; they never authorize a run.

A route request contains `schema: "sdp-project-governance/0.1"`, `categories`,
optional `intent` and `selection`. Maintenance additionally selects
`routineId: "maintenance"` and `routineVersion: "0.1"`. `inputRevision` can pin a
previous context. A single `status` category yields `no-work`; `maintenance` yields
`matched` only with complete coverage. Other dispositions are `missing`,
`ambiguous`, `incompatible`, `unavailable` and `conflict`. Categories are proposed
by the caller; this core does not classify natural language. These operations
create no binding, run, project document or operational directory.

## Bind and run

```sh
sdptool /path/to/project governance bind --state /private/state/checkout --json
sdptool /path/to/project governance start --state /private/state/checkout --input start.json --json
sdptool /path/to/project governance read --state /private/state/checkout --run RUN_UUID --json
sdptool /path/to/project governance transition --state /private/state/checkout --input transition.json --json
```

Use the same state directory throughout. The default locator is
`$XDG_STATE_HOME/sdptool/project-governance/<sha256-canonical-root>`, falling back to
`~/.local/state`. It is outside the project, private to the local OS account and
bound to the canonical checkout root, discovered project identity and a newly
created UUID. The root hash locates the store; the UUID is its checkout identity.
A copied store cannot be used at another root. `bind` is explicit and idempotent.

`start` takes the project/checkout identity from bind, fresh UUIDv4 `operationId`
and `requestId`, `inputRevision` from context, `routineDigest` from route,
`selection`, and an `assignment` with `outcome`, `writablePaths` and optional
`excludedPaths`. It creates a ready run at revision 1. Relative scope paths declare
assignment boundaries; PGI1 neither executes changes nor intercepts shell access.

`transition` takes project/checkout identity, a new UUIDv4 `operationId`, `runId`,
`expectedRevision`, the run's pinned `inputRevision`, `action` and nonempty `reason`.

| Action | Allowed source | Result |
| --- | --- | --- |
| start | ready | running |
| interrupt | running | interrupted |
| resume | interrupted | ready |
| report | running | verifying |
| fail | running, verifying | failed |
| retry | failed | ready |
| scope-change | any nonterminal state | replan-required |
| cancel | any nonterminal state | canceled |
| complete | any | refuses with evidence_missing |

`step` stays `implement`: `report` means an implementation result needs checking,
not that the verify/review/record steps ran. Resume preserves run, assignment and
Session identity. It does not represent reconnection to a host attempt. Replanning,
PGI4-C1 adds the controller-only evidence and failed-evidence rework operations described below. Observed review in a separate context is described below; checked publication and completion remain later gates.
The serialized definition and digest stay pinned on each run.

## Authority, consistency and recovery

The CLI is an explicit trusted local controller adapter. Core mutation needs an
adapter-supplied Caller capability; requests cannot provide role/owner approval.
A referenced authority document is pinned input, not parsed proof of approval.
The host adapter remains responsible for deciding whether an invocation is
actually authorized. This boundary does not defend against arbitrary same-user
shell access. The MCP adapter never grants LocalController from model text.

A checkout-wide nonblocking Linux flock serializes mutations. A transition checks
its run revision and pinned source revision before publication. Selected records,
root AGENTS, project manifest and the whole canonical management ledger are hashed;
therefore any ledger append conservatively invalidates advancing work in this
pilot. Missing optional AGENTS/manifest inputs are pinned as absent. Unrelated files
are not hashed. Changes in coverage block advancement; interrupt, cancel, failure
and scope-change can still stop work after ordinary input drift.

Every accepted operation saves its normalized request digest and exact result in
an immutable SHA-256 generation. Synced generation bytes precede atomic HEAD
replacement; reads use only HEAD. Restart ignores orphan generations. Repeating
an operation ID with identical typed content returns its original receipt, even
when current coverage later changes. Use `read` for current coverage. Changed
content with the same ID fails; use the original operation after an uncertain
write response. Repeating a request ID with a different operation ID refuses to
create another run. Missing/corrupt referenced generations require reconciliation;
PGI1 does not repair or garbage-collect them automatically.

Limits: 1 MiB UTF-8 request/selected Markdown, 16 MiB management ledger or snapshot,
32 JSON nesting levels, 64 KiB intent/outcome/reason, 64 combined scope paths,
128 runs and 1,024 operation receipts per checkout. Capacity exhaustion refuses
without removing historical receipts. Explicit null, duplicate/unknown fields
(including case aliases), trailing JSON, non-integer revision tokens, hidden or
escaping reference paths and symlinked source/state paths refuse. These limits
are a pilot boundary. PGI2 Events pages add revision-pinned cursors and count/byte
limits for observed attempt content.

[Schema definitions](schema.json) describe each CLI request and response. Select
the corresponding `$defs` entry; semantic coverage, byte limits and duplicate-key
rejection additionally run in Go. Errors go to stderr with exit 2; success uses
stdout/exit 0. Human output is default and `--json` selects machine output.

## Observed attempts (PGI2)

Only the controller can bind an observed idle thread, record completed visible
items, mark a disconnect unknown, or reconcile a mapped attempt from host history.
New connections pin Codex 0.161.0; historical 0.160.1 attempts remain readable; a connection epoch change requires reconciliation and
preserves an observation gap. A host turn completing never completes SDP work.
AssignedAgent is a launch-supplied run/attempt/epoch capability: only report,
fail and scope-change are allowed, and a stale/disconnected binding refuses.

Run results and immutable receipts include attempt metadata and eventCount, with
empty events arrays. Read content through Events: at most 256 observations and
240 KiB of encoded observations per page, with a cursor pinned to run revision.
A single encoded observation exceeding 240 KiB refuses admission. Text is at most
64 KiB; there are at most 64 bounded attachment references, 16 attempts per run and
1,024 observations per attempt, also subject to the checkout receipt/state limits.
Oversize input refuses rather than truncating visible messages.

The unpublished sdp-project-governance/0.1 pilot schema adds optional attempts
and eventCount. PGI2 reads PGI1 state; downgrading an updated store to the older
strict PGI1 reader is unsupported. This is not a released profile migration.

## Verification

```sh
go test ./governance
go test -race ./...
go vet ./...
```

Run from SDPTool with the selected local toolchain. `TestCompiledGovernanceJourney`
builds the actual executable and uses a separate process for each call. It checks
read-only routing, explicit binding, retry, interrupt/resume, completion refusal,
human output and byte preservation of every fixture project file. Set
`PGI_SCHEMA_EVIDENCE` to a temporary directory to export observed wire cases, then
validate them with `python3 governance/check_schema.py PATH/cli-wire-cases.json`
(requires the existing Python jsonschema package). Core tests cover stale inputs,
canonical record conflicts, lock contention, corruption, and injected interruption
before/after HEAD publication. They do not simulate physical power loss or prove
host execution, independent review or the full owner workflow.

PGI3 candidate state also adds optional owner inputs and last-observation timestamps.
The controller-only dispatch contract and recovery limits are documented in the
[owner controller](../governancecontroller/README.md). Old strict readers cannot
read stores updated with these unpublished fields.


PGI3-C2 persists an optional user-message clientId. `input.resolve` is controller-only
and requires a recorded exact clientId/payload/turn match in the bound attempt;
it acknowledges an uncertain launch/steer without sending it again. Conflicting
correlation refuses, and acknowledgment never advances the SDP routine. Older
observations without clientId are not retroactively promoted to correlation proof.


## Observed approval decisions — PGI3-C3

Controller-only approval.requested/resolved observations retain the bounded SDP MCP
request and its closure. Approval owner inputs preserve exact stdin text separately
from the typed outbound decision, use the existing prepare/claim/uncertain sequence,
and resolve only after observed closure. No new agent mutation or management ledger
is introduced. Approval does not imply tool execution, evidence acceptance or work
completion; the default launch still refuses host approval requests.

For an exact accepted agent transition, the core can retain the original operation
ID/digest while crossing only recorded host observations and that reply's bookkeeping.
Every intervening revision and unchanged decision/attempt/coverage state is checked.
Changed arguments, steering, lifecycle mutations, reconnect and controller requests
retain normal stale-revision refusal. This prevents the approval record itself
from invalidating the queued call without weakening the normal CAS contract.

PGI3-C4 separates optional `agentEpoch` from the observation `epoch` on attempts.
Legacy state falls back to the existing epoch. New binds initialize both; ordinary
history reconciliation still changes the active capability epoch. Only trusted
`host.reattach` with the originally bound nonempty `hostIdentity`, current latest
attempt/turn, source coverage and resolved outbound work preserves that capability
while replacing the observation epoch. Approval identities remain connection-bound.
Same request IDs in different epochs are distinct approval observations. A gap is
never erased by reattachment. These additive fields are development-profile state;
older strict-decoding pilot binaries cannot read the new state and no published
upgrade compatibility is claimed.

PGI3-C5 adds `context.observed`, a controller-only launch observation containing
bounded ContextProvenance. It is accepted only for the latest idle, unmapped-turn
attempt and matching assignment revision. It does not create role authority.
Provenance returns a read-only as-of projection; a changed epoch is explicitly
historical. Missing prior captures stay unavailable. Host reports, controller
source hashes and skill availability have separate meanings; actual skill use is
not observed by this profile. New context records require this development
candidate's strict schema; no published downgrade/upgrade contract is implied.


Explicit launch skill inputs are controller-owned immutable content, stored with
local source digests and omitted from ordinary run summaries. Snapshot drift blocks
prepare/claim; history resolution does not recheck later bytes. User-message input
fingerprints preserve ordered text/skill identity, with unsupported content marked
unavailable. Legacy text-only observations remain readable. Provenance exposes a
bounded invocation projection and distinguishes acknowledgment from observed host
messages; compliance remains unobserved. These additive preview fields do not imply
backward writer compatibility or a new published distribution.

## Registered work projection

`Service.Overview` reads one persisted generation and projects registered runs in
stable ID order. It reuses context/coverage and provenance checks and creates no
state. The standalone controller exposes readable and JSON output; see
[overview usage](../governancecontroller/README.md#registered-project-work-overview--pgl1)
and [wire schema](overview.schema.json). This is the first KB052 inspection slice,
not a new assignment registry or multi-agent scheduler. Operational state and
canonical Session/plan/card authority keep their existing owners.

## SAD and implementation preparation — PGL2-C1

The trusted controller can prepare an immutable assignment independently of an
execution run. `PreparationContext` reuses canonical coverage for SAD
DesignPlan/ArchitecturePlan and Master ImplementationPlan; legacy `Context` and
Maintenance routing retain their plan restrictions. `PrepareAssignment` requires
binding, exact generation/source digests, literal bounded paths and return criteria.
SAD pins NOW; Master pins NOW/TARGET/task/blueprint/design-disposition files. These
references are not provider validation or a decision made by the owner.

Prepared work shares the same atomic generation store and operation-ID namespace.
Overlapping local writable scopes are refused in both directions with nonterminal
legacy runs. Nonliteral legacy scopes require reconciliation when preparing new
work. No automatic expiry, cross-checkout reservation or external-shell exclusion
is claimed. Repeated operation/content returns the immutable original receipt;
`InspectAssignments` rechecks current inputs and reports stale records separately.
Changing the source never silently updates a pin or grants continuation.

First preparation upgrades only the selected store to operational
`sdp-project-governance-state/0.2`. This version can still read 0.1; old binaries
refuse the new generation. Use a separate preview store with consistent binary
versions. Existing owner-lab state is not upgraded by inspection.
[Application schema](preparation.schema.json) describes the preparation wire;
Maintenance/MCP wire stays unchanged. All prepared results explicitly report
`dispatch: unavailable` and `designAcceptance: not-assessed`. C2 owns authenticated
disposition, provider-aware revisions and retirement/revision; PGL3 owns execution.


PGL3-C5 adds trusted-controller `AcknowledgeReport`: exact report/current-generation
binding, durable retry and separate recorded processing, without advancing the run.
The first acknowledgement selects operational profile 0.8 and inbox profile 2;
untouched stores retain their profiles. See the
[controller workflow](../governancecontroller/README.md#pgl3-c5-report-acknowledgement).


PGL3-C6 `Coordination` combines existing work/report/attention projections with
one operational generation, stable item IDs, prioritized bounded pages and a
content-bound cursor. It is read-only; it adds no store profile or lifecycle.
[Coordinator MCP](../governancemcp/README.md#explicit-coordinator-adapter--pgl3-c6)
and the controller console reuse this projection and existing guarded writes.

## Project Leader delivery — operational profile 0.9

PGL3-C7 adds attributed owner messages and serialized decision dispatches to the
same immutable-generation store. Leader delivery receipts share the global operation
ID namespace, use their own monotonic revision and never mutate assignment state.
An enqueue uses expectedRevision=0 and stable operation/content retry; other writes
use leader revision CAS. Reads never upgrade the profile. Explicit leader writes
select 0.9; preparation, attention, continuation and acknowledgement preserve it.
Older readers refuse. Historical report receipts and earlier profiles remain readable.

A separate leader lease serializes external effects while owner messages and worker
controllers remain usable. Durable prepared/opening/bound/sending/running boundaries
precede host effects. Only prepared/bound may supersede. Config, selected skills,
message IDs, notice keys and exact payload are retained. Unique observed host inputs
are required to resolve uncertain sends; dispatch completion is delivery, not approval.

See [controller usage and bounds](../governancecontroller/README.md#serialized-project-leader-turns--pgl3-c7)
and the separate output-only [leader schema](leader.schema.json). The worker MCP
input descriptions do not acquire leader-only definitions or capabilities.


## Candidate evidence and rework — operational profile 0.10

`CaptureCandidate` reads a run's literal writable scope plus explicit read-only
dependencies, excluding unselected write-excluded paths. The sorted inventory pins
file bytes, exact executable bits, directory membership and absent roots alongside
checkout/run/authority identity. No Git commit, model or file mutation is required.
Dependencies do not grant write authority; their completeness is not inferred.
Hidden components, symlinks, globs and special files refuse. Limits are 64 dependency
paths, 512 entries, 4 MiB/file and 32 MiB total file content.

`AttachEvidence` is a trusted-controller operation. It pins the latest report, run
revision, captured candidate and the digest of an external artifact (up to 4 MiB,
outside candidate scope). Types are claim/test/review; verdicts are passed/failed/
inconclusive. These are controller-attributed records, not observed test execution
or authenticated independent review. Attachment preserves run state/revision.

`Evidence` exposes history and separately checks candidate, artifact, report and
authority applicability. Source drift does not erase a receipt; exact retries
return the original result. Read projections detect operational generation changes.
External writers must remain quiescent for coherent filesystem inspection; no
filesystem-wide transaction is promised. At most 16 evidence records per run fit
within existing 1024 receipt/16 MiB snapshot bounds.

`ReworkEvidence` accepts only a current failed test/review on a verifying run. It
requires current authority, no blocking attention, terminal observed attempts with
no gaps and no unresolved controller inputs. It returns the same assignment/Session
to ready, preserves report/evidence history and rotates the prior agent capability.
A fresh attempt must carry further agent work. Completion still refuses.

Profile 0.10 is selected only by attachment; existing current writers preserve it.
Earlier binaries refuse rather than dropping evidence. The common wire schema's
Candidate/EvidenceRequest/EvidenceReworkRequest/EvidenceView/Result definitions
cover the CLI. Existing worker/coordinator MCP tool authority is unchanged.


## Observed review — operational profile 0.11

`ReviewWrite` prepares one controller-issued review with a C1 candidate, latest report,
1..8 current test receipts, review focus, exact instructions, host version, MCP
executable and explicit skill snapshots including sdp-reviewer. It does not change
the implementation run. The payload includes selected governing records, assignment
and test references, without copying implementation conversation history.

The native controller opens a fresh persistent read-only thread, checks empty history
and a distinct thread ID, and submits explicit skill inputs with the review ID as
client correlation ID. Run-bound MCP exposes only context_read, route_evaluate and
run_read; no attempt/epoch or coordinator capability is granted. Account mode must
be managed ChatGPT. Exclusive tool inventory and the supported host/wire contract
are enforced; inherited additional tools cause refusal.

Prepared/opening/bound/sending/running states retain effect intent. Unknown opening
never creates a replacement; prepared/bound work can be superseded without sending.
After controller replacement, an unsent bound thread is superseded. Sending/running
recover only the same host/thread/input/turn with exactly one full-history turn.
Only final visible text is retained. A malformed/inconclusive/failed answer cannot
qualify as passing review. Private reasoning and commentary are not review results.

`ReviewRead` separates historical outcome from current applicability and qualification.
Qualification means a passed structured review, observed in a separate fresh context,
whose report/candidate/test artifacts/authority/skills still match. It does not mean
that tests were independently rerun, a human accepted work, native child hierarchy
was established, or publication/completion occurred. Attached C1 review prose remains
controller-attributed regardless of its wording. Reviewer skill acceptance does not
attest procedural compliance or consumed file bytes.

Limits: eight reviews/run, 64/store, 64 KiB payload/final answer, 32 findings with
4 KiB each, within existing 1024 receipts and 16 MiB state capacity. Current writers
preserve profile 0.11. Older binaries refuse rather than dropping reviews. Same-user
shell edits and coherent filesystem inspection retain the C1 boundary.

## Checked result publication — operational profile 0.12

`PreparePublication` accepts a trusted controller request naming the exact run,
revision, passed review, system prefix, work summary and remaining work. It requires
current passing test artifacts, no contradictory failed evidence/review, terminal
known attempts, acknowledged inputs and unchanged authority. It generates a
reviewable six-file proposal: one new evidence document, appended Session/card/plan
entries, two management events and one Traceability verification event. No arbitrary
patches, paths or acceptance decisions are accepted. Broader lifecycle and authored
roadmap stay unchanged; `complete` remains guarded separately.

`PublicationRead` is read-only. `ApplyPublication` persists intent before replacing
files and uses exact before/after hashes for recovery; retry never appends duplicate
events. `AbandonPublication` only retires an unapplied proposal, including an orphan
project claim when intent persistence failed before file effects. Once applying,
unknown preimages stop with `reconciliation_required`, and other store mutations
refuse until recovery. There is no automatic rollback over external edits.

Canonical `SDP/` shared-management records are required. Existing installation
history validation is reused with the project's recorded 0.1-to-0.2 typed adoption
allowed; installation behavior stays strict. PG additionally checks payload fields,
sequence/time, board index, references and current grouping. Generated ledger bytes
are append-only and the complete proposed record set is validated before effects.
Unsupported minimal fixtures refuse. No arbitrary repository-specific validators or
semantic prose checks are inferred. Preparation pins all inspected canonical record
files, so even another work item's management change conservatively needs a fresh
proposal. Candidate/test/skill changes also prevent publication.

The project-writer lock and pending-store pointer live at
`$XDG_STATE_HOME/sdptool/project-publication/<root-sha256>/` (default
`~/.local/state`). This coordinates PG stores for the same canonical root. The exact
intent remains in its private operational store. Files are individually replaced
and fsynced; the set is not atomic. External editors do not share this lock: observed
changes refuse, but this is not a filesystem CAS or a same-user security boundary.

At most 16 proposals/store, 1024 canonical management files, 2048 input pins and
4 MiB of generated file contents/proposal are supported within the existing 16 MiB
store bound. Current writers preserve profile 0.12; old readers refuse. Publication
changes canonical authority bytes, so the original review becomes historical under
normal run coverage. The publication view separately checks its accepted postimages
and candidate. This does not rebase the old assignment or accept a new candidate.

Coordination uses `sdp-governance-coordination/2` when reviews are present (otherwise
version 1). Compact notices expose verdict, summary, finding count, applicability,
context provenance and publication status/current applicability; full findings remain
in `ReviewRead` and published evidence. Terminal review notices enter the existing
serialized leader delivery, with content-key deduplication and renewed delivery when
applicability changes. The coordinator MCP capabilities are unchanged.

## Session instances — PGF1-C1

The additive `sdp-session-execution/0.1` contract binds an existing active Session
and live typed plan to selected steps and source cards. It stores a version-pinned
instance in operational profile 0.13; old runs remain intact. This increment
registers and inspects work. Every step stays at `admission` with an explicit
activation hold. It does not start agents, move cards or complete a Session.

Create `SDP/Sessions/execution/<SESSION-ID>.json` and link it from the authored
Session. Use [session.schema.json](session.schema.json) for the exact wire shapes.
A binding supplies Session and plan `{id,path}`, an authority document path,
workflow `{id,version,digest}`, and steps with IDs, existing plan milestones,
`cards: [{id,coverage}]`, `dependsOn`, literal `writablePaths`/`protectedPaths`,
`review` and `ownerAcceptance` policies. The current workflow is read from
`session-definition`; do not invent or automatically upgrade its digest.

```sh
sdptool /path/to/project governance session-definition --json
sdptool /path/to/project governance bind --state /private/state/checkout --json
sdptool /path/to/project governance session-context --state /private/state/checkout --input context.json --json
sdptool /path/to/project governance session-adopt --state /private/state/checkout --input adopt.json --json
sdptool /path/to/project governance session-read --state /private/state/checkout --instance INSTANCE_UUID --json
sdptool /path/to/project governance session-resume --state /private/state/checkout --input resume.json --json
sdptool /path/to/project governance session-list --state /private/state/checkout --json
```

`context.json` contains `bindingPath`. Adoption takes schema, projectRef,
checkoutId, a fresh UUID operationId, bindingPath, the observed inputRevision and
reason. Resume takes schema, projectRef, checkoutId, a new operationId, instanceId,
expectedRevision and reason. Retry the same request/operationId after a lost reply;
a different request under that ID is refused. Adoption returns revision 1 and
`registration: pending`; resume appends the exact registration journal and returns
revision 2/committed. Further resumes only checkpoint current coverage. They do not
reattach to a host or advance workflow stages.

The six corresponding `governance.session.*` JSON invocation actions share this
core/CLI boundary. Context/adopt/resume parameters use `requestFile`, plus optional
`state`; read takes `instance`, list takes optional `state`. Relative invocation
paths resolve from the project. Writes require a trusted injected controller
principal; local `invoke` can attribute it with `--as ACTOR --authority controller`.
This is local account authority, not owner acceptance. These actions are not added
to the restricted coordinator MCP capability set in this increment.

Coverage includes canonical cards/plan/Session/authority, binding, instructions,
manifest and the complete management ledger. C1 conservatively blocks resume on
any changed ledger, even an unrelated append; scoped reconciliation belongs to C2.
Only the exact own registration postimage is recoverable automatically. External
edits are preserved and reported as stale coverage. One instance per Session is
enforced in the selected local store, not across all clones or arbitrary stores.
Limits are 64 KiB binding, 32 steps, 16 cards per step/64 distinct cards, 32 instances
and the existing shared 1024-receipt capacity. Old readers refuse profile 0.13.

Tests cover strict decoding, DAG/version/record guards, stale revisions, retry and
crash recovery, untouched legacy records, actual subprocess CLI/JSON parity, schema
captures and a disposable pg_lab fixture. Full lifecycle/model/host recovery,
privileged leader MCP, independent product review and owner trial remain open.

## Checked step selection — PGF1-C2a

After registration, `session-select --input selection.json` prepares one selected
step and its canonical record effects. The request contains schema, projectRef,
checkoutId, a fresh operationId, instanceId, expectedRevision, stepId and reason.
The `governance.session.select` invocation action takes requestFile/state and
requires separate trusted controller attribution. Its result explicitly reports
`recordEffects: prepared`; no card has moved yet. Inspect with session-read, then
session-resume at the returned revision to apply the exact recorded effects.
Identical-operation retries return the original receipt, including its historical
prepared status; read reports current effects/status.

Selection moves backlog/queued cards to active/ready and activates a planned plan.
It preserves already ready/in-progress cards, records whole/contribution scope and
updates worklogs, current selected-document links, board index and canonical history.
Unsupported references outside the selected card/plan/Session/index set refuse the
move before effects; arbitrary historical documents are not silently rewritten.
A committed selection has `recordEffects: committed` and `selected: true` on its
step. Selection alone stays at admission; explicit assessment is described below.
Execution requires the later assessment/start operations below; completion remains
unavailable. A gate-review card cannot be taken over
by this operation.

Operational profile 0.14 stores exact preimages/postimages and selection coverage.
The shared canonical ledger lock coordinates current PG publication, Session effects
and projecthistory/Blueprint appenders. A durable Ledger.ndjson.pg-pending marker
blocks cooperating readers/appenders across an interruption. Recovery compares both
move paths and resumes each exact effect. Source removal follows destination creation;
all canonical postconditions must validate before the instance advances. Do not
remove the marker to bypass a conflict. Shell editors and older independent binaries
are outside this coordination protocol; external edits are detected optimistically,
not protected by an operating-system transaction.

After committed selection, a validated append-only unrelated ledger change can be
observed and checkpointed without accepting edits to selected documents. A changed
selected subject/reference, a changed historical prefix or new Blueprint history
holds continuation. Prepared effects still require their full exact input set; this
increment does not rebase a pending intent. Rebind/replan remains later work.

Selection publishes a scope claim under SDP/Sessions/execution/claims. Other current
selectors refuse overlapping writer/reader scope, card, Session or plan claims;
a second store cannot register the same selected Session. C2b2 also checks these claims in current legacy admission/dispatch. This is
cooperative coordination, not interception of arbitrary shell writers or other checkouts.
C2a permits one dependency-free selected step per instance. Claims persist without
expiry until later verified disposition/replan support; do not delete them to bypass
ownership. Limits: 128 claims, 128 effects/4 MiB content and 4096 reference-scan
Markdown files, in addition to existing store limits. This is a development checkpoint,
not the full C2 start/progress/disposition engine or privileged leader MCP integration.


## Explicit Session design assessment — PGF1-C2b1

After selection effects commit, `session-assess --input assessment.json` prepares
an assessment of that selected step. `session-resume` applies its canonical records
and commits the route. Both commands use the same revision/receipt/writer protocol
as selection. `governance.session.assess` is the corresponding trusted-controller
JSON action; restricted worker/leader MCP does not gain this authority.

The request extends the selection fields (`schema`, project/checkout, operation and
instance IDs, expected revision, step ID and reason) with:

- `criteria`: exactly one entry for contracts, data-model, responsibilities,
  state-flow, interfaces, persistence and deployment. Each has `impact`
  (`changed`, `unchanged` or `unknown`) and a nonempty `rationale`.
- `blueprint`: separate `requirement` (`required`, `not-required` or `unknown`)
  and rationale. A no-design decision does not waive a required Blueprint.
- `evidence`: one to eight project-relative path/SHA-256 snapshots.
- `skills`: project-local absolute SKILL.md paths/digests for sdp-project-leader
  and sdp-change-analysis or sdp-architect. Up to eight distinct skill snapshots
  use the existing SkillInput contract. Use canonical source paths.

Any changed criterion selects design. Explicitly writable .design/.sdui files or
SDP/SDL source paths also require design, regardless of negative prose criteria.
Otherwise unknown impact or unknown Blueprint applicability keeps admission held.
All-unchanged criteria route to implementation with an explicit no-model disposition;
required Blueprint handoff remains a hold. **No route dispatches an agent in this
checkpoint.** Implementation is a business stage, not an observed running host.
SDL WORK/design handoff and model release remain required for design-changing work.

Operational profile 0.15 retains at most eight assessment intents per instance.
Canonical records live at `SDP/Sessions/execution/assessments/<operation-id>.json`,
with Session, plan and selected-card worklogs and append-only management events.
CardState remains ready/in-progress: assessment is neither start nor completion.
The artifact contains the attributed request, exact binding digest and derived
impact. The controller supplies judgment; PG validates structure, pins and routing,
not semantic truth or actual host consumption of a skill.

A later assessment can resolve uncertainty against current pinned inputs, keeping
earlier records unchanged. Edited selected documents/skills/evidence require explicit
reconciliation; reassessment cannot silently rebind them. The eight-record bound
refuses further assessment without discarding history; rollover/replan is not yet
implemented. Pending effects use exact pre/postimage recovery and the canonical
reservation. Even a historical resume retry cannot clear a newer pending reservation.

Existing 0.13/0.14 stores remain readable and gain 0.15 only on assessment. Older
executables refuse a 0.15 store. The workflow definition/digest stays pinned. CLI
and action schemas expose committed assessment, stage and holds. Existing run/attempt
ownership remains unchanged; C2b2 integrates recorded start/progress as described
below. C2c disposition and PGF2/PGF3 leader/model
gates, independent review and whole-workflow acceptance remain open.


## Session-bound start and progress — PGF1-C2b2

`session-start --input start.json` prepares a selected, assessed assignment.
`session-resume` commits every selected card as in-progress, journals the Session
and plan, appends canonical history and creates exactly one ready run in the same
operational commit. No run exists while start effects are pending. Launch is still
an explicit operation of the existing controller; no host call runs under the
canonical writer lock. Whole cards and bounded contributions retain their distinct
binding and all selected cards are included in controller context.

Start uses schema/projectRef/checkoutId/operationId/instanceId/expectedRevision/reason,
plus the exact latest `assessmentId`, an `outcome`, and `skills` snapshots containing
sdp-master and sdp-worker. Each snapshot uses the existing SkillInput name, absolute
canonical project-local SKILL.md path and digest. The controller launch must supply
those same snapshots through its existing skill-input mechanism. Use the returned
step `runId` for existing run inspection, launch, reports and candidate evidence.
The run identity is the start operation ID; the `session-implementation` routine
uses the existing run/attempt engine with an immutable Session binding. It does not
relax legacy Maintenance or approved-Blueprint admission.

Only all-unchanged design criteria with Blueprint explicitly not-required may start.
Unknown/changed design and required/unknown Blueprint remain held for PGF3. Start
refuses stale assessment, changed skill/evidence/selected record pins, overlapping
current run/prepared scopes, and writes overlapping its own pinned governance input.
Current legacy admission/dispatch observes persistent Session claims optimistically;
this is not a universal cross-store or filesystem lock. Nonliteral conflicting
scopes require reconciliation. Claims still have no automatic expiry/release.

After the assigned run reports, capture candidate evidence using the existing
controller evidence API. `session-progress --input progress.json` takes the common
identity/revision/reason fields plus `evidenceId` and `expectedRunRevision`. It
requires that run's latest assigned report, exact evidence receipt/artifact/candidate,
known terminal host attempts and no unresolved outbound input or blocking attention.
A local controller cannot impersonate the assigned report. Preparation and each
canonical effect recheck the candidate. Resume commits progress and exposes the
verification business stage; cards remain in-progress. Test results can be failed:
submission for verification is not a passed review, acceptance or completion.

The JSON actions are `governance.session.start` and `governance.session.progress`,
using requestFile/state and separate trusted-controller attribution. Restricted
worker/leader MCP does not gain these controller actions. The strict wire decoder
rejects cross-verb fields, even empty ones. Identical retries return historical
receipts; inspect current state with session-read. Uncertain host sends remain
unknown and are not automatically repeated.

Operational profile 0.16 adds one start and one progress intent per instance;
older stores remain readable and older executables refuse the new store. Immutable
canonical records live in `SDP/Sessions/execution/work/<operation-id>.json`. Existing
pre/postimage recovery, reservations and input coverage apply, including unrelated
validated ledger appends. Full binding, run state, progress evidence and actual host
execution are separate view fields. Generic rework, further advance/relaunch after
progress and legacy single-card result publication cannot bypass the pending Session
disposition. C2c must add verified whole/contribution disposition, rework/dependencies
and claim release/replan; design/model release, privileged leader integration, native
whole-workflow trials and independent review remain required.
