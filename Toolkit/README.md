# SDP Toolkit

Toolkit-Version: 0.2.0 (unreleased)

This directory contains the reusable Toolkit boundary:

- [SDPTool](SDPTool/README.md) — source home and plan for the project-aware sdptool facade; Go facade delivered; see its current scope
- `profiles/five-phase.json` — authored current process configuration and reproducible artifact
- `SDP-install.manifest.json` — authoritative legacy install-v1 inventory and policy
- `schemas/` — installation, plan, manifest, traceability, record and build contracts
- `scripts/` — supported PowerShell installer, build metadata and validator
- `tests/` — deterministic Python and PowerShell fixtures
- `conformance/install-v1/` — language-neutral scenario and expected-outcome contract
- `payload/` — copied Toolkit-managed files only
- [Template](../Template/README.md) — neutral files that become project-owned when created
- [Skills](../Skills/README.md) — canonical versioned skills and shared references

`SDP.manifest.yaml` at repository root is authoritative for Toolkit release and
capability facts. It is not the installation inventory. The JSON installation
manifest explicitly lists every copied or generated target so PowerShell and
independent clients such as `gh-sdp` do not reconstruct hidden behavior.

For the current five-phase/shared-management profile, follow
[versioned process installation](docs/Process-Installation.md): build the artifact,
inspect -ProfileArtifact/-PlanJson, then consume that exact plan with -ApplyPlan.
Recovery, installed facts 2.0 and target Maintenance reports belong to that
contract. [V2 conformance](conformance/install-v2/README.md) includes interrupted
execution and the prebuilt SDPTool consumer.

The following retained interface is install-v1. Produce its mutation-free plan:

```powershell
.\Toolkit\scripts\Install-SDP.ps1 `
  -ProjectRoot C:\path\to\Project `
  -PlanJson
```

The normal human preview remains available with `-Preview`. A project may
already contain a non-empty `SDP/` directory. Project-owned content is preserved
under normal, forced and repeated installation; changed managed files are backed
up before replacement. Unsupported schemas and downgrades block before mutation.
Portable-path, case-collision and link/reparse containment checks run before
planning and mutation. The emitted plan carries exact `source`, `destination` and
`targetSource` facts and is validated semantically before apply. Its required
`migration-first-manifest-order-v1` policy puts target migrations first,
preserves manifest-array order and keeps backup/mutation pairs adjacent. AGENTS
migrations additionally carry an exact-byte SHA-256 and absent-destination
precondition.

Validate the Toolkit or an installed consuming project explicitly:

```powershell
python Toolkit\scripts\validate_sdp.py --mode toolkit
python Toolkit\scripts\validate_sdp.py --mode project --project-root C:\path\to\Project
python Toolkit\conformance\install-v1\run_conformance.py --validate-only
python Toolkit\conformance\install-v1\run_conformance.py --powershell powershell
```

Python regression fixtures include the supported `gh-sdp` closure offline, so
cross-project schema and traceability compatibility do not depend on network
availability. Independent clients can consume the conformance scenario index
and committed expected plans/failure classes without executing PowerShell.
Normal tests do not regenerate those authorities; candidate regeneration is an
explicit maintainer-only `--write-candidates` operation followed by diff review.

See `Toolkit/docs/Installation-Contract.md`, `Toolkit/docs/Installer-Migration.md` and
`Toolkit/docs/Validation.md`.
