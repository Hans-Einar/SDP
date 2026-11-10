# Session0014 — SDP2.2.1 equivalent to the local producer

```mermaid
gantt
    title Sequence only — release preparation and publication
    dateFormat YYYY-MM-DD
    Prepare Linux trust and revised candidate :active, a, 2026-11-10, 1d
    Publish and verify distribution :b, after a, 1d
```

| Step | State | Outcome |
| --- | --- | --- |
| R221-K1 | completed | Independent Linux signer, overlapping trust, tested bridge candidate and scoped review |
| R221-M1 | active | Final revised-candidate release checklist and predecessor-upgrade evidence |
| R221-M2 | planned | Signed publication and verified explicit gh-sdp selection |

Goal: publish the owner-selected2.2.1 corresponding to local34ee535, with truthful
remaining PG limitations. Direct authorized Maintenance assignment:
[MAINT-SDP-0019](../Maintenance/R221/Plan.md), no duplicate feature card.

T001, 2026-11-10, manual work summary; host turn ID unknown. Release, Versioning,
Planning and SDP entry skills loaded. Existing clean checkout reused; no new
worktree. Runtime implementation remains unchanged. Latest public version2.2.0
confirmed and original signed descriptor downloaded; predecessor digest added.
Publisher private key not located locally; owner asked for location only, not
secret contents. Preparation continues. Public identities remain unset until
observed. No ProjectGovernance completion or project-directory upgrade inferred.

T002, 2026-11-10, manual summary. Local preparation reviewed at c337712;
[Evidence](../Maintenance/R221/Evidence.md),
[Review](../Maintenance/R221/Review.md) and exact package identity retained.
All local suites and 22 actual viewpoint requests pass. Six signed predecessor
installations preserve owner files through unsigned-target rehearsal and repeat
no-op. Corrected a predecessor test expectation; runtime remains unchanged.
Portable static amd64 v1 package replaces the narrower initial build. PR57 is
pushed; exact CI remains pending. Owner is locating the separate publisher key,
possibly on an offline laptop. No publication, main merge or project upgrade is
claimed. Next: recover key, finish CI, sign/rehearse/publish/verify distribution.


T003, 2026-11-11, manual work summary (host turn ID unavailable). Owner requires
an independent Linux signing key, since the laptop is only intermittently online.
Read SDP entry/document workflow, Architect, Planning, Master and Worker skills,
local Git/Session/management rules and installation trust contract.
Preflight and bounded refinement: [Linux signing](../Maintenance/R221/Linux-Signing.md).
Generated Linux-only key; existing laptop key untouched. Added overlapping trust
and a public proof regression. R221's historical candidate/review no longer
describe the changed bootstrap; rebuild and release checks remain required.
Next: verify trust and prepare a pinned gh-sdp bridge; public release remains distinct.


T003 closeout work summary: R221-K1 completed with clean producer643ec77 and
clientb3258f9, public dependency pin and portable amd64/v1 package. Exact signed
candidate, seven operational checks and packaged tests pass; independent reviewer
approves bounded preparation. Linux key remains only on Linux, laptop key untouched.
See Linux-Candidate.json, Linux-Signing.md and Linux-Review.md. The original R221
release review is not reused for the trust change. No installed client replacement,
published release, main merge or real project upgrade. Next: R221-M1 final release
gates and bridge publication coordination, then R221-M2; full publication remains
distinct from this completed signing prerequisite. Traceability routine also loaded.
