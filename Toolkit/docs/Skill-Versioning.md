# Skill metadata and versioning

Canonical sources are root Skills/. Native metadata v2 is compatible with the
host skill format and carries distribution facts as string-valued metadata:

```yaml
---
name: sdp-versioning
description: Select release versions from public compatibility impact.
metadata:
  skillId: sdp-versioning
  skillVersion: 2.0.0
  minimumToolkitVersion: 0.2.0
  capabilities: sdp.versioning.select
  compatibilityNotes: Profile-aware procedure replacing the legacy role.
---
```

Quote numeric-looking values when needed so every metadata value is a string.
Separate multiple capabilities with commas; the SDP validator normalizes this
field to its established capability list. Native name must equal skillId and
description must be non-empty. Do not mix v1 top-level fields and v2 metadata.
The validator retains legacy v1 support for existing installed skill metadata and
frozen consumer fixtures; current canonical skills must use native metadata.

Skill IDs remain stable, including the bare sdp router. Skill versions are SemVer
independent of Toolkit, process/profile, language and project product versions.
Increment MAJOR for incompatible procedures, MINOR for compatible capabilities,
PATCH for corrections. Ten replacement procedures use 2.0.0; newer focused roles use 1.0.0.
The routing skill is 1.1.0 with Planning, whose initial version is 1.0.0. The Toolkit and install generator declare exact versions and the v2 metadata
capability. AGENTS contract 2.1.0 retains routing and adds adopted per-plan Git policy.
These are source distribution facts in the unreleased Toolkit, not a publication.

Installed files remain under .codex/skills and come from root Skills/ through the
explicit inventory, including references. Do not edit those copies. A same-Toolkit-
version skill inventory/version change requires reviewing and applying the existing
ForceManagedFiles option; unforced application fails before mutation rather than
claiming new versions while retaining old skills. Changed managed files are backed
up, project-owned files are preserved, and an ordinary repeat is idempotent.
