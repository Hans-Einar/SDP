# SDP vNow skill candidates

Status: **draft for evaluation; not installed or canonical**. All documents in
this collection are written in English. The directory name `skills_v2` denotes
a second skill proposal, not a new SDP product release or adoption of vNext.

## Purpose

Make existing SDP practice usable by agents receiving short owner prompts.
The agent must recover the relevant workflow, current design decisions and
evidence before turning a symptom into an implementation assignment. Preserve
vNow's project records, role separation and bounded execution; improve the
quality of the decisions made within them.

These files do not replace `Toolkit/skills`, change the installer, migrate any
project or authorize product changes. Existing project authorization and safety
boundaries still apply. A skill describes how to perform authorized work; it
does not grant additional authority.

## Branch evidence and adoption boundary

Inspected on 2026-09-15:

| Source | Exact head | Relationship and use |
|---|---|---|
| PR #4, `codex/sdp-install-contract-v1` | `d611b8bf72aeb30d86c5ef28902469462b06a803` | Working baseline for this proposal; retain its installation work. |
| PR #8, `codex/issue-7-provisional-vnext-pilot` | `ea9fcf1cdd3198aeac515b89b55398282c463838` | Separate lineage: PR #4 has 24 unique commits and PR #8 has 60. It is not based on PR #4. |
| Common ancestor | `e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15` | Explains why switching to PR #8 would not retain all PR #4 work. |

References: [PR #4](https://github.com/Hans-Einar/SDP/pull/4) and
[PR #8](https://github.com/Hans-Einar/SDP/pull/8). These are source observations,
not claims that either PR is accepted or merged.

Selected pilot ideas, expressed using existing vNow records:

- Keep durable product intent distinct from an individual assignment, branch
  and release. Reuse existing feature/requirement identifiers.
- Separate declared intent, observed implementation and accepted outcomes.
- Bound studies by a question; a study alone does not authorize implementation.
- Scale evidence to the consequence and uncertainty of the change.
- Declare ownership and shared touchpoints when concurrent work exists.
- Preserve accepted history; rework evidence must identify its new candidate.

Source material is the pinned pilot's
[WorkflowContract](https://github.com/Hans-Einar/SDP/blob/ea9fcf1cdd3198aeac515b89b55398282c463838/SDP-vnext-pilot/WorkflowContract.md),
[Profiles](https://github.com/Hans-Einar/SDP/blob/ea9fcf1cdd3198aeac515b89b55398282c463838/SDP-vnext-pilot/Profiles.md), and
[WorkDomains](https://github.com/Hans-Einar/SDP/blob/ea9fcf1cdd3198aeac515b89b55398282c463838/SDP-vnext-pilot/WorkDomains.md).

Deferred: domain UID registries, reservation epochs/digests, new cardinalities,
new schemas or ID grammar, cross-repository move protocols, mandatory Issue
creation, and removal of existing Sprint/Iteration/ledger requirements. The
pilot's optional-record profile does not waive an installed vNow contract.
PR #11's system-description-language research is outside this proposal.

## Index

| Skill | Primary responsibility |
|---|---|
| [sdp](sdp/SKILL.md) | Resolve project authority, route the request and apply document obligations. |
| [sdp-change-analysis](sdp-change-analysis/SKILL.md) | Translate a symptom into a system-grounded change boundary. |
| [sdp-architect](sdp-architect/SKILL.md) | Resolve material design choices and their consequences. |
| [sdp-master](sdp-master/SKILL.md) | Route and coordinate one authorized assignment. |
| [sdp-worker](sdp-worker/SKILL.md) | Implement the agreed outcome without silently redesigning it. |
| [sdp-reviewer](sdp-reviewer/SKILL.md) | Challenge assignment adequacy and implementation independently. |
| [sdp-verifier](sdp-verifier/SKILL.md) | Establish whether the evidence proves the claimed outcome. |
| [sdp-traceability](sdp-traceability/SKILL.md) | Keep decisions, work and evidence discoverable and consistent. |
| [sdp-steering](sdp-steering/SKILL.md) | Prepare owner decisions and record actual dispositions. |
| [sdp-auditor](sdp-auditor/SKILL.md) | Inspect installation and project consistency without changing it. |
| [sdp-vertical-refactor](sdp-vertical-refactor/SKILL.md) | Migrate complete workflows while preserving declared behavior. |
| [sdp-versioning](sdp-versioning/SKILL.md) | Select release versions from compatibility impact. |
| [sdp-release](sdp-release/SKILL.md) | Prepare and reconcile an explicitly authorized release. |

## Routing and proportionality

Follow the assigned role. Otherwise use Master when the project's SDP contract
designates it. For a symptom, unclear behavior or new capability, recover the
relevant context through change analysis before assigning implementation. For
a fully specified change, cite and reuse a still-valid analysis rather than
repeating it. Material alternatives belong with Architect; product choices
outside current authorization belong with the owner/Steering.

Not every skill means a separate agent, document or approval. Use separate
contexts where the installed contract requires independent review or where a
material design/evidence question needs an independent assessment. Do not ask
the owner to repeat authorization already given. Ask only about unresolved
choices that affect their intended outcome, scope or protected contracts.

Small work can use a concise section in its existing work record. Larger or
uncertain work needs a linked Study/Design record. Always preserve local vNow
requirements for Slices, reviews and traceability. Do not manufacture empty
documents simply to resemble a large project.

The chain to preserve is: owner observation -> user outcome -> applicable
decision -> proposed change -> exact candidate -> evidence -> disposition.
File changes are consequences of this chain, not the starting definition of
the user's problem.

## Skill purposes and boundaries

### SDP entrypoint

[sdp](sdp/SKILL.md) is the proposed `$sdp` entry for short owner prompts. It
selects relevant roles while preserving the assigned role and authorization.
Its [Document workflow](sdp/references/document-workflow.md) tells every role
which current SDP documents to read and maintain, how to capture newly
discovered requirements and how to distinguish capability from assignment.
Its [Source map](sdp/references/source-map.md) records historical provenance,
conflicting versions and deliberate adoption limits. It does not itself
install skills, introduce a new management role or adopt vNext schemas.

### Change analysis

[sdp-change-analysis](sdp-change-analysis/SKILL.md) produces a concise evidence
map of the affected user journey and its inherited constraints. It separates
facts, assumptions, derived obligations and product choices. It does not
implement, accept a new design or demand an exhaustive system model.

### Architect

[sdp-architect](sdp-architect/SKILL.md) compares viable approaches, records why a
choice fits the existing system and supplies implementation-ready boundaries.
It restores the original Architect emphasis on intent, workflow and state
ownership. It does not turn hypothetical future consumers into abstractions
or treat an agent's preference as an owner decision.

### Master

[sdp-master](sdp-master/SKILL.md) owns routing, scope and integration for the
current assignment. Its first gate is whether the problem is understood well
enough to delegate. It does not assume that a short prompt implies a trivial
code edit, nor automatically expand into the next assignment.

### Worker

[sdp-worker](sdp-worker/SKILL.md) implements a coherent bounded outcome,
preserves existing behavior and reports discoveries. It can resolve ordinary
implementation details within the contract. Changes to workflows, ownership
or accepted constraints return for a decision rather than becoming hidden
implementation choices.

### Reviewer

[sdp-reviewer](sdp-reviewer/SKILL.md) begins with owner intent and authoritative
context before implementation summaries. It reviews the adequacy of the
assignment as well as the diff. Review cannot establish owner acceptance or
compensate for missing runtime evidence with source inspection alone.

### Verifier

[sdp-verifier](sdp-verifier/SKILL.md) selects and checks evidence at the level
of the claim: function, service, workflow, application or release. It makes
untested and blocked cases visible. It does not weaken expected behavior to
make tests pass or claim final-commit evidence for uncommitted work.

### Traceability

[sdp-traceability](sdp-traceability/SKILL.md) keeps current decision rationale,
scope, exceptions and supersession navigable alongside work and evidence.
It uses installed schemas; it does not invent event types, rebuild history or
introduce a second authoritative system model.

### Steering

[sdp-steering](sdp-steering/SKILL.md) makes unresolved owner choices concrete
and records prompt, response, assessment and disposition distinctly. It
extracts the existing steering guidance from Architect into an explicit role.
It assists governance and does not impersonate the owner or approve its own
recommendation as a human disposition.

### Auditor

[sdp-auditor](sdp-auditor/SKILL.md) checks that skills are present, discoverable
and correctly referenced, and that declared project status agrees with
evidence. It reports gaps without installing, migrating or rewriting records.

### Vertical refactor

[sdp-vertical-refactor](sdp-vertical-refactor/SKILL.md) covers behavior-preserving
structural migration through runnable vertical outcomes. It identifies
temporary adapters and their removal conditions. It does not combine an
unrequested UI redesign with a structural migration.

### Versioning

[sdp-versioning](sdp-versioning/SKILL.md) separates release compatibility from
work coordinates and records a justified version target. It does not infer a
major version from code volume or treat this proposal's folder name as SemVer.

### Release

[sdp-release](sdp-release/SKILL.md) checks the selected contents and exact
candidate, prepares publication and records real publication results. It
does not publish without existing explicit authority or run release machinery
for ordinary development tasks.

## Mapping to the six proposed improvements

| Improvement | Candidate implementation |
|---|---|
| Correct format and distribution | Required `name`/`description` in every skill; candidate metadata; explicit discovery/adoption checks below. |
| Master pre-delegation gate | Master routes incomplete context to change analysis and material choices to Architect. |
| Concrete architectural understanding | Change analysis + Architect produce workflow, constraints, alternatives and decision rationale. |
| Review the assignment itself | Reviewer examines the owner outcome and governing contracts before the diff. |
| Preserve decision reasons | Traceability links current reasons and supersession; Steering records actual authority. |
| Verify at the correct level | Verifier distinguishes isolated tests from whole-application evidence and explains changed expectations. |

## Discovery, installation and compatibility

Reading these files for evaluation is not installation. Their existing SDP role
names are retained for eventual replacement, so do not expose both old and new
collections under the same names in one skill catalog.

Each candidate uses standard `name`, `description` and string-valued `metadata`.
`candidate-version` is proposal metadata, not canonical `skillVersion` and not
a claim that Toolkit's current validator accepts a new installed contract.
Promotion must reconcile metadata extraction, skill inventory and versions,
installer destinations, AGENTS paths, validation fixtures and references as
one reviewed change. Existing schemas/manifests remain untouched here.

Current [OpenAI skill documentation](https://learn.chatgpt.com/docs/build-skills)
describes repository/user `.agents/skills` discovery. The inspected vNow
installer and AGENTS template use project `.codex/skills`; bundled tooling in
this environment also uses a user `.codex/skills` location. Test the actual
supported host/version rather than assuming every path is interchangeable.
An AGENTS instruction can explicitly load a file but does not itself prove
native catalog discovery. A local installation does not prove availability
on another machine, container or cloud environment.

Adoption evidence must distinguish: file present; catalog discovery; selected
skill loaded; observable task behavior consistent with the skill. Skill text
alone cannot enforce workflow or guarantee compliance.

## Evaluation before promotion

Validate all frontmatter and Markdown links. Then give fresh contexts realistic
owner prompts and raw project evidence without the expected answer. Evaluate
at least: a UI reachability problem, a well-specified small fix, a discovery
that would alter workflow, and a stale-evidence handoff. Check that the agent
finds governing decisions, avoids unnecessary ceremony, distinguishes choices
from obligations, and reports unsupported claims accurately.

Keep evaluation outputs separate from project truth. Record failures and make
only evidence-backed revisions. No claim of production effectiveness follows
from a valid Markdown file or one successful trial.

### Initial evaluation, 2026-09-15

- All 12 entrypoints passed the bundled skill-creator `quick_validate.py`
  validator. All local Markdown links resolve, and each skill has both an index
  link and a link in its purpose/boundary section above.
- A fresh agent received only Change Analysis and Architect as initial skill
  instructions, the short prompt below, and access to raw Concept1 source and
  requirement/architecture/design records. Prior conversation conclusions and
  the previous APT study/verification records were excluded. No product writes,
  browser interaction or hardware actions were authorized in the trial.

> On the APT / Apteringsfil page I cannot reach the bottom rows. Please figure
> out what should be changed.

The agent found the enclosing fixed-HMI geometry and explicit inner-overflow
policy, distinguished ancestor clipping from matrix scrolling, preserved
import/load/draft boundaries, and returned a source-supported recommendation
with missing runtime and owner-decision evidence stated. Its proposed inner
scroll design is an evaluation output, not an accepted product solution.
It also identified an existing activation-panel removal as an unresolved
authority boundary because its approval evidence was deliberately unavailable.

Evaluated entrypoint SHA-256 values:

- Change Analysis: `db5a5fa647a8e71535ca670ef0ca3d2f359554429bbffd6ee9fc45067457e448`
- Architect: `23c964ca3e171a4f676083046b06e8f80aa1c558b3c7a853c06836b8da747b35`

Observed friction: workflow/history investigation can grow beyond a small
symptom, and the two skills should share one analysis handoff rather than
duplicate records. Test proportionality on a well-specified small fix before
adding more instructions. The other ten skills have format/content inspection
only, not independent behavioral validation. Native discovery, installation,
full role orchestration and runtime product behavior remain unverified.

### Document workflow revision, 2026-09-15

Draft.2 adds the `sdp` router and shared document obligations to all twelve
role candidates after reviewing all sixteen `docs/` Markdown files and the
bounded root/GitHub sources listed in the source map. Previous behavioral
results apply only to the hashes above, not to this revised collection.

The revision makes durable Requirements/Architecture/Design updates part of
the assignment and review, rather than treating ledger maintenance as enough.
It retains current installed records where Feature governance is not deployed.
Feature/functionality and user/system distinctions remain descriptive views;
Tier is explained as vertical capability, not a renamed Sprint or Iteration.

For eventual project adoption, project instructions should require agents to
load the installed SDP entrypoint for SDP-governed work. Owners may also use
`$sdp` explicitly; normal description-based selection remains possible. Neither
mechanism guarantees compliance. Test catalog discovery, reference loading,
role routing and document maintenance independently. All referenced roles and
supporting files must be installed as a compatible bundle. No installation or
project-instruction changes are included in this draft.

All thirteen draft.2 entrypoints passed `quick_validate.py`; local Markdown
links resolve and every skill has an index and purpose-section link.

A fresh, read-only context manually loaded `sdp` for an implementation-discovery
scenario: import and activation had been conflated; the project had existing
lifecycle documents but no Feature registry; the owner also asked about Tiers.
No expected answer or product documents were supplied. It loaded the two shared
references, Worker and Master, identified missing behavioral authority, mapped
consequences to existing documents, preserved the distinction between a ledger
event and current design, and avoided inventing Features or renaming Sprints.
It correctly reported that the supplied facts did not determine the intended
import/activation behavior. It treated the response as method explanation,
not an executed Worker assignment. No routing revision was justified by this
trial. This establishes a limited response-level result, not successful real
document edits, native activation or full implementation/review orchestration.

Evaluated collection digest:
`c95b4c3ce6144fe100d719fe132f1f489c39b4c7c44f385a5f06904271e7f594`.
SHA-256 input is each Markdown file except this README, sorted by relative path,
concatenated as UTF-8 relative path, NUL, exact file bytes, NUL.
