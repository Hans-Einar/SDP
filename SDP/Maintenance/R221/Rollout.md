# R221 publication and deployment

SDP2.2.1 and gh-sdp0.2.3 are published from the exact reviewed source commits.
Publication.json records producer tag, assets and download hashes. The client
record is SDP/Sprints/Sprint-010/Publication.json in Hans-Einar/gh-sdp.
The public descriptor was verified through the exact client with a fresh cache
before publishing that client. Exact producer CI42925278417 passed.

Host-Installations.json records Linux, Beinerud and labtop using client0.2.3
with the reviewed binary hash and public signed engine2.2.1. Previous extensions
are retained under ~/.local/state/sdp-rollout-2.2.1/gh-sdp-before-0.2.3 on each host.

## Consumer results

Fleet-Results.json records the Linux inventory. kitty, sdui_widget_lab and xfmd
are upgraded to signed2.2.1. Each operation used a saved public-default preview,
verified the final receipt and Git index, preserved the management-ledger byte
prefix and repeated as a no-op. Owner file hashes were checked against backups.
XFMD advanced through concurrent RPM work: its backup was refreshed before
applying; two packaging files changed during verification outside the installer
plan. Those edits were left intact and recorded, not claimed byte-identical.
No consumer edits were committed; pre-existing dirty work remains the owner's.

Four signed2.2.0 installations remain blocked by existing invalid history:

| Repository | Existing blocker |
| --- | --- |
| HEOS | Ledger line40 EVT-KB-HEOS-000021 uses to=in-progress with active/ path (also line47) |
| emuSA80535-N | Ledger line7 EVT-KB-EMU-000004 has the same state/path mismatch |
| ponsse | Ledger line43 EVT-PM-PONSSE-000006 uses unsupported x-management:activated (also a later state/path mismatch) |
| ponsse-service-continuity | Shared historical unsupported x-management:activated event |

The canonical ledger uses folder states for from/to, distinct from CardState.
No immutable historical events were rewritten and no validation bypass was used.
A history-preserving reconciliation strategy is needed before those upgrades.
The gh-sdp repository's legacy schema1 manual Toolkit is excluded from this
gh-sdp-installed-repository rollout; it needs explicit adoption separately.

R221-M1/M2 are complete. R221-M3 is partially complete with four blocked consumers;
the overall plan stays active. No main merge, application deployment or completion
of the separate ProjectGovernance design obligations is claimed.
