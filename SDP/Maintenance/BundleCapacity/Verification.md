# BC1 verification

| Field | Value |
| --- | --- |
| id | SDPTOOL-VER-BC1 |
| Plan | MAINT-SDP-0021 |

Baseline 23c5193 plus [candidate hashes](candidate.json) implements shared sealed
preview/source-reading output bounds of 512 files and 128 MiB. Both bounds are
inclusive. Overflow errors retain code `limit` and now state actual counts/bytes.
Source validation, freshness, protected output and publication ownership remain.

Focused Go tests pass ([log](tests.log), 0.148s): exact 512-file and 128-MiB acceptance,
overflow rejection,128-source composed-model publication through both Preview and
SelectProject, stale sources, failed-refresh preservation and existing source links.
The generated fixture needed correction to put all includes before semantic facts;
no product change was needed for that fixture correction. No full suite, separate
review or native agent execution was commissioned for this bounded parameter change.

Command: local Go 1.27.1, `go -C SDPTool test . -run
'TestPreviewBundleCapacity|TestPreviewAndSourceReadingBeyondOldBundleLimit|TestPreviewFreshSourceAndPreservation|TestPreviewFailuresPreserveSource|TestFailedRefreshKeepsPreviousBundle|TestReadingMembershipCollectionsAndSources'
-count=1 -v`, with the existing isolated GOCACHE/GOMODCACHE and GOTOOLCHAIN=local.

Source implementation is complete. Public release, installed engine replacement
and reproduction against the owner's unidentified original model are not claimed.
