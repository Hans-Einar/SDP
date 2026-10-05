# MG2 — bounded model artifact contract

Implemented contract 0.1, reconciled 2026-10-05 under PLAN-SDP-0019.
Owner-selected scope is in [Study](Study.md). See the [runtime contract](../../../../SDPTool/model/README.md)
for exact schema, limits, recovery commands and platform boundaries. This branch
delivery does not imply a published SDP release.

## 1. Commands and context

```text
sdptool model create work:NAME [from release:VERSION]
sdptool model create work:NAME --initial
sdptool model create work:Combined from work:A work:B
sdptool model commit work:NAME --message "Explain the change"
sdptool model restore work:NAME to commit:00003
sdptool model merge work:A into work:B
sdptool model create proposal:NAME from work:A
sdptool model create candidate:NAME from work:A
sdptool model create release:VERSION from candidate:NAME --evidence model-only
sdptool model snapshot work:NAME
sdptool model recover OPERATION-UUID resume|abort
sdptool model status work:NAME
sdptool model history work:NAME
```

The implementation exposes one explicit merge spelling;
context-relative pull is deferred, not an alias with different behavior. The release
version in the typed target avoids two competing version arguments. Existing
SDPTool project selection remains unchanged. Explicit paths may select an artifact
only after its metadata is validated; references resolve among direct sibling
artifacts of the selected model area. No recursive search through .merge/.commits.
A future standalone adapter uses the same library; no separate binary is required
for the first increment. The package must work with no SDP installation or Git.

Without `from`, WORK creation selects the single accepted head in that model area.
A release retains its predecessor UUID. One ancestry tip is selectable; independent
accepted tips or duplicate versions require explicit reconciliation, never highest
number/mtime. A single accepted release can be selected explicitly even when heads
are ambiguous; publishing a successor requires resolving that ambiguity first.
No releases means an actionable error suggesting --initial. --initial and from
are mutually exclusive. Names are case-insensitively unique for portability;
initially use ASCII letters/digits with internal hyphens/underscores, rejecting
reserved platform names, dots, slashes and leading hyphens. Versions are explicit
MAJOR.MINOR.PATCH in v1; automatic increment is deferred.

## 2. Artifact and source layout

WORK--NAME contains WORK--NAME.yaml, model sources, .commits and .merge.
PROPOSAL/CANDIDATE directory names append the final four UUID hex characters;
metadata holds the full UUID and human name. Fail without overwriting if a display path collides before publication. RELEASE--Vx.y.z is unique within the model area.
Its full UUID is metadata, not an extra name suffix. No existing artifact is replaced.
The metadata file stem always matches its enclosing directory stem.

A WORK with no history still has both local history directories. .commits/#00000
contains a baseline; an initial baseline has an empty file inventory. Git need not
preserve empty directories: the tool recreates missing empty .merge as equivalent.
Reject source paths that collide with reserved metadata/history/transaction names.
Source inventory owns all ordinary files in the artifact except those explicitly
reserved directories and metadata; no heuristic extension-based deletion of assets.
User-generated caches should live outside the artifact. Reject symlinks, nonregular
files, escaping paths and case collisions. Retained profile resources must be inside
the snapshot; no unpinned external dependencies for frozen artifacts.

## 3. Metadata shape and source identity

YAML is restricted to mappings/sequences/scalars, no aliases, custom tags or duplicate
keys. Reject unknown fields for schema 0.1. Read limits and inventory limits must be
explicit in implementation (proposed 128 MiB sources, 10000 source files, 16 MiB per
metadata file); no partial success on overflow. No timestamps from filesystem mtime.

The executable schema is the typed Go structures and domain validator in
SDPTool/model/records.go. WORK metadata has schema, UUID, kind, name, sequence,
head, sourceDigest, metadataDigest, optional baseRelease, and an embedded ledger
of complete records. Frozen artifacts add source-derived validationTargets;
release adds evidence, acceptedBy and optional predecessor. Test-created artifacts
are valid examples; there is no parallel hand-maintained registration schema.

`sourceDigest` is committed state, not a promise that mutable WORK is clean. status
reports liveDigest and dirty separately. Every commit record has id, kind
(baseline/commit/checkpoint/merge/restore), UTC recordedAt, recorded local author identity and original artifact name,
message, parents[], sourceDigest, inventory (path -> SHA-256), deleted[], and optional
restoredFrom/mergeBase. A commit ID uses the artifact UUID and a monotonically
increasing local counter; restore never reuses sequence numbers. files contains
full after-images stored under files/<relative-path>; baseline/checkpoint inventories
are complete. Other commit payloads are changes from their first parent.
Retain a full path/hash inventory for the resulting state as well for verification.

Digest algorithm: SHA-256 over UTF-8 domain `SDP-model-sources/1` then NUL, followed
by paths sorted by UTF-8 byte order; each path is uint64 big-endian byte length,
UTF-8 path bytes and its 32-byte SHA-256 content digest. Paths use forward slashes,
are relative and preserve case; source bytes are exact (no newline normalization).
Folders/modes/timestamps are not source content. Enforce portable normalized names;
no alternate Unicode-normalization aliases in schema 0.1. Hash manifests separately
from a deterministic sorted-key JSON encoding of their YAML data, setting
metadataDigest to an empty string. Parent links bind IDs to full embedded records; shared IDs with unequal records are rejected. Hashes detect
inconsistency; they do not authenticate an author or reviewer.

Frozen metadata replaces live commit references with a retained metadata-only
lineage graph and removes all payload references. Each origin includes identity,
original display name and source digest, not merely an absolute WORK path. Persist
all transitive metadata needed to browse origins, once per identity. Duplicate ID
with different immutable content is an error; lineage must be acyclic. Root schema
validation plus content-digest verification is mandatory before using frozen sources.
Exact machine-readable schema files/fixtures are MG3 inputs, not supplied by this
illustrative example alone.

## 4. Commit, restore and local publication

A commit saves changed after-images and deletions against the previous head. Empty
commit returns unchanged and does not create an entry. Rename is delete plus add.
Invalid/intermediate SDL may be committed for recovery; it cannot become candidate
until the configured validation targets pass. Commit records are immutable.

Acquire one cooperative writer lock per artifact; source editors are not assumed
to honor it. Build a staging directory outside source traversal, persist files and
records, verify expected head and captured source inventory again, then atomically
replace root metadata last. Journal prepared/committing/completed states and backup
paths permit explicit recovery. Never infer success from directory existence.
A crash before head replacement leaves an orphan to reconcile, not a visible commit.
Same-filesystem staging is required. Exact sync/rename portability must be exercised
on supported platforms before claiming power-loss durability. Initial acceptance
is Linux process interruption only; other platform guarantees remain unverified.

Restore reconstructs the chosen local state in staging, checks all hashes, and
saves a dirty pre-restore checkpoint if needed. Replace sources with a journaled
transaction and append a new restore record; failures retain backups and block
subsequent writers until explicit resume/abort recovery. Never recursively delete
unrecognized files. No selective branch undo or force-overwrite shortcut in v1.

## 5. Merge and archive handling

Use two retained input states and one unambiguous common ancestor. Initial v1
supports same-release divergence and common retained commit ancestry; missing or
multiple merge bases stop explicitly. Different release bases need an explicit
common retained ancestor; otherwise refuse rather than assume the newer release.
Repeated identical input integration reports unchanged by ancestry and content.

Capture exact dirty input bytes in the operation archive without changing the
source WORK's head. A target's dirty state gets a pre-merge recovery checkpoint.
Archive source/target YAML and both history directories before merge; stage outside
the trees and reject self-merge/ancestor-directory targets. Include captured current
sources as a checkpoint when required; YAML/history alone cannot recover dirty files.
Use identity-qualified archive names to avoid two unrelated same-name WORKs colliding.
Same archived immutable identity/digest may be reused; never deduplicate just by name.

Merge identical/one-sided file changes directly. Same-file text merge may use a
bounded three-way algorithm, exercised in MG3 before choosing dependency. Divergent
binary files, delete/modify, differing same-path additions and unresolved text edits
are conflicts. Retain three versions and a conflict inventory in WORK, never publish
conflicted data as candidate. Explicitly resolve and commit the result before making
a candidate. Pre-merge and final resolved merge checkpoints support whole-state
rollback; ancestor-branch commits are not selective edits to the integrated result.
Record both merge parents, chosen base and input digests. Revalidate combined sources.

## 6. Validation, promotion and acceptance

Validation targets are derived from captured source headers and dependencies at
promotion. The persisted list is a receipt, not a manually maintained registry.
Included SDL fragments must be reachable from a source-defined root; independent
roots are all checked. At least one supported SDL/SDUI entrypoint is required for
a frozen artifact. Empty initial WORK is allowed. Invalid intermediate WORK can
be committed for recovery but cannot be frozen.

Create candidate/proposal from a consistent read; do not silently include uncommitted
edits without reporting the captured digest. A frozen submission may retain captured
input as a new metadata-only lineage node; it need not add a persistent WORK commit.
Promotion omits all undo payloads, retains current sources and complete provenance,
then validates the staged artifact before publishing its new directory.

Release requires explicit integration-owner acceptance against the exact candidate
UUID/digest and a recorded evidence disposition: verified (named code digest and
checks) or model-only (implementation unverified). Neither code tags nor SDL parsing
prove implementation. Owner identity supplied in metadata is attribution, not access
control. No automatic authority gate is claimed without a trusted host. First release
from an initial candidate has no release parent; later release records its accepted
predecessor. Candidate content must be unchanged. Publication checks expected accepted
head and locks the local area; independent clones can still create competing heads.
Detect those on the next read before choosing a default. No release is a merge target.

## 7. Consumers and excluded work

The Go API returns immutable captured source maps, identity/digest, role, lineage and
payload availability. Read-only WORK preview can capture in memory without locking
WORK; recheck the inventory and retry boundedly or report changed-input. Do not claim
an atomic snapshot against arbitrary external editors; the result identifies the
bytes actually captured. Preliminary status remains visible.

SDPTool discovery must prune .commits/.merge/staging from ordinary source scans and
identify current/frozen artifact roots separately when feature support is installed.
Do not deploy an installer migration or alter current discovery before its scoped
implementation and compatibility tests. Future blueprint generation is KB050.
No permanent central object store, nested VCS, remote hosting, automatic distributed
locking, arbitrary script execution or separate mandatory PROPOSAL stage.

## MG3 evidence refinement — 2026-10-04

[Proof](Proof.md) confirms bounded after-image recovery and metadata-only promotion.
Presence checks must distinguish missing and empty files. File-level repeated merge
idempotence does not replace ancestry/event idempotence. The executable YAML subset
is mg-probe/0.1 only; the full sdp-model/0.1 machine schema remains a production
acceptance requirement. MG3 establishes one process-interruption boundary, not full
transaction durability or concurrent-writer safety. No new public command is shipped.

## MGI implementation refinement — 2026-10-05

Use `commit --resolved` to finalize a conflict after editing its live sources.
Conflict versions use base64 encoding so binary data remains readable under the
restricted YAML profile. Oversized metadata is rejected before publication.
Dirty-source capture IDs are operation-qualified, never future source commit IDs.
WORK-to-WORK merge is the bounded v1 surface; direct frozen-input integration and
semantic blueprints remain excluded. Full process recovery evidence and independent
review are in the [implementation evidence](../../../05--Implementation/SDPTool/ModelGovernance/Evidence.md).
