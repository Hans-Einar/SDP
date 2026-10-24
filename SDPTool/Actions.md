# Discoverable actions and JSON invocation

Implemented on the Actions development branch and integrated through Session0011.
Source availability does not imply that an installed gh-sdp selects this producer.
Query the actual selected engine; no new release is published by this delivery.

## Find operations

```sh
sdptool actions
sdptool /path/to/project actions --json
gh sdp /path/to/project actions --json
```

The `sdptool.actions/1` catalogue contains producer version/revision, a metadata
revision, invocation protocol versions and sorted action descriptors. Each has a
stable ID, label/group/icon key, mutation class, generated input schema, output
schema identity and contextual availability. No input registration file is needed.
Without an explicit project, supported actions remain visible as requires-input.
Invalid context returns diagnostics and unavailable entries, with exit code 1.

| ID | Behavior | Writes |
| --- | --- | --- |
| project.discover | Source-derived project navigation and capabilities, including optional lazy SDL expansion | none |
| project.model-tree | Expand one SDL model at its discovered revision | none |
| blueprint.generate | Preview or immutable retained blueprint from model refs | documents |
| blueprint.assess | Evaluate pinned evidence; preserve blocked results | none |
| blueprint.assignment.list | Project-history assignment projection | none |
| blueprint.assignment.apply | Revision-bound attributed transition | canonical history |
| governance.session.definition | Supported pinned Session workflow | none |
| governance.session.context | Validate binding and canonical coverage | none |
| governance.session.list | Inspect instances in the selected store | none |
| governance.session.read | Inspect one instance, stages and holds | none |
| governance.session.adopt | Persist explicit adoption/registration intent | operational state |
| governance.session.resume | Reconcile registration, prepared selection/assessment/start/progress/verification/review effects and coverage | selected canonical records and operational state |
| governance.session.select | Prepare checked step/ready-card selection | operational state |
| governance.session.assess | Prepare attributed design/Blueprint assessment and canonical records | operational state |
| governance.session.start | Prepare assessed Session assignment and in-progress records | operational state |
| governance.session.progress | Prepare candidate-bound report progress for verification | operational state |
| governance.session.verify | Prepare exact passing test intake for Session review | operational state |
| governance.session.review | Prepare exact observed passing review intake for disposition | operational state |
| governance.session.dispose | Prepare checked selected-step/card completion and claim release | operational state |

`project.discover` accepts `navigation: "eager"` (default) or `"lazy-sdl"`.
`project.model-tree` requires the discovered `model` and exact `revision`; stale
inputs refuse and require refreshed discovery. See the [producer contract](Contract.md#deferred-sdl-model-navigation--n15a).

Only implemented operations are advertised. Install/upgrade, program run and PG
controller dispatch are still available through their existing interfaces; they
are not implicitly executable catalogue entries.

## Invoke one operation

```sh
sdptool invoke --request request.json
gh sdp invoke --request - < request.json
```

`invoke` must be the first argument after sdptool/gh sdp. It always emits JSON;
`--json` is optional, `--json=false` is rejected. Request files must be regular files;
use `-` for stdin. One object, UTF-8, at most 1 MiB, depth at most 32, no trailing
JSON or duplicate/unknown/case-variant fields. There is no command-string execution.

```json
{
  "schema": "sdptool.invoke/1",
  "requestId": "xfmd-42",
  "action": "project.discover",
  "context": {"project": "/path/to/project"},
  "parameters": {}
}
```

Project may name the root or its SDP directory; empty context selects cwd. Relative
project paths resolve against cwd. Other path parameters resolve against the project
root, independent of where request.json is saved. `entry` is model-source-relative;
embedded assignment paths retain their SDP-area-relative domain convention. Catalogue
`x-pathBase` annotations declare these exceptions. Underlying path/revision checks
remain authoritative. Absolute local paths retain existing local CLI capabilities.

Generation parameters:

```json
{
  "modelArea": "SDP/Models",
  "from": "release:0.2.3",
  "to": "work:Calibration",
  "entry": "System.design",
  "task": "task.json",
  "catalogue": "SDP/Blueprints"
}
```

Supply exactly one of `catalogue` and `output`. Task/evidence remain authored files;
WORK previews stay preliminary. Assessment takes `bundle` and `evidence`. Assignment
list takes `{}`. Apply takes `{"request": ...}` containing the existing strict
[sdp-blueprint-assignment/1 request](../SDP/04--Design/SDPTool/Blueprints/Assignment-Lifecycle.md).
Local assignment attribution is separately supplied:

```sh
sdptool invoke --request transition.json --as controller --authority controller
```

Session adopt/resume/select also require controller attribution. Those flags are rejected
for other operations. They attribute a trusted local
process; they do not authenticate remote callers. JSON cannot grant caller authority.
A hosted adapter must establish its own principal before calling the Go API.

## Results and consumer behavior

The `sdptool.invoke/1` response echoes safe requestId/action, then status, exitCode
and either result or error. Result is the existing domain JSON, not a reinterpreted
GUI structure. Optional requestId is 1–128 printable ASCII characters and only
correlates responses; it is not mutation deduplication.

| Exit | Status | Meaning |
| --- | --- | --- |
| 0 | ok | Operation succeeded |
| 1 | error | Domain/context/authority/I/O failure |
| 2 | error | Invalid envelope, schema, action or typed parameters |
| 3 | blocked | Valid assessment with unmet readiness; result retained |
| 130 | canceled | Cooperative cancellation before a committed outcome |

Decode stdout even on nonzero exit; preserve diagnostics and the exit code. Malformed
input has null correlation fields when they cannot be trusted. Error codes distinguish
invalid-request, invalid-parameters, unsupported-schema, unknown-action,
invalid-context, authority-denied, stale-revision and execution-failed. Where an
existing domain error has no typed classification, execution-failed preserves its
message instead of guessing from text.

Normal handled paths emit one envelope. A broken stdout or abrupt process death may
prevent a complete response: treat mutation outcome as unknown and reconcile via
reads. Never retry arbitrary mutations automatically. Assignment retry uses its
existing event ID, expectedEvent and pinned inputs. Retained generation is immutable
and repeatable. Cancellation is cooperative; noncancelable domain commits are not
reported as undone. The executable interrupts incomplete stdin on SIGINT/SIGTERM;
library callers own interruption of their supplied Reader.

Native clients should retain personal toolbar IDs across catalogue refreshes. Hide
unsupported actions without deleting preferences; disable supported unavailable
ones with their reason. Unknown icon uses a generic icon. Unknown schema vocabulary
requires a supported form/structured editor, not omitted mandatory fields. Input
schemas use the emitted object/string/required/additionalProperties/oneOf vocabulary;
complex assignment requests can use a structured JSON editor. This does not implement
XFMD widgets, Preferences or menu code.

## Go boundary and development verification

`Actions(ctx, selected)` returns the catalogue; `Invoke(ctx, reader, ActionCaller)`
returns a typed response. `RunWithInput` is the argv/stdin/output adapter; existing
`Run` callers remain compatible with empty implicit stdin. The registry binds typed
handlers and derives parameter schemas from those same types. Existing human commands
share the included operation handlers. Human catalogue output uses the presentation
registry, with terminal-control escaping.

[Integration handoff](XFMD-Blueprint-Integration.md) and
[verification](../SDP/05--Implementation/SDPTool/Actions/Evidence.md) describe the
actual compiled producer, headless consumer and isolated gh-sdp signed-distribution
trial. No normal extension installation, project upgrade or native GUI acceptance
is inferred from those checks.


## Session invocation parameters

The Session operations use `sdp-session-execution/0.1` results and the same service
as the governance CLI. Definition takes `{}`. List accepts optional `state`; read
requires `instance` and accepts optional `state`. Context/adopt/resume/select require
`requestFile` and accept optional `state`; both paths are project-relative unless
absolute. The referenced request uses the strict [Session schema](governance/session.schema.json).

Controller attribution is a separate trusted-host input, not a request parameter.
Adoption returns pending registration. Selection returns prepared record effects;
resume reconciles the journal or commits those effects, including ready-card moves.
These operations do not launch agents, advance beyond admission or confer acceptance. Restricted coordinator MCP capabilities remain unchanged.
See the [core guide](governance/README.md#session-instances--pgf1-c1) for the full
sequence, input/receipt rules and conservative drift behavior.
