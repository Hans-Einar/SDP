# GIP — implement SDPTool installation and thin gh-sdp delegation

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0003 |
| project | SDP |
| state | active |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDP |
| source | KB-SDP-033; IPD-3 handoff from PLAN-SDP-0002 |

## Outcome and selection boundary

Implement a single Go installation engine inside root SDPTool and a thin gh-sdp
client that obtains/invokes it. Deliver the reviewed end-to-end upgrade of a
**disposable** XFMD copy first. Direct sdptool and gh sdp use the same policy and
produce equivalent plans/outcomes. No PowerShell is required by the new user path.

This is the implementation handoff from the completed
[IPD DesignPlan](../../../04--Design/SDPTool/Installation/Plan.md), not an instruction
to execute production work during the design assignment. The owner selected execution of GIP on 2026-09-27 after IPD delivery. No live rollout, merge, key
publication, signed production release or support claim is implicit.

Inputs: [working contract](../../../04--Design/SDPTool/Installation/Contract.md),
[acceptance cases](../../../04--Design/SDPTool/Installation/Scenarios.md),
[generated review](../../../04--Design/SDPTool/Installation/review/index.md),
[canonical model](../../../03--Architecture/SDPTool.design),
[REQ-SDPTOOL-007](../../../02--Requirements/SDPTool.md),
[KB-SDP-033](../../../KanBan/active/%23033--Study--XFMD-SDP-adoption-and-SDL-pilot.md).
This installation-specific plan supplies concrete milestones to the existing
[feature plan](../../SDPTool.md); it does not rewrite T0–T5 history or restart BP2.

## Implementation rules

Use Go packages under SDPTool for selection, release/adoption records, inspection,
planning, execution and journal/receipt handling. Map them to existing SDL units;
packages are not extra containers. Keep current navigation/preview commands usable.
Do not copy policy into gh-sdp, or create another standalone Go Toolkit executable.
Extract shared deterministic fixtures from retained Toolkit inputs. Preserve frozen
legacy fixtures and historical bytes; intentional behavior differences get explicit
fixtures/rationale rather than edited golden results to hide drift.

The public contract names in IPD are draft version identities. GIP-1 freezes exact
JSON/YAML schemas, bounded readers, canonical serialization and status/exit semantics.
It also specifies the new installed-receipt schema and upgrades discovery readers
before writers. Do not label old facts 2.0 with unrecognized provenance keys or
pretend today's PowerShell prerequisite is already engine-neutral.

## Phases and milestones

| Phase | Milestone | Runnable outcome and acceptance | Design / cases |
| --- | --- | --- | --- |
| GIP-1 — read-only planning | GIP-1-M1 | Freeze executable schemas and shared golden fixtures for release, adoption, receipt, plan and journal. Test required/unknown/duplicate keys, limits, path roots, identity and error codes. Add installation-specific project eligibility without breaking navigation discovery. | InspectInstallationBaseline; IC01/02/03/06/09 |
| GIP-1 — read-only planning | GIP-1-M2 | Direct sdptool install/upgrade preview works from a verified local fixture or explicit development artifact. Deterministic exact plans identify writes/deletes/preservation/conflicts; no project mutation. Known upgrade and manual adoption are distinguishable. Preserve current discovery/preview output contracts. | BuildInstallationPlan; IC01–06/08/10 |
| GIP-2 — bounded local execution | GIP-2-M1 | Apply a saved plan on temporary projects with locks, confinement, exact input binding, backups and journaled per-file replacement. Clean install and known upgrade run end to end; installed reader recognizes completion. No live consumer target. | ApplyInstallationPlan; IC01/02/04/05/06/12 |
| GIP-2 — bounded local execution | GIP-2-M2 | Resume every interrupted step/finalization boundary; preserve reserved bytes/IDs, history prefixes and truthful incomplete state. Post-failure edits block recovery. Repetition is no-change; old pending journals cannot be silently consumed. | InstallationRecoverySlice; IC07/08/11/12 |
| GIP-3 — distribution and client | GIP-3-M1 | Verify signed test descriptors and platform/engine/protocol compatibility, populate digest-keyed cache, exercise offline/corrupt/wrong-key inputs and fixed local child invocation. Freeze a release-neutral profile with explicit provenance/capabilities. Test keys remain clearly non-production. | ResolveInstallationRelease, RetrieveSdpToolBinary; IC09/10 |
| GIP-3 — distribution and client | GIP-3-M2 | In the gh-sdp repository, implement only bootstrap/delegation. `gh sdp upgrade --manifest ... --plan-output ...` and explicit --apply match direct SDPTool behavior; forward streams/exits and reject incompatible binaries. Update gh-sdp's old authority/model allocation without changing its history. | GhSdpClientSlice; IC01/03/09/10 |
| GIP-4 — XFMD adoption and retirement readiness | GIP-4-M1 | Refresh an exact XFMD baseline, author a reviewed adoption manifest, then preview/apply/resume/repeat on a disposable copy through both entry points. Prove original live bytes/status unchanged, project documents/root instructions preserved and KanBan/history migrated exactly. | ManualAdoptionTrial; IC03–08/11 |
| GIP-4 — XFMD adoption and retirement readiness | GIP-4-M2 | Deliver exact-candidate evidence, documented remaining platform gaps, compatibility/retirement map and a concrete live-upgrade/release proposal. Remove legacy execution only where replacement evidence covers actual consumers; retained schemas/assets migrate explicitly. | IC01–12; no live release/upgrade claim |

Milestone status is recorded in Execution progress below. The model's matching activities remain planned until
actual implementation evidence warrants a status change. GIP-1-M2 alone is not a
working installer; GIP-2 local fixture execution is not signed release support;
GIP-4 copy success is not authorization to apply to live XFMD.

## Verification and candidate identity

Run meaningful Go unit/integration/race checks at each engine milestone, plus
repository management/schema/document checks. Record exact candidate/source and
input/output hashes, baseline platform/toolchain and commands. Each public change
must be exercised through the actual packaged executable. Wrapper acceptance uses
its actual executable with a test release/cache and the real SDPTool child, not
only a mocked policy. Host scope starts Linux; Windows/macOS support requires native
execution of the same scenarios before being advertised. No browser/GUI is required
for this CLI workflow.

Shared fixture parity must cover normalized action ordering, AGENTS preservation
collision/idempotence, old/current/new managed edits, target path types, manual
unknown baselines, root files and append-only ledger finalization. Legacy PowerShell
results are comparison evidence, not an undocumented normative oracle. Record any
intentionally different policy and update the contract before changing expected
results. Every operation boundary gets a process-exit injection/recovery check;
power-loss durability and whole-tree rollback are outside the selected promise.

Before GIP-4 create the disposable fixture from current XFMD with read-only
provenance, not the dated IPD snapshot. Write no live files or release identity on
behalf of a manual baseline. Keep adoption mapping separate from trusted release
inventory and never reuse the copy's root-bound apply plan on the real worktree.

## Git, repository coordination and review

When selected, create stacked phase branches from the then-approved integration
base: sdp/install-gip-1, sdp/install-gip-2, sdp/install-gip-3,
sdp/install-gip-4 (choose an unused suffix if necessary). Commit each milestone with
its ID and concrete outcome; push completed phases and prepare reviewable PRs under
existing authorization. Do not merge/release without the corresponding selection.

GIP-3-M2 changes belong in gh-sdp, on a branch from its own current main. Coordinate
its local SDP instructions and review requirements before changing product code;
the SDP branch cannot replace that repository's ownership. Record the exact paired
SDPTool/gh-sdp candidates in both handoffs. The canonical engine contract remains
in SDP; links/compatible versions prevent duplicated policy prose.

The design's author walkthrough is not independent review or owner acceptance.
Review implementation and release claims against this plan's actual scenario
matrix. At the live-upgrade decision, present the selected verified release, fresh
manifest/plan, preservation evidence and remaining limitations as a concrete result.

## Remaining decisions and handoff

Working choices to review with the generated design: preview-only default with
explicit saved-plan apply; signed descriptor/key distribution; fixed local binary
delegation; forward-only recovery. Actual release number, production trust keys,
supported OS matrix and live XFMD rollout date remain unselected. These do not
block local schema/fixture design but must be resolved before publication/support.

IPD evidence remains linked rather than copied as Go implementation evidence.
KB-SDP-033 also retains the later practical XFMD SDL-modeling pilot; this plan
covers installation/adoption only. On completion, explicitly disposition that
remaining scope instead of leaving a delivered installer card ambiguously active.

## Execution progress


- GIP-1-M1: delivered. Executable closed records, five canonical golden fixtures, bounded strict JSON/YAML decoding, portable path/root inspection and receipt 3.0 reader delivered. Go package suite passes; management validator passes. No apply engine or signed distribution claimed.

- GIP-1-M2: delivered. Direct packaged install/upgrade previews delivered for explicit development artifacts, with deterministic plans, manual adoption, preserving relocation/link rebasing, old/current/new managed comparison and no-change detection. Go race suite passes. Packaged preview produced six actions and left its project empty. Signed resolution and apply remain later milestones; therefore the broader SDL InstallationPlanSlice stays planned.
