# Activity records and admission

The installation profile distributes versioned Activity descriptor/payload schemas
and the read-only `SDP/Framework/activities/activity_history.py` utility. The source
candidate advertises `sdp.activities.records.v0.2`; publication and actual installed
capabilities must be checked against the exact signed descriptor/receipt. This
source change does not change previously released Toolkit facts or adopt a project.

The canonical management ledger also carries the separate Activity domain. Retain
manual descriptor/payload 0.1 records and their original identities/timestamps.
Runtime operations use 0.2 records; prior descriptor bytes remain under
`Activities/.descriptors/<sha256>.yaml`, and Session creation snapshots under
`Activities/.sessions/<sha256>.md`. Mutable current documents are separate from
those immutable input snapshots. Lifecycle state comes from causal replay, never
from the directory name, descriptor metadata, file times or agent completion.

Use the installed SDPTool's Activity CLI/actions only when that running binary
supports the required schema. Check its advertised actions and exact version.
Older writers that reject Activity events cannot safely publish through them;
upgrade the writer or use a read-only compatible file view. Do not strip unsupported
events, edit historical bytes or infer enforcement from these installed documents.

Activity admission is a separate explicit adoption. Bind/initialize the selected PG
stores through their existing operation first, then inventory those exact stores
in the adoption request. No store is created automatically by a read or adoption.
Adoption pins the canonical Activity namespace in `x-activity-admission:adopted`
payload 0.1; `.admission/stores/` only locates actual PG implementation state. Arbitrary old/manual
executables, disconnected writers and unknown legacy state locations are outside
cooperating runtime enforcement. Missing/uncertain known stores block implementation disposition.

A pending hold/switch blocks new implementation admissions and outbound claims.
Existing exact sending/uncertain host inputs remain in-flight until actual
observations/reconciliation establish their disposition. Observation, interruption
and recovery remain possible; an interrupt request is not observed quiescence.
Do not retain the canonical filesystem lock or record reservation while waiting for
an agent. Read-only browsing, Context/Route/Read and KB capture do not require implementation
admission. Candidate inspection retains its existing source-revision checks.
Reviewer and Project Leader hosts retain their existing read-only sandbox/tool
contracts and are outside implementation quiescence; no universal host stop is claimed.

Run `python3 SDP/Framework/activities/activity_history.py SDP` for structural/replay
validation. It reads local records and schema files; it neither probes external
hosts nor proves execution quiescence, owner acceptance or publication. The Go
installer validates captured Activity records through the same runtime reducer,
preserving project content and immutable ledger bytes during install/upgrade.
