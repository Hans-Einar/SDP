# MPV1-M1 source relocation evidence

Baseline: 17c5dfb37d2e8c25df9c697aa3722f4d90db73e7, before this work.
[migration.json](migration.json) records each old path and its destination(s),
the baseline statement count and SHA-256 of sorted statements with multiplicity.
Only source-membership sentences, headers, comments and whitespace are excluded.
All 4,523 declarations/facts are unchanged. UI/Views.design splits into operator,
simulator and common definitions; 66 files become 68. Model identities, containment,
owners, contracts, scenario ordering and quoted obligations are retained.

Run from repository root:

```sh
python3 experiments/mvp1_sdl/verify_restructure.py
python3 experiments/mvp1_sdl/audit_inventory.py
python3 SDP/ProjectManagement/validate.py
```

[verification.json](verification.json) records the passing current inventory and
statement comparison. No external Ponsse source-hash recheck was performed; the
pinned source-inventory.json and original inventory-audit.json are unchanged.
This is an authoring inventory/preservation check, not grammar/type validation,
execution, implementation coverage or independent review.

The maintained README/coverage CSV now point to current source locations. Shared
responsibility guidance distinguishes MVP1 composition from reusable candidates.
Broader SDL source-set support and production migration remain KB-SDL-005 and
KB-SDP-020. The next phase creates static current-profile SDUI designs.
