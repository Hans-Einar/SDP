# MGI implementation evidence

Milestone checks are recorded as executed; platform scope is Linux process recovery.
Windows/macOS mutation and physical power-loss acceptance are not claimed.

## MGI1-M1

2026-10-05T00:16:32.082240+00:00: WORK creation/status and strict metadata foundation implemented, with staged area-locked publication. SDL/SDUI validation and recovery library scaffolding included but later CLI operations not yet exposed. Full SDPTool test suite passes.

## MGI2-M1

2026-10-05T00:17:57.226155+00:00: Commit/history/whole-state restore and explicit recover resume/abort CLI implemented. Race tests pass, including child-process interruption at prepared/backup/installed boundaries, dirty preservation, corrupt payload and external edit refusal. Backups retained; physical power-loss and non-Linux mutation not claimed.

## MGI3-M1

2026-10-05T00:20:56.033163+00:00: Native bounded Go three-way merge, combined WORK creation, persistent conflict inventory and resolved commits implemented. Tests cover disjoint and overlapping same-file edits, unrelated bases, dirty inputs, repeat integration without extra events, three archive generations and whole-state rollback. No external Git merge dependency.

## MGI4-M1

2026-10-05T00:28:20.994585+00:00: Frozen candidate/proposal CLI now validates real SDL source graphs and SDUI using existing parsers, preserves metadata lineage and drops undo payloads. Model tests pass. Independent review found dirty-capture identity and file/directory restore defects; regression fixes and stricter domain validation are included, with final re-review pending.

## MGI4-M2

2026-10-05T00:30:41.202655+00:00: Release promotion requires candidate validation and explicit model-only or verified evidence attribution. Default WORK resolves the unique accepted head; stale and competing releases fail closed. Model tests pass including actual two-clone Git transport. Review fixes add abort of unjournaled staging and base64 conflict values with bounded YAML round-trip validation before publication.

## MGI5-M1

2026-10-05T00:36:50.759403+00:00: Artifact-aware discovery exposes kind, UUID and preliminary role, prunes only owned histories and transaction staging, and preserves ordinary projects. Read-only snapshot returns captured bytes/digest without a commit. Full SDPTool tests pass, including compiled no-Git CLI lifecycle with real SDL/SDUI validation and copied-area inspection. Provenance now includes local author/acceptor attribution, original names and merge/restore references.

