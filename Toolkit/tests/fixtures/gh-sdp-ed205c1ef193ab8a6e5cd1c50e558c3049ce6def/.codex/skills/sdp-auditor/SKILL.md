---
skillId: sdp-auditor
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.audit.consistency
- sdp.audit.release-state
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Auditor

Use this skill for a read-only consistency audit.

1. Parse Toolkit, installed and project manifests plus all skill metadata.
2. Check Toolkit/Framework/AGENTS/skill versions and compatibility ranges agree.
3. Check CurrentIndex, Relations, Ledger, release records, notes, reviews and
   verification for stale, missing or contradictory state.
4. Detect forgotten active work, incomplete included Slices/Fixes and unsupported
   schemas.
5. Detect claims of tags, releases, commits or verification without evidence.
6. Check released note sections against the immutable baseline.
7. Report findings by severity with exact records and remediation direction.

This initial skill defines the executable audit contract; SDP-Analyzer may later
implement richer navigation. The auditor never fabricates evidence or mutates
state while auditing.
