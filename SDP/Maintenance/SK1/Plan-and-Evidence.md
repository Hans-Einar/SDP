# MAINT-SDP-0002 — canonical skills and project activation

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0002 |
| project | SDP |
| state | completed |
| source | SCRUM-SDP-0003; KB-SDP-027; owner conversation 2026-09-25 |
| Systems | SDP |

Selected by [Scrum-0003](../../Agents/Scrum/Scrum--%230003--Skills-and-agent-activation.md).
This record owns the future Maintenance delivery; the source Study closes when
its Scrum and plan are delivered. No wrapper card or Sprint is required. Current delivery and evidence are recorded in the milestone table and execution
sections below; the earlier inventory describes the pre-adoption baseline.

## Outcome and boundary

Maintain one English skill collection under root Skills/, with a small routing
entrypoint and independently loadable roles. Make that collection discoverable
and usable in this project, with portable references and repeatable installation.
Preserve useful behavior and remove competing maintained collections after their
consumers migrate. Keep archived evidence immutable. No SDL/SDUI syntax/runtime,
XFMD application changes, global user-skill installation or new plugin is needed.

Use the draft collection as the editing base, not as an already approved release.
Retain ten established role identities and the three useful new roles below.
A routing skill loads only relevant instructions; thirteen files must not mean
thirteen simultaneous roles or automatic multi-agent execution. Host/session
instructions and existing owner authorization govern delegation and publication.

## Inventory and adoption map

Inspected at c4aed09944235f5da1d247001688f6f162ec124d. Toolkit/skills contains ten
version-1.0.0 role files; all lack native name/description fields. Toolkit/skills_v2
contains thirteen draft.2 files with those fields and shared workflow/source-map
references. Neither collection is currently a root project discovery adapter.
Frozen bootstrap copies under SDP/History are evidence, not another editing base.

| Target under Skills/ | Source and required adaptation |
| --- | --- |
| sdp | Retain draft router; correct expansion to System Design Process; discover current local/installed profile and load the smallest relevant role set |
| sdp-change-analysis | Retain draft; route new findings into existing cards, distinguishing derived obligations from new choices |
| sdp-architect | Replace legacy role with adapted draft; honor three Systems and actual container boundaries, five phases without equating folders to abstraction levels |
| sdp-master | Replace legacy role with adapted draft; support authorized milestones, Maintenance and optional Sprints, not only mandatory Slice/Fix/delegation |
| sdp-worker | Replace legacy role with adapted draft; support the selected bounded work and milestone evidence without imposing older work-record conventions |
| sdp-reviewer | Replace legacy role with adapted draft; preserve independent-review meaning, label same-context checks honestly and respect host delegation constraints |
| sdp-verifier | Replace legacy role with adapted draft; preserve evidence levels, candidate identity and explicit unverified outcomes |
| sdp-traceability | Replace legacy role with adapted draft; distinguish management history from system change evidence and use supported system-prefixed IDs |
| sdp-steering | Retain draft; record actual owner decisions in existing records, without manufacturing approval or a new mandatory document tree |
| sdp-auditor | Replace legacy role with adapted draft; verify the actual declared profile, host discovery and source/install inventory separately |
| sdp-vertical-refactor | Replace legacy role with adapted draft; use current Refactor/Maintenance records and preserve runnable workflows |
| sdp-versioning | Replace legacy role with adapted draft; separate Toolkit, skill, process/profile, language and project product versions |
| sdp-release | Replace legacy role with adapted draft; preserve scoped publication authority and actual remote evidence; phase push is not merge/release acceptance |

Keep the shared document-workflow reference, updated for the current local profile,
English policy, KanBan, optional Scrum/Sprint/Maintenance and the two histories.
Replace the draft source map's old-current pointers with a concise authority map;
retain its earlier provenance as dated history. Older installed profiles may have
other supported conventions: detect the profile rather than assert that every
consumer already uses this repository's local process. Do not impose local
system names, branch names or file URLs on all consuming projects.

## Discovery decision and measured limits

[Official Codex skill documentation](https://learn.chatgpt.com/docs/build-skills)
requires name/description and describes repository discovery through .agents/skills,
including symlinked skill directories. It distinguishes initial metadata discovery
from loading full instructions. Therefore use per-role repository symlinks:
.agents/skills/sdp → ../../Skills/sdp, and equivalently for each selected role.
Root AGENTS.md explains routing and current authority; a root SKILL.md is not the
selected discovery mechanism. Other hosts need their own verified adapters later.

[Recorded probe](discovery-probe.json), Codex CLI 0.156.1 on this Linux host:

| Temporary fixture | Observed native catalog result |
| --- | --- |
| Draft router only in root Skills/ | No fixture skill discovered |
| Draft router with .agents/skills/sdp symlink | Exactly sdp discovered with repo scope; canonical source path returned |
| Legacy Master behind the same kind of adapter | No fixture skill; missing description error |
| Root SKILL.md plus root Skills/, no adapter | No fixture skill discovered |

Reproduce from this repository with:

```sh
python3 SDP/Maintenance/SK1/probe_discovery.py
```

The helper creates disposable Git repositories, copies inspected inputs, calls
initialize and skills/list in a fresh local app-server per case, filters output
to fixture paths, then removes the fixtures. It starts no thread/model turn and
changes no user configuration. The record proves catalog discovery only. It does
not prove instruction loading, routing behavior, all thirteen roles, IDE/desktop
support, nested worktree behavior or an installed package. Future results must
identify the then-current source commit; do not overwrite this dated result as
if it tested a future adopted collection.

## Coupled installer work

The current install manifest copies ten Toolkit/skills sources to .codex/skills.
Its generated AGENTS template explicitly reads that location. The validator checks
skillId, skillVersion, minimumToolkitVersion, capabilities and compatibilityNotes;
its skill-ID pattern excludes bare sdp. Draft candidate-version does not satisfy
that distribution contract. A path-only move would break installation or leave
old workflow instructions active.

SK1 must update source inventory, metadata validation, generated destinations,
AGENTS references, relevant allowlists and install fixtures together when adopting
the canonical collection. Retain version/capability facts in the existing release
and installed manifests; reconcile native name/description with distribution
metadata without creating a competing registry. Select versions from actual
compatibility impact at execution; draft.2 is not a released 2.0.0. For a platform
without usable symlinks, a generated install copy may be justified and hash-checked;
it must never become a separately authored skill collection.

[KB-SDP-028](../../KanBan/completed/%23028--Study--Installer-upgrade-and-versioned-layout.md)
owns the broader declarative layout compiler, versioned process upgrades and
manual-project adoption. The owner authorized SK1 execution before that Scrum on 2026-09-25. SK1
therefore decides its minimal distribution boundary now and hands the actual
result to the later installer Scrum. Do not expand
SK1 into migration of every SDP process directory. The broad
[Toolkit audit](../../KanBan/backlog/%23018--Study--Toolkit-audit-and-organization.md)
reuses this inventory. Existing ID compatibility failures remain KB-SDP-011.

## Implementation milestones

Use a new stacked phase branch sdp/phase-sk1-skills-activation, with one commit
per completed milestone, evidence here and push at phase completion.

| Milestone | Delivery and acceptance | State |
| --- | --- | --- |
| SK1-M1 | Apply the owner-selected maintenance-first ordering; finalize concise role/workflow content, compatible metadata and migration inventory, including every current source/consumer and removal condition | Delivered |
| SK1-M2 | Atomically adopt root Skills/, migrate manifest/template/validator consumers and project adapter/AGENTS, remove the two superseded maintained sources; verify links, metadata, no duplicate names and supported clean/repeat install plus preservation/failure cases | Delivered |
| SK1-M3 | Verify all adopted roles in a fresh host catalog and actual explicit loading/routing on representative tasks; record host/version, results and gaps, update evidence and close Maintenance only after acceptance | Delivered |

Milestone verification must preserve frozen hashes/history and unrelated drafts.
Run affected installer/conformance and skill metadata tests, compare any known
baseline failures explicitly, and run shared management/link checks. Assert that
installed relative references resolve and repeated install changes nothing.
Test conflicting/unmanaged adapter paths without overwriting them. Exercise
repository root and nested cwd; distinguish Git worktrees/subrepos by actual
host behavior instead of assuming identical discovery.

For actual loading, use an explicitly invoked adopted skill in a fresh authorized
session and record its resolved source. For routing, use one bounded implementation
prompt, a read-only audit, an owner-choice prompt and a trivial task needing no SDP
role. Observe which files were loaded and whether the relevant existing record
was used; a correct answer alone is weak loading evidence. Do not fabricate
independent review or spawn agents merely to satisfy a role name. If only catalog
checks are available, SK1-M3 remains incomplete with a concrete handoff.

## Recovery and completion

Before migration record the full consumer/file inventory and existing adapter
ownership. Perform installer tests in disposable projects; never overwrite user
configuration to make a test pass. If adoption validation fails, keep the current
source/adapter intact or restore the previous coherent milestone; do not leave
half the installer pointing at deleted sources. Preserve local customizations,
report conflicts, and use the existing installer plan/backup policy.

Completion requires one maintained collection, working references and installation,
observable discovery/loading/routing evidence, truthful version facts, and explicit
remaining host limitations. No parallel legacy implementation should remain active.
The current state is maintained in the metadata and milestone table. #027's
completed Study selected this job; its status does not substitute for execution.

## Execution authority — 2026-09-25

The owner requests this Maintenance now, before the next Scrum. This supersedes
the earlier recommended dependency order; no additional meeting is required.
Root Skills/ and the native repository adapter are in scope. Keep existing
consumer destinations when tested host support allows it, avoiding an unrelated
process-tree migration. Final distribution/version choices and evidence follow.

## SK1-M1 result

[Migration contract](Migration-contract.md) fixes the consumer inventory, metadata,
versions and maintenance-first ordering. Thirteen prepared entrypoints pass the
skill-creator validator. Their measured source hashes are in Content-inventory.json;
shared references now describe the actual profile and distinguish the two histories.
The active install source and repository catalog remain unchanged until M2.
[Baseline](toolkit-baseline.txt) records 38 existing Toolkit errors; no unrelated
Traceability repair is included. The earlier draft evaluation is retained as history.

## SK1-M2 result

The atomic adoption and installer migration are delivered. See [execution
evidence](Evidence.md) for metadata, preservation, native adapter and independent
review results. Root Skills/ is now the only authored collection; both former
maintained sources are removed. Installed paths remain .codex/skills as measured
on the current host, and all shared references are included. Two new conformance
scenarios make the same-version skill-transition boundary portable and explicit.
The broader installer Scrum remains unstarted; no external project was migrated.

## SK1-M3 result and closure

MAINT-SDP-0002 is complete. [Catalog and behavior evidence](Evidence.md) verifies
all thirteen roles on Codex CLI 0.156.1 plus four independent representative tasks.
The native catalog was rechecked on M2 commit 2d1c560; role bytes match the M1
inventory, and the original Scrum probe remains reproducible from baseline Git
bytes. Local/current profile adoption is complete within the tested host scope;
other hosts and implicit-selection generalization remain explicit limitations.

No parallel maintained legacy skill collection remains. No external repository,
process tree or global user configuration was migrated. #028 receives the measured
install boundary and version facts; #018 reuses the completed skills inventory;
#011 retains the exact 38-error pre-existing Traceability baseline. #010 remains
on its owner's existing non-blocking review. The next Scrum has not started.
