# ProjectGovernance — first workflow contract

Design candidate PGD1, owned by [PLAN-SDP-0017](Plan.md). This specifies the bounded
implementation handoff. PGI1/PGI2 implement the bounded core/CLI and protocol adapters; the
PGI3 controller candidate adds explicit covered-run submission and steering, while
its live workflow verification and installed process integration remain open.
Owner-selected scope and role responsibilities are retained. Storage locations,
operation names and implementation choices below are the recommended pilot design.

Current implementation details and limits are in the
[core contract](../../../../SDPTool/governance/README.md) and
[PGI plan](../../../05--Implementation/SDPTool/ProjectGovernance/Plan.md).
The application wire schema is local to the PGI pilot; it is not a published Toolkit profile.
PGD1's original evidence remains historical design/schema-inspection evidence.

## Current owner-facing role refinement

Owner Session0007 T014 selects [KB052 orchestration](Orchestration.md).
Project Leader replaces the earlier combined Steering/PM label for normal owner
dialogue; Steering remains strategic. Owner T039 additionally requires the
[mandatory admission and lifecycle contract](Orchestration.md#t039--mandatory-governance-admission-and-lifecycle-supervision):
card-derived Session/plan coverage, SDL WORK for design changes, supervised
integration and model-release disposition before whole-assignment closeout. The original single-Maintenance journey below
is retained as a foundation, not the complete new operating model.

## 1. Outcome and first user journey

The owner opens an explicitly selected project and its existing Session, submits
a bounded change, and sees which routine and work records cover it. A fresh task
Master receives that assignment and surrounding constraints. The owner can inspect
actual progress, correct scope, resume after interruption and see the evidence
behind the result. Project memory survives replacement of the Codex conversation.

The first executable routine is bounded Maintenance on an isolated SDP project.
Inspect/status, resume/reconcile and procedure-gap handling are supporting paths.
The sixteen-family [catalog](../../../02--Requirements/RoutineGovernance/Routine-Catalog.md)
remains a coverage map; this pilot does not claim to implement all families.

1. Select a physical project/checkout and recover current records and Session.
2. Submit a request to the project-facing Project Leader conversation. Record the
   actual owner input and a pending routing result; an agent proposes its category
   and basis. Context reads are allowed while selection remains unresolved.
3. Resolve existing card, plan, milestone and authorization. The core checks the
   proposal against those records; it does not pretend to understand free text
   through deterministic rules. Reuse covered work. A missing plan is an explicit
   prerequisite, not permission to fabricate one or silently begin code changes.
4. Create a version-pinned run and bounded assignment; launch a fresh Master
   attempt with deliberate context. The Master coordinates permitted execution,
   review and return. Native children are preferred, subject to measured host
   capability and assignment authority.
5. Show routine step, responsible assignment, actual Codex status, pending inputs,
   evidence and observation freshness. A tool event changes activity, not acceptance.
6. Return candidate-bound results, verification and applicable independent review;
   complete covered work under existing authority and reconcile the Session.

A pure status question uses the read path and creates no card, plan or routine run.
Its visible message may still belong to an existing Session transcript. An unrelated
question does not start a ProjectGovernance Session. Compound requests retain their
parts and coverage; unrelated new work is recorded separately without abandoning
the current objective.

## 2. Ownership and implementation boundaries

| Responsibility | Owner / proposed implementation | Existing behavior to preserve |
| --- | --- | --- |
| Method, selection and actual owner decisions | Existing SDP documents and owner | Agent statements are not new authority |
| Routing checks, routine versions, runs and transition validation | Go package under SDPTool, shared by adapters | Discovery and existing command output contracts remain compatible |
| Card/plan lifecycle | ProjectManagement documents and Ledger.ndjson | Existing event schemas, IDs and predecessor chains |
| Candidate evidence and design/code relationships | Traceability and referenced evidence | No management-only event duplication |
| Owner conversation and execution observations | CodexClient app-server adapter | Host sandbox/approval rules are not weakened |
| Agent access to context and bounded operations | Thin local SDPTool MCP adapter | No duplicate workflow rules or generic ledger append tool |
| Goal, steps and cross-turn continuity | Session record and its versioned core representation | Existing manual Sessions remain readable and usable |
| Model history and promotion | Separate ModelGovernance workstream | Reference source identity/digest; no model store implementation here |
| Timeline and diagram | Later KB043 projection | No timing, compliance or completion inferred from animation |

Start with one local controller and a text terminal client. Reuse SDPTool Go
packages for domain logic; provisionally place the client and app-server adapter
in the same Go module, with separate commands for client and MCP entry points.
An interactive framework is not needed for the first line-oriented command/status
view. A separate read-only terminal view reads snapshots. A full KanBan TUI and
native GUI are outside this pilot.

Alternative choices considered: a Codex fork increases upstream maintenance; a
passive monitor alone cannot reliably own incoming requests; a client-only workflow
engine duplicates policy across CLI/MCP. The selected candidate uses a standalone
client and shared core. No second service or network port is required initially.

## 3. Identity, revision and authority

The pilot application schema is proposed as `sdp-project-governance/0.1`, separate
from installed method, management payload, routine definition and transport versions.
Unknown major/profile versions refuse mutation. Read-only diagnostics may report
unsupported records without rewriting them.

| Identity / reference | Meaning and lifetime |
| --- | --- |
| projectRef | Declared project identity plus resolved SDP root; identity alone cannot authorize a path |
| checkoutId | Locally allocated UUID bound to a canonical physical root; different checkout means different binding, even for the same Git repository |
| sessionRef / stepId | Existing Session reference plus a stable local step; independent of Codex sessionId |
| requestId | UUID allocated when the controller receives one submitted/steered input |
| routineRef | ID, version and exact definition digest; immutable for an existing run |
| runId / revision | UUID and monotonically increasing core revision for one execution |
| assignmentId / assignmentRevision | Bounded scope, role, allowed work, references and return criteria |
| attemptId | One execution attempt; replacement thread creates a new attempt, not a new assignment |
| threadId / turnId / itemId | Opaque Codex identifiers observed from the selected protocol |
| operationId | Caller-generated UUID for durable retry correlation; unrelated to JSON-RPC request id |
| candidateRef | Relevant exact source/artifact inventory, optional commit, dirty-input hashes and evidence boundary |

Existing Session IDs are imported as explicit references without renaming documents.
Register canonical root and checkout binding before mutation. Reopening from another
path must resolve the same binding or require explicit selection; never guess from
the most recently touched repository. Git is useful candidate provenance, not a
requirement for Session, run or assignment identity. Symlink/path escape and ambiguous
project resolution reject the operation.

An assignment contains outcome, selected work IDs/milestone, governing source paths
and digests, writable boundary, preserved invariants, excluded work, required evidence,
review independence, authority references and procedure-gap/return behavior. Model
and blueprint references are optional enrichments with availability status, not
fabricated complete context. Refresh references when inputs change.

Owner input arrives through the controller's user channel; agent proposals arrive
through assignment-bound adapters. A `role=owner` payload, MCP tool result, Codex
approval or successful process exit never creates an owner decision. Existing
authorization is reused with its source and scope; the interface asks only about
an unresolved material choice or genuinely ungranted action.

The launch binding limits an agent to a project, assignment, role capabilities and
attempt. Credentials/bindings are never written into project documents or observer
output. The adapter must not accept a caller-selected identity in their place.
Native child identity and reviewer isolation require a real compatibility test.
If the chosen host cannot establish distinct trusted provenance, show that result
as an agent claim and refuse an independent-review transition. Do not silently
replace the desired role hierarchy or claim security against arbitrary same-user
shell access. The pilot governs its APIs; it cannot intercept every external edit.

Expose agent context, route proposals, assigned run reads, evidence claims and
permitted transitions through MCP. Keep owner-input recording, launch binding,
owner disposition and unrestricted project publication on the trusted controller
side. A review-capable binding still needs the selected review assignment and
candidate; it cannot review its own implementation assignment. Child completion
is aggregated only when every required child has a valid return disposition.
Failed or optional canceled children retain their explicit outcomes.

## 4. Routine and work states

Routine definitions contain ID/version/digest, purpose, applicability, required
inputs, ordered/conditional steps, transition guards, evidence rules, permitted
roles, exception paths and compatibility. A run retains the exact definition bytes
or a content-addressed immutable copy. An update affects future runs; migration
of an active run must be an explicit later capability, initially unsupported.

Initial route outcomes: `no-work`, `matched`, `ambiguous`, `missing`, `unavailable`,
`incompatible`, `conflict`. Persist the proposing actor and stated basis. Core checks
resolve referenced records and guards; agent language interpretation remains an
attributed judgment. Missing versus unreadable versus conflicting authority must
not collapse to an empty successful result.

| Step state | Entry / permitted exit |
| --- | --- |
| pending | Prerequisites absent; become ready only when required inputs and earlier steps satisfy the routine |
| ready | Valid assignment and authority; start -> running |
| running | Report result -> verifying; need-input -> waiting-input; interrupt -> interrupted; error -> failed; scope delta -> replan-required |
| waiting-input | Record the actual missing input; supplied input is reclassified before returning to ready |
| interrupted | Execution stopped or status uncertain; reconcile -> ready or running only after observing actual host state |
| replan-required | Scope/authority/candidate inputs changed; refresh the assignment revision before ready |
| verifying | Accept required candidate-bound evidence -> completed; findings -> ready with rework; owner decision needed -> waiting-input |
| failed | Record failure and any committed effects; explicit retry/reconcile -> ready with a new attempt |
| completed | Guards satisfied for the identified candidate; later relevant change invalidates the current evidence projection and requires a successor/rework step |
| canceled | Explicit authorized cancellation; preserve prior results and remaining-work disposition |

Terminal historical transitions are immutable. A current view can mark old evidence
stale without erasing the completion event or falsely claiming the revised work is
complete. Exception/skip is a recorded disposition with reason and authority, never
a successful verification result. CardState remains its existing separate vocabulary.

The initial Maintenance routine uses resolve coverage -> prepare assignment ->
implement bounded change -> verify candidate -> applicable independent review ->
record return/disposition. Review can send work back for rework. Routine termination
does not automatically move a card whose broader capability remains unfinished.
Release/merge steps are unsupported extensions in the first pilot; it can identify
the required routine and missing authority without performing publication.

## 5. Application operations

All adapters call the same application service. Names below define domain operations,
not raw app-server methods. JSON schemas and conformance fixtures are deliverables
of the first implementation slice. Reject unknown fields, duplicate JSON keys,
wrong types, unbounded collections and unsupported schema versions.

Common mutable envelope: `schema`, `projectRef`, `checkoutId`, `operationId`,
`runId` where applicable, `expectedRevision`, `inputRevision`, and operation-specific
data. Trusted caller context is supplied out-of-band. New-run creation uses an
expected current context revision instead of a nonexistent run revision.
Core inputRevision covers selected authority/work inputs, not arbitrary UI activity.

Use UUID strings for new operational identities, nonnegative integer revisions and
lowercase SHA-256 hex for digests. Hash normalized application requests after strict
decoding, with recursively sorted object keys, preserved array order and exact UTF-8
string contents; forbid floating-point values in revision/sequence fields. Omit
transport request IDs and credentials from the retry digest. Normalize omitted
optional fields to their schema-defined defaults before computing it.

| Operation | Inputs / output | Side effects and guards |
| --- | --- | --- |
| context.read | Explicit root -> resolved work, authority, source digests, capabilities and coverage | Read-only; reuse discovery; do not start processes or render views |
| route.evaluate | Bounded intent proposal, requestRef, context/catalog revisions -> route outcome and unmet prerequisites | Read-only evaluation; does not select/activate work |
| request.record | Observed owner input reference and route basis -> requestId / stored disposition | Controller-authorized operational record; no inferred owner approval |
| run.start | Selected work, Session/step, pinned routine, authority and assignment -> run snapshot | Reject unsupported profile, unresolved coverage or changed input revision |
| run.read / events.read | Scoped run and optional cursor/limit -> snapshot or ordered observations/transitions | Read-only; cursor coverage and gaps explicit |
| assignment.prepare | Run and bounded context -> assignmentRevision and launch binding | Master coordinates work within already granted delegation authority |
| attempt.bind | Assignment revision plus observed host identifiers -> attempt mapping | Controller operation; retain prior attempts |
| evidence.attach | CandidateRef, evidence type/location/digest, issuer and limits -> attributed evidence reference | Record claim versus observed test versus independent review separately |
| run.transition | Routine action, candidate/evidence/authority references -> accepted state or refusal | Compare revision and guards; caller cannot provide arbitrary state patches |
| run.reconcile | Run, host observations and current source revisions -> recovery disposition | No blind resend, model invocation or completed-state inference |
| session.project | Session source revision and accepted run outcome -> checked summary proposal | Rendering is separate from publication; preserve authored roadmap/history |
| records.publish | Explicit proposed changes to existing Session/card/plan/evidence records with preimage hashes -> receipt | Authorized bounded reconciliation; current schemas and ownership apply |

Read responses include schema, selected project/checkout, revision, observedAt,
sourceRefs and coverage (`known`, `stale`, `unavailable`, `unsupported`). Mutable
success returns operationId, accepted revision, disposition, references and event
cursor. Refusal returns code, retryable, currentRevision and unmetConditions.
Stable error codes include unauthorized, missing_coverage, stale_revision,
stale_candidate, evidence_missing, incompatible_profile, idempotency_conflict,
busy, observation_gap and reconciliation_required. Do not turn refusals into
successful empty results. Malformed transport is distinct from a domain refusal.

Initial bounds: 1 MiB application request, 64 referenced artifacts per operation,
256 events per page and 64 KiB per routed intent summary. Larger source/transcript
artifacts use validated references with explicit size/digest; never silently
truncate evidence. A local read can inspect a bounded authorized artifact, not
an arbitrary path supplied by a model.

## 6. Persistence, concurrency and publication

Use a private local state directory under the user's XDG state home (default
`~/.local/state/sdptool/project-governance/<checkoutId>/`), with 0700 directories
and 0600 files. It is durable operational data, not an evictable cache. Store
root binding, exact routine definitions, run snapshots, operation receipts,
attempt mappings and allowlisted observations. Persist Session references to
project documents. Moving to another machine recovers those documents, but full
in-flight runtime recovery requires transferring the state store; do not claim
automatic cross-machine or Git synchronization in this pilot.

This operational store owns run/attempt activity only. It references the existing
management ledger and Traceability for their facts; it cannot override a changed
card, plan or acceptance record. Do not append unsupported Session events to the
current management ledger. Storage/profile extension is explicit pilot work, not
silent migration of existing installations or old Session files.

Serialise local core mutations with a checkout-scoped OS lock shared by all adapters.
Under the lock: validate expected revision and relevant preimages, look up
operationId, prepare a complete next immutable generation, flush its files and
directory, atomically replace the head pointer, flush the parent, then acknowledge.
Each generation contains its predecessor, snapshot and accepted operation receipt.
Same operationId and normalized input digest returns the recorded outcome; different
content with the same key refuses. A crash before head publication leaves an orphan
generation; after publication the receipt permits retry without repeating effects.
Recovery validates the referenced generation before allowing further mutations.
Do not promise network-filesystem or multiple-user writer guarantees initially.

Project document/ledger reconciliation is a separate recoverable operation. Prepare
old/new byte hashes, exact append bytes, event IDs and resulting file contents before
writing. Recheck them under a project-writer lock and persist intent. Each step is
idempotent: old hash -> apply; new hash -> already done; other bytes -> conflict.
Verify exact ledger event presence, not merely its ID, before skipping an append.
Commit a publication receipt only when all required documents agree. On restart,
resume an incomplete publication once, or report reconciliation_required. No
multi-file atomicity is claimed; incomplete state is visible and blocks dependent
completion. Existing writers must participate in the lock for strict exclusion;
detect observed changes from external writers and never overwrite them knowingly.

Publication preserves manual Session content and historical ledger bytes. In the
pilot, write a bounded appended journal summary and explicit checked roadmap/card
changes; do not rewrite the whole document from chat. A moved card resolves through
stable ID plus the current ledger path. The operation validates the complete proposed
record set using current management/traceability rules before reporting success.

## 7. Codex integration, capture and recovery

Use one controller-owned app-server process over local stdio, explicit selected cwd,
and version-specific generated protocol bindings. The inspected baseline is
codex-cli 0.160.1 for PGI2, default schema surface, with MCP Go SDK v1.8.0
and MCP 2025-11-25. [PGI2 evidence](../../../05--Implementation/SDPTool/ProjectGovernance/Evidence-PGI2.md)
records actual initialization/account-mode inspection and ephemeral-thread MCP
discovery, plus fixture capture/recovery limits. PGD1 [Evidence.md](Evidence.md)
for 0.160.0 and earlier 0.158.0 evidence remain historical. Unsupported versions fail with
a compatibility diagnostic until tested. Do not require undocumented experimental
fields or interpret the documentation's newest example as the installed contract.

Owner Session0007 T029 selects PGC1 as the current admission rule: stable Codex
0.164.0 or later under the bounded existing wire contract. This supersedes exact
version admission, preserving the prior measurements as history. Executable and
initialize versions must agree; persist the observed connection version per attempt.
Required operation/response changes refuse explicitly, and future eligibility is
not a claim of native workflow verification. Keep old attempt histories readable;
recovery may not relabel them onto a different connected host version. See the
[PGC1 contract](../../../05--Implementation/SDPTool/ProjectGovernance/Codex-Compatibility.md).

Initialize once per connection, create/resume mapped threads, then submit/start or
steer the identified turn. Generate clientUserMessageId before persistence/submission
where the pinned schema supports it; it is correlation, not a proven server-side
exactly-once guarantee. Persist outbound intent before sending, and acknowledged
IDs afterward. After an ambiguous send, inspect available thread history and prompt
identity; if delivery cannot be established, show uncertain and require deliberate
reconciliation instead of issuing another turn automatically.

Managed ChatGPT authentication remains the requested mode. At pilot startup inspect
the supported account status without copying credentials; never silently switch to
API-key billing. Missing authentication should lead to the documented sign-in path.
No sign-in or model entitlement was exercised during PGD1. Native host permissions
and approval requests remain host responsibilities. The client's SDP decision UI
must clearly distinguish them from product review and owner-reserved decisions.

Capture only visible owner messages, assistant commentary/final messages, allowlisted
tool/agent activity metadata and completion/error state required by the workflow.
Exclude reasoning, credential material and authentication payloads from the Session
and observer even if the protocol exports them. Attachments retain references and
availability; unavailable content is not replaced by invented text. Unknown message
phase remains unknown, and an interrupted turn need not have a final response.

Use (threadId, turnId, itemId, event kind) to reconcile item lifecycle where supplied;
preserve actual item updates. Streaming deltas have no assumed globally stable
deduplication key: display them provisionally, consolidate from authoritative item
completion/history when available, and record gaps on reconnect. Never concatenate
replayed deltas into a second final response. Conflicting repeated completed items
are an integrity/coverage problem, not a silently overwritten transcript.

Record observedAt and local receivedAt separately. Use local monotonic offsets only
within one controller process epoch. Reconnects and clock resets break elapsed-time
continuity; mark that break. Host event times, ledger write times and active effort
are different facts. KB043 can later render this data without inventing intervals.

If the observer closes, execution continues. If the controller or server disconnects,
mark execution status unknown/interrupted as observed, preserve durable work state,
and reconcile before launch. A turn/completed notification may mean completed,
failed or interrupted; even success only records execution outcome. Scope-changing
input moves affected work to replan-required, requests an interrupt where necessary,
and reconciles any already committed effects before preparing a new assignment.
The client cannot roll back a running external command by changing a state label.

## 8. Delivery boundary and open compatibility evidence

No unresolved product choice blocks implementing the pure core and fixtures under
this candidate design. Runtime dependencies must be demonstrated before the integrated
pilot is called delivered: supported MCP/SDK negotiation, managed account mode,
native child/context provenance, independent review and disconnect/history recovery.
Keep adapters behind explicit capability results and do not substitute unverified
fallbacks silently. [Acceptance.md](Acceptance.md) gives the concrete scenarios;
the implementation handoff sequences the checks before wider execution.

This design does not introduce a workflow editor, distributed scheduler, all routines,
model-store backend, semantic blueprint engine, full timeline UI, remote controller,
Codex fork, release automation or repository-wide write interception.


## PGI3 implementation refinement — 2026-10-07

The first controller surface operates an existing covered Maintenance run; it does
not claim the complete Steering/PM conversation or automatic classification. Durable
owner input records exact text and the bounded assignment payload before sending,
using prepared/sending/acknowledged/uncertain states. Unresolved sends block new
launch/steer requests. A prior nonterminal host attempt requires reconciliation.
This is the no-blind-resend obligation above, not a server exactly-once claim.

The current host prohibits agent delegation and restricts network/filesystem access.
The candidate launches a fresh Master assignment without granting native delegation;
child/reviewer attribution remains unavailable and no independent review is claimed.
The live probe currently ends with EOF before initialization, so fixture success
does not deliver PGI3. History recovery records observed state only and cannot infer
or retry an unmapped ambiguous turn. See the [controller contract](../../../../SDPTool/governancecontroller/README.md).


## PGI3 measured refinement — 2026-10-07

The installed 0.161.0 host successfully executes the scoped Maintenance fixture
and receives steering. History includes a userMessage.clientId correlation value.
For a lost outbound acknowledgment, the selected recovery refinement is to bind
only a unique recorded user message whose clientId and exact payload match the
persisted input, within the already mapped thread/attempt. Persist that observed
message before resolving the input. Missing, conflicting or unmapped history stays
uncertain; never infer acceptance from text similarity or resubmit to discover it.
This derives from PG-A08/11/12 and the existing no-duplicate-dispatch contract.

Host approval refusals remain refusals. The live fixture's MCP report required
approval under the selected never policy. Do not change that policy or claim
routine acceptance from a host-completed turn. PGI3 must surface the limitation;
a future explicit owner-interaction path can handle requested approvals, while
PGI4 continues to own evidence and publication. No automatic approval is selected.

## PGI3-C3 refinement — explicit MCP approval, 2026-10-08

Owner continuation after C2's recommendation selects implementation of the explicit
request/response interaction. This authorizes building and testing that flow, not
blanket acceptance of future tool calls. Keep the default refusal mode; an explicit
`--approvals ask` opts a new client launch into owner-reviewed MCP requests.

The installed 0.161.0 host was measured before implementation: run_transition uses
mcpServer/elicitation/request, mode form, an empty object form, and metadata
codex_approval_kind=mcp_tool_call with exact tool_params. Support only that bounded
SDP server/tool shape, with exact thread, turn, transport epoch and typed request ID.
Present its message and full transition arguments as data; never execute display
text. Accept, decline and cancel are per request; omit persistent/session grants.
Unrelated permission, URL and secret-input requests remain unsupported. No account
or global configuration changes, native delegation or automatic approval review.

Persist observed requests and exact owner decisions before sending. Duplicate,
resolved, cross-attempt and old-epoch replies refuse. A written response alone is
not host acceptance: keep its outcome uncertain until serverRequest/resolved is
observed; that notification proves closure only, not successful tool execution or
SDP acceptance. Loss of transport never resends a response automatically.

Approval capture itself advances the run observation revision. A queued MCP call
therefore needs a narrow revision rule: only the exact observed and explicitly
accepted agent transition may cross revisions consisting exclusively of host
observations and this approval's input bookkeeping, with the same live attempt,
coverage and routine decision state. Any intervening steering, lifecycle mutation,
changed arguments or external coverage change refuses. Controller calls retain
strict revision matching. This is required to avoid invalidating every approved
call merely by recording its approval; it does not grant new agent actions.

Verification: reject unsupported shapes and stale/mismatched/duplicate answers;
exercise accept/decline and transport loss through the compiled controller; measure
an actual managed-account MCP report through the final binary. Driver approvals in
disposable tests are scenarios, not attributed owner conversation. PGI3's other
reconnection/provenance obligations and PGI4/PGI5 remain separate.

Compatibility refinement: actual initialization rejects the exported granular
approval shape without experimentalApi. Use the measured stable on-request policy
with approvalsReviewer=user for explicit ask mode; the adapter rejects all request
types outside the bounded SDP MCP form. Do not enable experimental API or silently
fall back. Default never remains unchanged; sandbox selection is independent.

The first actual response also exposed concurrent core lock contention between
MCP execution and approval observation. Retry only busy/stale local bookkeeping
that refused before mutation, bounded to 100 attempts with 20 ms delays. The MCP
adapter retries busy only, retaining exact arguments and operation ID. It does not
retry stale semantic operations or host sends. Persistent contention stays explicit.

## PGI3-C4 refinement — executing-turn reattachment, 2026-10-08

PG-A11 requires rejoining execution, not merely reading its eventual history. The
installed 0.161.0 host was measured with an isolated Unix WebSocket listener:
a command continued after client closure, and thread/resume subscribed a fresh
connection to the same in-progress turn and its eventual completion. The raw
app-server proxy does not translate JSONL into WebSocket frames. An idle thread
without a persisted turn could not be resumed; do not infer recoverability from
thread/start alone. Official transport reference:
https://learn.chatgpt.com/docs/app-server

Select an opt-in private local Unix socket, with an explicitly separate foreground
Codex server. Keep owned stdio as the default. Use coder/websocket v1.8.15 for
framing rather than adding a custom execution supervisor. No shared daemon, TCP
listener, global configuration, experimental API or automatic server restart is
selected. Losing this client closes only its connection. Server shutdown remains
an explicit execution boundary.

Bind new socket attempts to the local socket/peer process identity and measured
host version. Reject another endpoint, replacement server, mismatched thread/turn,
source drift or unresolved outbound input before reattachment. Subscribe using
thread/resume with only the mapped thread ID; never send turn/start or overrides.
A running snapshot establishes status only: incomplete historical items are not
completed messages. Retain the observation gap and replay completed history when
available; do not claim uninterrupted coverage.

Separate connection epoch from agent capability epoch for this measured path.
Only explicit same-host reattachment preserves the launch capability; ordinary
history recovery retains its existing revocation semantics. Existing attempts
without a host identity cannot use reattachment. The original agent remains scoped
to the same latest attempt, assignment and source revision. Old client observations
and approval replies remain fenced. Unresolved decisions require reconciliation;
never resend them on reconnect. New approval requests need fresh explicit replies.

Verify a real disposable Maintenance task through final binaries, disconnect while
executing, attach with a new process, exercise owner steering and a subsequent MCP
report, and establish one run/attempt/thread/turn with preserved excluded files.
Negative service/transport tests cover identity substitution and capability fencing.
This checkpoint does not establish child provenance or independent review.

Actual controller trials revealed that this host creates an absolute socket alias
in the requested directory, pointing into a private per-user runtime directory.
Both initial trials refused before any thread submission. The selected refinement
permits that direct canonical alias only when both parent directories and the
socket are private and user-owned; indirect aliases and symlink parents refuse.
Peer process lifetime and socket identity remain bound. No default daemon is used.


## PGI3-C5 refinement — bounded launch provenance, 2026-10-09

Owner continuation selects applicable role/skill/context evidence under PG-A04,
PG-A10/12/15. The installed 0.161.0 protocol exposes thread identity/source/native
role fields, thread/start instructionSources, and skills/list availability. These
are different evidence levels: controller-assigned Master is not a host-native
role, a host-reported instruction path is not a digest of bytes demonstrably used,
and a discovered enabled skill is not evidence of invocation or compliance.

Select a durable controller-owned context observation and a read-only provenance
projection. Capture the fresh thread's reported identity and instruction paths,
current governing source references and digests, and bounded project-local skill
availability from skills/list with forceReload. Hash readable project-local skill
files as a local snapshot; do not read external skill contents or infer activation.
Keep unavailable metadata, host discovery errors, external inventory counts and
unobserved skill use explicit. Observations confer no authority or review rights.

Preserve existing launch text, input correlation, retry and recovery semantics.
This increment does not add skill inputs, automatic skill selection, new MCP tools,
role switching, child dispatch or review acceptance. A snapshot remains as-of its
capture epoch after reconnection; do not relabel it as current runtime context.
Use actual project skills in the disposable trial and retain returned claims as
agent statements. Full skill invocation/activation telemetry remains unknown when
not exposed by the measured adapter. Native child execution remains unavailable
under the current no-delegation assignment; it must not be simulated as proof.

The alternative of parsing shell commands or model prose as proof of skill loading
is rejected: neither establishes exact consumed context or compliance. The explicit
skill-input protocol is available for a later selected invocation workflow, which
must also preserve lost-response correlation for non-text inputs. C5 does not
silently change the existing text-only dispatch contract to get stronger claims.

Verification: actual fresh managed-account Maintenance execution with installed SDP
skills, observed instruction sources and project-local inventory, bound context
receipt and unchanged exclusions. Negative tests cover forged roles/agent writes,
oversized records, unknown/missing metadata, invalid paths and inventory limits.
Read-only provenance must work without connecting to Codex. No independent review
or full role hierarchy is claimed.

## PGI3-C6 refinement — explicit skill input and correlation, 2026-10-09

Owner continuation selects the next C5 obligation under PG-A04/10/11/12/15/18.
An actual read-only Codex 0.161.0 probe returned the same ordered text and skill
name/path content with clientId in both completed notifications and full history.
The host adds empty text_elements to text. No loading/compliance telemetry was
observed; this is accepted input evidence, not consumed-byte attestation.

Select repeatable `run --skill PROJECT_RELATIVE_SKILL_MD` for an explicit bounded
list (at most eight) of enabled, uniquely named project skills. Resolve through
fresh skills/list, require a readable local digest and refuse missing, disabled,
ambiguous, duplicate or external selections before creating an intent/thread.
Default launches stay text-only; do not infer selected skills from prose, role
names or all enabled skills. The controller adds explicit name markers to the
outbound text while preserving the owner's exact text separately, and persists
selected name, canonical absolute path and local digest before host dispatch.
Recheck local bytes before claiming and immediately before turn/start. There is
no atomic host-consumed digest guarantee; later changes cannot invalidate history
of a previously accepted input or retroactively establish compliance.

Persist a fingerprint of the exact ordered input content on observed user messages.
Canonicalization supports one text followed by selected skills; absent or empty
text_elements are equivalent, nonempty/unknown content is not silently dropped.
Resolve an uncertain skill send only with unique clientId, exact ordered content
fingerprint and mapped/full history. Do not accept matching text with changed,
missing, reordered or additional skills. Preserve legacy text-only observations
and refuse uncertain resends. Skill content is immutable across input transitions;
read-only provenance projects selection and response acknowledgment/observed host
message separately from the availability snapshot. No agent claim establishes use.

A text-only marker was rejected as insufficient explicit skill identity. Parsing
model statements or shell reads as compliance remains rejected. Arbitrary external
skill loading and skill steering are outside this bounded increment. Existing
sandbox, permissions, scoped MCP tools, recovery gaps and no-delegation instructions
stay intact; native hierarchy and independent review remain unverified.

Verify malformed/missing selections, drift, duplicate/changed content and exact
lost-response recovery; test real skill-equipped Maintenance with explicit report
approval and preserved fixture files. Retain actual protocol and workflow evidence,
including a deliberate response-loss trial without resubmission.

## PGI4-C1 — Candidate-bound evidence and rework, T035

Refines existing PG-A06/07/08/16/17/18 within the selected completion route.
A report is a claim awaiting verification. The trusted local controller may capture
an exact source inventory and attach an attributed test, review or general claim;
an attachment is not proof that PG executed a test or observed an independent reviewer.

Use the existing atomic store/receipt namespace (operational profile 0.10), rather
than a second evidence ledger or Git HEAD as a proxy for dirty source bytes. Old
readers refuse the new profile; every current writer preserves later data. No
canonical project records are published in this increment.

Candidate capture is read-only and bound to checkout, run, input revision, exact
writable/excluded paths and explicitly selected additional dependency paths. Literal
files/directories only; no globs, symlinks, hidden components or special files.
Sorted inventory records files, executable bits, directories and missing roots.
Directory membership changes, missing/new/deleted files and relevant executable-bit
changes invalidate the digest. Unselected excluded work does not participate. Explicit read-only dependencies
may include write-excluded files (for example go.mod); this grants no write authority. Bounds: 64 dependencies, 512
inventory entries, 4 MiB/file, 32 MiB total content, 4096 bytes/path. Capture does not
infer that the selected dependency list covers every possible behavioral influence.
It observes the filesystem; external writers must be quiescent for a coherent
snapshot. Repeated capture before attachment detects observed drift, not atomic
filesystem isolation.

Attachment pins the latest exact report, run revision, authority/input revision,
candidate inventory and a project-relative artifact path/digest. The artifact must
be outside the inventory scope. Capture it before attaching; PG rechecks source and
artifact bytes before atomic publication. Retain at most 16 evidence receipts/run,
within the existing 1024 receipt/16 MiB store bounds. Exact retries return the original
receipt even after later drift; read projection separately reports current status.
Attribution is controller-supplied, never authenticated human or native reviewer
identity. Review independence remains unverified regardless of artifact prose.

Read projection keeps all historical records and separates current candidate,
artifact, authority and report applicability. A source change invalidates current
use without deleting the original result. Evidence attachment does not advance or
accept work. Failed current test/review evidence can support controller-only rework
from verifying to ready, preserving the same Session and original report/evidence.
Rework requires no open blocking attention, no uncertain input and no active/unknown
host attempt. Old agent capabilities cannot submit another report after rework;
subsequent work requires a fresh attempt. Completion still refuses until required
review and checked publication exist. Agent/coordinator MCP authority is unchanged.

## PGI4-C2 — Observed review context, T036

Implements PG-A06/07/08/11/15/16/18's independent-context review prerequisite.
Use an explicitly controller-assigned Reviewer in a fresh persistent read-only
thread, not a resumed/forked implementation conversation. Native child hierarchy
is still unverified and is not claimed by this route. A fresh-context host mapping
establishes separation within the measured local adapter, not reviewer competence,
model diversity, authenticated human approval or protection from same-user shell edits.

A review pins run/report/input identity, exact C1 Candidate, 1..8 current test
receipt IDs, bounded review focus, installed host version, local MCP executable and
explicit project reviewer skill snapshots and exact reviewer instructions. The payload carries governing selection,
assignment, candidate and referenced test records, not implementation conversation
history. At least one selected skill must be named sdp-reviewer. Other relevant
project skills may accompany it (maximum eight). Payload/result are bounded to
64 KiB; at most eight reviews/run and 64 reviews/store, under existing receipt/store
limits. Skills and source/evidence/authority must still match before opening/sending.

Profile 0.11 extends the existing operational store with review records. Standard
receipts retain exact retry semantics and shared UUID collisions. Review states are
prepared, opening, bound, sending, running, completed, failed, interrupted or
superseded. Each external effect follows a durable claim. Unknown opening has no
replacement; unsent prepared/bound work may be superseded safely. A new controller
cannot assume an unsent native bound thread survived. Sending/running must recover
the exact original host/thread/input/turn mapping, never resubmit. A per-run lease
serializes controller effects while owner messages and unrelated runs remain free.

The trusted adapter checks fresh idle thread provenance, empty history, expected
project/version and a thread ID distinct from all recorded implementation, leader
and review threads. Before sending, it checks the complete MCP inventory contains
only the configured read-only PG context/route/run tools, bound to this run without
an attempt/epoch mutation capability. Built-in tool sandbox is
read-only with approvals never; no native delegation or project publication is
requested. Native host role and controller-assigned Reviewer remain distinct.

Recovery requires exactly one full-history review turn with the unique review ID
and exact payload/skills input. Retain only the final visible answer, never private
reasoning. The final answer must decode as a bounded verdict (passed/failed/
inconclusive), summary and findings. Failed/interrupted/malformed outcomes are
retained without a passing review. Capture results even after source drift; current
qualification separately rechecks report, candidate, authority, selected test
artifacts and skill snapshots. Completed host status alone is not qualification.
A passed, applicable result establishes an observed review in a separate context;
canonical acceptance/publication remains unavailable. C1 attached prose remains
controller-attributed; it cannot manufacture this provenance.

## PGI4-C3 — Checked result publication and leader delivery, T037

The first publication operation records a passed current review and its attributed
passing tests, with a controller-supplied work summary and explicit remaining work.
It appends to the selected Session/card/plan, writes one immutable evidence document,
and appends canonical management updates and a verification event. Existing
lifecycle metadata and authored roadmaps remain unchanged: this is publication of
a bounded result, not overall completion or owner acceptance. A preview exposes
exact before/after hashes and proposed bytes before apply. No arbitrary paths,
patches, ledger events or owner decisions can be supplied by an agent.

Reuse the install package's existing management-history validation on the complete
proposed record set. Its supported shared-management profile is required; minimal
old lab history is not silently upgraded. Preserve all historical bytes. The core
checks selected identities, test/review applicability, source/skill/authority pins
and exact generated event contents. This validator reuse is an internal dependency;
it does not change installation policy or claim every repository-specific rule.

Profile 0.12 retains publication intent in the existing immutable operational store.
Apply serializes project writes across PG stores, persists applying state before
file effects and compares every relevant input against its preimage or generated
postimage. Unknown bytes stop with reconciliation_required; matching postimages
are replayed without duplication. Interrupted work remains inspectable. Other
run mutations in that store refuse while publication is applying. Atomic replacement
is per file; the full set is not atomic. Nonparticipating editors can race; exact
rechecks detect observed changes, not enforce security against same-user tools.

Read-only coordination adds compact review notices with current applicability and
observed provenance. The existing serialized leader delivery deduplicates their
stable content and delivers changes in applicability; no new mutation capability,
reviewer authority or host-native hierarchy is introduced. Publication and remaining
completion/owner disposition are explicitly distinct in these notices.
