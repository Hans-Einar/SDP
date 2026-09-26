---
skillId: sdp-reviewer
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.review.independent
- sdp.release.review
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Reviewer

Use this skill in a fresh context independent of implementation.

1. Read the Slice/Fix contract, release target and authoritative requirements,
   architecture and design before reading implementation summaries.
2. Inspect the actual diff and affected surrounding files.
3. Check correctness, regressions, scope discipline, compatibility, state
   ownership, error paths and security implications.
4. Inspect real verification evidence and rerun checks when practical.
5. Check manifests, release notes, traceability, implementation notes and handoff.
6. For release work, verify that SemVer is separate from development identity,
   released history is immutable and publication claims are truthful.
7. Report findings by blocking, high, medium, low or note with exact locations.
8. State approved, changes required or blocked.

Do not fix implementation from the Reviewer pass unless assigned a separate
remediation Slice.
