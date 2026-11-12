# Activity record domain

`activities.New(area, projectID)` takes the absolute SDP document directory and
its canonical uppercase project namespace (for example `SDP`). The package imports
canonical `projecthistory`, YAML and standard libraries; host adapters and execution
authority remain outside it. `Caller` is trusted host context supplied separately
from request JSON. The library does not establish a security boundary against an
arbitrary in-process caller or a hostile process with the same filesystem access.

`Read` and `Replay` validate canonical history and its Activity-specific causal
chain, exact descriptor pins, identity/project matches and single-active rule.
Invalid replay returns no partial catalog. Manual descriptor/payload 0.1 records
remain readable without rewriting their bytes. State comes from events. Legacy
Session references retain their declared paths; no relocation is inferred.

`Create` creates a planned Activity and mandate without dispatching work.
`CreateSession` creates an Activity-local working document, before any plan if
needed, and requires both current Activity event and descriptor revisions. Session
numbers include existing legacy and Activity homes, declared Session references,
and actual primary Session identities in other Markdown homes. Phase contributions
can span the six Activity phases; the legacy single phase remains supported.
Optional milestone references pin a real plan identity and child milestone
(`PLAN-…:A1`); predecessor and handoff inputs name actual available records. Their
metadata and source links are retained in the hashed Session creation snapshot.
`Recover(operationID, caller)` completes the same bounded operation after failure;
changed requests with an existing operation ID fail. `Decode` rejects duplicate,
unknown/noncanonical JSON properties and trailing values.

New descriptors use `sdp-activity/0.2`: the 0.1 fields, optional `scope` containing
portable project-relative path prefixes, optional `dependencies`, and possibly empty `sourceCards` for a
direct owner need. New payloads use `0.2`: all 0.1 fields plus `requestDigest` and,
for `x-activity:session-created`, `session: {id,path,sha256}`. Session events retain
Activity state while advancing its predecessor. The envelope remains canonical
history 1.0. New creation event IDs derive from operation UUIDs.

Immutable descriptor bytes live in `Activities/.descriptors/<sha256>.yaml`.
Session creation bytes live in `Activities/.sessions/<sha256>.md`; their working
logs may subsequently evolve. Historical descriptor pins resolve only to matching
live bytes or verified immutable snapshots.

Bounded `Activities/.operations/<uuid>.json` intents contain the request/caller
digest, allocated identities, fixed event bytes/time, file before/after images and
fixed result. They are recovery metadata, not an authoritative history. Under the
shared canonical lock, an operation retains its intent, reserves the canonical
writer, checks file images, publishes files, validates replay, appends an atomic
batch, and releases the reservation. Recovery rejects external file conflicts.
Exact committed retries verify their original canonical batch and return the fixed
result even after later history or legitimate working-document edits. File and
history publication use synced atomic replacements. Filesystem checks reject
symlinks; like the shared history library, they do not defend against a hostile
same-user process swapping ancestors between checks and use.

`Transition(request, caller, inspector)` implements activate, hold, resume, switch,
complete, cancel and reopen. Active exits require explicit admission adoption with
an execution-store inventory and an actual inspector. Unknown disposition refuses;
in-flight work retains a pending request and blocks new admissions without holding
a filesystem lock or canonical reservation. `RecoverLifecycle` rechecks actual
quiescence; plain `Recover` reports pending honestly. Exact request retries preserve
identity. A switch publishes one complete two-event batch; partial/forged batches
are invalid replay. Hold/switch require an explicit former-Activity resume condition.

`Update` changes descriptor metadata/scope while retaining both descriptor revisions
and a state-preserving event. Dependencies are stable `{activityId,
satisfiedByEventId?}` references with checked missing identities, self references
and cycles. Selection requires explicit satisfaction by the dependency's latest
canonical completion event and current completed state. Renaming does not change
the reference. An active dependent must be held before reopening its dependency.

Reference memberships may select `contribution` or `whole-scope`. Complete/cancel
records explicit dispositions for all selected memberships; linking does not change
the source record. Whole-scope completion additionally checks the referenced
record's actual completed disposition. Completion requires trusted owner authority,
explicit need-level outcome and available evidence; agent/plan status alone never
completes an Activity. Cancellation records completed/not-applicable contributions
or the return/transfer of remaining selected work. Unmet completion contributions
require an explicit scope update first.

Available local typed references validate their actual Markdown/YAML/JSON record
identity and applicable type/project; optional SHA256 revisions detect changed
reviewed input. Source cards resolve to exactly one actual local project card.
Foreign unavailable references must explicitly name their project and reason;
they remain navigation diagnostics and never become local Session membership or
completion evidence. Legacy 0.1 reference fields are unchanged.

`ReplayRecords(projectID, snapshot, records)` uses the same reducer on captured
SDP-area-relative source files without consulting the filesystem. Live and captured
replay reject missing, duplicate, unregistered or conflicting descriptor homes.
New creates scaffold all six Activity phases and the one root Sessions home;
this does not claim phase readiness.

Targeted verification: `go test ./activities`. Tests mutate disposable projects;
the real manual Activity0001/0002 compatibility test only reads repository records.
