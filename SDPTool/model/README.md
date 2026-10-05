# Local model governance — sdp-model/0.1

The Go `model` package and `sdptool model` adapter manage bounded SDL/SDUI source
artifacts without Git, an SDP installation, an external merge command or a central
object store. This is model lifecycle management, not semantic blueprint generation.

## First use

Build with `go build -o sdptool ./cmd/sdptool` from SDPTool. Select an existing model
area explicitly; the default is the current directory. Existing sources are never
migrated automatically. Copy the intended source tree into the created WORK,
excluding caches, version-control folders and external links.

```sh
sdptool /path/to/models model create work:First --initial
# Copy SDL/SDUI sources into /path/to/models/WORK--First, then edit normally.
sdptool /path/to/models model commit work:First --message "Initial design"
sdptool /path/to/models model create candidate:First from work:First
sdptool /path/to/models model create release:0.1.0 from candidate:First --evidence model-only
sdptool /path/to/models model create work:Next
```

`model help` lists commands. `--json` selects the stable machine envelope; default
output uses the common Go presentation registry. Failed operations return nonzero
and a diagnostic. `status` and `history` include source/metadata digests and lineage;
`snapshot` adds captured file bytes (base64 in JSON). WORK snapshots are preliminary,
even when clean; reading never implicitly commits or freezes a WORK.

```sh
sdptool model restore work:Next to commit:00001
sdptool model merge work:Other into work:Next
sdptool model create work:Combined from work:Next work:Other
sdptool model commit work:Next --message "Resolve integration" --resolved
sdptool model create proposal:Optional from work:Next
```

Only WORK inputs merge in this increment. An optional proposal is an immutable
submission snapshot; direct proposal/candidate integration is not implemented.
The tool uses bounded conservative line-based three-way merge implemented in Go,
with no added dependency or third-party merge source. Missing/multiple common bases
are errors. Conflicts retain base/ours/theirs as explicitly base64-encoded bytes in
metadata. Resolve source files and commit with `--resolved`; real language validation
must pass. Renames are deletion/addition, not inferred identities. Rollback restores
a whole recorded state and preserves current dirty bytes as a new checkpoint first.
Selective inverse merges and cherry-picks are outside this feature.

## Storage and executable schema

`records.go` defines the strict typed YAML schema and domain validator. Tests are
executable positive/negative schema and workflow fixtures. Unknown/duplicate fields,
aliases, anchors, custom tags, multiple documents, bad paths/IDs/hashes/counters and
invalid lineage fail before publication. The serialized representation is validated
again before writing. No hand-maintained source registration file is used.

WORK directories contain `<directory-name>.yaml`, sources, `.commits/#NNNNN`
after-images/deletions and `.merge` archived input trees. Commit records hold full
result inventories; their file payloads contain changes from the first parent, except
full baseline/checkpoint/restore payloads. Frozen lineage has no payload paths.
Records preserve author (local OS account attribution), original artifact name,
parents and explicit restore/merge-base references. Dirty source captures use fresh
operation identities plus `origin`; they never consume the source WORK's counter.
Names use portable ASCII; UUIDs are full v4 IDs in metadata. Only proposal/candidate
folders show a four-character UUID suffix. Releases use explicit MAJOR.MINOR.PATCH.

SHA-256 source identity binds sorted paths and exact file bytes, excluding directory
modes/timestamps. Metadata uses sorted-key JSON of normalized YAML with
`metadataDigest` set to the empty string. This detects accidental inconsistency,
not malicious forgery or authenticated approval. Shared lineage identities must
have identical logical records before an ancestry fast path can succeed.

Source limits: 128 MiB, 10,000 files; metadata: 16 MiB; lineage: 10,000 records and
256 ancestor levels; retained artifact: 1 GiB/100,000 entries. History hashing streams
file data independently of source limits. Text diff caps its LCS matrix at two million
cells; larger ambiguous changes are conflicts. Oversized inline conflicts fail before
publication, leaving original inputs intact. Cleanup/compaction is not automated.

## Validation and accepted releases

Candidate/proposal creation validates every SDL/SDUI source through existing Go
libraries. SDL 0.6 roots and reachable fragments determine validation entrypoints;
orphan fragments and invalid sources reject promotion. `validationTargets` is a
receipt derived from the captured sources, never an editable discovery registry.
Frozen artifacts drop undo payloads but retain complete metadata lineage. Deleting
origin WORK directories therefore does not erase their recorded provenance.

Release requires an unchanged candidate and either `--evidence model-only`, or
`--evidence verified --code-digest SHA256 --checks REFERENCE`. The latter records
an attributed claim; it does not run those checks or prove implementation. `acceptedBy`
is the local OS account. The first release bootstraps the area; subsequent releases
must descend from its unique accepted head. Competing clone releases survive as
ordinary files, but block automatic head selection; highest version is not authority.
Same-version Git conflicts require explicit human reconciliation, never auto-selection.

## Recovery, platform and consumers

Linux mutation uses an OS-released advisory area lock and same-filesystem staging,
with retained backups and journals under `.model-operations`. A pending operation
blocks further writes and gives its UUID. Inspect the diagnostic, then use:

```sh
sdptool model recover OPERATION-UUID resume
sdptool model recover OPERATION-UUID abort
```

Incomplete builds require abort. Unjournaled staging can be aborted and is retained
under `aborted-UUID`. A replaced/externally edited target cannot be blindly aborted.
All backups remain until an explicit future cleanup workflow; no silent expiry.
Do not modify reserved metadata/history/staging by hand. Cooperative tool writers
are serialized; external editors are detected with repeated capture and optimistic
whole-tree comparisons, not locked out. No atomic snapshot against arbitrary editors
or physical power-loss durability is claimed. Windows/macOS compile but mutation
returns unsupported; native acceptance is Linux only.

Within SDP discovery, recognized artifacts expose kind/UUID/preliminary metadata.
Only their `.commits`/`.merge` trees and recognized transaction areas are pruned.
Unrelated folders with those names remain ordinary sources. Views are regenerated
from sources; no navigation.json is created. Artifact-only areas use `model` APIs
without needing process templates. Installer/release packaging is separate work.
