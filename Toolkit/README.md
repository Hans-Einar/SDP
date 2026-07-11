# SDP Toolkit

This directory contains reusable, toolkit-managed assets:

- `skills/` - reusable Codex skills
- `scripts/` - installation and update tooling
- `payload/` - files copied into consuming projects

The numbered folders and lifecycle directories at the repository root are the
reference project structure. They are templates and documentation, not the
active SDP records for this toolkit repository.

Install into an existing project with:

```powershell
C:\Users\hanse\GIT\SDP\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\Users\hanse\GIT\GrassPhenology
```

The target project may already contain a non-empty `SDP/` directory. Existing
project-owned documents are preserved.