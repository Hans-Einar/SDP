# Legacy SDP Framework Payload

This root-level `payload/` tree is a legacy duplicate and is explicitly excluded
by `Toolkit/SDP-install.manifest.json`. It is not a canonical installation source
and must not be copied by PowerShell, `gh-sdp` or another conforming client.

Current Toolkit-managed sources live under `Toolkit/payload/` and
`Toolkit/skills/`. Neutral project-owned creation sources live under
`Toolkit/project-templates/`. Project-specific rules belong in the consuming
project's `SDP/Instructions/` or `SDP/AGENT-REMINDERS.md`.
