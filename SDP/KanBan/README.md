# SDP — shared KanBan

## Card index

| ID | Type | Status | Document |
| --- | --- | --- | --- |
| KB-SDL-001 | Proposal | backlog | [Stakeholders, actors, user stories and readable SDL](backlog/%23001--SDL--Proposal--Requirements-narrative.md) |
| KB-SDL-002 | Proposal | backlog | [Weak links and visual paths through nodes](backlog/%23002--SDL--Proposal--Links-through.md) |
| KB-SDL-003 | Ref | completed | [SDP tools as a consumer of SDL](completed/%23003--SDL--Ref--SDP--017--sdptool.md) |
| KB-SDL-004 | Change | superseded | [KB-SDL-004 — Organize language development by process phase](superseded/%23004--SDL--Change--Language-source-organization.md) |
| KB-SDL-005 | Change | backlog | [SDL System and explicit source-set contract](backlog/%23005--SDL--Change--System-and-source-sets.md) |
| KB-SDP-001 | Proposal | completed | [Project structure, Template and studies per phase](completed/%23001--Proposal--Project-structure.md) |
| KB-SDP-002 | Proposal | superseded | [sdptool: project discovery, implementation plan and viewer](superseded/%23002--Proposal--sdptool.md) |
| KB-SDP-003 | Idea | backlog | [KanBan graph with a time axis and progressive detail](backlog/%23003--Idea--KanBan-graph.md) |
| KB-SDP-004 | Proposal | backlog | [Traceability between SDL design, slices, code and evidence](backlog/%23004--Proposal--Design-traceability.md) |
| KB-SDP-005 | Ref | completed | [SDL: links through in SDP viewpoints](completed/%23005--Ref--SDL--002--Links-through.md) |
| KB-SDP-006 | Ref | completed | [SDL: requirements models in the SDP process](completed/%23006--Ref--SDL--001--Requirements-narrative.md) |
| KB-SDP-007 | Change | completed | [K1-M1: establish KanBan and preserve conversation proposals](completed/%23007--Change--KanBan-foundation.md) |
| KB-SDP-008 | Change | completed | [K2-M1: visible metadata in KanBan cards](completed/%23008--Change--Visible-card-metadata.md) |
| KB-SDP-009 | Change | completed | [K3-M1: traceable card merging and splitting](completed/%23009--Change--Card-merge-and-split.md) |
| KB-SDP-010 | Proposal | active | [Consolidate documentation and separate active contracts from history](active/%23010--Proposal--Document-consolidation.md) |
| KB-SDP-011 | Bug | backlog | [Resolve older Traceability IDs against the Toolkit contract](backlog/%23011--Bug--Traceability-id-conformance.md) |
| KB-SDP-012 | Change | completed | [K4-M1: living cards, revision logs and Git diff](completed/%23012--Change--Card-history.md) |
| KB-SDP-013 | Change | completed | [English documentation throughout SDP, SDL and SDUI](completed/%23013--Change--English-documentation.md) |
| KB-SDP-014 | Proposal | backlog | [Define the KanBan version contract and reusable distribution](backlog/%23014--Proposal--KanBan-version-contract-and-distribution.md) |
| KB-SDP-015 | Change | completed | [CardState and scoped shell commands](completed/%23015--Change--CardState-and-shell-cli.md) |
| KB-SDP-016 | Proposal | superseded | [Define SDP discovery, version and viewer capability metadata](superseded/%23016--Proposal--SDP-discovery-and-viewer-capabilities.md) |
| KB-SDP-017 | Proposal | active | [sdptool and integrated project navigation](active/%23017--Proposal--sdptool-and-project-navigation.md) |
| KB-SDP-018 | Study | backlog | [Audit Toolkit responsibilities, organization and distribution](backlog/%23018--Study--Toolkit-audit-and-organization.md) |
| KB-SDP-019 | Ref | completed | [KB-SDP-019 — Shared process adoption for SDL and SDUI](completed/%23019--Ref--SDL--004--Language-source-organization.md) |
| KB-SDP-020 | Change | backlog | [Organize shared system design sources by phase](backlog/%23020--Change--Shared-design-source-organization.md) |
| KB-SDUI-001 | Ref | completed | [SDUI as a subproject and library for SDP tools](completed/%23001--SDUI--Ref--SDP--017--sdptool.md) |
| KB-SDUI-002 | Ref | completed | [KB-SDUI-002 — Split the shared design into phase-owned sources](completed/%23002--SDUI--Ref--SDL--004--Language-source-organization.md) |
| KB-SDP-021 | Change | completed | [Saved-file design preview](completed/%23021--Change--SDPTool-saved-design-preview.md) |
| KB-SDP-022 | Change | completed | [Project recognition and delegation contract](completed/%23022--Change--SDPTool-discovery-contract.md) |
| KB-SDP-023 | Change | completed | [Project resolver and configured viewer bridge](completed/%23023--Change--SDPTool-project-viewer-bridge.md) |
| KB-SDP-024 | Change | completed | [Model-derived navigation and selected generation](completed/%23024--Change--SDPTool-model-navigation.md) |
| KB-SDP-025 | Change | completed | [KanBan and SDUI navigation services](completed/%23025--Change--SDPTool-KanBan-and-SDUI-inventory.md) |
| KB-SDP-026 | CodeReview | completed | [Consumer contract and end-to-end delivery review](completed/%23026--CodeReview--SDPTool-consumer-contract-review.md) |
| KB-SDP-027 | Study | completed | [Scrum for skills consolidation and project activation](completed/%23027--Study--Skills-review-and-project-activation.md) |
| KB-SDP-028 | Study | backlog | [Scrum for versioned SDP installation and upgrades](backlog/%23028--Study--Installer-upgrade-and-versioned-layout.md) |

## Purpose and authority

Capture ideas, questions and requested changes from conversations as cards, even
when details remain unresolved. Registration does not adopt a language rule or
authorize implementation. Distinguish owner decisions, agent recommendations and
open choices. Do not re-register delivered G phases as new requests.

Use one project board at SDP/KanBan. SDL and SDUI cards now live here with
stable namespaces; the [management ledger](../ProjectManagement/Ledger.ndjson)
owns all new transitions. [The local profile](../ProjectManagement/README.md)
adds optional Sprint/Scrum grouping and relates management to system Traceability.
Historical boards are frozen in its import manifest. This local migration does
not change XFMD's pinned workflow or distributed installer templates.

**All card and workflow documentation must be English**, including metadata
values, worklogs, queue explanations and new event descriptions. Preserve old
append-only event bytes and clearly identified verbatim source quotations.

## Types and tags

The primary `type` appears in the filename. Optional `tags` may identify topics
such as tooling, language, process or visualization. Type is not priority or status.

| Type | Purpose |
| --- | --- |
| Idea | A possibility to evaluate; no finished solution required |
| Proposal | A concrete proposal with purpose, open choices and acceptance criteria |
| Question | A clarification requiring a recorded answer |
| Study | A bounded investigation with questions and expected outcomes |
| Change | An authorized change with known scope |
| Bug | An observed deviation with expected behavior and reproduction |
| Decision | Decision, decision-maker, date, rationale and consequences; the type alone is not approval |
| Ref | Local impact and a link to one primary card in another project |
| CodeReview | Review work, with a linked review record and findings |
| Refactor | Behavior-preserving code/design improvement, with explicit evidence |

Use UserStory and other model concepts in SDL only when its profile supports
them. KanBan types are not new SDL keywords. Add types only for a concrete need.

## Identity and document format

Example filename: `#003--Idea--KanBan-graph.md`; reference card:
`#001--Ref--SDP--002--sdptool.md`. Allocate increasing numbers per project across
all statuses/types, with at least three digits. Never reuse closed/deleted IDs.
The stable ID is `KB-<PROJECT>-<number>`; filenames and locations may change.
Check every directory and ledger before allocation. Resolve parallel ID/event
collisions before committing; never overwrite another registration.

Use the [card template](Card-template.md). Metadata is a visible Markdown table
immediately below the title, with `Field` and `Value` columns. Keep field names
`id`, `project`, `type`, `created`, `source` and optional `next_review`, `primary`,
`tags`, `SprintId`, `ScrumId`, `Systems`, plus mandatory `CardState`. Do not duplicate metadata in YAML frontmatter. The ledger remains
JSON and owns event history.

Each primary card owns one coherent need. A Ref has its own ID/status, `primary`
containing the primary card ID, a clickable link and local impact. Link directly
to the primary card, never through a Ref chain. Prefer the nearest responsible
project, but do not duplicate a need merely to obtain perfect placement.
Completing a Ref does not complete its primary card. Link implementation slices
back to the card where needed. Include card IDs in related issues/PRs; GitHub
status does not automatically change local KanBan status.

## Lifecycle directories

| Directory | Meaning and requirements |
| --- | --- |
| backlog | Registered, awaiting prioritization or clarification; specify the next review |
| active | Selected, bounded work; specify scope, owner and completion criteria |
| onHold | Blocked or deferred; specify reason, restart condition and review date |
| completed | Agreed outcome achieved, with linked decision/delivery/evidence |
| canceled | Previously relevant or selected work deliberately stopped; explain why |
| superseded | Fully replaced, merged or split; name and link every successor |
| irrelevant | Reviewed as outside scope or no longer relevant; explain why |

Names are case-sensitive: use `superseded`, not `superseeded`. An Idea, Question
or Study can complete when its agreed evaluation is delivered; that does not
mean the proposed product feature is implemented. Implementation cards require
actual implementation and agreed verification, not just a plan. Link authorized
follow-up work and obtain implementation status from Traceability.

Any status may be reopened with a reason. Do not delete cards to empty backlog.
Move duplicates to superseded with a primary-card link. Follow [Lineage](Lineage.md)
for merge/split: new targets, preserved sources and explicit remaining work.
Only fully replaced sources close. Review related cards before selecting work.

## Work rhythm and scope

Use active cards as working documents, with current scope, next action and a
worklog. [History and diff](History.md) connects log entries, ledger events and
Git revisions, including reviews without a status change. Git stores content;
KanBan stores processing history; Traceability stores implementation evidence.

At startup, read the board and affected Refs, select bounded work and activate
it with explicit scope. Record scope-changing discoveries before changing focus.
At milestones, update outcomes, ledger and references. Review backlog/onHold
before the next phase and at the agreed review date. Decide the next work,
defer with a new date, cancel, supersede or mark irrelevant. The initial
2026-09-30 review date is not a delivery deadline. Age alone does not justify
deleting or rejecting an idea. Aim for a small, well-understood backlog.

## Moving cards — manual workflow

1. Read the card and its latest event; note ID, state and old path.
2. Record reason/outcome, successors if any, plan and evidence in the card.
3. Move the whole file, preserving its ID and normally its filename.
4. Append an event with actual time, actor and old/new status and paths.
5. Update the index and incoming Markdown links across registered boards. Search
   both stable ID and old filename; encode `#` as `%23` in URL targets. Do not
   rewrite historical ledger paths.
6. Verify links, replay and physical placement. Commit the document, ledger and
   indexes together. Investigate mismatches; do not infer state from mtime.

The durable reference is the card ID; its Markdown link is its current address.
Moves within one board preserve relative depth. For this explicit consolidation, retain the original namespace and prefix it
in filenames where needed. Future arbitrary cross-repository transfer still needs
an explicit contract; this migration does not imply distributed transactions.

## Ledger contract — payload 0.1 and 0.2

[board.json](board.json) points to the common management ledger. Use the
[SDP envelope](../../Toolkit/schemas/ledger-event.schema.json), and payload 0.2
for new x-kanban events. Preserve historical 0.1 events. New management record
events are defined separately in the [management contract](../ProjectManagement/README.md).
The rules below apply to the KanBan event subset of that one ledger.

- `eventId`: unique, increasing `EVT-KB-<PROJECT>-<number>` within the board.
- `eventType`: `x-kanban:created`, `x-kanban:moved` or `x-kanban:reviewed`.
- `subjectId`: stable card ID. `occurredAt`: actual RFC3339 timestamp with timezone
  (UTC here). `actor`: recorder. `commit` may be null; the introducing Git commit
  records the event. Never rewrite an event to add its own commit hash later.
- Payload: `schemaVersion`, `projectId`, `previousEventId` for the same card,
  `from`, `to`, `fromPath`, `toPath`, `reason`, `links`. Version 0.2 also permits
  typed `lineage`; see [its completeness rules](Lineage.md).
- Paths are literal UTF-8 board-relative paths, with no `..` or absolute prefix,
  not URLs. `links` contains stable card/slice IDs or board-relative evidence paths.
- created: first event; previous event/status/path are null. Start in backlog or
  active. Do not invent an earlier history.
- moved: predecessor matches the latest event; origin matches previous status/path;
  target represents a real move or rename.
- reviewed: same status/path, with outcome/reason and any new review date in the
  card. Correct state through a new explained event, not by rewriting history.
- completed/superseded transitions require an outcome/successor reference in
  `links`. Superseded must name every successor, each an existing primary card.

Replay follows line order and the predecessor chain, not timestamps alone.
Resolve new, unintegrated event-ID collisions before merging. The last event
must match exactly one physical card at the recorded location. Schema validation
alone does not prove chain rules, references, status transitions or placement.
No interactive graph is delivered by this contract. See the
[K1 evidence](completed/%23007--Change--KanBan-foundation.md) for its original checks.

## CardState — current work state

The card's visible metadata table owns **one** `CardState`. No gate file, queue
folder, duplicate YAML field or state column in this index. Directory/ledger
still own lifecycle placement/history; CardState refines work within a stage.

| CardState | Directory | Meaning |
| --- | --- | --- |
| backlog | backlog | Registered, not selected next |
| queued | backlog | Proposed next bounded work; add a Queue section with reason, predecessor, prerequisites and next step |
| ready | active | Selected, currently awaiting work |
| in-progress | active | An agent/person is working on the bounded scope |
| gate-review | active | Concrete result awaits owner review; explain the decision and evidence |
| onHold | onHold | Deferred/blocked, with restart condition and review date |
| completed | completed | Agreed outcome delivered |
| canceled | canceled | Deliberately stopped |
| superseded | superseded | Fully replaced, with successors |
| irrelevant | irrelevant | Reviewed as outside scope/no longer relevant |

Set in-progress when starting; return to ready only with a recorded pause/handoff.
Use gate-review when review is required, not as a synonym for unresolved future
implementation. Owner acceptance leads to completed; requested changes return to
in-progress. Do not infer acceptance from elapsed time. Existing explicitly
accepted or objectively authorized deliveries may close with evidence.

Prefer at most one queued primary card per board. A Queue section explains why
it is next and its prerequisites; the predecessor's Next selection section records
why it selected that card. Queued is prioritization, not automatic authorization.
Do not keep a completed predecessor active merely because it has a successor.

For state-only changes, update metadata/worklog and append a reviewed event whose
English reason records old → new state. For moves, update folder/state/index/links
and append moved. No new ledger schema or duplicate authoritative state is needed.
Git preserves exact document revisions; ledger events explain transitions.

For a quick text search: `rg '^\| CardState \|' backlog active`. K5's
[CLI](../../Toolkit/scripts/cli/README.md) lists grouped states without mutation.
Old cards without CardState are ignored by the lister and should gain metadata
when maintained; an empty result is an error. New maintained cards require it.
