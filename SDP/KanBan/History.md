# Worklogs, revisions and diffs for KanBan cards

A card is a working document while active. Keep one file with a stable ID through
its lifecycle. Git owns content revisions; the KanBan ledger owns processing
events; Traceability owns implementation and verification. This is K4's manual
workflow, not a new version-control system or an implemented sdptool command.

## While a card is active

The opening sections describe current needs, agreed scope, owner, next action and
completion criteria. Add a concise worklog with time, actor, event ID, result and
evidence. Update it for important findings, scope changes, decisions, checks and
milestones even when status does not change. It is not a chat transcript.

- Distinguish observations, agent proposals and owner decisions; name the decider.
- Use `x-kanban:reviewed` for substantive processing without a move, preserving
  path/status and linking the outcome. Group related small edits where sensible;
  spelling and link repairs do not each need an event.
- Update current text with new knowledge and explain replacements in a new log
  entry. Do not remove earlier rationale or rewrite ledger history.
- State remaining work at each milestone. A finished subtask does not complete
  the entire card. Preserve the original request and separate delivered/planned work.
- Commit card, related documents, index and ledger together at milestones. Include
  both milestone ID and card ID in the commit. Between milestones, local drafts
  are not durable Git snapshots; Git cannot show revisions never committed.

The worklog is a readable summary. Correct earlier entries through additions.
Do not create snapshot copies of cards in a revision directory. Do not insert a
commit's own hash into its contents. `commit: null` remains valid in the ledger:
the Git commit introducing its unique event ID identifies the revision. Later
log entries may refer to already existing commits.

## Maintenance and active cards

`SDP/Maintenance/<phase>/` may own a coherent phase plan, inventory and larger
verification evidence. The card links there and summarizes progress, next steps
and remaining work. A reader should understand its situation without searching
random notes. Do not duplicate whole plans or raw test logs inside the card.
Preserve cards and evidence on closure. Maintenance is a delivery location,
not another mandatory SDP phase or abstraction level.

## Inspect history with Git today

Run from the repository root. This concrete example uses KB-SDP-001:

```sh
git log --follow --date=iso-strict --format='%h %ad %s' -- 'SDP/KanBan/completed/#001--Proposal--Project-structure.md'
git log --follow -p -- 'SDP/KanBan/completed/#001--Proposal--Project-structure.md'
git diff -- 'SDP/KanBan/completed/#001--Proposal--Project-structure.md'
git diff --cached -- 'SDP/KanBan/completed/#001--Proposal--Project-structure.md'
```

The last two show unstaged and staged changes. For historical revisions, find
the commit and the path at that revision; it need not match today's path:

```sh
git show 'bb3728c:SDP/Agents/KanBan/backlog/#001--Proposal--Project-structure.md'
git diff 'bb3728c:SDP/Agents/KanBan/backlog/#001--Proposal--Project-structure.md' 'f722dc2:SDP/Agents/KanBan/active/#001--Proposal--Project-structure.md'
git log --format='%h %s' -G '"eventId"[[:space:]]*:[[:space:]]*"EVT-KB-SDP-000015"' -- SDP/ProjectManagement/Ledger.ndjson
```

`--follow` infers renames for one file. It is useful for ordinary moves, but is
not authoritative for card identity and does not follow semantic merge/split.
When uncertain, obtain historical paths from the ledger and compare exact
`commit:path` blobs. If the board itself moved, consult the repository migration map.

## Multiple sources, repositories and incomplete history

[Lineage](Lineage.md) governs merge/split. New targets and preserved sources have
separate revisions. A lifecycle view must show their branches, rather than pretend
there was one file. Partial transfer leaves the source's worklog active.

A future tool must resolve the project's Git repository through project registration,
then the stable card ID, ledger path and revision. Pair repository identity with
commit hashes; a hash alone is not a cross-repository address. Refs link the
primary card's history while retaining their own local processing history.

Label local drafts, unavailable repositories, shallow clones, missing objects and
broken history explicitly. Never reconstruct content from mtime or present a
missing revision as empty. A Markdown-only export includes the worklog, not full
Git history. Rebase/squash may change commit identities; preserve the owner's
phase/milestone commits. Card/event IDs remain stable.

Integrated history/diff remains in [sdptool](active/%23017--Proposal--sdptool-and-project-navigation.md)
and the [timeline idea](backlog/%23003--Idea--KanBan-graph.md).

The PM1 [import map](../ProjectManagement/History/import.json) preserves the exact
old board roots and event bytes. Historical commit:path examples above deliberately
retain their original addresses. Current readers use the common management ledger.
