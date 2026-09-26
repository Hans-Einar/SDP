# Short title

| Field | Value |
| --- | --- |
| id | KB-PROJECT-NNN |
| project | PROJECT |
| type | Idea |
| CardState | backlog |
| created | REPLACE_WITH_RFC3339_TIMESTAMP |
| source | REPLACE_WITH_CONVERSATION_OR_ISSUE_REFERENCE |
| next_review | YYYY-MM-DD |

This is a template, not a registered card. Create the file and created event
together. For Ref cards, add `primary` with `KB-OTHER-NNN` and a Markdown link
to that primary card. Optional `tags`: process, tooling. Lifecycle status follows
the directory and ledger. Write all prose and metadata values in English.

## Need and source

What problem should be solved? Distinguish the owner's request/decision from
recommendations and observations.

## Proposal, scope and open questions

What is known, what remains undecided, and which projects are affected?

## Next action and completion criteria

What must be decided or delivered? In active: owner, bounded scope and checks.
In onHold: reason, restart condition and next review.

## Outcome and references

Update before closing. Link decisions, plans/slices, evidence or successors.
Do not claim product implementation merely because a proposal was evaluated.

## Worklog and revisions

Update while active as well as at transitions. Add a row for each substantive
review; retain older entries and correct them through additions. Keep current
scope and next steps above. Do not insert the commit's own hash before committing;
the event ID connects the record to its introducing Git revision.
See [history and diff](History.md).

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| Actual registration time | Actor; EVT-KB-PROJECT-NNNNNN | Distinguish outcomes, proposals and owner decisions | Link and concrete next step |

## Lineage when merging or splitting

Use this section only for merge/split. Name the operationId and all source/target
IDs with clickable links. Optional metadata `sources` and `superseded_by` lists
IDs; the ledger's lineage field describes the complete operation.

| Source ID/link | Target ID | Transferred scope | Replaced statement and reason | Remaining work and home |
| --- | --- | --- | --- | --- |
| One row per relevant source/target | New target | What is carried forward | Decision, date and authority | In the source or a named target |

See [Lineage](Lineage.md). Never close a source while remaining work has no home.

## Queue (when CardState is queued)

Why this is the next bounded item, who selected it, relevant predecessors,
prerequisites and the first deliverable. Queue selection is not approval to start.

## Gate review (when CardState is gate-review)

Link the concrete reviewable result, checks, remaining choices and the owner's
requested decision. Acceptance closes the card; requested changes return it to
in-progress with a new worklog/review event. Do not infer approval from silence.

## Optional grouping

Add SprintId (SPR-SDP-NNNN), ScrumId (SCRUM-SDP-NNNN) and Systems when useful.
Omit unassigned fields. The metadata owns current grouping; record every assignment
change in the management ledger and keep the Sprint member list consistent.
Tags do not move cards or imply implementation. See
[project management](../ProjectManagement/README.md).
