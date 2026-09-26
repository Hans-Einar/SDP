# Collaboration and Git traceability

Owner decisions of 2026-09-22 and 2026-09-25 (KB-SDP-029):

- Follow the selected plan's explicit BranchPolicy and CommitPolicy. Small plans
  normally use the current working branch with commits per phase or milestone.
  Large implementation plans may select stacked phase branches and milestone
  commits. Preserve historical plans' already selected branch commitments.
- Do not develop directly on main. Create a working branch when needed. The
  owner-authorized MP1 integration on 2026-09-26 promoted sdp-vNow to main;
  sdp-vNow and the earlier phase branches now preserve development history.
- Include the phase/milestone ID and concrete delivery in commits. Update plan
  and evidence at each delivery; never label unfinished work delivered.
- Check Git status before branching/staging. Preserve unrelated changes and
  exclude caches. Generated SDL viewpoints must come from validated model facts.
- Existing owner authorization permits phase pushes to origin and a combined
  PR against main when agreed work is ready. Merging requires explicit owner
  authorization; MAINT-SDP-0005 received that authorization for its integration
  and closeout only. Do not infer authorization for later merges or releases.
- Read the [plan contract](SDP/ProjectManagement/Plans.md). For plan creation,
  selection or revision, load Skills/sdp-planning/SKILL.md after the SDP entrypoint.

The [development history](SDP/Development-Branch-Stack.md) describes the branch
stack and handling of earlier uncommitted work. Read applicable subdirectory
AGENTS.md files before making changes.

## Documentation language — English

Owner decision of 2026-09-24: **all maintained documentation must be in English**.
This includes READMEs, language specifications, studies, plans, decisions, agent
instructions, evidence narratives, KanBan cards, metadata values, worklogs,
queue/review explanations and new ledger event descriptions. This rule applies
to SDP, SDL, SDUI, Toolkit and distributed templates. Do not add Norwegian prose
because the owner discusses the work in Norwegian. Conversation replies may
remain in the owner's language.

Translate existing prose when maintaining it; preserve its meaning, authority,
identifiers, technical syntax and provenance. Clearly identify verbatim source
quotations or intentional localized UI example text. Frozen test fixtures,
append-only historical event bytes and hash-pinned machine evidence are records,
not editable prose: retain them and use English explanations for new records.
Generated documentation must be translated in the generator/source and rebuilt,
never hand-edited. Do not describe translation as new implementation evidence.

## Local file URLs in responses

Use plain absolute `file:///` URLs, one per line, without Markdown wrapping,
line-number suffixes or terminal control codes. Expand the home directory and
percent-encode spaces and reserved URL characters.

## KanBan for ideas and scope changes

Read the [KanBan workflow](SDP/KanBan/README.md) and relevant project board
before new work. One shared board retains SDP, SDL and SDUI identity namespaces.
Record new ideas and out-of-scope findings in a primary card, with Ref cards
where needed, before changing direction. Registration alone does not authorize
implementation. Update the card, append-only KanBan ledger and links together
when processing or moving a card. Review backlog/onHold at milestones and before
the next phase. Keep cards useful as working documents, including progress,
evidence and explicit remaining work.

KanBan records how proposals are handled; implementation and verification still
belong in Traceability. Do not migrate old process areas or adopt language rules
merely because a backlog card describes them.

## Card work state

Maintain exactly one visible CardState metadata row on each KanBan card. Read
README.md for allowed states and folder mapping. Select ready, use in-progress
while working, and gate-review only for a concrete owner review. Queued cards
remain in backlog with a Queue explanation and predecessor reference. Update
worklog and ledger at changes; do not leave delivered work active solely to point
to future work. Use kanban status/state for listing; there is no separate gate file.

## Systems and shared process ownership

Owner decision of 2026-09-25: SDL, SDUI and SDPTool are three separate software
systems documented in the shared root SDP area. Do not create a full parallel
process tree per language/system. The shared board is SDP/KanBan; earlier SDL/SDUI board histories are frozen
under SDP/ProjectManagement/History with a byte-preserving import map. Group
architecture/design sources by system and actual container responsibility;
shared libraries are not automatically runtime containers. System is a selected
design concept, not yet implemented syntax in design-core 0.5.
XFMD is a collaborating system with its own SDP area in its repository. Its
process bootstrap and adoption card do not authorize XFMD application changes
from this workstream. See SDP/03--Architecture/System-Boundaries-study.md.

## Project-management history and optional sprints

Owner decision of 2026-09-25: read SDP/ProjectManagement/README.md before lifecycle
work. Its Ledger.ndjson is the only writable history for cards, Scrum, Sprint,
Maintenance, CodeReview and Refactor. Do not append management-only transitions
to Traceability. System design/code evidence references management IDs/events;
new traceability IDs use the affected system prefix. Preserve old IDs/ledger bytes.
Scrum reviews live in SDP/Agents/Scrum. A review can select Maintenance directly;
no mandatory wrapper card or Sprint. Optional SprintId/ScrumId metadata groups
cards without copying/moving them into sprint-specific directories. Actual sprint
start activates its selected backlog cards as ready with recorded events; retain
existing CardState meanings. Current metadata and ledger/record membership must agree.

## Project skills

For SDP-governed study, implementation, review or Maintenance, load
Skills/sdp/SKILL.md and only the relevant roles it routes to. Preserve an explicitly
assigned role and current host/session delegation rules. A trivial factual question
or unrelated task does not need SDP lifecycle work. Root Skills/ is the only
maintained source; .agents/skills contains relative discovery symlinks. Do not
create a second collection or imply that discovery proves actual loading.
