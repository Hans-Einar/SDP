# Canonical project-history append — BPI3-M2a

This package provides a shared persistence primitive for
SDP/ProjectManagement/Ledger.ndjson. It implements no KanBan/assignment state machine,
actor authorization or counter allocation. The calling domain must validate its
payload, transition, authority and references before calling Append. No installed
CLI consumes this package yet; BPI3-M2b owns integration and validator extensions.

```go
snapshot, err := projecthistory.Read(ledgerPath)
// Validate domain state against snapshot.Bytes and construct one pinned event.
result, err := projecthistory.Append(ledgerPath, snapshot.Revision, eventJSON)
```

Read returns exact owned bytes and SHA-256, including when the ledger does not yet
exist (empty content hash). Its parent directory must already exist for mutation.
The reader rejects incomplete tails, duplicate IDs, unsafe paths, nonregular files,
invalid generic envelopes and streams larger than 16 MiB. Existing blank lines and
formatting are retained. Domain extension payloads are not interpreted here.

Append accepts one generic envelope 1.0 event, compacts only that new JSON object,
and appends exactly one newline. It preserves all previous bytes and existing file
permissions. An event ID repeated with identical compacted bytes is an idempotent
retry even when later events exist. Different bytes under an existing ID return
ErrConflict. Stale expected history returns ErrStale; a competing cooperating
writer returns ErrBusy. These conditions do not alter the authoritative stream.
IDs, timestamps and event bytes must be kept fixed when retrying. A different
property order is different request content; this is not semantic JSON merging.

On Linux, writers share a nonblocking flock on `<ledger>.lock`; the sidecar is
stable across ledger replacement and must not be deleted during operation. They
stage in the ledger directory, sync the complete candidate, recheck the original
stream, atomically rename and sync the directory. Before-rename failure leaves the
old stream; after-rename failure can mean committed-but-unacknowledged. Repeat the
exact request to reconcile. A process exit releases the lock. Orphaned stage files
are never authoritative and readers ignore them. Do not truncate or repair malformed
history automatically. Read works on other platforms; mutation returns ErrUnsupported.

This is logical append-only history implemented using atomic replacement. It is
not a second ledger, a private assignment snapshot or a full transaction across
other files. Tests prove process-exit behavior on the tested Linux filesystem;
they do not prove every storage device's power-loss semantics. Arbitrary editors
and existing manual writers do not acquire this lock. The last recheck detects
observed interference but cannot prevent an uncooperative same-user writer from
racing after it. Other programmatic writers must adopt this protocol before a
shared-writer safety claim is made.

Run `go test -race ./projecthistory` and `go vet ./projecthistory` from SDPTool.
Tests include the actual historical envelope set, preserved bytes/permissions,
retry after later events, CAS/conflicts, concurrent writers, process exits before
and after rename/sync, cross-process lock contention, malicious paths, capacity
and malformed streams. No test writes the repository's actual history.
