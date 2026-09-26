# Plan XFMD SDP adoption through gh-sdp and a practical SDL design pilot

| Field | Value |
| --- | --- |
| id | KB-SDP-033 |
| project | SDP |
| type | Study |
| CardState | backlog |
| Systems | SDP, SDL, SDPTOOL |
| created | 2026-09-25T22:20:56Z |
| source | Owner conversation 2026-09-26: SDL context, implementation drift and post-main XFMD adoption |
| next_review | After MAINT-SDP-0005 main consolidation |

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
[KB-SDP-032](%23032--Study--Viewpoint-navigation-feedback.md) and
[KB-SDP-018](%23018--Study--Toolkit-audit-and-organization.md).

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

Recommended boundary: reuse the existing installer; add release retrieval/verification
as a facade rather than another migration engine. Keep an externally authoritative
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
