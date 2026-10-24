# PGI4-C3 evidence

`checks.py` builds both final binaries and runs full SDPTool race/vet checks. The
Go publication tests cover prepare/apply, nine interruption points, failed intent
persistence, exact retries, conflicting source/evidence/authority/documents, symlinks,
contradictory results, shared project claims and read-only coordination/leader
selection. The optional canonical-record test also inspects this repository's
125 management files without writing them.

`lab_checks.py` repeats inspection/apply through the compiled client, validates the
canonical disposable fixture with the existing Python management validator, and
performs an actual stdio MCP negotiation/list/read. It also inspects the retained
C2 native review stores without modifying their heads. No new model runs occur.
The publication fixture's review is explicitly synthetic protocol evidence.
`schema_checks.py` validates actual returned wire objects. `old-reader.json` records
actual C2 refusal of profile 0.12; `candidate.json` pins final source bytes.

`preliminary-lab-wire.json` is an earlier disposable prototype trial; its temporary
project was removed, retaining the wire artifact. `preliminary/` retains the passing
full suite and compiled fixture before correction of an unresolved JSON-schema
reference. Final checks and final schema objects are in this directory. The two
retained lab fixture directories distinguish pre-correction and final execution.

Operational stores and binaries are outside Git. Exact source/evidence hashes and
canonical published snapshots are committed; this is not a distribution/release
claim. `manifest.json` excludes itself and the generated final record-check report.
