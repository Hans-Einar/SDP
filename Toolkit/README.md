# SDP Toolkit

Toolkit-Version: 0.2.0 (unreleased)

This directory contains reusable, Toolkit-managed assets:

- `skills/` — versioned Codex skills
- `schemas/` — manifest, release-event and build-identity contracts
- `scripts/` — safe installation, migration, build metadata and validation
- `tests/` — deterministic Python and PowerShell fixtures
- `payload/` — managed files and additive project templates

`SDP.manifest.yaml` is the authoritative release manifest. Installed projects
receive generated static facts in `SDP/Framework/installed-toolkit.manifest.yaml`
and a project-owned `SDP/SDP-project.manifest.yaml` when missing.

Preview installation:

```powershell
C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\path\to\Project `
  -Preview
```

The target may already contain a non-empty `SDP/` directory. Project-owned
content is preserved. Managed files are backed up before replacement. Unsupported
manifest schemas stop safely.
