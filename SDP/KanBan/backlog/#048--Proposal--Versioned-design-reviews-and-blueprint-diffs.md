# Versioned design reviews and semantic blueprint differences

| Field | Value |
| --- | --- |
| id | KB-SDP-048 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| Systems | SDL, SDUI, SDPTOOL |
| created | 2026-10-01 |
| source | Owner discussion after XFMD dynamic-tab implementation |
| tags | blueprint, model-history, baseline, design-review, storage |

Session: [SESSION-SDP-0005](../../Sessions/session-%230005--Blueprint_model_history.md).

## Owner outcome

The XFMD agent changed implementation before updating its SDL/SDUI design. The
owner proposes independent local version control for design sources: create a
model branch, design a change, review its proposed difference, generate a blueprint
showing the affected system, implement against that blueprint, verify code against
the target model, then integrate the verified design. The initial preference avoided a directory copy for every release; the owner now
explicitly proposes full snapshots as a simpler alternative (see below). Avoid
requiring a separate GitHub model repository. Preserve
both SDL and SDUI revisions together and let SDPTool manage the workflow.

This card captures a proposal, not approval of a storage backend, new commands,
language syntax, installer migration or automatic implementation proof. No local
model repository has been initialized. Existing source/installation layouts and
XFMD remain unchanged. The completed BP1 study is not reopened.

## Existing authority and current limits

[BP1 study](../../04--Design/SDPTool/Blueprints/Study.md) already calls for pinned
NOW/TARGET models, complete modeled impact context, constraint marks and separately
observed implementation evidence. [BP2 DesignPlan](../../04--Design/SDPTool/Blueprints/Plan.md)
is still planned. [KB-SDP-004](%23004--Proposal--Design-traceability.md) owns mappings
between model, code and evidence; no replacement ledger should be created here.
The source-composition and discovery work delivered after BP1 changes its historical
single-file assumptions, but does not itself implement semantic model differences,
assignment generation, local design-review requests or code-conformance proofs.

## Owner constraint — no separately maintained history export

Owner clarification, 2026-10-02 local date: reject git bundle and any workflow
whose durability depends on SDPTool or an agent remembering a separate export.
History must be durable in its authoritative native storage after each model
commit and capable of accompanying the project. The earlier local Git plus bundle
recommendation is withdrawn. Subversion is an owner-suggested candidate; no backend
is selected. Ordinary parent Git commits remain necessary to publish project data,
but a second model-history export must not be a hidden prerequisite.

## Git facts that constrain the proposal

Git supports local repositories without a network remote. Commits, branches,
tags and local clones are available; a hosted pull-request discussion/approval
record is not a native Git object. SDPTool would own a local design-review record.
A release can reference an immutable model commit/tag without copying source trees.

An ordinary repository embedded inside another does not automatically publish its
files/history through the parent. A submodule records a commit pointer; the objects
must be obtainable separately. If models must follow the parent as readable source,
keep ordinary files there and explicitly manage the separate model store/export.
Git bundle can transport selected reachable objects and refs, but not review records,
working-tree edits, reflogs or every repository setting. No-ref/garbage-collected
proposal history cannot be assumed recoverable.

Primary references:
- https://git-scm.com/docs/gitsubmodules
- https://git-scm.com/docs/git-bundle
- https://git-scm.com/docs/git-init
- https://docs.github.com/en/pull-requests/reference/pull-requests

## Storage alternatives to decide in BP2

| Option | What it provides | Consequence |
| --- | --- | --- |
| Project Git plus a separate design workflow | Native branch/PR merging and portable source history without export | Recommended for reconsideration under the parallel-merge requirement; changes the earlier preference for a physically separate repository |
| Local model Git plus bundle export | Independent history with portable export | Rejected by owner: requires another persistence/export operation |
| Normal submodule | Independent Git history with parent pin | Local-only objects do not follow the parent's push/clone |
| Local Subversion repository | Atomic tree revisions and reproducible differences within one repository | Independently changed copies of its repository cannot safely be merged as ordinary Git-tracked backend files |
| Fossil repository | Native single-file SQLite repository and its own history/merge operations | Git merging divergent database files does not perform a Fossil history merge; no automatic solution to the new requirement |
| RCS/CVS | Established predecessors supporting revision differences and history | RCS is per-file; coordinated model revisions and modern concurrent workflows need additional assessment |

Investigate native authoritative storage, not another dump/export that must be
remembered. SVN does not intrinsically eliminate the problem: versioning only its
working copy preserves current source, not the actual repository history. Native
repository storage could be included separately, but capturing a live repository
consistently is not equivalent to adding ordinary source files. Its hotcopy/dump
operations are backup mechanisms, not a selected workaround for the rejected extra
step.

The earlier Fossil suggestion addressed export because its repository itself is a
single SQLite file rather than an exported snapshot. A parent-tracked native
repository would not require a bundle step. That alone does not prove our workflow:
Git cannot text-merge divergent database copies, capture must not race a database
write/journal, and source/history correspondence must survive a fresh parent clone,
rollback, parallel Git worktrees and interrupted operations. Fossil repositories
also contain auxiliary configuration/user data; assess suitability before treating
that native database as a publishable project artifact. Do not introduce an extra
scrub/backup export and silently call it compliance with the owner's constraint.

The new parallel-parent-merge requirement below supersedes any inference that SVN
or Fossil native storage alone solves this problem. No storage layout is adopted.
Keep the semantic blueprint/review contract separate
from the storage adapter. Known version control supplies revisions and text diffs;
SDL/SDUI semantic impact and implementation evidence remain our tooling's job.

Additional primary references:
- https://svnbook.red-bean.com/en/1.7/svn-book.html
- https://svnbook.red-bean.com/en/1.8/svn.reposadmin.maint.html
- https://fossil-scm.org/home/doc/trunk/www/quickstart.wiki
- https://fossil-scm.org/home/doc/trunk/www/index.wiki
- https://www.gnu.org/s/rcs/manual/html_node/Overview.html

## Parallel parent merges — owner clarification and recommendation

Owner clarification on 2026-10-02: independent design commits on parallel project
Git branches must survive PR integration with safe merges. Any retained pair of
design revisions should be usable as reproducible blueprint inputs. A local store
that works only until two project branches diverge does not meet the requirement.

Analysis: a single local SVN repository can compare committed revisions and a
working copy. However, copying that repository into Git branches produces two
independently writable histories with a shared origin. Both can allocate revision
43 to different transactions. Git merging repository backend files does not reconcile
SVN revision identity, references and transactions. Similarly, choosing one of two
Fossil database files loses the other side unless a Fossil-level reconciliation is
performed. These are architectural deductions from the storage models; no SVN or
Fossil runtime test was performed in this turn (svn/svnadmin are unavailable).
Custom Git merge drivers would require installed configuration and actual invocation
at every merge location, including hosted PR merges, and semantic validation after
integration. A committed .gitattributes file alone does not install that mechanism.
They are not an implicit solution to ordinary GitHub PR merging.

Recommended change to the proposal, requiring owner disposition: separate the
**design lifecycle and view**, while using the project's existing Git object/history
store. SDL and SDUI remain ordinary versioned sources. SDPTool presents model-only
history, proposal workspaces and blueprint operations instead of asking users to
manage a second database. This is a deliberate tradeoff against the earlier
physical-repository separation preference, not a claim that preference was accepted
or silently superseded by tooling.

Proposed sequence:

1. Pin the checked implementation/model baseline at project commit B.
2. Create an ordinary project branch/worktree; commit only the intended model
   changes as D. Preserve unrelated staged work; the tool needs scoped staging or
   an isolated worktree/index, not an unrestricted Git commit.
3. Generate and review the semantic B-to-D blueprint before implementation.
4. Commit implementation and verification as I on that branch. Model changes after
   D require a fresh target/blueprint review; do not quietly reuse an old approval.
5. Integrate through ordinary Git merge. Revalidate the actual combined model and
   implementation, even if Git reports no textual conflicts. Approval of A and B
   separately does not prove their combined contracts.

A commit's model input is the complete resolved model source closure at that
commit, not only its changed lines or the working directory's current sources.
The root model subtree hash can identify unchanged source content across code-only
commits, but full revision identity also includes external/pinned inputs, language
profiles and generator/selection policy. Commit-specific code mappings/evidence
must remain associated with the corresponding code revision.

For the requirement that every reviewed design commit survive, use normal history-
preserving merges after review and retain reachable commits. Squash/rebase can
remove or rewrite intermediate identities. A SHA written in a document does not
by itself preserve an otherwise unreachable Git object. Define repository policy
and enforceable checks; no branch-policy configuration has been changed here.
A merge commit has multiple parents: comparison must select an explicit baseline
(first-parent integration view or each-parent views), not guess a single previous
model. A code-only commit can legitimately produce an empty model difference.

If physical storage separation remains non-negotiable, the viable alternative to
study is an immutable, content-addressed store in ordinary parent-tracked files:
unique snapshot/commit records, parent links and explicit conflicting branch heads.
Disjoint records can coexist after a Git merge, but model branch integration and
validation still need a designed protocol. This amounts to implementing substantial
version-control machinery and is not a recommended first delivery or a selected
new format. Raw copied VCS databases are not equivalent to that store.

### Bounded Git experiment — 2026-10-02

A disposable repository outside the project used plain textual .design fixtures;
no parser, blueprint generator or implementation-conformance test is claimed.
Git 2.52.0 passed all five checks:

- A design commit followed by a code-only commit retained the exact model subtree.
- Two branches changing different model files merged and retained both changes.
- A normal merge retained both parent histories.
- A fresh non-local-optimized clone retrieved exact baseline/design/code snapshots.
- Conflicting edits to the same model line stopped with an explicit merge conflict.

Temporary probe: /tmp/sdp-model-history-probe-e5l3okel/Evidence.json.
Baseline: 34a30bcc5a96579d72c5022f0896c7aac75066e0.
Design A: 6e7d0b42f762168866c1b8fbaaefd26043dc1a93.
Implementation A: 0677b697445766602d959d3f9646c6d3bfdc28c5.
Design B: e57b6632b173fd354ee4fe841cf4e3af792a5acf.
Merge: 6f0c8c531f12f9915f12efbea6766fc04b6add81.
These temporary identities demonstrate storage behavior only and are not product
release or accepted system design revisions. No project Git metadata was modified.

Primary sources for these mechanisms:
- https://git-scm.com/docs/git-merge
- https://git-scm.com/docs/git-worktree
- https://git-scm.com/docs/git-log
- https://git-scm.com/docs/gitattributes
- https://svnbook.red-bean.com/en/1.7/svn.basic.in-action.html

## Baseline and review semantics

Distinguish approved design, proposed design and verified implementation baseline.
A model commit or branch named implemented is not implementation evidence. Link:

- immutable baseline model revision (including SDUI and all resolved sources),
- immutable baseline implementation revision,
- immutable target model revision,
- actual resulting implementation revision,
- selected language/tool profiles, verification scope/results and reviewed exceptions.

Use existing Traceability for verification links and existing project management
history for review lifecycle. Reuse these authorities rather than introducing a
manually synchronized navigation/index file. Pin reviewed revisions; advancing a
branch invalidates approval/bundles affected by that change. Runtime store paths
are local implementation details, never portable identities.

A local DesignChange (working name, not an adopted schema) serves the PR role:
intent, exact base/target, generated blueprint, review state, constraints and evidence.
Design approval authorizes the intended change before coding. Implementation review
later confirms the checked result; only then advance the verified baseline. A
rejected/canceled proposal keeps the verified baseline unchanged. Release identity
records the model/code pair, not a second copy of the sources. Partial implementation
and incomplete verification remain explicitly scoped, not whole-system verified.

A model-store transaction is not automatically atomic with a parent Git commit.
Define recoverable promotion and a code/model receipt; concurrent proposals require
base validation/rebase, conflict resolution and renewed affected review. Never
silently overwrite changed accepted source or consume an unrelated worktree state.

## Owner refinement — standalone snapshot history, 2026-10-02

The owner identifies a remaining failure in the project-Git recommendation:
blueprints must work without a Git repository and without a project commit. The
owner proposes immutable release directories and editable WORK copies, initially
restricted to work starting from the latest release, with explicit file checkout,
commit and pull. This supersedes the earlier preference against release copies as
a constraint on investigation. No backend or command syntax is adopted yet.

**Recommended direction for the next design experiment:** ordinary-file snapshot
history owned by SDPTool, independently usable without Git. Parent Git transports
those authoritative files directly when used; no bundle/export/rebuild step owns
history. Start with full snapshots; deduplication can wait. This is a bounded
version-control subsystem, not merely a rename command.

### Minimum durable identities and operations

- Each immutable revision has a unique ID, parent revision IDs, a complete SDL/SDUI
  source snapshot and a manifest of paths/content digests. Include or pin the full
  source closure and parser/profile inputs needed for blueprint reproduction.
  UUID-style identifiers plus content verification are a possible implementation,
  not a chosen schema. A release number is a label, not the revision identity.
- WORK has its own ID and pinned base. Creating it from the current accepted
  release is a reasonable initial restriction. An existing WORK remains based on
  its original revision when a newer release appears; require explicit update and
  revalidation before release rather than pretending its base changed.
- A model commit freezes a complete candidate snapshot, even when the selected
  changes concern only a few files. It does not mean the implementation is verified
  and must not advance the accepted release implicitly. Blueprint generation can
  freeze a candidate automatically without making a project Git commit.
- Pull names an exact source revision. Importing another task's committed changes
  does not make them accepted or implemented. Record that dependency and invalidate
  affected earlier review/evidence. A common integration candidate should be
  explicit; never silently use whichever WORK last wrote a file.
- Release requires a validated candidate and the applicable implementation evidence.
  Preserve intermediate revisions and lineage; deleting/renaming the only WORK
  directory cannot be the history mechanism. Assign the next version only during
  successful publication, with local exclusive coordination and crash recovery.

### Bounded locking and merging

Local exclusive checkout may simplify the initial workflow. Key a lock by model
store identity and normalized logical source path, not an absolute path inside a
release directory. Track base revision/digest, WORK owner and recovery information.
Treat rename/delete and case-equivalent paths explicitly. Read-only permissions
are an aid, not enforcement: compare digests and refuse untracked writes or stale
bases before commit/release. Do not let canceled or crashed work leave unrecoverable
locks. Local locks do not establish distributed exclusivity across copied stores
or independent project Git clones; this boundary must be visible.

For v1, automatically combine disjoint file changes, preserve identical changes,
and stop on divergent changes to the same file (including delete/modify conflicts).
Three-way merge needs the common base, ours and theirs; a two-file diff/patch alone
cannot reliably distinguish concurrent edits from already incorporated changes.
Later text merging can improve convenience, but every combined model still needs
parsing and semantic validation: disjoint files may violate the same contract.
Never silently choose a side. The owner's combined merge/release command is a
possible convenience after these stages and failure behavior are defined.

### Surviving parent Git merges

Store immutable revisions at distinct identity-based paths. Independent additions
can then coexist after a normal project merge without merging database internals.
Human names such as SDL--V0.2.3 may be convenient projections, but two branches can
both allocate that name for different revisions. Retain both unique revisions and
report the ambiguous label; do not overwrite either or silently pick a latest head.
Record release claims/acceptance events independently, so semantic conflicts are
detectable even when Git reports no textual conflict. A changed revision with an
existing identity is corruption/conflict, never a normal history edit.

Mutable WORK data may conflict during project merges and requires explicit
resolution. Machine-local lock state cannot become authority in another clone.
The next SDPTool operation must validate imported history, competing release/head
claims and hashes before using a baseline. “Survives merge” means history is
preserved and ambiguity is exposed; it cannot mean every parallel design decision
is automatically compatible. User-facing discovery must avoid counting historical
snapshots as multiple live systems; source-root and cache/WORK handling need design.

### Proposed proof before implementation commitment

Exercise a project with no Git: freeze baseline and candidate, regenerate their
blueprint inputs, commit/pull across two WORK directories, reject stale writes,
and preserve canceled work history. Then use two project Git clones to create
independent revisions with the same release label, integrate and confirm both
histories survive while release ambiguity blocks promotion. Include interrupted
snapshot publication, changed immutable files, rename/delete, lock recovery and
cross-file semantic conflict. No such snapshot-store experiment has been run yet;
the earlier Git-only probe is not evidence for this design. BP2 remains planned.

## Lifecycle refinement — WORK / PROPOSAL / CANDIDATE / RELEASE

Owner discussion, 2026-10-02: consider browsable directories with UUID and content
SHA in YAML, parallel proposals integrated by one designated candidate/release
owner, and immutable releases. The owner also asks whether PROPOSAL can be removed.
The following is architectural advice, not an adopted storage schema.

Recommend retaining four workflow meanings while implementing only two storage
behaviors: mutable WORK and immutable snapshots. PROPOSAL freezes one contributor's
submission; CANDIDATE freezes the integrated target; RELEASE records acceptance of
that exact target. One proposal may form a candidate without merging. These roles
need not be four incompatible formats or require manual folder movement.
Removing PROPOSAL is viable if CANDIDATE also means contributor submission, but
then submission versus integrated acceptance target must still be distinguished.
A mutable WORK name alone cannot identify a reproducible blueprint input.

Suggested browsable convention (not a migration decision): WORK--<name>--<uuid>,
PROPOSAL--<name>--<uuid>, CANDIDATE--<name>--<uuid>, and
RELEASE--V<version>--<uuid>. Each contains model.yaml and a sources/ tree retaining
relative SDL/SDUI source paths. Directory names are labels; metadata owns identity.
This location must be designed alongside discovery to avoid treating all historical
copies as live systems. ZIP is optional transport, not authoritative storage.

A proposed manifest contains schema version, store ID, snapshot UUID, kind, parent
snapshot references, base release reference, language/profile inputs and a content
inventory. Define SHA-256 over a canonical sorted inventory of normalized relative
paths, file kinds and exact byte hashes. Exclude model.yaml itself and generated
caches from that content hash; separately digest the canonical immutable manifest
without its own digest field. Define path/case collision, symlink and external-source
rules before implementation (initially rejecting symlinks/unpinned external inputs
is simpler). UUID identifies a record; the digest detects changes, not authorship
or tamper-proof provenance. No mutable status update to a published manifest.

Promotion creates a new immutable record with lineage to its input; a RELEASE may
have the same source-content digest as its CANDIDATE while having a distinct UUID
and release metadata. Full copies are acceptable initially. Bind verification to
the candidate identity/content and relevant code revision or code artifact digest;
release confirms that exact content. Store later review and rejection decisions
separately, without rewriting the frozen proposal/candidate. A fix creates another
WORK and snapshot; an old candidate remains reproducible.

| Requested operation | Proposed meaning |
| --- | --- |
| WORK <- WORK | Freeze the source input, merge into the target WORK with a recorded base; preserve unresolved work and do not publish a candidate |
| CANDIDATE <- WORK | Freeze the WORK input, integrate in temporary WORK, validate and publish a new candidate UUID; old candidate remains unchanged |
| CANDIDATE <- CANDIDATE | Integrate both frozen inputs in WORK using their common ancestry, resolve and validate, then publish a third candidate with both parents |
| RELEASE <- anything | Forbidden; derive a new WORK/candidate and publish a new release |

Retain common ancestors. For the first supported multi-input integration, require
one unambiguous common base; stop for explicit handling of missing or multiple
merge bases rather than silently choosing one. Record input order and actual
resolution for multi-proposal integration. Same-file edits can be three-way merged
when non-overlapping; overlap must produce explicit conflicts. Delete/modify,
rename and added-path collisions need defined conservative handling. Textual success
is followed by SDL/SDUI parse, binding and contract checks plus applicable review;
it does not prove merged behavior. Existing review is stale for affected changed
content. Git merge-file demonstrates a standalone three-file merge primitive,
not a selected production dependency or a semantic model merger:
https://git-scm.com/docs/git-merge-file

The designated integrator is workflow authority, not distributed locking. Locally
serialize candidate/release publication and check the expected accepted head. After
Git combines independently produced histories, preserve both identities and reject
ambiguous release labels/competing accepted heads until explicitly reconciled.
Independent proposal paths should combine as file additions, not mutate a shared
central index; indexes may be rebuilt. UUID collision or changed frozen content
must fail validation. Accepted release is not inferred from highest version alone.

Blueprint comparison supports any two frozen proposals, candidates or releases.
Comparing WORK first freezes a snapshot; label ad-hoc mutable previews as such.
A comparison between unrelated proposals is meaningful as A-to-B difference but
is not automatically the implementation task from the accepted baseline. Pin the
chosen baseline explicitly and calculate impacted context from both models.

Next experiment should prove same-file disjoint merge, overlapping edit conflict,
candidate-to-candidate lineage, unchanged reviewed candidate promotion, corrupted
snapshot rejection and competing release claims after parent Git integration.
No implementation has been authorized or performed by this refinement.

## Semantic blueprint algorithm to design

1. Resolve and validate the complete baseline and target source graphs, including
   SDUI references. Pin includes/contracts/resources needed to reproduce each side.
2. Compare typed declarations, relations, contracts and relevant widget/layout or
   callback identities. Formatting and statement-order changes should not become
   semantic changes. Define stable identity and ambiguous rename handling; current
   statement-local fact ordinals are not adequate cross-revision identifiers.
3. Start at added/removed/changed facts and traverse relevant dependency directions
   in the union of baseline and target. Removed relations must not hide old consumers.
   Include enclosing system context, preserved boundary obligations and relevant
   indirect effects; arbitrary depth limits must expose their excluded frontier.
4. Generate before/after/overlay diagrams from those facts. Distinguish editable
   work area, affected neighbors, preserved contracts, additions/removals and unknown
   coverage. Do not claim complete actual-code impact from an incomplete model.
5. Freeze the assignment inputs and required checks. Link observed code mappings,
   structural checks, contract/scenario tests and reviewer judgment to exact revisions.
   SDL tags locate declared mappings; they do not prove behavior. AI review is useful
   evidence, not a guarantee of full semantic equivalence.

## Bounded pilot and acceptance

Select one small real XFMD change after reconciling an honest current baseline.
The already implemented dynamic-tab change can illustrate a retrospective delta,
but must not be reported as design approved before implementation. No XFMD edits
or external card transition is authorized by registering this proposal.

First demonstrate one revision-pinned before/after blueprint and an implementation
review. Test a candidate native model-history store against the no-export constraint; do not require a local PR UI, forks and release server
before validating that blueprint usefulness. The owner preference for separate
history is preserved; backend experimentation must not silently make project Git
permanently authoritative for the separate design lifecycle.

Acceptance must exercise moved/renamed sources, removed dependencies, unchanged
consumers, incomplete mappings, stale review, parallel proposals/worktrees, interrupted
promotion, canceled proposals, clean parent clone with no local cache, and recovery
of selected model branches/history from the portable artifact. No secret/local Git
configuration should be exported into the parent repository.

## Next decision

Select BP2 design work with this versioning input and KB-SDP-004's evidence contract.
Resolve native durable history without a separate export, revision identity/rename policy and the initial
XFMD change before implementation. No product implementation is selected by this card.

## Worklog

2026-10-01: Recorded owner proposal and bounded architectural analysis; linked BP1,
planned BP2 and existing Traceability proposal. State remains backlog.

2026-10-02 local date: Owner rejected bundle/export-dependent durability and asked
about SVN and pre-Git systems. Withdrew that recommendation; recorded SVN/Fossil
as unselected candidates and the native-storage/parallel-history acceptance gaps.
EVT-KB-SDP-000272; CardState remains backlog. No backend installed or implemented.

2026-10-02: Captured mandatory parallel parent-Git merge behavior and per-revision
blueprint reproduction. Corrected the SVN/Fossil direction: native store format
alone does not make independent histories Git-mergeable. Proposed project-Git-backed
model history for owner reconsideration; a five-case disposable Git storage probe
passed. EVT-KB-SDP-000273; no production command/backend or owner adoption claimed.

2026-10-02: Owner proposes standalone release snapshots and locked WORK directories.
Recorded Git-independent blueprint requirement, unique immutable revision identities,
local-lock limits, conservative merge rules and competing release-label detection.
EVT-KB-SDP-000274; proposal only, CardState remains backlog.

2026-10-02: Refined owner WORK/PROPOSAL/CANDIDATE/RELEASE proposal. Recommended
immutable submissions and integrated targets, identity-based browsable directories,
canonical content hashes, new-candidate merge semantics and designated integration
ownership. EVT-KB-SDP-000275; implementation remains unselected.

## Latest owner corrections — 2026-10-02

These take precedence over earlier recommendations in this card:

- Generate WORK blueprints directly without locking WORK or creating a persistent
  proposal/candidate. Mark them WORK / preliminary. Capture consistent bytes in
  memory or temporary storage and report source digests; detect edits during read
  and retry/fail rather than mix revisions. Such a preview is not approved assignment
  evidence. Reproduction requires retained inputs; hashes alone cannot restore them.
- WORK names have no UUID suffix. PROPOSAL/CANDIDATE names may use the last four
  alphanumeric characters of their UUID. Full UUID lives in tool-generated YAML.
  Short suffixes are display only: detect name collisions and regenerate the new
  UUID before publication or fail explicitly; never overwrite an existing snapshot.
- Keep optional PROPOSAL for now. Normal flow is WORK -> CANDIDATE -> RELEASE.
  A mandatory intermediate PROPOSAL is not required to freeze WORK.
- Final naming and standalone-tool ownership remain design questions, not
  implemented commands. A proposed common model namespace would expose work,
  proposal, candidate, merge, release and blueprint; it covers both SDL and SDUI.

Current command discussion (not executable today): prefer verb before object,
for example `create work`, `create blueprint`, `create candidate`, `create release`.
The owner is considering plain prepositions instead of flags. Recommended initial
spelling below uses `from`, `to` and `into` consistently; this spelling remains a
proposal, not approval of a parser or an additional alias set.

```text
sdptool model create work fixing-navigation
sdptool model create work fixing-navigation from release:0.2.3
sdptool model create blueprint from release:0.2.3 to work:fixing-navigation
sdptool model create candidate from work:fixing-navigation
sdptool model create work Combination from work:fix-something work:fix-navigation
sdptool model merge work:fix-navigation into work:fix-something
sdptool model create release from candidate:integrated --version 0.2.4
```

Omitted `from` on create work selects the current accepted release in the resolved
model store, records its exact identity and reports it. Never select by directory
mtime or silently choose between competing accepted heads. If no release exists,
report that an explicit initialization path is needed (bootstrap syntax remains
unselected), rather than inventing an implicit empty baseline.

`create work Combination from A B` means non-destructive three-way integration into
a new WORK, not directory concatenation or last-writer-wins copy. Source inputs
remain unchanged. Resolve each WORK consistently without locking/freezing it as a
published proposal; retain exact integration inputs/ancestry for later merges.
Require a known unambiguous common base for initial support, preserve source
provenance, and flag conflicts in the new WORK. Unresolved WORK cannot be promoted
to a valid candidate. Distinguish no-change, clean integration and unresolved
conflict in the result; a created folder alone is not merge success.

For existing WORK the owner also considers context-relative pull:

```text
# Run inside WORK--fix-something
sdptool model pull ../WORK--fix-navigation
```

This would integrate changes into the current WORK, preserving local edits and
recording incorporation so repeated pull does not reapply the same change.
Recommend `merge SOURCE into TARGET` as the explicit operation for now; whether
`pull SOURCE` is its context-relative convenience remains pending owner choice.
No network fetch or remote is implied. Reject inferred mutation when cwd does not
resolve uniquely to a mutable WORK; never infer a candidate/release as a writable
target. Record source expansion in provenance; the destination name need not be
renamed automatically. Plan/card scope changes remain explicit governance decisions.

A future standalone `sdl-model` command could share the implementation library
with `sdptool model`; naming and ownership are unselected. Current SDPTool imports
SDL/SDUI Go libraries directly and owns discovery/install/navigation orchestration;
it is not a generic dynamically discovered external-command plugin host.

2026-10-02: Corrected WORK preview and naming, registered Session0005 after owner
identified missing journal upkeep, and linked the bounded instruction correction
in MAINT-SDP-0014. EVT-KB-SDP-000276; feature remains backlog.

2026-10-02: Owner selects human-readable verb/object ordering and default latest
accepted release for WORK creation; discusses contextual pull and multi-source
create work. Recorded proposed consistent prepositions, explicit merge direction,
base/conflict rules and unresolved pull naming. Session0005 T004; event277.

## Typed creation targets and commit lineage — 2026-10-02, T005

Owner proposes using the same kind:name reference when creating and selecting an
artifact, YAML ledger history, explicit commit messages and reconstructable diffs.
Recommend the consistent creation grammar:

```text
sdptool model create work:Combination from work:fix-something work:fix-navigation
sdptool model create work:First --initial
sdptool model commit work:Combination --message "Combine navigation and document fixes"
```

These are proposed commands. Creation reserves a new name and fails if it exists;
source references resolve existing identities. Names with spaces need shell quotes.
The owner's invariant is a WORK anchored to a release or explicitly initial/empty;
combined WORK also records its exact source commit parents. Initial work is a root
with no release parent, not an implicitly selected missing release. The prior
latest-release default still applies when neither from nor initial is supplied.
Inputs based on different releases require explicit base reconciliation, not an
arbitrary choice of latest. Commit saves content state as well as a message, without
locking WORK or promoting it to candidate. Preliminary preview still needs no commit.

Recommend history as a directed acyclic graph, displayed as a tree with references
for shared ancestors. A commit contains unique identity, originating WORK identity,
parent commit IDs, timestamp/author/message and source snapshot digest. A merge
commit names both input heads and the resolved resulting snapshot. Copying input
ledgers together is set union by immutable event/commit ID with equality checks:
shared ancestors occur once; same ID with different content is an integrity error.
Do not interleave by timestamp and assume causality or flatten away original WORK
identity. Imported ancestry records contributions, not automatic acceptance of
all original code or behavior; resolved content and evidence establish the result.

The model.yaml ledger section can expose heads and references to durable immutable
records. Recommend one ordinary YAML record per commit/event as authoritative storage,
with all required content/history included in the portable model store. An inline
map keyed by IDs is an alternative, but a growing nested ledger in each manifest
creates duplication and parent-Git edit conflicts. Do not maintain two independently
editable copies. Exact paths/layout remain undecided; derived tree views are not
history authority and cannot depend on deleted WORK directories or external caches.
This model-history ledger is distinct from project-management lifecycle events.

The proposed delta algorithm has two different possible meanings: a cumulative
patch from release to each commit (no replay chain but repeated data), or incremental
patches from each preceding reconstructed commit (replay chain). Do not mix the two.
Both require retained bases, integrity checks and explicit deletion/rename semantics.
Recommendation: first save full logical snapshots; unchanged file versions can be
shared by content digest, or initially copied for simplicity. Generate diff on demand.
Delta compression is a later storage optimization, not a prerequisite for history.
Atomic commit publication must make content durable before advancing WORK head;
crash recovery must not leave a visible commit pointing at absent content.

Git comparison verified against official Git Internals documentation: commits point
to trees representing complete snapshots and parent commits; blobs hold file content.
Packfiles can delta-compress similar objects independently of logical parent history.
Thus Git's history is not defined as replaying a patch for every preceding commit.
Sources: https://git-scm.com/book/en/v2/Git-Internals-Git-Objects and
https://git-scm.com/book/en/v2/Git-Internals-Packfiles . This is reference analysis,
not a decision to use a nested Git repository or implement packfiles.

2026-10-02: Recorded typed creation targets, explicit initial WORK, commit content
and ledger DAG/union semantics, and snapshot-versus-delta alternatives. Session0005
T005; EVT-KB-SDP-000278. Storage schema and implementation remain unselected.

## Owner scope boundary — WORK-local undo, 2026-10-03, T006

This clarification supersedes recommendations for permanently retained full commit
content in a central model history store. The owner wants a deliberately limited
WORK-local history: numbered commits with messages and changed-file copies, nested
merge provenance, and recovery during model development. CANDIDATE/RELEASE retain
model content, commit messages and lineage but may discard WORK restoration payloads.
Do not require arbitrary historical WORK blueprint reconstruction after disposal.
Parent Git can retain checked-in history but cannot recover content never committed.
Surviving immutable candidates/releases remain valid blueprint inputs.

Owner proposed layout: a YAML file named after the enclosing artifact directory,
.commits and .merge always present logically, with numbered
commit directories; .merge archives each source/target YAML and its local histories,
including prior .merge. New root YAML refers to archived parents. Recommend matching
the directory stem exactly (WORK--Fix-Navigation/WORK--Fix-Navigation.yaml) instead
of the accidental single/double-hyphen discrepancy in the example. Full UUID remains
inside metadata; counters are local and display-oriented, not global identities.

Necessary recovery details for a bounded implementation:

- Changed-file commits store complete after-images at relative paths, plus explicit
  deleted-path records. A rename may initially be represented as delete plus add.
  Changes compare with the last saved WORK state, not repeatedly with its release.
- Reconstruction needs an available full starting state. Recommend .commits/#00000
  as a local baseline snapshot (empty for initial WORK), followed by numbered
  changed-file records; this avoids depending on an external release directory
  continuing to exist. History/metadata/cache directories are excluded from source
  snapshots. Commit identity is WORK UUID plus counter or another unique record ID.
- Restore means reconstruct the whole model at a chosen local commit. Preview the
  change and preserve current uncommitted work before applying. Restoration should
  append a new recovery record, not erase existing lineage. Selective cross-branch
  undo/cherry-pick and automatic inverse merge are explicitly outside initial scope.
- A merge archives both inputs as they stood before integration, resolving unsaved
  edits by an explicit saved checkpoint or stopping before mutation. Build archives
  in staging outside the destination traversal; never recursively copy the new
  destination into itself. Use identity-qualified archive entries if folder names
  collide; path names alone do not distinguish independent WORKs.
- Keep complete pre-merge and resolved post-merge checkpoints as bounded recovery
  anchors. Whole-WORK restoration before/after merge is possible even with overlapping
  files when those states exist. Restoring only one source's contribution while
  retaining other changes is not promised. Source-branch commits are provenance,
  not automatically restore points on the integrated WORK's current line.
- Recursive archive copying is a workable first representation but can duplicate
  common ancestry and grow rapidly. Preserve existing archives and references;
  reuse identical frozen archives by full identity/digest within the WORK when
  practical. No central history database is required. Missing payload is reported
  as unavailable restoration, not reconstructed from commit messages.
- Promotion copies current validated sources and a metadata-only lineage closure,
  including source/target identities, messages and parent relationships. It excludes
  .commits file payloads and archived source payloads from .merge. Do not leave
  provenance depending on soon-deleted WORK paths. Mark retained historic entries
  as metadata-only where restoration content has been dropped. Candidate/release
  immutability remains unchanged.

The root YAML may be replaced as the current WORK head description; archived YAML
records remain immutable and linked by identity. WORK file layout/history schema
is proposed, not implemented. No global VCS, packfiles, distributed locks, selective
rollback or permanent every-commit source archive is required by this scope.

2026-10-03: Recorded owner reduction to local undo history with changed-file copies,
recursive merge provenance and metadata-only promotion. Explained deletion/base
requirements and whole-state versus selective merge rollback. Session0005 T006;
EVT-KB-SDP-000279. BP2 remains planned; no implementation started.
