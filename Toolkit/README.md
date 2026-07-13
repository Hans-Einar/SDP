# SDP Toolkit

Toolkit-Version: 0.2.0 (unreleased)

This directory contains the reusable Toolkit boundary:

- `SDP-install.manifest.json` — authoritative installable inventory and policy
- `schemas/` — installation, plan, manifest, traceability, record and build contracts
- `scripts/` — supported PowerShell installer, build metadata and validator
- `tests/` — deterministic Python and PowerShell fixtures
- `payload/` — copied Toolkit-managed files only
- `project-templates/` — neutral files that become project-owned when created
- `skills/` — versioned Toolkit-managed Codex skills

`SDP.manifest.yaml` at repository root is authoritative for Toolkit release and
capability facts. It is not the installation inventory. The JSON installation
manifest explicitly lists every copied or generated target so PowerShell and
independent clients such as `gh-sdp` do not reconstruct hidden behavior.

Produce a portable, mutation-free plan:

```powershell
.\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\path\to\Project `
  -PlanJson
```

The normal human preview remains available with `-Preview`. A project may
already contain a non-empty `SDP/` directory. Project-owned content is preserved
under normal, forced and repeated installation; changed managed files are backed
up before replacement. Unsupported schemas and downgrades block before mutation.

Validate the Toolkit or an installed consuming project explicitly:

```powershell
python Toolkit\scripts\validate_sdp.py --mode toolkit
python Toolkit\scripts\validate_sdp.py --mode project --project-root C:\path\to\Project
```

See `docs/Installation-Contract.md`, `docs/Installer-Migration.md` and
`docs/Validation.md`.
