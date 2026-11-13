# One SDPTool installation and update path

| Field | Value |
| --- | --- |
| id | KB-SDP-063 |
| project | SDP |
| type | Change |
| CardState | in-progress |
| created | 2026-11-13 |
| Systems | SDPTOOL, gh-sdp, XFMD |
| source | Owner rejects XFMD-private SDPTool distribution and requests one gh-extension-owned update path |

Owner constraints: distribute/update SDP tooling through gh-sdp, expose SDP
operations through `gh sdp`, and let XFMD invoke standalone `sdptool` offline
without gh execution or GitHub authentication. Avoid independent installed engine
versions and competing update mechanisms. The existing signed XFMD0.10 package
remains the historical delivered artifact until a selected migration replaces it.

Analysis: a binary GitHub CLI extension is downloaded, not executed as an installer.
Recommend shipping the canonical SDPTool engine in that executable and providing a
stable standalone entrypoint to the same file. A one-time local registration step
is needed for that public entrypoint; later extension upgrades update the engine
without a second engine download/cache/pin. Details and migration are a design
proposal; no product implementation or replacement publication is selected here.

[Activity mandate](../../Activities/Activity-0004--Single-SDPTool-installation/01-Mandate/Mandate.md)
and [initial study](../../Activities/Activity-0004--Single-SDPTool-installation/02-Study/Study.md)
record the planned need and proposed architecture. SESSION-SDP-0020 completes the
bounded initial study; product implementation is not started.

Owner selects execution through PLAN-SDP-0029 and SESSION-SDP-0021; prior analysis-only wording is historical.
