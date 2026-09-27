# Installer migration

Use the compiled SDPTool, normally through `gh sdp`. The source-distributed legacy
installation scripts and profile-artifact builder are removed. Old release assets
remain historical; the current release does not execute their installation engines.

## New or registered projects

Run from the project root, or supply the project path before the operation:

```sh
gh sdp /path/to/project install --plan-output /tmp/sdp-install-plan.json --json
gh sdp /path/to/project install --apply /tmp/sdp-install-plan.json --json
gh sdp /path/to/project upgrade --plan-output /tmp/sdp-upgrade-plan.json --json
gh sdp /path/to/project upgrade --apply /tmp/sdp-upgrade-plan.json --json
```

Select a published release with `--release` when not using the client's compiled
default. Plan and output paths must be new and outside the target project as
required by the CLI. Inspect the actual actions/conflicts before apply. Apply
rechecks the exact source/target snapshot; a changed project requires a new plan.

## Manual or legacy installations

The Go engine reads earlier receipt formats, but migration requires a verified
original descriptor or explicit adoption manifest with the observed project
snapshot, preserving moves and managed refresh exceptions. Do not infer ownership
from file names. For an explicitly authored adoption input:

```sh
gh sdp /path/to/project upgrade --manifest adoption.yaml --plan-output /tmp/adoption-plan.json --json
gh sdp /path/to/project upgrade --apply /tmp/adoption-plan.json --json
```

A pending legacy operation blocks migration; do not delete its journal or silently
resume it with a different protocol. Recovery of such historical operations needs
a separate assessed migration/recovery task. Current Go operations use `--resume`.

See [record schemas and limits](../../SDPTool/install/Records.md). Project-owned
prose is preserved rather than overwritten by newer templates. Generated upgrade
Maintenance reports and receipt updates record the actual transaction.
