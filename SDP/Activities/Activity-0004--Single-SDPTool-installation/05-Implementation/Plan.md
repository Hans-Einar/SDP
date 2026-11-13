# PLAN-SDP-0029 — One shared SDPTool distribution

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0029 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Activity | ACTIVITY-SDP-0004 |
| Source | Owner instruction to execute Activity4; KB-SDP-063 |

Deliver the selected single-update-path need across canonical SDPTool, gh-sdp
and XFMD. The owner selects execution of the studied migration. Carry through
qualified release and successor RPM delivery under the continuing distribution
assignment; main merges and unrelated live project upgrades remain unselected.

| Milestone | Outcome / acceptance | State |
| --- | --- | --- |
| C0 | Canonical design/model and bounded contracts; Activity/Session selected | complete |
| C1 | Shared standalone entrypoint, diagnostics and exact engine identity; existing project operations preserved | in-progress |
| C2 | gh extension payload is the canonical engine, no runtime bootstrap child/cache; actual extension install/upgrade checks | in-progress |
| C3 | XFMD uses registered local runtime without gh/network/auth; removes private engine and reconciles all SDP helpers | in-progress |
| C4 | Integrated offline/migration/native checks and fresh independent review on exact candidates | planned |
| C5 | Signed canonical/shared extension publication then successor XFMD RPM; actual public update/install and truthful closeout | planned |

Use sdp/shared-tool-installation in all three repositories. SDP starts432c292;
gh-sdp startsb177e14; isolated XFMD starts140be2b. Preserve primary XFMD main
working files and existing runtime lock/cache artifacts. Do not modify HEOS.

Canonical engine release target2.3.0 adds Activity/reading and setup capabilities.
The extension distribution adopts the same2.3.0 identity because its executable
is now that exact engine, replacing the earlier0.2.x bootstrap-client product.
XFMD0.11 selects the changed runtime-installation contract from0.10. These are
public compatibility/version choices, not commit counters or SDL schema changes.

Verification: actual extension update selects new engine even before a gh sdp
invocation; stable registered path; source/hash/version equality across entrypoints;
no gh/auth/network for local native browsing; conflicts and GUI PATH; interrupted
updates; canonical signed framework previews/applies and unchanged project bytes;
prior RPM upgrade removes private helper; missing/incompatible engine diagnostics;
whole helper inventory and preserved supported workflows. Linux amd64 is the
existing native release boundary; other platform support must be stated honestly.
