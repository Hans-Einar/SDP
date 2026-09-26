# Small Fix Policy

SDP uses **Fix records**, not Micro Slices or free-floating revision records.

A normal Fix ID is `FIX-<target-version>-NNN`, for example `FIX-0.8.1-001`.
An emergency production correction may use `HOTFIX-<released-version>-NNN`.

Every Fix record states:

- why a full Iteration/Slice is disproportionate
- bounded scope and expected files
- invariants and non-goals
- relationship to the active or next release target
- required verification and review
- optional development revision
- ledger and Relations entries

A Fix is prohibited when the work adds a capability, changes architecture,
changes a public contract, requires a broad migration or contains unresolved
security impact. Such work requires a planned Iteration and Slice.

Emergency Fixes may start before the normal planning ceremony only to reduce
active harm. They still require a record immediately, complete verification,
independent review and a normal PATCH release transaction.
