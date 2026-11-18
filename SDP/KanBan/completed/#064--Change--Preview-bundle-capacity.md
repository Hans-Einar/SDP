# Raise generated preview bundle capacity

| Field | Value |
| --- | --- |
| id | KB-SDP-064 |
| project | SDP |
| type | Change |
| CardState | completed |
| created | 2026-11-17 |
| Systems | SDPTOOL |
| source | Owner reports the bundle limit and requests at least four times the capacity |

The owner encounters `SDP: limit: bundle exceeds 128 files or 32 MiB` after the
recent SDP/SDPTool/gh-sdp releases. Increase generated preview and source-reading
publication limits to 512 files / 128 MiB, retaining bounded output and existing
source/selection validation. [Plan](../../Maintenance/BundleCapacity/Plan.md).

## Worklog

2026-11-17: registered and selected for immediate bounded implementation. The
error originates in SDPTool publication, not SDL parsing or gh-sdp transport.

2026-11-17: implemented shared 512-file/128-MiB bounds; focused tests pass.
[Verification](../../Maintenance/BundleCapacity/Verification.md). Public distribution
is not claimed; the installed signed release is unchanged.

2026-11-17: subsequent owner-authorized publication completed as signed SDPTool
and gh-sdp 2.3.2. Public downloads and isolated ordinary extension upgrade pass.
[R232 evidence](../../Maintenance/R232/Verification.md). The owner's live extension
remains unchanged; no remaining work is selected for this card.
