# MG3-M1 — WORK-local history proof

2026-10-04. Bounded filesystem feasibility evidence. This does not deliver public
ModelGovernance commands or prove the entire MG2 contract.

## Candidate and reproduction

Experiment: experiments/model_governance, standalone Go module with yaml.v3 v3.0.1.
Source/test hashes and environment are in [proof/manifest.json](proof/manifest.json).
Raw Go test results are in [proof/test-results.ndjson](proof/test-results.ndjson).
Baseline repository commit is recorded in the manifest; the added experiment is
identified by hashes, not falsely described as already committed when tested.

```sh
go -C experiments/model_governance test -race -count=1 -json ./...
go -C experiments/model_governance vet ./...
```

Go 1.27.1 linux/amd64 was invoked from the local SDP toolchain because `go` is not
on the shell PATH. All fixture operations use temporary directories; no real model
WORK is initialized and no installer or XFMD state is changed.

## Observed results

13 top-level tests passed, including 5 conflict-matrix subcases; race execution
reported no failures and go vet succeeded. Race instrumentation on these tests is
not evidence of support for competing process writers.

| Behavior | Observation | Boundary |
| --- | --- | --- |
| Baseline and changed-file commits | Unchanged file not recopied; after-images and deletion list reconstruct expected model | In-memory expected maps compared to filesystem reconstruction |
| Rename/edit/delete | Delete plus add restores both pre/post states | No rename inference |
| Dirty state before restore | Dirty contents preserved in preceding checkpoint; restore adds a higher sequence | No live multi-file replacement transaction |
| Missing/corrupted content | Tampered after-image and missing empty file rejected | Not full hostile metadata testing |
| Initial empty WORK | Baseline exists; unchanged commit leaves head unchanged; existing destination refused | Frozen empty-model acceptance not implemented |
| Portable history | Copied store reopens after deleting original; no .git exists | No parent Git merge/release-label experiment in this milestone |
| Process interruption | Child exits after payload/pending write, before head replacement; reopened old head remains valid, next writer blocked, resume publishes checked pending state | Controlled Linux process exit, not arbitrary kill points/power loss |
| Corrupt pending payload | Resume refuses corrupted content and keeps old head | Full journal recovery policy remains implementation work |
| File merge matrix | Disjoint/identical changes succeed; overlap, delete/modify and differing additions produce conflicts | File-level primitive, not SDL semantic merge |
| Same-file three-way comparison | git merge-file combines distant edits in one file and reports conflicting edits | External Git binary used only for this comparison, no repository; Go dependency decision remains open |
| Merge archive and whole-state rollback | Both pre-merge stores archived outside target traversal; merged state saved; target can restore pre-merge state | One merge archive depth tested, not unbounded recursive history |
| Promotion | Current sources and 5 lineage records survive deletion of source WORKs; no .commits/.merge payload in candidate | Prototype YAML layout, no production publication/acceptance gate |
| Lineage identity | Duplicate equal record deduplicated, same ID with differing message rejected | Cycle/depth checks and immutable metadata hashes not implemented |
| Repeated integration | Repeating the file merge keeps resulting contents identical | Does not yet suppress a duplicate history event by ancestry |
| YAML | Unknown/duplicate fields rejected | Strict alias/tag/size/profile schema restrictions still required |

## Findings and design consequences

The bounded owner direction is feasible: full starting state plus changed-file
copies/deletions can provide local recovery; complete model snapshots are not
required for every commit. Source histories can be retained as metadata at promotion
without a permanent central repository.

A missing empty file must not be interpreted as an existing zero-byte file. The
probe now checks presence independently from content hash. Dirty merge inputs
require captured current bytes, not just YAML plus earlier commits. Whole-state
recovery is possible around merges; selective inverse merge remains excluded.

Payload-first, head-last publication protects the logical commit boundary at the
injected failure point. Production still needs a complete transaction journal,
writer ownership, stale-input detection, filesystem sync policy and retry/recovery
behavior for every mutation phase. No assertion of general crash safety follows.

## Deferred proof obligations for implementation

- Complete machine schema/fixtures for sdp-model/0.1, UUIDs, canonical metadata
  digests, limits, portable paths, aliases/tags/symlink rejection and lineage cycles.
- Transactional source-tree restore, dirty editor races, lock acquisition/recovery,
  corruption at other phases and Windows/macOS filesystem behavior.
- Multi-generation .merge copying/reuse, merge-base selection, ancestry idempotence,
  conflicting release labels after project Git merges and accepted-head resolution.
- SDL/SDUI validation at candidate creation, readonly payload stripping on actual
  artifact layout, evidence policy and integration-owner publication behavior.
- Public CLI/discovery compatibility, renderer output, application integration and
  independent review. The same agent authored and ran this proof.

These limitations must become ImplementationPlan acceptance criteria. They are not
silently marked verified by the passing experiment. MG4 is the next handoff step.
