# Community feedback after an understood internal pilot

| Field | Value |
| --- | --- |
| id | KB-SDP-035 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| created | 2026-09-27T21:01:36.384239+00:00 |
| source | Owner conversation 2026-09-27: MCP, agent context and community readiness |
| tags | community, adoption, process, evidence |
| next_review | After a repeatable internal MVP1 workflow and owner readiness review |

## Owner intent and timing

The owner wants SDP, SDL and SDUI to prove useful in our own projects and eventually
attract an open-source community. An MCP adapter around SDPTool, accompanied by
skills describing when and how to use it, is a promising direction for making
agent context and operations explicit and inspectable.

Owner decision: reach a somewhat more mature state before involving others.
We need to understand and use our own process first. Early feedback is welcome
once that foundation exists; presenting a confusing or unproven workflow risks
losing people's interest. This does not require completing every planned feature.

Registration authorizes no announcement, outreach, submission, application or
MCP implementation. Keep this card in backlog until readiness is reviewed.

## Proposed readiness evidence

- Complete and repeat a small real-project workflow, preferably in Ponsse/MVP1:
  inspect the design and surrounding contracts, define bounded work, make a change,
  and review its verification evidence.
- The owner can explain and navigate that workflow, including authoritative
  sources, current work and generated views, without reconstructing old chats.
- Provide a compact runnable example or walkthrough with a clear entry point,
  actual results, supported profiles and explicit limitations.
- Distinguish delivered capabilities from proposals. If demonstrating MCP-based
  context retrieval or drift prevention, demonstrate those operations and checks
  first. MCP availability and skill instructions alone do not prove compliance.

These are preparation proposals; the owner decides when the process is understood
well enough to invite others. There is no arbitrary maturity score or deadline.

## Later feedback approach

Start with a few builders who recognize architecture/context drift. Ask whether
they can understand the example, recover relevant context, and identify missing
or unnecessary steps. Seek related projects and alternatives; do not claim SDP
is unique without a proper comparison.

Possible thread title: **SDP: connecting system design, UI prototypes and
coding-agent workflows — looking for early feedback**. Explain the original
problem, show one working example, distinguish implementation from roadmap and
ask specific questions. This is draft positioning, not published material.

Candidate channels checked on 2026-09-27; recheck guidance before posting:

| Channel | Intended use |
| --- | --- |
| [OpenAI Developer Community](https://community.openai.com/) | One lasting project thread, concrete example and technical feedback |
| [OpenAI Discord](https://discord.com/invite/openai) | Find interested builders for a small trial; follow channel rules |
| [OpenAI Developer Showcase](https://developers.openai.com/showcase) | Submit a reproducible demonstration once ready |
| [Codex for Open Source](https://developers.openai.com/community/codex-for-oss) | Separately consider maintainer support; eligibility is not assumed |

[Forum project-sharing guidance](https://community.openai.com/t/where-can-we-showcase-our-ai-webapp-built-with-codex-5-2-as-a-high-value-add-tool/1379215)
recommends the Community category/project tag, one project thread for updates,
and emphasis on how it was built. The [official community hub](https://developers.openai.com/community)
links the forum, Discord and showcase.

## Dependencies and completion

Use [KB-SDP-032](%23032--Study--Viewpoint-navigation-feedback.md) for internal
viewpoint feedback and the [MVP1 navigation handoff](../../../experiments/mvp1_sdl/Navigation.md)
for current capability boundaries. [KB-SDL-005](../active/%23005--SDL--Change--System-and-source-sets.md)
owns source-set support. The [blueprint study](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md)
records assignment-context design; it does not deliver an MCP adapter.

Next action: revisit after an internal demonstration, record the owner's readiness
decision, and select proportionate preparation/outreach work. MCP implementation
needs its own bounded work selection rather than being hidden in this card.
Complete with the actual disposition and, if outreach is later authorized, links
to published material and recorded feedback with concrete follow-up owners.
A draft announcement alone does not complete the card.

## Worklog

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-09-27T21:01:36.384239+00:00 | codex; EVT-KB-SDP-000205 | Registered the owner's maturity-first community-feedback decision | Internal pilot and owner readiness review pending; nothing published |
