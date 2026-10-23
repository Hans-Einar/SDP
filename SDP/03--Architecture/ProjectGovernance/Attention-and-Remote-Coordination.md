# ProjectGovernance — autonomy, attention and remote coordination

Status: owner direction and proposed architecture, 2026-10-19 / Session0007 T025.
[KB-SDP-054](../../KanBan/backlog/%23054--Proposal--PG-attention-and-remote-coordination.md)
registers the multi-machine service and client extension. Existing local work stays
under [KB052](../../KanBan/active/%23052--Change--Project-Leader-and-SAD-orchestration.md)
and [PLAN-SDP-0018](../../05--Implementation/SDPTool/ProjectGovernance/Plan.md).
No remote service, deployment or installed app is delivered by this document.

## Owner decision and current gap

The owner selects as much automatic execution as possible, including status/result
reporting within an authorized assignment, so large projects can work for hours
without continuous supervision. PG must also call for attention through a durable
message/question interface, comparable to follow-up questions, across projects on
different PCs. A self-hosted service should expose oversight and control through
desktop/mobile web and an installable notification-capable app. The old SDP_Analyzer
is obsolete and is not the baseline or a dependency; no old project is modified.

Current PG has local run/attempt state, approved implementation activation and
per-run controller exclusion. It has no shared mailbox, cross-machine connection,
remote owner commands or installed web client. C2d performed real code work but its
MCP report was refused under never; C1 concurrent progress is protocol-test evidence.
The new direction resolves the pending reporting preference. It does not convert
that earlier failed report into success or override host-enforced restrictions.

Automatic reporting means permission to submit evidence, not to accept its own
result. Existing scope, design, review and explicit release/merge boundaries remain.
The runtime should continue authorized independent work while an affected dependency
waits for input. Missing input is not approval; alerts must identify what is waiting
and what can continue. Ordinary progress goes to activity/history, not a question
that requires an owner click. Long-running operation is controller scheduling with
durable state and recoverable agent turns, not an immortal model context.

## Proposed responsibilities

| Component | Responsibility | Authority boundary |
| --- | --- | --- |
| Local PG runtime on each machine | Execute existing services; persist runs, reports, decisions and outbound events; reconnect and resume within mandate | Authoritative for local execution and current source/assignment checks; credentials remain local |
| Self-hosted coordination service (working name: PG Hub) | Enroll machines/projects, persist the shared mailbox, route commands/replies, project activity and manage notification subscriptions | Authoritative for mailbox/delivery records; a command is pending until the target runtime acknowledges and applies it |
| Responsive web client, installable as a PWA | All-project overview, activity, attention inbox, contextual replies and targeted control | Authenticated user interface; closing it does not stop local work |
| Optional external adapters | Web Push and potentially ChatGPT MCP Events | Delivery surfaces for the same messages; no separate approval/state history |

This extends the ProjectGovernance ecosystem. Runtime component names here are
proposals, not new shipped binaries or automatic extraction of SDPTool libraries.
Existing SDPTool process rules and project ledgers remain authoritative. The hub
does not copy those rules into a competing workflow engine. Native XFMD changes
and migration of SDL/SDUI sources are outside this assignment.

Prefer an outbound authenticated connection from each enrolled runtime to the hub,
so ordinary use does not require exposing a local control port on every PC. Choose
transport, implementation stack, repository and identity provider in the bounded
follow-on design. Device credentials must be separately revocable; user permissions
must scope projects/actions. Expose governed actions, not a remote arbitrary-shell
endpoint. Do not send account tokens, hidden reasoning or whole repositories by
default; transfer the selected status and review evidence needed by the user.

Local authorized work continues through hub/network outages where dependencies
permit. Persist an outbox, retry with stable message/command IDs and reconcile
receipts. The hub must distinguish last observation, disconnected/unknown state,
queued command, received command and applied result. A disconnected PC cannot be
promised an immediate stop; show that pending control explicitly. Remote decisions
bind project, checkout, assignment and current question/design revision. Late,
duplicate, expired and superseded answers must never launch duplicate work or apply
to a newer package. No failover may create a second executor without an explicit
ownership/recovery protocol.

## One durable attention workflow

1. An agent or runtime emits a correlated report or attention request into local
   durable state. Reports are automatic inside the assignment's granted capability.
2. The runtime forwards the envelope to the hub, which deduplicates and stores it.
3. The inbox shows project, machine, Session/task, requester, reason, question,
   options/free text, evidence links and the exact blocked dependency. Priority and
   notification preferences determine whether to send an alert or batch a digest.
4. The owner opens the item and responds. Reading, acknowledging, snoozing and
   making an actual decision are distinct operations. Never choose a default answer
   merely because the owner does not respond.
5. The target runtime rechecks revision and authority, then applies or rejects the
   response and records the result. Project Leader and direct SAD/Master views see
   that same disposition. Closing the alert is not evidence that work resumed.

Proposed semantic categories: informational progress, nonblocking question, blocking
input/decision, and operational failure. A full typed wire schema is follow-on work.
Each needs stable event identity, origin/target binding, correlation, revision and
explicit acknowledgement/disposition. Critical unresolved items remain in the inbox
even if push delivery fails or the user silences notifications. Coalesce duplicate
alerts without erasing substantive updates. The inbox must be useful with push off.

For example: two projects continue their implementation work. A third needs a decision
about a changed data format. The phone shows a concise question with alternatives;
only the dependent task waits. After a current authorized reply is acknowledged by
that PC, the task continues automatically. If its proposal changed meanwhile, the
old answer is shown as stale and a revised question is retained.

## Client and integration recommendation

Start with a responsive website plus PWA installation. This provides one desktop/
mobile interface and an app-like installed entry point without committing to separate
native clients. Browser/OS support and permissions must be tested on the actual
phones. A native wrapper/client remains an option if measured needs such as stronger
OS integration cannot be met. Backend and inbox contracts must not depend on that
choice. A browser/service worker is not the long-running PG executor.

MDN documents PWA installation, manifest/HTTPS requirements and platform differences.
WebKit documents Web Push for installed Home Screen web apps on iOS/iPadOS 16.4+,
with user permission. The application/backend can be self-hosted while standard
mobile push delivery still uses browser/OS push services, including Apple's APNs;
self-hosting does not imply a completely independent push network. No push delivery
latency or uninterrupted background execution is guaranteed by this proposal.

OpenAI currently documents MCP Events: a user subscribes a Work chat to server
events delivered by webhook. The documented surfaces are Work on web, desktop Work
with Cloud selected, and dots; it requires MCP 2.0 / 2026-07-28. This is a candidate
adapter to investigate, not verified support in the owner's current app/account,
not an arbitrary push-to-any-ChatGPT-user API and not a reason to replace PG's inbox.
A received chat event is not proof of an OS notification or a durable PG reply.
Keep this adapter separate from the local pinned host/worker MCP protocol. Plugin
panels, deep links and forms are also documented, but integration remains untested.

Sources inspected for this proposal (2026-10-19):

- [PWA installation](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Guides/Making_PWAs_installable)
- [Push API](https://developer.mozilla.org/en-US/docs/Web/API/Push_API)
- [WebKit Home Screen Web Push](https://webkit.org/blog/13878/web-push-for-web-apps-on-ios-and-ipados/)
- [OpenAI MCP Events](https://developers.openai.com/plugins/build/mcp-events)
- [OpenAI plugin extensions](https://developers.openai.com/plugins/build/extensions)
- [ChatGPT notification controls](https://learn.chatgpt.com/docs/notifications)

## Delivery sequence and acceptance

First finish the local PGL3 reporting/inbox boundary with the owner's automatic
reporting decision, including a supported host path and measured reporting evidence.
Do not disable global approval controls or claim a callback bypasses a host refusal.
Add a typed attention request and reply using the same local authority and correlation
rules; make it usable locally before adding a network dependency.

Owner T026 assigns the hub and web application to a separate project. KB054 remains
the SDP-side handoff and future adapter coordination record until an external
repository/board is selected. That separate project selects its own bounded plan:
two enrolled machines and projects,
a persistent hub inbox and read-only overview first, followed by authenticated replies
and targeted controls, then PWA installation and push. ChatGPT integration is an
optional parallel investigation and not a prerequisite. Large-project unattended
operation requires an actual multi-hour trial with a declared assignment/concurrency/
resource budget, restart/reconnect, expired authority and one answered blocking item.
Record model/account rate limits and exhausted budgets as operational conditions;
do not promise unlimited execution.

Verify duplicate/out-of-order delivery, hub and runtime restart, offline response,
stale decision, wrong-project/device/user refusal and lost control acknowledgement.
An unauthorized or stale request must not execute. A queued stop must not look applied.
Verify one task waiting while independent authorized tasks progress; denied/missed
push must leave a usable inbox. Real desktop and mobile tests must establish installed
behavior, notification opt-in and the link back to the exact question. These are
future acceptance criteria, not claims about the current CLI pilot.

## Confirmed project boundary — T026

The owner explicitly separates the hub and web application from SDP/SDPTool. Local
execution, report persistence, attention requests/replies and the eventual transport
adapter belong here. Hub storage, authentication, multi-device aggregation, web/mobile
views and push delivery belong to the external project. The architecture above is
a handoff proposal for those components, not authorization to build them in SDP.
No new repository, service identity, deployment or public network protocol is selected.
