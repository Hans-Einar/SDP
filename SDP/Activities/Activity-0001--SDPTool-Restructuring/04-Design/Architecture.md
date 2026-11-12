# M1 architecture and implementation contract

Owner: ACTIVITY-SDP-0001; PLAN-SDP-0026. No new public capability. The whole-system
boundary inventory is ../02-Study/Inventory.md and Operations.json. The new
ArchitectureMigration SDL concern owns this cross-entrypoint responsibility map;
existing detailed domain models remain authoritative at their existing homes.

## Selected change and rationale

Introduce governanceapp as a typed application operation layer above governance.
It imports governance, not the root facade, CLI, presentation, MCP or controller.
An operation definition owns its name, input kind, local action metadata/exposure,
write/attribution policy and typed domain handler. Strict domain decoding remains
used. CLI parses flags; actions decode their existing parameter schema; both resolve
project/state and call the same operation. No internal argv or stdout/JSON bridge.
Registration metadata is transport description; governance still checks authority.

OperationRegistry keeps the 28 existing public actions with identical descriptions,
labels, schemas, mutation hints and catalogue revision. Internal action policies
carry caller attribution and platform availability beside registration; Invoke
must not contain a second action-name allowlist. Public structs/wire fields stay
unchanged. Blueprint handlers remain direct calls to blueprints services.

MCP retains its explicit capability-based operations and typed domain calls.
Its ordinary and Session response wrappers share one implementation parameterized
by existing schema/error-limit policy. Do not replace launch-bound guards with local
CLI owner/controller policy, expose new tools, or broaden owner-decision authority.

## Preserved process boundaries

Controller self-exec is retained: it creates independently leased child processes,
private logs and bounded concurrency; no CLI-result parsing is involved. gh-sdp
forwarding, installer binary version probes, XFMD subprocesses, SDUI program launch
and renderer execution are real process protocols. They are not prohibited by the
in-process adapter rule. No container is split merely to match a package folder.

## Bounded worker assignment: M2

Allowed code: SDPTool/governanceapp/, governance_cli.go, action_sessions.go,
new governance_service.go helper, action_registry.go, action_invoke.go,
action_blueprints.go and focused tests. These share dispatch registration; do not
change governance domain rules or mutate operational stores as a side effect of
construction. Root facade resolves existing project/state paths exactly as before.

Remove sessionActionCall's argv construction, CLI call, output buffers and JSON
reparse. Replace governance CLI's duplicated name/input/dispatch lists using the
operation definitions. Replace Invoke's attribution allowlist with registered
policy. Preserve same CLI argument errors, strict request validation, trusted
caller interpretation and context cancellation semantics. Durable mutation must
never be reported undone because cancellation arrived after its commit.

The public Invoke result is heterogeneous any, already typed for Blueprint actions.
Compare JSON wire meaning for Session parity rather than requiring an incidental
map representation produced by the deleted round trip. Do not introduce another
serialization round trip to satisfy that implementation-specific assertion.

## Removal and evidence matrix

| Removal/retention | Required check |
| --- | --- |
| Session action to CLI/output bridge removed | AST dependency check plus CLI/action result/error parity |
| Invoke action-ID attribution allowlist removed | Registry policy checks and denied/allowed caller behavior |
| Separate governance CLI dispatch/name classification removed | Operation inventory/CLI parity; domain handlers referenced once in governanceapp |
| Duplicate MCP retry/error framing removed | Shared wrapper tests for both schemas, busy/stale/cancel/output caps; existing capability tests |
| Existing operations and external process boundaries retained | Frozen catalogue and inventory comparison; compiled CLI/controller/MCP and native producer tests |

Architecture controls must detect a deliberately prohibited adapter dependency or
CLI bridge. Generic static rules cover packages/adapter imports; integration checks
cover behavior. Checks cannot claim full Go/SDL semantic equivalence.

## Blueprint and handoff

The retained whole-system Blueprint compares validated ObservedArchitecture and
TargetArchitecture, both preliminary. See Blueprints.md for exact references.
It expresses structural ownership changes, not detailed runtime proof. This
contract and the inventory supply missing behavioral/code context. New agents must
read those artifacts, explain the ownership/removals and run required checks.
Do not treat current transitional code as authority over this target.
