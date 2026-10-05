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

## MGI5-M2

2026-10-05T00:39:38.822402+00:00: Integrated candidate fd7033b passes SDPTool race suite and vet, SDL parser/sourcegraph and SDUI parser tests, compiled CLI lifecycle, and Windows amd64/macOS arm64 cross-builds. Independent fresh-context review approves bounded Linux implementation after regression fixes. Child-process recovery covers dirty restore at six boundaries. Canonical SDL activity is implemented and nine viewpoint diagrams were regenerated; product release/main merge remain excluded.

## Integrated verification and limits

Candidate: fd7033b, Linux, Go 1.27.1. Product source hashes: source-hashes.json.
Final closeout changes only documentation, ledger and SDL delivery projection.
Raw race output: race-results.txt. Independent disposition: Review.md.

Executed successfully:

- `go -C SDPTool test -race ./...`
- `go -C SDPTool vet ./...`
- `go -C SDL/go test ./parser ./sourcegraph`
- `go -C SDUI/go test ./parser`
- `GOOS=windows GOARCH=amd64 go -C SDPTool build ./cmd/sdptool`
- `GOOS=darwin GOARCH=arm64 go -C SDPTool build ./cmd/sdptool`
- `python3 SDP/ProjectManagement/validate.py`
- `python3 Toolkit/scripts/validate_sdp.py`

The compiled CLI test runs a full model journey in a fresh directory without Git,
including SDL/SDUI validation, conflict resolution, accepted release and copied-area
inspection. The separate transport test uses two real local Git clones, preserving
competing release identities while blocking an ambiguous default. Recovery tests
terminate child processes at unrecorded/building/payload/prepared/backup/installed
boundaries, including dirty restore. Cross-process writer contention is tested.

The strict Go structs/domain validator constitute the executable sdp-model/0.1
schema. No parallel JSON Schema registry or external source list was introduced.
Conservative merge limits, source/history size bounds, retained transaction backups,
local account attribution and unsupported non-Linux mutation are documented in
SDPTool/model/README.md. No native non-Linux or physical power-loss claim is made.
An arbitrary external editor is not locked out; optimistic comparisons detect
observed concurrent changes, without promising an atomic external-editor snapshot.

SDL check and static VP01/VP02/VP08 generation succeeded for revision
a4844c717e9ef27849ddda5efe01a27664359a07b910568f89ad3e6c73c99ae6:
9 diagrams / 38 generated files. Generated outputs were not edited by hand.
The activity status is now implemented; parsing a diagram alone is not runtime proof.

Backlog/onHold review retained KB050 as separate semantic blueprint work, and
KB036/037/042/043 as their existing governance/session scopes. No release packaging,
installer migration, XFMD change or new feature was silently added.
