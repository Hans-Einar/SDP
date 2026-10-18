# Discoverable actions and JSON invocation — proposed contract

Status: Session0011 T007 design proposal for the discoverable-actions KB-SDP-051.
No commands or schemas below are implemented yet. Product integration remains
behind the [concurrency gate](../../../Sessions/evidence/0011-concurrent-work/refresh-T007.md).
[ImplementationPlan](../../../05--Implementation/SDPTool/Actions/Plan.md) owns delivery.

## Outcome and boundary

XFMD can ask the selected SDPTool binary which supported operations exist, build
menus and forms from their metadata, and invoke a selected operation with JSON.
Human CLI and machine clients call the same typed operation handlers. Navigation
continues to describe project content; an action catalogue describes operations.
Neither requires a manually maintained navigation or registration file.

Use one compiled registry, not dynamic plugin loading, help-text scraping or JSON
translated into a shell command. The registry contains metadata, input validation
and the handler binding for each advertised action. Libraries own domain behavior;
transport owns framing, output and correlation. This is additive to existing CLI
commands. Do not advertise all CLI commands before their handlers are integrated.

MCP can later map this boundary; this increment is not an MCP server or daemon.
ProjectGovernance keeps controller/worker authority, dispatch and durable inbox
ownership. XFMD keeps GUI menus, personal toolbar configuration, icons and forms.
No host launches, installer mutations or arbitrary executable actions are exposed
by the first catalogue.

## Proposed commands

```sh
sdptool [PROJECT-OR-SDP] actions [--json]
sdptool invoke --request -
sdptool invoke --request request.json
```

`actions` defaults to readable output; `--json` selects its machine catalogue.
`invoke` is a machine entrypoint and always writes a JSON envelope; optional
`--json` is redundant and `--json=false` is rejected. No positional project or
operation flags may override the JSON request. Request-file location never changes
path resolution. Existing human commands and their output shapes remain compatible.
The executable passes stdin explicitly through a new input-aware entrypoint;
retain `Run(ctx, args, out, errs)` as a compatibility wrapper. Tests and library
callers must not depend on global stdin, cwd mutation or process environment changes.

For local assignment attribution only, invoke also accepts the existing-style
`--as ACTOR --authority controller|assignee|reviewer` outside the JSON envelope.
Other actions reject these flags. This has the same local trust boundary as today's
assignment CLI: a user already controlling the process can select attribution.
It is not authentication or an authorization mechanism for untrusted callers.
A future hosted adapter supplies a separately established principal in Go and must
not forward untrusted JSON into these flags.

## Request and path rules

Proposed schema: `sdptool.invoke/1`. One UTF-8 JSON object, followed only by
whitespace/EOF; no batches or streaming protocol. Example:

```json
{
  "schema": "sdptool.invoke/1",
  "requestId": "xfmd-42",
  "action": "project.discover",
  "context": {"project": "/path/to/project"},
  "parameters": {}
}
```

Require schema, action, context and parameters. requestId is optional (1–128
printable ASCII characters when supplied); it is correlation only. Reject unknown,
duplicate or wrongly cased fields at every typed object level, wrong types,
unexpected nulls, unsupported versions and trailing JSON before handler execution.
Limit the request to 1 MiB and nesting to 32 levels. Embedded assignment requests
obey their existing stricter domain schema as well. Limit failures are structured
errors and must not execute even a read handler. Values are data, never commands.

`context.project` accepts a project root or SDP area; omitted inside context means
invocation cwd. Resolve relative project paths once against the process cwd. All
path-valued parameters are then relative to the resolved project root (absolute
paths permitted by the existing local service policy). ModelArea explicitly selects
a ModelGovernance area; entry is a source-relative path inside each captured model,
not a path relative to cwd. Model selectors such as `work:Calibration` retain their
existing semantics. Nested binding paths inside assignment requests retain the
canonical SDP-area-relative lifecycle contract. Catalogue metadata must state the
base of each path field; saved requests must not depend on their file location.

Protocol path handling does not weaken the underlying symlink, containment,
protected-output or revision rules. Absolute local paths are not a remote file API.

## Catalogue and supported input vocabulary

Proposed schema: `sdptool.actions/1`. Include producer version/revision, invocation
protocol versions and a deterministically sorted `actions` array. Derive an action
catalogue revision from canonical metadata, excluding contextual availability;
clients invalidate their cache on producer or metadata revision changes.

Each action declares a stable id, label, description, group, iconKey,
mutation (`none`, `documents`, `history`), inputSchema, outputSchema and availability.
InputSchema uses a documented bounded JSON Schema vocabulary: object/properties,
required, additionalProperties=false, string, boolean, integer, enum, array/items,
and oneOf for exclusive parameter groups. Include bounds and explicit defaults;
no remote schema fetching, regex execution or executable default expressions.
Typed parameters and catalogue validators must derive from the same definitions,
with tests detecting drift from the domain request schemas.

UI annotations identify path bases, model selectors and selection bindings;
annotations are hints, never additional permissions. First-party reference forms
support scalar/path fields and structured JSON object editors for complex requests.
A client unable to support a schema disables that form with a reason rather than
silently omitting required fields. OutputSchema identifies the wrapped domain
result contract; results need not have one universal tree shape.

Availability is `available`, `unavailable` or `requires-input`, with stable reason
codes and text. Catalogue generation does not generate documents, create state,
run a renderer or change a ledger. Context-free catalogue output still lists
supported actions; absent selections yield requires-input. An invalid selected
project is an explicit diagnostic, not an empty catalogue pretending no support.
Availability is advisory: invocation always resolves current inputs and checks again.

Missing action on an older producer: hide it without deleting saved preferences.
Supported but unavailable action: disable it and show the reason. Unknown iconKey:
generic icon. Consumers choose display order; labels are not persistence keys.

## First operation set and service mapping

These IDs are proposed stable names; publish them only with executable handlers.

| Action | Inputs beyond context | Result/service | Mutation |
| --- | --- | --- | --- |
| project.discover | none | Existing Discover Project/navigation envelope | none |
| blueprint.generate | modelArea, from, to, entry, task; exactly one of output or catalogue | blueprints.Generate / GenerateRetained Result | documents |
| blueprint.assess | bundle, evidence | blueprints.Assess Assessment; blocked is a domain result | none |
| blueprint.assignment.list | none | AssignmentViews for resolved SDP area | none |
| blueprint.assignment.apply | request (embedded existing blueprint-assignment Request) | ApplyAssignment with separately supplied principal | history |

Task and evidence parameters identify existing files; do not add a second authored
task/evidence representation in this phase. Generation uses current capture,
preliminary WORK marking, immutable retention and source freshness rules. Listing
must preserve missing/stale entries and all diagnostics. Assignment requests retain
stable event IDs, expectedEvent, exact bundle/evidence pins and role checks.
A lost response must not prompt automatic replay of an arbitrary mutation.

Existing CLI parsers become adapters into these same operations. Do not call Run
recursively, serialize through temporary request files, or parse human stdout.
Noncatalogued commands continue through their current routes. Application services
can remain in the sdptool package to avoid import cycles; a small protocol package
may own transport types and decoding without importing its host. Choose file/package
boundaries from the integrated baseline, not by introducing another runtime container.

## Response, errors and retries

Proposed response shape (result varies by action):

```json
{
  "schema": "sdptool.invoke/1",
  "requestId": "xfmd-42",
  "action": "project.discover",
  "status": "ok",
  "exitCode": 0,
  "result": {"schema": "sdptool/0.1", "operation": "discover"}
}
```

The abbreviated result above illustrates wrapping, not a complete discovery fixture.
Exactly one envelope on stdout on normal success, blocked result or handled failure;
progress/log diagnostics only on stderr. Error envelopes have a stable error.code,
message and optional field/diagnostics; missing or unsafe correlation values are
null. Never echo the entire request. The following exit mapping belongs to invoke,
not a retrospective redefinition of existing CLI exit codes:

| Status / exit | Meaning |
| --- | --- |
| ok / 0 | Handler returned successfully |
| error / 2 | Invalid request, unsupported version/action, invalid typed parameters |
| error / 1 | Execution failure, stale revision, authority refusal or domain rejection; code distinguishes them |
| blocked / 3 | Valid assessment with unsatisfied readiness; preserve its complete result |
| canceled / 130 | Cooperative cancellation acknowledged before outcome commit |

Unknown action, unsupported schema, invalid-request, invalid-parameters,
stale-revision, authority-denied, execution-failed and canceled are stable codes.
More specific codes may be added compatibly. Do not classify errors by scraping text;
translate typed domain errors where available and preserve a generic domain failure
where no trustworthy category exists.

Encode a complete response before writing it. A broken stdout may prevent any
complete envelope: return nonzero, report on stderr, never emit a second envelope
or automatically rerun the mutation. Abrupt termination can leave the caller's
outcome unknown. Caller reconciles retained revision/event ID with existing read
operations; requestId alone never deduplicates work.

Pass context from process to service. The current blueprint CLI uses
context.Background for generation; replace it within the selected shared route.
Check cancellation before entering a noncancelable mutation, but do not convert a
successfully committed mutation into a canceled/no-change claim afterwards. Existing
service atomicity and recovery boundaries remain authoritative.

## Verification and integration gate

Exercise catalogue against actual handlers, strict parsing, CLI equivalence,
read-only no-write behavior, revision refusal, successful mutation and safe retry,
cancellation and output failure. Use a compiled binary with piped JSON and a
headless XFMD-like client that retains nonzero JSON stdout. Record wrapper and
engine versions separately for the actual gh-sdp test. No native GUI acceptance
is inferred from these tests.

Known blockers: shared output_arguments/CLI/presentation contracts, diverged
management histories and two distinct KB051 identities. PG is now actively working
on KB052. No code merge or ID rewrite is part of this design proposal. Integration
needs named ownership and ordering, preserved event provenance, a selected combined
baseline and a fresh preflight. Update authoritative SDL sources before implementation
on that baseline; the current shared System.design remains untouched meanwhile.
