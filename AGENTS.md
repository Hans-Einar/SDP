# Collaboration and Git traceability

Owner decision of 2026-09-22 for SDL/SDUI development:

- Use one branch per phase and a separate commit for every completed milestone.
- Create each phase branch from the last commit of the preceding phase. Preserve
  phase branches so their deliveries can be reviewed separately.
- The combined target branch is `sdp-vNow`, originally at
  `9ad432407004080dd7f4f0ab06d107523f4316fd`. Do not develop directly on it.
- Include the phase/milestone ID in commit messages, with the concrete delivery
  and relevant verification in the body. Do not label unfinished work delivered.
- Check Git status before branching and staging. Stage only the milestone's work;
  preserve unrelated local changes and exclude generated caches.
- Update the implementation plan and evidence at each milestone. Generated SDL
  viewpoints must come from the SDL tool and validated model facts.
- Open a combined PR against `sdp-vNow` when the agreed work is ready. Local phase
  branches and milestone commits are authorized; merging is not. The owner's
  addition of 2026-09-22 authorizes pushing to origin after each completed phase
  in the same session.

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

Read the [KanBan workflow](SDP/Agents/KanBan/README.md) and relevant project board
before new work. SDP, SDL and SDUI have their own boards registered there.
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
