# Plan XFMD SDP adoption through gh-sdp and a practical SDL design pilot

| Field | Value |
| --- | --- |
| id | KB-SDP-033 |
| project | SDP |
| type | Study |
| CardState | in-progress |
| PlanId | MAINT-SDP-0007 |
| Systems | SDP, SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | Concrete production release/live rollout proposal after completed GIP |

## Current disposition — RP1 release and live adoption authorized

[PLAN-SDP-0003](../../05--Implementation/SDPTool/Installation/Plan.md) has delivered
all selected phases: one Go engine, shared signed-test bootstrap, thin gh-sdp,
independent client review and fresh disposable XFMD adoption through both entry
points. Read the [evidence](../../05--Implementation/SDPTool/Installation/Evidence.md)
and [concrete rollout/retirement proposal](../../05--Implementation/SDPTool/Installation/Rollout-and-Retirement.md).

The owner selected [RP1](../../Maintenance/RP1/Plan.md) for publication and live
adoption on 2026-09-27. CardState is in-progress. The broader SDL-modeling pilot
remains a separate subsequent assignment.

## Owner intent and baseline

After consolidation of SDP-vNow into main, install or upgrade SDP in XFMD through
gh-sdp and begin describing XFMD with SDL. This creates practical experience for
viewpoints and assignment bundles. XFMD already has a bootstrapped SDP area and
KanBan history; inspect its current state rather than treating it as an empty
project. The owner reports no existing SDL model of XFMD; verify before authoring.

## Study and plan

Inspect gh-sdp's actual installation delegation, supported versions and project
recognition. Establish how it selects the intended main/distribution artifact;
merging to main does not publish a Toolkit release. Compare installed manifests,
capabilities and managed/project-owned paths. Prepare a dry-run upgrade with
versioned artifact provenance, backup/recovery and a Maintenance record. Preserve
XFMD's existing project documents and KanBan/ledger history. Coordinate existing
XFMD adoption cards rather than starting a duplicate bootstrap.

Study XFMD code and existing native blueprint documents critically. Describe
observed architecture in SDL with explicit uncertainty and distinguish it from
proposed improvements. Start with one bounded real workflow, then derive useful
views and propose a subsequent assignment-bundle pilot. A native blueprint is
input evidence, not an already accepted SDL blueprint format.

## Acceptance and responsibility

Deliver an executable adoption/pilot plan naming the actual gh-sdp path, artifact,
current/target installation identities, preservation checks and one end-to-end
XFMD workflow to model. Record missing gh-sdp capabilities explicitly. Route XFMD
code/model authoring and gh-sdp changes to their repository agents/cards; do not
perform live migration or product changes during this intake or main merge.

Depends on main consolidation and coordinates with
[KB-SDP-031](../completed/%23031--Study--SDL-assignment-bundles-and-blueprints.md),
[KB-SDP-032](../backlog/%23032--Study--Viewpoint-navigation-feedback.md) and
[KB-SDP-018](../canceled/%23018--Study--Toolkit-audit-and-organization.md).

## Worklog

2026-09-25T22:20:56Z: Registered in backlog before main integration; EVT-KB-SDP-000174.

## Owner proposal and read-only adoption assessment — 2026-09-26

The owner selects xfmd-sdl-navigation as the intended installation worktree and
proposes a release-bound external inventory, a small project-local version receipt,
and a custom manifest for its manually bootstrapped baseline. The present request
asks whether this is a sound upgrade method. No live apply or release publication
is performed by this assessment. Blueprint work remains deferred; this record is
on sdp/study-xfmd-install-adoption, stacked after BP1's captured owner clarification.

Observed XFMD: commit 456e095baa8a8809c11234a606e19987e4953bca, clean branch
sprint/009/phase/057-favorites. It has five phase homes and board schema 0.1 at
SDP/Agents/KanBan, but no installed-toolkit facts, project manifest, navigation
registration or shared ProjectManagement ledger. Bootstrap provenance is documented
in its SDP/Maintenance/SDP1/Plan-and-Evidence.md; do not invent an old release.

Current [process installer](../../../Toolkit/docs/Process-Installation.md) already
supports manual-five-phase-board adoption, deterministic exact plans, input hashes,
explicit relocations, backups, durable journaling and forward resume. Its external
profile artifact carries ownership, payload hashes and configuration identity.
Installed facts schema 2.0 records versions/profiles and configurationDigest; the
profile engine currently emits sourceCommit: null. This is not yet a trusted remote
release-catalog resolver. Internal digest consistency is not publisher authentication.

Read-only PlanJson runs against the actual target used the unchanged main profile
artifact with configurationDigest
da2e2ebe2002abf61b70ebbb7f176aa46190c421c8a4573f2d5ca3eea118e41a.
The normal plan reports unknown historical versions, recognizes manual-five-phase-board
and blocks on managed-refresh-requires-force: AGENTS.md. A second **preview only**
with ForceManagedFiles is applicable: 139 file actions, 193 preserved paths, no
conflicts. It proposes preserving project instructions as AGENTS-project.md while
installing managed AGENTS.md. The plan warns that scripts/non-Markdown references
need separate review after relocation. All 246 observed file hashes still match
after planning; no installation was applied. This is planning feasibility, not
successful migration or acceptance of every proposed write.

GitHub inspection: Hans-Einar/SDP has no published releases and its manifest still
says 0.2.0 unreleased. gh extension list has no gh-sdp. Remote gh-sdp main remains
32613734781bf39f2fce176db2acfb2284dfc92f and contains process documents rather than
CLI product code. Do not present gh-sdp release lookup/install as implemented.

Historical recommendation, superseded by the owner decision below: reuse the
existing installer and add release retrieval/verification as a facade rather than
another migration engine. Keep an externally authoritative
release inventory and a project receipt binding repository, version, immutable
artifact digest and process identity. A local cached copy is acceptable if verified;
not storing the inventory locally is not itself a tamper-proof guarantee. The
mutable receipt is a lookup hint until validated against trusted release facts and
actual managed file state. Inventory should cover all installer-owned paths,
including root instructions/skills outside SDP, and identify file/directory types.

Treat XFMD's custom manifest as an **observed adoption baseline**: exact project
commit, inspected files/hashes/types, ownership classification, unknown original
version and reviewed migration map. Store it outside the live target. A snapshot
proves what was observed, not that it was once a standard release. Validate against
that frozen baseline immediately before apply. Preserve project documents, cards,
append-only ledger bytes and legitimate additions; they are not immutable release
payload. Compare old managed payload, actual local bytes and new payload to detect
local edits; a missing file or collision requires a specific disposition. Inventories
alone do not define rename/move semantics: use explicit versioned migration steps.

Next proposed work: freeze/adapt the baseline and inspect all planned writes;
exercise migration and repeat/no-change on a disposable snapshot, including local
scripts and incoming links; select an exact target artifact/release; implement or
coordinate the missing gh-sdp facade in its own repository; only then apply the
reviewed operation to the selected live worktree and verify the receipt/history.
Use the existing installed-facts contract where possible; any additional release
identity fields require an explicit schema/reader change. No new manifest schema,
release version, migration approval or direct apply is selected by this note.

Known patterns: package ownership and preservation of locally edited configuration
([Debian policy](https://www.debian.org/doc/debian-policy/ch-files.html#configuration-files));
immutable release assets/tags and verifiable release attestations
([GitHub documentation](https://docs.github.com/en/code-security/concepts/supply-chain-security/immutable-releases)).
A hash from a mutable local manifest alone does not establish a trusted baseline.

2026-09-26T12:37:22Z: EVT-KB-SDP-000185. CardState remains backlog while adoption/release/facade boundaries
are clarified; this is a preflight observation, not completion of the card's wider
SDL-modeling pilot.


## Earlier Go-engine allocation — superseded by SDPTool clarification below

The owner selects a compiled Go installation engine owned by the SDP repository.
gh-sdp should be a thin client for accessing that engine without a separate manual
Toolkit installation. This supersedes any interpretation of the earlier facade
recommendation that would retain PowerShell as the required production engine.
Reuse the existing ownership, planning, preservation and recovery contracts and
verification cases; do not maintain a separate migration-policy implementation in
gh-sdp. The current PowerShell implementation remains the observed implementation,
not the selected target. Its replacement/removal requires a migration plan and
verified Go behavior; no executable is changed by this decision capture.

Clarification: GitHub CLI extensions execute locally. GitHub hosts source and
precompiled release assets; it does not execute a local project upgrade remotely.
Official references: [extension execution](https://cli.github.com/manual/gh_extension)
and [extension distribution](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions).

Recommended distribution, pending detailed design: gh-sdp resolves an exact,
compatible Toolkit release, verifies and caches its platform-specific Go executable
and release inventory, then delegates local plan/apply operations. This avoids a
manual Toolkit setup while preserving one engine. An alternative is linking the
same SDP-owned Go package into gh-sdp; it avoids a second runtime download but
couples engine updates to extension releases. Neither packaging choice changes
ownership. Define engine/protocol compatibility, provenance verification, cache
behavior and installed-version receipts before implementation.

The desired user entrypoint is `gh sdp upgrade --manifest xfmd-upgrade.yaml`, run
from the project repository. The custom manifest describes the observed manual
baseline and explicit adoption mapping; it does not replace the authoritative
release inventory. This command and custom-manifest schema are not implemented.

Refreshed gh-sdp inspection: origin/main at be990cc includes merged PR #3 and
accepts the historical Study with low findings. Requirements and Architecture
remain templates; there is no product executable. Its mandate explicitly rejects
required PowerShell on Linux/macOS. The old Study's upstream assessment must be
reconciled against the current Toolkit contract before using it as a current gap
list. The new owner allocation keeps the engine in SDP rather than assigning a
separate apply engine to gh-sdp. Coordinate that change in gh-sdp's governing
records during its next authorized work item.

Next scope: plan the Go engine migration and thin client as one runnable XFMD
adoption workflow, using a disposable copy first. Keep release publication and
live project mutation explicit. CardState remains backlog; the broader SDL pilot
is still outstanding.

2026-09-26T18:51:00Z: EVT-KB-SDP-000186. XA1-M2 records the owner-selected Go engine ownership and local extension execution; no implementation or installation claimed.

2026-09-26T19:03:42Z: EVT-KB-SDP-000188. Owner clarifies existing SDPTool owns installation and all common commands; withdraw separate Toolkit-engine direction. Preserve XFMD adoption and SDL pilot; root relocation tracked by MAINT-SDP-0006.

## Current owner direction: reuse root SDPTool

The owner identifies the existing Go SDPTool as the common entry point for all
SDP commands. Its new source home is [SDPTool](../../../SDPTool/README.md).
Installation/upgrade is a responsibility to implement inside that existing
system, not a new Toolkit executable or a gh-sdp apply engine. gh-sdp should
retrieve/invoke SDPTool as a thin client. The preceding standalone-Toolkit binary
recommendation is withdrawn; its release identity, integrity, preservation and
manual-baseline requirements still apply to SDPTool distribution and execution.

SDPTool currently reads installation facts/journals but does not implement
install/upgrade. Existing PowerShell behavior, schemas and fixtures remain in the
legacy/transition Toolkit until the Go migration is verified. The current
[SDPTool requirements](../../02--Requirements/SDPTool.md) distinguish that target
from implemented discovery/preview/navigation. [MAINT-SDP-0006](../../Maintenance/ST1/Plan.md)
only relocates source and reconciles responsibilities. Completed KB-SDP-028 and its
Maintenance evidence remain historical deliveries; they are not retroactively
canceled. The separate KB-SDP-018 audit is canceled at the owner's request.

Next selected implementation plan must implement the bounded workflow
`gh sdp upgrade --manifest xfmd-upgrade.yaml` through SDPTool on a disposable XFMD
copy before any live adoption. Detailed CLI/manifest/protocol design and release
publication remain outstanding. No new executable is delivered by this card update.

2026-09-26T22:16:08Z: EVT-KB-SDP-000190. Activated as ready under PLAN-SDP-0002; IPD-0-M1 records plan/baseline, installation modeling remains next.

2026-09-26T22:23:11Z: EVT-KB-SDP-000191. Owner starts execution of IPD-1 through IPD-3; CardState ready to in-progress. SDL model, contracts and generated design review only; no live migration or production installer implementation.

2026-09-26T22:24:18Z: EVT-KB-SDP-000192. IPD-1-M1 delivered installation responsibilities in SDL. CardState remains in-progress for contracts and generated review.

2026-09-26T22:31:14Z: EVT-KB-SDP-000193. IPD-2-M1 completes design contracts and scenario walkthrough; CardState remains in-progress for reproducible generated review and implementation handoff.

2026-09-26T22:39:32Z: EVT-KB-SDP-000194. CardState in-progress to gate-review for completed IPD design, generated views and concrete planned GIP implementation phases. Owner review concerns preview/apply, trust/bootstrap and recovery choices. Broader adoption and SDL pilot remain outstanding.

2026-09-26T23:05:54Z: EVT-KB-SDP-000195. CardState gate-review to in-progress; owner selects PLAN-SDP-0003 implementation after IPD design.

2026-09-26T23:54:50Z: EVT-KB-SDP-000196. CardState in-progress to gate-review for concrete production release/live-adoption proposal. GIP is complete; subsequent XFMD SDL modeling explicitly deferred.

2026-09-27T07:53:08.653575Z: EVT-KB-SDP-000197. Owner selects RP1; CardState gate-review to in-progress.
