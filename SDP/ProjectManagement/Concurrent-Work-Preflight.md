# Concurrent-work preflight

Adopted local procedure, owner Session0011 T002, 2026-10-09. This is a manual
execution prerequisite, not an implemented distributed lock or MCP routine engine.
Apply before starting/resuming product or SDL/SDUI implementation, before expanding
its write scope, and before integration. Reference the result from the active Session.

## Inputs and inspection

1. Recover current card, Session, plan, worktree and exact Git HEAD/status. State
   the intended system, planned write paths, shared contracts and management files.
2. Enumerate known local worktrees with `git worktree list --porcelain`. Read active
   cards and active/paused Sessions in relevant worktrees. An old checkout's active
   Session is a historical projection, not proof that an agent is currently working.
   Find the freshest governing record and its source; mark unresolved ownership.
3. Inspect each relevant branch's committed changes since the common base, plus
   staged, unstaged and untracked changes. A clean worktree can still contain
   unintegrated work. Include shared dispatchers, response schemas, discovery,
   presenters, dependency files, SDL models, ledgers and ID namespaces.
4. Compare intended write sets and behavioral contracts. Check newly allocated
   card/Session/plan/event IDs across known worktrees. Matching IDs with different
   subjects are a conflict even if Git can merge their files without a text conflict.
5. Where useful, run `git merge-tree --write-tree LEFT RIGHT` on exact committed
   heads. It creates Git objects but does not merge branches or alter worktrees.
   Retain exit status and conflicts. A clean forecast excludes uncommitted work and
   says nothing about semantic compatibility or edits made after inspection.
6. Record source paths, heads/status, affected owners/work records, exact overlap,
   verdict, permitted scope and restart conditions. Recheck immediately before the
   first write and whenever a participating head, working copy, scope or owner changes.

## Verdict and action

| Verdict | Meaning | Permitted next action |
| --- | --- | --- |
| Green | No unresolved relevant overlaps/identity collisions, or an explicit ownership/integration agreement covers them | Execute only the reviewed scope; retain the agreement and exact baseline |
| Yellow | Relevant owner, scope, baseline or freshness cannot be established | Continue read-only analysis and isolated coordination records; resolve unknowns before affected product edits |
| Red | Concrete overlapping edits, competing contracts or conflicting identities lack a disposition | Stop affected product implementation; resolve ownership/integration and rerun preflight |

Shared-system work is not automatically a conflict: independent packages can proceed
when their integration boundary is explicit. Conversely, separate branches do not
make a shared dispatcher or public contract independent. A green result is scoped,
time-bound evidence, never a guarantee against all future or remote conflicts.

Do not use an arbitrary observation timeout to infer another agent has finished.
Do not reset, stash, renumber, merge or rewrite another workstream to manufacture
a green result. Coordinate using available owner/agent channels within their actual
authorization. Agree an integrator, shared-file ownership and sequencing; retain
append-only history when reconciling IDs. Merge/release authority remains separate.
Unknown remote work is recorded as outside observation coverage, not absent.

## Durable result

Use an existing Session evidence directory and journal; no new mandatory database,
card or Sprint. Record stable card identity together with subject/path/revision when
an ID collision is discovered. Keep the Session roadmap and queue prerequisites
aligned. A blocked implementation check does not pause the Session automatically.
The owner may explicitly pause it; otherwise retain the gate and useful coordination
work. Link this prerequisite from the applicable assignment/plan before execution.

Distribution into templates/installed project skills is separate scoped work; this
local adoption does not claim every project already runs the procedure.
