# SDPTool MCP connection study

| Field | Value |
| --- | --- |
| Milestone | RGS2-B-M1 |
| Date | 2026-09-28 |
| Authority | [StudyPlan](StudyPlan.md), PLAN-SDP-0007; KB-SDP-036–038 |
| Status | Research recommendation; no adapter or routine engine implemented |
| Local candidate | `02f4e6d99e19f6945973c1a260135053b329a4c3`, branch `sdp/request-routine` |
| Evidence | [Source pins and inspection record](Evidence/RGS2-MCP-source-pins.json) |

**Recommendation:** add a thin MCP adapter to the existing Go SDPTool, with one
application service shared by CLI, MCP and observer/client queries. First prove
one bounded Maintenance routine. MCP should expose its contract, not contain a
second workflow engine. This refines [RGS1](Study.md) and the
[routine catalog](Routine-Catalog.md); it does not adopt new process authority.

## 1. What exists

Source inspection found the following boundaries. These are implementation
observations, not new runtime test results.

| Existing source | Capability and limit |
| --- | --- |
| `SDPTool/project.go`, `navigation.go`, `kanban.go` | `Discover`, `Navigation`, `BoardNodes`: explicit project bindings, declared capabilities, bounded navigation and board/history reading. Navigation includes an inventory digest and unavailable-service diagnostics. No generic work-transition API. |
| `preview.go`, `selection.go`, `sdui.go`, `viewer.go` | Revision-checked previews, saved bundles and viewer integration. Rendering writes files and can launch configured programs; it is not a read-only knowledge resource. |
| `install/plan.go`, `execute.go`, `records.go`, `history.go` | Installation plans, exact-input drift checks, journals, Apply/Resume, receipts and installation-specific management history. Installation preview can populate the artifact cache. These are not general routine transactions. |
| `install/lock_unix.go` | Nonblocking installer lock keyed by project-root hash in the user's cache. This is not yet a shared routine-writer lock, cross-user lock or distributed transaction. |
| `cli.go`, `go.mod` | CLI dispatch calls Go functions; Go 1.26.0 module with local SDL/SDUI/bootstrap modules. No MCP dependency or MCP command, routine service or blueprint implementation was found in the inspected Go tree. |

[The producer contract](../../../SDPTool/Contract.md) already separates navigation
eligibility from installed conformance. The current registration omits a project
manifest reference: distribution metadata at repository root must not be reported
as this checkout's installed receipt. Existing build-version defaults likewise
do not prove the version of any installed executable. Product implementation,
installed process profile, routine version and MCP protocol are separate facts.

## 2. Protocol and SDK baseline

RGS1's dated `2025-11-25` sources remain historical evidence. For this study the
official specification repository is pinned at
`ab3a39c13bd23be691c2760e1c6c5c15a64582e1` (2026-09-24), containing revision
`2026-07-28`. That revision uses per-request metadata and version discovery;
legacy `initialize` semantics apply through `2025-11-25`. Compatibility must be
selected explicitly, not inferred from a successful TCP connection.
[Versioning source](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/basic/versioning.mdx)

Official Go SDK `v1.8.0`, released 2026-09-14, resolves to
`3f3b699b2b67e1ed033a63d6651671dab53c2d32`; its Go requirement is 1.25.0.
Its documented APIs support both protocol eras, stdio and Streamable HTTP;
supported versions can be restricted. Modern HTTP requires its stateless option.
This is a plausible dependency candidate, not a dependency selection or build
compatibility result. The actual intersection with the sibling study's installed
Codex 0.158.0 candidate remains unverified; neither a current SDK nor an
app-server schema proves that host's MCP negotiation.
[SDK release](https://github.com/modelcontextprotocol/go-sdk/releases/tag/v1.8.0),
[pinned protocol guide](https://github.com/modelcontextprotocol/go-sdk/blob/3f3b699b2b67e1ed033a63d6651671dab53c2d32/docs/protocol.md),
[module](https://github.com/modelcontextprotocol/go-sdk/blob/3f3b699b2b67e1ed033a63d6651671dab53c2d32/go.mod)

## 3. Candidate interface

All names below are proposals. Resources provide retrievable context; tools
perform requested operations, including parameterized read-only queries. No
resource read should start a run, generate a preview or change a card.
[Resources specification](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/server/resources.mdx)

| Candidate surface | Input → result; implementation dependency |
| --- | --- |
| `sdp://project/{handle}/context` resource | Resolved root/worktree, current authority references, selected work, installed facts and freshness. Reuse discovery; add authority/work resolution. |
| `sdp://project/{handle}/routines/{id}/{version}` resource | Immutable definition, hash, compatibility and authority. Catalog storage is new. |
| `sdp://project/{handle}/runs/{runId}` resource | Revisioned snapshot, children, waiting reasons and evidence references. Routine state is new. |
| `sdp://project/{handle}/assignments/{id}` resource | Bounded assignment, candidate and BP2 bundle reference/coverage. Does not activate BP2 or invent bundle contents. |
| `sdp_context` read tool | Explicit selected project → same context resource envelope; useful where host resource discovery is weak. |
| `sdp_route` read tool | Bounded intent summary, existing work reference, catalog revision → match/no-work/missing/ambiguous/unavailable/incompatible/conflict plus reasons. Suggestions grant no authority. |
| `sdp_run_start` mutable tool | Selected work, pinned routine and assignment → durable run ID. Requires authorized work; no automatic card/plan proliferation. |
| `sdp_evidence_attach` mutable tool | Run, candidate-bound artifact digest, evidence level and provenance → recorded claim/reference, not automatic verification. |
| `sdp_transition` mutable tool | Run, action, expected revision, idempotency key and evidence references → accepted transition or explicit refusal. |
| `sdp_events` read tool | Project/run scope, opaque cursor, bounded limit → ordered events, next cursor and coverage/gap status. |

Use strict object input/output JSON Schemas, bounded strings/arrays and rejected
unknown fields. Common results should contain `schemaVersion`, `projectHandle`,
`worktreeId`, `revision`, `observedAt`, `sourceRefs`, `coverage`, and operation data.
Provenance names the source and digest; `coverage` distinguishes absent, declared,
validated, stale and unavailable. Return structured content and concise text.
Tool annotations describe behavior but do not implement security.
[Tools specification](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/server/tools.mdx)

For example, transition arguments are
`{projectHandle, worktreeId, runId, action, expectedRevision, idempotencyKey,
evidenceRefs}`. Action is a routine-defined enum, not arbitrary JSON patch or
shell text. Success returns `{operationId, revision, state, eventCursor}`;
refusal returns `{code, message, retryable, currentRevision, unmetConditions}`.
Exact schemas and compatibility policy belong to the next DesignPlan.

## 4. Authority, identity and failure

The core resolves handles to explicitly selected physical roots, preserving
existing path/symlink containment. A project ID is insufficient to distinguish
clones and worktrees. Bind a checkout identity plus Git candidate and relevant
dirty-input hashes; retain task/run IDs across replacement agent sessions.
Thread, turn, attempt and transport identities are separate correlation fields.
External project references remain unverified unless required evidence is supplied.

Validate caller identity and assignment permissions from a trusted launch binding
or authenticated request context. A supplied `role: owner`, client name, run ID
or thread ID grants no permission. Steering/PM carries project scope and recorded
owner decisions; Master coordinates the assigned task; Worker proposes completion;
Reviewer records review within assigned authority. Owner acceptance, merging and
publication require their actual authorization. Host confirmation and OAuth
access do not establish those project decisions.

For HTTP, transport authorization follows the MCP OAuth resource-server model;
stdio uses the local launch/environment trust boundary. Both still need
application checks per action. Start with read-only and assignment-scoped mutation
permissions; do not expose unrestricted file reads, subprocess options or generic
ledger append. Private evidence should be referenced minimally, not copied into
tool descriptions or telemetry.
[Authorization specification](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/basic/authorization/index.mdx)

The proposed core serializes conflicting writes and compares the complete
relevant input revision immediately before publication. Same idempotency key and
same request digest return the recorded outcome; key reuse with different
arguments fails. JSON-RPC request IDs are not durable deduplication keys.
After a timeout, inspect the operation or retry its key before starting another.
Cancellation requests work to stop; it cannot undo an already committed effect.

Malformed protocol requests use JSON-RPC errors. Domain refusals use a tool error
result with stable application codes such as `unauthorized`, `stale_revision`,
`evidence_missing`, `incompatible_profile`, `idempotency_conflict`, `busy`, and
`reconciliation_required`. Distinguish retryable contention from changed authority
that needs reassessment. Preserve diagnostics without leaking private content.
[Tool error distinction](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/server/tools.mdx#error-handling)

## 5. Transport, recovery and version ownership

| Choice | Benefit | Cost and boundary |
| --- | --- | --- |
| Stdio subprocess | Small local pilot; host controls process lifetime; no listening port | Separate clients may launch separate writers, so locking/storage must be shared. Keep stdout exclusively MCP; restart loses in-flight requests. |
| Streamable HTTP | Shared service for simultaneous client/observer access and possible remote use | Requires service ownership, access control and deployment policy; validate Origin and use localhost binding for a local service. Remote use adds network trust and private-data concerns. |

[Stdio contract](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/basic/transports/stdio.mdx),
[HTTP contract](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/basic/transports/streamable-http.mdx)

Recommend stdio for the first local acceptance candidate, conditional on the
chosen host's demonstrated protocol support. Defer network deployment. An
observer can query the same core through a read-only CLI initially; a later
application client may use a service facade. Neither should reconstruct
transition rules from events or scrape the agent conversation.

Modern MCP subscriptions require re-establishment and provide no durable
subscription state across reconnections. `2026-07-28` HTTP removes session IDs
and `Last-Event-ID` replay; legacy replay support must not become an SDP guarantee.
Use notifications as refresh hints. SDP snapshots plus application-owned ordered
event cursors provide recovery, duplicates/gaps and retention-expiry handling.
[Subscriptions](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/basic/patterns/subscriptions.mdx),
[HTTP revision](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/basic/transports/streamable-http.mdx)

Operational persistence should reference management events and Traceability
evidence without creating another management ledger. Multi-file updates need a
journal/reconciliation boundary: do not acknowledge a complete transition while
its card and history disagree. Installation recovery offers concrete patterns,
not an already reusable transaction engine. Reads during recovery must expose
incomplete/stale state. Signed installation receipts prove package provenance,
not owner acceptance or skill compliance. Pin existing runs to definition hashes;
incompatible upgrades require explicit migration or the original reader.

For modern cacheable responses, set private scope for project/user data and
appropriate freshness; public caching can otherwise cross caller boundaries.
Authorization and revisions still apply on every mutation.
[Caching contract](https://github.com/modelcontextprotocol/modelcontextprotocol/blob/ab3a39c13bd23be691c2760e1c6c5c15a64582e1/docs/specification/2026-07-28/server/utilities/caching.mdx)

## 6. Pilot and acceptance boundary

Skills should recover context, call the selected route and describe responsibilities.
MCP supplies structured facts and checked operations. Host integration binds real
calls/threads to work. An agent with shell/file access can bypass all three by
editing records directly or ignoring a tool refusal. An adapter cannot certify
unobserved activity or impose OS isolation; the observer must distinguish declared
use, observed invocation and verified outcome. SDK capability negotiation does
not prove instruction loading or compliance.

Propose one sandbox Maintenance assignment using RT01/RT08/RT10/RT11/RT12/RT14
and a bounded RT15 gap outcome. Build core/CLI and MCP together around the same
transition function, then a read-only observer. Exclude real publication,
installation and automatic procedure adoption. Before success is claimed, prove:

1. CLI and MCP return equivalent decisions for the same candidate and permissions.
2. Stale worktree/evidence, forged roles and missing authorization are rejected.
3. Two writers, retry after lost reply and interrupted card/history publication
   recover without duplicate accepted effects or falsely completed records.
4. Restarted observer reconstructs state; expired cursors expose a gap and resync.
5. Fresh Master/session resumes the same assignment; Worker escalation preserves scope.
6. Installed host/SDK negotiate a pinned common protocol; bypass attempts show
   actual coverage limits. Compare overhead, false blocking and owner rescue
   against ordinary matched tasks.

These are proposed acceptance cases, not executed tests. This study performed
source reads, targeted searches and official-document retrieval only; no MCP
server, host configuration/auth inspection, real model call or product test ran.
Remaining decisions are persistence/reconciliation schema, trusted caller binding,
host/protocol intersection, catalog compatibility and pilot acceptance thresholds.
The next DesignPlan must resolve these before treating the adapter as enforcement.
