# Installer Migration Guide

Run a preview first:

```powershell
C:/Users/hanse/GIT/SDP/Toolkit/scripts/Install-SDP.ps1 `
  -ProjectRoot C:/path/to/Project `
  -Preview
```

Apply the update:

```powershell
C:/Users/hanse/GIT/SDP/Toolkit/scripts/Install-SDP.ps1 `
  -ProjectRoot C:/path/to/Project
```

Use `-ForceManagedFiles` only when intentionally restoring Toolkit-managed
Framework or skill files that were locally modified. Project-owned release notes,
project manifest, AGENTS-project, lifecycle documents and populated traceability
are never replaced.

The installer:

- accepts empty and non-empty SDP directories
- detects the installed manifest schema before writing
- treats installations without a manifest as the supported pre-versioning format
- stops safely on unsupported schemas
- creates missing project manifest and release notes additively
- refreshes managed files on a Toolkit-version upgrade
- backs up changed managed files before replacement
- reports previewed and applied actions
- preserves the original timestamp when the same Toolkit version is reinstalled
- compares complete path segments, so sibling folders such as `SDP-Analyzer` are valid

Backups default to `SDP/.sdp-backups/<UTC timestamp>/` and may be redirected with
`-BackupRoot`.
