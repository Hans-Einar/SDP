# Stakeholders, actors, user stories and readable SDL

| Field | Value |
| --- | --- |
| id | KB-SDL-001 |
| CardState | backlog |
| project | SDL |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registered from the owner conversation on 2026-09-23. The timestamp records registration, not a reconstructed time for earlier discussions. The directory and ledger record lifecycle status.

## Owner proposal

Start requirements work directly in SDL: `01--Actors.design`, followed by `02--UserStories.design` and `03--UseCases.design`. List stakeholders first. Actors may be linked to stakeholders, but this is optional. [KB-SDP-001](../../../../../SDP/Agents/KanBan/completed/%23001--Proposal--Project-structure.md) owns process placement and studies.

A Stakeholder has an interest in the system. An Actor is a role or external entity that interacts with it. Distinguish human users from machines/software systems: P1000 is a machine actor, not a User. Define explicit actor types and consistent visual symbols; do not infer type from names or draw every actor as a person.

A UserStory has identity, free text and links to actors and, later, use cases. Its text is information, not automatically executable logic or inferred requirements. Proposed English spelling: `UserStory`; the owner's `userstorry` was a conversational sketch.

## Syntax sketches — not implemented SDL

```text
Stakeholder stakeholder1
Actor actor1 links to stakeholder1
UserStory story1 links to actor1
story1.story = "As an operator, I want ..."

UserStory story2 links to actor1
.story = "As an operator, I want ..."
.comment = "Description to display in the documentation."
```

The owner wants explicit object properties and consideration of a leading-dot shorthand. Recommendation for review: bind the shorthand only to the nearest preceding declaration in the same block, not any earlier reference use. Define behavior at imports, block boundaries, blank lines, mixed statements and end of file.

`comment` should provide documentable descriptions, distinct from source-only comments. Keyword versus property remains undecided. Compare multiline text (quotes/braces/brackets) and the owner's `#include 01--Actors.design` sketch with the actual grammar before choosing. This card does not introduce an include preprocessor or interpretation of narrative text. [KB-SDL-002](%23002--Proposal--Links-through.md) owns weak links and visual paths.

## Next work and acceptance

Inspect the active parser/profile before proposing grammar. Resolve names, casing, actor types, multiline/import rules and shorthand binding. Then deliver a versioned profile, AST/source positions, validation, negative examples and document projections. A readable requirements package must preserve text exactly and show stakeholder–actor–story–use-case traceability without treating a descriptive link as requirements satisfaction.

[SDL overview](../../../../README.md)

## Backlog review — 2026-09-25

Retained as a bounded requirements-language proposal. Coordinate basic weak-link semantics with KB-SDL-002, but narrative/profile delivery need not absorb visual traversal through architecture nodes.

Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDL-000010. CardState remains backlog.
