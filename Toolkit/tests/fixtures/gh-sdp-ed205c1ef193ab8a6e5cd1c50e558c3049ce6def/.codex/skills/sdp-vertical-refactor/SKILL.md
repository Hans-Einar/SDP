---
skillId: sdp-vertical-refactor
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.refactor.vertical
- sdp.compatibility.preserve
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Vertical Refactor

Use this skill for behavior-preserving architectural refactors.

1. Establish a verified baseline before product-code changes.
2. Define target architecture, compatibility range and measurable exit criteria.
3. Refactor one complete vertical workflow at a time.
4. Preserve a runnable application after every Slice.
5. Move behavior behind explicit contracts before deleting compatibility paths.
6. Keep state ownership singular and adapters explicitly scheduled for removal.
7. Verify typecheck, tests, build and rendered/manual behavior after each Slice.
8. Use a fresh Reviewer and stop at the Slice boundary.

A refactor does not independently determine a SemVer increment; compatibility
impact does. Never split files merely to reduce line count.
