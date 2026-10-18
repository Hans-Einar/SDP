# ProjectGovernance attention mailbox and remote coordination

| Field | Value |
| --- | --- |
| id | KB-SDP-054 |
| project | SDP |
| type | Proposal |
| CardState | backlog |
| Systems | SDPTOOL |
| created | 2026-10-18T22:33:35.255668+00:00 |
| source | Owner Session0007 T025: autonomous work, call attention and self-hosted multi-machine oversight |
| next_review | After KB052 PGL3 local report/inbox contract; select a bounded hub/client plan |
| tags | project-governance, autonomy, attention, mailbox, remote-control, web, mobile, pwa |

## Requested outcome and ownership

The owner wants projects on different PCs to work unattended for hours, while a
shared service calls for attention when input is actually needed. Provide a
self-hosted project overview and messaging/control interface on desktop and mobile,
with an installable notification-capable client. SDP_Analyzer is obsolete and is
not the baseline. ProjectGovernance / Session0007 owns this registration.

[Architecture proposal](../../03--Architecture/ProjectGovernance/Attention-and-Remote-Coordination.md)
records the workflow, responsibility split, alternatives, platform evidence and
remaining design choices. The working name PG Hub is provisional. No new runtime,
public protocol, hosting or deployment is selected by card registration alone.

## Relationship to existing work

[KB052](../active/%23052--Change--Project-Leader-and-SAD-orchestration.md) and
[PLAN-SDP-0018](../../05--Implementation/SDPTool/ProjectGovernance/Plan.md) keep the
local leader, report/inbox, authority and continuation implementation. The owner's
choice of automatic reporting applies there immediately as product direction.
This card owns the extension across machines and its web/mobile delivery. Select
its execution plan explicitly; do not silently make the entire remote service a
new completion requirement for the existing local Session milestones.

[KB038](../active/%23038--Proposal--Codex-app-server-development-client.md) remains
the client/runtime foundation. [KB043](%23043--Proposal--Event-derived-session-timelines.md)
retains the broader event-derived timeline. No competing project-management ledger
or revival of SDP_Analyzer is implied. The architecture identifies proposed service/
client responsibilities; Systems SDPTOOL records the existing integration owner,
not a claim that every future component must be linked into that binary.

## Acceptance criteria for the follow-on plan

- Two machines and projects appear with truthful last-observed/connected state;
  closing the UI or losing the hub connection does not stop independent authorized work.
- A correlated attention question survives restart, duplicate/out-of-order delivery
  and missed push. It supports options and free text, shows blocking scope and evidence,
  and distinguishes read/acknowledged from answered/applied.
- Authenticated owner replies and targeted control route to the correct checkout/
  assignment revision. Stale/revoked/wrong-project commands refuse; uncertain sends
  do not repeat work. Offline control is visibly queued until actually acknowledged.
- Status/result reporting is automatic within granted authority. Scope changes and
  reserved decisions still require current authorization; one waiting task does not
  unnecessarily stop the other projects.
- Desktop/mobile web works; actual supported-device tests establish PWA install,
  notification consent and opening the exact question. Inbox works without push.
- A declared multi-hour unattended trial demonstrates scheduling, recovery, bounded
  retries/resources and continuation after an attention response. Host/runtime
  limitations and incomplete evidence remain visible.
- Optional ChatGPT MCP Events is tested separately for account/client/protocol
  availability and notification/reply behavior; it is not a delivery prerequisite.

## Worklog

2026-10-19 / T025: Registered after scanning card/event identities across thirteen
local worktrees (highest prior SDP card 53). Recorded owner automation decision and
proposed local-runtime / self-hosted-hub / responsive-PWA boundaries. Official WebKit,
MDN and OpenAI sources establish platform possibilities, not tested PG integration.
Existing KB052 execution continues; implementation of this extension is not started.
