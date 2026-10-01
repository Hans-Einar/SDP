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
| tags | blueprint, model-history, baseline, design-review, Git |

## Owner outcome

The XFMD agent changed implementation before updating its SDL/SDUI design. The
owner proposes independent local version control for design sources: create a
model branch, design a change, review its proposed difference, generate a blueprint
showing the affected system, implement against that blueprint, verify code against
the target model, then integrate the verified design. Avoid a directory copy for
every release and avoid requiring a separate GitHub model repository. Preserve
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
| Project Git plus path-scoped model revisions | Existing portable history; straightforward exact model/code links | Model and implementation commits share one repository, contrary to the owner's preferred separation |
| Local model Git, isolated proposal worktrees, accepted source export into parent | Independent model branches and ordinary readable source files in the parent | History/review recovery needs an explicit durable transport and reconciliation contract |
| Normal submodule | Independent Git history with parent pin | Does not satisfy local-only history automatically following the parent's push/clone |

Recommended direction to prototype: the second option. Keep the local Git database
outside the visible source tree; do not insert a .git directory/file into the
tracked SDP/SDL tree or dual-control the same working files with two active indexes.
Use isolated model proposal worktrees and publish an accepted snapshot into the
project only through a checked operation. Default discovery sees the accepted source;
proposal preview must explicitly select its isolated workspace/revision.

If full local model history must survive an ordinary parent clone, export a complete,
verified bundle of required refs/objects plus portable review/baseline records as
project-owned data. The parent stores ordinary accepted SDL/SDUI files as well.
This is synchronization to implement, not something Git does implicitly. The bundle
is opaque to parent text merges and may grow; evaluate deterministic export,
round-trip restoration, incremental/full choices and concurrent update conflicts
before selecting its on-disk format or making it an installer default. Merely
reconstructing a new repository from current source loses proposal history/identities.

Do not store runtime-only history in a nonportable directory and then claim it is
backed up by the parent. Separate Git worktrees need isolated work state, correct
Git common-directory discovery and an explicit synchronization/locking policy.

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

There is no atomic Git commit spanning the two repositories. Define a recoverable
promotion/export transaction and code/model receipt; concurrent proposals require
base validation/rebase, conflict resolution and renewed affected review. Never
silently overwrite changed accepted source or consume an unrelated worktree state.

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
review using existing Git history. Then test the independent local-history backend
against that same contract; do not require a local PR UI, forks and release server
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
Resolve durable-history transport, revision identity/rename policy and the initial
XFMD change before implementation. No product implementation is selected by this card.

## Worklog

2026-10-01: Recorded owner proposal and bounded architectural analysis; linked BP1,
planned BP2 and existing Traceability proposal. State remains backlog.
