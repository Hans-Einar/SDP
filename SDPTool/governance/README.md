# ProjectGovernance core — PGI1

This package implements the first slice of
[PLAN-SDP-0018](../../SDP/05--Implementation/SDPTool/ProjectGovernance/Plan.md):
read-only work inspection, one pinned Maintenance routine and durable local run
transitions. Linux supports mutation; other platforms explicitly refuse writes.
The CLI uses the same core that future adapters will call. No agent execution,
MCP server, app-server connection, automatic Session journal or delivery gate is
implemented in this slice.

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
evidence submission and further routine steps belong to later PGI milestones.
The serialized definition and digest stay pinned on each run.

## Authority, consistency and recovery

The CLI is an explicit trusted local controller adapter. Core mutation needs an
adapter-supplied Caller capability; requests cannot provide role/owner approval.
A referenced authority document is pinned input, not parsed proof of approval.
The host adapter remains responsible for deciding whether an invocation is
actually authorized. This boundary does not defend against arbitrary same-user
shell access. Future agent adapters must not grant LocalController from model text.

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
are a pilot boundary, not a pagination API. No cursors/events API exists yet.

[Schema definitions](schema.json) describe each CLI request and response. Select
the corresponding `$defs` entry; semantic coverage, byte limits and duplicate-key
rejection additionally run in Go. Errors go to stderr with exit 2; success uses
stdout/exit 0. Human output is default and `--json` selects machine output.

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
