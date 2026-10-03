# ModelGovernance — initial design

Status: working design, MG2 in progress; no production schema or CLI implemented.
[Study](Study.md) identifies owner decisions; recommendations below are not silently
promoted to accepted language rules. [Plan](Plan.md) owns delivery milestones.

## Responsibility and boundaries

SDPTool coordinates model artifact creation, local checkpoints, recovery, merge,
validation and promotion. A reusable Go package should own these operations; CLI
argument parsing and presentation are adapters. Standalone executable packaging is
undecided; do not duplicate logic or introduce runtime command plugins.
SDL and SDUI libraries own their syntax and semantic validation. ModelGovernance
owns revision identity and consistent source access, not a new language parser.
Blueprints consume model source views plus provenance; they own semantic differences,
affected context, diagrams and assignment constraints. ProjectManagement owns work
lifecycle; the model ledger owns model edits. They reference rather than replace one another.

## Artifact lifecycle

| Artifact | Behavior | Promotion / modification |
| --- | --- | --- |
| WORK | Mutable sources and local undo payloads | Commit/restore/merge; create candidate or optional proposal |
| PROPOSAL | Optional frozen submission | Input to a new WORK or candidate; never edited |
| CANDIDATE | Frozen integration target | New candidate for changes; accepted exact content can become release |
| RELEASE | Immutable accepted model | Source for new WORK; never a merge destination |

A release is acceptance of a particular model with explicitly scoped evidence, not
an automatic claim that all modeled behavior is implemented. Record model/code
artifact digests and actual checks when evidence exists; no Git dependency is
introduced for these identities. Unknown or absent implementation evidence stays
visible. Release evidence policy and initial empty-project bootstrap are MG2 decisions.

## Proposed directory and metadata contract

```text
WORK--Fix-Navigation/
  WORK--Fix-Navigation.yaml
  <model source tree>
  .commits/
    #00000/   # full initial source snapshot
    #00001/   # changed-file after-images and deletion record
  .merge/
    <archive identity>/
      <original artifact YAML>
      .commits/
      .merge/
```

YAML shares the artifact directory stem. Full UUID identifies the artifact; names
are labels. Metadata should include schema version, kind, base release reference,
current commit head, parent identities, source profile inputs and ledger references.
Commit records carry artifact identity plus local sequence, message, parents,
source digest and changed/deleted paths. Counters do not identify records globally.
WORK names omit UUID. Frozen proposal/candidate names may use four final UUID
characters; collision detection must never overwrite an existing artifact.

Sources exclude history, metadata and generated caches. Propose rejecting symlinks,
escaping paths and portable case collisions initially. SHA-256 source identity is
computed from a canonical sorted inventory of normalized relative paths and exact
content hashes. Exclude the manifest from its source-content hash; separately cover
immutable metadata without self-referential hashing. Specify encoding in MG2.

## Operations and failure behavior

| Operation | Required behavior |
| --- | --- |
| create work:NAME | Resolve one accepted release or explicit --initial; copy baseline, persist identity; fail on existing destination |
| commit work:NAME | Save after-images and deletions versus last head; publish content before head; retain WORK editability |
| restore | Reconstruct whole local model at selected checkpoint; preserve dirty state; record recovery without erasing history |
| create work:Combined from A B | Read stable input states, establish common base, merge in new WORK; preserve inputs |
| merge SOURCE into WORK | Preserve pre-merge state; conflicts remain explicit in WORK; no immutable artifact overwritten |
| create candidate from WORK | Require consistent source view and successful supported validation; freeze new identity/digest |
| create release from CANDIDATE | Check identity, acceptance policy and expected current release; publish a new immutable record |

Command spelling above is provisional except owner-selected verb-first and typed
WORK target direction. Context-relative pull, restore syntax, empty commit policy
and automatic minor/patch selection remain open. Resolve references in one declared
model area; never mutate an ambiguous cwd target or silently pick competing heads.

## Recovery and merge

Keep a full initial WORK baseline and complete pre/post-merge recovery checkpoints.
Numbered commits store whole changed files and explicit deletions; rename can begin
as delete plus add. No patch engine is necessary for recovery. Restore appends a
new record preserving prior history. Old incoming branch commits are browsable
provenance; they are not automatically selective restore points on the merged line.

Use a known common ancestor for three-way merge. Disjoint changes may combine;
overlap, delete/modify, ambiguous rename, missing/multiple bases and added-path
collisions require explicit handling. Same-file text merging is a bounded experiment,
not proof of semantic correctness. Revalidate the combined model and evidence.
Persist exact merge input states/IDs to make repeated incorporation detectable.

Archive source and destination YAML/.commits/.merge from before integration. Stage
outside the copied trees, reject self/descendant traversal, and reuse identical
archived identities where practical. Same identity with different bytes is an
integrity failure. Dirty input checkpoint policy must preserve edits rather than
silently discard them. Crash recovery and concurrent local writers require a small
transaction protocol; no distributed lock authority is promised.

## Promotion and lineage retention

Copy only current sources plus complete metadata-only provenance into frozen
artifacts. Preserve commit messages, parent edges, origin WORK identities and merge
resolution identities. Drop .commits after-images and archived source payloads.
History must not point only at deletable WORK paths. Label non-restorable history
explicitly. Tree displays reference shared ancestors; underlying lineage is a DAG.

Independent artifact identities survive parent Git as ordinary files. Duplicate
human release labels or competing accepted heads are validation conflicts, even
when Git merges text cleanly. Do not use highest version or modification time as
acceptance authority. Parent Git cannot restore uncommitted discarded payloads.

## Blueprint and navigation boundary

Expose consistent source bytes, identity/digest, role, provenance and availability
to a later blueprint consumer. WORK preview captures a stable read without locking
WORK or mandatory persisted commit; mark preliminary and detect concurrent edits.
Retained candidate/release comparisons remain possible after WORK payload removal;
arbitrary discarded commit comparisons do not. Blueprint generation itself belongs
to KB050. Discovery must distinguish live roots from .commits/.merge copies so
history does not appear as duplicate systems. Existing project discovery must not
be changed implicitly by merely adding this design document.

## Design completion and verification cases

MG2 must resolve the schema and recovery decisions above, model the feature using
supported SDL in the canonical source tree, and produce an ImplementationPlan.
MG3's proof must cover add/edit/delete/restore, dirty state, two-source integration,
same-file conflict, candidate promotion without undo payloads, retained lineage,
interrupted publication, repeated merge, competing release labels and clean-copy
operation without project Git. No probe result is claimed in this document.
