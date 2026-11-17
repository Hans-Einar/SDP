# Generated preview capacity correction

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0021 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

Owner authority, 2026-11-17: investigate the reported 128-file/32-MiB error and
increase the limit at least fourfold. Execute BC1 on sdp/bundle-capacity from
23c5193 (published 2.3.1 records), preserving the persistent management lock.
No merge, public release or replacement of the installed signed engine is selected.

Analysis: preview.go and source_reading.go independently cap sealed generated
bundles. This bounds published output resources, including retained source copies;
it is not an SDL language constraint. A valid 128-file composed source already
needs additional manifest/presentation files. No measured rationale for the exact
old values has been found; the contract records a conservative facade limit.
Use one shared 512-file/128-MiB check, reporting actual counts on failure.
Parameter tuning preserves workflow, ownership, model facts, wire format and SDL
semantics; no SDL structural design change or WORK revision is required.

| Milestone | State | Acceptance |
| --- | --- | --- |
| BC1 | completed | Preview and source reading accept >128 generated files; 512 files / 128 MiB inclusive bounds and overflow rejection tested; contract updated |

Preflight: sole working checkout is clean apart from the persistent ledger lock.
Detached /tmp/sdp-activity4-release at 774e150 and /tmp/sdp-r231-source at 84d5e46
are clean completed-release snapshots; no concurrent local write overlap. Current
Activity catalogue has no active Activity. KB064/MAINT0021 are unused. Existing
PG/Blueprint work remains outside this fix. Remote writers are outside observation.
Green for preview/source-reading capacity helper, focused tests, Contract.md and
linked management records. Recheck if competing edits appear.

Verification: boundary checks and one composed-source regression through both
publication paths, plus existing preview/source-reading tests. No new agent,
separate review or broad/native test matrix for this parameter correction.

BC1 is implemented and verified; [evidence](Verification.md) pins the candidate.
A future authorized release must carry the change to installed SDPTool/gh-sdp.
