# Activity commands and SDPTool source ownership

Originally selected as design work in Session0016 T3; Activity0002/PLAN-SDP-0027
now implements the contract. Current executable usage is maintained in
[SDPTool Activities](../../../../../SDPTool/Activities.md). The authored SDL WORK lives in
[ActivityWorkflow](../../../../SDL/ProjectGovernance/SDPTool/ActivityWorkflow/README.md).

## Source responsibilities (updated after Activity0001)

| Area under SDPTool | Responsibility |
| --- | --- |
| cmd/sdptool | Small executable entry; delegates argv/input/output to the root package |
| cli.go and *_cli.go | CLI routing and domain adapters; they are not separate service processes |
| action_registry.go, action_*.go | Typed action catalogue, input schemas, JSON invocation and calls into shared handlers; not every existing CLI command is catalogued |
| project.go, discovery.go, navigation.go, tree.go, selection.go, preview.go | Project/source discovery, final navigation, selected source-model view and document coordination |
| activities/ | Activity lifecycle, record recovery, reference/dependency validation, projection and admission policy |
| governanceapp/ | Shared typed Session application operations used by CLI/actions |
| presentation/ | Human-readable output/JSON routing for command results; not XFMD's graphics renderer |
| model/ | Mutable WORK, model checkpoints, proposals/candidates/releases, snapshots and recovery; distinct from project Git |
| blueprints/, blueprintstate/ | Design-change context/bundles and assignment lifecycle reduction/evidence |
| projecthistory/ | Canonical ledger persistence, expected-revision append and transaction coordination; no domain-specific Activity rules |
| governance/ | PG context, assignments/runs, Session execution binding, guards, evidence/review/publication and recovery |
| governancecontroller/, governancehost/, governancemcp/ | Coordination, host/app-server integration and capability-bound MCP adapters to that core |
| cmd/sdp-controller, cmd/sdp-mcp | Separate runnable controller and MCP entrypoints using shared packages |
| install/, bootstrap/, codexcompat/ | Toolkit installation/upgrades; thin client download/verification; host compatibility respectively |

The Go module imports SDL/go and SDUI/go through local replacements. Those own
language parsing/projections/runtime behavior. gh-sdp delegates to an engine;
XFMD consumes its output. Directory count does not imply microservices.

Existing governance.session.* actions operate on governed execution bindings or
instances. session-successor validates a successor context and its prepared effects;
it is not a generic Activity-root Session document creation API. Reuse the existing
execution/recovery path when a created Session is later adopted, not at creation.

## Selected facade and actions

Keep a single activity command family, including Activity-scoped Session records:

| CLI under `sdptool PROJECT` | Action | Effect |
| --- | --- | --- |
| activity create --request FILE | activity.create | Create planned Activity, standard layout, mandate from supplied need and one canonical creation event |
| activity session create --request FILE | activity.session.create | Create/link one Session beneath Activity-root Sessions; no execution binding or agent start |
| activity list / read / history / timeline / validate | activity.list/read/history/timeline/validate | Read/replay records or derive selected content |
| activity view --uri URI --revision HASH --output DIR | activity.view | On-demand Markdown/SVG bundle, with source-bound selection and caller-owned output |
| activity tree --uri URI --revision HASH | activity.tree | Bounded producer-owned Activity subtree, separate from view selection |
| activity update --request FILE | activity.update | Retain descriptor versions and append metadata/scope/reference changes |
| activity adopt --request FILE | activity.adopt | Explicit admission opt-in with initialized, known quiescent execution store inventory |
| activity recover --operation UUID | activity.recover | Recover exact retained effects or resume pending quiescence |
| activity activate / hold / resume / complete / cancel / reopen --request FILE | activity.<operation> | Checked lifecycle transition with current evidence/authority |
| activity switch --request FILE | activity.switch | Atomic hold/activate selection with in-flight work disposition |

CLI and JSON/actions call the same typed application service. A later authorized
MCP adapter uses that service; do not expose write operations in an inspection-only
worker MCP connection. Advertise actions only after implementation is executable.
Do not add parallel state-changing session-create paths in governance CLI.

## Creation and Session contracts

Activity create input: operation identity, expected project-history revision,
title, owner attribution, authored need/desired outcome/constraints, optional source
cards and dependencies. Actor/authority comes from trusted invocation context,
not a freely asserted owner role in JSON. Allocate a project Activity identity and
safe folder name under the canonical transaction. No reuse after deletion/cancel.
Scaffold six phases, standard README and root Sessions. Keep the need narrative
only in Mandate.md. Initial state is planned; create never implicitly activates.

Session create input: operation identity, expected history and Activity revision,
Activity reference, Session objective, phase contributions, selected plan/milestone
references where already available, and optional predecessor/handoff reference.
Require real linked records; future intended outputs are not fabricated identities.
Allocate a project-unique Session identity across legacy and Activity-local records.
For planned/on-hold Activities, creation is record-only preparation, never dispatch.
Completed/canceled Activities require reopening before adding new work Sessions;
explicit historical adoption remains a separate operation with honest timestamps.

Session creation publishes its document and Activity membership atomically using
the canonical history transaction. Reuse the Activity domain's linked event for
membership; do not invent an implemented Session lifecycle ledger. A Session may
be created for study without an ImplementationPlan. The actual PG execution gate
still checks its chosen routine/plan/card coverage before adoption/start.

Both create operations return disposition, allocated identity/path, resulting
history/descriptor revision and readable navigation target. An identical retry
returns the existing result; different content under the same operation identity
conflicts. Missing mandate, broken references, stale revisions, unsafe paths,
unsupported schema, ID collision and pending recovery are explicit failures.
Partial folder creation must not appear as a registered Activity/Session.

## Implemented ownership boundary

[Architecture assessment](ArchitectureReview.md) records existing adapter and
catalogue gaps and the recommended boundary for Activity implementation.

The `activities/` package owns typed Activity application operations, pure
lifecycle checks, descriptor/Session-document assembly, domain history replay and
projections. Use projecthistory for canonical persistence. root activity_cli.go
and action_activities.go remain thin adapters; presentation gains result formatters.

The existing governance publication machinery is partly bound to PG runs. Do not
claim it can simply be called as a generic creation API: reuse its underlying
projecthistory primitives and recovery invariants, sharing a general primitive only
where inspection demonstrates one. Activity services must not depend on host/MCP
adapters. governance consumes a narrow Activity scope/admission interface; any
in-flight run check is an injected port or orchestration concern, avoiding cycles.

Discovery asks the Activity service for the ready tab/tree. Projection code owns
classification, ordering and source attribution; XFMD does not calculate these.
Session creation and later execution use the same canonical identity and a current
revision-bound handoff. Existing independent model/Blueprint stores stay separate.

## Model and verification boundaries

The first SDL increment models create Activity, create Session, guarded state
change, implementation admission and projections, with explicit request/result
and record-effect boundaries. Accepted creation and refusal scenarios demonstrate
ordering and result types; prose defines the actual predicates. Runtime atomicity,
quiescence, permission and timeline attribution need implementation tests later.

The original design-only turn ended at a validated WORK checkpoint. Activity0002
now supplies implementation, retained Blueprint handoff and executable verification.
The Blueprint's structural result does not claim runtime conformance; authored
code mapping and candidate evidence supply that assessment. No accepted model
release, product release or live installation is inferred from implementation.

Public writes require `--authority-source` (an existing SDP-relative owner decision)
or an attributed trusted owner invocation. Command and action metadata share
`activities.Operations`; CLI parsing is not an internal service interface.

## Historical design-only preflight — Session0016 T3

SDP has one registered worktree on sdp/sdl-browsing-planning at a96affd. Dirty files
are this conversation's planning records; there are no existing SDL or Go edits.
SessionLifecycle WORK is clean at :00003. Its prior PG ownership and model history
remain intact. New write scope is ActivityWorkflow, this command/design handoff,
the System README index and current KB061/Session0016/evidence records. Existing
navigation/model/PG implementation is read-only. Scoped verdict: green for this
new bounded model concern; unknown remote work and runtime integration remain
outside this check. No independent review or implementation readiness inferred.
