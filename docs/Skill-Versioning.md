# Skill Versioning

Every canonical `Toolkit/skills/*/SKILL.md` begins with YAML front matter:

```yaml
---
skillId: sdp-versioning
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
  - sdp.versioning.select
compatibilityNotes: Initial formal version.
---
```

Skill IDs are stable. Skill versions use SemVer independently from the Toolkit.
The Toolkit manifest lists the exact distributed skill versions. Validation
fails when front matter and the Toolkit manifest disagree.

Increment MAJOR for incompatible procedure/contract changes, MINOR for new
backward-compatible workflow capability and PATCH for clarification or correction.
Do not change skill metadata in installed consuming-project copies; update the
canonical Toolkit skill and reinstall.
