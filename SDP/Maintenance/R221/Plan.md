# SDP 2.2.1 — publish the local producer capability set

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0019 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| SessionId | SESSION-SDP-0014 |

Owner explicitly requests SDP2.2.1 equivalent to the local producer built from
34ee535fbf6257b73d2fa0c7fbe4f9d5131aaf5e. Prepare, verify, independently review and
publish the signed Linux amd64 release. The original scope preserved runtime implementation bytes;
change release metadata, notes, generated release log and predecessor inventory
only. Final release provenance differs from the local development build.
Use branch sdp/release-2.2.1 in the existing checkout; no new worktree.

Version2.2.1 is an explicit owner selection. Relative to public2.2.0 this includes
additive viewpoint/governance changes, so do not describe it as solely a bugfix
or infer complete governance acceptance. Existing PG obligations remain open.
No unrelated feature implementation or project upgrade. The 2026-11-11 owner
refinement below authorizes Linux key provisioning and bootstrap/client trust
preparation; it does not itself publish a new client.

Use the existing checklist and a single final independent release review. Check
exact source, complete product tests and actual published2.2.0 upgrade path,
preserved owner documents, signatures, downloaded asset identity and explicit
gh-sdp selection. Retain earlier supported predecessors. If the publisher key is
unavailable, finish a reviewable candidate but never substitute an untrusted key
or claim signed/public completion. Ask only for the existing key location.

| Milestone | State | Done criteria |
| --- | --- | --- |
| R221-M1 | active | Exact source/package, targeted/full release gates and independent review |
| R221-M2 | planned | Signed tag/release/download verification and recorded consumer selection |

Source checkout initially clean; only one SDP worktree exists. Previous PG branch
is preserved. No concurrent dirty work observed. Release scope adds no runtime
code changes to the local producer baseline. Publication is authorized by the
current owner instruction; update main through normal integration after validation
under the preceding consolidation direction, without force or unrelated edits.


## Owner refinement — Linux publisher, 2026-11-11

Linux must build and sign without the laptop being online. Add a second compiled
trust anchor, retaining the existing key and exact signature format. Provision
the new secret only on Linux; use a public signed probe for regression tests.
This supersedes the runtime-byte preservation constraint only for bootstrap
trust. Existing candidate c337712, package and review remain historical evidence;
the release needs a new clean-source package and final review after this edit.

R221-K1 (completed; see Linux-Signing.md and Linux-Review.md): Linux key, overlapping trust, regression proof, operator guide
and gh-sdp bridge candidate pinned to the new immutable bootstrap. Keep the
2.2.0 default until a separately checked 2.2.1 release is available. Complete
independent review and tests, without claiming published client adoption.
BranchPolicy remains current (sdp/release-2.2.1), CommitPolicy milestone.
The client uses its own working branch and record. R221-M2 remains pending.
