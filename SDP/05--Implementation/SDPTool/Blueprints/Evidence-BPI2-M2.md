# BPI2-M2 — retained catalogue and consumer evidence

Baseline: 3d00521 on sdp/blueprint-implementation, isolated worktree
/tmp/sdp-blueprint-implementation. Candidate source/test and compiled executable
hashes are recorded in Pilot-BPI2-M2.json. English documentation and process
records accompany the code; no release or main integration is claimed.

## Delivered behavior

The diagnostic generator accepts --catalogue instead of --output. It captures and
checks inputs, retains complete bytes atomically in a content-addressed revision,
and refuses to overwrite changed retained content. Repeat identical publication is
idempotent. Existing blueprint revision means the same input/analysis identity;
retainedRevision additionally binds the complete file map, including manifest.
Retained bytes are logically immutable, not authenticated or write-protected.

Discovery finds SDP/Blueprints without registration, groups System/task/revision,
and publishes an additive Blueprints tab. Seven typed open targets expose generated
Markdown and Mermaid. Invalid entries have diagnostics and no targets. Historical
source captures do not register as live SDL/SDUI. Refresh detects byte edits even
with preserved timestamps. Limits apply across the catalogue, including invalid
entries: 256 task/revision candidates, 10,000 directory entries/file reads, 64 MiB
read bytes and bounded depth; enumeration permits one overflow probe before stopping.
Assignment state is unknown and live-source freshness is not evaluated.

## Reproduction and consumer fixture

Use Go 1.27.1, build from SDPTool:

```sh
go build -o /tmp/bpi2m2-final-sdptool ./cmd/sdptool
/tmp/bpi2m2-final-sdptool /tmp/bpi2-preview-pilot/models model create blueprint from work:Before to work:After --entry System.design --task /tmp/bpi2-preview-pilot/task.json --catalogue /tmp/bpi2-catalogue-pilot/SDP/Blueprints --json
/tmp/bpi2m2-final-sdptool /tmp/bpi2-catalogue-pilot discover --json
/tmp/bpi2m2-final-sdptool /tmp/bpi2-catalogue-pilot tree --json
```

The temporary model inputs reuse the reduced MVP1 Before/After artifacts from M1.
The compiled CLI regression creates fresh equivalent inputs without Git or installed
SDP. Consumer-BPI2-M2.json is actual generated discover output with local temporary
paths, not a portable registration file. A headless consumer reconstructed the tab,
confirmed discover/tree equality, followed every child ID, checked all seven open
paths against their target hashes, and confirmed no retained source was registered
as a live model. Pilot-BPI2-M2.json pins this fixture and the generated revision.

## Verification scope

Tests cover idempotence and multiple revisions, tampering, source/Markdown edits with
rewritten manifest hashes, missing/unmanaged files, wrong directory identities,
unsupported metadata, symlinks, bounded enumeration/reads, compiled command behavior,
refresh and discover/tree agreement. The existing navigation test's root count
increases only for the selected additive tab. The SDPTool SDL source model passes
the real parser/checker; no experimental syntax or new container is introduced.

## Independent review

Read-only agent /root/bpi2_catalogue_review identified three material gaps:
captured source hashes were not checked against source snapshot identities; directory
limits did not account for all invalid/overflow candidates; Mermaid files lacked
typed open targets. Source digests and a full retained-byte identity now protect
retention, both enumeration layers reserve their shared budgets and stop on overflow,
and both Mermaid files are exposed. Regression tests cover these corrections.
Final independent verdict: approved for bounded BPI2-M2, no remaining material
findings. The reviewer also reproduced the original rehashed-source and oversized
catalogue cases independently; both are rejected after fixes. Final reviewer race
regression passed in 242.963 seconds. This is not owner acceptance.

## Remaining boundaries

This is a headless producer/consumer contract, not a native XFMD integration test.
Readiness, assignment/implementation progress, independent generated-bundle trials
and evidence-backed lifecycle are BPI3. No installed gh-sdp upgrade or release was
performed. The unchanged installer race timeout recorded in BPI2-M1 remains an
outstanding final-integration check; no full-repository pass is claimed here.


## Verification results

- `go test -race . ./blueprints ./presentation -timeout 8m`: PASS. This full
  relevant-package run began before the final enumeration-budget corrections;
  the final focused rerun below covers those changes and their consumer regression.
- `go test -race . ./blueprints -run 'Test(Retained|Catalogue|BlueprintDiscovery|CompiledBlueprint)' -count=1`: PASS on final code (root 23.791 s; blueprints 255.279 s).
- `go vet . ./blueprints ./presentation`: PASS on the final code.
- `/tmp/bp2-sdl check SDP/SDL/ProjectGovernance/SDPTool/System.design`: PASS.
- Compiled retention/discover/tree and all seven headless targets: PASS.
- ProjectManagement and Toolkit validators and Git whitespace check: PASS;
  repeated after ledger closeout.

Raw race results are retained in Tests-SDPTool-BPI2-M2.txt. Exact tested source/test
hashes are recorded in Pilot-BPI2-M2.json; documentation-only closeout does not
change the compiled candidate.
