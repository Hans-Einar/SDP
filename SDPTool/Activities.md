# Activity commands

Activity commands create and manage need-driven records in a project's SDP area.
They use the canonical `SDP/ProjectManagement/Ledger.ndjson`; there is no separate
Activity ledger. Creating an Activity leaves it **planned**. Creating a Session
writes a **planned discussion/work document**, without dispatching an agent or
creating a ProjectGovernance execution binding.

Use a disposable project when rehearsing these commands. The examples use
`jq` and lowercase UUIDs from `uuidgen`. Replace the project and authority paths
with your actual project and an existing document recording the owner's decision.
Do not create a real follow-on Activity merely to test the tool.

```sh
project=/absolute/path/to/project
authority=owner-decision.md
sdptool "$project" activity list --json > catalog.json
```

`authority` is relative to the SDP directory, so this example refers to
`$project/SDP/owner-decision.md`. Every mutation requires `--authority-source`.
The local CLI supplies its actor and owner invocation context. Connected action
hosts must supply a trusted owner principal separately; actor/authority fields in
request JSON cannot grant permission.

The canonical project namespace comes from `SDP/KanBan/board.json` and existing
Activity policy/history. Keep that namespace consistent. A browsing project slug
is not a replacement for a conflicting canonical identity.

## Create an Activity

Capture the current history revision and write a request file:

```sh
history=$(jq -r '.data.historyRevision' catalog.json)
jq -n --arg operation "$(uuidgen | tr 'A-Z' 'a-z')" --arg history "$history" '{
  operationId: $operation,
  expectedHistoryRevision: $history,
  title: "Selected customer need",
  owner: "project-owner",
  mandate: "Need: describe the problem.\n\nDesired outcome: describe acceptance.\n\nConstraints: describe the boundaries.",
  reason: "Owner selects this need for registration",
  sourceCards: [],
  systems: [],
  references: [],
  dependencies: [],
  scope: ["Source"]
}' > create.json
sdptool "$project" activity create --request create.json \
  --authority-source "$authority" --json > created.json
activity=$(jq -r '.data.activityId' created.json)
```

`scope` contains portable paths relative to the project root that an explicitly
adopted Activity may authorize for implementation. Select the actual allowed
prefixes; the example's `Source` is not a universal default. Source cards must
resolve to actual, unambiguous project cards. The tool allocates the identity and
folder, writes the mandate, and scaffolds six phases and one root Sessions home.

Keep the request file and operation UUID. Retrying the same request returns its
original result, including after later events. Changing content under that UUID
is a conflict. A new operation must use fresh expected revisions.

## Create a Session document

Read current Activity pins before writing a new request:

```sh
sdptool "$project" activity read --activity "$activity" --json > selected.json
sdptool "$project" activity list --json > catalog.json
jq -n --arg operation "$(uuidgen | tr 'A-Z' 'a-z')" \
  --arg history "$(jq -r '.data.historyRevision' catalog.json)" \
  --arg activity "$activity" \
  --arg event "$(jq -r '.data.lastEventId' selected.json)" \
  --arg descriptor "$(jq -r '.data.descriptorRef.sha256' selected.json)" '{
  operationId: $operation,
  expectedHistoryRevision: $history,
  activityId: $activity,
  expectedActivityEventId: $event,
  expectedDescriptorRevision: $descriptor,
  title: "Study the selected need",
  objective: "Investigate the need and identify the next decision",
  phase: "02-Study",
  reason: "Owner requests a study Session"
}' > session.json
sdptool "$project" activity session create --request session.json \
  --authority-source "$authority" --json
```

`session-create` is the equivalent flat command. Valid phases are `01-Mandate`,
`02-Study`, `03-Requirements`, `04-Design`, `05-Implementation` and
`06-Verification`. A plan is optional; when already selected, supply
`plan: {role:"plan", id:"PLAN-…", path:"…"}` naming a real record relative to SDP.
Plan references can also name actual Maintenance plans. For a Session spanning
phases, use `phaseContributions: ["02-Study", "04-Design"]`. The legacy `phase`
can be omitted when contributions are supplied; when both are present, it must
belong to that list. At least one phase is required, and duplicates refuse.

Optional `milestoneRefs` identify actual selected milestones using the parent
plan identity and milestone key, for example
`{role:"milestone", id:"PLAN-TEST-0001:A1", path:"Plans/Study.md"}`. Both the plan's
actual identity and its A1 milestone must exist. Optional `predecessor` names a real
Session with role `session` or `predecessor`; optional `handoff` names an actual
identified handoff/Session record. These fields use the same typed reference
shape and available-record checks. Add them only after those inputs exist; they
cannot invent future outputs. The generated document retains all contributions,
selected refs and relative source links under the exact creation snapshot hash.

Session numbers include
legacy Sessions and Activity-local homes. Completed/canceled Activities must be
reopened before adding work Sessions.

The generated Session starts with `Status: planned`. Subsequent actual execution
uses the existing governed Session context, selection, preparation and binding
workflow; a document's creation result is not an execution-ready PG instance.
See [ProjectGovernance](governance/README.md) for that separate workflow.

## Explicitly adopt execution admission

Reading an old manual catalog does not enable policy. Adoption explicitly names
all operational stores for the canonical project and requires known quiescence.
Use the same external PG state directory for binding, adoption and later execution.
Initialize each selected store first:

```sh
state=/absolute/private/path/outside-the-project/checkout-state
sdptool "$project" governance bind --state "$state" --json
sdptool "$project" activity list --json > catalog.json
jq -n --arg operation "$(uuidgen | tr 'A-Z' 'a-z')" \
  --arg history "$(jq -r '.data.historyRevision' catalog.json)" \
  --arg state "$state" '{
  operationId: $operation,
  expectedHistoryRevision: $history,
  stores: [{statePath: $state}]
}' > adopt.json
sdptool "$project" activity adopt --request adopt.json \
  --authority-source "$authority" --json
```

An uninitialized, unavailable, foreign or uncertain store refuses adoption.
Inventory every applicable store; an empty or guessed inventory does not prove
that no work is running. Adoption enables explicit Activity membership and scope
checks at governed implementation admission and dispatch. It does not select an
Activity, stop an existing process or authorize unlisted work.

Manual 0.1 descriptors remain readable. Before using them for scoped execution,
explicitly update their scope and selected memberships. `activity update` accepts
the current three revision pins, `activityId`, `operationId`, `reason` and the
complete intended `title`, `owner`, `scope`, `references`, `sourceCards`, `systems`
and `dependencies` fields. It preserves identity, mandate and exact prior
versions instead of overwriting historical descriptor evidence.

## Select and change lifecycle

For each new operation, refresh `catalog.json` and `selected.json` as above, then
construct the common transition request:

```sh
jq -n --arg operation "$(uuidgen | tr 'A-Z' 'a-z')" \
  --arg history "$(jq -r '.data.historyRevision' catalog.json)" \
  --arg activity "$activity" \
  --arg event "$(jq -r '.data.lastEventId' selected.json)" \
  --arg descriptor "$(jq -r '.data.descriptorRef.sha256' selected.json)" '{
  operationId: $operation,
  expectedHistoryRevision: $history,
  activityId: $activity,
  expectedActivityEventId: $event,
  expectedDescriptorRevision: $descriptor,
  reason: "Owner selects this Activity for current work"
}' > activate.json
sdptool "$project" activity activate --request activate.json \
  --authority-source "$authority" --json
```

Use the corresponding command with a newly pinned request:

| Command | Required starting state / extra input |
| --- | --- |
| `activate` | Planned; no other active Activity; dependencies explicitly satisfied |
| `hold` | Active; `resumeCondition` explains when the held need may resume |
| `resume` | On hold; active slot and dependencies available |
| `switch` | Active source; `targetActivityId`, `expectedTargetEventId`, `expectedTargetDescriptorRevision`, and former Activity `resumeCondition` |
| `complete` | Active; explicit need-level `outcome`, available typed `evidence`, and selected contribution dispositions |
| `cancel` | Planned/active/on hold; reason and dispositions for selected remaining work |
| `reopen` | Completed/canceled; explicit reason for the same need |

Active exit operations require the adopted store inventory and actual execution
inspection. If work is still in flight, the result says `pending: true`, preserves
the active state and blocks new admissions. Existing run observation, interruption
and recovery remain under PG. Retrying the exact lifecycle request, or recovering
its operation, checks quiescence again; a hold event never pretends to kill a run.
A switch holds its source and activates its target in one atomic event batch.

A completion request can add:

```json
{
  "outcome": "Owner accepts the evidenced need-level result",
  "evidence": [
    {"role":"verification", "id":"EVIDENCE-TEST-0001", "path":"Evidence/Result.md"}
  ],
  "contributionDispositions": [
    {"referenceId":"KB-TEST-001", "disposition":"completed", "outcome":"The selected contribution is accepted"}
  ]
}
```

These are additions to the pinned transition request, not a complete request.
Evidence must identify actual records. A completed plan or successful agent alone
does not complete the need. A no-code study can complete with its accepted result
and actual evidence.

References may select `membership: "contribution"` or `"whole-scope"`. Every
selected membership needs an explicit completion/cancellation disposition.
Contribution completion leaves a shared card's independent lifecycle alone;
whole-scope completion also checks the referenced record's actual completed state.
Cancellation can record `returned-to-backlog` or `transferred` with an outcome.
Unmet completion scope requires an explicit scope update first.

Dependencies are stable `{activityId, satisfiedByEventId?}` references.
Missing/self/cyclic dependencies refuse. Activation requires the dependency to be
completed and its explicit satisfaction event to name the latest actual completion;
renaming a title cannot change that identity. Foreign unavailable references must
name `projectId`, `availability: "unavailable"` and `unavailableReason`. They remain
visible diagnostics and cannot stand in for available completion evidence or local
Session membership. An optional reference `revision` pins SHA256 input bytes.

## Recover and browse

```sh
operation=$(jq -r '.operationId' activate.json)
sdptool "$project" activity recover --operation "$operation" \
  --authority-source "$authority" --json
sdptool "$project" activity history --activity "$activity" --json
sdptool "$project" activity timeline --activity "$activity" --json
sdptool "$project" activity validate --json
sdptool "$project" tree --json
```

Recovery uses retained exact file images and event bytes; external conflicts are
diagnostics. Preserve recovery metadata and repair/reconcile the actual conflict;
do not delete reservations to manufacture success. Completed retries return the
original result even after later history and legitimate Session-log edits.

Projects with Activity catalogs expose an Activities tab containing fixed state
and phase order, one Sessions home per Activity, and unassigned legacy Sessions.
They do not add a second standalone Sessions tab. Tree expansion targets contain
an `activity-tree://` URI and source revision. Pass both unchanged:

```sh
sdptool "$project" activity tree --uri "$expansion_uri" \
  --revision "$source_revision" --json
sdptool "$project" activity view --uri "$view_uri" \
  --revision "$source_revision" --output /private/tmp/activity-view --json
```

Use the actual producer-advertised URIs and revisions; follow continuation nodes
for every page. A changed source produces a stale diagnostic requiring refresh.
View output must be outside SDP and must not replace an ancestor of SDP. The
returned bundle is caller-owned and can be deleted after the consumer releases it.
Views link to actual source files. The timeline selection supports marker filters
and day/week/month zoom through its producer links; calendar duration is not effort.

All `activity.*` action definitions come from the same operation table as the CLI.
Actions accept their advertised parameter schema, such as `requestFile` plus
`authoritySource`; request files contain the same typed domain request shown here.

After admission adoption, hold an active Activity and wait for known implementation
quiescence before updating its scope, references or dependencies. Metadata-only
updates and equivalent array reorderings remain possible while active. The initial
manual scope/membership mapping is allowed before adoption. Updating a descriptor
cannot revoke permissions already issued to a host.
