# SDPTool application boundaries — Activity design review

Session0016 T4, 2026-11-12; primary tracking card KB-SDP-061. Owner asks which Go
architecture improvements are worthwhile, especially CLI/action coverage. This is
source inspection and recommendations, not independent review or authorization
for a broad refactor. Related: [command design](Commands.md).

## Foundations to preserve

Keep domain packages, separate SDL/SDUI language ownership, shared project-history
primitives and capability-limited host adapters. The registry already derives
input schemas from typed inputs and advertises implemented handlers. Blueprint
handlers directly call domain services. Missing actions do not imply a broken CLI.

## Findings

| Observation and source | Recommendation |
| --- | --- |
| `action_sessions.go:sessionActionCall` builds argv, calls `governanceCommandAs`, buffers output and decodes JSON into `any`. The domain service is shared indirectly through a CLI adapter. | Both adapters should call a typed application service directly. Preserve strict decoding, paths, caller semantics and wire compatibility. Use this direction for Activity immediately; migrate existing Session adapters as bounded follow-up. |
| `action_invoke.go` maintains a caller-attribution action-ID allowlist; `action_sessions.go` separately classifies writes and roles; `action_registry.go` special-cases platform availability. | Declare attribution requirements and availability hooks beside each registered operation. Domain authorization remains authoritative; mutation metadata never grants authority. |
| Session bridge checks cancellation on entry, but `governanceCommandAs` accepts no context. | Carry context through applicable application operations, with explicit durable commit boundaries. Cancellation cannot imply a committed write was undone; retain reconciliation. |
| `actionRegistry` covers discovery/model-tree, Blueprint and Session actions. CLI model create/commit/status and select/preview are examples outside it. | Maintain explicit coverage and exclusions by consumer. Prioritize Activity, model and selection operations needed by XFMD/agents. Do not expose privileged controller or arbitrary program operations through every adapter. |
| Session schemas describe `requestFile`; catalogue availability mostly checks project validity and required input; `outputSchema` is an identity. | Give Activity inspectable input/output contracts and readable prerequisites. File input remains a CLI convenience. Recheck current state and authority at invocation. Schema identity alone supplies neither schema content nor a human overview. |
| `projecthistory/` supplies append and transaction reservations, not Activity rules or full multi-file recovery. PG publication is run/review-bound. | Activity owns its lifecycle and record-operation recovery. Reuse persistence primitives; extract shared record machinery only after checking actual reuse. No second ledger or fabricated PG run. |

The root package combines routing, discovery, navigation and view coordination.
Extract responsibilities incrementally when changing those workflows; file count
alone does not justify package migration or separate services. Navigation cleanup
remains in queued KB-SDP-060. MCP retains capability-specific exposure even if the
broader local action catalogue expands.

## Recommended scope

T5 supersession: the owner rejects unscheduled migration follow-up. The proposed
[complete migration plan](../../../../05--Implementation/SDPTool/ArchitectureMigration/Plan.md)
replaces the delivery strategy below with one finite inventory and mandatory
removal inside the same plan. T6 owner scope separates restructuring (Activity0001)
from Activity feature implementation (future second Activity). Findings remain
applicable; no new product functionality is permitted during restructuring.

For Activity, establish typed service calls, CLI/action parity, operation policy,
explicit record recovery and producer-owned projections. Complete payload/guard
design before implementation. Broader Session adapter/catalogue migration remains
a proposed follow-up tracked here under KB-SDP-061 until separately scoped; it is
not all a prerequisite for Activity delivery.

Future verification should compare adapter results and failures, caller rejection,
stale revisions, retries, interrupted publication/recovery and wire compatibility.
No tests were run for this source-only assessment; these are future obligations,
not demonstrated runtime failures. No product code or SDL changed.
