Need: eliminate independently updated SDPTool installations after the owner rejects bundling SDPTool inside XFMD. The normal installation/update route is gh extension install/upgrade sdp. All SDP tool operations must be reachable through gh sdp. XFMD must execute sdptool locally, without invoking gh, requiring GitHub authentication or needing a network for local browsing.

Outcome: one authoritative installed engine/version per user for the CLI and native consumers; coherent tool-suite distribution, explicit compatibility and migration from old PATH/cached/private installations. Preserve user projects/settings and retain published release history.

Current assignment: read-only architecture analysis and need capture. Proposed one-time registration, binary packaging shape and migration acceptance still need detailed design. This planned Activity does not start product implementation, main integration or replacement publication.


Execution selection: the owner now instructs execution of Activity4. The prior
analysis-only boundary is superseded for this need. Complete design, implementation,
qualification and the studied shared-runtime/successor-package delivery. Preserve
separate main-merge and unrelated live project-upgrade boundaries.
