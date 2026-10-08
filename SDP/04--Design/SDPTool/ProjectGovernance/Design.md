# ProjectGovernance — first workflow contract

Design candidate PGD1, owned by [PLAN-SDP-0017](Plan.md). This specifies the bounded
implementation handoff; no new runtime or installed process schema exists yet.
Owner-selected scope and role responsibilities are retained. Storage locations,
operation names and implementation choices below are the recommended pilot design.

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
2. Submit a request to the project-facing Steering/PM conversation. Record the
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
codex-cli 0.160.0, default schema surface; [Evidence.md](Evidence.md) records its
limits. Earlier 0.158.0 evidence remains historical. Unsupported versions fail with
a compatibility diagnostic until tested. Do not require undocumented experimental
fields or interpret the documentation's newest example as the installed contract.

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
