# R232 — Publish the preview capacity correction

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0022 |
| project | SDP |
| state | active |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | milestone |

Owner explicitly requests publication so `gh extension upgrade sdp` installs
BC1. Publish canonical SDPTool and identical gh-sdp engine as2.3.2, a compatible
patch to output capacity. Retain all supported predecessors and add the actual
signed2.3.1 descriptor. No main merge, real local extension replacement or consumer
Framework migration is selected; the owner will upgrade the extension.

Use sdp/release-2.3.2 from4d12b27 and the corresponding clean gh-sdp release branch.
The preflight is green: root has only the preserved management lock; both detached
release snapshots and gh-sdp54f830a are clean. New writes are release metadata,
bootstrap default/predecessor inventory, package provenance and lifecycle/evidence.
The single gh-sdp worker owns its repository; canonical signing and publication
remain here. Other PG/Blueprint work and unknown remote writers are outside scope.

| Milestone | State | Acceptance |
| --- | --- | --- |
| R232-M1 | active | Exact clean signed candidate, focused tests and release CI, upgrade/extracted-source checks, one combined independent release review |
| R232-M2 | planned | Canonical and extension publication downloaded and verified; isolated real gh upgrade from2.3.1 delivers the fixed same-byte engine |

Reuse BC1 targeted verification and existing evidence for unchanged native behavior.
Run changed profile/bootstrap and packaging checks, all supported signed predecessor
upgrades with preservation/conflict/no-op, and an isolated public install/upgrade.
One final combined release review suffices. No private signing material is retained
in evidence. Published notes/assets/history remain immutable.
