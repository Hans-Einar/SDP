# Sessions release preparation evidence

## Installed consumer baseline

Read-only invocation from XFMD: `gh sdp --version --json` reports SDPTool 2.0.0,
schema sdptool/0.2, source d304261c90066a86b2d8ffcaa2115517ea339b05.
`gh sdp . discover --json` returns roots files/sdl/kanban/sdui. Resolving IDs
through navigation.nodes yields directory/tab/tab/tab. This confirms XFMD can
adopt dynamic tabs before the Sessions-enabled tool is published. No live project
files or globally installed extension were changed.

## Preparation boundary

Product implementation 2419853 is tested under PLAN-SDP-0015. SR1 freezes additive
2.1.0 notes and generated log, retains Framework 2.0.0 and the identical payload
inventory, adds the original signed 2.0.0 descriptor as a supported predecessor,
and selects the future 2.1.0 bootstrap URL. That URL is not yet a published asset.
No tag/publication or consumer installation is claimed. Exact candidate evidence
will be appended after the clean package, independent review and CI checks.

Independent review identified an inherited duplicate-node path when WalkDir
reports an unreadable directory twice. SR1 corrects that within the accepted
unavailable/unique-node contract and adds real permission regressions for Sessions
root/nested directories. The predecessor expectation adds the exact original
2.0.0 digest; all previous digests remain checked.

The corrected full SDPTool Go suite passes. Candidate packaging and publication
checks remain separate; this commit establishes the source for those checks.
