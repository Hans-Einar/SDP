# Session 0005 — Blueprint model history

## Session roadmap

Latest recorded turn: T006. Earlier discussion is reconstructed below, not a
complete transcript. Sequence only; synthetic dates do not measure elapsed time.

```mermaid
gantt
    title Blueprint model history - sequence only
    dateFormat YYYY-MM-DD
    section Route
    DONE S1 Constraints :done,s1,2000-01-01,1d
    ACTIVE S2 Lifecycle :active,s2,after s1,1d
    NEXT S3 Design and experiment plan :s3,after s2,1d
    PLANNED S4 Implement and verify pilot :s4,after s3,1d
```

| State | Step | Work | Authority / evidence |
| --- | --- | --- | --- |
| completed | S1 | Recover standalone history requirements | KB048; no Git/export prerequisite |
| on-going | S2 | Refine snapshots, merging, blueprints and command ownership | Owner discussion; KB048 recommendations and corrections |
| next | S3 | Select/revise bounded design and experiment plan | BP2 remains planned; implementation not selected |
| planned | S4 | Implement and verify an agreed pilot | Requires selected scope/plan; no backend delivered |

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0005 |
| Status | active |
| Primary card | KB-SDP-048 |
| Snapshot date | 2026-10-03 |
| Current step | S2 |
| Proposed next step | S3 after settling the minimal lifecycle and tool boundary |
| Execution authority | Owner authorizes discussion and immediate Session-instruction correction; model-store implementation remains unselected |

## Goal

Define a practical SDL/SDUI history and blueprint workflow independent of project
Git, whose immutable records survive transport through parent Git merges. Preserve
reviewed targets and distinguish model acceptance from implementation evidence.

## Affected cards

| Card | Role | Initial state | Planned final | Current | Actual final |
| --- | --- | --- | --- | --- | --- |
| [KB048](../KanBan/backlog/%23048--Proposal--Versioned-design-reviews-and-blueprint-diffs.md) | Primary | backlog at late registration | Bounded design with explicit implementation handoff | backlog | Pending |
| [KB004](../KanBan/backlog/%23004--Proposal--Design-traceability.md) | Context: evidence mappings | backlog | Not disposed by this discussion | backlog | Pending |

## Plan register

| Plan | Document readiness | Canonical lifecycle | Outcome |
| --- | --- | --- | --- |
| [BP2 / PLAN-SDP-0001](../04--Design/SDPTool/Blueprints/Plan.md) | ready, requires revision if selected | planned | Existing blueprint contract successor; not activated here |
| [MAINT-SDP-0014](../Maintenance/SC1/Plan.md) | completed | completed | Immediate Session-governance correction |

## Route changes and decisions

| Date / entry | Change | Authority / consequence |
| --- | --- | --- |
| 2026-10-02 / retrospective | Reject export/bundle dependency and project-Git prerequisite; investigate full snapshots | Owner constraint; earlier Git recommendation not adopted |
| 2026-10-02 / T001 | Mutable WORK, immutable history, base/parent identities | Design recommendation; no implementation |
| 2026-10-02 / T002 | Consider four workflow roles and same-file/candidate merging | Owner proposal; release cannot receive changes |
| 2026-10-02 / T003 | Optional PROPOSAL; default WORK -> CANDIDATE -> RELEASE; preliminary WORK blueprint without locking | Owner correction supersedes mandatory-freeze recommendation |
| 2026-10-02 / T003 | No UUID in WORK names; four-character suffix for PROPOSAL/CANDIDATE; full identity in YAML | Owner naming requirement; collision handling still required |

| 2026-10-02 / T004 | Verb-first create commands; default accepted release; multi-source WORK creation and contextual pull considered | Owner preference; merge/pull naming remains open |

| 2026-10-02 / T005 | Typed creation targets; YAML ledger and commit history; snapshot/delta comparison | Owner proposal; storage representation remains open |

| 2026-10-03 / T006 | Limit content history to WORK; promotion keeps messages/lineage and final model, not undo payloads | Owner scope reduction supersedes permanent commit archive recommendation |

## Turn journal

All entries are manually summarized. Host turn IDs and precise timestamps are
unknown. Local T numbers are journal entries, not recovered host turn numbers.
No automatic routine engine or exact transcript capture is claimed.

### Retrospective context — before T001

The owner sought separate design history portable with parent Git and rejected
remembered bundle exports. SVN/Fossil raw-store merging did not satisfy parallel
history requirements. The assistant proposed project Git, then the owner clarified
that blueprints must work without Git/commits. Earlier analysis and bounded Git
storage-only probe are recorded in KB048; they do not prove a model-store backend.
Skills reported in earlier work: SDP change analysis and Architect. Missing original
Session upkeep is acknowledged, not retroactively represented as timely recording.

### T001 — Full snapshots and local checkout

Owner input summary: propose release directories, multiple WORK copies and local
file checkout/commit/pull. Assistant outcome summary: immutable revision identities,
optional local locks, conservative merging and explicit conflicting release labels.
Updated KB048/Study and event274. No model-store implementation. Session was omitted.

### T002 — Snapshot lifecycle and same-file merging

Owner input summary: UUID and SHA in YAML, browsable directories, WORK/PROPOSAL/
CANDIDATE/RELEASE, one integrator, and whether PROPOSAL is needed. Assistant outcome
summary: recommend frozen proposals/candidates and new identity after every merge;
record canonical hashing and three-way merge limits. Updated event275 and committed
analysis as c768564. Skills: SDP/Architect (agent-reported). Session still omitted.

### T003 — Correct Session upkeep, preview and naming; inventory tools

Owner input summary: immediately require Session updates in AGENTS and project
governance skill; allow preliminary WORK blueprints without locking; remove UUID
from WORK names and shorten proposal/candidate suffixes; default directly to
CANDIDATE; explain current SDL commands and SDPTool delegation.

No separate skill named
project_governance was found in the canonical collection or local skills search.
Updated the shared sdp entrypoint instead of adding a competing governance skill.

Request categories: process maintenance, architecture refinement, capability
inventory. Skills loaded/reused: sdp, sdp-architect, sdp-change-analysis,
sdp-planning and skill-creator; provenance is agent-reported, routine run ID unknown.

Work summary (not a captured final response): created this late-registered Session;
made per-turn upkeep explicit in AGENTS, sdp skill and Session guide; reconciled
skill patch metadata in the install inventory. Recorded latest owner corrections
in KB048/Study. Inspected SDL command entrypoints and SDPTool imports; standalone
SDL commands exist but SDPTool is not a universal plugin router. Snapshot/blueprint
commands are proposed only. Maintenance validation is recorded in SC1.

Next: refine S2 tool ownership and minimum command contract, then select S3. No
backend implementation, published release or consuming-project upgrade is claimed.

Additional owner steering in T003: identified SDL/go/cmd as the command location.
Confirmed: these are standalone Go program entrypoints; shared packages implement
the behavior and SDPTool imports several directly. This does not imply every
SDPTool operation has a separate standalone executable.

### T004 — Human-readable command grammar and integration

Owner input summary: put create before its noun, default create work to the latest
release, consider pull into the current WORK and create a new WORK from multiple
WORK sources. Asked whether create is a verb and blueprint a noun.

Request category: architecture/command design discussion. Skills reused: sdp and
sdp-architect (agent-reported); routine run ID and host IDs unknown. Steps: S2.
Work summary, not a captured final response: confirmed verb/direct-object grammar;
updated KB048 examples to verb-first syntax. Recommended consistent from/to/into
prepositions, multi-source create work as a non-destructive merge, and explicit
merge SOURCE into TARGET for existing WORK. Contextual pull remains an optional
alternative awaiting selection. Default source is an unambiguous accepted release;
missing baseline and conflicting heads must be explicit. No implementation.

Next: decide minimal command vocabulary and context rules in S2, then select S3's
bounded design/proof plan. Card remains backlog; BP2 remains planned.

### T005 — Typed targets, commits and history representation

Owner input summary: use create work:Combination consistently; include a ledger in
YAML and preserve merged WORK histories as a browsable tree; add a commit message
command; anchor WORK in a release or explicit empty initial state; consider storing
successive diffs and ask how Git represents revisions.

Request category: architecture/design discussion. Skills reused: sdp and
sdp-architect (agent-reported); routine execution and host IDs unknown. Step S2.
Work summary, not a captured final response: recorded typed-target recommendation,
commit versus preview distinction, DAG ancestry and deduplicated immutable ledger
records. Compared cumulative versus incremental patches; recommended full logical
snapshots first, optional deduplicated content storage later. Consulted official
Git objects/packfiles references; Git snapshots and pack compression are distinct.
No backend, command or automatic blueprint generator implemented.

Next: S2 must settle minimal persistent format and ownership before S3 design/pilot
selection. Ledger placement is a recommendation, not assumed owner acceptance.

### T006 — Bounded local undo and nested merge provenance

Owner input summary: avoid a full version-control system or central store. Keep
changed-file copies in numbered .commits directories inside WORK, preserve source
and target YAML/.commits/.merge in merge archives, and name YAML after the artifact.
Promotion can discard undo payloads while retaining messages and lineage. Asked
whether overlap prevents rollback and whether the proposed structure can work.

Request category: architecture/scope refinement. Skills reused: sdp and
sdp-architect (agent-reported); routine/host IDs unknown. Step S2.
Work summary, not a captured final response: accepted the bounded scope as current
owner direction; documented baseline and deletion requirements for after-image
reconstruction, metadata-only promotion and whole-state rollback versus selective
undo. Recommended pre/post merge checkpoints and safe staged archive copying;
recursive ancestry duplication needs limits/reuse, not a central repository.
Corrected permanent history expectations in KB048 and the Study. Validation checks
cover documents only; no storage, restore or merge implementation has been tested.

Next: settle minimal local schema and recovery operations, then select S3's bounded
plan. Keep permanent VCS features and selective merge undo outside the pilot.

## Closeout

Open. Model history/blueprint goal is not achieved. The immediate Session upkeep
correction is complete separately from the feature discussion.
