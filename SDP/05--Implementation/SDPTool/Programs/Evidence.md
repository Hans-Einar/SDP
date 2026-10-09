# RSP1 evidence

Baseline: 04f88ff918e7c25c9fa1883f061ef20b8c3923a2, isolated sdp/runnable-programs.
RSP1-M1 implements explicit declaration, additive discovery and foreground launch.
Initial full SDPTool suite passed (root 7.333s, install 14.075s, model 6.212s).
Focused final program race tests passed (2.378s), including literal argv/root cwd,
non-executing discovery, invalid/duplicate/removed catalog, blocked source/entry/
command, stale selection and owned descendant cancellation. Full race and actual
gh/native candidate verification remain pending; no independent review claimed.
