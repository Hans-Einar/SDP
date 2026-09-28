# RGS1 evidence and limits

Research date: 2026-09-28. This is a bounded process study under
[PLAN-SDP-0006](Plan.md), not a product audit of the portfolio. Read
[Study.md](Study.md) for recommendations. [Source-pins.json](Source-pins.json)
records sampled file hashes and source revisions. A hash identifies the fetched
file, not proof that every line was reviewed or that its claims were reproduced.

## Method and confidence

The study combines current local records, selected GitHub files at exact commits,
GitHub PR observations, earlier usage studies and official interface documentation.
Large records were sampled around authority, work status and review disposition;
products were not built, run or tested. No independent agent review or routine
engine trial was performed. No external repository was changed.

The owner's report of repeated divergence is a research input. The sampled records
support concrete problems with stale status, context boundaries and evidence
handoff. They do not establish that agents caused every reported product defect,
or quantify the incidence of drift across the portfolio.

Private raw documents remain outside this repository. Findings below describe
process mechanisms, without reproducing business records, device data or private
implementation excerpts. Private source links require the reader's own access.
Temporary research caches are not a distribution dependency.

## E01 — SDP itself

Baseline main: `9e4c173`; study activation: `3d969cf` on `sdp/request-routine`.
The local source sample includes skills, current plan/management contracts,
SDPTool CLI, installer records/profile and the blueprint study. Unrelated
untracked `SDL/go/sourceinput` and Node files were excluded.

What is implemented: project discovery/navigation, selection and previews through
the Go SDPTool; signed installation/upgrade mechanisms; canonical role skills;
KanBan, plans and two distinct histories. Inspection found no generic routine
engine or SDP MCP adapter in the current command interface. Installer operation
journals are specific recovery machinery, not a generic workflow engine.

What still needs integration: [BP2](../../04--Design/SDPTool/Blueprints/Plan.md)
remains a planned blueprint contract/pilot. The
[BP1 study](../../04--Design/SDPTool/Blueprints/Study.md) distinguishes generated
model facts from authored intent and observed evidence. Do not describe it as
an implemented assignment compiler.

A concrete documentation inconsistency: the inspected SDPTool contract and skill
source map still refer to `Toolkit/SDP-install.manifest.json`, while the current
Go profile selects `SDPTool/profiles/payload.json`. This study records the stale
reference; it does not repair application/process authority during research.
A routine resolver must expose such inconsistencies rather than silently choosing
whichever text happens to enter an agent's context first.

[SK1 evidence](../../Maintenance/SK1/Evidence.md) already separates native skill
discovery from actual loading and bounded behavioral trials. Those historical
trials are useful precedent, not proof of this proposed routine system. The
current host reports Codex CLI 0.157.1; hooks were not configured or exercised.

## E02 — XFMD: installation, product authority and application evidence

Local worktree: `/home/warloc/git/xfmd-sdl-navigation`, HEAD
`41fa4948c5377edd792ba3f503e6b92b0fd4717a`. Its SDP 1.0.0 installation and process
updates are uncommitted alongside pre-existing card work. HEAD alone therefore
does not identify all inspected process bytes; local hashes accompany the pins.

The sampled `docs/working-method.md` distinguishes requirements, blueprint,
contracts/plumbing, application implementation and verification. It distinguishes
building from installing and scope-specific authorization for merge/install.
Native ownership spans application, interpreter and renderer. Installing a new
SDP payload does not by itself supersede all project-specific working rules.

Implication: capture both the installed process identity and the project's
current adoption/authority state. A native sidepanel change and a renderer fix
need different impact boundaries even when both appear as a small UI request.
No XFMD application verification was rerun for this study.

## E03 — Ponsse: nested authority and the difference between stopping and finishing

Local Concept1: branch `codex/concept1-gh65-manual-test-ready`, HEAD
`882ad7c0457a723550ea6cce5bf027230997ec2a`, with unrelated dirty application and
local document files. Sample: root AGENTS, Concept1 CurrentIndex and WORK-101
handoff. The records distinguish a bounded stop, independent evidence, outstanding
owner/manual acceptance and permission for later work. Root instructions preserve
offline/no-physical-I/O boundaries and explicit confidence levels.

MVP1 was sampled at branch `steering/mvp1-wave1-integration`, commit
`7589a5811a21ec6bb13979612060ceff0f690df8`, through
[its SDP entry](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/SDP/README.md).
This is a separate development line with its own authority precedence and
pre-implementation review disposition. Repeated corrective records are not proof
that the latest candidate is accepted. The default branch was separately resolved
at `9cc5a691c60c977bdf4fea493a13437017013fa7`; it is not substituted for MVP1.

The local [MVP1 SDL experiment](../../../experiments/mvp1_sdl/README.md) models
nine containers across multiple files and separates BuckingUI/SimulatorUI. Its
experimental SDL profile is not supported by the released Go parser. Static SDUI
previews do not establish executable SDL coverage of the complete system.

Implication: “continue Ponsse” needs a project/worktree/assignment resolution step.
Unknown model coverage must remain visible in a blueprint. A routine must never
turn offline investigation into physical machine operation implicitly.

## E04 — HSX: default branch is not the whole design context

Sampled default main `e374da88f4dd470bad2d8ec6a2a14f1ce367e40e`:
[agent instructions](https://github.com/Hans-Einar/HSX/blob/e374da88f4dd470bad2d8ec6a2a14f1ce367e40e/agents.md)
and the documented TUI/debugger integration gaps. Those gaps are assertions in
project records; they were not independently reproduced against running code.

Separately inspected `codex/dbg-rf-002-003` at
`69a54aeb3394d3cd4792bce620748e15bab69f1f`, especially
[shared process](https://github.com/Hans-Einar/HSX/blob/69a54aeb3394d3cd4792bce620748e15bab69f1f/SDP/Shared/Process.md).
That branch defines separate tracks and study/refactor ownership. Other branches
exist, including dbg-rf-004; the sampled branch is not asserted to be the latest
authoritative candidate. A live PR check also found an open draft targeting
`Implementation/vscode`, rather than main.

Implication: dependency and authority resolution must include the selected branch,
not infer the current assignment from the repository's default page. Studies and
product-integration findings must feed scope before a worker receives a local fix.

## E05 — agro-crm: local correctness can miss another consumer

Default main: `4e6c01ac328779c45e5e5707bbd0e839b2119a3e`. Sampled AGENTS,
[StateUpdate review](https://github.com/Hans-Einar/agro-crm/blob/4e6c01ac328779c45e5e5707bbd0e839b2119a3e/SDP/CodeReview/001-StateUpdateLogic.md)
and a coverage-planning iteration record. Instructions distinguish the project's
authority from example SDP content. The review reports differing update semantics
between derived views. Iteration notes show evolving UI/algorithm work and
cancellation/restart handling.

Implication: impact analysis needs consumers, derived state and cancellation
semantics, not only the edited function. Review observations are evidence of a
reported problem at their recorded candidate, not a claim that current main still
contains the defect. No private example data is reproduced here.

## E06 — TerrainAnalyzer: stored remote status can be demonstrably stale

Default main: `9c422f7d02e6e1d1b0ac212abf2010f883783f31`. Sampled instructions,
CurrentIndex and
[CurrentAssignment](https://github.com/Hans-Einar/TerrainAnalyzer/blob/9c422f7d02e6e1d1b0ac212abf2010f883783f31/SDP/Steering/CurrentAssignment.yaml).
The assignment stores PR 21 as `draft_open` and work as
`complete_pending_steering`. A fresh GitHub query on the research date reports
[PR 21](https://github.com/Hans-Einar/TerrainAnalyzer/pull/21) merged at
2026-08-19T17:53:15Z, head `8163c74d6643d8f65fb1f4d9523b3935573bb474`.

This is a verified metadata contradiction. It does not establish who authorized
the merge or resolve all remaining Steering conditions. An observer should show
both the stored disposition and the newer remote observation, then request
reconciliation through the appropriate routine. It must not rewrite approval
history based on “merged”.

## E07 — GrassPhenology: a review is a delivery, not permission to repair

Default master: `8ba7567aaaf0524d86bdb4675a09232b41496f36`. Sampled SDP README and
[RFI-004A](https://github.com/Hans-Einar/GrassPhenology/blob/8ba7567aaaf0524d86bdb4675a09232b41496f36/SDP/Refactors/Refactor-001-POCWorkbenchArchitecture/07--Iterations/RFI-004A-WorkspacePatternReview.md).
The review explicitly classifies ownership, pins its evidence and excludes fixes
from its selected scope. It prepares a later slice without authorizing it.

Implication: completed analysis can return a larger capability to backlog.
A workflow must distinguish “review found a necessary change” from “that change
has been selected for implementation”. No current application behavior was tested.

## E08 — processor emulator and debugger: compatible peers and real evidence

Emulator master: `b43fe36b0965b6ac8628677bb6fcc16513d1f567`.
[EMU-DEBUG-ARCH-002](https://github.com/Hans-Einar/emuSA80535-N/blob/b43fe36b0965b6ac8628677bb6fcc16513d1f567/SDP/04--Architecture/EMU-DEBUG-ARCH-002.md)
separates the generic emulator/debug interface from client and board-specific
policy. Instructions preserve that boundary. The study concerns this repository
and `emuSA80535-DAP`, not an independently identified third simulator project.

DAP main: `2e4c2dae2bb3637c4f3dbf803b83cc3f7e246301`.
[Its review](https://github.com/Hans-Einar/emuSA80535-DAP/blob/2e4c2dae2bb3637c4f3dbf803b83cc3f7e246301/SDP/CodeReview/DAP-STEERING-REV-001.md)
distinguishes useful fake-runtime tests from required compatible real-runtime
evidence. The README names a paired emulator revision, demonstrating why a lone
repository HEAD is insufficient integration identity. These checks were not rerun.

The sampled README still calls PR 4 unmerged; the live
[PR 4 observation](https://github.com/Hans-Einar/emuSA80535-DAP/pull/4) reports a
merge at 2026-09-02T08:30:12Z. This is a second verified status contradiction.
Implication: preserve consumer/provider contracts and paired candidates; do not
promote fake-only results to system verification or unverified external card
references to resolved dependencies.

## E09 — earlier research is input, not a newly adopted process

The [cross-repository study](../../Studies/UsageAnalysis/CrossRepositoryPatterns.md),
[proposed workflow](../../Studies/UsageAnalysis/ProposedSDPWorkflow.md) and
[Feature governance proposal](../../../docs/process/Feature-Governance-And-SDP-2.0.md)
already discuss duplicated current state, assignments, review and durable feature
identity. Their Issue-first mechanisms and older SDP 2.0 organization are not
silently reinstated. Current KanBan, optional Scrum/Sprint and typed plans remain
authoritative. This study does not analyze all old issue-7 commits again.

## E10 — documented interface capabilities, not a configured integration

OpenAI's [Docs MCP](https://developers.openai.com/learn/docs-mcp) provides read-only
documentation search and content. This is a useful model for SDP knowledge access,
not an example of a project execution engine.

Official [Codex MCP documentation](https://learn.chatgpt.com/docs/extend/mcp?surface=cli)
describes STDIO and Streamable HTTP connections. The
[Go MCP SDK](https://go.sdk.modelcontextprotocol.io/) provides a candidate Go
implementation basis; this study selects no dependency version. The explicit
[MCP resources baseline](https://modelcontextprotocol.io/specification/2025-11-25/server/resources)
is 2025-11-25, not a claim that it is the newest protocol. Subscriptions require
capability negotiation and are not a substitute for a durable execution journal.

[Codex hooks documentation](https://learn.chatgpt.com/docs/hooks) describes
UserPromptSubmit and PreToolUse. The latter covers selected command, edit and
MCP/local-function tools, but is not a complete interception boundary: some paths
are excluded and write_stdin does not repeat it for later input. This makes a
host adapter worth testing; it does not prove every request or mutation can be
guarded in this session. No hook setup, bypass test or live monitor was performed.
