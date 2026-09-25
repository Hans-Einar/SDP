# Versioned process installation — contract 2.0

Authority: MAINT-SDP-0003, executed by owner instruction on 2026-09-25.
This contract extends Install-SDP.ps1; the existing unqualified install-v1
interface remains supported. There is no second installer executable or Go
installation engine.

## Artifact and ownership

Toolkit/profiles/five-phase.json is the authored process configuration.
A build command emits a self-contained JSON artifact with schemaVersion 2.0,
profile sdp-five-phase/0.1, managementProfile sdp-project-management/0.1,
configurationDigest, release facts, file inventory (source, destination,
ownership, SHA-256 and base64 bytes), ordered relocations and prerequisites.
Build identity excludes timestamps. The digest covers canonical configuration
and every input byte, including release facts. No repository lifecycle records
are copied. The legacy inventory remains authoritative only for install-v1.

Files copied from Skills and Toolkit payload are managed; neutral process
templates are initialize-if-missing. Project documents and histories are owned
by the project. Managed replacements need explicit ForceManagedFiles and
byte backups, even at the same unreleased Toolkit version. Unknown AGENTS.md
is preserved as AGENTS-project.md before installing managed instructions; an
unequal existing preservation file is a conflict.

Installed facts use explicit schema 2.0 at the existing YAML manifest path.
They retain the release/skills fields and add processProfile, managementProfile
and configurationDigest. The schema 1.0 reader must reject these facts.
New readers accept both versions. Toolkit version stays 0.2.0 unreleased;
process profile, schema versions and capabilities are independent identities.

## Baselines and transitions

| Observed baseline | Supported transition |
| --- | --- |
| No SDP area | Initialize neutral five-phase and shared management areas |
| Legacy seven-phase layout | Explicit phase relocation; archive Study and DesignAnalysis; preserve relative Markdown links |
| Five-phase local shared board | Preserve sources/history and adopt installed facts |
| XFMD manual five-phase, Agents/KanBan board 0.1 | Move board to KanBan; transfer exact ledger bytes to shared history; update board declaration and relative Markdown links |
| Installed schema 1.0 or 2.0 | Validate facts, reject downgrade, apply selected profile and managed refresh policy |
| Mixed old/new phase destinations, duplicate boards/histories, unsupported schemas | Conflict; no mutation; resolve an explicit mapping first |

Folder observations classify a layout only. Missing version facts remain
unknown. Empty old directories left after file relocation are not authorities.
Historical ledger bytes and IDs are never rewritten. Project-specific scripts
and non-Markdown references are preserved and identified for operator review;
the installer does not claim semantic migration of arbitrary executable code.
No historical document is split into Systems/containers automatically.

## Plan and apply

PowerShell 7.4+ on Linux is the initial verified host; Windows execution is
experimental until its conformance job passes. Build needs Python 3 and PyYAML;
using a built artifact needs only PowerShell. Go tools are prebuilt separately.

Install-SDP.ps1 -ProfileArtifact ARTIFACT -ProjectRoot ROOT -PlanJson produces
a read-only deterministic schema 2.0 plan. Paths, baseline, observed file hashes,
artifact hash, actions, preserved files and conflicts are explicit.
Apply uses -ApplyPlan PLAN with the same artifact and force setting. An exact
recomputed plan must match before mutation. Source and destination drift,
unsupported filesystem objects, overlapping roots and case collisions fail.
A plan is not a transferable authorization token; its project root is bound.

Execution uses one exclusive OS file lock and a durable journal per operation
under SDP/.sdp-operations. Each write/delete has before/after hashes, a backup
and an atomic per-file replacement. Journals retain the exact plan and generated
finalization bytes. After interruption, -ResumeOperation ID verifies completed
steps and permits only the recorded before/after bytes for the pending step.
Other plans are blocked while an incomplete operation exists. No whole-tree
atomicity or automatic rollback is claimed. Recovery is forward resume;
post-failure edits stop recovery instead of being overwritten.

Facts, Maintenance report and management history finalize as journaled steps.
Management IDs are allocated from the target history and report paths, never
copied from this repository. Bootstrap records created(active) and completed
using management payload 0.1; prior history remains an exact byte prefix.
A journal is completed only after all final bytes verify. A crash between
history publication and completion resumes without appending a duplicate event.
No-change creates neither report nor event. A failed journal is the durable
failure report until successful recovery.

## Navigation and tools

Create only an empty source registration for a new project: no invented model,
implementation plan or viewer. Existing navigation bindings stay project-owned;
known board path relocation may update that path. SDPTool reads declared
installed facts without claiming full installation verification. Installation
does not build Go tools during documentation viewing.

## Verification boundary

Disposable fixtures cover each transition, corruption, source/target drift,
force/preservation, every journal boundary, retry/no-change, lock contention,
unsupported schemas, unknown versions, links and case collisions. Baseline
XFMD provenance: cf11709e4ec9d6925b0d17d95c71f72fbf011959, clean worktree,
inspected 2026-09-25; board schema 0.1 in SDP/Agents/KanBan, no installed facts.
No live project upgrade is authorized by this implementation assignment.

## Executing the built profile

From a Toolkit checkout or packaged source distribution, using an existing
separate project directory:

~~~sh
python3 Toolkit/scripts/build_process_profile.py --output /tmp/sdp-profile.json
pwsh -NoProfile -File Toolkit/scripts/Install-SDP.ps1 -ProjectRoot /path/to/project -ProfileArtifact /tmp/sdp-profile.json -PlanJson > /tmp/sdp-plan.json
pwsh -NoProfile -File Toolkit/scripts/Install-SDP.ps1 -ProjectRoot /path/to/project -ProfileArtifact /tmp/sdp-profile.json -ApplyPlan /tmp/sdp-plan.json
~~~

Review canApply, conflicts, warnings, before/after hashes and preserved paths.
A managed refresh requires -ForceManagedFiles on both plan and apply; this does
not authorize overwriting project content. Keep the reviewed artifact and plan.
Use -ResumeOperation with the returned install-<digest> ID and original artifact
after an interruption. The journal's error/status and backups are under
SDP/.sdp-operations/<ID>. Do not edit the journal or remove its lock file while
another process may be using it. A failed operation blocks another installation.
Restore a conflicting post-failure edit deliberately before resuming; the engine
will not discard it. Missing dependencies fail explicitly; no tool is downloaded.

Markdown rebasing covers inline links (including balanced/escaped parentheses,
titles and angle destinations), images and reference definitions. It scans SDP
and incoming project Markdown. Git metadata, node_modules, .venv, vendor, build,
.cache and symlink trees are excluded from incoming-reference discovery; review
these and non-Markdown code/config references when a relocation is planned.
An existing broken link is not evidence of successful content repair. Historical
NDJSON bytes are transferred without reinterpretation or rewriting.

The recovery guarantee covers process interruption at journaled boundaries.
Per-file flush/rename is used; power-loss durability of the filesystem and
whole-operation rollback are not claimed. Verification must distinguish injected
process exits from hardware/power failure.
