# SDP 2.3.1 release and HEOS rollout

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0020 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

Owner requests a new installable SDP release after main integration, with HEOS
as the first project upgraded. This selects production publication of the canonical
engine and the same-byte gh extension, local extension update/registration and
HEOS Framework preview/rehearsal/apply. No other project's Framework is selected.
Use sdp/release-2.3.1 from integrated main cef86e4. Product behavior is unchanged
from2.3.0; patch2.3.1 records the integrated documentation/model correction and adds
signed2.3.0 to the supported predecessor inventory. Language/wire schemas and
Framework payload ownership stay unchanged. No XFMD release is required.

| Milestone | State | Acceptance |
| --- | --- | --- |
| R231-M1 | active | Clean exact signed package, supported predecessor matrix, HEOS copy rehearsal, independent review and applicable CI |
| R231-M2 | planned | Public canonical and extension bytes/signatures verified; local tool update and actual HEOS upgrade preserve owner files |

Preflight: SDP is clean except its persistent management lock; detached Activity4
release worktree is unchanged. gh-sdp is clean on its completed delivery branch.
HEOS main has existing dirty browsing study, card, Session, overview, README and
management history; preserve their bytes, do not stash or commit them. Snapshot
before rehearsal and recheck before live apply. Verdict green for isolated release
metadata and the installer-owned paths, conditional on a conflict-free live preview.

Reuse2.3.0 source/native qualification for unchanged implementation. Run profile,
installation and bootstrap tests for changed release selection, all signed supported
predecessor upgrades with conflicts/preservation/no-op, exact package/source checks,
public install/update and explicit real HEOS upgrade. Existing exact CI remains
historical; query the new candidate's actual CI before publication. Independent
release review is required. No private key enters evidence. Historical releases,
ledger bytes and unrelated dirty HEOS work remain intact.
